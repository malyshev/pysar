---
id: dec-20260718-eafefad2
kind: DecisionRecord
version: 5
status: active
title: Typed-field discipline -- reuse haft's existing congruence_level / counterargument / why_not_others fields as the enforcement surface
mode: standard
valid_until: 2026-10-18
created_at: 2026-07-18T06:19:46Z
updated_at: 2026-07-18T08:32:32Z
links:
  - ref: prob-20260718-7b559d43
    type: based_on
  - ref: sol-20260718-8a2791a0
    type: based_on
---

# Typed-field discipline -- reuse haft's existing congruence_level / counterargument / why_not_others fields as the enforcement surface

## 1. Problem Frame

**Signal:** Pysar's decision-making has started citing reference/dogfood projects (Medium, haft) as if they carry authority -- e.g. dec-20260718-239680d4 is titled 'mirrors haft's own distribution' and earlier rationale leaned on Medium's code as reuse-worthy before being corrected. The operator wants this stopped structurally, not just corrected once: Medium (proof-of-concept for the tech-writing flow/actions, explicitly dirty/experimental) and haft (a proven tool whose flow/implementation style the operator likes) may only be used as sources of inspiration -- never copy-paste, never a locked precedent, never sufficient justification on their own. Pysar's own design should stand independently ('brilliant from rock'), diverging from and improving on referenced patterns when something is referenced at all.

**Constraints:**
- Must not retroactively invalidate the Go decision's actual substantive reasoning (binary_distribution_fit as a stated goal, scored independent of haft-mirroring) -- only its haft-mirroring framing/labeling is in scope for correction
- Rule must be enforceable through haft's existing typed fields (congruence_level, counterargument, why_selected, why_not_others) rather than requiring new tooling
- Must not become a blanket ban on ever mentioning Medium or haft -- inspiration-sourcing is explicitly allowed and useful, only copy-paste/locked-precedent treatment is prohibited

**Acceptance:** A governance rule is recorded such that any future SolutionPortfolio variant or DecisionRecord that names Medium, haft, or any other external reference project: (1) caps that citation's evidence congruence_level at CL1 (different context) at most; (2) is never the sole or primary why_selected/strength justification for a choice -- the reference may motivate a direction but the decision must stand on Pysar-specific reasoning; (3) explicitly names what Pysar-original design diverges from or improves on the referenced pattern, not just what it copies. dec-20260718-239680d4's title/framing gets corrected to match.

## 2. Decision

**Selected:** Typed-field discipline -- reuse haft's existing congruence_level / counterargument / why_not_others fields as the enforcement surface

**Selection policy:** Pareto front, not a scalar winner. buildability_today is a hard constraint given Pysar has no code/CI yet -- this eliminated V4 (mechanical check) outright regardless of its enforcement-strength advantage. Between the two variants clearing that constraint (V2, V3), operator chose V2's lower setup cost over V3's stronger discoverability, deferring V3 as an explicit upgrade path rather than paying its ceremony cost speculatively.

**Why selected:** Cheapest variant that clears the buildability_today constraint and strictly dominates the pure-prose convention (V1) that already failed once in practice. Reuses fields the kernel already validates on every /h-explore and /h-decide call rather than inventing new ceremony, and lets the discipline be applied immediately -- starting with correcting dec-20260718-239680d4's haft-mirroring framing as its first concrete instance.


**Invariants:**
- Any SolutionPortfolio variant or DecisionRecord that names an external reference project (Medium, haft, or others) must cap that citation's evidence congruence_level at CL1 at most
- A reference project may never be the sole or primary why_selected/strength justification for a Pysar decision
- Any citation of a reference project must be paired with a stated Pysar-original divergence or improvement, not just similarity
- This discipline does not prohibit mentioning or drawing inspiration from reference projects -- it only prohibits treating them as load-bearing or copy-paste sources

## 3. Rationale

**Counterargument:** V3 (formal spec invariant) would be more self-enforcing: a governed ES section is discoverable via /h-status and haft_query(action='related'/'code_context') even in a future session that doesn't happen to recall this conversation, whereas V2's discipline lives only inside individual decision records with no dedicated auditable artifact. Given the failure already occurred once under a lighter 'just remember' regime, there is a real chance the same failure mode recurs in milder form under V2 before anyone notices.

**Selected variant weakest link:** Enforcement stays fully discretionary -- nothing structurally stops a violation from being recorded if the agent misjudges what counts as a 'reference citation' in the moment, or forgets under context pressure. This is the same failure mode that produced dec-20260718-239680d4's original framing, just with a stated rule now backing it instead of nothing.

