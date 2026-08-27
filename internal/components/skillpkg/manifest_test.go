package skillpkg

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func writeTestZip(t *testing.T, path, content string) {
	t.Helper()
	writeTestFile(t, path, content)
}

func readTestFile(t *testing.T, path string) string {
	t.Helper()
	content, err := os.ReadFile(path) //nolint:gosec // G304: test-controlled fixture path.
	if err != nil {
		t.Fatalf("reading %s: %v", path, err)
	}
	return string(content)
}

func TestWriteManifestWritesChecksums(t *testing.T) {
	dir := t.TempDir()
	paths := ManifestPaths{
		PublicDir: dir,
		Name:      "rgb-specialist",
		Version:   "v0.1",
		Manifest:  filepath.Join(dir, "manifest.json"),
		Checksums: filepath.Join(dir, "SHA256SUMS"),
	}
	for _, name := range expectedSkillArtifactFiles(paths.Name, paths.Version) {
		writeTestZip(t, filepath.Join(dir, name), strings.Repeat(name, 40))
	}

	if err := WriteManifest(paths); err != nil {
		t.Fatalf("WriteManifest returned error: %v", err)
	}

	var manifest skillManifest
	if err := json.Unmarshal([]byte(readTestFile(t, paths.Manifest)), &manifest); err != nil {
		t.Fatalf("decoding manifest: %v", err)
	}
	if manifest.Schema != manifestSchema {
		t.Fatalf("unexpected schema: %s", manifest.Schema)
	}
	if manifest.Name != paths.Name || manifest.Version != paths.Version {
		t.Fatalf("unexpected manifest identity: %#v", manifest)
	}
	if len(manifest.Artifacts) != 2 {
		t.Fatalf("expected 2 artifacts, got %d", len(manifest.Artifacts))
	}
	checksums := readTestFile(t, paths.Checksums)
	for _, artifact := range manifest.Artifacts {
		if artifact.SHA256 == "" || artifact.Bytes == 0 {
			t.Fatalf("invalid artifact metadata: %#v", artifact)
		}
		if !strings.Contains(checksums, artifact.File) {
			t.Fatalf("checksums file missing %s", artifact.File)
		}
	}
}

func TestWriteManifestRequiresAllPaths(t *testing.T) {
	dir := t.TempDir()
	if err := WriteManifest(ManifestPaths{PublicDir: dir}); err == nil {
		t.Fatal("expected error for incomplete manifest paths")
	}
}

func TestWriteManifestFailsWhenArtifactMissing(t *testing.T) {
	dir := t.TempDir()
	paths := ManifestPaths{
		PublicDir: dir,
		Name:      "rgb-specialist",
		Version:   "v0.1",
		Manifest:  filepath.Join(dir, "manifest.json"),
		Checksums: filepath.Join(dir, "SHA256SUMS"),
	}
	if err := WriteManifest(paths); err == nil {
		t.Fatal("expected error when the skill zip artifacts are missing")
	}
}

func TestCheckManifestAcceptsMatchingArtifacts(t *testing.T) {
	dir := t.TempDir()
	paths := ManifestPaths{
		PublicDir: dir,
		Name:      "rgb-specialist",
		Version:   "v0.1",
		Manifest:  filepath.Join(dir, "manifest.json"),
		Checksums: filepath.Join(dir, "SHA256SUMS"),
	}
	for _, name := range expectedSkillArtifactFiles(paths.Name, paths.Version) {
		writeTestZip(t, filepath.Join(dir, name), strings.Repeat(name, 40))
	}
	if err := WriteManifest(paths); err != nil {
		t.Fatalf("WriteManifest returned error: %v", err)
	}

	if err := CheckManifest(paths); err != nil {
		t.Fatalf("CheckManifest returned error: %v", err)
	}
}

func TestCheckManifestRejectsChecksumMismatch(t *testing.T) {
	dir := t.TempDir()
	paths := ManifestPaths{
		PublicDir: dir,
		Name:      "rgb-specialist",
		Version:   "v0.1",
		Manifest:  filepath.Join(dir, "manifest.json"),
		Checksums: filepath.Join(dir, "SHA256SUMS"),
	}
	for _, name := range expectedSkillArtifactFiles(paths.Name, paths.Version) {
		writeTestZip(t, filepath.Join(dir, name), strings.Repeat(name, 40))
	}
	if err := WriteManifest(paths); err != nil {
		t.Fatalf("WriteManifest returned error: %v", err)
	}

	tamperedPath := filepath.Join(dir, "rgb-specialist-latest.zip")
	writeTestZip(t, tamperedPath, "tampered content")

	if err := CheckManifest(paths); err == nil {
		t.Fatal("expected CheckManifest to reject a tampered artifact")
	}
}

func TestCheckManifestRejectsMissingManifestFile(t *testing.T) {
	dir := t.TempDir()
	paths := ManifestPaths{
		PublicDir: dir,
		Name:      "rgb-specialist",
		Version:   "v0.1",
		Manifest:  filepath.Join(dir, "missing-manifest.json"),
		Checksums: filepath.Join(dir, "SHA256SUMS"),
	}
	if err := CheckManifest(paths); err == nil {
		t.Fatal("expected CheckManifest to fail when the manifest file is missing")
	}
}

