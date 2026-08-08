package main

import (
	"fmt"
	"os"
	"path/filepath"
)

// hostAdapter is one AI-editor host's scaffold dialect (dec-20260808-f3001106).
// Skill bodies live in the shared skillAssets corpus; adapters own only install
// paths, MCP config shape, and host-specific permission/settings files.
type hostAdapter interface {
	Name() string
	Scaffold(dir string, force bool) error
}

// hostRegistry is the single registration point for real hosts. Adding a host
// means implementing hostAdapter and one map entry — not a new flag-dispatch
// shape (dec-20260718-8278c494 consequence; dec-20260808-f3001106).
var hostRegistry = map[string]hostAdapter{
	"claude": claudeHost{},
	"cursor": cursorHost{},
}

func resolveHost(claude, cursor, codex bool) (hostAdapter, error) {
	switch {
	case cursor:
		return hostRegistry["cursor"], nil
	case codex:
		return nil, fmt.Errorf("pysar init --codex: not yet supported")
	default:
		// --claude and no-flag both resolve here (dec-20260718-8278c494).
		_ = claude
		return hostRegistry["claude"], nil
	}
}

type claudeHost struct{}

func (claudeHost) Name() string { return "claude" }

func (claudeHost) Scaffold(dir string, force bool) error {
	alreadySetUp, err := writeProjectManifest(dir, "claude")
	if err != nil {
		return fmt.Errorf("pysar init: %w", err)
	}

	claudePath := filepath.Join(dir, "CLAUDE.md")
	if _, _, err := writeIfAbsent(claudePath, claudeMDTemplate, force); err != nil {
		return fmt.Errorf("pysar init: %w", err)
	}

	home, err := homeDir()
	if err != nil {
		return fmt.Errorf("pysar init: locate home directory for global skill install: %w", err)
	}

	skillOut, err := writeHostSkills(home, filepath.Join(".claude", "skills"), force)
	if err != nil {
		return fmt.Errorf("pysar init: %w", err)
	}
	templateOut, err := writeBuiltinTemplates(home, force)
	if err != nil {
		return fmt.Errorf("pysar init: %w", err)
	}

	settingsPath := filepath.Join(dir, ".claude", "settings.json")
	if err := os.MkdirAll(filepath.Dir(settingsPath), 0o755); err != nil {
		return fmt.Errorf("pysar init: %w", err)
	}
	if _, _, err := writeIfAbsent(settingsPath, claudeSettingsJSON, force); err != nil {
		return fmt.Errorf("pysar init: %w", err)
	}

	mcpPath := filepath.Join(dir, ".mcp.json")
	if _, _, err := writeIfAbsent(mcpPath, claudeMCPJSON, force); err != nil {
		return fmt.Errorf("pysar init: %w", err)
	}

	printInitSummary(dir, alreadySetUp, skillOut, templateOut, force)
	return nil
}

type cursorHost struct{}

func (cursorHost) Name() string { return "cursor" }

func (cursorHost) Scaffold(dir string, force bool) error {
	alreadySetUp, err := writeProjectManifest(dir, "cursor")
	if err != nil {
		return fmt.Errorf("pysar init: %w", err)
	}

	home, err := homeDir()
	if err != nil {
		return fmt.Errorf("pysar init: locate home directory for global skill install: %w", err)
	}

	// Same skill corpus as Claude; Cursor discovers skills under ~/.cursor/skills
	// (haft CL1 / dec-20260808-f3001106).
	skillOut, err := writeHostSkills(home, filepath.Join(".cursor", "skills"), force)
	if err != nil {
		return fmt.Errorf("pysar init: %w", err)
	}
	templateOut, err := writeBuiltinTemplates(home, force)
	if err != nil {
		return fmt.Errorf("pysar init: %w", err)
	}

	mcpPath := filepath.Join(dir, ".cursor", "mcp.json")
	if err := os.MkdirAll(filepath.Dir(mcpPath), 0o755); err != nil {
		return fmt.Errorf("pysar init: %w", err)
	}
	if _, _, err := writeIfAbsent(mcpPath, cursorMCPJSON, force); err != nil {
		return fmt.Errorf("pysar init: %w", err)
	}

	printInitSummary(dir, alreadySetUp, skillOut, templateOut, force)
	return nil
}

type syncResult struct {
	installed, autoRefreshed, refreshed, ambiguous []string
}

func writeHostSkills(home, skillsRel string, force bool) (syncResult, error) {
	skillsDir := filepath.Join(home, skillsRel)
	installed, auto, refreshed, ambiguous, err := syncManagedTree(skillAssets, "assets/skills", skillsDir, force)
	return syncResult{installed, auto, refreshed, ambiguous}, err
}

func writeBuiltinTemplates(home string, force bool) (syncResult, error) {
	voiceInstalled, voiceAuto, voiceRefreshed, voiceAmbiguous, err := writeVoiceTemplates(home, force)
	if err != nil {
		return syncResult{}, err
	}
	styleInstalled, styleAuto, styleRefreshed, styleAmbiguous, err := writeStyleTemplates(home, force)
	if err != nil {
		return syncResult{}, err
	}
	return syncResult{
		installed:     append(voiceInstalled, styleInstalled...),
		autoRefreshed: append(voiceAuto, styleAuto...),
		refreshed:     append(voiceRefreshed, styleRefreshed...),
		ambiguous:     append(voiceAmbiguous, styleAmbiguous...),
	}, nil
}

func printInitSummary(dir string, alreadySetUp bool, skills, templates syncResult, force bool) {
	if alreadySetUp {
		fmt.Printf("pysar init: %s is already set up\n", projectLabel(dir))
	} else {
		fmt.Printf("pysar init: %s is set up and ready\n", projectLabel(dir))
	}
	if len(skills.installed) > 0 {
		fmt.Println("pysar init: installed agentic skills (shared across all pysar projects on this machine)")
	}
	if len(skills.autoRefreshed) > 0 {
		fmt.Println("pysar init: updated agentic skills to the current version")
	}
	if len(skills.refreshed) > 0 {
		fmt.Println("pysar init: refreshed agentic skills to the current version (--force)")
	}
	if len(templates.installed) > 0 {
		fmt.Println("pysar init: installed built-in templates (shared across all pysar projects on this machine)")
	}
	if len(templates.autoRefreshed) > 0 {
		fmt.Println("pysar init: updated built-in templates to the current version")
	}
	if len(templates.refreshed) > 0 {
		fmt.Println("pysar init: refreshed built-in templates to the current version (--force)")
	}
	if !force {
		switch {
		case len(skills.ambiguous) > 0 && len(templates.ambiguous) > 0:
			fmt.Println("pysar init: some agentic skills and built-in templates look customized or from an unrecognized version -- left untouched; rerun with --force to overwrite them")
		case len(skills.ambiguous) > 0:
			fmt.Println("pysar init: some agentic skills look customized or from an unrecognized version -- left untouched; rerun with --force to overwrite them")
		case len(templates.ambiguous) > 0:
			fmt.Println("pysar init: some built-in templates look customized or from an unrecognized version -- left untouched; rerun with --force to overwrite them")
		}
	}
}
