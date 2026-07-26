package sharpen

import (
	"pysar/internal/draft"
)

// WriteToPiece writes the revised content to sharpen.md -- never to
// draft.md or staff-edit.md, which stay as their own untouched stages.
// Wholesale-replaces sharpen.md on each call (same contract the earlier
// passes use -- one current sharpen revision, not an accumulating list)
// and appends a sharpen-changelog.md entry recording what changed.
func WriteToPiece(dir string, b Bundle) (words int, err error) {
	return draft.WriteRevision(dir, "sharpen.md", "sharpen-changelog.md", b.RevisedMD, b.Checks, b.Mode)
}
