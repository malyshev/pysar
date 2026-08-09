package main

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

// cursorPluginRoot is the canonical Cursor Plugin package
// (dec-20260809-cursor-marketplace-v1-dual-discovery-8b748a7a).
const cursorPluginRoot = "plugins/pysar"

func TestCursorPluginManifestAndMCP(t *testing.T) {
	root := repoRoot(t)
	manifestPath := filepath.Join(root, cursorPluginRoot, ".cursor-plugin", "plugin.json")
	raw, err := os.ReadFile(manifestPath)
	if err != nil {
		t.Fatalf("read plugin manifest: %v", err)
	}
	var manifest struct {
		Name        string `json:"name"`
		Description string `json:"description"`
		Version     string `json:"version"`
		License     string `json:"license"`
		Logo        string `json:"logo"`
		Author      struct {
			Name string `json:"name"`
		} `json:"author"`
	}
	if err := json.Unmarshal(raw, &manifest); err != nil {
		t.Fatalf("parse plugin manifest: %v", err)
	}
	if manifest.Name != "pysar" {
		t.Fatalf("plugin name = %q, want pysar", manifest.Name)
	}
	if manifest.Description == "" || manifest.Version == "" || manifest.License == "" || manifest.Author.Name == "" {
		t.Fatalf("plugin manifest missing required metadata: %+v", manifest)
	}
	if manifest.Logo == "" {
		t.Fatal("plugin manifest missing logo")
	}
	logoPath := filepath.Join(root, cursorPluginRoot, filepath.FromSlash(manifest.Logo))
	if _, err := os.Stat(logoPath); err != nil {
		t.Fatalf("logo path %s: %v", logoPath, err)
	}

	pluginMCP, err := os.ReadFile(filepath.Join(root, cursorPluginRoot, "mcp.json"))
	if err != nil {
		t.Fatalf("read plugin mcp.json: %v", err)
	}
	assetMCP, err := os.ReadFile(filepath.Join(root, "cmd/pysar/assets/cursor/mcp.json"))
	if err != nil {
		t.Fatalf("read assets cursor mcp.json: %v", err)
	}
	if string(pluginMCP) != string(assetMCP) {
		t.Fatal("plugins/pysar/mcp.json must match cmd/pysar/assets/cursor/mcp.json (same spawn contract)")
	}
	if !json.Valid(pluginMCP) {
		t.Fatal("plugin mcp.json is not valid JSON")
	}
}

func TestCursorPluginSkillsMatchCorpus(t *testing.T) {
	root := repoRoot(t)
	src := filepath.Join(root, "cmd/pysar/assets/skills")
	dst := filepath.Join(root, cursorPluginRoot, "skills")
	err := filepath.WalkDir(src, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		rel, err := filepath.Rel(src, path)
		if err != nil {
			return err
		}
		target := filepath.Join(dst, rel)
		if d.IsDir() {
			st, err := os.Stat(target)
			if err != nil || !st.IsDir() {
				t.Errorf("missing plugin skill dir %s (run scripts/sync-cursor-plugin-skills.sh)", rel)
			}
			return nil
		}
		want, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		got, err := os.ReadFile(target)
		if err != nil {
			t.Errorf("missing plugin skill file %s (run scripts/sync-cursor-plugin-skills.sh): %v", rel, err)
			return nil
		}
		if string(got) != string(want) {
			t.Errorf("plugin skill drift at %s — edit assets/skills and re-sync", rel)
		}
		return nil
	})
	if err != nil {
		t.Fatalf("walk skill corpus: %v", err)
	}
}

func TestMarketplaceManifestPointsAtPlugin(t *testing.T) {
	root := repoRoot(t)
	raw, err := os.ReadFile(filepath.Join(root, ".cursor-plugin/marketplace.json"))
	if err != nil {
		t.Fatalf("read marketplace.json: %v", err)
	}
	var m struct {
		Name    string `json:"name"`
		Plugins []struct {
			Name   string `json:"name"`
			Source string `json:"source"`
		} `json:"plugins"`
		Metadata struct {
			PluginRoot string `json:"pluginRoot"`
		} `json:"metadata"`
	}
	if err := json.Unmarshal(raw, &m); err != nil {
		t.Fatalf("parse marketplace.json: %v", err)
	}
	if m.Name != "pysar" || m.Metadata.PluginRoot != "plugins" {
		t.Fatalf("marketplace identity unexpected: name=%q pluginRoot=%q", m.Name, m.Metadata.PluginRoot)
	}
	if len(m.Plugins) != 1 || m.Plugins[0].Name != "pysar" || m.Plugins[0].Source != "pysar" {
		t.Fatalf("marketplace plugins = %+v, want single pysar source", m.Plugins)
	}
	if _, err := os.Stat(filepath.Join(root, "plugins", m.Plugins[0].Source, ".cursor-plugin", "plugin.json")); err != nil {
		t.Fatalf("marketplace source missing plugin.json: %v", err)
	}
}

func repoRoot(t *testing.T) string {
	t.Helper()
	wd, err := os.Getwd()
	if err != nil {
		t.Fatalf("getcwd: %v", err)
	}
	// Tests run from cmd/pysar; walk up until plugins/pysar exists.
	dir := wd
	for {
		if _, err := os.Stat(filepath.Join(dir, cursorPluginRoot, ".cursor-plugin", "plugin.json")); err == nil {
			return dir
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			t.Fatalf("repo root not found from %s", wd)
		}
		dir = parent
	}
}
