package mcpserver

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"pysar/internal/editorial"
)

// registerCheckOnboardingStatus wires check_onboarding_status
// (dec-20260718-ebca6318): /ps-onboard's only tool, and the sole place any
// onboarding orchestration logic lives. It contains no hardcoded knowledge
// of "voice" or "style" specifically -- it constructs a State from what
// exists on disk, then asks every registered editorial.Pass its own
// Precondition, exactly as dec-20260718-ebca6318's invariant requires
// ("discovers outstanding Passes generically via their Precondition").
// Adding a future pass (e.g. a third onboarding kind) requires no change
// here at all.
func (s *Server) registerCheckOnboardingStatus() {
	s.register(
		tool{
			Name:        "check_onboarding_status",
			Description: "Report which onboarding passes (voice, style, and any future kind) are still outstanding for this project. Discovered generically via each pass's own precondition, never hardcoded per kind.",
			InputSchema: map[string]interface{}{"type": "object", "properties": map[string]interface{}{}},
		},
		s.callCheckOnboardingStatus,
	)
}

func (s *Server) callCheckOnboardingStatus(args json.RawMessage) callToolResult {
	var produced []editorial.ArtifactKind
	if fileExists(filepath.Join(s.baseDir, ".pysar", "voice.md")) {
		produced = append(produced, editorial.ArtifactVoiceProfile)
	}
	if fileExists(filepath.Join(s.baseDir, ".pysar", "style.md")) {
		produced = append(produced, editorial.ArtifactStyleProfile)
	}

	// SurfaceBlog is a neutral default: voicePass and stylePass's own
	// Precondition never consults surface, so this choice is inert for
	// onboarding specifically. A future surface-conditional onboarding pass
	// would need this to be a real project surface, not this placeholder.
	state := editorial.NewState(editorial.SurfaceBlog, produced...)

	var b strings.Builder
	count := 0
	for _, p := range editorial.Registered() {
		op, ok := p.(editorial.OnboardingPass)
		if !ok || !op.IsOnboardingPass() {
			continue // per-piece editorial pass (e.g. idea-scaffold), not onboarding
		}
		count++
		if err := p.Precondition(state); err == nil {
			fmt.Fprintf(&b, "- %s: outstanding\n", p.Name())
		} else {
			fmt.Fprintf(&b, "- %s: done (%s)\n", p.Name(), err.Error())
		}
	}
	if count == 0 {
		return textResult("No onboarding passes are registered.")
	}
	return textResult("%s", b.String())
}

func fileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}
