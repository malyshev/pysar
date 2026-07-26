---
name: ps-draft
description: |
  Write phase: draft draft.md from an existing piece's brief/outline/
  angles/sources, or from a rough brief given directly instead of a path
  (invokes ps-intake first to scaffold one, same as any other intake entry).
  draft.md is channel-agnostic writing, not a finished web/print artifact --
  citations go through [^shortname] markers matching real sources from
  ps-research and stay unresolved; turning those into inline links (or any
  other channel-specific packaging: SEO metadata, a formatted letter) is a
  later, optional, separate concern this skill doesn't own. Adapted from a
  prior writing POC's article-writer discipline (citation hygiene,
  consumability, worked-example rule, thesis-consistency), agent-agnostic:
  no author registry (reuses this project's own voice.md/style.md via
  read_author_defaults), no Medium-specific publish ceremony (SEO fields,
  tags, cover image, Boost optimization, cross-article corpus-variance
  scanning) -- this project targets more than one surface and doesn't own
  that ceremony.
when_to_use: |
  Operator runs /ps-draft, points at an existing piece to write its first
  full draft, or says "write it from the brief" / gives a rough brief
  directly with no piece yet.
argument-hint: "[@path/to/piece or path, e.g. .pysar/pieces/<name>/] | [rough brief text]"
allowed-tools: mcp__pysar__save_draft_bundle mcp__pysar__read_author_defaults mcp__pysar__save_intake_bundle Read Skill
---

# /ps-draft — first full draft exists

**Author-facing outcome: "First full draft exists."** This pass writes the
piece's first full `draft.md` from its brief/outline/angles/sources. If the
operator says "write it from the brief," this is that command.

`draft.md` is **channel-agnostic writing**, not a finished web/print
artifact — every other piece file is named for its content, not a channel
(`brief.md`, `outline.md`, `angles.md`), and this follows the same
convention rather than assuming the eventual output is a web article.
Citation markers stay as `[^shortname]`, never resolved to inline links —
turning them into links (or adding SEO metadata, or reformatting for a
letter) is later, optional, channel-specific packaging this skill doesn't
do.

This skill is **host-agnostic**. Prefer pysar MCP tools for anything
mechanical (profile defaults, citation checking, disk writes, run-log). Do
**not** use Bash/`ls`/`cat` or raw Write for piece files — that burns tokens
and triggers permission prompts the MCP tools exist to avoid.

Prior writing POCs are inspiration only (CL1); do not copy their author
registry or Medium-specific publish ceremony. This project has no author
registry — `voice.md`/`style.md` (via `read_author_defaults`, the same tool
`/ps-intake` uses) is the whole mechanism. This project also targets more
than one surface (blog, letter, and future kinds) — it doesn't own a
Medium-UI checklist (SEO title, tags, cover image, member-only flag,
distribution target) or a cross-piece corpus-variance scanner; neither is
built here.

## Two modes

- **Piece-anchored**: operator gives a path to an existing piece (e.g.
  `.pysar/pieces/docker-for-developers-abc123/`) — never call this a
  "slug." Referencing it via Claude Code's native `@` file-picker is normal
  flow, same as `/ps-research` — `@` commonly resolves to a specific file
  inside the piece rather than the directory itself; pass whatever it
  resolves to straight through as `piece_path`. `save_draft_bundle` accepts
  either form.
- **Standalone / free-form brief**: operator gives a rough brief directly
  ("write it from the brief") instead of a path — no piece exists yet.
  Invoke `ps-intake` via the Skill tool first, treating the operator's own
  words as the idea/draft input, to scaffold a real piece the same way any
  other intake entry would. Then draft against the resulting piece. Never
  invent a piece directory yourself, and never skip straight to `Write`ing
  a `draft.md` with no piece behind it — a draft with no brief has nothing
  for later passes (or the operator) to check it against.

If the operator gives neither a path nor enough to scaffold a brief from,
ask once which they mean — same genuine-fork exception `/ps-research` uses.

## Step 1 — load context

Read the piece's `brief.md`, `outline.md`, `angles.md`, and `sources.md`
(note the `Shortname` on each entry) once. The outline is the skeleton; the
angles are the spine; the brief is the contract with the reader. Call
`read_author_defaults` for register/tone — the same tool and the same
profile `/ps-intake` already reads. If the profile has structural
preferences (a preferred opener style, a closing habit), follow them; if
it's silent, use plain judgment — don't invent a fixed archetype taxonomy
the operator's own profile doesn't have.

## Step 2 — front-load the hook

