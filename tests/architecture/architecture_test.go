// Package architecture enforces clean-architecture import boundaries (M001).
// It is a guardrail test (M016): it must keep passing as the codebase grows,
// and any relaxation of a rule requires an explicit design decision (ADR).
package architecture

import (
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

const modulePath = "github.com/SergioLacerda/rpg-system-rgb"

// repoRoot resolves the repository root relative to this test file.
func repoRoot(t *testing.T) string {
	t.Helper()
	root, err := filepath.Abs(filepath.Join("..", ".."))
	if err != nil {
		t.Fatalf("resolving repo root: %v", err)
	}
	return root
}

// goFilesUnder returns all .go files below dir, relative to the repo root.
func goFilesUnder(t *testing.T, root, dir string) []string {
	t.Helper()
	var files []string
	base := filepath.Join(root, dir)
	err := filepath.WalkDir(base, func(path string, entry os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if entry.IsDir() {
			return nil
		}
		if strings.HasSuffix(path, ".go") {
			rel, relErr := filepath.Rel(root, path)
			if relErr != nil {
				return relErr
			}
			files = append(files, rel)
		}
		return nil
	})
	if err != nil {
		t.Fatalf("walking %s: %v", dir, err)
	}
	return files
}

// importsOf parses only the import block of a file.
func importsOf(t *testing.T, root, relFile string) []string {
	t.Helper()
	fset := token.NewFileSet()
	parsed, err := parser.ParseFile(fset, filepath.Join(root, relFile), nil, parser.ImportsOnly)
	if err != nil {
		t.Fatalf("parsing %s: %v", relFile, err)
	}
	var imports []string
	for _, imp := range parsed.Imports {
		imports = append(imports, strings.Trim(imp.Path.Value, `"`))
	}
	return imports
}

// componentOf returns the component name for files under internal/components,
// or "" for the shared contract package itself.
func componentOf(relFile string) string {
	prefix := filepath.Join("internal", "components") + string(filepath.Separator)
	rest := strings.TrimPrefix(relFile, prefix)
	if rest == relFile {
		return ""
	}
	parts := strings.SplitN(rest, string(filepath.Separator), 2)
	if len(parts) < 2 {
		return "" // file directly in internal/components (shared contract)
	}
	return parts[0]
}

// TestComponentsDoNotImportOuterLayers keeps domain components free of
// application, entrypoint, and orchestration dependencies.
func TestComponentsDoNotImportOuterLayers(t *testing.T) {
	root := repoRoot(t)
	for _, file := range goFilesUnder(t, root, filepath.Join("internal", "components")) {
		for _, imp := range importsOf(t, root, file) {
			if imp == modulePath+"/internal/app" || strings.HasPrefix(imp, modulePath+"/cmd") {
				t.Errorf("%s imports outer layer %q; components must not depend on app or cmd", file, imp)
			}
		}
	}
}

// TestComponentsDoNotImportSiblings keeps components decoupled from each
// other. Shared types belong in internal/components (the contract package).
func TestComponentsDoNotImportSiblings(t *testing.T) {
	root := repoRoot(t)
	componentsPkg := modulePath + "/internal/components"
	for _, file := range goFilesUnder(t, root, filepath.Join("internal", "components")) {
		owner := componentOf(file)
		if owner == "" {
			continue
		}
		for _, imp := range importsOf(t, root, file) {
			if !strings.HasPrefix(imp, componentsPkg+"/") {
				continue
			}
			target := strings.TrimPrefix(imp, componentsPkg+"/")
			targetComponent := strings.SplitN(target, "/", 2)[0]
			if targetComponent != owner {
				t.Errorf("%s (component %q) imports sibling component %q; share contracts via internal/components instead", file, owner, targetComponent)
			}
		}
	}
}

// TestSharedContractStaysLeaf keeps internal/components (the shared contract
// package) free of any project-internal dependency.
func TestSharedContractStaysLeaf(t *testing.T) {
	root := repoRoot(t)
	for _, file := range goFilesUnder(t, root, filepath.Join("internal", "components")) {
		if componentOf(file) != "" {
			continue
		}
		for _, imp := range importsOf(t, root, file) {
			if strings.HasPrefix(imp, modulePath) {
				t.Errorf("%s imports %q; the shared contract package must only use the standard library", file, imp)
			}
		}
	}
}

// TestEntrypointsGoThroughApp keeps cmd binaries thin: they may only reach
// project code through the application layer.
func TestEntrypointsGoThroughApp(t *testing.T) {
	root := repoRoot(t)
	for _, file := range goFilesUnder(t, root, "cmd") {
		for _, imp := range importsOf(t, root, file) {
			if !strings.HasPrefix(imp, modulePath) {
				continue
			}
			if imp != modulePath+"/internal/app" {
				t.Errorf("%s imports %q; entrypoints must depend only on %s/internal/app", file, imp, modulePath)
			}
		}
	}
}
