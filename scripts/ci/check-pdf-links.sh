#!/usr/bin/env bash
set -euo pipefail

base="${1:?usage: check-pdf-links.sh <base_url>}"
base="${base%/}/"

for locale in en pt-br; do
  url="${base}downloads/rgb-system-core-v2-latest-${locale}.pdf"
  echo "checking ${url}"
  size=$(curl -sSfL -o /tmp/check.pdf -w '%{size_download}' "${url}")
  if [ "${size}" -lt 1000 ]; then
    echo "::error::${url} returned suspiciously small content (${size} bytes)"
    exit 1
  fi
  echo "OK: ${url} (${size} bytes)"
done
