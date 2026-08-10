---
id: dec-20260808-codex-mcp-approve-dial-cd11a72c
kind: DecisionRecord
version: 12
status: active
title: Ship default_tools_approval_mode = approve
mode: standard
valid_until: 2026-11-08
created_at: 2026-08-08T14:19:55Z
updated_at: 2026-08-10T16:08:23Z
links:
  - ref: prob-20260808-2edd57d2
    type: based_on
  - ref: sol-20260808-e01ffb9a
    type: based_on
---

# Ship default_tools_approval_mode = approve

## 1. Problem Frame

**Signal:** dec-20260808-codex-mcp-auto-approve-d9ebb3a1 shipped default_tools_approval_mode=auto; trusted test-codex still shows Allow the pysar MCP server to run tool prompts. Codex treats missing annotations as risky under auto; pysar tools have no annotations. Operator wants Claude-like no-prompt for pysar MCP on Codex.

**Constraints:**
- Stay inside project .codex/config.toml mcp_servers.pysar packaging or MCP tool metadata
- Do not set global approval_policy=never as product default
- Claude/Cursor scaffolds unchanged unless required
- Reversible via init --force + config line change

**Acceptance:** After pysar init --codex refresh on a trusted project, ordinary pysar MCP tool calls from skills do not show per-tool Allow prompts; go test asserts shipped .codex/config.toml uses the chosen dial; docs match.

## 2. Decision

**Selected:** Ship default_tools_approval_mode = approve

**Selection policy:** Maximize claude_friction_parity and codex_semantic_correctness under scope_discipline >=4. Prefer smallest packaging change. Operator commissioned solve after auto measure failed.

**Why selected:** Live trusted Codex still prompted under auto because unannotated MCP tools require approval. Codex AppToolApproval::Approve is the skip-review dial; Haft operators use approval_mode=approve for pre-allow. Operator asked to solve the issue after that diagnosis.


**Invariants:**
- Shipped assets/codex/config.toml uses default_tools_approval_mode = "approve" under [mcp_servers.pysar]
- Claude and Cursor scaffolds remain unchanged by this dial
- No product default of approval_policy = "never"

**Pre-conditions:**
- [ ] Codex project-scoped MCP config applies for trusted projects

**Post-conditions:**
- [ ] go test asserts approve in shipped config
- [ ] docs describe approve (not auto) as the Claude-parity dial
- [ ] test-codex refreshable via pysar init --codex --force

**Admissibility:**
- NOT: Leaving auto as the shipped Claude-parity claim after measure failed
- NOT: Relying on trust alone without the correct dial

## 3. Rationale

**Counterargument:** Approve naming confuses readers; ARC/guardian on ChatGPT desktop might still interrupt despite Approve.

**Selected variant weakest link:** ChatGPT desktop guardian/ARC may still surface a prompt even when Approve skips the MCP tool approval path.

**Rejected alternatives:**
| Variant | Verdict | Reason |
|---------|---------|--------|
| Ship default_tools_approval_mode = approve | **Selected** | Live trusted Codex still prompted under auto because unan... |
| Emit MCP tool annotations only; keep auto | Rejected | Does not stop save-tool prompts; fails Claude parity acceptance for this turn. |
| approve dial + honest annotations | Rejected | Larger blast than needed to stop prompts now; deferred as follow-up. |

**Evidence requirements:**
- go test ./cmd/pysar for Codex config golden
- codex mcp get pysar shows approve after refresh
- live smoke: no per-tool Allow prompt for check_onboarding_status

**Predictions:**
| Claim | Observable | Threshold |
|-------|------------|-----------|
| Shipped Codex MCP config uses approve | assets/codex/config.toml and init golden contain default_tools_approval_mode = "approve" | go test ./cmd/pysar passes |
| Trusted Codex stops per-tool Allow prompts for pysar tools after refresh | skill-driven check_onboarding_status / read_author_defaults do not show Allow the pysar MCP server prompts | no per-tool allow prompt in trusted session |

