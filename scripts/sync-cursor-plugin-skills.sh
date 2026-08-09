#!/usr/bin/env bash
# Sync host-neutral skill corpus into the canonical Cursor Plugin package.
# Source of truth: cmd/pysar/assets/skills
# Projection:      plugins/pysar/skills
# (dec-20260809-cursor-marketplace-v1-dual-discovery-8b748a7a)
set -euo pipefail
root="$(cd "$(dirname "$0")/.." && pwd)"
src="$root/cmd/pysar/assets/skills"
dst="$root/plugins/pysar/skills"
if [[ ! -d "$src" ]]; then
  echo "sync-cursor-plugin-skills: missing $src" >&2
  exit 1
fi
mkdir -p "$dst"
rsync -a --delete "$src/" "$dst/"
echo "sync-cursor-plugin-skills: $src -> $dst"
