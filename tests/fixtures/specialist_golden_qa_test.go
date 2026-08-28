package fixtures

import (
	"encoding/json"
	"os"
	"regexp"
	"strings"
	"testing"
)

type goldenQAEntry struct {
	id              string
	category        string
	semanticIDs     []string
	answerTraits    []string
	ambiguityNote   string
	hasSemanticIDs  bool // distinguishes "semantic_ids: []" (present, empty) from missing
	hasAnswerTraits bool // distinguishes "expected_answer_traits: []" (present, empty) from missing
}

var (
	goldenQAField = regexp.MustCompile(`(?m)^    (\w+):\s*(.*)$`)
)

// parseBracketList parses a flow-style YAML list such as `[a, b]` or
// `["a", "b"]` into its trimmed, unquoted elements. Returns nil for an
// empty or absent list.
func parseBracketList(raw string) []string {
	inner := strings.TrimSuffix(strings.TrimPrefix(raw, "["), "]")
	inner = strings.TrimSpace(inner)
	if inner == "" {
		return nil
	}
	var items []string
	for _, part := range strings.Split(inner, ",") {
		item := strings.TrimSpace(part)
		item = strings.Trim(item, `"`)
		items = append(items, item)
	}
	return items
}

// parseGoldenQA reads benchmark/golden-qa.yaml from the skills/specialist
// package and extracts each entry. Line-oriented reader, not a full YAML
// parser, matching this repo's other hand-authored fixture files (see
// parity_test.go).
func parseGoldenQA(t *testing.T) []goldenQAEntry {
	t.Helper()
	body, err := os.ReadFile("../../skills/specialist/benchmark/golden-qa.yaml")
	if err != nil {
		t.Fatal(err)
	}

	blocks := strings.Split(string(body), "\n  - id:")
	var entries []goldenQAEntry
	for i, block := range blocks {
		if i == 0 {
			continue // header/comment content before the first entry
		}
		full := "  id:" + block
		fields := map[string]string{}
		for _, match := range goldenQAField.FindAllStringSubmatch(full, -1) {
			fields[match[1]] = strings.TrimSpace(match[2])
		}
		id := strings.TrimSpace(strings.SplitN(full, "\n", 2)[0])
		id = strings.TrimPrefix(id, "id:")
		id = strings.TrimSpace(id)

		rawIDs, hasIDs := fields["semantic_ids"]
		rawTraits, hasTraits := fields["expected_answer_traits"]

		entries = append(entries, goldenQAEntry{
			id:              id,
			category:        fields["category"],
			semanticIDs:     parseBracketList(rawIDs),
			answerTraits:    parseBracketList(rawTraits),
			ambiguityNote:   fields["ambiguity_note"],
			hasSemanticIDs:  hasIDs,
			hasAnswerTraits: hasTraits,
		})
	}
	return entries
}

func loadSemanticIndexIDs(t *testing.T) map[string]bool {
	t.Helper()
	body, err := os.ReadFile("../../docs/core/semantic/core-v2.index.json")
	if err != nil {
		t.Fatal(err)
	}
	var index struct {
		Units []struct {
			ID string `json:"id"`
		} `json:"units"`
	}
	if err := json.Unmarshal(body, &index); err != nil {
		t.Fatal(err)
	}
	ids := make(map[string]bool, len(index.Units))
	for _, unit := range index.Units {
		ids[unit.ID] = true
	}
	return ids
}

// mirrors: skills/specialist/benchmark/golden-qa.yaml
// TestGoldenQANormativeEntriesCiteRealSemanticIDs is the dataset-level
// enforcement of design.md's must_not: answer_normatively_without_source.
// Every entry outside category "ambiguous" must cite at least one
// semantic ID, and every cited ID must actually exist in the semantic
// index — a golden entry that cites a stale or invented ID is exactly the
// defect this benchmark exists to catch before any runtime Specialist
// answer does the same thing.
func TestGoldenQANormativeEntriesCiteRealSemanticIDs(t *testing.T) {
	entries := parseGoldenQA(t)
	if len(entries) == 0 {
		t.Fatal("no entries parsed from golden-qa.yaml")
	}
	knownIDs := loadSemanticIndexIDs(t)

	seenCategories := map[string]bool{}
	for _, entry := range entries {
		t.Run(entry.id, func(t *testing.T) {
			if entry.category == "" {
				t.Fatalf("%s: missing category", entry.id)
			}
			seenCategories[entry.category] = true

			if entry.category == "ambiguous" {
				if strings.TrimSpace(entry.ambiguityNote) == "" {
					t.Fatalf("%s: category=ambiguous requires a non-empty ambiguity_note", entry.id)
				}
				return
			}

			if !entry.hasSemanticIDs || len(entry.semanticIDs) == 0 {
				t.Fatalf("%s: category=%s is normative and must cite at least one semantic_ids entry", entry.id, entry.category)
			}
			for _, id := range entry.semanticIDs {
				if !knownIDs[id] {
					t.Fatalf("%s: cites semantic id %q, which does not exist in core-v2.index.json", entry.id, id)
				}
			}
		})
	}

	requiredCategories := []string{
		"rule_lookup", "action_classification", "ambiguous",
		"example_validation", "en_pt_br_parity",
	}
	for _, category := range requiredCategories {
		if !seenCategories[category] {
			t.Errorf("golden-qa.yaml has no entry for required category %q", category)
		}
	}
}

// mirrors: skills/specialist/benchmark/golden-qa.yaml
// TestGoldenQANormativeEntriesHaveExpectedAnswerTraits is the golden-layer
// placeholder oracle for design.md's "expected answer traits": every
// normative entry (outside category "ambiguous") must declare at least one
// short keyword/phrase fragment a correct answer must contain, ahead of
// any runtime Specialist answer to check the traits against.
func TestGoldenQANormativeEntriesHaveExpectedAnswerTraits(t *testing.T) {
	entries := parseGoldenQA(t)
	if len(entries) == 0 {
		t.Fatal("no entries parsed from golden-qa.yaml")
	}

	for _, entry := range entries {
		t.Run(entry.id, func(t *testing.T) {
			if entry.category == "ambiguous" {
				return
			}
			if !entry.hasAnswerTraits || len(entry.answerTraits) == 0 {
				t.Fatalf("%s: category=%s is normative and must declare at least one expected_answer_traits entry", entry.id, entry.category)
			}
			for _, trait := range entry.answerTraits {
				if strings.TrimSpace(trait) == "" {
					t.Fatalf("%s: expected_answer_traits contains an empty entry", entry.id)
				}
			}
		})
	}
}
