---
title: Init a project
slug: init
nav_order: 20
section: journey
---

# Init a project

`pysar init` scaffolds a writing project for a host agent. Default host is
Claude Code. Cursor is supported. Codex is not yet supported.

## Claude Code (default)

```bash
mkdir my-piece && cd my-piece
pysar init
# same as:
pysar init --claude
```

Or pass a directory:

```bash
pysar init --claude ./my-piece
```

Typical layout:

| Path | Role |
|------|------|
| `.pysar/project` | Project manifest (host + project identity) |
| `CLAUDE.md` | Project instructions for the agent |
| `.claude/settings.json` | Permission allowlist for `mcp__pysar__*` tools |
| `.mcp.json` | MCP server entry that runs `pysar serve` |
| Skills under `~/.claude/skills/ps-*` | Shared `ps-*` skill corpus (global install) |

Open the folder in Claude Code. MCP is pre-approved by the scaffolded settings.

## Cursor

```bash
mkdir my-piece && cd my-piece
pysar init --cursor
```

Typical layout:

| Path | Role |
|------|------|
| `.pysar/project` | Project manifest |
| `.cursor/mcp.json` | MCP server with `PYSAR_PROJECT_ROOT=${workspaceFolder}` |
| Skills under `~/.cursor/skills/ps-*` | Same shared skill corpus as Claude |

Open the folder in Cursor and enable the `pysar` MCP server if prompted.
Cursor does not ship a Claude-style settings allowlist; you may see
permission prompts for MCP tools on first use.

## Flags

| Flag | Meaning |
|------|---------|
| `--claude` | Scaffold for Claude Code (default if no host flag) |
| `--cursor` | Scaffold for Cursor |
| `--codex` | Not yet supported (exits with an error) |
| `--force` | Refresh host project files, skills, and MCP/settings to the shipped version — **never** overwrites `.pysar/` piece data |

Host flags are mutually exclusive.

## After init

Optional voice/style onboarding (~2 minutes):

```text
/ps-onboard
```

Or skip and start writing:

```text
/ps a habit that actually helped our team ship faster
```

See [Run the pipeline](./pipeline.md).
