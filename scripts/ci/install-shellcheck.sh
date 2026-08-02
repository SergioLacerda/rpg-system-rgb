#!/usr/bin/env bash
set -euo pipefail

version="${1:?shellcheck version is required}"
dest_dir="${2:?destination bin directory is required}"

case "$(uname -s)" in
  Linux) os="linux" ;;
  Darwin) os="darwin" ;;
  *) echo "unsupported shellcheck OS: $(uname -s)" >&2; exit 1 ;;
esac

case "$(uname -m)" in
  x86_64|amd64) arch="x86_64" ;;
  arm64|aarch64) arch="aarch64" ;;
  *) echo "unsupported shellcheck architecture: $(uname -m)" >&2; exit 1 ;;
esac

archive="shellcheck-v${version}.${os}.${arch}.tar.xz"
url="https://github.com/koalaman/shellcheck/releases/download/v${version}/${archive}"
tmp_dir="$(mktemp -d)"
trap 'rm -rf "$tmp_dir"' EXIT

mkdir -p "$dest_dir"
curl -fsSL "$url" -o "$tmp_dir/$archive"
tar -xJf "$tmp_dir/$archive" -C "$tmp_dir"
install -m 0755 "$tmp_dir/shellcheck-v${version}/shellcheck" "$dest_dir/shellcheck"
"$dest_dir/shellcheck" --version
