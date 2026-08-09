// Package humanize validates and persists /ps-humanize output: a revised
// humanize.md for an existing piece, plus a changelog of what changed and
// why. Neither draft.md, staff-edit.md, sharpen.md, nor seo.md (whichever
// ran) is touched -- each stage keeps its own file, same reasoning as
// staffedit and sharpen.
//
// Hedge-stack removal, stock-transition and throat-clearing strip, section/
// paragraph/bullet symmetry variance, sentence-rhythm variance -- these are
// genuine writing-quality tells regardless of whether the draft came from a
// person or a model, so removing them is legitimate editorial craft, not
// detection evasion.
//
// Deliberately excludes techniques (including a since-declined draft prompt
// with that shape) that serve no reader and exist only to defeat
// AI-detection classifiers: no deliberate grammar-breaking, no leaving
// thoughts incomplete, no "embracing messiness" as a goal, no skipping
// explanations for effect -- not even as an opt-in. The author-facing goal
// is "sounds like you, not the machine" -- voice authenticity, not
// statistical fingerprint scrambling -- and that goal does not require
// making the prose worse to read.
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
	// to draft.md, staff-edit.md, sharpen.md, or seo.md, which stay
	// untouched. Replaces the whole humanize.md file on each call
	// (wholesale, same contract the earlier passes use), not a diff. If
	// /ps-seo ran, its [anchor](url) links must survive untouched -- edit
	// prose around them, never inside them. If /ps-seo did not run,
	// citation markers ([^shortname]) must still resolve to real sources;
	// never invented, never resolved to a link here (that's /ps-seo's job).
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

// Validate checks citation integrity for humanize.md. When /ps-seo ran
// first, humanize's input already has [^shortname] markers resolved into
// real [anchor](url) links (internal/seo's own inverted rule) -- those
// links must survive untouched, so this cannot simply reuse
// draft.ValidateContent, which treats any raw URL (including one inside a
// resolved link) as a violation. validShortnames still allows an
// unresolved [^shortname] marker when /ps-seo did not run (same rule
// draft.ValidateContent enforces); validSourceURLs allows a resolved
// [anchor](url) link when its url matches a source this piece's research
// actually recorded. A bare URL belonging to neither form is still
// rejected -- humanize.md is never the place a citation gets invented or
// silently downgraded to plain text.
func Validate(b Bundle, validShortnames, validSourceURLs map[string]bool) error {
	var missing []string

	if strings.TrimSpace(b.PiecePath) == "" {
		missing = append(missing, "piece_path")
	}
	missing = append(missing, validation.RequireNonEmptyChecks(b.Checks)...)
	missing = append(missing, validateContent(b.RevisedMD, validShortnames, validSourceURLs)...)

	if len(missing) > 0 {
		return &validation.Error{Kind: "humanize bundle", Missing: missing}
	}
	return nil
}

func validateContent(content string, validShortnames, validSourceURLs map[string]bool) []string {
	body := strings.TrimSpace(content)
	if body == "" {
		return []string{"revised_md"}
	}

	missing := draft.ValidateCitations("revised_md", body, validShortnames)

	linkFreeBody := draft.MarkdownLinkRe.ReplaceAllStringFunc(body, func(link string) string {
		m := draft.MarkdownLinkRe.FindStringSubmatch(link)
		if len(m) == 2 && validSourceURLs[m[1]] {
			return "" // a resolved link matching a real source is not a raw-URL violation
		}
		return link // an unmatched link's URL still gets flagged below, same as a bare one
	})
	if draft.HasRawURL(linkFreeBody) {
		missing = append(missing, "revised_md (raw URL found in prose that isn't a [^shortname] citation or a [anchor](url) link matching a source this piece's research recorded)")
	}

	return missing
}
