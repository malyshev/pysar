---
title: Export
slug: export
nav_order: 50
section: journey
---

# Export

Export copies the piece’s most refined Markdown into the **project root** as a
standalone file you can edit, commit, or publish elsewhere. Pysar never posts
for you.

## When it runs

`/ps` exports automatically after the final stage it runs (usually humanize).

You can also export via the MCP tool `export_piece_to_root` (skills invoke this;
you rarely call it by hand).

## What gets exported

Export picks the newest revision that exists for the piece, in this priority:

1. `humanize.md`
2. `seo.md` (if humanize has not run)
3. `sharpen.md`
4. `staff-edit.md`
5. `draft.md`

A draft is required. Earlier stage files are not deleted or rewritten by
export. Re-running export overwrites the previous root file for that piece.

Export currently copies that revision **as-is**. Research citation markers like
`[^shortname]` may still appear in the root file until mechanical resolve-at-
export ships. Do **not** treat `/ps --seo` as the required cleanup for those
markers — SEO packaging is a separate opt-in stage.

## Piece IDs and filenames

- Working files: `.pysar/pieces/<piece-id>/`
- Exported file: `<project-root>/<piece-id>.md` (same basename as the piece
  directory; the MCP tool response reports the path and word count)

`<piece-id>` is a machine-made **Latin** slug from the piece title. Non-Latin
titles (for example Ukrainian, Japanese, Arabic) are transliterated so you do
not get a bare `template-*` id. Authors do not invent the Latin slug by hand.

## Next

If export fails or picks the wrong revision, see
[Troubleshooting](./troubleshooting.md).
