package mcpserver

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestSaveExportBundlePicksHumanizeWhenEveryStageRan(t *testing.T) {
	dir := t.TempDir()
	s := New("pysar", "test", dir, t.TempDir(), nil, &bytes.Buffer{})
	piecePath := createTestPiece(t, s)

	runLines(t, s, toolCallLine(t, 2, "save_draft_bundle", map[string]interface{}{
		"piece_path": piecePath,
		"draft_md":   "# Draft Title\n\n*Subtitle.*\n\ndraft body\n",
	}))
	runLines(t, s, toolCallLine(t, 3, "save_staff_edit_bundle", map[string]interface{}{
		"piece_path": piecePath,
		"revised_md": "# Staff-Edited Title\n\n*Subtitle.*\n\nstaff-edited body\n",
		"checks":     []string{"[stakes] tightened the opener"},
	}))
	runLines(t, s, toolCallLine(t, 4, "save_sharpen_bundle", map[string]interface{}{
		"piece_path": piecePath,
		"revised_md": "# Sharpened Title\n\n*Subtitle.*\n\nsharpened body\n",
		"checks":     []string{"[opener] sharpened further"},
	}))
	runLines(t, s, toolCallLine(t, 5, "save_humanize_bundle", map[string]interface{}{
		"piece_path": piecePath,
		"revised_md": "# Humanized Title\n\n*Subtitle.*\n\nhumanized body\n",
		"checks":     []string{"[hedge-stack] dropped a redundant qualifier"},
	}))

	exportResp := runLines(t, s, toolCallLine(t, 6, "export_piece_to_root", map[string]interface{}{
		"piece_path": piecePath,
	}))
	exportResult := toolResultMap(t, exportResp[5])
	if exportResult["isError"] == true {
		t.Fatalf("export_piece_to_root error: %v", exportResult)
	}
	text := toolText(t, exportResult)
	if !strings.Contains(text, "humanize.md") {
		t.Fatalf("expected export to report humanize.md as the source, got %s", text)
	}

	slug := filepath.Base(filepath.Join(dir, filepath.FromSlash(piecePath)))
	got, err := os.ReadFile(filepath.Join(dir, slug+".md"))
	if err != nil {
		t.Fatalf("expected %s.md at project root: %v", slug, err)
	}
	if !strings.Contains(string(got), "Humanized Title") {
		t.Fatalf("expected the exported file to contain humanize.md's content, got:\n%s", got)
	}

	pieceDir := filepath.Join(dir, filepath.FromSlash(piecePath))
	runLog, _ := os.ReadFile(filepath.Join(pieceDir, "run-log.jsonl"))
	if !strings.Contains(string(runLog), `"pass":"export"`) {
		t.Fatalf("run-log missing export entry: %s", runLog)
	}
}

func TestSaveExportBundlePicksSEOWhenHumanizeHasNotRun(t *testing.T) {
	// Regression: export must not silently fall back to sharpen.md's
	// unpackaged prose when seo.md exists but humanize hasn't run yet --
	// draft.RevisionPriority is the single source of truth this and
	// tools_humanize.go's own revised_from tracking both read from.
	dir := t.TempDir()
	s := New("pysar", "test", dir, t.TempDir(), nil, &bytes.Buffer{})
	piecePath := createTestPiece(t, s)

	runLines(t, s, toolCallLine(t, 2, "save_draft_bundle", map[string]interface{}{
		"piece_path": piecePath,
		"draft_md":   "# Title\n\n*Subtitle.*\n\nAn opening claim, no citation yet.\n",
	}))
	runLines(t, s, toolCallLine(t, 3, "save_research_bundle", map[string]interface{}{
		"piece_path":  piecePath,
		"expert_lens": "distributed-systems / platform engineering",
		"sources":     []map[string]interface{}{validResearchSource("retry-budget")},
	}))
	runLines(t, s, toolCallLine(t, 4, "save_sharpen_bundle", map[string]interface{}{
		"piece_path": piecePath,
		"revised_md": "# Sharpened Title\n\n*Subtitle.*\n\nA claim backed by research[^retry-budget].\n",
		"checks":     []string{"[opener] tightened the hook"},
	}))
	runLines(t, s, toolCallLine(t, 5, "save_seo_bundle", validSEOArgs(piecePath,
		"# SEO Packaged Title\n\n*Subtitle.*\n\nA claim backed by [research](https://example.com/retry-budget).\n")))

	exportResp := runLines(t, s, toolCallLine(t, 6, "export_piece_to_root", map[string]interface{}{
		"piece_path": piecePath,
	}))
	exportResult := toolResultMap(t, exportResp[5])
	if exportResult["isError"] == true {
		t.Fatalf("export_piece_to_root error: %v", exportResult)
	}
	text := toolText(t, exportResult)
	if !strings.Contains(text, "seo.md") {
		t.Fatalf("expected export to report seo.md as the source, got %s", text)
	}

	slug := filepath.Base(filepath.Join(dir, filepath.FromSlash(piecePath)))
	got, err := os.ReadFile(filepath.Join(dir, slug+".md"))
	if err != nil {
		t.Fatalf("expected %s.md at project root: %v", slug, err)
	}
	if !strings.Contains(string(got), "SEO Packaged Title") {
		t.Fatalf("expected the exported file to contain seo.md's content, not an earlier stage's, got:\n%s", got)
	}
}

