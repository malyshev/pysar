---
name: ps
description: |
  Runs the full pipeline -- /ps-intake -> /ps-draft -> /ps-staff-edit ->
  /ps-sharpen -> /ps-humanize -- back to back on one piece, then exports
  the result to the project root. Autopilot by default: no stopping
  between stages. Pass --review to stop after each stage and wait for the
  operator's explicit go-ahead before continuing. Resumes correctly from
  an existing piece at whatever stage it already reached -- it does not
  restart a piece from intake just because /ps is what's invoked.
  Research is deliberately not in this chain (it's optional and
  ps-intake/ps-draft already invoke it where each needs it); this command
  exists only to remove the "stop and manually invoke the next skill"
  friction between the five stages above, nothing more.
when_to_use: |
  Operator types /ps with a new idea or an existing piece path, wanting
  the whole pipeline run without babysitting each stage themselves.
argument-hint: "[idea text | @path/to/piece] [--review]"
allowed-tools: Skill mcp__pysar__read_author_defaults mcp__pysar__save_intake_bundle mcp__pysar__save_draft_bundle mcp__pysar__save_staff_edit_bundle mcp__pysar__save_sharpen_bundle mcp__pysar__save_humanize_bundle mcp__pysar__export_piece_to_root Read WebSearch WebFetch
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

## Step 1 — figure out where this piece actually is

**A new idea (no existing piece path given):** start at `ps-intake`. There
is no earlier stage to detect.

**An existing piece path given** (a directory under `.pysar/pieces/`, or a
file inside one — `@`-reference or plain path, same as every other
piece-anchored command): determine the latest stage already reached by
attempting to `Read` these files in this order, stopping at the first one
that exists — the same priority order `/ps-sharpen` and `/ps-humanize`
already use for "which file to revise from":

1. `humanize.md` exists → every stage already ran. Nothing left to do but
   export (Step 4).
2. `sharpen.md` exists → next stage is `ps-humanize`.
3. `staff-edit.md` exists → next stage is `ps-sharpen`.
4. `draft.md` exists → next stage is `ps-staff-edit`.
5. `brief.md` exists (none of the above) → next stage is `ps-draft`.
6. None exist → treat like a new idea; start at `ps-intake` with whatever
   text the operator gave as the idea.

Do not guess or infer from the operator's phrasing which stage to start
at — check the files.

## Step 2 — run the chain, in order, starting from Step 1's answer

The fixed order is: `ps-intake` → `ps-draft` → `ps-staff-edit` →
`ps-sharpen` → `ps-humanize`. For each stage from the determined starting
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

## Step 4 — export to project root

Call `export_piece_to_root` with the piece's path. This is the one step
that writes outside `.pysar/pieces/<piece>/` — it copies whichever file is
currently the piece's most-refined revision to `<project root>/<piece
name>.md`. Call it whenever this run's chain ends, whether that's after a
full autopilot run, or wherever the operator stopped in `--review` mode —
`export_piece_to_root` works on a draft alone if that's as far as the
piece got; it does not require every stage to have run.

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
- Do not ask the operator to confirm continuing when `--review` was NOT
  given — autopilot means no stops, not "ask before each stage."
- Do not call `export_piece_to_root` more than once, or before the chain
  for this run has actually ended.
- Do not `Read` a stage's `SKILL.md` file directly — always go through the
  Skill tool.
- Do not invoke `/ps-research` as part of this chain — it is deliberately
  excluded; `ps-intake` and `ps-draft` already reach for grounding via
  `ps-factcheck`/research where each needs it, on their own terms.
