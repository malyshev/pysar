package mcpserver

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func validSEOArgs(piecePath, revisedMD string) map[string]interface{} {
	return map[string]interface{}{
		"piece_path":       piecePath,
		"revised_md":       revisedMD,
		"checks":           []string{"[citation] resolved [^retry-budget] to an inline link"},
		"tags":             []string{"reliability"},
		"meta_title":       "A concrete title for search",
		"meta_description": "A concrete, specific description written for a reader who hasn't opened the piece yet.",
		"url_slug":         "a-concrete-title",
	}
}

func TestSaveSEOBundleEndToEndAfterSharpen(t *testing.T) {
	dir := t.TempDir()
	s := New("pysar", "test", dir, t.TempDir(), nil, &bytes.Buffer{})
	piecePath := createTestPiece(t, s)

	draftResp := runLines(t, s, toolCallLine(t, 2, "save_draft_bundle", map[string]interface{}{
		"piece_path": piecePath,
		"draft_md":   "# Title\n\n*Subtitle.*\n\nAn opening claim, no citation yet.\n",
	}))
	if toolResultMap(t, draftResp[1])["isError"] == true {
		t.Fatalf("save_draft_bundle error: %v", draftResp[1])
	}
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

	revisedMD := "# Sharpened Title\n\n*Subtitle.*\n\nA claim backed by [research](https://example.com/retry-budget).\n"
	seoResp := runLines(t, s, toolCallLine(t, 5, "save_seo_bundle", validSEOArgs(piecePath, revisedMD)))
	seoResult := toolResultMap(t, seoResp[4])
	if seoResult["isError"] == true {
		t.Fatalf("save_seo_bundle error: %v", seoResult)
	}
	text := toolText(t, seoResult)
	if !strings.Contains(text, "seo done") {
		t.Fatalf("expected the author-facing outcome message, got %s", text)
	}

	pieceDir := filepath.Join(dir, filepath.FromSlash(piecePath))
	for _, f := range []string{"sharpen.md", "seo.md", "seo-changelog.md", "seo-checklist.md"} {
		if _, err := os.Stat(filepath.Join(pieceDir, f)); err != nil {
			t.Fatalf("missing %s: %v", f, err)
		}
	}

	gotSharpen, _ := os.ReadFile(filepath.Join(pieceDir, "sharpen.md"))
	if !strings.Contains(string(gotSharpen), "[^retry-budget]") {
		t.Fatalf("sharpen.md must remain untouched by seo, got:\n%s", gotSharpen)
	}

	gotSEO, _ := os.ReadFile(filepath.Join(pieceDir, "seo.md"))
	if strings.Contains(string(gotSEO), "[^retry-budget]") {
		t.Fatalf("seo.md must not contain unresolved citation markers, got:\n%s", gotSEO)
	}
	if !strings.Contains(string(gotSEO), "https://example.com/retry-budget") {
		t.Fatalf("seo.md missing the resolved link, got:\n%s", gotSEO)
	}

	checklist, _ := os.ReadFile(filepath.Join(pieceDir, "seo-checklist.md"))
	if !strings.Contains(string(checklist), "a-concrete-title") {
		t.Fatalf("seo-checklist.md missing the url slug, got:\n%s", checklist)
	}

	runLog, _ := os.ReadFile(filepath.Join(pieceDir, "run-log.jsonl"))
	if !strings.Contains(string(runLog), `"pass":"discoverability"`) {
		t.Fatalf("run-log missing discoverability entry: %s", runLog)
	}
}

func TestSaveSEOBundleRejectsUnresolvedCitation(t *testing.T) {
	dir := t.TempDir()
	s := New("pysar", "test", dir, t.TempDir(), nil, &bytes.Buffer{})
	piecePath := createTestPiece(t, s)

	runLines(t, s, toolCallLine(t, 2, "save_draft_bundle", map[string]interface{}{
		"piece_path": piecePath,
		"draft_md":   "# Title\n\n*Subtitle.*\n\nbody\n",
	}))

	resp := runLines(t, s, toolCallLine(t, 3, "save_seo_bundle", validSEOArgs(piecePath, "# Title\n\n*Subtitle.*\n\nStill cites [^unresolved].\n")))
	result := toolResultMap(t, resp[2])
	if result["isError"] != true {
		t.Fatalf("expected error for an unresolved citation marker, got %v", result)
	}
}

