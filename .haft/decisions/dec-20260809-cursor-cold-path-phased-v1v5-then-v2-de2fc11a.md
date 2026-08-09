---
id: dec-20260809-cursor-cold-path-phased-v1v5-then-v2-de2fc11a
kind: DecisionRecord
version: 1
status: active
title: Init rewrites Cursor mcp.json with a Cursor-visible binary path
mode: standard
valid_until: 2026-11-09
created_at: 2026-08-09T08:55:45Z
updated_at: 2026-08-09T08:55:45Z
links:
  - ref: prob-20260809-5a4334ac
    type: based_on
  - ref: sol-20260809-ee230a49
    type: based_on
---

# Init rewrites Cursor mcp.json with a Cursor-visible binary path

## 1. Problem Frame

**Signal:** Live dogfood on Cursor 3.15.6 (Customize → MCPs): after pysar init --cursor, project .cursor/mcp.json is present with command=pysar / serve / PYSAR_PROJECT_ROOT=${workspaceFolder}, but the pysar server stays status=disconnected (Cursor logs: project-0-test-cursor-pysar none→disconnected; createClient only for Cloudflare marketplace plugins; no stdio spawn). Customize 'Connected' lists only Plugin MCPs. macOS Dock-launched Cursor GUI PATH lacks ~/.local/bin, so bare 'pysar' cannot start even when toggled. End-author recovery today requires absolute command path, workspace-scope filter hunting, and/or MCP Logs — not a product path. This is not 'operator forgot a step'; it is the cold path our Cursor host DRR (dec-20260808-f3001106) claims to ship.

**Constraints:**
- Preserve hostAdapter + shared skill corpus (dec-20260808-f3001106); Claude and Codex scaffolds must stay green
- Piece persistence remains via pysar MCP tools (ps-* mechanical I/O contract), not agent raw Write/Bash of .pysar/**
- No requirement that the author be a developer who understands MCP, PATH, or Cursor Customize scopes
- Do not solve by shipping a 'open MCP Logs and fix absolute path' runbook as the primary product answer

**Acceptance:** DRAFT for operator edit: On a clean macOS machine with pysar installed only via the documented install path (binary on ~/.local/bin), run `pysar init --cursor <dir>`, open <dir> in Cursor launched from Dock (not from a shell with enriched PATH), enable any single first-run Cursor prompt if shown, then invoke `/ps` with a one-sentence idea. Observable success: pysar MCP tools are available to the agent AND a provisional piece appears under <dir>/.pysar/pieces/*/brief.md written via MCP — with zero hand-edits to mcp.json, zero Customize filter archaeology, and zero MCP Logs debugging.

## 2. Decision

**Selected:** Init rewrites Cursor mcp.json with a Cursor-visible binary path

**Selection policy:** PASS only if mcp_io_contract=pass and spawn_reliability≥3. Among survivors maximize discoverability; if discoverability gap ≥2 ordinal steps prefer higher discoverability even at worse ship_latency; else tie-break by lower ship_latency. Operator override of single-variant V2 advisory: bind near-term complementary stack V1+V5 now, commit V2 as explicit next phase (not rejected).

**Why selected:** Operator bind 2026-08-09: phased contract. Phase A (now): ship V1 (init writes Cursor-visible binary path into .cursor/mcp.json — prefer ${userHome}/.local/bin/pysar or resolved LookPath) together with V5 (init emits one explicit next action: Cursor MCP install deeplink and/or one-line Customize enable locus — no filter archaeology). Phase B (next): V2 Marketplace / install-link so Pysar appears in Customize Connected like other plugins, after we learn Cursor plugin packaging. Compare advisory preferred V2 alone for max discoverability; operator accepts that as the durable end-state but refuses to block cold-path dogfood on marketplace lead time. Preserves MCP mechanical I/O and hostAdapter Cursor scaffold.


**Invariants:**
- Piece persistence remains via pysar MCP tools (not CLI dual-path as primary cold path)
- Claude and Codex init scaffolds remain green; Phase A changes are Cursor-host scoped unless a shared helper is extracted
- Phase A must ship V1 and V5 together — path rewrite without one-step locus, or one-step without spawnable path, is incomplete
- V2 is committed next phase, not optional forever backlog without waive
- Primary product answer is not an MCP Logs troubleshooting essay

**Pre-conditions:**
- [ ] Portfolio sol-20260809-ee230a49 compared; operator chose phased bind
- [ ] pysar binary install path remains documented as ~/.local/bin (or init resolves LookPath)

