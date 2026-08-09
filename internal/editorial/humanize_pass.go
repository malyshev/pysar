package editorial

import "pysar/internal/stagereq"

// humanizePass is the sixth real per-piece editorial pass (after intake,
// research, draft, staff-edit, sharpen). File I/O lives in internal/
// humanize + MCP; this Pass only owns precondition + produced-artifact
// state updates.
//
// Precondition requires ArtifactDraft, not ArtifactSharpen or
// ArtifactStaffEdit -- both are optional (same reasoning as research being
// optional for draft), so humanize must work whether or not either ran,
// reading whichever is the piece's latest revision (seo.md if present,
// else sharpen.md, else staff-edit.md, else draft.md). When the piece
// lists seo in required_stages (dec-20260809-701b59d3),
// ArtifactDiscoverability (seo.md) is also required before humanize may
// proceed — precondition awareness only, not SEO content rewriting.
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
	if s.Requires(stagereq.StageSEO) && !s.Has(ArtifactDiscoverability) {
		return &PreconditionError{
			Missing: []ArtifactKind{ArtifactDiscoverability},
			Hint:    "this piece requires SEO packaging before humanize (required_stages includes seo / /ps --seo) -- run /ps-seo and save_seo_bundle so seo.md exists",
		}
	}
	return nil
}

func (humanizePass) Body(s *State) (*State, error) {
	return s.WithProduced(ArtifactHumanize), nil
}
