package main

import (
	"bytes"
	"encoding/json"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// captureStdout redirects os.Stdout for the duration of fn and returns
// whatever was printed, so tests can assert on message noise (or silence).
func captureStdout(t *testing.T, fn func()) string {
	t.Helper()
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatalf("pipe: %v", err)
	}
	orig := os.Stdout
	os.Stdout = w

	fn()

	os.Stdout = orig
	if err := w.Close(); err != nil {
		t.Fatalf("close pipe writer: %v", err)
	}
	var buf bytes.Buffer
	if _, err := io.Copy(&buf, r); err != nil {
		t.Fatalf("read captured output: %v", err)
	}
	return buf.String()
}

// withFakeHome points homeDir at a temp directory for the duration of the
// test, so no test ever writes into a real developer's actual
// ~/.claude/skills/. Returns the fake home path.
func withFakeHome(t *testing.T) string {
	t.Helper()
	fake := t.TempDir()
	orig := homeDir
	homeDir = func() (string, error) { return fake, nil }
	t.Cleanup(func() { homeDir = orig })
	return fake
}

func TestInitDefaultWritesManifestOnly(t *testing.T) {
	withFakeHome(t)
	dir := t.TempDir()
	cmd := newRootCmd()
	cmd.SetArgs([]string{"init", dir})
	if err := cmd.Execute(); err != nil {
		t.Fatalf("init failed: %v", err)
	}

	entries, err := os.ReadDir(dir)
	if err != nil {
		t.Fatalf("read dir: %v", err)
	}
	names := map[string]bool{}
	for _, e := range entries {
		names[e.Name()] = true
	}
	if len(names) != 4 || !names[".pysar"] || !names["CLAUDE.md"] || !names[".claude"] || !names[".mcp.json"] {
		t.Fatalf("unexpected entries at project root: %v", names)
	}

	manifestPath := filepath.Join(dir, ".pysar", "project")
	data, err := os.ReadFile(manifestPath)
	if err != nil {
		t.Fatalf("read manifest: %v", err)
	}
	var m projectManifest
	if err := json.Unmarshal(data, &m); err != nil {
		t.Fatalf("parse manifest: %v", err)
	}
	if m.SchemaVersion != schemaVersion || m.Host != "claude" {
		t.Fatalf("unexpected manifest: %+v", m)
	}
}

func TestInitScaffoldsMCPConfig(t *testing.T) {
	withFakeHome(t)
	dir := t.TempDir()
	cmd := newRootCmd()
	cmd.SetArgs([]string{"init", dir})
	if err := cmd.Execute(); err != nil {
		t.Fatalf("init failed: %v", err)
	}

	mcpPath := filepath.Join(dir, ".mcp.json")
	content, err := os.ReadFile(mcpPath)
	if err != nil {
		t.Fatalf("expected .mcp.json scaffolded: %v", err)
	}
	if !strings.Contains(string(content), `"command": "pysar"`) || !strings.Contains(string(content), `"serve"`) {
		t.Fatalf(".mcp.json missing expected pysar serve entry, got: %s", content)
	}
}

func TestInitLeavesExistingMCPConfigUntouched(t *testing.T) {
	withFakeHome(t)
	dir := t.TempDir()
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatalf("seed dir: %v", err)
	}
	const sentinel = `{"mcpServers":{"author-customized":{}}}`
	mcpPath := filepath.Join(dir, ".mcp.json")
	if err := os.WriteFile(mcpPath, []byte(sentinel), 0o644); err != nil {
		t.Fatalf("seed .mcp.json: %v", err)
	}

	cmd := newRootCmd()
	cmd.SetArgs([]string{"init", dir})
	if err := cmd.Execute(); err != nil {
		t.Fatalf("init failed: %v", err)
	}

	got, err := os.ReadFile(mcpPath)
	if err != nil {
		t.Fatalf("read .mcp.json: %v", err)
	}
	if string(got) != sentinel {
		t.Fatalf(".mcp.json was overwritten: got %q, want %q", got, sentinel)
	}
}

func TestInitClaudeFlagMatchesDefault(t *testing.T) {
	withFakeHome(t)
	dir := t.TempDir()
	cmd := newRootCmd()
	cmd.SetArgs([]string{"init", "--claude", dir})
	if err := cmd.Execute(); err != nil {
		t.Fatalf("init --claude failed: %v", err)
	}
	if _, err := os.Stat(filepath.Join(dir, ".pysar", "project")); err != nil {
		t.Fatalf("expected manifest: %v", err)
	}
	if _, err := os.Stat(filepath.Join(dir, "CLAUDE.md")); err != nil {
		t.Fatalf("expected CLAUDE.md: %v", err)
	}
}

