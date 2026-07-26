package features

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// readFeatureFile reads a .feature file relative to tests/features.
func readFeatureFile(t *testing.T, file string) string {
	t.Helper()
	body, err := os.ReadFile(file)
	if err != nil {
		t.Fatal(err)
	}
	return string(body)
}

// discoverFeatureFiles walks tests/features for every .feature file, so new
// files are picked up without editing this test (structural review F2).
func discoverFeatureFiles(t *testing.T) []string {
	t.Helper()
	files, err := filepath.Glob("*/*.feature")
	if err != nil {
		t.Fatal(err)
	}
	if len(files) == 0 {
		t.Fatal("no .feature files found under tests/features")
	}
	return files
}

func TestCoreV2FeatureFilesExist(t *testing.T) {
	for _, file := range discoverFeatureFiles(t) {
		t.Run(file, func(t *testing.T) {
			text := readFeatureFile(t, file)
			for _, required := range []string{"Feature:", "Scenario"} {
				if !strings.Contains(text, required) {
					t.Fatalf("%s missing %q", file, required)
				}
			}
		})
	}
}
