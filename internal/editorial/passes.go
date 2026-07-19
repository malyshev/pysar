package editorial

// Seam stubs exercising the mechanism (ordering + target-surface
// applicability) ahead of real editorial pass bodies (dec-20260718-ffbfb04b).

type ideaScaffoldPass struct{}

var _ Pass = (*ideaScaffoldPass)(nil)

func (ideaScaffoldPass) Name() string { return "idea-scaffold" }

// AllowEmpty: idea-scaffold is the legitimate first pass -- it applies to a
// piece with nothing produced yet.
func (ideaScaffoldPass) AllowEmpty() bool { return true }

func (ideaScaffoldPass) Precondition(_ *State) error { return nil }

func (ideaScaffoldPass) Body(s *State) (*State, error) {
	return s.WithProduced(ArtifactStake), nil
}

type draftEditPass struct{}

var _ Pass = (*draftEditPass)(nil)

func (draftEditPass) Name() string { return "draft-edit" }

func (draftEditPass) Precondition(s *State) error {
	if !s.Has(ArtifactStake) {
		return &PreconditionError{
			Missing: []ArtifactKind{ArtifactStake},
			Hint:    "run idea-scaffold first to produce a stake",
		}
	}
	return nil
}

func (draftEditPass) Body(s *State) (*State, error) {
	return s.WithProduced(ArtifactDraft), nil
}

type discoverabilityPass struct{}

var _ Pass = (*discoverabilityPass)(nil)

func (discoverabilityPass) Name() string { return "discoverability" }

func (discoverabilityPass) Precondition(s *State) error {
	if s.Surface() == SurfaceLetter {
		return &PreconditionError{
			Hint: "discoverability does not apply to a paper-letter target surface",
		}
	}
	if !s.Has(ArtifactDraft) {
		return &PreconditionError{
			Missing: []ArtifactKind{ArtifactDraft},
			Hint:    "run draft-edit first to produce a draft",
		}
	}
	return nil
}

func (discoverabilityPass) Body(s *State) (*State, error) {
	return s, nil
}

func init() {
	Register(ideaScaffoldPass{})
	Register(draftEditPass{})
	Register(discoverabilityPass{})
}
