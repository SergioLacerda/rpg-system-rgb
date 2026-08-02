package fixtures

import (
	"testing"

	"github.com/SergioLacerda/rpg-system-rgb/internal/components/core"
)

func TestCharactersReturnIndependentCopies(t *testing.T) {
	first := Characters()
	second := Characters()
	first["red-vanguard"].Name = "mutated"
	if second["red-vanguard"].Name == "mutated" {
		t.Fatal("Characters must return independent character copies")
	}
	if strongestVector(core.Vectors{R: 3, G: 2, B: 1}) != core.VectorR {
		t.Fatal("R should win strongest vector ties by priority")
	}
	if strongestVector(core.Vectors{R: 1, G: 3, B: 2}) != core.VectorG {
		t.Fatal("G should win when greater than B")
	}
	if strongestVector(core.Vectors{R: 1, G: 2, B: 3}) != core.VectorB {
		t.Fatal("B should win when highest")
	}
}

func TestEncountersExposeValidCanonicalLoops(t *testing.T) {
	characters := Characters()
	for _, encounter := range Encounters() {
		if encounter.ID == "" || encounter.Name == "" || len(encounter.Actions) == 0 {
			t.Fatalf("invalid encounter fixture: %+v", encounter)
		}
		if _, err := core.RunEncounter(characters, encounter); err != nil {
			t.Fatalf("%s should run: %v", encounter.ID, err)
		}
	}
}
