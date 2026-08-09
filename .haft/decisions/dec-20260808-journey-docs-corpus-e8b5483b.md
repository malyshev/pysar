---
id: dec-20260808-journey-docs-corpus-e8b5483b
kind: DecisionRecord
version: 15
status: active
title: Journey-first docs/ tree with SSG-ready frontmatter
context: docs
mode: standard
valid_until: 2026-11-08T00:00:00Z
created_at: 2026-08-08T11:37:16Z
updated_at: 2026-08-09T09:56:28Z
links:
  - ref: prob-20260808-39f48c3d
    type: based_on
  - ref: sol-20260808-084f2966
    type: based_on
---

# Journey-first docs/ tree with SSG-ready frontmatter

## 1. Problem Frame

**Signal:** Operator intent (note-20260808-bf54fd6d): Pysar user guidance is concentrated in README.md and host-specific scaffolds; there is no dedicated in-repo user-documentation tree that covers each usage aspect (install, init/hosts, pipeline stages, MCP/skills, export, troubleshooting) as readable, task-oriented pages. Desired dual use: the same Markdown corpus must serve GitHub readers today and dogfood a future Pysar documentation website readable outside GitHub — without a second parallel docs source later. Path name, IA detail, and static-site stack are not yet chosen.

**Constraints:**
- Single corpus: no parallel 'website docs' source that diverges from the repo tree
- User-facing Markdown versioned with the product (not only chat/README dump)
- Document shipped behavior only — do not invent CLI/MCP/pipeline features
- Out of scope for this problem: building/hosting the website, SSG choice, marketing pages, Go API/internal reference, i18n
- Must not contradict active install/init/host decisions (e.g. Cursor/Claude host adapters, GoReleaser distribution)

**Acceptance:** DRAFT — edit if wrong: (1) A dedicated docs tree exists in the repo with a clear entry index that navigates install → init/hosts → pipeline stages → MCP/skills → export → troubleshooting. (2) Each of those aspects has at least one task-oriented page a first-time operator can follow without reading Go source. (3) Page layout/front-matter (or equivalent convention) is stable enough that a future static site can ingest the tree without rewriting IA. (4) README.md points to that tree as the user-guide entry (README stays short; docs carry depth). (5) Cold-path check: install + init for one supported host is completable from the docs alone.

## 2. Decision

**Selected:** Journey-first docs/ tree with SSG-ready frontmatter

**Selection policy:** Prefer the variant that (1) satisfies draft acceptance as written (dedicated docs tree, journey coverage, short README pointer, cold-path install+init), (2) keeps a single Markdown corpus with stable slug/section/nav_order for future site ingest without choosing an SSG now, (3) minimizes new runtime/CI product surface. Operator-bound after explore without a scored compare — policy declared as acceptance-fit + deferred-site readiness + low infra cost.

**Why selected:** V1 maps directly onto the framed acceptance and the original note intent: human-readable journey pages under docs/ today, with frontmatter that a future website can ingest without rewriting IA. It avoids premature SSG lock-in (V4), Diátaxis authoring overhead (V3), generator/CI as a second product (V2), and acceptance failure of a root hub (V5).


**Invariants:**
- Single user-docs corpus in-repo — no parallel website-only docs tree that diverges
- Document shipped behavior only — do not invent CLI/MCP/pipeline features
- Journey IA remains the primary nav (install → init/hosts → pipeline → MCP/skills → export → troubleshoot) until a superseding decision
- README is entry/pointer, not a second full user guide
- SSG/hosting choice remains out of scope of this decision — frontmatter must stay tool-agnostic enough to adapt

**Pre-conditions:**
- [ ] ProblemCard prob-20260808-39f48c3d and portfolio sol-20260808-084f2966 remain the governing frame
- [ ] Shipped install/init/host behavior is readable from README, release notes, and host adapters (no docs invention)

**Post-conditions:**
- [ ] docs/ exists with journey index and per-aspect task pages
- [ ] README points to docs/ as user-guide entry
- [ ] Pages carry consistent YAML frontmatter for future ingest

**Admissibility:**
- NOT: Parallel docs source for a future website
- NOT: Expanding README as the only deep guide while claiming docs/ is complete
- NOT: Committing an undeployed SSG as a required part of this decision
- NOT: Documenting unreleased CLI/MCP behavior

## 3. Rationale

**Counterargument:** Without a mechanical drift gate (V2) or Diátaxis type discipline (V3), journey docs will rot after the next host/CLI change and the 'site-ready' frontmatter will be cargo-cult — so binding V1 now buys structure without guaranteeing the acceptance cold-path stays true.

