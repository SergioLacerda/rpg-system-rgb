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
	deps := makeTargetClosure(t, string(body), "check-fast")
	if !hasMakeDependency(deps, "public-artifact-check") {
		t.Fatalf("check-fast missing public-artifact-check in %v", deps)
	}
}

func TestCheckFastIncludesCoverageInventory(t *testing.T) {
	body, err := os.ReadFile(filepath.Join(repoRoot(t), "Makefile"))
	if err != nil {
		t.Fatal(err)
	}
	deps := makeTargetClosure(t, string(body), "check-fast")
	if !hasMakeDependency(deps, "coverage-inventory") {
		t.Fatalf("check-fast missing coverage-inventory in %v", deps)
	}
}

// CI runs one job per gate group; check-fast must stay the union of them so the
// local gate and the PR gate cannot drift apart.
func TestCheckFastIsUnionOfGateGroups(t *testing.T) {
	body, err := os.ReadFile(filepath.Join(repoRoot(t), "Makefile"))
	if err != nil {
		t.Fatal(err)
	}
	deps := makeTargetDependencies(t, string(body), "check-fast")
	for _, group := range []string{
		"check-go",
		"public-artifact-check",
		"check-generated-drift",
		"check-hygiene",
	} {
		if !hasMakeDependency(deps, group) {
			t.Fatalf("check-fast missing gate group %q in %v", group, deps)
		}
	}
}

func TestGateGroupsIncludeConcreteGates(t *testing.T) {
	body, err := os.ReadFile(filepath.Join(repoRoot(t), "Makefile"))
	if err != nil {
		t.Fatal(err)
	}
	groups := map[string][]string{
		"check-go": {
			"fmt-check", "vet", "lint", "test", "test-arch", "validate",
			"cover-check", "mutation-core", "coverage-inventory",
		},
		"check-hygiene": {"check-publication-runtime", "check-governance-files"},
	}
	for group, required := range groups {
		deps := makeTargetDependencies(t, string(body), group)
		for _, gate := range required {
			if !hasMakeDependency(deps, gate) {
				t.Fatalf("%s missing required gate %q in %v", group, gate, deps)
			}
		}
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
	deps, ok := lookupMakeTargetDependencies(body, target)
	if !ok {
		t.Fatalf("expected %s target with FORCE sentinel", target)
	}
	return deps
}

// makeTargetClosure returns every target reachable from target through rules
// defined in the Makefile itself (targets from make/*.mk are leaves).
func makeTargetClosure(t *testing.T, body, target string) []string {
	t.Helper()
	if _, ok := lookupMakeTargetDependencies(body, target); !ok {
		t.Fatalf("expected %s target with FORCE sentinel", target)
	}
	seen := map[string]bool{}
	var order []string
	queue := []string{target}
	for len(queue) > 0 {
		current := queue[0]
		queue = queue[1:]
		deps, _ := lookupMakeTargetDependencies(body, current)
		for _, dep := range deps {
			if seen[dep] {
				continue
			}
			seen[dep] = true
			order = append(order, dep)
			queue = append(queue, dep)
		}
	}
	return order
}

func lookupMakeTargetDependencies(body, target string) ([]string, bool) {
	pattern := regexp.MustCompile(`(?m)^` + regexp.QuoteMeta(target) + `:\s+(?P<deps>.*?)\s+FORCE\s+##`)
	match := pattern.FindStringSubmatch(body)
	if len(match) != 2 {
		return nil, false
	}
	return strings.Fields(match[1]), true
}

func hasMakeDependency(deps []string, required string) bool {
	for _, dep := range deps {
		if dep == required {
			return true
		}
	}
	return false
}
