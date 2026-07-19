package mcpserver

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"pysar/internal/onboarding"
)

// voiceProfileProperties is the JSON-Schema property set shared by
// save_voice_profile and save_voice_template -- both persist the exact same
// Profile shape (dec-20260718-ab150a73), differing only in destination and,
// for the template tool, an added display name.
func voiceProfileProperties() map[string]interface{} {
	return map[string]interface{}{
		"tone":            map[string]string{"type": "string", "description": "How the author wants to come across, e.g. 'warm, direct, a little dry'"},
		"formality":       map[string]string{"type": "string", "description": "Where on the conversational-to-formal range, and how"},
		"sentence_length": map[string]string{"type": "string", "description": "Sentence length/rhythm preference"},
		"register":        map[string]string{"type": "string", "description": "Vocabulary register, e.g. 'plain, international-standard English'"},
		"banned_phrases": map[string]interface{}{
			"type":        "array",
			"items":       map[string]string{"type": "string"},
			"description": "Optional -- phrases/constructions the author never wants to see",
		},
		"notes": map[string]string{"type": "string", "description": "Optional -- texture the structured fields don't capture"},
		"goldens": map[string]interface{}{
			"type": "array",
			"items": map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"label": map[string]string{"type": "string", "description": "Short caption for what this passage demonstrates"},
					"text":  map[string]string{"type": "string", "description": "The reference passage itself"},
				},
				"required": []string{"label", "text"},
			},
			"description": fmt.Sprintf("At least %d real reference passages the author has confirmed", onboarding.MinGoldens),
		},
	}
}

var saveVoiceProfileSchema = map[string]interface{}{
	"type":       "object",
	"properties": voiceProfileProperties(),
	"required":   []string{"tone", "formality", "sentence_length", "register", "goldens"},
}

var saveVoiceTemplateSchema = func() map[string]interface{} {
	props := voiceProfileProperties()
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
	var input struct {
		Tone           string                     `json:"tone"`
		Formality      string                     `json:"formality"`
		SentenceLength string                     `json:"sentence_length"`
		Register       string                     `json:"register"`
		BannedPhrases  []string                   `json:"banned_phrases"`
		Notes          string                     `json:"notes"`
		Goldens        []onboarding.GoldenExample `json:"goldens"`
	}
	if err := json.Unmarshal(args, &input); err != nil {
		return errorResult("invalid arguments: %s", err)
	}

	profile := onboarding.Profile{
		Kind:           onboarding.KindVoice,
		Tone:           input.Tone,
		Formality:      input.Formality,
		SentenceLength: input.SentenceLength,
		Register:       input.Register,
		BannedPhrases:  input.BannedPhrases,
		Notes:          input.Notes,
		Goldens:        input.Goldens,
	}

	if err := onboarding.Validate(profile); err != nil {
		return errorResult("%s", err.Error())
	}

	rendered, err := onboarding.Render(profile)
	if err != nil {
		return errorResult("render profile: %s", err)
	}

	path := filepath.Join(s.baseDir, ".pysar", "voice.md")
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return errorResult("create .pysar directory: %s", err)
	}
	if err := os.WriteFile(path, []byte(rendered), 0o644); err != nil {
		return errorResult("write voice.md: %s", err)
	}

	return textResult("Saved voice profile to %s", path)
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
	var input struct {
		Name           string                     `json:"name"`
		Slug           string                     `json:"slug"`
		Tone           string                     `json:"tone"`
		Formality      string                     `json:"formality"`
		SentenceLength string                     `json:"sentence_length"`
		Register       string                     `json:"register"`
		BannedPhrases  []string                   `json:"banned_phrases"`
		Notes          string                     `json:"notes"`
		Goldens        []onboarding.GoldenExample `json:"goldens"`
	}
	if err := json.Unmarshal(args, &input); err != nil {
		return errorResult("invalid arguments: %s", err)
	}
	if strings.TrimSpace(input.Name) == "" {
		return errorResult("template name is required")
	}

	profile := onboarding.Profile{
		Kind:           onboarding.KindVoice,
		Tone:           input.Tone,
		Formality:      input.Formality,
		SentenceLength: input.SentenceLength,
		Register:       input.Register,
		BannedPhrases:  input.BannedPhrases,
		Notes:          input.Notes,
		Goldens:        input.Goldens,
	}

	if err := onboarding.Validate(profile); err != nil {
		return errorResult("%s", err.Error())
	}

	rendered, err := onboarding.Render(profile)
	if err != nil {
		return errorResult("render profile: %s", err)
	}

	// dec-20260719-b12539fa: name is template-only metadata, wrapped ahead
	// of the exact Profile content Render produces -- Profile itself never
	// carries a name field. slug is the stable key: pass it back unchanged
	// to update/rename a template in place, or omit it on first save to
	// derive one from the name.
	wrapped, err := onboarding.WrapTemplate(input.Name, rendered)
	if err != nil {
		return errorResult("wrap template: %s", err)
	}

	dir := onboarding.TemplatesDir(s.homeDir, onboarding.KindVoice)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return errorResult("create templates directory: %s", err)
	}
	slug := strings.TrimSpace(input.Slug)
	if slug == "" {
		slug = onboarding.Slug(input.Name)
	}
	path := filepath.Join(dir, slug+".md")
	if err := os.WriteFile(path, []byte(wrapped), 0o644); err != nil {
		return errorResult("write template: %s", err)
	}

	return textResult("Saved reusable template %q (slug: %s) to %s", input.Name, slug, path)
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
	dir := onboarding.TemplatesDir(s.homeDir, onboarding.KindVoice)
	entries, err := os.ReadDir(dir)
	if err != nil {
		if os.IsNotExist(err) {
			return textResult("No voice templates available yet.")
		}
		return errorResult("list templates: %s", err)
	}

	var b strings.Builder
	count := 0
	for _, e := range entries {
		if e.IsDir() || !strings.HasSuffix(e.Name(), ".md") {
			continue
		}
		content, err := os.ReadFile(filepath.Join(dir, e.Name()))
		if err != nil {
			return errorResult("read template %s: %s", e.Name(), err)
		}
		slug := strings.TrimSuffix(e.Name(), ".md")

		// dec-20260719-b12539fa: name lives in a wrapper ahead of the
		// profile content. A file saved without one (e.g. hand-authored,
		// or predating this decision) degrades gracefully -- it's still
		// listed, just under its slug as a fallback name, rather than
		// silently dropped from the list.
		name, profileContent, err := onboarding.UnwrapTemplate(string(content))
		if err != nil {
			name = slug
			profileContent = string(content)
		}

		fmt.Fprintf(&b, "### %s (slug: %s)\n\n%s\n\n---\n\n", name, slug, profileContent)
		count++
	}
	if count == 0 {
		return textResult("No voice templates available yet.")
	}
	return textResult("%s", b.String())
}