**Heading shape depends on the piece's target surface — this project isn't
Medium-only, so this isn't one fixed rule.** For a blog-shaped piece: line 1
is the title (`# `), the line right under it is an italic subtitle (`*...*`)
that sets the stake, not a restatement of the title. For a piece with no
heading hierarchy at all (e.g. a letter), skip this entirely — don't impose
a web-page title/subtitle shape where none belongs. Either way, the first
paragraph front-loads the hook — no throat-clearing, no "in this piece we'll
cover..." bridge.

## Step 3 — draft the body

- Follow `outline.md`'s structure unless there's a concrete reason to
  deviate — mention the deviation when you hand off.
- Engage every contrarian/misconception angle in `angles.md` — at least one
  section should land an opinion the reader could disagree with.
- Land every `killer_section` from `brief.md` with a concrete worked
  example (a real stack, a real config, a real scenario — never a generic
  placeholder), and every `counterintuitive` finding without burying it —
  each earns at least a heading or a one-line punchline, not a passing
  mention.
- One idea per body paragraph. Split an overloaded paragraph as you write
  it — don't draft it flat and rely on a later pass to fix the rhythm.
- One governing metaphor. Extend it if you have one; don't stack a second.
- Close by landing the point once — don't restate the thesis two or three
  different ways in the last paragraphs.
- **Cite via `[^shortname]` markers** whenever a claim needs backing — the
  shortname must match a real entry in `sources.md`. Never a raw URL in
  prose; if nothing in `sources.md` backs a claim you want to make, either
  drop the claim's implied authority (state it as your own reasoning) or
  run `/ps-research` first.
- For a blog-shaped piece: `##` / `###` only (no deeper), and `###` only
  when a `##` section has 2+ subpoints worth separating. For a piece with
  no heading hierarchy, use paragraph breaks and transitions instead — this
  is judgment, not a tool-enforced rule (see Step 6).

## Step 4 — length

Let the topic's actual scope set the length — real examples and trade-offs
earn space, recap paragraphs don't. There's no fixed target-length field to
check against here; if the brief's `topic_scope` clearly implies a size far
off from what you're about to write, say so once before drafting instead of
silently over- or under-shooting.

## Step 5 — self-check before persisting

- Every `[^shortname]` marker matches a real source in `sources.md` — the
  tool enforces this, but checking it yourself first avoids a wasted
  round trip.
- Zero raw URLs in the prose — the tool enforces this too.
- Heading shape matches the piece's target surface (Step 2/3) — the tool
  does **not** check this; it's your judgment call to get right.
- Every `killer_section` and `counterintuitive` finding from `brief.md` is
  genuinely landed in the body, not just mentioned in passing.
- The close doesn't reproduce what `brief.md`'s `Thesis` critiques.

## Step 6 — persist (MCP)

Call `save_draft_bundle` with `piece_path` and `draft_md` (plus an optional
one-line `notes` for the changelog, e.g. `"first full draft"` or `"revision
after a staff-edit pass"`). Never use `Edit`, `Write`, or Bash on `draft.md`
directly, including for a small fix after the tool flags a problem — the
tool validates citation integrity (every `[^shortname]` resolves to a real
source; no raw URL in prose), so nothing here risks shipping a broken or
invented citation. It does **not** resolve citation markers to links or
check heading/title shape — link-resolution and heading shape both depend
on the piece's eventual target channel, so they stay your judgment call
(Step 2/3) or a later pass's job, not something this tool can generically
enforce. A redraft replaces `draft.md` wholesale — expected behavior, since
a draft isn't an accumulating list the way `sources.md` is.

On tool error: fix exactly what it names (an unmatched citation, a raw
URL). Never relabel or drop a citation just to make validation pass —
either fix the shortname or fetch the source via
`/ps-research`.

## Step 7 — summary (short)

**"First full draft exists."** — word count, the saved path, and anything
you deviated from the outline on. No phase jargon, no permission theater.
Stop.

## Do not

- Use `Edit`, `Write`, or Bash on `draft.md` — always `save_draft_bundle`
- Invent a fixed opener/closing-archetype taxonomy the operator's profile
  doesn't have — use `voice.md`/`style.md` via `read_author_defaults`, the
  same source `/ps-intake` already reads
- Build or reuse an author registry — this project has none
- Write a Medium-specific publish checklist (SEO title, tags, cover image,
  member-only flag, distribution target) — this project targets more than
  one surface and doesn't own that ceremony
- Scan a cross-piece corpus for opener/closing variance — no such mechanism
  exists here; draft each piece on its own merits
- Put a raw URL in prose — cite via `[^shortname]`
- Skip an angle from `angles.md` or a `killer_section` from `brief.md` just
  because covering it is more work
- Invent a piece directory yourself in standalone mode — invoke `ps-intake`
  via the Skill tool and let it allocate one the normal way
- Call a piece path a "slug"
