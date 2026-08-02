#!/usr/bin/env bash
set -euo pipefail

public_dir="${1:?usage: check-release-supply-chain.sh <public_dir> <basename> <version> <manifest> <checksums> <sbom> <provenance>}"
basename="${2:?usage: check-release-supply-chain.sh <public_dir> <basename> <version> <manifest> <checksums> <sbom> <provenance>}"
version="${3:?usage: check-release-supply-chain.sh <public_dir> <basename> <version> <manifest> <checksums> <sbom> <provenance>}"
manifest="${4:?usage: check-release-supply-chain.sh <public_dir> <basename> <version> <manifest> <checksums> <sbom> <provenance>}"
checksums="${5:?usage: check-release-supply-chain.sh <public_dir> <basename> <version> <manifest> <checksums> <sbom> <provenance>}"
sbom="${6:?usage: check-release-supply-chain.sh <public_dir> <basename> <version> <manifest> <checksums> <sbom> <provenance>}"
provenance="${7:?usage: check-release-supply-chain.sh <public_dir> <basename> <version> <manifest> <checksums> <sbom> <provenance>}"

required_files=(
  "${manifest}"
  "${checksums}"
  "${sbom}"
  "${provenance}"
  "${public_dir}/${basename}-latest-en.pdf"
  "${public_dir}/${basename}-${version}-en.pdf"
  "${public_dir}/${basename}-latest-pt-br.pdf"
  "${public_dir}/${basename}-${version}-pt-br.pdf"
)

for file in "${required_files[@]}"; do
  if [[ ! -s "${file}" ]]; then
    printf 'missing or empty release supply-chain file: %s\n' "${file}" >&2
    exit 1
  fi
done

if ! grep -q '"spdxVersion": "SPDX-2.3"' "${sbom}"; then
  printf 'SBOM does not declare SPDX-2.3: %s\n' "${sbom}" >&2
  exit 1
fi

if ! grep -q '"schema": "rgb-system-release-provenance/0.1"' "${provenance}"; then
  printf 'provenance schema mismatch: %s\n' "${provenance}" >&2
  exit 1
fi

for file in "${basename}-${version}-en.pdf" "${basename}-${version}-pt-br.pdf"; do
  if ! grep -q "\"fileName\": \"${file}\"" "${sbom}"; then
    printf 'SBOM missing versioned artifact: %s\n' "${file}" >&2
    exit 1
  fi
done

if ! grep -q "\"git_tag\": \"${version}\"" "${provenance}"; then
  printf 'provenance does not record release tag %s\n' "${version}" >&2
  exit 1
fi

printf 'release supply-chain metadata valid\n'
