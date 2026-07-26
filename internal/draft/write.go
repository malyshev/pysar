package draft

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"pysar/internal/intake"
)

// WriteToPiece writes draft.md to dir and appends a draft-changelog.md
// entry, returning the draft's word count so callers (the run-log entry,
// the hand-off summary) reuse this one computation instead of each
// re-scanning DraftMD themselves. Unlike research's sources.md, a redraft
// is expected to replace the previous draft.md wholesale -- there is no
// accumulating list to merge, so overwriting here is correct, not a repeat
// of research's overwrite bug.
func WriteToPiece(dir string, b Bundle) (words int, err error) {
	words = WordCount(b.DraftMD)

	body := strings.TrimSpace(b.DraftMD) + "\n"
	if err := os.WriteFile(filepath.Join(dir, "draft.md"), []byte(body), 0o644); err != nil {
		return 0, fmt.Errorf("write draft.md: %w", err)
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
