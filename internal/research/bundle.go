// Package research validates and persists /ps-research output: real,
// tiered sources added to an existing piece, or gathered standalone before
// a piece exists. Source tiers, the authority floor, citation hygiene, and
// the stop criterion are adapted from a prior writing POC's
// research-sourcing discipline (CL1) -- kept as proven, not watered down.
// Never touches thesis/killer_sections/counterintuitive: "we found
// sources; your take stayed yours."
package research

import (
	"fmt"
	"regexp"
	"strings"

	"pysar/internal/validation"
)

// Tier is a source's authority level. "Recognized-source" tier (primary or
// secondary) is what counts toward the authority floor -- community-only
// sources don't, matching research-sourcing.mdc's own definition.
type Tier string

const (
	TierPrimary   Tier = "primary"
	TierSecondary Tier = "secondary"
	TierCommunity Tier = "community"
)

// Source is one real, fetched reference. Shortname must be unique within
// the bundle (citation hygiene: writers cite via [^shortname]). RawExcerpt
// is the verbatim quoted text for raw/<shortname>.md -- Notes paraphrases,
// RawExcerpt quotes exactly, matching "quote exact wording in raw/
// excerpts; paraphrase only in sources.md summaries."
type Source struct {
	Shortname  string   `json:"shortname"`
	URL        string   `json:"url"`
	Tier       Tier     `json:"tier"`
	Accessed   string   `json:"accessed"` // YYYY-MM-DD, the day it was actually fetched
	KeyClaims  []string `json:"key_claims"`
	Notes      string   `json:"notes"`
	RawExcerpt string   `json:"raw_excerpt"`
}

// Competitor is one optional competitor-scan entry (--competitors=).
type Competitor struct {
	URL              string `json:"url"`
	Angle            string `json:"angle"`
	StrongestSection string `json:"strongest_section"`
	Gap              string `json:"gap"`
}

// Bundle is the structured research payload the agent assembles and the
// MCP tool validates before writing. Mechanical fields are validated in Go
// so the agent never hand-writes Markdown to disk.
type Bundle struct {
	// PiecePath, if set, anchors this research to an existing piece
	// directory (.pysar/pieces/<name>/) -- sources.md/raw/ get written
	// there and brief.md/angles.md get safely appended to, never
	// overwritten or replaced. Leave empty for standalone research (no
	// piece exists yet).
	PiecePath string `json:"piece_path,omitempty"`

	// Topic is required when PiecePath is empty (standalone mode) --
	// what the research is about, used to name the output directory.
	Topic string `json:"topic,omitempty"`

	// ExpertLens is the discipline/practitioner viewpoint judging source
	// authority tier for this topic -- reused from the piece's own
	// brief.md when PiecePath is set, determined fresh otherwise. This
	// is the agent-agnostic replacement for a prior POC's fixed
	// per-topic-family authority table.
	ExpertLens string `json:"expert_lens"`

	// TopicFamilyNote is set only when ExpertLens doesn't map cleanly to
	// an obvious authority standard for this topic -- names the mapping
	// chosen, so a reviewer can see the standard applied. Leave empty
	// when the mapping is obvious; most topics don't need this.
	TopicFamilyNote string `json:"topic_family_note,omitempty"`

	Sources     []Source     `json:"sources"`               // >=1 required; 6-12 typical, 8-14 without competitors (budget, not a hard cap)
	Competitors []Competitor `json:"competitors,omitempty"` // optional, only with --competitors=

	// KeyQuestionsAdditions/AnglesMisconceptions/AnglesContrarian are new
	// items this research surfaced. In piece-anchored mode they're
	// appended to the existing piece's sections, never replacing what's
	// there. In standalone mode they seed a fresh research-summary.md.
	// Every item here must be source-backed -- this is what research
	// adds, never a rewrite of the operator's own thesis or take.
	KeyQuestionsAdditions []string `json:"key_questions_additions,omitempty"`
	AnglesMisconceptions  []string `json:"angles_misconceptions,omitempty"`
	AnglesContrarian      []string `json:"angles_contrarian,omitempty"`
}

// ValidationError names every missing/invalid field at once.
type ValidationError = validation.Error

