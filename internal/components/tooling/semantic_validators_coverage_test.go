package tooling

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func TestDescriptor(t *testing.T) {
	descriptor := Descriptor()
	if descriptor.ID != "tooling" || descriptor.Name == "" || descriptor.Description == "" {
		t.Fatalf("unexpected descriptor: %+v", descriptor)
	}
}

func TestIndexValidationHelpers(t *testing.T) {
	repoRoot := semanticRepoFixture(t)
	unit := validIndexUnit()
	authorityTypes := map[string]bool{"canonical_semantic": true, "generated_artifact": true, "translation": true}
	statuses := map[string]bool{"canonical": true, "derived": true, "current": true}
	kinds := map[string]bool{"rule": true, "translation": true}
	surfaces := map[string]bool{"markdown_en": true, "markdown_pt_br": true}
	consumers := map[string]bool{"specialist": true}
	seen := map[string]bool{}

	if err := validateIndexTopLevelFields(SemanticIndex{Schema: expectedSemanticIndexSchema, SourceLocale: "en", DefaultLocalizedLocale: "PT-br", Units: []SemanticIndexUnit{unit}}); err != nil {
		t.Fatal(err)
	}
	for _, index := range []SemanticIndex{
		{Schema: "wrong", SourceLocale: "en", DefaultLocalizedLocale: "PT-br", Units: []SemanticIndexUnit{unit}},
		{Schema: expectedSemanticIndexSchema, DefaultLocalizedLocale: "PT-br", Units: []SemanticIndexUnit{unit}},
		{Schema: expectedSemanticIndexSchema, SourceLocale: "en", Units: []SemanticIndexUnit{unit}},
		{Schema: expectedSemanticIndexSchema, SourceLocale: "en", DefaultLocalizedLocale: "PT-br"},
	} {
		if err := validateIndexTopLevelFields(index); err == nil {
			t.Fatalf("expected invalid index top level to fail: %+v", index)
		}
	}
	if _, err := requiredSet("field", nil); err == nil {
		t.Fatal("expected missing required set to fail")
	}
	if _, err := requiredSet("field", []string{""}); err == nil {
		t.Fatal("expected empty required set value to fail")
	}
	if err := validateIndexUnit(repoRoot, unit, authorityTypes, statuses, kinds, surfaces, consumers, seen); err != nil {
		t.Fatal(err)
	}

	mutations := []func(*SemanticIndexUnit){
		func(u *SemanticIndexUnit) { u.ID = "" },
		func(u *SemanticIndexUnit) { u.Kind = "unknown" },
		func(u *SemanticIndexUnit) { u.Locale = "" },
		func(u *SemanticIndexUnit) { u.AuthorityType = "unknown" },
		func(u *SemanticIndexUnit) { u.SourceStatus = "unknown" },
		func(u *SemanticIndexUnit) { u.Title = "" },
		func(u *SemanticIndexUnit) { u.SourcePath = "" },
		func(u *SemanticIndexUnit) { u.SourcePath = "/absolute" },
		func(u *SemanticIndexUnit) { u.SourcePath = ".git" },
		func(u *SemanticIndexUnit) { u.ProjectionPaths = nil },
		func(u *SemanticIndexUnit) { u.ProjectionPaths = map[string]string{"unknown": "docs/adr/decision.md"} },
		func(u *SemanticIndexUnit) { u.Relationships = nil },
		func(u *SemanticIndexUnit) { u.Index.RetrievalSummary = "" },
		func(u *SemanticIndexUnit) { u.Index.Track = nil },
		func(u *SemanticIndexUnit) { u.Index.Track = []string{"unknown"} },
		func(u *SemanticIndexUnit) { u.Index.Tags = nil },
		func(u *SemanticIndexUnit) { u.Provenance = nil },
		func(u *SemanticIndexUnit) {
			u.Provenance = map[string]any{"decision_refs": []any{"docs/adr/decision.md"}}
		},
		func(u *SemanticIndexUnit) {
			u.Provenance = map[string]any{"source_revision": "rev", "decision_refs": "bad"}
		},
		func(u *SemanticIndexUnit) { u.ComponentConsumers = nil },
		func(u *SemanticIndexUnit) { u.ComponentConsumers = []string{"unknown"} },
	}
	for _, mutate := range mutations {
		candidate := validIndexUnit()
		mutate(&candidate)
		if err := validateIndexUnit(repoRoot, candidate, authorityTypes, statuses, kinds, surfaces, consumers, map[string]bool{}); err == nil {
			t.Fatalf("expected mutated unit to fail: %+v", candidate)
		}
	}
	bridge := validIndexUnit()
	bridge.AuthorityType = "canonical_markdown_bridge"
	bridge.ProjectionPaths = map[string]string{"markdown_pt_br": "docs/core/PT-br/topic.md"}
	if err := validateIndexUnitProjections(repoRoot, bridge, surfaces); err == nil {
		t.Fatal("expected bridge without markdown_en projection to fail")
	}
	translation := validIndexUnit()
	translation.Kind = "translation"
	translation.AuthorityType = "translation"
	if err := validateIndexTranslationContract(translation, statuses); err == nil {
		t.Fatal("expected incomplete translation contract to fail")
	}
	translation.SourceUnit = "core.test"
	translation.TranslationStatus = "current"
	translation.ProjectionPaths["markdown_pt_br"] = "docs/core/PT-br/topic.md"
	if err := validateIndexTranslationContract(translation, statuses); err != nil {
		t.Fatal(err)
	}
	if err := validateIndexRelationships(SemanticIndexUnit{ID: "u", Relationships: map[string][]string{"depends_on": {"missing"}}}, map[string]bool{"u": true}); err == nil {
		t.Fatal("expected unknown relationship target to fail")
	}
	if err := validateIndexRelationships(SemanticIndexUnit{ID: "u", SourceUnit: "missing", Relationships: map[string][]string{}}, map[string]bool{"u": true}); err == nil {
		t.Fatal("expected unknown source_unit to fail")
	}
}

