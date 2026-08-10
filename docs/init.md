---
title: Init a project
slug: init
nav_order: 20
section: journey
---

# Init a project

`pysar init` scaffolds a writing project for a host agent. Default host is
Claude Code. Cursor and Codex CLI/App are also supported.

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

Install the Pysar binary first (`~/.local/bin/pysar`), then scaffold a project
and load the **Pysar Cursor plugin** for `/ps` skills.

```bash
mkdir my-piece && cd my-piece
pysar init --cursor
```

Typical layout after init:

| Path | Role |
|------|------|
| `.pysar/project` | Project manifest |
| `.cursor/mcp.json` | MCP server (`command` = `${userHome}/.local/bin/pysar`, `PYSAR_PROJECT_ROOT=${workspaceFolder}`) — same spawn contract as the plugin |

`/ps` skills come from the Cursor plugin package (`plugins/pysar` in the repo),
not from `~/.cursor/skills`.

**Two different installs (do not conflate them)**

| Path | What you get |
|------|----------------|
| **Marketplace** or **local plugin copy** | Full package: `/ps` skills **and** plugin MCP |
| getpysar.com **Install in Cursor** deeplink | **MCP only** (stdio server spawn). Skills still need Marketplace or a local plugin copy |

**Plugin install options**

1. **Marketplace** — Customize → search **pysar** (after the listing is live).
2. **Local dogfood** — from a **pysar monorepo** clone (not the writing project):

   ```bash
   mkdir -p ~/.cursor/plugins/local
   rsync -a --delete "$(pwd)/plugins/pysar/" ~/.cursor/plugins/local/pysar/
   ```

   Then **Developer: Reload Window**. Cursor rejects symlinks whose target is
   outside `~/.cursor/plugins/local` — use a real copy and re-`rsync` after skill
   sync.
3. **MCP deeplink** — [getpysar.com](https://getpysar.com) **Install in Cursor**
   (and the link printed by `pysar init --cursor`) enables the server in
   Customize → MCPs. That is not a substitute for the plugin skills package.

`pysar init --cursor` also registers user-scope MCP (`~/.cursor/mcp.json`) for
Customize → Connected while the plugin path settles. Reload MCP or restart
Cursor once, then open the project folder.

## Codex CLI / App

```bash
mkdir my-piece && cd my-piece
pysar init --codex
```

Typical layout:

| Path | Role |
|------|------|
| `.pysar/project` | Project manifest |
| `.codex/config.toml` | MCP server entry that runs `pysar serve` (`default_tools_approval_mode = "approve"`) |
| Skills under `~/.agents/skills/ps-*` | Same shared skill corpus, Codex-packaged (`$ps-*` + `agents/openai.yaml`) |

Only the orchestrator skill (`ps`) allows implicit invocation; stage skills are
explicit (`$ps-intake`, …). Open the folder in Codex after init.

Project `.codex/config.toml` MCP settings (including tool approval mode) apply
when Codex treats the project as **trusted**. If you still see per-tool
“Allow the pysar MCP server…” prompts, trust the project and re-run
`pysar init --codex --force`.

## Flags

| Flag | Meaning |
|------|---------|
| `--claude` | Scaffold for Claude Code (default if no host flag) |
| `--cursor` | Scaffold for Cursor |
| `--codex` | Scaffold for Codex CLI / App |
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
