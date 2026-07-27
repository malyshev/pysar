package humanize

import (
	"pysar/internal/draft"
)

// WriteToPiece writes the revised content to humanize.md -- never to
// draft.md, staff-edit.md, or sharpen.md, which stay as their own
// untouched stages. Wholesale-replaces humanize.md on each call (same
// contract the earlier passes use -- one current humanize revision, not
// an accumulating list) and appends a humanize-changelog.md entry
// recording what changed.
func WriteToPiece(dir string, b Bundle) (words int, err error) {
	return draft.WriteRevision(dir, "humanize.md", "humanize-changelog.md", b.RevisedMD, b.Checks, b.Mode)
}
