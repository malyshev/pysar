---
name: ps-seo
description: |
  Discoverability packaging pass on an existing draft (sharpen.md if
  present, else staff-edit.md, else draft.md): resolves every [^shortname]
  citation marker into a real inline link, tunes the title (click-through)
  and subtitle (read-completion) as two distinct jobs, checks scannability
  and outbound-link authority, and produces a platform-neutral
  discoverability checklist (tags, meta title, meta description, URL
  slug). Runs before /ps-humanize, never after -- humanize deliberately
  roughens the transitions this pass smooths, so reversing the order lets
  the later pass undo the earlier one (dec-20260804-e3234e50). Writes the
  revision to a separate seo.md plus seo-checklist.md -- earlier stages
  (draft.md, staff-edit.md, sharpen.md) are never touched, same reasoning
  as every other pass keeping its own file. Adapted from a prior writing
  POC's SEO-optimizer discipline (citation resolution, title/subtitle as
  two distinct jobs, outbound-link authority floor) but not mirrored:
  that POC's checklist is Medium-specific (a rigid 5-slot tag taxonomy,
  Medium's own SEO/curator model); this pass targets whatever platform the
  author is posting to, with no forced slot count or platform-specific
  field names. Opt-in only -- most pieces never need this; it exists for
  pieces headed to a blog/web surface where search and link-preview
  discovery matter.
when_to_use: |
  Operator runs /ps-seo, says "make it discoverable" or "I need this to
  rank," or wants citations resolved into real links and a discoverability
  checklist filled before humanizing a piece headed to the web.
argument-hint: "[@path/to/piece or path, e.g. .pysar/pieces/<name>/]"
allowed-tools: mcp__pysar__save_seo_bundle mcp__pysar__read_author_defaults Read
---

# /ps-seo — links resolved, discoverable, ready for a final voice pass

**Author-facing outcome: "Links filled in, title earns the click, and it's
ready to be found."** If the operator says "make it discoverable" or "I
need this to rank," this is that command. Every earlier pass wrote and
edited prose; this pass packages that prose for search and link-preview
surfaces -- resolving citations into real links, sharpening the title and
subtitle for two different jobs, and filling in the metadata a platform
actually reads.

**This is opt-in, not part of the default chain.** Most pieces (a paper
letter, an internal note, a piece with no citations) never need this pass.
Run it when the piece is headed to a blog/web surface and discoverability
matters -- `/ps` supports it via `--seo`.

**Why this must run before `/ps-humanize`, never after.** Humanize
deliberately de-sterilizes prose -- it breaks overly clean transitions and
restores sentence-rhythm variance. This pass does the opposite: it
resolves citations and tightens packaging, which naturally smooths
prose. Running it after humanize would silently undo humanize's work,
handing the author something that reads sterile again despite humanize
having already run. `/ps-humanize` reads `seo.md` first when it exists,
precisely so this ordering is structural, not a reminder someone has to
remember.

**The revision goes to `seo.md` and `seo-checklist.md`, never `draft.md`,
`staff-edit.md`, or `sharpen.md`.** `seo.md` is reader-facing article
prose -- the same kind of file `/ps-humanize` reads and export copies to
the project root. `seo-checklist.md` is operator-facing packaging
metadata (tags, meta title/description, URL slug) -- the author copies
these into whatever platform they're posting to; it never gets exported
to the project root, the same way `sources.md` and the changelogs stay in
the piece directory.

This skill is **host-agnostic**. Prefer pysar MCP tools for anything
mechanical (profile defaults, citation checking, disk writes, run-log). Do
**not** use Bash/`ls`/`cat` or raw Write for piece files.

Prior writing POCs are inspiration only (CL1); do not copy their Medium-
specific 5-slot tag taxonomy, curator-eligibility framing, or FAQ-schema
tricks. This project has no author registry -- `voice.md`/`style.md` (via
`read_author_defaults`) is the whole mechanism, and title/subtitle tuning
must preserve whatever voice forms it already encodes, not flatten them
into generic SEO-template phrasing.

