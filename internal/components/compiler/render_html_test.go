package compiler

import (
	"strings"
	"testing"
)

func TestRenderHTMLHeadingAndParagraph(t *testing.T) {
	blocks := Parse("# Title\n\nA paragraph with **bold** text.\n")
	html := RenderHTML("Title", blocks)
	for _, want := range []string{"<h1>Title</h1>", "<strong>bold</strong>", "<p>A paragraph with"} {
		if !strings.Contains(html, want) {
			t.Fatalf("output missing %q:\n%s", want, html)
		}
	}
}

func TestRenderHTMLEscapesUnsafeText(t *testing.T) {
	blocks := Parse("Contains <script>alert(1)</script> raw text.\n")
	html := RenderHTML("Title", blocks)
	if strings.Contains(html, "<script>") {
		t.Fatalf("raw HTML was not escaped:\n%s", html)
	}
	if !strings.Contains(html, "&lt;script&gt;") {
		t.Fatalf("expected escaped script tag:\n%s", html)
	}
}

func TestRenderHTMLList(t *testing.T) {
	blocks := Parse("- first\n- second\n")
	html := RenderHTML("Title", blocks)
	for _, want := range []string{"<ul>", "<li>first</li>", "<li>second</li>", "</ul>"} {
		if !strings.Contains(html, want) {
			t.Fatalf("output missing %q:\n%s", want, html)
		}
	}
}

func TestRenderHTMLTable(t *testing.T) {
	source := "| A | B |\n| --- | --- |\n| 1 | 2 |\n"
	blocks := Parse(source)
	html := RenderHTML("Title", blocks)
	for _, want := range []string{"<table>", "<th>A</th>", "<td>1</td>", "</table>"} {
		if !strings.Contains(html, want) {
			t.Fatalf("output missing %q:\n%s", want, html)
		}
	}
}

func TestRenderHTMLCodeFence(t *testing.T) {
	blocks := Parse("```text\nR -> attack\n```\n")
	html := RenderHTML("Title", blocks)
	if !strings.Contains(html, "<pre><code") || !strings.Contains(html, "R -&gt; attack") {
		t.Fatalf("code fence not rendered as expected:\n%s", html)
	}
}

func TestRenderHTMLLinkRewritesMarkdownToHTML(t *testing.T) {
	blocks := Parse("See [Armor](armor.md) and [Shields](../equipment/shields.md) for details.\n")
	html := RenderHTML("Title", blocks)
	if !strings.Contains(html, `<a href="armor.html">Armor</a>`) {
		t.Fatalf("relative .md link was not rewritten to .html:\n%s", html)
	}
	if !strings.Contains(html, `<a href="../equipment/shields.html">Shields</a>`) {
		t.Fatalf("relative .md link with path was not rewritten to .html:\n%s", html)
	}
}

func TestRenderHTMLLinkLeavesExternalAndAnchorLinksUnchanged(t *testing.T) {
	blocks := Parse("See [external](https://example.com/page.md) and [section](#health).\n")
	html := RenderHTML("Title", blocks)
	if !strings.Contains(html, `<a href="https://example.com/page.md">external</a>`) {
		t.Fatalf("external link was rewritten unexpectedly:\n%s", html)
	}
	if !strings.Contains(html, `<a href="#health">section</a>`) {
		t.Fatalf("anchor link was rewritten unexpectedly:\n%s", html)
	}
}
