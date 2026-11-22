package maker

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/SergioLacerda/rpg-system-rgb/internal/components"
)

// BatchInput is the structured JSON payload produced by the Maker LLM layer
// for one input batch.
type BatchInput struct {
	ReportID           string              `json:"report_id"`
	InputBatch         []string            `json:"input_batch"`
	EntityPackages     []EntityPackage     `json:"entity_packages"`
	ConflictReports    []ConflictReport    `json:"conflict_reports,omitempty"`
	VisualObservations []VisualObservation `json:"visual_observations,omitempty"`
	OpenQuestions      []string            `json:"open_questions,omitempty"`
}

// BatchResult holds the routed inputs, validated session report, original
// parsed batch input, and combined human-readable Markdown output.
type BatchResult struct {
	Report       MakerReport `json:"report"`
	RoutedInputs []InputRef  `json:"routed_inputs"`
	Input        BatchInput  `json:"-"`
	Markdown     string      `json:"markdown"`
}

// ProcessBatch routes raw inputs, validates all packages (and canonical
// dependencies when index is provided), aggregates the MakerReport, and
// renders the combined Markdown artifacts.
func ProcessBatch(rawJSON []byte, index *components.SemanticIndex) (BatchResult, error) {
	var in BatchInput
	if err := json.Unmarshal(rawJSON, &in); err != nil {
		return BatchResult{}, fmt.Errorf("invalid maker batch JSON: %w", err)
	}

	routed, allErrs := RouteInputs(in.InputBatch)
	allErrs = append(allErrs, validateBatchPackages(in, index)...)

	report, aggErrs := AggregateMakerReport(AggregateInput(in))
	allErrs = append(allErrs, aggErrs...)
	allErrs = append(allErrs, ValidateMakerReport(report)...)

	if len(allErrs) > 0 {
		return BatchResult{}, errors.Join(allErrs...)
	}

	return BatchResult{
		Report:       report,
		RoutedInputs: routed,
		Input:        in,
		Markdown:     renderBatchMarkdown(report, in),
	}, nil
}

// WriteBatchFiles writes maker-report.md and individual Markdown files under
// outDir/entities/, outDir/visual-observations/, outDir/conflicts/ (open),
// and outDir/conflicts/resolved/ (archived decisions), pruning stale .md files.
func WriteBatchFiles(outDir string, res BatchResult) error {
	if strings.TrimSpace(outDir) == "" {
		return fmt.Errorf("maker: outDir must not be empty")
	}
	if err := os.MkdirAll(outDir, 0o750); err != nil {
		return err
	}
	if err := os.WriteFile(filepath.Join(outDir, "maker-report.md"), []byte(RenderMakerReport(res.Report)), 0o644); err != nil { //nolint:gosec // G306: rendered Markdown package artifact.
		return err
	}
	if err := writePackageMarkdowns(filepath.Join(outDir, "entities"), res.Input.EntityPackages, func(p EntityPackage) string { return p.ID }, RenderEntityPackage); err != nil {
		return err
	}
	if err := writePackageMarkdowns(filepath.Join(outDir, "visual-observations"), res.Input.VisualObservations, func(o VisualObservation) string { return o.ID }, RenderVisualObservation); err != nil {
		return err
	}
	return writeConflictMarkdowns(filepath.Join(outDir, "conflicts"), res.Input.ConflictReports)
}

func writeConflictMarkdowns(conflictsDir string, reports []ConflictReport) error {
	var openReports, resolvedReports []ConflictReport
	for _, r := range reports {
		if r.Resolution == "unresolved" {
			openReports = append(openReports, r)
			continue
		}
		resolvedReports = append(resolvedReports, r)
	}
	if err := writePackageMarkdowns(conflictsDir, openReports, func(r ConflictReport) string { return r.ID }, RenderConflictReport); err != nil {
		return err
	}
	return writePackageMarkdowns(filepath.Join(conflictsDir, "resolved"), resolvedReports, func(r ConflictReport) string { return r.ID }, RenderConflictReport)
}

func writePackageMarkdowns[T any](dir string, items []T, getID func(T) string, render func(T) string) error {
	if err := os.MkdirAll(dir, 0o750); err != nil {
		return err
	}
	keep := make(map[string]bool, len(items))
	for _, item := range items {
		id := getID(item)
		if err := validatePackageID("package", id); err != nil {
			return err
		}
		name := id + ".md"
		keep[name] = true
		if err := os.WriteFile(filepath.Join(dir, name), []byte(render(item)), 0o644); err != nil { //nolint:gosec // G306: rendered Markdown package artifact.
			return err
		}
	}
	return pruneStaleMarkdowns(dir, keep)
}

func pruneStaleMarkdowns(dir string, keep map[string]bool) error {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return err
	}
	for _, entry := range entries {
		if !isStaleMarkdown(entry, keep) {
			continue
		}
		if err := os.Remove(filepath.Join(dir, entry.Name())); err != nil {
			return err
		}
	}
	return nil
}

func isStaleMarkdown(entry os.DirEntry, keep map[string]bool) bool {
	if entry.IsDir() {
		return false
	}
	return strings.HasSuffix(entry.Name(), ".md") && !keep[entry.Name()]
}

func validateBatchPackages(in BatchInput, index *components.SemanticIndex) []error {
	var errs []error
	for _, p := range in.EntityPackages {
		errs = append(errs, ValidateEntityPackage(p)...)
		if index != nil {
			errs = append(errs, ValidateCanonicalDependencies(p, *index)...)
		}
	}
	for _, o := range in.VisualObservations {
		errs = append(errs, ValidateVisualObservation(o)...)
	}
	for _, r := range in.ConflictReports {
		errs = append(errs, ValidateConflictReport(r)...)
	}
	return errs
}

func renderBatchMarkdown(report MakerReport, in BatchInput) string {
	sections := []string{strings.TrimSpace(RenderMakerReport(report))}
	for _, o := range in.VisualObservations {
		sections = append(sections, strings.TrimSpace(RenderVisualObservation(o)))
	}
	for _, r := range in.ConflictReports {
		sections = append(sections, strings.TrimSpace(RenderConflictReport(r)))
	}
	for _, p := range in.EntityPackages {
		sections = append(sections, strings.TrimSpace(RenderEntityPackage(p)))
	}
	return strings.Join(sections, "\n\n---\n\n") + "\n"
}
