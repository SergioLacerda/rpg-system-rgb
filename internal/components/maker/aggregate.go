package maker

import "fmt"

// AggregateInput holds the already-validated packages emitted by the Maker
// LLM layer for one input batch.
type AggregateInput struct {
	ReportID           string
	InputBatch         []string
	EntityPackages     []EntityPackage
	ConflictReports    []ConflictReport
	VisualObservations []VisualObservation
	OpenQuestions      []string
}

// AggregateMakerReport builds a MakerReport from a batch of validated
// packages, checking that every referenced ID corresponds to an emitted
// package and rejecting duplicate package IDs.
//
// The returned MakerReport is always in status "draft" — promotion to
// needs_review or stale is a user decision.
func AggregateMakerReport(in AggregateInput) (MakerReport, []error) {
	var errs []error

	if in.ReportID == "" {
		errs = append(errs, fmt.Errorf("maker report: id must not be empty"))
	}
	if len(in.InputBatch) == 0 {
		errs = append(errs, fmt.Errorf("maker report %q: input_batch must not be empty", in.ReportID))
	}

	epIDs, epSet, epErrs := collectBatchIDs("entity_package", in.EntityPackages, func(p EntityPackage) string { return p.ID })
	crIDs, crSet, crErrs := collectBatchIDs("conflict_report", in.ConflictReports, func(r ConflictReport) string { return r.ID })
	voIDs, _, voErrs := collectBatchIDs("visual_observation", in.VisualObservations, func(o VisualObservation) string { return o.ID })

	errs = append(errs, epErrs...)
	errs = append(errs, crErrs...)
	errs = append(errs, voErrs...)
	errs = append(errs, validateEntityConflicts(in.EntityPackages, crSet)...)
	errs = append(errs, validateConflictLinks(in.ConflictReports, epSet)...)
	errs = append(errs, validateVisualLinks(in.VisualObservations, epSet)...)

	report := MakerReport{
		ID:                 in.ReportID,
		Status:             "draft",
		InputBatch:         in.InputBatch,
		EntityPackages:     epIDs,
		ConflictReports:    crIDs,
		VisualObservations: voIDs,
		OpenQuestions:      in.OpenQuestions,
	}

	return report, errs
}

func collectBatchIDs[T any](kind string, items []T, getID func(T) string) ([]string, map[string]bool, []error) {
	ids := make([]string, 0, len(items))
	set := make(map[string]bool, len(items))
	var errs []error
	for _, item := range items {
		id := getID(item)
		ids = append(ids, id)
		if set[id] {
			errs = append(errs, fmt.Errorf("duplicate %s id %q in batch", kind, id))
			continue
		}
		set[id] = true
	}
	return ids, set, errs
}

func validateEntityConflicts(pkgs []EntityPackage, crSet map[string]bool) []error {
	var errs []error
	for _, p := range pkgs {
		for _, cid := range p.Conflicts {
			if !crSet[cid] {
				errs = append(errs, fmt.Errorf("entity package %q references conflict %q which is not in the emitted conflict reports", p.ID, cid))
			}
		}
	}
	return errs
}

func validateConflictLinks(reports []ConflictReport, epSet map[string]bool) []error {
	var errs []error
	for _, r := range reports {
		if r.LinkedEntity != "" && !epSet[r.LinkedEntity] {
			errs = append(errs, fmt.Errorf("conflict report %q is linked to entity %q which is not in the emitted entity packages", r.ID, r.LinkedEntity))
		}
	}
	return errs
}

func validateVisualLinks(obs []VisualObservation, epSet map[string]bool) []error {
	var errs []error
	for _, o := range obs {
		if o.LinkedEntity != "" && !epSet[o.LinkedEntity] {
			errs = append(errs, fmt.Errorf("visual observation %q is linked to entity %q which is not in the emitted entity packages", o.ID, o.LinkedEntity))
		}
	}
	return errs
}
