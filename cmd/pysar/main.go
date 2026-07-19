package main

import (
	"crypto/sha256"
	"embed"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"

	"pysar/internal/mcpserver"
	"pysar/internal/onboarding"
)

const version = "0.0.1-dev" // Dev version string until release tagging exists

const schemaVersion = 1

// claudeMDTemplate is intentionally minimal and provisional (dec-20260718-fe926229) --
// host-file conventions have not been decided as their own DecisionRecord yet.
const claudeMDTemplate = `# Pysar project

Initialized by ` + "`pysar init --claude`" + `. Host-file conventions here are
provisional and may change as Pysar's Claude Code integration develops.
`

// claudeSkillAssets holds every ps-* skill shipped for the Claude Code host
// surface (dec-20260718-a2a24f32), embedded into the binary so a single
// static executable can scaffold them into any project (dec-20260718-239680d4).
// New skills (e.g. ps-style) need only a new directory under assets/claude/skills --
// scaffoldClaude walks whatever is here, no per-skill wiring.
//
//go:embed assets/claude/skills
var claudeSkillAssets embed.FS

// claudeSettingsJSON pre-approves the pysar MCP server's tools plus reading
// .pysar/** (dec-20260719-fa0366dd) -- ps-* skills persist author content by
// calling save_voice_profile (and future save_* tools) on pysar serve, never
// via a raw client-side Write/Edit tool call, so the whole class of Claude
// Code permission surfaces this project hit (Bash gate, Edit/Write rule-name
// mismatch, file-creation-vs-edit gate) collapses to one: "may Claude call
// mcp__pysar__*".
//
//go:embed assets/claude/settings.json
var claudeSettingsJSON string

// claudeMCPJSON registers pysar serve as this project's local MCP server so
// Claude Code launches it automatically -- mirrors this repo's own .mcp.json
// entry for haft, proven working throughout this project's own development.
//
//go:embed assets/claude/mcp.json
var claudeMCPJSON string

// templateAssets holds every built-in reusable content template (starting
// with the one "generic" voice default), seeded into the operator's
// cross-project template store (dec-20260719-3e36577e) the same way
// claudeSkillAssets seeds ~/.claude/skills/ -- same embed-and-sync shape,
// different target directory, host-agnostic since templates are pysar's own
// data, not a Claude-Code-specific config surface.
//
//go:embed assets/templates
var templateAssets embed.FS

type projectManifest struct {
	SchemaVersion int    `json:"schema_version"`
	Host          string `json:"host"`
}

func main() {
	if err := newRootCmd().Execute(); err != nil {
		os.Exit(1)
	}
}

func newRootCmd() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "pysar",
		Short: "An author-directed editorial engine for your writing projects",
		Long: `Pysar — an author-directed editorial engine for your writing projects.

Bring your take (an idea or a draft). Pysar helps you shape it into a ship-ready
piece you trust — without forcing pipeline jargon on you, and without posting
on your behalf.

This binary is the console CLI. Agentic slash commands (ps-* / /ps) are a
separate host-agent surface when installed; they are not invoked by typing pysar.`,
		Example: `  pysar init              # scaffold a Claude Code project in .
  pysar init --claude ./my-piece
  pysar --version`,
		Version:      version,
		SilenceUsage: true,
	}
	rootCmd.SetVersionTemplate("pysar {{.Version}}\n")
	rootCmd.AddCommand(newInitCmd())
	rootCmd.AddCommand(newServeCmd())
	return rootCmd
}

func newInitCmd() *cobra.Command {
	var claude, cursor, codex, force bool

	cmd := &cobra.Command{
		Use:   "init [dir]",
		Short: "Scaffold a writing project for a host agent (default: Claude Code)",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			dir := "."
			if len(args) > 0 {
				dir = args[0]
			}

			switch {
			case cursor:
				return fmt.Errorf("pysar init --cursor: not yet supported")
			case codex:
				return fmt.Errorf("pysar init --codex: not yet supported")
			default:
				return scaffoldClaude(dir, force)
			}
		},
	}

	cmd.Flags().BoolVar(&claude, "claude", false, "scaffold for Claude Code (default)")
	cmd.Flags().BoolVar(&cursor, "cursor", false, "scaffold for Cursor (not yet supported)")
	cmd.Flags().BoolVar(&codex, "codex", false, "scaffold for Codex (not yet supported)")
	cmd.Flags().BoolVar(&force, "force", false, "refresh CLAUDE.md, agentic skills, and permission/MCP config to the current shipped version, even if already installed -- never touches .pysar/ project data")
	cmd.MarkFlagsMutuallyExclusive("claude", "cursor", "codex")

	return cmd
}

