package compiler

import (
	"fmt"
	"html"
	"strings"
)

const pageTemplate = `<!doctype html>
<html lang="%s">
<head>
<meta charset="utf-8">
<title>%s</title>
<meta name="viewport" content="width=device-width, initial-scale=1">
<style>
body { font-family: system-ui, sans-serif; max-width: 48rem; margin: 2rem auto; padding: 0 1rem; line-height: 1.5; }
pre { background: #f4f4f4; padding: 0.75rem; overflow-x: auto; }
table { border-collapse: collapse; }
th, td { border: 1px solid #ccc; padding: 0.4rem 0.6rem; text-align: left; }
</style>
</head>
<body>
%s</body>
</html>
`

// RenderHTML renders blocks into a complete, self-contained HTML document.
// Title is HTML-escaped; block content is escaped per span.
func RenderHTML(title string, blocks []Block) string {
	return renderHTMLPage("en", title, blocks)
}

// renderHTMLPage renders blocks into a complete HTML document tagged with
// the given lang attribute.
func renderHTMLPage(lang, title string, blocks []Block) string {
	var body strings.Builder
	for _, block := range blocks {
		writeBlockHTML(&body, block)
	}
	return fmt.Sprintf(pageTemplate, html.EscapeString(lang), html.EscapeString(title), body.String())
}

func writeBlockHTML(w *strings.Builder, block Block) {
	switch b := block.(type) {
	case Heading:
		fmt.Fprintf(w, "<h%d>%s</h%d>\n", b.Level, spansToHTML(b.Inline), b.Level)
	case Paragraph:
		fmt.Fprintf(w, "<p>%s</p>\n", spansToHTML(b.Inline))
	case CodeFence:
		writeCodeFenceHTML(w, b)
	case List:
		writeListHTML(w, b)
	case Table:
		writeTableHTML(w, b)
	case HorizontalRule:
		w.WriteString("<hr>\n")
	}
}

func writeCodeFenceHTML(w *strings.Builder, fence CodeFence) {
	class := ""
	if fence.Lang != "" {
		class = fmt.Sprintf(` class="language-%s"`, html.EscapeString(fence.Lang))
	}
	fmt.Fprintf(w, "<pre><code%s>%s</code></pre>\n", class, html.EscapeString(strings.Join(fence.Lines, "\n")))
}

func writeListHTML(w *strings.Builder, list List) {
	tag := "ul"
	if list.Ordered {
		tag = "ol"
	}
	fmt.Fprintf(w, "<%s>\n", tag)
	for _, item := range list.Items {
		fmt.Fprintf(w, "<li>%s</li>\n", spansToHTML(item))
	}
	fmt.Fprintf(w, "</%s>\n", tag)
}

func writeTableHTML(w *strings.Builder, table Table) {
	w.WriteString("<table>\n<thead>\n<tr>")
	for _, cell := range table.Header {
		fmt.Fprintf(w, "<th>%s</th>", html.EscapeString(cell))
	}
	w.WriteString("</tr>\n</thead>\n<tbody>\n")
	for _, row := range table.Rows {
		w.WriteString("<tr>")
		for _, cell := range row {
			fmt.Fprintf(w, "<td>%s</td>", html.EscapeString(cell))
		}
		w.WriteString("</tr>\n")
	}
	w.WriteString("</tbody>\n</table>\n")
}

// spansToHTML renders inline spans to an escaped HTML fragment.
func spansToHTML(spans []InlineSpan) string {
	var w strings.Builder
	for _, span := range spans {
		text := html.EscapeString(span.Text)
		switch {
		case span.LinkURL != "":
			fmt.Fprintf(&w, `<a href="%s">%s</a>`, html.EscapeString(rewriteMarkdownLink(span.LinkURL)), text)
		case span.Bold:
			fmt.Fprintf(&w, "<strong>%s</strong>", text)
		case span.Italic:
			fmt.Fprintf(&w, "<em>%s</em>", text)
		case span.Code:
			fmt.Fprintf(&w, "<code>%s</code>", text)
		default:
			w.WriteString(text)
		}
	}
	return w.String()
}

// rewriteMarkdownLink rewrites a relative .md cross-reference into the
// .html sibling that RenderHTMLTree produces for it, so links between
// rendered pages stay valid on the static site. External links (containing
// "://") and in-page anchors (starting with "#") are left unchanged.
func rewriteMarkdownLink(url string) string {
	if strings.Contains(url, "://") || strings.HasPrefix(url, "#") {
		return url
	}
	if strings.HasSuffix(url, ".md") {
		return strings.TrimSuffix(url, ".md") + ".html"
	}
	return url
}
