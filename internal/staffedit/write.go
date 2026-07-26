package staffedit

import (
	"pysar/internal/draft"
)

// WriteToPiece writes the revised content to staff-edit.md -- never to
// draft.md, which stays as the untouched first-draft artifact so the two
// can be compared. Wholesale-replaces staff-edit.md on each call (same
// contract /ps-draft's own WriteToPiece uses for draft.md -- there's one
// current staff-edit revision, not an accumulating list) and appends a
// staff-edit-changelog.md entry recording what changed.
func WriteToPiece(dir string, b Bundle) (words int, err error) {
	return draft.WriteRevision(dir, "staff-edit.md", "staff-edit-changelog.md", b.RevisedMD, b.Checks, b.Mode)
}