// newServeCmd runs the Pysar MCP server over stdio (dec-20260719-fa0366dd).
// A pysar-init'd project's .mcp.json launches this automatically; it is not
// meant to be run by hand. Project root comes from PYSAR_PROJECT_ROOT (set
// by .mcp.json's ${PWD:-.} substitution, mirroring this repo's own working
// .mcp.json entry for haft) and falls back to the process's own working
// directory if unset.
func newServeCmd() *cobra.Command {
	return &cobra.Command{
		Use:    "serve",
		Short:  "Run the Pysar MCP server (stdio) that ps-* skills use to persist author content",
		Args:   cobra.NoArgs,
		Hidden: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			dir := os.Getenv("PYSAR_PROJECT_ROOT")
			if dir == "" {
				wd, err := os.Getwd()
				if err != nil {
					return fmt.Errorf("pysar serve: %w", err)
				}
				dir = wd
			}
			home, err := homeDir()
			if err != nil {
				return fmt.Errorf("pysar serve: locate home directory: %w", err)
			}
			srv := mcpserver.New("pysar", version, dir, home, os.Stdin, os.Stdout)
			return srv.Run()
		},
	}
}

// homeDir resolves the user's home directory. It is a package-level var,
// not a direct os.UserHomeDir() call, so tests can override it and avoid
// ever writing into a real developer's actual ~/.claude/skills/.
var homeDir = os.UserHomeDir

// projectLabel turns an init target into a name a human actually reads as a
// word -- the project's own folder name, never the raw argument the author
// typed. This matters most for the default, most common invocation (`pysar
// init` with no argument, dir="."), where printing the argument verbatim
// would read as broken output (". is already set up"), not a real sentence.
func projectLabel(dir string) string {
	abs, err := filepath.Abs(dir)
	if err != nil {
		return "this project"
	}
	name := filepath.Base(abs)
	if name == "" || name == "." || name == string(filepath.Separator) {
		return "this project"
	}
	return name
}

