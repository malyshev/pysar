package stagereq

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestRequireMergesAndIdempotent(t *testing.T) {
	dir := t.TempDir()
	brief := "---\npiece: x\nresearch_mode: none\naudience: \"devs\"\n---\n\n# Brief\n"
	if err := os.WriteFile(filepath.Join(dir, "brief.md"), []byte(brief), 0o644); err != nil {
		t.Fatal(err)
	}

	if err := Require(dir, StageResearch); err != nil {
		t.Fatal(err)
	}
	got, err := Load(dir)
	if err != nil {
		t.Fatal(err)
	}
	if len(got) != 1 || got[0] != StageResearch {
		t.Fatalf("got %v", got)
	}

	if err := Require(dir, StageSEO, StageResearch); err != nil {
		t.Fatal(err)
	}
	got, err = Load(dir)
	if err != nil {
		t.Fatal(err)
	}
	if len(got) != 2 {
		t.Fatalf("expected both stages, got %v", got)
	}

	raw, _ := os.ReadFile(filepath.Join(dir, "brief.md"))
	if c := strings.Count(string(raw), "required_stages:"); c != 1 {
		t.Fatalf("expected one required_stages line, got %d:\n%s", c, raw)
	}
	if !strings.Contains(string(raw), `audience: "devs"`) {
		t.Fatalf("brief body/frontmatter must stay intact:\n%s", raw)
	}
}

func TestRequireRejectsUnknown(t *testing.T) {
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "brief.md"), []byte("---\nresearch_mode: none\n---\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := Require(dir, "factcheck"); err == nil {
		t.Fatal("expected unknown stage error")
	}
}

func TestResearchComplete(t *testing.T) {
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "brief.md"), []byte("---\nresearch_mode: targeted\n---\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	ok, err := ResearchComplete(dir)
	if err != nil || ok {
		t.Fatalf("targeted must not complete: ok=%v err=%v", ok, err)
	}
	if err := os.WriteFile(filepath.Join(dir, "brief.md"), []byte("---\nresearch_mode: full\n---\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	ok, err = ResearchComplete(dir)
	if err != nil || !ok {
		t.Fatalf("full must complete: ok=%v err=%v", ok, err)
	}
}

func TestSEOComplete(t *testing.T) {
	dir := t.TempDir()
	ok, err := SEOComplete(dir)
	if err != nil || ok {
		t.Fatalf("missing seo.md: ok=%v err=%v", ok, err)
	}
	if err := os.WriteFile(filepath.Join(dir, "seo.md"), []byte("# seo\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	ok, err = SEOComplete(dir)
	if err != nil || !ok {
		t.Fatalf("seo.md present: ok=%v err=%v", ok, err)
	}
}
