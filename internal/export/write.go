package export

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"pysar/internal/draft"
)

const projectManifestName = "project"

// WriteToRoot copies the piece's most-refined revision (draft.RevisionPriority
// order: humanize.md if it ran, else seo.md, else sharpen.md, else
// staff-edit.md, else draft.md) to <exportDir>/<slug>.md under the project,
// overwriting any previous export of the same piece. exportDir is the
// effective destination (MCP override or project default); empty means the
// project root. Returns the destination path, which source file was copied,
// and its word count.
func WriteToRoot(projectRoot, pieceDir, exportDir string) (destPath, sourceFile string, words int, err error) {
	sourceFile = draft.LatestRevisionFile(pieceDir, draft.RevisionPriority...)

	content, err := os.ReadFile(filepath.Join(pieceDir, sourceFile))
	if err != nil {
		return "", "", 0, fmt.Errorf("read %s: %w", sourceFile, err)
	}

	destDir, err := ResolveExportDirectory(projectRoot, exportDir)
	if err != nil {
		return "", "", 0, err
	}
	if err := os.MkdirAll(destDir, 0o755); err != nil {
		return "", "", 0, fmt.Errorf("mkdir %s: %w", destDir, err)
	}

	slug := filepath.Base(pieceDir)
	destPath = filepath.Join(destDir, slug+".md")
	if err := os.WriteFile(destPath, content, 0o644); err != nil {
		return "", "", 0, fmt.Errorf("write %s: %w", destPath, err)
	}

	words = draft.WordCount(string(content))
	return destPath, sourceFile, words, nil
}

// ResolveExportDirectory joins a project-relative exportDir under projectRoot
// (empty → project root) and rejects paths that escape the project.
func ResolveExportDirectory(projectRoot, exportDir string) (string, error) {
	root, err := filepath.Abs(filepath.Clean(projectRoot))
	if err != nil {
		return "", fmt.Errorf("resolve project root: %w", err)
	}

	trimmed := strings.TrimSpace(exportDir)
	if trimmed == "" {
		return root, nil
	}
	if filepath.IsAbs(trimmed) {
		return "", fmt.Errorf("export_dir must be project-relative, got absolute path %q", trimmed)
	}

	cleaned := filepath.Clean(filepath.FromSlash(trimmed))
	if cleaned == "." {
		return root, nil
	}
	if cleaned == ".." || strings.HasPrefix(cleaned, ".."+string(os.PathSeparator)) {
		return "", fmt.Errorf("export_dir escapes project root: %q", trimmed)
	}

	dest := filepath.Clean(filepath.Join(root, cleaned))
	rel, err := filepath.Rel(root, dest)
	if err != nil || rel == ".." || strings.HasPrefix(rel, ".."+string(os.PathSeparator)) {
		return "", fmt.Errorf("export_dir escapes project root: %q", trimmed)
	}
	return dest, nil
}

// ValidateRelativeExportDir checks an init-time --export-dir value before it
// is written to .pysar/project.
func ValidateRelativeExportDir(exportDir string) error {
	trimmed := strings.TrimSpace(exportDir)
	if trimmed == "" {
		return nil
	}
	if filepath.IsAbs(trimmed) {
		return fmt.Errorf("export_dir must be project-relative, got absolute path %q", trimmed)
	}
	cleaned := filepath.Clean(filepath.FromSlash(trimmed))
	if cleaned == ".." || strings.HasPrefix(cleaned, ".."+string(os.PathSeparator)) {
		return fmt.Errorf("export_dir escapes project root: %q", trimmed)
	}
	return nil
}

// EffectiveExportDir applies V4 precedence: non-empty override wins, else
// projectDefault (from .pysar/project), else project root (empty string).
func EffectiveExportDir(projectDefault, override string) string {
	if strings.TrimSpace(override) != "" {
		return strings.TrimSpace(override)
	}
	return strings.TrimSpace(projectDefault)
}

type projectManifestExport struct {
	ExportDir string `json:"export_dir"`
}

// ReadProjectExportDir returns export_dir from .pysar/project, or "" when the
// file is missing, unreadable, or the field is absent/empty.
func ReadProjectExportDir(projectRoot string) string {
	path := filepath.Join(projectRoot, ".pysar", projectManifestName)
	data, err := os.ReadFile(path)
	if err != nil {
		return ""
	}
	var m projectManifestExport
	if err := json.Unmarshal(data, &m); err != nil {
		return ""
	}
	return strings.TrimSpace(m.ExportDir)
}