// scaffoldClaude writes the .pysar/project machine manifest (dec-20260718-4b222794),
// a minimal, provisional CLAUDE.md host file (dec-20260718-fe926229), installs
// every embedded ps-* skill into the user's global ~/.claude/skills/, and
// writes project-local .claude/settings.json and .mcp.json so pysar serve is
// registered and its tools pre-approved (dec-20260719-fa0366dd) -- ps-* skills
// persist author content through it, never through a raw client-side
// Write/Edit tool call.
//
// Skills are installed globally, not per-project: Claude Code only discovers
// skills from ~/.claude/skills/, never from a project-local .claude/skills/
// (confirmed by direct testing against a real pysar-init'd project; see
// upstream https://github.com/anthropics/claude-code/issues/33733, closed as
// not planned). This also fits Pysar's own shape: a skill's prompt is
// project-agnostic -- only the data it produces (.pysar/voice.md) is
// project-scoped, and that already lives under the project itself.
func scaffoldClaude(dir string, force bool) error {
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
	installed, autoRefreshed, refreshed, ambiguous, err := writeClaudeSkills(home, force)
	if err != nil {
		return fmt.Errorf("pysar init: %w", err)
	}

	voiceInstalled, voiceAutoRefreshed, voiceRefreshed, voiceAmbiguous, err := writeVoiceTemplates(home, force)
	if err != nil {
		return fmt.Errorf("pysar init: %w", err)
	}
	styleInstalled, styleAutoRefreshed, styleRefreshed, styleAmbiguous, err := writeStyleTemplates(home, force)
	if err != nil {
		return fmt.Errorf("pysar init: %w", err)
	}
	templatesInstalled := append(voiceInstalled, styleInstalled...)
	templatesAutoRefreshed := append(voiceAutoRefreshed, styleAutoRefreshed...)
	templatesRefreshed := append(voiceRefreshed, styleRefreshed...)
	templatesAmbiguous := append(voiceAmbiguous, styleAmbiguous...)

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

	if alreadySetUp {
		fmt.Printf("pysar init: %s is already set up\n", projectLabel(dir))
	} else {
		fmt.Printf("pysar init: %s is set up and ready\n", projectLabel(dir))
	}
	if len(installed) > 0 {
		fmt.Println("pysar init: installed agentic skills (shared across all pysar projects on this machine)")
	}
	if len(autoRefreshed) > 0 {
		fmt.Println("pysar init: updated agentic skills to the current version")
	}
	if len(refreshed) > 0 {
		fmt.Println("pysar init: refreshed agentic skills to the current version (--force)")
	}
	if len(templatesInstalled) > 0 {
		fmt.Println("pysar init: installed built-in templates (shared across all pysar projects on this machine)")
	}
	if len(templatesAutoRefreshed) > 0 {
		fmt.Println("pysar init: updated built-in templates to the current version")
	}
	if len(templatesRefreshed) > 0 {
		fmt.Println("pysar init: refreshed built-in templates to the current version (--force)")
	}
	// Skills and templates share one "rerun with --force" nudge: printing it
	// once when either (or both) look customized says everything the other
	// copy would have, so a run with both ambiguous never repeats itself.
	if !force {
		switch {
		case len(ambiguous) > 0 && len(templatesAmbiguous) > 0:
			fmt.Println("pysar init: some agentic skills and built-in templates look customized or from an unrecognized version -- left untouched; rerun with --force to overwrite them")
		case len(ambiguous) > 0:
			fmt.Println("pysar init: some agentic skills look customized or from an unrecognized version -- left untouched; rerun with --force to overwrite them")
		case len(templatesAmbiguous) > 0:
			fmt.Println("pysar init: some built-in templates look customized or from an unrecognized version -- left untouched; rerun with --force to overwrite them")
		}
	}
	// Deliberately silent, without --force, when a project-local file (CLAUDE.md,
	// .claude/settings.json, .mcp.json) differs from what pysar would write: with
	// no update mechanism for those files, "this differs" was a fact the operator
	// couldn't act on -- pure noise. --force is that mechanism now; see
	// note-20260719-d883dd47 and note-20260719-82ea61c0. Global skills instead get
	// automatic per-file refresh (dec-20260719-25712417); the "ambiguous" nudge
	// above is that mechanism's own deliberately narrow exception.
	return nil
}

// skillManifestFile records, per skill file, the hash of the shipped content
// that pysar itself last wrote there. It lives inside home/.claude/skills/
// alongside the skills it tracks -- system-owned state for system-owned
// files, never read or written by an author.
const skillManifestFile = ".pysar-manifest.json"

// skillManifest is the on-disk shape of skillManifestFile. Files maps a
// skill's embed-relative path (e.g. "ps-voice/SKILL.md") to the sha256 hex
// digest of the shipped content pysar last wrote at that path.
type skillManifest struct {
	Files map[string]string `json:"files"`
}

// loadSkillManifest treats a missing OR corrupt manifest the same way: start
// from empty. The manifest is a cache pysar itself owns to detect which
// shipped files an author has customized -- losing it just means every
// skill file is (safely) treated as needing a fresh write on this run, not a
// permanent failure. A hard error here would otherwise brick every future
// `pysar init` for an operator whose manifest was left truncated by an
// interrupted write, with no way to recover short of finding and deleting it
// by hand.
func loadSkillManifest(path string) (skillManifest, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return skillManifest{Files: map[string]string{}}, nil
		}
		return skillManifest{}, err
	}
	var m skillManifest
	if err := json.Unmarshal(data, &m); err != nil {
		return skillManifest{Files: map[string]string{}}, nil
	}
	if m.Files == nil {
		m.Files = map[string]string{}
	}
	return m, nil
}

// saveSkillManifest writes via temp-file-then-rename so a process killed
// mid-write (crash, OOM, ctrl-C) leaves either the old manifest or the new
// one intact on disk, never a truncated file in between.
func saveSkillManifest(path string, m skillManifest) error {
	data, err := json.MarshalIndent(m, "", "  ")
	if err != nil {
		return err
	}
	data = append(data, '\n')

	tmp, err := os.CreateTemp(filepath.Dir(path), filepath.Base(path)+".tmp-*")
	if err != nil {
		return err
	}
	tmpPath := tmp.Name()
	defer os.Remove(tmpPath) // no-op once the rename below succeeds

	if _, err := tmp.Write(data); err != nil {
		tmp.Close()
		return err
	}
	if err := tmp.Close(); err != nil {
		return err
	}
	if err := os.Chmod(tmpPath, 0o644); err != nil {
		return err
	}
	return os.Rename(tmpPath, path)
}

