---
title: Run the pipeline
slug: pipeline
nav_order: 30
section: journey
---

# Run the pipeline

The writing workflow lives in the host agent as slash skills. The usual entry
is `/ps`: it runs the full chain and exports when finished.

## Full run: `/ps`

```text
/ps <your idea in a sentence or two>
```

Default stage order:

1. `/ps-intake` — clarify the brief
2. `/ps-draft` — first draft
3. `/ps-staff-edit` — structural edit
4. `/ps-sharpen` — tighten claims and language
5. `/ps-humanize` — final voice pass
6. Export to the configured export directory — project root by default (see [Export](./export.md))

Optional flags on `/ps`:

| Flag | Effect |
|------|--------|
| `--research` | Run `/ps-research` after intake and before draft. Persists a piece precondition so draft save fails closed until research is full. |
| `--seo` | Insert `/ps-seo` between sharpen and humanize (discoverability packaging for web/blog). Ordering is fixed: SEO never runs after humanize. Persists a piece precondition so humanize save fails closed until `seo.md` exists. |
| `--review` | Stop after each stage and wait for your go-ahead |

With no idea and no piece path, `/ps` prints a short plain-language
orientation instead of guessing.

Resume: if a piece already exists under `.pysar/pieces/`, `/ps` continues from
the stage already reached; it does not restart from intake just because you
invoked `/ps`.

## Optional first step: `/ps-onboard`

Teaches Pysar your voice and style (~2 minutes). Skipable — Pysar works with
a general-audience default if you skip it.

## Stage skills (manual)

You can run stages one at a time:

| Skill | Role |
|-------|------|
| `/ps-intake` | Brief / intake bundle |
| `/ps-draft` | Draft |
| `/ps-staff-edit` | Staff edit |
| `/ps-sharpen` | Sharpen |
| `/ps-seo` | Optional SEO packaging |
| `/ps-humanize` | Humanize |
| `/ps-research` | Optional research; required in-chain when `/ps --research` |
| `/ps-factcheck` | Optional fact-check |
| `/ps-voice` / `/ps-style` | Voice and style helpers used with onboarding |

Piece files live under `.pysar/pieces/<piece-id>/` (Latin slug from the title;
non-Latin titles are transliterated). Skills persist through Pysar MCP tools —
not by hand-editing those paths. Details: [MCP and skills](./mcp-and-skills.md).

Export copies the latest revision as-is today. Leftover research `[^shortname]`
markers are not fixed by adding `--seo`; see [Export](./export.md).
