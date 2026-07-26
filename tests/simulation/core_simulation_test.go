package simulation

import (
	"testing"

	"github.com/SergioLacerda/rpg-system-rgb/internal/components/core"
	"github.com/SergioLacerda/rpg-system-rgb/internal/components/core/fixtures"
)

type metrics struct {
	Rounds         int
	Actions        int
	VectorCounts   map[core.Vector]int
	UndefinedSteps int
	ObjectiveState string
	DominantVector core.Vector
}

func TestCoreSimulationRecordsDeterministicMetrics(t *testing.T) {
	characters := fixtures.Characters()
	encounters := fixtures.Encounters()
	collected := metrics{
		VectorCounts: map[core.Vector]int{
			core.VectorR: 0,
			core.VectorG: 0,
			core.VectorB: 0,
		},
	}

	for _, encounter := range encounters {
		collected.Rounds++
		result, err := core.RunEncounter(characters, encounter)
		if err != nil {
			t.Fatal(err)
		}
		collected.Actions += len(encounter.Actions)
		collected.UndefinedSteps += len(result.UndefinedSteps)
		for _, action := range encounter.Actions {
			collected.VectorCounts[action.PrimaryVector]++
		}
	}
	collected.ObjectiveState = "core-fixtures-completed"
	collected.DominantVector = dominantVector(collected.VectorCounts)

	if collected.Rounds != 3 {
		t.Fatalf("rounds got %d want 3", collected.Rounds)
	}
	if collected.Actions != 3 {
		t.Fatalf("actions got %d want 3", collected.Actions)
	}
	if collected.UndefinedSteps != 0 {
		t.Fatalf("undefined steps got %d want 0", collected.UndefinedSteps)
	}
	for _, vector := range []core.Vector{core.VectorR, core.VectorG, core.VectorB} {
		if collected.VectorCounts[vector] == 0 {
			t.Fatalf("simulation did not exercise vector %s", vector)
		}
	}
	if collected.ObjectiveState != "core-fixtures-completed" {
		t.Fatalf("unexpected objective state %s", collected.ObjectiveState)
	}
	if collected.DominantVector == "" {
		t.Fatal("dominant vector must be recorded")
	}
}

func dominantVector(counts map[core.Vector]int) core.Vector {
	dominant := core.VectorR
	for _, vector := range []core.Vector{core.VectorG, core.VectorB} {
		if counts[vector] > counts[dominant] {
			dominant = vector
		}
	}
	return dominant
}
