---
id: dec-20260808-codex-host-v4-ac3eae46
kind: DecisionRecord
version: 17
status: active
title: Full Codex dialect including openai.yaml / invocation policy
mode: standard
valid_until: 2026-11-08T00:00:00Z
created_at: 2026-08-08T13:28:55Z
updated_at: 2026-08-09T17:23:05Z
links:
  - ref: prob-20260808-7ba11dda
    type: based_on
  - ref: sol-20260808-971c1818
    type: based_on
---

# Full Codex dialect including openai.yaml / invocation policy

## 1. Problem Frame

**Signal:** Operator wants Pysar usable inside OpenAI Codex (ChatGPT coding agent) as the next AI-editor host after Claude Code and Cursor. Today: pysar init --codex is intentionally stubbed (dec-20260718-8278c494) with 'not yet supported'; real host adapters exist only for Claude and Cursor (dec-20260808-f3001106). Shared ps-* skill corpus and MCP persistence already ship; a Codex author who installs the pysar binary still has no host-native skill install path, MCP/config dialect, or documented way to invoke the pipeline. Prior Cursor work explicitly left Codex stubbed. 'Add Codex support' is the proposed solution label — the broken condition is missing Codex-host scaffolding for the already-shipped agentic surface.

**Constraints:**
- Preserve dec-20260718-8278c494 flag surface (--codex stays a boolean; mutual exclusivity with --claude/--cursor)
- Extend existing hostAdapter registry from dec-20260808-f3001106 — do not invent a third dispatch shape
- One host-neutral skill corpus remains the source of truth; Codex packaging may differ (paths, MCP TOML/dialect, skill install dir) but editorial skill bodies must not fork into Codex-only prose
- Piece I/O remains via pysar MCP (dec-20260719-fa0366dd) — no regression to agent Write/Bash for .pysar/**
- Codex packaging must be grounded in verified Codex/haft CL1 conventions (paths, config format) before inventing a dialect — no Medium copy-paste, no guess-from-memory
- Do not change editorial Pass semantics, piece file formats, or Claude/Cursor golden init behavior

**Acceptance:** From a clean scratch directory, `pysar init --codex` exits 0 and writes a Codex-host project such that: (1) Codex can launch `pysar serve` as an MCP server for that project without hand-editing config (using Codex's native config dialect, not Claude's .mcp.json or Cursor's .cursor/mcp.json); (2) the operator can invoke an agentic entry equivalent to Claude/Cursor /ps (and at least stage skills intake through humanize + onboarding) with instructions that match the host-agnostic skill bodies already shipped; (3) piece persistence still goes through pysar MCP tools (not raw Write/Bash for .pysar/**); (4) `pysar init --claude` and `pysar init --cursor` behavior is unchanged; (5) a short smoke: open/run the project in Codex with a one-sentence idea and observe a provisional piece directory under .pysar/pieces/ with brief.md written via MCP.

## 2. Decision

**Selected:** Full Codex dialect including openai.yaml / invocation policy

**Selection policy:** Prefer the variant that matches verified Haft stable-Codex packaging CL1 from /Users/malyshev/Projects/Me/haft (project .codex/config.toml MCP + shared skill corpus under ~/.agents/skills + agents/openai.yaml implicit-invocation policy + Codex invocation-dialect rewrite as needed), subject to acceptance_fit >=4 and claude_cursor_regression_risk <=2. Maximize chance that Codex authors can invoke the pipeline without hand-editing. Do not optimize away codex_dialect_coupling (observation). Do not fork skill bodies. Extend hostAdapter registry only — no third dispatch shape.

**Why selected:** V4 is what Haft actually ships for full haft init --codex: native Codex MCP TOML, skills publication under .agents/skills, per-skill openai.yaml gating implicit invocation (only orchestrator implicit), and /h- -> $h- style rewrite for Codex dialect. Pysar already has hostAdapter + shared corpus from Cursor (dec-20260808-f3001106); V4 extends that adapter with the Codex-specific surfaces Haft proved necessary after prompts were deprecated. Operator explicitly bound V4 after CL1 review of Me/haft and Pareto comparison with V1.


**Invariants:**
- Flag surface remains --claude/--cursor/--codex with mutual exclusivity (dec-20260718-8278c494)
- Codex support extends hostAdapter registry from dec-20260808-f3001106 — no third dispatch shape
- One host-neutral skill corpus remains source of truth; Codex may transform packaging (paths, $ps- rewrite, openai.yaml) but must not fork editorial skill bodies into a second corpus
- pysar init --codex writes project .codex/config.toml launching pysar serve with project-root env (Codex native MCP carrier)
- Skills install to Codex skill root (~/.agents/skills by default, or documented --local equivalent if added later) via syncManagedTree
- Each installed Codex skill gets agents/openai.yaml; only the orchestrator skill (ps) allows implicit invocation; stage/onboarding skills are explicit-only
- Piece I/O remains via pysar MCP tools (dec-20260719-fa0366dd)
- pysar init --claude and pysar init --cursor golden behavior unchanged

**Pre-conditions:**
- [ ] prob-20260808-7ba11dda and sol-20260808-971c1818 remain the governing frame/portfolio
- [ ] hostAdapter registry and shared assets/skills corpus already shipped (dec-20260808-f3001106)
- [ ] Codex packaging details re-checked against /Users/malyshev/Projects/Me/haft before inventing fields

**Post-conditions:**
- [ ] pysar init --codex exits 0 with real Codex scaffold (not not-yet-supported)
- [ ] TestInitCodexNotYetSupported replaced by positive Codex scaffold tests including MCP TOML + skill install + openai.yaml policy shape
- [ ] docs/init.md and troubleshooting updated for Codex support
- [ ] Baseline snapshot taken for drift detection

**Admissibility:**
- NOT: NOT: Inventing a third host dispatch shape outside hostAdapter
- NOT: NOT: Forking ps-* skill bodies into Codex-only editorial prose
- NOT: NOT: Writing MCP to shared ~/.codex/config.toml instead of project .codex/config.toml without a superseding decision
- NOT: NOT: Allowing all stage skills implicit invocation by default
- NOT: NOT: Changing editorial Pass/MCP tool contracts or piece file formats

## 3. Rationale

**Counterargument:** V1 (Haft-CL1 Codex adapter without openai.yaml/policy) may be enough if current Codex discovers SKILL.md without policy files — V4 then overbuilds coupling and maintenance for zero acceptance gain, paying the observation dimension we said not to optimize.

**Selected variant weakest link:** Codex CLI/App version churn on openai.yaml, prompts-vs-skills, and invocation naming — Haft's own changelog already shows these surfaces moving; a Pysar adapter that mirrors today's dialect can rot when Codex changes.

**Rejected alternatives:**
| Variant | Verdict | Reason |
|---------|---------|--------|
| Full Codex dialect including openai.yaml / invocation policy | **Selected** | V4 is what Haft actually ships for full haft init --codex... |
| Haft-CL1 Codex adapter on existing hostAdapter | Rejected | On the Pareto front for reuse/locality, but under-ships vs verified full Haft --codex (missing openai.yaml / implicit policy / dialect rewrite). Operator chose CL1-faithful Codex support over thinner first slice. |
| MCP TOML only; skills documented not installed | Rejected | Fails acceptance: no installed skills means no agentic /ps-equivalent entry for non-technical authors. |
| Project-local Codex skills (.agents/skills in repo) | Rejected | Diverges from Pysar Claude/Cursor global skill install and from Haft's default Codex global .agents/skills locality. |
| Fork Codex-only skill bodies from shared corpus | Rejected | Violates host-neutral skill corpus invariant from dec-20260808-f3001106. |

**Evidence requirements:**
- go test ./cmd/pysar Codex + Claude + Cursor goldens (CL3)
- CL1 re-read of Me/haft Codex init paths at implement time
- Live Codex intake smoke for prediction 2 (CL3 when available)

**Predictions:**
| Claim | Observable | Threshold |
|-------|------------|-----------|
| pysar init --codex scaffolds MCP + skills + openai.yaml without error | go test ./cmd/pysar covering --codex: exits 0; writes .codex/config.toml with pysar serve; installs shared skill corpus under Codex skills root; each skill has agents/openai.yaml with only ps allow_implicit_invocation true | All Codex golden tests pass; Claude and Cursor goldens still pass |
| Codex can run intake smoke via MCP after init | From a Codex-inited scratch project, invoke orchestrator/intake equivalent with a one-sentence idea; provisional .pysar/pieces/*/brief.md appears via MCP | brief.md present and written through pysar MCP (not raw agent Write to .pysar) |

