package tooling

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func withFakePDFTools(t *testing.T, runner func(name string, args ...string) ([]byte, error)) {
	t.Helper()
	originalLookPath := lookPath
	t.Cleanup(func() { lookPath = originalLookPath })
	lookPath = func(name string) (string, error) { return "/fake/bin/" + name, nil }

	originalRunner := runExternalCommand
	t.Cleanup(func() { runExternalCommand = originalRunner })
	runExternalCommand = runner
}

func fakePDFRegressionRunner(t *testing.T, pages map[string]int, text map[string]string, links map[string]int) func(string, ...string) ([]byte, error) {
	t.Helper()
	return func(name string, args ...string) ([]byte, error) {
		switch name {
		case "pdfinfo":
			pdf := args[len(args)-1]
			return []byte(fmt.Sprintf("Pages: %d\n", pages[pdf])), nil
		case "pdftotext":
			pdf := args[len(args)-2]
			out := args[len(args)-1]
			return nil, os.WriteFile(out, []byte(text[pdf]), 0o644)
		case "pdftohtml":
			pdf := args[len(args)-2]
			out := args[len(args)-1] + ".xml"
			body := strings.Repeat(`<a href="#x">x</a>`, links[pdf])
			return nil, os.WriteFile(out, []byte("<page>"+body+"</page>"), 0o644)
		case "pdftoppm":
			prefix := args[len(args)-1]
			n := pages[args[len(args)-2]]
			for i := 1; i <= n; i++ {
				writeReadablePNG(t, fmt.Sprintf("%s-%d.png", prefix, i))
			}
			return nil, nil
		default:
			return nil, fmt.Errorf("unexpected command %s", name)
		}
	}
}

func TestCheckPDFRegressionSelfComparisonPasses(t *testing.T) {
	dir := t.TempDir()
	doc := filepath.Join(dir, "doc.pdf")
	writeTestPDF(t, doc, strings.Repeat("doc", 400))

	withFakePDFTools(t, fakePDFRegressionRunner(t,
		map[string]int{doc: 3},
		map[string]string{doc: "Line one\nLine two\n2026-08-28\nabc1234\nLine three\n"},
		map[string]int{doc: 2},
	))

	if err := CheckPDFRegression(doc, doc); err != nil {
		t.Fatalf("expected self-comparison to pass, got %v", err)
	}
}

func TestCheckPDFRegressionMissingToolsFails(t *testing.T) {
	originalLookPath := lookPath
	t.Cleanup(func() { lookPath = originalLookPath })
	lookPath = func(name string) (string, error) { return "", fmt.Errorf("missing %s", name) }

	if err := CheckPDFRegression("a.pdf", "b.pdf"); err == nil {
		t.Fatal("expected missing PDF tool to fail")
	}
}

func TestCheckPDFRegressionPageCountRegression(t *testing.T) {
	dir := t.TempDir()
	baseline := filepath.Join(dir, "baseline.pdf")
	candidate := filepath.Join(dir, "candidate.pdf")
	writeTestPDF(t, baseline, "x")
	writeTestPDF(t, candidate, "y")

	withFakePDFTools(t, fakePDFRegressionRunner(t,
		map[string]int{baseline: 4, candidate: 5},
		nil, nil,
	))

	err := CheckPDFRegression(baseline, candidate)
	if err == nil || !strings.Contains(err.Error(), "page count regression") {
		t.Fatalf("expected page count regression error, got %v", err)
	}
}

func TestCheckPDFRegressionTextRegression(t *testing.T) {
	dir := t.TempDir()
	baseline := filepath.Join(dir, "baseline.pdf")
	candidate := filepath.Join(dir, "candidate.pdf")
	writeTestPDF(t, baseline, "x")
	writeTestPDF(t, candidate, "y")

	withFakePDFTools(t, fakePDFRegressionRunner(t,
		map[string]int{baseline: 2, candidate: 2},
		map[string]string{
			baseline:  "Rule A applies\nSecond line\n",
			candidate: "Rule A does not apply\nSecond line\n",
		},
		nil,
	))

	err := CheckPDFRegression(baseline, candidate)
	if err == nil || !strings.Contains(err.Error(), "extracted text regression") {
		t.Fatalf("expected text regression error, got %v", err)
	}
}

