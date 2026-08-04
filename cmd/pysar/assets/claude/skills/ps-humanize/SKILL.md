---
name: ps-humanize
description: |
  Final voice-authenticity pass on an existing draft (seo.md if present,
  else sharpen.md, else staff-edit.md, else draft.md): strips genuine
  AI-tell patterns -- hedge stacks, stock transitions, throat-clearing,
  too-clean section/paragraph/bullet symmetry, uniform sentence rhythm --
  that read as generic regardless of who or what wrote them. Writes the
  revision to a separate humanize.md -- earlier stages are never touched,
  same reasoning as staff-edit, sharpen, and seo. When seo.md exists, its
  resolved [anchor](url) links are preserved verbatim -- this pass edits
  prose around them, never inside them. Adapted from a prior writing POC's
  text-humanizer discipline, deliberately excluding the parts of that
  POC's own aggressive framing (and a since-declined prompt with the same
  shape) that serve no reader and exist only to defeat AI-detection
  classifiers: no deliberate grammar-breaking, no leaving thoughts
  incomplete, no "embracing messiness" as a goal, no skipping explanations
  for effect. That POC's own default mode already rules these out; this
  skill doesn't implement them even as an opt-in. The goal is voice
  authenticity, not statistical fingerprint scrambling, and it doesn't
  require making the prose worse to read. Agent-agnostic: reuses the
  operator's own voice.md/style.md banned_phrases instead of a hardcoded
  list, no Medium-specific curator-eligibility framing.
when_to_use: |
  Operator runs /ps-humanize, says "it reads like AI," or wants a final
  voice-authenticity pass before calling a piece done.
argument-hint: "[@path/to/piece or path, e.g. .pysar/pieces/<name>/]"
allowed-tools: mcp__pysar__save_humanize_bundle mcp__pysar__read_author_defaults Read
---

# /ps-humanize — sounds like you, not the machine

**Author-facing outcome: "Sounds like you, not the machine."** If the
operator says "it reads like AI," this is that command. Everything
upstream (`/ps-staff-edit`, `/ps-sharpen`) already checked substance and
reader experience; this pass checks something narrower and more
mechanical: does the prose carry the generic statistical fingerprints of
model output — hedge stacks, stock connective tissue, suspiciously uniform
structure — that make a piece sound like nobody in particular wrote it.

**What this is not.** This is not an AI-detector-evasion tool, and it will
not become one. A prior writing POC's own HUMANIZE phase — and a draft
prompt this project considered and declined — both lean toward techniques
whose only function is confusing detection classifiers: deliberately
broken grammar, incomplete thoughts, "embracing messiness." None of that
serves a reader; broken grammar reads worse regardless of who wrote it.
This skill fixes things that are genuinely bad writing independent of
source — a hedge stack (*"may potentially"*) is weak whether a nervous
human or a model wrote it — and stops there. If the operator's own voice
is genuinely mannered, hedged, or fragment-heavy, that's what `voice.md`
already encodes; this pass never fights the operator's real voice, only
generic-model tells that aren't anyone's voice.

**The revision goes to `humanize.md`, not an earlier-stage file.** Same
convention as every pass since `/ps-draft`: `brief.md`, `outline.md`,
`angles.md`, `draft.md`, `staff-edit.md`, `sharpen.md`, `seo.md` (and its
`seo-checklist.md`) all stay put once written.

This skill is **host-agnostic**. Prefer pysar MCP tools for anything
mechanical (profile defaults, citation checking, disk writes, run-log). Do
**not** use Bash/`ls`/`cat` or raw Write for piece files.

## Args

A path to an existing piece — never call this a "slug." Referencing it via
Claude Code's native `@` file-picker is normal flow, same as the earlier
passes.

If `draft.md` doesn't exist yet, tell the operator to run `/ps-draft`
first — do not invent one.

## Step 1 — load context

Read the piece's `brief.md` and `sources.md` once (for citation context).
Call `read_author_defaults` for the operator's own `banned_phrases` and
voice — the same tool and profile every other pass uses. Don't build or
reuse a separate banned-phrase list; `voice.md`/`style.md` is the whole
mechanism.

**Which file to revise from:** read `seo.md` if it exists (the piece went
through `/ps-seo` -- its resolved `[anchor](url)` links must survive this
pass untouched, see Hard rules), else `sharpen.md` (the latest, most-
refined revision), else `staff-edit.md`, else `draft.md`.

## Step 2 — the checks

Apply in order. Produce a fully revised version in your own context as you
go — this is a rewrite pass, not a comment pass. The revision is submitted
whole in Step 3, never edited into an existing file directly (see Hard
rules).

1. **Stock transitions and throat-clearing.** Generic connective tissue —
   "however," "moreover," "notably," "ultimately," "in conclusion" — used
   as filler rather than because two ideas actually contrast or build.
   Density is the tell, not any single use: a genuine "however" contrasting
   two real claims is fine; six paragraphs each opening on one is not. Cut
   or replace with the actual logical relationship (cause, objection,
   example) stated in plain words, or just delete it and let the
   paragraph break do the work.
2. **Hedge stacks.** Two or more hedges qualifying the same claim
   ("may potentially," "might generally," "could possibly"). A person
   commits to one hedge or none. Drop the weaker one. If both feel
   load-bearing, the claim itself is too vague — tighten the claim instead
   so it only needs one qualifier, or none.
3. **Synthetic-enthusiasm and banned phrases.** Cross-check against the
   operator's own `voice.md`/`style.md` `banned_phrases` (from Step 1) --
   these already cover most of this ("leverage," "seamless," "cutting-edge"
   and whatever else the operator has flagged). Fix every hit.