// TestInitOnAlreadySetUpProjectSucceeds locks in a correction: re-running
// init on an already-initialized project is a normal, successful outcome
// for the author (exit 0, friendly message), not an error -- while the
// manifest itself is still genuinely never overwritten.
func TestInitOnAlreadySetUpProjectSucceeds(t *testing.T) {
	withFakeHome(t)
	dir := t.TempDir()

	first := newRootCmd()
	first.SetArgs([]string{"init", dir})
	if err := first.Execute(); err != nil {
		t.Fatalf("first init failed: %v", err)
	}
	original, err := os.ReadFile(filepath.Join(dir, ".pysar", "project"))
	if err != nil {
		t.Fatalf("read manifest after first init: %v", err)
	}

	output := captureStdout(t, func() {
		second := newRootCmd()
		second.SetArgs([]string{"init", dir})
		if err := second.Execute(); err != nil {
			t.Fatalf("expected second init on an already set up project to succeed, got: %v", err)
		}
	})

	if strings.Contains(output, "Error") || strings.Contains(output, "refusing") {
		t.Fatalf("expected a success-shaped message, not an error, got: %q", output)
	}
	if !strings.Contains(output, "already set up") {
		t.Fatalf("expected a plain confirmation the project is already set up, got: %q", output)
	}

	got, err := os.ReadFile(filepath.Join(dir, ".pysar", "project"))
	if err != nil || string(got) != string(original) {
		t.Fatalf("manifest was touched by re-init: got %q, want %q (err=%v)", got, original, err)
	}
}

// TestInitDefaultArgNeverPrintsRawDot is the exact regression the operator
// hit: bare `pysar init` (no argument) resolves dir to ".", and printing
// that literally reads as broken output ("is already set up") rather than
// a real sentence. Runs the actual default-arg code path via Chdir, not a
// simulated dir="." arg, so it reproduces the real bug precisely.
func TestInitDefaultArgNeverPrintsRawDot(t *testing.T) {
	withFakeHome(t)
	dir := t.TempDir()

	origWD, err := os.Getwd()
	if err != nil {
		t.Fatalf("getwd: %v", err)
	}
	if err := os.Chdir(dir); err != nil {
		t.Fatalf("chdir: %v", err)
	}
	t.Cleanup(func() { _ = os.Chdir(origWD) })

	runBareInit := func() string {
		return captureStdout(t, func() {
			cmd := newRootCmd()
			cmd.SetArgs([]string{"init"})
			if err := cmd.Execute(); err != nil {
				t.Fatalf("init failed: %v", err)
			}
		})
	}

	first := runBareInit()
	for _, word := range strings.Fields(first) {
		if word == "." || word == "./" {
			t.Fatalf("first-run output contains a raw '.' token, reads as broken: %q", first)
		}
	}

	second := runBareInit()
	for _, word := range strings.Fields(second) {
		if word == "." || word == "./" {
			t.Fatalf("re-init output contains a raw '.' token, reads as broken: %q", second)
		}
	}
	if !strings.Contains(second, "already set up") {
		t.Fatalf("expected a plain already-set-up confirmation, got: %q", second)
	}
}

// TestInitForceRefreshesStaleGlobalSkill is the direct fix for the operator's
// repeated real-world problem: a global skill installed once, then never
// updated when pysar's shipped skill content changes, silently leaving an
// old version in place indefinitely. --force is the actual answer.
func TestInitForceRefreshesStaleGlobalSkill(t *testing.T) {
	fakeHome := withFakeHome(t)
	skillPath := filepath.Join(fakeHome, ".claude", "skills", "ps-voice", "SKILL.md")
	if err := os.MkdirAll(filepath.Dir(skillPath), 0o755); err != nil {
		t.Fatalf("seed skill dir: %v", err)
	}
	const stale = "old pre-MCP skill content\n"
	if err := os.WriteFile(skillPath, []byte(stale), 0o644); err != nil {
		t.Fatalf("seed stale SKILL.md: %v", err)
	}

	dir := t.TempDir()

	// Without --force: stays stale (already-established, correct behavior).
	withoutForce := newRootCmd()
	withoutForce.SetArgs([]string{"init", dir})
	if err := withoutForce.Execute(); err != nil {
		t.Fatalf("init without --force failed: %v", err)
	}
	got, err := os.ReadFile(skillPath)
	if err != nil || string(got) != stale {
		t.Fatalf("expected skill still stale without --force: got %q, err %v", got, err)
	}

	// With --force: actually refreshes to the current shipped content.
	withForce := newRootCmd()
	withForce.SetArgs([]string{"init", "--force", dir})
	if err := withForce.Execute(); err != nil {
		t.Fatalf("init --force failed: %v", err)
	}
	got, err = os.ReadFile(skillPath)
	if err != nil {
		t.Fatalf("read skill after --force: %v", err)
	}
	if string(got) == stale {
		t.Fatal("expected --force to refresh the stale skill, but it is unchanged")
	}
	if !strings.Contains(string(got), "name: ps-voice") {
		t.Fatalf("expected refreshed content to be the real shipped skill, got: %q", got)
	}
}

