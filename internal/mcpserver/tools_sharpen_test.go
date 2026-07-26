package mcpserver

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestSaveSharpenBundleEndToEndDirectlyAfterDraft(t *testing.T) {
	// Staff-edit is optional -- sharpen must work directly after draft, the
	// same way draft itself works without research having run.
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

	draftMD := "# Fix Docker in Three Settings\n\n" +
		"*Most Docker pain is a missing default, not a Docker flaw.*\n\n" +
		"Developers blame Docker for problems its own defaults would have caught[^docker-docs-init].\n\n" +
		"## The setting everyone skips\n\n" +
		"Run with --init and most zombie-process complaints disappear.\n"
	draftResp := runLines(t, s, toolCallLine(t, 3, "save_draft_bundle", map[string]interface{}{
		"piece_path": piecePath,
		"draft_md":   draftMD,
	}))
	if toolResultMap(t, draftResp[2])["isError"] == true {
		t.Fatalf("save_draft_bundle error: %v", draftResp[2])
	}

	revisedMD := "# The Zombie Process Nobody Warns You About\n\n" +
		"*A container's zombie process is the specific failure that makes people blame Docker.*\n\n" +
		"Developers blame Docker for problems its own defaults would have caught[^docker-docs-init].\n\n" +
		"## The setting everyone skips\n\n" +
		"Run with --init and most zombie-process complaints disappear -- the exact fix for the failure the opening named.\n"
	sharpenResp := runLines(t, s, toolCallLine(t, 4, "save_sharpen_bundle", map[string]interface{}{
		"piece_path": piecePath,
		"revised_md": revisedMD,
		"checks":     []string{"[opener] renamed the title and subtitle to lead with the zombie-process hook instead of a generic settings claim", "[arc] closed the piece by tying the fix back to the exact failure the opener named"},
		"mode":       "delta",
	}))
	sharpenResult := toolResultMap(t, sharpenResp[3])
	if sharpenResult["isError"] == true {
		t.Fatalf("save_sharpen_bundle error: %v", sharpenResult)
	}
	text := toolText(t, sharpenResult)
	if !strings.Contains(text, "sharpen done") {
		t.Fatalf("expected the author-facing outcome message, got %s", text)
	}

	pieceDir := filepath.Join(dir, filepath.FromSlash(piecePath))
	for _, f := range []string{"draft.md", "sharpen.md", "sharpen-changelog.md", "run-log.jsonl"} {
		if _, err := os.Stat(filepath.Join(pieceDir, f)); err != nil {
			t.Fatalf("missing %s: %v", f, err)
		}
	}
	if _, err := os.Stat(filepath.Join(pieceDir, "staff-edit.md")); err == nil {
		t.Fatal("staff-edit.md should not exist -- this test never ran staff-edit")
	}

	gotDraft, _ := os.ReadFile(filepath.Join(pieceDir, "draft.md"))
	if !strings.Contains(string(gotDraft), "Fix Docker in Three Settings") {
		t.Fatalf("draft.md must remain untouched by sharpen, got:\n%s", gotDraft)
	}

	gotSharpen, _ := os.ReadFile(filepath.Join(pieceDir, "sharpen.md"))
	if !strings.Contains(string(gotSharpen), "The Zombie Process Nobody Warns You About") {
		t.Fatalf("sharpen.md missing the revised content:\n%s", gotSharpen)
	}

	changelog, _ := os.ReadFile(filepath.Join(pieceDir, "sharpen-changelog.md"))
	if !strings.Contains(string(changelog), "[opener]") || !strings.Contains(string(changelog), "[arc]") {
		t.Fatalf("expected both recorded checks in the changelog:\n%s", changelog)
	}

	runLog, _ := os.ReadFile(filepath.Join(pieceDir, "run-log.jsonl"))
	if !strings.Contains(string(runLog), `"pass":"sharpen"`) {
		t.Fatalf("run-log missing sharpen entry: %s", runLog)
	}

	brief, _ := os.ReadFile(filepath.Join(pieceDir, "brief.md"))
	if !strings.Contains(string(brief), "Most Docker pain comes from ignoring the defaults, not Docker itself") {
		t.Fatalf("thesis was touched by sharpen:\n%s", brief)
	}
}

func TestSaveSharpenBundleEndToEndAfterStaffEdit(t *testing.T) {
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
	sharpenResp := runLines(t, s, toolCallLine(t, 4, "save_sharpen_bundle", map[string]interface{}{
		"piece_path": piecePath,
		"revised_md": "# Sharpened Title\n\n*Subtitle.*\n\nsharpened body\n",
		"checks":     []string{"[opener] sharpened further"},
	}))
	if toolResultMap(t, sharpenResp[3])["isError"] == true {
		t.Fatalf("save_sharpen_bundle error: %v", sharpenResp[3])
	}

	pieceDir := filepath.Join(dir, filepath.FromSlash(piecePath))
	gotStaffEdit, _ := os.ReadFile(filepath.Join(pieceDir, "staff-edit.md"))
	if !strings.Contains(string(gotStaffEdit), "Staff-Edited Title") {
		t.Fatalf("staff-edit.md must remain untouched by sharpen, got:\n%s", gotStaffEdit)
	}
	gotSharpen, _ := os.ReadFile(filepath.Join(pieceDir, "sharpen.md"))
	if !strings.Contains(string(gotSharpen), "Sharpened Title") {
		t.Fatalf("sharpen.md missing the revised content:\n%s", gotSharpen)
	}
}

