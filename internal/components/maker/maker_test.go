package maker

import (
	"strings"
	"testing"

	"github.com/SergioLacerda/rpg-system-rgb/internal/components"
)

func TestDescriptor(t *testing.T) {
	descriptor := Descriptor()
	if descriptor.ID != "maker" || descriptor.Name == "" || descriptor.Description == "" {
		t.Fatalf("unexpected descriptor: %+v", descriptor)
	}
}

func validSampleEntityPackage() EntityPackage {
	return EntityPackage{
		ID:         "npc.harbor-quartermaster",
		EntityType: "npc",
		Status:     "needs_review",
		SourceRefs: []string{"Note 1 (session prep)", "Note 2 (session recap)", "harbor-npc-sketch.png"},
		ExplicitFacts: []ExplicitFact{
			{Field: "role", Value: "harbor quartermaster", SourceRef: "Note 1 (session prep)"},
			{Field: "missing_fingers", Value: "two fingers on left hand", SourceRef: "Note 1 (session prep)"},
		},
		InferredSuggestions: []InferredSuggestion{
			{Field: "attitude_toward_party", Value: "transactional", Confidence: "medium", Rationale: "dock access requires earning favor"},
		},
		UnresolvedFields: []UnresolvedField{
			{Field: "name", Reason: "disputed — see conflict.harbor-quartermaster-name"},
		},
		Conflicts: []string{"conflict.harbor-quartermaster-name"},
	}
}