// TestInitForceNeverTouchesProjectManifest guards the one hard boundary:
// --force refreshes tooling/scaffold defaults, never author project data.
func TestInitForceNeverTouchesProjectManifest(t *testing.T) {
	withFakeHome(t)
	dir := t.TempDir()

	first := newRootCmd()
	first.SetArgs([]string{"init", dir})
	if err := first.Execute(); err != nil {
		t.Fatalf("first init failed: %v", err)
	}
	original, err := os.ReadFile(filepath.Join(dir, ".pysar", "project"))
	if err != nil {
		t.Fatalf("read manifest: %v", err)
	}

	forced := newRootCmd()
	forced.SetArgs([]string{"init", "--force", dir})
	if err := forced.Execute(); err != nil {
		t.Fatalf("init --force failed: %v", err)
	}
	got, err := os.ReadFile(filepath.Join(dir, ".pysar", "project"))
	if err != nil || string(got) != string(original) {
		t.Fatalf("--force touched the project manifest: got %q, want %q (err=%v)", got, original, err)
	}
}

// TestInitBackfillsNewScaffoldFilesOnExistingProject covers the real bug the
// "always error" behavior masked: a project initialized before pysar added
// more scaffold files (skills, settings.json, .mcp.json) must still receive
// them on a later init, not be permanently stuck without them.
func TestInitBackfillsNewScaffoldFilesOnExistingProject(t *testing.T) {
	fakeHome := withFakeHome(t)
	dir := t.TempDir()

	// Simulate an "old" pysar-init'd project: only the manifest exists,
	// nothing else pysar now also scaffolds.
	pysarDir := filepath.Join(dir, ".pysar")
	if err := os.MkdirAll(pysarDir, 0o755); err != nil {
		t.Fatalf("seed .pysar dir: %v", err)
	}
	if err := os.WriteFile(filepath.Join(pysarDir, "project"), []byte(`{"schema_version":1,"host":"claude"}`+"\n"), 0o644); err != nil {
		t.Fatalf("seed manifest: %v", err)
	}

	cmd := newRootCmd()
	cmd.SetArgs([]string{"init", dir})
	if err := cmd.Execute(); err != nil {
		t.Fatalf("init on a manifest-only project failed: %v", err)
	}

	for _, rel := range []string{"CLAUDE.md", ".mcp.json", filepath.Join(".claude", "settings.json")} {
		if _, err := os.Stat(filepath.Join(dir, rel)); err != nil {
			t.Fatalf("expected %s to be backfilled, got: %v", rel, err)
		}
	}
	if _, err := os.Stat(filepath.Join(fakeHome, ".claude", "skills", "ps-voice", "SKILL.md")); err != nil {
		t.Fatalf("expected global ps-voice skill to be backfilled, got: %v", err)
	}
}

func TestInitInstallsPsVoiceSkillGlobally(t *testing.T) {
	fakeHome := withFakeHome(t)
	dir := t.TempDir()
	cmd := newRootCmd()
	cmd.SetArgs([]string{"init", dir})
	if err := cmd.Execute(); err != nil {
		t.Fatalf("init failed: %v", err)
	}

	skillPath := filepath.Join(fakeHome, ".claude", "skills", "ps-voice", "SKILL.md")
	content, err := os.ReadFile(skillPath)
	if err != nil {
		t.Fatalf("expected ps-voice skill installed to home dir: %v", err)
	}
	if !strings.Contains(string(content), "name: ps-voice") {
		t.Fatalf("installed skill missing expected frontmatter, got: %q", content[:min(len(content), 100)])
	}
}

func TestInitSecondProjectStaysQuietAboutAlreadyInstalledSkill(t *testing.T) {
	withFakeHome(t)

	dirA := t.TempDir()
	cmdA := newRootCmd()
	cmdA.SetArgs([]string{"init", dirA})
	if err := cmdA.Execute(); err != nil {
		t.Fatalf("first init failed: %v", err)
	}

	dirB := t.TempDir()
	output := captureStdout(t, func() {
		cmdB := newRootCmd()
		cmdB.SetArgs([]string{"init", dirB})
		if err := cmdB.Execute(); err != nil {
			t.Fatalf("second init failed: %v", err)
		}
	})

	if strings.Contains(output, "ps-voice") || strings.Contains(output, "already exists") {
		t.Fatalf("expected no noise about an already-up-to-date global skill on a second project's init, got: %q", output)
	}
}