func TestSaveSEOBundleRejectsFabricatedLink(t *testing.T) {
	dir := t.TempDir()
	s := New("pysar", "test", dir, t.TempDir(), nil, &bytes.Buffer{})
	piecePath := createTestPiece(t, s)

	runLines(t, s, toolCallLine(t, 2, "save_draft_bundle", map[string]interface{}{
		"piece_path": piecePath,
		"draft_md":   "# Title\n\n*Subtitle.*\n\nbody\n",
	}))

	resp := runLines(t, s, toolCallLine(t, 3, "save_seo_bundle", validSEOArgs(piecePath, "# Title\n\n*Subtitle.*\n\nA claim citing [an invented source](https://not-real.example.com).\n")))
	result := toolResultMap(t, resp[2])
	if result["isError"] != true {
		t.Fatalf("expected error for a link with no matching recorded source, got %v", result)
	}
}

func TestSaveSEOBundleRejectsWithoutDraft(t *testing.T) {
	dir := t.TempDir()
	s := New("pysar", "test", dir, t.TempDir(), nil, &bytes.Buffer{})
	piecePath := createTestPiece(t, s)

	resp := runLines(t, s, toolCallLine(t, 2, "save_seo_bundle", validSEOArgs(piecePath, "# Title\n\n*Subtitle.*\n\nbody\n")))
	result := toolResultMap(t, resp[1])
	if result["isError"] != true {
		t.Fatalf("expected error when no draft.md exists yet, got %v", result)
	}
}

func TestSaveSEOBundleRejectsAfterHumanizeAlreadyRan(t *testing.T) {
	// Ordering invariant (dec-20260804-e3234e50): seo must run BEFORE
	// humanize, never after. Enforced by discoverabilityPass's own
	// Precondition, exercised here at the MCP surface.
	dir := t.TempDir()
	s := New("pysar", "test", dir, t.TempDir(), nil, &bytes.Buffer{})
	piecePath := createTestPiece(t, s)

	runLines(t, s, toolCallLine(t, 2, "save_draft_bundle", map[string]interface{}{
		"piece_path": piecePath,
		"draft_md":   "# Title\n\n*Subtitle.*\n\nbody\n",
	}))
	runLines(t, s, toolCallLine(t, 3, "save_humanize_bundle", map[string]interface{}{
		"piece_path": piecePath,
		"revised_md": "# Humanized Title\n\n*Subtitle.*\n\nhumanized body\n",
		"checks":     []string{"[hedge-stack] dropped a redundant qualifier"},
	}))

	resp := runLines(t, s, toolCallLine(t, 4, "save_seo_bundle", validSEOArgs(piecePath, "# Title\n\n*Subtitle.*\n\nbody\n")))
	result := toolResultMap(t, resp[3])
	if result["isError"] != true {
		t.Fatalf("expected error when humanize.md already exists, got %v", result)
	}
}

func TestSaveSEOBundleRejectsMissingChecklistFields(t *testing.T) {
	dir := t.TempDir()
	s := New("pysar", "test", dir, t.TempDir(), nil, &bytes.Buffer{})
	piecePath := createTestPiece(t, s)

	runLines(t, s, toolCallLine(t, 2, "save_draft_bundle", map[string]interface{}{
		"piece_path": piecePath,
		"draft_md":   "# Title\n\n*Subtitle.*\n\nbody\n",
	}))

	args := map[string]interface{}{
		"piece_path": piecePath,
		"revised_md": "# Title\n\n*Subtitle.*\n\nbody\n",
		"checks":     []string{"no changes needed"},
	}
	resp := runLines(t, s, toolCallLine(t, 3, "save_seo_bundle", args))
	result := toolResultMap(t, resp[2])
	if result["isError"] != true {
		t.Fatalf("expected error for missing tags/meta_title/meta_description/url_slug, got %v", result)
	}
}
