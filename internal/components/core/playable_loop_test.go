package core_test

import (
	"testing"

	"github.com/SergioLacerda/rpg-system-rgb/internal/components/core"
	"github.com/SergioLacerda/rpg-system-rgb/internal/components/core/fixtures"
)

func TestFourCharactersCanRunThreeCoreEncounters(t *testing.T) {
	characters := fixtures.Characters()
	if len(characters) != 4 {
		t.Fatalf("expected four characters, got %d", len(characters))
	}
	for id, character := range characters {
		if character.ID != id {
			t.Fatalf("character map key %s does not match character ID %s", id, character.ID)
		}
		if len(character.Abilities) == 0 {
			t.Fatalf("%s must include at least one fixture ability", id)
		}
	}

	encounters := fixtures.Encounters()
	if len(encounters) != 3 {
		t.Fatalf("expected three encounters, got %d", len(encounters))
	}
	for _, encounter := range encounters {
		result, err := core.RunEncounter(characters, encounter)
		if err != nil {
			t.Fatalf("%s failed: %v", encounter.ID, err)
		}
		if len(result.UndefinedSteps) > 0 {
			t.Fatalf("%s has undefined rule steps: %v", encounter.ID, result.UndefinedSteps)
		}
		if len(result.ActionResults) != len(encounter.Actions) {
			t.Fatalf("%s action results got %d want %d", encounter.ID, len(result.ActionResults), len(encounter.Actions))
		}
	}
}
