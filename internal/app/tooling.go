package app

import "github.com/SergioLacerda/rpg-system-rgb/internal/components/tooling"

// ValidateDocs runs every semantic-documentation validator against
// repoRoot, exposing internal/components/tooling to cmd/rgb-tooling
// without letting the entrypoint import the component directly.
func ValidateDocs(repoRoot string) error {
	return tooling.Validate(repoRoot)
}

// GenerateProjections regenerates the derived generated/**.json projection
// artifacts under repoRoot.
func GenerateProjections(repoRoot string) error {
	return tooling.GenerateDefault(repoRoot)
}