func TestInitStaysQuietWhenExistingFileAlreadyMatches(t *testing.T) {
	withFakeHome(t)
	dir := t.TempDir()

	claudePath := filepath.Join(dir, "CLAUDE.md")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatalf("seed dir: %v", err)
	}
	if err := os.WriteFile(claudePath, []byte(claudeMDTemplate), 0o644); err != nil {
		t.Fatalf("seed matching CLAUDE.md: %v", err)
	}

	output := captureStdout(t, func() {
		cmd := newRootCmd()
		cmd.SetArgs([]string{"init", dir})
		if err := cmd.Execute(); err != nil {
			t.Fatalf("init failed: %v", err)
		}
	})

	if strings.Contains(output, "CLAUDE.md") {
		t.Fatalf("expected no noise about CLAUDE.md when its content already matches the default, got: %q", output)
	}
}

func TestInitScaffoldsPermissionSettings(t *testing.T) {
	withFakeHome(t)
	dir := t.TempDir()
	cmd := newRootCmd()
	cmd.SetArgs([]string{"init", dir})
	if err := cmd.Execute(); err != nil {
		t.Fatalf("init failed: %v", err)
	}

	settingsPath := filepath.Join(dir, ".claude", "settings.json")
	content, err := os.ReadFile(settingsPath)
	if err != nil {
		t.Fatalf("expected .claude/settings.json scaffolded: %v", err)
	}
	if !strings.Contains(string(content), `mcp__pysar__*`) {
		t.Fatalf("settings.json missing expected mcp__pysar__* allow rule, got: %s", content)
	}
}

func TestInitLeavesExistingSettingsJSONUntouched(t *testing.T) {
	withFakeHome(t)
	dir := t.TempDir()
	claudeDir := filepath.Join(dir, ".claude")
	if err := os.MkdirAll(claudeDir, 0o755); err != nil {
		t.Fatalf("seed .claude dir: %v", err)
	}
	const sentinel = `{"permissions":{"allow":["author-customized"]}}`
	settingsPath := filepath.Join(claudeDir, "settings.json")
	if err := os.WriteFile(settingsPath, []byte(sentinel), 0o644); err != nil {
		t.Fatalf("seed settings.json: %v", err)
	}

	cmd := newRootCmd()
	cmd.SetArgs([]string{"init", dir})
	if err := cmd.Execute(); err != nil {
		t.Fatalf("init failed: %v", err)
	}

	got, err := os.ReadFile(settingsPath)
	if err != nil {
		t.Fatalf("read settings.json: %v", err)
	}
	if string(got) != sentinel {
		t.Fatalf("settings.json was overwritten: got %q, want %q", got, sentinel)
	}
}

func TestInitLeavesExistingGlobalSkillUntouched(t *testing.T) {
	fakeHome := withFakeHome(t)
	skillDir := filepath.Join(fakeHome, ".claude", "skills", "ps-voice")
	if err := os.MkdirAll(skillDir, 0o755); err != nil {
		t.Fatalf("seed skill dir: %v", err)
	}
	const sentinel = "author-customized skill\n"
	skillPath := filepath.Join(skillDir, "SKILL.md")
	if err := os.WriteFile(skillPath, []byte(sentinel), 0o644); err != nil {
		t.Fatalf("seed SKILL.md: %v", err)
	}

	dir := t.TempDir()
	cmd := newRootCmd()
	cmd.SetArgs([]string{"init", dir})
	if err := cmd.Execute(); err != nil {
		t.Fatalf("init failed: %v", err)
	}

	got, err := os.ReadFile(skillPath)
	if err != nil {
		t.Fatalf("read SKILL.md: %v", err)
	}
	if string(got) != sentinel {
		t.Fatalf("SKILL.md was overwritten: got %q, want %q", got, sentinel)
	}
}

