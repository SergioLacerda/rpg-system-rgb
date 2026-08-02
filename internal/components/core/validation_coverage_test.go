package core

import "testing"

func validCoverageAbility() Ability {
	return Ability{
		ID:           "coverage",
		Name:         "Coverage",
		Vector:       VectorR,
		Tier:         1,
		Requirements: map[Vector]int{VectorR: 1},
		Timing:       TimingAction,
		Cost:         1,
		Range:        1,
		Duration:     "instant",
		Effects:      []string{"effect"},
		Tags:         []string{"tag"},
	}
}

func TestAbilityValidationRejectsEveryRequiredField(t *testing.T) {
	cases := []struct {
		name   string
		mutate func(*Ability)
	}{
		{"empty ID", func(a *Ability) { a.ID = "" }},
		{"empty name", func(a *Ability) { a.Name = "" }},
		{"invalid vector", func(a *Ability) { a.Vector = "X" }},
		{"invalid tier", func(a *Ability) { a.Tier = 0 }},
		{"invalid requirement vector", func(a *Ability) { a.Requirements = map[Vector]int{"X": 1} }},
		{"negative requirement", func(a *Ability) { a.Requirements = map[Vector]int{VectorR: -1} }},
		{"invalid timing", func(a *Ability) { a.Timing = "instant" }},
		{"negative cost", func(a *Ability) { a.Cost = -1 }},
		{"negative range", func(a *Ability) { a.Range = -1 }},
		{"empty duration", func(a *Ability) { a.Duration = "" }},
		{"empty effects", func(a *Ability) { a.Effects = nil }},
		{"empty tags", func(a *Ability) { a.Tags = nil }},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			ability := validCoverageAbility()
			tc.mutate(&ability)
			if err := ability.Validate(); err == nil {
				t.Fatalf("expected %s to fail validation", tc.name)
			}
		})
	}
}

func validCoverageAction() Action {
	return Action{
		Actor:             "actor",
		Target:            "target",
		PrimaryVector:     VectorR,
		SecondaryVector:   VectorG,
		Intent:            "do a thing",
		TargetStateChange: StateExposed,
		Procedure:         ProcedureAttack,
		Consequence:       "thing happens",
	}
}

func TestActionValidationRejectsEveryRequiredField(t *testing.T) {
	cases := []struct {
		name   string
		mutate func(*Action)
	}{
		{"negative round", func(a *Action) { a.Round = -1 }},
		{"empty actor", func(a *Action) { a.Actor = "" }},
		{"empty target", func(a *Action) { a.Target = "" }},
		{"invalid primary vector", func(a *Action) { a.PrimaryVector = "X" }},
		{"invalid secondary vector", func(a *Action) { a.SecondaryVector = "X" }},
		{"empty intent", func(a *Action) { a.Intent = "" }},
		{"invalid target state", func(a *Action) { a.TargetStateChange = "unknown" }},
		{"invalid procedure", func(a *Action) { a.Procedure = "unknown" }},
		{"empty consequence", func(a *Action) { a.Consequence = "" }},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			action := validCoverageAction()
			tc.mutate(&action)
			if err := action.Validate(); err == nil {
				t.Fatalf("expected %s to fail validation", tc.name)
			}
		})
	}
}

func TestCharacterCreationValidationAndStartingBudget(t *testing.T) {
	if _, err := NewCharacter("", "Hero", Vectors{R: 1, G: 1, B: 1}, nil); err == nil {
		t.Fatal("expected empty character ID to fail")
	}
	if _, err := NewCharacter("hero", "", Vectors{R: 1, G: 1, B: 1}, nil); err == nil {
		t.Fatal("expected empty character name to fail")
	}
	if _, err := NewCharacter("hero", "Hero", Vectors{R: -1}, nil); err == nil {
		t.Fatal("expected invalid vectors to fail")
	}
	invalidAbility := validCoverageAbility()
	invalidAbility.ID = ""
	if _, err := NewCharacter("hero", "Hero", Vectors{R: 1, G: 1, B: 1}, []Ability{invalidAbility}); err == nil {
		t.Fatal("expected invalid ability to fail character creation")
	}
	if _, err := NewStartingCharacter("hero", "Hero", Vectors{R: 1, G: 1, B: 1}, nil); err == nil {
		t.Fatal("expected wrong starting budget to fail")
	}
	character, err := NewStartingCharacter("hero", "Hero", Vectors{R: 3, G: 2, B: 2}, nil)
	if err != nil {
		t.Fatal(err)
	}
	if character.Level != 1 || !character.HasState(StateHealthy) {
		t.Fatalf("unexpected starting character: %+v", character)
	}
}

func TestDefenseValidationAndValues(t *testing.T) {
	procedures := []Procedure{
		ProcedureAttack, ProcedureEvade, ProcedureReposition, ProcedureBlock,
		ProcedureSustain, ProcedureInterrupt, ProcedureCounterpressure,
	}
	for _, procedure := range procedures {
		if err := procedure.Validate(); err != nil {
			t.Fatalf("%s should validate: %v", procedure, err)
		}
		if procedure.Vector() == "" {
			t.Fatalf("%s should have vector owner", procedure)
		}
	}
	if err := Procedure("unknown").Validate(); err == nil {
		t.Fatal("expected unknown procedure to fail validation")
	}
	if Procedure("unknown").Vector() != "" {
		t.Fatal("unknown procedure should not have vector owner")
	}
	if _, err := DefenseValue(Character{Vectors: Vectors{R: 1}}, "unknown"); err == nil {
		t.Fatal("expected unknown defense procedure to fail")
	}
}

