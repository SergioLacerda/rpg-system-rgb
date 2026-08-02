package tooling

import (
	"encoding/json"
	"image"
	"image/color"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestWriteReleaseArtifactManifestWritesChecksums(t *testing.T) {
	dir := t.TempDir()
	paths := ReleaseArtifactPaths{
		PublicDir: dir,
		Basename:  "rgb-system-core-v2",
		Version:   "v0.2",
		Manifest:  filepath.Join(dir, "manifest.json"),
		Checksums: filepath.Join(dir, "SHA256SUMS"),
	}
	for _, name := range expectedReleaseArtifactFiles(paths.Basename, paths.Version) {
		writeTestPDF(t, filepath.Join(dir, name), strings.Repeat(name, 80))
	}

	if err := WriteReleaseArtifactManifest(paths); err != nil {
		t.Fatalf("WriteReleaseArtifactManifest returned error: %v", err)
	}

	var manifest releaseArtifactManifest
	readJSONFile(t, paths.Manifest, &manifest)
	if manifest.Schema != releaseArtifactSchema {
		t.Fatalf("unexpected schema: %s", manifest.Schema)
	}
	if manifest.Basename != paths.Basename || manifest.Version != paths.Version {
		t.Fatalf("unexpected manifest identity: %#v", manifest)
	}
	if len(manifest.Artifacts) != 4 {
		t.Fatalf("expected 4 artifacts, got %d", len(manifest.Artifacts))
	}
	checksums := readTextFile(t, paths.Checksums)
	for _, artifact := range manifest.Artifacts {
		if artifact.SHA256 == "" || artifact.Bytes < 1000 {
			t.Fatalf("invalid artifact metadata: %#v", artifact)
		}
		if !strings.Contains(checksums, artifact.File) {
			t.Fatalf("SHA256SUMS missing %s", artifact.File)
		}
	}
	for _, locale := range []string{"en", "pt-br"} {
		versionedChecksum := filepath.Join(dir, paths.Basename+"-"+paths.Version+"-"+locale+".pdf.sha256")
		if !strings.Contains(readTextFile(t, versionedChecksum), paths.Basename+"-"+paths.Version+"-"+locale+".pdf") {
			t.Fatalf("versioned checksum missing locale %s", locale)
		}
	}
}

func TestValidateManifestRejectsChecksumMismatch(t *testing.T) {
	dir := t.TempDir()
	paths := ReleaseArtifactPaths{
		PublicDir: dir,
		Basename:  "rgb-system-core-v2",
		Version:   "v0.2",
		Manifest:  filepath.Join(dir, "manifest.json"),
		Checksums: filepath.Join(dir, "SHA256SUMS"),
	}
	for _, name := range expectedReleaseArtifactFiles(paths.Basename, paths.Version) {
		writeTestPDF(t, filepath.Join(dir, name), strings.Repeat(name, 80))
	}
	if err := WriteReleaseArtifactManifest(paths); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, paths.Basename+"-latest-en.pdf"), []byte("%PDF-corrupted"), 0o644); err != nil {
		t.Fatal(err)
	}

	err := validateManifest(paths)
	if err == nil || !strings.Contains(err.Error(), "checksum mismatch") {
		t.Fatalf("expected checksum mismatch, got %v", err)
	}
}

func TestValidateRasterImageAcceptsReadableLightPage(t *testing.T) {
	page := image.NewRGBA(image.Rect(0, 0, 200, 200))
	fill(page, color.RGBA{R: 255, G: 255, B: 255, A: 255})
	for y := 80; y < 130; y++ {
		for x := 50; x < 150; x++ {
			if x%3 == 0 || y%7 == 0 {
				page.Set(x, y, color.RGBA{R: 20, G: 20, B: 20, A: 255})
			}
		}
	}

	if err := validateRasterImage("page.png", page); err != nil {
		t.Fatalf("expected readable light page, got %v", err)
	}
}

func TestValidateRasterImageRejectsBlankPage(t *testing.T) {
	page := image.NewRGBA(image.Rect(0, 0, 200, 200))
	fill(page, color.RGBA{R: 255, G: 255, B: 255, A: 255})

	err := validateRasterImage("blank.png", page)
	if err == nil || !strings.Contains(err.Error(), "blank") {
		t.Fatalf("expected blank-page error, got %v", err)
	}
}

func writeTestPDF(t *testing.T, path, body string) {
	t.Helper()
	content := "%PDF-" + body
	for len(content) < 1200 {
		content += body
	}
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
}

func readJSONFile(t *testing.T, path string, target any) {
	t.Helper()
	bytes, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if err := json.Unmarshal(bytes, target); err != nil {
		t.Fatal(err)
	}
}

func readTextFile(t *testing.T, path string) string {
	t.Helper()
	bytes, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	return string(bytes)
}

func fill(img *image.RGBA, c color.RGBA) {
	for y := img.Bounds().Min.Y; y < img.Bounds().Max.Y; y++ {
		for x := img.Bounds().Min.X; x < img.Bounds().Max.X; x++ {
			img.Set(x, y, c)
		}
	}
}
