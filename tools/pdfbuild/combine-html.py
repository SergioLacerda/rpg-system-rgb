#!/usr/bin/env python3
"""Combine MkDocs pages into one WeasyPrint input for manual PDF authoring."""

from __future__ import annotations

import re
import sys
from dataclasses import dataclass
from html import escape, unescape
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
HEADING_RE = re.compile(r"<h(?P<level>[23])[^>]*>(?P<body>.*?)</h[23]>", re.IGNORECASE | re.DOTALL)
TEXT_CODE_RE = re.compile(
    r'<pre><code class="language-text">(?P<body>.*?)</code></pre>',
    re.IGNORECASE | re.DOTALL,
)
LIST_ITEM_RE = re.compile(r"<li\b[^>]*>.*?</li>", re.IGNORECASE | re.DOTALL)
PARAGRAPH_RE = re.compile(r"<p>(?P<body>.*?)</p>", re.IGNORECASE | re.DOTALL)
EXAMPLE_RE = re.compile(
    r"<p>(?P<label>Example:|Examples:|Exemplo:|Exemplos:)</p>\s*(?P<body><(?:ul|ol|pre|table)\b.*?</(?:ul|ol|pre|table)>)",
    re.IGNORECASE | re.DOTALL,
)
SEE_ALSO_RE = re.compile(
    r"<p>(?P<label>See also:|Veja:|Ver também:)</p>\s*(?P<body><ul\b.*?</ul>)",
    re.IGNORECASE | re.DOTALL,
)
OPTIONAL_RULE_RE = re.compile(
    r"(?P<heading><h[23][^>]*>(?:Optional Rule|Regra Opcional):.*?</h[23]>)\s*<p>(?P<body>.*?)</p>",
    re.IGNORECASE | re.DOTALL,
)


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


def contains_adr_reference(text: str) -> bool:
    lowered = text.lower()
    return any(marker in lowered for marker in ("docs/adr/", "../adr/", "../../adr/", "/adr/", "adr-"))


def strip_public_engineering_html(html: str) -> str:
    html = re.sub(
        r"<h2[^>]*>(?:Architecture Decisions|Decisões de Arquitetura)</h2>.*?(?=<h2|\Z)",
        "",
        html,
        flags=re.IGNORECASE | re.DOTALL,
    )
    html = LIST_ITEM_RE.sub(lambda match: "" if contains_adr_reference(match.group(0)) else match.group(0), html)
    html = re.sub(r'<a\s+href="[^"]*(?:docs/adr/|\.\./adr/|\.\./\.\./adr/|adr-)[^"]*">(.*?)</a>', r"\1", html, flags=re.IGNORECASE | re.DOTALL)
    return html


def vector_label(lang: str, level: str, vector: str) -> str:
    if level == "rule":
        noun = "Regra" if lang.lower().startswith("pt") else "Rule"
        return f"{vector.upper()} {noun}"
    if level == "exception":
        noun = "Exceção" if lang.lower().startswith("pt") else "Exception"
        return f"{vector.upper()}! {noun}"
    if level == "example":
        return "Na mesa" if lang.lower().startswith("pt") else "At the table"
    return vector.upper()


def editorial_box(level: str, vector: str, label: str, body: str) -> str:
    return (
        f'<div class="vector-box vector-box--{level} vector-box--{vector}">'
        f'<p class="vector-box__label">{escape(label)}</p>{body}</div>'
    )


def pipe_table_cells(line: str) -> list[str]:
    return [cell.strip() for cell in line.strip().strip("|").split("|")]


def is_pipe_table_separator(line: str) -> bool:
    normalized = line.replace(" ", "")
    if not normalized or set(normalized) <= {"|", "-", ":"}:
        return True
    cells = pipe_table_cells(line)
    if not cells:
        return False
    return all(
        cell == "" or re.fullmatch(r":?-{3,}:?", cell.replace(" ", ""))
        for cell in cells
    )


def parse_pipe_table_line(line: str) -> list[str] | None:
    if "|" not in line:
        return None
    cells = pipe_table_cells(line)
    cells = [cell for cell in cells if cell != ""]
    return cells if len(cells) >= 2 else None


def is_spaced_table_separator(line: str) -> bool:
    return bool(re.fullmatch(r"[-\s]{6,}", line)) and "---" in line


def parse_spaced_table_line(line: str) -> list[str] | None:
    cells = [cell.strip() for cell in re.split(r"\s{2,}", line.strip())]
    cells = [cell for cell in cells if cell != ""]
    return cells if len(cells) >= 2 else None


def normalize_text_table_lines(code: str) -> list[str]:
    lines: list[str] = []
    for line in code.splitlines():
        line = unescape(line).strip()
        if not line:
            continue
        glued = re.search(r"-{2,}:?(?P<row>[^\s-].*)$", line)
        if glued:
            lines.append(line[: glued.start("row")].strip())
            lines.append(glued.group("row").strip())
            continue
        lines.append(line)
    return lines


