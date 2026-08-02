package components

import (
	"os"
	"path/filepath"
	"testing"
)

func TestComponentContractAllowsBoundariesToDeclareIdentity(t *testing.T) {
	component := Component{ID: "core", Name: "Core", Description: "Rules engine"}
	if component.ID != "core" || component.Name == "" || component.Description == "" {
		t.Fatalf("unexpected component contract: %+v", component)
	}
}

func TestLoadSemanticIndexReadsValidJSONAndRejectsInvalidInput(t *testing.T) {
	dir := t.TempDir()
	indexFile := filepath.Join(dir, "index.json")
	if err := os.WriteFile(indexFile, []byte(`{"schema":"schema","units":[{"id":"unit-1"}]}`), 0o644); err != nil {
		t.Fatal(err)
	}
	index, err := LoadSemanticIndex(indexFile)
	if err != nil {
		t.Fatal(err)
	}
	if index.Schema != "schema" || len(index.Units) != 1 || index.Units[0].ID != "unit-1" {
		t.Fatalf("unexpected index: %+v", index)
	}

	if _, err := LoadSemanticIndex(filepath.Join(dir, "missing.json")); err == nil {
		t.Fatal("expected missing index to fail")
	}
	if err := os.WriteFile(indexFile, []byte(`{`), 0o644); err != nil {
		t.Fatal(err)
	}
	if _, err := LoadSemanticIndex(indexFile); err == nil {
		t.Fatal("expected invalid JSON to fail")
	}
}