func TestL10nValidationHelpers(t *testing.T) {
	repoRoot := semanticRepoFixture(t)
	manifest := l10nManifest{
		Schema:              expectedL10nSchema,
		SourceLocale:        "en",
		LocalizedLocales:    []string{"PT-br"},
		AuthorityTypes:      []string{"translation"},
		TranslationStatuses: []string{"current", "missing", "not_applicable"},
		AuthorityDecision:   "docs/adr/decision.md",
		Description:         "desc",
		Documents:           []l10nDocument{validL10nDocument()},
	}
	if err := validateL10nManifestTopLevel(manifest); err != nil {
		t.Fatal(err)
	}
	for _, mutate := range []func(*l10nManifest){
		func(m *l10nManifest) { m.Schema = "wrong" },
		func(m *l10nManifest) { m.SourceLocale = "pt" },
		func(m *l10nManifest) { m.LocalizedLocales = nil },
		func(m *l10nManifest) { m.Description = "" },
	} {
		candidate := manifest
		mutate(&candidate)
		if err := validateL10nManifestTopLevel(candidate); err == nil {
			t.Fatalf("expected invalid l10n manifest to fail: %+v", candidate)
		}
	}
	if _, _, err := buildL10nLookupSets(l10nManifest{}); err == nil {
		t.Fatal("expected missing l10n lookup sets to fail")
	}
	if err := validateL10nDocuments(repoRoot, manifest.Documents, map[string]bool{"translation": true}, map[string]bool{"current": true}); err != nil {
		t.Fatal(err)
	}
	if err := validateL10nDocuments(repoRoot, nil, map[string]bool{}, map[string]bool{}); err == nil {
		t.Fatal("expected empty l10n documents to fail")
	}
	for _, mutate := range []func(*l10nDocument){
		func(d *l10nDocument) { d.Source = "" },
		func(d *l10nDocument) { d.Source = "docs/core/PT-br/topic.md" },
		func(d *l10nDocument) { d.Locale = "en" },
		func(d *l10nDocument) { d.AuthorityType = "unknown" },
		func(d *l10nDocument) { d.TranslationStatus = "unknown" },
		func(d *l10nDocument) { d.SourceRevision = "" },
		func(d *l10nDocument) { d.Localized = "" },
		func(d *l10nDocument) { d.Localized = "docs/core/en/topic.md" },
		func(d *l10nDocument) { d.Localized = "docs/core/PT-br/missing.md" },
	} {
		doc := validL10nDocument()
		mutate(&doc)
		if err := validateL10nDocument(repoRoot, doc, map[string]bool{"translation": true}, map[string]bool{"current": true}, map[string]bool{}); err == nil {
			t.Fatalf("expected invalid l10n document to fail: %+v", doc)
		}
	}
	missingDoc := validL10nDocument()
	missingDoc.TranslationStatus = "missing"
	missingDoc.Localized = ""
	if err := validateL10nDocumentLocalized(repoRoot, missingDoc); err != nil {
		t.Fatal(err)
	}
	if containsString([]string{"a"}, "b") {
		t.Fatal("containsString should reject absent value")
	}
	if hasSlashPrefix("short", "long-prefix") {
		t.Fatal("hasSlashPrefix should reject too-short value")
	}
}