// TestInitStaysQuietForProjectFilesThatDiffer locks in a deliberate
// correction: a project-local file (CLAUDE.md) differing from what pysar
// would write is NOT announced, on top of an already-matching file not being
// announced. With no update mechanism for project-local files, "this
// differs" is a fact the author cannot act on -- noise, not information
// (note-20260719-d883dd47). No raw filesystem paths either, as a secondary
// guard even though the messages themselves are gone. This invariant is
// unchanged by dec-20260719-25712417, which governs only the global skill
// tree, not project-local files.
func TestInitStaysQuietForProjectFilesThatDiffer(t *testing.T) {
	withFakeHome(t)
	dir := t.TempDir()
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatalf("seed dir: %v", err)
	}
	if err := os.WriteFile(filepath.Join(dir, "CLAUDE.md"), []byte("customized\n"), 0o644); err != nil {
		t.Fatalf("seed CLAUDE.md: %v", err)
	}

	output := captureStdout(t, func() {
		cmd := newRootCmd()
		cmd.SetArgs([]string{"init", dir})
		if err := cmd.Execute(); err != nil {
			t.Fatalf("init failed: %v", err)
		}
	})

	if strings.Contains(output, "CLAUDE.md") {
		t.Fatalf("expected no noise about a differing project-local file, got: %q", output)
	}
	if strings.Contains(output, dir+string(filepath.Separator)) {
		t.Fatalf("expected no raw filesystem paths in output, got: %q", output)
	}

	got, err := os.ReadFile(filepath.Join(dir, "CLAUDE.md"))
	if err != nil || string(got) != "customized\n" {
		t.Fatalf("CLAUDE.md was touched: content=%q err=%v", got, err)
	}
}

// TestInitNudgesForAmbiguousGlobalSkillWithoutForce covers the decided
// fallback policy for dec-20260719-25712417's one open question: a global
// skill file whose disk content pysar has no hash record for (never
// installed by this pysar, or installed before the manifest existed) is
// ambiguous -- it might be a real author customization. The decided policy
// is to never guess: never auto-overwrite, but nudge (generically, no skill
// name, no raw path) rather than stay silent like the old fully-manual
// behavior, since there's now an actual action (--force) for the operator to
// take.
func TestInitNudgesForAmbiguousGlobalSkillWithoutForce(t *testing.T) {
	fakeHome := withFakeHome(t)
	skillDir := filepath.Join(fakeHome, ".claude", "skills", "ps-voice")
	if err := os.MkdirAll(skillDir, 0o755); err != nil {
		t.Fatalf("seed skill dir: %v", err)
	}
	if err := os.WriteFile(filepath.Join(skillDir, "SKILL.md"), []byte("customized\n"), 0o644); err != nil {
		t.Fatalf("seed SKILL.md: %v", err)
	}

	dir := t.TempDir()
	output := captureStdout(t, func() {
		cmd := newRootCmd()
		cmd.SetArgs([]string{"init", dir})
		if err := cmd.Execute(); err != nil {
			t.Fatalf("init failed: %v", err)
		}
	})

	if !strings.Contains(output, "--force") {
		t.Fatalf("expected a nudge pointing at --force for the ambiguous skill, got: %q", output)
	}
	if strings.Contains(output, "ps-voice") {
		t.Fatalf("expected no skill name leaked into the nudge, got: %q", output)
	}
	if strings.Contains(output, fakeHome) {
		t.Fatalf("expected no raw filesystem paths in output, got: %q", output)
	}

	got, err := os.ReadFile(filepath.Join(skillDir, "SKILL.md"))
	if err != nil || string(got) != "customized\n" {
		t.Fatalf("ambiguous SKILL.md was touched: content=%q err=%v", got, err)
	}
}

// TestInitAutoRefreshesSkillWithoutForceWhenOnlyShippedContentChanged is the
// core of dec-20260719-25712417: once pysar has installed a skill and
// recorded its shipped hash, a later init must refresh that file to the
// current shipped content with NO --force, as long as the file on disk still
// matches the hash pysar itself last wrote there (i.e. the operator never
// touched it -- only the shipped content changed, a normal pysar upgrade).
func TestInitAutoRefreshesSkillWithoutForceWhenOnlyShippedContentChanged(t *testing.T) {
	fakeHome := withFakeHome(t)
	skillDir := filepath.Join(fakeHome, ".claude", "skills", "ps-voice")
	if err := os.MkdirAll(skillDir, 0o755); err != nil {
		t.Fatalf("seed skill dir: %v", err)
	}

	const oldShippedContent = "old shipped ps-voice content, as if from a previous pysar version\n"
	skillPath := filepath.Join(skillDir, "SKILL.md")
	if err := os.WriteFile(skillPath, []byte(oldShippedContent), 0o644); err != nil {
		t.Fatalf("seed old shipped SKILL.md: %v", err)
	}

	manifestPath := filepath.Join(skillDir, "..", skillManifestFile)
	manifest := skillManifest{Files: map[string]string{"ps-voice/SKILL.md": hashContent([]byte(oldShippedContent))}}
	if err := saveSkillManifest(filepath.Clean(manifestPath), manifest); err != nil {
		t.Fatalf("seed manifest: %v", err)
	}

	dir := t.TempDir()
	output := captureStdout(t, func() {
		cmd := newRootCmd()
		cmd.SetArgs([]string{"init", dir})
		if err := cmd.Execute(); err != nil {
			t.Fatalf("init failed: %v", err)
		}
	})

	if strings.Contains(output, "--force") {
		t.Fatalf("expected auto-refresh to need no --force mention, got: %q", output)
	}

	got, err := os.ReadFile(skillPath)
	if err != nil {
		t.Fatalf("read skill after auto-refresh: %v", err)
	}
	if string(got) == oldShippedContent {
		t.Fatal("expected the skill to be auto-refreshed to the current shipped content, but it is unchanged")
	}
	if !strings.Contains(string(got), "name: ps-voice") {
		t.Fatalf("expected refreshed content to be the real shipped skill, got: %q", got)
	}

	updated, err := loadSkillManifest(filepath.Clean(manifestPath))
	if err != nil {
		t.Fatalf("read manifest after auto-refresh: %v", err)
	}
	if updated.Files["ps-voice/SKILL.md"] != hashContent(got) {
		t.Fatalf("expected manifest to record the new shipped hash, got: %+v", updated)
	}
}

