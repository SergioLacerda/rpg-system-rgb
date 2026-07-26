package compiler

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

// RenderHTMLTree renders every Markdown file under docs/core/{en,PT-br}/**
// into a matching .html file under generated/library/html/{locale}/**,
// mirroring the source tree 1:1 (docs/core/en/core/README.md ->
// generated/library/html/en/core/README.html). It returns the number of
// files rendered.
func RenderHTMLTree(repoRoot string) (int, error) {
	rendered := 0
	for _, locale := range []string{"en", "PT-br"} {
		sourceDir := filepath.Join(repoRoot, "docs", "core", locale)
		outputDir := filepath.Join(repoRoot, "generated", "library", "html", locale)
		count, err := renderLocaleTree(sourceDir, outputDir, locale)
		if err != nil {
			return rendered, err
		}
		rendered += count
	}
	return rendered, nil
}

func renderLocaleTree(sourceDir, outputDir, locale string) (int, error) {
	count := 0
	err := filepath.WalkDir(sourceDir, func(path string, entry fs.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		if entry.IsDir() || filepath.Ext(path) != ".md" {
			return nil
		}
		rel, err := filepath.Rel(sourceDir, path)
		if err != nil {
			return err
		}
		outputPath := filepath.Join(outputDir, strings.TrimSuffix(rel, ".md")+".html")
		if err := renderFile(path, outputPath, locale); err != nil {
			return fmt.Errorf("%s: %w", rel, err)
		}
		count++
		return nil
	})
	return count, err
}

func renderFile(sourcePath, outputPath, locale string) error {
	body, err := os.ReadFile(sourcePath) //nolint:gosec // G304: path comes from a WalkDir over the repo's own canonical docs tree, by design
	if err != nil {
		return err
	}
	blocks := Parse(string(body))
	title := titleOf(blocks, sourcePath)
	html := renderHTMLPage(localeTag(locale), title, blocks)

	if err := os.MkdirAll(filepath.Dir(outputPath), 0o750); err != nil {
		return err
	}
	return os.WriteFile(outputPath, []byte(html), 0o644) //nolint:gosec // G306: generated HTML pages are intended to be readable
}

// titleOf returns the first heading's text as a page title, falling back
// to the source file's base name when no heading is present.
func titleOf(blocks []Block, sourcePath string) string {
	for _, block := range blocks {
		if heading, ok := block.(Heading); ok {
			return spansToText(heading.Inline)
		}
	}
	base := filepath.Base(sourcePath)
	return strings.TrimSuffix(base, filepath.Ext(base))
}

func spansToText(spans []InlineSpan) string {
	var w strings.Builder
	for _, span := range spans {
		w.WriteString(span.Text)
	}
	return w.String()
}

// localeTag converts the repo's locale directory name into an HTML lang
// attribute value.
func localeTag(locale string) string {
	if locale == "PT-br" {
		return "pt-BR"
	}
	return locale
}
