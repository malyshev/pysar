// Package validation holds the "collect every missing/invalid field, then
// report them all at once" error shape every bundle validator in the
// project needs (intake, research) -- one type instead of a hand-copied
// struct+Error() method per package.
package validation

import "fmt"

// Error names every missing/invalid field a bundle validator found, all at
// once rather than failing on the first one -- so the agent assembling a
// bundle sees every problem in one round-trip instead of playing
// whack-a-mole against save_*_bundle.
type Error struct {
	Kind    string // e.g. "intake bundle", "research bundle"
	Missing []string
}

func (e *Error) Error() string {
	return fmt.Sprintf("%s incomplete: missing %v", e.Kind, e.Missing)
}
