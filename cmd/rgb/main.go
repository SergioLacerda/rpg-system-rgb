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
	case "docs":
		return runDocs(args[1:])
	default:
		return fmt.Errorf("unknown subcommand %q (want validate|generate|bundle|release|docs)", args[0])
	}
}

func optionalRepoRoot(args []string) string {
	if len(args) >= 1 {
		return args[0]
	}
	return "."
}

//nolint:gocyclo // Mirrors the top-level CLI dispatch shape for docs subcommands.
func runDocs(args []string) error {
	if len(args) < 1 {
		return fmt.Errorf("missing docs subcommand (want library|pdf|check)")
	}
	switch args[0] {
	case "library":
		options, err := parseLibraryFlags(args[1:])
		if err != nil {
			return err
		}
		return app.BuildLibrary(options)
	case "pdf":
		options, err := parsePDFFlags(args[1:])
		if err != nil {
			return err
		}
		return app.PublishPDFs(options)
	case "check":
		options, err := parseDocsCheckFlags(args[1:])
		if err != nil {
			return err
		}
		return app.CheckPublication(options)
	default:
		return fmt.Errorf("unknown docs subcommand %q (want library|pdf|check)", args[0])
	}
}

func parseLibraryFlags(args []string) (app.LibraryOptions, error) {
	var options app.LibraryOptions
	flags := flag.NewFlagSet("docs library", flag.ContinueOnError)
	flags.SetOutput(io.Discard)
	flags.StringVar(&options.SourceDir, "source", "docs", "documentation source directory")
	flags.StringVar(&options.OutDir, "out", "web/landing/public/library", "public Library output directory")
	if err := flags.Parse(args); err != nil {
		return app.LibraryOptions{}, err
	}
	if flags.NArg() != 0 {
		return app.LibraryOptions{}, fmt.Errorf("unexpected docs library arguments: %v", flags.Args())
	}
	return options, nil
}

func parsePDFFlags(args []string) (app.PDFOptions, error) {
	var options app.PDFOptions
	flags := flag.NewFlagSet("docs pdf", flag.ContinueOnError)
	flags.SetOutput(io.Discard)
	flags.StringVar(&options.PublicDir, "public-dir", "web/landing/public/downloads", "public downloads directory")
	flags.StringVar(&options.Basename, "basename", "rgb-system-core-v2", "release artifact basename")
	flags.StringVar(&options.Version, "version", "v0.2", "release version")
	flags.StringVar(&options.SourceEN, "source-en", "", "reviewed English PDF source")
	flags.StringVar(&options.SourcePT, "source-pt-br", "", "reviewed Portuguese PDF source")
	if err := flags.Parse(args); err != nil {
		return app.PDFOptions{}, err
	}
	if flags.NArg() != 0 {
		return app.PDFOptions{}, fmt.Errorf("unexpected docs pdf arguments: %v", flags.Args())
	}
	return options, nil
}

func parseDocsCheckFlags(args []string) (app.PublicationCheckOptions, error) {
	var options app.PublicationCheckOptions
	flags := flag.NewFlagSet("docs check", flag.ContinueOnError)
	flags.SetOutput(io.Discard)
	flags.StringVar(&options.LibraryDir, "library", "web/landing/public/library", "public Library directory")
	flags.StringVar(&options.PublicDir, "public-dir", "web/landing/public/downloads", "public downloads directory")
	flags.StringVar(&options.Basename, "basename", "rgb-system-core-v2", "release artifact basename")
	flags.StringVar(&options.Version, "version", "v0.2", "release version")
	if err := flags.Parse(args); err != nil {
		return app.PublicationCheckOptions{}, err
	}
	if flags.NArg() != 0 {
		return app.PublicationCheckOptions{}, fmt.Errorf("unexpected docs check arguments: %v", flags.Args())
	}
	return options, nil
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
