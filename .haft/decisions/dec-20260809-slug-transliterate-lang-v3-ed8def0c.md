---
id: dec-20260809-slug-transliterate-lang-v3-ed8def0c
kind: DecisionRecord
version: 8
status: active
title: Separate Transliterate(lang) API; Slug calls it when script ≠ Latin
mode: standard
valid_until: 2026-11-09
created_at: 2026-08-09T09:21:59Z
updated_at: 2026-08-09T14:38:12Z
links:
  - ref: prob-20260809-59175eae
    type: based_on
  - ref: sol-20260809-808d8c8f
    type: based_on
---

# Separate Transliterate(lang) API; Slug calls it when script ≠ Latin

## 1. Problem Frame

**Signal:** Operator wrote a Ukrainian-language piece; on-disk piece dir / export file resolved to template-<hex> instead of a readable Latin slug of the title/idea. Verified cause: onboarding.Slug (internal/onboarding/profile.go) only keeps ASCII [a-z0-9] and returns the literal fallback "template" when the input has no Latin letters — Cyrillic/other scripts are stripped entirely. intake.AllocateUniqueName then suffixes that fallback → template-{12-hex}. Export uses the piece directory basename as <slug>.md at project root. Same ASCII-only class of failure as intake.Degenerate before dec-20260809-degenerate-unicode-i18n-5008f0e6, now on the naming/export path. Operator wants a solid transcription/transliteration mechanism so any-language titles yield stable Latin-letter slugs.

**Constraints:**
- No third-party network call at slug time (offline, deterministic)
- Slug remains machine-facing; authors are never asked to invent a Latin slug in conversation
- Must not break rename-stable template slug contract from dec-20260719-b12539fa for already-shipped ASCII template keys (e.g. generic)
- Filesystem-safe on macOS/Linux/Windows common cases (no spaces, no path separators)

**Acceptance:** DRAFT: Given a Ukrainian title (and at least one other non-Latin script fixture, e.g. Japanese or Arabic), AllocateUniqueName / export produce a piece dir and root export filename whose readable prefix is non-empty Latin [a-z0-9-] derived from the title (not the bare fallback "template"), remains filesystem-safe, and unit tests lock Cyrillic→Latin mapping for a fixed fixture. Existing English/ASCII slug behavior for Latin-only names stays unchanged.

## 2. Decision

**Selected:** Separate Transliterate(lang) API; Slug calls it when script ≠ Latin

**Selection policy:** PASS only if readable_latin_from_ua≥3 AND determinism≥4. Among survivors maximize any_script_coverage; tie-break by lower ship_cost. wrong_romanization_risk observation-only. Operator override of compare advisory (V2): bind V3 for extensible official-per-locale transliteration.

**Why selected:** Operator bind 2026-08-09 after clarifying that V3 is the extensible contract for official/national transliteration schemes per language (UA now; other langs as standards are implemented), not a single global ad-hoc table. Compare front kept V1/V2/V3; advisory was V2 on ship_cost tie-break. Operator prioritizes locale-standard fidelity and multi-language official paths over minimal ship cost. Slug remains the ASCII collapse gate after Transliterate; authors never type Latin slugs.


**Invariants:**
- Slug/collapse remains machine-facing; authors are never asked to invent a Latin slug
- Deterministic offline: no network/LLM at slug time
- Missing/unknown lang uses documented fallback that never panics and never invents officialness
- Already-shipped ASCII template keys (e.g. generic) keep stable slugs under rename contract dec-20260719-b12539fa
- English/Latin-only inputs keep prior Slug behavior

**Pre-conditions:**
- [ ] Portfolio sol-20260809-808d8c8f compared; operator chose V3
- [ ] Slug is the shared naming gate for pieces/templates/export basename

**Post-conditions:**
- [ ] Transliterate(lang) (or equivalent) exists and Slug uses it for non-Latin script
- [ ] UA standard (or clearly named UA table documented as interim official) covers common Ukrainian letters used in titles
- [ ] Tests lock UA→Latin fixture and ASCII regression

**Admissibility:**
- NOT: LLM-proposed slugs as the primary naming path
- NOT: Requiring authors to type Latin slugs in intake conversation
- NOT: Silent wrong-lang official claims without fallback documentation

## 3. Rationale

**Counterargument:** V3 adds API and detection surface that V1 would avoid; for the immediate Ukrainian dogfood pain, a Cyrillic table inside Slug would ship faster and might be enough for months — V3 risks over-building before a second locale exists.

**Selected variant weakest link:** Language/script must be known or detected; wrong lang → wrong official romanization; languages without an implemented standard still need an honest fallback (fold/table/template) or they regress to opaque names.

**Rejected alternatives:**
| Variant | Verdict | Reason |
|---------|---------|--------|
| Separate Transliterate(lang) API; Slug calls it when script ≠ Latin | **Selected** | Operator bind 2026-08-09 after clarifying that V3 is the ... |
| Extend Slug with Unicode letters then ASCII-fold (NFKD + strip marks) | Rejected | Passes UA with zero deps but is not an official-per-locale contract; adding languages becomes table sprawl inside Slug without a named standard hook. |
| stdlib-only go-runes + ICU-style via golang.org/x/text/transform | Rejected | Compare advisory on coverage+ship tie-break; still needs hand tables for Cyrillic and does not give a first-class Transliterate(lang) surface for national standards. |
| Hash-prefix only for non-Latin; keep readable English when present | Rejected | Fails readable_latin_from_ua constraint — operator wants readable Latin from Ukrainian titles. |
| Agent-proposed Latin slug at intake (LLM) persisted as piece name | Rejected | Fails determinism constraint — not an offline stable naming contract. |

