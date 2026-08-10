---
id: dec-20260810-export-resolve-citations-6afd8c31
kind: DecisionRecord
version: 1
status: active
title: Mechanical resolve-at-export from sources.md
context: pipeline
mode: standard
valid_until: 2026-11-10
created_at: 2026-08-10T08:18:17Z
updated_at: 2026-08-10T08:18:17Z
links:
  - ref: prob-20260810-70678bfc
    type: based_on
  - ref: sol-20260810-b527959e
    type: based_on
---

# Mechanical resolve-at-export from sources.md

## 1. Problem Frame

**Signal:** Operator dogfood of /ps --research --seo: final piece text still contains raw footnote-style markers such as [^getpysar-llms-txt][^getpysar-init] inline after sentences. These markers are internal research/SEO citation carriers and must not appear in author-facing export/result text. Prior Claude dogfood (note-20260809-8deb4552) already observed seo.md retaining leftover [^…] markers while the SEO stage/gate still passed — same failure class now visible in final text.

**Constraints:**
- Do not break [^shortname] as the in-piece research citation convention before SEO
- SEO-before-humanize ordering (dec-20260804-e3234e50) remains binding when --seo is set
- Hard gate for --seo currently only requires seo.md existence (dec-20260809-701b59d3) — do not silently weaken that DRR; amend/supersede if completeness must include zero markers
- Default /ps without --seo must not invent SEO packaging
- No detector-evasion / fake cleanup in humanize

**Acceptance:** After /ps --seo (and after humanize/export of that chain), the exported root .md body contains zero [^shortname] citation markers; markers are either resolved to real inline links/citations during SEO or stripped before export. Observable via grep on the exported file and on the revision humanize/export reads.

## 2. Decision

**Selected:** Mechanical resolve-at-export from sources.md

**Selection policy:** HARD after operator correction (note-20260810-63369d03): opt_in_seo_surface_ok must be true AND the recommended path must not require the author to invoke /ps-seo / discoverability packaging. Eliminate citation_honesty < 4 or export_cleanliness < 4. Among survivors maximize citation_honesty then standalone_humanize_safe. Prefer automatic clear of carriers over fail-closed that only teaches SEO. impl_coupling observation only. Operator bound V1.

**Why selected:** Operator explicitly bound V1. Export is the author-facing surface: mechanically replace every [^shortname] with a real [anchor](url) from sources.md (fail closed if unknown shortname; never invent URLs). Clears research-carrier garbage without forcing --seo packaging. Stage files (draft/sharpen/humanize) may keep markers; only the exported root file is cleaned. When SEO already ran, zero markers remain and resolve is a no-op.


**Invariants:**
- Exported root .md contains zero [^shortname] markers after a successful export
- Resolution URLs come only from the piece's sources.md / recorded research — never invented
- Missing shortname fails closed with an actionable error naming the marker
- Does not require --seo or write seo.md
- Does not rewrite draft.md / staff-edit.md / sharpen.md / humanize.md / seo.md as part of export resolve
- When seo.md already resolved links, export resolve is a no-op on markers
- --seo remains opt-in discoverability packaging (dec-20260809-701b59d3)

**Pre-conditions:**
- [ ] sources.md / research shortname→URL map is loadable for the piece
- [ ] export already selects the most-refined revision (humanize > seo > sharpen > …)
- [ ] draft.FindCitationMarkers (or equivalent) available for detection

**Post-conditions:**
- [ ] Go tests: export with markers + sources → root file has links and zero markers; unknown shortname → error
- [ ] export_piece_to_root / internal/export path performs resolve before write
- [ ] Docs/skills do not instruct SEO as the fix for research-carrier garbage

**Admissibility:**
- NOT: NOT: Tell the author to run /ps-seo just to clear research carriers
- NOT: NOT: Strip markers without linking when a source URL exists
- NOT: NOT: Invent URLs or anchors not backed by sources.md
- NOT: NOT: Silently leave [^…] in the exported root file
- NOT: NOT: Make --seo mandatory for research pieces

## 3. Rationale

**Counterargument:** Leaving markers in humanize.md means the 'final' in-piece revision still shows garbage if the author opens humanize.md instead of the export; V5 (resolve-only earlier in the chain) would clean both surfaces. Export-only resolve also creates two bodies (humanize vs export) that can confuse diffs.

**Selected variant weakest link:** Anchor-text choice without the SEO pass is heuristic (nearby words / claim span) and can produce awkward links; if sources.md lacks a URL for a marker, export fails and agents may strip markers instead of fixing research — need a clear error naming the shortname.

**Rejected alternatives:**
| Variant | Verdict | Reason |
|---------|---------|--------|
| Mechanical resolve-at-export from sources.md | **Selected** | Operator explicitly bound V1. Export is the author-facing... |
| Export fails closed if any [^shortname] remains | Rejected | Operator rejected: practical next step is run SEO; conflates citation cleanup with packaging they did not order. |
| Strip markers at export (remove carriers, leave prose) | Rejected | Destroys grounding the research pass paid for; fails citation_honesty. |
| Always run citation-resolve (SEO link step) before humanize when markers exist | Rejected | As written pulls authors into SEO-shaped stage; fails corrected opt_in_seo_surface_ok unless a separate non-SEO resolve product exists. |
| /ps --research implies citation resolve before export (orchestrator) | Rejected | Survives on Pareto but weaker standalone_humanize_safe — export-time resolve covers any path to export_piece_to_root without depending on /ps orchestration. |

**Evidence requirements:**
- Unit tests for resolve success, unknown shortname fail-closed, and no-op when no markers
- Dogfood re-export of test-ps as-an-author piece shows zero [^…] in root .md without running --seo

**Predictions:**
| Claim | Observable | Threshold |
|-------|------------|-----------|
| Successful export of a research piece with [^shortname] markers yields a root .md with zero markers and markdown links whose URLs match sources.md | Go test: fixture humanize/sharpen with markers + sources.md → export Write → grep zero [^…] and links match recorded URLs | Test PASS; no unresolved markers in destination |
| Export fails closed when a marker has no matching source shortname | Go test: marker [^missing] with empty/unknown sources → export error names missing | Error returned; no partial dirty export written |
| Export without any [^shortname] markers (including post-SEO pieces) is unchanged in citation behavior | Go test: already-resolved [anchor](url) body exports identically regarding links; no forced seo.md | PASS; no --seo required |

## 4. Consequences

**Rollback plan:**
Triggers:
- Export invents or mismatches URLs vs sources.md
- Export breaks pieces with intentional literal [^…] in code fences
- Agents start stripping markers before export to dodge fail-closed
Steps:
1. Revert export resolve helper and MCP wiring
2. Restore prior export copy-as-is behavior
3. If needed, temporarily document manual /ps-seo only as optional packaging — not as required cleanup
Blast radius: internal/export (+ shared resolve helper if extracted), export MCP tool, possibly internal/draft citation helpers, docs mentioning export cleanliness

**Refresh triggers:**
- Authors complain exported anchors are wrong/awkward
- Need resolve earlier so humanize.md itself is clean (revisit V5)
- sources.md schema change for shortname→URL

**Affected files:** internal/export/bundle.go, internal/export/write.go, internal/mcpserver/tools_export.go, internal/draft/bundle.go, internal/research/write.go, docs/pipeline.md, docs/export.md
