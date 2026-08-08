---
id: dec-20260808-codex-mcp-auto-approve-d9ebb3a1
kind: DecisionRecord
version: 6
status: deprecated
title: Ship default_tools_approval_mode = auto
mode: standard
valid_until: 2026-11-08T00:00:00Z
created_at: 2026-08-08T13:49:35Z
updated_at: 2026-08-08T14:33:10Z
links:
  - ref: prob-20260808-f9cd1683
    type: based_on
  - ref: sol-20260808-26f9b3d6
    type: based_on
---

# Ship default_tools_approval_mode = auto

## 1. Problem Frame

**Signal:** Live Codex smoke after pysar init --codex: Codex asks 'Allow the pysar MCP server to run tool read_author_defaults?' (and will ask again for other save_* tools). Claude init ships .claude/settings.json with mcp__pysar__* pre-approval so authors do not hit that class of prompt. Codex scaffold today only registers [mcp_servers.pysar] in .codex/config.toml without default_tools_approval_mode. OpenAI Codex docs support default_tools_approval_mode = auto|prompt|writes|approve on the MCP server block (project .codex/config.toml; trusted projects only). Operator wants Claude-like allowance for Codex.

**Constraints:**
- Extend Codex packaging only — do not change Claude .claude/settings.json contract or Cursor gap in this problem
- Stay inside project .codex/config.toml MCP server table (dec-20260808-codex-host-v4-ac3eae46 hostAdapter surface)
- Do not set global Codex approval_policy=never or sandbox_mode=danger-full-access as the product default
- Document trusted-project requirement for project-scoped .codex/config.toml
- Piece I/O remains via pysar MCP (dec-20260719-fa0366dd)

**Acceptance:** After pysar init --codex (or --force refresh) on a trusted Codex project, calling pysar MCP tools from a skill does not require a per-tool allow prompt for ordinary pysar tools; go test asserts the shipped .codex/config.toml contains the chosen approval setting; docs state the trusted-project caveat. Claude/Cursor scaffolds unchanged.

## 2. Decision

**Selected:** Ship default_tools_approval_mode = auto

**Selection policy:** Maximize claude_friction_parity and acceptance_fit under constraints (acceptance_fit >=4, scope_discipline >=4). Prefer lower config_maintenance_cost. Do not optimize write_tool_prompt_safety (observation). Prefer Codex-native server-level default_tools_approval_mode over per-tool inventory. Stay inside project .codex/config.toml mcp_servers.pysar — no global approval_policy=never product default.

**Why selected:** Operator asked for Claude-like MCP tool allowance after live Codex prompted on read_author_defaults. Claude pre-approves mcp__pysar__* (including saves). Codex's isomorphic dial is default_tools_approval_mode = auto on [mcp_servers.pysar]. One-line packaging change, documented by OpenAI Codex MCP config, minimal maintenance vs per-tool lists. Operator explicitly bound V1 after compare.


**Invariants:**
- Shipped assets/codex/config.toml [mcp_servers.pysar] includes default_tools_approval_mode = "auto"
- No product default of global approval_policy = never or sandbox_mode = danger-full-access for this decision
- Claude .claude/settings.json and Cursor scaffold unchanged by this decision
- Docs state that project .codex/config.toml MCP settings apply for trusted Codex projects
- Piece I/O remains via pysar MCP tools (dec-20260719-fa0366dd)
- Packaging stays on hostAdapter Codex surface from dec-20260808-codex-host-v4-ac3eae46

**Pre-conditions:**
- [ ] prob-20260808-f9cd1683 and sol-20260808-26f9b3d6 govern this choice
- [ ] Codex host scaffold already ships .codex/config.toml (dec-20260808-codex-host-v4-ac3eae46)

**Post-conditions:**
- [ ] config.toml embed contains default_tools_approval_mode = "auto"
- [ ] Init/Codex golden test asserts the key is present
- [ ] docs mention trusted-project caveat for MCP auto-approval
- [ ] Baseline snapshot taken

**Admissibility:**
- NOT: NOT: Setting global Codex approval_policy=never as the init default
- NOT: NOT: Changing Claude settings.json or Cursor MCP JSON in this decision
- NOT: NOT: Per-tool inventory that must track every new MCP tool name
- NOT: NOT: Claiming prompts disappear on untrusted projects without evidence

## 3. Rationale

**Counterargument:** V2 (writes) would keep human clicks on save/export and is safer if auto-approving write tools is undesirable — V1 trades that safety for full Claude friction parity.

