// Package sharpen validates and persists /ps-sharpen output: a revised
// sharpen.md for an existing piece, plus a changelog of what changed and
// why. Neither draft.md nor staff-edit.md (if present) is touched -- each
// stage keeps its own file, same reasoning as staffedit.
//
// Opener hook, key-insight elevation, edge-section expansion, worked-example
// check, thesis/arc consistency. Consumability (metaphor budget, paragraph
// density, rhythm) and opener concreteness already live in /ps-staff-edit;
// what's left here is the piece's reader-experience arc specifically: does
// the opening hook, are the piece's own best findings elevated instead of
// buried, does the ending resolve what the opening promised -- staff-edit
// already established that the argument is sound; this pass checks that it
// lands as a read.
//
// Also deliberately absent here: SEO/citation-resolution/title-for-search
// tuning -- that's internal/seo's job, a separate opt-in pass that runs
// after this one. AI-detection-evasion rhythm-breaking is not a quality
// goal this project has any interest in.
package sharpen

import (
	"strings"

	"pysar/internal/draft"
	"pysar/internal/validation"
)

// Bundle is the structured sharpen payload the agent assembles and the MCP
// tool validates before writing.
type Bundle struct {
	// PiecePath is the existing piece this pass belongs to, relative to the
	// project root -- resolved via research.ResolvePieceDir, same as
	// /ps-draft and /ps-staff-edit.
	PiecePath string `json:"piece_path"`

	// RevisedMD is the FULL revised text, written to sharpen.md -- never to
	// draft.md or staff-edit.md, which stay untouched. Replaces the whole
	// sharpen.md file on each call (wholesale, same contract the earlier
	// passes use), not a diff. Citation markers ([^shortname]) must still
	// resolve to real sources; never invented, never resolved to a link
	// here.
	RevisedMD string `json:"revised_md"`

	// Checks is >=1 one-line notes on what changed and why, e.g.
	// "[opener] tightened the hook so paragraph 1 promises a specific
	// payoff" or "[elevate] promoted the contrarian finding on X to its
	// own heading". Record "no changes needed" as a single entry when a
	// check genuinely required none.
	Checks []string `json:"checks"`

	// Mode is an optional one-word note on the edit's scope, e.g. "delta"
	// (surgical) or "rewrite" (structural) -- informational only, folded
	// into the changelog entry.
	Mode string `json:"mode,omitempty"`
}

// Validate reuses draft.ValidateContent for citation integrity (sharpen.
// md's content is still, mechanically, a piece of draft prose -- the same
// citation/raw-URL rules apply, so this doesn't re-implement them a fourth
// time, and calling it with this package's own "revised_md" field name
// keeps error messages pointing at an argument that actually exists in
// save_sharpen_bundle's schema) and validation.RequireNonEmptyChecks for
// the one check specific to this pass.
func Validate(b Bundle, validShortnames map[string]bool) error {
	var missing []string

	if strings.TrimSpace(b.PiecePath) == "" {
		missing = append(missing, "piece_path")
	}
	missing = append(missing, validation.RequireNonEmptyChecks(b.Checks)...)
	missing = append(missing, draft.ValidateContent("revised_md", b.RevisedMD, validShortnames)...)

	if len(missing) > 0 {
		return &validation.Error{Kind: "sharpen bundle", Missing: missing}
	}
	return nil
}
