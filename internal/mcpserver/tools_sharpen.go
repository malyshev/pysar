package mcpserver

import (
	"encoding/json"
	"fmt"

	"pysar/internal/editorial"
	"pysar/internal/research"
	"pysar/internal/sharpen"
)

var saveSharpenBundleSchema = map[string]interface{}{
	"type": "object",
	"properties": map[string]interface{}{
		"piece_path": map[string]string{"type": "string", "description": "Path to an existing piece, relative to the project root -- not a slug. Accepts the piece directory itself or a path to a file inside it (e.g. brief.md, what a Claude Code @-reference commonly resolves to, or a raw/ excerpt) -- either form resolves to the same piece directory."},
		"revised_md": map[string]string{"type": "string", "description": "The FULL revised text, written to sharpen.md -- never to draft.md or staff-edit.md, which stay untouched. Replaces the whole sharpen.md file, not a diff/patch. Citation markers ([^shortname]) must still resolve to a real source this piece's research recorded -- never invented, never resolved to a link here."},
		"checks": map[string]interface{}{
			"type":        "array",
			"description": "At least 1 one-line note on what changed and why, e.g. '[opener] tightened the hook so paragraph 1 promises a specific payoff' or '[elevate] promoted the contrarian finding on X to its own heading'. Record a single 'no changes needed' entry when a check genuinely required none -- omitting checks entirely isn't a completed pass.",
			"items":       map[string]string{"type": "string"},
		},
		"mode": map[string]string{"type": "string", "description": "Optional, informational only: 'delta' (surgical) or 'rewrite' (structural). Folded into the changelog entry, not otherwise enforced."},
	},
	"required": []string{"piece_path", "revised_md", "checks"},
}

func (s *Server) registerSaveSharpenBundle() {
	s.register(
		tool{
			Name: "save_sharpen_bundle",
			Description: "Validate and persist a /ps-sharpen pass: writes the revision to sharpen.md and appends sharpen-changelog.md. Neither draft.md nor staff-edit.md is touched -- each stage keeps its own file. " +
				"Reuses draft.Validate for citation integrity (sharpen.md's content is still, mechanically, a piece of draft prose) -- no raw URL in prose, every [^shortname] resolved against the piece's actual research output, never accepted on trust. Requires >=1 recorded check -- an edit pass that logged nothing isn't a completed pass. " +
				"A re-run replaces sharpen.md wholesale, same as the earlier passes' own writes -- expected, not data loss. Never touches brief.md, outline.md, angles.md, or sources.md. Prefer this over Write/Bash so no extra filesystem permissions are needed.",
			InputSchema: saveSharpenBundleSchema,
		},
		s.callSaveSharpenBundle,
	)
}

func (s *Server) callSaveSharpenBundle(args json.RawMessage) callToolResult {
	var b sharpen.Bundle
	if err := json.Unmarshal(args, &b); err != nil {
		return errorResult("invalid arguments: %s", err)
	}

	pieceDir, err := s.resolveAnchoredPass("sharpen", b.PiecePath)
	if err != nil {
		return errorResult("%s", err.Error())
	}

	validShortnames, err := research.LoadShortnames(pieceDir)
	if err != nil {
		return errorResult("%s", err)
	}
	if err := sharpen.Validate(b, validShortnames); err != nil {
		return errorResult("%s", err)
	}

	words, err := sharpen.WriteToPiece(pieceDir, b)
	if err != nil {
		return errorResult("%s", err)
	}
	if err := editorial.AppendRunLog(pieceDir, editorial.RunLogEntry{
		Pass:    "sharpen",
		Summary: fmt.Sprintf("piece=%s words=%d checks=%d", b.PiecePath, words, len(b.Checks)),
	}); err != nil {
		return errorResult("append run-log: %s", err)
	}

	return textResult("sharpen done: %s/sharpen.md (%d words, %d checks recorded). opening and arc -- yours to have landed; check them before handoff.", b.PiecePath, words, len(b.Checks))
}
