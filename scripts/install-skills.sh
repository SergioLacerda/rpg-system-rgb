#!/usr/bin/env bash
# Install or update packaged RGB skills (rgb-specialist, rgb-maker, and rpg-maker alias)
# into local agent skill directories (~/.gemini/config/skills and ~/.agents/skills).
set -euo pipefail

repo_root="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
default_targets=(
  "${HOME}/.gemini/config/skills"
  "${HOME}/.agents/skills"
)

if [ "$#" -gt 0 ]; then
  targets=("$@")
else
  targets=("${default_targets[@]}")
fi

sync_skill() {
  local src_dir="$1"
  local dst_root="$2"
  local skill_name="$3"
  local rewrite_name="${4:-}"
  local dst_dir="${dst_root}/${skill_name}"

  mkdir -p "$dst_dir"
  rm -rf "${dst_dir:?}/"*
  cp -R "${src_dir}/." "${dst_dir}/"

  if [ -n "$rewrite_name" ] && [ -f "${dst_dir}/SKILL.md" ]; then
    sed -i "s/^name: .*/name: ${rewrite_name}/" "${dst_dir}/SKILL.md"
  fi
  printf 'installed %-16s -> %s\n' "$skill_name" "$dst_dir"
}

for target in "${targets[@]}"; do
  mkdir -p "$target"
  sync_skill "${repo_root}/skills/specialist" "$target" "rgb-specialist"
  sync_skill "${repo_root}/skills/maker" "$target" "rgb-maker"
  sync_skill "${repo_root}/skills/maker" "$target" "rpg-maker" "rpg-maker"
done
