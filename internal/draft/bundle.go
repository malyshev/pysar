// Package draft validates and persists /ps-draft output: the first full
// draft.md for an existing piece. Content craft covers voice, structure,
// citation hygiene, consumability, and the worked-example rule.
// Author-registry resolution and platform-specific publish ceremony (SEO
// fields, tags, cover image, cross-article variance scanning) are
// deliberately absent: this project has no author registry (voice.md/
// style.md via /ps-onboard is the whole mechanism) and targets more than
// one surface, so it doesn't own a platform publish checklist.
//
// draft.md is channel-agnostic writing, not a finished channel artifact --
// its [^shortname] markers stay unresolved here on purpose. Resolving them
// into inline links (or any other channel-specific packaging: SEO
// metadata, a formatted letter, whatever a given surface needs) is a later,
// optional, separate concern this package doesn't own.
package draft

import (
	"fmt"
	"regexp"
	"strings"

	"pysar/internal/validation"
)

// Bundle is the structured draft payload the agent assembles and the MCP
// tool validates before writing.
type Bundle struct {
	// PiecePath is the existing piece this draft belongs to, relative to
	// the project root -- resolved via research.ResolvePieceDir, so a
	// directory or any file inside it (brief.md, a raw/ excerpt) works.
	PiecePath string `json:"piece_path"`

	// DraftMD is the full draft: channel-agnostic prose. Citations go
	// through [^shortname] markers matching a real source recorded by
	// /ps-research -- never a raw URL, and never resolved to an inline
	// link here -- that's a later, channel-specific packaging concern.
	DraftMD string `json:"draft_md"`

	// Notes is an optional one-line changelog note (e.g. "first full
	// draft", "revision after staff-edit pass"). Defaults to "draft".
	Notes string `json:"notes,omitempty"`
}

var (
	citationRe  = regexp.MustCompile(`\[\^([a-zA-Z0-9-]+)\]`)
	rawURLRe    = regexp.MustCompile(`https?://|www\.`)
	codeFenceRe = regexp.MustCompile("(?s)```.*?```")
	inlineRe    = regexp.MustCompile("`[^`\n]*`")

	// MarkdownLinkRe matches a resolved [anchor](url) link, capturing the
	// URL. Allows one level of nested parens inside the URL (e.g.
	// https://en.wikipedia.org/wiki/Foo_(bar), a common shape on
	// Wikipedia/doc sites) -- a naive `\(([^)]+)\)` stops at the first
	// `)`, truncating the captured URL and making a correctly-resolved
	// link fail every validSourceURLs lookup that compares against the
	// full recorded URL. Exported so internal/humanize and internal/seo,
	// which both need to recognize an already-resolved link, share this
	// one definition instead of each re-declaring a copy that can drift
	// (and did: two independent copies existed with the same truncation
	// bug before this was extracted).
	MarkdownLinkRe = regexp.MustCompile(`\[[^\]]+\]\(([^()]*(?:\([^()]*\)[^()]*)*)\)`)
)

// StripCode removes both fenced (```...```) and inline (`...`) code from
// body, once, so callers apply the same code exclusion consistently -- a
// code fence that shows citation or link syntax as a literal example
// (e.g. "cite claims like `[^example]`") must not be checked against real
// shortnames/sources, the same way it already shields a URL shown as an
// example from the raw-URL check. Exported so internal/seo (which needs
// the same exclusion for its own inverted citation rule) doesn't
// re-implement it.
func StripCode(body string) string {
	return inlineRe.ReplaceAllString(codeFenceRe.ReplaceAllString(body, ""), "")
}

