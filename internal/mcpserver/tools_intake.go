package mcpserver

import (
	"encoding/json"
	"fmt"
	"strings"

	"pysar/internal/editorial"
	"pysar/internal/intake"
)

var saveIntakeBundleSchema = map[string]interface{}{
	"type": "object",
	"properties": map[string]interface{}{
		"name":        map[string]string{"type": "string", "description": "Optional readable prefix for the piece directory. A unique suffix is always appended automatically — never ask the author about this. If omitted, derived from idea."},
		"idea":        map[string]string{"type": "string", "description": "Author idea text (required for entry_mode=idea)."},
		"entry_mode":  map[string]string{"type": "string", "description": "idea | draft"},
		"pov_source":  map[string]string{"type": "string", "description": "author-practitioner | author-draft"},
		"restatement": map[string]string{"type": "string", "description": "One concise line: audience + register + topic-scope (always shown, never a confirmation gate)."},
		"audience":    map[string]string{"type": "string"},
		"register":    map[string]string{"type": "string"},
		"expert_lens": map[string]string{"type": "string", "description": "The discipline/practitioner viewpoint whose rigor standard governs this piece's claims -- distinct from audience (who you write FOR). E.g. 'experimental science / methodology', 'horticulture', 'distributed-systems security'. Drives what counts as a claim worth grounding in Step 4."},
		"topic_scope": map[string]string{"type": "string"},
		"thesis":      map[string]string{"type": "string"},
		"promise":     map[string]string{"type": "string"},
		"killer_sections": map[string]interface{}{
			"type":        "array",
			"description": "1-2 sections with a genuine informational edge over existing coverage. A bare title is not enough — each needs a stated edge and a concrete example.",
			"items": map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"title":   map[string]string{"type": "string"},
					"edge":    map[string]string{"type": "string", "description": "Why this section beats existing/competing coverage of the topic."},
					"example": map[string]string{"type": "string", "description": "A concrete worked example to include — not a vague description."},
				},
				"required": []string{"title", "edge", "example"},
			},
		},
		"counterintuitive": map[string]interface{}{
			"type":        "array",
			"description": "1-3 myth-breaking claims, each with its contradiction named explicitly so it lands as a real section or punchline.",
			"items": map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"claim":         map[string]string{"type": "string"},
					"contradiction": map[string]string{"type": "string", "description": "What makes this claim counterintuitive, stated explicitly."},
				},
				"required": []string{"claim", "contradiction"},
			},
		},
		"key_questions": map[string]interface{}{
			"type":        "array",
			"description": "At least 1 question the piece must answer. Not optional — prevents scope drift once drafting starts.",
			"items":       map[string]string{"type": "string"},
		},
		"non_goals": map[string]interface{}{
			"type":        "array",
			"description": "At least 1 thing this piece will explicitly NOT cover. Not optional — prevents scope inflation once drafting starts.",
			"items":       map[string]string{"type": "string"},
		},
		"outline_md": map[string]string{"type": "string", "description": "Full outline.md Markdown body: sections and subsections, written as ## / ### headings."},
		"angles_md":  map[string]string{"type": "string", "description": "Full angles.md Markdown including a Contrarian section."},
		"draft_path": map[string]string{"type": "string", "description": "Optional path when entry_mode=draft."},
		"sources": map[string]interface{}{
			"type":        "array",
			"description": "Optional. Real sources actually fetched (WebSearch/WebFetch) to ground an uncertain technical claim. Empty is the common, legitimate default — most topics don't need this. Never list a source you didn't fetch.",
			"items": map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"url":  map[string]string{"type": "string"},
					"note": map[string]string{"type": "string", "description": "One line: what this source confirmed."},
				},
				"required": []string{"url", "note"},
			},
		},
	},
	"required": []string{
		"idea", "entry_mode", "pov_source", "restatement", "audience", "register", "expert_lens",
		"topic_scope", "thesis", "promise", "killer_sections", "counterintuitive",
		"key_questions", "non_goals", "outline_md", "angles_md",
	},
}

func (s *Server) registerSaveIntakeBundle() {
	s.register(
		tool{
			Name: "save_intake_bundle",
			Description: "Validate and persist intake scaffolding (brief/outline/angles/sources stub) under a provisional, uniquely-named piece directory (.pysar/pieces/<name>-<random-suffix>/). " +
				"Directory naming and any collision are handled automatically and silently — never surfaced as something the author needs to decide. " +
				"Runs the editorial intake Pass, appends run-log.jsonl, and writes STORAGE.md marking the layout provisional. " +
				"Rejects incomplete or degenerate idea input with a structured error. Never invents web sources. Prefer this over Write/Bash so no extra filesystem permissions are needed.",
			InputSchema: saveIntakeBundleSchema,
		},
		s.callSaveIntakeBundle,
	)
}

func (s *Server) callSaveIntakeBundle(args json.RawMessage) callToolResult {
	var b intake.Bundle
	if err := json.Unmarshal(args, &b); err != nil {
		return errorResult("invalid arguments: %s", err)
	}

	if intake.Degenerate(b.Idea) && b.EntryMode != intake.EntryDraft {
		return errorResult("degenerate idea input — ask the author once for a real idea or draft; do not scaffold")
	}

	pass := findPass("intake")
	if pass == nil {
		return errorResult("intake pass is not registered")
	}
	if _, err := editorial.Run(pass, editorial.NewState(editorial.SurfaceBlog)); err != nil {
		return errorResult("%s", err.Error())
	}

	dir, err := intake.WriteBundle(s.baseDir, b)
	if err != nil {
		return errorResult("%s", err)
	}

	if err := editorial.AppendRunLog(dir, editorial.RunLogEntry{
		Pass:    "intake",
		Summary: fmt.Sprintf("piece=%s entry=%s thesis=%q", b.Name, b.EntryMode, strings.TrimSpace(b.Thesis)),
	}); err != nil {
		return errorResult("append run-log: %s", err)
	}

	return textResult("intake saved under provisional path %s (storage_status=provisional). restatement: %s", dir, strings.TrimSpace(b.Restatement))
}

func (s *Server) registerReadAuthorDefaults() {
	s.register(
		tool{
			Name: "read_author_defaults",
			Description: "Read .pysar/voice.md and .pysar/style.md and return audience/register/tone defaults for intake. " +
				"Use this instead of Read/Bash on those files to avoid permission prompts and save tokens. " +
				"Missing profiles return an explicit generalist-fallback — never invent a silent generic-public audience when a profile exists.",
			InputSchema: map[string]interface{}{"type": "object", "properties": map[string]interface{}{}},
		},
		s.callReadAuthorDefaults,
	)
}

func (s *Server) callReadAuthorDefaults(args json.RawMessage) callToolResult {
	d, err := intake.LoadAuthorDefaults(s.baseDir)
	if err != nil {
		return errorResult("load author defaults: %s", err)
	}
	raw, err := json.Marshal(d)
	if err != nil {
		return errorResult("encode defaults: %s", err)
	}
	return textResult("%s", string(raw))
}

func findPass(name string) editorial.Pass {
	for _, p := range editorial.Registered() {
		if p.Name() == name {
			return p
		}
	}
	return nil
}
