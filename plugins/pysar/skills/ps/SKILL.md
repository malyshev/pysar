---
name: ps
description: |
  Runs the full pipeline -- /ps-intake -> [/ps-research if --research] ->
  /ps-draft -> /ps-staff-edit -> /ps-sharpen -> [/ps-seo if --seo] ->
  /ps-humanize -- back to back on one piece, then exports the result to
  the configured export directory (project root by default). Autopilot by
  default: no stopping between stages.
  Pass --review to stop after each stage and wait for the operator's
  explicit go-ahead before continuing. Pass --research to require full
  /ps-research after intake and before draft (hard-gated via
  require_piece_stages + save_draft_bundle; dec-20260809-701b59d3).
  Pass --seo to insert /ps-seo between /ps-sharpen and /ps-humanize
  (hard-gated via require_piece_stages + save_humanize_bundle;
  dec-20260804-e3234e50 ordering — never after /ps-humanize). Without
  those flags, research and SEO stay out of the chain. Resumes correctly
  from an existing piece at whatever stage it already reached. When run
  with genuinely nothing -- no idea text, no piece path -- shows a short,
  plain-language orientation instead of guessing at an idea.
when_to_use: |
  Operator types /ps with a new idea or an existing piece path, wanting
  the whole pipeline run without babysitting each stage themselves. Typing
  /ps with nothing at all shows plain-language orientation instead of
  guessing at an idea.
argument-hint: "[idea text | @path/to/piece] [--review] [--research] [--seo]"
allowed-tools: Skill mcp__pysar__read_author_defaults mcp__pysar__check_onboarding_status mcp__pysar__save_intake_bundle mcp__pysar__require_piece_stages mcp__pysar__save_draft_bundle mcp__pysar__save_staff_edit_bundle mcp__pysar__save_sharpen_bundle mcp__pysar__save_seo_bundle mcp__pysar__save_humanize_bundle mcp__pysar__export_piece_to_root Read WebSearch WebFetch
---

# /ps — full pipeline, autopilot by default

**Every message you send is governed by this project's "Pysar output discipline" section (CLAUDE.md, dec-20260719-7d675b61)** — apply its 3-question check before sending anything.

This command is an orchestrator, not a sixth conversation of its own. It
contains no hardcoded knowledge of what intake, draft, staff-edit, sharpen,
or humanize each actually check — that knowledge lives entirely in their
own `SKILL.md` files, which this command invokes via the Skill tool and
follows exactly as written, same pattern `/ps-onboard` already established.

**Always invoke a stage via the Skill tool, never via `Read` on its
`SKILL.md` file.** `Skill(ps-intake)`, `Skill(ps-draft)`, etc. are
pre-approved by `pysar init`'s shipped `settings.json`; a raw `Read` on a
skill file's absolute path is not.

**Do not discover skills via the filesystem.** Never `Bash`/`ls`/`find`/
`grep` on `~/.claude/skills`, `~/.cursor/skills`, or any skills directory,
and never ask the operator for permission to do so. The stage names are
already in this file's chain — invoke them with `Skill(...)`. If a Skill
call fails, say that; do not inventory the disk.

## Step 0 — nothing to work with: show orientation, don't guess

If the operator ran `/ps` with genuinely nothing to work with — no idea
text, no `@path/to/piece`, nothing but flags or literally nothing at
all — do not hand an empty string to `ps-intake` and do not guess at an
idea. Show orientation instead, then stop.

Any `--review`, `--research`, or `--seo` flag given on this invocation
still applies once the operator's next message supplies the idea or path
— this is one continuous exchange, not two separate invocations. Don't
ask them to repeat a flag they already gave.

1. Call `check_onboarding_status` once, silently — no need to narrate
   the check itself. **On tool error:** skip straight to the orientation
   message below without the onboarding paragraph — don't surface the
   raw error to the operator, and don't drop the whole message over one
   failed side-check.
2. Print exactly this shape, in plain language, no phase jargon (never
   say intake, staff-edit, sharpen, seo, or humanize here):

   > **Pysar turns an idea or a rough draft into a piece you'd actually
   > publish.** You bring the take; I handle turning it into something
   > finished.
   >
   > **To start, just tell me what it's about:**
   > - `/ps a habit that actually helped our team ship faster`
   > - `/ps --from-draft=@notes/half-finished-draft.md` — if you already
   >   have something written
   >
   > That's it — no setup required. I'll ask if I genuinely need more
   > from you; otherwise I'll get on with it and show you the result.

   If `check_onboarding_status` reports voice or style outstanding, add
   one more short paragraph: an optional one-time step, `/ps-onboard`
   (about 2 minutes), teaches Pysar to sound like the operator from the
   start — framed as optional, not a blocker; skipping it and just
   writing now with a general-audience default works fine too. If both
   are already done, don't mention onboarding at all — irrelevant noise
   in that branch.