func TestValidateDocsL10nManifestPublicFlow(t *testing.T) {
	root := semanticRepoFixture(t)
	manifestFile := filepath.Join(root, "docs", "core", "semantic", "l10n-manifest.v0.1.json")
	manifest := l10nManifest{
		Schema:              expectedL10nSchema,
		SourceLocale:        "en",
		LocalizedLocales:    []string{"PT-br"},
		AuthorityTypes:      []string{"translation"},
		TranslationStatuses: []string{"current", "missing", "not_applicable"},
		AuthorityDecision:   "docs/adr/decision.md",
		Description:         "desc",
		Documents:           []l10nDocument{validL10nDocument()},
	}
	writeJSONForValidation(t, manifestFile, manifest)
	if err := ValidateDocsL10nManifest(manifestFile); err != nil {
		t.Fatal(err)
	}

	writeValidationFile(t, filepath.Join(root, "docs", "core", "en", "missing.md"), "content")
	if err := ValidateDocsL10nManifest(manifestFile); err == nil {
		t.Fatal("expected missing source coverage to fail")
	}
	if err := os.Remove(filepath.Join(root, "docs", "core", "en", "missing.md")); err != nil {
		t.Fatal(err)
	}

	manifest.AuthorityDecision = "docs/adr/missing.md"
	writeJSONForValidation(t, manifestFile, manifest)
	if err := ValidateDocsL10nManifest(manifestFile); err == nil {
		t.Fatal("expected missing authority decision to fail")
	}
	writeValidationFile(t, manifestFile, "{")
	if err := ValidateDocsL10nManifest(manifestFile); err == nil {
		t.Fatal("expected invalid l10n manifest JSON to fail")
	}
}

