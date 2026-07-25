package editorial

// ArtifactKind names something a Pass can produce for a piece.
type ArtifactKind string

const (
	ArtifactStake        ArtifactKind = "stake"
	ArtifactBrief        ArtifactKind = "brief"
	ArtifactOutline      ArtifactKind = "outline"
	ArtifactAngles       ArtifactKind = "angles"
	ArtifactSourcesStub  ArtifactKind = "sources_stub"
	ArtifactSourcesFull  ArtifactKind = "sources_full"
	ArtifactDraft        ArtifactKind = "draft"
	ArtifactVoiceLock    ArtifactKind = "voice_lock"
	ArtifactVoiceProfile ArtifactKind = "voice_profile"
	ArtifactStyleProfile ArtifactKind = "style_profile"
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
// exist for it, plus its target surface. It is never a single named phase
// value assuming one fixed universal sequence (dec-20260718-20a08f83 invariant).
type State struct {
	produced map[ArtifactKind]bool
	surface  TargetSurface
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

// WithProduced returns a new State with k added, leaving s unmodified.
func (s *State) WithProduced(k ArtifactKind) *State {
	ns := NewState(s.surface)
	for existing := range s.produced {
		ns.produced[existing] = true
	}
	ns.produced[k] = true
	return ns
}
