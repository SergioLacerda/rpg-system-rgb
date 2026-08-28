package fixtures

import (
	"encoding/json"
	"os"
	"testing"
)

// semanticUnitOracle is the golden expectation for one canonical semantic
// unit: its authority/status classification stays stable and it remains
// covered by at least one projection, per design.md's golden-layer
// guidance for "canonical documentation ... semantic IDs, projection
// manifests". This is a behavior-level oracle over a curated ID list,
// distinct from the whole-index structural validators in
// internal/components/tooling.
type semanticUnitOracle struct {
	id            string
	authorityType string
	sourceStatus  string
}

var goldenSemanticUnits = []semanticUnitOracle{
	{id: "core.combat.attack-margin", authorityType: "canonical_semantic", sourceStatus: "canonical"},
	{id: "core.damage.flow", authorityType: "canonical_semantic", sourceStatus: "canonical"},
	{id: "core.damage.armor-reduction", authorityType: "canonical_semantic", sourceStatus: "canonical"},
	{id: "core.damage.shield-absorption", authorityType: "canonical_semantic", sourceStatus: "canonical"},
}

type semanticIndexUnit struct {
	ID            string `json:"id"`
	AuthorityType string `json:"authority_type"`
	SourceStatus  string `json:"source_status"`
}

func loadSemanticIndexUnits(t *testing.T) map[string]semanticIndexUnit {
	t.Helper()
	body, err := os.ReadFile("../../docs/core/semantic/core-v2.index.json")
	if err != nil {
		t.Fatal(err)
	}
	var index struct {
		Units []semanticIndexUnit `json:"units"`
	}
	if err := json.Unmarshal(body, &index); err != nil {
		t.Fatal(err)
	}
	units := make(map[string]semanticIndexUnit, len(index.Units))
	for _, unit := range index.Units {
		units[unit.ID] = unit
	}
	return units
}

type consumerContract struct {
	Component               string   `json:"component"`
	ForbiddenAuthorityTypes []string `json:"forbidden_authority_types"`
}

// loadConsumerContractForbiddenTypes returns, per component, the set of
// authority types that component's consumer contract forbids ingesting.
func loadConsumerContractForbiddenTypes(t *testing.T) map[string]map[string]bool {
	t.Helper()
	body, err := os.ReadFile("../../docs/core/semantic/consumer-contracts.v0.1.json")
	if err != nil {
		t.Fatal(err)
	}
	var doc struct {
		Contracts []consumerContract `json:"contracts"`
	}
	if err := json.Unmarshal(body, &doc); err != nil {
		t.Fatal(err)
	}
	forbidden := make(map[string]map[string]bool, len(doc.Contracts))
	for _, contract := range doc.Contracts {
		types := make(map[string]bool, len(contract.ForbiddenAuthorityTypes))
		for _, authorityType := range contract.ForbiddenAuthorityTypes {
			types[authorityType] = true
		}
		forbidden[contract.Component] = types
	}
	return forbidden
}

func loadProjectionSourceUnits(t *testing.T) map[string]bool {
	t.Helper()
	body, err := os.ReadFile("../../docs/core/semantic/projection-manifest.v0.1.json")
	if err != nil {
		t.Fatal(err)
	}
	var manifest struct {
		Projections []struct {
			SourceUnits []string `json:"source_units"`
		} `json:"projections"`
	}
	if err := json.Unmarshal(body, &manifest); err != nil {
		t.Fatal(err)
	}
	covered := map[string]bool{}
	for _, projection := range manifest.Projections {
		for _, id := range projection.SourceUnits {
			covered[id] = true
		}
	}
	return covered
}

func TestGoldenSemanticUnitsKeepStableProjectionCoverage(t *testing.T) {
	indexUnits := loadSemanticIndexUnits(t)
	projectionCoverage := loadProjectionSourceUnits(t)

	for _, want := range goldenSemanticUnits {
		t.Run(want.id, func(t *testing.T) {
			unit, ok := indexUnits[want.id]
			if !ok {
				t.Fatalf("semantic id %q no longer exists in core-v2.index.json", want.id)
			}
			if unit.AuthorityType != want.authorityType {
				t.Fatalf("%s: authority_type = %q, want %q", want.id, unit.AuthorityType, want.authorityType)
			}
			if unit.SourceStatus != want.sourceStatus {
				t.Fatalf("%s: source_status = %q, want %q", want.id, unit.SourceStatus, want.sourceStatus)
			}
			if !projectionCoverage[want.id] {
				t.Fatalf("%s: no projection in projection-manifest.v0.1.json lists it as a source_unit", want.id)
			}
		})
	}
}

// TestGoldenSemanticUnitsRemainConsumable guards against a consumer contract
// accidentally forbidding the authority_type these canonical golden units
// carry (e.g. Core itself losing the ability to ingest canonical_semantic
// content it owns), per design.md's golden-layer guidance to extend
// coverage into docs/core/semantic/consumer-contracts.v0.1.json.
func TestGoldenSemanticUnitsRemainConsumable(t *testing.T) {
	indexUnits := loadSemanticIndexUnits(t)
	forbiddenByComponent := loadConsumerContractForbiddenTypes(t)
	if len(forbiddenByComponent) == 0 {
		t.Fatal("no consumer contracts loaded from consumer-contracts.v0.1.json")
	}

	for _, want := range goldenSemanticUnits {
		t.Run(want.id, func(t *testing.T) {
			unit, ok := indexUnits[want.id]
			if !ok {
				t.Fatalf("semantic id %q no longer exists in core-v2.index.json", want.id)
			}
			for component, forbiddenTypes := range forbiddenByComponent {
				if forbiddenTypes[unit.AuthorityType] {
					t.Fatalf("%s: consumer contract %q forbids authority_type %q, which this canonical unit carries", want.id, component, unit.AuthorityType)
				}
			}
		})
	}
}