func TestCheckManifestRejectsNameAndVersionMismatch(t *testing.T) {
	dir := t.TempDir()
	paths := ManifestPaths{
		PublicDir: dir,
		Name:      "rgb-specialist",
		Version:   "v0.1",
		Manifest:  filepath.Join(dir, "manifest.json"),
		Checksums: filepath.Join(dir, "SHA256SUMS"),
	}
	for _, name := range expectedSkillArtifactFiles(paths.Name, paths.Version) {
		writeTestZip(t, filepath.Join(dir, name), strings.Repeat(name, 40))
	}
	if err := WriteManifest(paths); err != nil {
		t.Fatalf("WriteManifest returned error: %v", err)
	}

	mismatchedName := paths
	mismatchedName.Name = "other-skill"
	if err := CheckManifest(mismatchedName); err == nil {
		t.Fatal("expected CheckManifest to reject a name mismatch")
	}

	mismatchedVersion := paths
	mismatchedVersion.Version = "v0.2"
	if err := CheckManifest(mismatchedVersion); err == nil {
		t.Fatal("expected CheckManifest to reject a version mismatch")
	}
}

func TestCheckManifestRejectsMissingArtifactEntry(t *testing.T) {
	dir := t.TempDir()
	paths := ManifestPaths{
		PublicDir: dir,
		Name:      "rgb-specialist",
		Version:   "v0.1",
		Manifest:  filepath.Join(dir, "manifest.json"),
		Checksums: filepath.Join(dir, "SHA256SUMS"),
	}
	writeTestFile(t, paths.Manifest, `{"schema":"rgb-system-skill-artifacts/0.1","name":"rgb-specialist","version":"v0.1","artifacts":[]}`)
	writeTestFile(t, paths.Checksums, "")

	if err := CheckManifest(paths); err == nil {
		t.Fatal("expected CheckManifest to reject a manifest missing the expected artifacts")
	}
}

func TestValidateManifestArtifactRejectsEmptyChecksumAndByteMismatch(t *testing.T) {
	dir := t.TempDir()
	name := "rgb-specialist-latest.zip"
	writeTestZip(t, filepath.Join(dir, name), "content")

	if err := validateManifestArtifact(dir, name, skillArtifact{SHA256: "", Bytes: 7}); err == nil {
		t.Fatal("expected validateManifestArtifact to reject an empty checksum")
	}
	if err := validateManifestArtifact(dir, name, skillArtifact{SHA256: "deadbeef", Bytes: 999}); err == nil {
		t.Fatal("expected validateManifestArtifact to reject a checksum mismatch")
	}
}

func TestValidateChecksumFileRejectsMissingEmptyAndMalformed(t *testing.T) {
	dir := t.TempDir()
	if err := validateChecksumFile(filepath.Join(dir, "missing"), dir); err == nil {
		t.Fatal("expected validateChecksumFile to fail for a missing file")
	}

	empty := filepath.Join(dir, "empty")
	writeTestFile(t, empty, "")
	if err := validateChecksumFile(empty, dir); err == nil {
		t.Fatal("expected validateChecksumFile to fail for an empty file")
	}

	malformed := filepath.Join(dir, "malformed")
	writeTestFile(t, malformed, "not-a-valid-line")
	if err := validateChecksumFile(malformed, dir); err == nil {
		t.Fatal("expected validateChecksumFile to fail for a malformed line")
	}
}

func TestValidateChecksumLineRejectsMissingTargetFile(t *testing.T) {
	dir := t.TempDir()
	if err := validateChecksumLine("checksums", dir, "deadbeef  missing.zip"); err == nil {
		t.Fatal("expected validateChecksumLine to fail when the referenced file is missing")
	}
}

func TestWriteManifestMetadataFailsForUnwritableManifestPath(t *testing.T) {
	dir := t.TempDir()
	paths := ManifestPaths{
		PublicDir: dir,
		Name:      "rgb-specialist",
		Version:   "v0.1",
		Manifest:  filepath.Join(dir, "missing-dir", "manifest.json"),
		Checksums: filepath.Join(dir, "SHA256SUMS"),
	}
	if err := writeManifestMetadata(paths, skillManifest{Schema: manifestSchema}, ""); err == nil {
		t.Fatal("expected writeManifestMetadata to fail for an unwritable manifest path")
	}
}

func TestReadManifestJSONFailsForMissingAndInvalidContent(t *testing.T) {
	dir := t.TempDir()
	var manifest skillManifest
	if err := readManifestJSON(filepath.Join(dir, "missing.json"), &manifest); err == nil {
		t.Fatal("expected readManifestJSON to fail for a missing file")
	}

	invalid := filepath.Join(dir, "invalid.json")
	writeTestFile(t, invalid, "not json")
	if err := readManifestJSON(invalid, &manifest); err == nil {
		t.Fatal("expected readManifestJSON to fail for invalid JSON")
	}
}

func TestCheckManifestRejectsSchemaMismatch(t *testing.T) {
	dir := t.TempDir()
	paths := ManifestPaths{
		PublicDir: dir,
		Name:      "rgb-specialist",
		Version:   "v0.1",
		Manifest:  filepath.Join(dir, "manifest.json"),
		Checksums: filepath.Join(dir, "SHA256SUMS"),
	}
	for _, name := range expectedSkillArtifactFiles(paths.Name, paths.Version) {
		writeTestZip(t, filepath.Join(dir, name), strings.Repeat(name, 40))
	}
	if err := WriteManifest(paths); err != nil {
		t.Fatalf("WriteManifest returned error: %v", err)
	}
	writeTestFile(t, paths.Manifest, `{"schema":"other/0.1","name":"rgb-specialist","version":"v0.1","artifacts":[]}`)

	if err := CheckManifest(paths); err == nil {
		t.Fatal("expected CheckManifest to reject an unexpected schema")
	}
}
