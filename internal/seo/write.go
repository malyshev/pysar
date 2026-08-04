package seo

import (
	"fmt"
	"strings"

	"pysar/internal/draft"
)

// WriteToPiece writes the revised content to seo.md -- never to draft.md,
// staff-edit.md, or sharpen.md, which stay as their own untouched stages --
// appends a seo-changelog.md entry, and (re)writes seo-checklist.md, the
// piece's discoverability fields. seo-checklist.md is a separate file from
// seo.md on purpose: seo.md is reader-facing article prose (what export
// copies to the project root once humanize has run against it);
// seo-checklist.md is operator-facing packaging metadata (tags, meta
// title/description, URL slug) the author copies into whatever platform
// they're posting to -- the same distinction the reference implementation
// this is adapted from draws between article.md and publish.md, kept here
// as two files rather than reviving a combined Status-enum publish
// checklist file.
func WriteToPiece(dir string, b Bundle) (words int, err error) {
	words, err = draft.WriteRevision(dir, "seo.md", "seo-changelog.md", b.RevisedMD, b.Checks, b.Mode)
	if err != nil {
		return 0, err
	}
	if err := draft.WriteArticleFile(dir, "seo-checklist.md", renderChecklist(b)); err != nil {
		return 0, err
	}
	return words, nil
}

func renderChecklist(b Bundle) string {
	var out strings.Builder
	out.WriteString("# Discoverability checklist\n\n")
	fmt.Fprintf(&out, "- Meta title: %s\n", strings.TrimSpace(b.MetaTitle))
	fmt.Fprintf(&out, "- Meta description: %s\n", strings.TrimSpace(b.MetaDescription))
	fmt.Fprintf(&out, "- URL slug: %s\n", strings.TrimSpace(b.URLSlug))
	out.WriteString("- Tags:\n")
	for _, t := range b.Tags {
		fmt.Fprintf(&out, "  - %s\n", strings.TrimSpace(t))
	}
	return out.String()
}