## Args

A path to an existing piece (e.g. `.pysar/pieces/docker-for-developers-
abc123/`) -- never call this a "slug." Referencing it via Claude Code's
native `@` file-picker is normal flow, same as every other piece-anchored
skill.

**If invoked with no path at all**, don't guess which piece is meant.
Show this, then stop and wait:

> **Resolves citations into real links, makes the title earn the
> click, fills a discoverability checklist** -- opt-in, for a piece
> headed to a blog/web surface.
>
> **To start:**
> - `/ps-seo @.pysar/pieces/docker-for-developers-abc123/`
>
> If this is part of a longer run, `/ps --seo` handles the whole
> pipeline in one go instead of stage by stage.

If `draft.md` doesn't exist yet for the given piece, tell the operator to
run `/ps-draft` first -- do not invent one.

## Step 1 -- load context

Read the piece's `brief.md` (for keywords and intent) and `sources.md`
(the only source of truth for resolving citations -- never invent a URL).
Call `read_author_defaults` for voice/register -- title and subtitle
tuning must stay recognizable as the author's own voice forms, not
generic search-optimized phrasing.

**Which file to revise from:** if `sharpen.md` exists, read that -- it's
the latest edited revision. Otherwise `staff-edit.md`, otherwise
`draft.md` directly (staff-edit and sharpen are both optional, same
precedence every later pass already uses).

## Step 2 -- citation resolution (always first)

Before any packaging work:

1. Scan the revision for every `[^shortname]` marker.
2. For each marker, look up `shortname` in `sources.md`.
3. Pick 1-3 anchor words from the surrounding sentence -- never "click
   here," never the URL itself, never the source's title verbatim.
4. Replace `<text>[^shortname]` with the anchor words turned into
   `[anchor](url)`, using the URL `sources.md` actually recorded. The
   `[^shortname]` marker is removed.
5. If a marker's shortname isn't in `sources.md`, stop and surface it as a
   blocker -- do **not** invent a link. `save_seo_bundle` rejects both an
   unresolved marker and a link whose URL doesn't match a recorded
   source, so guessing here just costs a round-trip.

After resolution, the revision must contain **zero** `[^shortname]`
markers -- that's this pass's one non-negotiable floor.

## Step 3 -- the packaging checks

1. **Title -- earns the click.** The article's own `# Title` line is the
   feed/reader-facing title; leave it as-is unless it's genuinely weak for
   a keyword match. Concrete, specific, primary keyword present if one
   exists naturally -- no vague "An Introduction to X" phrasing. Preserve
   whatever voice form `voice.md` encodes (a two-part dash title, a named-
   concept title, a declarative claim) -- don't flatten it into generic
   SEO phrasing.
2. **Subtitle -- converts the click to a read.** The `*subtitle*` line
   under the title has a different job: making someone who already saw
   the title commit to reading. Distinct angle from the title, not a
   restatement of it -- if you could delete either without losing
   information, the piece is failing one of the two jobs.
3. **Scannability.** Short paragraphs, sane heading hierarchy (H2/H3
   only), no wall-of-text sections. Don't restructure content that's
   already scannable -- this is a floor check, not a rewrite mandate.
4. **Outbound-link authority.** Every resolved link's underlying source
   should already meet whatever authority floor `sources.md` recorded
   (primary/secondary preferred). This pass doesn't re-vet sources --
   that's research's job -- but don't add a new outbound link beyond
   what citation resolution already produced.
5. **Inline-link density.** After resolution, no paragraph should carry
   more than 2 inline links -- three blue-underlined links in one
   paragraph reads as a research paper, not a piece. If resolution
   produced that, cluster two citations onto one anchor span when they
   support the same claim, or flag it in the summary for a follow-up
   `/ps-sharpen` paragraph-split rather than silently shipping it.

