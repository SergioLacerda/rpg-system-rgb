// Command rgb-tooling validates and generates RGB System documentation
// artifacts. It replaces the standalone scripts/*.go module.
//
// Deprecated: cmd/rgb is now the canonical CLI (see ADR-006). rgb-tooling
// is kept working for existing callers (tests/semantic_docs,
// docs/core/semantic/README.md) but new usage should prefer
// `rgb validate|generate|bundle`.
package main

import (
	"fmt"
	"io"
	"os"

	"github.com/SergioLacerda/rpg-system-rgb/internal/app"
)

func main() {
	if err := run(cliArgs()); err != nil {
		fatal(err.Error())
	}
}

func run(args []string) error {
	if len(args) < 1 {
		return fmt.Errorf("usage: rgb-tooling <validate|generate|bundle> [repo-root]")
	}

	repoRoot := "."
	if len(args) >= 2 {
		repoRoot = args[1]
	}

	switch args[0] {
	case "validate":
		return app.ValidateDocs(repoRoot)
	case "generate":
		return app.GenerateProjections(repoRoot)
	case "bundle":
		return app.BuildBundle(repoRoot)
	default:
		return fmt.Errorf("unknown subcommand %q", args[0])
	}
}

func fatal(message string) {
	_, _ = fmt.Fprintln(stderr, message)
	exitProcess(1)
}

var (
	cliArgs               = func() []string { return os.Args[1:] }
	stderr      io.Writer = os.Stderr
	exitProcess           = os.Exit
)
