package core

import (
	"reflect"
	"testing"
)

func TestArmorCategories(t *testing.T) {
	cases := []struct {
		category   ArmorCategory
		protection int
		penalty    int
	}{
		{ArmorCategoryNone, 0, 0},
		{ArmorCategoryLight, 2, 1},
		{ArmorCategoryMedium, 4, 2},
		{ArmorCategoryHeavy, 6, 3},
	}
	for _, tc := range cases {
		if err := tc.category.Validate(); err != nil {
			t.Fatalf("%s should validate: %v", tc.category, err)
		}
		protection, err := tc.category.Protection()
		if err != nil {
			t.Fatalf("%s protection failed: %v", tc.category, err)
		}
		if protection != tc.protection {
			t.Fatalf("%s protection got %d want %d", tc.category, protection, tc.protection)
		}
		penalty, err := tc.category.MobilityPenalty()
		if err != nil {
			t.Fatalf("%s mobility failed: %v", tc.category, err)
		}
		if penalty != tc.penalty {
			t.Fatalf("%s mobility got %d want %d", tc.category, penalty, tc.penalty)
		}
	}

	unknown := ArmorCategory("unknown")
	if err := unknown.Validate(); err == nil {
		t.Fatal("expected unknown armor category to fail validation")
	}
	if _, err := unknown.Protection(); err == nil {
		t.Fatal("expected unknown armor category protection to fail")
	}
	if _, err := unknown.MobilityPenalty(); err == nil {
		t.Fatal("expected unknown armor category mobility to fail")
	}
}

func TestShieldCategories(t *testing.T) {
	cases := []struct {
		category   ShieldCategory
		protection int
		penalty    int
	}{
		{ShieldCategoryNone, 0, 0},
		{ShieldCategoryLight, 1, 0},
		{ShieldCategoryMedium, 2, 1},
		{ShieldCategoryHeavy, 3, 2},
	}
	for _, tc := range cases {
		if err := tc.category.Validate(); err != nil {
			t.Fatalf("%s should validate: %v", tc.category, err)
		}
		protection, err := tc.category.Protection()
		if err != nil {
			t.Fatalf("%s protection failed: %v", tc.category, err)
		}
		if protection != tc.protection {
			t.Fatalf("%s protection got %d want %d", tc.category, protection, tc.protection)
		}
		penalty, err := tc.category.MobilityPenalty()
		if err != nil {
			t.Fatalf("%s mobility failed: %v", tc.category, err)
		}
		if penalty != tc.penalty {
			t.Fatalf("%s mobility got %d want %d", tc.category, penalty, tc.penalty)
		}
	}

	unknown := ShieldCategory("unknown")
	if err := unknown.Validate(); err == nil {
		t.Fatal("expected unknown shield category to fail validation")
	}
	if _, err := unknown.Protection(); err == nil {
		t.Fatal("expected unknown shield category protection to fail")
	}
	if _, err := unknown.MobilityPenalty(); err == nil {
		t.Fatal("expected unknown shield category mobility to fail")
	}
}

func TestImpactAndPenetrationTables(t *testing.T) {
	impactCases := []struct {
		source      ImpactSource
		weaponValue int
		attackerR   int
		want        int
	}{
		{ImpactSourceMelee, 4, 3, 7},
		{ImpactSourceFirearm, 4, 3, 4},
		{ImpactSourceExplosive, 6, 9, 6},
		{ImpactSourceAbility, 5, 9, 5},
	}
	for _, tc := range impactCases {
		got, err := ComputeImpact(tc.source, tc.weaponValue, tc.attackerR)
		if err != nil {
			t.Fatalf("%s impact failed: %v", tc.source, err)
		}
		if got != tc.want {
			t.Fatalf("%s impact got %d want %d", tc.source, got, tc.want)
		}
	}
	if _, err := ComputeImpact("unknown", 1, 1); err == nil {
		t.Fatal("expected unknown impact source to fail")
	}

	penetrationCases := map[WeaponCategory]int{
		WeaponCategoryPistol:       1,
		WeaponCategorySMG:          2,
		WeaponCategoryCarbine:      2,
		WeaponCategoryAssaultRifle: 3,
		WeaponCategorySniperRifle:  4,
		WeaponCategoryMachineGun:   3,
	}
	for category, want := range penetrationCases {
		got, err := category.DefaultPenetration()
		if err != nil {
			t.Fatalf("%s penetration failed: %v", category, err)
		}
		if got != want {
			t.Fatalf("%s penetration got %d want %d", category, got, want)
		}
	}
	if _, err := WeaponCategory("unknown").DefaultPenetration(); err == nil {
		t.Fatal("expected unknown weapon category to fail")
	}
}