func TestCheckPDFRegressionTextIgnoresStampedDateAndCommit(t *testing.T) {
	dir := t.TempDir()
	baseline := filepath.Join(dir, "baseline.pdf")
	candidate := filepath.Join(dir, "candidate.pdf")
	writeTestPDF(t, baseline, "x")
	writeTestPDF(t, candidate, "y")

	withFakePDFTools(t, fakePDFRegressionRunner(t,
		map[string]int{baseline: 2, candidate: 2},
		map[string]string{
			baseline:  "RGB System\n2026-08-27\nhttps://example.com @ abc1234\n",
			candidate: "RGB System\n2026-08-28\nhttps://example.com @ fed5678\n",
		},
		map[string]int{baseline: 1, candidate: 1},
	))

	if err := CheckPDFRegression(baseline, candidate); err != nil {
		t.Fatalf("expected stamped date/commit lines to be ignored, got %v", err)
	}
}

func TestNormalizePDFRegressionTextStripsStampedLines(t *testing.T) {
	lines := normalizePDFRegressionText("Title\n2026-08-28\nabc1234\nBody text")
	if len(lines) != 2 || lines[0] != "Title" || lines[1] != "Body text" {
		t.Fatalf("unexpected normalized lines: %#v", lines)
	}
}

func TestCheckPDFRegressionLinksRegressionFails(t *testing.T) {
	dir := t.TempDir()
	baseline := filepath.Join(dir, "baseline.pdf")
	candidate := filepath.Join(dir, "candidate.pdf")
	writeTestPDF(t, baseline, "x")
	writeTestPDF(t, candidate, "y")

	withFakePDFTools(t, fakePDFRegressionRunner(t,
		map[string]int{baseline: 2, candidate: 2},
		map[string]string{baseline: "same\n", candidate: "same\n"},
		map[string]int{baseline: 5, candidate: 3},
	))

	err := CheckPDFRegression(baseline, candidate)
	if err == nil || !strings.Contains(err.Error(), "hyperlink count regression") {
		t.Fatalf("expected hyperlink regression error, got %v", err)
	}
}

func TestCheckPDFRegressionLinksIncreaseDoesNotFail(t *testing.T) {
	dir := t.TempDir()
	baseline := filepath.Join(dir, "baseline.pdf")
	candidate := filepath.Join(dir, "candidate.pdf")
	writeTestPDF(t, baseline, "x")
	writeTestPDF(t, candidate, "y")

	withFakePDFTools(t, fakePDFRegressionRunner(t,
		map[string]int{baseline: 2, candidate: 2},
		map[string]string{baseline: "same\n", candidate: "same\n"},
		map[string]int{baseline: 3, candidate: 5},
	))

	if err := CheckPDFRegression(baseline, candidate); err != nil {
		t.Fatalf("expected an increase in links not to fail, got %v", err)
	}
}

func TestCheckPDFRegressionRasterRegression(t *testing.T) {
	dir := t.TempDir()
	baseline := filepath.Join(dir, "baseline.pdf")
	candidate := filepath.Join(dir, "candidate.pdf")
	writeTestPDF(t, baseline, "x")
	writeTestPDF(t, candidate, "y")

	withFakePDFTools(t, func(name string, args ...string) ([]byte, error) {
		switch name {
		case "pdfinfo":
			return []byte("Pages: 1\n"), nil
		case "pdftotext":
			return nil, os.WriteFile(args[len(args)-1], []byte("same\n"), 0o644)
		case "pdftohtml":
			return nil, os.WriteFile(args[len(args)-1]+".xml", []byte(`<page><a href="#x">x</a></page>`), 0o644)
		case "pdftoppm":
			pdf := args[len(args)-2]
			prefix := args[len(args)-1]
			if strings.Contains(pdf, "candidate") {
				writeDarkPNG(t, fmt.Sprintf("%s-1.png", prefix))
			} else {
				writeReadablePNG(t, fmt.Sprintf("%s-1.png", prefix))
			}
			return nil, nil
		default:
			return nil, fmt.Errorf("unexpected command %s", name)
		}
	})

	err := CheckPDFRegression(baseline, candidate)
	if err == nil || !strings.Contains(err.Error(), "rendered diff") {
		t.Fatalf("expected rendered diff regression error, got %v", err)
	}
}