func TestApplyDamageValidationAndDownedConsequences(t *testing.T) {
	if _, err := ApplyDamage(nil, DamageInput{}); err == nil {
		t.Fatal("expected nil target to fail")
	}
	target, err := NewCharacter("target", "Target", Vectors{R: 1, G: 1, B: 1}, nil)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := ApplyDamage(&target, DamageInput{Impact: -1}); err == nil {
		t.Fatal("expected negative impact to fail")
	}
	if _, err := ApplyDamage(&target, DamageInput{Penetration: -1}); err == nil {
		t.Fatal("expected negative penetration to fail")
	}
	target.Resources.Armor = 10
	target.Resources.CurrentShield = 0
	result, err := ApplyDamage(&target, DamageInput{Impact: 2})
	if err != nil {
		t.Fatal(err)
	}
	if result.HealthDamage != 0 || !target.HasState(StateHealthy) {
		t.Fatalf("fully armored target should remain healthy: result=%+v states=%+v", result, target.States)
	}

	target.Resources.CurrentHealth = 1
	target.Resources.Armor = 0
	result, err = ApplyDamage(&target, DamageInput{Impact: 10})
	if err != nil {
		t.Fatal(err)
	}
	if result.HealthDamage != 10 || !target.HasState(StateDowned) || target.HasState(StateHealthy) {
		t.Fatalf("lethal damage should mark downed: result=%+v states=%+v", result, target.States)
	}
}

func TestEncounterValidationAndMissingTargets(t *testing.T) {
	valid := Encounter{ID: "enc", Name: "Encounter", Actions: []Action{validCoverageAction()}}
	if err := validateEncounterShape(valid); err != nil {
		t.Fatal(err)
	}
	for _, encounter := range []Encounter{
		{Name: "Encounter", Actions: []Action{validCoverageAction()}},
		{ID: "enc", Actions: []Action{validCoverageAction()}},
		{ID: "enc", Name: "Encounter"},
		{ID: "enc", Name: "Encounter", Actions: []Action{validCoverageAction()}, Objective: &Objective{}},
	} {
		if err := validateEncounterShape(encounter); err == nil {
			t.Fatalf("expected invalid encounter to fail: %+v", encounter)
		}
	}

	actor, err := NewCharacter("actor", "Actor", Vectors{R: 3, G: 2, B: 1}, nil)
	if err != nil {
		t.Fatal(err)
	}
	target, err := NewCharacter("target", "Target", Vectors{R: 1, G: 2, B: 3}, nil)
	if err != nil {
		t.Fatal(err)
	}
	action := validCoverageAction()
	if _, _, err := resolveActionValues(map[string]*Character{"target": &target}, action); err == nil {
		t.Fatal("expected missing actor to fail")
	}
	if _, _, err := resolveActionValues(map[string]*Character{"actor": &actor}, action); err == nil {
		t.Fatal("expected missing target to fail")
	}
	action.PrimaryVector = "X"
	if _, _, err := resolveActionValues(map[string]*Character{"actor": &actor, "target": &target}, action); err == nil {
		t.Fatal("expected invalid action vector to fail")
	}
	action.PrimaryVector = VectorR
	action.Procedure = "unknown"
	if _, _, err := resolveActionValues(map[string]*Character{"actor": &actor, "target": &target}, action); err == nil {
		t.Fatal("expected invalid defense procedure to fail")
	}
}

func TestProgressionRequirementEdges(t *testing.T) {
	character, err := NewCharacter("hero", "Hero", Vectors{R: 1, G: 1, B: 1}, nil)
	if err != nil {
		t.Fatal(err)
	}
	invalidRequirement := validCoverageAbility()
	invalidRequirement.Requirements = map[Vector]int{"X": 1}
	if character.MeetsRequirements(invalidRequirement) {
		t.Fatal("invalid requirement vector should not be met")
	}
	discounted := validCoverageAbility()
	discounted.Requirements = map[Vector]int{VectorR: 1}
	character.Specializations = map[Vector]bool{VectorR: true}
	if !character.MeetsRequirements(discounted) {
		t.Fatal("specialization discount should satisfy one-point requirement")
	}
	if err := character.AdvanceSpecialization("X"); err == nil {
		t.Fatal("expected invalid specialization vector to fail")
	}

	invalidAbility := validCoverageAbility()
	invalidAbility.ID = ""
	if err := character.AdvanceNewAbility(invalidAbility); err == nil {
		t.Fatal("expected invalid advancement ability to fail")
	}
	known := validCoverageAbility()
	character.Abilities = []Ability{known}
	replacement := known
	replacement.Name = ""
	if err := character.ReplaceAbility(replacement); err == nil {
		t.Fatal("expected invalid replacement ability to fail")
	}
}
