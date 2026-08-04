package mcpserver

import (
	"encoding/json"
	"fmt"

	"pysar/internal/draft"
	"pysar/internal/editorial"
	"pysar/internal/research"
	"pysar/internal/seo"
)

var saveSEOBundleSchema = map[string]interface{}{
	"type": "object",
	"properties": map[string]interface{}{
		"piece_path": map[string]string{"type": "string", "description": "Path to an existing piece, relative to the project root -- not a slug. Accepts the piece directory itself or a path to a file inside it (e.g. brief.md, what a Claude Code @-reference commonly resolves to, or a raw/ excerpt) -- either form resolves to the same piece directory."},
		"revised_md": map[string]string{"type": "string", "description": "The FULL revised text, written to seo.md -- never to draft.md, staff-edit.md, or sharpen.md, which stay untouched. Replaces the whole seo.md file, not a diff/patch. Unlike every earlier pass, this content must contain ZERO [^shortname] markers -- every one must already be resolved into a real [anchor](url) link whose url matches a source this piece's research actually recorded. Never invent a link."},
		"checks": map[string]interface{}{
			"type":        "array",
			"description": "At least 1 one-line note on what changed and why, e.g. '[citation] resolved [^src-a] to an inline link on \"the retry budget\"' or '[title] tightened the title for a concrete keyword match'. Record a single 'no changes needed' entry when a check genuinely required none -- omitting checks entirely isn't a completed pass.",
			"items":       map[string]string{"type": "string"},
		},
		"mode": map[string]string{"type": "string", "description": "Optional, informational only: 'delta' (surgical) or 'rewrite' (structural). Folded into the changelog entry, not otherwise enforced."},
		"tags": map[string]interface{}{
			"type":        "array",
			"description": "At least 1 discovery keyword/topic tag for whatever platform this piece is headed to. No fixed slot count or platform-specific taxonomy -- pick what actually helps discovery.",
			"items":       map[string]string{"type": "string"},
		},
		"meta_title":       map[string]string{"type": "string", "description": "The search/share title -- distinct from the article's own '# Title' line (tuned for a search-results or link-preview context, not the feed/reader context)."},
		"meta_description": map[string]string{"type": "string", "description": "The search/share summary -- concrete, distinct from the article's own subtitle, written for a reader who hasn't opened the piece yet."},
		"url_slug":         map[string]string{"type": "string", "description": "A kebab-case, keyword-bearing suggested URL path segment."},
	},
	"required": []string{"piece_path", "revised_md", "checks", "tags", "meta_title", "meta_description", "url_slug"},
}

func (s *Server) registerSaveSEOBundle() {
	s.register(
		tool{
			Name: "save_seo_bundle",
			Description: "Validate and persist a /ps-seo pass: writes the revision to seo.md, appends seo-changelog.md, and (re)writes seo-checklist.md (tags, meta title/description, URL slug). Neither draft.md, staff-edit.md, nor sharpen.md is touched -- each stage keeps its own file. " +
				"Runs before humanize, never after (dec-20260804-e3234e50) -- humanize deliberately roughens the transitions this pass smooths, so reversing the order lets the later pass undo the earlier one. " +
				"Inverts every earlier pass's citation rule: revised_md must contain ZERO [^shortname] markers -- each one must already be resolved into a real [anchor](url) link matching a URL this piece's research actually recorded, never invented. Requires >=1 recorded check. " +
				"A re-run replaces seo.md and seo-checklist.md wholesale, same as the earlier passes' own writes -- expected, not data loss. Never touches brief.md, outline.md, angles.md, sources.md, draft.md, staff-edit.md, or sharpen.md. Prefer this over Write/Bash so no extra filesystem permissions are needed.",
			InputSchema: saveSEOBundleSchema,
		},
		s.callSaveSEOBundle,
	)
}

func (s *Server) callSaveSEOBundle(args json.RawMessage) callToolResult {
	var b seo.Bundle
	if err := json.Unmarshal(args, &b); err != nil {
		return errorResult("invalid arguments: %s", err)
	}

	pieceDir, err := s.resolveAnchoredPass("discoverability", b.PiecePath)
	if err != nil {
		return errorResult("%s", err.Error())
	}

	validSourceURLs, err := research.LoadSourceURLs(pieceDir)
	if err != nil {
		return errorResult("%s", err)
	}
	if err := seo.Validate(b, validSourceURLs); err != nil {
		return errorResult("%s", err)
	}

	words, err := seo.WriteToPiece(pieceDir, b)
	if err != nil {
		return errorResult("%s", err)
	}

	// Discoverability has the same "which file did this revise from"
	// ambiguity as sharpen/humanize (sharpen.md if present, else
	// staff-edit.md, else draft.md, per the skill's own instruction) --
	// same honest-proxy reasoning as their own revised_from tracking.
	revisedFrom := draft.LatestRevisionFile(pieceDir, "sharpen.md", "staff-edit.md", "draft.md")
	if err := editorial.AppendRunLog(pieceDir, editorial.RunLogEntry{
		Pass:    "discoverability",
		Summary: fmt.Sprintf("piece=%s words=%d checks=%d tags=%d revised_from=%s", b.PiecePath, words, len(b.Checks), len(b.Tags), revisedFrom),
	}); err != nil {
		return errorResult("append run-log: %s", err)
	}

	return textResult("seo done: %s/seo.md (%d words, %d checks recorded), %s/seo-checklist.md written. next: /ps-humanize.", b.PiecePath, words, len(b.Checks), b.PiecePath)
}