func hashContent(content []byte) string {
	sum := sha256.Sum256(content)
	return hex.EncodeToString(sum[:])
}

// fileSyncOutcome names what syncManagedFile actually did to one file.
type fileSyncOutcome int

const (
	syncUnchanged fileSyncOutcome = iota
	syncInstalled
	syncAutoRefreshed
	syncForceRefreshed
	syncAmbiguous
)

// syncManagedFile writes content to target if needed, using manifest to
// distinguish three cases (dec-20260719-25712417):
//   - disk already matches the current shipped content: nothing to do.
//   - disk matches the hash pysar itself shipped last time (manifest[key]),
//     but the shipped content has since changed: safe to auto-refresh --
//     the operator never touched this file since pysar wrote it.
//   - disk matches neither: ambiguous. It may be a genuine author
//     customization, or a file installed before this manifest existed. The
//     decided fallback policy is to never guess -- leave it untouched unless
//     force is set, same as any other force-only overwrite.
//
// Shared by writeClaudeSkills and writeVoiceTemplates (dec-20260719-3e36577e
// prediction: the template mechanism reuses this exact logic, not a second
// implementation) -- two managed directories, one set of rules.
func syncManagedFile(target string, content []byte, manifest *skillManifest, key string, force bool) (fileSyncOutcome, error) {
	if err := os.MkdirAll(filepath.Dir(target), 0o755); err != nil {
		return syncUnchanged, err
	}
	newHash := hashContent(content)

	existing, readErr := os.ReadFile(target)
	if readErr != nil {
		if !os.IsNotExist(readErr) {
			return syncUnchanged, readErr
		}
		if err := os.WriteFile(target, content, 0o644); err != nil {
			return syncUnchanged, err
		}
		manifest.Files[key] = newHash
		return syncInstalled, nil
	}

	diskHash := hashContent(existing)
	if diskHash == newHash {
		manifest.Files[key] = newHash
		return syncUnchanged, nil
	}

	lastShippedHash, tracked := manifest.Files[key]
	safeToAutoRefresh := tracked && lastShippedHash == diskHash

	switch {
	case safeToAutoRefresh:
		if err := os.WriteFile(target, content, 0o644); err != nil {
			return syncUnchanged, err
		}
		manifest.Files[key] = newHash
		return syncAutoRefreshed, nil
	case force:
		if err := os.WriteFile(target, content, 0o644); err != nil {
			return syncUnchanged, err
		}
		manifest.Files[key] = newHash
		return syncForceRefreshed, nil
	default:
		return syncAmbiguous, nil
	}
}

// syncManagedTree walks every file under embedRoot in assets, syncing each
// into targetDir (preserving the embedded tree's relative structure) via
// syncManagedFile, using one manifest stored at targetDir/skillManifestFile.
func syncManagedTree(assets fs.FS, embedRoot, targetDir string, force bool) (installed, autoRefreshed, refreshed, ambiguous []string, err error) {
	// Not every Kind ships a built-in default yet -- e.g. Style's generic
	// template is deliberately withheld until it comes from a real
	// operator-run /ps-style conversation (dec-20260719-a1ac8959's own
	// generic.md was fabricated content, not dogfooded, and was removed).
	// A missing embedded subtree is a legitimate empty state, not an error.
	if _, statErr := fs.Stat(assets, embedRoot); statErr != nil {
		return nil, nil, nil, nil, nil
	}

	manifestPath := filepath.Join(targetDir, skillManifestFile)
	manifest, err := loadSkillManifest(manifestPath)
	if err != nil {
		return nil, nil, nil, nil, err
	}

	walkErr := fs.WalkDir(assets, embedRoot, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}

		rel, err := filepath.Rel(embedRoot, path)
		if err != nil {
			return err
		}
		target := filepath.Join(targetDir, rel)

		content, err := fs.ReadFile(assets, path)
		if err != nil {
			return err
		}

		outcome, err := syncManagedFile(target, content, &manifest, rel, force)
		if err != nil {
			return err
		}
		switch outcome {
		case syncInstalled:
			installed = append(installed, target)
		case syncAutoRefreshed:
			autoRefreshed = append(autoRefreshed, target)
		case syncForceRefreshed:
			refreshed = append(refreshed, target)
		case syncAmbiguous:
			ambiguous = append(ambiguous, target)
		}
		return nil
	})
	if walkErr != nil {
		return nil, nil, nil, nil, walkErr
	}

	if err := saveSkillManifest(manifestPath, manifest); err != nil {
		return nil, nil, nil, nil, err
	}
	return installed, autoRefreshed, refreshed, ambiguous, nil
}

