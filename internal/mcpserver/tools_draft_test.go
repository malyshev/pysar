package mcpserver

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func createTestPiece(t *testing.T, s *Server) (piecePath string) {
	t.Helper()
	intakeArgs := map[string]interface{}{
		"name":        "docker-for-developers",
		"idea":        "write something about docker for developers who just need it to work",
		"entry_mode":  "idea",
		"pov_source":  "author-practitioner",
		"restatement": "For working developers, plain language, practical Docker defaults",
		"audience":    "developers",
		"register":    "plain language",
		"expert_lens": "distributed-systems / platform engineering",
		"topic_scope": "practical Docker defaults",
		"thesis":      "Most Docker pain comes from ignoring the defaults, not Docker itself",
		"promise":     "Three settings that fix most Docker complaints",
		"killer_sections": []map[string]interface{}{
			{"title": "The setting everyone skips", "edge": "Nobody documents this default", "example": "docker run with --init shown side by side"},
		},
		"counterintuitive": []map[string]interface{}{
			{"claim": "More flags make Docker safer", "contradiction": "most flags just paper over a missing default"},
		},
		"key_questions": []string{"What should I set by default?"},
		"non_goals":     []string{"Kubernetes"},
		"outline_md":    "## Why\n### Scene\n## Fixes\n",
		"angles_md":     "## Misconceptions\n- x\n\n## Contrarian / under-discussed\n- Defaults over flags\n\n## Trade-offs\n- y\n\n## Edge cases\n- z\n",
	}
	resp := runLines(t, s, toolCallLine(t, 1, "save_intake_bundle", intakeArgs))
	result := toolResultMap(t, resp[0])
	if result["isError"] == true {
		t.Fatalf("save_intake_bundle error: %v", result)
	}
	entries, err := os.ReadDir(filepath.Join(s.baseDir, ".pysar", "pieces"))
	if err != nil || len(entries) != 1 {
		t.Fatalf("expected exactly one piece dir, got %v, err=%v", entries, err)
	}
	return ".pysar/pieces/" + entries[0].Name()
}

func TestSaveDraftBundleEndToEndWithResearch(t *testing.T) {
	dir := t.TempDir()
	s := New("pysar", "test", dir, t.TempDir(), nil, &bytes.Buffer{})
	piecePath := createTestPiece(t, s)

	researchArgs := map[string]interface{}{
		"piece_path":  piecePath,
		"expert_lens": "distributed-systems / platform engineering",
		"sources":     []map[string]interface{}{validResearchSource("docker-docs-init")},
	}
	researchResp := runLines(t, s, toolCallLine(t, 2, "save_research_bundle", researchArgs))
	if toolResultMap(t, researchResp[1])["isError"] == true {
		t.Fatalf("save_research_bundle error: %v", researchResp[1])
	}

	draftMD := "# Fix Docker in three settings\n\n" +
		"*Most Docker pain is a missing default, not a Docker flaw.*\n\n" +
		"Developers blame Docker for problems its own defaults would have caught[^docker-docs-init].\n\n" +
		"## The setting everyone skips\n\n" +
		"Run with --init and most zombie-process complaints disappear.\n"

	draftArgs := map[string]interface{}{
		"piece_path": piecePath,
		"draft_md":   draftMD,
		"notes":      "first full draft",
	}
	draftResp := runLines(t, s, toolCallLine(t, 3, "save_draft_bundle", draftArgs))
	draftResult := toolResultMap(t, draftResp[2])
	if draftResult["isError"] == true {
		t.Fatalf("save_draft_bundle error: %v", draftResult)
	}
	text := toolText(t, draftResult)
	if !strings.Contains(text, "first full draft exists") {
		t.Fatalf("expected the author-facing outcome message, got %s", text)
	}

	pieceDir := filepath.Join(dir, filepath.FromSlash(piecePath))
	for _, f := range []string{"draft.md", "draft-changelog.md", "run-log.jsonl"} {
		if _, err := os.Stat(filepath.Join(pieceDir, f)); err != nil {
			t.Fatalf("missing %s: %v", f, err)
		}
	}
	got, _ := os.ReadFile(filepath.Join(pieceDir, "draft.md"))
	if !strings.Contains(string(got), "Fix Docker in three settings") {
		t.Fatalf("draft.md missing expected content:\n%s", got)
	}
	runLog, _ := os.ReadFile(filepath.Join(pieceDir, "run-log.jsonl"))
	if !strings.Contains(string(runLog), `"pass":"draft"`) {
		t.Fatalf("run-log missing draft entry: %s", runLog)
	}

	// Brief/outline/angles/sources must never be touched by a draft pass.
	brief, _ := os.ReadFile(filepath.Join(pieceDir, "brief.md"))
	if !strings.Contains(string(brief), "Most Docker pain comes from ignoring the defaults, not Docker itself") {
		t.Fatalf("thesis was touched by draft:\n%s", brief)
	}
}

