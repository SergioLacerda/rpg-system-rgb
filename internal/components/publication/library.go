package publication

import (
	"bytes"
	"fmt"
	"html"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
)

// LibraryOptions configures static HTML Library generation.
type LibraryOptions struct {
	SourceDir string
	OutDir    string
}

type docPage struct {
	Title   string
	RelPath string
	URLPath string
	Lang    string
	Body    string
}

var markdownLinkPattern = regexp.MustCompile(`\[([^\]]+)\]\(([^)]+)\)`)
var inlineCodePattern = regexp.MustCompile("`([^`]+)`")
var boldPattern = regexp.MustCompile(`\*\*([^*]+)\*\*`)

// BuildLibrary renders docs/core/** Markdown into a static HTML Library.
//
//nolint:gocyclo // Coordinates validation, discovery, cleanup, rendering, and index writing.
func BuildLibrary(options LibraryOptions) error {
	if options.SourceDir == "" || options.OutDir == "" {
		return fmt.Errorf("source and out are required")
	}
	pages, err := collectPages(options.SourceDir)
	if err != nil {
		return err
	}
	if len(pages) == 0 {
		return fmt.Errorf("no docs/core pages found under %s", options.SourceDir)
	}
	if err := os.RemoveAll(options.OutDir); err != nil {
		return err
	}
	if err := os.MkdirAll(options.OutDir, 0o750); err != nil {
		return err
	}
	if err := writeLibraryAssets(options.OutDir); err != nil {
		return err
	}
	for _, page := range pages {
		if err := writePage(options.OutDir, page, pages); err != nil {
			return err
		}
	}
	return writeLibraryIndex(options.OutDir, pages)
}

//nolint:gocyclo // WalkDir callback keeps discovery decisions close to path handling.
func collectPages(sourceDir string) ([]docPage, error) {
	var pages []docPage
	coreDir := filepath.Join(sourceDir, "core")
	err := filepath.WalkDir(coreDir, func(path string, entry os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".md") {
			return nil
		}
		rel, err := filepath.Rel(sourceDir, path)
		if err != nil {
			return err
		}
		rel = filepath.ToSlash(rel)
		if !strings.HasPrefix(rel, "core/en/") && !strings.HasPrefix(rel, "core/PT-br/") {
			return nil
		}
		content, err := os.ReadFile(path) //nolint:gosec // G304: source path is controlled by the CLI/Make target.
		if err != nil {
			return err
		}
		page := docPage{
			Title:   titleFromMarkdown(content, rel),
			RelPath: rel,
			URLPath: outputURL(rel),
			Lang:    langFromRel(rel),
			Body:    renderMarkdown(content),
		}
		pages = append(pages, page)
		return nil
	})
	if err != nil {
		return nil, err
	}
	sort.Slice(pages, func(i, j int) bool {
		if pages[i].Lang != pages[j].Lang {
			return pages[i].Lang < pages[j].Lang
		}
		return pages[i].URLPath < pages[j].URLPath
	})
	return pages, nil
}

func titleFromMarkdown(content []byte, fallback string) string {
	for _, line := range strings.Split(string(content), "\n") {
		if strings.HasPrefix(line, "# ") {
			return strings.TrimSpace(strings.TrimPrefix(line, "# "))
		}
	}
	return strings.TrimSuffix(filepath.Base(fallback), filepath.Ext(fallback))
}

func langFromRel(rel string) string {
	if strings.HasPrefix(rel, "core/PT-br/") {
		return "pt-br"
	}
	return "en"
}

func outputURL(rel string) string {
	withoutExt := strings.TrimSuffix(rel, ".md")
	withoutExt = strings.TrimSuffix(withoutExt, "/README")
	return "/" + withoutExt + "/"
}

func outputFile(outDir string, urlPath string) string {
	clean := strings.Trim(urlPath, "/")
	if clean == "" {
		return filepath.Join(outDir, "index.html")
	}
	return filepath.Join(outDir, filepath.FromSlash(clean), "index.html")
}

