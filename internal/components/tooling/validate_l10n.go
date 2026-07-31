package tooling

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
)

const expectedL10nSchema = "rgb-docs-l10n/0.1"

type l10nManifest struct {
	Schema              string         `json:"schema"`
	SourceLocale        string         `json:"source_locale"`
	LocalizedLocales    []string       `json:"localized_locales"`
	AuthorityTypes      []string       `json:"authority_types"`
	TranslationStatuses []string       `json:"translation_statuses"`
	AuthorityDecision   string         `json:"authority_decision"`
	Description         string         `json:"description"`
	Documents           []l10nDocument `json:"documents"`
}

type l10nDocument struct {
	Source            string   `json:"source"`
	Locale            string   `json:"locale"`
	Localized         string   `json:"localized"`
	AuthorityType     string   `json:"authority_type"`
	TranslationStatus string   `json:"translation_status"`
	SourceRevision    string   `json:"source_revision"`
	Notes             []string `json:"notes"`
}

// ValidateDocsL10nManifest validates docs/core/semantic/l10n-manifest.v0.1.json,
// migrated from scripts/validate_docs_l10n_manifest.go.
func ValidateDocsL10nManifest(manifestFile string) error {
	manifest, err := loadL10nManifest(manifestFile)
	if err != nil {
		return err
	}
	if err := validateL10nManifestTopLevel(manifest); err != nil {
		return err
	}

	repoRoot, err := resolveL10nRepoRoot(manifestFile, manifest.AuthorityDecision)
	if err != nil {
		return err
	}

	authorityTypes, statuses, err := buildL10nLookupSets(manifest)
	if err != nil {
		return err
	}

	if err := validateL10nDocuments(repoRoot, manifest.Documents, authorityTypes, statuses); err != nil {
		return err
	}

	fmt.Printf("docs-l10n validation passed: %s (%d documents)\n", manifestFile, len(manifest.Documents))
	return nil
}

// loadL10nManifest reads and parses the l10n manifest file.
func loadL10nManifest(manifestFile string) (l10nManifest, error) {
	bytes, err := os.ReadFile(manifestFile) //nolint:gosec // G304: path is a caller-supplied validation target, by design
	if err != nil {
		return l10nManifest{}, err
	}
	var manifest l10nManifest
	if err := json.Unmarshal(bytes, &manifest); err != nil {
		return l10nManifest{}, fmt.Errorf("invalid JSON: %w", err)
	}
	return manifest, nil
}

// validateL10nManifestTopLevel validates the manifest's own schema/locale/
// description fields, independent of any individual document entry.
func validateL10nManifestTopLevel(manifest l10nManifest) error {
	if manifest.Schema != expectedL10nSchema {
		return fmt.Errorf("schema must be `%s`", expectedL10nSchema)
	}
	if manifest.SourceLocale != "en" {
		return errors.New("source_locale must be `en`")
	}
	if !containsString(manifest.LocalizedLocales, "PT-br") {
		return errors.New("localized_locales must include `PT-br`")
	}
	if manifest.Description == "" {
		return errors.New("description must be non-empty")
	}
	return nil
}

// resolveL10nRepoRoot resolves the repository root from the manifest file
// path and validates the manifest's authority decision reference.
func resolveL10nRepoRoot(manifestFile, authorityDecision string) (string, error) {
	repoRoot, err := repoRootFromFile(manifestFile)
	if err != nil {
		return "", err
	}
	if err := validateAcceptedDecision(repoRoot, authorityDecision); err != nil {
		return "", err
	}
	return repoRoot, nil
}

// buildL10nLookupSets builds the authority-type and translation-status
// allowed-value sets used to validate every document.
func buildL10nLookupSets(manifest l10nManifest) (authorityTypes, statuses map[string]bool, err error) {
	authorityTypes = stringSet(manifest.AuthorityTypes)
	if len(authorityTypes) == 0 {
		return nil, nil, errors.New("authority_types must be non-empty")
	}
	statuses = stringSet(manifest.TranslationStatuses)
	if len(statuses) == 0 {
		return nil, nil, errors.New("translation_statuses must be non-empty")
	}
	return authorityTypes, statuses, nil
}

// validateL10nDocuments validates every document entry, then checks that
// every source Markdown file under docs/core/en has a manifest entry.
func validateL10nDocuments(repoRoot string, documents []l10nDocument, authorityTypes, statuses map[string]bool) error {
	if len(documents) == 0 {
		return errors.New("documents must be non-empty")
	}
	sources := map[string]bool{}
	for _, document := range documents {
		if err := validateL10nDocument(repoRoot, document, authorityTypes, statuses, sources); err != nil {
			return err
		}
		sources[document.Source] = true
	}
	return validateL10nSourceCoverage(repoRoot, sources)
}

