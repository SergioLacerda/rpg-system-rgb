package maker

import (
	"fmt"
	"strings"

	"github.com/SergioLacerda/rpg-system-rgb/internal/components"
)

var forbiddenAuthorityTypes = map[string]bool{"projection": true, "generated_artifact": true}

// ValidateConflictReport checks a ConflictReport against conflict-report.schema.yaml.
func ValidateConflictReport(r ConflictReport) []error {
	errs := validateConflictHeader(r)
	for i, c := range r.Claims {
		errs = append(errs, validateConflictClaim(r.ID, i, c)...)
	}
	return append(errs, validateConflictResolution(r)...)
}

func isSafePackageID(id string) bool {
	return strings.TrimSpace(id) != "" && !strings.Contains(id, "..") && !strings.ContainsAny(id, `/\`)
}

func validatePackageID(kind, id string) error {
	if strings.TrimSpace(id) == "" {
		return fmt.Errorf("%s: id must not be empty", kind)
	}
	if !isSafePackageID(id) {
		return fmt.Errorf("%s %q: id must not contain path separators or '..'", kind, id)
	}
	return nil
}

func validateConflictHeader(r ConflictReport) []error {
	var errs []error
	if err := validatePackageID("conflict report", r.ID); err != nil {
		errs = append(errs, err)
	}
	if strings.TrimSpace(r.Field) == "" {
		errs = append(errs, fmt.Errorf("conflict report %q: field must not be empty", r.ID))
	}
	if !validStatus[r.Status] {
		errs = append(errs, fmt.Errorf("conflict report %q: invalid status %q", r.ID, r.Status))
	}
	if !validResolution[r.Resolution] {
		errs = append(errs, fmt.Errorf("conflict report %q: invalid resolution %q", r.ID, r.Resolution))
	}
	if len(r.Claims) < 2 {
		errs = append(errs, fmt.Errorf("conflict report %q: must have at least 2 claims, got %d", r.ID, len(r.Claims)))
	}
	return errs
}

func validateConflictClaim(reportID string, index int, c Claim) []error {
	var errs []error
	if strings.TrimSpace(c.Value) == "" {
		errs = append(errs, fmt.Errorf("conflict report %q: claim[%d] has empty value", reportID, index))
	}
	if strings.TrimSpace(c.SourceRef) == "" {
		errs = append(errs, fmt.Errorf("conflict report %q: claim[%d] has empty source_ref", reportID, index))
	}
	if !validSourceKind[c.SourceKind] {
		errs = append(errs, fmt.Errorf("conflict report %q: claim[%d] has invalid source_kind %q (must be text|image)", reportID, index, c.SourceKind))
	}
	return errs
}

func validateConflictResolution(r ConflictReport) []error {
	var errs []error
	if r.Resolution == "resolved_by_authority_rule" {
		errs = append(errs, validateAuthorityRuleResolution(r)...)
	}
	if r.Resolution == "resolved_by_user" && strings.TrimSpace(r.ResolutionNote) == "" {
		errs = append(errs, fmt.Errorf("conflict report %q: resolution_note is required when resolution is resolved_by_user", r.ID))
	}
	return errs
}

func validateAuthorityRuleResolution(r ConflictReport) []error {
	var errs []error
	if !hasClaimKind(r.Claims, "image") || !hasClaimKind(r.Claims, "text") {
		errs = append(errs, fmt.Errorf("conflict report %q: resolution is resolved_by_authority_rule but claims must include both text and image — the authority rule only applies to image-vs-text disputes", r.ID))
	}
	if strings.TrimSpace(r.ResolutionNote) == "" {
		errs = append(errs, fmt.Errorf("conflict report %q: resolution_note is required when resolution is resolved_by_authority_rule", r.ID))
	}
	return errs
}

func hasClaimKind(claims []Claim, kind string) bool {
	for _, c := range claims {
		if c.SourceKind == kind {
			return true
		}
	}
	return false
}

// ValidateMakerReport checks a MakerReport against maker-report.schema.yaml.
func ValidateMakerReport(r MakerReport) []error {
	errs := validateMakerReportHeader(r)
	errs = append(errs, validateReportStrings(r.ID, "input_batch", r.InputBatch)...)
	errs = append(errs, validateUniqueIDs(r.ID, "entity_packages", r.EntityPackages)...)
	errs = append(errs, validateUniqueIDs(r.ID, "conflict_reports", r.ConflictReports)...)
	errs = append(errs, validateUniqueIDs(r.ID, "visual_observations", r.VisualObservations)...)
	errs = append(errs, validateReportStrings(r.ID, "open_questions", r.OpenQuestions)...)
	return append(errs, validateReportStrings(r.ID, "pending_decisions", r.PendingDecisions)...)
}

func validateMakerReportHeader(r MakerReport) []error {
	var errs []error
	if err := validatePackageID("maker report", r.ID); err != nil {
		errs = append(errs, err)
	}
	if !validStatus[r.Status] {
		errs = append(errs, fmt.Errorf("maker report %q: invalid status %q", r.ID, r.Status))
	}
	if len(r.InputBatch) == 0 {
		errs = append(errs, fmt.Errorf("maker report %q: input_batch must not be empty", r.ID))
	}
	if r.EntityPackages == nil {
		errs = append(errs, fmt.Errorf("maker report %q: entity_packages must not be nil (schema requires an explicit list, even when empty)", r.ID))
	}
	return errs
}

func validateReportStrings(reportID, field string, items []string) []error {
	var errs []error
	for i, item := range items {
		if strings.TrimSpace(item) == "" {
			errs = append(errs, fmt.Errorf("maker report %q: %s[%d] must not be empty", reportID, field, i))
		}
	}
	return errs
}

func validateUniqueIDs(reportID, field string, ids []string) []error {
	var errs []error
	seen := make(map[string]bool, len(ids))
	for i, id := range ids {
		if strings.TrimSpace(id) == "" {
			errs = append(errs, fmt.Errorf("maker report %q: %s[%d] must not be empty", reportID, field, i))
			continue
		}
		if seen[id] {
			errs = append(errs, fmt.Errorf("maker report %q: duplicate %s id %q", reportID, field, id))
			continue
		}
		seen[id] = true
	}
	return errs
}

// ValidateCanonicalDependencies validates canonical_dependencies against idx.
func ValidateCanonicalDependencies(p EntityPackage, idx components.SemanticIndex) []error {
	units := make(map[string]components.SemanticIndexUnit, len(idx.Units))
	for _, u := range idx.Units {
		units[u.ID] = u
	}
	var errs []error
	for _, dep := range p.CanonicalDeps {
		unit, ok := units[dep]
		if !ok {
			errs = append(errs, fmt.Errorf("entity package %q: unknown semantic ID %q in canonical_dependencies", p.ID, dep))
			continue
		}
		if forbiddenAuthorityTypes[unit.AuthorityType] {
			errs = append(errs, fmt.Errorf("entity package %q: canonical dependency %q has forbidden authority_type %q", p.ID, dep, unit.AuthorityType))
		}
	}
	return errs
}