4. **Section/paragraph/bullet symmetry.** Read the structure the way a
   skim-reader would. Every `##` section with exactly the same number of
   sub-points, every list with exactly the same item count, every
   paragraph nearly the same length, 3+ consecutive bullets sharing the
   same opening shape and length — these read as a template, even when
   each individual sentence is fine. Vary at least one place where the
   uniformity isn't earned by the content itself: drop a sub-point that
   doesn't pull weight, add one where a section is genuinely deeper, merge
   two bullets that are really saying the same thing.
5. **Sentence-rhythm variance.** Distinct from `/ps-staff-edit`'s
   readability check (which fixes *overloaded* sentences/paragraphs) —
   this checks for *uniformity*: a run of sentences that are all roughly
   the same length, even when none of them individually carries too much.
   Vary by splitting or combining at a natural boundary — never by
   cramming extra clauses into a sentence just to make it longer, and
   never by leaving a sentence incomplete to make it shorter.

## Hard rules

- **`[^shortname]` markers preserved** (only relevant when `/ps-seo` did
  not run, so markers are still present). Don't convert, remove, or
  relocate a marker away from the claim it anchors. Resolving markers to
  links is not this pass's job.
- **Resolved `[anchor](url)` links preserved verbatim** (only relevant
  when `/ps-seo` did run). Edit the prose around a link -- never split,
  rewrite, or relocate the anchor text or URL, and never add a new link of
  your own. That's `/ps-seo`'s job, not this pass's.
- **No new factual claims, stats, versions, or names.** This pass changes
  phrasing and structure, never substance.
- **Preserve verbatim** code blocks, commands, error messages, and version
  pins shown in code.
- **No banned phrases** from the operator's own `style.md`/`voice.md`.
- **Never deliberately break grammar, leave a thought incomplete, or
  introduce a non-sequitur "for authenticity."** These read as errors, not
  humanity — see "What this is not" above. If a genuine person would
  actually write a sentence fragment here for effect (not to seem
  imperfect, but because it lands), that's normal craft judgment, not a
  goal to manufacture.
- **Never touch `brief.md`, `outline.md`, `angles.md`, `sources.md`,
  `draft.md`, `staff-edit.md`, `sharpen.md`, `seo.md`, or
  `seo-checklist.md`.** This pass only ever writes `humanize.md`.
- **Never use `Edit`, `Write`, or Bash on any piece file — not even for a
  small change.** Read the current file into context, make the change
  there, submit the whole revised text back through
  `save_humanize_bundle`. The tool re-validates citation integrity the
  same way every earlier pass's save does; a direct `Edit` bypasses that
  check entirely.

## Step 3 — persist (MCP)

Call `save_humanize_bundle` with `piece_path`, `revised_md` (the **full**
revised text — replaces the whole `humanize.md` file, never touches an
earlier-stage file), and `checks` (one line per real change, e.g.
`"[hedge-stack] dropped the weaker of two hedges on the compose claim"`,
`"[symmetry] varied the skip-list's third item so it doesn't read as a
template"`). If a check genuinely required no change, record that as a
single line too — an edit pass that logs nothing isn't a completed pass.
`mode` is optional and informational (`"delta"` or `"rewrite"`).

On tool error: fix exactly what it names. If it's a missing `draft.md`,
tell the operator to run `/ps-draft` first.

## Step 4 — summary (short, plain language for the author)

**"Sounds like you, not the machine."** — then 1-3 sentences describing
what actually changed, in plain language. Say what changed in the prose
itself ("cut a couple of hedge stacks and varied a too-uniform list") —
never recite this skill's own internal check names (no "stock transitions,"
"hedge stacks," "symmetry," "rhythm variance" as category labels — describe
the actual change, not the check that found it).

**Don't enumerate what needed no change.** Silence on a non-issue is the
right amount of information, per this project's own output discipline
(CLAUDE.md).

End with the saved path (`humanize.md`) — nothing else. Don't explain that
earlier files were left untouched, or why. No phase jargon, no permission
theater. Stop.

## Do not

- Deliberately break grammar, leave a thought incomplete, or "embrace
  messiness" as a goal — these are detector-evasion techniques, not voice
  authenticity, and this project has no interest in the former
- Skip an explanation the reader actually needs "for effect" — that's
  withholding information, not humanity
- Recite Step 2's internal check names back to the operator — describe
  what changed in the prose itself, in plain language
- Enumerate which checks needed no change — silence on a non-issue, not a
  line explaining the absence
- Use `Edit`, `Write`, or Bash on any piece file — always
  `save_humanize_bundle` with the full revised `revised_md`
- Overwrite or edit `draft.md`, `staff-edit.md`, `sharpen.md`, `seo.md`,
  or `seo-checklist.md` — the revision goes to `humanize.md`
- Split, rewrite, or relocate a resolved `[anchor](url)` link's anchor
  text or URL, or add a new link of your own — that's `/ps-seo`'s job
- Hardcode a banned-phrase list — read the operator's own `voice.md`/
  `style.md` via `read_author_defaults`
- Re-check stakes, brief-alignment, failure modes, honest scope, technical
  sanity, water/filler, readability, opener hook, key-insight elevation,
  killer-section weight, worked examples, or arc completion — that's
  `/ps-staff-edit`'s and `/ps-sharpen`'s job
- Build or reuse a Medium-specific curator-eligibility density gate — no
  such mechanism exists in this project
- Call a piece path a "slug"
