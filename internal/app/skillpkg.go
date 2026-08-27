package app

import "github.com/SergioLacerda/rpg-system-rgb/internal/components/skillpkg"

// SkillOptions configures skill package .zip publication.
type SkillOptions = skillpkg.Options

// SkillManifestPaths identifies the public skill .zip artifact metadata files.
type SkillManifestPaths = skillpkg.ManifestPaths

// PackageSkill zips a skill directory into the public downloads directory.
func PackageSkill(options SkillOptions) error {
	return skillpkg.Package(options)
}

// WriteSkillManifest writes the skill .zip release manifest and checksums.
func WriteSkillManifest(paths SkillManifestPaths) error {
	return skillpkg.WriteManifest(paths)
}

// CheckSkillManifest validates published skill .zip release artifacts.
func CheckSkillManifest(paths SkillManifestPaths) error {
	return skillpkg.CheckManifest(paths)
}
