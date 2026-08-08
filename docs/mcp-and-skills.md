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
3. **Skills** — `ps-*` instructions the host agent follows (`/ps`, `/ps-intake`, …)

## MCP config

After init, the host launches:

```text
pysar serve
```

with the project root set so piece I/O lands in that project’s `.pysar/`.

| Host | Config file | Project root env |
|------|-------------|------------------|
| Claude Code | `.mcp.json` at project root | `PYSAR_PROJECT_ROOT=${PWD:-.}` |
| Cursor | `.cursor/mcp.json` | `PYSAR_PROJECT_ROOT=${workspaceFolder}` |

You do not need to run `pysar serve` in a separate terminal for normal use —
the host agent starts it as an MCP server.

## Skills install location

`pysar init` installs the shared skill corpus globally:

| Host | Skills directory |
|------|------------------|
| Claude Code | `~/.claude/skills/ps-*` |
| Cursor | `~/.cursor/skills/ps-*` |

Bodies are the same across hosts. Only paths and MCP/settings files differ.

Refresh skills and host config without touching piece data:

```bash
pysar init --force          # current host (default Claude)
pysar init --cursor --force
```

## Persistence rule

Author content under `.pysar/**` is written through MCP tools
(for example `save_intake_bundle`, `save_draft_bundle`, …,
`export_piece_to_root`). Skills are written to call those tools — do not
bypass them with raw file writes into `.pysar/pieces/`.

## Related skills

Orchestrator: `/ps`  
Onboarding: `/ps-onboard`  
Stages and helpers: see [Run the pipeline](./pipeline.md).
