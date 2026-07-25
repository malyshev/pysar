package mcpserver

import (
	"encoding/json"
	"fmt"
	"path/filepath"
	"strings"

	"pysar/internal/editorial"
	"pysar/internal/research"
)

var saveResearchBundleSchema = map[string]interface{}{
	"type": "object",
	"properties": map[string]interface{}{
		"piece_path":        map[string]string{"type": "string", "description": "Path to an existing piece, relative to the project root -- not a slug. Accepts either the piece directory itself (e.g. .pysar/pieces/<name>/) or a path to a file inside it (e.g. .pysar/pieces/<name>/brief.md, what a Claude Code @-reference commonly resolves to) -- either form resolves to the same piece directory. Leave empty for standalone research with no piece yet; topic is required in that case."},
		"topic":             map[string]string{"type": "string", "description": "Required when piece_path is empty (standalone mode). What the research is about."},
		"expert_lens":       map[string]string{"type": "string", "description": "The discipline/practitioner viewpoint judging source authority for this topic -- reuse the piece's own expert_lens from its brief.md when piece_path is set (read it first), or determine one fresh for standalone research."},
		"topic_family_note": map[string]string{"type": "string", "description": "Optional. Only set when expert_lens doesn't map cleanly to an obvious source-authority standard -- names the mapping chosen so a reviewer can see it. Leave empty when the mapping is obvious (the common case)."},
		"sources": map[string]interface{}{
			"type":        "array",
			"description": "The actual research. >=1 required; 6-12 typical, 8-14 without competitors. >=60% must be primary or secondary tier (community-tier sources don't count toward that floor). Every field verified against ps-factcheck's Job 2 discipline -- the note must match what was actually fetched, not memory of it.",
			"items": map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"shortname": map[string]string{"type": "string", "description": "kebab-case, unique within this research pass -- e.g. 'docker-buildkit-docs'."},
					"url":       map[string]string{"type": "string"},
					"tier":      map[string]string{"type": "string", "description": "primary | secondary | community"},
					"accessed":  map[string]string{"type": "string", "description": "YYYY-MM-DD, the day this was actually fetched."},
					"key_claims": map[string]interface{}{
						"type":        "array",
						"description": "2-5 bullets: what this source actually supports.",
						"items":       map[string]string{"type": "string"},
					},
					"notes":       map[string]string{"type": "string", "description": "Paraphrase: caveats, bias, staleness. Re-checked against the fetched content, not memory of it."},
					"raw_excerpt": map[string]string{"type": "string", "description": "Verbatim quoted text for raw/<shortname>.md. Quote exactly; paraphrase goes in notes, not here."},
				},
				"required": []string{"shortname", "url", "tier", "accessed", "key_claims", "notes", "raw_excerpt"},
			},
		},
		"competitors": map[string]interface{}{
			"type":        "array",
			"description": "Optional (--competitors=). One entry per competitor URL supplied.",
			"items": map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"url":               map[string]string{"type": "string"},
					"angle":             map[string]string{"type": "string"},
					"strongest_section": map[string]string{"type": "string"},
					"gap":               map[string]string{"type": "string"},
				},
				"required": []string{"url", "angle", "strongest_section", "gap"},
			},
		},
		"key_questions_additions": map[string]interface{}{
			"type":        "array",
			"description": "New, source-backed questions this research answers. Piece-anchored: appended to the existing key_questions, never replacing them. Standalone: seeds research-summary.md.",
			"items":       map[string]string{"type": "string"},
		},
		"angles_misconceptions": map[string]interface{}{
			"type":        "array",
			"description": "New, source-backed misconceptions this research surfaced. Appended, never replacing the operator's own angles.",
			"items":       map[string]string{"type": "string"},
		},
		"angles_contrarian": map[string]interface{}{
			"type":        "array",
			"description": "New, source-backed contrarian/under-discussed angles this research surfaced. Appended, never replacing the operator's own angles.",
			"items":       map[string]string{"type": "string"},
		},
	},
	"required": []string{"expert_lens", "sources"},
}

func (s *Server) registerSaveResearchBundle() {
	s.register(
		tool{
			Name: "save_research_bundle",
			Description: "Validate and persist a /ps-research pass. Piece-anchored (piece_path set): writes real sources.md + raw/ excerpts + optional competitors.md, sets the piece's research_mode to 'full', and appends (never replaces) key_questions/angles.md entries -- thesis, killer_sections, and counterintuitive are never touched. " +
				"Standalone (piece_path empty, topic required): writes the same sources.md/raw/competitors.md under a fresh .pysar/research/<topic>-<suffix>/ directory plus a research-summary.md consumable by /ps-intake --from-draft=. " +
				"Enforces the authority floor (>=60% primary/secondary tier) and citation hygiene (unique shortname per source) server-side. Never invents a source. Prefer this over Write/Bash so no extra filesystem permissions are needed.",
			InputSchema: saveResearchBundleSchema,
		},
		s.callSaveResearchBundle,
	)
}

func (s *Server) callSaveResearchBundle(args json.RawMessage) callToolResult {
	var b research.Bundle
	if err := json.Unmarshal(args, &b); err != nil {
		return errorResult("invalid arguments: %s", err)
	}

	if strings.TrimSpace(b.PiecePath) == "" {
		dir, err := research.WriteStandalone(s.baseDir, b)
		if err != nil {
			return errorResult("%s", err)
		}
		return textResult("standalone research saved under %s. no piece was created -- run /ps-intake --from-draft=%s/research-summary.md when ready to turn this into a piece.", dir, dir)
	}

	pass := findPass("research")
	if pass == nil {
		return errorResult("research pass is not registered")
	}
	pieceDir := research.ResolvePieceDir(s.baseDir, b.PiecePath)
	state := editorial.NewState(editorial.SurfaceBlog)
	if fileExists(filepath.Join(pieceDir, "brief.md")) {
		state = state.WithProduced(editorial.ArtifactStake).WithProduced(editorial.ArtifactBrief)
	}
	if _, err := editorial.Run(pass, state); err != nil {
		return errorResult("%s", err.Error())
	}

	if err := research.WriteToPiece(s.baseDir, b); err != nil {
		return errorResult("%s", err)
	}
	if err := editorial.AppendRunLog(pieceDir, editorial.RunLogEntry{
		Pass:    "research",
		Summary: fmt.Sprintf("sources=%d competitors=%d", len(b.Sources), len(b.Competitors)),
	}); err != nil {
		return errorResult("append run-log: %s", err)
	}

	return textResult("research saved to %s (research_mode=full, %d sources). thesis and killer sections untouched -- your take stayed yours.", b.PiecePath, len(b.Sources))
}