## 4. Consequences

**Rollback plan:**
Triggers:
- Live Codex smoke fails because openai.yaml or $ps- rewrite is wrong/unsupported on the Codex version under test
- Codex packaging causes Claude/Cursor golden init regressions
- Operator finds V4 coupling cost exceeds benefit and wants V1 thin adapter
Steps:
1. Revert codexHost adapter and assets/codex additions
2. Restore resolveHost --codex not-yet-supported stub and TestInitCodexNotYetSupported
3. Restore docs stub wording
4. If keeping MCP-only interim, supersede this DRR toward V1/V2 via /h-decide rather than silent partial ship
Blast radius: cmd/pysar host adapter + Codex assets + docs; no piece format or MCP tool schema change

**Refresh triggers:**
- Codex CLI/App changes skill root, MCP TOML schema, or openai.yaml semantics
- Haft Me/haft changes its Codex adapter in a way that invalidates CL1 paths
- Second live smoke failure on Codex after a successful first measure

**Affected files:** cmd/pysar/host.go, cmd/pysar/main.go, cmd/pysar/init_test.go, cmd/pysar/assets/codex/, cmd/pysar/assets/skills/, docs/init.md, docs/troubleshooting.md, docs/index.md, README.md

## Impact Measurement (2026-08-08)

**Verdict:** partial

