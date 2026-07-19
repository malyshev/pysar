package mcpserver

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// toolCallLine builds a single-line (no embedded newlines) tools/call
// request via json.Marshal, since the server frames one message per line.
func toolCallLine(t *testing.T, id int, toolName string, args interface{}) string {
	t.Helper()
	argBytes, err := json.Marshal(args)
	if err != nil {
		t.Fatalf("marshal args: %v", err)
	}
	req := struct {
		JSONRPC string `json:"jsonrpc"`
		Method  string `json:"method"`
		ID      int    `json:"id"`
		Params  struct {
			Name      string          `json:"name"`
			Arguments json.RawMessage `json:"arguments"`
		} `json:"params"`
	}{JSONRPC: "2.0", Method: "tools/call", ID: id}
	req.Params.Name = toolName
	req.Params.Arguments = argBytes

	line, err := json.Marshal(req)
	if err != nil {
		t.Fatalf("marshal request: %v", err)
	}
	return string(line)
}

func TestSaveVoiceProfileWritesCompleteProfile(t *testing.T) {
	dir := t.TempDir()
	s := New("pysar", "0.0.0-test", dir, t.TempDir(), nil, &bytes.Buffer{})

	args := map[string]interface{}{
		"tone":            "warm, direct",
		"formality":       "neutral",
		"sentence_length": "varied",
		"register":        "plain international English",
		"goldens": []map[string]string{
			{"label": "opening", "text": "Bring your take."},
			{"label": "explanation", "text": "The gap is going from idea to shaped piece."},
			{"label": "closing", "text": "You decide when it is ready."},
		},
	}
	resp := runLines(t, s, toolCallLine(t, 1, "save_voice_profile", args))

	result := resp[0]["result"].(map[string]interface{})
	if result["isError"] == true {
		t.Fatalf("expected success, got error result: %v", result)
	}

	voicePath := filepath.Join(dir, ".pysar", "voice.md")
	content, err := os.ReadFile(voicePath)
	if err != nil {
		t.Fatalf("expected .pysar/voice.md written: %v", err)
	}
	if !strings.Contains(string(content), "tone: warm, direct") {
		t.Fatalf("unexpected voice.md content: %s", content)
	}
	if !strings.Contains(string(content), "## Golden: opening") {
		t.Fatalf("expected golden section, got: %s", content)
	}
}

func TestSaveVoiceProfileRejectsIncompleteProfile(t *testing.T) {
	dir := t.TempDir()
	s := New("pysar", "0.0.0-test", dir, t.TempDir(), nil, &bytes.Buffer{})

	// Missing register and only one golden -- deliberately incomplete.
	args := map[string]interface{}{
		"tone":            "warm",
		"formality":       "neutral",
		"sentence_length": "varied",
		"goldens":         []map[string]string{{"label": "opening", "text": "Bring your take."}},
	}
	resp := runLines(t, s, toolCallLine(t, 1, "save_voice_profile", args))

	result := resp[0]["result"].(map[string]interface{})
	if result["isError"] != true {
		t.Fatalf("expected isError true for incomplete profile, got %v", result)
	}
	content := result["content"].([]interface{})[0].(map[string]interface{})["text"].(string)
	if !strings.Contains(content, "register") || !strings.Contains(content, "goldens") {
		t.Fatalf("expected error to name both missing register and insufficient goldens, got: %s", content)
	}

	if _, err := os.Stat(filepath.Join(dir, ".pysar", "voice.md")); err == nil {
		t.Fatal("expected no file written for an incomplete profile")
	}
}

func TestSaveVoiceProfileOverwritesPriorProfile(t *testing.T) {
	dir := t.TempDir()
	s := New("pysar", "0.0.0-test", dir, t.TempDir(), nil, &bytes.Buffer{})

	goldens := []map[string]string{
		{"label": "a", "text": "one"},
		{"label": "b", "text": "two"},
		{"label": "c", "text": "three"},
	}
	first := toolCallLine(t, 1, "save_voice_profile", map[string]interface{}{
		"tone": "first version", "formality": "neutral", "sentence_length": "varied", "register": "plain", "goldens": goldens,
	})
	second := toolCallLine(t, 2, "save_voice_profile", map[string]interface{}{
		"tone": "re-tuned version", "formality": "neutral", "sentence_length": "varied", "register": "plain", "goldens": goldens,
	})
	runLines(t, s, first, second)

	content, err := os.ReadFile(filepath.Join(dir, ".pysar", "voice.md"))
	if err != nil {
		t.Fatalf("read voice.md: %v", err)
	}
	if !strings.Contains(string(content), "re-tuned version") || strings.Contains(string(content), "first version") {
		t.Fatalf("expected re-tune to replace the profile, got: %s", content)
	}
}

