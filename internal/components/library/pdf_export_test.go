package library

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestWritePDFExportInstructions(t *testing.T) {
	scratch := t.TempDir()
	if err := WritePDFExportInstructions(scratch); err != nil {
		t.Fatalf("WritePDFExportInstructions failed: %v", err)
	}
	body, err := os.ReadFile(filepath.Join(scratch, "generated", "library", "PDF_EXPORT.md"))
	if err != nil {
		t.Fatalf("PDF_EXPORT.md not written: %v", err)
	}
	text := string(body)
	for _, want := range []string{
		"core-v2-rules-en.html",
		"core-v2-rules-PT-br.html",
		"Print to PDF",
		"make compile",
	} {
		if !strings.Contains(text, want) {
			t.Fatalf("PDF_EXPORT.md missing %q:\n%s", want, text)
		}
	}
}