**Findings:**
codexHost registered; init --codex scaffolds project .codex/config.toml, installs Codex-packaged shared skills to ~/.agents/skills with openai.yaml policies, docs updated. Unit/integration goldens for prediction 1 pass. Live Codex MCP intake smoke (prediction 2) not run in this session.

**Criteria met:**
- [x] pysar init --codex exits 0 with real Codex scaffold (not not-yet-supported)
- [x] TestInitCodexNotYetSupported replaced by positive Codex scaffold tests including MCP TOML + skill install + openai.yaml policy shape
- [x] docs/init.md and troubleshooting updated for Codex support

**Criteria NOT met:**
- [ ] Live Codex intake smoke writing .pysar/pieces/*/brief.md via MCP

**Measurements:**
- go test ./... exit 0
- TestInitCodexScaffoldsMCPAndSkills pass
- TestInitCodexSkillsShareCorpusEditorialWithClaude pass

## Impact Measurement (2026-08-08)

**Verdict:** accepted

**Findings:**
/h-verify: prediction 1 (scaffold goldens + openai.yaml policy shape) still holds; prediction 2 (live Codex via MCP) holds — intake brief.md present and full pipeline run-log shows intake→draft→staff-edit→sharpen→humanize on test-codex piece, plus root export.

**Criteria met:**
- [x] pysar init --codex scaffolds MCP + skills + openai.yaml without error
- [x] Codex can run intake smoke via MCP after init

**Measurements:**
- go test Codex goldens ok
- ps openai.yaml implicit true; stage skills false
- brief.md + run-log passes on testing-what-matters-ed5e514799b0

## Impact Measurement (2026-08-08)

**Verdict:** accepted

**Findings:**
/h-verify on released pysar 0.3.0: both predictions hold. Release binary init --codex scaffolds MCP+skills+openai.yaml; goldens green; Claude/Cursor init unaffected. Live Codex MCP pipeline artifacts on test-codex still satisfy prediction 2. Re-baselining after incidental drift.

**Criteria met:**
- [x] pysar init --codex scaffolds MCP + skills + openai.yaml without error
- [x] Codex can run intake smoke via MCP after init

**Measurements:**
- pysar 0.3.0 ~/.local/bin
- scratch init --codex: approve + 12 skills + openai.yaml policies
- go test Init*|Codex*|ResolveHost ok
- test-codex brief.md + run-log + export present

## Impact Measurement (2026-08-08)

**Verdict:** accepted

**Findings:**
/h-verify re-run: both predictions hold on released pysar 0.3.0 and prior live Codex smoke artifacts. Re-baselined for incidental drift.

**Criteria met:**
- [x] pysar init --codex scaffolds MCP + skills + openai.yaml without error
- [x] Codex can run intake smoke via MCP after init

**Measurements:**
- pysar 0.3.0 init --codex: approve + 12 skills + openai.yaml
- Claude/Cursor init exit 0
- go test InitCodex*|InitClaude*|InitCursor*|ResolveHost ok
- test-codex brief.md + run-log + export present

## Impact Measurement (2026-08-08)

**Verdict:** accepted

**Findings:**
Bundle drift verify: both predictions hold. init --codex + approve dial intact; prior live Codex smoke stands. Drift incidental site/.pnpm-store.

**Criteria met:**
- [x] init --codex scaffolds MCP+skills+openai.yaml
- [x] prior Codex MCP intake smoke

**Measurements:**
- scratch init --codex exit 0
- approve in assets/codex/config.toml
- pysar 0.3.0

## Impact Measurement (2026-08-08)

**Verdict:** accepted

**Findings:**
Both predictions hold. Drift is undeployed site/llms.txt only — incidental.

**Criteria met:**
- [x] init --codex scaffolds
- [x] prior live Codex smoke

**Measurements:**
- scratch init --codex exit 0
- approve in config.toml

## Impact Measurement (2026-08-09)

**Verdict:** accepted

**Findings:**
claim-001 Codex init goldens PASS. Shared-file drift incidental. claim-002 live Codex intake smoke not due until 2026-09-08 and not re-run; prior support left standing.

**Criteria met:**
- [x] Codex scaffold goldens

**Measurements:**
- claim-001: TestInitCodex* PASS

## Impact Measurement (2026-08-09)

**Verdict:** accepted

**Findings:**
Incidental docs + false-positive stagereq file adds under broad baseline. Codex dialect packaging (openai.yaml, config.toml, hostAdapter) not altered by this drift.

**Criteria met:**
- [x] Flag surface mutual exclusivity untouched
- [x] Piece I/O remains MCP tools

**Measurements:**
- docs/mcp-and-skills.md Persistence rule only for Codex-relevant surface
- no Codex scaffold file changes in drift set
