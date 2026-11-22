package app

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

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

// MakerWorkspaceStats collects and renders image and filtered-text metrics for
// a Maker story workspace directory.
func MakerWorkspaceStats(dir string) (string, error) {
	stats, err := maker.CollectWorkspaceStats(dir)
	if err != nil {
		return "", err
	}
	return maker.RenderWorkspaceStats(stats), nil
}

// ExportMakerChatPack locates the batch JSON inside workspaceDir/packages/,
// generates compact chat-loading Markdown files under outDir, and returns a summary.
func ExportMakerChatPack(workspaceDir, outDir string) (string, error) {
	batchPath, err := findWorkspaceBatchJSON(workspaceDir)
	if err != nil {
		return "", err
	}
	raw, err := os.ReadFile(batchPath) //nolint:gosec // G304: caller-supplied workspace directory.
	if err != nil {
		return "", err
	}
	var batch maker.BatchInput
	if err := json.Unmarshal(raw, &batch); err != nil {
		return "", fmt.Errorf("invalid maker batch JSON %q: %w", batchPath, err)
	}
	res, err := maker.ExportChatPack(batch, outDir)
	if err != nil {
		return "", err
	}
	return maker.RenderExportSummary(res), nil
}

func findWorkspaceBatchJSON(workspaceDir string) (string, error) {
	pkgDir := filepath.Join(workspaceDir, "packages")
	entries, err := os.ReadDir(pkgDir)
	if err != nil {
		return "", err
	}
	for _, entry := range entries {
		if !entry.IsDir() && strings.HasSuffix(entry.Name(), ".json") {
			return filepath.Join(pkgDir, entry.Name()), nil
		}
	}
	return "", fmt.Errorf("maker export: no batch .json file found in %q", pkgDir)
}
