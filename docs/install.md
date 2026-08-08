---
title: Install
slug: install
nav_order: 10
section: journey
---

# Install

Get the `pysar` binary on your `PATH`, then go to [Init a project](./init.md).

## Recommended: install script (macOS / Linux)

```bash
curl -fsSL https://getpysar.com/install.sh | bash
```

The script detects `darwin`/`linux` and `amd64`/`arm64`, downloads the matching
`.tar.gz` from the
[latest GitHub release](https://github.com/malyshev/pysar/releases/latest),
and installs to the first writable directory among:

1. `~/.local/bin`
2. `/usr/local/bin`

If the chosen directory is not on your `PATH`, the script prints the
`export PATH=...` line to add to your shell profile.

Confirm:

```bash
pysar --version
pysar --help
```

## Windows

Download the Windows `.zip` for your architecture from the
[latest release](https://github.com/malyshev/pysar/releases/latest), extract
`pysar.exe`, and place it on your `PATH`. The install script does not support
Windows.

## Build from source

Requires Go matching the version in [`go.mod`](../go.mod):

```bash
git clone https://github.com/malyshev/pysar.git
cd pysar
go build -o pysar ./cmd/pysar
# optional: install into $(go env GOPATH)/bin
go install ./cmd/pysar
```

If you already have an older `pysar` earlier on `PATH` (for example
`~/.local/bin` vs `~/go/bin`), check which one wins with `which pysar` and
`pysar --version`.

## Next

[Scaffold a writing project](./init.md) with `pysar init`.