//nolint:gocyclo // Small, explicit Markdown subset renderer; each branch maps to one supported block.
func renderMarkdown(content []byte) string {
	lines := strings.Split(string(content), "\n")
	var out strings.Builder
	var paragraph []string
	inCode := false

	flushParagraph := func() {
		if len(paragraph) == 0 {
			return
		}
		fmt.Fprintf(&out, "<p>%s</p>\n", renderInline(strings.Join(paragraph, " ")))
		paragraph = nil
	}

	for _, raw := range lines {
		line := strings.TrimRight(raw, " \t\r")
		if strings.HasPrefix(line, "```") {
			flushParagraph()
			if inCode {
				out.WriteString("</code></pre>\n")
				inCode = false
			} else {
				out.WriteString("<pre><code>")
				inCode = true
			}
			continue
		}
		if inCode {
			out.WriteString(html.EscapeString(line))
			out.WriteByte('\n')
			continue
		}
		if strings.TrimSpace(line) == "" {
			flushParagraph()
			continue
		}
		if strings.HasPrefix(line, "#") {
			flushParagraph()
			level := headingLevel(line)
			text := strings.TrimSpace(strings.TrimLeft(line, "#"))
			fmt.Fprintf(&out, "<h%d id=%q>%s</h%d>\n", level, slug(text), renderInline(text), level)
			continue
		}
		if strings.HasPrefix(line, "- ") {
			flushParagraph()
			out.WriteString("<ul>\n")
			fmt.Fprintf(&out, "<li>%s</li>\n", renderInline(strings.TrimSpace(strings.TrimPrefix(line, "- "))))
			out.WriteString("</ul>\n")
			continue
		}
		if strings.HasPrefix(line, "|") && strings.HasSuffix(line, "|") {
			flushParagraph()
			renderTableRow(&out, line)
			continue
		}
		paragraph = append(paragraph, strings.TrimSpace(line))
	}
	flushParagraph()
	if inCode {
		out.WriteString("</code></pre>\n")
	}
	return out.String()
}

func headingLevel(line string) int {
	level := 0
	for level < len(line) && line[level] == '#' {
		level++
	}
	if level < 1 {
		return 1
	}
	if level > 6 {
		return 6
	}
	return level
}

func renderTableRow(out *strings.Builder, line string) {
	cells := strings.Split(strings.Trim(line, "|"), "|")
	if isSeparatorRow(cells) {
		return
	}
	out.WriteString("<table><tr>")
	for _, cell := range cells {
		fmt.Fprintf(out, "<td>%s</td>", renderInline(strings.TrimSpace(cell)))
	}
	out.WriteString("</tr></table>\n")
}

func isSeparatorRow(cells []string) bool {
	for _, cell := range cells {
		clean := strings.Trim(strings.TrimSpace(cell), ":-")
		if clean != "" {
			return false
		}
	}
	return true
}

func renderInline(text string) string {
	escaped := html.EscapeString(text)
	escaped = markdownLinkPattern.ReplaceAllStringFunc(escaped, func(match string) string {
		parts := markdownLinkPattern.FindStringSubmatch(match)
		if len(parts) != 3 {
			return match
		}
		return fmt.Sprintf(`<a href="%s">%s</a>`, html.EscapeString(rewriteMarkdownHref(parts[2])), parts[1])
	})
	escaped = inlineCodePattern.ReplaceAllString(escaped, "<code>$1</code>")
	return boldPattern.ReplaceAllString(escaped, "<strong>$1</strong>")
}

func rewriteMarkdownHref(href string) string {
	if strings.HasPrefix(href, "http://") || strings.HasPrefix(href, "https://") || strings.HasPrefix(href, "#") {
		return href
	}
	anchor := ""
	if idx := strings.Index(href, "#"); idx >= 0 {
		anchor = href[idx:]
		href = href[:idx]
	}
	if strings.HasSuffix(href, ".md") {
		href = strings.TrimSuffix(href, ".md") + "/"
	}
	return href + anchor
}

func slug(text string) string {
	var buf bytes.Buffer
	lastDash := false
	for _, r := range strings.ToLower(text) {
		if (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') {
			buf.WriteRune(r)
			lastDash = false
			continue
		}
		if !lastDash {
			buf.WriteByte('-')
			lastDash = true
		}
	}
	return strings.Trim(buf.String(), "-")
}

