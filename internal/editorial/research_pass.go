package editorial

// researchPass is the second real per-piece editorial pass (after intake),
// wired into dec-20260718-20a08f83's Precondition mechanism the same way
// intakePass is. File I/O lives in internal/research + MCP; this Pass only
// owns precondition + produced-artifact state updates.
//
// This governs the piece-anchored /ps-research flow only. Standalone
// research (no existing piece) has no piece identity to attach a
// Precondition/run-log to and is deliberately not routed through here --
// see internal/research.WriteStandalone.
type researchPass struct{}

var _ Pass = (*researchPass)(nil)

func (researchPass) Name() string { return "research" }

func (researchPass) Precondition(s *State) error {
	if !s.Has(ArtifactBrief) {
		return &PreconditionError{
			Missing: []ArtifactKind{ArtifactBrief},
			Hint:    "no brief.md found for this piece -- either intake hasn't produced one yet, or piece_path doesn't point at a real piece (check for a typo)",
		}
	}
	return nil
}

func (researchPass) Body(s *State) (*State, error) {
	return s.WithProduced(ArtifactSourcesFull), nil
}
