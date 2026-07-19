package onboarding

import (
	"pysar/internal/editorial"
	"testing"
)

func TestStylePassPreconditionDeclinesWhenAlreadyProduced(t *testing.T) {
	s := editorial.NewState(editorial.SurfaceBlog, editorial.ArtifactStyleProfile)
	if err := (stylePass{}).Precondition(s); err == nil {
		t.Fatal("expected Precondition to decline once style profile is already produced")
	}
}

func TestStylePassRunProducesArtifactOnEmptyState(t *testing.T) {
	s := editorial.NewState(editorial.SurfaceBlog)
	next, err := editorial.Run(stylePass{}, s)
	if err != nil {
		t.Fatalf("expected stylePass to run on empty state, got %v", err)
	}
	if !next.Has(editorial.ArtifactStyleProfile) {
		t.Fatal("expected resulting state to have ArtifactStyleProfile")
	}
}

func TestStylePassIsRegistered(t *testing.T) {
	for _, p := range editorial.Registered() {
		if p.Name() == "ps-style" {
			return
		}
	}
	t.Fatal("expected ps-style to be registered in the editorial Pass registry")
}

func TestStylePassDoesNotInterfereWithVoicePass(t *testing.T) {
	// Producing one artifact must never satisfy the other pass's precondition.
	s := editorial.NewState(editorial.SurfaceBlog, editorial.ArtifactVoiceProfile)
	if err := (stylePass{}).Precondition(s); err != nil {
		t.Fatalf("expected stylePass to still be outstanding when only voice is produced, got %v", err)
	}

	s2 := editorial.NewState(editorial.SurfaceBlog, editorial.ArtifactStyleProfile)
	if err := (voicePass{}).Precondition(s2); err != nil {
		t.Fatalf("expected voicePass to still be outstanding when only style is produced, got %v", err)
	}
}

// TestOnboardingPassMarkerDistinguishesFromPieceLevelPasses locks in the
// filter check_onboarding_status relies on: voicePass/stylePass mark
// themselves via editorial.OnboardingPass; the pre-existing per-piece seam
// stubs (idea-scaffold, draft-edit, discoverability) do not implement it at
// all, so a type assertion correctly excludes them without naming them.
func TestOnboardingPassMarkerDistinguishesFromPieceLevelPasses(t *testing.T) {
	var onboardingNames, otherNames []string
	for _, p := range editorial.Registered() {
		if op, ok := p.(editorial.OnboardingPass); ok && op.IsOnboardingPass() {
			onboardingNames = append(onboardingNames, p.Name())
		} else {
			otherNames = append(otherNames, p.Name())
		}
	}

	wantOnboarding := map[string]bool{"ps-voice": true, "ps-style": true}
	for _, name := range onboardingNames {
		if !wantOnboarding[name] {
			t.Errorf("unexpected pass marked as onboarding: %s", name)
		}
		delete(wantOnboarding, name)
	}
	if len(wantOnboarding) != 0 {
		t.Errorf("expected passes not found as onboarding: %v", wantOnboarding)
	}

	for _, name := range otherNames {
		if name == "ps-voice" || name == "ps-style" {
			t.Errorf("%s incorrectly excluded from onboarding", name)
		}
	}
}
