---
name: ps-research
description: |
  Optional research pass: real, tiered sources for an existing piece, or a
  standalone research bundle when no piece exists yet. Never rewrites the
  operator's thesis, killer sections, or counterintuitive findings -- only
  adds source-backed key questions and angles. Adapted from a prior writing
  POC's research-sourcing discipline (source tiers, authority floor,
  citation hygiene, stop criterion), agent-agnostic: reuses expert_lens
  instead of a fixed topic-family table, reuses ps-factcheck for citation
  accuracy instead of a third copy of that check.
when_to_use: |
  Operator runs /ps-research, points at an existing piece for real sources,
  or wants to learn broadly about a topic before committing to a brief.
argument-hint: "[@path/to/piece or path, e.g. .pysar/pieces/<name>/] | [topic] [--competitors=url1,url2,...]"
allowed-tools: mcp__pysar__save_research_bundle Read WebSearch WebFetch Skill
---

# /ps-research — real sources, author's take stays theirs

**Author-facing outcome: "We found sources; your take stayed yours."** This
pass adds citations and source-backed context. It never rewrites the
operator's thesis, killer sections, or counterintuitive findings to match
what the web says — those are the operator's, decided at intake, and stay
that way. If the operator says "I need citations," this is that command.

This follows the same research-grounding direction `/ps-intake` already
applies (`dec-20260725-35fa2d24`'s amendment) — research grounds facts and
examples, never the take. That amendment's formal supersession is still
pending an operator `/h-decide` (tracked by `prob-20260725-ae764bb8`); this
skill anticipates it consistently with `/ps-intake`'s own current behavior,
not as a settled decision in its own right.

## Two modes

- **Piece-anchored**: operator gives a path to an existing piece (e.g.
  `.pysar/pieces/docker-for-developers-abc123/`) — never call this a
  "slug." **Referencing it via Claude Code's native `@` file-picker is
  normal flow, not a special case** — `@` commonly resolves to a specific
  file inside the piece (most likely `brief.md`) rather than the
  directory itself; pass whatever it resolves to straight through as
  `piece_path`. `save_research_bundle` accepts either form (a directory or
  a file inside it) and resolves to the containing piece directory either
  way — don't strip the filename yourself first. Adds real
  `sources.md`/`raw/`/optional `competitors.md` to it, sets
  `research_mode: full`, and appends (never replaces) `key_questions` and
  `angles.md`'s Misconceptions/Contrarian sections.
- **Standalone**: operator gives a topic instead of a path — no piece
  exists yet. Produces the same sources/raw/competitors under a fresh
  `.pysar/research/<topic>-<random-suffix>/` directory, plus a
  `research-summary.md`. This is for "I want to learn everything around a
  topic before I honestly say I know what I'm writing about" — background
  before a brief exists, not a brief itself. No thesis is committed here.

If the operator gives neither a path nor a topic, ask once which mode they
mean — this is the one genuine fork with no sensible default, unlike a
broad idea (which always has one).

## Step 1 — load context

- **Piece-anchored**: Read the piece's `brief.md` and `angles.md` once.
  Take `expert_lens` from there — don't ask, don't re-derive one. Restate
  the thesis to yourself in one sentence to confirm you're about to source
  *around* it, not replace it. Note existing `key_questions`/angles so you
  don't duplicate them.
- **Standalone**: determine `expert_lens` fresh from the topic, same
  reasoning `/ps-intake` Step 2 uses (the discipline/practitioner viewpoint
  whose rigor standard governs this topic's claims).

## Step 2 — source budget and authority floor

- **Budget**: 6-12 sources, 3-5 competitors if `--competitors=` was given.
  Without competitors: 8-14 sources, and at least one extra `angles`
  entry to compensate for the missing competitive contrast.
- **Tier each source** primary / secondary / community, judged by
  `expert_lens` — what would a working practitioner in that discipline
  call a primary source? (Official docs/specs/peer-reviewed for
  engineering; ag-extension/university sources for horticulture; SEP/named
  academic sources for methodology or craft topics — the pattern from
  `/ps-intake`'s own grounding, at research scale.) Set
  `topic_family_note` only if the mapping genuinely isn't obvious; leave
  it empty otherwise — most topics don't need it.
- **Authority floor**: >=60% of sources must be primary or secondary tier.
  `save_research_bundle` enforces this server-side and rejects a bundle
  that doesn't meet it — fix by fetching a better source, never by
  relabeling a community source's tier to make the number work.

## Step 3 — search and fetch (WebSearch/WebFetch)

For each source: a unique kebab-case `shortname`, the URL, tier, the date
actually fetched, 2-5 `key_claims` it actually supports, a `notes`
paraphrase, and a verbatim `raw_excerpt` (quote exactly — paraphrase goes
in notes, never in the excerpt). Never record a source you didn't actually
fetch.

If `--competitors=` gave URLs: for each, note its angle, strongest section,
and the gap this piece can fill.

## Step 4 — verify before persisting (invoke `ps-factcheck`)

Invoke `ps-factcheck` via the Skill tool and apply its **Job 2** to every
source you're about to submit: re-read what was actually fetched, confirm
each `notes`/`key_claims` entry matches it (not memory or a rounded
paraphrase), and don't let a source's citation imply coverage broader than
what it actually addresses. This is the same discipline `/ps-intake` uses
for its own lighter grounding — don't reimplement it here, invoke it.

## Step 5 — what to add (source-backed only)

- `key_questions_additions` — new questions this research answers that
  weren't already covered.
- `angles_misconceptions` / `angles_contrarian` — new, source-backed
  entries. Augment, don't replace: never delete or reword the operator's
  own existing angles.
- **Never touch thesis, killer_sections, or counterintuitive.** If a
  source contradicts the operator's own stated take, that's a note for the
  operator to weigh, not something this pass resolves by rewriting their
  thesis to match the web.

## Step 6 — stop criterion

Stop researching when: every `key_questions` entry (existing + new) is
answerable from `sources`, at least one contrarian/under-discussed angle
is documented (two if no competitors were supplied), the source count is
within budget, and the authority floor is met. Past this point, more
research is procrastination — persist and hand off.

## Step 7 — persist (MCP)

Call `save_research_bundle` with `piece_path` (piece-anchored) or `topic`
(standalone), `expert_lens`, `sources`, and whatever `competitors`/
`key_questions_additions`/`angles_*` apply. Never use `Edit`, `Write`, or
Bash for any of this, including on the piece's existing `brief.md`/
`angles.md`/`sources.md` — the tool does the mechanical file work so
nothing here risks touching content it shouldn't.

## Step 8 — summary (short)

- Piece-anchored: **"We found sources; your take stayed yours."** — then
  the count and the saved path. Suggest drafting next.
- Standalone: the saved research directory, and that running
  `/ps-intake --from-draft=<path>/research-summary.md` turns it into a
  piece when ready.

No phase jargon. No permission theater. Stop.

## Do not

- Use `Edit`, `Write`, or Bash on any piece file — always `save_research_bundle`
- Rewrite the operator's thesis, killer_sections, or counterintuitive to
  match what sources say — research adds citations, intake owns the take
- Invent a source, a URL, or a raw excerpt that wasn't actually fetched
- Reimplement `ps-factcheck`'s citation-accuracy check inline instead of
  invoking it
- Build or reuse a fixed per-topic authority table — use `expert_lens`
  instead, the way `/ps-intake` already does
- Call a piece path a "slug"
- Ask which mode (piece-anchored vs standalone) when a path or topic was
  already given — only ask when genuinely neither was
- Skip the authority floor or relabel a source's tier to satisfy it
- Delete or reword any of the operator's own existing angles when
  appending new ones
