package skillpkg

import (
	"archive/zip"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

const manifestSchema = "rgb-system-skill-artifacts/0.1"

// ManifestPaths identifies the public skill .zip artifact and its metadata
// files.
type ManifestPaths struct {
	PublicDir string
	Name      string
	Version   string
	Manifest  string
	Checksums string
}

type skillManifest struct {
	Artifacts []skillArtifact `json:"artifacts"`
	Name      string          `json:"name"`
	Schema    string          `json:"schema"`
	Version   string          `json:"version"`
}

type skillArtifact struct {
	Bytes  int64  `json:"bytes"`
	File   string `json:"file"`
	Path   string `json:"path"`
	SHA256 string `json:"sha256"`
}

// WriteManifest writes the JSON manifest and SHA256SUMS for the expected
// public skill .zip artifacts (<name>-latest.zip and <name>-<version>.zip).
// It computes only SHA256 and byte count — unlike
// internal/components/tooling's PDF release manifest, it performs no
// editorial validation, since a skill archive has none to perform (see
// ADR-013 §2).
func WriteManifest(paths ManifestPaths) error {
	if err := validateManifestPaths(paths); err != nil {
		return err
	}
	if err := os.MkdirAll(paths.PublicDir, 0o750); err != nil {
		return err
	}

	files := expectedSkillArtifactFiles(paths.Name, paths.Version)
	manifest := skillManifest{
		Artifacts: make([]skillArtifact, 0, len(files)),
		Name:      paths.Name,
		Schema:    manifestSchema,
		Version:   paths.Version,
	}
	var checksumLines strings.Builder
	for _, name := range files {
		artifact, line, err := skillArtifactMetadata(paths.PublicDir, name)
		if err != nil {
			return err
		}
		manifest.Artifacts = append(manifest.Artifacts, artifact)
		checksumLines.WriteString(line)
	}
	return writeManifestMetadata(paths, manifest, checksumLines.String())
}

// CheckManifest validates the expected public skill .zip artifacts against
// their manifest and SHA256SUMS — checksum and byte-count only, with no
// PDF-style editorial validation (see ADR-013 §2).
func CheckManifest(paths ManifestPaths) error {
	if err := validateManifestPaths(paths); err != nil {
		return err
	}
	var manifest skillManifest
	if err := readManifestJSON(paths.Manifest, &manifest); err != nil {
		return err
	}
	if err := validateManifestHeader(manifest, paths); err != nil {
		return err
	}
	if err := validateManifestArtifacts(manifest, paths); err != nil {
		return err
	}
	if err := validateChecksumFile(paths.Checksums, paths.PublicDir); err != nil {
		return err
	}
	return validatePackagedSkillPublicContent(paths)
}

func validateManifestArtifacts(manifest skillManifest, paths ManifestPaths) error {
	actual := make(map[string]skillArtifact, len(manifest.Artifacts))
	for _, artifact := range manifest.Artifacts {
		actual[artifact.File] = artifact
	}
	for _, name := range expectedSkillArtifactFiles(paths.Name, paths.Version) {
		artifact, ok := actual[name]
		if !ok {
			return fmt.Errorf("manifest missing artifacts: %s", name)
		}
		if err := validateManifestArtifact(paths.PublicDir, name, artifact); err != nil {
			return err
		}
	}
	return nil
}

func validateManifestHeader(manifest skillManifest, paths ManifestPaths) error {
	if manifest.Schema != manifestSchema {
		return fmt.Errorf("manifest has unexpected schema %s", manifest.Schema)
	}
	if manifest.Name != paths.Name {
		return fmt.Errorf("manifest name mismatch: %s", manifest.Name)
	}
	if manifest.Version != paths.Version {
		return fmt.Errorf("manifest version mismatch: %s", manifest.Version)
	}
	return nil
}

func writeManifestMetadata(paths ManifestPaths, manifest skillManifest, checksumLines string) error {
	data, err := json.MarshalIndent(manifest, "", "  ")
	if err != nil {
		return err
	}
	data = append(data, '\n')
	if err := os.WriteFile(paths.Manifest, data, 0o644); err != nil { //nolint:gosec // G306: manifest is public release metadata.
		return err
	}
	return os.WriteFile(paths.Checksums, []byte(checksumLines), 0o644) //nolint:gosec // G306: checksums are public release metadata.
}

func skillArtifactMetadata(publicDir, name string) (skillArtifact, string, error) {
	path := filepath.Join(publicDir, name)
	content, err := os.ReadFile(path) //nolint:gosec // G304: skill artifact path is provided by the Make target.
	if err != nil {
		return skillArtifact{}, "", fmt.Errorf("missing artifact: %s", path)
	}
	digest := sha256.Sum256(content)
	sha := hex.EncodeToString(digest[:])
	artifact := skillArtifact{
		Bytes:  int64(len(content)),
		File:   name,
		Path:   "/downloads/" + name,
		SHA256: sha,
	}
	return artifact, fmt.Sprintf("%s  %s\n", sha, name), nil
}

func validateManifestArtifact(publicDir, name string, artifact skillArtifact) error {
	if artifact.SHA256 == "" || artifact.Bytes == 0 {
		return fmt.Errorf("invalid artifact metadata: %s", name)
	}
	content, err := os.ReadFile(filepath.Join(publicDir, name)) //nolint:gosec // G304: skill artifact path is provided by the Make target.
	if err != nil {
		return err
	}
	digest := sha256.Sum256(content)
	if hex.EncodeToString(digest[:]) != artifact.SHA256 {
		return fmt.Errorf("manifest checksum mismatch: %s", name)
	}
	if int64(len(content)) != artifact.Bytes {
		return fmt.Errorf("manifest byte count mismatch: %s", name)
	}
	return nil
}

func validateChecksumFile(checksumFile, baseDir string) error {
	content, err := os.ReadFile(checksumFile) //nolint:gosec // G304: checksum path is provided by the Make target.
	if err != nil {
		return fmt.Errorf("missing checksum file: %s", checksumFile)
	}
	lines := strings.Split(strings.TrimSpace(string(content)), "\n")
	if len(lines) == 0 || lines[0] == "" {
		return fmt.Errorf("empty checksum file: %s", checksumFile)
	}
	for _, line := range lines {
		if err := validateChecksumLine(checksumFile, baseDir, line); err != nil {
			return err
		}
	}
	return nil
}

func validateChecksumLine(checksumFile, baseDir, line string) error {
	parts := strings.Fields(line)
	if len(parts) != 2 {
		return fmt.Errorf("invalid checksum line in %s: %q", checksumFile, line)
	}
	content, err := os.ReadFile(filepath.Join(baseDir, parts[1])) //nolint:gosec // G304: checksum target is release metadata controlled by Make.
	if err != nil {
		return err
	}
	digest := sha256.Sum256(content)
	if hex.EncodeToString(digest[:]) != parts[0] {
		return fmt.Errorf("checksum mismatch: %s", parts[1])
	}
	return nil
}

func validatePackagedSkillPublicContent(paths ManifestPaths) error {
	for _, name := range expectedSkillArtifactFiles(paths.Name, paths.Version) {
		if err := validateSkillArchivePublicContent(filepath.Join(paths.PublicDir, name)); err != nil {
			return err
		}
	}
	return nil
}

func validateSkillArchivePublicContent(path string) error {
	reader, err := zip.OpenReader(path)
	if err != nil {
		return err
	}
	defer func() {
		_ = reader.Close()
	}()
	for _, file := range reader.File {
		if !isPublicSkillTextFile(file.Name) {
			continue
		}
		content, err := readZipText(file)
		if err != nil {
			return err
		}
		if containsSkillEngineeringContent(string(content)) {
			return fmt.Errorf("%s contains engineering-only ADR content in %s", path, file.Name)
		}
	}
	return nil
}

func containsSkillEngineeringContent(content string) bool {
	for _, raw := range strings.Split(content, "\n") {
		line := strings.TrimSpace(raw)
		if containsADRReference(line) || packagedSkillEngineeringHeading.MatchString(line) {
			return true
		}
	}
	return false
}

func isPublicSkillTextFile(name string) bool {
	switch strings.ToLower(filepath.Ext(name)) {
	case ".md", ".yaml", ".yml", ".txt":
		return true
	default:
		return false
	}
}

func readZipText(file *zip.File) ([]byte, error) {
	reader, err := file.Open()
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = reader.Close()
	}()
	return io.ReadAll(reader)
}

func expectedSkillArtifactFiles(name, version string) []string {
	return []string{
		fmt.Sprintf("%s-latest.zip", name),
		fmt.Sprintf("%s-%s.zip", name, version),
	}
}

func validateManifestPaths(paths ManifestPaths) error {
	if paths.PublicDir == "" || paths.Name == "" || paths.Version == "" || paths.Manifest == "" || paths.Checksums == "" {
		return errors.New("public-dir, name, version, manifest, and checksums are required")
	}
	return nil
}

func readManifestJSON(path string, target any) error {
	content, err := os.ReadFile(path) //nolint:gosec // G304: manifest path is provided by the Make target.
	if err != nil {
		return err
	}
	return json.Unmarshal(content, target)
}
