#!/usr/bin/env bash
# Compose the Pages site: ILLUSIONISTA RGB surface + RGB Library + RGB downloads.
# usage: compose-site.sh <site_dir> <artifacts_dir>
#   site_dir       staging directory to create (must not exist or be empty)
#   artifacts_dir  directory holding the downloaded CI artifacts as
#                  <artifacts_dir>/pdfs and <artifacts_dir>/skills
# Required env: ILLUSIONIST_REPOSITORY, ILLUSIONIST_RGB_VERSION,
#               ILLUSIONIST_RGB_ARTIFACT, ILLUSIONIST_RGB_CHECKSUM,
#               ILLUSIONIST_RGB_SHA256, GH_TOKEN
set -euo pipefail
site_dir="${1:?usage: compose-site.sh <site_dir> <artifacts_dir>}"
artifacts_dir="${2:?usage: compose-site.sh <site_dir> <artifacts_dir>}"

bash scripts/ci/build-illusionist-site.sh "$site_dir"
if [ -e "$site_dir/downloads" ] && [ -n "$(ls -A "$site_dir/downloads")" ]; then
  echo 'site kit unexpectedly supplied downloads; refusing to overwrite' >&2
  exit 1
}

[ -n "${GH_TOKEN:-}" ] || fail "GH_TOKEN is required to download the private ILUSIONISTA release"

case "${version}" in
  latest | main | "") fail "ILLUSIONIST_RGB_VERSION must be an exact tag, got '${version}'" ;;
esac

if [ -e "${site_dir}" ] && [ -n "$(ls -A "${site_dir}")" ]; then
  fail "staging directory ${site_dir} is not empty"
fi
mkdir -p "${site_dir}"

download_dir="$(mktemp -d)"
trap 'rm -rf "${download_dir}"' EXIT

echo "downloading ${repo} release ${version}"
gh release download "${version}" --repo "${repo}" \
  --pattern "${artifact}" --pattern "${checksum}" --dir "${download_dir}"

# Integrity: the release checksum asset and the digest pinned in ci.yml must both match.
(cd "${download_dir}" && sha256sum -c "${checksum}")
actual_sha256="$(sha256sum "${download_dir}/${artifact}" | cut -d' ' -f1)"
if [ "${actual_sha256}" != "${expected_sha256}" ]; then
  fail "${artifact} digest ${actual_sha256} does not match pinned ${expected_sha256}"
fi

tar -xzf "${download_dir}/${artifact}" -C "${site_dir}"

# Surface sanity checks (derived from the real rgb-v0.1.0 archive).
for f in index.html pt-br/index.html en/index.html; do
  [ -f "${site_dir}/${f}" ] || fail "surface is missing ${f}"
done
[ -d "${site_dir}/_astro" ] || fail "surface is missing _astro/"

# Collision guard: RGB composition points must not carry real surface content.
for d in library downloads; do
  if [ -e "${site_dir}/${d}" ] && [ -n "$(ls -A "${site_dir}/${d}")" ]; then
    fail "surface already provides non-empty ${d}/; contract problem, refusing to overwrite"
  fi
done

# Library: RGB-owned generator, written straight into staging.
make docs-build LIBRARY_DIR="${site_dir}/library"

# Downloads: tracked baseline files first, then the CI artifacts on top.
mkdir -p "${site_dir}/downloads"
cp -R "${tracked_downloads}/." "${site_dir}/downloads/"
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
