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
	},
	"required": []string{"piece_path"},
}

func (s *Server) registerSaveExportBundle() {
	s.register(
		tool{
			Name: "export_piece_to_root",
			Description: "Copies a piece's most-refined revision (humanize.md if it ran, else sharpen.md, else staff-edit.md, else draft.md) to <project_root>/<slug>.md -- the one step that writes outside the piece's own .pysar/pieces/ directory. " +
				"Does not touch draft.md, staff-edit.md, sharpen.md, or humanize.md; does not require every earlier stage to have run -- only a draft. Re-running overwrites the previous export of the same piece, same as every other save_*_bundle's wholesale-replace contract. " +
				"Not a publish ceremony: no SEO/tags/cover-image/Status handling -- this project deliberately excludes that (see CLAUDE.md). Prefer this over Write/Bash so no extra filesystem permissions are needed.",
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

	destPath, sourceFile, words, err := export.WriteToRoot(s.baseDir, pieceDir)
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