func TestValidateEntityPackage(t *testing.T) {
	if errs := ValidateEntityPackage(validSampleEntityPackage()); len(errs) != 0 {
		t.Fatalf("expected valid sample package to pass, got: %v", errs)
	}

	cases := []struct {
		name   string
		mutate func(*EntityPackage)
		want   string
	}{
		{
			name:   "empty id",
			mutate: func(p *EntityPackage) { p.ID = "" },
			want:   "id must not be empty",
		},
		{
			name:   "invalid entity_type",
			mutate: func(p *EntityPackage) { p.EntityType = "spell" },
			want:   "invalid entity_type",
		},
		{
			name:   "canonical status rejected",
			mutate: func(p *EntityPackage) { p.Status = "canonical" },
			want:   "invalid status",
		},
		{
			name:   "empty source_refs",
			mutate: func(p *EntityPackage) { p.SourceRefs = nil },
			want:   "source_refs must not be empty",
		},
		{
			name:   "empty item in source_refs",
			mutate: func(p *EntityPackage) { p.SourceRefs = append(p.SourceRefs, "  ") },
			want:   "source_refs[3] must not be empty",
		},
		{
			name:   "empty item in conflicts",
			mutate: func(p *EntityPackage) { p.Conflicts = []string{""} },
			want:   "conflicts[0] must not be empty",
		},
		{
			name:   "empty item in canonical_dependencies",
			mutate: func(p *EntityPackage) { p.CanonicalDeps = []string{" "} },
			want:   "canonical_dependencies[0] must not be empty",
		},
		{
			name:   "empty explicit_facts",
			mutate: func(p *EntityPackage) { p.ExplicitFacts = nil },
			want:   "explicit_facts must not be empty",
		},
		{
			name: "empty explicit fact field",
			mutate: func(p *EntityPackage) {
				p.ExplicitFacts = append(p.ExplicitFacts, ExplicitFact{Field: "", Value: "x", SourceRef: "Note 1 (session prep)"})
			},
			want: "explicit_facts entry has empty field",
		},
		{
			name: "empty explicit fact value",
			mutate: func(p *EntityPackage) {
				p.ExplicitFacts[0].Value = "   "
			},
			want: "has empty value",
		},
		{
			name: "empty explicit fact source_ref",
			mutate: func(p *EntityPackage) {
				p.ExplicitFacts[0].SourceRef = ""
			},
			want: "has empty source_ref",
		},
		{
			name: "undeclared explicit fact source_ref",
			mutate: func(p *EntityPackage) {
				p.ExplicitFacts[0].SourceRef = "Unknown Note"
			},
			want: "references undeclared source_ref",
		},
		{
			name: "empty inferred suggestion field",
			mutate: func(p *EntityPackage) {
				p.InferredSuggestions = append(p.InferredSuggestions, InferredSuggestion{Field: "", Value: "x", Confidence: "low", Rationale: "y"})
			},
			want: "inferred_suggestions entry has empty field",
		},
		{
			name: "invalid confidence in inferred suggestion",
			mutate: func(p *EntityPackage) {
				p.InferredSuggestions[0].Confidence = "certain"
			},
			want: "invalid confidence",
		},
		{
			name: "empty inferred suggestion value",
			mutate: func(p *EntityPackage) {
				p.InferredSuggestions[0].Value = ""
			},
			want: "inferred_suggestions field \"attitude_toward_party\" has empty value",
		},
		{
			name: "empty inferred suggestion rationale",
			mutate: func(p *EntityPackage) {
				p.InferredSuggestions[0].Rationale = " "
			},
			want: "has empty rationale",
		},
		{
			name: "empty unresolved field",
			mutate: func(p *EntityPackage) {
				p.UnresolvedFields = append(p.UnresolvedFields, UnresolvedField{Field: "", Reason: "missing"})
			},
			want: "unresolved_fields entry has empty field",
		},
		{
			name: "empty unresolved field reason",
			mutate: func(p *EntityPackage) {
				p.UnresolvedFields[0].Reason = "  "
			},
			want: "has empty reason",
		},
		{
			name: "duplicate field in explicit and inferred",
			mutate: func(p *EntityPackage) {
				p.InferredSuggestions = append(p.InferredSuggestions, InferredSuggestion{
					Field: "role", Value: "clerk", Confidence: "low", Rationale: "dup",
				})
			},
			want: "appears in both explicit_facts and inferred_suggestions",
		},
		{
			name: "duplicate field in inferred and unresolved",
			mutate: func(p *EntityPackage) {
				p.UnresolvedFields = append(p.UnresolvedFields, UnresolvedField{
					Field: "attitude_toward_party", Reason: "dup",
				})
			},
			want: "appears in both inferred_suggestions and unresolved_fields",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			pkg := validSampleEntityPackage()
			tc.mutate(&pkg)
			errs := ValidateEntityPackage(pkg)
			if !containsError(errs, tc.want) {
				t.Fatalf("expected error containing %q, got %v", tc.want, errs)
			}
		})
	}
}

func validSampleVisualObservation() VisualObservation {
	return VisualObservation{
		ID:        "observation.harbor-npc-sketch",
		Status:    "draft",
		SourceRef: "harbor-npc-sketch.png",
		AppearanceFacts: []AppearanceFact{
			{Field: "clothing", Value: "weathered coat", Confidence: "high"},
			{Field: "hands", Value: "both hands appear undamaged", Confidence: "medium"},
		},
		LinkedEntity:         "npc.harbor-quartermaster",
		NonAuthoritativeNote: "This observation may describe appearance only.",
	}
}