**Post-conditions:**
- [ ] pysar init --cursor writes a Cursor-spawnable command path and prints exactly one enable locus/deeplink
- [ ] Docs/journey Cursor path describe that one step, not Customize archaeology
- [ ] Phase B tracked (commission/problem note) after Phase A lands

**Admissibility:**
- NOT: Bare command=pysar as the sole Cursor spawn contract for Dock GUI
- NOT: Runbook-primary cold path (MCP Logs / workspace filter hunting)
- NOT: Shipping V1 without V5 companion in Phase A
- NOT: Silent drop of Phase B without waive

## 3. Rationale

**Counterargument:** Phasing risks baking project-mcp.json + printed-toggle as the permanent Cursor story and delaying V2 indefinitely, leaving authors on a second-class surface while Cursor pushes marketplace Connected as the real UX — the same discoverability failure recurs in support.

**Selected variant weakest link:** Phase A still depends on one human enable action in Cursor UI; if deeplink/install-link format is wrong for Cursor 3.15+ or project MCP remains invisible even after the printed step, authors still cannot find Pysar and path rewrite alone does not fix Connected discoverability.

**Rejected alternatives:**
| Variant | Verdict | Reason |
|---------|---------|--------|
| Init rewrites Cursor mcp.json with a Cursor-visible binary path | **Selected** | Operator bind 2026-08-09: phased contract. Phase A (now):... |
| Cursor Marketplace / install-link first-class connect | Rejected | Not rejected — deferred as Phase B. Doing V2 alone now blocks near-term dogfood while we reverse-engineer Cursor plugin packaging. |
| CLI mechanical I/O so cold path does not need live MCP | Rejected | Fails mcp_io_contract constraint; reopens host-agnostic MCP-only persistence without a supersede DRR. |
| Install packaging puts pysar on GUI-visible PATH | Rejected | Dominated by V1 on discoverability and ship_latency; does not tell authors where to look in Customize. |
| Init emits Cursor deeplink + one-toggle enable copy (no archaeology) | Rejected | Not rejected — required companion in Phase A. Alone fails spawn_reliability (bare pysar on Dock PATH). |

**Evidence requirements:**
- Cursor init golden tests for mcp.json command form
- Live Dock Cursor smoke: printed one-step → MCP tools → brief.md
- Phase B kickoff artifact or waive note before valid_until

**Predictions:**
| Claim | Observable | Threshold |
|-------|------------|-----------|
| Phase A: after pysar init --cursor on a clean dir with binary at ~/.local/bin/pysar, Dock-launched Cursor following only the init-printed one-step can use pysar MCP tools and write a provisional piece via intake | scratch init --cursor; open in Dock Cursor; follow printed step only; /ps one-sentence idea → .pysar/pieces/*/brief.md via MCP | success without hand-editing mcp.json and without MCP Logs / Customize filter archaeology |
| Phase A goldens: shipped Cursor mcp scaffold no longer relies on bare command pysar as the only spawn form | go test ./cmd/pysar Cursor init goldens + assets/cursor/mcp.json content | command uses ${userHome}/.local/bin/pysar and/or init-time resolved absolute path; tests green |
| Phase B started: Marketplace plugin or official install-link work is either in-repo WIP with a ProblemCard/commission or explicitly waived with rationale before valid_until | haft status / repo for Cursor marketplace or install-link carrier | non-empty WIP (plugin manifest or deeplink generator + docs) OR note waiving Phase B with operator rationale |

## 4. Consequences

**Rollback plan:**
Triggers:
- Phase A ship: Dock Cursor + printed one-step still leaves pysar disconnected / not findable for a clean init --cursor scratch
- Phase A causes committed machine-specific absolute paths that break teammate clones without a portable ${userHome} form
- Phase B abandoned past valid_until with no marketplace/deeplink progress and no waive note
Steps:
1. Revert cursor mcp asset / init rewrite to prior bare command=pysar scaffold
2. Revert init stdout/docs one-step copy
3. Reopen problem or supersede this DRR with alternate bind (V2-first or V3 with explicit MCP-I/O supersede)
Blast radius: cmd/pysar Cursor host scaffold + init UX copy + Cursor journey docs; no Claude/Codex MCP dialect change required in Phase A

**Refresh triggers:**
- Cursor changes Customize MCP / install-link / plugin packaging again
- Documented install path moves off ~/.local/bin
- Phase A smoke fails twice on clean machines

**Affected files:** cmd/pysar/assets/cursor/mcp.json, cmd/pysar/host.go, cmd/pysar/init_test.go, docs/init.md, docs/troubleshooting.md, docs/mcp-and-skills.md, site/engineering/human-setup.md
