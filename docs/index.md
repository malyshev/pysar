---
title: Pysar user guide
slug: index
nav_order: 0
section: journey
---

# Pysar user guide

Pysar is an author-directed editorial engine. You bring an idea or a rough draft;
it helps you shape a piece you trust. It does not post on your behalf.

The binary is a small CLI. The writing workflow runs inside a host agent
(Claude Code, Cursor, or Codex) through `ps-*` skills and a local MCP server
(`pysar serve`).

## Start here

1. [Install](./install.md) — get `pysar` on your `PATH`
2. [Init a project](./init.md) — scaffold Claude Code, Cursor, or Codex
3. [Run the pipeline](./pipeline.md) — `/ps` from idea to finished piece
4. [MCP and skills](./mcp-and-skills.md) — how the agent talks to Pysar
5. [Export](./export.md) — finished Markdown in your project root
6. [Troubleshooting](./troubleshooting.md) — common failures

## What you need

- macOS or Linux for the install script (Windows: download a release zip)
- A host agent: **Claude Code** (default), **Cursor**, or **Codex**
- Network only for install and optional research/fact-check steps the agent runs

Contributor build, tests, and release mechanics stay in the repo
[README](../README.md).
