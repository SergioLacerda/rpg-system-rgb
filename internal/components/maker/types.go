// Package maker implements the RGB Maker component boundary: input routing,
// schema validation, and report aggregation for campaign/content structuring.
//
// The LLM instruction layer (skills/maker/**) handles semantic extraction and
// classification. This package handles deterministic behavior only: schema
// validation, field-invariant checking, structured-fact diffing, and
// aggregation. See ADR-015 for the runtime split rationale.
package maker

// --- Entity Package (entity-package.schema.yaml) ---

// EntityPackage is one extracted entity (NPC, location, item, faction)
// produced by the Maker LLM layer. Mirrors rgb-maker-entity-package/0.1.
type EntityPackage struct {
	ID                  string               `json:"id"`
	EntityType          string               `json:"entity_type"`
	Status              string               `json:"status"`
	SourceRefs          []string             `json:"source_refs"`
	ExplicitFacts       []ExplicitFact       `json:"explicit_facts"`
	InferredSuggestions []InferredSuggestion `json:"inferred_suggestions,omitempty"`
	UnresolvedFields    []UnresolvedField    `json:"unresolved_fields,omitempty"`
	Conflicts           []string             `json:"conflicts,omitempty"`
	CanonicalDeps       []string             `json:"canonical_dependencies,omitempty"`
	Variants            []EntityVariant      `json:"variants,omitempty"`
}

// EntityVariant captures one alternative version, reboot, or remake of an
// entity preserved alongside shared facts instead of silently overwriting them.
type EntityVariant struct {
	Label     string         `json:"label"`
	SourceRef string         `json:"source_ref"`
	Summary   string         `json:"summary"`
	Facts     []ExplicitFact `json:"facts,omitempty"`
}

// ExplicitFact is a fact stated directly in the source text.
type ExplicitFact struct {
	Field     string `json:"field"`
	Value     string `json:"value"`
	SourceRef string `json:"source_ref"`
}

// InferredSuggestion is a plausible fact Maker inferred but the source did
// not state directly.
type InferredSuggestion struct {
	Field      string `json:"field"`
	Value      string `json:"value"`
	Confidence string `json:"confidence"` // low | medium | high
	Rationale  string `json:"rationale"`
}

// UnresolvedField is a field the entity type normally has but the source
// did not provide and Maker did not infer.
type UnresolvedField struct {
	Field  string `json:"field"`
	Reason string `json:"reason"`
}

// --- Visual Observation (visual-observation.schema.yaml) ---

// VisualObservation records appearance-only facts read from an image.
// Per the authority rule (SKILL.md): image = appearance, text = facts.
// Mirrors rgb-maker-visual-observation/0.1.
type VisualObservation struct {
	ID                   string           `json:"id"`
	Status               string           `json:"status"`
	SourceRef            string           `json:"source_ref"`
	AppearanceFacts      []AppearanceFact `json:"appearance_facts"`
	LinkedEntity         string           `json:"linked_entity,omitempty"`
	NonAuthoritativeNote string           `json:"non_authoritative_note"`
}

// AppearanceFact is one visual property read from an image.
type AppearanceFact struct {
	Field      string `json:"field"`
	Value      string `json:"value"`
	Confidence string `json:"confidence"` // low | medium | high
}

// --- Conflict Report (conflict-report.schema.yaml) ---

// ConflictReport covers exactly one field in dispute between two or more
// sources that Maker cannot resolve on its own. Mirrors
// rgb-maker-conflict-report/0.1.
type ConflictReport struct {
	ID             string  `json:"id"`
	Status         string  `json:"status"`
	Field          string  `json:"field"`
	LinkedEntity   string  `json:"linked_entity,omitempty"`
	Claims         []Claim `json:"claims"`
	Resolution     string  `json:"resolution"` // unresolved | resolved_by_authority_rule | resolved_by_user
	ResolutionNote string  `json:"resolution_note,omitempty"`
}

// Claim is one competing value in a conflict report.
type Claim struct {
	Value      string `json:"value"`
	SourceRef  string `json:"source_ref"`
	SourceKind string `json:"source_kind"` // text | image
}

// --- Maker Report (maker-report.schema.yaml) ---

// MakerReport is the session-level summary of everything Maker extracted
// from one input batch. Mirrors rgb-maker-report/0.1.
type MakerReport struct { //nolint:revive // MakerReport mirrors the maker-report schema name distinct from ConflictReport.
	ID                 string   `json:"id"`
	Status             string   `json:"status"`
	InputBatch         []string `json:"input_batch"`
	EntityPackages     []string `json:"entity_packages"`
	ConflictReports    []string `json:"conflict_reports,omitempty"`
	VisualObservations []string `json:"visual_observations,omitempty"`
	OpenQuestions      []string `json:"open_questions,omitempty"`
	PendingDecisions   []string `json:"pending_decisions,omitempty"`
}
