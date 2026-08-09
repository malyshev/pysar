// Package stagereq is the shared opt-in stage-precondition substrate
// (dec-20260809-701b59d3): durable piece policy for required pipeline
// stages, distinct from outcomes (research_mode, presence of seo.md).
package stagereq

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// Stage names persisted on the piece and checked by draft/humanize gates.
const (
	StageResearch = "research" // completed when brief research_mode is full
	StageSEO      = "seo"      // completed when seo.md exists
)

var (
	requiredStagesLineRe = regexp.MustCompile(`(?m)^required_stages:\s*(.*)$`)
	researchModeLineRe   = regexp.MustCompile(`(?m)^research_mode:\s*(\S+)`)
	flowListItemRe       = regexp.MustCompile(`[A-Za-z][A-Za-z0-9_-]*`)
)

// Known reports whether name is a supported required stage.
func Known(name string) bool {
	switch strings.TrimSpace(name) {
	case StageResearch, StageSEO:
		return true
	default:
		return false
	}
}

// Load returns the piece's required_stages from brief.md. Missing brief or
// missing field yields an empty list, not an error — unmarked pieces stay
// ungated (dec-20260809-701b59d3).
func Load(pieceDir string) ([]string, error) {
	brief, err := os.ReadFile(filepath.Join(pieceDir, "brief.md"))
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, fmt.Errorf("read brief.md: %w", err)
	}
	return parseStages(string(brief)), nil
}

// Require merges stages into brief.md's required_stages (idempotent union).
// Unknown stage names are rejected.
func Require(pieceDir string, stages ...string) error {
	normalized, err := normalize(stages)
	if err != nil {
		return err
	}
	if len(normalized) == 0 {
		return fmt.Errorf("stages: at least one of %q or %q is required", StageResearch, StageSEO)
	}

	briefPath := filepath.Join(pieceDir, "brief.md")
	raw, err := os.ReadFile(briefPath)
	if err != nil {
		return fmt.Errorf("read brief.md: %w", err)
	}
	content := string(raw)
	merged := union(parseStages(content), normalized)
	updated, err := setStagesLine(content, merged)
	if err != nil {
		return err
	}
	if updated == content {
		return nil
	}
	if err := os.WriteFile(briefPath, []byte(updated), 0o644); err != nil {
		return fmt.Errorf("write brief.md: %w", err)
	}
	return nil
}

// ResearchComplete is true when brief.md declares research_mode: full.
func ResearchComplete(pieceDir string) (bool, error) {
	brief, err := os.ReadFile(filepath.Join(pieceDir, "brief.md"))
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, fmt.Errorf("read brief.md: %w", err)
	}
	m := researchModeLineRe.FindStringSubmatch(string(brief))
	if len(m) != 2 {
		return false, nil
	}
	return m[1] == "full", nil
}

// SEOComplete is true when seo.md exists in the piece directory.
func SEOComplete(pieceDir string) (bool, error) {
	_, err := os.Stat(filepath.Join(pieceDir, "seo.md"))
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

func normalize(stages []string) ([]string, error) {
	var out []string
	seen := map[string]bool{}
	for _, s := range stages {
		s = strings.TrimSpace(s)
		if s == "" {
			continue
		}
		if !Known(s) {
			return nil, fmt.Errorf("unknown required stage %q (want %q or %q)", s, StageResearch, StageSEO)
		}
		if !seen[s] {
			seen[s] = true
			out = append(out, s)
		}
	}
	return out, nil
}

func parseStages(brief string) []string {
	m := requiredStagesLineRe.FindStringSubmatch(brief)
	if len(m) != 2 {
		return nil
	}
	raw := strings.TrimSpace(m[1])
	if raw == "" || raw == "[]" {
		return nil
	}
	items := flowListItemRe.FindAllString(raw, -1)
	var out []string
	seen := map[string]bool{}
	for _, item := range items {
		if !Known(item) || seen[item] {
			continue
		}
		seen[item] = true
		out = append(out, item)
	}
	return out
}

func union(existing, add []string) []string {
	seen := map[string]bool{}
	var out []string
	for _, s := range existing {
		if !seen[s] {
			seen[s] = true
			out = append(out, s)
		}
	}
	for _, s := range add {
		if !seen[s] {
			seen[s] = true
			out = append(out, s)
		}
	}
	return out
}

func formatStagesLine(stages []string) string {
	return "required_stages: [" + strings.Join(stages, ", ") + "]"
}

func setStagesLine(brief string, stages []string) (string, error) {
	line := formatStagesLine(stages)
	if requiredStagesLineRe.MatchString(brief) {
		return requiredStagesLineRe.ReplaceAllString(brief, line), nil
	}
	// Insert after research_mode so policy sits next to the research outcome field.
	if researchModeLineRe.MatchString(brief) {
		return researchModeLineRe.ReplaceAllStringFunc(brief, func(m string) string {
			return m + "\n" + line
		}), nil
	}
	// Brief without research_mode is unexpected for a real piece; still persist.
	if strings.HasPrefix(brief, "---\n") {
		return "---\n"+line+"\n"+strings.TrimPrefix(brief, "---\n"), nil
	}
	return "", fmt.Errorf("brief.md has no YAML frontmatter to hold required_stages")
}
