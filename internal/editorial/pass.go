package editorial

import "fmt"

// PreconditionError is returned by a Pass's Precondition when it is not
// satisfied. It is always structured and actionable -- a Pass never fails
// silently or without explanation (dec-20260718-20a08f83 invariant).
type PreconditionError struct {
	Missing []ArtifactKind
	Hint    string
}

func (e *PreconditionError) Error() string {
	if len(e.Missing) == 0 {
		return fmt.Sprintf("precondition not met: %s", e.Hint)
	}
	return fmt.Sprintf("precondition not met: missing %v -- %s", e.Missing, e.Hint)
}

// Pass is the unit every author-facing editorial action implements. Each Pass
// owns its own independent Precondition check; no central table or object
// owns "the one correct sequence" for all passes (dec-20260718-20a08f83).
// A type missing Precondition cannot satisfy Pass and so cannot be Registered
// -- the compiler enforces this, not convention (dec-20260718-ffbfb04b).
type Pass interface {
	Name() string
	Precondition(s *State) error
	Body(s *State) (*State, error)
}

// allowsEmpty is an optional interface a Pass implements when it is
// legitimately applicable to a piece with no produced artifacts yet (e.g. the
// first pass in a piece's life). Passes not implementing it are assumed to
// require something already produced.
type allowsEmpty interface {
	AllowEmpty() bool
}

// OnboardingPass is an optional interface a Pass implements to mark itself as
// author-level onboarding (voice, style, and any future kind) rather than a
// per-piece editorial pass (idea-scaffold, draft-edit, discoverability). It
// is the property check_onboarding_status filters the shared registry on --
// never a hardcoded pass name (dec-20260718-ebca6318 invariant: "/ps-onboard
// contains no hardcoded knowledge of specific questionnaire names").
type OnboardingPass interface {
	IsOnboardingPass() bool
}

// Run gates a Pass's Body behind its own Precondition. The body never
// executes when Check fails (dec-20260718-20a08f83).
func Run(p Pass, s *State) (*State, error) {
	if err := p.Precondition(s); err != nil {
		return nil, err
	}
	return p.Body(s)
}

var registry []Pass

// Register adds p to the package registry. Passing a value that does not
// satisfy Pass (e.g. missing Precondition) fails to compile at the call site.
func Register(p Pass) {
	registry = append(registry, p)
}

// Registered returns a copy of the currently registered passes.
func Registered() []Pass {
	out := make([]Pass, len(registry))
	copy(out, registry)
	return out
}