def text_table_html(code: str) -> str | None:
    lines = normalize_text_table_lines(code)
    if len(lines) < 2:
        return None

    rows: list[list[str]] = []
    if "|" in lines[0]:
        for line in lines:
            if is_pipe_table_separator(line):
                continue
            row = parse_pipe_table_line(line)
            if row is None:
                return None
            rows.append(row)
    elif len(lines) >= 3 and is_spaced_table_separator(lines[1]):
        rows.append(parse_spaced_table_line(lines[0]) or [])
        for line in lines[2:]:
            row = parse_spaced_table_line(line)
            if row is None:
                return None
            rows.append(row)
    else:
        return None

    if len(rows) < 2:
        return None
    column_count = len(rows[0])
    if column_count < 2 or any(len(row) != column_count for row in rows):
        return None

    numeric_cell_re = re.compile(r"[−+-]?\d+(\.\d+)?(\s?[A-Za-z])?$")
    body_rows = rows[1:]
    numeric_columns = {
        col
        for col in range(1, column_count)
        if all(numeric_cell_re.fullmatch(row[col]) for row in body_rows)
    }

    def cell(tag: str, col: int, value: str) -> str:
        attr = ' style="text-align: right"' if col in numeric_columns else ""
        return f"<{tag}{attr}>{escape(value)}</{tag}>"

    head = "".join(cell("th", col, value) for col, value in enumerate(rows[0]))
    body = "\n".join(
        "    <tr>"
        + "".join(cell("td", col, value) for col, value in enumerate(row))
        + "</tr>"
        for row in body_rows
    )
    return f'<table class="value-table value-table--text"><thead><tr>{head}</tr></thead><tbody>\n{body}\n  </tbody></table>'


def convert_text_tables(html: str) -> str:
    def replacement(match: re.Match[str]) -> str:
        table = text_table_html(match.group("body"))
        return table if table else match.group(0)

    return TEXT_CODE_RE.sub(replacement, html)


def annotate_editorial_blocks(html: str, lang: str, vector: str) -> str:
    html = strip_public_engineering_html(html)
    html = convert_text_tables(html)

    def optional_rule(match: re.Match[str]) -> str:
        body = match.group("body")
        box = editorial_box("exception", vector, vector_label(lang, "exception", vector), f"<p>{body}</p>")
        return f"{match.group('heading')}\n{box}"

    html = OPTIONAL_RULE_RE.sub(optional_rule, html)

    def example(match: re.Match[str]) -> str:
        return editorial_box("example", "g", vector_label(lang, "example", "g"), match.group("body"))

    html = EXAMPLE_RE.sub(example, html)

    def see_also(match: re.Match[str]) -> str:
        return editorial_box("marginal", "g", vector_label(lang, "marginal", "g"), match.group("body"))

    html = SEE_ALSO_RE.sub(see_also, html)

    def paragraph(match: re.Match[str]) -> str:
        body = match.group("body").strip()
        normalized = re.sub(r"<.*?>", "", body).strip()
        if contains_adr_reference(body):
            return ""
        if re.match(r"^(Do not|Não|Nunca|Never|Examples cannot|Exemplos não)\b", normalized, re.IGNORECASE):
            return editorial_box("rule", vector, vector_label(lang, "rule", vector), f"<p>{body}</p>")
        return match.group(0)

    return PARAGRAPH_RE.sub(paragraph, html)


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


def chapter_epigraph(lang: str, index: int) -> str:
    pt = lang.lower().startswith("pt")
    epigraphs_en = {
        1: "Three numbers. Every outcome the table will ever need.",
        2: "Red presses. Green answers. Blue endures.",
        3: "A fight is a question asked in pressure, answered in cost.",
        4: "Nothing here changes who you are. Only what you can do next.",
        5: "A weapon is Red given a shape.",
        6: "For the table, mid-session, under pressure.",
    }
    epigraphs_pt = {
        1: "Três números. Tudo que a mesa vai precisar.",
        2: "Vermelho pressiona. Verde responde. Azul resiste.",
        3: "Uma luta é uma pergunta feita em pressão, respondida em custo.",
        4: "Nada aqui muda quem você é. Só o que você pode fazer a seguir.",
        5: "Uma arma é o Vermelho com uma forma.",
        6: "Para a mesa, no meio da sessão, sob pressão.",
    }
    table = epigraphs_pt if pt else epigraphs_en
    return table.get(index, "")


def reference_sheet_label(lang: str) -> str:
    return "Folhas de referência" if lang.lower().startswith("pt") else "Reference Sheets"


def index_label(lang: str) -> str:
    return "Índice remissivo" if lang.lower().startswith("pt") else "Index"


def strip_tags(html: str) -> str:
    text = re.sub(r"<[^>]+>", "", html)
    return " ".join(unescape(text).split())


def collect_index_entries(lang: str, site: Path, pages: list[BookPage]) -> list[tuple[str, str]]:
    entries: dict[str, tuple[str, str]] = {}
    skip_titles = {"overview", "visão geral", "appendices", "apêndices"}
    for page in pages:
        label = display_label(lang, page.title)
        if label.strip().lower() not in skip_titles:
            entries.setdefault(label.casefold(), (label, page_anchor(site, page.path)))
        if page.path.stem != "glossary":
            continue
        for match in HEADING_RE.finditer(extract_body(page.path)):
            term = strip_tags(match.group("body"))
            if term:
                entries.setdefault(term.casefold(), (term, page_anchor(site, page.path)))
    return sorted(entries.values(), key=lambda item: item[0].casefold())


