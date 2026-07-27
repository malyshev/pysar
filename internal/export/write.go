package export

import (
	"fmt"
	"os"
	"path/filepath"

	"pysar/internal/draft"
)

// WriteToRoot copies the piece's most-refined revision (humanize.md if it
// ran, else sharpen.md, else staff-edit.md, else draft.md -- same priority
// draft.LatestRevisionFile already uses for sharpen/humanize's own
// revised_from tracking) to <projectRoot>/<slug>.md, overwriting any
// previous export of the same piece. Returns the destination path, which
// source file was copied, and its word count.
func WriteToRoot(projectRoot, pieceDir string) (destPath, sourceFile string, words int, err error) {
	sourceFile = draft.LatestRevisionFile(pieceDir, "humanize.md", "sharpen.md", "staff-edit.md", "draft.md")

	content, err := os.ReadFile(filepath.Join(pieceDir, sourceFile))
	if err != nil {
		return "", "", 0, fmt.Errorf("read %s: %w", sourceFile, err)
	}

	slug := filepath.Base(pieceDir)
	destPath = filepath.Join(projectRoot, slug+".md")
	if err := os.WriteFile(destPath, content, 0o644); err != nil {
		return "", "", 0, fmt.Errorf("write %s: %w", destPath, err)
	}

	words = draft.WordCount(string(content))
	return destPath, sourceFile, words, nil
}
