package editorial

// intakePass is the first real per-piece editorial pass
// (dec-20260725-35fa2d24). File I/O lives in internal/intake + MCP; this Pass
// only owns precondition + produced-artifact state updates.
//
// Precondition is unconditional: every intake call allocates a fresh,
// uniquely-suffixed piece directory (internal/intake.AllocateUniqueName), so
// there is no existing-piece state to block against here. Directory naming
// and collisions are storage mechanics resolved silently in internal/intake,
// never a decision this pass (or the author) needs to make.
type intakePass struct{}

var _ Pass = (*intakePass)(nil)

func (intakePass) Name() string { return "intake" }

func (intakePass) AllowEmpty() bool { return true }

func (intakePass) Precondition(_ *State) error { return nil }

func (intakePass) Body(s *State) (*State, error) {
	ns := s.WithProduced(ArtifactStake)
	ns = ns.WithProduced(ArtifactBrief)
	ns = ns.WithProduced(ArtifactOutline)
	ns = ns.WithProduced(ArtifactAngles)
	ns = ns.WithProduced(ArtifactSourcesStub)
	return ns, nil
}
