// Package fixtures provides reusable Core V2 characters and encounters for
// tests, simulations, and playtests. It mirrors tests/fixtures/*.yaml; see
// tests/fixtures/fixture_files_test.go for the parity guardrail.
package fixtures

import "github.com/SergioLacerda/rpg-system-rgb/internal/components/core"

// Characters returns the four canonical Core V2 test characters, keyed by ID.
func Characters() map[string]*core.Character {
	characters := map[string]*core.Character{}
	for _, character := range []core.Character{
		mustCharacter("red-vanguard", "Red Vanguard", core.Vectors{R: 5, G: 2, B: 2}),
		mustCharacter("green-runner", "Green Runner", core.Vectors{R: 2, G: 5, B: 2}),
		mustCharacter("blue-warden", "Blue Warden", core.Vectors{R: 2, G: 2, B: 5}),
		mustCharacter("rgb-balanced", "RGB Balanced", core.Vectors{R: 3, G: 3, B: 3}),
	} {
		characterCopy := character
		characters[characterCopy.ID] = &characterCopy
	}
	return characters
}

// Encounters returns the canonical Core V2 test encounters that exercise
// the four fixture characters across attack, movement, and damage loops.
func Encounters() []core.Encounter {
	return []core.Encounter{
		{
			ID:   "attack-defense-loop",
			Name: "Attack And Defense Loop",
			Actions: []core.Action{
				{
					Actor:             "red-vanguard",
					Target:            "green-runner",
					PrimaryVector:     core.VectorR,
					Intent:            "pressure target with a clean strike",
					TargetStateChange: core.StateVulnerable,
					Procedure:         core.ProcedureEvade,
					Consequence:       "target becomes vulnerable when pressure beats evasion",
				},
			},
		},
		{
			ID:   "movement-state-loop",
			Name: "Movement And State Loop",
			Actions: []core.Action{
				{
					Actor:             "green-runner",
					Target:            "blue-warden",
					PrimaryVector:     core.VectorG,
					Intent:            "reposition into cover",
					TargetStateChange: core.StateCovered,
					Procedure:         core.ProcedureSustain,
					Consequence:       "actor changes the relation to pressure",
				},
			},
		},
		{
			ID:   "damage-mitigation-loop",
			Name: "Damage And Mitigation Loop",
			Actions: []core.Action{
				{
					Actor:             "rgb-balanced",
					Target:            "red-vanguard",
					PrimaryVector:     core.VectorB,
					Intent:            "test mitigation under pressure",
					TargetStateChange: core.StateInjured,
					Procedure:         core.ProcedureBlock,
					Consequence:       "damage pipeline decides health consequence",
				},
			},
		},
	}
}

func mustCharacter(id, name string, vectors core.Vectors) core.Character {
	ability := core.Ability{
		ID:       id + "-signature",
		Name:     name + " Signature",
		Vector:   strongestVector(vectors),
		Tier:     1,
		Timing:   core.TimingAction,
		Cost:     0,
		Range:    1,
		Duration: "instant",
		Effects:  []string{"declare_vector_identity"},
		Tags:     []string{"fixture"},
	}
	character, err := core.NewCharacter(id, name, vectors, []core.Ability{ability})
	if err != nil {
		panic(err)
	}
	return character
}

func strongestVector(vectors core.Vectors) core.Vector {
	if vectors.R >= vectors.G && vectors.R >= vectors.B {
		return core.VectorR
	}
	if vectors.G >= vectors.B {
		return core.VectorG
	}
	return core.VectorB
}
