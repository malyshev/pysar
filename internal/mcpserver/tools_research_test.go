package mcpserver

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func validResearchSource(shortname string) map[string]interface{} {
	return map[string]interface{}{
		"shortname":   shortname,
		"url":         "https://example.com/" + shortname,
		"tier":        "primary",
		"accessed":    "2026-07-25",
		"key_claims":  []string{"a real claim this source backs"},
		"notes":       "a real note about this source",
		"raw_excerpt": "a verbatim quoted sentence from the source",
	}
}

func TestSaveResearchBundleStandaloneEndToEnd(t *testing.T) {
	dir := t.TempDir()
	s := New("pysar", "test", dir, t.TempDir(), nil, &bytes.Buffer{})

	args := map[string]interface{}{
		"topic":       "tomato blossom end rot",
		"expert_lens": "horticulture",
		"sources":     []map[string]interface{}{validResearchSource("uconn-ber")},
	}
	resp := runLines(t, s, toolCallLine(t, 1, "save_research_bundle", args))
	result := toolResultMap(t, resp[0])
	if result["isError"] == true {
		t.Fatalf("save_research_bundle error: %v", result)
	}
	text := toolText(t, result)
	if !strings.Contains(text, "/ps-intake --from-draft=") {
		t.Fatalf("expected hand-off hint to /ps-intake, got %s", text)
	}

	entries, err := os.ReadDir(filepath.Join(dir, ".pysar", "research"))
	if err != nil {
		t.Fatalf("read research dir: %v", err)
	}
	if len(entries) != 1 || !strings.HasPrefix(entries[0].Name(), "tomato-blossom-end-rot-") {
		t.Fatalf("expected one suffixed research dir, got %v", entries)
	}
}

func TestSaveResearchBundlePieceAnchoredEndToEnd(t *testing.T) {
	dir := t.TempDir()
	s := New("pysar", "test", dir, t.TempDir(), nil, &bytes.Buffer{})

	intakeArgs := map[string]interface{}{
		"name":        "ai-safety-for-parents",
		"idea":        "write something about AI safety for parents who feel overwhelmed",
		"entry_mode":  "idea",
		"pov_source":  "author-practitioner",
		"restatement": "For overwhelmed parents, plain language, household AI safety defaults",
		"audience":    "parents",
		"register":    "plain language",
		"expert_lens": "child-safety / parenting practitioner",
		"topic_scope": "household AI safety defaults",
		"thesis":      "Defaults beat lectures for family AI safety",
		"promise":     "Three rules you can apply tonight",
		"killer_sections": []map[string]interface{}{
			{"title": "The setting that matters more than the brand", "edge": "Compares default safety settings, not brands", "example": "Router-level DNS filter setup screenshots"},
		},
		"counterintuitive": []map[string]interface{}{
			{"claim": "More filters can hide risk", "contradiction": "hiding the conversation removes the actual safety net: a kid talking to a parent"},
		},
		"key_questions": []string{"What first?"},
		"non_goals":     []string{"Benchmarks"},
		"outline_md":    "## Why\n### Scene\n## Rules\n",
		"angles_md":     "## Misconceptions\n- x\n\n## Contrarian / under-discussed\n- Defaults\n\n## Trade-offs\n- y\n\n## Edge cases\n- z\n",
	}
	intakeResp := runLines(t, s, toolCallLine(t, 1, "save_intake_bundle", intakeArgs))
	intakeResult := toolResultMap(t, intakeResp[0])
	if intakeResult["isError"] == true {
		t.Fatalf("save_intake_bundle error: %v", intakeResult)
	}

	entries, err := os.ReadDir(filepath.Join(dir, ".pysar", "pieces"))
	if err != nil || len(entries) != 1 {
		t.Fatalf("expected exactly one piece dir, got %v, err=%v", entries, err)
	}
	piecePath := ".pysar/pieces/" + entries[0].Name()

	researchArgs := map[string]interface{}{
		"piece_path":              piecePath,
		"expert_lens":             "child-safety / parenting practitioner",
		"sources":                 []map[string]interface{}{validResearchSource("aap-screen-time")},
		"key_questions_additions": []string{"What does the AAP actually recommend?"},
		"angles_misconceptions":   []string{"Parental control apps alone solve the problem -- they don't"},
	}
	// runLines re-parses the server's whole output buffer each call (it
	// isn't reset between calls), so this second call's response is at
	// index 1 -- index 0 is still the first call's (intake) response.
	researchResp := runLines(t, s, toolCallLine(t, 2, "save_research_bundle", researchArgs))
	researchResult := toolResultMap(t, researchResp[1])
	if researchResult["isError"] == true {
		t.Fatalf("save_research_bundle error: %v", researchResult)
	}
	text := toolText(t, researchResult)
	if !strings.Contains(text, "your take stayed yours") {
		t.Fatalf("expected the author-facing outcome message, got %s", text)
	}

	pieceDir := filepath.Join(dir, filepath.FromSlash(piecePath))
	brief, _ := os.ReadFile(filepath.Join(pieceDir, "brief.md"))
	briefStr := string(brief)
	if !strings.Contains(briefStr, "research_mode: full") {
		t.Fatalf("expected research_mode: full:\n%s", briefStr)
	}
	if !strings.Contains(briefStr, "Defaults beat lectures for family AI safety") {
		t.Fatalf("thesis was touched:\n%s", briefStr)
	}
	if !strings.Contains(briefStr, "What does the AAP actually recommend?") {
		t.Fatalf("key question not appended:\n%s", briefStr)
	}

	angles, _ := os.ReadFile(filepath.Join(pieceDir, "angles.md"))
	if !strings.Contains(string(angles), "Parental control apps alone solve the problem") {
		t.Fatalf("misconception not appended:\n%s", angles)
	}

	for _, f := range []string{"sources.md", "raw/aap-screen-time.md", "run-log.jsonl"} {
		if _, err := os.Stat(filepath.Join(pieceDir, f)); err != nil {
			t.Fatalf("missing %s: %v", f, err)
		}
	}
	runLog, _ := os.ReadFile(filepath.Join(pieceDir, "run-log.jsonl"))
	if !strings.Contains(string(runLog), `"pass":"research"`) {
		t.Fatalf("run-log missing research entry: %s", runLog)
	}
}

func TestSaveResearchBundleRejectsMissingPiece(t *testing.T) {
	dir := t.TempDir()
	s := New("pysar", "test", dir, t.TempDir(), nil, &bytes.Buffer{})

	args := map[string]interface{}{
		"piece_path":  ".pysar/pieces/does-not-exist",
		"expert_lens": "x",
		"sources":     []map[string]interface{}{validResearchSource("s1")},
	}
	resp := runLines(t, s, toolCallLine(t, 1, "save_research_bundle", args))
	result := toolResultMap(t, resp[0])
	if result["isError"] != true {
		t.Fatalf("expected error for a piece path with no brief.md, got %v", result)
	}
}