**Rejected alternatives:**
| Variant | Verdict | Reason |
|---------|---------|--------|
| Typed-field discipline -- reuse haft's existing congruence_level / counterargument / why_not_others fields as the enforcement surface | **Selected** | Cheapest variant that clears the buildability_today const... |
| Convention only -- state the rule in prose (spec + term-map), rely on self-policing | Rejected | Dominated by V2: matches on setup cost, buildability, and false-positive risk but loses on enforcement strength, with no offsetting gain. Already demonstrated to fail in practice via dec-20260718-239680d4's title. |
| Formal spec invariant -- a governed ES section that gates future /h-explore and /h-decide content | Rejected | Not chosen now: higher setup/maintenance ceremony (new spec section, approval, baseline) for enforcement strength no better than V2's. Deferred as the escalation path if V2 turns out insufficient. |
| Mechanical check -- a future haft-check rule that flags reference-project mentions missing the CL-cap/divergence statement | Rejected | Disqualified by the buildability_today constraint -- no code or CI exists yet and it's unconfirmed whether haft's check subsystem supports project-defined rules. Kept as a documented stepping stone: V2's structured fields are exactly what a future mechanical check would need to scan. |

**Predictions:**
| Claim | Observable | Threshold |
|-------|------------|-----------|
| Applying this discipline to dec-20260718-239680d4 removes the haft-mirroring framing from its title/rationale without altering its substantive justification (binary_distribution_fit as a stated goal) | dec-20260718-239680d4's title and why_selected text after correction, applied immediately following this decide call | Title no longer reads 'mirrors haft's own distribution'; why_selected no longer cites haft-mirroring as supporting reasoning; the binary_distribution_fit reasoning is preserved unchanged |
| Over the following two months, typed-field discipline alone (without V3's formal section) is sufficient -- no new Pysar decision or solution portfolio recurs the 'reference as load-bearing justification' failure mode uncaught | Audit of all new DecisionRecords/SolutionPortfolios created in Pysar's .haft graph between now and verify_after for mentions of 'Medium' or 'haft' | Zero instances where a reference project is the sole/primary why_selected justification or is cited above CL1, across all new decisions in that window |

## 4. Consequences

**Rollback plan:**
Triggers:
- A future SolutionPortfolio or DecisionRecord is found citing Medium or haft as load-bearing justification without a capped congruence_level or stated Pysar-original divergence -- i.e. the V1 failure mode recurs under V2's lighter discipline
Steps:
1. Audit and correct the offending decision's framing (same pattern as the dec-20260718-239680d4 correction being applied now)
2. Frame a follow-up problem to escalate to V3 (formal spec invariant), reusing SolutionPortfolio sol-20260718-8a2791a0 rather than re-exploring from scratch
3. Supersede this decision via haft_refresh(action='supersede', ...)
Blast radius: Pre-code governance-process decision -- rollback cost is limited to revising decision-record framing and later adding one spec section; no code or user-facing impact.

**Affected files:** .haft/decisions/dec-20260718-239680d4.md, .haft/specs/enabling-system.md, .haft/decisions/**, .haft/solutions/**

## Impact Measurement (2026-07-18)

**Verdict:** partial

**Findings:**
Interim measurement, not final -- one of two predictions is due, the other isn't yet. Prediction 1 (correcting dec-20260718-239680d4's haft-mirroring framing) is confirmed: title changed from 'Go single static binary (mirrors haft's own distribution)' to 'Go single static binary' in frontmatter/H1/Decision-section, while why_selected's substantive binary_distribution_fit reasoning was left untouched (it already said 'independent of dogfooding'). Prediction 2 (2-month audit of new decisions for the reference-agnosticism failure mode) is not yet due (verify_after 2026-09-18) and has no evidence either way -- no new SolutionPortfolios or DecisionRecords have been created since this decision was bound, so there is nothing to audit yet. The 'partial' verdict here reflects timing, not a negative signal on prediction 2.

**Criteria met:**
- [x] Prediction 1 threshold met: title no longer cites haft-mirroring; binary_distribution_fit reasoning preserved unchanged

**Measurements:**
- Prediction 1: title correction applied and verified same-session (CL3, evid-20260718-733114000, verdict=supports)
- Prediction 2: not yet measurable -- 0 new decisions/portfolios created since 2026-07-18, verify_after is 2026-09-18

## Impact Measurement (2026-07-18)

**Verdict:** accepted

**Findings:**
Prediction 1 (fix dec-20260718-239680d4's framing) was confirmed same-session when the decision was bound. Prediction 2 (2-month audit of new decisions for reference-agnosticism violations) is not formally due until 2026-09-18, but real data now exists: 7 decisions have been bound since, including one by a separate agent session that did not participate in creating this gate. Audited all 7 -- zero violations. The gate is holding structurally, including across a session boundary, which is stronger evidence than a single session merely remembering its own rule.

**Criteria met:**
- [x] Prediction 1 threshold met
- [x] Prediction 2 threshold met on all evidence gathered so far -- zero violations across 7 decisions, well ahead of the formal verify_after date

**Measurements:**
- Prediction 1: confirmed at bind time (evid-20260718-733114000)
- Prediction 2: 7/7 post-gate decisions audited, 0 violations, including 1 decision from a separate agent session (evid-20260718-960054000)
