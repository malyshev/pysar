package editorial

// humanizePass is the sixth real per-piece editorial pass (after intake,
// research, draft, staff-edit, sharpen). File I/O lives in internal/
// humanize + MCP; this Pass only owns precondition + produced-artifact
// state updates.
//
// Precondition requires ArtifactDraft, not ArtifactSharpen or
// ArtifactStaffEdit -- both are optional (same reasoning as research being
// optional for draft), so humanize must work whether or not either ran,
// reading whichever is the piece's latest revision (sharpen.md if
// present, else staff-edit.md, else draft.md).
type humanizePass struct{}

var _ Pass = (*humanizePass)(nil)

func (humanizePass) Name() string { return "humanize" }

func (humanizePass) Precondition(s *State) error {
	if !s.Has(ArtifactDraft) {
		return &PreconditionError{
			Missing: []ArtifactKind{ArtifactDraft},
			Hint:    "no draft.md found for this piece -- either /ps-draft hasn't run yet, or piece_path doesn't point at a real piece (check for a typo)",
		}
	}
	return nil
}

func (humanizePass) Body(s *State) (*State, error) {
	return s.WithProduced(ArtifactHumanize), nil
}
