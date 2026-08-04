package editorial

import (
	"errors"
	"testing"
)

func TestDraftBlockedWithoutBrief(t *testing.T) {
	s := NewState(SurfaceBlog)
	_, err := Run(draftPass{}, s)

	var pe *PreconditionError
	if !errors.As(err, &pe) {
		t.Fatalf("expected *PreconditionError, got %v", err)
	}
	if len(pe.Missing) != 1 || pe.Missing[0] != ArtifactBrief {
		t.Fatalf("expected missing=[brief], got %v", pe.Missing)
	}
}

func TestDraftProceedsWithBrief(t *testing.T) {
	s := NewState(SurfaceBlog, ArtifactStake, ArtifactBrief)
	ns, err := Run(draftPass{}, s)
	if err != nil {
		t.Fatalf("expected draft to proceed once a brief exists, got %v", err)
	}
	if !ns.Has(ArtifactDraft) {
		t.Fatal("expected ArtifactDraft to be produced")
	}
}

func TestStaffEditBlockedWithoutDraft(t *testing.T) {
	s := NewState(SurfaceBlog, ArtifactStake, ArtifactBrief)
	_, err := Run(staffEditPass{}, s)

	var pe *PreconditionError
	if !errors.As(err, &pe) {
		t.Fatalf("expected *PreconditionError, got %v", err)
	}
	if len(pe.Missing) != 1 || pe.Missing[0] != ArtifactDraft {
		t.Fatalf("expected missing=[draft], got %v", pe.Missing)
	}
}

func TestStaffEditProceedsWithDraft(t *testing.T) {
	s := NewState(SurfaceBlog, ArtifactStake, ArtifactBrief, ArtifactDraft)
	ns, err := Run(staffEditPass{}, s)
	if err != nil {
		t.Fatalf("expected staff-edit to proceed once a draft exists, got %v", err)
	}
	if !ns.Has(ArtifactStaffEdit) {
		t.Fatal("expected ArtifactStaffEdit to be produced")
	}
}

func TestSharpenBlockedWithoutDraft(t *testing.T) {
	s := NewState(SurfaceBlog, ArtifactStake, ArtifactBrief)
	_, err := Run(sharpenPass{}, s)

	var pe *PreconditionError
	if !errors.As(err, &pe) {
		t.Fatalf("expected *PreconditionError, got %v", err)
	}
	if len(pe.Missing) != 1 || pe.Missing[0] != ArtifactDraft {
		t.Fatalf("expected missing=[draft], got %v", pe.Missing)
	}
}

func TestSharpenProceedsWithDraftAlone(t *testing.T) {
	// Staff-edit is optional -- sharpen must proceed on a draft alone, the
	// same way draft itself proceeds without research having run.
	s := NewState(SurfaceBlog, ArtifactStake, ArtifactBrief, ArtifactDraft)
	ns, err := Run(sharpenPass{}, s)
	if err != nil {
		t.Fatalf("expected sharpen to proceed once a draft exists, got %v", err)
	}
	if !ns.Has(ArtifactSharpen) {
		t.Fatal("expected ArtifactSharpen to be produced")
	}
}

func TestSharpenProceedsWithStaffEditToo(t *testing.T) {
	s := NewState(SurfaceBlog, ArtifactStake, ArtifactBrief, ArtifactDraft, ArtifactStaffEdit)
	if _, err := Run(sharpenPass{}, s); err != nil {
		t.Fatalf("expected sharpen to proceed when staff-edit also ran, got %v", err)
	}
}

func TestHumanizeBlockedWithoutDraft(t *testing.T) {
	s := NewState(SurfaceBlog, ArtifactStake, ArtifactBrief)
	_, err := Run(humanizePass{}, s)

	var pe *PreconditionError
	if !errors.As(err, &pe) {
		t.Fatalf("expected *PreconditionError, got %v", err)
	}
	if len(pe.Missing) != 1 || pe.Missing[0] != ArtifactDraft {
		t.Fatalf("expected missing=[draft], got %v", pe.Missing)
	}
}

func TestHumanizeProceedsWithDraftAlone(t *testing.T) {
	// Staff-edit and sharpen are both optional -- humanize must proceed on
	// a draft alone, the same way sharpen proceeds without staff-edit.
	s := NewState(SurfaceBlog, ArtifactStake, ArtifactBrief, ArtifactDraft)
	ns, err := Run(humanizePass{}, s)
	if err != nil {
		t.Fatalf("expected humanize to proceed once a draft exists, got %v", err)
	}
	if !ns.Has(ArtifactHumanize) {
		t.Fatal("expected ArtifactHumanize to be produced")
	}
}

func TestHumanizeProceedsWithStaffEditAndSharpenToo(t *testing.T) {
	s := NewState(SurfaceBlog, ArtifactStake, ArtifactBrief, ArtifactDraft, ArtifactStaffEdit, ArtifactSharpen)
	if _, err := Run(humanizePass{}, s); err != nil {
		t.Fatalf("expected humanize to proceed when staff-edit and sharpen also ran, got %v", err)
	}
}