func TestSaveDraftBundleRejectsCitationWithNoResearch(t *testing.T) {
	dir := t.TempDir()
	s := New("pysar", "test", dir, t.TempDir(), nil, &bytes.Buffer{})
	piecePath := createTestPiece(t, s)

	draftMD := "# A title\n\n*Subtitle.*\n\nA claim that cites a source that was never fetched[^ghost-source].\n"
	draftArgs := map[string]interface{}{"piece_path": piecePath, "draft_md": draftMD}
	resp := runLines(t, s, toolCallLine(t, 2, "save_draft_bundle", draftArgs))
	result := toolResultMap(t, resp[1])
	if result["isError"] != true {
		t.Fatalf("expected error for a citation with no research behind it, got %v", result)
	}
}

func TestSaveDraftBundleRejectsRawURL(t *testing.T) {
	dir := t.TempDir()
	s := New("pysar", "test", dir, t.TempDir(), nil, &bytes.Buffer{})
	piecePath := createTestPiece(t, s)

	draftMD := "# A title\n\n*Subtitle.*\n\nSee https://example.com directly.\n"
	draftArgs := map[string]interface{}{"piece_path": piecePath, "draft_md": draftMD}
	resp := runLines(t, s, toolCallLine(t, 2, "save_draft_bundle", draftArgs))
	result := toolResultMap(t, resp[1])
	if result["isError"] != true {
		t.Fatalf("expected error for a raw URL in prose, got %v", result)
	}
}

func TestSaveDraftBundleRejectsMissingPiece(t *testing.T) {
	dir := t.TempDir()
	s := New("pysar", "test", dir, t.TempDir(), nil, &bytes.Buffer{})

	draftArgs := map[string]interface{}{
		"piece_path": ".pysar/pieces/does-not-exist",
		"draft_md":   "# A title\n\n*Subtitle.*\n\nbody\n",
	}
	resp := runLines(t, s, toolCallLine(t, 1, "save_draft_bundle", draftArgs))
	result := toolResultMap(t, resp[0])
	if result["isError"] != true {
		t.Fatalf("expected error for a piece path with no brief.md, got %v", result)
	}
}

func TestSaveDraftBundleRedraftReplacesDraftWholesale(t *testing.T) {
	dir := t.TempDir()
	s := New("pysar", "test", dir, t.TempDir(), nil, &bytes.Buffer{})
	piecePath := createTestPiece(t, s)

	first := map[string]interface{}{
		"piece_path": piecePath,
		"draft_md":   "# First draft\n\n*Subtitle.*\n\nFirst body.\n",
	}
	firstResp := runLines(t, s, toolCallLine(t, 2, "save_draft_bundle", first))
	if toolResultMap(t, firstResp[1])["isError"] == true {
		t.Fatalf("first draft error: %v", firstResp[1])
	}

	second := map[string]interface{}{
		"piece_path": piecePath,
		"draft_md":   "# Second draft\n\n*Subtitle.*\n\nSecond body.\n",
		"notes":      "revision",
	}
	secondResp := runLines(t, s, toolCallLine(t, 3, "save_draft_bundle", second))
	if toolResultMap(t, secondResp[2])["isError"] == true {
		t.Fatalf("second draft error: %v", secondResp[2])
	}

	pieceDir := filepath.Join(dir, filepath.FromSlash(piecePath))
	got, _ := os.ReadFile(filepath.Join(pieceDir, "draft.md"))
	if strings.Contains(string(got), "First draft") || !strings.Contains(string(got), "Second draft") {
		t.Fatalf("expected redraft to replace draft.md wholesale, got:\n%s", got)
	}
	changelog, _ := os.ReadFile(filepath.Join(pieceDir, "draft-changelog.md"))
	if strings.Count(string(changelog), "\n") != 2 {
		t.Fatalf("expected two changelog lines (one per draft call), got:\n%s", changelog)
	}
}