func TestContractProjectionAndGeneratedHelpers(t *testing.T) {
	repoRoot := semanticRepoFixture(t)
	indexFile := filepath.Join(repoRoot, "docs", "core", "semantic", "index.json")
	contract := validConsumerContract()
	contracts := consumerContractsFile{Schema: expectedContractsSchema, SourceIndex: repoRelative(repoRoot, indexFile), Description: "desc", Contracts: []consumerContract{contract}}
	if err := validateContractsTopLevel(contracts, repoRoot, indexFile); err != nil {
		t.Fatal(err)
	}
	if err := validateContract(repoRoot, contract, map[string]bool{}, map[string]bool{"canonical_semantic": true}, map[string]bool{"canonical": true}, map[string]bool{"markdown_en": true}, map[string]bool{"specialist": true}, map[string]bool{"id": true}); err != nil {
		t.Fatal(err)
	}
	for _, mutate := range []func(*consumerContract){
		func(c *consumerContract) { c.ID = "" },
		func(c *consumerContract) { c.Component = "unknown" },
		func(c *consumerContract) { c.Status = "unknown" },
		func(c *consumerContract) { c.Description = "" },
		func(c *consumerContract) { c.AllowedInputs = []string{"unknown"} },
		func(c *consumerContract) { c.RequiredUnitFields = []string{""} },
		func(c *consumerContract) { c.AllowedProjectionSurfaces = []string{"unknown"} },
		func(c *consumerContract) { c.ForbiddenAuthorityTypes = []string{"unknown"} },
		func(c *consumerContract) { c.RequiredDisclosures = nil },
		func(c *consumerContract) { c.Outputs = nil },
		func(c *consumerContract) { c.ValidationChecks = nil },
		func(c *consumerContract) { c.SourceRefs = []string{"/absolute"} },
	} {
		candidate := validConsumerContract()
		mutate(&candidate)
		if err := validateContract(repoRoot, candidate, map[string]bool{}, map[string]bool{"canonical_semantic": true}, map[string]bool{"canonical": true}, map[string]bool{"markdown_en": true}, map[string]bool{"specialist": true}, map[string]bool{"id": true}); err == nil {
			t.Fatalf("expected invalid contract to fail: %+v", candidate)
		}
	}

	proj := validProjection()
	statuses := map[string]bool{"derived": true}
	surfaces := map[string]bool{"markdown_en": true}
	components := map[string]bool{"specialist": true}
	units := map[string]bool{"core.test": true}
	allowed := map[string]map[string]bool{"specialist": {"markdown_en": true}}
	if err := validateProjectionEntry(repoRoot, proj, map[string]bool{}, statuses, surfaces, components, units, allowed); err != nil {
		t.Fatal(err)
	}
	for _, mutate := range []func(*projection){
		func(p *projection) { p.ID = "" },
		func(p *projection) { p.Surface = "unknown" },
		func(p *projection) { p.Owner = "unknown" },
		func(p *projection) { p.Status = "unknown" },
		func(p *projection) { p.Description = "" },
		func(p *projection) { p.SourceUnits = nil },
		func(p *projection) { p.SourceUnits = []string{""} },
		func(p *projection) { p.SourceUnits = []string{"missing"} },
		func(p *projection) { p.OutputPath = "" },
		func(p *projection) { p.OutputPath = "/absolute" },
		func(p *projection) { p.Provenance = nil },
		func(p *projection) { p.Provenance = map[string]any{"source_revision": ""} },
		func(p *projection) {
			p.Provenance = map[string]any{"source_revision": "rev", "decision_refs": []any{123}}
		},
		func(p *projection) { p.RequiredDisclosures = nil },
		func(p *projection) { p.GenerationGate = []string{""} },
	} {
		candidate := validProjection()
		mutate(&candidate)
		if err := validateProjectionEntry(repoRoot, candidate, map[string]bool{}, statuses, surfaces, components, units, allowed); err == nil {
			t.Fatalf("expected invalid projection to fail: %+v", candidate)
		}
	}

	generatedProj := generatedValidationManifestProjection{ID: "p", Surface: "markdown_en", Owner: "specialist", SourceUnits: []string{"core.test"}, OutputPath: "generated/out.json"}
	generatedOut := generatedValidationOutput{Schema: "rgb-docs-generated-projection/0.1", ProjectionID: "p", Surface: "markdown_en", Owner: "specialist", AuthorityType: "generated_artifact", GeneratedFrom: []string{"core.test"}, Units: []struct {
		ID string `json:"id"`
	}{{ID: "core.test"}}}
	if err := validateGeneratedProjectionOutput(generatedProj, generatedOut); err != nil {
		t.Fatal(err)
	}
	for _, mutate := range []func(*generatedValidationOutput){
		func(o *generatedValidationOutput) { o.Schema = "wrong" },
		func(o *generatedValidationOutput) { o.ProjectionID = "wrong" },
		func(o *generatedValidationOutput) { o.Surface = "wrong" },
		func(o *generatedValidationOutput) { o.Owner = "wrong" },
		func(o *generatedValidationOutput) { o.AuthorityType = "canonical" },
		func(o *generatedValidationOutput) { o.GeneratedFrom = nil },
		func(o *generatedValidationOutput) { o.Units = nil },
	} {
		out := generatedOut
		mutate(&out)
		if err := validateGeneratedProjectionOutput(generatedProj, out); err == nil {
			t.Fatalf("expected invalid generated output to fail: %+v", out)
		}
	}
}

func semanticRepoFixture(t *testing.T) string {
	t.Helper()
	root := t.TempDir()
	for _, path := range []string{
		"README.md",
		"docs/adr/decision.md",
		"docs/core/en/topic.md",
		"docs/core/PT-br/topic.md",
		"generated/out.json",
	} {
		content := "content"
		if path == "docs/adr/decision.md" {
			content = "## Status\n\nAccepted.\n"
		}
		writeValidationFile(t, filepath.Join(root, path), content)
	}
	return root
}