// TestInitSelfHealsCorruptSkillManifest guards against a truncated
// .pysar-manifest.json (e.g. left behind by a process killed mid-write)
// permanently bricking every future `pysar init` -- it must fall back to
// treating the manifest as empty rather than hard-failing.
func TestInitSelfHealsCorruptSkillManifest(t *testing.T) {
	fakeHome := withFakeHome(t)
	skillsDir := filepath.Join(fakeHome, ".claude", "skills")
	if err := os.MkdirAll(skillsDir, 0o755); err != nil {
		t.Fatalf("seed skills dir: %v", err)
	}
	manifestPath := filepath.Join(skillsDir, skillManifestFile)
	if err := os.WriteFile(manifestPath, []byte(`{"files": {"ps-voice/SKILL`), 0o644); err != nil {
		t.Fatalf("seed corrupt manifest: %v", err)
	}

	dir := t.TempDir()
	cmd := newRootCmd()
	cmd.SetArgs([]string{"init", dir})
	if err := cmd.Execute(); err != nil {
		t.Fatalf("expected init to self-heal past a corrupt manifest, got: %v", err)
	}

	updated, err := loadSkillManifest(manifestPath)
	if err != nil {
		t.Fatalf("read manifest after self-heal: %v", err)
	}
	if len(updated.Files) == 0 {
		t.Fatal("expected the manifest to be rewritten with real entries after self-heal")
	}
}

// TestInitSeedsBuiltInVoiceTemplate covers dec-20260719-3e36577e: pysar init
// seeds the built-in "generic" voice template into the operator's
// cross-project template store, not just the project itself.
func TestInitSeedsBuiltInVoiceTemplate(t *testing.T) {
	fakeHome := withFakeHome(t)
	dir := t.TempDir()
	cmd := newRootCmd()
	cmd.SetArgs([]string{"init", dir})
	if err := cmd.Execute(); err != nil {
		t.Fatalf("init failed: %v", err)
	}

	templatePath := filepath.Join(fakeHome, ".pysar", "templates", "voice", "generic.md")
	content, err := os.ReadFile(templatePath)
	if err != nil {
		t.Fatalf("expected built-in generic template seeded at %s: %v", templatePath, err)
	}
	if !strings.Contains(string(content), "kind: voice") {
		t.Fatalf("seeded template missing expected frontmatter, got: %q", content[:min(len(content), 100)])
	}

	if _, err := os.Stat(filepath.Join(dir, "templates")); err == nil {
		t.Fatal("expected the template store to live under the home directory, not the project")
	}
}

