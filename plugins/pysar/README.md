# Pysar — Cursor Plugin

Canonical Cursor Plugin package for [Pysar](https://getpysar.com)
(`dec-20260809-cursor-marketplace-v1-dual-discovery-8b748a7a`).

This directory is the **single** plugin product. Cursor Marketplace and the
getpysar.com **Install in Cursor** path both expose this package — not two
different products.

## Contents

| Path | Role |
|------|------|
| `.cursor-plugin/plugin.json` | Cursor Plugin manifest (`name`: `pysar`) |
| `mcp.json` | stdio MCP → `${userHome}/.local/bin/pysar serve` |
| `skills/ps-*/SKILL.md` | Projection of `cmd/pysar/assets/skills` (host-neutral corpus) |
| `assets/pysar-logo-mark-512.png` | Marketplace logo (default / light UI) |
| `assets/pysar-logo-mark-dark-512.png` | Marketplace logo for dark UI |
| `assets/pysar-logo-mark-*.svg` | Vector sources (light / dark) |

## Prerequisites

1. Install the Pysar binary (documented path `~/.local/bin/pysar`):

   ```bash
   curl -fsSL https://getpysar.com/install.sh | bash
   ```

2. Scaffold a writing project:

   ```bash
   pysar init --cursor
   ```

## Install this plugin

### Local dogfood (`~/.cursor/plugins/local`)

From the **pysar monorepo root** (not a writing project):

```bash
mkdir -p ~/.cursor/plugins/local
rsync -a --delete "$(pwd)/plugins/pysar/" ~/.cursor/plugins/local/pysar/
```

Cursor currently **rejects** symlinks whose target is outside
`~/.cursor/plugins/local` (silent skip; Plugins log:
`loadUserLocalPlugin pysar rejected: symlink target … is outside …/local`).
Use a real copy for dogfood; re-run `rsync` after skill sync.

Reload Cursor (**Developer: Reload Window**). **Customize → Plugins** should
list **Pysar** (Installed). Enable if prompted. Then check **MCPs** for Plugin
MCP `pysar`, and run `/ps` with a one-sentence idea.

Note: `pysar init --cursor` prints an MCP deeplink that only installs the
stdio server — that is **not** the same as loading this plugin package.

### Marketplace

Submit this repository at
[cursor.com/marketplace/publish](https://cursor.com/marketplace/publish).
The multi-plugin manifest is `.cursor-plugin/marketplace.json` with source
`plugins/pysar`.

### Install in Cursor (site deeplink)

getpysar.com ships an MCP install deeplink whose config matches this
`mcp.json` (server name `pysar`). That is a discovery path onto the same
spawn contract — not a second package. Full skills still come from this
plugin (local or Marketplace).

## Keep skills in sync

```bash
./scripts/sync-cursor-plugin-skills.sh
go test ./cmd/pysar -run TestCursorPlugin
```

Do not edit `plugins/pysar/skills` by hand — edit `cmd/pysar/assets/skills`
and re-sync.
