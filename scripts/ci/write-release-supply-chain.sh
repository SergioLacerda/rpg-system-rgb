#!/usr/bin/env bash
set -euo pipefail

public_dir="${1:?usage: write-release-supply-chain.sh <public_dir> <basename> <version> <manifest> <checksums> <sbom> <provenance>}"
basename="${2:?usage: write-release-supply-chain.sh <public_dir> <basename> <version> <manifest> <checksums> <sbom> <provenance>}"
version="${3:?usage: write-release-supply-chain.sh <public_dir> <basename> <version> <manifest> <checksums> <sbom> <provenance>}"
manifest="${4:?usage: write-release-supply-chain.sh <public_dir> <basename> <version> <manifest> <checksums> <sbom> <provenance>}"
checksums="${5:?usage: write-release-supply-chain.sh <public_dir> <basename> <version> <manifest> <checksums> <sbom> <provenance>}"
sbom="${6:?usage: write-release-supply-chain.sh <public_dir> <basename> <version> <manifest> <checksums> <sbom> <provenance>}"
provenance="${7:?usage: write-release-supply-chain.sh <public_dir> <basename> <version> <manifest> <checksums> <sbom> <provenance>}"

mkdir -p "${public_dir}"

commit="$(git rev-parse HEAD)"
dirty="false"
if ! git diff --quiet || ! git diff --cached --quiet; then
  dirty="true"
fi
tag_present="false"
if git rev-parse -q --verify "refs/tags/${version}" >/dev/null; then
  tag_present="true"
fi
generated_at="$(date -u '+%Y-%m-%dT%H:%M:%SZ')"

artifact_files=(
  "${basename}-latest-en.pdf"
  "${basename}-${version}-en.pdf"
  "${basename}-latest-pt-br.pdf"
  "${basename}-${version}-pt-br.pdf"
  "$(basename "${manifest}")"
  "$(basename "${checksums}")"
  "${basename}-${version}-en.pdf.sha256"
  "${basename}-${version}-pt-br.pdf.sha256"
)

{
  printf '{\n'
  printf '  "schema": "rgb-system-release-provenance/0.1",\n'
  printf '  "basename": "%s",\n' "${basename}"
  printf '  "version": "%s",\n' "${version}"
  printf '  "git_commit": "%s",\n' "${commit}"
  printf '  "git_tag": "%s",\n' "${version}"
  printf '  "git_tag_present": %s,\n' "${tag_present}"
  printf '  "worktree_dirty": %s,\n' "${dirty}"
  printf '  "generated_at": "%s",\n' "${generated_at}"
  printf '  "sbom": "%s",\n' "$(basename "${sbom}")"
  printf '  "manifest": "%s",\n' "$(basename "${manifest}")"
  printf '  "checksums": "%s"\n' "$(basename "${checksums}")"
  printf '}\n'
} >"${provenance}"

{
  printf '{\n'
  printf '  "spdxVersion": "SPDX-2.3",\n'
  printf '  "dataLicense": "CC0-1.0",\n'
  printf '  "SPDXID": "SPDXRef-DOCUMENT",\n'
  printf '  "name": "%s-%s-release-artifacts",\n' "${basename}" "${version}"
  printf '  "documentNamespace": "https://github.com/SergioLacerda/rpg-system-rgb/releases/%s/sbom",\n' "${version}"
  printf '  "creationInfo": {\n'
  printf '    "created": "%s",\n' "${generated_at}"
  printf '    "creators": ["Tool: scripts/ci/write-release-supply-chain.sh"]\n'
  printf '  },\n'
  printf '  "files": [\n'
  first="true"
  for file in "${artifact_files[@]}"; do
    path="${public_dir}/${file}"
    if [[ ! -f "${path}" ]]; then
      continue
    fi
    sha="$(sha256sum "${path}" | awk '{print $1}')"
    if [[ "${first}" == "true" ]]; then
      first="false"
    else
      printf ',\n'
    fi
    printf '    {"SPDXID": "SPDXRef-File-%s", "fileName": "%s", "checksums": [{"algorithm": "SHA256", "checksumValue": "%s"}]}' \
      "$(printf '%s' "${file}" | tr -c 'A-Za-z0-9' '-')" "${file}" "${sha}"
  done
  printf '\n  ]\n'
  printf '}\n'
} >"${sbom}"

printf 'wrote %s\n' "${sbom}"
printf 'wrote %s\n' "${provenance}"
