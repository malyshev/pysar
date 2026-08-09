---
id: dec-20260809-cursor-marketplace-v1-dual-discovery-8b748a7a
kind: DecisionRecord
version: 8
status: active
title: In-monorepo Cursor Plugin package (skills+MCP), binary preinstalled
context: cursor-marketplace
mode: standard
valid_until: 2026-11-09T00:00:00Z
created_at: 2026-08-09T09:20:30Z
updated_at: 2026-08-09T17:23:01Z
links:
  - ref: prob-20260809-d4ae8601
    type: based_on
  - ref: sol-20260809-b964ff5f
    type: based_on
---

# In-monorepo Cursor Plugin package (skills+MCP), binary preinstalled

## 1. Problem Frame

**Signal:** Operator wants Pysar published to the Cursor Marketplace. Cursor already supports Pysar via pysar init --cursor (skills under ~/.cursor/skills + project/user mcp.json pointing at ${userHome}/.local/bin/pysar). Live dogfood on Cursor 3.15 showed Customize Connected lists Plugin MCPs preferentially; Phase A of dec-20260809-cursor-cold-path-phased-v1v5-then-v2-de2fc11a shipped spawnable path + one-step enable locus. Phase B (Marketplace / install-link) is committed but has no ProblemCard and no in-repo plugin package: no .cursor-plugin/plugin.json, no plugin-root skills/mcp layout, no local ~/.cursor/plugins/local dogfood, no marketplace submission. note-20260809-df7e8ede tracks Phase B as open through 2026-11-09.

**Constraints:**
- Preserve hostAdapter + one host-neutral skill corpus (dec-20260808-f3001106); do not fork editorial skill bodies for marketplace
- Piece persistence remains via pysar MCP tools, not agent raw Write/Bash of .pysar/**
- Claude and Codex init scaffolds stay green; marketplace work must not break pysar init --claude/--codex goldens
- Non-developer author cold path: no MCP Logs or Customize filter archaeology as the primary product answer
- Marketplace plugin must not ship secrets; any plugin variables declare names only
- Do not put CLAUDE.md / haft / site engineering harness into the end-author plugin package
- Phase B of dec-20260809-cursor-cold-path-phased-v1v5-then-v2-de2fc11a must not be silently dropped — waive with rationale if postponed past valid_until

**Acceptance:** On a clean machine with the documented pysar binary install (~/.local/bin/pysar), an author installs the Pysar Cursor plugin (Marketplace listing after review, or equivalent local load from ~/.cursor/plugins/local that uses the same package shape), enables it once in Customize, opens any writing project, runs /ps with a one-sentence idea, and gets a provisional piece under .pysar/pieces/*/brief.md written via pysar MCP — with zero hand-edits to mcp.json and without MCP Logs / workspace-filter archaeology. Plugin package also passes Cursor's submission checklist (valid .cursor-plugin/plugin.json, relative logo, README covering install + components) and is submitted at cursor.com/marketplace/publish.

## 2. Decision

**Selected:** In-monorepo Cursor Plugin package (skills+MCP), binary preinstalled

**Selection policy:** PASS only if connected_plugin_mcp=pass. Among survivors maximize skill_corpus_integrity; if tied, prefer lower marketplace_review_risk (long-term shippability), then higher author_cold_path_completeness; final tie-break lower distribution_ops_burden. Operator bind 2026-08-09: V1 is the sole canonical package; Cursor Marketplace and getpysar.com Install in Cursor deeplink are dual discovery/install channels onto that same package (note-20260809-1ce7b2b0) — not rival products. V3 shim is not a second SKU; if pursued later it upgrades inside the V1 carrier.

**Why selected:** Operator bind via /h-decide: V1 + dual discovery. In-monorepo Cursor Plugin is the durable Connected-plugin packaging with one skill corpus and shippable review profile. Marketplace listing and website Install in Cursor deeplink both expose that same V1 package, widening discovery without forking packaging (recovers Phase B install-link half of dec-20260809-cursor-cold-path without a second product). Binary remains preinstalled at ${userHome}/.local/bin/pysar for MCP spawn.


**Invariants:**
- Exactly one canonical Cursor Plugin package in the Pysar monorepo (plugins/ or equivalent); Marketplace and site deeplink MUST point at that same package shape
- One host-neutral skill corpus remains source of truth (dec-20260808-f3001106); plugin skills are projections/sync of that corpus, not a fork
- Piece persistence remains via pysar MCP tools
- Plugin mcp.json spawns documented binary path (${userHome}/.local/bin/pysar or equivalent); no silent bare-command-only contract for Dock GUI
- Claude and Codex init scaffolds stay green
- V3 install shim is not shipped as a separate Marketplace product/SKU
- Do not put CLAUDE.md / haft / site engineering harness into the end-author plugin package

