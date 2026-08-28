package app

import "github.com/SergioLacerda/rpg-system-rgb/internal/components/tooling"

// PDFRegressionPaths identifies the baseline and candidate PDF directories
// and release identity used for the regression check. BaselineDir and
// CandidateDir may be the same directory — that compares the current
// release against itself, proving the harness before it judges a
// genuinely different renderer's output.
type PDFRegressionPaths struct {
	BaselineDir  string
	CandidateDir string
	Basename     string
	Version      string
}

// CheckPDFRegression validates candidate release PDFs against baseline
// release PDFs across page count, extracted text, hyperlink count, and
// rasterized page diff, for both the en and pt-br locales.
func CheckPDFRegression(paths PDFRegressionPaths) error {
	return tooling.CheckPDFRegressionForRelease(paths.BaselineDir, paths.CandidateDir, paths.Basename, paths.Version)
}