func TestRasterDiffRatioDimensionMismatch(t *testing.T) {
	dir := t.TempDir()
	small := filepath.Join(dir, "small.png")
	big := filepath.Join(dir, "big.png")
	writeSizedPNG(t, small, 50, 50)
	writeSizedPNG(t, big, 80, 80)

	ratio, err := rasterDiffRatio(small, big)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if ratio != 1 {
		t.Fatalf("expected full mismatch ratio for differing dimensions, got %v", ratio)
	}
}

func TestRasterizedPagesByNumberParsesPadded(t *testing.T) {
	dir := t.TempDir()
	writeReadablePNG(t, filepath.Join(dir, "page-01.png"))
	writeReadablePNG(t, filepath.Join(dir, "page-02.png"))

	pages, err := rasterizedPagesByNumber(dir, "page")
	if err != nil {
		t.Fatal(err)
	}
	if len(pages) != 2 || pages[1] == "" || pages[2] == "" {
		t.Fatalf("unexpected pages map: %#v", pages)
	}
}

func TestCheckPDFRegressionForReleaseWrapsLocale(t *testing.T) {
	dir := t.TempDir()
	baselineDir := filepath.Join(dir, "baseline")
	candidateDir := filepath.Join(dir, "candidate")
	for _, sub := range []string{baselineDir, candidateDir} {
		if err := os.MkdirAll(sub, 0o750); err != nil {
			t.Fatal(err)
		}
	}
	basename, version := "rgb-system-core-v2", "v2.0"
	for _, locale := range []string{"en", "pt-br"} {
		name := fmt.Sprintf("%s-%s-%s.pdf", basename, version, locale)
		writeTestPDF(t, filepath.Join(baselineDir, name), locale)
		writeTestPDF(t, filepath.Join(candidateDir, name), locale)
	}

	enPDF := filepath.Join(baselineDir, fmt.Sprintf("%s-%s-en.pdf", basename, version))
	enCandidatePDF := filepath.Join(candidateDir, fmt.Sprintf("%s-%s-en.pdf", basename, version))
	ptPDF := filepath.Join(baselineDir, fmt.Sprintf("%s-%s-pt-br.pdf", basename, version))
	ptCandidatePDF := filepath.Join(candidateDir, fmt.Sprintf("%s-%s-pt-br.pdf", basename, version))

	withFakePDFTools(t, fakePDFRegressionRunner(t,
		map[string]int{enPDF: 2, enCandidatePDF: 2, ptPDF: 2, ptCandidatePDF: 3},
		map[string]string{enPDF: "same\n", enCandidatePDF: "same\n", ptPDF: "same\n", ptCandidatePDF: "same\n"},
		map[string]int{enPDF: 1, enCandidatePDF: 1, ptPDF: 1, ptCandidatePDF: 1},
	))

	err := CheckPDFRegressionForRelease(baselineDir, candidateDir, basename, version)
	if err == nil || !strings.Contains(err.Error(), "pdf regression (pt-br)") {
		t.Fatalf("expected pt-br locale to fail with wrapped error, got %v", err)
	}
}