// Validate checks citation integrity for a durable draft write: every
// [^shortname] marker resolves to a real source, and nothing in the prose
// bypasses that convention with a raw URL instead. validShortnames is the
// piece's current research output (research.LoadShortnames) -- empty when
// no research has run, which correctly makes any citation marker invalid
// rather than silently accepted.
//
// This deliberately does not check markdown heading structure (title-line
// count, heading depth) or any other rendering-shape constraint -- those
// are properties of a *blog-shaped* target surface, and this project also
// targets SurfaceLetter (internal/editorial/state.go), which has no
// heading hierarchy at all. A mechanical "exactly one H1" gate here would
// just be a blog-surface structural constraint re-smuggled in as "universal"
// -- the same class of surface-specific ceremony this package's own doc
// comment says is deliberately excluded. Heading structure for blog-shaped
// pieces is SKILL.md-level guidance instead, same tier as the
// worked-example rule below.
//
// This also does not check voice/structure/craft quality (worked examples,
// thesis-consistency, killer-section coverage, consumability) -- those are
// semantic judgment calls the SKILL.md's own self-check step applies; a
// naive mechanical check here would produce false positives that erode
// trust in the tool, the same reasoning intake.Validate leaves
// killer_sections'/counterintuitive's depth to agent judgment beyond
// presence.
func Validate(b Bundle, validShortnames map[string]bool) error {
	var missing []string

	if strings.TrimSpace(b.PiecePath) == "" {
		missing = append(missing, "piece_path")
	}
	missing = append(missing, ValidateContent("draft_md", b.DraftMD, validShortnames)...)

	if len(missing) > 0 {
		return &validation.Error{Kind: "draft bundle", Missing: missing}
	}
	return nil
}

// ValidateContent checks citation integrity for a piece of draft prose --
// every [^shortname] marker resolves to a real source, nothing bypasses
// that convention with a raw URL, and the content isn't empty. fieldName
// names the field in the CALLER's own schema -- "draft_md" here, but
// internal/staffedit and internal/sharpen call this directly with their
// own field name ("revised_md") for content that's mechanically the same
// draft prose under a different name, so their error messages point at an
// argument that actually exists in the tool being validated for, instead
// of a hardcoded "draft_md" that isn't in their schema at all.
func ValidateContent(fieldName, content string, validShortnames map[string]bool) []string {
	body := strings.TrimSpace(content)
	if body == "" {
		return []string{fieldName}
	}

	var missing []string
	if HasRawURL(body) {
		missing = append(missing, fmt.Sprintf("%s (raw URL found in prose -- cite via [^shortname] instead)", fieldName))
	}
	missing = append(missing, ValidateCitations(fieldName, body, validShortnames)...)
	return missing
}

// ValidateCitations checks every [^shortname] marker in content resolves
// against validShortnames. Split out from ValidateContent so a caller
// with a different raw-URL policy (internal/humanize, once /ps-seo has
// resolved some markers into real links) can reuse just the citation half
// instead of re-implementing it.
func ValidateCitations(fieldName, content string, validShortnames map[string]bool) []string {
	var missing []string
	for _, name := range FindCitationMarkers(content) {
		if !validShortnames[name] {
			missing = append(missing, fmt.Sprintf("%s ([^%s] does not match any source this piece's research has recorded)", fieldName, name))
		}
	}
	return missing
}

// FindCitationMarkers returns every [^shortname] marker's shortname found
// in content, code-fence-stripped so a marker shown as a literal example
// doesn't get treated as a real citation. Exported so a caller with a
// different acceptance rule than ValidateCitations' "must resolve" (
// internal/seo's inverted "must NOT remain" rule) can reuse the same scan
// instead of re-declaring the citation regex.
func FindCitationMarkers(content string) []string {
	matches := citationRe.FindAllStringSubmatch(StripCode(content), -1)
	if len(matches) == 0 {
		return nil
	}
	names := make([]string, len(matches))
	for i, m := range matches {
		names[i] = m[1]
	}
	return names
}

// HasRawURL reports whether content contains a URL outside of code fences/
// inline code -- exported so a caller with a different tolerance for
// *where* a URL is allowed (internal/humanize allows one inside a resolved
// [anchor](url) link; this package's own ValidateContent allows none) can
// reuse the same detection instead of re-implementing the regex.
func HasRawURL(content string) bool {
	return rawURLRe.MatchString(StripCode(content))
}

// WordCount is a plain whitespace-delimited word count of markdown source --
// good enough for a changelog line and a hand-off summary, not a claim of
// precise prose-word counting (code fences, headings, and citation markers
// all count as their literal tokens).
func WordCount(draftMD string) int {
	return len(strings.Fields(draftMD))
}