## Step 4 -- the checklist fields

Fill in, platform-neutral (no fixed slot count, no Medium-specific
taxonomy):

- **`tags`** -- >=1 discovery keyword/topic tag, whatever actually helps
  this piece get found on the platform it's headed to.
- **`meta_title`** -- the search/share title, a **third** job distinct
  from both the article title (feed/reader CTR) and subtitle (read-
  completion): what shows in a search result or link preview. Roughly 60
  characters is a reasonable target, not a hard cutoff.
- **`meta_description`** -- the search/share summary, concrete, written
  for someone who hasn't opened the piece yet. Roughly 155 characters is
  a reasonable target, not a hard cutoff.
- **`url_slug`** -- kebab-case, keyword-bearing, short.

## Hard rules

- **Zero `[^shortname]` markers may remain.** This is the one thing every
  earlier pass explicitly deferred to this pass -- `save_seo_bundle`
  rejects a submission that still has one.
- **Never invent a link.** Every `[anchor](url)` must trace to a URL
  `sources.md` actually recorded; `save_seo_bundle` rejects one that
  doesn't.
- **No new factual claims, stats, or attributions** beyond what
  `sources.md` already supports.
- **Preserve verbatim** code blocks, commands, error messages, and
  version pins shown in code.
- **No banned phrases** from the operator's own `style.md`/`voice.md`.
- **Never touch `brief.md`, `outline.md`, `angles.md`, `sources.md`,
  `draft.md`, `staff-edit.md`, or `sharpen.md`.** This pass only ever
  writes `seo.md` and `seo-checklist.md`.
- **Never use `Edit`, `Write`, or Bash on any piece file.** Read the
  current file into context, make the change there, submit the whole
  revised text back through `save_seo_bundle`.

## Step 5 -- persist (MCP)

Call `save_seo_bundle` with `piece_path`, `revised_md` (the **full**
revised text, zero `[^shortname]` markers remaining), `checks` (one line
per real change, e.g. `"[citation] resolved [^retry-budget] to an inline
link on 'the retry budget'"`, `"[title] tightened the title for a
concrete keyword match"`), `tags`, `meta_title`, `meta_description`, and
`url_slug`. `mode` is optional and informational (`"delta"` or
`"rewrite"`).

On tool error: fix exactly what it names. A rejected link almost always
means a shortname wasn't actually in `sources.md` -- re-check there, never
guess a URL.

## Step 6 -- summary (short, plain language for the author)

**"Links filled in, title earns the click, discoverability checklist
ready."** -- then 2-3 sentences on what actually changed, in plain
language: how many citations got resolved, whether the title/subtitle
changed and why. Never recite this skill's own internal step names back to
the operator.

**Don't enumerate what needed no change.** If a check found nothing to
fix, say nothing about it.

End with the saved paths (`seo.md`, `seo-checklist.md`) and the next step:
`/ps-humanize`. No phase jargon, no permission theater. Stop.

## Do not

- Leave any `[^shortname]` marker unresolved, or invent a link that
  doesn't trace to `sources.md`
- Run this pass after `/ps-humanize` -- always before, per the ordering
  this skill exists to protect
- Use `Edit`, `Write`, or Bash on any piece file -- always
  `save_seo_bundle` with the full revised `revised_md`
- Overwrite or edit `draft.md`, `staff-edit.md`, or `sharpen.md` -- the
  revision goes to `seo.md`
- Force a fixed tag-slot count or copy a platform-specific taxonomy --
  this pass is platform-neutral by design
- Flatten the title/subtitle into generic SEO-template phrasing --
  preserve the voice forms `voice.md` encodes
- Add deliberate grammar-breaking, incomplete thoughts, or "AI-detection
  evasion" techniques of any kind -- that is not this pass's job, and it
  is not this project's job anywhere (see `/ps-humanize`'s own
  discipline)
- Call a piece path a "slug"
