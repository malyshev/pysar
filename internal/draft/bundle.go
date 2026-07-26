// Package draft validates and persists /ps-draft output: the first full
// draft.md for an existing piece. Content craft (voice, structure, citation
// hygiene, consumability, worked-example rule) is adapted from a prior
// writing POC's article-writer discipline (CL1) -- author-registry
// resolution and Medium-specific publish ceremony (SEO fields, tags, cover
// image, Boost optimization, cross-article variance scanning) are
// deliberately absent: this project has no author registry (voice.md/
// style.md via /ps-onboard is the whole mechanism) and targets more than
// one surface, so it doesn't own Medium's UI checklist.
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
	citationRe   = regexp.MustCompile(`\[\^([a-zA-Z0-9-]+)\]`)
	rawURLRe     = regexp.MustCompile(`https?://|www\.`)
	codeFenceRe  = regexp.MustCompile("(?s)```.*?```")
	inlineCodeRe = regexp.MustCompile("`[^`\n]*`")
)

// stripCode removes both fenced (```...```) and inline (`...`) code from
// body, once, so the raw-URL check and the citation check apply the same
// code exclusion consistently -- a code fence that shows citation syntax
// as a literal example (e.g. "cite claims like `[^example]`") must not be
// checked against real shortnames, the same way it already shields a URL
// shown as an example from the raw-URL check.
func stripCode(body string) string {
	return inlineCodeRe.ReplaceAllString(codeFenceRe.ReplaceAllString(body, ""), "")
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
// just be Medium's own structural constraint re-smuggled in as "universal"
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
	codeFree := stripCode(body)
	if rawURLRe.MatchString(codeFree) {
		missing = append(missing, fmt.Sprintf("%s (raw URL found in prose -- cite via [^shortname] instead)", fieldName))
	}
	for _, m := range citationRe.FindAllStringSubmatch(codeFree, -1) {
		name := m[1]
		if !validShortnames[name] {
			missing = append(missing, fmt.Sprintf("%s ([^%s] does not match any source this piece's research has recorded)", fieldName, name))
		}
	}
	return missing
}

// WordCount is a plain whitespace-delimited word count of markdown source --
// good enough for a changelog line and a hand-off summary, not a claim of
// precise prose-word counting (code fences, headings, and citation markers
// all count as their literal tokens).
func WordCount(draftMD string) int {
	return len(strings.Fields(draftMD))
}
