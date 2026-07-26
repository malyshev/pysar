package editorial

// draftPass is the third real per-piece editorial pass (after intake,
// research) -- the real implementation of the "draft-edit" seam stub named
// in dec-20260718-ffbfb04b's own Impact Measurement, not a new pass added
// alongside it. File I/O lives in internal/draft + MCP; this Pass only owns
// precondition + produced-artifact state updates.
//
// Precondition only requires ArtifactBrief, not ArtifactStake separately --
// intakePass always produces both together in one Body call (intake_pass.go),
// so Brief's presence already implies Stake's. Research is deliberately not
// required: it's an optional pass (researchPass's own doc comment), so draft
// must work whether or not it ran.
type draftPass struct{}

var _ Pass = (*draftPass)(nil)

func (draftPass) Name() string { return "draft" }

func (draftPass) Precondition(s *State) error {
	if !s.Has(ArtifactBrief) {
		return &PreconditionError{
			Missing: []ArtifactKind{ArtifactBrief},
			Hint:    "no brief.md found for this piece -- either intake hasn't produced one yet, or piece_path doesn't point at a real piece (check for a typo)",
		}
	}
	return nil
}

func (draftPass) Body(s *State) (*State, error) {
	return s.WithProduced(ArtifactDraft), nil
}
