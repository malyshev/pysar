package editorial

// Seam stubs for later passes; intake is a real pass in intake_pass.go
// (dec-20260718-ffbfb04b, dec-20260725-35fa2d24).

type draftEditPass struct{}

var _ Pass = (*draftEditPass)(nil)

func (draftEditPass) Name() string { return "draft-edit" }

func (draftEditPass) Precondition(s *State) error {
	if !s.Has(ArtifactStake) && !s.Has(ArtifactBrief) {
		return &PreconditionError{
			Missing: []ArtifactKind{ArtifactStake, ArtifactBrief},
			Hint:    "run intake first to produce a stake/brief",
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
	Register(intakePass{})
	Register(draftEditPass{})
	Register(discoverabilityPass{})
}
