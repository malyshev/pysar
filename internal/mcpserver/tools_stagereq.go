package mcpserver

import (
	"encoding/json"
	"fmt"
	"strings"

	"pysar/internal/editorial"
	"pysar/internal/research"
	"pysar/internal/stagereq"
)

var requirePieceStagesSchema = map[string]interface{}{
	"type": "object",
	"properties": map[string]interface{}{
		"piece_path": map[string]string{"type": "string", "description": "Path to an existing piece, relative to the project root -- not a slug. Accepts the piece directory itself or a path to a file inside it."},
		"stages": map[string]interface{}{
			"type":        "array",
			"description": "Opt-in stages this piece must complete before later saves (dec-20260809-701b59d3). Allowed values: \"research\" (blocks save_draft_bundle until research_mode is full), \"seo\" (blocks save_humanize_bundle until seo.md exists). Merges with any stages already on the piece; never clears existing ones.",
			"items":       map[string]string{"type": "string"},
		},
	},
	"required": []string{"piece_path", "stages"},
}

func (s *Server) registerRequirePieceStages() {
	s.register(
		tool{
			Name: "require_piece_stages",
			Description: "Persist opt-in stage preconditions on a piece's brief.md (required_stages). Call when /ps is invoked with --research and/or --seo so draft/humanize MCP gates can fail closed if the agent skips the stage. " +
				"Merges stages idempotently; does not run research or SEO itself. Prefer this over hand-editing brief.md.",
			InputSchema: requirePieceStagesSchema,
		},
		s.callRequirePieceStages,
	)
}

func (s *Server) callRequirePieceStages(args json.RawMessage) callToolResult {
	var in struct {
		PiecePath string   `json:"piece_path"`
		Stages    []string `json:"stages"`
	}
	if err := json.Unmarshal(args, &in); err != nil {
		return errorResult("invalid arguments: %s", err)
	}
	if strings.TrimSpace(in.PiecePath) == "" {
		return errorResult("piece_path is required")
	}
	if len(in.Stages) == 0 {
		return errorResult("stages must include at least one of %q or %q", stagereq.StageResearch, stagereq.StageSEO)
	}

	pieceDir, found := research.ResolvePieceDir(s.baseDir, in.PiecePath)
	if !found {
		return errorResult("no piece found at %q -- run intake first", in.PiecePath)
	}
	if err := stagereq.Require(pieceDir, in.Stages...); err != nil {
		return errorResult("%s", err)
	}
	stages, err := stagereq.Load(pieceDir)
	if err != nil {
		return errorResult("%s", err)
	}
	if err := editorial.AppendRunLog(pieceDir, editorial.RunLogEntry{
		Pass:    "require-stages",
		Summary: fmt.Sprintf("required_stages=%v", stages),
	}); err != nil {
		return errorResult("append run-log: %s", err)
	}
	return textResult("required_stages set on %s: %v. draft/humanize will refuse until each listed stage's outcome is present.", in.PiecePath, stages)
}
