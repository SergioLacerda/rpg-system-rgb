package app

import (
	"fmt"

	"github.com/SergioLacerda/rpg-system-rgb/internal/components/compiler"
	"github.com/SergioLacerda/rpg-system-rgb/internal/components/library"
)

// Compile renders the canonical docs into the print-ready pages and PDF
// export instructions, and, when includeHTML is true, also the static
// HTML site and its navigation. Compile assumes
// generated/pdf/core-v2-rules.manifest.json already exists (run
// GenerateProjections first, e.g. via `make generate`).
func Compile(repoRoot string, includeHTML bool) error {
	if includeHTML {
		rendered, err := compiler.RenderHTMLTree(repoRoot)
		if err != nil {
			return fmt.Errorf("rendering HTML tree: %w", err)
		}
		fmt.Printf("rendered %d HTML pages\n", rendered)

		if err := library.BuildSite(repoRoot); err != nil {
			return fmt.Errorf("building site navigation: %w", err)
		}
		fmt.Println("built site navigation")
	}

	if err := compiler.RenderPrintTree(repoRoot); err != nil {
		return fmt.Errorf("rendering print pages: %w", err)
	}
	fmt.Println("rendered print-ready pages")

	if err := library.WritePDFExportInstructions(repoRoot); err != nil {
		return fmt.Errorf("writing PDF export instructions: %w", err)
	}
	fmt.Println("wrote PDF export instructions")

	return nil
}
