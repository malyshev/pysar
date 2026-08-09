---
name: ps-staff-edit
description: |
  Senior-editor pass on an existing draft.md: stakes, brief alignment,
  failure modes and honest scope, close memorability, technical sanity,
  water/filler removal, and readability. Writes the revision to a separate
  staff-edit.md --
  draft.md is never touched, so the first draft and the edited version can
  be compared (or shown side by side in a future UI), the same reasoning
  that keeps brief/outline/angles/draft as separate files instead of one
  evolving document. Runs after /ps-draft. Stakes, brief alignment, failure
  modes, honest scope, close memorability, technical sanity, water/filler,
  and readability (including consumability and sentence rhythm) live here
  as this pass's own checks. Preserves [^shortname] citation markers
  (never resolves them to links); never touches
  brief.md/outline.md/angles.md/sources.md/draft.md. Agent-agnostic: no
  author registry, no platform-specific publish ceremony, no cross-article
  corpus-variance scanning.
when_to_use: |
  Operator runs /ps-staff-edit, asks "is this argument sound?", or wants a
  senior-editor pass on an existing draft before calling it done.
argument-hint: "[@path/to/piece or path, e.g. .pysar/pieces/<name>/]"
allowed-tools: mcp__pysar__save_staff_edit_bundle mcp__pysar__read_author_defaults Read
---

# /ps-staff-edit — stakes, spine, honest scope

**Author-facing outcome: "Stakes, spine, honest scope."** This pass answers
*"Is this argument sound?"* — not by resolving citations or chasing
platform-engagement metrics, but by pressure-testing the draft's substance
and producing a revised version where it's thin. If the operator asks
whether the argument holds up, this is that command.

This is the **last check before the piece is presentable** — in the sense a
real staff editor means it: not a style nitpick pass, the point where
someone with editorial judgment reads the whole thing end to end and either
signs off or fixes it. Readability is part of that signoff, not a separate
later command: a piece that's substantively sound but reads like a washing-
machine manual — flat, enumerated, no rhythm — hasn't actually passed staff
edit. Readability is part of this signoff, not deferred to a later
engagement-only pass.

**The revision goes to `staff-edit.md`, not `draft.md`.** Every other piece
file is named for its content, not overwritten in place as it evolves —
`brief.md`, `outline.md`, `angles.md`, `draft.md` each stay put once
written. This pass follows the same convention: `draft.md` survives exactly
as `/ps-draft` left it, so the operator (or a future UI) can compare the
first draft against the edited version, not just trust that an in-place
edit was an improvement.

This skill is **host-agnostic**. Prefer pysar MCP tools for anything
mechanical (profile defaults, citation checking, disk writes, run-log). Do
**not** use Bash/`ls`/`cat` or raw Write for piece files.

Do not invent an author registry or a platform publish checklist. This
project has no author registry — `voice.md`/`style.md` (via
`read_author_defaults`, the same tool `/ps-draft` and `/ps-intake` use) is
the whole mechanism, including its own `banned_phrases` — reuse those,
don't hardcode a different list. No cross-piece corpus-variance scanning
(opener/closing forbiddance windows against recent pieces) — no such
mechanism exists here; edit each piece on its own merits.

## Args

A path to an existing piece (e.g. `.pysar/pieces/docker-for-developers-
abc123/`) — never call this a "slug." Referencing it via Claude Code's
native `@` file-picker is normal flow, same as `/ps-draft` and
`/ps-research` — `@` commonly resolves to a specific file inside the piece
rather than the directory itself; pass whatever it resolves to straight
through as `piece_path`.

**If invoked with no path at all**, don't guess which piece is meant.
Show this, then stop and wait:

> **Pressure-tests a draft's substance** — stakes, brief-alignment,
> failure modes, honest scope — before calling the argument sound.
>
> **To start:**
> - `/ps-staff-edit @.pysar/pieces/docker-for-developers-abc123/`
>
> If this is part of a longer run, `/ps` handles the whole pipeline in
> one go instead of stage by stage.

If `draft.md` doesn't exist yet for the given piece, tell the operator to
run `/ps-draft` first — do not invent one.

## Step 1 — load context

Read the piece's `brief.md`, `outline.md`, `angles.md`, and `sources.md`
once. `brief.md`'s `Thesis`, `Promise to the reader`, `Killer sections`,
and `Non-goals` are the contract this pass checks the draft against. Call
`read_author_defaults` for the operator's own `banned_phrases` and
register/tone — don't invent or reuse a different list.

