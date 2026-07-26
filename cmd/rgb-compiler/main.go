// Command rgb-compiler renders the canonical RGB System documentation into
// a static HTML site and print-ready pages for manual PDF export.
package main

import (
	"fmt"
	"os"

	"github.com/SergioLacerda/rpg-system-rgb/internal/app"
)

func main() {
	repoRoot := "."
	if len(os.Args) >= 2 {
		repoRoot = os.Args[1]
	}
	if err := app.Compile(repoRoot); err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
}
