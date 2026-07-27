package editorial

// Seam stub for a later pass; intake, research, draft, staff-edit,
// sharpen, humanize, and export are real passes in their own files
// (dec-20260718-ffbfb04b, dec-20260725-35fa2d24).

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
			Hint:    "run draft first to produce a draft",
		}
	}
	return nil
}

func (discoverabilityPass) Body(s *State) (*State, error) {
	return s, nil
}

func init() {
	Register(intakePass{})
	Register(researchPass{})
	Register(draftPass{})
	Register(staffEditPass{})
	Register(sharpenPass{})
	Register(humanizePass{})
	Register(exportPass{})
	Register(discoverabilityPass{})
}
