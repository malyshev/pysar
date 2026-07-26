// Package staffedit validates and persists /ps-staff-edit output: a revised
// staff-edit.md for an existing piece, plus a changelog of what changed and
// why. draft.md itself is never touched -- it stays the untouched first-
// draft artifact so the two can be compared (or shown side by side in a
// future UI), the same reasoning that keeps intake's brief.md/outline.md/
// angles.md as separate files instead of one another. Adapted from a prior
// writing POC's staff-editor discipline
// (stakes/thesis pressure, signal-vs-noise against the brief, failure
// modes and honest scope, close memorability, technical sanity) -- folded
// in here rather than split into a separate downstream engagement pass,
// since this project's pipeline doesn't have one (no SHARPEN/HUMANIZE
// equivalent exists): readability -- sentence rhythm, one idea per
// paragraph, cutting the flat step-by-step "manual" cadence where the
// piece isn't a step-by-step manual -- is one of this pass's own checks,
// not a separate command.
package staffedit

import (
	"fmt"
	"strings"

	"pysar/internal/draft"
	"pysar/internal/validation"
)

// Bundle is the structured staff-edit payload the agent assembles and the
// MCP tool validates before writing.
type Bundle struct {
	// PiecePath is the existing piece this edit belongs to, relative to the
	// project root -- resolved via research.ResolvePieceDir, same as
	// /ps-draft.
	PiecePath string `json:"piece_path"`

	// RevisedMD is the FULL revised draft, written to staff-edit.md --
	// never to draft.md, which stays untouched. This replaces the whole
	// staff-edit.md file on each call (wholesale, same contract /ps-draft
	// itself uses for draft.md), not a diff. Citation markers ([^shortname])
	// must still resolve to real sources; never invented, never resolved
	// to a link here.
	RevisedMD string `json:"revised_md"`

	// Checks is >=1 one-line notes on what changed and why, e.g.
	// "[stakes] named the specific failure mode in paragraph 2" or
	// "[readability] split the three-clause sentence in the compose
	// section". Record "no changes needed" as a single entry when a check
	// genuinely required none -- an edit pass that logged nothing isn't a
	// completed pass, but "nothing needed changing, here's why" is.
	Checks []string `json:"checks"`

	// Mode is an optional one-word note on the edit's scope, e.g. "delta"
	// (surgical) or "rewrite" (structural) -- informational only, folded
	// into the changelog entry.
	Mode string `json:"mode,omitempty"`
}

// Validate reuses draft.Validate for citation integrity (staff-edit.md's
// content is still, mechanically, a piece of draft prose -- the same
// citation/raw-URL rules apply, so this doesn't re-implement them a third
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
		return &validation.Error{Kind: "staff-edit bundle", Missing: missing}
	}
	return nil
}
