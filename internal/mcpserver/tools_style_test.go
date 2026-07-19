package mcpserver

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func completeStyleArgs() map[string]interface{} {
	return map[string]interface{}{
		"rules": []string{
			"Put the main point first",
			"Prefer active voice",
			"Cut what doesn't earn its place",
		},
		"goldens": []map[string]string{
			{"label": "opening", "text": "The deadline moved to Friday."},
			{"label": "structure", "text": "Three things changed. First, ... Second, ... Third, ..."},
			{"label": "cutting", "text": "We tested it. It broke. We fixed it."},
		},
	}
}

func TestSaveStyleProfileWritesCompleteProfile(t *testing.T) {
	dir := t.TempDir()
	s := New("pysar", "0.0.0-test", dir, t.TempDir(), nil, &bytes.Buffer{})

	resp := runLines(t, s, toolCallLine(t, 1, "save_style_profile", completeStyleArgs()))

	result := resp[0]["result"].(map[string]interface{})
	if result["isError"] == true {
		t.Fatalf("expected success, got error result: %v", result)
	}

	stylePath := filepath.Join(dir, ".pysar", "style.md")
	content, err := os.ReadFile(stylePath)
	if err != nil {
		t.Fatalf("expected .pysar/style.md written: %v", err)
	}
	if !strings.Contains(string(content), "# Style Profile") {
		t.Fatalf("expected a Style Profile heading, got: %s", content)
	}
	if !strings.Contains(string(content), "- Put the main point first") {
		t.Fatalf("expected rules rendered as list items, got: %s", content)
	}
	if strings.Contains(string(content), "tone:") {
		t.Fatalf("expected no tone field in a style profile with no tone given, got: %s", content)
	}
}

func TestSaveStyleProfileRejectsTooFewRules(t *testing.T) {
	dir := t.TempDir()
	s := New("pysar", "0.0.0-test", dir, t.TempDir(), nil, &bytes.Buffer{})

	args := completeStyleArgs()
	args["rules"] = []string{"Only one rule"}
	resp := runLines(t, s, toolCallLine(t, 1, "save_style_profile", args))

	result := resp[0]["result"].(map[string]interface{})
	if result["isError"] != true {
		t.Fatalf("expected isError true for too few rules, got %v", result)
	}
	content := result["content"].([]interface{})[0].(map[string]interface{})["text"].(string)
	if !strings.Contains(content, "rules") {
		t.Fatalf("expected the error to name 'rules' as missing, got: %s", content)
	}

	if _, err := os.Stat(filepath.Join(dir, ".pysar", "style.md")); err == nil {
		t.Fatal("expected no file written for an incomplete style profile")
	}
}

func TestSaveStyleProfileDoesNotTouchVoiceProfile(t *testing.T) {
	dir := t.TempDir()
	s := New("pysar", "0.0.0-test", dir, t.TempDir(), nil, &bytes.Buffer{})

	runLines(t, s, toolCallLine(t, 1, "save_style_profile", completeStyleArgs()))

	if _, err := os.Stat(filepath.Join(dir, ".pysar", "voice.md")); err == nil {
		t.Fatal("expected save_style_profile to never write .pysar/voice.md")
	}
}

func TestSaveStyleTemplateWritesToStyleSubdirNotVoice(t *testing.T) {
	home := t.TempDir()
	s := New("pysar", "0.0.0-test", t.TempDir(), home, nil, &bytes.Buffer{})

	args := completeStyleArgs()
	args["name"] = "Plain English -- GOV.UK style"
	resp := runLines(t, s, toolCallLine(t, 1, "save_style_template", args))

	result := resp[0]["result"].(map[string]interface{})
	if result["isError"] == true {
		t.Fatalf("expected success, got error result: %v", result)
	}

	stylePath := filepath.Join(home, ".pysar", "templates", "style", "plain-english-gov-uk-style.md")
	if _, err := os.Stat(stylePath); err != nil {
		t.Fatalf("expected style template at %s: %v", stylePath, err)
	}
	if _, err := os.Stat(filepath.Join(home, ".pysar", "templates", "voice", "plain-english-gov-uk-style.md")); err == nil {
		t.Fatal("expected the style template to never land under templates/voice/")
	}
}

func TestListStyleTemplatesDoesNotReturnVoiceTemplates(t *testing.T) {
	home := t.TempDir()
	s := New("pysar", "0.0.0-test", t.TempDir(), home, nil, &bytes.Buffer{})

	voiceArgs := map[string]interface{}{
		"name": "A Voice Template", "tone": "warm", "formality": "neutral", "sentence_length": "varied", "register": "plain",
		"goldens": []map[string]string{{"label": "a", "text": "1"}, {"label": "b", "text": "2"}, {"label": "c", "text": "3"}},
	}
	styleArgs := completeStyleArgs()
	styleArgs["name"] = "A Style Template"

	resp := runLines(t, s,
		toolCallLine(t, 1, "save_voice_template", voiceArgs),
		toolCallLine(t, 2, "save_style_template", styleArgs),
		toolCallLine(t, 3, "list_style_templates", map[string]interface{}{}),
	)
	text := resp[2]["result"].(map[string]interface{})["content"].([]interface{})[0].(map[string]interface{})["text"].(string)
	if !strings.Contains(text, "A Style Template") {
		t.Fatalf("expected the style template listed, got: %q", text)
	}
	if strings.Contains(text, "A Voice Template") {
		t.Fatalf("expected list_style_templates to never include a voice template, got: %q", text)
	}
}

func TestSaveStyleTemplateRenameUpdatesInPlace(t *testing.T) {
	home := t.TempDir()
	s := New("pysar", "0.0.0-test", t.TempDir(), home, nil, &bytes.Buffer{})

	first := completeStyleArgs()
	first["name"] = "Original Style Name"
	first["slug"] = "generic"
	renamed := completeStyleArgs()
	renamed["name"] = "Renamed Style"
	renamed["slug"] = "generic"

	runLines(t, s,
		toolCallLine(t, 1, "save_style_template", first),
		toolCallLine(t, 2, "save_style_template", renamed),
	)

	templateDir := filepath.Join(home, ".pysar", "templates", "style")
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
	if len(mdFiles) != 1 || mdFiles[0] != "generic.md" {
		t.Fatalf("expected exactly one file (generic.md) after a rename, got: %v", mdFiles)
	}

	content, err := os.ReadFile(filepath.Join(templateDir, "generic.md"))
	if err != nil {
		t.Fatalf("read renamed template: %v", err)
	}
	if !strings.Contains(string(content), "name: Renamed Style") || strings.Contains(string(content), "Original Style Name") {
		t.Fatalf("expected the file to reflect the new name only, got: %q", content)
	}
}