func TestCheckPDFRegressionPageCountToolFailure(t *testing.T) {
	dir := t.TempDir()
	baseline := filepath.Join(dir, "baseline.pdf")
	candidate := filepath.Join(dir, "candidate.pdf")
	writeTestPDF(t, baseline, "x")
	writeTestPDF(t, candidate, "y")

	withFakePDFTools(t, func(name string, args ...string) ([]byte, error) {
		if name != "pdfinfo" {
			return nil, fmt.Errorf("unexpected command %s", name)
		}
		pdf := args[len(args)-1]
		if strings.Contains(pdf, "baseline") {
			return []byte("no page field\n"), nil
		}
		return []byte("Pages: 3\n"), nil
	})

	if err := checkPDFRegressionPageCount(baseline, candidate); err == nil {
		t.Fatal("expected baseline page-count parse failure to propagate")
	}
	if err := checkPDFRegressionPageCount(candidate, baseline); err == nil {
		t.Fatal("expected candidate page-count parse failure to propagate")
	}
}

func TestExtractPDFTextCommandFailure(t *testing.T) {
	withFakePDFTools(t, func(name string, _ ...string) ([]byte, error) {
		return nil, fmt.Errorf("%s failed", name)
	})
	if _, err := extractPDFText("doc.pdf", filepath.Join(t.TempDir(), "out.txt")); err == nil {
		t.Fatal("expected pdftotext failure to propagate")
	}
}

func TestCheckPDFRegressionTextCandidateExtractFailure(t *testing.T) {
	dir := t.TempDir()
	baseline := filepath.Join(dir, "baseline.pdf")
	candidate := filepath.Join(dir, "candidate.pdf")

	withFakePDFTools(t, func(name string, args ...string) ([]byte, error) {
		if name != "pdftotext" {
			return nil, fmt.Errorf("unexpected command %s", name)
		}
		pdf := args[len(args)-2]
		if strings.Contains(pdf, "candidate") {
			return nil, fmt.Errorf("pdftotext failed")
		}
		return nil, os.WriteFile(args[len(args)-1], []byte("ok\n"), 0o644)
	})

	if err := checkPDFRegressionText(baseline, candidate, dir); err == nil {
		t.Fatal("expected candidate text extraction failure to propagate")
	}
}

func TestCountPDFLinksCommandFailure(t *testing.T) {
	withFakePDFTools(t, func(name string, _ ...string) ([]byte, error) {
		return nil, fmt.Errorf("%s failed", name)
	})
	if _, err := countPDFLinks("doc.pdf", filepath.Join(t.TempDir(), "links")); err == nil {
		t.Fatal("expected pdftohtml failure to propagate")
	}
}

func TestCheckPDFRegressionLinksCandidateFailure(t *testing.T) {
	dir := t.TempDir()
	baseline := filepath.Join(dir, "baseline.pdf")
	candidate := filepath.Join(dir, "candidate.pdf")

	withFakePDFTools(t, func(name string, args ...string) ([]byte, error) {
		if name != "pdftohtml" {
			return nil, fmt.Errorf("unexpected command %s", name)
		}
		pdf := args[len(args)-2]
		if strings.Contains(pdf, "candidate") {
			return nil, fmt.Errorf("pdftohtml failed")
		}
		return nil, os.WriteFile(args[len(args)-1]+".xml", []byte(`<page><a href="#x">x</a></page>`), 0o644)
	})

	if err := checkPDFRegressionLinks(baseline, candidate, dir); err == nil {
		t.Fatal("expected candidate link extraction failure to propagate")
	}
}

func TestCheckPDFRegressionRasterNoPagesAvailable(t *testing.T) {
	dir := t.TempDir()
	baseline := filepath.Join(dir, "baseline.pdf")
	candidate := filepath.Join(dir, "candidate.pdf")

	withFakePDFTools(t, func(name string, args ...string) ([]byte, error) {
		if name != "pdfinfo" {
			return nil, fmt.Errorf("unexpected command %s", name)
		}
		return []byte("Pages: 0\n"), nil
	})

	err := checkPDFRegressionRaster(baseline, candidate, dir)
	if err == nil || !strings.Contains(err.Error(), "no pages available") {
		t.Fatalf("expected no-pages-available error, got %v", err)
	}
}

