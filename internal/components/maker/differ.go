package maker

import (
	"fmt"
	"strings"
)

type entityFieldKey struct {
	entityID string
	field    string
}

// DiffEntityPackages detects structured-fact conflicts across EntityPackage
// drafts sharing the same entity ID and explicit fact field name.
//
// When two or more distinct values are recorded for the same (entity, field):
//   - If claims span both text and image on a non-visual fact field, the
//     conflict is marked resolved_by_authority_rule per SKILL.md.
//   - Otherwise, resolution remains "unresolved" for the GM to decide.
func DiffEntityPackages(pkgs []EntityPackage) []ConflictReport {
	order, grouped := collectEntityClaims(pkgs)
	var reports []ConflictReport
	for _, key := range order {
		claims := grouped[key]
		if !hasDistinctValues(claims) {
			continue
		}
		reports = append(reports, buildConflictReport(key, claims))
	}
	return reports
}

func collectEntityClaims(pkgs []EntityPackage) ([]entityFieldKey, map[entityFieldKey][]Claim) {
	var order []entityFieldKey
	grouped := make(map[entityFieldKey][]Claim)
	seenKeys := make(map[entityFieldKey]bool)

	for _, p := range pkgs {
		if strings.TrimSpace(p.ID) == "" {
			continue
		}
		recordPackageFacts(p, &order, grouped, seenKeys)
	}
	return order, grouped
}

func recordPackageFacts(
	p EntityPackage,
	order *[]entityFieldKey,
	grouped map[entityFieldKey][]Claim,
	seenKeys map[entityFieldKey]bool,
) {
	for _, f := range p.ExplicitFacts {
		if strings.TrimSpace(f.Field) == "" || strings.TrimSpace(f.Value) == "" {
			continue
		}
		key := entityFieldKey{entityID: p.ID, field: f.Field}
		if !seenKeys[key] {
			seenKeys[key] = true
			*order = append(*order, key)
		}
		claim := Claim{
			Value:      f.Value,
			SourceRef:  f.SourceRef,
			SourceKind: string(detectType(f.SourceRef)),
		}
		grouped[key] = appendUniqueClaim(grouped[key], claim)
	}
}

func appendUniqueClaim(claims []Claim, candidate Claim) []Claim {
	for _, c := range claims {
		if c.Value == candidate.Value && c.SourceRef == candidate.SourceRef {
			return claims
		}
	}
	return append(claims, candidate)
}

func hasDistinctValues(claims []Claim) bool {
	if len(claims) < 2 {
		return false
	}
	first := claims[0].Value
	for _, c := range claims[1:] {
		if c.Value != first {
			return true
		}
	}
	return false
}

func isAuthorityRuleFactField(field string) bool {
	return nonVisualFields[normalizeFieldName(field)]
}

func buildConflictReport(key entityFieldKey, claims []Claim) ConflictReport {
	report := ConflictReport{
		ID:           fmt.Sprintf("conflict.%s-%s", slugify(key.entityID), slugify(key.field)),
		Status:       "needs_review",
		Field:        key.field,
		LinkedEntity: key.entityID,
		Claims:       claims,
		Resolution:   "unresolved",
	}
	if hasClaimKind(claims, "text") && hasClaimKind(claims, "image") && isAuthorityRuleFactField(key.field) {
		report.Resolution = "resolved_by_authority_rule"
		report.ResolutionNote = "text source is authoritative over non-visual facts per SKILL.md authority rule"
	}
	return report
}
