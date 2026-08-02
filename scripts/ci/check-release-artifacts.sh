#!/usr/bin/env bash
set -euo pipefail

public_dir="${1:?usage: check-release-artifacts.sh <public_dir> <basename> <version> <manifest> <checksums>}"
basename="${2:?usage: check-release-artifacts.sh <public_dir> <basename> <version> <manifest> <checksums>}"
version="${3:?usage: check-release-artifacts.sh <public_dir> <basename> <version> <manifest> <checksums>}"
manifest="${4:?usage: check-release-artifacts.sh <public_dir> <basename> <version> <manifest> <checksums>}"
checksums="${5:?usage: check-release-artifacts.sh <public_dir> <basename> <version> <manifest> <checksums>}"

: "${GOCACHE:=/tmp/go-cache}"
export GOCACHE

go run ./cmd/rgb release check \
  --public-dir "${public_dir}" \
  --basename "${basename}" \
  --version "${version}" \
  --manifest "${manifest}" \
  --checksums "${checksums}"
