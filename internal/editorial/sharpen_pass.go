package editorial

// sharpenPass is the fifth real per-piece editorial pass (after intake,
// research, draft, staff-edit). File I/O lives in internal/sharpen + MCP;
// this Pass only owns precondition + produced-artifact state updates.
//
// Precondition requires ArtifactDraft, not ArtifactStaffEdit -- staff-edit
// is optional (same reasoning as research being optional for draft), so
// sharpen must work whether or not it ran, reading whichever is the
// piece's latest revision (staff-edit.md if present, else draft.md).
type sharpenPass struct{}

var _ Pass = (*sharpenPass)(nil)

func (sharpenPass) Name() string { return "sharpen" }

func (sharpenPass) Precondition(s *State) error {
	if !s.Has(ArtifactDraft) {
		return &PreconditionError{
			Missing: []ArtifactKind{ArtifactDraft},
			Hint:    "no draft.md found for this piece -- either /ps-draft hasn't run yet, or piece_path doesn't point at a real piece (check for a typo)",
		}
	}
	return nil
}

func (sharpenPass) Body(s *State) (*State, error) {
	return s.WithProduced(ArtifactSharpen), nil
}