**Which file to revise from:** if `staff-edit.md` already exists (a prior
staff-edit run), read that — it's the latest revision, and this pass keeps
moving it forward. Otherwise read `draft.md` — this is the first staff-edit
pass for this piece.

## Step 2 — the checks

Apply in order. Earlier checks inform later ones. Produce a fully revised
version in your own context as you go — this is a **rewrite pass**, not a
comment pass; when a check finds a real problem, fix it in the text, don't
just note that it exists. The revision is submitted whole in Step 3, never
edited into an existing file directly (see Hard rules).

1. **Stakes and reader thesis.** After the first screen (opener + first
   couple paragraphs), could a skeptical reader name the specific mistake
   or decision this piece helps them avoid, without reading further? If
   not, tighten the opener so the stake is concrete and falsifiable — no
   added throat-clearing to get there.
2. **Brief alignment — signal vs. noise.** Cross-read the draft against
   `brief.md`'s `Thesis`/`Promise`/`Killer sections`/`Non-goals`. Cut or
   demote anything that doesn't serve the promise or thesis. A
   `killer_section` must read like the competitive edge the brief says it
   is, not like equal-weight boilerplate.
3. **Failure modes, trade-offs, and honest limits.** A piece that's mostly
   happy-path tutorial needs failure modes, footguns, or trade-offs in
   proportion to how much happy path it runs — drawn from the draft's own
   sources and domain consistency, never invented war stories. If the
   piece is tutorial-heavy and this is thin, borrow space from redundant
   steps rather than padding the piece.
4. **Honest scope.** Check the opener and any summary for implied
   guarantees the brief or the cited sources don't actually support. Tighten
   language so scope matches what the material actually establishes.
5. **Close memorability.** The close should land one sharp point the reader
   remembers tomorrow, not restate the thesis two or three different ways.
   If it currently does, compress to the strongest single line.
6. **Technical sanity.** Skim code fences, config keys, CLI commands, and
   version numbers against `sources.md` and the draft's own internal
   consistency. Fix an obvious inconsistency only when the fix is a
   deletion or a wording tightening — never fabricate behavior to make
   something consistent.
7. **Water and filler.** For every sentence, ask: does removing this lose
   the reader information, nuance, or a concrete example they need — or
   does it just restate a point already made, hedge without narrowing
   anything real, or pad the piece toward length without adding substance?
   Cut anything that fails the test. Concrete signs: two sentences saying
   the same thing in different words; a hedge that doesn't actually qualify
   anything ("in many cases," "it's worth noting," "generally speaking"
   used as filler, not genuine qualification); a transition that exists
   only to announce the next section instead of adding content; an example
   or point restated in slightly different phrasing later in the same
   section. **This is not a length target** — a genuinely earned example,
   trade-off, or failure mode stays even if it makes the piece longer; the
   test is whether each sentence earns its place, not the word count. A
   piece running long because it covers more real ground is fine. A piece
   running long because the same three points are each said twice is not —
   cut it here, don't carry senseless length forward to a pass that doesn't
   exist yet.
8. **Readability.** Read the piece the way a reader would, straight
   through. Vary sentence rhythm — a run of same-length declarative
   sentences reads like an instruction manual even when every sentence is
   individually correct. Keep one idea per paragraph; split any paragraph
   carrying two genuinely distinct ideas. (If a paragraph's second half
   turns out to be restating its first half rather than a separate idea,
   that's water — check 7 should already have cut it, not left it here to
   split.) This is the "make it pretty, not a washing-machine manual"
   check — substance can be perfect and a piece can still fail this one.

## Hard rules

- **`[^shortname]` markers preserved.** Don't convert, remove, or relocate
  a marker away from the claim it anchors. If you split a sentence that
  carried a marker, keep the marker on the clause that still needs it.
  Resolving markers to links is not this pass's job (nor any pass's yet).
- **No new factual claims, stats, versions, or names** not already implied
  by the draft or `sources.md`. You may delete unsupported fluff; you may
  rephrase toward honesty ("the docs don't promise X") when the draft or
  sources already ground that. If a real gap needs a genuinely new source,
  say so and suggest `/ps-research` — don't invent the citation yourself.
- **Preserve verbatim** code blocks, commands, error messages, and version
  pins shown in code — byte-identical, even when tightening surrounding
  prose.
