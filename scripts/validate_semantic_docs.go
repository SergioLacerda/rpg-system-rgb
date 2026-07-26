package main

import (
	"fmt"
	"os"
	"os/exec"
)

type validationCommand struct {
	name string
	args []string
}

func main() {
	commands := []validationCommand{
		{
			name: "project paths",
			args: []string{
				"run",
				"scripts/validate_project_paths.go",
			},
		},
		{
			name: "semantic index",
			args: []string{
				"run",
				"scripts/validate_semantic_index.go",
				"docs/core/semantic/core-v2.index.json",
			},
		},
		{
			name: "docs l10n manifest",
			args: []string{
				"run",
				"scripts/validate_docs_l10n_manifest.go",
				"docs/core/semantic/l10n-manifest.v0.1.json",
			},
		},
		{
			name: "consumer contracts",
			args: []string{
				"run",
				"scripts/validate_semantic_contracts.go",
				"docs/core/semantic/consumer-contracts.v0.1.json",
				"docs/core/semantic/core-v2.index.json",
			},
		},
		{
			name: "semantic source",
			args: []string{
				"run",
				"scripts/validate_semantic_source.go",
				"docs/core/semantic/source/core-v2-rules.v0.1.json",
				"docs/core/semantic/core-v2.index.json",
			},
		},
		{
			name: "projection manifest",
			args: []string{
				"run",
				"scripts/validate_semantic_projections.go",
				"docs/core/semantic/projection-manifest.v0.1.json",
				"docs/core/semantic/core-v2.index.json",
				"docs/core/semantic/consumer-contracts.v0.1.json",
			},
		},
		{
			name: "generated projections",
			args: []string{
				"run",
				"scripts/validate_generated_projections.go",
				"docs/core/semantic/projection-manifest.v0.1.json",
			},
		},
	}

	for _, command := range commands {
		if err := run(command); err != nil {
			fmt.Fprintf(os.Stderr, "semantic docs validation failed at %s: %v\n", command.name, err)
			os.Exit(1)
		}
	}

	fmt.Println("semantic docs validation passed")
}

func run(command validationCommand) error {
	fmt.Printf("validating %s...\n", command.name)
	cmd := exec.Command("go", command.args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
