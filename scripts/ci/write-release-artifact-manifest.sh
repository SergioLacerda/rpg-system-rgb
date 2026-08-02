#!/usr/bin/env bash
set -euo pipefail

public_dir="${1:?usage: write-release-artifact-manifest.sh <public_dir> <basename> <version> <manifest> <checksums>}"
basename="${2:?usage: write-release-artifact-manifest.sh <public_dir> <basename> <version> <manifest> <checksums>}"
version="${3:?usage: write-release-artifact-manifest.sh <public_dir> <basename> <version> <manifest> <checksums>}"
manifest="${4:?usage: write-release-artifact-manifest.sh <public_dir> <basename> <version> <manifest> <checksums>}"
checksums="${5:?usage: write-release-artifact-manifest.sh <public_dir> <basename> <version> <manifest> <checksums>}"

mkdir -p "${public_dir}"

: "${GOCACHE:=/tmp/go-cache}"
export GOCACHE

go run ./cmd/rgb release manifest \
  --public-dir "${public_dir}" \
  --basename "${basename}" \
  --version "${version}" \
  --manifest "${manifest}" \
  --checksums "${checksums}"
