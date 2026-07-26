package main

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

func main() {
	if err := validateProjectPaths(); err != nil {
		fmt.Fprintf(os.Stderr, "project-path validation failed: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("project-path validation passed")
}

func validateProjectPaths() error {
	requiredDirs := []string{
		"docs/core",
		"docs/adr",
		"generated",
	}
	for _, dir := range requiredDirs {
		if !pathDirExists(dir) {
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
		if pathExists(dir) {
			return fmt.Errorf("forbidden legacy directory exists: %s", dir)
		}
	}

	publicRoots := []string{
		"README.md",
		"docs",
		"generated",
	}
	for _, root := range publicRoots {
		if err := rejectRuntimeRefs(root); err != nil {
			return err
		}
	}
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
	bytes, err := os.ReadFile(path)
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
