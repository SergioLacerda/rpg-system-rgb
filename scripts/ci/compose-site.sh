#!/usr/bin/env bash
# Compose the Pages site: ILLUSIONISTA RGB surface + RGB Library + RGB downloads.
# usage: compose-site.sh <site_dir> <artifacts_dir>
#   site_dir       staging directory to create (must not exist or be empty)
#   artifacts_dir  directory holding the downloaded CI artifacts as
#                  <artifacts_dir>/pdfs and <artifacts_dir>/skills
# Required env: ILLUSIONIST_REPOSITORY, ILLUSIONIST_RGB_VERSION,
#               ILLUSIONIST_RGB_ARTIFACT, ILLUSIONIST_RGB_CHECKSUM,
#               ILLUSIONIST_RGB_SHA256
set -euo pipefail

site_dir="${1:?usage: compose-site.sh <site_dir> <artifacts_dir>}"
artifacts_dir="${2:?usage: compose-site.sh <site_dir> <artifacts_dir>}"
repo="${ILLUSIONIST_REPOSITORY:?}"
version="${ILLUSIONIST_RGB_VERSION:?}"
artifact="${ILLUSIONIST_RGB_ARTIFACT:?}"
checksum="${ILLUSIONIST_RGB_CHECKSUM:?}"
expected_sha256="${ILLUSIONIST_RGB_SHA256:?}"
tracked_downloads="web/landing/public/downloads"

fail() {
  printf '::error::%s\n' "$1"
  exit 1
}

case "${version}" in
  latest | main | "") fail "ILLUSIONIST_RGB_VERSION must be an exact tag, got '${version}'" ;;
esac

if [ -e "${site_dir}" ] && [ -n "$(ls -A "${site_dir}")" ]; then
  fail "staging directory ${site_dir} is not empty"
fi
mkdir -p "${site_dir}"

download_dir="$(mktemp -d)"
trap 'rm -rf "${download_dir}"' EXIT

base_url="https://github.com/${repo}/releases/download/${version}"
echo "downloading ${base_url}/${artifact}"
curl -fsSL -o "${download_dir}/${artifact}" "${base_url}/${artifact}"
curl -fsSL -o "${download_dir}/${checksum}" "${base_url}/${checksum}"

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
  [ -d "${artifacts_dir}/${sub}" ] || fail "missing artifact directory ${artifacts_dir}/${sub}"
  cp -R "${artifacts_dir}/${sub}/." "${site_dir}/downloads/"
done

# Composed tree checks: everything the surface links to must exist.
required=(
  library/index.html
  downloads/rgb-system-core-v2-latest-pt-br.pdf
  downloads/rgb-system-core-v2-latest-en.pdf
  downloads/rgb-specialist-latest.zip
  downloads/SHA256SUMS
  downloads/rgb-specialist-SHA256SUMS
)
for f in "${required[@]}"; do
  [ -f "${site_dir}/${f}" ] || fail "composed site is missing ${f}"
done

# Baseline parity: every tracked download must still be published.
while IFS= read -r f; do
  [ -f "${site_dir}/downloads/${f}" ] || fail "composed site lost tracked download ${f}"
done < <(cd "${tracked_downloads}" && find . -type f | sed 's|^\./||')

echo "composed site ready at ${site_dir}"
