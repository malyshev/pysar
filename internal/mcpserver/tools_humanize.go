package mcpserver

import (
	"encoding/json"
	"fmt"

	"pysar/internal/draft"
	"pysar/internal/editorial"
	"pysar/internal/humanize"
	"pysar/internal/research"
)

var saveHumanizeBundleSchema = map[string]interface{}{
	"type": "object",
	"properties": map[string]interface{}{
		"piece_path": map[string]string{"type": "string", "description": "Path to an existing piece, relative to the project root -- not a slug. Accepts the piece directory itself or a path to a file inside it (e.g. brief.md, what a Claude Code @-reference commonly resolves to, or a raw/ excerpt) -- either form resolves to the same piece directory."},
		"revised_md": map[string]string{"type": "string", "description": "The FULL revised text, written to humanize.md -- never to draft.md, staff-edit.md, or sharpen.md, which stay untouched. Replaces the whole humanize.md file, not a diff/patch. Citation markers ([^shortname]) must still resolve to a real source this piece's research recorded -- never invented, never resolved to a link here."},
		"checks": map[string]interface{}{
			"type":        "array",
			"description": "At least 1 one-line note on what changed and why, e.g. '[hedge-stack] dropped the weaker of two hedges on the compose claim' or '[symmetry] varied the skip-list's third item so it doesn't read as a template'. Record a single 'no changes needed' entry when a check genuinely required none -- omitting checks entirely isn't a completed pass.",
			"items":       map[string]string{"type": "string"},
		},
		"mode": map[string]string{"type": "string", "description": "Optional, informational only: 'delta' (surgical) or 'rewrite' (structural). Folded into the changelog entry, not otherwise enforced."},
	},
	"required": []string{"piece_path", "revised_md", "checks"},
}

func (s *Server) registerSaveHumanizeBundle() {
	s.register(
		tool{
			Name: "save_humanize_bundle",
			Description: "Validate and persist a /ps-humanize pass: writes the revision to humanize.md and appends humanize-changelog.md. Neither draft.md, staff-edit.md, nor sharpen.md is touched -- each stage keeps its own file. " +
				"Reuses draft.ValidateContent for citation integrity (humanize.md's content is still, mechanically, a piece of draft prose) -- no raw URL in prose, every [^shortname] resolved against the piece's actual research output, never accepted on trust. Requires >=1 recorded check -- an edit pass that logged nothing isn't a completed pass. " +
				"This tool has no way to detect or enforce 'sounds human' -- that's the skill's own judgment; it only enforces the same mechanical citation/structure floor every other save_*_bundle tool does. " +
				"A re-run replaces humanize.md wholesale, same as the earlier passes' own writes -- expected, not data loss. Never touches brief.md, outline.md, angles.md, or sources.md. Prefer this over Write/Bash so no extra filesystem permissions are needed.",
			InputSchema: saveHumanizeBundleSchema,
		},
		s.callSaveHumanizeBundle,
	)
}

func (s *Server) callSaveHumanizeBundle(args json.RawMessage) callToolResult {
	var b humanize.Bundle
	if err := json.Unmarshal(args, &b); err != nil {
		return errorResult("invalid arguments: %s", err)
	}

	pieceDir, err := s.resolveAnchoredPass("humanize", b.PiecePath)
	if err != nil {
		return errorResult("%s", err.Error())
	}

	validShortnames, err := research.LoadShortnames(pieceDir)
	if err != nil {
		return errorResult("%s", err)
	}
	if err := humanize.Validate(b, validShortnames); err != nil {
		return errorResult("%s", err)
	}

	words, err := humanize.WriteToPiece(pieceDir, b)
	if err != nil {
		return errorResult("%s", err)
	}

	// Humanize has a three-way "which file did this revise from"
	// ambiguity (sharpen.md if present, else staff-edit.md, else
	// draft.md, per the skill's own instruction) -- same honest-proxy
	// reasoning as sharpen's own revised_from tracking.
	revisedFrom := draft.LatestRevisionFile(pieceDir, "sharpen.md", "staff-edit.md", "draft.md")
	if err := editorial.AppendRunLog(pieceDir, editorial.RunLogEntry{
		Pass:    "humanize",
		Summary: fmt.Sprintf("piece=%s words=%d checks=%d revised_from=%s", b.PiecePath, words, len(b.Checks), revisedFrom),
	}); err != nil {
		return errorResult("append run-log: %s", err)
	}

	return textResult("sounds like you, not the machine: %s/humanize.md (%d words, %d checks recorded). check them before handoff.", b.PiecePath, words, len(b.Checks))
}
