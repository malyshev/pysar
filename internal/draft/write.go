package draft

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"pysar/internal/intake"
)

// WriteArticleFile writes content to dir/filename wholesale. Exported so
// other passes that also persist a full revision of a piece's prose (e.g.
// internal/staffedit's staff-edit.md, a separate file from draft.md so the
// original draft survives for comparison) reuse the exact same write
// mechanics instead of duplicating it -- each pass still owns its own
// changelog entry and filename, since those genuinely differ per pass.
func WriteArticleFile(dir, filename, content string) error {
	body := strings.TrimSpace(content) + "\n"
	if err := os.WriteFile(filepath.Join(dir, filename), []byte(body), 0o644); err != nil {
		return fmt.Errorf("write %s: %w", filename, err)
	}
	return nil
}

// WriteToPiece writes draft.md to dir and appends a draft-changelog.md
// entry, returning the draft's word count so callers (the run-log entry,
// the hand-off summary) reuse this one computation instead of each
// re-scanning DraftMD themselves. Unlike research's sources.md, a redraft
// is expected to replace the previous draft.md wholesale -- there is no
// accumulating list to merge, so overwriting here is correct, not a repeat
// of research's overwrite bug.
func WriteToPiece(dir string, b Bundle) (words int, err error) {
	words = WordCount(b.DraftMD)

	if err := WriteArticleFile(dir, "draft.md", b.DraftMD); err != nil {
		return 0, err
	}

	note := strings.TrimSpace(b.Notes)
	if note == "" {
		note = "draft"
	}
	line := intake.FormatChangelogLine(fmt.Sprintf("%s (%d words)", note, words))
	if err := intake.AppendChangelogLine(dir, "draft-changelog.md", line); err != nil {
		return 0, err
	}
	return words, nil
}