// writeClaudeSkills installs every embedded ps-* skill into home/.claude/skills/
// and keeps them current on every init via per-file content-hash comparison
// (dec-20260719-25712417) instead of requiring the operator to remember
// --force. installed lists freshly-written paths (first install);
// autoRefreshed lists paths safely auto-updated; refreshed lists paths
// force-updated (--force, including force-overwritten ambiguous files);
// ambiguous lists paths left untouched pending an explicit --force.
func writeClaudeSkills(home string, force bool) (installed, autoRefreshed, refreshed, ambiguous []string, err error) {
	skillsDir := filepath.Join(home, ".claude", "skills")
	return syncManagedTree(claudeSkillAssets, "assets/claude/skills", skillsDir, force)
}

// writeVoiceTemplates installs/refreshes every embedded built-in voice
// template into the operator's cross-project template store
// (onboarding.TemplatesDir, dec-20260719-3e36577e), using the exact same
// syncManagedFile mechanism as writeClaudeSkills -- one manifest format, one
// set of rules, applied to a second managed directory.
func writeVoiceTemplates(home string, force bool) (installed, autoRefreshed, refreshed, ambiguous []string, err error) {
	templatesDir := onboarding.TemplatesDir(home, onboarding.KindVoice)
	return syncManagedTree(templateAssets, "assets/templates/voice", templatesDir, force)
}

// writeStyleTemplates mirrors writeVoiceTemplates for the built-in style
// template (dec-20260719-a1ac8959) -- same mechanism, same manifest format,
// a second Kind's subdirectory under the same cross-project template store.
func writeStyleTemplates(home string, force bool) (installed, autoRefreshed, refreshed, ambiguous []string, err error) {
	templatesDir := onboarding.TemplatesDir(home, onboarding.KindStyle)
	return syncManagedTree(templateAssets, "assets/templates/style", templatesDir, force)
}

// writeProjectManifest writes .pysar/project if it doesn't already exist.
// It never overwrites an existing manifest (dec-20260718-fe926229 invariant),
// but an already-initialized project is a normal, successful outcome for the
// author, not a failure -- the returned bool tells the caller which greeting
// applies. It is never an error on its own; callers must not treat it as one
// or stop the rest of scaffolding because of it (a project initialized
// before a newer pysar version added more scaffold files must still get
// those files backfilled on the next init).
func writeProjectManifest(dir, host string) (alreadyExists bool, err error) {
	pysarDir := filepath.Join(dir, ".pysar")
	if err := os.MkdirAll(pysarDir, 0o755); err != nil {
		return false, err
	}

	manifestPath := filepath.Join(pysarDir, "project")
	if _, statErr := os.Stat(manifestPath); statErr == nil {
		return true, nil
	} else if !os.IsNotExist(statErr) {
		return false, statErr
	}

	data, err := json.MarshalIndent(projectManifest{SchemaVersion: schemaVersion, Host: host}, "", "  ")
	if err != nil {
		return false, err
	}
	if err := os.WriteFile(manifestPath, append(data, '\n'), 0o644); err != nil {
		return false, err
	}
	return false, nil
}

// writeIfAbsent writes content to path if it doesn't already exist. If it
// exists with matching content, nothing happens (matches=true). If it exists
// with different content: without force, it's left untouched (the operator
// has no way to act on knowing this, so callers stay silent -- see
// note-20260719-82ea61c0); with force, it's overwritten (wrote=true), which
// is the actual answer to "how do I get pysar's latest scaffold content."
func writeIfAbsent(path, content string, force bool) (wrote, matches bool, err error) {
	existing, readErr := os.ReadFile(path)
	if readErr == nil {
		if string(existing) == content {
			return false, true, nil
		}
		if !force {
			return false, false, nil
		}
		if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
			return false, false, err
		}
		return true, false, nil
	}
	if !os.IsNotExist(readErr) {
		return false, false, readErr
	}
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		return false, false, err
	}
	return true, false, nil
}
