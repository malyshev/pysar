// Package onboarding defines the schema for author content profiles (voice,
// style, and future kinds) and their completeness validation. The Go layer
// only defines schema and validation/completeness checks -- it never conducts
// the questionnaire conversation itself; that happens on the agentic surface
// (dec-20260718-5570376d).
package onboarding

import (
	"fmt"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

// ProfileKind names what a Profile describes.
type ProfileKind string

const (
	KindVoice ProfileKind = "voice"
	KindStyle ProfileKind = "style"
)

// GoldenExample is a concrete reference passage an AI can be steered by,
// alongside a short caption naming what it demonstrates.
type GoldenExample struct {
	Label string `json:"label"`
	Text  string `json:"text"`
}

// MinGoldens is the minimum number of golden examples a complete Profile
// must carry (dec-20260718-ab150a73 -- goldens are not optional).
const MinGoldens = 3

// Profile is the shared schema shape for any onboarding-produced content
// profile: structured attributes plus required golden examples
// (dec-20260718-ab150a73). VoiceProfile and StyleProfile are both this same
// type, distinguished only by Kind -- no divergent format
// (dec-20260718-ab150a73 invariant).
type Profile struct {
	Kind           ProfileKind     `json:"kind"`
	Tone           string          `json:"tone"`
	Formality      string          `json:"formality"`
	SentenceLength string          `json:"sentence_length"`
	Register       string          `json:"register"`
	BannedPhrases  []string        `json:"banned_phrases,omitempty"`
	Notes          string          `json:"notes,omitempty"`
	Goldens        []GoldenExample `json:"goldens"`
}

// ValidationError reports every missing piece of a Profile at once, not just
// the first -- so an author or skill can fix a profile in one pass.
type ValidationError struct {
	Missing []string
}

func (e *ValidationError) Error() string {
	return fmt.Sprintf("profile incomplete: missing %v", e.Missing)
}

// Validate checks that p satisfies the completeness bar: every structured
// attribute field is non-empty and at least MinGoldens golden examples, each
// with non-empty text, are present.
func Validate(p Profile) error {
	var missing []string

	if p.Tone == "" {
		missing = append(missing, "tone")
	}
	if p.Formality == "" {
		missing = append(missing, "formality")
	}
	if p.SentenceLength == "" {
		missing = append(missing, "sentence_length")
	}
	if p.Register == "" {
		missing = append(missing, "register")
	}
	if len(p.Goldens) < MinGoldens {
		missing = append(missing, fmt.Sprintf("goldens (need >= %d, have %d)", MinGoldens, len(p.Goldens)))
	}
	for i, g := range p.Goldens {
		if g.Text == "" {
			missing = append(missing, fmt.Sprintf("goldens[%d].text", i))
		}
	}

	if len(missing) > 0 {
		return &ValidationError{Missing: missing}
	}
	return nil
}

// frontmatter is Profile's scalar/attribute fields only (Goldens render as
// Markdown body sections, not frontmatter) -- kept as its own type so yaml.v3
// owns escaping/quoting for author-supplied text, never hand-rolled string
// concatenation that could produce invalid YAML.
type frontmatter struct {
	Kind           ProfileKind `yaml:"kind"`
	Tone           string      `yaml:"tone"`
	Formality      string      `yaml:"formality"`
	SentenceLength string      `yaml:"sentence_length"`
	Register       string      `yaml:"register"`
	BannedPhrases  []string    `yaml:"banned_phrases,omitempty"`
	Notes          string      `yaml:"notes,omitempty"`
}

// TemplatesDir is the cross-project, host-agnostic directory pysar stores
// reusable Kind templates in: ~/.pysar/templates/<kind>/ (dec-20260719-3e36577e).
// A template file at this path is a plain Profile rendered exactly like any
// onboarding-produced profile -- same schema, same Render/Validate, no
// divergent template-only format.
func TemplatesDir(home string, kind ProfileKind) string {
	return filepath.Join(home, ".pysar", "templates", string(kind))
}

// Slug turns an arbitrary template display name (e.g. "Measured plain
// English -- speakable, understated, general audience") into a safe
// filename stem: lowercased, runs of non [a-z0-9] characters collapsed to a
// single '-', leading/trailing '-' trimmed. Never returns an empty string.
func Slug(name string) string {
	var b strings.Builder
	lastDash := true // suppresses a leading '-'
	for _, r := range strings.ToLower(name) {
		switch {
		case r >= 'a' && r <= 'z' || r >= '0' && r <= '9':
			b.WriteRune(r)
			lastDash = false
		case !lastDash:
			b.WriteRune('-')
			lastDash = true
		}
	}
	slug := strings.TrimRight(b.String(), "-")
	if slug == "" {
		return "template"
	}
	return slug
}

// Render serializes a Profile into the Markdown+YAML-frontmatter format
// ps-voice (and future ps-*) skills persist -- structured fields in
// frontmatter, goldens as real Markdown sections (dec-20260718-ab150a73).
// Callers should Validate before Render; Render does not itself enforce
// completeness.
func Render(p Profile) (string, error) {
	fm := frontmatter{
		Kind:           p.Kind,
		Tone:           p.Tone,
		Formality:      p.Formality,
		SentenceLength: p.SentenceLength,
		Register:       p.Register,
		BannedPhrases:  p.BannedPhrases,
		Notes:          p.Notes,
	}
	fmBytes, err := yaml.Marshal(fm)
	if err != nil {
		return "", fmt.Errorf("render frontmatter: %w", err)
	}

	var b strings.Builder
	b.WriteString("---\n")
	b.Write(fmBytes)
	b.WriteString("---\n\n# Voice Profile\n")
	for _, g := range p.Goldens {
		b.WriteString(fmt.Sprintf("\n## Golden: %s\n\n%s\n", g.Label, g.Text))
	}
	return b.String(), nil
}
