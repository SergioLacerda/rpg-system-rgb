package tooling

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"image"
	"image/png"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

const releaseArtifactSchema = "rgb-system-release-artifacts/0.1"

// ReleaseArtifactPaths identifies the PDF artifacts and metadata files used by
// the public release surface.
type ReleaseArtifactPaths struct {
	PublicDir string
	Basename  string
	Version   string
	Manifest  string
	Checksums string
}

type releaseArtifactManifest struct {
	Artifacts []releaseArtifact `json:"artifacts"`
	Basename  string            `json:"basename"`
	Schema    string            `json:"schema"`
	Version   string            `json:"version"`
}

type releaseArtifact struct {
	Bytes  int64  `json:"bytes"`
	File   string `json:"file"`
	Path   string `json:"path"`
	SHA256 string `json:"sha256"`
}

// WriteReleaseArtifactManifest writes the JSON release manifest, SHA256SUMS,
// and per-version .sha256 files for the expected public PDF artifacts.
func WriteReleaseArtifactManifest(paths ReleaseArtifactPaths) error {
	if err := validateReleasePaths(paths); err != nil {
		return err
	}
	if err := os.MkdirAll(paths.PublicDir, 0o750); err != nil {
		return err
	}
	manifest, checksumLines, err := buildReleaseArtifactManifest(paths)
	if err != nil {
		return err
	}
	return writeReleaseArtifactMetadata(paths, manifest, checksumLines)
}

func buildReleaseArtifactManifest(paths ReleaseArtifactPaths) (releaseArtifactManifest, string, error) {
	manifest := releaseArtifactManifest{
		Artifacts: make([]releaseArtifact, 0, len(expectedReleaseArtifactFiles(paths.Basename, paths.Version))),
		Basename:  paths.Basename,
		Schema:    releaseArtifactSchema,
		Version:   paths.Version,
	}
	var checksumLines strings.Builder
	for _, name := range expectedReleaseArtifactFiles(paths.Basename, paths.Version) {
		artifact, line, err := releaseArtifactMetadata(paths.PublicDir, name)
		if err != nil {
			return releaseArtifactManifest{}, "", err
		}
		manifest.Artifacts = append(manifest.Artifacts, artifact)
		checksumLines.WriteString(line)
		if err := writeVersionChecksum(paths.PublicDir, name, paths.Version, line); err != nil {
			return releaseArtifactManifest{}, "", err
		}
	}
	return manifest, checksumLines.String(), nil
}

func writeReleaseArtifactMetadata(paths ReleaseArtifactPaths, manifest releaseArtifactManifest, checksumLines string) error {
	bytes, err := json.MarshalIndent(manifest, "", "  ")
	if err != nil {
		return err
	}
	bytes = append(bytes, '\n')
	if err := os.WriteFile(paths.Manifest, bytes, 0o644); err != nil { //nolint:gosec // G306: manifest is public release metadata.
		return err
	}
	return os.WriteFile(paths.Checksums, []byte(checksumLines), 0o644) //nolint:gosec // G306: checksums are public release metadata.
}

func releaseArtifactMetadata(publicDir, name string) (releaseArtifact, string, error) {
	path := filepath.Join(publicDir, name)
	content, err := os.ReadFile(path) //nolint:gosec // G304: release artifact path is provided by the Make target.
	if err != nil {
		return releaseArtifact{}, "", fmt.Errorf("missing artifact: %s", path)
	}
	digest := sha256.Sum256(content)
	sha := hex.EncodeToString(digest[:])
	artifact := releaseArtifact{
		Bytes:  int64(len(content)),
		File:   name,
		Path:   "/downloads/" + name,
		SHA256: sha,
	}
	return artifact, fmt.Sprintf("%s  %s\n", sha, name), nil
}

func writeVersionChecksum(publicDir, name, version, line string) error {
	if !strings.Contains(name, "-"+version+"-") {
		return nil
	}
	return os.WriteFile(filepath.Join(publicDir, name+".sha256"), []byte(line), 0o644) //nolint:gosec // G306: checksum files are public release metadata.
}

// CheckReleaseArtifacts validates the expected public PDF files, manifest,
// checksums, TOC links, and rasterized pages without depending on Python.
func CheckReleaseArtifacts(paths ReleaseArtifactPaths) error {
	if err := validateReleasePaths(paths); err != nil {
		return err
	}
	if err := validateRequiredPDFTools(); err != nil {
		return err
	}
	if err := validateReleaseArtifactMetadata(paths); err != nil {
		return err
	}
	return validateReleaseArtifactEditorial(paths)
}

