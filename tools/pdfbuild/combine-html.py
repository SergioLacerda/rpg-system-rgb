#!/usr/bin/env python3
"""Combine MkDocs pages into one WeasyPrint input for manual PDF authoring."""

from __future__ import annotations

import re
import sys
from dataclasses import dataclass
from html import escape
from pathlib import Path

import yaml


ARTICLE_RE = re.compile(
    r"<article[^>]*class=\"[^\"]*md-content__inner[^\"]*\"[^>]*>(?P<body>.*?)</article>",
    re.IGNORECASE | re.DOTALL,
)
BODY_RE = re.compile(r"<body[^>]*>(?P<body>.*)</body>", re.IGNORECASE | re.DOTALL)
TITLE_RE = re.compile(r"<title>(?P<title>.*?)</title>", re.IGNORECASE | re.DOTALL)
ID_RE = re.compile(r'id="([^"]+)"')
HREF_RE = re.compile(r'href="#([^"]+)"')


@dataclass(frozen=True)
class BookPage:
    path: Path
    title: str
    group: str
    group_index: int
    level: int


def prefix_anchors(html: str, anchor_prefix: str) -> str:
    html = ID_RE.sub(lambda match: f'id="{anchor_prefix}-{match.group(1)}"', html)
    return HREF_RE.sub(lambda match: f'href="#{anchor_prefix}-{match.group(1)}"', html)


def extract_body(path: Path, anchor_prefix: str = "") -> str:
    text = path.read_text(encoding="utf-8")
    article = ARTICLE_RE.search(text)
    if article:
        body = article.group("body")
        return prefix_anchors(body, anchor_prefix) if anchor_prefix else body
    match = BODY_RE.search(text)
    return match.group("body") if match else text


def page_title(path: Path) -> str:
    text = path.read_text(encoding="utf-8")
    match = TITLE_RE.search(text)
    if not match:
        return path.stem
    title = " ".join(match.group("title").split())
    return re.sub(r"\s+-\s+RGB System$", "", title)


def html_path_for_markdown(site_dir: Path, markdown_path: str) -> Path:
    path = Path(markdown_path)
    rel = path.with_name("index.html") if path.name == "README.md" else path.with_suffix(".html")
    return site_dir / rel


def collect_nav_pages(
    items: list[object],
    site_dir: Path,
    group: str,
    level: int = 1,
) -> list[BookPage]:
    pages: list[BookPage] = []
    for item in items:
        if isinstance(item, str):
            path = html_path_for_markdown(site_dir, item)
            if path.exists():
                pages.append(BookPage(path, page_title(path), group, 0, level))
        elif isinstance(item, dict):
            for label, value in item.items():
                if isinstance(value, str):
                    path = html_path_for_markdown(site_dir, value)
                    if path.exists():
                        pages.append(BookPage(path, label, group, 0, level))
                elif isinstance(value, list):
                    pages.extend(collect_nav_pages(value, site_dir, label, level + 1))
    return pages


def book_pages(site_dir: Path, config_path: Path) -> list[BookPage]:
    config = yaml.safe_load(config_path.read_text(encoding="utf-8")) or {}
    nav_pages: list[BookPage] = []
    for item in config.get("nav", []):
        if isinstance(item, str):
            path = html_path_for_markdown(site_dir, item)
            if path.exists():
                nav_pages.append(BookPage(path, page_title(path), "Foundations", 0, 2))
        elif isinstance(item, dict):
            for label, value in item.items():
                if isinstance(value, str):
                    path = html_path_for_markdown(site_dir, value)
                    if path.exists():
                        group = "Foundations" if value == "README.md" or value.startswith("introduction/") else label
                        nav_pages.append(BookPage(path, label, group, 0, 2))
                elif isinstance(value, list):
                    nav_pages.extend(collect_nav_pages(value, site_dir, label, 2))
    nav_paths = [page.path for page in nav_pages]
    fallback_pages = [
        path
        for path in sorted(site_dir.rglob("*.html"))
        if path.name != "404.html" and "search" not in path.parts and path not in nav_paths
    ]
    appendices = [
        BookPage(path, page_title(path), "Appendices", len(nav_pages) + 1, 2)
        for path in fallback_pages
    ]
    return nav_pages + appendices


def page_anchor(site_dir: Path, page: Path) -> str:
    rel = page.relative_to(site_dir).with_suffix("")
    return "page-" + "-".join(part for part in rel.parts if part)


def group_anchor(group: str) -> str:
    normalized = re.sub(r"[^a-z0-9]+", "-", group.lower()).strip("-")
    return f"chapter-{normalized or 'section'}"


def group_vector(index: int) -> str:
    return ("b", "r", "g")[index % 3]


