package maker

import (
	"fmt"
	"strings"
)

// nonVisualFields is a curated set of field names belonging to entity facts or
// inferences, never to a visual observation.
var nonVisualFields = map[string]bool{
	"loyalty": true, "backstory": true, "name": true, "stats": true,
	"relationships": true, "canon_status": true, "true_identity": true,
	"alignment": true, "motivation": true, "faction": true,
}

var validStatus = map[string]bool{"draft": true, "needs_review": true, "stale": true}
var validConfidence = map[string]bool{"low": true, "medium": true, "high": true}
var validResolution = map[string]bool{"unresolved": true, "resolved_by_authority_rule": true, "resolved_by_user": true}
var validSourceKind = map[string]bool{"text": true, "image": true}
var validEntityType = map[string]bool{"npc": true, "location": true, "item": true, "faction": true, "other": true}
var fieldNormalizer = strings.NewReplacer("-", "_", " ", "_")

func normalizeFieldName(s string) string {
	return fieldNormalizer.Replace(strings.ToLower(strings.TrimSpace(s)))
}

func validateNonEmptyStrings(pkgID, field string, items []string) []error {
	var errs []error
	for i, item := range items {
		if strings.TrimSpace(item) == "" {
			errs = append(errs, fmt.Errorf("entity package %q: %s[%d] must not be empty", pkgID, field, i))
		}
	}
	return errs
}

func toStringSet(items []string) map[string]bool {
	set := make(map[string]bool, len(items))
	for _, item := range items {
		if strings.TrimSpace(item) != "" {
			set[item] = true
		}
	}
	return set
}

// ValidateEntityPackage checks an EntityPackage against entity-package.schema.yaml:
// non-empty explicit_facts, disjoint fields across layers, provenance integrity, and valid enums.
func ValidateEntityPackage(p EntityPackage) []error {
	errs := validateEntityHeader(p)
	errs = append(errs, validateNonEmptyStrings(p.ID, "source_refs", p.SourceRefs)...)
	errs = append(errs, validateNonEmptyStrings(p.ID, "conflicts", p.Conflicts)...)
	errs = append(errs, validateNonEmptyStrings(p.ID, "canonical_dependencies", p.CanonicalDeps)...)
	sourceSet := toStringSet(p.SourceRefs)
	seen := make(map[string]string)
	errs = append(errs, validateExplicitFacts(p.ID, p.ExplicitFacts, sourceSet, seen)...)
	errs = append(errs, validateInferredSuggestions(p.ID, p.InferredSuggestions, seen)...)
	errs = append(errs, validateUnresolvedFields(p.ID, p.UnresolvedFields, seen)...)
	errs = append(errs, validateEntityVariants(p.ID, p.Variants, sourceSet)...)
	return errs
}

func validateEntityHeader(p EntityPackage) []error {
	var errs []error
	if err := validatePackageID("entity package", p.ID); err != nil {
		errs = append(errs, err)
	}
	if !validEntityType[p.EntityType] {
		errs = append(errs, fmt.Errorf("entity package %q: invalid entity_type %q", p.ID, p.EntityType))
	}
	if !validStatus[p.Status] {
		errs = append(errs, fmt.Errorf("entity package %q: invalid status %q (must be draft|needs_review|stale)", p.ID, p.Status))
	}
	if len(p.SourceRefs) == 0 {
		errs = append(errs, fmt.Errorf("entity package %q: source_refs must not be empty", p.ID))
	}
	if len(p.ExplicitFacts) == 0 {
		errs = append(errs, fmt.Errorf("entity package %q: explicit_facts must not be empty (a package with zero explicit facts is pure invention)", p.ID))
	}
	return errs
}

func recordLayerField(pkgID, layer, field string, seen map[string]string) error {
	if strings.TrimSpace(field) == "" {
		return fmt.Errorf("entity package %q: %s entry has empty field", pkgID, layer)
	}
	if prior, ok := seen[field]; ok {
		return fmt.Errorf("entity package %q: field %q appears in both %s and %s — each field must be in exactly one layer", pkgID, field, prior, layer)
	}
	seen[field] = layer
	return nil
}

func validateExplicitFacts(pkgID string, facts []ExplicitFact, sourceSet map[string]bool, seen map[string]string) []error {
	var errs []error
	for _, f := range facts {
		errs = append(errs, validateExplicitFactItem(pkgID, f, sourceSet, seen)...)
	}
	return errs
}

