package main

import (
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
