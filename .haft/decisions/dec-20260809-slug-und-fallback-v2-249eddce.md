---
id: dec-20260809-slug-und-fallback-v2-249eddce
kind: DecisionRecord
version: 4
status: active
title: Add go-unidecode / similar any-script approx as fallback scheme und
context: slug-i18n
mode: standard
valid_until: 2026-11-09
created_at: 2026-08-09T09:29:26Z
updated_at: 2026-08-09T09:33:08Z
links:
  - ref: prob-20260809-19bd7e90
    type: based_on
  - ref: sol-20260809-2a547783
    type: based_on
---

# Add go-unidecode / similar any-script approx as fallback scheme und

## 1. Problem Frame

**Signal:** dec-20260809-slug-transliterate-lang-v3-ed8def0c shipped Transliterate(lang) with only uk (CMU 55) filled. Cyrillic titles now get readable Latin piece dirs; Japanese and Arabic titles still strip to empty ASCII and AllocateUniqueName yields template-{hex}. Original ProblemCard DRAFT acceptance for prob-20260809-59175eae explicitly wanted at least one other non-Latin script fixture (e.g. Japanese or Arabic); that gap remains. Operator 2026-08-09: also need Japanese/Arabic. Extend the V3 registry — do not reopen the UA contract unless a variant requires superseding it.

**Constraints:**
- Offline deterministic — no network/LLM at slug time
- Preserve Transliterate(lang) registry + Slug(name) omit-lang contract from dec-20260809-slug-transliterate-lang-v3-ed8def0c
- Authors never invent Latin slugs in conversation
- Unknown/unsupported script must not invent officialness; document honest limits (especially Japanese kanji)
- Zero or explicit third-party deps only if compare selects them; prefer stdlib when quality holds
- Must not break uk CMU 55 or ASCII goldens

**Acceptance:** Fixed Japanese fixture (kana and/or mixed as scoped by the chosen variant) and fixed Arabic fixture each produce AllocateUniqueName / Slug prefixes that are non-empty Latin [a-z0-9-] derived from the title (not bare template), offline and deterministic, with unit tests locking the mapping; ASCII and uk behavior remain unchanged.

## 2. Decision

**Selected:** Add go-unidecode / similar any-script approx as fallback scheme und

**Selection policy:** PASS only if ja_readable_prefix≥3 AND ar_readable_prefix≥3 AND offline_determinism≥4 AND no_false_officialness≥4 (und must be labeled non-official). Among survivors maximize any_script_coverage; tie-break by lower ship_cost vs kanji dictionary. Operator direct bind of V2 without formal /h-compare scoring — prioritizes JA+AR+future scripts in one stroke under V3 registry.

**Why selected:** Operator bind 2026-08-09: cover Japanese and Arabic (and other non-Latin scripts) without waiting for per-locale official tables or a Japanese dictionary. Keeps uk as the official CMU 55 scheme; und is an explicit non-official any-script approximator used when detectTransliterationLang finds no named scheme (or as last resort after named schemes leave non-Latin). Satisfies DRAFT gap from prior problem while preserving never-invent-officialness via labeling.


**Invariants:**
- uk remains the official CMU 55 scheme; und must not claim officialness
- Offline deterministic — no network/LLM at slug time
- Slug(name) omit-lang contract preserved; authors never invent Latin slugs
- ASCII-only and uk behaviors unchanged
- Unknown named lang stays identity; und is the labeled any-script fallback only

**Pre-conditions:**
- [ ] dec-20260809-slug-transliterate-lang-v3-ed8def0c Transliterate registry exists
- [ ] Portfolio sol-20260809-2a547783; operator chose V2

**Post-conditions:**
- [ ] und (or equivalent) registered and invoked when non-Latin remains after named schemes / detection miss
- [ ] JA and AR fixtures locked in tests with non-template Latin prefixes
- [ ] Dependency chosen and pinned; license acceptable for Pysar distribution

**Admissibility:**
- NOT: Presenting und output as a national/official transliteration standard
- NOT: Requiring authors to type Latin slugs
- NOT: Breaking uk or ASCII slug goldens
- NOT: Network call at slug time

## 3. Rationale

**Counterargument:** V1 (hand ja+ar) keeps zero deps and clearer standards story for Arabic; V3 dictionary is better Japanese quality. V2 trades fidelity and dep hygiene for breadth, and may look like inventing romanization if und is not loudly non-official in code/docs.

