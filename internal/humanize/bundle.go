// Package humanize validates and persists /ps-humanize output: a revised
// humanize.md for an existing piece, plus a changelog of what changed and
// why. Neither draft.md, staff-edit.md, nor sharpen.md (whichever ran) is
// touched -- each stage keeps its own file, same reasoning as staffedit
// and sharpen.
//
// Adapted from a prior writing POC's text-humanizer discipline: hedge-
// stack removal, stock-transition and throat-clearing strip, section/
// paragraph/bullet symmetry variance, sentence-rhythm variance -- these
// are genuine writing-quality tells regardless of whether the draft came
// from a person or a model, so removing them is legitimate editorial
// craft, not detection evasion.
//
// Deliberately excludes the parts of that POC's own aggressive framing
// (and a since-declined draft prompt with the same shape) that serve no
// reader and exist only to defeat AI-detection classifiers: no deliberate
// grammar-breaking, no leaving thoughts incomplete, no "embracing
// messiness" as a goal, no skipping explanations for effect. The POC's own
// default mode already rules these out ("no deliberate grammar-breaking");
// this package doesn't implement them even as an opt-in mode. The
// author-facing goal is "sounds like you, not the machine" -- voice
// authenticity, not statistical fingerprint scrambling -- and that goal
// does not require making the prose worse to read.
package humanize

import (
	"strings"

	"pysar/internal/draft"
	"pysar/internal/validation"
)

// Bundle is the structured humanize payload the agent assembles and the
// MCP tool validates before writing.
type Bundle struct {
	// PiecePath is the existing piece this pass belongs to, relative to the
	// project root -- resolved via research.ResolvePieceDir, same as
	// /ps-draft, /ps-staff-edit, and /ps-sharpen.
	PiecePath string `json:"piece_path"`

	// RevisedMD is the FULL revised text, written to humanize.md -- never
	// to draft.md, staff-edit.md, or sharpen.md, which stay untouched.
	// Replaces the whole humanize.md file on each call (wholesale, same
	// contract the earlier passes use), not a diff. Citation markers
	// ([^shortname]) must still resolve to real sources; never invented,
	// never resolved to a link here.
	RevisedMD string `json:"revised_md"`

	// Checks is >=1 one-line notes on what changed and why, e.g.
	// "[hedge-stack] dropped the weaker of two hedges on the compose
	// claim" or "[symmetry] varied the skip-list's third item so it
	// doesn't read as a template". Record "no changes needed" as a single
	// entry when a check genuinely required none.
	Checks []string `json:"checks"`

	// Mode is an optional one-word note on the edit's scope, e.g. "delta"
	// (surgical) or "rewrite" (structural) -- informational only, folded
	// into the changelog entry.
	Mode string `json:"mode,omitempty"`
}

// Validate reuses draft.ValidateContent for citation integrity (humanize.
// md's content is still, mechanically, a piece of draft prose -- the same
// citation/raw-URL rules apply, so this doesn't re-implement them a fifth
// time, and calling it with this package's own "revised_md" field name
// keeps error messages pointing at an argument that actually exists in
// save_humanize_bundle's schema) and validation.RequireNonEmptyChecks for
// the one check specific to this pass.
func Validate(b Bundle, validShortnames map[string]bool) error {
	var missing []string

	if strings.TrimSpace(b.PiecePath) == "" {
		missing = append(missing, "piece_path")
	}
	missing = append(missing, validation.RequireNonEmptyChecks(b.Checks)...)
	missing = append(missing, draft.ValidateContent("revised_md", b.RevisedMD, validShortnames)...)

	if len(missing) > 0 {
		return &validation.Error{Kind: "humanize bundle", Missing: missing}
	}
	return nil
}
