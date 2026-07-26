---
name: ps-intake
description: |
  Scaffold brief, outline, angles, and a sources stub from an author's idea or draft.
  Profile-driven audience/register defaults; asks only on genuinely degenerate input.
  Persists via pysar MCP to a provisional, uniquely-named .pysar/pieces/ directory.
  Host-agnostic (Claude Code, Cursor, Codex, Gemini, and any MCP-capable agent).
  Use when the operator runs /ps-intake, brings an idea, or supplies a draft to shape.
when_to_use: |
  Operator types /ps-intake, pastes an idea, or points at a draft to turn into piece scaffolding.
argument-hint: "[idea text] [--from-draft=path] [--name=kebab-case]"
allowed-tools: mcp__pysar__read_author_defaults mcp__pysar__save_intake_bundle Read WebSearch WebFetch Skill
---

# /ps-intake — idea or draft → brief scaffold

You are running Pysar **intake**: turn the author's take into durable scaffolding.
**No fabricated sources. No author registry. No word-count cap.**

This skill is **host-agnostic**. Prefer pysar MCP tools for anything mechanical
(profile defaults, directory naming, disk writes, run-log). Do **not** use
Bash/`ls`/`mkdir` or raw Write for piece files — that burns tokens and
triggers permission prompts the MCP tools exist to avoid. `settings` from
`pysar init` pre-allow `mcp__pysar__*`.

Governing decision: `dec-20260725-35fa2d24` — amended: intake may now research
to ground its own technical claims (see Step 4, delegated to `ps-factcheck`).
The operator's idea/thesis stays the framing authority either way; research
only verifies facts and examples, it never replaces the operator's take.
Prior writing POCs are inspiration only (CL1); do not copy their
author-registry or `--length` behavior.

**Output discipline:** before every message, apply this project's shared Pysar output-discipline check (CLAUDE.md / host contract). Prefer silent progress + one short restatement over questionnaires.

**Directory naming is never the author's problem.** Every piece gets its own
directory with a random suffix appended automatically
(`.pysar/pieces/<readable-prefix>-<random-suffix>/`) — this is storage
mechanics, not a decision. Never ask the author to name, rename, or resolve a
collision for a piece directory. Running intake twice on the same idea is
fine; it just produces two pieces.

## Goal: end-to-end without interruption

Complete the whole flow in one pass when the idea is usable:

1. Load defaults via MCP
2. Assemble brief/outline/angles in your head, grounded via `ps-factcheck`
3. Persist via `save_intake_bundle`
4. Print a short summary

**Ask the author only if input is degenerate** (empty, gibberish, a single stray word, or explicitly self-contradictory instructions). Broad or vague ideas are **not** grounds to ask — pick the best default from the profile, state it in the restatement line, and proceed.

## Args

- Idea text (positional / freeform) — required unless `--from-draft=` supplies prose
- `--from-draft=<path>` — extract thesis/outline/angles **only from the author's words** (use Read on that path once)
- `--name=` — optional override for the piece directory's readable prefix. A unique suffix is always appended automatically regardless — this flag never needs mentioning to the author unless they ask.

## Step 1 — mechanical prep (MCP, not chat)

Call `read_author_defaults` — returns audience/register/tone and whether a profile exists (`source: profile` vs `generalist-fallback`).

Do **not** re-read `.pysar/voice.md` / `.pysar/style.md` with Read unless the MCP tool failed.

## Step 2 — resolve defaults and frame the expert lens (never ask "who is this for?")

From `read_author_defaults`:

- If `source: profile` → use its `audience_hint`, `register`, `tone` as the piece defaults. Guessing a generic-public audience when a profile exists is a misread.
- If `source: generalist-fallback` → use general-audience explainer register as the piece default. **Do not mention in the restatement that no profile exists or that this is a "fallback"** — that's internal reasoning, not actionable for the author, and narrating it is exactly the absence-explaining noise this project's output discipline forbids (CLAUDE.md). State the resulting audience/register as a plain fact, same as you would with a real profile.

Topic-scope within a broad domain: pick the single most likely angle from the idea + profile. State it as an assumption in the restatement. Do not ask.

**Expert lens (`expert_lens`) — distinct from audience.** Audience is who you're writing *for*; expert lens is whose rigor standard governs what you're about to *claim*. Determine the single discipline/practitioner viewpoint the topic's core claims actually belong to — "experimental science / methodology" for a science piece, "horticulture" for a gardening piece, "distributed-systems security" for an infra piece — and hold that lens while assembling Step 5's content: what would a working practitioner in that discipline scrutinize, what precision would they expect, what's the well-known pitfall they'd catch immediately. This is what makes a piece survive an expert reader's review instead of just reading fluently. Never ask the author to confirm the lens — pick the obvious one and proceed, same as audience/register. Hand this to `ps-factcheck` in Step 4 as the rigor standard it judges claims against.

## Step 3 — draft path (only when `--from-draft=`)