func validateExplicitFactItem(pkgID string, f ExplicitFact, sourceSet map[string]bool, seen map[string]string) []error {
	var errs []error
	if err := recordLayerField(pkgID, "explicit_facts", f.Field, seen); err != nil {
		errs = append(errs, err)
	}
	if strings.TrimSpace(f.Value) == "" {
		errs = append(errs, fmt.Errorf("entity package %q: explicit_facts field %q has empty value", pkgID, f.Field))
	}
	if strings.TrimSpace(f.SourceRef) == "" {
		errs = append(errs, fmt.Errorf("entity package %q: explicit_facts field %q has empty source_ref", pkgID, f.Field))
	} else if len(sourceSet) > 0 && !sourceSet[f.SourceRef] {
		errs = append(errs, fmt.Errorf("entity package %q: explicit_facts field %q references undeclared source_ref %q", pkgID, f.Field, f.SourceRef))
	}
	return errs
}

func validateInferredSuggestions(pkgID string, items []InferredSuggestion, seen map[string]string) []error {
	var errs []error
	for _, f := range items {
		errs = append(errs, validateInferredItem(pkgID, f, seen)...)
	}
	return errs
}

func validateInferredItem(pkgID string, f InferredSuggestion, seen map[string]string) []error {
	var errs []error
	if strings.TrimSpace(f.Field) != "" && !validConfidence[f.Confidence] {
		errs = append(errs, fmt.Errorf("entity package %q: inferred_suggestions field %q has invalid confidence %q", pkgID, f.Field, f.Confidence))
	}
	if err := recordLayerField(pkgID, "inferred_suggestions", f.Field, seen); err != nil {
		errs = append(errs, err)
	}
	if strings.TrimSpace(f.Value) == "" {
		errs = append(errs, fmt.Errorf("entity package %q: inferred_suggestions field %q has empty value", pkgID, f.Field))
	}
	if strings.TrimSpace(f.Rationale) == "" {
		errs = append(errs, fmt.Errorf("entity package %q: inferred_suggestions field %q has empty rationale", pkgID, f.Field))
	}
	return errs
}

func validateUnresolvedFields(pkgID string, items []UnresolvedField, seen map[string]string) []error {
	var errs []error
	for _, f := range items {
		if err := recordLayerField(pkgID, "unresolved_fields", f.Field, seen); err != nil {
			errs = append(errs, err)
		}
		if strings.TrimSpace(f.Reason) == "" {
			errs = append(errs, fmt.Errorf("entity package %q: unresolved_fields field %q has empty reason", pkgID, f.Field))
		}
	}
	return errs
}

// ValidateVisualObservation checks a VisualObservation against visual-observation.schema.yaml:
// appearance_facts must exclude non-visual fields, non_authoritative_note is required, and confidence is bounded.
func ValidateVisualObservation(o VisualObservation) []error {
	errs := validateVisualHeader(o)
	for _, f := range o.AppearanceFacts {
		errs = append(errs, validateAppearanceFact(o.ID, f)...)
	}
	return errs
}

func validateAppearanceFact(obsID string, f AppearanceFact) []error {
	var errs []error
	if strings.TrimSpace(f.Field) == "" {
		errs = append(errs, fmt.Errorf("visual observation %q: appearance_facts entry has empty field", obsID))
	} else if nonVisualFields[normalizeFieldName(f.Field)] {
		errs = append(errs, fmt.Errorf("visual observation %q: appearance_facts contains non-visual field %q — this belongs in an entity-package, not a visual observation", obsID, f.Field))
	}
	if strings.TrimSpace(f.Value) == "" {
		errs = append(errs, fmt.Errorf("visual observation %q: appearance fact %q has empty value", obsID, f.Field))
	}
	if !validConfidence[f.Confidence] {
		errs = append(errs, fmt.Errorf("visual observation %q: appearance fact %q has invalid confidence %q", obsID, f.Field, f.Confidence))
	}
	return errs
}

func validateVisualHeader(o VisualObservation) []error {
	var errs []error
	if err := validatePackageID("visual observation", o.ID); err != nil {
		errs = append(errs, err)
	}
	if !validStatus[o.Status] {
		errs = append(errs, fmt.Errorf("visual observation %q: invalid status %q", o.ID, o.Status))
	}
	if strings.TrimSpace(o.SourceRef) == "" {
		errs = append(errs, fmt.Errorf("visual observation %q: source_ref must not be empty", o.ID))
	}
	if strings.TrimSpace(o.NonAuthoritativeNote) == "" {
		errs = append(errs, fmt.Errorf("visual observation %q: non_authoritative_note must be present", o.ID))
	}
	if len(o.AppearanceFacts) == 0 {
		errs = append(errs, fmt.Errorf("visual observation %q: appearance_facts must not be empty", o.ID))
	}
	return errs
}