func writeLibraryAssets(outDir string) error {
	css := `:root{color-scheme:dark;--bg:#111418;--panel:#181d23;--text:#f2eee7;--muted:#b5ada2;--accent:#df4f3d;--line:#343b45}*{box-sizing:border-box}body{margin:0;background:var(--bg);color:var(--text);font:16px/1.6 system-ui,-apple-system,BlinkMacSystemFont,"Segoe UI",sans-serif}a{color:#7dc7ff}header{border-bottom:1px solid var(--line);background:#0d1014;position:sticky;top:0}nav{max-width:1180px;margin:0 auto;padding:14px 24px;display:flex;gap:16px;align-items:center;flex-wrap:wrap}.brand{font-weight:700;color:var(--text);text-decoration:none}main{max-width:1180px;margin:0 auto;padding:28px 24px;display:grid;grid-template-columns:260px minmax(0,1fr);gap:28px}aside{border-right:1px solid var(--line);padding-right:20px}aside a{display:block;padding:4px 0;color:var(--muted);text-decoration:none}aside a.active{color:var(--text);font-weight:700}article{min-width:0}h1,h2,h3{line-height:1.2}h1{font-size:34px}h2{margin-top:34px;border-top:1px solid var(--line);padding-top:22px}pre,code{background:#0c0f13;border:1px solid var(--line);border-radius:4px}code{padding:1px 4px}pre{padding:14px;overflow:auto}table{border-collapse:collapse;margin:12px 0;width:100%}td,th{border:1px solid var(--line);padding:6px 8px;vertical-align:top}.lang-switch{margin-left:auto}.index-grid{display:grid;grid-template-columns:repeat(auto-fit,minmax(220px,1fr));gap:16px}.index-card{border:1px solid var(--line);background:var(--panel);border-radius:6px;padding:16px;text-decoration:none;color:var(--text)}@media(max-width:820px){main{display:block}aside{border-right:0;border-bottom:1px solid var(--line);padding:0 0 18px;margin-bottom:22px}}`
	return os.WriteFile(filepath.Join(outDir, "styles.css"), []byte(css), 0o644) //nolint:gosec // G306: public static asset.
}

func writePage(outDir string, page docPage, pages []docPage) error {
	path := outputFile(outDir, page.URLPath)
	if err := os.MkdirAll(filepath.Dir(path), 0o750); err != nil {
		return err
	}
	html := renderPage(page, pages)
	return os.WriteFile(path, []byte(html), 0o644) //nolint:gosec // G306: public static asset.
}

func renderPage(page docPage, pages []docPage) string {
	var nav strings.Builder
	for _, candidate := range pages {
		if candidate.Lang != page.Lang {
			continue
		}
		class := ""
		if candidate.URLPath == page.URLPath {
			class = ` class="active"`
		}
		fmt.Fprintf(&nav, `<a%s href="%s">%s</a>`, class, relHref(page.URLPath, candidate.URLPath), html.EscapeString(candidate.Title))
		nav.WriteByte('\n')
	}
	return fmt.Sprintf(`<!doctype html>
<html lang="%s">
<head>
  <meta charset="utf-8">
  <meta name="viewport" content="width=device-width, initial-scale=1">
  <title>%s</title>
  <link rel="stylesheet" href="%s">
</head>
<body>
  <header><nav><a class="brand" href="%s">RGB System Library</a><a href="%s">English</a><a href="%s">Portuguese</a></nav></header>
  <main>
    <aside>%s</aside>
    <article>%s</article>
  </main>
</body>
</html>
`, page.Lang, html.EscapeString(page.Title), relHref(page.URLPath, "/styles.css"), relHref(page.URLPath, "/"), relHref(page.URLPath, "/core/en/"), relHref(page.URLPath, "/core/PT-br/"), nav.String(), page.Body)
}

func writeLibraryIndex(outDir string, pages []docPage) error {
	var cards strings.Builder
	for _, lang := range []struct {
		Path  string
		Title string
		Desc  string
	}{
		{Path: "/core/en/", Title: "English Library", Desc: "Official RGB System Core reference."},
		{Path: "/core/PT-br/", Title: "Biblioteca em Portuguese", Desc: "Projecao localizada das regras Core."},
	} {
		fmt.Fprintf(&cards, `<a class="index-card" href="%s"><h2>%s</h2><p>%s</p></a>`, strings.TrimPrefix(lang.Path, "/"), lang.Title, lang.Desc)
	}
	index := fmt.Sprintf(`<!doctype html>
<html lang="en">
<head>
  <meta charset="utf-8">
  <meta name="viewport" content="width=device-width, initial-scale=1">
  <title>RGB System Library</title>
  <link rel="stylesheet" href="styles.css">
</head>
<body>
  <header><nav><a class="brand" href="./">RGB System Library</a></nav></header>
  <main>
    <article>
      <h1>RGB System Library</h1>
      <div class="index-grid">%s</div>
    </article>
  </main>
</body>
</html>
`, cards.String())
	_ = pages
	return os.WriteFile(filepath.Join(outDir, "index.html"), []byte(index), 0o644) //nolint:gosec // G306: public static asset.
}

func relHref(fromURL, toURL string) string {
	fromDir := strings.Trim(strings.TrimSuffix(fromURL, "/"), "/")
	depth := 0
	if fromDir != "" {
		depth = len(strings.Split(fromDir, "/"))
	}
	prefix := strings.Repeat("../", depth)
	target := strings.TrimPrefix(toURL, "/")
	if target == "" {
		return prefix
	}
	return prefix + target
}