**Pre-conditions:**
- [ ] sol-20260809-b964ff5f compared; operator chose V1 + dual discovery
- [ ] Phase A Cursor path rewrite (dec-20260809-cursor-cold-path) remains the spawn contract for binary path
- [ ] Documented binary install path remains ~/.local/bin (or init-resolved equivalent)

**Post-conditions:**
- [ ] In-repo Cursor Plugin package exists with valid .cursor-plugin/plugin.json, skills/, mcp.json, logo, plugin README
- [ ] Package dogfooded via ~/.cursor/plugins/local before Marketplace submit
- [ ] Package submitted to cursor.com/marketplace/publish (or equivalent Cursor review path)
- [ ] getpysar.com exposes an Install in Cursor deeplink/button that installs/enables the same V1 package (not a divergent mcp.json-only paste)
- [ ] pysar init --cursor does not create a competing skill dual-home that drifts from the plugin corpus (demote global skill sync or document plugin as primary Cursor skill carrier)
- [ ] Phase B note note-20260809-df7e8ede satisfied by this WIP/decision rather than silent drop

**Admissibility:**
- NOT: NOT: Separate Marketplace product vs separate deeplink product with divergent package contents
- NOT: NOT: Skills-only Marketplace listing as the Phase B answer
- NOT: NOT: V3 shim as a second named product/SKU alongside V1
- NOT: NOT: MCP Logs / Customize filter archaeology as primary author path
- NOT: NOT: Forking editorial skill bodies for the plugin

## 3. Rationale

**Counterargument:** Shipping Marketplace and a site deeplink against one package still leaves binary preinstall as a second step — V3's bootstrap would close more of the non-developer cold path in one product surface. Separately, a thin dedicated plugin repo (V2) might clear Cursor review faster than a monorepo marketplace.json entry, delaying listing while we perfect dual discovery in-tree.

**Selected variant weakest link:** Authors who install the plugin (Marketplace or deeplink) without the pysar binary still cannot spawn MCP; Connected discoverability is fixed but cold path remains two-step (binary + plugin enable). Also: if Cursor install-link format does not load the same plugin package as Marketplace, dual discovery silently forks into two carriers.

**Rejected alternatives:**
| Variant | Verdict | Reason |
|---------|---------|--------|
| In-monorepo Cursor Plugin package (skills+MCP), binary preinstalled | **Selected** | Operator bind via /h-decide: V1 + dual discovery. In-mono... |
| Marketplace plugin with install shim that ensures binary then serves | Rejected | Not rejected as forever-impossible; rejected as a second product. Auto-download under MCP spawn is high review/sandbox risk. If ever pursued, it upgrades inside the V1 carrier — same package, not a rival SKU. |
| Separate public plugin repo vendoring skills+MCP | Rejected | Dominated on skill_corpus_integrity and ops: cross-repo sync lag with no Connected gain over in-monorepo V1. |
| Install-link / deeplink distribution without Marketplace listing first | Rejected | Demoted from rival packaging to a discovery channel on V1. Deeplink alone without the V1 package does not satisfy Connected Plugin MCP acceptance; with V1 it is part of dual discovery, not an alternative end-state. |
| Skills-only Marketplace plugin; MCP stays init/user mcp.json | Rejected | Fails connected_plugin_mcp constraint — does not put pysar in Customize Connected as Plugin MCP. |

**Evidence requirements:**
- Local ~/.cursor/plugins/local dogfood: Connected + /ps → brief.md
- cursor.com/marketplace/publish submission record or listing URL
- Site CTA/deeplink points at same plugin identity as manifest name
- Init --cursor golden/docs updated for plugin-primary Cursor skill carrier

