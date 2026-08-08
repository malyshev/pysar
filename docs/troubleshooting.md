---
title: Troubleshooting
slug: troubleshooting
nav_order: 60
section: journey
---

# Troubleshooting

## `pysar: command not found`

- Re-run the [install](./install.md) script, or download a release binary.
- Ensure the install directory is on `PATH` (`~/.local/bin` is common).
- Check for a shadowing binary: `which -a pysar` and `pysar --version`.

## Wrong or old version

Release installs often land in `~/.local/bin`. A `go install` binary may land
in `$(go env GOPATH)/bin`. Whichever directory appears first on `PATH` wins.

## Host flags rejected together

`--claude`, `--cursor`, and `--codex` are mutually exclusive. Pass only one.

## MCP server not connecting

1. Confirm `pysar` is on `PATH` inside the host (same binary the terminal uses).
2. Claude: check project `.mcp.json` and that Claude Code loaded the server.
3. Cursor: check `.cursor/mcp.json`, enable the `pysar` server if prompted, and
   confirm `PYSAR_PROJECT_ROOT` is `${workspaceFolder}`.
4. Codex: check `.codex/config.toml` for `[mcp_servers.pysar]`,
   `PYSAR_PROJECT_ROOT = "."`, and `default_tools_approval_mode = "approve"`.
5. From the project directory, you can smoke-test the binary with
   `pysar --help` (stdio MCP is started by the host, not by a long-running
   daemon you manage yourself).

## Codex keeps asking to allow pysar MCP tools

Shipped Codex config sets `default_tools_approval_mode = "approve"` (Codex
skip-review dial for `pysar` MCP tools; `auto` is not enough when tools lack
annotations). Refresh with:

```bash
pysar init --codex --force
```

Project `.codex/config.toml` applies only when Codex treats the project as
**trusted**. If prompts continue after refresh, trust the project in Codex and
retry. Untrusted projects may still prompt per tool.

## Skills missing (`/ps` or `$ps` unknown)

Re-run init with `--force` for your host so global skills refresh:

```bash
pysar init --force
# or
pysar init --cursor --force
pysar init --codex --force
```

Restart the host agent after skill install so it rediscovers `~/.claude/skills`,
`~/.cursor/skills`, or `~/.agents/skills` (Codex).

## Cursor keeps asking to approve MCP tools

Expected for now: Cursor scaffold does not include a Claude-style settings
allowlist. Approve the `pysar` tools when prompted, or adjust Cursor’s MCP
permissions for the project.

## Piece files not appearing under `.pysar/pieces/`

- Confirm the MCP server is connected for this workspace.
- Skills must use MCP save tools — not manual writes into `.pysar/`.
- Check the agent’s tool errors for `save_*_bundle` / path anchoring failures.

## Export missing or stale

- Export needs at least a draft for the piece.
- `/ps` exports at the end of a successful run; a stopped `--review` run may
  not have reached export yet.
- Re-run export (or finish `/ps`); re-export overwrites the previous root file.

## Still stuck

Open an issue at [github.com/malyshev/pysar](https://github.com/malyshev/pysar)
with `pysar --version`, host (Claude or Cursor), and the exact command or
skill you ran.
