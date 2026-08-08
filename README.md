# Pysar

An author-directed editorial engine for your writing projects.

Bring your take — an idea or a rough draft. Pysar helps you shape it into a
piece you trust, without forcing pipeline jargon on you and without posting
on your behalf. It runs as a CLI plus an MCP server and a set of `ps-*`
skills for Claude Code.

## Install

```bash
curl -fsSL https://raw.githubusercontent.com/malyshev/pysar/master/install.sh | bash
```

This installs the `pysar` binary for macOS and Linux (amd64/arm64) from the
[latest release](https://github.com/malyshev/pysar/releases/latest). Windows
users: download the `.zip` for your platform from that page and extract
`pysar.exe` onto your `PATH`.

Building from source (requires Go, matching the version in `go.mod`):

```bash
go build -o pysar ./cmd/pysar
```

## Quick start

```bash
pysar init --claude   # scaffold a Claude Code project in the current directory
pysar init --cursor   # scaffold a Cursor project (.cursor/mcp.json + skills)
```

Then open the project in your host agent (Claude Code: MCP is pre-approved by
the scaffolded `.claude/settings.json`; Cursor: enable the `pysar` MCP server
from `.cursor/mcp.json` if prompted) and either:

- Run `/ps-onboard` first to teach Pysar your voice and style (~2 minutes,
  optional — Pysar works with a general-audience default if you skip it), or
- Just run `/ps <your idea in a sentence or two>` and let it walk the whole
  pipeline (intake → draft → staff-edit → sharpen → optional SEO packaging
  with `--seo` → humanize) and export a finished piece to your project root.

Run `/ps` with no arguments for a plain-language walkthrough.

## Development

```bash
go build ./cmd/pysar     # build
go install ./cmd/pysar   # build + install to $(go env GOPATH)/bin
go test ./...            # run tests
```

Project governance, decision history, and contribution discipline live in
[CLAUDE.md](./CLAUDE.md) and the `.haft/` directory (this project uses
[haft](https://github.com/m0n0x41d/haft) for structured decision recording).

## Releases

Releases are built and published by CI ([`.github/workflows/release.yml`](.github/workflows/release.yml))
via [GoReleaser](https://goreleaser.com/) whenever a `v*` tag is pushed. Tests
must pass before a build runs — see [`.goreleaser.yaml`](.goreleaser.yaml)
for the build matrix.

## License

[MIT](./LICENSE)