func TestMovementGrappleInitiativeAndRecovery(t *testing.T) {
	if MovementDistance(4) != 8 {
		t.Fatalf("movement distance got %d want 8", MovementDistance(4))
	}
	effective, err := EffectiveG(5, ArmorCategoryMedium, ShieldCategoryHeavy)
	if err != nil {
		t.Fatal(err)
	}
	if effective != 1 {
		t.Fatalf("effective G got %d want 1", effective)
	}
	effective, err = EffectiveG(2, ArmorCategoryHeavy, ShieldCategoryHeavy)
	if err != nil {
		t.Fatal(err)
	}
	if effective != 0 {
		t.Fatalf("effective G must floor at zero, got %d", effective)
	}
	if _, err := EffectiveG(5, ArmorCategory("unknown"), ShieldCategoryNone); err == nil {
		t.Fatal("expected invalid armor to fail EffectiveG")
	}
	if _, err := EffectiveG(5, ArmorCategoryNone, ShieldCategory("unknown")); err == nil {
		t.Fatal("expected invalid shield to fail EffectiveG")
	}

	if !ResolveGrapple(3, 4, 3, 2) {
		t.Fatal("grapple should succeed when both G and R meet target")
	}
	if ResolveGrapple(2, 4, 3, 2) || ResolveGrapple(3, 1, 3, 2) {
		t.Fatal("grapple should fail when either G or R is below target")
	}

	order := InitiativeOrder([]InitiativeEntry{
		{ActorID: "slow", G: 1},
		{ActorID: "surprised", G: 0, Surprise: true},
		{ActorID: "fast", G: 5},
		{ActorID: "tie-a", G: 3},
		{ActorID: "tie-b", G: 3},
	})
	wantOrder := []string{"surprised", "fast", "tie-a", "tie-b", "slow"}
	if !reflect.DeepEqual(order, wantOrder) {
		t.Fatalf("initiative order got %v want %v", order, wantOrder)
	}

	character, err := NewCharacter("medic", "Medic", Vectors{R: 1, G: 2, B: 3}, nil)
	if err != nil {
		t.Fatal(err)
	}
	if err := character.AddState(StateInjured); err != nil {
		t.Fatal(err)
	}
	if err := Stabilize(&character); err != nil {
		t.Fatal(err)
	}
	if character.HasState(StateInjured) || !character.HasState(StateStabilized) {
		t.Fatalf("stabilize states mismatch: %+v", character.States)
	}
	if err := Stabilize(nil); err == nil {
		t.Fatal("expected nil character stabilization to fail")
	}
}

