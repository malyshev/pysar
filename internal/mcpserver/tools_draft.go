package mcpserver

import (
	"encoding/json"
	"fmt"

	"pysar/internal/draft"
	"pysar/internal/editorial"
	"pysar/internal/research"
)

var saveDraftBundleSchema = map[string]interface{}{
	"type": "object",
	"properties": map[string]interface{}{
		"piece_path": map[string]string{"type": "string", "description": "Path to an existing piece, relative to the project root -- not a slug. Accepts the piece directory itself or a path to a file inside it (e.g. brief.md, what a Claude Code @-reference commonly resolves to, or a raw/ excerpt) -- either form resolves to the same piece directory."},
		"draft_md":   map[string]string{"type": "string", "description": "The full draft: channel-agnostic prose, not a finished web/print artifact. Cite via [^shortname] markers matching a real source this piece's research recorded -- never a raw URL, and never resolve a marker to an inline link here -- that's a later, optional, channel-specific packaging step this tool doesn't own. Heading/title shape (e.g. one '# ' title + italic subtitle for a blog-shaped piece) is not enforced here -- it depends on the piece's target surface, so it's the skill's own judgment call, not this tool's."},
		"notes":      map[string]string{"type": "string", "description": "Optional one-line changelog note, e.g. 'first full draft' or 'revision after staff-edit pass'. Defaults to 'draft'."},
	},
	"required": []string{"piece_path", "draft_md"},
}

func (s *Server) registerSaveDraftBundle() {
	s.register(
		tool{
			Name: "save_draft_bundle",
			Description: "Validate and persist a /ps-draft pass: writes draft.md (channel-agnostic prose, not a finished web/print artifact) to an existing piece and appends draft-changelog.md. " +
				"Enforces citation integrity server-side -- no raw URL in prose, and every [^shortname] citation resolved against the piece's actual research output (research.LoadShortnames), never accepted on trust. Does not resolve citation markers to inline links or enforce heading/title shape -- both depend on the piece's eventual target channel (e.g. a letter has no heading hierarchy at all, and link-resolution is a later, optional, channel-specific packaging step), so both are the skill's own judgment call, not this tool's. " +
				"When the piece's required_stages includes research (set via require_piece_stages / /ps --research), refuses until research_mode is full (dec-20260809-701b59d3). " +
				"A redraft replaces draft.md wholesale -- expected, not data loss, since a draft isn't an accumulating list like sources.md. Never touches brief.md, outline.md, angles.md, or sources.md. Prefer this over Write/Bash so no extra filesystem permissions are needed.",
			InputSchema: saveDraftBundleSchema,
		},
		s.callSaveDraftBundle,
	)
}

func (s *Server) callSaveDraftBundle(args json.RawMessage) callToolResult {
	var b draft.Bundle
	if err := json.Unmarshal(args, &b); err != nil {
		return errorResult("invalid arguments: %s", err)
	}

	pieceDir, err := s.resolveAnchoredPass("draft", b.PiecePath)
	if err != nil {
		return errorResult("%s", err.Error())
	}

	validShortnames, err := research.LoadShortnames(pieceDir)
	if err != nil {
		return errorResult("%s", err)
	}
	if err := draft.Validate(b, validShortnames); err != nil {
		return errorResult("%s", err)
	}

	words, err := draft.WriteToPiece(pieceDir, b)
	if err != nil {
		return errorResult("%s", err)
	}
	if err := editorial.AppendRunLog(pieceDir, editorial.RunLogEntry{
		Pass:    "draft",
		Summary: fmt.Sprintf("piece=%s words=%d", b.PiecePath, words),
	}); err != nil {
		return errorResult("append run-log: %s", err)
	}

	return textResult("first full draft exists: %s/draft.md (%d words). thesis, killer sections, and counterintuitive findings are yours to have landed -- check them before handoff.", b.PiecePath, words)
}