func validIndexUnit() SemanticIndexUnit {
	return SemanticIndexUnit{
		ID:              "core.test",
		Kind:            "rule",
		Locale:          "en",
		AuthorityType:   "canonical_semantic",
		SourceStatus:    "canonical",
		Title:           "Test",
		SourcePath:      "docs/core/en/topic.md",
		ProjectionPaths: map[string]string{"markdown_en": "docs/core/en/topic.md"},
		Relationships:   map[string][]string{},
		Index: SemanticIndexUnitIndex{
			Track:            []string{"specialist"},
			Tags:             []string{"test"},
			RetrievalSummary: "summary",
		},
		Provenance:         map[string]any{"source_revision": "rev", "decision_refs": []any{"docs/adr/decision.md"}},
		ComponentConsumers: []string{"specialist"},
	}
}

func validL10nDocument() l10nDocument {
	return l10nDocument{
		Source:            "docs/core/en/topic.md",
		Locale:            "PT-br",
		Localized:         "docs/core/PT-br/topic.md",
		AuthorityType:     "translation",
		TranslationStatus: "current",
		SourceRevision:    "rev",
		Notes:             []string{"note"},
	}
}

func validConsumerContract() consumerContract {
	return consumerContract{
		ID:                        "contract",
		Component:                 "specialist",
		Status:                    "canonical",
		Description:               "desc",
		AllowedInputs:             []string{"canonical_semantic"},
		RequiredUnitFields:        []string{"id"},
		AllowedProjectionSurfaces: []string{"markdown_en"},
		ForbiddenAuthorityTypes:   []string{"canonical_semantic"},
		RequiredDisclosures:       []string{"disclosure"},
		Outputs:                   []string{"output"},
		ValidationChecks:          []string{"check"},
		SourceRefs:                []string{"docs/adr/decision.md"},
	}
}

func validProjection() projection {
	return projection{
		ID:                  "projection",
		Surface:             "markdown_en",
		Owner:               "specialist",
		Status:              "derived",
		Description:         "desc",
		SourceUnits:         []string{"core.test"},
		OutputPath:          "generated/out.json",
		Provenance:          map[string]any{"source_revision": "rev", "decision_refs": []any{"docs/adr/decision.md"}},
		RequiredDisclosures: []string{"disclosure"},
		GenerationGate:      []string{"gate"},
	}
}

func writeJSONForValidation(t *testing.T, path string, value any) {
	t.Helper()
	bytes, err := json.Marshal(value)
	if err != nil {
		t.Fatal(err)
	}
	writeValidationFile(t, path, string(bytes))
}

func writeValidationFile(t *testing.T, path, content string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0o750); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(content), 0o640); err != nil {
		t.Fatal(err)
	}
}

func semanticRepoFiles(t *testing.T) (root, semantic, indexFile, contractsFile, manifestFile, sourceFile string) {
	t.Helper()
	root = semanticRepoFixture(t)
	semantic = filepath.Join(root, "docs", "core", "semantic")
	indexFile = filepath.Join(semantic, "core-v2.index.json")
	contractsFile = filepath.Join(semantic, "consumer-contracts.v0.1.json")
	manifestFile = filepath.Join(semantic, "projection-manifest.v0.1.json")
	sourceFile = filepath.Join(semantic, "source", "core-v2-rules.v0.1.json")
	return root, semantic, indexFile, contractsFile, manifestFile, sourceFile
}