// TestInitTemplateAutoRefreshesWithoutForce mirrors
// TestInitAutoRefreshesSkillWithoutForceWhenOnlyShippedContentChanged for the
// template mechanism -- both reuse the exact same syncManagedFile logic
// (dec-20260719-3e36577e's own prediction), so both must behave identically.
func TestInitTemplateAutoRefreshesWithoutForce(t *testing.T) {
	fakeHome := withFakeHome(t)
	templateDir := filepath.Join(fakeHome, ".pysar", "templates", "voice")
	if err := os.MkdirAll(templateDir, 0o755); err != nil {
		t.Fatalf("seed template dir: %v", err)
	}

	const oldShippedContent = "old shipped generic template, as if from a previous pysar version\n"
	templatePath := filepath.Join(templateDir, "generic.md")
	if err := os.WriteFile(templatePath, []byte(oldShippedContent), 0o644); err != nil {
		t.Fatalf("seed old shipped template: %v", err)
	}

	manifestPath := filepath.Join(templateDir, skillManifestFile)
	manifest := skillManifest{Files: map[string]string{"generic.md": hashContent([]byte(oldShippedContent))}}
	if err := saveSkillManifest(manifestPath, manifest); err != nil {
		t.Fatalf("seed manifest: %v", err)
	}

	dir := t.TempDir()
	output := captureStdout(t, func() {
		cmd := newRootCmd()
		cmd.SetArgs([]string{"init", dir})
		if err := cmd.Execute(); err != nil {
			t.Fatalf("init failed: %v", err)
		}
	})

	if strings.Contains(output, "--force") {
		t.Fatalf("expected template auto-refresh to need no --force mention, got: %q", output)
	}
	got, err := os.ReadFile(templatePath)
	if err != nil {
		t.Fatalf("read template after auto-refresh: %v", err)
	}
	if string(got) == oldShippedContent {
		t.Fatal("expected the template to be auto-refreshed to the current shipped content, but it is unchanged")
	}
	if !strings.Contains(string(got), "kind: voice") {
		t.Fatalf("expected refreshed content to be the real shipped template, got: %q", got)
	}
}

// TestInitLeavesAmbiguousTemplateUntouchedWithoutForce mirrors the skill
// mechanism's decided fallback policy (dec-20260719-25712417, reused as-is
// by dec-20260719-3e36577e): an untracked or genuinely mismatched template
// is never silently overwritten.
func TestInitLeavesAmbiguousTemplateUntouchedWithoutForce(t *testing.T) {
	fakeHome := withFakeHome(t)
	templateDir := filepath.Join(fakeHome, ".pysar", "templates", "voice")
	if err := os.MkdirAll(templateDir, 0o755); err != nil {
		t.Fatalf("seed template dir: %v", err)
	}
	if err := os.WriteFile(filepath.Join(templateDir, "generic.md"), []byte("author-customized template\n"), 0o644); err != nil {
		t.Fatalf("seed customized template: %v", err)
	}

	dir := t.TempDir()
	cmd := newRootCmd()
	cmd.SetArgs([]string{"init", dir})
	if err := cmd.Execute(); err != nil {
		t.Fatalf("init failed: %v", err)
	}

	got, err := os.ReadFile(filepath.Join(templateDir, "generic.md"))
	if err != nil || string(got) != "author-customized template\n" {
		t.Fatalf("ambiguous template was touched: content=%q err=%v", got, err)
	}
}

func TestInitLeavesExistingClaudeMDUntouched(t *testing.T) {
	withFakeHome(t)
	dir := t.TempDir()
	claudePath := filepath.Join(dir, "CLAUDE.md")
	const sentinel = "author-customized content\n"
	if err := os.WriteFile(claudePath, []byte(sentinel), 0o644); err != nil {
		t.Fatalf("seed CLAUDE.md: %v", err)
	}

	cmd := newRootCmd()
	cmd.SetArgs([]string{"init", dir})
	if err := cmd.Execute(); err != nil {
		t.Fatalf("init failed: %v", err)
	}

	got, err := os.ReadFile(claudePath)
	if err != nil {
		t.Fatalf("read CLAUDE.md: %v", err)
	}
	if string(got) != sentinel {
		t.Fatalf("CLAUDE.md was overwritten: got %q, want %q", got, sentinel)
	}
}

func TestInitCursorNotYetSupported(t *testing.T) {
	withFakeHome(t)
	dir := t.TempDir()
	cmd := newRootCmd()
	cmd.SetArgs([]string{"init", "--cursor", dir})
	if err := cmd.Execute(); err == nil {
		t.Fatal("expected --cursor to fail")
	}
	entries, _ := os.ReadDir(dir)
	if len(entries) != 0 {
		t.Fatalf("expected no files written for --cursor, got %v", entries)
	}
}

func TestInitCodexNotYetSupported(t *testing.T) {
	withFakeHome(t)
	dir := t.TempDir()
	cmd := newRootCmd()
	cmd.SetArgs([]string{"init", "--codex", dir})
	if err := cmd.Execute(); err == nil {
		t.Fatal("expected --codex to fail")
	}
	entries, _ := os.ReadDir(dir)
	if len(entries) != 0 {
		t.Fatalf("expected no files written for --codex, got %v", entries)
	}
}

func TestInitMutuallyExclusiveHostFlags(t *testing.T) {
	withFakeHome(t)
	dir := t.TempDir()
	cmd := newRootCmd()
	cmd.SetArgs([]string{"init", "--claude", "--cursor", dir})
	if err := cmd.Execute(); err == nil {
		t.Fatal("expected mutually exclusive flags to fail")
	}
}