func validateReleaseArtifactEditorial(paths ReleaseArtifactPaths) error {
	tmpDir, err := os.MkdirTemp("", "rgb-pdf-check-*")
	if err != nil {
		return err
	}
	defer func() {
		_ = os.RemoveAll(tmpDir)
	}()

	for _, locale := range []string{"en", "pt-br"} {
		if err := validateLocaleEditorialPages(paths, tmpDir, locale); err != nil {
			return err
		}
	}
	return validateRasterizedPages(tmpDir)
}

func validateReleaseArtifactMetadata(paths ReleaseArtifactPaths) error {
	for _, locale := range []string{"en", "pt-br"} {
		if err := validateLocaleReleaseArtifacts(paths, locale); err != nil {
			return err
		}
	}
	if err := validateManifest(paths); err != nil {
		return err
	}
	return validateChecksumFile(paths.Checksums, paths.PublicDir)
}

func validateRequiredPDFTools() error {
	for _, tool := range []string{"pdfinfo", "pdftotext", "pdftohtml", "pdftoppm"} {
		if _, err := exec.LookPath(tool); err != nil {
			return fmt.Errorf("::error::%s is required for PDF editorial checks", tool)
		}
	}
	return nil
}

func validateLocaleReleaseArtifacts(paths ReleaseArtifactPaths, locale string) error {
	for _, channel := range []string{"latest", paths.Version} {
		pdf := filepath.Join(paths.PublicDir, fmt.Sprintf("%s-%s-%s.pdf", paths.Basename, channel, locale))
		if err := validatePDFArtifact(pdf); err != nil {
			return err
		}
	}
	latest := filepath.Join(paths.PublicDir, fmt.Sprintf("%s-latest-%s.pdf", paths.Basename, locale))
	versioned := filepath.Join(paths.PublicDir, fmt.Sprintf("%s-%s-%s.pdf", paths.Basename, paths.Version, locale))
	if err := validateEqualFiles(latest, versioned, locale); err != nil {
		return err
	}
	checksum := filepath.Join(paths.PublicDir, fmt.Sprintf("%s-%s-%s.pdf.sha256", paths.Basename, paths.Version, locale))
	return validateChecksumFile(checksum, paths.PublicDir)
}

func validateLocaleEditorialPages(paths ReleaseArtifactPaths, tmpDir, locale string) error {
	pdf := filepath.Join(paths.PublicDir, fmt.Sprintf("%s-%s-%s.pdf", paths.Basename, paths.Version, locale))
	if err := validatePDFTOC(pdf, filepath.Join(tmpDir, locale+".txt")); err != nil {
		return err
	}
	if err := validatePDFLinks(pdf, filepath.Join(tmpDir, locale+"-links")); err != nil {
		return err
	}
	pages, err := pdfPageCount(pdf)
	if err != nil {
		return err
	}
	return rasterizePDF(pdf, filepath.Join(tmpDir, locale+"-page"), min(pages, 8))
}

func validateReleasePaths(paths ReleaseArtifactPaths) error {
	if paths.PublicDir == "" || paths.Basename == "" || paths.Version == "" || paths.Manifest == "" || paths.Checksums == "" {
		return errors.New("public-dir, basename, version, manifest, and checksums are required")
	}
	return nil
}

func expectedReleaseArtifactFiles(basename, version string) []string {
	return []string{
		fmt.Sprintf("%s-latest-en.pdf", basename),
		fmt.Sprintf("%s-%s-en.pdf", basename, version),
		fmt.Sprintf("%s-latest-pt-br.pdf", basename),
		fmt.Sprintf("%s-%s-pt-br.pdf", basename, version),
	}
}

func validatePDFArtifact(pdf string) error {
	if err := validatePDFHeader(pdf); err != nil {
		return err
	}
	output, err := runCommand("pdfinfo", pdf)
	if err != nil {
		return err
	}
	if err := validatePDFMetadata(pdf, output); err != nil {
		return err
	}
	pages, err := parsePDFPages(output)
	if err != nil {
		return err
	}
	if pages < 3 {
		return fmt.Errorf("::error::%s has too few pages (%d)", pdf, pages)
	}
	return nil
}

func validatePDFHeader(pdf string) error {
	info, err := os.Stat(pdf)
	if err != nil {
		return fmt.Errorf("::error::missing %s", pdf)
	}
	if info.Size() < 1000 {
		return fmt.Errorf("::error::%s is suspiciously small (%d bytes)", pdf, info.Size())
	}
	content, err := os.ReadFile(pdf) //nolint:gosec // G304: release artifact path is provided by the Make target.
	if err != nil {
		return err
	}
	if len(content) < 5 || string(content[:5]) != "%PDF-" {
		return fmt.Errorf("::error::%s does not start with a PDF header", pdf)
	}
	return nil
}

