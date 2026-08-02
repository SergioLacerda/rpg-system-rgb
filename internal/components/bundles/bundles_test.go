package bundles

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func TestDescriptor(t *testing.T) {
	descriptor := Descriptor()
	if descriptor.ID != "bundles" || descriptor.Name == "" || descriptor.Description == "" {
		t.Fatalf("unexpected descriptor: %+v", descriptor)
	}
}

func TestBuildProjectsSemanticIndexUnits(t *testing.T) {
	repoRoot := validBundleRepo(t)
	bundle, err := Build(repoRoot)
	if err != nil {
		t.Fatal(err)
	}
	if bundle.Schema != bundleSchema {
		t.Fatalf("unexpected schema %s", bundle.Schema)
	}
	if bundle.GeneratedAt == "" {
		t.Fatal("GeneratedAt must be populated")
	}
	if len(bundle.Units) != 1 {
		t.Fatalf("expected one unit, got %+v", bundle.Units)
	}
	unit := bundle.Units[0]
	if unit.ID != "core.test" || unit.Kind != "rule" || unit.Title != "Test Rule" || unit.SourceStatus != "canonical" {
		t.Fatalf("unexpected projected unit: %+v", unit)
	}
	if unit.Relationships["relates_to"][0] != "core.other" {
		t.Fatalf("relationships not preserved: %+v", unit.Relationships)
	}
}

func TestBuildRejectsMissingIndex(t *testing.T) {
	if _, err := Build(t.TempDir()); err == nil {
		t.Fatal("expected missing semantic index to fail")
	}
}

func TestWriteCreatesBundleArtifact(t *testing.T) {
	repoRoot := validBundleRepo(t)
	if err := Write(repoRoot); err != nil {
		t.Fatal(err)
	}
	outPath := filepath.Join(repoRoot, "generated", "bundle", "rgb.bundle.json")
	bytes, err := os.ReadFile(outPath)
	if err != nil {
		t.Fatal(err)
	}
	var bundle Bundle
	if err := json.Unmarshal(bytes, &bundle); err != nil {
		t.Fatal(err)
	}
	if len(bundle.Units) != 1 || bundle.Units[0].ID != "core.test" {
		t.Fatalf("unexpected written bundle: %+v", bundle)
	}
}

func TestWriteSurfacesBuildErrors(t *testing.T) {
	if err := Write(t.TempDir()); err == nil {
		t.Fatal("expected missing semantic index to fail write")
	}
}

func TestWriteSurfacesOutputDirectoryErrors(t *testing.T) {
	repoRoot := validBundleRepo(t)
	if err := os.WriteFile(filepath.Join(repoRoot, "generated"), []byte("not a dir"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := Write(repoRoot); err == nil {
		t.Fatal("expected generated file to block bundle output directory creation")
	}
}

func validBundleRepo(t *testing.T) string {
	t.Helper()
	repoRoot := t.TempDir()
	indexPath := filepath.Join(repoRoot, "docs", "core", "semantic", "core-v2.index.json")
	if err := os.MkdirAll(filepath.Dir(indexPath), 0o755); err != nil {
		t.Fatal(err)
	}
	index := map[string]any{
		"schema": "rgb-docs-semantic-index/0.1",
		"units": []map[string]any{
			{
				"id":            "core.test",
				"kind":          "rule",
				"title":         "Test Rule",
				"source_status": "canonical",
				"relationships": map[string][]string{"relates_to": {"core.other"}},
			},
		},
	}
	bytes, err := json.Marshal(index)
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(indexPath, bytes, 0o644); err != nil {
		t.Fatal(err)
	}
	return repoRoot
}
