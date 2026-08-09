---
id: dec-20260809-degenerate-unicode-i18n-5008f0e6
kind: DecisionRecord
version: 11
status: active
title: Unicode letters + Unicode word split (not only Fields)
context: intake
mode: standard
valid_until: 2026-11-09
created_at: 2026-08-09T08:20:27Z
updated_at: 2026-08-09T17:23:09Z
links:
  - ref: prob-20260809-976cff71
    type: based_on
  - ref: sol-20260809-6cc79891
    type: based_on
---

# Unicode letters + Unicode word split (not only Fields)

## 1. Problem Frame

**Signal:** save_intake_bundle (MCP) calls intake.Degenerate before scaffolding. Live non-English author ideas (e.g. Ukrainian/Russian sentences with multiple whitespace-separated words, zero or few Latin letters) are rejected with 'degenerate idea input — ask the author once for a real idea or draft'. Root cause already located: Degenerate counts only ASCII a-zA-Z as letters (internal/intake/bundle.go), so letters < 3 for pure Cyrillic/CJK even when the idea is a real multi-word sentence. Secondary trap: len(strings.Fields)==1 rejects scriptio continua (Japanese/Chinese without spaces). This violates dec-20260725-35fa2d24's invariant that mechanical block is only for no extractable semantic content (empty/gibberish/single stray word), not for non-Latin script. English ideas still pass.

**Constraints:**
- Keep mechanical-only gate — do not add NLP/language detection or LLM calls in Degenerate
- Preserve ask-only-when-degenerate invariant from dec-20260725-35fa2d24 — broad/vague English must still non-block
- EntryMode draft path remains exempt from Degenerate as today
- Fix must not require author to paste Latin tokens to pass validation
- Stay in intake.Degenerate (+ tests); no host-skill-only workaround as the sole fix

**Acceptance:** Unit tests: (1) a multi-word Ukrainian idea with zero Latin letters is NOT Degenerate; (2) a multi-word Japanese idea with unicode letters is NOT Degenerate (or an explicitly documented alternate rule for no-space scripts); (3) empty, whitespace-only, single Latin token gibberish like 'asdfasdfasdfasdf', and sub-3-rune input remain Degenerate; (4) existing English case 'write something about AI security for parents' still passes. go test ./internal/intake ./internal/mcpserver covering Degenerate/save path green.

## 2. Decision

**Selected:** Unicode letters + Unicode word split (not only Fields)

**Selection policy:** Long-term: maximize script_coverage then long_term_maintain; hard constraints mech_only_ok=pass and gibberish_reject>=3; impl_complexity observation-only. Declared before scoring in compare on sol-20260809-6cc79891.

**Why selected:** Operator chose V4 after compare: it is the durable fix for both failure modes in Degenerate — ASCII-only letter counting (breaks Cyrillic) and whitespace-only Fields word count (breaks no-space CJK). Stays a pure local Go heuristic aligned with dec-20260725-35fa2d24 mechanical-gate intent, without LLM/language-ID or skill-only workarounds.


**Invariants:**
- Degenerate remains mechanical-only — no LLM, network language detection, or host-skill-only sole fix
- Letter detection uses Unicode letter categories (unicode.IsLetter or equivalent), never ASCII a-z alone
- Single-token / word-count rule must not rely solely on ASCII whitespace Fields when rejecting ideas — CJK/kana/ideograph runs count as word-like units
- Empty, whitespace-only, and sub-minimum-length input stay degenerate
- EntryMode draft remains exempt from Degenerate as today
- Ask-only-when-degenerate invariant from dec-20260725-35fa2d24 preserved for broad/vague non-gibberish input

**Pre-conditions:**
- [ ] prob-20260809-976cff71 and sol-20260809-6cc79891 remain the governing frame
- [ ] Current Degenerate + TestDegenerate behavior understood (ASCII letter probe)

**Post-conditions:**
- [ ] Multi-word Cyrillic and representative no-space CJK ideas are not Degenerate
- [ ] Classic English gibberish single-token cases still are
- [ ] go test ./internal/intake (and MCP intake path if covered) green

**Admissibility:**
- NOT: NOT: Keep ASCII-only letter counting
- NOT: NOT: Skill workaround as the only fix
- NOT: NOT: External language-ID or translation service in the gate

## 3. Rationale

**Counterargument:** V1 (unicode.IsLetter + keep Fields>=2) would unclog the immediate Ukrainian/Russian failure with a tiny diff and leave CJK for later — binding V4 now may overbuild relative to today's author base if almost all non-English ideas are space-separated.

**Selected variant weakest link:** A hand-rolled Unicode word/run splitter can mis-segment edge scripts or accept/reject wrongly; complexity is higher than V1 and may need follow-up tuning when real non-spaced author ideas arrive.

