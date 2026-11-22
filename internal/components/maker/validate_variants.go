package maker

import (
	"fmt"
	"strings"
)

func validateEntityVariants(pkgID string, variants []EntityVariant, sourceSet map[string]bool) []error {
	var errs []error
	for i, v := range variants {
		errs = append(errs, validateVariantItem(pkgID, i, v, sourceSet)...)
	}
	return errs
}

func validateVariantItem(pkgID string, idx int, v EntityVariant, sourceSet map[string]bool) []error {
	var errs []error
	if strings.TrimSpace(v.Label) == "" {
		errs = append(errs, fmt.Errorf("entity package %q: variants[%d] has empty label", pkgID, idx))
	}
	if strings.TrimSpace(v.SourceRef) == "" {
		errs = append(errs, fmt.Errorf("entity package %q: variants[%d] has empty source_ref", pkgID, idx))
	} else if len(sourceSet) > 0 && !sourceSet[v.SourceRef] {
		errs = append(errs, fmt.Errorf("entity package %q: variant %q references undeclared source_ref %q", pkgID, v.Label, v.SourceRef))
	}
	if strings.TrimSpace(v.Summary) == "" {
		errs = append(errs, fmt.Errorf("entity package %q: variant %q has empty summary", pkgID, v.Label))
	}
	return append(errs, validateVariantFacts(pkgID, v.Label, v.Facts)...)
}

func validateVariantFacts(pkgID, label string, facts []ExplicitFact) []error {
	var errs []error
	seen := make(map[string]bool, len(facts))
	for _, f := range facts {
		errs = append(errs, validateVariantFact(pkgID, label, f, seen)...)
	}
	return errs
}

func validateVariantFact(pkgID, label string, f ExplicitFact, seen map[string]bool) []error {
	var errs []error
	if strings.TrimSpace(f.Field) == "" {
		errs = append(errs, fmt.Errorf("entity package %q: variant %q has fact with empty field", pkgID, label))
	} else if seen[f.Field] {
		errs = append(errs, fmt.Errorf("entity package %q: variant %q has duplicate field %q", pkgID, label, f.Field))
	} else {
		seen[f.Field] = true
	}
	if strings.TrimSpace(f.Value) == "" {
		errs = append(errs, fmt.Errorf("entity package %q: variant %q fact %q has empty value", pkgID, label, f.Field))
	}
	return errs
}