**Selected variant weakest link:** Project-scoped .codex/config.toml applies only for trusted Codex projects — untrusted projects may still prompt despite auto; also Codex version skew if the key is ignored or renamed.

**Rejected alternatives:**
| Variant | Verdict | Reason |
|---------|---------|--------|
| Ship default_tools_approval_mode = auto | **Selected** | Operator asked for Claude-like MCP tool allowance after l... |
| Ship default_tools_approval_mode = writes | Rejected | Weaker Claude friction parity; fails acceptance_fit >=4 under the declared policy; write_prompt_safety is observation-only and does not justify under-shipping init allowance. |
| Docs-only: tell authors to paste auto | Rejected | Does not ship allowance at init — non-technical authors still hit per-tool prompts. |
| Per-tool tools.*.approval_mode = auto for every pysar tool | Rejected | Same runtime posture as V1 with higher maintenance — every new MCP tool must update TOML or silently re-prompts. |

**Evidence requirements:**
- go test Codex config golden (CL3)
- Live trusted-Codex smoke for prediction 2 (CL3 when available)
- OpenAI Codex MCP docs for default_tools_approval_mode (CL1)

**Predictions:**
| Claim | Observable | Threshold |
|-------|------------|-----------|
| Shipped Codex MCP config includes default_tools_approval_mode auto | go test ./cmd/pysar Codex init golden reads .codex/config.toml and finds default_tools_approval_mode = "auto" under [mcp_servers.pysar] | Test passes; Claude/Cursor goldens unchanged |
| Trusted Codex project stops per-tool allow prompts for pysar MCP tools after refresh | On a trusted Codex project after pysar init --codex --force, skill-driven read_author_defaults / save_* calls do not show Allow the pysar MCP server to run tool prompts | No per-tool allow prompt for ordinary pysar MCP tools in that trusted session |

## 4. Consequences

**Rollback plan:**
Triggers:
- Live trusted Codex still prompts per tool despite auto (key ignored or trust gate)
- Operator rejects auto-approving write tools and wants writes mode instead
- Codex renames/removes default_tools_approval_mode
Steps:
1. Remove default_tools_approval_mode from assets/codex/config.toml
2. Revert related test/docs lines
3. If switching to writes, supersede this DRR via /h-decide rather than silent partial change
Blast radius: cmd/pysar/assets/codex/config.toml + Codex init tests + docs; no Pass/MCP schema change

**Refresh triggers:**
- Codex changes MCP approval config keys or trusted-project rules
- Live smoke shows prompts persist on trusted projects
- Operator requests writes-mode split after experience with auto

**Affected files:** cmd/pysar/assets/codex/config.toml, cmd/pysar/init_test.go, docs/init.md, docs/troubleshooting.md, docs/mcp-and-skills.md

## Impact Measurement (2026-08-08)

**Verdict:** partial

**Findings:**
Scaffold ships default_tools_approval_mode=auto; goldens and docs updated; binary rebuilt/installed. Live trusted-Codex confirmation that per-tool Allow prompts are gone remains for the operator.

**Criteria met:**
- [x] config.toml embed contains default_tools_approval_mode = "auto"
- [x] Init/Codex golden test asserts the key is present
- [x] docs mention trusted-project caveat for MCP auto-approval

**Criteria NOT met:**
- [ ] Live trusted Codex session confirms no per-tool allow prompts

**Measurements:**
- go test ./cmd/pysar -run Codex|ClaudeAndCursor exit 0
- config.toml line default_tools_approval_mode = "auto"

## Impact Measurement (2026-08-08)

**Verdict:** failed

**Findings:**
Scaffold and docs for auto shipped and tests pass, but prediction 2 failed on trusted live Codex: per-tool Allow prompts continue. Root cause: auto is not Claude-parity allow-all for unannotated custom MCP tools after Codex #15519; pysar tools have no annotations. Haft never used auto for this purpose.

**Criteria met:**
- [x] config.toml embed contains default_tools_approval_mode = "auto"
- [x] docs mention trusted-project caveat for MCP auto-approval

**Criteria NOT met:**
- [ ] Live trusted Codex session confirms no per-tool allow prompts

**Measurements:**
- codex mcp get pysar: default_tools_approval_mode=auto
- tools/list: 16/16 annotations=None
- UI prompt on check_onboarding_status despite trusted+auto

## Deprecated (2026-08-08)

**Reason:** Measure failed; superseded in practice by dec-20260808-codex-mcp-approve-dial-cd11a72c (approve dial). auto is incorrect Claude-parity claim for unannotated MCP tools.