**Selected variant weakest link:** Unidecode-class output for kanji/Arabic is approximate and library-defined — wrong-looking but stable Latin slugs; dep license/churn sits on the cold naming path.

**Rejected alternatives:**
| Variant | Verdict | Reason |
|---------|---------|--------|
| Add go-unidecode / similar any-script approx as fallback scheme und | **Selected** | Operator bind 2026-08-09: cover Japanese and Arabic (and ... |
| Hand-register ja + ar schemes (kana Hepburn + Arabic letter table); kanji omitted honestly | Rejected | Fails kanji-heavy Japanese titles (still template); operator asked for Japanese coverage that kana-only cannot deliver. |
| Vendor Japanese dictionary romanizer (e.g. kakasi/kuroshiro-class) + Arabic table | Rejected | Higher fidelity for JA but heavier embed/deps and wrong-reading risk; operator chose breadth approximator over dictionary. |
| Script-tiered: Arabic full table + Japanese kana now; kanji → short content hash suffix with kana/latin stem when any exists | Rejected | Kanji-only names stay opaque (hash); V2 gives approximate readable Latin instead. |
| Require profile/project language tag; SlugLang(lang) only — no script detection for ja/ar | Rejected | Does not by itself implement JA/AR Latinization; adds author/config cold-path without closing template gap. |

**Evidence requirements:**
- JA+AR unit tests via Slug/AllocateUniqueName
- uk+ASCII regression green
- und labeled non-official in code

**Predictions:**
| Claim | Observable | Threshold |
|-------|------------|-----------|
| Fixed Japanese fixture (incl. kanji and/or kana) yields non-empty Latin Slug/AllocateUniqueName prefix, not template | unit test on fixed JA title through Slug and AllocateUniqueName | prefix matches locked Latin approx for fixture; basename not template-only |
| Fixed Arabic fixture yields non-empty Latin Slug/AllocateUniqueName prefix, not template | unit test on fixed AR title through Slug and AllocateUniqueName | prefix matches locked Latin approx for fixture; basename not template-only |
| und is labeled non-official and uk CMU 55 + ASCII Slug goldens remain unchanged | code/docs comment or scheme naming + existing uk/ASCII tests | uk and ASCII tests PASS; und documented as approximator not national standard |

## 4. Consequences

**Rollback plan:**
Triggers:
- JA or AR fixture still yields template-* after ship
- und path presented or documented as a national/official standard
- Chosen library license or maintenance blocks release
- uk or ASCII Slug goldens regress
Steps:
1. Remove und scheme / Slug fallback call
2. Keep uk scheme and prior tests
3. Re-bind V1 or V3 if dogfood needs official AR or dictionary JA
Blast radius: internal/onboarding transliterate registry + Slug fallback; go.mod dep; tests — no MCP protocol change

**Refresh triggers:**
- Library abandoned or license change
- Dogfood demands dictionary-quality Japanese
- National ja/ar schemes later preferred over und

**Affected files:** internal/onboarding/transliterate.go, internal/onboarding/transliterate_test.go, internal/onboarding/profile.go, internal/intake/bundle_test.go, go.mod, go.sum

## Impact Measurement (2026-08-09)

**Verdict:** accepted

**Findings:**
All three predictions hold under CL3 unit tests 2026-08-09 (early vs verify_after Aug 16/23). claim-001: JA kana+kanji fixtures yield locked Latin Slug/AllocateUniqueName prefixes via und. claim-002: AR fixture mql-n-l-dt (lossy but non-template). claim-003: und labeled non-official in schemes/docs; uk CMU 55 and ASCII goldens PASS; Київ still kyiv not und Kiyiv. Library pinned mozillazg/go-unidecode v0.2.0 (note-20260809-7007c9d3).

**Criteria met:**
- [x] Japanese fixture non-empty Latin prefix not template
- [x] Arabic fixture non-empty Latin prefix not template
- [x] und labeled non-official; uk+ASCII unchanged
- [x] Offline deterministic und via go-unidecode

**Measurements:**
- claim-001: TestSlugJapaneseAndArabicNotTemplate + AllocateUniqueName JA cases PASS — korehatesutodesu / dong-jing-noxi-guan
- claim-002: AR AllocateUniqueName prefix mql-n-l-dt- PASS
- claim-003: TestTransliterateUKOfficialExamples + TestSlugNormalizesArbitraryNames + und non-official comments PASS