**Evidence requirements:**
- Unit tests for UA transliteration + ASCII regression
- AllocateUniqueName/export path uses new Slug behavior
- Extension point for a second locale visible in code

**Predictions:**
| Claim | Observable | Threshold |
|-------|------------|-----------|
| Ukrainian title/idea produces non-empty Latin piece-dir prefix (not bare template) via AllocateUniqueName | unit test: Slug/Transliterate on fixed UA fixture + intake name allocation | prefix matches expected Latin for fixture; no template-only basename |
| ASCII-only English names unchanged vs current Slug behavior | existing onboarding Slug tests + children-of-war style fixtures | all prior ASCII golden slugs still pass |
| API allows a second locale standard to be added without changing Slug's public collapse contract | Transliterate registry/table seam + one non-UA stub or documented extension point in code/tests | second lang hook exists (even if only UA fully filled) without Slug signature change for callers that omit lang |

## 4. Consequences

**Rollback plan:**
Triggers:
- Ukrainian fixture still yields template-* piece/export names after ship
- Transliterate API forces authors to supply lang in conversation
- Wrong-lang detection produces systematically wrong UA slugs in dogfood
Steps:
1. Revert Transliterate + Slug wiring to prior ASCII-only Slug
2. Keep unit fixtures as documentation of the gap
3. Re-bind V1 if only UA table is needed
Blast radius: internal/onboarding.Slug call sites (templates, intake AllocateUniqueName, export basename); new transliteration helper; tests — no MCP protocol change required

**Refresh triggers:**
- Second real locale (non-UA) needed in dogfood
- National UA standard revision
- SEO public URL slug diverges from piece-dir naming

**Affected files:** internal/onboarding/profile.go, internal/onboarding/profile_test.go, internal/intake/write.go, internal/intake/bundle_test.go

## Impact Measurement (2026-08-09)

**Verdict:** accepted

**Findings:**
All three declared predictions hold under CL3 unit tests run 2026-08-09. claim-001: Ukrainian AllocateUniqueName prefix is Latin (napyshy-stattiu-…), not template-*. claim-002: prior ASCII Slug goldens unchanged. claim-003: schemes registry + SupportedTransliterationLangs + unchanged Slug(name) satisfy the second-locale hook with only uk filled. Early verify relative to verify_after dates (claim-001 2026-08-23, claim-002 2026-08-16, claim-003 2026-09-09) — evidence is already available. Out of scope of predictions: DRAFT problem acceptance still mentions a Japanese/Arabic fixture; not required by post_conditions or claims; unsupported scripts still collapse to template by design.

**Criteria met:**
- [x] Ukrainian title/idea produces non-empty Latin piece-dir prefix via AllocateUniqueName
- [x] ASCII-only English names unchanged vs prior Slug behavior
- [x] Second locale hook exists without changing Slug omit-lang contract
- [x] Transliterate(lang) exists and Slug uses it for non-Latin (Cyrillic→uk)
- [x] UA CMU 55 table covers common Ukrainian title letters with locked fixtures

**Measurements:**
- claim-001: TestAllocateUniqueNameUkrainianPrefixNotTemplate PASS — prefix napyshy-stattiu-pro-zvychku-deploity-shchodnia (not template)
- claim-002: TestSlugNormalizesArbitraryNames + TestSlugNeverProducesLeadingOrTrailingDash PASS — ASCII goldens intact
- claim-003: schemes map + SupportedTransliterationLangs includes uk; Slug(name string) signature unchanged; TestSupportedTransliterationLangsIncludesUK PASS

## Impact Measurement (2026-08-09)

**Verdict:** accepted

**Findings:**
Drift from und fallback wiring (SlugLang + tests) is compatible extension: uk official path unchanged; claim-001/002/003 still hold under re-run tests. Und is labeled non-official and does not replace the Transliterate registry contract.

**Criteria met:**
- [x] UA Latin piece prefix
- [x] ASCII unchanged
- [x] Registry extensible without Slug(name) omit-lang break

**Measurements:**
- claim-001: TestAllocateUniqueNameUkrainianPrefixNotTemplate PASS
- claim-002: ASCII Slug goldens PASS
- claim-003: schemes seam still present (uk + und)

## Impact Measurement (2026-08-09)

**Verdict:** accepted

**Findings:**
Drift on internal/intake/write.go is MaxPieceNameLength comment only. UA Latin prefix, ASCII goldens, and Transliterate registry seam still hold — intake + onboarding tests PASS.

**Criteria met:**
- [x] UA Latin prefixes still work
- [x] ASCII unchanged
- [x] Registry seam intact

**Measurements:**
- claim-001: UA Slug/AllocateUniqueName still covered by green tests
- claim-002: ASCII Slug goldens still PASS
- claim-003: Transliterate(lang) registry (uk+und) unchanged
