package core

import (
	"reflect"
	"testing"
)

func TestVectorsExposeTeachingAndNormativeLabels(t *testing.T) {
	vectors := Vectors{R: 3, G: 2, B: 1}
	if err := vectors.Validate(); err != nil {
		t.Fatal(err)
	}
	if value, err := vectors.Value(VectorR); err != nil || value != 3 {
		t.Fatalf("R value mismatch: got %d err %v", value, err)
	}
	if VectorR.TeachingLabel() != "attack" {
		t.Fatalf("unexpected R teaching label: %s", VectorR.TeachingLabel())
	}
	if VectorB.NormativeLabel() != "preserve continuity under pressure" {
		t.Fatalf("unexpected B normative label: %s", VectorB.NormativeLabel())
	}
}

func TestCharacterCreationDerivesResourcesFromRAndB(t *testing.T) {
	character, err := NewCharacter("warden", "Warden", Vectors{R: 1, G: 2, B: 5}, nil)
	if err != nil {
		t.Fatal(err)
	}
	if character.Resources.MaxHealth != 10 {
		t.Fatalf("health must use base + R + B, got %d", character.Resources.MaxHealth)
	}
	if character.Resources.MaxShield != 5 {
		t.Fatalf("shield must derive from B, got %d", character.Resources.MaxShield)
	}
	if !character.HasState(StateHealthy) {
		t.Fatal("new character must start healthy")
	}
	if !character.HasState(StateShielded) {
		t.Fatal("character with B-derived shield must start shielded")
	}
}

func TestActionContractValidation(t *testing.T) {
	action := Action{
		Actor:             "actor",
		Target:            "target",
		PrimaryVector:     VectorR,
		SecondaryVector:   VectorG,
		Intent:            "suppress while advancing",
		TargetStateChange: StateSuppressed,
		Procedure:         ProcedureEvade,
		Consequence:       "target loses initiative",
	}
	if err := action.Validate(); err != nil {
		t.Fatal(err)
	}
}

func TestMarginResolutionBands(t *testing.T) {
	cases := map[int]Outcome{
		4:  OutcomeStrongSuccess,
		1:  OutcomeSuccess,
		0:  OutcomeSuccessWithCost,
		-1: OutcomeFailureWithOpportunity,
		-4: OutcomeClearFailure,
	}
	for margin, expected := range cases {
		if got := ClassifyMargin(margin); got != expected {
			t.Fatalf("margin %d classified as %s, want %s", margin, got, expected)
		}
	}

	resolution := Resolve(5, 1, 3)
	if resolution.Margin != 3 || resolution.Outcome != OutcomeSuccess {
		t.Fatalf("unexpected resolution: %#v", resolution)
	}
}

func TestDefenseProcedureVectorOwnership(t *testing.T) {
	cases := map[Procedure]Vector{
		ProcedureEvade:           VectorG,
		ProcedureReposition:      VectorG,
		ProcedureBlock:           VectorB,
		ProcedureSustain:         VectorB,
		ProcedureInterrupt:       VectorR,
		ProcedureCounterpressure: VectorR,
	}
	for procedure, expected := range cases {
		if got := procedure.Vector(); got != expected {
			t.Fatalf("%s vector owner got %s want %s", procedure, got, expected)
		}
	}

	character, err := NewCharacter("runner", "Runner", Vectors{R: 1, G: 5, B: 2}, nil)
	if err != nil {
		t.Fatal(err)
	}
	value, err := DefenseValue(character, ProcedureEvade)
	if err != nil {
		t.Fatal(err)
	}
	if value != 5 {
		t.Fatalf("evade must use G, got %d", value)
	}
}

func TestDamagePipelineOrderAndStateConsequences(t *testing.T) {
	target, err := NewCharacter("target", "Target", Vectors{R: 2, G: 2, B: 3}, nil)
	if err != nil {
		t.Fatal(err)
	}
	target.Resources.Armor = 2

	result, err := ApplyDamage(&target, DamageInput{Impact: 7, Penetration: 1})
	if err != nil {
		t.Fatal(err)
	}

	var stepNames []string
	for _, step := range result.Steps {
		stepNames = append(stepNames, step.Name)
	}
	expectedSteps := []string{"impact", "penetration", "armor_reduction", "shield_absorption", "health_consequence"}
	if !reflect.DeepEqual(stepNames, expectedSteps) {
		t.Fatalf("damage pipeline order got %v want %v", stepNames, expectedSteps)
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
	if !target.HasState(StateInjured) {
		t.Fatal("health damage must mark target injured")
	}
}

func TestStateTransitions(t *testing.T) {
	character, err := NewCharacter("scout", "Scout", Vectors{R: 1, G: 4, B: 1}, nil)
	if err != nil {
		t.Fatal(err)
	}
	if err := character.AddState(StateHidden); err != nil {
		t.Fatal(err)
	}
	if !character.HasState(StateHidden) {
		t.Fatal("hidden state should be active")
	}
	character.RemoveState(StateHidden)
	if character.HasState(StateHidden) {
		t.Fatal("hidden state should be removed")
	}
	if err := character.AddState("unknown"); err == nil {
		t.Fatal("unknown state must fail validation")
	}
}

func TestAbilityContractValidation(t *testing.T) {
	ability := Ability{
		ID:           "phantom-step",
		Name:         "Phantom Step",
		Vector:       VectorG,
		Tier:         2,
		Requirements: map[Vector]int{VectorG: 4},
		Timing:       TimingReaction,
		Cost:         2,
		Range:        4,
		Duration:     "instant",
		Effects:      []string{"move_without_reaction"},
		Limits:       []string{"twice_per_scene"},
		Tags:         []string{"mobility", "supernatural"},
	}
	if err := ability.Validate(); err != nil {
		t.Fatal(err)
	}
}
