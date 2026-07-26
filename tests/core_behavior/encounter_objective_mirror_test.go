package corebehavior

import (
	"testing"

	"github.com/SergioLacerda/rpg-system-rgb/internal/components/core"
)

const evacuationFailureReason = "reactor collapse"

func evacuationCharacters(t *testing.T) map[string]*core.Character {
	t.Helper()
	researcher, err := core.NewCharacter("researcher", "Researcher", core.Vectors{R: 1, G: 3, B: 2}, nil)
	if err != nil {
		t.Fatal(err)
	}
	warden, err := core.NewCharacter("warden", "Warden", core.Vectors{R: 2, G: 2, B: 5}, nil)
	if err != nil {
		t.Fatal(err)
	}
	runner, err := core.NewCharacter("runner", "Runner", core.Vectors{R: 2, G: 5, B: 2}, nil)
	if err != nil {
		t.Fatal(err)
	}
	return map[string]*core.Character{
		"researcher": &researcher,
		"warden":     &warden,
		"runner":     &runner,
	}
}

func evacuationObjective() *core.Objective {
	return &core.Objective{
		DeadlineRounds: 4,
		Target:         "researcher",
		RequiredState:  core.StateDistant,
		FailureReason:  evacuationFailureReason,
	}
}

// mirrors: tests/features/encounters/laboratory_evacuation.feature#Objective met before the deadline succeeds
func TestObjectiveMetBeforeDeadlineSucceedsFeatureExample(t *testing.T) {
	characters := evacuationCharacters(t)
	encounter := core.Encounter{
		ID:   "laboratory-evacuation",
		Name: "Laboratory Evacuation",
		Actions: []core.Action{
			{
				Actor:             "runner",
				Target:            "researcher",
				PrimaryVector:     core.VectorG,
				Intent:            "open an alternate route",
				TargetStateChange: core.StateDistant,
				Procedure:         core.ProcedureReposition,
				Consequence:       "researcher reaches the exit route",
				Round:             3,
			},
		},
		Objective: evacuationObjective(),
	}

	result, err := core.RunEncounter(characters, encounter)
	if err != nil {
		t.Fatal(err)
	}
	if len(result.UndefinedSteps) > 0 {
		t.Fatalf("undefined steps: %v", result.UndefinedSteps)
	}
	if !result.Objective.Declared {
		t.Fatal("objective must be declared")
	}
	if !result.Objective.Succeeded {
		t.Fatal("objective must succeed when the researcher reaches the exit before the deadline")
	}
	if result.Objective.ResolvedRound != 3 {
		t.Fatalf("resolved round got %d want 3", result.Objective.ResolvedRound)
	}
}

// mirrors: tests/features/encounters/laboratory_evacuation.feature#Deadline reached without the objective fails the encounter
func TestDeadlineReachedWithoutObjectiveFailsFeatureExample(t *testing.T) {
	characters := evacuationCharacters(t)
	encounter := core.Encounter{
		ID:   "laboratory-evacuation",
		Name: "Laboratory Evacuation",
		Actions: []core.Action{
			{
				Actor:             "warden",
				Target:            "researcher",
				PrimaryVector:     core.VectorB,
				Intent:            "protect the researcher",
				TargetStateChange: core.StateGuarded,
				Procedure:         core.ProcedureSustain,
				Consequence:       "researcher remains protected but has not reached the exit",
				Round:             4,
			},
		},
		Objective: evacuationObjective(),
	}

	result, err := core.RunEncounter(characters, encounter)
	if err != nil {
		t.Fatal(err)
	}
	if !result.Objective.Declared {
		t.Fatal("objective must be declared")
	}
	if result.Objective.Succeeded {
		t.Fatal("objective must fail when the deadline passes without the researcher reaching the exit")
	}
	if result.Objective.ResolvedRound != 4 {
		t.Fatalf("resolved round got %d want 4", result.Objective.ResolvedRound)
	}
	if result.Objective.FailureReason != evacuationFailureReason {
		t.Fatalf("failure reason got %q want %q", result.Objective.FailureReason, evacuationFailureReason)
	}
}
