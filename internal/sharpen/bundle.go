// Package sharpen validates and persists /ps-sharpen output: a revised
// sharpen.md for an existing piece, plus a changelog of what changed and
// why. Neither draft.md nor staff-edit.md (if present) is touched -- each
// stage keeps its own file, same reasoning as staffedit.
//
// Adapted from a prior writing POC's engagement-editor discipline (opener
// hook, key-insight elevation, edge-section expansion, worked-example
// check, thesis/arc consistency) -- deliberately narrower than that POC's
// own SHARPEN phase, since this project's /ps-staff-edit already absorbed
// what that POC calls "consumability" (metaphor budget, paragraph density,
// rhythm) and opener concreteness. What's left here is the piece's reader-
// experience arc specifically: does the opening hook, are the pieces's own
// best findings elevated instead of buried, does the ending resolve what
// the opening promised -- staff-edit already established that the argument
// is sound; this pass checks that it lands as a read.
//
// Also deliberately absent: SEO/title/tag tuning (a prior writing POC's
// OPTIMIZE phase) and AI-detection-evasion rhythm-breaking (that POC's
// HUMANIZE phase) -- neither exists in this project, and the latter isn't
// a quality goal this project has any interest in building.
package sharpen

import (
	"fmt"
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

// Validate reuses draft.Validate for citation integrity (sharpen.md's
// content is still, mechanically, a piece of draft prose -- the same
// citation/raw-URL rules apply, so this doesn't re-implement them a fourth
// time) and adds the one check specific to this pass: at least one
// recorded check.
func Validate(b Bundle, validShortnames map[string]bool) error {
	var missing []string

	if len(b.Checks) == 0 {
		missing = append(missing, "checks (>=1 required -- even 'no changes needed' is worth recording as one entry)")
	}
	for i, c := range b.Checks {
		if strings.TrimSpace(c) == "" {
			missing = append(missing, fmt.Sprintf("checks[%d] (empty)", i))
		}
	}

	if err := draft.Validate(draft.Bundle{PiecePath: b.PiecePath, DraftMD: b.RevisedMD}, validShortnames); err != nil {
		if ve, ok := err.(*validation.Error); ok {
			missing = append(missing, ve.Missing...)
		} else {
			missing = append(missing, err.Error())
		}
	}

	if len(missing) > 0 {
		return &validation.Error{Kind: "sharpen bundle", Missing: missing}
	}
	return nil
}