func validatePDFMetadata(pdf string, output []byte) error {
	for _, field := range []string{"Title:", "Subject:", "Producer:"} {
		if !regexp.MustCompile(`(?m)^` + regexp.QuoteMeta(field)).Match(output) {
			name := strings.TrimSuffix(strings.ToLower(field), ":")
			return fmt.Errorf("::error::%s is missing PDF %s metadata", pdf, name)
		}
	}
	return nil
}

func validateEqualFiles(latest, versioned, locale string) error {
	left, err := os.ReadFile(latest) //nolint:gosec // G304: release artifact path is provided by the Make target.
	if err != nil {
		return err
	}
	right, err := os.ReadFile(versioned) //nolint:gosec // G304: release artifact path is provided by the Make target.
	if err != nil {
		return err
	}
	if !bytes.Equal(left, right) {
		return fmt.Errorf("::error::latest and versioned PDFs differ for %s", locale)
	}
	return nil
}

func validateManifest(paths ReleaseArtifactPaths) error {
	var manifest releaseArtifactManifest
	if err := readGeneratedJSON(paths.Manifest, &manifest); err != nil {
		return err
	}
	if err := validateManifestHeader(manifest, paths); err != nil {
		return err
	}

	expected := stringSet(expectedReleaseArtifactFiles(paths.Basename, paths.Version))
	actual := map[string]releaseArtifact{}
	for _, artifact := range manifest.Artifacts {
		actual[artifact.File] = artifact
	}
	for name := range expected {
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

func validateManifestHeader(manifest releaseArtifactManifest, paths ReleaseArtifactPaths) error {
	if manifest.Schema != releaseArtifactSchema {
		return fmt.Errorf("manifest has unexpected schema %s", manifest.Schema)
	}
	if manifest.Basename != paths.Basename {
		return fmt.Errorf("manifest basename mismatch: %s", manifest.Basename)
	}
	if manifest.Version != paths.Version {
		return fmt.Errorf("manifest version mismatch: %s", manifest.Version)
	}
	return nil
}

func validateManifestArtifact(publicDir, name string, artifact releaseArtifact) error {
	if artifact.SHA256 == "" || artifact.Bytes < 1000 {
		return fmt.Errorf("invalid artifact metadata: %s", name)
	}
	content, err := os.ReadFile(filepath.Join(publicDir, name)) //nolint:gosec // G304: release artifact path is provided by the Make target.
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
		return fmt.Errorf("::error::missing checksum file: %s", checksumFile)
	}
	lines := strings.Split(strings.TrimSpace(string(content)), "\n")
	if len(lines) == 0 || lines[0] == "" {
		return fmt.Errorf("::error::empty checksum file: %s", checksumFile)
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
	path := filepath.Join(baseDir, parts[1])
	bytes, err := os.ReadFile(path) //nolint:gosec // G304: checksum target is release metadata controlled by Make.
	if err != nil {
		return err
	}
	digest := sha256.Sum256(bytes)
	if hex.EncodeToString(digest[:]) != parts[0] {
		return fmt.Errorf("checksum mismatch: %s", parts[1])
	}
	return nil
}

func validatePDFTOC(pdf, textPath string) error {
	if _, err := runCommand("pdftotext", "-layout", "-f", "1", "-l", "3", pdf, textPath); err != nil {
		return err
	}
	text, err := os.ReadFile(textPath) //nolint:gosec // G304: temporary path created by this process.
	if err != nil {
		return err
	}
	if regexp.MustCompile(`(^|\s)0\s*$`).Match(text) {
		return fmt.Errorf("::error::%s contains a page-zero table-of-contents entry", pdf)
	}
	return nil
}

func validatePDFLinks(pdf, outputPrefix string) error {
	if _, err := runCommand("pdftohtml", "-xml", "-i", "-f", "1", "-l", "3", pdf, outputPrefix); err != nil {
		return err
	}
	content, err := os.ReadFile(outputPrefix + ".xml") //nolint:gosec // G304: temporary path created by this process.
	if err != nil {
		return err
	}
	if !bytes.Contains(content, []byte("<a href=")) {
		return fmt.Errorf("::error::%s has no extractable TOC links on the critical pages", pdf)
	}
	return nil
}

func rasterizePDF(pdf, outputPrefix string, firstContentPage int) error {
	_, err := runCommand("pdftoppm", "-r", "72", "-f", "1", "-l", strconv.Itoa(firstContentPage), "-png", pdf, outputPrefix)
	return err
}

func validateRasterizedPages(root string) error {
	pngs, err := filepath.Glob(filepath.Join(root, "*-page-*.png"))
	if err != nil {
		return err
	}
	if len(pngs) == 0 {
		return errors.New("no rasterized PDF pages found")
	}
	for _, path := range pngs {
		if err := validateRasterPNG(path); err != nil {
			return err
		}
	}
	return nil
}

func validateRasterPNG(path string) error {
	file, err := os.Open(path) //nolint:gosec // G304: temporary raster path created by this process.
	if err != nil {
		return err
	}
	defer func() {
		_ = file.Close()
	}()
	imageData, err := png.Decode(file)
	if err != nil {
		return err
	}
	if err := validateRasterImage(filepath.Base(path), imageData); err != nil {
		return err
	}
	return nil
}

func validateRasterImage(name string, imageData image.Image) error {
	bounds := imageData.Bounds()
	width := bounds.Dx()
	height := bounds.Dy()
	if width <= 0 || height <= 0 {
		return fmt.Errorf("%s appears blank or missing readable text", name)
	}

	avg, readableRatio, veryLightRatio := rasterLuminanceStats(imageData)
	if avg < 180 || veryLightRatio < 0.55 {
		return fmt.Errorf("%s appears too dark for editorial PDF output", name)
	}
	if readableRatio < 0.003 {
		return fmt.Errorf("%s appears blank or missing readable text", name)
	}

	margin := max(4, min(width, height)*2/100)
	if rasterEdgeDarkRatio(imageData, margin) > 0.015 {
		return fmt.Errorf("%s has excessive dark marks near page edges", name)
	}
	return nil
}

func rasterLuminanceStats(imageData image.Image) (float64, float64, float64) {
	bounds := imageData.Bounds()
	var total, readableMarks, veryLight int64
	var lumaTotal float64
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			luma := luminance(imageData.At(x, y).RGBA())
			lumaTotal += luma
			total++
			if luma < 180 {
				readableMarks++
			}
			if luma > 210 {
				veryLight++
			}
		}
	}
	return lumaTotal / float64(total), float64(readableMarks) / float64(total), float64(veryLight) / float64(total)
}