def display_label(lang: str, label: str) -> str:
    if not lang.lower().startswith("pt"):
        return label
    labels = {
        "Appendices": "Apêndices",
        "Armor": "Armaduras",
        "Attack and Defense": "Ataque e defesa",
        "Attributes": "Atributos",
        "Character Creation": "Criação de personagem",
        "Character Sheet": "Ficha de personagem",
        "Combat": "Combate",
        "Combat Decision Model": "Modelo de decisão de combate",
        "Combat Example": "Exemplo de combate",
        "Combat Walkthrough": "Passo a passo de combate",
        "Core": "Sistema base",
        "Damage Interaction Model": "Modelo de interação de dano",
        "Damage Model": "Modelo de dano",
        "Equipment": "Equipamentos",
        "Extra": "Extra",
        "Foundations": "Fundamentos",
        "Game Play Example": "Exemplo de jogo",
        "Gameplay Example": "Exemplo de jogo",
        "Gameplay Loop": "Loop de jogo",
        "Gear": "Equipamentos gerais",
        "Glossary": "Glossário",
        "Movement": "Movimento",
        "One Page Rules": "Regras em uma página",
        "Overview": "Visão geral",
        "Progression": "Progressão",
        "Quick Start": "Início rápido",
        "Reference": "Referência",
        "Skills and Abilities": "Habilidades",
        "System Engine": "Engine do sistema",
        "System Overview": "Visão geral do sistema",
        "Weapons": "Armas",
    }
    return labels.get(label, label)


def chapter_label(lang: str, index: int) -> str:
    noun = "Capítulo" if lang.lower().startswith("pt") else "Chapter"
    return f"{noun} {index:02d}"


def in_section_label(lang: str) -> str:
    return "Nesta seção" if lang.lower().startswith("pt") else "In this section"


def main() -> int:
    if len(sys.argv) != 6:
        print(
            "usage: combine-html.py <lang> <version> <cover-html> <site-dir> <mkdocs-yml>",
            file=sys.stderr,
        )
        return 2

    lang = sys.argv[1]
    version = sys.argv[2]
    cover = Path(sys.argv[3])
    site = Path(sys.argv[4])
    config = Path(sys.argv[5])
    pages = book_pages(site, config)

    print("<!doctype html>")
    print(f'<html lang="{lang}">')
    print("<head>")
    print('  <meta charset="utf-8">')
    print("  <title>RGB System Core Rules</title>")
    print("  <meta name=\"author\" content=\"RGB System contributors\">")
    print(f"  <meta name=\"subject\" content=\"RGB System Core Rules - {version}\">")
    print(f"  <meta name=\"description\" content=\"RGB System Core Rules - {version}\">")
    print("  <meta name=\"generator\" content=\"MkDocs + WeasyPrint\">")
    print("</head>")
    print(f'<body lang="{lang}" data-version="{version}">')
    print(extract_body(cover))
    print('<nav class="toc">')
    print('  <div class="toc__head">')
    print("    <h2>Contents</h2>")
    print(f'    <span class="toc__version">{version}</span>')
    print("  </div>")
    print("  <ol>")
    current_group = ""
    current_group_index = 0
    for page in pages:
        if page.group != current_group:
            current_group = page.group
            current_group_index += 1
            vector = group_vector(current_group_index)
            group_label = display_label(lang, page.group)
            print(f'    <li class="toc__group vector-{vector}"><a href="#{group_anchor(page.group)}">{escape(group_label)}</a></li>')
        anchor = page_anchor(site, page.path)
        vector = group_vector(current_group_index)
        print(f'    <li class="lvl-{page.level} vector-{vector}"><a href="#{anchor}">{escape(display_label(lang, page.title))}</a></li>')
    print("  </ol>")
    print("</nav>")

    current_group = ""
    current_group_index = 0
    for page in pages:
        if page.group != current_group:
            current_group = page.group
            current_group_index += 1
            vector = group_vector(current_group_index)
            group_label = display_label(lang, page.group)
            print(f'<section id="{group_anchor(page.group)}" class="chapter-opener chapter-opener--{vector}">')
            print(f'  <div class="chapter-opener__rule"></div>')
            print(f'  <p class="chapter-opener__label">{chapter_label(lang, current_group_index)}</p>')
            print(f'  <h1>{escape(group_label)}</h1>')
            print('  <div class="chapter-opener__contents">')
            print(f"    <p>{in_section_label(lang)}</p>")
            print("    <ol>")
            for child in [candidate for candidate in pages if candidate.group == page.group][:6]:
                child_anchor = page_anchor(site, child.path)
                child_label = display_label(lang, child.title)
                print(f'      <li><a href="#{child_anchor}">{escape(child_label)}</a></li>')
            print("    </ol>")
            print("  </div>")
            print("</section>")

        anchor = page_anchor(site, page.path)
        print(f'<section id="{anchor}" class="pdf-page-source">')
        print(extract_body(page.path, anchor))
        print("</section>")
    print("</body>")
    print("</html>")
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