func TestValidateVisualObservation(t *testing.T) {
	if errs := ValidateVisualObservation(validSampleVisualObservation()); len(errs) != 0 {
		t.Fatalf("expected valid visual observation to pass, got: %v", errs)
	}

	cases := []struct {
		name   string
		mutate func(*VisualObservation)
		want   string
	}{
		{
			name:   "empty id",
			mutate: func(o *VisualObservation) { o.ID = "" },
			want:   "id must not be empty",
		},
		{
			name:   "invalid status",
			mutate: func(o *VisualObservation) { o.Status = "canonical" },
			want:   "invalid status",
		},
		{
			name:   "empty source_ref",
			mutate: func(o *VisualObservation) { o.SourceRef = "" },
			want:   "source_ref must not be empty",
		},
		{
			name:   "empty non_authoritative_note",
			mutate: func(o *VisualObservation) { o.NonAuthoritativeNote = "" },
			want:   "non_authoritative_note must be present",
		},
		{
			name:   "empty appearance_facts",
			mutate: func(o *VisualObservation) { o.AppearanceFacts = nil },
			want:   "appearance_facts must not be empty",
		},
		{
			name: "empty appearance fact field",
			mutate: func(o *VisualObservation) {
				o.AppearanceFacts = []AppearanceFact{{Field: "  ", Value: "cloak", Confidence: "low"}}
			},
			want: "appearance_facts entry has empty field",
		},
		{
			name: "empty appearance fact value",
			mutate: func(o *VisualObservation) {
				o.AppearanceFacts = []AppearanceFact{{Field: "clothing", Value: " ", Confidence: "low"}}
			},
			want: "has empty value",
		},
		{
			name: "non-visual field loyalty rejected",
			mutate: func(o *VisualObservation) {
				o.AppearanceFacts = []AppearanceFact{{Field: "loyalty", Value: "guild", Confidence: "low"}}
			},
			want: "non-visual field \"loyalty\"",
		},
		{
			name: "non-visual field Loyalty normalized and rejected",
			mutate: func(o *VisualObservation) {
				o.AppearanceFacts = []AppearanceFact{{Field: "Loyalty", Value: "guild", Confidence: "low"}}
			},
			want: "non-visual field \"Loyalty\"",
		},
		{
			name: "non-visual field true-identity normalized and rejected",
			mutate: func(o *VisualObservation) {
				o.AppearanceFacts = []AppearanceFact{{Field: "true-identity", Value: "spy", Confidence: "low"}}
			},
			want: "non-visual field \"true-identity\"",
		},
		{
			name: "non-visual field name rejected",
			mutate: func(o *VisualObservation) {
				o.AppearanceFacts = []AppearanceFact{{Field: "name", Value: "Yssa", Confidence: "low"}}
			},
			want: "non-visual field \"name\"",
		},
		{
			name: "invalid confidence rejected",
			mutate: func(o *VisualObservation) {
				o.AppearanceFacts[0].Confidence = "sure"
			},
			want: "invalid confidence",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			obs := validSampleVisualObservation()
			tc.mutate(&obs)
			errs := ValidateVisualObservation(obs)
			if !containsError(errs, tc.want) {
				t.Fatalf("expected error containing %q, got %v", tc.want, errs)
			}
		})
	}
}

func validSampleConflictReport() ConflictReport {
	return ConflictReport{
		ID:           "conflict.harbor-quartermaster-name",
		Status:       "needs_review",
		Field:        "name",
		LinkedEntity: "npc.harbor-quartermaster",
		Claims: []Claim{
			{Value: "Yssa Kell", SourceRef: "Note 1 (session prep)", SourceKind: "text"},
			{Value: "Yssa Voss", SourceRef: "Note 2 (session recap)", SourceKind: "text"},
		},
		Resolution: "unresolved",
	}
}

