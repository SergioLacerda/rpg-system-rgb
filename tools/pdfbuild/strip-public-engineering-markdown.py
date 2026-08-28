#!/usr/bin/env python3
"""Remove engineering-only ADR references from public PDF source copies."""

from __future__ import annotations

import re
import sys
from pathlib import Path


ADR_LINK_RE = re.compile(r"\[([^\]]+)\]\((?:\.\./\.\./adr/|\.\./adr/|docs/adr/)[^)]+\)")
ENGINEERING_SECTION_RE = re.compile(r"^(Architecture Decisions|Decisões de Arquitetura)$", re.IGNORECASE)


def contains_adr_reference(text: str) -> bool:
    lowered = text.lower()
    return any(marker in lowered for marker in ("docs/adr/", "../adr/", "../../adr/", "/adr/", "adr-"))


def strip_public_engineering(content: str) -> str:
    lines: list[str] = []
    skip_section = False
    for raw in content.splitlines():
        line = raw.strip()
        if line.startswith("## "):
            skip_section = bool(ENGINEERING_SECTION_RE.match(line.removeprefix("## ").strip()))
        elif line.startswith("#"):
            skip_section = False
        if skip_section:
            continue
        raw = ADR_LINK_RE.sub(r"\1", raw)
        if contains_adr_reference(raw):
            continue
        lines.append(raw)
    return "\n".join(lines).rstrip() + "\n"


def main() -> int:
    if len(sys.argv) != 2:
        print("usage: strip-public-engineering-markdown.py <docs-dir>", file=sys.stderr)
        return 2
    root = Path(sys.argv[1])
    for path in sorted(root.rglob("*.md")):
        path.write_text(strip_public_engineering(path.read_text(encoding="utf-8")), encoding="utf-8")
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
