// Package bundles marks the machine-consumable bundle output boundary.
package bundles

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/SergioLacerda/rpg-system-rgb/internal/components"
)

// Descriptor identifies this component to the application layer.
func Descriptor() components.Component {
	return components.Component{
		ID:          "bundles",
		Name:        "RGB Bundles",
		Description: "Machine-consumable bundle output boundary.",
	}
}

const bundleSchema = "rgb-system-bundle/0.1"

// Bundle is the machine-consumable projection of the semantic index units,
// written to generated/bundle/rgb.bundle.json.
type Bundle struct {
	Schema      string       `json:"schema"`
	GeneratedAt string       `json:"generated_at"`
	Units       []BundleUnit `json:"units"`
}

// BundleUnit is one unit's bundle-relevant fields, projected from the
// shared components.SemanticIndexUnit.
type BundleUnit struct {
	ID            string              `json:"id"`
	Kind          string              `json:"kind"`
	Title         string              `json:"title"`
	SourceStatus  string              `json:"source_status"`
	Relationships map[string][]string `json:"relationships,omitempty"`
}

// Build reads the semantic index at repoRoot and projects it into a Bundle.
// It depends only on internal/components (the shared contract package),
// never on a sibling component such as tooling — see
// tests/architecture/architecture_test.go's TestComponentsDoNotImportSiblings.
func Build(repoRoot string) (Bundle, error) {
	indexPath := filepath.Join(repoRoot, "docs", "core", "semantic", "core-v2.index.json")
	index, err := components.LoadSemanticIndex(indexPath)
	if err != nil {
		return Bundle{}, fmt.Errorf("loading semantic index: %w", err)
	}

	units := make([]BundleUnit, 0, len(index.Units))
	for _, unit := range index.Units {
		units = append(units, BundleUnit{
			ID:            unit.ID,
			Kind:          unit.Kind,
			Title:         unit.Title,
			SourceStatus:  unit.SourceStatus,
			Relationships: unit.Relationships,
		})
	}

	return Bundle{
		Schema:      bundleSchema,
		GeneratedAt: time.Now().UTC().Format(time.RFC3339),
		Units:       units,
	}, nil
}

// Write builds the bundle for repoRoot and writes it to
// generated/bundle/rgb.bundle.json.
func Write(repoRoot string) error {
	bundle, err := Build(repoRoot)
	if err != nil {
		return err
	}

	outPath := filepath.Join(repoRoot, "generated", "bundle", "rgb.bundle.json")
	if err := os.MkdirAll(filepath.Dir(outPath), 0o750); err != nil {
		return fmt.Errorf("creating bundle output directory: %w", err)
	}

	data, err := json.MarshalIndent(bundle, "", "  ")
	if err != nil {
		return fmt.Errorf("encoding bundle: %w", err)
	}
	data = append(data, '\n')

	if err := os.WriteFile(outPath, data, 0o644); err != nil { //nolint:gosec // G306: generated artifact, world-readable is expected
		return fmt.Errorf("writing bundle: %w", err)
	}

	fmt.Printf("wrote %d units to generated/bundle/rgb.bundle.json\n", len(bundle.Units))
	return nil
}
