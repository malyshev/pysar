package mcpserver

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestCheckOnboardingStatusReportsBothOutstandingOnFreshProject(t *testing.T) {
	dir := t.TempDir()
	s := New("pysar", "0.0.0-test", dir, t.TempDir(), nil, &bytes.Buffer{})

	resp := runLines(t, s, toolCallLine(t, 1, "check_onboarding_status", map[string]interface{}{}))
	text := resp[0]["result"].(map[string]interface{})["content"].([]interface{})[0].(map[string]interface{})["text"].(string)

	if !strings.Contains(text, "ps-voice: outstanding") {
		t.Fatalf("expected ps-voice outstanding on a fresh project, got: %q", text)
	}
	if !strings.Contains(text, "ps-style: outstanding") {
		t.Fatalf("expected ps-style outstanding on a fresh project, got: %q", text)
	}
}

func TestCheckOnboardingStatusReportsVoiceDoneAfterSave(t *testing.T) {
	dir := t.TempDir()
	s := New("pysar", "0.0.0-test", dir, t.TempDir(), nil, &bytes.Buffer{})

	saveVoice := toolCallLine(t, 1, "save_voice_profile", map[string]interface{}{
		"tone": "warm", "formality": "neutral", "sentence_length": "varied", "register": "plain",
		"goldens": []map[string]string{{"label": "a", "text": "1"}, {"label": "b", "text": "2"}, {"label": "c", "text": "3"}},
	})
	status := toolCallLine(t, 2, "check_onboarding_status", map[string]interface{}{})
	resp := runLines(t, s, saveVoice, status)

	text := resp[1]["result"].(map[string]interface{})["content"].([]interface{})[0].(map[string]interface{})["text"].(string)
	if !strings.Contains(text, "ps-voice: done") {
		t.Fatalf("expected ps-voice done after a real save, got: %q", text)
	}
	if !strings.Contains(text, "ps-style: outstanding") {
		t.Fatalf("expected ps-style to remain outstanding, unaffected by the voice save, got: %q", text)
	}
}

func TestCheckOnboardingStatusReportsNothingOutstandingWhenBothDone(t *testing.T) {
	dir := t.TempDir()
	s := New("pysar", "0.0.0-test", dir, t.TempDir(), nil, &bytes.Buffer{})

	saveVoice := toolCallLine(t, 1, "save_voice_profile", map[string]interface{}{
		"tone": "warm", "formality": "neutral", "sentence_length": "varied", "register": "plain",
		"goldens": []map[string]string{{"label": "a", "text": "1"}, {"label": "b", "text": "2"}, {"label": "c", "text": "3"}},
	})
	saveStyle := toolCallLine(t, 2, "save_style_profile", map[string]interface{}{
		"rules":   []string{"Put the main point first", "Prefer active voice", "Cut what doesn't earn its place"},
		"goldens": []map[string]string{{"label": "a", "text": "1"}, {"label": "b", "text": "2"}, {"label": "c", "text": "3"}},
	})
	status := toolCallLine(t, 3, "check_onboarding_status", map[string]interface{}{})
	resp := runLines(t, s, saveVoice, saveStyle, status)

	text := resp[2]["result"].(map[string]interface{})["content"].([]interface{})[0].(map[string]interface{})["text"].(string)
	if strings.Contains(text, "outstanding") {
		t.Fatalf("expected nothing outstanding once both are saved, got: %q", text)
	}
	if !strings.Contains(text, "ps-voice: done") || !strings.Contains(text, "ps-style: done") {
		t.Fatalf("expected both passes reported done, got: %q", text)
	}
}

func TestFileExistsHelper(t *testing.T) {
	dir := t.TempDir()
	missing := filepath.Join(dir, "nope.md")
	if fileExists(missing) {
		t.Fatal("expected fileExists to be false for a missing file")
	}
	present := filepath.Join(dir, "here.md")
	if err := os.WriteFile(present, []byte("x"), 0o644); err != nil {
		t.Fatalf("seed file: %v", err)
	}
	if !fileExists(present) {
		t.Fatal("expected fileExists to be true for a file that exists")
	}
}
