#!/usr/bin/env bash
# Compose the verified Illusionist site with RGB's PDF and skill downloads.
set -euo pipefail
site_dir="${1:?usage: compose-site.sh <site_dir> <artifacts_dir>}"
artifacts_dir="${2:?usage: compose-site.sh <site_dir> <artifacts_dir>}"

bash scripts/ci/build-illusionist-site.sh "$site_dir"
if [ -e "$site_dir/downloads" ] && [ -n "$(ls -A "$site_dir/downloads")" ]; then
  echo 'site kit unexpectedly supplied downloads; refusing to overwrite' >&2
  exit 1
fi
mkdir -p "$site_dir/downloads"
cp -R generated/downloads/. "$site_dir/downloads/"
for sub in pdfs skills; do
  [ -d "$artifacts_dir/$sub" ] || { echo "missing artifact directory $artifacts_dir/$sub" >&2; exit 1; }
  cp -R "$artifacts_dir/$sub/." "$site_dir/downloads/"
done
for file in index.html en/index.html library/index.html en/library/index.html \
  library/search-index.json assets/library.js \
  downloads/rgb-system-core-v2-latest-pt-br.pdf downloads/rgb-system-core-v2-latest-en.pdf \
  downloads/rgb-specialist-latest.zip downloads/SHA256SUMS downloads/rgb-specialist-SHA256SUMS; do
  [ -s "$site_dir/$file" ] || { echo "missing site artifact: $file" >&2; exit 1; }
done
while IFS= read -r file; do
  [ -f "$site_dir/downloads/$file" ] || { echo "lost tracked download: $file" >&2; exit 1; }
done < <(cd generated/downloads && find . -type f | sed 's|^\./||')
