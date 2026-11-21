package maker

import (
	"encoding/json"
	"errors"
	"fmt"
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

// BatchResult holds the routed inputs, validated session report, and combined
// human-readable Markdown output for one Maker batch.
type BatchResult struct {
	Report       MakerReport `json:"report"`
	RoutedInputs []InputRef  `json:"routed_inputs"`
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
		Markdown:     renderBatchMarkdown(report, in),
	}, nil
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
