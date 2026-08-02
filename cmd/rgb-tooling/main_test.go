package main

import (
	"bytes"
	"strings"
	"testing"
)

func TestRunRejectsMissingAndUnknownCommands(t *testing.T) {
	if err := run(nil); err == nil || !strings.Contains(err.Error(), "usage") {
		t.Fatalf("expected usage error, got %v", err)
	}
	if err := run([]string{"unknown"}); err == nil || !strings.Contains(err.Error(), "unknown") {
		t.Fatalf("expected unknown command error, got %v", err)
	}
}

func TestRunDispatchesCommandsAndSurfacesErrors(t *testing.T) {
	repoRoot := t.TempDir()
	for _, command := range []string{"validate", "generate", "bundle"} {
		t.Run(command, func(t *testing.T) {
			if err := run([]string{command, repoRoot}); err == nil {
				t.Fatalf("expected %s to fail against an incomplete repo root", command)
			}
		})
	}
}

func TestMainReportsRunErrors(t *testing.T) {
	originalArgs := cliArgs
	originalStderr := stderr
	originalExit := exitProcess
	t.Cleanup(func() {
		cliArgs = originalArgs
		stderr = originalStderr
		exitProcess = originalExit
	})

	var output bytes.Buffer
	var exitCode int
	cliArgs = func() []string { return []string{"unknown"} }
	stderr = &output
	exitProcess = func(code int) { exitCode = code }

	main()

	if exitCode != 1 {
		t.Fatalf("exit code got %d want 1", exitCode)
	}
	if !strings.Contains(output.String(), "unknown subcommand") {
		t.Fatalf("stderr missing command error: %q", output.String())
	}
}

func TestFatalWritesMessageAndExits(t *testing.T) {
	originalStderr := stderr
	originalExit := exitProcess
	t.Cleanup(func() {
		stderr = originalStderr
		exitProcess = originalExit
	})

	var output bytes.Buffer
	var exitCode int
	stderr = &output
	exitProcess = func(code int) { exitCode = code }

	fatal("boom")

	if exitCode != 1 {
		t.Fatalf("exit code got %d want 1", exitCode)
	}
	if strings.TrimSpace(output.String()) != "boom" {
		t.Fatalf("stderr got %q want boom", output.String())
	}
}
