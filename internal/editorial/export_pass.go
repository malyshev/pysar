package editorial

// exportPass is the seventh real per-piece editorial pass (after intake,
// research, draft, staff-edit, sharpen, humanize). File I/O lives in
// internal/export + MCP; this Pass only owns precondition + produced-
// artifact state updates.
//
// Precondition requires ArtifactDraft only -- staff-edit, sharpen, and
// humanize are all optional (same reasoning throughout this pipeline), so
// /ps's autopilot orchestrator can export whatever stage a piece actually
// reached, not just a fully-piped one.
type exportPass struct{}

var _ Pass = (*exportPass)(nil)

func (exportPass) Name() string { return "export" }

func (exportPass) Precondition(s *State) error {
	if !s.Has(ArtifactDraft) {
		return &PreconditionError{
			Missing: []ArtifactKind{ArtifactDraft},
			Hint:    "no draft.md found for this piece -- either /ps-draft hasn't run yet, or piece_path doesn't point at a real piece (check for a typo)",
		}
	}
	return nil
}

func (exportPass) Body(s *State) (*State, error) {
	return s.WithProduced(ArtifactExported), nil
}
