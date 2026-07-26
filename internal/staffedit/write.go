package staffedit

import (
	"strings"

	"pysar/internal/draft"
	"pysar/internal/intake"
)

// WriteToPiece writes the revised content to staff-edit.md -- never to
// draft.md, which stays as the untouched first-draft artifact so the two
// can be compared. Wholesale-replaces staff-edit.md on each call (same
// contract /ps-draft's own WriteToPiece uses for draft.md -- there's one
// current staff-edit revision, not an accumulating list) and appends a
// staff-edit-changelog.md entry recording what changed. Returns the
// revised content's word count for the caller to reuse (run-log entry,
// hand-off summary) instead of re-scanning it.
func WriteToPiece(dir string, b Bundle) (words int, err error) {
	words = draft.WordCount(b.RevisedMD)

	if err := draft.WriteArticleFile(dir, "staff-edit.md", b.RevisedMD); err != nil {
		return 0, err
	}

	summary := strings.Join(b.Checks, "; ")
	if mode := strings.TrimSpace(b.Mode); mode != "" {
		summary = "(" + mode + ") " + summary
	}
	line := intake.FormatChangelogLine(summary)
	if err := intake.AppendChangelogLine(dir, "staff-edit-changelog.md", line); err != nil {
		return 0, err
	}
	return words, nil
}
