package mcpserver

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestSaveHumanizeBundleEndToEndDirectlyAfterDraft(t *testing.T) {
	// Staff-edit and sharpen are both optional -- humanize must work
	// directly after draft, the same way sharpen works without staff-edit.
	dir := t.TempDir()
	s := New("pysar", "test", dir, t.TempDir(), nil, &bytes.Buffer{})
	piecePath := createTestPiece(t, s)

	runLines(t, s, toolCallLine(t, 2, "save_draft_bundle", map[string]interface{}{
		"piece_path": piecePath,
		"draft_md":   "# Fix Docker in Three Settings\n\n*Subtitle.*\n\nDevelopers blame Docker for problems its own defaults would have caught.\n",
	}))

	revisedMD := "# Fix Docker Without the Ceremony\n\n" +
		"*Turns out most Docker complaints trace back to three defaults nobody bothered to change.*\n\n" +
		"Developers blame Docker for problems its own defaults would have caught -- and honestly, most of the time it's an easy fix.\n"
	humanizeResp := runLines(t, s, toolCallLine(t, 3, "save_humanize_bundle", map[string]interface{}{
		"piece_path": piecePath,
		"revised_md": revisedMD,
		"checks":     []string{"[hedge-stack] dropped a redundant qualifier in the opener", "[rhythm] varied two consecutive same-length sentences"},
		"mode":       "delta",
	}))
	humanizeResult := toolResultMap(t, humanizeResp[2])
	if humanizeResult["isError"] == true {
		t.Fatalf("save_humanize_bundle error: %v", humanizeResult)
	}
	text := toolText(t, humanizeResult)
	if !strings.Contains(text, "sounds like you, not the machine") {
		t.Fatalf("expected the author-facing outcome message, got %s", text)
	}

	pieceDir := filepath.Join(dir, filepath.FromSlash(piecePath))
	for _, f := range []string{"draft.md", "humanize.md", "humanize-changelog.md", "run-log.jsonl"} {
		if _, err := os.Stat(filepath.Join(pieceDir, f)); err != nil {
			t.Fatalf("missing %s: %v", f, err)
		}
	}
	for _, f := range []string{"staff-edit.md", "sharpen.md"} {
		if _, err := os.Stat(filepath.Join(pieceDir, f)); err == nil {
			t.Fatalf("%s should not exist -- this test never ran that pass", f)
		}
	}

	gotDraft, _ := os.ReadFile(filepath.Join(pieceDir, "draft.md"))
	if !strings.Contains(string(gotDraft), "Fix Docker in Three Settings") {
		t.Fatalf("draft.md must remain untouched by humanize, got:\n%s", gotDraft)
	}

	gotHumanize, _ := os.ReadFile(filepath.Join(pieceDir, "humanize.md"))
	if !strings.Contains(string(gotHumanize), "Fix Docker Without the Ceremony") {
		t.Fatalf("humanize.md missing the revised content:\n%s", gotHumanize)
	}

	changelog, _ := os.ReadFile(filepath.Join(pieceDir, "humanize-changelog.md"))
	if !strings.Contains(string(changelog), "[hedge-stack]") || !strings.Contains(string(changelog), "[rhythm]") {
		t.Fatalf("expected both recorded checks in the changelog:\n%s", changelog)
	}

	runLog, _ := os.ReadFile(filepath.Join(pieceDir, "run-log.jsonl"))
	if !strings.Contains(string(runLog), `"pass":"humanize"`) {
		t.Fatalf("run-log missing humanize entry: %s", runLog)
	}
	if !strings.Contains(string(runLog), "revised_from=draft.md") {
		t.Fatalf("run-log should record revised_from=draft.md when neither staff-edit.md nor sharpen.md exists: %s", runLog)
	}

	brief, _ := os.ReadFile(filepath.Join(pieceDir, "brief.md"))
	if !strings.Contains(string(brief), "Most Docker pain comes from ignoring the defaults, not Docker itself") {
		t.Fatalf("thesis was touched by humanize:\n%s", brief)
	}
}

