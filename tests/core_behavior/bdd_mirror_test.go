package corebehavior

import (
	"testing"

	"github.com/SergioLacerda/rpg-system-rgb/internal/components/core"
)

func TestAttackEvasionFeatureExamples(t *testing.T) {
	cases := []struct {
		name         string
		attack       int
		evasion      int
		outcome      core.Outcome
		damageStarts bool
	}{
		{name: "below", attack: 2, evasion: 4, outcome: core.OutcomeFailureWithOpportunity, damageStarts: false},
		{name: "equal", attack: 3, evasion: 3, outcome: core.OutcomeSuccessWithCost, damageStarts: true},
		{name: "above", attack: 5, evasion: 3, outcome: core.OutcomeSuccess, damageStarts: true},
		{name: "strong", attack: 7, evasion: 3, outcome: core.OutcomeStrongSuccess, damageStarts: true},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			resolution := core.Resolve(tc.attack, 0, tc.evasion)
			if resolution.Outcome != tc.outcome {
				t.Fatalf("outcome got %s want %s", resolution.Outcome, tc.outcome)
			}
			if resolution.Outcome.Successful() != tc.damageStarts {
				t.Fatalf("damage start flag got %v want %v", resolution.Outcome.Successful(), tc.damageStarts)
			}
		})
	}
}

func TestDamageFeatureArmorShieldHealthFlow(t *testing.T) {
	target, err := core.NewCharacter("target", "Target", core.Vectors{R: 2, G: 2, B: 3}, nil)
	if err != nil {
		t.Fatal(err)
	}
	target.Resources.Armor = 2

	result, err := core.ApplyDamage(&target, core.DamageInput{Impact: 7, Penetration: 1})
	if err != nil {
		t.Fatal(err)
	}
	if result.ArmorReduced != 1 {
		t.Fatalf("armor reduced got %d want 1", result.ArmorReduced)
	}
	if result.ShieldAbsorbed != 3 {
		t.Fatalf("shield absorbed got %d want 3", result.ShieldAbsorbed)
	}
	if result.HealthDamage != 3 {
		t.Fatalf("health damage got %d want 3", result.HealthDamage)
	}
	if !target.HasState(core.StateInjured) {
		t.Fatal("target must become injured")
	}
}

func TestRGBObstacleApproachFeatureExamples(t *testing.T) {
	cases := []struct {
		name      string
		vector    core.Vector
		procedure core.Procedure
		state     core.State
	}{
		{name: "R transforms source", vector: core.VectorR, procedure: core.ProcedureCounterpressure, state: core.StateVulnerable},
		{name: "G changes relation", vector: core.VectorG, procedure: core.ProcedureReposition, state: core.StateCovered},
		{name: "B preserves action", vector: core.VectorB, procedure: core.ProcedureSustain, state: core.StateGuarded},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.procedure.Vector() != tc.vector {
				t.Fatalf("%s owner got %s want %s", tc.procedure, tc.procedure.Vector(), tc.vector)
			}
			actor, err := core.NewCharacter("actor", "Actor", core.Vectors{R: 4, G: 4, B: 4}, nil)
			if err != nil {
				t.Fatal(err)
			}
			if err := actor.AddState(tc.state); err != nil {
				t.Fatal(err)
			}
			if !actor.HasState(tc.state) {
				t.Fatalf("actor missing state %s", tc.state)
			}
		})
	}
}

func TestEncounterFeatureSucceedsWithoutDefeatingAllOpponents(t *testing.T) {
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
	mercenary, err := core.NewCharacter("mercenary", "Mercenary", core.Vectors{R: 3, G: 3, B: 2}, nil)
	if err != nil {
		t.Fatal(err)
	}

	characters := map[string]*core.Character{
		"researcher": &researcher,
		"warden":     &warden,
		"runner":     &runner,
		"mercenary":  &mercenary,
	}
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
				Consequence:       "researcher remains protected",
			},
			{
				Actor:             "runner",
				Target:            "researcher",
				PrimaryVector:     core.VectorG,
				Intent:            "open an alternate route",
				TargetStateChange: core.StateDistant,
				Procedure:         core.ProcedureReposition,
				Consequence:       "researcher reaches the exit route",
			},
		},
	}

	result, err := core.RunEncounter(characters, encounter)
	if err != nil {
		t.Fatal(err)
	}
	if len(result.UndefinedSteps) > 0 {
		t.Fatalf("undefined steps: %v", result.UndefinedSteps)
	}
	if mercenary.Resources.IsDown() {
		t.Fatal("encounter success must not require defeating every mercenary")
	}
	if !researcher.HasState(core.StateGuarded) || !researcher.HasState(core.StateDistant) {
		t.Fatalf("researcher objective states missing: guarded=%v distant=%v", researcher.HasState(core.StateGuarded), researcher.HasState(core.StateDistant))
	}
}
