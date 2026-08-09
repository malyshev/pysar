// Package intake validates and persists intake scaffolding for a piece
// (dec-20260725-35fa2d24). Content shape is inspired by a prior writing POC
// (brief/outline/angles/sources-stub) at CL1 only — Author fields and
// --length caps are deliberately absent.
package intake

import (
	"fmt"
	"strings"
	"unicode"
	"unicode/utf8"

	"pysar/internal/validation"
)

// EntryMode records how the thesis entered intake.
type EntryMode string

const (
	EntryIdea  EntryMode = "idea"
	EntryDraft EntryMode = "draft"
)

// POVSource records where the thesis came from (no research-synthesized at intake).
type POVSource string

const (
	POVAuthorPractitioner POVSource = "author-practitioner"
	POVAuthorDraft        POVSource = "author-draft"
)

// KillerSection is a section where the piece has a genuine informational
// edge over existing coverage. Edge and Example are required, not optional
// color -- a bare title gives a drafter nothing to write from. Shape
// inspired by a prior writing POC's brief-assembly contract (CL1), adapted
// here without its author/length fields.
type KillerSection struct {
	Title   string `json:"title"`
	Edge    string `json:"edge"`    // why this beats existing/competing coverage
	Example string `json:"example"` // a concrete worked example to include
}

// CounterintuitiveFinding is a myth-breaking claim with its contradiction
// named explicitly, so it lands as a real section or punchline instead of a
// vague "this is surprising" line.
type CounterintuitiveFinding struct {
	Claim         string `json:"claim"`
	Contradiction string `json:"contradiction"` // what makes this claim counterintuitive, stated
}

// Source is one real, actually-fetched reference the agent used to ground a
// technical claim during intake (dec-20260725-35fa2d24, amended). Never a
// citation the agent didn't verify -- Validate rejects an empty URL.
// Grounding facts and examples this way is distinct from letting research
// override the operator's own thesis/angle, which stays forbidden.
type Source struct {
	URL  string `json:"url"`
	Note string `json:"note"` // one line: what this source confirmed
}

// Bundle is the structured intake payload the agent assembles and the MCP
// tool validates before writing. Mechanical fields are validated in Go so
// the agent never hand-writes Markdown to disk.
type Bundle struct {
	// Name is an optional readable prefix for the piece directory. A random
	// suffix is always appended on write (see write.go) — directory naming
	// and collisions are storage mechanics, never a decision surfaced to
	// the author. Leave empty to derive one from Idea.
	Name        string    `json:"name"`
	Idea        string    `json:"idea"`
	EntryMode   EntryMode `json:"entry_mode"`
	POVSource   POVSource `json:"pov_source"`
	Restatement string    `json:"restatement"` // one line: audience, register, topic-scope
	Audience    string    `json:"audience"`
	Register    string    `json:"register"`
	// ExpertLens is the discipline/practitioner viewpoint whose rigor
	// standard governs the piece's claims -- distinct from Audience (who
	// you write FOR). E.g. "experimental science / methodology" for a
	// science piece, "horticulture" for a gardening piece. Drives both
	// content depth (Step 5) and what counts as a claim worth grounding
	// (Step 4): a claim a practitioner in this discipline would expect
	// verified, not just one the agent feels confident about.
	ExpertLens       string                    `json:"expert_lens"`
	TopicScope       string                    `json:"topic_scope"`
	Thesis           string                    `json:"thesis"`
	Promise          string                    `json:"promise"`
	KillerSections   []KillerSection           `json:"killer_sections"`  // 1–2 required
	Counterintuitive []CounterintuitiveFinding `json:"counterintuitive"` // 1–3 required
	KeyQuestions     []string                  `json:"key_questions"`
	NonGoals         []string                  `json:"non_goals"`
	OutlineMD        string                    `json:"outline_md"`
	AnglesMD         string                    `json:"angles_md"`
	DraftPath        string                    `json:"draft_path,omitempty"`
	// Sources are real, fetched references used to ground technical claims
	// (Step 4). Optional and empty by default -- most topics don't need
	// grounding; this is not a lesser outcome than having sources.
	Sources []Source `json:"sources,omitempty"`
}

// ValidationError names every missing/invalid field at once.
type ValidationError = validation.Error

// Degenerate reports mechanical cases where intake must block (empty,
// whitespace-only, a single word-like unit, or too few extractable letters).
// Letter detection is Unicode-wide (dec-20260809-degenerate-unicode-i18n-5008f0e6),
// not ASCII a–z. Word-like units are whitespace-separated letter runs for
// alphabetic scripts, and each Han/Hiragana/Katakana/Hangul letter for
// scriptio continua — so Japanese/Chinese without spaces are not false
// rejects. Broader semantic emptiness / self-contradiction remains an agent
// judgment before calling save_intake_bundle.
func Degenerate(idea string) bool {
	trimmed := strings.TrimSpace(idea)
	if trimmed == "" {
		return true
	}
	if utf8.RuneCountInString(trimmed) < 3 {
		return true
	}
	letters := 0
	for _, r := range trimmed {
		if unicode.IsLetter(r) {
			letters++
		}
	}
	if letters < 3 {
		return true
	}
	return ideaUnits(trimmed) < 2
}