func TestInitInstallsPsStyleSkillGlobally(t *testing.T) {
	fakeHome := withFakeHome(t)
	dir := t.TempDir()
	cmd := newRootCmd()
	cmd.SetArgs([]string{"init", dir})
	if err := cmd.Execute(); err != nil {
		t.Fatalf("init failed: %v", err)
	}

	skillPath := filepath.Join(fakeHome, ".claude", "skills", "ps-style", "SKILL.md")
	content, err := os.ReadFile(skillPath)
	if err != nil {
		t.Fatalf("expected ps-style skill installed to home dir: %v", err)
	}
	if !strings.Contains(string(content), "name: ps-style") {
		t.Fatalf("installed skill missing expected frontmatter, got: %q", content[:min(len(content), 100)])
	}
	if !strings.Contains(string(content), "save_style_profile") {
		t.Fatalf("installed skill missing expected tool reference, got: %q", content[:min(len(content), 200)])
	}
}

// TestInitInstallsPsOnboardSkillGlobally covers dec-20260718-ebca6318's
// implementation: /ps-onboard ships alongside /ps-voice and /ps-style.
func TestInitInstallsPsOnboardSkillGlobally(t *testing.T) {
	fakeHome := withFakeHome(t)
	dir := t.TempDir()
	cmd := newRootCmd()
	cmd.SetArgs([]string{"init", dir})
	if err := cmd.Execute(); err != nil {
		t.Fatalf("init failed: %v", err)
	}

	skillPath := filepath.Join(fakeHome, ".claude", "skills", "ps-onboard", "SKILL.md")
	content, err := os.ReadFile(skillPath)
	if err != nil {
		t.Fatalf("expected ps-onboard skill installed to home dir: %v", err)
	}
	if !strings.Contains(string(content), "name: ps-onboard") {
		t.Fatalf("installed skill missing expected frontmatter, got: %q", content[:min(len(content), 100)])
	}
	if !strings.Contains(string(content), "check_onboarding_status") {
		t.Fatalf("installed skill missing expected tool reference, got: %q", content[:min(len(content), 200)])
	}
}

// TestInitSeedsBuiltInStyleTemplate covers the promoted, genuinely dogfooded
// style default: unlike an earlier, agent-fabricated attempt at this file
// (removed after the operator caught it -- see note-20260719-99c8f3a0), this
// content came from the operator's own real /ps-style conversation and was
// explicitly promoted with their approval, the same process already
// established for voice's own generic default.
func TestInitSeedsBuiltInStyleTemplate(t *testing.T) {
	fakeHome := withFakeHome(t)
	dir := t.TempDir()
	cmd := newRootCmd()
	cmd.SetArgs([]string{"init", dir})
	if err := cmd.Execute(); err != nil {
		t.Fatalf("init failed: %v", err)
	}

	templatePath := filepath.Join(fakeHome, ".pysar", "templates", "style", "generic.md")
	content, err := os.ReadFile(templatePath)
	if err != nil {
		t.Fatalf("expected built-in generic style template seeded at %s: %v", templatePath, err)
	}
	if !strings.Contains(string(content), "kind: style") {
		t.Fatalf("seeded template missing expected frontmatter, got: %q", content[:min(len(content), 100)])
	}
	if !strings.Contains(string(content), "Lead with the main point") {
		t.Fatalf("seeded template missing the real dogfooded rules content, got: %q", content[:min(len(content), 300)])
	}

	// Voice and Style templates must never collide.
	voicePath := filepath.Join(fakeHome, ".pysar", "templates", "voice", "generic.md")
	if _, err := os.Stat(voicePath); err != nil {
		t.Fatalf("expected voice's own generic template to also still be seeded: %v", err)
	}
}

// TestInitMessagingIsKindAgnosticForTemplates locks in the simplified
// messaging design: with two profile kinds now shipping built-in templates
// (voice, style), the init summary reports on "templates" collectively
// rather than naming each kind separately.
func TestInitMessagingIsKindAgnosticForTemplates(t *testing.T) {
	withFakeHome(t)
	dir := t.TempDir()
	output := captureStdout(t, func() {
		cmd := newRootCmd()
		cmd.SetArgs([]string{"init", dir})
		if err := cmd.Execute(); err != nil {
			t.Fatalf("init failed: %v", err)
		}
	})

	if !strings.Contains(output, "installed built-in templates") {
		t.Fatalf("expected a Kind-agnostic templates-installed message, got: %q", output)
	}
	if strings.Contains(output, "voice template") || strings.Contains(output, "style template") {
		t.Fatalf("expected no Kind-specific template wording, got: %q", output)
	}
}