func TestSaveVoiceTemplateWritesToHomeNotProject(t *testing.T) {
	projectDir := t.TempDir()
	home := t.TempDir()
	s := New("pysar", "0.0.0-test", projectDir, home, nil, &bytes.Buffer{})

	args := map[string]interface{}{
		"name":            "Measured plain English",
		"tone":            "measured, plain",
		"formality":       "middle-to-formal",
		"sentence_length": "varied",
		"register":        "general and accessible",
		"goldens": []map[string]string{
			{"label": "a", "text": "one"},
			{"label": "b", "text": "two"},
			{"label": "c", "text": "three"},
		},
	}
	resp := runLines(t, s, toolCallLine(t, 1, "save_voice_template", args))

	result := resp[0]["result"].(map[string]interface{})
	if result["isError"] == true {
		t.Fatalf("expected success, got error result: %v", result)
	}

	templatePath := filepath.Join(home, ".pysar", "templates", "voice", "measured-plain-english.md")
	content, err := os.ReadFile(templatePath)
	if err != nil {
		t.Fatalf("expected template written to %s: %v", templatePath, err)
	}
	if !strings.Contains(string(content), "tone: measured, plain") {
		t.Fatalf("unexpected template content: %s", content)
	}

	if _, err := os.Stat(filepath.Join(projectDir, ".pysar", "voice.md")); err == nil {
		t.Fatal("expected save_voice_template to never touch the project's own .pysar/voice.md")
	}
}

func TestSaveVoiceTemplateRejectsMissingName(t *testing.T) {
	s := New("pysar", "0.0.0-test", t.TempDir(), t.TempDir(), nil, &bytes.Buffer{})

	args := map[string]interface{}{
		"tone": "warm", "formality": "neutral", "sentence_length": "varied", "register": "plain",
		"goldens": []map[string]string{
			{"label": "a", "text": "one"}, {"label": "b", "text": "two"}, {"label": "c", "text": "three"},
		},
	}
	resp := runLines(t, s, toolCallLine(t, 1, "save_voice_template", args))

	result := resp[0]["result"].(map[string]interface{})
	if result["isError"] != true {
		t.Fatalf("expected isError true when name is missing, got %v", result)
	}
}

func TestSaveVoiceTemplateRejectsIncompleteProfile(t *testing.T) {
	home := t.TempDir()
	s := New("pysar", "0.0.0-test", t.TempDir(), home, nil, &bytes.Buffer{})

	args := map[string]interface{}{
		"name": "Incomplete", "tone": "warm", "formality": "neutral", "sentence_length": "varied",
		"goldens": []map[string]string{{"label": "a", "text": "one"}},
	}
	resp := runLines(t, s, toolCallLine(t, 1, "save_voice_template", args))

	result := resp[0]["result"].(map[string]interface{})
	if result["isError"] != true {
		t.Fatalf("expected isError true for incomplete profile, got %v", result)
	}
	if _, err := os.Stat(filepath.Join(home, ".pysar", "templates")); err == nil {
		t.Fatal("expected no template directory created for a rejected save")
	}
}

func TestSaveVoiceTemplateOverwritesSameName(t *testing.T) {
	home := t.TempDir()
	s := New("pysar", "0.0.0-test", t.TempDir(), home, nil, &bytes.Buffer{})

	goldens := []map[string]string{
		{"label": "a", "text": "one"}, {"label": "b", "text": "two"}, {"label": "c", "text": "three"},
	}
	first := toolCallLine(t, 1, "save_voice_template", map[string]interface{}{
		"name": "Generic", "tone": "first version", "formality": "neutral", "sentence_length": "varied", "register": "plain", "goldens": goldens,
	})
	second := toolCallLine(t, 2, "save_voice_template", map[string]interface{}{
		"name": "Generic", "tone": "second version", "formality": "neutral", "sentence_length": "varied", "register": "plain", "goldens": goldens,
	})
	runLines(t, s, first, second)

	content, err := os.ReadFile(filepath.Join(home, ".pysar", "templates", "voice", "generic.md"))
	if err != nil {
		t.Fatalf("read template: %v", err)
	}
	if !strings.Contains(string(content), "second version") || strings.Contains(string(content), "first version") {
		t.Fatalf("expected re-saving the same name to overwrite, got: %s", content)
	}
}

