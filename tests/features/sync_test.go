package features

import (
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
)

var scenarioNamePattern = regexp.MustCompile(`^\s*Scenario(?: Outline)?:\s*(.+)$`)

var mirrorAnchorPattern = regexp.MustCompile(`//\s*mirrors:\s*(\S+)#(.+)$`)

// repoRoot resolves the repository root relative to this test file
// (tests/features is two directories below the root).
func repoRoot(t *testing.T) string {
	t.Helper()
	root, err := filepath.Abs(filepath.Join("..", ".."))
	if err != nil {
		t.Fatalf("resolving repo root: %v", err)
	}
	return root
}

// scenarioNamesIn returns every Scenario and Scenario Outline title declared
// in a .feature file, in file order.
func scenarioNamesIn(t *testing.T, file string) []string {
	t.Helper()
	body := readFeatureFile(t, file)
	var names []string
	for _, line := range strings.Split(body, "\n") {
		if match := scenarioNamePattern.FindStringSubmatch(line); match != nil {
			names = append(names, strings.TrimSpace(match[1]))
		}
	}
	return names
}

// mirroredScenarios scans every Go source file under tests/ (excluding
// tests/features itself, which declares scenarios rather than mirroring
// them) for `// mirrors: <feature path>#<scenario name>` anchor comments.
func mirroredScenarios(t *testing.T, root string) map[string]bool {
	t.Helper()
	mirrored := map[string]bool{}
	testsDir := filepath.Join(root, "tests")
	featuresDir := filepath.Join(testsDir, "features")

	err := filepath.WalkDir(testsDir, func(path string, entry os.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		if entry.IsDir() {
			return nil
		}
		if !strings.HasSuffix(path, ".go") {
			return nil
		}
		if strings.HasPrefix(path, featuresDir+string(filepath.Separator)) {
			return nil
		}
		body, readErr := os.ReadFile(path)
		if readErr != nil {
			return readErr
		}
		for _, line := range strings.Split(string(body), "\n") {
			match := mirrorAnchorPattern.FindStringSubmatch(line)
			if match == nil {
				continue
			}
			mirrored[match[1]+"#"+strings.TrimSpace(match[2])] = true
		}
		return nil
	})
	if err != nil {
		t.Fatalf("scanning tests/ for mirror anchors: %v", err)
	}
	return mirrored
}

// TestEveryScenarioHasAMirror keeps Gherkin specifications and Go tests in
// sync (structural review F1): a scenario edited or added in a .feature
// file must be matched by a `// mirrors:` anchor somewhere under tests/, or
// this guardrail fails with the exact missing key.
func TestEveryScenarioHasAMirror(t *testing.T) {
	root := repoRoot(t)
	mirrored := mirroredScenarios(t, root)

	for _, file := range discoverFeatureFiles(t) {
		relFeaturePath := filepath.Join("tests", "features", file)
		for _, scenario := range scenarioNamesIn(t, file) {
			key := relFeaturePath + "#" + scenario
			if !mirrored[key] {
				t.Errorf("scenario %q in %s has no matching `// mirrors: %s` anchor under tests/", scenario, file, key)
			}
		}
	}
}
