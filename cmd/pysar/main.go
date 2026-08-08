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

// version defaults to a dev string for `go build`/`go install`; release builds
// override it via -ldflags "-X main.version=..." (dec-20260805-db745b4a).
var version = "0.0.1-dev"

const schemaVersion = 1

// claudeMDTemplate is intentionally minimal and provisional (dec-20260718-fe926229) --
// host-file conventions have not been decided as their own DecisionRecord yet.
const claudeMDTemplate = `# Pysar project

Initialized by ` + "`pysar init --claude`" + `. Host-file conventions here are
provisional and may change as Pysar's Claude Code integration develops.
`

// skillAssets is the host-neutral ps-* skill corpus (dec-20260808-f3001106).
// Host adapters install these bytes into Claude/Cursor skill directories; Codex
// applies an install-time packaging transform ($ps- rewrite + openai.yaml)
// without forking the corpus (dec-20260808-codex-host-v4-ac3eae46). New
// skills need only a directory under assets/skills -- no per-host wiring.
//
//go:embed assets/skills
var skillAssets embed.FS

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

// cursorMCPJSON registers pysar serve for Cursor (haft CL1 path shape:
// .cursor/mcp.json + ${workspaceFolder}; dec-20260808-f3001106).
//
//go:embed assets/cursor/mcp.json
var cursorMCPJSON string

// codexMCPTOML registers pysar serve for Codex CLI/App (haft CL1 path shape:
// project .codex/config.toml; dec-20260808-codex-host-v4-ac3eae46).
//
//go:embed assets/codex/config.toml
var codexMCPTOML string

// templateAssets holds every built-in reusable content template (starting
// with the one "generic" voice default), seeded into the operator's
// cross-project template store (dec-20260719-3e36577e) the same way
// skillAssets seed host skill dirs -- same embed-and-sync shape,
// different target directory, host-agnostic since templates are pysar's own
// data, not a host-specific config surface.
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
  pysar init --cursor ./my-piece
  pysar init --codex ./my-piece
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
			host, err := resolveHost(claude, cursor, codex)
			if err != nil {
				return err
			}
			return host.Scaffold(dir, force)
		},
	}

	cmd.Flags().BoolVar(&claude, "claude", false, "scaffold for Claude Code (default)")
	cmd.Flags().BoolVar(&cursor, "cursor", false, "scaffold for Cursor")
	cmd.Flags().BoolVar(&codex, "codex", false, "scaffold for Codex CLI / App")
	cmd.Flags().BoolVar(&force, "force", false, "refresh host project files, agentic skills, and permission/MCP config to the current shipped version, even if already installed -- never touches .pysar/ project data")
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

// skillManifestFile records, per skill file, the hash of the shipped content
// that pysar itself last wrote there. It lives inside the host's skills
// directory alongside the skills it tracks -- system-owned state for
// system-owned files, never read or written by an author.
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
// Shared by writeHostSkills and writeVoiceTemplates (dec-20260719-3e36577e
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

// writeVoiceTemplates installs/refreshes every embedded built-in voice
// template into the operator's cross-project template store
// (onboarding.TemplatesDir, dec-20260719-3e36577e), using the exact same
// syncManagedFile mechanism as writeHostSkills -- one manifest format, one
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
