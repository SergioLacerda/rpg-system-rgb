package tooling

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestWriteReleaseArtifactManifestErrorStages(t *testing.T) {
	if err := WriteReleaseArtifactManifest(ReleaseArtifactPaths{}); err == nil {
		t.Fatal("expected invalid release paths to fail")
	}

	dir := t.TempDir()
	publicFile := filepath.Join(dir, "public")
	if err := os.WriteFile(publicFile, []byte("not a dir"), 0o640); err != nil {
		t.Fatal(err)
	}
	if err := WriteReleaseArtifactManifest(ReleaseArtifactPaths{PublicDir: publicFile, Basename: "rgb", Version: "v1", Manifest: filepath.Join(dir, "manifest.json"), Checksums: filepath.Join(dir, "SHA256SUMS")}); err == nil {
		t.Fatal("expected public-dir mkdir failure")
	}

	publicDir := filepath.Join(dir, "downloads")
	paths := ReleaseArtifactPaths{PublicDir: publicDir, Basename: "rgb", Version: "v1", Manifest: filepath.Join(publicDir, "manifest.json"), Checksums: filepath.Join(publicDir, "SHA256SUMS")}
	if err := os.MkdirAll(publicDir, 0o750); err != nil {
		t.Fatal(err)
	}
	for _, name := range expectedReleaseArtifactFiles(paths.Basename, paths.Version) {
		writeTestPDF(t, filepath.Join(publicDir, name), strings.Repeat(name, 80))
	}
	if err := os.Mkdir(filepath.Join(publicDir, "rgb-v1-en.pdf.sha256"), 0o750); err != nil {
		t.Fatal(err)
	}
	if _, _, err := buildReleaseArtifactManifest(paths); err == nil {
		t.Fatal("expected version checksum write failure")
	}
}

func TestReleasePathAndMetadataValidationEdges(t *testing.T) {
	if err := validateReleasePaths(ReleaseArtifactPaths{}); err == nil {
		t.Fatal("expected empty release paths to fail")
	}
	dir := t.TempDir()
	paths := ReleaseArtifactPaths{PublicDir: dir, Basename: "rgb", Version: "v1", Manifest: filepath.Join(dir, "manifest.json"), Checksums: filepath.Join(dir, "SHA256SUMS")}
	if err := validateReleasePaths(paths); err != nil {
		t.Fatal(err)
	}
	if _, _, err := releaseArtifactMetadata(dir, "missing.pdf"); err == nil {
		t.Fatal("expected missing release artifact metadata to fail")
	}
	if _, _, err := buildReleaseArtifactManifest(paths); err == nil {
		t.Fatal("expected missing release artifacts to fail manifest build")
	}
	if err := writeReleaseArtifactMetadata(ReleaseArtifactPaths{Manifest: filepath.Join(dir, "missing", "manifest.json")}, releaseArtifactManifest{}, ""); err == nil {
		t.Fatal("expected manifest write error")
	}
	if err := writeReleaseArtifactMetadata(ReleaseArtifactPaths{Manifest: filepath.Join(dir, "manifest.json"), Checksums: filepath.Join(dir, "missing", "SHA256SUMS")}, releaseArtifactManifest{}, ""); err == nil {
		t.Fatal("expected checksums write error")
	}
}