func TestSaveExportBundleFallsBackToDraftAlone(t *testing.T) {
	dir := t.TempDir()
	s := New("pysar", "test", dir, t.TempDir(), nil, &bytes.Buffer{})
	piecePath := createTestPiece(t, s)

	runLines(t, s, toolCallLine(t, 2, "save_draft_bundle", map[string]interface{}{
		"piece_path": piecePath,
		"draft_md":   "# Only A Draft\n\n*Subtitle.*\n\ndraft body\n",
	}))

	exportResp := runLines(t, s, toolCallLine(t, 3, "export_piece_to_root", map[string]interface{}{
		"piece_path": piecePath,
	}))
	exportResult := toolResultMap(t, exportResp[2])
	if exportResult["isError"] == true {
		t.Fatalf("export_piece_to_root error: %v", exportResult)
	}
	text := toolText(t, exportResult)
	if !strings.Contains(text, "draft.md") {
		t.Fatalf("expected export to fall back to draft.md, got %s", text)
	}

	slug := filepath.Base(filepath.Join(dir, filepath.FromSlash(piecePath)))
	got, err := os.ReadFile(filepath.Join(dir, slug+".md"))
	if err != nil {
		t.Fatalf("expected %s.md at project root: %v", slug, err)
	}
	if !strings.Contains(string(got), "Only A Draft") {
		t.Fatalf("expected the exported file to contain draft.md's content, got:\n%s", got)
	}
}

func TestSaveExportBundleRejectsWithoutDraft(t *testing.T) {
	dir := t.TempDir()
	s := New("pysar", "test", dir, t.TempDir(), nil, &bytes.Buffer{})
	piecePath := createTestPiece(t, s)

	resp := runLines(t, s, toolCallLine(t, 2, "export_piece_to_root", map[string]interface{}{
		"piece_path": piecePath,
	}))
	result := toolResultMap(t, resp[1])
	if result["isError"] != true {
		t.Fatalf("expected error when no draft.md exists yet, got %v", result)
	}
}

func TestSaveExportBundleRerunOverwritesRootFile(t *testing.T) {
	dir := t.TempDir()
	s := New("pysar", "test", dir, t.TempDir(), nil, &bytes.Buffer{})
	piecePath := createTestPiece(t, s)

	runLines(t, s, toolCallLine(t, 2, "save_draft_bundle", map[string]interface{}{
		"piece_path": piecePath,
		"draft_md":   "# First Draft\n\n*Subtitle.*\n\nfirst body\n",
	}))
	runLines(t, s, toolCallLine(t, 3, "export_piece_to_root", map[string]interface{}{
		"piece_path": piecePath,
	}))
	runLines(t, s, toolCallLine(t, 4, "save_draft_bundle", map[string]interface{}{
		"piece_path": piecePath,
		"draft_md":   "# Second Draft\n\n*Subtitle.*\n\nsecond body\n",
	}))
	runLines(t, s, toolCallLine(t, 5, "export_piece_to_root", map[string]interface{}{
		"piece_path": piecePath,
	}))

	slug := filepath.Base(filepath.Join(dir, filepath.FromSlash(piecePath)))
	got, err := os.ReadFile(filepath.Join(dir, slug+".md"))
	if err != nil {
		t.Fatalf("expected %s.md at project root: %v", slug, err)
	}
	if strings.Contains(string(got), "First Draft") || !strings.Contains(string(got), "Second Draft") {
		t.Fatalf("expected re-export to overwrite the root file wholesale, got:\n%s", got)
	}
}