func TestCheckPDFRegressionRasterizeCandidateFailure(t *testing.T) {
	dir := t.TempDir()
	baseline := filepath.Join(dir, "baseline.pdf")
	candidate := filepath.Join(dir, "candidate.pdf")

	withFakePDFTools(t, func(name string, args ...string) ([]byte, error) {
		switch name {
		case "pdfinfo":
			return []byte("Pages: 1\n"), nil
		case "pdftoppm":
			pdf := args[len(args)-2]
			if strings.Contains(pdf, "candidate") {
				return nil, fmt.Errorf("pdftoppm failed")
			}
			writeReadablePNG(t, fmt.Sprintf("%s-1.png", args[len(args)-1]))
			return nil, nil
		default:
			return nil, fmt.Errorf("unexpected command %s", name)
		}
	})

	if err := checkPDFRegressionRaster(baseline, candidate, dir); err == nil {
		t.Fatal("expected candidate rasterization failure to propagate")
	}
}

func TestCompareRasterizedPagePairsMissingRaster(t *testing.T) {
	dir := t.TempDir()
	writeReadablePNG(t, filepath.Join(dir, "baseline-page-1.png"))
	writeReadablePNG(t, filepath.Join(dir, "candidate-page-1.png"))

	err := compareRasterizedPagePairs(dir, 2)
	if err == nil || !strings.Contains(err.Error(), "missing baseline raster for page 2") {
		t.Fatalf("expected missing baseline raster error, got %v", err)
	}

	writeReadablePNG(t, filepath.Join(dir, "baseline-page-2.png"))
	err = compareRasterizedPagePairs(dir, 2)
	if err == nil || !strings.Contains(err.Error(), "missing candidate raster for page 2") {
		t.Fatalf("expected missing candidate raster error, got %v", err)
	}
}

func TestDecodePNGInvalidFile(t *testing.T) {
	dir := t.TempDir()
	bad := filepath.Join(dir, "not-a-png.png")
	if err := os.WriteFile(bad, []byte("not a png"), 0o644); err != nil {
		t.Fatal(err)
	}
	if _, err := decodePNG(bad); err == nil {
		t.Fatal("expected invalid PNG to fail decoding")
	}
	if _, err := decodePNG(filepath.Join(dir, "missing.png")); err == nil {
		t.Fatal("expected missing PNG file to fail")
	}
}

func TestCheckPDFRegressionForReleaseBothLocalesPass(t *testing.T) {
	dir := t.TempDir()
	baselineDir := filepath.Join(dir, "baseline")
	candidateDir := filepath.Join(dir, "candidate")
	for _, sub := range []string{baselineDir, candidateDir} {
		if err := os.MkdirAll(sub, 0o750); err != nil {
			t.Fatal(err)
		}
	}
	basename, version := "rgb-system-core-v2", "v2.0"
	pages := map[string]int{}
	text := map[string]string{}
	links := map[string]int{}
	for _, dirPair := range []string{baselineDir, candidateDir} {
		for _, locale := range []string{"en", "pt-br"} {
			path := filepath.Join(dirPair, fmt.Sprintf("%s-%s-%s.pdf", basename, version, locale))
			pages[path] = 2
			text[path] = "same\n"
			links[path] = 1
		}
	}

	withFakePDFTools(t, fakePDFRegressionRunner(t, pages, text, links))

	if err := CheckPDFRegressionForRelease(baselineDir, candidateDir, basename, version); err != nil {
		t.Fatalf("expected both locales to pass, got %v", err)
	}
}

func writeDarkPNG(t *testing.T, path string) {
	t.Helper()
	page := image.NewRGBA(image.Rect(0, 0, 200, 200))
	fill(page, color.RGBA{R: 10, G: 10, B: 10, A: 255})
	file, err := os.Create(path)
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = file.Close() }()
	if err := png.Encode(file, page); err != nil {
		t.Fatal(err)
	}
}

func writeSizedPNG(t *testing.T, path string, w, h int) {
	t.Helper()
	page := image.NewRGBA(image.Rect(0, 0, w, h))
	fill(page, color.RGBA{R: 255, G: 255, B: 255, A: 255})
	file, err := os.Create(path)
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = file.Close() }()
	if err := png.Encode(file, page); err != nil {
		t.Fatal(err)
	}
}