func TestVectorsStatesResourcesAndObjectives(t *testing.T) {
	vectors := Vectors{R: 1, G: 2, B: 3}
	for _, vector := range []Vector{VectorR, VectorG, VectorB} {
		if err := vector.Validate(); err != nil {
			t.Fatalf("%s should validate: %v", vector, err)
		}
		if _, err := vectors.Value(vector); err != nil {
			t.Fatalf("%s value failed: %v", vector, err)
		}
		if vector.TeachingLabel() == "unknown" || vector.NormativeLabel() == "unknown" {
			t.Fatalf("%s should expose known labels", vector)
		}
	}
	if err := Vector("X").Validate(); err == nil {
		t.Fatal("expected unknown vector to fail validation")
	}
	if _, err := vectors.Value("X"); err == nil {
		t.Fatal("expected unknown vector value lookup to fail")
	}
	if Vector("X").TeachingLabel() != "unknown" || Vector("X").NormativeLabel() != "unknown" {
		t.Fatal("unknown vector should use unknown labels")
	}
	if err := (Vectors{R: -1}).Validate(); err == nil {
		t.Fatal("expected negative vectors to fail validation")
	}

	states := AllStates()
	if len(states) != len(allStates) {
		t.Fatalf("AllStates got %d states want %d", len(states), len(allStates))
	}
	states[0] = "mutated"
	if allStates[0] == "mutated" {
		t.Fatal("AllStates must return a defensive copy")
	}
	for _, state := range allStates {
		lifecycle, ok := Lifecycle(state)
		if !ok {
			t.Fatalf("%s missing lifecycle", state)
		}
		if lifecycle.ProducedBy == "" || lifecycle.ClearedBy == "" {
			t.Fatalf("%s lifecycle must be complete: %+v", state, lifecycle)
		}
	}
	if _, ok := Lifecycle("unknown"); ok {
		t.Fatal("unknown state must not have lifecycle")
	}

	if !((Resources{CurrentHealth: 0}).IsDown()) {
		t.Fatal("zero health should be down")
	}
	if (Resources{CurrentHealth: 1}).IsDown() {
		t.Fatal("positive health should not be down")
	}

	validObjective := Objective{DeadlineRounds: 2, Target: "target", RequiredState: StateCovered, FailureReason: "too late"}
	if err := validObjective.Validate(); err != nil {
		t.Fatal(err)
	}
	for _, objective := range []Objective{
		{DeadlineRounds: 0, Target: "target", RequiredState: StateCovered, FailureReason: "too late"},
		{DeadlineRounds: 1, RequiredState: StateCovered, FailureReason: "too late"},
		{DeadlineRounds: 1, Target: "target", RequiredState: "unknown", FailureReason: "too late"},
		{DeadlineRounds: 1, Target: "target", RequiredState: StateCovered},
	} {
		if err := objective.Validate(); err == nil {
			t.Fatalf("expected invalid objective to fail: %+v", objective)
		}
	}
}

func TestEncounterObjectiveOutcomes(t *testing.T) {
	target, err := NewCharacter("target", "Target", Vectors{R: 1, G: 1, B: 1}, nil)
	if err != nil {
		t.Fatal(err)
	}
	objective := Objective{DeadlineRounds: 3, Target: "target", RequiredState: StateCovered, FailureReason: "exposed"}
	target.States[StateCovered] = true
	succeeded := evaluateObjective(map[string]*Character{"target": &target}, objective, 2)
	if !succeeded.Succeeded || succeeded.ResolvedRound != 2 || succeeded.FailureReason != "exposed" {
		t.Fatalf("unexpected success outcome: %+v", succeeded)
	}

	delete(target.States, StateCovered)
	failed := evaluateObjective(map[string]*Character{"target": &target}, objective, 4)
	if failed.Succeeded || failed.ResolvedRound != 4 || failed.FailureReason != "exposed" {
		t.Fatalf("unexpected failure outcome: %+v", failed)
	}
	missing := evaluateObjective(map[string]*Character{}, objective, 0)
	if missing.Succeeded || missing.ResolvedRound != 3 {
		t.Fatalf("missing target should fail at deadline: %+v", missing)
	}

	invalidResult, err := RunEncounter(map[string]*Character{}, Encounter{ID: "bad", Name: "Bad", Actions: []Action{{}}})
	if err != nil {
		t.Fatal(err)
	}
	if len(invalidResult.UndefinedSteps) != 1 {
		t.Fatalf("invalid action should surface as undefined step, got %+v", invalidResult)
	}
	result, err := RunEncounter(map[string]*Character{}, Encounter{
		ID: "partial", Name: "Partial", Actions: []Action{{
			Actor: "missing", Target: "target", PrimaryVector: VectorR, SecondaryVector: VectorG,
			Intent: "act", Procedure: ProcedureEvade, Consequence: "none",
		}},
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(result.UndefinedSteps) != 1 {
		t.Fatalf("expected one undefined step, got %+v", result)
	}
}
