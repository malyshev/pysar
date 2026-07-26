---
name: ps-sharpen
description: |
  Reader-experience pass on an existing draft (staff-edit.md if present,
  else draft.md): does the opening hook, are the piece's own thesis/killer-
  section/counterintuitive-finding ideas positioned prominently AND worded
  as sharply as they deserve, do killer sections get real structural
  weight, does every named framework have a worked example, and does the
  ending resolve what the opening promised. Deliberately not a punchline
  quota or rhythm floor -- optimizes for the piece's real, already-
  identified important ideas landing hard, never for manufacturing
  standalone one-liners to hit a count (a quota target gets quota-
  satisfying prose, not necessarily better writing). Writes the revision
  to a separate sharpen.md -- earlier stages (draft.md, staff-edit.md) are
  never touched, same reasoning as staff-edit keeping its own file.
  Deliberately narrower than a prior writing POC's SHARPEN phase: that
  POC's "consumability" audits (metaphor budget, paragraph density) and
  opener concreteness already live in this project's /ps-staff-edit
  (water/filler + readability + stakes checks); that POC's numeric rhythm
  floor is explicitly not reproduced here, replaced by the expression
  check above. Agent-agnostic: no author registry, no SEO/title/tag tuning
  (a prior writing POC's OPTIMIZE phase, unbuilt here), no AI-detection-
  evasion rhythm-breaking (that POC's HUMANIZE phase -- not built, and not
  a goal this project has any interest in).
when_to_use: |
  Operator runs /ps-sharpen, says "the intro doesn't land," or wants a
  final reader-experience pass before calling a piece complete.
argument-hint: "[@path/to/piece or path, e.g. .pysar/pieces/<name>/]"
allowed-tools: mcp__pysar__save_sharpen_bundle mcp__pysar__read_author_defaults Read
---

# /ps-sharpen — opening hits; piece completes

**Author-facing outcome: "Opening hits; piece completes."** If the operator
says "the intro doesn't land," this is that command. Staff-edit (if it ran)
already settled whether the argument is sound; this pass checks whether the
piece works as an experience, start to finish — does the opening actually
hook a reader, are the piece's own best findings given real weight instead
of sitting buried mid-paragraph, and does the ending deliver on what the
opening promised.

**Why this isn't redundant with `/ps-staff-edit`, even though both revise
prose:** staff-edit asks *"is this true and worth saying"* — stakes,
brief-alignment, failure modes, honest scope, technical sanity, water/
filler, sentence-level readability. This pass asks *"does it land"* — a
different, narrower question, at the whole-piece level: hook power, which
findings get elevated, which sections get real structural weight, whether
the arc resolves. If staff-edit already did its job, most of the heavy
lifting here is verification, not rewriting.

**The revision goes to `sharpen.md`, not `draft.md` or `staff-edit.md`.**
Each stage keeps its own file — the same convention `/ps-staff-edit`
established and this project uses throughout (`brief.md`, `outline.md`,
`angles.md`, `draft.md` each stay put once written).

This skill is **host-agnostic**. Prefer pysar MCP tools for anything
mechanical (profile defaults, citation checking, disk writes, run-log). Do
**not** use Bash/`ls`/`cat` or raw Write for piece files.

Prior writing POCs are inspiration only (CL1); do not copy their author
registry, SEO/title-tuning ceremony, or AI-detection-evasion rhythm work.
This project has no author registry — `voice.md`/`style.md` (via
`read_author_defaults`, the same tool `/ps-draft` and `/ps-staff-edit` use)
is the whole mechanism. No cross-piece corpus-variance scanning. No SEO
fields, tags, or title/subtitle tuning for a publish checklist that doesn't
exist here. No stripping rhythm patterns to evade AI-detection — that's not
a quality goal, and this project has no interest in building it regardless
of what a prior POC's pipeline did.

## Args

A path to an existing piece (e.g. `.pysar/pieces/docker-for-developers-
abc123/`) — never call this a "slug." Referencing it via Claude Code's
native `@` file-picker is normal flow, same as `/ps-draft` and
`/ps-staff-edit`.

If `draft.md` doesn't exist yet, tell the operator to run `/ps-draft`
first — do not invent one.

## Step 1 — load context

Read the piece's `brief.md`, `outline.md`, `angles.md`, and `sources.md`
once. `brief.md`'s `Counterintuitive findings to elevate` and `Killer
sections` are what checks 2 and 3 below check against; the `Thesis` and
`Promise to the reader` are what check 5 checks the ending against. Call
`read_author_defaults` for register/tone — don't invent a different
source.

**Which file to revise from:** if `staff-edit.md` exists, read that — it's
the latest substantively-checked revision. Otherwise read `draft.md`
directly (staff-edit is optional, same as research is optional for draft).

## Step 2 — the checks

Apply in order. Produce a fully revised version in your own context as you
go — this is a rewrite pass, not a comment pass. The revision is submitted
whole in Step 3, never edited into an existing file directly (see Hard
rules).

1. **Opener hook power.** Read the first 2-3 paragraphs as a real reader
   would, cold. Would they actually want to keep reading — not "is this
   technically concrete" (staff-edit's stakes check already covers that),
   but "does this create a pull"? A factually-concrete but flat opener
   still fails this check. Tighten in place if it doesn't hook; don't
   swap the whole approach the piece already committed to, sharpen it.
2. **Key-insight elevation and expression.** Cross-read `brief.md`'s
   `Thesis`, `Counterintuitive findings to elevate`, and each `Killer
   section`'s core edge claim against the piece — this is the fixed list
   of ideas that actually matter; nothing outside it is in scope for this
   check. For each idea on that list, ask two separate questions:
   - **Position.** Is it sitting as a section heading, a standalone line,
     or otherwise given real prominence — or is it buried as supporting
     detail inside a paragraph about something else? If buried, elevate
     it: promote it to its own heading, or land it as a standalone line
     right after the evidence that supports it.
   - **Expression.** Independent of position, is its wording the sharpest
     available phrasing for what it's actually claiming — or is it
     hedged, diluted, or buried in clause structure that blunts it?
     Tighten the wording if a sharper phrasing exists for the *same*
     claim.

   **This is not a quota.** There's no target count of punchy sentences to
   hit, and no rhythm floor to satisfy — a piece where every idea on the
   list already lands hard, positioned and worded well, passes this check
   with zero edits. The goal is that the piece's *actual* important ideas
   hit as hard as they deserve to; it is never license to manufacture a
   standalone one-liner for something that isn't already one of these
   real, identified ideas. A quota target gets quota-satisfying prose —
   short sentences inserted to hit a number, whether or not they carry
   real weight. This check has no number to hit, on purpose.
3. **Killer-section structural weight.** Cross-read `brief.md`'s `Killer
   sections` against the piece's actual section lengths and headings. A
   killer section should read as visibly more substantial than the
   piece's average section, with a heading that names the specific edge,
   not a generic label. (This is a *structural* check — space and heading
   prominence — distinct from staff-edit's brief-alignment check, which is
   about whether the section's *content* earns its billing.)
4. **Worked-example completeness.** Scan for any named framework, rubric,
   or pattern the piece introduces (a numbered list of decisions, a named
   heuristic, a scoring system). Each needs a concrete worked example
   within a paragraph or two, from a real domain the audience would
   recognize — not a generic placeholder. Add one if missing, drawn from
   `sources.md` or the piece's own domain, never invented from nothing.
5. **Arc completion.** Read straight through, start to end. Does the
   close resolve the *specific* stake or question the opening raised — or
   does the piece wander onto a related-but-different note by the end?
   This is different from staff-edit's close-memorability check (is the
   close's own point sharp) — this checks whether it's the *right* point,
   the one the opening actually promised. If the ending doesn't resolve
   the opening, tighten the close to answer it directly, or (if the
   middle drifted) tighten the throughline so the piece stays on the
   promise it made in paragraph 1.

## Hard rules

- **`[^shortname]` markers preserved.** Don't convert, remove, or relocate
  a marker away from the claim it anchors. Resolving markers to links is
  not this pass's job.
- **No new factual claims, stats, versions, or names** not already implied
  by the piece or `sources.md`. Elevating a finding means giving it more
  prominence, not inventing new support for it.
- **Preserve verbatim** code blocks, commands, error messages, and version
  pins shown in code.
- **No banned phrases** from the operator's own `style.md`/`voice.md`.
- **No new sections.** This pass elevates, expands, and reorders what's
  already there. A genuinely missing section is a draft-stage gap — name
  it in the summary, don't invent one here.
- **Never touch `brief.md`, `outline.md`, `angles.md`, `sources.md`,
  `draft.md`, or `staff-edit.md`.** This pass only ever writes
  `sharpen.md`.
- **Never use `Edit`, `Write`, or Bash on any piece file — not even for a
  small change.** Read the current file into context, make the change
  there, submit the whole revised text back through
  `save_sharpen_bundle`. The tool re-validates citation integrity the
  same way `/ps-draft` and `/ps-staff-edit` do; a direct `Edit` bypasses
  that check entirely.

## Step 3 — persist (MCP)

Call `save_sharpen_bundle` with `piece_path`, `revised_md` (the **full**
revised text — replaces the whole `sharpen.md` file, never touches
`draft.md` or `staff-edit.md`), and `checks` (one line per real change,
e.g. `"[opener] tightened the hook so paragraph 1 promises a specific
payoff"`, `"[elevate] promoted the contrarian finding on X to its own
heading"`). If a check genuinely required no change, record that as a
single line too — an edit pass that logs nothing isn't a completed pass.
`mode` is optional and informational (`"delta"` or `"rewrite"`).

On tool error: fix exactly what it names. If it's a missing `draft.md`,
tell the operator to run `/ps-draft` first.

## Step 4 — summary (short, plain language for the author)

**"Opening hits; piece completes."** — then 2-3 sentences describing what
actually changed, in plain language, for the author. Say what changed in
the prose itself ("tightened the opening so it hooks on the specific
zombie-process failure instead of a generic settings claim," "moved the
contrarian finding about X into its own heading") — never recite this
skill's own internal check names back to the operator (no "opener hook
power," "key-insight elevation and expression," "killer-section structural
weight," "arc completion" — Step 2's numbered checks are this pass's own
working process).

**Don't enumerate what needed no change.** If a check found nothing to
fix, say nothing about it — silence on a non-issue is the right amount of
information, per this project's own output discipline (CLAUDE.md).

End with the saved path (`sharpen.md`) — nothing else. Don't explain that
earlier files were left untouched, or why — that's this pass's own
internal storage choice, not a feature the author asked about. No phase
jargon, no permission theater. Stop.

## Do not

- Recite Step 2's internal check names (opener hook power, key-insight
  elevation and expression, killer-section structural weight, worked-
  example completeness, arc completion) back to the operator — describe
  what changed in the prose itself, in plain language
- Enumerate which checks needed no change — silence on a non-issue, not a
  line explaining the absence
- Manufacture a standalone one-line "punchline" to hit some notion of a
  rhythm quota — wording tightened under check 2 applies only to the
  piece's already-identified important ideas (thesis, killer sections,
  counterintuitive findings), never invented as decoration for a count
- Use `Edit`, `Write`, or Bash on any piece file — always
  `save_sharpen_bundle` with the full revised `revised_md`
- Overwrite or edit `draft.md` or `staff-edit.md` — the revision goes to
  `sharpen.md`
- Resolve `[^shortname]` markers into links, tune SEO fields/tags/titles,
  or strip rhythm for AI-detection avoidance — none of these exist in this
  project, and the last one isn't a goal this project has
- Re-check stakes, brief-alignment, failure modes, honest scope, technical
  sanity, water/filler, or sentence-level readability — that's
  `/ps-staff-edit`'s job; re-litigating it here is duplicated work
- Build or reuse an author registry, or scan a cross-piece corpus for
  opener/closing variance — neither mechanism exists in this project
- Invent a new section to fix a missing-section problem — that's a
  draft-stage gap; name it, don't paper over it here
- Call a piece path a "slug"
