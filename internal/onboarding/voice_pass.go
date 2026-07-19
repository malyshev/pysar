package onboarding

import "pysar/internal/editorial"

// voicePass registers voice capture in the same Pass registry as editorial
// passes (dec-20260718-ebca6318 invariant: "never a bespoke onboarding-only
// orchestration mechanism"). Its Precondition expresses "still needed": it
// succeeds (may run) until a voice profile has been produced, then declines.
//
// Body is a schema-layer stub. The actual conversation that fills a
// VoiceProfile happens on the agentic surface, not in Go
// (dec-20260718-5570376d); wiring a real, conversationally-filled profile's
// persistence into .pysar/project is the /ps-voice skill's job, not this
// package's.
type voicePass struct{}

var _ editorial.Pass = (*voicePass)(nil)

func (voicePass) Name() string { return "ps-voice" }

// AllowEmpty: voice capture is an onboarding action, meaningful even on a
// piece/project with nothing else produced yet.
func (voicePass) AllowEmpty() bool { return true }

// IsOnboardingPass marks voicePass as author-level onboarding, not a
// per-piece editorial pass -- see editorial.OnboardingPass.
func (voicePass) IsOnboardingPass() bool { return true }

func (voicePass) Precondition(s *editorial.State) error {
	if s.Has(editorial.ArtifactVoiceProfile) {
		return &editorial.PreconditionError{Hint: "voice profile already captured"}
	}
	return nil
}

func (voicePass) Body(s *editorial.State) (*editorial.State, error) {
	return s.WithProduced(editorial.ArtifactVoiceProfile), nil
}

func init() {
	editorial.Register(voicePass{})
}
