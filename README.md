# Pysar

An author-directed editorial engine for your writing projects.

Bring your take — an idea or a rough draft. Pysar helps you shape it into a
piece you trust, without forcing pipeline jargon on you and without posting
on your behalf. It runs as a CLI plus an MCP server and a set of `ps-*`
skills for Claude Code, Cursor, and Codex.

**User guide:** start at [docs/index.md](./docs/index.md)
(install → init → pipeline → MCP/skills → export → troubleshooting).

```bash
curl -fsSL https://raw.githubusercontent.com/malyshev/pysar/master/install.sh | bash
pysar init          # Claude Code (default)
pysar init --cursor # Cursor
pysar init --codex  # Codex CLI / App
```

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
