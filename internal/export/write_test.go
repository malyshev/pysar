package export

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestResolveExportDirectoryRootDefault(t *testing.T) {
	root := t.TempDir()
	got, err := ResolveExportDirectory(root, "")
	if err != nil {
		t.Fatalf("ResolveExportDirectory: %v", err)
	}
	want, err := filepath.Abs(root)
	if err != nil {
		t.Fatalf("abs: %v", err)
	}
	if got != want {
		t.Fatalf("got %q, want %q", got, want)
	}
}

func TestResolveExportDirectoryConfigured(t *testing.T) {
	root := t.TempDir()
	got, err := ResolveExportDirectory(root, "published")
	if err != nil {
		t.Fatalf("ResolveExportDirectory: %v", err)
	}
	want := filepath.Join(mustAbs(t, root), "published")
	if got != want {
		t.Fatalf("got %q, want %q", got, want)
	}
}

func TestResolveExportDirectoryRejectsEscape(t *testing.T) {
	root := t.TempDir()
	for _, dir := range []string{"..", "../outside", "foo/../../outside"} {
		if _, err := ResolveExportDirectory(root, dir); err == nil {
			t.Fatalf("expected escape error for %q", dir)
		}
	}
	if _, err := ResolveExportDirectory(root, "/tmp/abs"); err == nil {
		t.Fatal("expected absolute path error")
	}
}

func TestEffectiveExportDirPrecedence(t *testing.T) {
	if got := EffectiveExportDir("published", "outbox"); got != "outbox" {
		t.Fatalf("override should win: got %q", got)
	}
	if got := EffectiveExportDir("published", ""); got != "published" {
		t.Fatalf("project default: got %q", got)
	}
	if got := EffectiveExportDir("", ""); got != "" {
		t.Fatalf("empty: got %q", got)
	}
}

func TestWriteToRootUsesExportDirAndMkdirAll(t *testing.T) {
	root := t.TempDir()
	pieceDir := filepath.Join(root, ".pysar", "pieces", "my-piece")
	if err := os.MkdirAll(pieceDir, 0o755); err != nil {
		t.Fatal(err)
	}
	draftPath := filepath.Join(pieceDir, "draft.md")
	if err := os.WriteFile(draftPath, []byte("# Title\n\nbody\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	dest, source, words, err := WriteToRoot(root, pieceDir, "published/out")
	if err != nil {
		t.Fatalf("WriteToRoot: %v", err)
	}
	if source != "draft.md" {
		t.Fatalf("source: got %q", source)
	}
	if words < 1 {
		t.Fatalf("expected positive word count, got %d", words)
	}
	want := filepath.Join(mustAbs(t, root), "published", "out", "my-piece.md")
	if dest != want {
		t.Fatalf("dest: got %q, want %q", dest, want)
	}
	got, err := os.ReadFile(want)
	if err != nil {
		t.Fatalf("read export: %v", err)
	}
	if !strings.Contains(string(got), "# Title") {
		t.Fatalf("unexpected content: %s", got)
	}
}

func TestWriteToRootRejectsEscape(t *testing.T) {
	root := t.TempDir()
	pieceDir := filepath.Join(root, ".pysar", "pieces", "my-piece")
	if err := os.MkdirAll(pieceDir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(pieceDir, "draft.md"), []byte("x\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if _, _, _, err := WriteToRoot(root, pieceDir, ".."); err == nil {
		t.Fatal("expected escape error")
	}
}

func TestReadProjectExportDir(t *testing.T) {
	root := t.TempDir()
	if got := ReadProjectExportDir(root); got != "" {
		t.Fatalf("missing manifest: got %q", got)
	}
	pysar := filepath.Join(root, ".pysar")
	if err := os.MkdirAll(pysar, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(pysar, "project"), []byte(`{"schema_version":1,"host":"claude","export_dir":"published"}`+"\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if got := ReadProjectExportDir(root); got != "published" {
		t.Fatalf("got %q, want published", got)
	}
}

func mustAbs(t *testing.T, p string) string {
	t.Helper()
	abs, err := filepath.Abs(p)
	if err != nil {
		t.Fatal(err)
	}
	return abs
}
