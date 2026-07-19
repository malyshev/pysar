package onboarding

import "pysar/internal/editorial"

// stylePass registers style capture in the same Pass registry as editorial
// passes and voicePass (dec-20260718-ebca6318 invariant: "never a bespoke
// onboarding-only orchestration mechanism"). Its Precondition mirrors
// voicePass exactly, substituting ArtifactStyleProfile.
//
// Body is a schema-layer stub, same as voicePass -- the actual conversation
// happens on the agentic surface (/ps-style), not in Go (dec-20260718-5570376d).
type stylePass struct{}

var _ editorial.Pass = (*stylePass)(nil)

func (stylePass) Name() string { return "ps-style" }

// AllowEmpty: style capture is an onboarding action, meaningful even on a
// piece/project with nothing else produced yet.
func (stylePass) AllowEmpty() bool { return true }

// IsOnboardingPass marks stylePass as author-level onboarding, not a
// per-piece editorial pass -- see editorial.OnboardingPass.
func (stylePass) IsOnboardingPass() bool { return true }

func (stylePass) Precondition(s *editorial.State) error {
	if s.Has(editorial.ArtifactStyleProfile) {
		return &editorial.PreconditionError{Hint: "style profile already captured"}
	}
	return nil
}

func (stylePass) Body(s *editorial.State) (*editorial.State, error) {
	return s.WithProduced(editorial.ArtifactStyleProfile), nil
}

func init() {
	editorial.Register(stylePass{})
}
