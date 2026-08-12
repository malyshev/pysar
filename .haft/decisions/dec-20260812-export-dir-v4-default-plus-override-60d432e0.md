---
id: dec-20260812-export-dir-v4-default-plus-override-60d432e0
kind: DecisionRecord
version: 3
status: active
title: Project default export_dir + optional per-call override
context: pipeline
mode: standard
valid_until: 2026-11-12
created_at: 2026-08-12T12:58:39Z
updated_at: 2026-08-12T15:19:27Z
links:
  - ref: prob-20260812-906a271d
    type: based_on
  - ref: sol-20260812-4f937d99
    type: based_on
---

# Project default export_dir + optional per-call override

## 1. Problem Frame

**Signal:** Finished pieces always land as <project root>/<slug>.md via WriteToRoot / export_piece_to_root. Operators want to choose where finished Markdown lands: (1) set an export directory at pysar init, (2) change that setting later in project config without re-init, (3) MCP must expose the resolved landing path so agents write/export to the configured place instead of assuming repo root.

**Constraints:**
- Default remains project root when unset (backward compatible with existing projects and docs that say root)
- Export destination is a project-level setting, not a global ~/.pysar preference
- Path must be constrained to inside the project (no escape outside project root) unless an explicit later decision allows otherwise
- Compose with pending mechanical citation resolve at export (dec-20260810-export-resolve-citations-6afd8c31) — destination change must not block that work
- Init stays non-interactive flags (Cobra flags), consistent with existing --claude/--cursor/--codex style
- Pysar output discipline: no noise explaining absences; MCP/skill text must tell agents the configured path, not narrate skips

**Acceptance:** After init with an export-dir option (or after editing the project config to a non-root relative path), export_piece_to_root (or successor) writes <configured-dir>/<slug>.md, creates the directory if missing, returns the absolute or project-relative destination path in the tool result, and default/omitted config still exports to project root. A post-init config change is honored on the next export without re-init. Observable via Go tests + one MCP call result path.

## 2. Decision

**Selected:** Project default export_dir + optional per-call override

**Selection policy:** Pareto among variants with operator_set_and_forget>=4, mcp_agent_landing_clarity>=4, backward_compat>=4; optimize implementation_surface; human binds via /h-decide. Operator chose V4 over advisory V1 explicitly to keep AI-adjustable per-call path later without a separate config MCP.

**Why selected:** Operator bound V4 after confirming the per-call MCP override is the path for later AI-adjusted landing without requiring a separate config tool first. Keeps durable export_dir in .pysar/project (init --export-dir + later JSON edit) and adds optional export_dir on the export MCP tool with project-default precedence when omitted.


**Invariants:**
- Absent or empty export_dir means project root (backward compatible)
- export_dir and per-call override are project-relative; resolved path must stay under project root
- MCP export tool returns the resolved destination path on success
- MkdirAll creates missing export_dir components before write
- Composes with citation resolve-at-export when that ships — destination choice does not skip resolve
- Init --export-dir is optional; omitting it leaves root default

**Pre-conditions:**
- [ ] projectManifest remains the durable project settings carrier (.pysar/project)
- [ ] WriteToRoot / export_piece_to_root remain the single write path for finished Markdown

**Post-conditions:**
- [ ] pysar init --export-dir PATH writes export_dir into .pysar/project
- [ ] Editing .pysar/project export_dir changes the next export without re-init
- [ ] export MCP tool accepts optional export_dir override; omit uses project default
- [ ] Successful export response includes the resolved dest path
- [ ] Go tests cover root default, configured dir, override precedence, and path-escape rejection

**Admissibility:**
- NOT: Export destination escaping outside the project root
- NOT: Silent default change away from project root for existing manifests without export_dir
- NOT: Requiring --seo or any pipeline stage just to choose a landing path
- NOT: Agents inventing absolute paths outside the project without an explicit later decision

## 3. Rationale

**Counterargument:** The optional per-call override adds API surface and a skill-discipline failure mode: agents may ignore the project default and pass ad-hoc paths every time, defeating set-and-forget. V1 avoids that by having only one authority (the manifest).