**Selected variant weakest link:** Hand-maintained journey pages and nav_order drift from shipped install/init/host behavior, so the cold-path 'docs alone' acceptance fails silently while the tree still looks complete on GitHub.

**Rejected alternatives:**
| Variant | Verdict | Reason |
|---------|---------|--------|
| Journey-first docs/ tree with SSG-ready frontmatter | **Selected** | V1 maps directly onto the framed acceptance and the origi... |
| Docs-as-code generation from shipped surfaces | Rejected | Protects reference drift but does not deliver first-time task narrative; generator+CI cost is premature before a human corpus exists. Can layer later on V1 for reference slices. |
| Diátaxis tree under docs/ (tutorials · how-to · explanation · reference) | Rejected | Stronger long-term IA for mixed audiences, but higher authoring/cross-link cost and cold-reader bounce risk before content exists; journey IA matches acceptance path more directly. |
| Docs tree plus undeployed SSG scaffold | Rejected | Proves ingest only if the scaffold is exercised; risks premature SSG dialect lock-in and bitrot while frame puts SSG choice out of scope. |
| Root hub — expand README or DOCUMENTATION.md, defer docs/ | Rejected | Fails draft acceptance (tree + per-aspect pages + short README pointer) unless the ProblemCard is rewritten; postpones the dual-use file boundaries the website dogfood needs. |

**Evidence requirements:**
- Cold-path install+init walkthrough transcript or checklist against docs alone
- Inventory that each acceptance aspect has a page
- Spot-check frontmatter key uniformity across docs pages

**Predictions:**
| Claim | Observable | Threshold |
|-------|------------|-----------|
| A docs/ tree with index + pages for install, init/hosts, pipeline, MCP/skills, export, and troubleshooting will exist, and README will point to it as the user-guide entry | Presence of docs/index (or equivalent) linking those sections; README contains a clear docs entry link and is not the sole deep guide | All six aspects have at least one task-oriented page; README deep content is relocated or linked, not duplicated as the only source |
| Cold-path: install + init for one supported host is completable from the docs alone without reading Go source | Operator or agent walkthrough using only docs/ (+ release install instructions therein) for Claude or Cursor init | Walkthrough reaches successful pysar init for that host with zero consultation of .go files |
| Frontmatter (title, slug, nav_order, section) is uniform enough that a future static site can map nav without rewriting page IA | Sample of docs pages share the same frontmatter keys; slugs/nav_order form a total order under the journey sections | ≥90% of docs pages use the declared keys; no page requires path renames solely to satisfy a hypothetical SSG sidebar |

## 4. Consequences

