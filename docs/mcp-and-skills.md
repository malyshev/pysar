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

| Host | Config file | Project root env |
|------|-------------|------------------|
| Claude Code | `.mcp.json` at project root | `PYSAR_PROJECT_ROOT=${PWD:-.}` |
| Cursor | `.cursor/mcp.json` | `PYSAR_PROJECT_ROOT=${workspaceFolder}` |
| Codex CLI / App | `.codex/config.toml` | `PYSAR_PROJECT_ROOT = "."` (`default_tools_approval_mode = "approve"`) |

Claude pre-approves `mcp__pysar__*` via `.claude/settings.json`. Codex uses
`default_tools_approval_mode = "approve"` on the MCP server block for the same
friction goal (`auto` still prompts when tools lack risk annotations); project
`.codex/config.toml` applies for **trusted** projects.

You do not need to run `pysar serve` in a separate terminal for normal use —
the host agent starts it as an MCP server.

## Skills install location

`pysar init` installs the shared skill corpus globally:

| Host | Skills directory |
|------|------------------|
| Claude Code | `~/.claude/skills/ps-*` |
| Cursor | `~/.cursor/skills/ps-*` |
| Codex CLI / App | `~/.agents/skills/ps-*` (Codex packaging: `$ps-*` + `agents/openai.yaml`) |

Editorial skill bodies share one corpus. Claude/Cursor install the corpus
bytes as-is; Codex applies an install-time packaging transform. Only paths,
MCP/settings, and that packaging differ.

Refresh skills and host config without touching piece data:

```bash
pysar init --force          # current host (default Claude)
pysar init --cursor --force
pysar init --codex --force
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