func TestSaveHumanizeBundleEndToEndAfterStaffEditAndSharpen(t *testing.T) {
	dir := t.TempDir()
	s := New("pysar", "test", dir, t.TempDir(), nil, &bytes.Buffer{})
	piecePath := createTestPiece(t, s)

	runLines(t, s, toolCallLine(t, 2, "save_draft_bundle", map[string]interface{}{
		"piece_path": piecePath,
		"draft_md":   "# Title\n\n*Subtitle.*\n\nbody\n",
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
	humanizeResp := runLines(t, s, toolCallLine(t, 5, "save_humanize_bundle", map[string]interface{}{
		"piece_path": piecePath,
		"revised_md": "# Humanized Title\n\n*Subtitle.*\n\nhumanized body\n",
		"checks":     []string{"[hedge-stack] dropped a redundant qualifier"},
	}))
	if toolResultMap(t, humanizeResp[4])["isError"] == true {
		t.Fatalf("save_humanize_bundle error: %v", humanizeResp[4])
	}

	pieceDir := filepath.Join(dir, filepath.FromSlash(piecePath))
	gotSharpen, _ := os.ReadFile(filepath.Join(pieceDir, "sharpen.md"))
	if !strings.Contains(string(gotSharpen), "Sharpened Title") {
		t.Fatalf("sharpen.md must remain untouched by humanize, got:\n%s", gotSharpen)
	}
	gotStaffEdit, _ := os.ReadFile(filepath.Join(pieceDir, "staff-edit.md"))
	if !strings.Contains(string(gotStaffEdit), "Staff-Edited Title") {
		t.Fatalf("staff-edit.md must remain untouched by humanize, got:\n%s", gotStaffEdit)
	}
	gotHumanize, _ := os.ReadFile(filepath.Join(pieceDir, "humanize.md"))
	if !strings.Contains(string(gotHumanize), "Humanized Title") {
		t.Fatalf("humanize.md missing the revised content:\n%s", gotHumanize)
	}

	runLog, _ := os.ReadFile(filepath.Join(pieceDir, "run-log.jsonl"))
	if !strings.Contains(string(runLog), "revised_from=sharpen.md") {
		t.Fatalf("run-log should record revised_from=sharpen.md when it's the most-refined file present: %s", runLog)
	}
}

func TestSaveHumanizeBundleRejectsWithoutDraft(t *testing.T) {
	dir := t.TempDir()
	s := New("pysar", "test", dir, t.TempDir(), nil, &bytes.Buffer{})
	piecePath := createTestPiece(t, s)

	resp := runLines(t, s, toolCallLine(t, 2, "save_humanize_bundle", map[string]interface{}{
		"piece_path": piecePath,
		"revised_md": "# Title\n\n*Subtitle.*\n\nbody\n",
		"checks":     []string{"[hedge-stack] dropped a redundant qualifier"},
	}))
	result := toolResultMap(t, resp[1])
	if result["isError"] != true {
		t.Fatalf("expected error when no draft.md exists yet, got %v", result)
	}
}

func TestSaveHumanizeBundleRejectsZeroChecks(t *testing.T) {
	dir := t.TempDir()
	s := New("pysar", "test", dir, t.TempDir(), nil, &bytes.Buffer{})
	piecePath := createTestPiece(t, s)

	runLines(t, s, toolCallLine(t, 2, "save_draft_bundle", map[string]interface{}{
		"piece_path": piecePath,
		"draft_md":   "# Title\n\n*Subtitle.*\n\nbody\n",
	}))

	resp := runLines(t, s, toolCallLine(t, 3, "save_humanize_bundle", map[string]interface{}{
		"piece_path": piecePath,
		"revised_md": "# Title\n\n*Subtitle.*\n\nrevised body\n",
		"checks":     []string{},
	}))
	result := toolResultMap(t, resp[2])
	if result["isError"] != true {
		t.Fatalf("expected error for zero recorded checks, got %v", result)
	}
}

func TestSaveHumanizeBundleRejectsUnmatchedCitation(t *testing.T) {
	dir := t.TempDir()
	s := New("pysar", "test", dir, t.TempDir(), nil, &bytes.Buffer{})
	piecePath := createTestPiece(t, s)

	runLines(t, s, toolCallLine(t, 2, "save_draft_bundle", map[string]interface{}{
		"piece_path": piecePath,
		"draft_md":   "# Title\n\n*Subtitle.*\n\nbody\n",
	}))

	resp := runLines(t, s, toolCallLine(t, 3, "save_humanize_bundle", map[string]interface{}{
		"piece_path": piecePath,
		"revised_md": "# Title\n\n*Subtitle.*\n\nA new unsupported claim[^ghost-source].\n",
		"checks":     []string{"[symmetry] varied a too-uniform list"},
	}))
	result := toolResultMap(t, resp[2])
	if result["isError"] != true {
		t.Fatalf("expected error for a citation with no research behind it, got %v", result)
	}
}

func TestSaveHumanizeBundleRerunReplacesHumanizeFileWholesale(t *testing.T) {
	dir := t.TempDir()
	s := New("pysar", "test", dir, t.TempDir(), nil, &bytes.Buffer{})
	piecePath := createTestPiece(t, s)

	runLines(t, s, toolCallLine(t, 2, "save_draft_bundle", map[string]interface{}{
		"piece_path": piecePath,
		"draft_md":   "# Original Draft\n\n*Subtitle.*\n\nOriginal body.\n",
	}))
	runLines(t, s, toolCallLine(t, 3, "save_humanize_bundle", map[string]interface{}{
		"piece_path": piecePath,
		"revised_md": "# First Humanize\n\n*Subtitle.*\n\nFirst humanized body.\n",
		"checks":     []string{"[hedge-stack] first pass"},
	}))
	secondResp := runLines(t, s, toolCallLine(t, 4, "save_humanize_bundle", map[string]interface{}{
		"piece_path": piecePath,
		"revised_md": "# Second Humanize\n\n*Subtitle.*\n\nSecond humanized body.\n",
		"checks":     []string{"[rhythm] second pass"},
	}))
	if toolResultMap(t, secondResp[3])["isError"] == true {
		t.Fatalf("second humanize error: %v", secondResp[3])
	}

	pieceDir := filepath.Join(dir, filepath.FromSlash(piecePath))
	gotHumanize, _ := os.ReadFile(filepath.Join(pieceDir, "humanize.md"))
	if strings.Contains(string(gotHumanize), "First Humanize") || !strings.Contains(string(gotHumanize), "Second Humanize") {
		t.Fatalf("expected re-run to replace humanize.md wholesale, got:\n%s", gotHumanize)
	}

	gotDraft, _ := os.ReadFile(filepath.Join(pieceDir, "draft.md"))
	if !strings.Contains(string(gotDraft), "Original Draft") {
		t.Fatalf("draft.md must survive both humanize runs untouched, got:\n%s", gotDraft)
	}
}
