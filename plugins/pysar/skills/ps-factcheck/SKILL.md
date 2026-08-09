---
name: ps-factcheck
description: |
  Reusable claim-grounding and citation-accuracy check for any Pysar writing
  command. Decides when a claim needs real verification (genuine uncertainty,
  or naming a specific person/framework/standard), fetches real sources when
  it does, and — separately — audits already-written prose against what was
  actually fetched so a citation never says more than its source supports.
  Never overrides the author's own thesis or take; only grounds supporting
  facts and examples. Invoked by other Pysar writing skills (ps-intake today;
  future drafting/research/editing passes), not typically run standalone.
when_to_use: |
  Another Pysar writing skill is about to state a factual/technical claim, or
  has already drafted prose making one, and needs it grounded or checked
  against real sources before finalizing. Also directly invocable via
  /ps-factcheck on an existing piece of prose if the operator wants an
  ad hoc citation-accuracy check.
argument-hint: "[optional: path or pasted text to check]"
allowed-tools: WebSearch WebFetch Read
---

# ps-factcheck — ground claims, then verify they were grounded correctly

Two separate, sequential jobs. Both use whatever **expert lens** (the
discipline/practitioner viewpoint governing this content's claims — e.g.
"experimental science," "horticulture," "distributed-systems security") the
calling skill already determined; if none was supplied, pick the obvious one
from the content itself before starting.

## Job 1 — ground claims before they're written

For each claim about to be stated, check it against two independent
triggers — they catch different failure modes, not the same thing twice:

1. **Accuracy uncertainty.** Stable, foundational knowledge (well-established
   semantics, mechanisms unlikely to have changed) → state it directly, no
   research; this is the common case. Genuinely uncertain (version-specific
   behavior, fast-moving specifics, a benchmark/statistic, a mechanism you'd
   be guessing at) → a small number of targeted searches/fetches (2-4, not a
   bibliography) to verify before stating it as fact.
2. **Named-authority claims — check this even when confident on accuracy.**
   About to invoke a specific named person, framework, methodology, standard,
   or study — being *correct* about it isn't the trigger, *naming* it is.
   Default to verifying it with one real fetch. Naming a real concept and
   declining to ground it is the weaker output, not the safer one — a named
   framework is usually the sharpest, most differentiating insight available,
   not decoration to trim under time pressure. Only drop the name entirely
   (no fetch needed) when it's genuinely decorative — when removing it costs
   the reader nothing the rephrased version doesn't already say.

Rules for both triggers:

- Never let this override the author's own thesis, angle, or stated take —
  it grounds supporting facts and examples only.
- Do this silently — no "let me research this" narration.
- Read for grounding, not for citation-display — this is not building a full
  bibliography.
- Track every real source fetched as `{url, note}` — `note` is one line on
  exactly what it confirmed. This list is the output the calling skill
  persists (e.g. into a `sources` field).

## Job 2 — verify what got written matches what got fetched

After the calling skill drafts its content using Job 1's grounding, re-check
the *written prose* against the *actual fetched sources*. This is a separate
check because drift happens here even when Job 1 was done correctly: a
fetched source can say one thing, and the written summary/note can drift from
it even when the final prose ends up accidentally correct (or incorrect).
Citing a real URL is not the same as citing it *accurately*.

- **Re-read what you fetched, not your memory of it.** For each `{url, note}`
  pair, confirm `note` states only what the source actually says. A rounded,
  paraphrased, or assumed number is not the source's number — if they differ,
  fix the note (or the claim), don't let them silently diverge.
- **Trace every research-flavored claim back to a source.** Anything in the
  drafted content that reads as backed by research — a specific number,
  date, named study, or phrase like "per X guidance" — must actually be
  anchored to a `{url, note}` pair from Job 1. If nothing supports it, either
  fetch something that does or drop the implied authority and state it as
  the writer's own reasoning.
- **Don't let a citation cover more than it covers.** If a source only
  addresses one case and the drafted content also makes a claim about a
  related but different case the source never mentions, don't let the
  citation imply it covers both. Say so honestly, or fetch something that
  actually supports the second case.
- This is a self-audit of existing work, not a new research round — no new
  fetches unless the check surfaces a real gap.

## Do not

- Invent a citation or URL that wasn't actually fetched
- Let research override or "improve" the author's own thesis, angle, or
  stated take
- Research by default for stable, well-established knowledge
- Strip a named framework/authority/standard to avoid the one-fetch cost of
  grounding it, when that name was the content's sharpest insight
- Write a source `note` from memory of what was skimmed instead of
  re-checking the actual fetched content
- Let a citation imply broader coverage than what its source actually
  addresses

## Output contract

A list of `{url, note}` pairs for every source actually fetched and used,
each verified against Job 2's checklist — hand this to whatever the calling
command uses to persist citations (e.g. Pysar's `save_intake_bundle` →
`sources` field). Empty is a legitimate, common output — most claims don't
need grounding.