3. Stop. Do not invoke `ps-intake` or any other stage this turn — wait
   for the operator's next message.

The two example lines above are illustrative, not a fixed script —
adapt the wording if it reads more naturally, but keep both forms (a
bare idea, and an `@path` to an existing draft) and keep it this short.

## Step 1 — figure out where this piece actually is

Only reachable once real input exists (this invocation had one from the
start, or Step 0's stop was followed by the operator supplying one).

**A new idea (no existing piece path given):** start at `ps-intake`. There
is no earlier stage to detect.

**An existing piece path given** (a directory under `.pysar/pieces/`, or a
file inside one — `@`-reference or plain path, same as every other
piece-anchored command): determine the latest stage already reached by
attempting to `Read` these files in this order, stopping at the first one
that exists — the same priority order `/ps-seo` and `/ps-humanize`
already use for "which file to revise from":

1. `humanize.md` exists → every stage already ran. Nothing left to do but
   export (Step 4).
2. `seo.md` exists → next stage is `ps-humanize` (the piece already went
   through `ps-seo`, whether or not `--seo` was passed *this* invocation
   — an existing `seo.md` always means humanize is next, never re-run
   `ps-seo` on a piece that already has one).
3. `sharpen.md` exists → next stage is `ps-seo` if `--seo` was given,
   else `ps-humanize`.
4. `staff-edit.md` exists → next stage is `ps-sharpen`.
5. `draft.md` exists → next stage is `ps-staff-edit`.
6. `brief.md` exists (none of the above) → if `--research` was given and
   brief frontmatter does **not** show `research_mode: full`, next stage
   is `ps-research`; otherwise next stage is `ps-draft`.
7. None exist:
   - **The operator gave actual idea text or `--from-draft=`** (not a
     piece path): treat like a new idea; start at `ps-intake` with the
     operator's own words, exactly as the no-existing-piece-path branch
     above does.
   - **The operator gave a piece path** (a directory under
     `.pysar/pieces/`, or an `@`-reference to one) **but it has no
     artifacts at all**: don't silently pass that path string to
     `ps-intake` as if it were idea text — it isn't one. Say so plainly
     ("that piece directory exists but has nothing in it yet — want to
     start intake there, or did you mean a different path?") and wait.

Do not guess or infer from the operator's phrasing which stage to start
at — check the files.

## Step 1b — persist opt-in stage preconditions (hard gates)

Once a piece path exists (after intake on a new idea, or when an existing
piece path was given), if this invocation has `--research` and/or `--seo`,
call `require_piece_stages` with that piece path and the matching stage
names (`research` and/or `seo`) **before** invoking draft or humanize.
Do this even on resume — the tool merges idempotently. This is what arms
the MCP fail-closed gates (dec-20260809-701b59d3). Without the flags, do
not call it and do not clear stages already on the piece.

## Step 2 — run the chain, in order, starting from Step 1's answer

The fixed order is: `ps-intake` → [`ps-research` if `--research`] →
`ps-draft` → `ps-staff-edit` → `ps-sharpen` → [`ps-seo` if `--seo`] →
`ps-humanize`. `ps-seo` is never inserted after `ps-humanize` — that
ordering is fixed regardless of when `--seo` was noticed
(dec-20260804-e3234e50). For each stage from the determined starting
point onward:

1. Invoke that stage via the Skill tool, passing it the piece path (once
   intake has produced one) exactly as that stage's own `argument-hint`
   expects. Conduct its entire conversation exactly as that skill defines
   it — same steps, same questions, same tools, same checks. Do not
   summarize, thin out, or shortcut it just because it's running from
   inside `/ps` — running it here must produce the exact same real result
   (a validated, saved file) as the operator typing that command directly.
2. Once that stage's own save tool has succeeded, apply the **review
   gate**, below, before touching the next stage.

**Review gate:**

- **`--review` not given (default, autopilot):** give one short,
  plain-language progress line (what just landed, e.g. "draft.md written,
  1,200 words — moving to staff-edit") and continue immediately to the
  next stage. Do not ask the operator anything.
- **`--review` given:** give that same one-line progress note, then
  **stop and wait for the operator's next message** before invoking the
  next stage. Do not invoke it in the same turn, and do not ask a
  yes/no question that invites an immediate auto-continue on your own —
  actually stop. If the operator's next message is anything other than an
  explicit go-ahead (a correction, a question, a request to edit the
  file), address that first; only proceed to the next stage once they
  clearly signal to continue.

If the operator gave `--review` and stops the chain partway (declines to
continue, or the conversation moves on to something else), treat wherever
the chain currently stands as final for this run — do not chase them for
a decision. Step 4 (export) still applies to whatever was reached; see
below.

## Step 3 — after humanize

Once `ps-humanize` completes (or Step 1 found `humanize.md` already
present), the chain is done. Move to Step 4.

## Step 4 — export

Call `export_piece_to_root` with the piece's path only — do **not** pass
`export_dir` unless the author explicitly asked for a one-off landing path
for this export. Omitting `export_dir` uses `.pysar/project`'s `export_dir`
when set, otherwise the project root. The tool result reports the resolved
destination path; cite that path in Step 5.

This is the one step that writes outside `.pysar/pieces/<piece>/` — it
copies whichever file is currently the piece's most-refined revision to
`<export_dir>/<piece name>.md`. Call it whenever this run's chain ends,
whether that's after a full autopilot run, or wherever the operator stopped
in `--review` mode — `export_piece_to_root` works on a draft alone if
that's as far as the piece got; it does not require every stage to have
run.

Do not call `export_piece_to_root` more than once per run of `/ps`, and do
not call it mid-chain — only once Step 2/3 has genuinely finished (either
the full chain completed, or the operator ended the review session).

## Step 5 — summary (short)

State, in 2-3 sentences: which stages ran this invocation (not which
stages exist in the abstract — the ones that actually ran, or that Step 1
found already done), and the final exported path from Step 4. No phase
jargon, no recap of each stage's own internal checks — those stages
already gave their own summaries as they ran. Stop.

## What NOT to do

- Do not hardcode what each stage checks or how it works — follow exactly
  what that stage's own `SKILL.md` defines, invoked live via the Skill
  tool. A future change to any one stage must work here with zero changes
  to this file.
- Do not restart a piece from `ps-intake` when an existing piece path was
  given — detect the real stage via Step 1's file checks first.
- Do not summarize or shortcut an individual stage's own steps — Description
  ≠ Work applies here too, same as `/ps-onboard`.
- Do not silently skip the review gate when `--review` was given — actually
  stop and wait for the operator's next message, not just a rhetorical
  pause before continuing anyway.
- Do not ask the operator to confirm continuing *between pipeline stages*
  when `--review` was NOT given — autopilot means no stops there, not
  "ask before each stage." This does not apply to Step 0's stop, which
  is unconditional (there's no stage to run yet at all when it fires).
- Do not call `export_piece_to_root` more than once, or before the chain
  for this run has actually ended.
- Do not `Read` a stage's `SKILL.md` file directly — always go through the
  Skill tool.
- Do not `Bash`/`ls`/`find`/`grep` skill install directories (including
  `~/.claude/skills` and `~/.cursor/skills`) to "see what's installed" —
  that is not part of this pipeline and must not trigger a permission
  prompt. Call `Skill(ps-…)` by the names in Step 2.
- Do not skip `ps-research` when `--research` was given (unless
  `research_mode: full` already), and do not invoke it when `--research`
  was not given.
- Do not stop the chain to ask the operator how to "handle" `--research`
  (skip vs broaden vs invent) — `ps-research` already owns the default
  non-fabrication path; invoke it and continue.
- Do not skip calling `require_piece_stages` when `--research` or `--seo`
  was given once a piece path exists — that call arms the hard gates.
- Do not invoke `ps-seo` after `ps-humanize`, ever, even if `--seo` is
  given late or the operator asks for it in that order mid-run — the
  ordering is fixed (dec-20260804-e3234e50); explain why if asked, don't
  silently reorder.
- Do not skip `ps-seo` when `--seo` was given, and do not run it when it
  was not — Step 1's `seo.md` check already tells you whether it ran on a
  resumed piece; don't re-derive that from the operator's phrasing.
- Do not treat a genuinely empty invocation as an idea and hand it to
  `ps-intake` — show Step 0's orientation and stop instead of guessing.
- Do not list this project's internal pass names (intake, staff-edit,
  sharpen, seo, humanize) in the Step 0 orientation message — an author
  running `/ps` for the first time has no reason to know these exist.
- Do not mention `/ps-onboard` in Step 0 when `check_onboarding_status`
  reports both voice and style already done.