func rasterEdgeDarkRatio(imageData image.Image, margin int) float64 {
	horizontalTotal, horizontalDark := rasterHorizontalEdgeStats(imageData, margin)
	verticalTotal, verticalDark := rasterVerticalEdgeStats(imageData, margin)
	edgeTotal := horizontalTotal + verticalTotal
	edgeDark := horizontalDark + verticalDark
	return float64(edgeDark) / float64(edgeTotal)
}

func rasterHorizontalEdgeStats(imageData image.Image, margin int) (int64, int64) {
	bounds := imageData.Bounds()
	var total, dark int64
	for x := bounds.Min.X; x < bounds.Max.X; x++ {
		for y := 0; y < margin; y++ {
			if luminance(imageData.At(x, bounds.Min.Y+y).RGBA()) < 120 {
				dark++
			}
			if luminance(imageData.At(x, bounds.Max.Y-1-y).RGBA()) < 120 {
				dark++
			}
			total += 2
		}
	}
	return total, dark
}

func rasterVerticalEdgeStats(imageData image.Image, margin int) (int64, int64) {
	bounds := imageData.Bounds()
	var total, dark int64
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := 0; x < margin; x++ {
			if luminance(imageData.At(bounds.Min.X+x, y).RGBA()) < 120 {
				dark++
			}
			if luminance(imageData.At(bounds.Max.X-1-x, y).RGBA()) < 120 {
				dark++
			}
			total += 2
		}
	}
	return total, dark
}

func luminance(r, g, b, _ uint32) float64 {
	rr := float64(r >> 8)
	gg := float64(g >> 8)
	bb := float64(b >> 8)
	return (rr*299 + gg*587 + bb*114) / 1000
}

func pdfPageCount(pdf string) (int, error) {
	output, err := runCommand("pdfinfo", pdf)
	if err != nil {
		return 0, err
	}
	return parsePDFPages(output)
}

func parsePDFPages(output []byte) (int, error) {
	for _, line := range strings.Split(string(output), "\n") {
		if strings.HasPrefix(line, "Pages:") {
			value := strings.TrimSpace(strings.TrimPrefix(line, "Pages:"))
			pages, err := strconv.Atoi(value)
			if err != nil {
				return 0, fmt.Errorf("invalid PDF page count %q: %w", value, err)
			}
			return pages, nil
		}
	}
	return 0, errors.New("missing PDF page count")
}

func runCommand(name string, args ...string) ([]byte, error) {
	cmd := exec.Command(name, args...) //nolint:gosec // G204: command names are fixed by caller; args are file paths from Make targets.
	output, err := cmd.CombinedOutput()
	if err != nil {
		return output, fmt.Errorf("%s failed: %w\n%s", name, err, string(output))
	}
	return output, nil
}
