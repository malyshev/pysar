package mcpserver

import (
	"encoding/json"
	"path/filepath"

	"pysar/internal/onboarding"
)

var saveVoiceProfileSchema = map[string]interface{}{
	"type":       "object",
	"properties": profileProperties(),
	"required":   []string{"tone", "formality", "sentence_length", "register", "goldens"},
}

var saveVoiceTemplateSchema = func() map[string]interface{} {
	props := profileProperties()
	props["name"] = map[string]string{"type": "string", "description": "A short, memorable name for this reusable template, e.g. 'Measured plain English -- speakable, understated, general audience'. Can be changed later without changing the template's slug -- pass the existing slug to rename in place."}
	props["slug"] = map[string]string{"type": "string", "description": "Optional -- an existing template's stable machine key (from list_voice_templates), e.g. 'generic'. Provide this to update/rename that exact template in place. Omit when creating a brand-new template; a slug is derived from the name automatically."}
	return map[string]interface{}{
		"type":       "object",
		"properties": props,
		"required":   []string{"name", "tone", "formality", "sentence_length", "register", "goldens"},
	}
}()

// registerSaveVoiceProfile wires the save_voice_profile tool: accepts
// structured VoiceProfile fields, validates via the real
// internal/onboarding.Validate (dec-20260719-fa0366dd invariant --
// completeness checking is never re-implemented as an LLM-only check), and
// on success renders and writes .pysar/voice.md server-side. No client-side
// Write/Edit tool call is ever made for this content.
func (s *Server) registerSaveVoiceProfile() {
	s.register(
		tool{
			Name:        "save_voice_profile",
			Description: "Validate and persist an author's voice profile to .pysar/voice.md. Rejects an incomplete profile with a structured error naming exactly what's missing -- never silently accepts one.",
			InputSchema: saveVoiceProfileSchema,
		},
		s.callSaveVoiceProfile,
	)
}

func (s *Server) callSaveVoiceProfile(args json.RawMessage) callToolResult {
	path := filepath.Join(s.baseDir, ".pysar", "voice.md")
	return s.saveProfileToPath(args, onboarding.KindVoice, path)
}

// registerSaveVoiceTemplate wires the save_voice_template tool: the operator
// chose it to reuse save_voice_profile's exact validate-then-persist
// rationale for a second, cross-project destination (dec-20260719-3e36577e)
// -- a template is a plain Profile, same schema, same completeness check,
// only the storage location and an added display name differ.
func (s *Server) registerSaveVoiceTemplate() {
	s.register(
		tool{
			Name:        "save_voice_template",
			Description: "Validate and persist a completed voice profile as a named, reusable template usable across projects (~/.pysar/templates/voice/). Same completeness rules as save_voice_profile -- rejects an incomplete profile with a structured error naming exactly what's missing.",
			InputSchema: saveVoiceTemplateSchema,
		},
		s.callSaveVoiceTemplate,
	)
}

func (s *Server) callSaveVoiceTemplate(args json.RawMessage) callToolResult {
	return s.saveTemplate(args, onboarding.KindVoice)
}

// registerListVoiceTemplates wires the list_voice_templates tool. Added
// after a real Claude Code session showed Glob("~/.pysar/templates/voice/*.md")
// falling back to a raw Bash `ls` (Glob does not expand `~`), which
// triggered the exact unwanted permission prompt this project has fought to
// eliminate all session -- refuting dec-20260719-3e36577e's prediction that
// listing could safely stay a direct filesystem read. This tool sidesteps
// path expansion entirely (the server already knows homeDir) and returns
// each template's full content inline, so the skill never needs a second
// per-template Read either.
func (s *Server) registerListVoiceTemplates() {
	s.register(
		tool{
			Name:        "list_voice_templates",
			Description: "List every reusable voice template available across all pysar projects (~/.pysar/templates/voice/), including the built-in 'generic' default. Returns each template's display name, its stable slug (pass this back to save_voice_template to rename in place), and full content.",
			InputSchema: map[string]interface{}{"type": "object", "properties": map[string]interface{}{}},
		},
		s.callListVoiceTemplates,
	)
}

func (s *Server) callListVoiceTemplates(args json.RawMessage) callToolResult {
	return s.listTemplates(onboarding.KindVoice)
}
