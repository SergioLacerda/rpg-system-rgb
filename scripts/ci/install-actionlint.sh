#!/usr/bin/env bash
set -euo pipefail

version="${1:?actionlint version is required}"
dest_dir="${2:?destination bin directory is required}"

case "$(uname -s)" in
  Linux) os="linux" ;;
  Darwin) os="darwin" ;;
  *) echo "unsupported actionlint OS: $(uname -s)" >&2; exit 1 ;;
esac

case "$(uname -m)" in
  x86_64|amd64) arch="amd64" ;;
  arm64|aarch64) arch="arm64" ;;
  *) echo "unsupported actionlint architecture: $(uname -m)" >&2; exit 1 ;;
esac

archive="actionlint_${version}_${os}_${arch}.tar.gz"
url="https://github.com/rhysd/actionlint/releases/download/v${version}/${archive}"
tmp_dir="$(mktemp -d)"
trap 'rm -rf "$tmp_dir"' EXIT

mkdir -p "$dest_dir"
curl -fsSL "$url" -o "$tmp_dir/$archive"
tar -xzf "$tmp_dir/$archive" -C "$tmp_dir"
install -m 0755 "$tmp_dir/actionlint" "$dest_dir/actionlint"
"$dest_dir/actionlint" -version