def emit_reference_sheets(lang: str, version: str) -> None:
    pt = lang.lower().startswith("pt")
    title = reference_sheet_label(lang)
    vectors = [
        ("R", "Pressão" if pt else "Pressure", "impacto, força, interrupção" if pt else "impact, force, interruption"),
        ("G", "Relação" if pt else "Relation", "movimento, timing, esquiva" if pt else "movement, timing, evasion"),
        ("B", "Preservação" if pt else "Preservation", "escudo, bloqueio, resistência" if pt else "shield, block, resistance"),
    ]
    damage_steps = [
        "Verificação de acerto ou contato" if pt else "Hit or contact check",
        "Fonte de impacto" if pt else "Impact source",
        "Penetração" if pt else "Penetration",
        "Redução por armadura" if pt else "Armor reduction",
        "Absorção por escudo" if pt else "Shield absorption",
        "Dano restante ou consequência" if pt else "Remaining damage or consequence",
    ]
    defense_rows = [
        ("Esquiva" if pt else "Evade", "G", "evitar ou alterar contato" if pt else "avoid or alter contact"),
        ("Reposicionar" if pt else "Reposition", "G", "mudar alcance, cobertura ou engajamento" if pt else "change range, cover, or engagement"),
        ("Bloquear" if pt else "Block", "B", "receber pressão de forma intencional" if pt else "intentionally receive pressure"),
        ("Interromper" if pt else "Interrupt", "R", "parar ação com pressão primeiro" if pt else "stop an action by applying pressure first"),
    ]

    print('<section id="reference-sheets" class="reference-sheets">')
    print(f"  <h1>{escape(title)}</h1>")
    print(f'  <p class="reference-sheets__version">{escape(version)}</p>')
    print(f"  <h2>{escape('Vetores' if pt else 'Vectors')}</h2>")
    print('  <table class="value-table value-table--vectors"><thead><tr>')
    print(f"    <th>{escape('Vetor' if pt else 'Vector')}</th><th>{escape('Nome' if pt else 'Name')}</th><th>{escape('Uso em mesa' if pt else 'At the table')}</th>")
    print("  </tr></thead><tbody>")
    for vector, name, use in vectors:
        print(f'    <tr><td class="value-table__key">{vector}</td><td>{escape(name)}</td><td>{escape(use)}</td></tr>')
    print("  </tbody></table>")

    print(f"  <h2>{escape('Fluxo de dano' if pt else 'Damage Flow')}</h2>")
    print('  <ol class="reference-flow">')
    for step in damage_steps:
        print(f"    <li>{escape(step)}</li>")
    print("  </ol>")

    print(f"  <h2>{escape('Procedimentos defensivos' if pt else 'Defensive Procedures')}</h2>")
    print('  <table class="value-table"><thead><tr>')
    print(f"    <th>{escape('Procedimento' if pt else 'Procedure')}</th><th>{escape('Vetor' if pt else 'Vector')}</th><th>{escape('Função' if pt else 'Purpose')}</th>")
    print("  </tr></thead><tbody>")
    for procedure, vector, purpose in defense_rows:
        print(f'    <tr><td>{escape(procedure)}</td><td class="value-table__key">{vector}</td><td>{escape(purpose)}</td></tr>')
    print("  </tbody></table>")

    print(f"  <div class=\"vector-box vector-box--rule vector-box--r\"><p class=\"vector-box__label\">{escape(vector_label(lang, 'rule', 'r'))}</p>")
    print(f"    <p>{escape('Exemplos e atalhos não substituem a regra declarada.' if pt else 'Examples and shortcuts do not replace the declared rule.')}</p></div>")
    print("</section>")


def emit_book_index(lang: str, site: Path, pages: list[BookPage]) -> None:
    print('<section id="book-index" class="book-index">')
    print(f"  <h1>{escape(index_label(lang))}</h1>")
    print("  <ol>")
    for label, anchor in collect_index_entries(lang, site, pages):
        print(f'    <li><a href="#{anchor}">{escape(label)}</a></li>')
    print("  </ol>")
    print("</section>")


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
    print(f"  <meta name=\"version\" content=\"{version}\">")
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
    print(f'    <li class="toc__group vector-r"><a href="#reference-sheets">{escape(reference_sheet_label(lang))}</a></li>')
    print(f'    <li class="toc__group vector-g"><a href="#book-index">{escape(index_label(lang))}</a></li>')
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
            epigraph = chapter_epigraph(lang, current_group_index)
            if epigraph:
                print(f'  <p class="chapter-opener__epigraph">{escape(epigraph)}</p>')
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
        print(annotate_editorial_blocks(extract_body(page.path, anchor), lang, group_vector(page.group_index or current_group_index)))
        print("</section>")
    emit_reference_sheets(lang, version)
    emit_book_index(lang, site, pages)
    print("</body>")
    print("</html>")
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
