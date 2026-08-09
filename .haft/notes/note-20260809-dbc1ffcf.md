---
id: note-20260809-dbc1ffcf
kind: Note
version: 1
status: active
title: Cursor rejects local plugin symlinks outside ~/.cursor/plugins/local
context: cursor-marketplace
mode: note
valid_until: 2026-11-07T11:27:49Z
created_at: 2026-08-09T11:27:49Z
updated_at: 2026-08-09T11:27:49Z
---

# Cursor rejects local plugin symlinks outside ~/.cursor/plugins/local

## Observations

- 2026-08-09 dogfood: symlink ~/.cursor/plugins/local/pysar → monorepo plugins/pysar was rejected by Cursor; Cursor Plugins.log: loadUserLocalPlugin pysar rejected: symlink target … is outside …/plugins/local; 0 plugins loaded.
- MCP tab can still show pysar from project .cursor/mcp.json or init MCP deeplink — that is not evidence the plugin package loaded under Customize → Plugins.
- Workaround that works: rsync -a --delete plugins/pysar/ → ~/.cursor/plugins/local/pysar/ (real directory). plugins/pysar/README.md updated accordingly.


## Anchors

- about `dec-20260809-cursor-marketplace-v1-dual-discovery-8b748a7a`
- relates_to `note-20260809-9d2e730b`

## Affected Files

- `plugins/pysar/README.md`
