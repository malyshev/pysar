package editorial

// staffEditPass is the fourth real per-piece editorial pass (after intake,
// research, draft). File I/O lives in internal/staffedit + MCP; this Pass
// only owns precondition + produced-artifact state updates.
//
// Precondition requires ArtifactDraft, not just ArtifactBrief -- staff-edit
// revises an existing draft.md, it doesn't start one.
type staffEditPass struct{}

var _ Pass = (*staffEditPass)(nil)

func (staffEditPass) Name() string { return "staff-edit" }

func (staffEditPass) Precondition(s *State) error {
	if !s.Has(ArtifactDraft) {
		return &PreconditionError{
			Missing: []ArtifactKind{ArtifactDraft},
			Hint:    "no draft.md found for this piece -- either /ps-draft hasn't run yet, or piece_path doesn't point at a real piece (check for a typo)",
		}
	}
	return nil
}

func (staffEditPass) Body(s *State) (*State, error) {
	return s.WithProduced(ArtifactStaffEdit), nil
}
