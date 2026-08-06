#!/bin/bash
# Pysar installer
#
# Installs the pysar binary globally from the latest GitHub Release.
# After installation, run `pysar init` in your writing project.
#
# Usage:
#   curl -fsSL https://raw.githubusercontent.com/malyshev/pysar/main/install.sh | bash

set -e

BOLD='\033[1m'
DIM='\033[2m'
RESET='\033[0m'
RED='\033[31m'
GREEN='\033[32m'
CYAN='\033[36m'

REPO="malyshev/pysar"
BIN_NAME="pysar"
BIN_DIRS=("$HOME/.local/bin" "/usr/local/bin")

get_os_arch() {
    local os arch
    os=$(uname -s | tr '[:upper:]' '[:lower:]')
    arch=$(uname -m)
    case "$os" in
        linux|darwin) ;;
        *)
            printf "${RED}✗ Unsupported OS: %s${RESET}\n" "$os"
            printf "  Windows users: download the .zip from https://github.com/%s/releases/latest and extract %s.exe onto your PATH.\n" "$REPO" "$BIN_NAME"
            exit 1
            ;;
    esac
    case "$arch" in
        x86_64) arch="amd64" ;;
        arm64|aarch64) arch="arm64" ;;
        *)
            printf "${RED}✗ Unsupported architecture: %s${RESET}\n" "$arch"
            exit 1
            ;;
    esac
    echo "${os}-${arch}"
}

find_bin_dir() {
    local dir
    for dir in "${BIN_DIRS[@]}"; do
        if [[ -d "$dir" && -w "$dir" ]]; then
            echo "$dir"
            return 0
        fi
    done
    mkdir -p "$HOME/.local/bin"
    echo "$HOME/.local/bin"
}

main() {
    printf "${BOLD}Pysar installer${RESET}\n\n"

    local os_arch bin_dir download_url tmp_dir archive_name archive_path
    os_arch=$(get_os_arch)
    bin_dir=$(find_bin_dir)

    printf "  Detected platform: ${CYAN}%s${RESET}\n" "$os_arch"

    local api_url="https://api.github.com/repos/${REPO}/releases/latest"
    download_url=$(curl -fsSL "$api_url" \
        | grep "browser_download_url.*${os_arch}\.tar\.gz" \
        | head -n1 \
        | sed -E 's/.*"browser_download_url": "([^"]+)".*/\1/')

    if [[ -z "$download_url" ]]; then
        printf "${RED}✗ No release asset found for %s.${RESET}\n" "$os_arch"
        printf "  Check https://github.com/%s/releases/latest for available downloads.\n" "$REPO"
        exit 1
    fi

    tmp_dir=$(mktemp -d)
    trap 'rm -rf "$tmp_dir"' EXIT
    archive_name=$(basename "$download_url")
    archive_path="$tmp_dir/$archive_name"

    printf "  Downloading %s...\n" "$archive_name"
    curl -fsSL "$download_url" -o "$archive_path"

    tar -xzf "$archive_path" -C "$tmp_dir"

    if [[ ! -f "$tmp_dir/$BIN_NAME" ]]; then
        printf "${RED}✗ Expected binary '%s' not found in downloaded archive.${RESET}\n" "$BIN_NAME"
        exit 1
    fi

    install -m 0755 "$tmp_dir/$BIN_NAME" "$bin_dir/$BIN_NAME"
    printf "  ${GREEN}✓${RESET} Installed to ${BOLD}%s/%s${RESET}\n\n" "$bin_dir" "$BIN_NAME"

    if ! command -v "$BIN_NAME" >/dev/null 2>&1; then
        printf "${DIM}%s is not on your PATH yet. Add this to your shell profile:${RESET}\n" "$bin_dir"
        printf "  export PATH=\"%s:\$PATH\"\n\n" "$bin_dir"
    fi

    printf "Next: run ${BOLD}pysar init${RESET} in your writing project.\n"
}

main
