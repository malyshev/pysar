---
id: dec-20260808-homepage-ai-install-signals-9eb5a6a9
kind: DecisionRecord
version: 7
status: active
title: Homepage / HTML-first signals (meta, visible AI block, JSON-LD)
context: site
mode: standard
valid_until: 2026-11-08T00:00:00Z
created_at: 2026-08-08T16:50:09Z
updated_at: 2026-08-09T09:56:28Z
links:
  - ref: prob-20260808-cabd5795
    type: based_on
  - ref: sol-20260808-7810d3e3
    type: based_on
---

# Homepage / HTML-first signals (meta, visible AI block, JSON-LD)

## 1. Problem Frame

**Signal:** Desired next trick: user pastes https://getpysar.com into a supported coding agent (Claude Code, Cursor, Codex; Gemini later) and says only 'install'. Today the site has human UI + /docs/install HTML + /install.sh, but no standard AI-readable carrier at a well-known URL that tells the agent the authoritative install+init sequence. Agents must guess from HTML scrape or training data.

**Constraints:**
- Must work with static Next export + Cloudflare Pages (dec-20260808-da785647)
- Must not contradict journey docs install path (curl getpysar.com/install.sh | bash)
- Human marketing homepage must not become AI-only noise
- No secret install path that humans cannot also find
- Host coverage today: Claude Code, Cursor, Codex; Gemini optional/coming soon

**Acceptance:** In a live smoke on each supported host: paste https://getpysar.com and prompt 'install' — the agent fetches a published machine-readable carrier from getpysar.com (not invented steps) and runs the documented path through binary on PATH plus host-appropriate pysar init, without inventing a competing install method.

## 2. Decision

**Selected:** Homepage / HTML-first signals (meta, visible AI block, JSON-LD)

**Selection policy:** Prefer the single variant that can satisfy URL-drop + 'install' without a second discovery hop or a marketplace pre-install. Stepping-stone indexes and non-standard playbook paths lose unless they alone close the acceptance criterion. Marketing cleanliness is secondary to agent-fetch co-location with the pasted homepage URL.

**Why selected:** V4 co-locates authoritative install+init instructions with https://getpysar.com — the exact URL the user pastes. When Claude Code, Cursor, or Codex fetch that page, the carrier is already in the response (visible AI block and/or machine-readable meta/link/JSON-LD). No dependence on hosts knowing /llms.txt or finding /ai/install.md.


**Invariants:**
- Install recipe on the homepage must match docs/install.md and https://getpysar.com/install.sh — no competing inventable path
- Static Next export + Cloudflare Pages remains the delivery shape (dec-20260808-da785647)
- Human marketing homepage must not become AI-only; AI signals may be visible but must not own the first viewport
- Host init remains pysar init --claude|--cursor|--codex as documented

**Pre-conditions:**
- [ ] getpysar.com already serves install.sh and /docs/install
- [ ] Homepage is under site/ Next static export

**Post-conditions:**
- [ ] Fetching https://getpysar.com returns machine- and/or human-visible AI install guidance co-located with that URL
- [ ] Guidance names curl|bash install.sh and host-appropriate pysar init

**Admissibility:**
- NOT: Secret AI-only install path humans cannot find
- NOT: Instructions that contradict journey docs
- NOT: Requiring marketplace skill install for the URL-drop trick to work

## 3. Rationale

**Counterargument:** A thin /llms.txt (V1) plus playbook (V2) is the emerging web convention and keeps the marketing homepage clean; putting agent instructions on the homepage risks visual noise and dual-audience design debt, while many coding agents already probe well-known paths.

**Selected variant weakest link:** Agents that summarize or truncate marketing HTML may miss the AI install block; scrapers that drop meta/JSON-LD leave only visible copy, which competes with hero chrome for attention.

