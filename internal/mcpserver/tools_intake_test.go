package mcpserver

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestSaveIntakeBundleEndToEnd(t *testing.T) {
	dir := t.TempDir()
	s := New("pysar", "test", dir, t.TempDir(), nil, &bytes.Buffer{})

	args := map[string]interface{}{
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
	// Re-running intake with the identical idea/name is never blocked or
	// surfaced as a decision to the author -- each call silently allocates
	// its own uniquely-suffixed piece directory.
	first := toolCallLine(t, 1, "save_intake_bundle", args)
	second := toolCallLine(t, 2, "save_intake_bundle", args)
	resp := runLines(t, s, first, second)
	if len(resp) != 2 {
		t.Fatalf("expected 2 responses, got %d", len(resp))
	}
	result := toolResultMap(t, resp[0])
	if result["isError"] == true {
		t.Fatalf("save_intake_bundle error: %v", result)
	}
	result2 := toolResultMap(t, resp[1])
	if result2["isError"] == true {
		t.Fatalf("expected second save (same idea, same name) to succeed with its own directory, got error: %v", result2)
	}

	piecesRoot := filepath.Join(dir, ".pysar", "pieces")
	entries, err := os.ReadDir(piecesRoot)
	if err != nil {
		t.Fatalf("read pieces dir: %v", err)
	}
	var matches []string
	for _, e := range entries {
		if strings.HasPrefix(e.Name(), "ai-safety-for-parents-") {
			matches = append(matches, e.Name())
		}
	}
	if len(matches) != 2 {
		t.Fatalf("expected 2 distinct suffixed piece dirs for the same idea, got %v", matches)
	}
	if matches[0] == matches[1] {
		t.Fatalf("expected two different directory names, got the same twice: %s", matches[0])
	}
	// Directory names must never fall back to the bare, un-suffixed name.
	if _, err := os.Stat(filepath.Join(piecesRoot, "ai-safety-for-parents")); err == nil {
		t.Fatalf("bare un-suffixed piece dir should never exist")
	}

	for _, name := range matches {
		piece := filepath.Join(piecesRoot, name)
		if _, err := os.Stat(filepath.Join(piece, "STORAGE.md")); err != nil {
			t.Fatalf("missing STORAGE.md: %v", err)
		}
		if _, err := os.Stat(filepath.Join(piece, "run-log.jsonl")); err != nil {
			t.Fatalf("missing run-log: %v", err)
		}
		logb, _ := os.ReadFile(filepath.Join(piece, "run-log.jsonl"))
		if !strings.Contains(string(logb), `"pass":"intake"`) {
			t.Fatalf("run-log missing intake: %s", logb)
		}
		brief, _ := os.ReadFile(filepath.Join(piece, "brief.md"))
		if strings.Contains(string(brief), "slug:") {
			t.Fatalf("brief.md must not use the word 'slug': %s", brief)
		}
	}
}

func TestReadAuthorDefaultsGeneralist(t *testing.T) {
	dir := t.TempDir()
	out := &bytes.Buffer{}
	s := New("pysar", "test", dir, t.TempDir(), nil, out)
	resp := runLines(t, s, toolCallLine(t, 1, "read_author_defaults", map[string]interface{}{}))
	result := toolResultMap(t, resp[0])
	if result["isError"] == true {
		t.Fatalf("error: %v", result)
	}
	text := toolText(t, result)
	if !strings.Contains(text, "generalist-fallback") {
		t.Fatalf("expected generalist-fallback, got %s", text)
	}
}

func toolResultMap(t *testing.T, resp map[string]interface{}) map[string]interface{} {
	t.Helper()
	result, ok := resp["result"].(map[string]interface{})
	if !ok {
		t.Fatalf("expected result object: %v", resp)
	}
	return result
}

func toolText(t *testing.T, result map[string]interface{}) string {
	t.Helper()
	content, ok := result["content"].([]interface{})
	if !ok || len(content) == 0 {
		t.Fatalf("no content: %v", result)
	}
	item := content[0].(map[string]interface{})
	return item["text"].(string)
}
