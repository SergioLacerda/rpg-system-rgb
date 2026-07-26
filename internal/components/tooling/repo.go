package tooling

import (
	"fmt"
	"os"
	"path/filepath"
)

// findRepoRoot walks up from start looking for the repository root,
// identified by a README.md alongside a docs/ directory. Shared by every
// validator and the generator in this package (previously duplicated
// verbatim across scripts/validate_*.go).
func findRepoRoot(start string) (string, error) {
	for dir := start; ; dir = filepath.Dir(dir) {
		if fileExists(filepath.Join(dir, "README.md")) && dirExists(filepath.Join(dir, "docs")) {
			return dir, nil
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			return "", fmt.Errorf("could not find repository root from %s", start)
		}
	}
}

// repoRootFromFile resolves the repository root starting from the
// directory containing the given file path.
func repoRootFromFile(file string) (string, error) {
	abs, err := filepath.Abs(file)
	if err != nil {
		return "", err
	}
	return findRepoRoot(filepath.Dir(abs))
}

func fileExists(path string) bool {
	info, err := os.Stat(path)
	return err == nil && !info.IsDir()
}

func dirExists(path string) bool {
	info, err := os.Stat(path)
	return err == nil && info.IsDir()
}

// repoRelative converts an absolute or CWD-relative path into its
// repository-relative, slash-form representation, matching the literal
// path strings stored inside the JSON documents themselves (e.g.
// `source_index: "docs/core/semantic/core-v2.index.json"`). Used wherever
// a path argument must be compared against or written into JSON content,
// as opposed to used for actual file I/O.
func repoRelative(repoRoot, path string) string {
	abs, err := filepath.Abs(path)
	if err != nil {
		return path
	}
	rel, err := filepath.Rel(repoRoot, abs)
	if err != nil {
		return path
	}
	return filepath.ToSlash(rel)
}

func stringSet(values []string) map[string]bool {
	result := make(map[string]bool, len(values))
	for _, value := range values {
		result[value] = true
	}
	return result
}
