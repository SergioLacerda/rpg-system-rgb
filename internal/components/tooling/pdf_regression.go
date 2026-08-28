package tooling

import (
	"errors"
	"fmt"
	"image"
	"image/png"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

// pdfRegressionRasterTolerance bounds the fraction of differing pixels
// allowed per compared page before CheckPDFRegression fails. High enough to
// absorb anti-aliasing jitter between two independent renders of the same
// PDF, low enough to catch a moved cover glyph or a recolored callout box.
const pdfRegressionRasterTolerance = 0.01

// pdfRegressionMaxComparedPages bounds worst-case runtime: rasterizing and
// diffing every page of a hundred-page book is unnecessary to catch an
// authoring regression, since cover/TOC/callout issues concentrate in the
// early pages.
const pdfRegressionMaxComparedPages = 12

// pdfRegressionPixelChannelThreshold is the per-channel (0-255) intensity
// difference above which a pixel counts as "differing" in the rasterized
// page diff — high enough to ignore anti-aliasing noise, low enough to
// catch an actual color or layout change.
const pdfRegressionPixelChannelThreshold = 24

// CheckPDFRegression compares a candidate PDF against a baseline PDF across
// four checks — page count, extracted text, hyperlink count on the first
// three pages, and rasterized page diff — and returns an error describing
// the first regression found. Calling it with the same path for baseline
// and candidate is a valid, expected use: it proves the harness itself is
// correct (all four checks pass with a zero diff) before it is trusted to
// judge a genuinely different renderer's output.
func CheckPDFRegression(baseline, candidate string) error {
	if err := validateRequiredPDFTools(); err != nil {
		return err
	}
	tmpDir, err := os.MkdirTemp("", "rgb-pdf-regression-*")
	if err != nil {
		return err
	}
	defer func() {
		_ = os.RemoveAll(tmpDir)
	}()

	if err := checkPDFRegressionPageCount(baseline, candidate); err != nil {
		return err
	}
	if err := checkPDFRegressionText(baseline, candidate, tmpDir); err != nil {
		return err
	}
	if err := checkPDFRegressionLinks(baseline, candidate, tmpDir); err != nil {
		return err
	}
	return checkPDFRegressionRaster(baseline, candidate, tmpDir)
}

// CheckPDFRegressionForRelease runs CheckPDFRegression for both release
// locales (en, pt-br), comparing each versioned PDF under baselineDir
// against its counterpart under candidateDir. Passing the same directory
// for both compares the current release against itself.
func CheckPDFRegressionForRelease(baselineDir, candidateDir, basename, version string) error {
	for _, locale := range []string{"en", "pt-br"} {
		name := fmt.Sprintf("%s-%s-%s.pdf", basename, version, locale)
		baseline := filepath.Join(baselineDir, name)
		candidate := filepath.Join(candidateDir, name)
		if err := CheckPDFRegression(baseline, candidate); err != nil {
			return fmt.Errorf("pdf regression (%s): %w", locale, err)
		}
	}
	return nil
}

func checkPDFRegressionPageCount(baseline, candidate string) error {
	basePages, err := pdfPageCount(baseline)
	if err != nil {
		return err
	}
	candPages, err := pdfPageCount(candidate)
	if err != nil {
		return err
	}
	if basePages != candPages {
		return fmt.Errorf("::error::page count regression: baseline=%d candidate=%d (%s vs %s)", basePages, candPages, baseline, candidate)
	}
	return nil
}

// pdfRegressionStampedLine matches lines the cover template stamps at build
// time (tools/pdfbuild/cover.html: {{ date }}, {{ commit_sha }}, and the
// colophon's "Source" line combining both) — an ISO date on its own line, a
// bare short/long git SHA on its own line, or a line ending in "@ <sha>".
var pdfRegressionStampedLine = regexp.MustCompile(`(?i)^(\d{4}-\d{2}-\d{2}|[0-9a-f]{7,40})$|@\s*[0-9a-f]{7,40}\s*$`)

func checkPDFRegressionText(baseline, candidate, tmpDir string) error {
	baseText, err := extractPDFText(baseline, filepath.Join(tmpDir, "baseline.txt"))
	if err != nil {
		return err
	}
	candText, err := extractPDFText(candidate, filepath.Join(tmpDir, "candidate.txt"))
	if err != nil {
		return err
	}
	baseLines := normalizePDFRegressionText(baseText)
	candLines := normalizePDFRegressionText(candText)
	if diff := firstTextLineDiff(baseLines, candLines); diff != "" {
		return fmt.Errorf("::error::extracted text regression: %s (%s vs %s)", diff, baseline, candidate)
	}
	return nil
}

func extractPDFText(pdf, outPath string) (string, error) {
	if _, err := runCommand("pdftotext", "-layout", pdf, outPath); err != nil {
		return "", err
	}
	content, err := os.ReadFile(outPath) //nolint:gosec // G304: temporary path created by this process.
	if err != nil {
		return "", err
	}
	return string(content), nil
}

func normalizePDFRegressionText(content string) []string {
	rawLines := strings.Split(content, "\n")
	lines := make([]string, 0, len(rawLines))
	for _, raw := range rawLines {
		trimmed := strings.TrimRight(raw, " \t\r")
		if pdfRegressionStampedLine.MatchString(strings.TrimSpace(trimmed)) {
			continue
		}
		lines = append(lines, trimmed)
	}
	return lines
}

func firstTextLineDiff(base, cand []string) string {
	for i := 0; i < len(base) && i < len(cand); i++ {
		if base[i] != cand[i] {
			return fmt.Sprintf("line %d differs: %q vs %q", i+1, base[i], cand[i])
		}
	}
	if len(base) != len(cand) {
		return fmt.Sprintf("line count differs: baseline=%d candidate=%d", len(base), len(cand))
	}
	return ""
}

func checkPDFRegressionLinks(baseline, candidate, tmpDir string) error {
	baseLinks, err := countPDFLinks(baseline, filepath.Join(tmpDir, "baseline-links"))
	if err != nil {
		return err
	}
	candLinks, err := countPDFLinks(candidate, filepath.Join(tmpDir, "candidate-links"))
	if err != nil {
		return err
	}
	if candLinks < baseLinks {
		return fmt.Errorf("::error::hyperlink count regression on pages 1-3: baseline=%d candidate=%d (%s vs %s)", baseLinks, candLinks, baseline, candidate)
	}
	if candLinks > baseLinks {
		fmt.Printf("pdf regression: hyperlink count increased on pages 1-3 (baseline=%d candidate=%d) — not a failure, worth a human look\n", baseLinks, candLinks)
	}
	return nil
}

func countPDFLinks(pdf, outputPrefix string) (int, error) {
	if _, err := runCommand("pdftohtml", "-xml", "-i", "-f", "1", "-l", "3", pdf, outputPrefix); err != nil {
		return 0, err
	}
	content, err := os.ReadFile(outputPrefix + ".xml") //nolint:gosec // G304: temporary path created by this process.
	if err != nil {
		return 0, err
	}
	return strings.Count(string(content), "<a href="), nil
}

func checkPDFRegressionRaster(baseline, candidate, tmpDir string) error {
	basePages, err := pdfPageCount(baseline)
	if err != nil {
		return err
	}
	candPages, err := pdfPageCount(candidate)
	if err != nil {
		return err
	}
	pages := min(basePages, candPages, pdfRegressionMaxComparedPages)
	if pages < 1 {
		return errors.New("::error::no pages available for rasterized comparison")
	}
	if err := rasterizePDF(baseline, filepath.Join(tmpDir, "baseline-page"), pages); err != nil {
		return err
	}
	if err := rasterizePDF(candidate, filepath.Join(tmpDir, "candidate-page"), pages); err != nil {
		return err
	}
	return compareRasterizedPagePairs(tmpDir, pages)
}

var pdfRegressionPageNumber = regexp.MustCompile(`-(\d+)\.png$`)

func rasterizedPagesByNumber(dir, prefix string) (map[int]string, error) {
	paths, err := filepath.Glob(filepath.Join(dir, prefix+"-*.png"))
	if err != nil {
		return nil, err
	}
	pages := make(map[int]string, len(paths))
	for _, path := range paths {
		match := pdfRegressionPageNumber.FindStringSubmatch(path)
		if match == nil {
			continue
		}
		number, err := strconv.Atoi(match[1])
		if err != nil {
			continue
		}
		pages[number] = path
	}
	return pages, nil
}

func compareRasterizedPagePairs(tmpDir string, pages int) error {
	basePages, err := rasterizedPagesByNumber(tmpDir, "baseline-page")
	if err != nil {
		return err
	}
	candPages, err := rasterizedPagesByNumber(tmpDir, "candidate-page")
	if err != nil {
		return err
	}
	for i := 1; i <= pages; i++ {
		basePath, ok := basePages[i]
		if !ok {
			return fmt.Errorf("::error::missing baseline raster for page %d", i)
		}
		candPath, ok := candPages[i]
		if !ok {
			return fmt.Errorf("::error::missing candidate raster for page %d", i)
		}
		if err := comparePagePair(i, basePath, candPath); err != nil {
			return err
		}
	}
	return nil
}

func comparePagePair(page int, basePath, candPath string) error {
	ratio, err := rasterDiffRatio(basePath, candPath)
	if err != nil {
		return err
	}
	if ratio > pdfRegressionRasterTolerance {
		return fmt.Errorf("::error::page %d rendered diff %.4f exceeds tolerance %.4f (%s vs %s)", page, ratio, pdfRegressionRasterTolerance, basePath, candPath)
	}
	return nil
}

func rasterDiffRatio(basePath, candPath string) (float64, error) {
	baseImg, err := decodePNG(basePath)
	if err != nil {
		return 0, err
	}
	candImg, err := decodePNG(candPath)
	if err != nil {
		return 0, err
	}
	baseBounds := baseImg.Bounds()
	candBounds := candImg.Bounds()
	if baseBounds.Dx() != candBounds.Dx() || baseBounds.Dy() != candBounds.Dy() {
		return 1, nil
	}

	var total, diff int64
	for y := 0; y < baseBounds.Dy(); y++ {
		for x := 0; x < baseBounds.Dx(); x++ {
			br, bg, bb, _ := baseImg.At(baseBounds.Min.X+x, baseBounds.Min.Y+y).RGBA()
			cr, cg, cb, _ := candImg.At(candBounds.Min.X+x, candBounds.Min.Y+y).RGBA()
			total++
			if pixelDiffers(br, bg, bb, cr, cg, cb) {
				diff++
			}
		}
	}
	return float64(diff) / float64(total), nil
}

func pixelDiffers(br, bg, bb, cr, cg, cb uint32) bool {
	return channelDiff(br, cr) > pdfRegressionPixelChannelThreshold ||
		channelDiff(bg, cg) > pdfRegressionPixelChannelThreshold ||
		channelDiff(bb, cb) > pdfRegressionPixelChannelThreshold
}

func channelDiff(a, b uint32) int {
	av, bv := int(a>>8), int(b>>8)
	if av > bv {
		return av - bv
	}
	return bv - av
}

func decodePNG(path string) (image.Image, error) {
	file, err := os.Open(path) //nolint:gosec // G304: temporary raster path created by this process.
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = file.Close()
	}()
	return png.Decode(file)
}
