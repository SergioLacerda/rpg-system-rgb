package fixtures

import (
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
)

var semanticIDToken = regexp.MustCompile(`core\.[a-zA-Z0-9_.-]+`)

// extractCitedSemanticIDs scans text for semantic-ID-shaped tokens
// (core.xxx.yyy) and returns the ones that look like concrete citations,
// excluding wildcard-prefix mentions such as "core.term.*" or
// "core.translation.pt-br.*" (identifiable because the character
// immediately following the match is "*", never part of a real ID).
func extractCitedSemanticIDs(text string) []string {
	var ids []string
	for _, loc := range semanticIDToken.FindAllStringIndex(text, -1) {
		token := strings.TrimRight(text[loc[0]:loc[1]], ".")
		if strings.HasPrefix(text[loc[1]:], "*") {
			continue // wildcard-prefix mention, not a citation
		}
		ids = append(ids, token)
	}
	return ids
}

// mirrors: .analysis/pending/20260802-skill-contract-deterministic-validation-analysis.md
// TestSpecialistCitationsAreRealSemanticIDs scans every Specialist
// procedure and worked example for semantic-ID-shaped tokens and fails if
// any of them does not exist in the semantic index. golden-qa.yaml
// already checks this for the benchmark dataset's structured semantic_ids
// field (specialist_golden_qa_test.go); this test covers the prose in
// procedures/ and examples/ themselves, which had no equivalent check.
func TestSpecialistCitationsAreRealSemanticIDs(t *testing.T) {
	knownIDs := loadSemanticIndexIDs(t)

	// examples/ are concrete worked instances — each must ground at least
	// one claim in a real semantic ID. procedures/ are sometimes purely
	// methodological (e.g. identify-ambiguity.md, locate-authority.md
	// describe a process, not a fixed rule) and are exempt from the
	// "at least one" requirement, but any core.* token they DO mention
	// still must be real.
	dirs := []struct {
		path              string
		requireAtLeastOne bool
	}{
		{"../../skills/specialist/procedures", false},
		{"../../skills/specialist/examples", true},
	}
	checked := 0
	for _, dir := range dirs {
		entries, err := os.ReadDir(dir.path)
		if err != nil {
			t.Fatal(err)
		}
		for _, entry := range entries {
			if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".md") {
				continue
			}
			path := filepath.Join(dir.path, entry.Name())
			body, err := os.ReadFile(path) //nolint:gosec // G304: path from ReadDir over the repo's own skills/specialist tree
			if err != nil {
				t.Fatal(err)
			}
			ids := extractCitedSemanticIDs(string(body))
			if dir.requireAtLeastOne && len(ids) == 0 {
				t.Errorf("%s: cites no semantic IDs — every Specialist worked example must ground at least one claim", path)
				continue
			}
			for _, id := range ids {
				checked++
				if !knownIDs[id] {
					t.Errorf("%s: cites semantic id %q, which does not exist in core-v2.index.json", path, id)
				}
			}
		}
	}
	if checked == 0 {
		t.Fatal("no semantic ID citations were checked — directories empty or misconfigured")
	}
}