func TestValidateConflictReport(t *testing.T) {
	if errs := ValidateConflictReport(validSampleConflictReport()); len(errs) != 0 {
		t.Fatalf("expected valid unresolved conflict report to pass, got: %v", errs)
	}

	authorityResolved := ConflictReport{
		ID:     "conflict.harbor-quartermaster-injury",
		Status: "needs_review",
		Field:  "missing_fingers",
		Claims: []Claim{
			{Value: "two fingers missing", SourceRef: "Note 1", SourceKind: "text"},
			{Value: "undamaged hands", SourceRef: "harbor-npc-sketch.png", SourceKind: "image"},
		},
		Resolution:     "resolved_by_authority_rule",
		ResolutionNote: "text is authoritative for injury history; image covers appearance only",
	}
	if errs := ValidateConflictReport(authorityResolved); len(errs) != 0 {
		t.Fatalf("expected valid authority-rule conflict report to pass, got: %v", errs)
	}

	userResolved := validSampleConflictReport()
	userResolved.Resolution = "resolved_by_user"
	userResolved.ResolutionNote = "GM confirmed Yssa Voss"
	if errs := ValidateConflictReport(userResolved); len(errs) != 0 {
		t.Fatalf("expected valid user-resolved conflict report to pass, got: %v", errs)
	}

	cases := []struct {
		name   string
		mutate func(*ConflictReport)
		want   string
	}{
		{
			name:   "empty id",
			mutate: func(r *ConflictReport) { r.ID = "" },
			want:   "id must not be empty",
		},
		{
			name:   "empty field",
			mutate: func(r *ConflictReport) { r.Field = "" },
			want:   "field must not be empty",
		},
		{
			name:   "invalid status",
			mutate: func(r *ConflictReport) { r.Status = "canonical" },
			want:   "invalid status",
		},
		{
			name:   "invalid resolution",
			mutate: func(r *ConflictReport) { r.Resolution = "auto" },
			want:   "invalid resolution",
		},
		{
			name: "fewer than 2 claims",
			mutate: func(r *ConflictReport) {
				r.Claims = r.Claims[:1]
			},
			want: "at least 2 claims",
		},
		{
			name: "empty claim value",
			mutate: func(r *ConflictReport) {
				r.Claims[0].Value = " "
			},
			want: "claim[0] has empty value",
		},
		{
			name: "empty claim source_ref",
			mutate: func(r *ConflictReport) {
				r.Claims[1].SourceRef = ""
			},
			want: "claim[1] has empty source_ref",
		},
		{
			name: "invalid source_kind",
			mutate: func(r *ConflictReport) {
				r.Claims[1].SourceKind = "audio"
			},
			want: "invalid source_kind",
		},
		{
			name: "resolved_by_authority_rule without image claim",
			mutate: func(r *ConflictReport) {
				r.Resolution = "resolved_by_authority_rule"
				r.ResolutionNote = "tried authority rule on text vs text"
			},
			want: "claims must include both text and image",
		},
		{
			name: "resolved_by_authority_rule image vs image rejected",
			mutate: func(r *ConflictReport) {
				r.Claims[0].SourceKind = "image"
				r.Claims[1].SourceKind = "image"
				r.Resolution = "resolved_by_authority_rule"
				r.ResolutionNote = "tried authority rule on image vs image"
			},
			want: "claims must include both text and image",
		},
		{
			name: "resolved_by_authority_rule without resolution_note",
			mutate: func(r *ConflictReport) {
				r.Claims[1].SourceKind = "image"
				r.Resolution = "resolved_by_authority_rule"
				r.ResolutionNote = ""
			},
			want: "resolution_note is required when resolution is resolved_by_authority_rule",
		},
		{
			name: "resolved_by_user without resolution_note",
			mutate: func(r *ConflictReport) {
				r.Resolution = "resolved_by_user"
				r.ResolutionNote = ""
			},
			want: "resolution_note is required when resolution is resolved_by_user",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			rep := validSampleConflictReport()
			tc.mutate(&rep)
			errs := ValidateConflictReport(rep)
			if !containsError(errs, tc.want) {
				t.Fatalf("expected error containing %q, got %v", tc.want, errs)
			}
		})
	}
}

