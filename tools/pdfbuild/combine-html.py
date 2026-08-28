#!/usr/bin/env python3
"""Combine MkDocs pages into one WeasyPrint input for manual PDF authoring."""

from __future__ import annotations

import re
import sys
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


def nav_markdown_paths(items: list[object]) -> list[str]:
    paths: list[str] = []
    for item in items:
        if isinstance(item, str):
            paths.append(item)
        elif isinstance(item, dict):
            for value in item.values():
                if isinstance(value, str):
                    paths.append(value)
                elif isinstance(value, list):
                    paths.extend(nav_markdown_paths(value))
    return paths


def page_paths(site_dir: Path, config_path: Path) -> list[Path]:
    config = yaml.safe_load(config_path.read_text(encoding="utf-8")) or {}
    nav_pages = [
        html_path_for_markdown(site_dir, path)
        for path in nav_markdown_paths(config.get("nav", []))
    ]
    nav_pages = [path for path in nav_pages if path.exists()]
    fallback_pages = [
        path
        for path in sorted(site_dir.rglob("*.html"))
        if path.name != "404.html" and "search" not in path.parts and path not in nav_pages
    ]
    return nav_pages + fallback_pages


def page_anchor(site_dir: Path, page: Path) -> str:
    rel = page.relative_to(site_dir).with_suffix("")
    return "page-" + "-".join(part for part in rel.parts if part)


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
    pages = page_paths(site, config)

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
    vector_classes = ("b", "r", "g", "b")
    for index, page in enumerate(pages):
        anchor = page_anchor(site, page)
        vector = vector_classes[index % len(vector_classes)]
        print(f'    <li class="lvl-1 vector-{vector}"><a href="#{anchor}">{page_title(page)}</a></li>')
    print("  </ol>")
    print("</nav>")
    for page in pages:
        anchor = page_anchor(site, page)
        print(f'<section id="{anchor}" class="pdf-page-source">')
        print(extract_body(page, anchor))
        print("</section>")
    print("</body>")
    print("</html>")
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