func writeSemanticValidationInputs(t *testing.T) (root, indexFile, contractsFile, manifestFile, sourceFile string) {
	t.Helper()
	root, _, indexFile, contractsFile, manifestFile, sourceFile = semanticRepoFiles(t)
	sourceIndex := SemanticIndex{
		Schema:                 expectedSemanticIndexSchema,
		SourceLocale:           "en",
		DefaultLocalizedLocale: "PT-br",
		AuthorityTypes:         []string{"canonical_semantic", "generated_artifact", "translation"},
		SourceStatuses:         []string{"canonical", "derived", "current"},
		Kinds:                  []string{"rule", "translation"},
		ProjectionSurfaces:     []string{"markdown_en", "markdown_pt_br"},
		ComponentConsumers:     []string{"specialist"},
		Units:                  []SemanticIndexUnit{validIndexUnit()},
	}
	contracts := consumerContractsFile{
		Schema:      expectedContractsSchema,
		SourceIndex: repoRelative(root, indexFile),
		Description: "desc",
		Contracts:   []consumerContract{validConsumerContract()},
	}
	manifest := projectionManifest{
		Schema:            expectedProjectionSchema,
		SourceIndex:       repoRelative(root, indexFile),
		ConsumerContracts: repoRelative(root, contractsFile),
		Description:       "desc",
		Projections:       []projection{validProjection()},
	}
	source := generationSource{Units: []generationSourceUnit{{
		ID:              "core.test",
		Kind:            "rule",
		Locale:          "en",
		AuthorityType:   "canonical_semantic",
		SourceStatus:    "canonical",
		Title:           "Test",
		Statement:       "Resolved statement.",
		ProjectionPaths: map[string]string{"markdown_en": "docs/core/en/topic.md"},
		Relationships:   map[string][]string{},
		Provenance:      map[string]any{"source_revision": "rev", "decision_refs": []any{"docs/adr/decision.md"}},
	}}}

	writeJSONForValidation(t, indexFile, sourceIndex)
	writeJSONForValidation(t, contractsFile, contracts)
	writeJSONForValidation(t, manifestFile, manifest)
	writeJSONForValidation(t, sourceFile, source)
	return root, indexFile, contractsFile, manifestFile, sourceFile
}

func TestPublicSemanticValidatorsHappyPaths(t *testing.T) {
	root, indexFile, contractsFile, manifestFile, sourceFile := writeSemanticValidationInputs(t)

	if err := ValidateSemanticIndex(indexFile); err != nil {
		t.Fatal(err)
	}
	if err := ValidateConsumerContracts(contractsFile, indexFile); err != nil {
		t.Fatal(err)
	}
	if err := ValidateProjectionManifest(manifestFile, indexFile, contractsFile); err != nil {
		t.Fatal(err)
	}
	if err := Generate(root, manifestFile, indexFile, sourceFile); err != nil {
		t.Fatal(err)
	}
	if err := ValidateGeneratedProjections(manifestFile); err != nil {
		t.Fatal(err)
	}
}

func TestPublicSemanticValidatorsTopLevelErrors(t *testing.T) {
	_, indexFile, _, _, _ := writeSemanticValidationInputs(t)

	index := SemanticIndex{
		Schema:                 expectedSemanticIndexSchema,
		SourceLocale:           "en",
		DefaultLocalizedLocale: "PT-br",
		AuthorityTypes:         []string{"canonical_semantic"},
		SourceStatuses:         []string{"canonical"},
		Kinds:                  []string{"rule"},
		ProjectionSurfaces:     []string{"markdown_en"},
		ComponentConsumers:     []string{"specialist"},
		Units:                  []SemanticIndexUnit{validIndexUnit()},
	}
	index.AuthorityTypes = nil
	writeJSONForValidation(t, indexFile, index)
	if err := ValidateSemanticIndex(indexFile); err == nil {
		t.Fatal("expected missing authority types to fail")
	}

	root, indexFile, contractsFile, _, _ := writeSemanticValidationInputs(t)
	contracts := consumerContractsFile{Schema: "wrong", SourceIndex: repoRelative(root, indexFile), Description: "desc", Contracts: []consumerContract{validConsumerContract()}}
	writeJSONForValidation(t, contractsFile, contracts)
	if err := ValidateConsumerContracts(contractsFile, indexFile); err == nil {
		t.Fatal("expected bad contracts schema to fail")
	}

	root, indexFile, contractsFile, manifestFile, _ := writeSemanticValidationInputs(t)
	manifest := projectionManifest{Schema: "wrong", SourceIndex: repoRelative(root, indexFile), ConsumerContracts: repoRelative(root, contractsFile), Description: "desc", Projections: []projection{validProjection()}}
	writeJSONForValidation(t, manifestFile, manifest)
	if err := ValidateProjectionManifest(manifestFile, indexFile, contractsFile); err == nil {
		t.Fatal("expected bad projection schema to fail")
	}

	_, _, _, manifestFile, _ = writeSemanticValidationInputs(t)
	writeJSONForValidation(t, manifestFile, generatedValidationManifest{})
	if err := ValidateGeneratedProjections(manifestFile); err == nil {
		t.Fatal("expected empty generated manifest to fail")
	}
}

