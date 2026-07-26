package library

import (
	"fmt"
	"html"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

type navEntry struct {
	Title string
	Href  string
}

// BuildSite writes a navigation index page at
// generated/library/html/{locale}/index.html linking every rendered page
// under generated/library/html/{locale}/**, grouped by top-level directory
// (the same grouping docs/core/{locale}/** already uses). Walking the
// rendered tree directly guarantees no orphan pages by construction,
// without needing to resolve the many-to-one mapping between semantic
// units (docs/core/semantic/core-v2.index.json) and source files.
func BuildSite(repoRoot string) error {
	for _, locale := range []string{"en", "PT-br"} {
		if err := buildLocaleNav(repoRoot, locale); err != nil {
			return err
		}
	}
	return nil
}

func buildLocaleNav(repoRoot, locale string) error {
	htmlDir := filepath.Join(repoRoot, "generated", "library", "html", locale)
	groups := map[string][]navEntry{}
	var groupOrder []string

	err := filepath.WalkDir(htmlDir, func(path string, entry fs.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		if entry.IsDir() || filepath.Ext(path) != ".html" {
			return nil
		}
		rel, err := filepath.Rel(htmlDir, path)
		if err != nil {
			return err
		}
		rel = filepath.ToSlash(rel)
		if rel == "index.html" {
			return nil
		}
		group := "Overview"
		if parts := strings.SplitN(rel, "/", 2); len(parts) > 1 {
			group = titleCase(parts[0])
		}
		if _, seen := groups[group]; !seen {
			groupOrder = append(groupOrder, group)
		}
		groups[group] = append(groups[group], navEntry{Title: titleFromFilename(path), Href: rel})
		return nil
	})
	if err != nil {
		return err
	}

	sortNavGroups(groupOrder, groups)
	body := renderNavHTML(locale, groupOrder, groups)
	return os.WriteFile(filepath.Join(htmlDir, "index.html"), []byte(body), 0o644) //nolint:gosec // G306: generated HTML pages are intended to be readable
}

// sortNavGroups sorts group names alphabetically, with "Overview" always
// first, and sorts each group's entries by link target.
func sortNavGroups(groupOrder []string, groups map[string][]navEntry) {
	sort.Slice(groupOrder, func(i, j int) bool {
		if groupOrder[i] == "Overview" {
			return true
		}
		if groupOrder[j] == "Overview" {
			return false
		}
		return groupOrder[i] < groupOrder[j]
	})
	for _, entries := range groups {
		sort.Slice(entries, func(i, j int) bool { return entries[i].Href < entries[j].Href })
	}
}

func renderNavHTML(locale string, groupOrder []string, groups map[string][]navEntry) string {
	var w strings.Builder
	fmt.Fprintf(&w, "<!doctype html>\n<html lang=\"%s\">\n<head>\n<meta charset=\"utf-8\">\n<title>RGB Library</title>\n", html.EscapeString(localeTag(locale)))
	w.WriteString("<style>body{font-family:system-ui,sans-serif;max-width:40rem;margin:2rem auto;padding:0 1rem;line-height:1.5;}ul{padding-left:1.2rem;}</style>\n</head>\n<body>\n<h1>RGB Library</h1>\n")
	for _, group := range groupOrder {
		fmt.Fprintf(&w, "<h2>%s</h2>\n<ul>\n", html.EscapeString(group))
		for _, entry := range groups[group] {
			fmt.Fprintf(&w, "<li><a href=\"%s\">%s</a></li>\n", html.EscapeString(entry.Href), html.EscapeString(entry.Title))
		}
		w.WriteString("</ul>\n")
	}
	w.WriteString("</body>\n</html>\n")
	return w.String()
}

// titleFromFilename derives a display title from a rendered page's file
// name, without re-parsing its HTML or Markdown source (keeping this
// component decoupled from internal/components/compiler).
func titleFromFilename(path string) string {
	base := strings.TrimSuffix(filepath.Base(path), filepath.Ext(path))
	base = strings.ReplaceAll(base, "_", " ")
	base = strings.ReplaceAll(base, "-", " ")
	return titleCase(base)
}

func titleCase(text string) string {
	words := strings.Fields(text)
	for i, word := range words {
		runes := []rune(word)
		if len(runes) > 0 {
			runes[0] = []rune(strings.ToUpper(string(runes[0])))[0]
		}
		words[i] = string(runes)
	}
	return strings.Join(words, " ")
}

// localeTag converts the repo's locale directory name into an HTML lang
// attribute value.
func localeTag(locale string) string {
	if locale == "PT-br" {
		return "pt-BR"
	}
	return locale
}