## 4. Consequences

**Rollback plan:**
Triggers:
- Live trusted Codex still prompts after approve refresh
- Approve dial rejected by Codex version
Steps:
1. Revert config.toml to previous value
2. Update tests/docs
3. Rebuild binary
Blast radius: Codex scaffold + docs + smoke project config only

**Refresh triggers:**
- Live smoke still prompts under approve
- Codex renames Approve semantics
- Annotations follow-up ships

**Affected files:** cmd/pysar/assets/codex/config.toml, cmd/pysar/init_test.go, docs/init.md, docs/mcp-and-skills.md, docs/troubleshooting.md, .haft/decisions/dec-20260808-codex-mcp-auto-approve-d9ebb3a1.md

## Impact Measurement (2026-08-08)

**Verdict:** partial

**Findings:**
Packaging shipped: assets/codex/config.toml uses approve; go test ./cmd/pysar passes; test-codex refreshed via pysar init --codex --force; codex mcp get pysar reports default_tools_approval_mode: approve. Live ChatGPT desktop smoke (no per-tool Allow prompt) still needs operator confirmation.

**Criteria met:**
- [x] go test asserts approve in shipped config
- [x] docs describe approve (not auto) as the Claude-parity dial
- [x] test-codex refreshable via pysar init --codex --force

**Criteria NOT met:**
- [ ] Live trusted Codex session confirms no per-tool allow prompts

**Measurements:**
- go test ./cmd/pysar ok
- codex mcp get pysar: default_tools_approval_mode=approve
- test-codex/.codex/config.toml contains approve

## Impact Measurement (2026-08-08)

**Verdict:** accepted

**Findings:**
Packaging (approve dial) + trusted live Codex smoke both hold. Operator confirmed full pipeline completed without per-tool Allow prompts; piece exported.

**Criteria met:**
- [x] go test asserts approve in shipped config
- [x] docs describe approve (not auto) as the Claude-parity dial
- [x] test-codex refreshable via pysar init --codex --force
- [x] Live trusted Codex session confirms no per-tool allow prompts

**Measurements:**
- go test ./cmd/pysar ok
- codex mcp get pysar: approve
- operator live smoke: pipeline completed, no MCP allow prompts

## Impact Measurement (2026-08-08)

**Verdict:** accepted

**Findings:**
/h-verify: both predictions hold. Packaging approve dial present and golden-tested; trusted live Codex full pipeline completed with no per-tool Allow prompts; export and piece artifacts present.

**Criteria met:**
- [x] Shipped Codex MCP config uses approve
- [x] Trusted Codex stops per-tool Allow prompts for pysar tools after refresh
- [x] go test asserts approve in shipped config
- [x] docs describe approve
- [x] live smoke confirmed

**Measurements:**
- go test ./cmd/pysar -run TestInitCodex|TestCodex: ok
- codex mcp get pysar: approve
- operator smoke: no MCP allow prompts; testing-what-matters-ed5e514799b0.md exported

## Impact Measurement (2026-08-09)

**Verdict:** accepted

**Findings:**
claim-001 approve dial still in asset+tests. claim-002 live trusted session not re-run; dial contract unchanged by drift.

**Criteria met:**
- [x] approve dial shipped

**Measurements:**
- claim-001: config.toml approve + init_test assert

## Impact Measurement (2026-08-09)

**Verdict:** accepted

**Findings:**
assets/codex/config.toml still approve; docs drift incidental.

**Criteria met:**
- [x] Shipped Codex config uses approve

**Measurements:**
- assets/codex/config.toml: default_tools_approval_mode = "approve"

## Impact Measurement (2026-08-10)

**Verdict:** accepted

**Findings:**
Docs drift only; approve dial still shipped.

**Criteria met:**
- [x] Shipped Codex config uses approve

**Measurements:**
- assets/codex/config.toml: approve
