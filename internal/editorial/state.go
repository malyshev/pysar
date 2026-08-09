package editorial

// ArtifactKind names something a Pass can produce for a piece.
type ArtifactKind string

const (
	ArtifactStake           ArtifactKind = "stake"
	ArtifactBrief           ArtifactKind = "brief"
	ArtifactOutline         ArtifactKind = "outline"
	ArtifactAngles          ArtifactKind = "angles"
	ArtifactSourcesStub     ArtifactKind = "sources_stub"
	ArtifactSourcesFull     ArtifactKind = "sources_full"
	ArtifactDraft           ArtifactKind = "draft"
	ArtifactStaffEdit       ArtifactKind = "staff_edit"
	ArtifactSharpen         ArtifactKind = "sharpen"
	ArtifactHumanize        ArtifactKind = "humanize"
	ArtifactDiscoverability ArtifactKind = "discoverability"
	ArtifactExported        ArtifactKind = "exported"
	ArtifactVoiceLock       ArtifactKind = "voice_lock"
	ArtifactVoiceProfile    ArtifactKind = "voice_profile"
	ArtifactStyleProfile    ArtifactKind = "style_profile"
)

// TargetSurface names what a piece is ultimately headed for. Passes use it to
// decide whether they apply at all (dec-20260718-20a08f83's conditional-applicability
// property) -- e.g. discoverability is meaningless for a paper letter.
type TargetSurface string

const (
	SurfaceBlog   TargetSurface = "blog"
	SurfaceLetter TargetSurface = "letter"
)

// State is a piece's produced-artifact record: which ArtifactKinds already
// exist for it, plus its target surface, plus opt-in required stages
// (dec-20260809-701b59d3). It is never a single named phase value assuming
// one fixed universal sequence (dec-20260718-20a08f83 invariant).
type State struct {
	produced map[ArtifactKind]bool
	surface  TargetSurface
	required map[string]bool
}

func NewState(surface TargetSurface, produced ...ArtifactKind) *State {
	s := &State{produced: make(map[ArtifactKind]bool), surface: surface}
	for _, k := range produced {
		s.produced[k] = true
	}
	return s
}

func (s *State) Has(k ArtifactKind) bool {
	return s.produced[k]
}

func (s *State) Surface() TargetSurface {
	return s.surface
}

// Requires reports whether the piece policy demands stage (e.g. "research", "seo").
func (s *State) Requires(stage string) bool {
	return s.required[stage]
}

// WithRequired returns a new State with the named stages marked required.
func (s *State) WithRequired(stages ...string) *State {
	ns := s.clone()
	if ns.required == nil {
		ns.required = make(map[string]bool)
	}
	for _, stage := range stages {
		if stage == "" {
			continue
		}
		ns.required[stage] = true
	}
	return ns
}

// WithProduced returns a new State with k added, leaving s unmodified.
func (s *State) WithProduced(k ArtifactKind) *State {
	ns := s.clone()
	ns.produced[k] = true
	return ns
}

func (s *State) clone() *State {
	ns := NewState(s.surface)
	for existing := range s.produced {
		ns.produced[existing] = true
	}
	if len(s.required) > 0 {
		ns.required = make(map[string]bool, len(s.required))
		for k, v := range s.required {
			ns.required[k] = v
		}
	}
	return ns
}