**Rejected alternatives:**
| Variant | Verdict | Reason |
|---------|---------|--------|
| Unicode letters + Unicode word split (not only Fields) | **Selected** | Operator chose V4 after compare: it is the durable fix fo... |
| unicode.IsLetter + keep Fields>=2 | Rejected | Pareto stepping-stone only — leaves scriptio continua broken; operator explicitly wanted the long-term option. |
| unicode.IsLetter + rune budget instead of Fields | Rejected | Fails gibberish constraint: long single-token Latin gibberish passes a letter-count budget. |
| Remove letter-script check; empty/short/Fields only | Rejected | Weaker mechanical filter permanently; still fails no-space CJK via Fields==1. |
| Skill-only: tell agent to ignore Degenerate / pad Latin | Rejected | Fails mech_only_ok; leaves MCP broken for non-English authors. |

**Evidence requirements:**
- Extended TestDegenerate table with Cyrillic + CJK cases
- Confirm no new third-party i18n dependency unless justified in impl brief

**Predictions:**
| Claim | Observable | Threshold |
|-------|------------|-----------|
| After the change, a multi-word Ukrainian idea with zero Latin letters is not Degenerate, and a representative Japanese idea without ASCII spaces is not Degenerate | Unit tests in internal/intake calling Degenerate on fixed fixtures (Ukrainian sentence; Japanese/CJK run) | Both fixtures return false; empty, 'ai', and long single Latin gibberish token still return true |
| save_intake_bundle no longer returns degenerate idea input for those same non-English fixtures when other required fields are valid | MCP/intake integration or unit test covering callSaveIntakeBundle / Validate path with Ukrainian idea | No degenerate error for the Ukrainian fixture; English control still scaffolds |

## 4. Consequences

**Rollback plan:**
Triggers:
- Non-English false rejects persist after fix
- English gibberish false-accept regresses materially
- Segmenter complexity blocks a clean review
Steps:
1. Revert Degenerate to prior commit
2. Or fall back to V1 (unicode.IsLetter + Fields>=2) via superseding decide if CJK is deferred
Blast radius: internal/intake/bundle.go (+ tests); MCP save_intake_bundle error path for idea entry — no piece schema change

**Refresh triggers:**
- Author reports false degenerate on a real script we mis-segment
- False accept of gibberish rises after ship
- Need for a shared segmenter dependency vs hand-rolled rule

**Affected files:** internal/intake/bundle.go, internal/intake/bundle_test.go, internal/mcpserver/tools_intake.go

## Impact Measurement (2026-08-09)

**Verdict:** accepted

**Findings:**
V4 shipped in intake.Degenerate: unicode.IsLetter + ideaUnits (alphabetic runs + per-letter CJK/Hangul). Both predictions hold via unit tests. MCP path unchanged structurally — still calls Degenerate; non-English fixtures no longer trip the gate.

**Criteria met:**
- [x] Ukrainian multi-word not Degenerate
- [x] Japanese no-space not Degenerate
- [x] English gibberish single-token still Degenerate
- [x] go test intake+mcpserver green

**Measurements:**
- TestDegenerate + TestIdeaUnitsScriptioContinua PASS
- go test ./internal/intake ./internal/mcpserver -count=1 ok

## Impact Measurement (2026-08-09)

**Verdict:** accepted

**Findings:**
/h-verify: both predictions hold. Unicode letter + scriptio-continua unit split remains in place; Cyrillic and Japanese fixtures pass Degenerate/Validate; English gibberish still blocked; no third-party i18n dependency. Drift on other site/umbrella decisions is unrelated.

**Criteria met:**
- [x] Ukrainian + Japanese not Degenerate
- [x] Validate path no degenerate error for Ukrainian
- [x] Gibberish/empty still Degenerate
- [x] No new i18n deps

**Measurements:**
- TestDegenerate PASS
- TestIdeaUnitsScriptioContinua PASS
- TestValidateUkrainianIdeaNotDegenerateError PASS
- go test ./internal/mcpserver ok

## Impact Measurement (2026-08-09)

**Verdict:** accepted

**Findings:**
Drift on bundle_test.go is incidental test additions from slug/und work. claim-001/002 still PASS (TestDegenerate, TestValidateUkrainianIdeaNotDegenerateError).

**Criteria met:**
- [x] Ukrainian/Japanese not Degenerate
- [x] Validate path no degenerate for UA fixture

**Measurements:**
- claim-001: TestDegenerate PASS
- claim-002: TestValidateUkrainianIdeaNotDegenerateError PASS

## Impact Measurement (2026-08-09)

**Verdict:** accepted

**Findings:**
Drift on internal/intake/bundle.go is package-doc comment scrub only. Degenerate Unicode behavior and Validate path still hold — intake tests PASS.

**Criteria met:**
- [x] Ukrainian/CJK Degenerate contracts still hold
- [x] Validate path still green

**Measurements:**
- claim-001: Degenerate Unicode fixtures still covered by passing intake tests
- claim-002: Validate/save path still PASS; no new degenerate false positives from comment edit

## Impact Measurement (2026-08-09)

**Verdict:** accepted

**Findings:**
Incidental mcpserver stagereq adds; Degenerate unicode behavior not in diff. go test ./internal/intake PASS.

**Criteria met:**
- [x] Unicode letter/word-split Degenerate unchanged by this drift

**Measurements:**
- go test ./internal/intake PASS
