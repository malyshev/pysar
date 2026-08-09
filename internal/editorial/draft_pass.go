package editorial

import "pysar/internal/stagereq"

// draftPass is the third real per-piece editorial pass (after intake,
// research) -- the real implementation of the "draft-edit" seam stub named
// in dec-20260718-ffbfb04b's own Impact Measurement, not a new pass added
// alongside it. File I/O lives in internal/draft + MCP; this Pass only owns
// precondition + produced-artifact state updates.
//
// Precondition only requires ArtifactBrief, not ArtifactStake separately --
// intakePass always produces both together in one Body call (intake_pass.go),
// so Brief's presence already implies Stake's. Research is optional by
// default; when the piece lists research in required_stages
// (dec-20260809-701b59d3), ArtifactSourcesFull (research_mode: full) is
// also required before draft may proceed.
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
	if s.Requires(stagereq.StageResearch) && !s.Has(ArtifactSourcesFull) {
		return &PreconditionError{
			Missing: []ArtifactKind{ArtifactSourcesFull},
			Hint:    "this piece requires full research before draft (required_stages includes research / /ps --research) -- run /ps-research and save_research_bundle so research_mode is full",
		}
	}
	return nil
}

func (draftPass) Body(s *State) (*State, error) {
	return s.WithProduced(ArtifactDraft), nil
}