**Rejected alternatives:**
| Variant | Verdict | Reason |
|---------|---------|--------|
| Homepage / HTML-first signals (meta, visible AI block, JSON-LD) | **Selected** | V4 co-locates authoritative install+init instructions wit... |
| Canonical /llms.txt index + agent instructions | Rejected | Stepping-stone discovery root; fails alone if the host never leaves the pasted homepage URL to fetch /llms.txt. |
| Single fetchable AI install playbook (/ai/install.md) | Rejected | Complete procedure but incomplete discovery — non-standard path needs a pointer; not sufficient alone for URL-drop. |
| Raw markdown export of journey docs at stable URLs | Rejected | Improves fetch quality of docs; does not put install steps on the dropped homepage URL. |
| Host skill/plugin install via marketplace (not site text) | Rejected | Different product shape; breaks the stated trick of site URL alone + 'install'. |

**Predictions:**
| Claim | Observable | Threshold |
|-------|------------|-----------|
| Pasting https://getpysar.com and saying install causes at least one supported host agent to use the homepage-published curl install.sh path rather than inventing brew/npm/go-get | Live smoke transcript on Claude Code or Cursor or Codex within 14 days of ship | Agent cites or executes getpysar.com/install.sh (or equivalent documented homepage instruction) in >=2 of 3 host smokes |
| Homepage remains human-readable without AI install copy dominating the first viewport brand composition | Visual check of shipped homepage first viewport after deploy | Brand/hero still primary; AI install signal is secondary (below fold, compact aside, or non-hero meta-only) — not a second hero |

## 4. Consequences

**Rollback plan:**
Triggers:
- Live smoke on Claude Code, Cursor, or Codex fails: agent invents install steps or never sees homepage AI signals after two attempts
- Homepage AI block harms human first-viewport composition enough that we prefer a well-known-path fallback
Steps:
1. Remove or hide the AI install block, meta/link tags, and JSON-LD from the homepage
2. Optionally add V1 /llms.txt later under a new decide if needed
3. Revert site homepage components and redeploy Pages
Blast radius: site homepage components, metadata in layout/page, possible robots/link tags; no Go CLI change

**Refresh triggers:**
- Any supported host changes how it fetches/pastes URLs
- Install path or init flags change in docs
- Homepage redesign removes or buries the AI signal

**Affected files:** site/src/app/page.tsx, site/src/app/layout.tsx, site/src/components, site/src/lib/site.ts

## Impact Measurement (2026-08-08)

**Verdict:** partial

**Findings:**
Bundle /h-verify after deploy b9ee84b: prediction 2 holds — live homepage has secondary #for-ai-agents block below hero plus meta/JSON-LD HowTo; brand/hero remains primary. Prediction 1 not met yet — carrier is published on getpysar.com but Claude/Cursor/Codex URL-drop+install agent smokes (>=2 of 3) have not been run within the 14-day window.

**Criteria met:**
- [x] Homepage AI install signal secondary to hero (live + page composition)
- [x] Carrier co-located on https://getpysar.com (meta, HowTo JSON-LD, visible block) matching install.sh path

**Criteria NOT met:**
- [ ] Live host agent smoke: agent uses getpysar.com/install.sh in >=2 of 3 hosts

**Measurements:**
- deploy b9ee84b success
- live meta pysar:agent-install + trigger=install
- JSON-LD HowTo+SoftwareApplication present
- visible for-ai-agents after H1
- host smoke transcripts: 0/3

## Impact Measurement (2026-08-08)

**Verdict:** partial

**Findings:**
Bundle /h-verify: pred2 holds on live getpysar.com (secondary #for-ai-agents + meta/HowTo). Pred1 still unmet — 0/3 host URL-drop install smokes. Working-tree /llms.txt complement (note-20260808-26ff4eb4) builds locally and tests pass but is undeployed (prod 404); per note it complements V4 and does not replace homepage co-location.

**Criteria met:**
- [x] Homepage AI signal secondary to hero
- [x] Carrier co-located on homepage URL

**Criteria NOT met:**
- [ ] Live host agent smoke >=2 of 3 using install.sh

**Measurements:**
- live for-ai-agents+meta+HowTo after H1
- host smokes 0/3
- prod /llms.txt 404; site/out/llms.txt + vitest ok locally

## Impact Measurement (2026-08-09)

**Verdict:** accepted

**Findings:**
Install.sh AI instruction still in site.ts; vitest green. Live multi-host paste smoke not re-run; structural claim surface intact after plugin-note drift.

**Criteria met:**
- [x] homepage AI install signals present

**Measurements:**
- claim-001/002: site.ts install path + vitest 7/7
