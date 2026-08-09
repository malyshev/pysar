package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
)

// cursorPysarCommand is the portable Cursor spawn path for project mcp.json
// (dec-20260809-cursor-cold-path-phased-v1v5-then-v2-de2fc11a). Bare "pysar"
// fails when Dock-launched Cursor omits ~/.local/bin from PATH.
const cursorPysarCommand = `${userHome}/.local/bin/pysar`

type cursorMCPServerConfig struct {
	Type    string            `json:"type"`
	Command string            `json:"command"`
	Args    []string          `json:"args"`
	Env     map[string]string `json:"env"`
}

type cursorMCPFile struct {
	MCPServers map[string]cursorMCPServerConfig `json:"mcpServers"`
}

func cursorPysarMCPServerConfig(command string) cursorMCPServerConfig {
	return cursorMCPServerConfig{
		Type:    "stdio",
		Command: command,
		Args:    []string{"serve"},
		Env: map[string]string{
			"PYSAR_PROJECT_ROOT": "${workspaceFolder}",
		},
	}
}

// resolveCursorPysarCommand prefers a real on-disk binary for user MCP / deeplinks
// (absolute path Cursor can spawn), then LookPath, then the portable ${userHome} form.
func resolveCursorPysarCommand() string {
	if home, err := homeDir(); err == nil {
		candidate := filepath.Join(home, ".local", "bin", "pysar")
		if st, err := os.Stat(candidate); err == nil && !st.IsDir() {
			return candidate
		}
	}
	if p, err := exec.LookPath("pysar"); err == nil {
		return p
	}
	return cursorPysarCommand
}

// ensureUserCursorMCP merges the pysar server into ~/.cursor/mcp.json so it
// appears on Cursor's Customize → MCPs Connected surface (user scope). Project
// .cursor/mcp.json alone stays disconnected/invisible on Cursor 3.15 without
// that surface — deeplinks are unreliable.
func ensureUserCursorMCP(home string) (wrote bool, err error) {
	path := filepath.Join(home, ".cursor", "mcp.json")
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return false, err
	}

	want := cursorPysarMCPServerConfig(resolveCursorPysarCommand())
	file := cursorMCPFile{MCPServers: map[string]cursorMCPServerConfig{}}

	existing, readErr := os.ReadFile(path)
	if readErr == nil {
		if err := json.Unmarshal(existing, &file); err != nil {
			return false, fmt.Errorf("parse %s: %w", path, err)
		}
		if file.MCPServers == nil {
			file.MCPServers = map[string]cursorMCPServerConfig{}
		}
		if cur, ok := file.MCPServers["pysar"]; ok && cursorMCPServerConfigEqual(cur, want) {
			return false, nil
		}
	} else if !os.IsNotExist(readErr) {
		return false, readErr
	}

	file.MCPServers["pysar"] = want
	out, err := json.MarshalIndent(file, "", "  ")
	if err != nil {
		return false, err
	}
	out = append(out, '\n')
	if err := os.WriteFile(path, out, 0o644); err != nil {
		return false, err
	}
	return true, nil
}

func cursorMCPServerConfigEqual(a, b cursorMCPServerConfig) bool {
	if a.Type != b.Type || a.Command != b.Command || len(a.Args) != len(b.Args) {
		return false
	}
	for i := range a.Args {
		if a.Args[i] != b.Args[i] {
			return false
		}
	}
	if len(a.Env) != len(b.Env) {
		return false
	}
	for k, v := range a.Env {
		if b.Env[k] != v {
			return false
		}
	}
	return true
}

// cursorMCPInstallDeeplink builds a Cursor MCP install deeplink
// (cursor.com/docs/mcp/install-links). Fallback when user mcp.json cannot be written.
func cursorMCPInstallDeeplink() (string, error) {
	raw, err := json.Marshal(cursorPysarMCPServerConfig(resolveCursorPysarCommand()))
	if err != nil {
		return "", err
	}
	encoded := base64.StdEncoding.EncodeToString(raw)
	return "cursor://anysphere.cursor-deeplink/mcp/install?name=" +
		url.QueryEscape("pysar") +
		"&config=" + url.QueryEscape(encoded), nil
}

func printCursorMCPNextStep(userMCPWrote bool) {
	// ensureUserCursorMCP returns wrote=false when ~/.cursor/mcp.json already
	// has the correct pysar entry. That is success, not a failed registration —
	// do not open the Cursor MCP install deeplink (that was unexpected on re-init).
	if !userMCPWrote {
		return
	}
	fmt.Println("pysar init: registered Pysar MCP for Cursor — reload MCP or restart Cursor")
}

// printCursorPluginNextStep tells authors where /ps skills come from after
// init stopped dual-writing ~/.cursor/skills
// (dec-20260809-cursor-marketplace-v1-dual-discovery-8b748a7a).
func printCursorPluginNextStep() {
	fmt.Println("pysar init: install the Pysar Cursor plugin for /ps skills (Marketplace, getpysar.com Install in Cursor, or symlink plugins/pysar → ~/.cursor/plugins/local/pysar)")
}
