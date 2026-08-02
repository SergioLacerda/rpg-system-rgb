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
