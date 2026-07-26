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
		if err := rejectRuntimeRefs(filepath.Join(repoRoot, filepath.FromSlash(root))); err != nil {
			return err
		}
	}
	fmt.Println("project-path validation passed")
	return nil
}

func rejectRuntimeRefs(root string) error {
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
			return nil
		}
		return rejectRuntimeRefsInFile(path)
	})
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
