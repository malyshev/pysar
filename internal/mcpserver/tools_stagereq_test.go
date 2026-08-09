package mcpserver

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestRequirePieceStagesBlocksDraftUntilResearchFull(t *testing.T) {
	dir := t.TempDir()
	s := New("pysar", "test", dir, t.TempDir(), nil, &bytes.Buffer{})
	piecePath := createTestPiece(t, s)

	reqResp := runLines(t, s, toolCallLine(t, 2, "require_piece_stages", map[string]interface{}{
		"piece_path": piecePath,
		"stages":     []string{"research"},
	}))
	if toolResultMap(t, reqResp[1])["isError"] == true {
		t.Fatalf("require_piece_stages error: %v", reqResp[1])
	}

	brief, err := os.ReadFile(filepath.Join(dir, filepath.FromSlash(piecePath), "brief.md"))
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(brief), "required_stages: [research]") {
		t.Fatalf("brief missing required_stages:\n%s", brief)
	}

	blocked := runLines(t, s, toolCallLine(t, 3, "save_draft_bundle", map[string]interface{}{
		"piece_path": piecePath,
		"draft_md":   "# Title\n\n*Subtitle.*\n\nBody without research yet.\n",
	}))
	if toolResultMap(t, blocked[2])["isError"] != true {
		t.Fatalf("expected draft blocked until research_mode full, got %v", blocked[2])
	}
	if msg := toolText(t, toolResultMap(t, blocked[2])); !strings.Contains(msg, "research") {
		t.Fatalf("error should name research gate, got %s", msg)
	}

	runLines(t, s, toolCallLine(t, 4, "save_research_bundle", map[string]interface{}{
		"piece_path":  piecePath,
		"expert_lens": "distributed-systems / platform engineering",
		"sources":     []map[string]interface{}{validResearchSource("docker-docs-init")},
	}))

	ok := runLines(t, s, toolCallLine(t, 5, "save_draft_bundle", map[string]interface{}{
		"piece_path": piecePath,
		"draft_md":   "# Title\n\n*Subtitle.*\n\nClaim grounded after research[^docker-docs-init].\n",
	}))
	if toolResultMap(t, ok[4])["isError"] == true {
		t.Fatalf("draft after research should succeed: %v", ok[4])
	}
}

func TestRequirePieceStagesBlocksHumanizeUntilSEO(t *testing.T) {
	dir := t.TempDir()
	s := New("pysar", "test", dir, t.TempDir(), nil, &bytes.Buffer{})
	piecePath := createTestPiece(t, s)

	runLines(t, s, toolCallLine(t, 2, "require_piece_stages", map[string]interface{}{
		"piece_path": piecePath,
		"stages":     []string{"seo"},
	}))
	runLines(t, s, toolCallLine(t, 3, "save_draft_bundle", map[string]interface{}{
		"piece_path": piecePath,
		"draft_md":   "# Title\n\n*Subtitle.*\n\nAn opening claim.\n",
	}))
	runLines(t, s, toolCallLine(t, 4, "save_research_bundle", map[string]interface{}{
		"piece_path":  piecePath,
		"expert_lens": "distributed-systems / platform engineering",
		"sources":     []map[string]interface{}{validResearchSource("retry-budget")},
	}))
	runLines(t, s, toolCallLine(t, 5, "save_sharpen_bundle", map[string]interface{}{
		"piece_path": piecePath,
		"revised_md": "# Title\n\n*Subtitle.*\n\nA claim[^retry-budget].\n",
		"checks":     []string{"[opener] tightened"},
	}))

	blocked := runLines(t, s, toolCallLine(t, 6, "save_humanize_bundle", map[string]interface{}{
		"piece_path": piecePath,
		"revised_md": "# Title\n\n*Subtitle.*\n\nA claim[^retry-budget].\n",
		"checks":     []string{"[hedge] none needed"},
	}))
	if toolResultMap(t, blocked[5])["isError"] != true {
		t.Fatalf("expected humanize blocked until seo.md, got %v", blocked[5])
	}

	runLines(t, s, toolCallLine(t, 7, "save_seo_bundle", validSEOArgs(piecePath,
		"# Title\n\n*Subtitle.*\n\nA claim backed by [research](https://example.com/retry-budget).\n")))

	ok := runLines(t, s, toolCallLine(t, 8, "save_humanize_bundle", map[string]interface{}{
		"piece_path": piecePath,
		"revised_md": "# Title\n\n*Subtitle.*\n\nA claim backed by [research](https://example.com/retry-budget).\n",
		"checks":     []string{"[hedge] none needed"},
	}))
	if toolResultMap(t, ok[7])["isError"] == true {
		t.Fatalf("humanize after seo should succeed: %v", ok[7])
	}
}

func TestUnmarkedPieceStillDraftsWithoutResearch(t *testing.T) {
	dir := t.TempDir()
	s := New("pysar", "test", dir, t.TempDir(), nil, &bytes.Buffer{})
	piecePath := createTestPiece(t, s)

	resp := runLines(t, s, toolCallLine(t, 2, "save_draft_bundle", map[string]interface{}{
		"piece_path": piecePath,
		"draft_md":   "# Title\n\n*Subtitle.*\n\nNo research required on unmarked pieces.\n",
	}))
	if toolResultMap(t, resp[1])["isError"] == true {
		t.Fatalf("unmarked draft must succeed: %v", resp[1])
	}
}