func TestReleasePDFHeaderMetadataAndPageParsing(t *testing.T) {
	dir := t.TempDir()
	missing := filepath.Join(dir, "missing.pdf")
	if err := validatePDFHeader(missing); err == nil {
		t.Fatal("expected missing PDF to fail")
	}
	small := filepath.Join(dir, "small.pdf")
	if err := os.WriteFile(small, []byte("%PDF-"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := validatePDFHeader(small); err == nil {
		t.Fatal("expected small PDF to fail")
	}
	notPDF := filepath.Join(dir, "not.pdf")
	if err := os.WriteFile(notPDF, []byte(strings.Repeat("x", 1200)), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := validatePDFHeader(notPDF); err == nil {
		t.Fatal("expected missing PDF header to fail")
	}
	good := filepath.Join(dir, "good.pdf")
	writeTestPDF(t, good, strings.Repeat("good", 300))
	if err := validatePDFHeader(good); err != nil {
		t.Fatal(err)
	}

	metadata := []byte("Title: RGB\nSubject: Core\nProducer: Test\nPages: 4\n")
	if err := validatePDFMetadata(good, metadata); err != nil {
		t.Fatal(err)
	}
	for _, output := range [][]byte{
		[]byte("Subject: Core\nProducer: Test\n"),
		[]byte("Title: RGB\nProducer: Test\n"),
		[]byte("Title: RGB\nSubject: Core\n"),
	} {
		if err := validatePDFMetadata(good, output); err == nil {
			t.Fatalf("expected incomplete metadata to fail: %s", output)
		}
	}
	pages, err := parsePDFPages(metadata)
	if err != nil {
		t.Fatal(err)
	}
	if pages != 4 {
		t.Fatalf("pages got %d want 4", pages)
	}
	if _, err := parsePDFPages([]byte("Pages: nope\n")); err == nil {
		t.Fatal("expected invalid page count to fail")
	}
	if _, err := parsePDFPages([]byte("Title: RGB\n")); err == nil {
		t.Fatal("expected missing page count to fail")
	}
}

func TestReleaseChecksumAndManifestValidationEdges(t *testing.T) {
	dir := t.TempDir()
	pdf := filepath.Join(dir, "artifact.pdf")
	writeTestPDF(t, pdf, strings.Repeat("artifact", 200))
	content, err := os.ReadFile(pdf)
	if err != nil {
		t.Fatal(err)
	}
	sum := sha256.Sum256(content)
	sha := hex.EncodeToString(sum[:])
	checksum := filepath.Join(dir, "SHA256SUMS")
	if err := os.WriteFile(checksum, []byte(sha+"  artifact.pdf\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := validateChecksumFile(checksum, dir); err != nil {
		t.Fatal(err)
	}
	if err := validateChecksumLine(checksum, dir, "not-enough-fields"); err == nil {
		t.Fatal("expected malformed checksum line to fail")
	}
	if err := validateChecksumLine(checksum, dir, sha+"  missing.pdf"); err == nil {
		t.Fatal("expected missing checksum target to fail")
	}
	if err := validateChecksumLine(checksum, dir, strings.Repeat("0", 64)+"  artifact.pdf"); err == nil {
		t.Fatal("expected checksum mismatch to fail")
	}
	empty := filepath.Join(dir, "empty")
	if err := os.WriteFile(empty, nil, 0o644); err != nil {
		t.Fatal(err)
	}
	if err := validateChecksumFile(empty, dir); err == nil {
		t.Fatal("expected empty checksum file to fail")
	}
	if err := validateChecksumFile(filepath.Join(dir, "missing"), dir); err == nil {
		t.Fatal("expected missing checksum file to fail")
	}

	paths := ReleaseArtifactPaths{PublicDir: dir, Basename: "rgb", Version: "v1", Manifest: filepath.Join(dir, "manifest.json"), Checksums: checksum}
	if err := validateManifestHeader(releaseArtifactManifest{Schema: "wrong", Basename: "rgb", Version: "v1"}, paths); err == nil {
		t.Fatal("expected manifest schema mismatch")
	}
	if err := validateManifestHeader(releaseArtifactManifest{Schema: releaseArtifactSchema, Basename: "wrong", Version: "v1"}, paths); err == nil {
		t.Fatal("expected manifest basename mismatch")
	}
	if err := validateManifestHeader(releaseArtifactManifest{Schema: releaseArtifactSchema, Basename: "rgb", Version: "wrong"}, paths); err == nil {
		t.Fatal("expected manifest version mismatch")
	}
	if err := validateManifestArtifact(dir, "artifact.pdf", releaseArtifact{File: "artifact.pdf"}); err == nil {
		t.Fatal("expected invalid manifest artifact metadata")
	}
	if err := validateManifestArtifact(dir, "artifact.pdf", releaseArtifact{File: "artifact.pdf", Bytes: int64(len(content)), SHA256: sha}); err != nil {
		t.Fatal(err)
	}
	if err := validateManifestArtifact(dir, "artifact.pdf", releaseArtifact{File: "artifact.pdf", Bytes: int64(len(content)) + 1, SHA256: sha}); err == nil {
		t.Fatal("expected byte count mismatch")
	}
}

func TestReleaseFileEqualityAndRasterValidationEdges(t *testing.T) {
	dir := t.TempDir()
	latest := filepath.Join(dir, "latest.pdf")
	versioned := filepath.Join(dir, "versioned.pdf")
	writeTestPDF(t, latest, strings.Repeat("same", 300))
	writeTestPDF(t, versioned, strings.Repeat("same", 300))
	if err := validateEqualFiles(latest, versioned, "en"); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(versioned, []byte("%PDF-different"+strings.Repeat("x", 1200)), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := validateEqualFiles(latest, versioned, "en"); err == nil {
		t.Fatal("expected unequal files to fail")
	}
	if err := validateEqualFiles(filepath.Join(dir, "missing.pdf"), versioned, "en"); err == nil {
		t.Fatal("expected missing latest file to fail")
	}

	if err := validateRasterizedPages(dir); err == nil {
		t.Fatal("expected no rasterized pages to fail")
	}
	pngPath := filepath.Join(dir, "en-page-1.png")
	page := image.NewRGBA(image.Rect(0, 0, 200, 200))
	fill(page, color.RGBA{R: 255, G: 255, B: 255, A: 255})
	for y := 70; y < 130; y++ {
		for x := 50; x < 150; x += 4 {
			page.Set(x, y, color.RGBA{R: 20, G: 20, B: 20, A: 255})
		}
	}
	file, err := os.Create(pngPath)
	if err != nil {
		t.Fatal(err)
	}
	if err := png.Encode(file, page); err != nil {
		t.Fatal(err)
	}
	if err := file.Close(); err != nil {
		t.Fatal(err)
	}
	if err := validateRasterPNG(pngPath); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "bad-page-1.png"), []byte("not png"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := validateRasterPNG(filepath.Join(dir, "bad-page-1.png")); err == nil {
		t.Fatal("expected invalid PNG to fail")
	}
	if err := validateRasterPNG(filepath.Join(dir, "missing-page-1.png")); err == nil {
		t.Fatal("expected missing PNG to fail")
	}
}

func TestCheckReleaseArtifactsWithFakePDFTools(t *testing.T) {
	dir := t.TempDir()
	paths := ReleaseArtifactPaths{
		PublicDir: dir,
		Basename:  "rgb",
		Version:   "v1",
		Manifest:  filepath.Join(dir, "manifest.json"),
		Checksums: filepath.Join(dir, "SHA256SUMS"),
	}
	for _, locale := range []string{"en", "pt-br"} {
		body := strings.Repeat(locale, 600)
		writeTestPDF(t, filepath.Join(dir, fmt.Sprintf("%s-latest-%s.pdf", paths.Basename, locale)), body)
		writeTestPDF(t, filepath.Join(dir, fmt.Sprintf("%s-%s-%s.pdf", paths.Basename, paths.Version, locale)), body)
	}
	if err := WriteReleaseArtifactManifest(paths); err != nil {
		t.Fatal(err)
	}

	originalRunner := runExternalCommand
	t.Cleanup(func() { runExternalCommand = originalRunner })
	runExternalCommand = func(name string, args ...string) ([]byte, error) {
		switch name {
		case "pdfinfo":
			return []byte("Title: RGB\nSubject: Core\nProducer: Test\nPages: 4\n"), nil
		case "pdftotext":
			if len(args) == 0 {
				return nil, fmt.Errorf("missing pdftotext args")
			}
			return nil, os.WriteFile(args[len(args)-1], []byte("Table of contents\n1\n2\n3\n"), 0o644)
		case "pdftohtml":
			if len(args) == 0 {
				return nil, fmt.Errorf("missing pdftohtml args")
			}
			return nil, os.WriteFile(args[len(args)-1]+".xml", []byte(`<page><a href="#section">Section</a></page>`), 0o644)
		case "pdftoppm":
			if len(args) < 2 {
				return nil, fmt.Errorf("missing pdftoppm args")
			}
			prefix := args[len(args)-1]
			for i := 1; i <= 4; i++ {
				writeReadablePNG(t, fmt.Sprintf("%s-%d.png", prefix, i))
			}
			return nil, nil
		default:
			return nil, fmt.Errorf("unexpected command %s", name)
		}
	}

	if err := CheckReleaseArtifacts(paths); err != nil {
		t.Fatal(err)
	}
}

func TestCheckReleaseArtifactsFailureStages(t *testing.T) {
	if err := CheckReleaseArtifacts(ReleaseArtifactPaths{}); err == nil {
		t.Fatal("expected invalid paths to fail before tool lookup")
	}
	originalLookPath := lookPath
	t.Cleanup(func() { lookPath = originalLookPath })
	lookPath = func(name string) (string, error) {
		return "", fmt.Errorf("missing %s", name)
	}
	if err := CheckReleaseArtifacts(ReleaseArtifactPaths{PublicDir: t.TempDir(), Basename: "rgb", Version: "v1", Manifest: "manifest.json", Checksums: "SHA256SUMS"}); err == nil {
		t.Fatal("expected missing PDF tool to fail")
	}
	lookPath = originalLookPath

	dir := t.TempDir()
	paths := ReleaseArtifactPaths{
		PublicDir: dir,
		Basename:  "rgb",
		Version:   "v1",
		Manifest:  filepath.Join(dir, "manifest.json"),
		Checksums: filepath.Join(dir, "SHA256SUMS"),
	}
	originalRunner := runExternalCommand
	t.Cleanup(func() { runExternalCommand = originalRunner })
	runExternalCommand = func(name string, _ ...string) ([]byte, error) {
		switch name {
		case "pdfinfo":
			return []byte("Title: RGB\nSubject: Core\nProducer: Test\nPages: 4\n"), nil
		default:
			return nil, fmt.Errorf("unexpected command %s", name)
		}
	}
	if err := validateReleaseArtifactMetadata(paths); err == nil {
		t.Fatal("expected missing release metadata to fail")
	}

	for _, locale := range []string{"en", "pt-br"} {
		body := strings.Repeat(locale, 600)
		writeTestPDF(t, filepath.Join(dir, fmt.Sprintf("%s-latest-%s.pdf", paths.Basename, locale)), body)
		writeTestPDF(t, filepath.Join(dir, fmt.Sprintf("%s-%s-%s.pdf", paths.Basename, paths.Version, locale)), body)
	}
	if err := WriteReleaseArtifactManifest(paths); err != nil {
		t.Fatal(err)
	}
	if err := validateReleaseArtifactMetadata(paths); err != nil {
		t.Fatal(err)
	}

	runExternalCommand = func(name string, args ...string) ([]byte, error) {
		switch name {
		case "pdftotext":
			return nil, os.WriteFile(args[len(args)-1], []byte("1\n2\n3\n"), 0o644)
		case "pdftohtml":
			return nil, os.WriteFile(args[len(args)-1]+".xml", []byte(`<page><a href="#x">x</a></page>`), 0o644)
		case "pdfinfo":
			return []byte("Title: RGB\nSubject: Core\nProducer: Test\nPages: 1\n"), nil
		case "pdftoppm":
			return nil, nil
		default:
			return nil, fmt.Errorf("unexpected command %s", name)
		}
	}
	if err := validateReleaseArtifactEditorial(paths); err == nil {
		t.Fatal("expected missing rasterized pages to fail")
	}
	if err := validatePDFArtifact(filepath.Join(dir, fmt.Sprintf("%s-%s-en.pdf", paths.Basename, paths.Version))); err == nil {
		t.Fatal("expected too few PDF pages to fail")
	}
}

func TestEditorialHelpersRejectBadFakeToolOutput(t *testing.T) {
	originalRunner := runExternalCommand
	t.Cleanup(func() { runExternalCommand = originalRunner })
	runExternalCommand = func(name string, args ...string) ([]byte, error) {
		switch name {
		case "pdftotext":
			return nil, os.WriteFile(args[len(args)-1], []byte("bad entry 0\n"), 0o644)
		case "pdftohtml":
			return nil, os.WriteFile(args[len(args)-1]+".xml", []byte(`<page></page>`), 0o644)
		case "pdfinfo":
			return []byte("Pages: 2\n"), nil
		case "pdftoppm":
			return []byte("boom"), fmt.Errorf("boom")
		default:
			return nil, fmt.Errorf("unexpected command %s", name)
		}
	}
	dir := t.TempDir()
	pdf := filepath.Join(dir, "doc.pdf")
	writeTestPDF(t, pdf, strings.Repeat("doc", 400))
	if err := validatePDFTOC(pdf, filepath.Join(dir, "toc.txt")); err == nil {
		t.Fatal("expected page-zero TOC to fail")
	}
	if err := validatePDFLinks(pdf, filepath.Join(dir, "links")); err == nil {
		t.Fatal("expected missing links to fail")
	}
	if err := rasterizePDF(pdf, filepath.Join(dir, "page"), 2); err == nil {
		t.Fatal("expected rasterize command failure")
	}
	if _, err := pdfPageCount(pdf); err != nil {
		t.Fatal(err)
	}
	if _, err := runCommand("pdftoppm"); err == nil {
		t.Fatal("expected runCommand to wrap command failure")
	}
}

func writeReadablePNG(t *testing.T, path string) {
	t.Helper()
	page := image.NewRGBA(image.Rect(0, 0, 200, 200))
	fill(page, color.RGBA{R: 255, G: 255, B: 255, A: 255})
	for y := 70; y < 130; y++ {
		for x := 50; x < 150; x += 4 {
			page.Set(x, y, color.RGBA{R: 20, G: 20, B: 20, A: 255})
		}
	}
	file, err := os.Create(path)
	if err != nil {
		t.Fatal(err)
	}
	if err := png.Encode(file, page); err != nil {
		t.Fatal(err)
	}
	if err := file.Close(); err != nil {
		t.Fatal(err)
	}
}