// TestListVoiceTemplatesReturnsNothingWhenStoreIsEmpty covers the state
// before any pysar init has ever seeded the built-in default.
func TestListVoiceTemplatesReturnsNothingWhenStoreIsEmpty(t *testing.T) {
	s := New("pysar", "0.0.0-test", t.TempDir(), t.TempDir(), nil, &bytes.Buffer{})

	resp := runLines(t, s, toolCallLine(t, 1, "list_voice_templates", map[string]interface{}{}))
	result := resp[0]["result"].(map[string]interface{})
	if result["isError"] == true {
		t.Fatalf("expected success (empty, not error) for a missing templates dir, got: %v", result)
	}
	text := result["content"].([]interface{})[0].(map[string]interface{})["text"].(string)
	if !strings.Contains(text, "No voice templates") {
		t.Fatalf("expected an explicit empty-store message, got: %q", text)
	}
}

// TestListVoiceTemplatesReturnsFullContentInline is the direct fix for the
// live regression: the skill must never need Glob or Bash to discover
// templates -- one tool call returns every template's full content, so the
// skill can summarize tone/formality/register without any further file
// access.
func TestListVoiceTemplatesReturnsFullContentInline(t *testing.T) {
	home := t.TempDir()
	s := New("pysar", "0.0.0-test", t.TempDir(), home, nil, &bytes.Buffer{})

	goldens := []map[string]string{
		{"label": "a", "text": "one"}, {"label": "b", "text": "two"}, {"label": "c", "text": "three"},
	}
	save := toolCallLine(t, 1, "save_voice_template", map[string]interface{}{
		"name": "Generic", "tone": "measured, plain", "formality": "neutral", "sentence_length": "varied", "register": "plain", "goldens": goldens,
	})
	list := toolCallLine(t, 2, "list_voice_templates", map[string]interface{}{})
	resp := runLines(t, s, save, list)

	listResult := resp[1]["result"].(map[string]interface{})
	if listResult["isError"] == true {
		t.Fatalf("expected success, got error result: %v", listResult)
	}
	text := listResult["content"].([]interface{})[0].(map[string]interface{})["text"].(string)
	if !strings.Contains(text, "generic") {
		t.Fatalf("expected the template's name in the listing, got: %q", text)
	}
	if !strings.Contains(text, "tone: measured, plain") {
		t.Fatalf("expected the template's full content (not just its name) inline, got: %q", text)
	}
}

// TestListVoiceTemplatesExcludesManifest ensures the manifest file
// (dec-20260719-25712417's hash tracking, reused by dec-20260719-3e36577e)
// never leaks into what the skill sees as an available template.
func TestListVoiceTemplatesExcludesManifest(t *testing.T) {
	home := t.TempDir()
	templateDir := filepath.Join(home, ".pysar", "templates", "voice")
	if err := os.MkdirAll(templateDir, 0o755); err != nil {
		t.Fatalf("seed template dir: %v", err)
	}
	if err := os.WriteFile(filepath.Join(templateDir, ".pysar-manifest.json"), []byte(`{"files":{}}`), 0o644); err != nil {
		t.Fatalf("seed manifest: %v", err)
	}

	s := New("pysar", "0.0.0-test", t.TempDir(), home, nil, &bytes.Buffer{})
	resp := runLines(t, s, toolCallLine(t, 1, "list_voice_templates", map[string]interface{}{}))
	result := resp[0]["result"].(map[string]interface{})
	text := result["content"].([]interface{})[0].(map[string]interface{})["text"].(string)
	if !strings.Contains(text, "No voice templates") {
		t.Fatalf("expected the manifest file to be excluded and reported as no templates, got: %q", text)
	}
}

// TestListVoiceTemplatesSurfacesDisplayNameAndSlugSeparately is the direct
// fix for the live regression (dec-20260719-b12539fa): the operator saw
// only the bare slug "generic" offered, never the meaningful display name
// they chose. Both must now be visible, and distinguishable from each other.
func TestListVoiceTemplatesSurfacesDisplayNameAndSlugSeparately(t *testing.T) {
	home := t.TempDir()
	s := New("pysar", "0.0.0-test", t.TempDir(), home, nil, &bytes.Buffer{})

	const displayName = "Measured plain English -- speakable, understated, general audience"
	goldens := []map[string]string{
		{"label": "a", "text": "one"}, {"label": "b", "text": "two"}, {"label": "c", "text": "three"},
	}
	save := toolCallLine(t, 1, "save_voice_template", map[string]interface{}{
		"name": displayName, "slug": "generic", "tone": "measured, plain", "formality": "neutral",
		"sentence_length": "varied", "register": "plain", "goldens": goldens,
	})
	list := toolCallLine(t, 2, "list_voice_templates", map[string]interface{}{})
	resp := runLines(t, s, save, list)

	text := resp[1]["result"].(map[string]interface{})["content"].([]interface{})[0].(map[string]interface{})["text"].(string)
	if !strings.Contains(text, displayName) {
		t.Fatalf("expected the persisted display name in the listing, got: %q", text)
	}
	if !strings.Contains(text, "slug: generic") {
		t.Fatalf("expected the stable slug labeled separately from the name, got: %q", text)
	}
}

