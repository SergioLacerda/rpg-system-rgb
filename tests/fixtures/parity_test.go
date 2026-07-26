package fixtures

import (
	"os"
	"regexp"
	"strconv"
	"strings"
	"testing"

	corefixtures "github.com/SergioLacerda/rpg-system-rgb/internal/components/core/fixtures"
)

var yamlCharacterBlock = regexp.MustCompile(`(?s)- id:\s*(\S+).*?R:\s*(\d+).*?G:\s*(\d+).*?B:\s*(\d+)`)

type yamlVectors struct {
	r, g, b int
}

// parseYAMLCharacters extracts id and R/G/B vectors from
// characters/core-v2.yaml by splitting on each "- id:" entry. It is a
// line-oriented reader rather than a full YAML parser, matching the fixture
// file's regular, hand-authored structure.
func parseYAMLCharacters(t *testing.T) map[string]yamlVectors {
	t.Helper()
	body, err := os.ReadFile("characters/core-v2.yaml")
	if err != nil {
		t.Fatal(err)
	}

	characters := map[string]yamlVectors{}
	blocks := strings.Split(string(body), "\n  - id:")
	for i, block := range blocks {
		if i == 0 {
			continue // header content before the first character entry
		}
		match := yamlCharacterBlock.FindStringSubmatch("  - id:" + block)
		if match == nil {
			t.Fatalf("could not parse character block starting with: %.40s", block)
		}
		r, _ := strconv.Atoi(match[2])
		g, _ := strconv.Atoi(match[3])
		b, _ := strconv.Atoi(match[4])
		characters[match[1]] = yamlVectors{r: r, g: g, b: b}
	}
	return characters
}

// mirrors: tests/features/fixtures/fixture_parity.feature#Character vectors match across fixture sources
// TestFixtureParityBetweenYAMLAndGo keeps the YAML fixture (authoring
// surface) and the Go fixture (test/simulation surface) from silently
// drifting apart (structural review F5): every character ID and R/G/B
// triple must match across both sources.
func TestFixtureParityBetweenYAMLAndGo(t *testing.T) {
	yamlCharacters := parseYAMLCharacters(t)
	goCharacters := corefixtures.Characters()

	if len(yamlCharacters) != len(goCharacters) {
		t.Fatalf("character count mismatch: yaml=%d go=%d", len(yamlCharacters), len(goCharacters))
	}
	for id, want := range yamlCharacters {
		got, ok := goCharacters[id]
		if !ok {
			t.Fatalf("YAML character %q has no Go fixture counterpart in internal/components/core/fixtures", id)
		}
		if got.Vectors.R != want.r || got.Vectors.G != want.g || got.Vectors.B != want.b {
			t.Fatalf(
				"%s vectors mismatch: yaml=R%d/G%d/B%d go=R%d/G%d/B%d",
				id, want.r, want.g, want.b, got.Vectors.R, got.Vectors.G, got.Vectors.B,
			)
		}
	}
	for id := range goCharacters {
		if _, ok := yamlCharacters[id]; !ok {
			t.Fatalf("Go fixture character %q has no YAML fixture counterpart in characters/core-v2.yaml", id)
		}
	}
}
