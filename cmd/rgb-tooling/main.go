// Command rgb-tooling validates and generates RGB System documentation
// artifacts. It replaces the standalone scripts/*.go module.
package main

import (
	"fmt"
	"os"

	"github.com/SergioLacerda/rpg-system-rgb/internal/app"
)

func main() {
	if len(os.Args) < 2 {
		fatal("usage: rgb-tooling <validate|generate|bundle> [repo-root]")
	}

	repoRoot := "."
	if len(os.Args) >= 3 {
		repoRoot = os.Args[2]
	}

	var err error
	switch os.Args[1] {
	case "validate":
		err = app.ValidateDocs(repoRoot)
	case "generate":
		err = app.GenerateProjections(repoRoot)
	case "bundle":
		err = app.BuildBundle(repoRoot)
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
