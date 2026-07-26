package library

import (
	"os"
	"path/filepath"
)

const pdfExportInstructions = `# RGB Library — PDF Export

PDF v1 is a documented manual export, not an automated build step. This
keeps the project's zero-external-dependency footprint intact for the
first release; automated PDF generation is an explicit future decision,
not part of this scope.

## Steps

1. Run ` + "`make compile`" + ` to (re)generate the print-ready pages.
2. Open the print-ready page for the locale you want to export:
   - English: ` + "`generated/library/print/core-v2-rules-en.html`" + `
   - Portuguese (Brazil): ` + "`generated/library/print/core-v2-rules-PT-br.html`" + `
3. In your browser, use Print (Ctrl/Cmd+P) and choose "Save as PDF" (or
   "Print to PDF").
4. Save the output as ` + "`RGB-Core-V2-Rules-<locale>.pdf`" + `.

The print-ready page already includes print CSS (` + "`@media print`" + `)
that starts each rule unit on its own page, in the same order as
` + "`generated/pdf/core-v2-rules.manifest.json`" + `.
`

// WritePDFExportInstructions writes the manual PDF export instructions to
// generated/library/PDF_EXPORT.md.
func WritePDFExportInstructions(repoRoot string) error {
	outputDir := filepath.Join(repoRoot, "generated", "library")
	if err := os.MkdirAll(outputDir, 0o750); err != nil {
		return err
	}
	return os.WriteFile(filepath.Join(outputDir, "PDF_EXPORT.md"), []byte(pdfExportInstructions), 0o644) //nolint:gosec // G306: generated documentation is intended to be readable
}