func TestValidateMakerReport(t *testing.T) {
	valid := MakerReport{
		ID:                 "report.2026-08-01-harbor-session",
		Status:             "draft",
		InputBatch:         []string{"Note 1"},
		EntityPackages:     []string{"npc.harbor-quartermaster"},
		ConflictReports:    []string{"conflict.harbor-quartermaster-name"},
		VisualObservations: []string{"observation.harbor-npc-sketch"},
		OpenQuestions:      []string{"Confirm Yssa's surname."},
	}
	if errs := ValidateMakerReport(valid); len(errs) != 0 {
		t.Fatalf("expected valid maker report to pass, got %v", errs)
	}

	emptyEntities := MakerReport{
		ID:             "report.empty",
		Status:         "draft",
		InputBatch:     []string{"Note 1"},
		EntityPackages: []string{},
	}
	if errs := ValidateMakerReport(emptyEntities); len(errs) != 0 {
		t.Fatalf("expected empty non-nil entity_packages to pass, got %v", errs)
	}

	cases := []struct {
		name   string
		mutate func(*MakerReport)
		want   string
	}{
		{
			name:   "empty id",
			mutate: func(r *MakerReport) { r.ID = "" },
			want:   "id must not be empty",
		},
		{
			name:   "invalid status",
			mutate: func(r *MakerReport) { r.Status = "canonical" },
			want:   "invalid status",
		},
		{
			name:   "empty input_batch",
			mutate: func(r *MakerReport) { r.InputBatch = nil },
			want:   "input_batch must not be empty",
		},
		{
			name:   "empty input_batch element",
			mutate: func(r *MakerReport) { r.InputBatch = []string{" "} },
			want:   "input_batch[0] must not be empty",
		},
		{
			name:   "nil entity_packages",
			mutate: func(r *MakerReport) { r.EntityPackages = nil },
			want:   "entity_packages must not be nil",
		},
		{
			name:   "empty entity_packages item",
			mutate: func(r *MakerReport) { r.EntityPackages = []string{""} },
			want:   "entity_packages[0] must not be empty",
		},
		{
			name:   "duplicate entity_packages item",
			mutate: func(r *MakerReport) { r.EntityPackages = []string{"npc.a", "npc.a"} },
			want:   "duplicate entity_packages id",
		},
		{
			name:   "empty open_questions item",
			mutate: func(r *MakerReport) { r.OpenQuestions = []string{""} },
			want:   "open_questions[0] must not be empty",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			rep := valid
			tc.mutate(&rep)
			errs := ValidateMakerReport(rep)
			if !containsError(errs, tc.want) {
				t.Fatalf("expected error containing %q, got %v", tc.want, errs)
			}
		})
	}
}

func TestValidateCanonicalDependencies(t *testing.T) {
	idx := components.SemanticIndex{
		Units: []components.SemanticIndexUnit{
			{ID: "core.ability.contract", AuthorityType: "canonical_semantic"},
			{ID: "core.damage.flow", AuthorityType: "canonical_markdown_bridge"},
			{ID: "projection.pack", AuthorityType: "projection"},
			{ID: "generated.pack", AuthorityType: "generated_artifact"},
		},
	}

	pkg := validSampleEntityPackage()
	pkg.CanonicalDeps = []string{"core.ability.contract", "core.damage.flow"}
	if errs := ValidateCanonicalDependencies(pkg, idx); len(errs) != 0 {
		t.Fatalf("expected canonical dependencies to pass, got %v", errs)
	}

	pkg.CanonicalDeps = []string{"core.missing", "projection.pack", "generated.pack"}
	errs := ValidateCanonicalDependencies(pkg, idx)
	if !containsError(errs, "unknown semantic ID \"core.missing\"") {
		t.Fatalf("expected unknown semantic ID error, got %v", errs)
	}
	if !containsError(errs, "forbidden authority_type \"projection\"") {
		t.Fatalf("expected forbidden projection error, got %v", errs)
	}
	if !containsError(errs, "forbidden authority_type \"generated_artifact\"") {
		t.Fatalf("expected forbidden generated_artifact error, got %v", errs)
	}
}