func TestGenerateErrors(t *testing.T) {
	root, indexFile, _, manifestFile, sourceFile := writeSemanticValidationInputs(t)
	if err := Generate(root, filepath.Join(root, "missing.json"), indexFile, sourceFile); err == nil {
		t.Fatal("expected missing manifest to fail")
	}
	writeValidationFile(t, manifestFile, "{")
	if err := Generate(root, manifestFile, indexFile, sourceFile); err == nil {
		t.Fatal("expected invalid manifest JSON to fail")
	}
	_, _, _, manifestFile, sourceFile = writeSemanticValidationInputs(t)
	writeValidationFile(t, indexFile, "{")
	if err := Generate(root, manifestFile, indexFile, sourceFile); err == nil {
		t.Fatal("expected invalid index JSON to fail")
	}
	root, indexFile, _, manifestFile, sourceFile = writeSemanticValidationInputs(t)
	writeValidationFile(t, sourceFile, "{")
	if err := Generate(root, manifestFile, indexFile, sourceFile); err == nil {
		t.Fatal("expected invalid source JSON to fail")
	}
	root, indexFile, _, manifestFile, sourceFile = writeSemanticValidationInputs(t)
	manifest := generationManifest{Projections: []generationProjection{validGenerationProjection()}}
	manifest.Projections[0].SourceUnits = []string{"missing"}
	writeJSONForValidation(t, manifestFile, manifest)
	if err := Generate(root, manifestFile, indexFile, sourceFile); err == nil {
		t.Fatal("expected unknown source unit to fail")
	}
	root, indexFile, _, manifestFile, sourceFile = writeSemanticValidationInputs(t)
	manifest.Projections[0].SourceUnits = []string{"core.test"}
	manifest.Projections[0].OutputPath = "generated/blocked/out.json"
	writeJSONForValidation(t, manifestFile, manifest)
	writeValidationFile(t, filepath.Join(root, "generated", "blocked"), "x")
	if err := Generate(root, manifestFile, indexFile, sourceFile); err == nil {
		t.Fatal("expected blocked output path to fail")
	}
}

func validGenerationProjection() generationProjection {
	return generationProjection(validProjection())
}

func TestReadGeneratedJSONErrors(t *testing.T) {
	dir := t.TempDir()
	if err := readGeneratedJSON(filepath.Join(dir, "missing.json"), &generatedValidationManifest{}); err == nil {
		t.Fatal("expected missing JSON to fail")
	}
	path := filepath.Join(dir, "bad.json")
	writeValidationFile(t, path, "{")
	if err := readGeneratedJSON(path, &generatedValidationManifest{}); err == nil {
		t.Fatal("expected invalid JSON to fail")
	}
}

func TestLoadInputJSONErrors(t *testing.T) {
	dir := t.TempDir()
	contracts := filepath.Join(dir, "contracts.json")
	index := filepath.Join(dir, "index.json")
	manifest := filepath.Join(dir, "manifest.json")
	writeValidationFile(t, contracts, "{")
	writeValidationFile(t, index, "{}")
	if _, _, err := loadConsumerContractsAndIndex(contracts, index); err == nil {
		t.Fatal("expected invalid contracts JSON")
	}
	writeValidationFile(t, contracts, "{}")
	writeValidationFile(t, index, "{")
	if _, _, err := loadConsumerContractsAndIndex(contracts, index); err == nil {
		t.Fatal("expected invalid index JSON")
	}
	writeValidationFile(t, manifest, "{")
	if _, _, _, err := loadProjectionInputs(manifest, index, contracts); err == nil {
		t.Fatal("expected invalid projection manifest JSON")
	}
	writeJSONForValidation(t, manifest, projectionManifest{})
	if _, _, _, err := loadProjectionInputs(manifest, index, contracts); err == nil {
		t.Fatal("expected invalid projection index JSON")
	}
	writeValidationFile(t, index, "{}")
	writeValidationFile(t, contracts, "{")
	if _, _, _, err := loadProjectionInputs(manifest, index, contracts); err == nil {
		t.Fatal("expected invalid projection contracts JSON")
	}
}
