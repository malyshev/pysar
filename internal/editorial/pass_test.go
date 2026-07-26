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
	if got := len(Registered()); got != 4 {
		t.Fatalf("expected 4 registered passes (intake, research, draft, discoverability), got %d", got)
	}
}
