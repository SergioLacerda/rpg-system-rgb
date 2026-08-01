package app

import "github.com/SergioLacerda/rpg-system-rgb/internal/components/bundles"

// BuildBundle writes the machine-consumable bundle
// (generated/bundle/rgb.bundle.json) for repoRoot.
func BuildBundle(repoRoot string) error {
	return bundles.Write(repoRoot)
}