// scriptioContinuaLetter is true for scripts that commonly write without
// spaces between words. Each such letter counts as its own idea unit.
func scriptioContinuaLetter(r rune) bool {
	return unicode.Is(unicode.Han, r) ||
		unicode.Is(unicode.Hiragana, r) ||
		unicode.Is(unicode.Katakana, r) ||
		unicode.Is(unicode.Hangul, r)
}

// ideaUnits counts mechanical word-like units (dec-20260809-degenerate-unicode-i18n-5008f0e6).
func ideaUnits(s string) int {
	units := 0
	inAlphabetic := false
	for _, r := range s {
		switch {
		case scriptioContinuaLetter(r):
			units++
			inAlphabetic = false
		case unicode.IsLetter(r):
			if !inAlphabetic {
				units++
				inAlphabetic = true
			}
		default:
			inAlphabetic = false
		}
	}
	return units
}

// Validate checks completeness for a durable intake write. Name is
// intentionally not required — an empty Name is auto-allocated at write
// time (write.go), not treated as a missing field.
func Validate(b Bundle) error {
	var missing []string
	if Degenerate(b.Idea) && b.EntryMode != EntryDraft {
		missing = append(missing, "idea (degenerate or empty — ask only when input has no extractable content)")
	}
	if b.EntryMode != EntryIdea && b.EntryMode != EntryDraft {
		missing = append(missing, "entry_mode (idea|draft)")
	}
	if b.POVSource != POVAuthorPractitioner && b.POVSource != POVAuthorDraft {
		missing = append(missing, "pov_source (author-practitioner|author-draft)")
	}
	for _, pair := range []struct {
		name, val string
	}{
		{"restatement", b.Restatement},
		{"audience", b.Audience},
		{"register", b.Register},
		{"expert_lens", b.ExpertLens},
		{"topic_scope", b.TopicScope},
		{"thesis", b.Thesis},
		{"promise", b.Promise},
		{"outline_md", b.OutlineMD},
		{"angles_md", b.AnglesMD},
	} {
		if strings.TrimSpace(pair.val) == "" {
			missing = append(missing, pair.name)
		}
	}
	if len(b.KillerSections) < 1 {
		missing = append(missing, "killer_sections (need >= 1)")
	}
	for i, k := range b.KillerSections {
		if strings.TrimSpace(k.Title) == "" {
			missing = append(missing, fmt.Sprintf("killer_sections[%d].title", i))
		}
		if strings.TrimSpace(k.Edge) == "" {
			missing = append(missing, fmt.Sprintf("killer_sections[%d].edge (why this beats existing coverage)", i))
		}
		if strings.TrimSpace(k.Example) == "" {
			missing = append(missing, fmt.Sprintf("killer_sections[%d].example (a concrete worked example)", i))
		}
	}
	if len(b.Counterintuitive) < 1 {
		missing = append(missing, "counterintuitive (need >= 1)")
	}
	for i, c := range b.Counterintuitive {
		if strings.TrimSpace(c.Claim) == "" {
			missing = append(missing, fmt.Sprintf("counterintuitive[%d].claim", i))
		}
		if strings.TrimSpace(c.Contradiction) == "" {
			missing = append(missing, fmt.Sprintf("counterintuitive[%d].contradiction (name what makes it counterintuitive)", i))
		}
	}
	if !strings.Contains(strings.ToLower(b.AnglesMD), "contrarian") {
		missing = append(missing, "angles_md (must include a Contrarian section)")
	}
	if len(b.KeyQuestions) < 1 {
		missing = append(missing, "key_questions (need >= 1 — not optional, prevents scope drift once drafting starts)")
	}
	if len(b.NonGoals) < 1 {
		missing = append(missing, "non_goals (need >= 1 — not optional, prevents scope inflation once drafting starts)")
	}
	for i, s := range b.Sources {
		if strings.TrimSpace(s.URL) == "" {
			missing = append(missing, fmt.Sprintf("sources[%d].url (a source with no URL was not actually fetched)", i))
		}
		if strings.TrimSpace(s.Note) == "" {
			missing = append(missing, fmt.Sprintf("sources[%d].note (what did this source confirm)", i))
		}
	}
	if len(missing) > 0 {
		return &ValidationError{Kind: "intake bundle", Missing: missing}
	}
	return nil
}
