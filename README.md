<p align="center">
  <img src="./plugins/pysar/assets/pysar-logo-mark-512.png" width="160" alt="Pysar logo">
</p>

<p align="center">
  <a href="https://getpysar.com"><img src="https://img.shields.io/website?down_color=red&down_message=offline&style=flat&url=https%3A%2F%2Fgetpysar.com" alt="Website Status"></a>
  <a href="./LICENSE"><img src="https://img.shields.io/github/license/malyshev/pysar" alt="License"></a>
</p>

# Pysar

An author-directed editorial engine for your writing projects.

Bring your take — an idea or a rough draft. Pysar helps you shape it into a
piece you trust, without forcing pipeline jargon on you and without posting
on your behalf. It runs as a CLI plus an MCP server and a set of `ps-*`
skills for Claude Code, Cursor, and Codex.

**User guide:** start at [docs/index.md](./docs/index.md)
(install → init → pipeline → MCP/skills → export → troubleshooting).

```bash
curl -fsSL https://getpysar.com/install.sh | bash
pysar init          # Claude Code (default)
pysar init --cursor # Cursor (project scaffold; install the Pysar Cursor plugin for /ps skills)
pysar init --codex  # Codex CLI / App
```

Cursor plugin package (Marketplace / Install in Cursor / local dogfood):
[`plugins/pysar`](./plugins/pysar/).

## Development

```bash
go build ./cmd/pysar     # build
go install ./cmd/pysar   # build + install to $(go env GOPATH)/bin
go test ./...            # run tests
```

Public site (getpysar.com) lives under [`site/`](./site/) — Next.js static
export, separate from the Go tool. See [`site/README.md`](./site/README.md).

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