func TestSaveSharpenBundleRejectsWithoutDraft(t *testing.T) {
	dir := t.TempDir()
	s := New("pysar", "test", dir, t.TempDir(), nil, &bytes.Buffer{})
	piecePath := createTestPiece(t, s)

	resp := runLines(t, s, toolCallLine(t, 2, "save_sharpen_bundle", map[string]interface{}{
		"piece_path": piecePath,
		"revised_md": "# Title\n\n*Subtitle.*\n\nbody\n",
		"checks":     []string{"[opener] tightened the hook"},
	}))
	result := toolResultMap(t, resp[1])
	if result["isError"] != true {
		t.Fatalf("expected error when no draft.md exists yet, got %v", result)
	}
}

func TestSaveSharpenBundleRejectsZeroChecks(t *testing.T) {
	dir := t.TempDir()
	s := New("pysar", "test", dir, t.TempDir(), nil, &bytes.Buffer{})
	piecePath := createTestPiece(t, s)

	runLines(t, s, toolCallLine(t, 2, "save_draft_bundle", map[string]interface{}{
		"piece_path": piecePath,
		"draft_md":   "# Title\n\n*Subtitle.*\n\nbody\n",
	}))

	resp := runLines(t, s, toolCallLine(t, 3, "save_sharpen_bundle", map[string]interface{}{
		"piece_path": piecePath,
		"revised_md": "# Title\n\n*Subtitle.*\n\nrevised body\n",
		"checks":     []string{},
	}))
	result := toolResultMap(t, resp[2])
	if result["isError"] != true {
		t.Fatalf("expected error for zero recorded checks, got %v", result)
	}
}

func TestSaveSharpenBundleRejectsUnmatchedCitation(t *testing.T) {
	dir := t.TempDir()
	s := New("pysar", "test", dir, t.TempDir(), nil, &bytes.Buffer{})
	piecePath := createTestPiece(t, s)

	runLines(t, s, toolCallLine(t, 2, "save_draft_bundle", map[string]interface{}{
		"piece_path": piecePath,
		"draft_md":   "# Title\n\n*Subtitle.*\n\nbody\n",
	}))

	resp := runLines(t, s, toolCallLine(t, 3, "save_sharpen_bundle", map[string]interface{}{
		"piece_path": piecePath,
		"revised_md": "# Title\n\n*Subtitle.*\n\nA new unsupported claim[^ghost-source].\n",
		"checks":     []string{"[elevate] added a citation"},
	}))
	result := toolResultMap(t, resp[2])
	if result["isError"] != true {
		t.Fatalf("expected error for a citation with no research behind it, got %v", result)
	}
}

func TestSaveSharpenBundleRerunReplacesSharpenFileWholesale(t *testing.T) {
	dir := t.TempDir()
	s := New("pysar", "test", dir, t.TempDir(), nil, &bytes.Buffer{})
	piecePath := createTestPiece(t, s)

	runLines(t, s, toolCallLine(t, 2, "save_draft_bundle", map[string]interface{}{
		"piece_path": piecePath,
		"draft_md":   "# Original Draft\n\n*Subtitle.*\n\nOriginal body.\n",
	}))
	runLines(t, s, toolCallLine(t, 3, "save_sharpen_bundle", map[string]interface{}{
		"piece_path": piecePath,
		"revised_md": "# First Sharpen\n\n*Subtitle.*\n\nFirst sharpened body.\n",
		"checks":     []string{"[opener] first pass"},
	}))
	secondResp := runLines(t, s, toolCallLine(t, 4, "save_sharpen_bundle", map[string]interface{}{
		"piece_path": piecePath,
		"revised_md": "# Second Sharpen\n\n*Subtitle.*\n\nSecond sharpened body.\n",
		"checks":     []string{"[arc] second pass"},
	}))
	if toolResultMap(t, secondResp[3])["isError"] == true {
		t.Fatalf("second sharpen error: %v", secondResp[3])
	}

	pieceDir := filepath.Join(dir, filepath.FromSlash(piecePath))
	gotSharpen, _ := os.ReadFile(filepath.Join(pieceDir, "sharpen.md"))
	if strings.Contains(string(gotSharpen), "First Sharpen") || !strings.Contains(string(gotSharpen), "Second Sharpen") {
		t.Fatalf("expected re-run to replace sharpen.md wholesale, got:\n%s", gotSharpen)
	}

	gotDraft, _ := os.ReadFile(filepath.Join(pieceDir, "draft.md"))
	if !strings.Contains(string(gotDraft), "Original Draft") {
		t.Fatalf("draft.md must survive both sharpen runs untouched, got:\n%s", gotDraft)
	}
}
