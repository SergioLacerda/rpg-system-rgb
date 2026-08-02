package app

import "github.com/SergioLacerda/rpg-system-rgb/internal/components/tooling"

// ReleaseArtifactPaths identifies the public PDF artifact metadata files.
type ReleaseArtifactPaths = tooling.ReleaseArtifactPaths

// WriteReleaseArtifactManifest writes the PDF release manifest and checksums.
func WriteReleaseArtifactManifest(paths ReleaseArtifactPaths) error {
	return tooling.WriteReleaseArtifactManifest(paths)
}

// CheckReleaseArtifacts validates published PDF release artifacts.
func CheckReleaseArtifacts(paths ReleaseArtifactPaths) error {
	return tooling.CheckReleaseArtifacts(paths)
}
