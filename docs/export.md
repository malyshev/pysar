---
title: Export
slug: export
nav_order: 50
section: journey
---

# Export

Export copies the piece’s most refined Markdown into a standalone file you can
edit, commit, or publish elsewhere. Pysar never posts for you.

By default that file lands in the **project root**. You can choose a different
project-relative directory at init (`pysar init --export-dir published`) or by
editing `export_dir` in `.pysar/project`. The MCP tool
`export_piece_to_root` uses that setting; an optional `export_dir` argument
overrides it for one call only.

## When it runs

`/ps` exports automatically after the final stage it runs (usually humanize).

You can also export via the MCP tool `export_piece_to_root` (skills invoke this;
you rarely call it by hand). Skills omit `export_dir` unless you asked for a
one-off path — the tool result reports where the file landed.

## What gets exported

Export picks the newest revision that exists for the piece, in this priority:

1. `humanize.md`
2. `seo.md` (if humanize has not run)
3. `sharpen.md`
4. `staff-edit.md`
5. `draft.md`

A draft is required. Earlier stage files are not deleted or rewritten by
export. Re-running export overwrites the previous file for that piece at the
resolved destination.

Export currently copies that revision **as-is**. Research citation markers like
`[^shortname]` may still appear in the exported file until mechanical resolve-at-
export ships. Do **not** treat `/ps --seo` as the required cleanup for those
markers — SEO packaging is a separate opt-in stage.

## Piece IDs and filenames

- Working files: `.pysar/pieces/<piece-id>/`
- Exported file: `<export_dir>/<piece-id>.md` (or `<project-root>/<piece-id>.md`
  when `export_dir` is unset; same basename as the piece directory; the MCP tool
  response reports the path and word count)

`<piece-id>` is a machine-made **Latin** slug from the piece title. Non-Latin
titles (for example Ukrainian, Japanese, Arabic) are transliterated so you do
not get a bare `template-*` id. Authors do not invent the Latin slug by hand.

## Next

If export fails or picks the wrong revision, see
[Troubleshooting](./troubleshooting.md).
