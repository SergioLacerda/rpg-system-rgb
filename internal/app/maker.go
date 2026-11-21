package app

import (
	"os"

	"github.com/SergioLacerda/rpg-system-rgb/internal/components"
	"github.com/SergioLacerda/rpg-system-rgb/internal/components/maker"
)

// ProcessMakerBatch reads a JSON Maker batch file, optionally validates
// canonical_dependencies against indexFile, writes per-package Markdown files
// when outDir is non-empty, and returns the combined rendered Markdown.
func ProcessMakerBatch(batchFile, indexFile, outDir string) (string, error) {
	raw, err := os.ReadFile(batchFile) //nolint:gosec // G304: path is a caller-supplied batch file, by design
	if err != nil {
		return "", err
	}
	var idxPtr *components.SemanticIndex
	if indexFile != "" {
		idx, err := components.LoadSemanticIndex(indexFile)
		if err != nil {
			return "", err
		}
		idxPtr = &idx
	}
	res, err := maker.ProcessBatch(raw, idxPtr)
	if err != nil {
		return "", err
	}
	if outDir != "" {
		if err := maker.WriteBatchFiles(outDir, res); err != nil {
			return "", err
		}
	}
	return res.Markdown, nil
}
