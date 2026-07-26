package compiler

import (
	"encoding/json"
	"fmt"
	"html"
	"os"
	"path/filepath"
	"strings"
)

const printPageTemplate = `<!doctype html>
<html lang="%s">
<head>
<meta charset="utf-8">
<title>%s</title>
<style>
body { font-family: system-ui, sans-serif; line-height: 1.5; }
pre { background: #f4f4f4; padding: 0.75rem; overflow-x: auto; }
table { border-collapse: collapse; }
th, td { border: 1px solid #ccc; padding: 0.4rem 0.6rem; text-align: left; }
.print-unit + .print-unit { margin-top: 2rem; }
@media print {
  .print-unit { page-break-before: always; }
  .print-unit:first-child { page-break-before: avoid; }
}
</style>
</head>
<body>
%s</body>
</html>
`

type printManifest struct {
	Units []printManifestUnit `json:"units"`
}

type printManifestUnit struct {
	ID              string            `json:"id"`
	ProjectionPaths map[string]string `json:"projection_paths"`
}

// RenderPrintPage builds a single, concatenated, print-ready HTML page for
// locale by reading generated/pdf/core-v2-rules.manifest.json (produced by
// internal/components/tooling) and rendering each referenced unit's
// Markdown source in manifest order. This is the source document for the
// PDF v1 manual export step (proposal.md Decision UN-02) — opening this
// page in a browser and using "Print to PDF" produces the release PDF.
func RenderPrintPage(repoRoot, locale string) (string, error) {
	manifestPath := filepath.Join(repoRoot, "generated", "pdf", "core-v2-rules.manifest.json")
	manifestBody, err := os.ReadFile(manifestPath) //nolint:gosec // G304: path is a fixed, repo-root-joined constant
	if err != nil {
		return "", err
	}
	var manifest printManifest
	if err := json.Unmarshal(manifestBody, &manifest); err != nil {
		return "", fmt.Errorf("%s: invalid JSON: %w", manifestPath, err)
	}

	projectionKey := "markdown_en"
	if locale == "PT-br" {
		projectionKey = "markdown_pt_br"
	}

	var sections strings.Builder
	for _, unit := range manifest.Units {
		relPath, ok := unit.ProjectionPaths[projectionKey]
		if !ok {
			return "", fmt.Errorf("%s: unit %s has no %s projection path", manifestPath, unit.ID, projectionKey)
		}
		sourcePath := filepath.Join(repoRoot, filepath.FromSlash(relPath))
		sourceBody, err := os.ReadFile(sourcePath) //nolint:gosec // G304: path is repo-root-joined from the generated manifest, by design
		if err != nil {
			return "", fmt.Errorf("%s: %w", unit.ID, err)
		}
		sections.WriteString(`<section class="print-unit">` + "\n")
		for _, block := range Parse(string(sourceBody)) {
			writeBlockHTML(&sections, block)
		}
		sections.WriteString("</section>\n")
	}

	title := "RGB Core V2 Rules"
	return fmt.Sprintf(printPageTemplate, html.EscapeString(localeTag(locale)), html.EscapeString(title), sections.String()), nil
}

// RenderPrintTree writes the print-ready page for both locales to
// generated/library/print/core-v2-rules-{locale}.html.
func RenderPrintTree(repoRoot string) error {
	outputDir := filepath.Join(repoRoot, "generated", "library", "print")
	if err := os.MkdirAll(outputDir, 0o750); err != nil {
		return err
	}
	for _, locale := range []string{"en", "PT-br"} {
		page, err := RenderPrintPage(repoRoot, locale)
		if err != nil {
			return err
		}
		outputPath := filepath.Join(outputDir, "core-v2-rules-"+locale+".html")
		if err := os.WriteFile(outputPath, []byte(page), 0o644); err != nil { //nolint:gosec // G306: generated HTML pages are intended to be readable
			return err
		}
	}
	return nil
}
