---
title: MCP and skills
slug: mcp-and-skills
nav_order: 40
section: journey
---

# MCP and skills

Pysar splits into three surfaces:

1. **CLI** — `pysar`, `pysar init`, `pysar serve`, `pysar --version`
2. **MCP server** — `pysar serve` (stdio), started by the host from project MCP config
3. **Skills** — `ps-*` instructions the host agent follows (`/ps` or `$ps`, …)

## MCP config

After init, the host launches:

```text
pysar serve
```

with the project root set so piece I/O lands in that project’s `.pysar/`.

| Host | Config file | Command / project root |
|------|-------------|------------------------|
| Claude Code | `.mcp.json` at project root | `pysar` / `PYSAR_PROJECT_ROOT=${PWD:-.}` |
| Cursor | `.cursor/mcp.json` | `${userHome}/.local/bin/pysar` / `PYSAR_PROJECT_ROOT=${workspaceFolder}` |
| Codex CLI / App | `.codex/config.toml` | `pysar` / `PYSAR_PROJECT_ROOT = "."` (`default_tools_approval_mode = "approve"`) |

Cursor uses `${userHome}/.local/bin/pysar` so Dock-launched Cursor can spawn the
server without inheriting your shell `PATH`. `pysar init --cursor` also prints a
one-step Cursor install link to enable the server in Customize → MCPs.

Claude pre-approves `mcp__pysar__*` via `.claude/settings.json`. Codex uses
`default_tools_approval_mode = "approve"` on the MCP server block for the same
friction goal (`auto` still prompts when tools lack risk annotations); project
`.codex/config.toml` applies for **trusted** projects.

You do not need to run `pysar serve` in a separate terminal for normal use —
the host agent starts it as an MCP server.

## Skills install location

Editorial skill bodies share one host-neutral corpus (`cmd/pysar/assets/skills`).
How each host receives them differs:

| Host | Skills carrier |
|------|----------------|
| Claude Code | `pysar init` → `~/.claude/skills/ps-*` |
| Cursor | **Pysar Cursor plugin** (`plugins/pysar`) via Marketplace or a local copy under `~/.cursor/plugins/local/pysar` — not `~/.cursor/skills`. Site **Install in Cursor** is MCP-only (see [Init](./init.md)). |
| Codex CLI / App | `pysar init --codex` → `~/.agents/skills/ps-*` (Codex packaging: `$ps-*` + `agents/openai.yaml`) |

Claude installs corpus bytes as-is. Cursor loads the same bytes from the
plugin package (kept in sync via `scripts/sync-cursor-plugin-skills.sh`).
Codex applies an install-time packaging transform.

Refresh skills and host config without touching piece data:

```bash
pysar init --force          # current host (default Claude)
pysar init --cursor --force
pysar init --codex --force
```

## Persistence rule

Author content under `.pysar/**` is written through MCP tools
(for example `save_intake_bundle`, `require_piece_stages`,
`save_draft_bundle`, …, `export_piece_to_root`). Skills are written to
call those tools — do not bypass them with raw file writes into
`.pysar/pieces/`. Finished-piece export lands under `export_dir` from
`.pysar/project` (or the project root when unset); `export_piece_to_root`
accepts an optional `export_dir` override for one call and returns the
resolved path (see [Export](./export.md)).

When `/ps --research` or `/ps --seo` is set, the orchestrator arms
`require_piece_stages` so later saves fail closed until that stage file exists
(see [Pipeline](./pipeline.md)).

## Related skills

Orchestrator: `/ps`  
Onboarding: `/ps-onboard`  
Stages and helpers: see [Run the pipeline](./pipeline.md).
