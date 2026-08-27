package main

import (
	"strings"
	"testing"
)

func TestRunWithNoArgsPrintsStatus(t *testing.T) {
	if err := run(nil); err != nil {
		t.Fatalf("run(nil) returned error: %v", err)
	}
}

func TestRunRejectsUnknownSubcommand(t *testing.T) {
	err := run([]string{"bogus"})
	if err == nil {
		t.Fatal("expected error for unknown subcommand")
	}
	if !strings.Contains(err.Error(), "bogus") {
		t.Fatalf("error should mention the unknown subcommand, got: %v", err)
	}
}

func TestRunDispatchesRepoRootCommandsAndSurfacesErrors(t *testing.T) {
	root := t.TempDir()
	for _, command := range []string{"validate", "generate", "bundle"} {
		t.Run(command, func(t *testing.T) {
			err := run([]string{command, root})
			if err == nil {
				t.Fatalf("expected %s to fail against incomplete repo root", command)
			}
		})
	}
}

func TestRunRejectsRetiredCompileSubcommand(t *testing.T) {
	err := run([]string{"compile", "all"})
	if err == nil {
		t.Fatal("expected error: compile was retired (ADR-008), publication runs through rgb docs")
	}
	if !strings.Contains(err.Error(), "compile") {
		t.Fatalf("error should mention the retired subcommand, got: %v", err)
	}
}

func TestRunRejectsUnknownDocsSubcommand(t *testing.T) {
	err := run([]string{"docs", "bogus"})
	if err == nil {
		t.Fatal("expected error for unknown docs subcommand")
	}
	if !strings.Contains(err.Error(), "bogus") {
		t.Fatalf("error should mention the unknown docs subcommand, got: %v", err)
	}
}

func TestRunRejectsMissingDocsSubcommand(t *testing.T) {
	err := run([]string{"docs"})
	if err == nil || !strings.Contains(err.Error(), "missing docs subcommand") {
		t.Fatalf("expected missing docs subcommand error, got %v", err)
	}
}

func TestRunDocsSubcommandsSurfaceComponentErrors(t *testing.T) {
	root := t.TempDir()
	cases := [][]string{
		{"docs", "library", "--source", root, "--out", root + "/out"},
		{"docs", "pdf", "--public-dir", root + "/downloads", "--basename", "rgb", "--version", "v1"},
		{"docs", "skill", "--source", root + "/missing-skill", "--out", root + "/downloads", "--name", "rgb-specialist", "--version", "v1"},
		{"docs", "check", "--library", root + "/library", "--public-dir", root + "/downloads", "--basename", "rgb", "--version", "v1"},
	}
	for _, args := range cases {
		if err := run(args); err == nil {
			t.Fatalf("expected %v to fail against incomplete fixtures", args)
		}
	}
}

func TestParseDocsFlags(t *testing.T) {
	library, err := parseLibraryFlags([]string{"--source", "src", "--out", "out"})
	if err != nil {
		t.Fatal(err)
	}
	if library.SourceDir != "src" || library.OutDir != "out" {
		t.Fatalf("unexpected library options: %+v", library)
	}
	pdf, err := parsePDFFlags([]string{
		"--public-dir", "public",
		"--basename", "rgb",
		"--version", "v1",
		"--source-en", "en.pdf",
		"--source-pt-br", "pt.pdf",
	})
	if err != nil {
		t.Fatal(err)
	}
	if pdf.PublicDir != "public" || pdf.Basename != "rgb" || pdf.Version != "v1" || pdf.SourceEN != "en.pdf" || pdf.SourcePT != "pt.pdf" {
		t.Fatalf("unexpected pdf options: %+v", pdf)
	}
	check, err := parseDocsCheckFlags([]string{"--library", "library", "--public-dir", "public", "--basename", "rgb", "--version", "v1"})
	if err != nil {
		t.Fatal(err)
	}
	if check.LibraryDir != "library" || check.PublicDir != "public" || check.Basename != "rgb" || check.Version != "v1" {
		t.Fatalf("unexpected check options: %+v", check)
	}
}

func TestParseDocsFlagsRejectInvalidAndPositionalArgs(t *testing.T) {
	if _, err := parseLibraryFlags([]string{"--bad"}); err == nil {
		t.Fatal("expected invalid library flag to fail")
	}
	if _, err := parsePDFFlags([]string{"--bad"}); err == nil {
		t.Fatal("expected invalid pdf flag to fail")
	}
	if _, err := parsePDFFlags([]string{"extra"}); err == nil {
		t.Fatal("expected positional pdf arg to fail")
	}
	if _, err := parseSkillFlags([]string{"--bad"}); err == nil {
		t.Fatal("expected invalid skill flag to fail")
	}
	if _, err := parseSkillFlags([]string{"extra"}); err == nil {
		t.Fatal("expected positional skill arg to fail")
	}
	if _, err := parseDocsCheckFlags([]string{"--bad"}); err == nil {
		t.Fatal("expected invalid docs check flag to fail")
	}
	if _, err := parseDocsCheckFlags([]string{"extra"}); err == nil {
		t.Fatal("expected positional docs check arg to fail")
	}
}

