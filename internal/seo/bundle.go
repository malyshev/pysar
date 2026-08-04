// Package seo validates and persists /ps-seo output: a revised seo.md for
// an existing piece, a discoverability checklist (seo-checklist.md), and a
// changelog of what changed and why. Neither draft.md, staff-edit.md, nor
// sharpen.md (whichever ran) is touched -- each stage keeps its own file,
// same reasoning as staffedit, sharpen, and humanize.
//
// This is the pass every earlier stage's own doc comments named and
// deliberately deferred: internal/draft says resolving [^shortname]
// markers into inline links "is a later, optional, separate concern this
// package doesn't own"; internal/sharpen names "SEO/title/tag tuning" as
// deliberately absent from that pass specifically because it belongs here
// instead. Where every earlier stage REQUIRES [^shortname] markers to
// still resolve against sources.md (citations stay unresolved on purpose,
// for a later channel-specific pass to handle), this package inverts that
// rule: a completed seo.md must contain ZERO [^shortname] markers -- every
// one of them has been resolved into a real [anchor](url) link, or the
// pass isn't done.
//
// Adapted from a prior writing POC's SEO-optimizer discipline (citation
// resolution, title/subtitle as two distinct jobs -- CTR vs. read-
// completion --, scannability, outbound-link authority floor) but not
// mirrored: that POC's checklist is Medium-specific (a rigid 5-slot tag
// strategy, Medium's own SEO/curator model). Pysar targets more than one
// platform, so the checklist here is platform-neutral -- tags and meta
// fields the author can hand to whatever platform they're posting to,
// with no forced slot count or platform-specific field names.
package seo

import (
	"fmt"
	"regexp"
	"strings"

	"pysar/internal/draft"
	"pysar/internal/validation"
)

// Bundle is the structured SEO payload the agent assembles and the MCP
// tool validates before writing.
type Bundle struct {
	// PiecePath is the existing piece this pass belongs to, relative to the
	// project root -- resolved via research.ResolvePieceDir, same as every
	// other piece-anchored pass.
	PiecePath string `json:"piece_path"`

	// RevisedMD is the FULL revised text, written to seo.md -- never to
	// draft.md, staff-edit.md, or sharpen.md, which stay untouched.
	// Replaces the whole seo.md file on each call (wholesale, same
	// contract the earlier passes use), not a diff. Unlike every earlier
	// stage, citation markers must NOT remain: every [^shortname] must
	// already be resolved into a real [anchor](url) link pointing at a URL
	// this piece's research actually recorded -- never invented.
	RevisedMD string `json:"revised_md"`

	// Checks is >=1 one-line notes on what changed and why, e.g.
	// "[citation] resolved [^src-a] to an inline link on 'the retry
	// budget'" or "[title] tightened the title for a concrete keyword
	// match". Record "no changes needed" as a single entry when a check
	// genuinely required none.
	Checks []string `json:"checks"`

	// Mode is an optional one-word note on the edit's scope, e.g. "delta"
	// (surgical) or "rewrite" (structural) -- informational only, folded
	// into the changelog entry.
	Mode string `json:"mode,omitempty"`

	// Tags is >=1 discovery keyword/topic tag for whatever platform the
	// piece is headed to -- no fixed slot count or Medium-specific
	// taxonomy, unlike the reference implementation this is adapted from.
	Tags []string `json:"tags"`

	// MetaTitle is the search/share title -- distinct from the article's
	// own `# Title` line, tuned for a search-results or link-preview
	// context rather than the feed/reader context the article title
	// itself is tuned for (same "title vs subtitle, two different jobs"
	// distinction the reference implementation draws, extended one layer
	// further: article title vs. meta title are a third job again).
	MetaTitle string `json:"meta_title"`

	// MetaDescription is the search/share summary -- concrete, distinct
	// from the article's own subtitle, written for a reader who hasn't
	// opened the piece yet.
	MetaDescription string `json:"meta_description"`

	// URLSlug is a kebab-case, keyword-bearing suggested URL path segment.
	URLSlug string `json:"url_slug"`
}

// slugRe is the one regex genuinely specific to this package -- every
// other content-pattern check (citation markers, resolved links, bare
// URLs, code-fence stripping) reuses internal/draft's exported helpers
// instead of re-declaring near-duplicate regexes that can silently drift
// from the originals.
var slugRe = regexp.MustCompile(`^[a-z0-9]+(-[a-z0-9]+)*$`)

// Validate checks that this pass actually finished its one job (zero
// citation markers remain, every resolved link is real) plus the
// checklist fields a discoverability pass exists to produce.
// validSourceURLs is the piece's current research output
// (research.LoadSourceURLs) -- empty when no research has run, which
// correctly makes any resolved link invalid rather than silently accepted.
func Validate(b Bundle, validSourceURLs map[string]bool) error {
	var missing []string

	if strings.TrimSpace(b.PiecePath) == "" {
		missing = append(missing, "piece_path")
	}
	missing = append(missing, validation.RequireNonEmptyChecks(b.Checks)...)
	missing = append(missing, validateContent(b.RevisedMD, validSourceURLs)...)

	if len(b.Tags) == 0 {
		missing = append(missing, "tags (>=1 required)")
	}
	for i, t := range b.Tags {
		if strings.TrimSpace(t) == "" {
			missing = append(missing, fmt.Sprintf("tags[%d] (empty)", i))
		}
	}
	if strings.TrimSpace(b.MetaTitle) == "" {
		missing = append(missing, "meta_title")
	}
	if strings.TrimSpace(b.MetaDescription) == "" {
		missing = append(missing, "meta_description")
	}
	if strings.TrimSpace(b.URLSlug) == "" {
		missing = append(missing, "url_slug")
	} else if !slugRe.MatchString(b.URLSlug) {
		missing = append(missing, fmt.Sprintf("url_slug (%q must be kebab-case: lowercase letters, digits, single hyphens)", b.URLSlug))
	}

	if len(missing) > 0 {
		return &validation.Error{Kind: "seo bundle", Missing: missing}
	}
	return nil
}

// validateContent enforces this pass's inverted citation rule (zero
// [^shortname] markers may remain; every resolved link must trace to a
// real recorded source URL) plus the same no-bare-URL-outside-a-link
// discipline every earlier stage's prose already follows -- built on
// internal/draft's exported primitives (FindCitationMarkers,
// MarkdownLinkRe, HasRawURL, StripCode) rather than a second, independently
// maintained copy of the same regex family.
func validateContent(content string, validSourceURLs map[string]bool) []string {
	body := strings.TrimSpace(content)
	if body == "" {
		return []string{"revised_md"}
	}

	var missing []string
	if names := draft.FindCitationMarkers(body); len(names) > 0 {
		missing = append(missing, fmt.Sprintf("revised_md (unresolved citation marker(s) remain: %v -- every [^shortname] must be resolved into a real [anchor](url) link before this pass is done)", names))
	}

	codeFree := draft.StripCode(body)
	for _, m := range draft.MarkdownLinkRe.FindAllStringSubmatch(codeFree, -1) {
		url := m[1]
		if !validSourceURLs[url] {
			missing = append(missing, fmt.Sprintf("revised_md (link to %q does not match any source this piece's research recorded -- never invent a link, resolve only from sources.md)", url))
		}
	}

	bareURLBody := draft.MarkdownLinkRe.ReplaceAllString(codeFree, "")
	if draft.HasRawURL(bareURLBody) {
		missing = append(missing, "revised_md (raw URL found outside a [anchor](url) link -- every URL must be a resolved markdown link, never bare text)")
	}

	return missing
}