func TestSaveExportBundleUsesProjectExportDir(t *testing.T) {
	dir := t.TempDir()
	if err := os.MkdirAll(filepath.Join(dir, ".pysar"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, ".pysar", "project"), []byte(`{"schema_version":1,"host":"claude","export_dir":"published"}`+"\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	s := New("pysar", "test", dir, t.TempDir(), nil, &bytes.Buffer{})
	piecePath := createTestPiece(t, s)
	runLines(t, s, toolCallLine(t, 2, "save_draft_bundle", map[string]interface{}{
		"piece_path": piecePath,
		"draft_md":   "# Configured Dest\n\n*Subtitle.*\n\nbody\n",
	}))
	exportResp := runLines(t, s, toolCallLine(t, 3, "export_piece_to_root", map[string]interface{}{
		"piece_path": piecePath,
	}))
	exportResult := toolResultMap(t, exportResp[2])
	if exportResult["isError"] == true {
		t.Fatalf("export error: %v", exportResult)
	}
	text := toolText(t, exportResult)
	slug := filepath.Base(filepath.Join(dir, filepath.FromSlash(piecePath)))
	wantRel := filepath.Join("published", slug+".md")
	if !strings.Contains(text, wantRel) && !strings.Contains(text, filepath.Join(dir, wantRel)) {
		t.Fatalf("expected result to cite published destination, got %s", text)
	}
	got, err := os.ReadFile(filepath.Join(dir, "published", slug+".md"))
	if err != nil {
		t.Fatalf("expected file under published/: %v", err)
	}
	if !strings.Contains(string(got), "Configured Dest") {
		t.Fatalf("unexpected content: %s", got)
	}
	if _, err := os.Stat(filepath.Join(dir, slug+".md")); !os.IsNotExist(err) {
		t.Fatalf("should not also write project-root export when export_dir is set")
	}
}

func TestSaveExportBundleOverrideBeatsProjectDefault(t *testing.T) {
	dir := t.TempDir()
	if err := os.MkdirAll(filepath.Join(dir, ".pysar"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, ".pysar", "project"), []byte(`{"schema_version":1,"host":"claude","export_dir":"published"}`+"\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	s := New("pysar", "test", dir, t.TempDir(), nil, &bytes.Buffer{})
	piecePath := createTestPiece(t, s)
	runLines(t, s, toolCallLine(t, 2, "save_draft_bundle", map[string]interface{}{
		"piece_path": piecePath,
		"draft_md":   "# Override Dest\n\n*Subtitle.*\n\nbody\n",
	}))
	exportResp := runLines(t, s, toolCallLine(t, 3, "export_piece_to_root", map[string]interface{}{
		"piece_path": piecePath,
		"export_dir": "outbox",
	}))
	exportResult := toolResultMap(t, exportResp[2])
	if exportResult["isError"] == true {
		t.Fatalf("export error: %v", exportResult)
	}
	slug := filepath.Base(filepath.Join(dir, filepath.FromSlash(piecePath)))
	got, err := os.ReadFile(filepath.Join(dir, "outbox", slug+".md"))
	if err != nil {
		t.Fatalf("expected override dest: %v", err)
	}
	if !strings.Contains(string(got), "Override Dest") {
		t.Fatalf("unexpected content: %s", got)
	}
	if _, err := os.Stat(filepath.Join(dir, "published", slug+".md")); !os.IsNotExist(err) {
		t.Fatal("override must not also write the project default dir")
	}

	// Next call without override returns to project default.
	runLines(t, s, toolCallLine(t, 4, "save_draft_bundle", map[string]interface{}{
		"piece_path": piecePath,
		"draft_md":   "# Back To Default\n\n*Subtitle.*\n\nbody\n",
	}))
	runLines(t, s, toolCallLine(t, 5, "export_piece_to_root", map[string]interface{}{
		"piece_path": piecePath,
	}))
	got, err = os.ReadFile(filepath.Join(dir, "published", slug+".md"))
	if err != nil {
		t.Fatalf("expected project default after omit: %v", err)
	}
	if !strings.Contains(string(got), "Back To Default") {
		t.Fatalf("unexpected content: %s", got)
	}
}

func TestSaveExportBundleRejectsEscapingExportDir(t *testing.T) {
	dir := t.TempDir()
	s := New("pysar", "test", dir, t.TempDir(), nil, &bytes.Buffer{})
	piecePath := createTestPiece(t, s)
	runLines(t, s, toolCallLine(t, 2, "save_draft_bundle", map[string]interface{}{
		"piece_path": piecePath,
		"draft_md":   "# Escape\n\n*Subtitle.*\n\nbody\n",
	}))
	resp := runLines(t, s, toolCallLine(t, 3, "export_piece_to_root", map[string]interface{}{
		"piece_path": piecePath,
		"export_dir": "..",
	}))
	result := toolResultMap(t, resp[2])
	if result["isError"] != true {
		t.Fatalf("expected escape error, got %v", result)
	}
}