- **No banned phrases** from the operator's own `style.md`/`voice.md`.
- **Never touch `brief.md`, `outline.md`, `angles.md`, `sources.md`, or
  `draft.md`.** This pass only ever writes `staff-edit.md` — never edit
  `draft.md` in place, even to "fix" it after seeing the revision; the
  point of a separate file is that `draft.md` stays exactly what `/ps-draft`
  produced.
- **Never use `Edit`, `Write`, or Bash on any piece file — not even for a
  small, surgical ("delta") change.** This is the one `ps-*` pass that
  revises existing prose rather than writing brand-new content, which makes
  `Edit` the tool that *looks* like the natural fit here — it is not. Even
  a one-sentence tweak still goes through `save_staff_edit_bundle` with the
  **full** revised text: read the current file into context, make the
  change there, submit the whole thing back. The tool re-validates citation
  integrity the same way `/ps-draft`'s own save does, so a slip here can't
  silently ship a broken or invented citation — a direct `Edit` bypasses
  that check entirely, which is the actual reason this rule exists, not a
  stylistic preference.

## Step 3 — persist (MCP)

Call `save_staff_edit_bundle` with `piece_path`, `revised_md` (the **full**
revised text — this replaces the whole `staff-edit.md` file, not a diff,
and never touches `draft.md`), and `checks` (one line per real change, e.g.
`"[stakes] named the specific failure mode in paragraph 2"`,
`"[readability] split the three-clause sentence in the compose section"`).
If a check genuinely required no change, record that as a single line too
(`"no changes needed -- <which checks, why>"`) rather than omitting
`checks` entirely — an edit pass that logs nothing isn't a completed pass.
`mode` is optional and informational (`"delta"` for surgical, `"rewrite"`
for structural).

On tool error: fix exactly what it names. If it's a missing `draft.md`,
tell the operator to run `/ps-draft` first — don't create one yourself.

## Step 4 — summary (short, plain language for the author)

**"Stakes, spine, honest scope."** — then 2-3 sentences describing what
actually changed in the piece, in plain language, for the author, not for
another editor. Say what changed in the prose itself ("sharpened the
opening so it names a concrete problem instead of staying abstract,"
"tightened two paragraphs that were each trying to do too much at once")
— never recite this skill's own internal check names back to the operator
(no "brief alignment," "failure-mode coverage," "citation integrity,"
"close memorability" — Step 2's numbered checks are this pass's own working
process, not vocabulary the author needs translated back to them).

**Don't enumerate what needed no change.** If three of the eight checks
found nothing to fix, say nothing about those three — silence on a
non-issue is the right amount of information, per this project's own
output discipline (CLAUDE.md): explaining an absence is noise, not signal.

End with the saved path (`staff-edit.md`) — nothing else. **Don't explain
that `draft.md` was left untouched, or why** — that's this pass's own
internal storage choice (see the note above Step 1), not a feature the
author asked about or that this project promises them; narrating it every
run is exactly the kind of noise this project's own output discipline
(CLAUDE.md) exists to cut. If an author later asks where their original
draft went, answer then — don't pre-explain it on every save. No phase
jargon, no permission theater. Stop.

## Do not

- Recite Step 2's internal check names (stakes, brief alignment,
  failure-mode coverage, honest scope, close memorability, technical
  sanity, water/filler, readability) back to the operator — describe what
  changed in the prose itself, in plain language
- Cut a genuinely earned example, trade-off, or failure mode just to
  shorten the piece — water/filler removal targets sentences that add no
  information, not length itself
- Enumerate which checks needed no change — silence on a non-issue, not a
  line explaining the absence
- Use `Edit`, `Write`, or Bash on any piece file — always
  `save_staff_edit_bundle` with the full revised `revised_md`
- Overwrite or edit `draft.md` — it stays untouched; the revision goes to
  `staff-edit.md`
- Resolve `[^shortname]` markers into links — not this pass's job
- Invent a new fact, statistic, version, or citation not already grounded
  in the draft or `sources.md`
- Hardcode a banned-phrase list — read the operator's own `voice.md`/
  `style.md` via `read_author_defaults`
- Build or reuse an author registry, or scan a cross-piece corpus for
  opener/closing variance — neither mechanism exists in this project
- Write a platform-specific status field or publish checklist — this project
  has neither
- Touch `brief.md`, `outline.md`, `angles.md`, or `sources.md`
- Treat readability or water/filler removal as optional or defer either to
  a pass that doesn't exist — both are this pass's own checks
- Call a piece path a "slug"
