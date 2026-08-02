// Command rgb is the unified CLI entrypoint for RGB System V2: it
// dispatches to the validate, generate, and bundle use cases exposed by
// internal/app. With no subcommand, it reports scaffold status.
package main

import (
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/SergioLacerda/rpg-system-rgb/internal/app"
)

func main() {
	if err := run(os.Args[1:]); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run(args []string) error {
	if len(args) < 1 {
		fmt.Println(app.Hello())
		return nil
	}

	switch args[0] {
	case "validate":
		repoRoot := optionalRepoRoot(args[1:])
		return app.ValidateDocs(repoRoot)
	case "generate":
		repoRoot := optionalRepoRoot(args[1:])
		return app.GenerateProjections(repoRoot)
	case "bundle":
		repoRoot := optionalRepoRoot(args[1:])
		return app.BuildBundle(repoRoot)
	case "release":
		return runRelease(args[1:])
	default:
		return fmt.Errorf("unknown subcommand %q (want validate|generate|bundle|release)", args[0])
	}
}

func optionalRepoRoot(args []string) string {
	if len(args) >= 1 {
		return args[0]
	}
	return "."
}

func runRelease(args []string) error {
	if len(args) < 1 {
		return fmt.Errorf("missing release subcommand (want manifest|check)")
	}
	switch args[0] {
	case "manifest":
		paths, err := parseReleaseFlags("release manifest", args[1:])
		if err != nil {
			return err
		}
		return app.WriteReleaseArtifactManifest(paths)
	case "check":
		paths, err := parseReleaseFlags("release check", args[1:])
		if err != nil {
			return err
		}
		return app.CheckReleaseArtifacts(paths)
	default:
		return fmt.Errorf("unknown release subcommand %q (want manifest|check)", args[0])
	}
}

func parseReleaseFlags(name string, args []string) (app.ReleaseArtifactPaths, error) {
	var paths app.ReleaseArtifactPaths
	flags := flag.NewFlagSet(name, flag.ContinueOnError)
	flags.SetOutput(io.Discard)
	flags.StringVar(&paths.PublicDir, "public-dir", "", "public downloads directory")
	flags.StringVar(&paths.Basename, "basename", "", "release artifact basename")
	flags.StringVar(&paths.Version, "version", "", "release version")
	flags.StringVar(&paths.Manifest, "manifest", "", "release manifest path")
	flags.StringVar(&paths.Checksums, "checksums", "", "SHA256SUMS path")
	if err := flags.Parse(args); err != nil {
		return app.ReleaseArtifactPaths{}, err
	}
	if flags.NArg() != 0 {
		return app.ReleaseArtifactPaths{}, fmt.Errorf("unexpected release arguments: %v", flags.Args())
	}
	return paths, nil
}