**Rollback plan:**
Triggers:
- Cold-path install+init cannot be completed from docs alone after the corpus ships
- Frontmatter/slug convention forces a full IA rewrite when a real SSG is chosen
- Docs tree and README diverge into two conflicting user guides
Steps:
1. Revert README pointer changes
2. Delete or archive the docs/ tree (or restore pre-decision README content from git)
3. If only IA failed: flatten or remount pages under a revised layout without inventing a second corpus
4. Re-open problem or supersede this decision if a different variant is chosen
Blast radius: README.md and docs/** only — no runtime binary, MCP, or release pipeline change required by this decision

**Refresh triggers:**
- First real docs website / SSG choice
- New host adapter or material init UX change
- Cold-path walkthrough fails
- Operator reports docs/README divergence

**Affected files:** README.md, docs/, docs/index.md, docs/install.md, docs/init.md, docs/pipeline.md, docs/mcp-and-skills.md, docs/export.md, docs/troubleshooting.md

## Impact Measurement (2026-08-08)

**Verdict:** accepted

**Findings:**
All three predictions held. Docs tree + README pointer present with six aspect pages and journey index. Frontmatter keys uniform on 7/7 pages with total nav_order. Cold-path Cursor init succeeded from docs alone (existing PATH binary confirmed per install.md; full curl install not re-run). No material gaps vs acceptance; residual risk remains content drift over time (decision weakest_link), not structure.

**Criteria met:**
- [x] Dedicated docs tree with journey index linking install→init→pipeline→MCP/skills→export→troubleshooting
- [x] Each aspect has a task-oriented page
- [x] Frontmatter uniform for future static ingest (≥90%)
- [x] README points to docs as user-guide entry
- [x] Cold-path install confirm + Cursor init completable from docs alone

**Measurements:**
- docs pages: 7 (index + 6 aspects); README links docs/index.md
- frontmatter key coverage: 100% (threshold ≥90%); nav_order 0..60 step 10
- cold-path: pysar init --cursor exit 0 in scratch dir; .pysar/project + .cursor/mcp.json + global ps skill; zero .go reads

## Impact Measurement (2026-08-08)

**Verdict:** accepted

**Findings:**
All three DRR predictions still hold after motherhome-related drift. Journey docs corpus intact (7/7 pages, README pointer, 100% frontmatter keys). Cold-path Cursor init from docs steps succeeds in scratch. Drift is incidental site/deploy churn from dec-20260808-da785647; docs/*.md unchanged. Single-corpus invariant strengthened: getpysar.com/docs ingests this tree without a parallel website-only docs source. Re-baseline after measure.

**Criteria met:**
- [x] Dedicated docs tree with journey index and six aspect pages
- [x] README points to docs as user-guide entry
- [x] Frontmatter uniform for static ingest (≥90%)
- [x] Cold-path install confirm + Cursor init from docs alone
- [x] Single corpus: no parallel website-only docs tree

**Measurements:**
- docs inventory: index + 6 aspects; README → docs/index.md — PASS
- frontmatter title/slug/nav_order/section: 7/7 (100%) — PASS
- cold-path: pysar init --cursor scratch exit 0; .pysar + .cursor/mcp.json + skills — PASS
- drift: site/motherhome + README baseline hash; docs/ clean vs HEAD — incidental
- dual-use: https://getpysar.com/docs and /docs/install → 200 from same corpus — PASS (bonus vs out-of-scope SSG)

## Impact Measurement (2026-08-08)

**Verdict:** accepted

**Findings:**
/h-verify: all three predictions hold. Tree+README pointer intact; frontmatter 100%; cold-path Cursor and Codex init succeed from docs steps without .go reads. Drift is same-corpus content updates (Codex approve dial) plus site packaging — incidental/expected refresh, not dual-corpus divergence. Soft note: docs/index journey blurb still says Claude or Cursor (Codex covered in init/mcp pages) — cosmetic IA copy lag, not prediction failure.

**Criteria met:**
- [x] docs tree with index + six aspect pages; README pointer
- [x] cold-path init from docs alone
- [x] frontmatter uniform ≥90%

**Measurements:**
- 7/7 pages; frontmatter 100%; nav_order 0..60
- README → docs/index.md
- init --cursor and --codex scratch exit 0
- getpysar.com/docs + /docs/install 200

## Impact Measurement (2026-08-08)

**Verdict:** accepted

**Findings:**
/h-verify: all three predictions hold. Journey corpus + README pointer intact; frontmatter 100%; cold-path Cursor and Codex init succeed from docs steps without .go reads; live /docs ingest 200. Drift (2 modified + 18 added) is same-corpus host/docs updates plus site packaging — expected refresh, not dual-corpus divergence. Soft: index intro/step2 still omit Codex while What you need lists it.

**Criteria met:**
- [x] docs tree + six aspect pages + README pointer
- [x] cold-path install confirm + Cursor/Codex init from docs alone
- [x] frontmatter ≥90% uniform (100%)

**Measurements:**
- 7/7 pages with title/slug/nav_order/section; nav_order 0..60
- README → docs/index.md
- pysar 0.3.0 init --cursor/--codex scratch exit 0
- getpysar.com/docs* all 200

## Impact Measurement (2026-08-08)

**Verdict:** accepted

**Findings:**
Bundle drift verify: all three predictions hold. Docs corpus/frontmatter/README pointer intact; live /docs 200. Drift is site homepage AI + .pnpm-store — not docs IA.

**Criteria met:**
- [x] docs tree + README pointer
- [x] frontmatter 100%
- [x] cold-path prior + docs unchanged

**Measurements:**
- fm 7/7
- README→docs/index.md
- /docs /docs/install 200

## Impact Measurement (2026-08-08)

**Verdict:** accepted

**Findings:**
All three predictions hold. install.md inspectable path is same-corpus content improvement; frontmatter/README/live docs intact.

**Criteria met:**
- [x] docs tree + README
- [x] frontmatter 100%
- [x] cold-path install+init

**Measurements:**
- fm 7/7
- live /docs /docs/install 200
- install.md adds download-then-bash

## Impact Measurement (2026-08-09)

**Verdict:** accepted

**Findings:**
Six journey pages + index still present; out/docs builds. Wording drift incidental.

**Criteria met:**
- [x] journey docs corpus

**Measurements:**
- claim-001/002/003: docs tree + built routes
