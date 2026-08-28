package architecture

import (
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
)

func TestCheckFastIncludesPublicArtifactGate(t *testing.T) {
	body, err := os.ReadFile(filepath.Join(repoRoot(t), "Makefile"))
	if err != nil {
		t.Fatal(err)
	}
	deps := makeTargetDependencies(t, string(body), "check-fast")
	if !hasMakeDependency(deps, "public-artifact-check") {
		t.Fatalf("check-fast missing public-artifact-check in %v", deps)
	}
}

func TestCheckFastIncludesCoverageInventory(t *testing.T) {
	body, err := os.ReadFile(filepath.Join(repoRoot(t), "Makefile"))
	if err != nil {
		t.Fatal(err)
	}
	deps := makeTargetDependencies(t, string(body), "check-fast")
	if !hasMakeDependency(deps, "coverage-inventory") {
		t.Fatalf("check-fast missing coverage-inventory in %v", deps)
	}
}

func TestPublicArtifactCheckIncludesConcreteGates(t *testing.T) {
	body, err := os.ReadFile(filepath.Join(repoRoot(t), "Makefile"))
	if err != nil {
		t.Fatal(err)
	}
	deps := makeTargetDependencies(t, string(body), "public-artifact-check")
	for _, required := range []string{
		"pdf-author-unit-test",
		"docs-check",
		"pdf-editorial-check",
		"release-skill-check",
	} {
		if !hasMakeDependency(deps, required) {
			t.Fatalf("check-fast missing required public artifact gate %q in %v", required, deps)
		}
	}
}

func makeTargetDependencies(t *testing.T, body, target string) []string {
	t.Helper()
	pattern := regexp.MustCompile(`(?m)^` + regexp.QuoteMeta(target) + `:\s+(?P<deps>.*?)\s+FORCE\s+##`)
	match := pattern.FindStringSubmatch(body)
	if len(match) != 2 {
		t.Fatalf("expected %s target with FORCE sentinel", target)
	}
	return strings.Fields(match[1])
}

func hasMakeDependency(deps []string, required string) bool {
	for _, dep := range deps {
		if dep == required {
			return true
		}
	}
	return false
}