func TestDiffEntityPackages(t *testing.T) {
	matching := []EntityPackage{
		{
			ID:         "npc.yssa",
			EntityType: "npc",
			Status:     "draft",
			SourceRefs: []string{"note-1.txt"},
			ExplicitFacts: []ExplicitFact{
				{Field: "role", Value: "quartermaster", SourceRef: "note-1.txt"},
			},
		},
		{
			ID:         "npc.yssa",
			EntityType: "npc",
			Status:     "draft",
			SourceRefs: []string{"note-2.txt"},
			ExplicitFacts: []ExplicitFact{
				{Field: "role", Value: "quartermaster", SourceRef: "note-2.txt"},
			},
		},
	}
	if reports := DiffEntityPackages(matching); len(reports) != 0 {
		t.Fatalf("expected no conflicts when explicit facts agree, got %+v", reports)
	}

	textDispute := []EntityPackage{
		{
			ID:         "npc.yssa",
			EntityType: "npc",
			Status:     "draft",
			SourceRefs: []string{"note-1.txt"},
			ExplicitFacts: []ExplicitFact{
				{Field: "name", Value: "Yssa Kell", SourceRef: "note-1.txt"},
				{Field: "", Value: "ignored", SourceRef: "note-1.txt"},
			},
		},
		{
			ID: "", // skipped
		},
		{
			ID:         "npc.yssa",
			EntityType: "npc",
			Status:     "draft",
			SourceRefs: []string{"note-2.txt"},
			ExplicitFacts: []ExplicitFact{
				{Field: "name", Value: "Yssa Voss", SourceRef: "note-2.txt"},
				{Field: "name", Value: "Yssa Voss", SourceRef: "note-2.txt"}, // deduplicated
			},
		},
	}
	reports := DiffEntityPackages(textDispute)
	if len(reports) != 1 {
		t.Fatalf("expected 1 conflict report, got %d (%+v)", len(reports), reports)
	}
	if reports[0].Resolution != "unresolved" || len(reports[0].Claims) != 2 {
		t.Fatalf("unexpected text-vs-text conflict report: %+v", reports[0])
	}
	if errs := ValidateConflictReport(reports[0]); len(errs) != 0 {
		t.Fatalf("generated conflict report failed validation: %v", errs)
	}

	authorityDispute := []EntityPackage{
		{
			ID:         "npc.yssa",
			EntityType: "npc",
			Status:     "draft",
			SourceRefs: []string{"note-1.txt"},
			ExplicitFacts: []ExplicitFact{
				{Field: "loyalty", Value: "harbor guild", SourceRef: "note-1.txt"},
			},
		},
		{
			ID:         "npc.yssa",
			EntityType: "npc",
			Status:     "draft",
			SourceRefs: []string{"sketch.png"},
			ExplicitFacts: []ExplicitFact{
				{Field: "loyalty", Value: "smugglers", SourceRef: "sketch.png"},
			},
		},
	}
	authReports := DiffEntityPackages(authorityDispute)
	if len(authReports) != 1 || authReports[0].Resolution != "resolved_by_authority_rule" {
		t.Fatalf("expected resolved_by_authority_rule for text-vs-image non-visual fact, got %+v", authReports)
	}
	if errs := ValidateConflictReport(authReports[0]); len(errs) != 0 {
		t.Fatalf("authority-rule conflict report failed validation: %v", errs)
	}
}