func TestExportBlockedWithoutDraft(t *testing.T) {
	s := NewState(SurfaceBlog, ArtifactStake, ArtifactBrief)
	_, err := Run(exportPass{}, s)

	var pe *PreconditionError
	if !errors.As(err, &pe) {
		t.Fatalf("expected *PreconditionError, got %v", err)
	}
	if len(pe.Missing) != 1 || pe.Missing[0] != ArtifactDraft {
		t.Fatalf("expected missing=[draft], got %v", pe.Missing)
	}
}

func TestExportProceedsWithDraftAlone(t *testing.T) {
	// Staff-edit, sharpen, and humanize are all optional -- export must
	// work on whatever stage a piece actually reached, even just a draft.
	s := NewState(SurfaceBlog, ArtifactStake, ArtifactBrief, ArtifactDraft)
	ns, err := Run(exportPass{}, s)
	if err != nil {
		t.Fatalf("expected export to proceed once a draft exists, got %v", err)
	}
	if !ns.Has(ArtifactExported) {
		t.Fatal("expected ArtifactExported to be produced")
	}
}

func TestResearchBlockedWithoutBrief(t *testing.T) {
	s := NewState(SurfaceBlog)
	_, err := Run(researchPass{}, s)

	var pe *PreconditionError
	if !errors.As(err, &pe) {
		t.Fatalf("expected *PreconditionError, got %v", err)
	}
	if len(pe.Missing) != 1 || pe.Missing[0] != ArtifactBrief {
		t.Fatalf("expected missing=[brief], got %v", pe.Missing)
	}
}

func TestResearchProceedsWithBrief(t *testing.T) {
	s := NewState(SurfaceBlog, ArtifactStake, ArtifactBrief)
	ns, err := Run(researchPass{}, s)
	if err != nil {
		t.Fatalf("expected research to proceed once a brief exists, got %v", err)
	}
	if !ns.Has(ArtifactSourcesFull) {
		t.Fatal("expected ArtifactSourcesFull to be produced")
	}
}

func TestDiscoverabilityDeclinesForLetter(t *testing.T) {
	s := NewState(SurfaceLetter, ArtifactStake, ArtifactDraft)
	if _, err := Run(discoverabilityPass{}, s); err == nil {
		t.Fatal("expected discoverability to decline for SurfaceLetter")
	}
}

func TestDiscoverabilityProceedsForBlogWithDraft(t *testing.T) {
	s := NewState(SurfaceBlog, ArtifactStake, ArtifactDraft)
	if _, err := Run(discoverabilityPass{}, s); err != nil {
		t.Fatalf("expected discoverability to proceed for blog with draft, got %v", err)
	}
}

func TestDiscoverabilityDeclinesAfterHumanizeAlreadyRan(t *testing.T) {
	// Ordering invariant (dec-20260804-e3234e50): discoverability must run
	// BEFORE humanize, never after -- enforced here via the same
	// Precondition mechanism every pass uses, not left to prose alone.
	s := NewState(SurfaceBlog, ArtifactStake, ArtifactDraft, ArtifactHumanize)
	if _, err := Run(discoverabilityPass{}, s); err == nil {
		t.Fatal("expected discoverability to decline once humanize.md already exists")
	}
}

// passHandlesEmptyCorrectly reports whether p either explicitly allows an
// empty piece state (AllowEmpty) or correctly rejects it via a non-nil
// Precondition error. Shared by the registry-wide meta-test and the
// meta-test's own self-check below.
func passHandlesEmptyCorrectly(p Pass, empty *State) bool {
	allow := false
	if ae, ok := p.(allowsEmpty); ok {
		allow = ae.AllowEmpty()
	}
	err := p.Precondition(empty)
	if allow {
		return err == nil
	}
	return err != nil
}

func TestRegisteredPassesRejectEmptyUnlessAllowEmpty(t *testing.T) {
	empty := NewState(SurfaceBlog)
	for _, p := range Registered() {
		if !passHandlesEmptyCorrectly(p, empty) {
			t.Errorf("%s: does not correctly handle empty state (AllowEmpty mismatch or trivial precondition)", p.Name())
		}
	}
}

// brokenNilPass is a deliberately-broken stub: not AllowEmpty, but its
// Precondition trivially returns nil regardless of state. It exists only to
// prove the meta-check catches this class of bug -- it is never registered.
type brokenNilPass struct{}

var _ Pass = (*brokenNilPass)(nil)

func (brokenNilPass) Name() string                  { return "broken-nil-pass" }
func (brokenNilPass) Precondition(_ *State) error   { return nil }
func (brokenNilPass) Body(s *State) (*State, error) { return s, nil }

func TestMetaCheckCatchesTrivialNilPrecondition(t *testing.T) {
	empty := NewState(SurfaceBlog)

	if passHandlesEmptyCorrectly(brokenNilPass{}, empty) {
		t.Fatal("expected meta-check to flag a non-AllowEmpty pass with a trivially-nil Precondition")
	}

	if !passHandlesEmptyCorrectly(draftPass{}, empty) {
		t.Fatal("expected a correctly implemented pass to pass the meta-check")
	}
}

func TestRegisteredLength(t *testing.T) {
	if got := len(Registered()); got != 8 {
		t.Fatalf("expected 8 registered passes (intake, research, draft, staff-edit, sharpen, humanize, export, discoverability), got %d", got)
	}
}
