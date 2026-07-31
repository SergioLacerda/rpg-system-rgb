package tooling

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

var forbiddenPublicPathMarkers = []string{
	".analysis/",
	".sdd/",
	".strategist/",
}

// publicPathWalkExclusions lists, per publicRoots entry, immediate
// subdirectories that are not produced by this repository's own build and
// must not be scanned for leaked internal-path markers. generated/client/
// is written by the governance/Strategist tooling's own client-build step —
// not by anything under cmd/ or internal/.
var publicPathWalkExclusions = map[string][]string{
	"generated": {"client"},
}

// ValidateProjectPaths checks the repository's public-facing directory
// shape and rejects internal path markers leaking into public surfaces,
// migrated from scripts/validate_project_paths.go. Unlike the original
// script, this takes repoRoot explicitly instead of assuming the process
// working directory is the repository root.
func ValidateProjectPaths(repoRoot string) error {
	requiredDirs := []string{
		"docs/core",
		"docs/adr",
		"generated",
	}
	for _, dir := range requiredDirs {
		if !pathDirExists(filepath.Join(repoRoot, filepath.FromSlash(dir))) {
			return fmt.Errorf("required directory missing: %s", dir)
		}
	}

	forbiddenDirs := []string{
		"docs/generated",
		"docs/en",
		"docs/PT-br",
		"docs/semantic",
	}
	for _, dir := range forbiddenDirs {
		if pathExists(filepath.Join(repoRoot, filepath.FromSlash(dir))) {
			return fmt.Errorf("forbidden legacy directory exists: %s", dir)
		}
	}

	publicRoots := []string{
		"README.md",
		"docs",
		"generated",
	}
	for _, root := range publicRoots {
		if err := rejectRuntimeRefs(filepath.Join(repoRoot, filepath.FromSlash(root)), publicPathWalkExclusions[root]); err != nil {
			return err
		}
	}
	fmt.Println("project-path validation passed")
	return nil
}

func rejectRuntimeRefs(root string, excludeDirs []string) error {
	info, err := os.Stat(root)
	if err != nil {
		return err
	}
	if !info.IsDir() {
		return rejectRuntimeRefsInFile(root)
	}
	return filepath.WalkDir(root, func(path string, entry fs.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		if entry.IsDir() {
			if isExcludedWalkDir(root, path, excludeDirs) {
				return filepath.SkipDir
			}
			return nil
		}
		return rejectRuntimeRefsInFile(path)
	})
}

// isExcludedWalkDir reports whether path (relative to root) matches one of
// excludeDirs — an immediate, named subdirectory exclusion, not a glob.
func isExcludedWalkDir(root, path string, excludeDirs []string) bool {
	rel, err := filepath.Rel(root, path)
	if err != nil {
		return false
	}
	rel = filepath.ToSlash(rel)
	for _, excluded := range excludeDirs {
		if rel == excluded {
			return true
		}
	}
	return false
}

func rejectRuntimeRefsInFile(path string) error {
	bytes, err := os.ReadFile(path) //nolint:gosec // G304: path comes from a WalkDir over the repo's own public tree, by design
	if err != nil {
		return err
	}
	content := string(bytes)
	for _, marker := range forbiddenPublicPathMarkers {
		if strings.Contains(content, marker) {
			return fmt.Errorf("%s references runtime/internal path marker %s", path, marker)
		}
	}
	return nil
}

func pathExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

func pathDirExists(path string) bool {
	info, err := os.Stat(path)
	return err == nil && info.IsDir()
}
