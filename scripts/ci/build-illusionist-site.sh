#!/usr/bin/env bash
# Build the whole RGB site using the pinned Illusionist release kit.
set -euo pipefail

out="${1:?usage: build-illusionist-site.sh <empty-output-directory>}"
archive=illusionist-rgb-site-v1.0.0.tar.gz
digest=acbdb7175b8b7da15918f4ad6884d9aba90bab880bced96e48a69a91a1c80eb1
url="https://github.com/SergioLacerda/illusionist-web/releases/download/rgb-site-v1.0.0"

if [ -e "$out" ] && [ -n "$(ls -A "$out")" ]; then
  echo "output directory is not empty: $out" >&2
  exit 1
fi
work="$(mktemp -d)"
trap 'rm -rf "$work"' EXIT
if [ -n "${ILLUSIONIST_KIT_DIR:-}" ]; then
  cp "$ILLUSIONIST_KIT_DIR/$archive" "$ILLUSIONIST_KIT_DIR/$archive.sha256" "$work/"
else
  curl -fsSL "$url/$archive" -o "$work/$archive"
  curl -fsSL "$url/$archive.sha256" -o "$work/$archive.sha256"
fi
(cd "$work" && sha256sum -c "$archive.sha256")
printf '%s  %s\n' "$digest" "$work/$archive" | sha256sum -c -
mkdir "$work/kit"
tar -xzf "$work/$archive" -C "$work/kit"
mkdir -p "$work/kit/docs"
cp -R docs/core "$work/kit/docs/core"

"$work/kit/illusionist" validate "$work/kit/illusionist.yaml"
"$work/kit/illusionist" build --out "$out" "$work/kit/illusionist.yaml"
