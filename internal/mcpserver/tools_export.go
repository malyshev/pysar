package mcpserver

import (
	"encoding/json"
	"fmt"

	"pysar/internal/editorial"
	"pysar/internal/export"
)

var saveExportBundleSchema = map[string]interface{}{
	"type": "object",
	"properties": map[string]interface{}{
		"piece_path": map[string]string{"type": "string", "description": "Path to an existing piece, relative to the project root -- not a slug. Accepts the piece directory itself or a path to a file inside it (e.g. brief.md, what a Claude Code @-reference commonly resolves to, or a raw/ excerpt) -- either form resolves to the same piece directory."},
		"export_dir": map[string]string{"type": "string", "description": "Optional project-relative directory for this export only. When omitted, uses export_dir from .pysar/project; when that is also empty, writes to the project root. Must stay inside the project (dec-20260812-export-dir-v4-default-plus-override-60d432e0). Prefer omitting this unless the author asked for a one-off landing path."},
	},
	"required": []string{"piece_path"},
}

func (s *Server) registerSaveExportBundle() {
	s.register(
		tool{
			Name: "export_piece_to_root",
			Description: "Copies a piece's most-refined revision (humanize.md if it ran, else seo.md, else sharpen.md, else staff-edit.md, else draft.md) to <export_dir>/<slug>.md — default export_dir is from .pysar/project, else the project root (dec-20260812-export-dir-v4-default-plus-override-60d432e0). Optional export_dir argument overrides the project default for this call only. " +
				"Does not touch draft.md, staff-edit.md, sharpen.md, seo.md, seo-checklist.md, or humanize.md; does not require every earlier stage to have run -- only a draft. Re-running overwrites the previous export of the same piece at the resolved destination, same as every other save_*_bundle's wholesale-replace contract. " +
				"Not a publish ceremony: no cover-image or Status-enum handling (dec-20260718-20a08f83). SEO/discoverability packaging is a real, separate opt-in capability (save_seo_bundle, dec-20260804-e3234e50) -- this tool copies whatever seo.md already produced, it does not generate or edit SEO fields itself. Prefer this over Write/Bash so no extra filesystem permissions are needed. The tool result includes the resolved destination path.",
			InputSchema: saveExportBundleSchema,
		},
		s.callSaveExportBundle,
	)
}

func (s *Server) callSaveExportBundle(args json.RawMessage) callToolResult {
	var b export.Bundle
	if err := json.Unmarshal(args, &b); err != nil {
		return errorResult("invalid arguments: %s", err)
	}
	if err := export.Validate(b); err != nil {
		return errorResult("%s", err)
	}

	pieceDir, err := s.resolveAnchoredPass("export", b.PiecePath)
	if err != nil {
		return errorResult("%s", err.Error())
	}

	effective := export.EffectiveExportDir(export.ReadProjectExportDir(s.baseDir), b.ExportDir)
	destPath, sourceFile, words, err := export.WriteToRoot(s.baseDir, pieceDir, effective)
	if err != nil {
		return errorResult("%s", err)
	}

	if err := editorial.AppendRunLog(pieceDir, editorial.RunLogEntry{
		Pass:    "export",
		Summary: fmt.Sprintf("piece=%s dest=%s source=%s words=%d", b.PiecePath, destPath, sourceFile, words),
	}); err != nil {
		return errorResult("append run-log: %s", err)
	}

	return textResult("exported %s (%d words) to %s.", sourceFile, words, destPath)
}
