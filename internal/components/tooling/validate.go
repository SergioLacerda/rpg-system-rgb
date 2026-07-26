package tooling

import (
	"fmt"
	"path/filepath"
)

// Validate runs every semantic-documentation validator in sequence against
// the standard docs/core/semantic/** paths under repoRoot, migrated from
// scripts/validate_semantic_docs.go. The original orchestrated separate
// `go run scripts/X.go` subprocesses; this calls the equivalent exported
// functions directly in the same process.
func Validate(repoRoot string) error {
	semantic := filepath.Join(repoRoot, "docs", "core", "semantic")
	index := filepath.Join(semantic, "core-v2.index.json")
	l10nManifestPath := filepath.Join(semantic, "l10n-manifest.v0.1.json")
	contracts := filepath.Join(semantic, "consumer-contracts.v0.1.json")
	source := filepath.Join(semantic, "source", "core-v2-rules.v0.1.json")
	projections := filepath.Join(semantic, "projection-manifest.v0.1.json")

	steps := []struct {
		name string
		run  func() error
	}{
		{"project paths", func() error { return ValidateProjectPaths(repoRoot) }},
		{"semantic index", func() error { return ValidateSemanticIndex(index) }},
		{"docs l10n manifest", func() error { return ValidateDocsL10nManifest(l10nManifestPath) }},
		{"consumer contracts", func() error { return ValidateConsumerContracts(contracts, index) }},
		{"semantic source", func() error { return ValidateSemanticSource(source, index) }},
		{"projection manifest", func() error { return ValidateProjectionManifest(projections, index, contracts) }},
		{"generated projections", func() error { return ValidateGeneratedProjections(projections) }},
	}

	for _, step := range steps {
		fmt.Printf("validating %s...\n", step.name)
		if err := step.run(); err != nil {
			return fmt.Errorf("semantic docs validation failed at %s: %w", step.name, err)
		}
	}

	fmt.Println("semantic docs validation passed")
	return nil
}
