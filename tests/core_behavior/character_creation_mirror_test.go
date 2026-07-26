package corebehavior

import (
	"testing"

	"github.com/SergioLacerda/rpg-system-rgb/internal/components/core"
)

// mirrors: tests/features/core/character_creation.feature#Health and shield derivation
func TestHealthAndShieldDerivationFeatureExamples(t *testing.T) {
	cases := []struct {
		r, g, b        int
		health, shield int
	}{
		{r: 3, g: 2, b: 2, health: 9, shield: 6},
		{r: 2, g: 2, b: 3, health: 9, shield: 9},
		{r: 0, g: 7, b: 0, health: 4, shield: 0},
	}
	for _, tc := range cases {
		character, err := core.NewCharacter("fixture", "Fixture", core.Vectors{R: tc.r, G: tc.g, B: tc.b}, nil)
		if err != nil {
			t.Fatal(err)
		}
		if character.Resources.MaxHealth != tc.health {
			t.Fatalf("R%d/G%d/B%d health got %d want %d", tc.r, tc.g, tc.b, character.Resources.MaxHealth, tc.health)
		}
		if character.Resources.MaxShield != tc.shield {
			t.Fatalf("R%d/G%d/B%d shield got %d want %d", tc.r, tc.g, tc.b, character.Resources.MaxShield, tc.shield)
		}
	}
}

// mirrors: tests/features/core/character_creation.feature#No preservation investment means no shield state
func TestNoPreservationMeansNoShieldStateFeatureExample(t *testing.T) {
	character, err := core.NewCharacter("fixture", "Fixture", core.Vectors{R: 4, G: 3, B: 0}, nil)
	if err != nil {
		t.Fatal(err)
	}
	if character.HasState(core.StateShielded) {
		t.Fatal("a character with B 0 must not start shielded")
	}
	if character.Resources.MaxShield != 0 {
		t.Fatalf("max shield got %d want 0", character.Resources.MaxShield)
	}
}