1. Read the draft file once with Read.
2. Extract thesis, outline, angles from **author text only**.
3. Set `entry_mode=draft`, `pov_source=author-draft`.
4. Do not "improve" the author's own stated take with outside knowledge — that's still forbidden (Step 4 governs *facts and examples*, never the thesis/angle itself).

## Step 4 — ground technical claims (invoke `ps-factcheck`)

Invoke `ps-factcheck` via the Skill tool before writing killer_sections/counterintuitive/examples in Step 5. Pass it `expert_lens` (Step 2) as the rigor standard to judge claims against — it decides when a claim needs grounding (genuine accuracy uncertainty, or naming a specific person/framework/standard) and fetches real sources when it does. Follow its **Job 1** now; its **Job 2** (verification) applies later, once Step 5's content actually exists to check against what got fetched — don't skip back to it early.

Whatever `ps-factcheck` fetches, carry its `{url, note}` output into `sources` when you persist in Step 6. Empty is the common, legitimate outcome — most claims don't need grounding.

## Step 5 — assemble content (in memory)

Write as if `expert_lens` (Step 2) will review it before it ships — precision and terminology a practitioner in that discipline would actually use, not generic tutorial language. Produce:

- **Restatement** — one concise line covering audience + register + chosen topic-scope (always shown in the tool result / summary; never a confirm/deny question)
- **Thesis**, **promise**
- **killer_sections** (1-2) — each needs `title` + `edge` (why this beats existing/competing coverage, not just "this is important") + `example` (a concrete worked example: a real stack, a real config snippet, a real scenario — never a vague placeholder)
- **counterintuitive** (1-3) — each needs `claim` + `contradiction` (name explicitly what makes it counterintuitive, don't just assert "surprisingly...")
- **key_questions** (≥1) — required, not optional. Questions the piece must actually answer; without at least one, drafting has no scope boundary and tends to drift.
- **non_goals** (≥1) — required, not optional. What this piece explicitly will NOT cover; without at least one, drafting tends to scope-inflate.
- **outline_md** — a section-and-subsection plan tied to the idea, written as Markdown headings (`##` for a section, `###` for a subsection). Think in editorial terms — sections, subsections — never "H2/H3"; that's web-page vocabulary, meaningless to an author thinking about a written piece, not a webpage.
- **angles_md** — must include Misconceptions, **Contrarian / under-discussed** (≥1 real contrarian), Trade-offs, Edge cases

Internal artifact voice: matter-of-fact scaffolding log, not the article's final prose voice.

Omit: Author / Author-name fields, target length / `--length`. Real URLs in `sources` are fine when Step 4 actually fetched them — never invent one you didn't fetch.

**Before moving to Step 6, apply `ps-factcheck`'s Job 2** (already in context from Step 4's invocation) to what you just wrote — cross-check killer_sections/counterintuitive/examples against the sources actually fetched. This catches drift between what was grounded and what got written; it's a separate check from Job 1, not the same thing repeated.

## Step 6 — persist (MCP)

Call `save_intake_bundle` with the structured fields (`name` only if `--name=` was explicitly given; otherwise omit it and the tool derives one from the idea, always with a unique suffix; `sources` only if `ps-factcheck` actually fetched something). Mention in one short sentence that you are saving.

On tool error: fix exactly what it names (missing field, degenerate idea). Directory naming and collisions are never a tool error the author needs to resolve — the tool handles that silently. Retry the tool. Never Write or Edit files yourself to work around it, including a quick fix to an already-written brief.md/outline.md/angles.md — re-run the tool with the corrected fields instead.

On success the tool writes under a **provisional**, uniquely-named piece directory:

`.pysar/pieces/<prefix>-<random-suffix>/` → `STORAGE.md`, `brief.md`, `outline.md`, `angles.md`, `sources.md` (stub, or real sources when Step 4 grounded the piece), `intake-changelog.md`, `run-log.jsonl`

Treat that path as temporary topology pending a future piece-layout decision.

## Step 7 — summary (short)

Print:

- Restatement line
- The saved path (call it "the piece" or "the directory" — never "slug")
- Thesis one-liner
- Suggested next, as one decisive recommendation, not a flat "or":
  - If `read_author_defaults` returned `source: generalist-fallback`: recommend running `/ps-voice` next (one-time, ~2 minutes) so this and future pieces use the author's own voice automatically — mention drafting now with today's generalist tone as the fallback option, not the primary suggestion.
  - If `source: profile`: just suggest drafting the piece. Do not mention onboarding at all — it's irrelevant noise in this branch.

No phase jargon for the author. No permission theater. Stop.

## Do not

- Ask clarifying questions for merely broad ideas
- Ask the author to name, rename, or resolve a collision for a piece directory
- Say "slug" — call it the piece, the directory, or the piece name
- Use Bash to create directories or write Markdown, or `Edit`/`Write` on a piece file directly
- Copy prior-POC author-registry or length-cap behavior
- Expand the restatement into a confirmation gate
- Mention that no voice/style profile exists as a caveat on the restatement — state the resulting default as a plain fact instead
- Reimplement grounding/citation-accuracy rules inline instead of invoking `ps-factcheck` — that duplication is exactly what drifts out of sync across skills over time