// shortnameRe enforces the kebab-case contract the tool schema already
// documents ("shortname: kebab-case, unique within this research pass").
// Beyond documentation, this is load-bearing: Shortname becomes a bare
// filename component (raw/<shortname>.md) with no further sanitization at
// write time, so anything outside [a-z0-9-] -- especially "/" or ".." --
// must never reach Validate as a pass.
var shortnameRe = regexp.MustCompile(`^[a-z0-9]+(-[a-z0-9]+)*$`)

// Validate checks completeness for a durable research write, including the
// authority floor (research-sourcing.mdc: >=60% of sources must be
// primary/secondary tier, not community-only) and citation hygiene
// (unique Shortname per source).
func Validate(b Bundle) error {
	var missing []string

	if strings.TrimSpace(b.PiecePath) == "" && strings.TrimSpace(b.Topic) == "" {
		missing = append(missing, "topic (required for standalone research -- no piece_path was given)")
	}
	if strings.TrimSpace(b.ExpertLens) == "" {
		missing = append(missing, "expert_lens")
	}

	if len(b.Sources) < 1 {
		missing = append(missing, "sources (need >= 1 -- research that found nothing is not a completed pass)")
	}

	shortnames := map[string]bool{}
	recognized := 0
	for i, s := range b.Sources {
		if strings.TrimSpace(s.Shortname) == "" {
			missing = append(missing, fmt.Sprintf("sources[%d].shortname", i))
		} else if !shortnameRe.MatchString(s.Shortname) {
			missing = append(missing, fmt.Sprintf("sources[%d].shortname (%q must be kebab-case: lowercase letters, digits, single hyphens -- it becomes a raw/<shortname>.md filename)", i, s.Shortname))
		} else if shortnames[s.Shortname] {
			missing = append(missing, fmt.Sprintf("sources[%d].shortname (duplicate %q -- every source needs a unique shortname)", i, s.Shortname))
		} else {
			shortnames[s.Shortname] = true
		}
		if strings.TrimSpace(s.URL) == "" {
			missing = append(missing, fmt.Sprintf("sources[%d].url", i))
		}
		switch s.Tier {
		case TierPrimary, TierSecondary, TierCommunity:
			if s.Tier != TierCommunity {
				recognized++
			}
		default:
			missing = append(missing, fmt.Sprintf("sources[%d].tier (must be primary|secondary|community)", i))
		}
		if strings.TrimSpace(s.Accessed) == "" {
			missing = append(missing, fmt.Sprintf("sources[%d].accessed (date actually fetched)", i))
		}
		if len(s.KeyClaims) < 1 {
			missing = append(missing, fmt.Sprintf("sources[%d].key_claims (need >= 1 -- what does this source actually support)", i))
		}
		if strings.TrimSpace(s.Notes) == "" {
			missing = append(missing, fmt.Sprintf("sources[%d].notes", i))
		}
		if strings.TrimSpace(s.RawExcerpt) == "" {
			missing = append(missing, fmt.Sprintf("sources[%d].raw_excerpt (verbatim quote for raw/%s.md)", i, s.Shortname))
		}
	}
	if len(b.Sources) > 0 && recognized*100 < 60*len(b.Sources) {
		missing = append(missing, fmt.Sprintf("sources (authority floor not met: %d/%d primary+secondary, need >= 60%% -- community-tier sources don't count toward the floor)", recognized, len(b.Sources)))
	}

	for i, c := range b.Competitors {
		if strings.TrimSpace(c.URL) == "" {
			missing = append(missing, fmt.Sprintf("competitors[%d].url", i))
		}
		if strings.TrimSpace(c.Angle) == "" {
			missing = append(missing, fmt.Sprintf("competitors[%d].angle", i))
		}
		if strings.TrimSpace(c.StrongestSection) == "" {
			missing = append(missing, fmt.Sprintf("competitors[%d].strongest_section", i))
		}
		if strings.TrimSpace(c.Gap) == "" {
			missing = append(missing, fmt.Sprintf("competitors[%d].gap", i))
		}
	}

	if len(missing) > 0 {
		return &ValidationError{Kind: "research bundle", Missing: missing}
	}
	return nil
}
