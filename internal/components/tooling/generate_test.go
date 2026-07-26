package tooling

import (
	"bytes"
	"encoding/json"
	"io/fs"
	"os"
	"path/filepath"
	"testing"
)

// TestGenerateDefaultMatchesCommittedOutput regenerates every projection
// declared in the projection manifest into a scratch copy of the repo's
// docs/core/semantic tree and compares each one byte-for-byte against the
// committed output. This is the migration's correctness check: the new
// package must produce identical output to the retired
// scripts/generate_semantic_projections.go. It only compares the
// projection manifest's own declared outputs, not everything under
// generated/ — that directory also holds output from
// internal/components/compiler and internal/components/library, which
// GenerateDefault does not own.
func TestGenerateDefaultMatchesCommittedOutput(t *testing.T) {
	root := repoRoot(t)
	scratch := t.TempDir()

	semanticSrc := filepath.Join(root, "docs", "core", "semantic")
	semanticDst := filepath.Join(scratch, "docs", "core", "semantic")
	mustCopyTree(t, semanticSrc, semanticDst)

	if err := GenerateDefault(scratch); err != nil {
		t.Fatalf("GenerateDefault failed: %v", err)
	}

	for _, outputPath := range declaredProjectionOutputs(t, root) {
		want, err := os.ReadFile(filepath.Join(root, filepath.FromSlash(outputPath)))
		if err != nil {
			t.Fatal(err)
		}
		got, err := os.ReadFile(filepath.Join(scratch, filepath.FromSlash(outputPath)))
		if err != nil {
			t.Errorf("regenerated output missing for %s: %v", outputPath, err)
			continue
		}
		if !bytes.Equal(want, got) {
			t.Errorf("regenerated output for %s does not match committed output", outputPath)
		}
	}
}

// declaredProjectionOutputs reads the projection manifest's own
// output_path declarations, so the comparison stays in sync with whatever
// GenerateDefault actually produces.
func declaredProjectionOutputs(t *testing.T, root string) []string {
	t.Helper()
	manifestPath := filepath.Join(root, "docs", "core", "semantic", "projection-manifest.v0.1.json")
	body, err := os.ReadFile(manifestPath)
	if err != nil {
		t.Fatal(err)
	}
	var manifest struct {
		Projections []struct {
			OutputPath string `json:"output_path"`
		} `json:"projections"`
	}
	if err := json.Unmarshal(body, &manifest); err != nil {
		t.Fatal(err)
	}
	outputs := make([]string, len(manifest.Projections))
	for i, projection := range manifest.Projections {
		outputs[i] = projection.OutputPath
	}
	return outputs
}

func mustCopyTree(t *testing.T, src, dst string) {
	t.Helper()
	err := filepath.WalkDir(src, func(path string, entry fs.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		rel, err := filepath.Rel(src, path)
		if err != nil {
			return err
		}
		target := filepath.Join(dst, rel)
		if entry.IsDir() {
			return os.MkdirAll(target, 0o755)
		}
		bytes, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		return os.WriteFile(target, bytes, 0o644)
	})
	if err != nil {
		t.Fatal(err)
	}
}
