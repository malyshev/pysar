package sharpen

import (
	"strings"

	"pysar/internal/draft"
	"pysar/internal/intake"
)

// WriteToPiece writes the revised content to sharpen.md -- never to
// draft.md or staff-edit.md, which stay as their own untouched stages.
// Wholesale-replaces sharpen.md on each call (same contract the earlier
// passes use -- one current sharpen revision, not an accumulating list)
// and appends a sharpen-changelog.md entry recording what changed. Returns
// the revised content's word count for the caller to reuse (run-log entry,
// hand-off summary) instead of re-scanning it.
func WriteToPiece(dir string, b Bundle) (words int, err error) {
	words = draft.WordCount(b.RevisedMD)

	if err := draft.WriteArticleFile(dir, "sharpen.md", b.RevisedMD); err != nil {
		return 0, err
	}

	summary := strings.Join(b.Checks, "; ")
	if mode := strings.TrimSpace(b.Mode); mode != "" {
		summary = "(" + mode + ") " + summary
	}
	line := intake.FormatChangelogLine(summary)
	if err := intake.AppendChangelogLine(dir, "sharpen-changelog.md", line); err != nil {
		return 0, err
	}
	return words, nil
}
