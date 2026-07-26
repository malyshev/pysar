package mcpserver

import (
	"encoding/json"
	"fmt"

	"pysar/internal/editorial"
	"pysar/internal/research"
	"pysar/internal/staffedit"
)

var saveStaffEditBundleSchema = map[string]interface{}{
	"type": "object",
	"properties": map[string]interface{}{
		"piece_path": map[string]string{"type": "string", "description": "Path to an existing piece, relative to the project root -- not a slug. Accepts the piece directory itself or a path to a file inside it (e.g. brief.md, what a Claude Code @-reference commonly resolves to, or a raw/ excerpt) -- either form resolves to the same piece directory."},
		"revised_md": map[string]string{"type": "string", "description": "The FULL revised draft, written to staff-edit.md -- never to draft.md, which stays untouched so the two can be compared. Replaces the whole staff-edit.md file, not a diff/patch. Citation markers ([^shortname]) must still resolve to a real source this piece's research recorded -- never invented, never resolved to a link here."},
		"checks": map[string]interface{}{
			"type":        "array",
			"description": "At least 1 one-line note on what changed and why, e.g. '[stakes] named the specific failure mode in paragraph 2' or '[readability] split the three-clause sentence in the compose section'. Record a single 'no changes needed' entry when a check genuinely required none -- omitting checks entirely isn't a completed pass.",
			"items":       map[string]string{"type": "string"},
		},
		"mode": map[string]string{"type": "string", "description": "Optional, informational only: 'delta' (surgical) or 'rewrite' (structural). Folded into the changelog entry, not otherwise enforced."},
	},
	"required": []string{"piece_path", "revised_md", "checks"},
}

func (s *Server) registerSaveStaffEditBundle() {
	s.register(
		tool{
			Name: "save_staff_edit_bundle",
			Description: "Validate and persist a /ps-staff-edit pass: writes the revision to staff-edit.md and appends staff-edit-changelog.md. draft.md is never touched -- it stays the untouched first-draft artifact for comparison (or a future UI showing both). " +
				"Reuses draft.Validate for citation integrity (staff-edit.md's content is still, mechanically, a piece of draft prose) -- no raw URL in prose, every [^shortname] resolved against the piece's actual research output, never accepted on trust. Requires >=1 recorded check -- an edit pass that logged nothing isn't a completed pass. " +
				"A re-run replaces staff-edit.md wholesale, same as /ps-draft's own write for draft.md -- expected, not data loss. Never touches brief.md, outline.md, angles.md, or sources.md. Prefer this over Write/Bash so no extra filesystem permissions are needed.",
			InputSchema: saveStaffEditBundleSchema,
		},
		s.callSaveStaffEditBundle,
	)
}

func (s *Server) callSaveStaffEditBundle(args json.RawMessage) callToolResult {
	var b staffedit.Bundle
	if err := json.Unmarshal(args, &b); err != nil {
		return errorResult("invalid arguments: %s", err)
	}

	pieceDir, err := s.resolveAnchoredPass("staff-edit", b.PiecePath)
	if err != nil {
		return errorResult("%s", err.Error())
	}

	validShortnames, err := research.LoadShortnames(pieceDir)
	if err != nil {
		return errorResult("%s", err)
	}
	if err := staffedit.Validate(b, validShortnames); err != nil {
		return errorResult("%s", err)
	}

	words, err := staffedit.WriteToPiece(pieceDir, b)
	if err != nil {
		return errorResult("%s", err)
	}
	if err := editorial.AppendRunLog(pieceDir, editorial.RunLogEntry{
		Pass:    "staff-edit",
		Summary: fmt.Sprintf("piece=%s words=%d checks=%d", b.PiecePath, words, len(b.Checks)),
	}); err != nil {
		return errorResult("append run-log: %s", err)
	}

	return textResult("staff edit done: %s/staff-edit.md (%d words, %d checks recorded). stakes, spine, and honest scope -- yours to have landed; check them before handoff.", b.PiecePath, words, len(b.Checks))
}
