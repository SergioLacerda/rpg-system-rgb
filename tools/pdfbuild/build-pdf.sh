#!/usr/bin/env bash
set -euo pipefail

usage() {
  echo "usage: tools/pdfbuild/build-pdf.sh <en|pt-br> <version>" >&2
}

lang_code="${1:-}"
version="${2:-}"
if [ -z "${lang_code}" ] || [ -z "${version}" ]; then
  usage
  exit 2
fi

case "${lang_code}" in
  en)
    source_lang="en"
    html_lang="en"
    lang_name="English"
    subtitle="Core Rules"
    kicker="Tabletop Roleplaying System"
    colophon_title="Colophon"
    ;;
  pt-br)
    source_lang="PT-br"
    html_lang="pt-BR"
    lang_name="Portuguese (Brazil)"
    subtitle="Regras Centrais"
    kicker="Sistema de RPG de Mesa"
    colophon_title="Ficha tecnica"
    ;;
  *)
    usage
    echo "unsupported language: ${lang_code}" >&2
    exit 2
    ;;
esac

script_dir="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
repo_root="$(cd "${script_dir}/../.." && pwd)"
work_dir="${repo_root}/.pdfbuild/${lang_code}/${version}"
out_dir="${script_dir}/out"
venv_dir="${script_dir}/.venv"
stamp="$(date -u +%Y-%m-%d)"
commit_sha="$(git -C "${repo_root}" rev-parse --short HEAD)"

mkdir -p "${work_dir}/docs" "${work_dir}/styles" "${work_dir}/fonts" "${out_dir}"
rm -rf "${work_dir}/docs" "${work_dir}/site"
mkdir -p "${work_dir}/docs" "${work_dir}/styles"

cp -R "${repo_root}/docs/core/${source_lang}/." "${work_dir}/docs/"
python3 "${script_dir}/strip-public-engineering-markdown.py" "${work_dir}/docs"
cp "${repo_root}/docs/styles/rgb-pdf.css" "${work_dir}/styles/rgb-pdf.css"
rm -rf "${work_dir}/fonts"
cp -R "${script_dir}/fonts" "${work_dir}/fonts"
cp "${script_dir}/mkdocs.pdf.yml" "${work_dir}/mkdocs.yml"

python3 -m venv "${venv_dir}"
# shellcheck disable=SC1091
source "${venv_dir}/bin/activate"
python -m pip install --quiet --require-hashes -r "${script_dir}/requirements.txt"

mkdocs build --strict --config-file "${work_dir}/mkdocs.yml"

cover_html="${work_dir}/cover.html"
sed \
  -e "s|{{ lang_code }}|${html_lang}|g" \
  -e "s|{{ lang_name }}|${lang_name}|g" \
  -e "s|{{ title }}|RGB System|g" \
  -e "s|{{ subtitle }}|${subtitle}|g" \
  -e "s|{{ author }}|RGB System contributors|g" \
  -e "s|{{ version }}|${version}|g" \
  -e "s|{{ date }}|${stamp}|g" \
  -e "s|{{ license }}|MIT|g" \
  -e "s|{{ commit_sha }}|${commit_sha}|g" \
  -e "s|{{ kicker }}|${kicker}|g" \
  -e "s|{{ colophon_title }}|${colophon_title}|g" \
  -e "s|{{ repo_url }}|https://github.com/SergioLacerda/rpg-system-rgb|g" \
  -e "s|{{ generator }}|MkDocs + WeasyPrint|g" \
  "${script_dir}/cover.html" > "${cover_html}"

book_html="${work_dir}/book.html"
python3 "${script_dir}/combine-html.py" "${html_lang}" "${version}" "${cover_html}" "${work_dir}/site" "${work_dir}/mkdocs.yml" > "${book_html}"

pdf="${out_dir}/rgb-system-core-${version}-${lang_code}.pdf"
weasyprint \
  --encoding utf-8 \
  --stylesheet "${work_dir}/styles/rgb-pdf.css" \
  --pdf-identifier "rgb-system-core-${version}-${lang_code}" \
  --pdf-variant pdf/ua-1 \
  --custom-metadata \
  "${book_html}" \
  "${pdf}"

pdfinfo "${pdf}"
pdftotext -f 1 -l 3 -layout "${pdf}" - | head -40
echo "generated: ${pdf}"