// TestSaveVoiceTemplateRenameUpdatesInPlace is dec-20260719-b12539fa's core
// acceptance criterion: passing the same explicit slug under a new name
// updates the existing file, never creates an orphaned duplicate.
func TestSaveVoiceTemplateRenameUpdatesInPlace(t *testing.T) {
	home := t.TempDir()
	s := New("pysar", "0.0.0-test", t.TempDir(), home, nil, &bytes.Buffer{})

	goldens := []map[string]string{
		{"label": "a", "text": "one"}, {"label": "b", "text": "two"}, {"label": "c", "text": "three"},
	}
	first := toolCallLine(t, 1, "save_voice_template", map[string]interface{}{
		"name": "Original Name", "slug": "generic", "tone": "measured", "formality": "neutral",
		"sentence_length": "varied", "register": "plain", "goldens": goldens,
	})
	renamed := toolCallLine(t, 2, "save_voice_template", map[string]interface{}{
		"name": "Renamed", "slug": "generic", "tone": "measured", "formality": "neutral",
		"sentence_length": "varied", "register": "plain", "goldens": goldens,
	})
	runLines(t, s, first, renamed)

	templateDir := filepath.Join(home, ".pysar", "templates", "voice")
	entries, err := os.ReadDir(templateDir)
	if err != nil {
		t.Fatalf("read template dir: %v", err)
	}
	var mdFiles []string
	for _, e := range entries {
		if strings.HasSuffix(e.Name(), ".md") {
			mdFiles = append(mdFiles, e.Name())
		}
	}
	if len(mdFiles) != 1 {
		t.Fatalf("expected exactly one template file after a rename, got: %v", mdFiles)
	}
	if mdFiles[0] != "generic.md" {
		t.Fatalf("expected the rename to keep the same slug/filename, got: %s", mdFiles[0])
	}

	content, err := os.ReadFile(filepath.Join(templateDir, "generic.md"))
	if err != nil {
		t.Fatalf("read renamed template: %v", err)
	}
	if !strings.Contains(string(content), "name: Renamed") || strings.Contains(string(content), "Original Name") {
		t.Fatalf("expected the file to reflect the new name only, got: %q", content)
	}
}

// TestSaveVoiceTemplateSanitizesExplicitSlug guards against a path-traversal
// slug (e.g. "../../etc/evil") escaping ~/.pysar/templates/voice/ via a
// caller-supplied slug -- the auto-derived slug already went through
// onboarding.Slug; the explicit-slug branch must too.
func TestSaveVoiceTemplateSanitizesExplicitSlug(t *testing.T) {
	home := t.TempDir()
	s := New("pysar", "0.0.0-test", t.TempDir(), home, nil, &bytes.Buffer{})

	goldens := []map[string]string{
		{"label": "a", "text": "one"}, {"label": "b", "text": "two"}, {"label": "c", "text": "three"},
	}
	save := toolCallLine(t, 1, "save_voice_template", map[string]interface{}{
		"name": "Escape Attempt", "slug": "../../../../tmp/pysar-escape-test", "tone": "measured",
		"formality": "neutral", "sentence_length": "varied", "register": "plain", "goldens": goldens,
	})
	runLines(t, s, save)

	if _, err := os.Stat(filepath.Join(os.TempDir(), "pysar-escape-test.md")); err == nil {
		t.Fatalf("slug escaped the templates directory onto disk at %s", filepath.Join(os.TempDir(), "pysar-escape-test.md"))
	}

	templateDir := filepath.Join(home, ".pysar", "templates", "voice")
	entries, err := os.ReadDir(templateDir)
	if err != nil {
		t.Fatalf("read template dir: %v", err)
	}
	var mdFiles []string
	for _, e := range entries {
		if strings.HasSuffix(e.Name(), ".md") {
			mdFiles = append(mdFiles, e.Name())
		}
	}
	if len(mdFiles) != 1 {
		t.Fatalf("expected the sanitized slug to land inside the templates dir as one file, got: %v", mdFiles)
	}
}