func TestAggregateMakerReport(t *testing.T) {
	ep := validSampleEntityPackage()
	cr := validSampleConflictReport()
	vo := validSampleVisualObservation()

	report, errs := AggregateMakerReport(AggregateInput{
		ReportID:           "report.2026-08-01-harbor-session",
		InputBatch:         []string{"Note 1 (session prep)", "Note 2 (session recap)", "harbor-npc-sketch.png"},
		EntityPackages:     []EntityPackage{ep},
		ConflictReports:    []ConflictReport{cr},
		VisualObservations: []VisualObservation{vo},
		OpenQuestions:      []string{"Confirm Yssa's surname."},
	})
	if len(errs) != 0 {
		t.Fatalf("expected clean aggregation, got errors: %v", errs)
	}
	if report.Status != "draft" {
		t.Fatalf("expected draft status, got %q", report.Status)
	}
	if verrs := ValidateMakerReport(report); len(verrs) != 0 {
		t.Fatalf("aggregated report failed ValidateMakerReport: %v", verrs)
	}

	_, errs = AggregateMakerReport(AggregateInput{})
	if !containsError(errs, "id must not be empty") || !containsError(errs, "input_batch must not be empty") {
		t.Fatalf("expected empty id and input_batch errors, got: %v", errs)
	}

	_, errs = AggregateMakerReport(AggregateInput{
		ReportID:       "r1",
		InputBatch:     []string{"note.txt"},
		EntityPackages: []EntityPackage{ep},
	})
	if !containsError(errs, "not in the emitted conflict reports") {
		t.Fatalf("expected missing conflict reference error, got: %v", errs)
	}

	_, errs = AggregateMakerReport(AggregateInput{
		ReportID:        "r1b",
		InputBatch:      []string{"note.txt"},
		ConflictReports: []ConflictReport{cr},
	})
	if !containsError(errs, "conflict report \"conflict.harbor-quartermaster-name\" is linked to entity") {
		t.Fatalf("expected missing conflict linked entity error, got: %v", errs)
	}

	_, errs = AggregateMakerReport(AggregateInput{
		ReportID:           "r2",
		InputBatch:         []string{"img.png"},
		VisualObservations: []VisualObservation{vo},
	})
	if !containsError(errs, "not in the emitted entity packages") {
		t.Fatalf("expected missing linked entity error, got: %v", errs)
	}

	_, errs = AggregateMakerReport(AggregateInput{
		ReportID:       "r3",
		InputBatch:     []string{"note.txt"},
		EntityPackages: []EntityPackage{ep, ep},
		ConflictReports: []ConflictReport{
			cr, cr,
		},
	})
	if !containsError(errs, "duplicate entity_package id") || !containsError(errs, "duplicate conflict_report id") {
		t.Fatalf("expected duplicate batch ID errors, got: %v", errs)
	}
}

func TestRouteInputAndRouteInputs(t *testing.T) {
	cases := []struct {
		raw      string
		wantSlug string
		wantType InputType
		wantErr  bool
	}{
		{raw: "harbor-npc-sketch.png", wantSlug: "harbor-npc-sketch", wantType: InputTypeImage},
		{raw: "dir/Sub/Map_01.JPG", wantSlug: "map-01", wantType: InputTypeImage},
		{raw: "Note 1 (session prep)", wantSlug: "note-1-session-prep", wantType: InputTypeText},
		{raw: "notes.txt", wantSlug: "notes", wantType: InputTypeText},
		{raw: "---", wantSlug: "input", wantType: InputTypeText},
		{raw: "   ", wantErr: true},
	}

	for _, tc := range cases {
		got, err := RouteInput(tc.raw)
		if tc.wantErr {
			if err == nil {
				t.Fatalf("RouteInput(%q): expected error, got %+v", tc.raw, got)
			}
			continue
		}
		if err != nil {
			t.Fatalf("RouteInput(%q): unexpected error: %v", tc.raw, err)
		}
		if got.Slug != tc.wantSlug || got.Type != tc.wantType || got.Raw != tc.raw {
			t.Fatalf("RouteInput(%q) = %+v, want slug=%q type=%q", tc.raw, got, tc.wantSlug, tc.wantType)
		}
	}

	refs, errs := RouteInputs([]string{"harbor-npc-sketch.png", " ", "session.md"})
	if len(refs) != 2 || len(errs) != 1 {
		t.Fatalf("RouteInputs mixed: got %d refs and %d errs, want 2 and 1", len(refs), len(errs))
	}
	if refs[0].Type != InputTypeImage || refs[1].Type != InputTypeText {
		t.Fatalf("unexpected routed types: %+v", refs)
	}
}

func containsError(errs []error, substr string) bool {
	for _, err := range errs {
		if strings.Contains(err.Error(), substr) {
			return true
		}
	}
	return false
}
