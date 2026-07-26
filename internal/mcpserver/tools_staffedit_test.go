package mcpserver

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestSaveStaffEditBundleEndToEnd(t *testing.T) {
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
	draftResp := runLines(t, s, toolCallLine(t, 3, "save_draft_bundle", map[string]interface{}{
		"piece_path": piecePath,
		"draft_md":   draftMD,
	}))
	if toolResultMap(t, draftResp[2])["isError"] == true {
		t.Fatalf("save_draft_bundle error: %v", draftResp[2])
	}

	revisedMD := "# The Docker Default That Actually Matters\n\n" +
		"*Most Docker pain is a missing default, not a Docker flaw.*\n\n" +
		"A zombie process wedged inside a container is a specific, common failure a working developer hits within the first week[^docker-docs-init].\n\n" +
		"## The setting everyone skips\n\n" +
		"Run with --init and most zombie-process complaints disappear.\n"
	staffEditResp := runLines(t, s, toolCallLine(t, 4, "save_staff_edit_bundle", map[string]interface{}{
		"piece_path": piecePath,
		"revised_md": revisedMD,
		"checks":     []string{"[stakes] named the specific failure mode (zombie process) instead of a vague claim", "[readability] shortened the opener sentence"},
		"mode":       "delta",
	}))
	staffEditResult := toolResultMap(t, staffEditResp[3])
	if staffEditResult["isError"] == true {
		t.Fatalf("save_staff_edit_bundle error: %v", staffEditResult)
	}
	text := toolText(t, staffEditResult)
	if !strings.Contains(text, "staff edit done") {
		t.Fatalf("expected the author-facing outcome message, got %s", text)
	}

	pieceDir := filepath.Join(dir, filepath.FromSlash(piecePath))
	for _, f := range []string{"draft.md", "staff-edit.md", "staff-edit-changelog.md", "run-log.jsonl"} {
		if _, err := os.Stat(filepath.Join(pieceDir, f)); err != nil {
			t.Fatalf("missing %s: %v", f, err)
		}
	}

	// draft.md must stay exactly as /ps-draft left it -- staff-edit writes
	// its own file, so the original survives for comparison.
	gotDraft, _ := os.ReadFile(filepath.Join(pieceDir, "draft.md"))
	if !strings.Contains(string(gotDraft), "Fix Docker in three settings") {
		t.Fatalf("draft.md must remain untouched by staff-edit, got:\n%s", gotDraft)
	}
	if strings.Contains(string(gotDraft), "The Docker Default That Actually Matters") {
		t.Fatalf("staff-edit must not overwrite draft.md, got:\n%s", gotDraft)
	}

	gotStaffEdit, _ := os.ReadFile(filepath.Join(pieceDir, "staff-edit.md"))
	if !strings.Contains(string(gotStaffEdit), "The Docker Default That Actually Matters") {
		t.Fatalf("staff-edit.md missing the revised content:\n%s", gotStaffEdit)
	}

	changelog, _ := os.ReadFile(filepath.Join(pieceDir, "staff-edit-changelog.md"))
	if !strings.Contains(string(changelog), "[stakes]") || !strings.Contains(string(changelog), "[readability]") {
		t.Fatalf("expected both recorded checks in the changelog:\n%s", changelog)
	}

	runLog, _ := os.ReadFile(filepath.Join(pieceDir, "run-log.jsonl"))
	if !strings.Contains(string(runLog), `"pass":"staff-edit"`) {
		t.Fatalf("run-log missing staff-edit entry: %s", runLog)
	}

	// Brief/outline/angles/sources must never be touched by a staff-edit pass.
	brief, _ := os.ReadFile(filepath.Join(pieceDir, "brief.md"))
	if !strings.Contains(string(brief), "Most Docker pain comes from ignoring the defaults, not Docker itself") {
		t.Fatalf("thesis was touched by staff-edit:\n%s", brief)
	}
}

func TestSaveStaffEditBundleRejectsWithoutDraft(t *testing.T) {
	dir := t.TempDir()
	s := New("pysar", "test", dir, t.TempDir(), nil, &bytes.Buffer{})
	piecePath := createTestPiece(t, s)

	resp := runLines(t, s, toolCallLine(t, 2, "save_staff_edit_bundle", map[string]interface{}{
		"piece_path": piecePath,
		"revised_md": "# Title\n\n*Subtitle.*\n\nbody\n",
		"checks":     []string{"[stakes] tightened the opener"},
	}))
	result := toolResultMap(t, resp[1])
	if result["isError"] != true {
		t.Fatalf("expected error when no draft.md exists yet, got %v", result)
	}
}

