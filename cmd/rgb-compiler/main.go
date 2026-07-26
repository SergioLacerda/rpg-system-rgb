// Command rgb-compiler renders the canonical RGB System documentation into
// print-ready pages for manual PDF export, and optionally a static HTML
// site.
package main

import (
	"fmt"
	"os"

	"github.com/SergioLacerda/rpg-system-rgb/internal/app"
)

func main() {
	if len(os.Args) < 2 {
		fatal("usage: rgb-compiler <all|no-html> [repo-root]")
	}

	repoRoot := "."
	if len(os.Args) >= 3 {
		repoRoot = os.Args[2]
	}

	var err error
	switch os.Args[1] {
	case "all":
		err = app.Compile(repoRoot, true)
	case "no-html":
		err = app.Compile(repoRoot, false)
	default:
		fatal(fmt.Sprintf("unknown subcommand %q", os.Args[1]))
	}
	if err != nil {
		fatal(err.Error())
	}
}

func fatal(message string) {
	fmt.Fprintln(os.Stderr, message)
	os.Exit(1)
}
