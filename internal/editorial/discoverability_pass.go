package editorial

// discoverabilityPass is the SEO/discoverability packaging pass -- an
// opt-in checklist (citation resolution, title/subtitle tuning, tags, meta
// description, URL slug, scannability) that runs before humanize, never
// after (dec-20260804-e3234e50, superseding dec-20260718-e9f5b5e6's
// earlier full removal of platform publish-checklist scope). File I/O
// lives in internal/seo + MCP; this Pass only owns precondition and
// produced-artifact state updates, same shape as every other real pass.
//
// Kept as "discoverability" (not renamed to "seo") to preserve the
// already-verified precondition behavior recorded against
// dec-20260718-20a08f83: declines for SurfaceLetter (a discoverability
// checklist is meaningless for a paper letter) via the same Precondition
// mechanism every pass uses, no separate skip-flag system. The author-
// facing name is /ps-seo and --seo; this Go type's name is a Pysar-
// internal detail, not a promise to the author surface.
type discoverabilityPass struct{}

var _ Pass = (*discoverabilityPass)(nil)

func (discoverabilityPass) Name() string { return "discoverability" }

func (discoverabilityPass) Precondition(s *State) error {
	if s.Surface() == SurfaceLetter {
		return &PreconditionError{
			Hint: "discoverability does not apply to a paper-letter target surface",
		}
	}
	if !s.Has(ArtifactDraft) {
		return &PreconditionError{
			Missing: []ArtifactKind{ArtifactDraft},
			Hint:    "no draft.md found for this piece -- either /ps-draft hasn't run yet, or piece_path doesn't point at a real piece (check for a typo)",
		}
	}
	if s.Has(ArtifactHumanize) {
		return &PreconditionError{
			Hint: "humanize.md already exists for this piece -- discoverability must run BEFORE humanize, never after (dec-20260804-e3234e50): running it now would silently smooth prose humanize deliberately roughened",
		}
	}
	return nil
}

func (discoverabilityPass) Body(s *State) (*State, error) {
	return s.WithProduced(ArtifactDiscoverability), nil
}