**Selected variant weakest link:** Precedence and skill instructions must stay unambiguous — if /ps and export_piece_to_root do not default to the project export_dir and only pass override when the author asked, agents will scatter finished pieces and the durable config becomes dead letter.

**Rejected alternatives:**
| Variant | Verdict | Reason |
|---------|---------|--------|
| Project default export_dir + optional per-call override | **Selected** | Operator bound V4 after confirming the per-call MCP overr... |
| V1 | Rejected | Meets init+config+MCP but has no AI one-off redirect without editing the manifest; operator explicitly wants the later AI-adjust path that V4's override provides. |
| V2 | Rejected | Config CLI is dominated on compare scores by V4 for post-init UX without introducing a new CLI surface for one field; revisit if JSON edit proves painful. |
| V3 | Rejected | No durable project setting — fails set-at-init and change-later requirements. |
| V5 | Rejected | Fixed convention cannot be specified or changed; breaks root default for everyone. |

**Predictions:**
| Claim | Observable | Threshold |
|-------|------------|-----------|
| Init with --export-dir writes that relative path into .pysar/project and subsequent export without override lands under that directory | Go test: init --export-dir published → manifest has export_dir; export fixture → published/<slug>.md exists | PASS |
| MCP optional export_dir overrides the project default for that call only and returns the overridden dest path | Go/MCP test: manifest export_dir=published, call with export_dir=outbox → file at outbox/<slug>.md; response cites that path; next call without override returns to published/ | PASS |
| Existing projects without export_dir still export to project root | Go test: empty/missing export_dir → <projectRoot>/<slug>.md | PASS |

## 4. Consequences

**Rollback plan:**
Triggers:
- Path-escape bugs land files outside the project
- Override precedence confuses agents so exports scatter unpredictably
- Manifest schema change proves costlier than the benefit vs V1
Steps:
1. Omit reading export_dir / override; restore WriteToRoot(projectRoot, pieceDir) join to project root only
2. Remove --export-dir from init and optional MCP arg
3. Leave stale export_dir keys in manifests harmless (ignored) or document delete
4. Revert skill/docs wording to project root
Blast radius: cmd/pysar init flag + projectManifest field; internal/export destination resolution; MCP export schema/docs/skills. Existing projects with export_dir set would need field ignored or removed to restore root-only behavior.

**Refresh triggers:**
- Operators ask for pysar config CLI because JSON edit is painful
- Agents routinely ignore project default and invent per-call paths
- Need for absolute/outside-project export destinations
- Tool rename away from export_piece_to_root becomes load-bearing

**Affected files:** cmd/pysar/main.go, cmd/pysar/init_test.go, internal/export/write.go, internal/export/write_test.go, internal/mcpserver/tools_export.go, internal/mcpserver/tools_export_test.go, cmd/pysar/assets/skills/ps/SKILL.md, plugins/pysar/skills/ps/SKILL.md, docs/export.md, docs/init.md, docs/pipeline.md, docs/mcp-and-skills.md

## Impact Measurement (2026-08-12)

**Verdict:** accepted

**Findings:**
All three DRR predictions hold under Go tests. Drift vs pre-impl baseline is expected implementation of V4 (WriteToRoot signature, init flag, MCP schema, skills/docs) — re-baseline after this measure. Weakest-link mitigation present in /ps SKILL.md (omit export_dir unless author asked). Citation resolve-at-export remains a separate pending decision and was not required for these predictions.

**Criteria met:**
- [x] pysar init --export-dir PATH writes export_dir into .pysar/project
- [x] Editing .pysar/project export_dir (or init patch) changes next export without full re-scaffold
- [x] export MCP accepts optional export_dir; omit uses project default
- [x] Successful export response includes resolved dest path
- [x] Go tests cover root default, configured dir, override precedence, path-escape rejection

**Measurements:**
- TestInitExportDirWritesManifestField: PASS
- TestSaveExportBundleUsesProjectExportDir: PASS (published/<slug>.md)
- TestSaveExportBundleOverrideBeatsProjectDefault: PASS (outbox then published)
- TestSaveExportBundleFallsBackToDraftAlone / RerunOverwritesRootFile: PASS (project root)
- TestSaveExportBundleRejectsEscapingExportDir + Resolve/Write escape tests: PASS