func validateL10nDocument(repoRoot string, document l10nDocument, authorityTypes, statuses map[string]bool, sources map[string]bool) error {
	if err := validateL10nDocumentSource(repoRoot, document, authorityTypes, statuses, sources); err != nil {
		return err
	}
	return validateL10nDocumentLocalized(repoRoot, document)
}

func validateL10nDocumentSource(repoRoot string, document l10nDocument, authorityTypes, statuses map[string]bool, sources map[string]bool) error {
	if err := validateL10nSourcePath(repoRoot, document, sources); err != nil {
		return err
	}
	return validateL10nSourceMetadata(document, authorityTypes, statuses)
}

// validateL10nSourcePath validates the document's source field: presence,
// uniqueness, Markdown-under-docs/core/en shape, and on-disk existence.
func validateL10nSourcePath(repoRoot string, document l10nDocument, sources map[string]bool) error {
	if document.Source == "" {
		return errors.New("documents[].source must be non-empty")
	}
	if sources[document.Source] {
		return fmt.Errorf("duplicate source entry: %s", document.Source)
	}
	if !hasSlashPrefix(document.Source, "docs/core/en/") || filepath.Ext(document.Source) != ".md" {
		return fmt.Errorf("%s: source must be a Markdown path under docs/core/en", document.Source)
	}
	return validateExistingRepoPath(repoRoot, document.Source, document.Source, "source")
}

// validateL10nSourceMetadata validates the document's locale, authority
// type, translation status, and source revision fields.
func validateL10nSourceMetadata(document l10nDocument, authorityTypes, statuses map[string]bool) error {
	if document.Locale != "PT-br" {
		return fmt.Errorf("%s: locale must be `PT-br`", document.Source)
	}
	if !authorityTypes[document.AuthorityType] {
		return fmt.Errorf("%s: unknown authority_type `%s`", document.Source, document.AuthorityType)
	}
	if !statuses[document.TranslationStatus] {
		return fmt.Errorf("%s: unknown translation_status `%s`", document.Source, document.TranslationStatus)
	}
	if document.SourceRevision == "" {
		return fmt.Errorf("%s: source_revision must be non-empty", document.Source)
	}
	return nil
}

func validateL10nDocumentLocalized(repoRoot string, document l10nDocument) error {
	switch document.TranslationStatus {
	case "missing", "not_applicable":
		if document.Localized == "" {
			return nil
		}
	default:
		if document.Localized == "" {
			return fmt.Errorf("%s: localized must be non-empty for status `%s`", document.Source, document.TranslationStatus)
		}
		if !hasSlashPrefix(document.Localized, "docs/core/PT-br/") || filepath.Ext(document.Localized) != ".md" {
			return fmt.Errorf("%s: localized must be a Markdown path under docs/core/PT-br", document.Source)
		}
		if err := validateExistingRepoPath(repoRoot, document.Localized, document.Source, "localized"); err != nil {
			return err
		}
	}
	return nil
}

func validateL10nSourceCoverage(repoRoot string, sources map[string]bool) error {
	return filepath.WalkDir(filepath.Join(repoRoot, "docs/core/en"), func(path string, entry fs.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		if entry.IsDir() || filepath.Ext(path) != ".md" {
			return nil
		}
		rel, err := filepath.Rel(repoRoot, path)
		if err != nil {
			return err
		}
		source := filepath.ToSlash(rel)
		if !sources[source] {
			return fmt.Errorf("missing manifest entry for source: %s", source)
		}
		return nil
	})
}

func validateExistingRepoPath(repoRoot, pathValue, source, field string) error {
	if filepath.IsAbs(pathValue) {
		return fmt.Errorf("%s: %s must be repository-relative", source, field)
	}
	fullPath := filepath.Join(repoRoot, filepath.FromSlash(pathValue))
	if _, err := os.Stat(fullPath); err != nil {
		return fmt.Errorf("%s: %s does not exist: %s", source, field, pathValue)
	}
	return nil
}

func containsString(values []string, target string) bool {
	for _, value := range values {
		if value == target {
			return true
		}
	}
	return false
}

func hasSlashPrefix(pathValue, prefix string) bool {
	return len(pathValue) >= len(prefix) && pathValue[:len(prefix)] == prefix
}
