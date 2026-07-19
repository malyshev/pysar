package mcpserver

import (
	"encoding/json"
	"path/filepath"

	"pysar/internal/onboarding"
)

var saveStyleProfileSchema = map[string]interface{}{
	"type":       "object",
	"properties": profileProperties(),
	"required":   []string{"rules", "goldens"},
}

var saveStyleTemplateSchema = func() map[string]interface{} {
	props := profileProperties()
	props["name"] = map[string]string{"type": "string", "description": "A short, memorable name for this reusable template, e.g. 'Plain English -- GOV.UK style'. Can be changed later without changing the template's slug -- pass the existing slug to rename in place."}
	props["slug"] = map[string]string{"type": "string", "description": "Optional -- an existing template's stable machine key (from list_style_templates), e.g. 'generic'. Provide this to update/rename that exact template in place. Omit when creating a brand-new template; a slug is derived from the name automatically."}
	return map[string]interface{}{
		"type":       "object",
		"properties": props,
		"required":   []string{"name", "rules", "goldens"},
	}
}()

// registerSaveStyleProfile wires the save_style_profile tool
// (dec-20260719-a1ac8959): mirrors save_voice_profile exactly, but requires
// Rules instead of the voice-flavored scalar fields, since style content
// (structure, sentences, word choice, formatting) doesn't fit tone/formality.
func (s *Server) registerSaveStyleProfile() {
	s.register(
		tool{
			Name:        "save_style_profile",
			Description: "Validate and persist an author's style profile to .pysar/style.md. Style is the craft standard the writing must meet (structure, sentence construction, word choice, formatting) -- distinct from voice (how it sounds). Rejects an incomplete profile with a structured error naming exactly what's missing.",
			InputSchema: saveStyleProfileSchema,
		},
		s.callSaveStyleProfile,
	)
}

func (s *Server) callSaveStyleProfile(args json.RawMessage) callToolResult {
	path := filepath.Join(s.baseDir, ".pysar", "style.md")
	return s.saveProfileToPath(args, onboarding.KindStyle, path)
}

// registerSaveStyleTemplate wires the save_style_template tool, mirroring
// save_voice_template for a second, cross-project destination
// (~/.pysar/templates/style/).
func (s *Server) registerSaveStyleTemplate() {
	s.register(
		tool{
			Name:        "save_style_template",
			Description: "Validate and persist a completed style profile as a named, reusable template usable across projects (~/.pysar/templates/style/). Same completeness rules as save_style_profile -- rejects an incomplete profile with a structured error naming exactly what's missing.",
			InputSchema: saveStyleTemplateSchema,
		},
		s.callSaveStyleTemplate,
	)
}

func (s *Server) callSaveStyleTemplate(args json.RawMessage) callToolResult {
	return s.saveTemplate(args, onboarding.KindStyle)
}

// registerListStyleTemplates wires the list_style_templates tool, mirroring
// list_voice_templates -- avoids Glob/Bash path-expansion entirely.
func (s *Server) registerListStyleTemplates() {
	s.register(
		tool{
			Name:        "list_style_templates",
			Description: "List every reusable style template available across all pysar projects (~/.pysar/templates/style/), including the built-in 'generic' default. Returns each template's display name, its stable slug (pass this back to save_style_template to rename in place), and full content.",
			InputSchema: map[string]interface{}{"type": "object", "properties": map[string]interface{}{}},
		},
		s.callListStyleTemplates,
	)
}

func (s *Server) callListStyleTemplates(args json.RawMessage) callToolResult {
	return s.listTemplates(onboarding.KindStyle)
}
