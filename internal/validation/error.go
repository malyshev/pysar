// Package validation holds the "collect every missing/invalid field, then
// report them all at once" error shape every bundle validator in the
// project needs (intake, research) -- one type instead of a hand-copied
// struct+Error() method per package.
package validation

import (
	"fmt"
	"strings"
)

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

// RequireNonEmptyChecks validates the "checks" field shared by internal/
// staffedit and internal/sharpen bundles: at least one entry required, and
// no entry may be blank. Returns missing-field messages for the caller to
// append to its own accumulated list, the same pattern draft.ValidateContent
// uses -- callers merge this alongside their other checks rather than
// treating it as a standalone error.
func RequireNonEmptyChecks(checks []string) []string {
	var missing []string
	if len(checks) == 0 {
		missing = append(missing, "checks (>=1 required -- even 'no changes needed' is worth recording as one entry)")
	}
	for i, c := range checks {
		if strings.TrimSpace(c) == "" {
			missing = append(missing, fmt.Sprintf("checks[%d] (empty)", i))
		}
	}
	return missing
}