func TestParseSkillFlags(t *testing.T) {
	options, err := parseSkillFlags([]string{
		"--source", "skills/specialist",
		"--out", "downloads",
		"--name", "rgb-specialist",
		"--version", "v0.1",
	})
	if err != nil {
		t.Fatal(err)
	}
	if options.SourceDir != "skills/specialist" || options.OutDir != "downloads" || options.Name != "rgb-specialist" || options.Version != "v0.1" {
		t.Fatalf("unexpected skill options: %+v", options)
	}
}

func TestParseLibraryFlagsRejectsPositionalArgs(t *testing.T) {
	_, err := parseLibraryFlags([]string{"extra"})
	if err == nil {
		t.Fatal("expected positional docs library argument to be rejected")
	}
}

func TestRunRejectsUnknownReleaseSubcommand(t *testing.T) {
	err := run([]string{"release", "bogus"})
	if err == nil {
		t.Fatal("expected error for unknown release subcommand")
	}
	if !strings.Contains(err.Error(), "bogus") {
		t.Fatalf("error should mention the unknown release subcommand, got: %v", err)
	}
}

func TestRunRejectsMissingReleaseSubcommand(t *testing.T) {
	err := run([]string{"release"})
	if err == nil || !strings.Contains(err.Error(), "missing release subcommand") {
		t.Fatalf("expected missing release subcommand error, got %v", err)
	}
}

func TestRunReleaseSubcommandsSurfaceComponentErrors(t *testing.T) {
	root := t.TempDir()
	cases := [][]string{
		{"release", "manifest", "--public-dir", root + "/downloads", "--basename", "rgb", "--version", "v1"},
		{"release", "check", "--public-dir", root + "/downloads", "--basename", "rgb", "--version", "v1"},
		{"release", "skill-manifest", "--public-dir", root + "/downloads", "--name", "rgb-specialist", "--version", "v1"},
		{"release", "skill-check", "--public-dir", root + "/downloads", "--name", "rgb-specialist", "--version", "v1"},
	}
	for _, args := range cases {
		if err := run(args); err == nil {
			t.Fatalf("expected %v to fail against incomplete fixtures", args)
		}
	}
}

func TestParseReleaseFlagsAcceptsAllFieldsAndRejectsInvalidFlags(t *testing.T) {
	paths, err := parseReleaseFlags("release manifest", []string{
		"--public-dir", "public",
		"--basename", "rgb",
		"--version", "v1",
		"--manifest", "manifest.json",
		"--checksums", "SHA256SUMS",
	})
	if err != nil {
		t.Fatal(err)
	}
	if paths.PublicDir != "public" || paths.Basename != "rgb" || paths.Version != "v1" || paths.Manifest != "manifest.json" || paths.Checksums != "SHA256SUMS" {
		t.Fatalf("unexpected release paths: %+v", paths)
	}
	if _, err := parseReleaseFlags("release manifest", []string{"--bad"}); err == nil {
		t.Fatal("expected invalid release flag to fail")
	}
}

func TestParseReleaseFlagsRejectsPositionalArgs(t *testing.T) {
	_, err := parseReleaseFlags("release manifest", []string{
		"--public-dir", "public",
		"--basename", "rgb-system-core-v2",
		"--version", "v0.2",
		"--manifest", "manifest.json",
		"--checksums", "SHA256SUMS",
		"extra",
	})
	if err == nil {
		t.Fatal("expected positional release argument to be rejected")
	}
}

func TestRunRejectsUnknownReleaseSkillSubcommand(t *testing.T) {
	err := run([]string{"release", "skill-bogus"})
	if err == nil {
		t.Fatal("expected error for unknown release subcommand")
	}
	if !strings.Contains(err.Error(), "skill-bogus") {
		t.Fatalf("error should mention the unknown release subcommand, got: %v", err)
	}
}

func TestParseSkillReleaseFlagsAcceptsAllFieldsAndRejectsInvalidFlags(t *testing.T) {
	paths, err := parseSkillReleaseFlags("release skill-manifest", []string{
		"--public-dir", "public",
		"--name", "rgb-specialist",
		"--version", "v0.1",
		"--manifest", "skill-manifest.json",
		"--checksums", "SKILL-SHA256SUMS",
	})
	if err != nil {
		t.Fatal(err)
	}
	if paths.PublicDir != "public" || paths.Name != "rgb-specialist" || paths.Version != "v0.1" || paths.Manifest != "skill-manifest.json" || paths.Checksums != "SKILL-SHA256SUMS" {
		t.Fatalf("unexpected skill release paths: %+v", paths)
	}
	if _, err := parseSkillReleaseFlags("release skill-manifest", []string{"--bad"}); err == nil {
		t.Fatal("expected invalid skill release flag to fail")
	}
}

func TestParseSkillReleaseFlagsRejectsPositionalArgs(t *testing.T) {
	_, err := parseSkillReleaseFlags("release skill-manifest", []string{
		"--public-dir", "public",
		"--name", "rgb-specialist",
		"--version", "v0.1",
		"--manifest", "skill-manifest.json",
		"--checksums", "SKILL-SHA256SUMS",
		"extra",
	})
	if err == nil {
		t.Fatal("expected positional skill release argument to be rejected")
	}
}