**Predictions:**
| Claim | Observable | Threshold |
|-------|------------|-----------|
| Local load of the V1 package from ~/.cursor/plugins/local makes pysar appear as a Plugin MCP in Customize Connected and can write a provisional piece via /ps after binary is installed | Copy/symlink plugins/pysar (or chosen path) into ~/.cursor/plugins/local; Dock Cursor; enable plugin; /ps one-sentence idea → .pysar/pieces/*/brief.md via MCP | success without hand-editing mcp.json and without MCP Logs archaeology |
| Marketplace submission uses the same V1 package directory that local dogfood used (no second packaging tree) | Submitted repo/path and local plugin root share the same .cursor-plugin/plugin.json + skills + mcp.json content hash lineage | single package path documented; no parallel plugins/pysar-deeplink tree |
| Site Install in Cursor control targets the same V1 package (deeplink or Cursor install URL), not a hand-rolled divergent mcp.json | getpysar.com Cursor install CTA href/docs resolve to V1 package install, verified against plugin manifest name | one plugin name/identity across Marketplace listing and site CTA |

## 4. Consequences

**Rollback plan:**
Triggers:
- Local dogfood fails twice: plugin does not appear as Connected Plugin MCP or /ps cannot use MCP tools with binary present
- Cursor review rejects the monorepo plugin layout with no acceptable marketplace.json fix within one iteration
- Site deeplink installs a divergent MCP/skills set from the V1 package
- Init --cursor + plugin dual-home causes skill drift incidents in dogfood
Steps:
1. Unpublish or withdraw Marketplace submission if live
2. Remove or quarantine plugins/ package from default docs CTAs
3. Revert site Install in Cursor CTA to prior host guidance
4. Keep Phase A init --cursor + user mcp path as sole Cursor cold path
5. Waive or supersede this DRR with alternate bind if V2 thin repo becomes mandatory for review
Blast radius: plugins/ (or chosen plugin dir), site Cursor install CTA, docs Cursor journey, possibly cmd/pysar Cursor hostAdapter skill-sync behavior; Claude/Codex unchanged if invariants hold

**Refresh triggers:**
- Cursor changes plugin install-link or Marketplace packaging format
- Marketplace review rejects monorepo/subdirectory plugin source
- Documented binary install path moves off ~/.local/bin
- Dual discovery channels diverge in package contents

**Affected files:** plugins/pysar/.cursor-plugin/plugin.json, plugins/pysar/mcp.json, plugins/pysar/skills, plugins/pysar/README.md, plugins/pysar/assets/logo.svg, .cursor-plugin/marketplace.json, cmd/pysar/host.go, cmd/pysar/assets/cursor/mcp.json, cmd/pysar/assets/skills, docs/init.md, docs/mcp-and-skills.md, site/src/components/home-integrations.tsx, site/src/lib/site.ts

## Impact Measurement (2026-08-09)

**Verdict:** partial

**Findings:**
claim-002 single V1 package + plugin tests PASS. claim-003 site CTA identity points at plugins/pysar. claim-001 live local plugin /ps smoke not re-run — blocks full accept.

**Criteria met:**
- [x] single package path
- [x] site discovery identity

**Criteria NOT met:**
- [ ] fresh local plugin Connected /ps smoke

**Measurements:**
- claim-002/003: package + site CURSOR_PLUGIN identity
- claim-001: no fresh live plugin smoke

## Impact Measurement (2026-08-09)

**Verdict:** partial

**Findings:**
Local dogfood prediction holds after real-directory copy into ~/.cursor/plugins/local/pysar (symlink rejected by Cursor). Plugins lists Pysar; MCPs shows pysar Plugin tag green/connected; editorial workflow produced first full draft (~1200 words) and continued to editorial pass; Try in chat overview worked. Marketplace submit and site full-plugin Install deeplink still open.

**Criteria met:**
- [x] Local ~/.cursor/plugins/local dogfood: Connected Plugin MCP + writing workflow via skills/MCP
- [x] Same V1 package shape used for dogfood (plugins/pysar content)

**Criteria NOT met:**
- [ ] Package submitted to cursor.com/marketplace/publish
- [ ] getpysar.com Install in Cursor as same full plugin install path (not only MCP deeplink)

**Measurements:**
- note-20260809-839bb4dd dogfood success
- note-20260809-dbc1ffcf symlink rejection
- Plugins UI: Pysar installed
- MCPs: pysar Plugin connected
- draft ~1200 words then editorial pass
- Try in chat responded with Pysar overview + workflow intent

## Impact Measurement (2026-08-09)

**Verdict:** accepted

**Findings:**
site.ts drift incidental. claim-001 closed by note-20260809-839bb4dd local dogfood (Plugin MCP connected; editorial workflow via plugin skills+MCP). claim-002/003: single plugins/pysar package identity matches site CTA.

**Criteria met:**
- [x] Local Plugin MCP dogfood without hand mcp.json
- [x] Single V1 package
- [x] Site CTA same identity

**Measurements:**
- claim-001: local plugin dogfood succeeded (note-20260809-839bb4dd)
- claim-002: single package path plugins/pysar
- claim-003: site cursorPlugin.name=pysar

## Impact Measurement (2026-08-09)

**Verdict:** accepted

**Findings:**
Only incidental docs/mcp-and-skills.md Persistence-rule edit. Plugin package invariants untouched.

**Criteria met:**
- [x] Single Cursor plugin package shape
- [x] Skill corpus not forked by this drift

**Measurements:**
- docs/mcp-and-skills.md +4/-3 require_piece_stages mention only
