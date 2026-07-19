package mcpserver

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"pysar/internal/onboarding"
)

// profileProperties is the JSON-Schema property set shared by every
// save_<kind>_profile / save_<kind>_template tool -- all persist the exact
// same Profile shape (dec-20260718-ab150a73, dec-20260719-a1ac8959).
// Each tool's own schema declares which of these are required for its Kind.
func profileProperties() map[string]interface{} {
	return map[string]interface{}{
		"tone":            map[string]string{"type": "string", "description": "How the author wants to come across, e.g. 'warm, direct, a little dry'. Voice-only -- leave empty for a style profile."},
		"formality":       map[string]string{"type": "string", "description": "Where on the conversational-to-formal range, and how. Voice-only -- leave empty for a style profile."},
		"sentence_length": map[string]string{"type": "string", "description": "Sentence length/rhythm preference. Voice-only -- leave empty for a style profile."},
		"register":        map[string]string{"type": "string", "description": "Vocabulary register, e.g. 'plain, international-standard English'. Voice-only -- leave empty for a style profile."},
		"banned_phrases": map[string]interface{}{
			"type":        "array",
			"items":       map[string]string{"type": "string"},
			"description": "Optional -- phrases/constructions the author never wants to see",
		},
		"notes": map[string]string{"type": "string", "description": "Optional -- texture the structured fields don't capture"},
		"rules": map[string]interface{}{
			"type":        "array",
			"items":       map[string]string{"type": "string"},
			"description": fmt.Sprintf("Short, individually actionable imperative statements, e.g. 'Prefer active voice'. Required for a style profile (at least %d); leave empty for a voice profile.", onboarding.MinRules),
		},
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

// profileInput is the wire shape every save_<kind>_profile tool accepts.
type profileInput struct {
	Tone           string                     `json:"tone"`
	Formality      string                     `json:"formality"`
	SentenceLength string                     `json:"sentence_length"`
	Register       string                     `json:"register"`
	BannedPhrases  []string                   `json:"banned_phrases"`
	Notes          string                     `json:"notes"`
	Rules          []string                   `json:"rules"`
	Goldens        []onboarding.GoldenExample `json:"goldens"`
}

func (in profileInput) toProfile(kind onboarding.ProfileKind) onboarding.Profile {
	return onboarding.Profile{
		Kind:           kind,
		Tone:           in.Tone,
		Formality:      in.Formality,
		SentenceLength: in.SentenceLength,
		Register:       in.Register,
		BannedPhrases:  in.BannedPhrases,
		Notes:          in.Notes,
		Rules:          in.Rules,
		Goldens:        in.Goldens,
	}
}

// saveProfileToPath validates and persists a profile of the given kind to
// path -- shared by save_voice_profile and save_style_profile
// (dec-20260719-a1ac8959 prediction: style tools mirror voice tools exactly,
// factored once here rather than duplicated per Kind).
func (s *Server) saveProfileToPath(args json.RawMessage, kind onboarding.ProfileKind, path string) callToolResult {
	var in profileInput
	if err := json.Unmarshal(args, &in); err != nil {
		return errorResult("invalid arguments: %s", err)
	}

	profile := in.toProfile(kind)
	if err := onboarding.Validate(profile); err != nil {
		return errorResult("%s", err.Error())
	}

	rendered, err := onboarding.Render(profile)
	if err != nil {
		return errorResult("render profile: %s", err)
	}

	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return errorResult("create directory: %s", err)
	}
	if err := os.WriteFile(path, []byte(rendered), 0o644); err != nil {
		return errorResult("write profile: %s", err)
	}

	return textResult("Saved %s profile to %s", kind, path)
}

// templateInput is the wire shape every save_<kind>_template tool accepts:
// a profileInput plus the template-only name/slug pair
// (dec-20260719-b12539fa).
type templateInput struct {
	Name string `json:"name"`
	Slug string `json:"slug"`
	profileInput
}

// saveTemplate validates a profile of the given kind and persists it as a
// named, reusable, cross-project template -- shared by save_voice_template
// and save_style_template.
func (s *Server) saveTemplate(args json.RawMessage, kind onboarding.ProfileKind) callToolResult {
	var in templateInput
	if err := json.Unmarshal(args, &in); err != nil {
		return errorResult("invalid arguments: %s", err)
	}
	if strings.TrimSpace(in.Name) == "" {
		return errorResult("template name is required")
	}

	profile := in.profileInput.toProfile(kind)
	if err := onboarding.Validate(profile); err != nil {
		return errorResult("%s", err.Error())
	}

	rendered, err := onboarding.Render(profile)
	if err != nil {
		return errorResult("render profile: %s", err)
	}

	wrapped, err := onboarding.WrapTemplate(in.Name, rendered)
	if err != nil {
		return errorResult("wrap template: %s", err)
	}

	dir := onboarding.TemplatesDir(s.homeDir, kind)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return errorResult("create templates directory: %s", err)
	}
	slug := strings.TrimSpace(in.Slug)
	if slug == "" {
		slug = onboarding.Slug(in.Name)
	} else {
		slug = onboarding.Slug(slug)
	}
	path := filepath.Join(dir, slug+".md")
	if err := os.WriteFile(path, []byte(wrapped), 0o644); err != nil {
		return errorResult("write template: %s", err)
	}

	return textResult("Saved reusable template %q (slug: %s) to %s", in.Name, slug, path)
}

// listTemplates returns every reusable template of the given kind, name and
// slug and full content inline -- shared by list_voice_templates and
// list_style_templates. Sidesteps Glob/Bash path-expansion entirely
// (dec-20260719-3e36577e's refuted prediction, fixed once here for every
// Kind rather than per-Kind).
func (s *Server) listTemplates(kind onboarding.ProfileKind) callToolResult {
	dir := onboarding.TemplatesDir(s.homeDir, kind)
	entries, err := os.ReadDir(dir)
	if err != nil {
		if os.IsNotExist(err) {
			return textResult("No %s templates available yet.", kind)
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

		name, profileContent, err := onboarding.UnwrapTemplate(string(content))
		if err != nil {
			name = slug
			profileContent = string(content)
		}

		fmt.Fprintf(&b, "### %s (slug: %s)\n\n%s\n\n---\n\n", name, slug, profileContent)
		count++
	}
	if count == 0 {
		return textResult("No %s templates available yet.", kind)
	}
	return textResult("%s", b.String())
}