func TestSaveStaffEditBundleRejectsNonexistentPieceEvenWithUnrelatedRootDraft(t *testing.T) {
	// A piece_path that doesn't resolve to any real piece must be rejected
	// even if the project root happens to contain an unrelated file named
	// draft.md (e.g. a leftover, a template) -- resolveAnchoredPass must not
	// mistake that coincidental file for this piece's draft.
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "draft.md"), []byte("unrelated root-level file"), 0o644); err != nil {
		t.Fatal(err)
	}
	s := New("pysar", "test", dir, t.TempDir(), nil, &bytes.Buffer{})

	resp := runLines(t, s, toolCallLine(t, 1, "save_staff_edit_bundle", map[string]interface{}{
		"piece_path": ".pysar/pieces/does-not-exist",
		"revised_md": "# Title\n\n*Subtitle.*\n\nbody\n",
		"checks":     []string{"[stakes] tightened the opener"},
	}))
	result := toolResultMap(t, resp[0])
	if result["isError"] != true {
		t.Fatalf("expected error for a piece_path that doesn't resolve to a real piece, got %v", result)
	}

	got, _ := os.ReadFile(filepath.Join(dir, "draft.md"))
	if string(got) != "unrelated root-level file" {
		t.Fatalf("the unrelated root-level draft.md must not be touched, got:\n%s", got)
	}
}

func TestSaveStaffEditBundleRejectsZeroChecks(t *testing.T) {
	dir := t.TempDir()
	s := New("pysar", "test", dir, t.TempDir(), nil, &bytes.Buffer{})
	piecePath := createTestPiece(t, s)

	runLines(t, s, toolCallLine(t, 2, "save_draft_bundle", map[string]interface{}{
		"piece_path": piecePath,
		"draft_md":   "# Title\n\n*Subtitle.*\n\nbody\n",
	}))

	resp := runLines(t, s, toolCallLine(t, 3, "save_staff_edit_bundle", map[string]interface{}{
		"piece_path": piecePath,
		"revised_md": "# Title\n\n*Subtitle.*\n\nrevised body\n",
		"checks":     []string{},
	}))
	result := toolResultMap(t, resp[2])
	if result["isError"] != true {
		t.Fatalf("expected error for zero recorded checks, got %v", result)
	}
}

func TestSaveStaffEditBundleRejectsUnmatchedCitation(t *testing.T) {
	dir := t.TempDir()
	s := New("pysar", "test", dir, t.TempDir(), nil, &bytes.Buffer{})
	piecePath := createTestPiece(t, s)

	runLines(t, s, toolCallLine(t, 2, "save_draft_bundle", map[string]interface{}{
		"piece_path": piecePath,
		"draft_md":   "# Title\n\n*Subtitle.*\n\nbody\n",
	}))

	resp := runLines(t, s, toolCallLine(t, 3, "save_staff_edit_bundle", map[string]interface{}{
		"piece_path": piecePath,
		"revised_md": "# Title\n\n*Subtitle.*\n\nA new unsupported claim[^ghost-source].\n",
		"checks":     []string{"[tech] added a citation"},
	}))
	result := toolResultMap(t, resp[2])
	if result["isError"] != true {
		t.Fatalf("expected error for a citation with no research behind it, got %v", result)
	}
}

func TestSaveStaffEditBundleRerunReplacesStaffEditFileWholesaleWithoutTouchingDraft(t *testing.T) {
	dir := t.TempDir()
	s := New("pysar", "test", dir, t.TempDir(), nil, &bytes.Buffer{})
	piecePath := createTestPiece(t, s)

	runLines(t, s, toolCallLine(t, 2, "save_draft_bundle", map[string]interface{}{
		"piece_path": piecePath,
		"draft_md":   "# Original Draft\n\n*Subtitle.*\n\nOriginal body.\n",
	}))
	runLines(t, s, toolCallLine(t, 3, "save_staff_edit_bundle", map[string]interface{}{
		"piece_path": piecePath,
		"revised_md": "# First Revision\n\n*Subtitle.*\n\nFirst revised body.\n",
		"checks":     []string{"[stakes] first pass"},
	}))
	secondResp := runLines(t, s, toolCallLine(t, 4, "save_staff_edit_bundle", map[string]interface{}{
		"piece_path": piecePath,
		"revised_md": "# Second Revision\n\n*Subtitle.*\n\nSecond revised body.\n",
		"checks":     []string{"[readability] second pass"},
	}))
	if toolResultMap(t, secondResp[3])["isError"] == true {
		t.Fatalf("second staff-edit error: %v", secondResp[3])
	}

	pieceDir := filepath.Join(dir, filepath.FromSlash(piecePath))
	gotStaffEdit, _ := os.ReadFile(filepath.Join(pieceDir, "staff-edit.md"))
	if strings.Contains(string(gotStaffEdit), "First Revision") || !strings.Contains(string(gotStaffEdit), "Second Revision") {
		t.Fatalf("expected re-run to replace staff-edit.md wholesale, got:\n%s", gotStaffEdit)
	}

	gotDraft, _ := os.ReadFile(filepath.Join(pieceDir, "draft.md"))
	if !strings.Contains(string(gotDraft), "Original Draft") {
		t.Fatalf("draft.md must survive both staff-edit runs untouched, got:\n%s", gotDraft)
	}
}
