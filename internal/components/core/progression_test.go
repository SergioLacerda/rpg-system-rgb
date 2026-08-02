package core

import "testing"

func TestNewCharacterStartsAtLevelOne(t *testing.T) {
	character, err := NewCharacter("hero", "Hero", Vectors{R: 3, G: 2, B: 2}, nil)
	if err != nil {
		t.Fatal(err)
	}
	if character.Level != 1 {
		t.Fatalf("new character must start at level 1, got %d", character.Level)
	}
}

func TestAdvanceVectorGrowthDistributesBudgetAndAdvancesLevel(t *testing.T) {
	character, err := NewCharacter("hero", "Hero", Vectors{R: 3, G: 2, B: 2}, nil)
	if err != nil {
		t.Fatal(err)
	}
	if err := character.AdvanceVectorGrowth(Vectors{R: 1, B: 1}); err != nil {
		t.Fatal(err)
	}
	if character.Vectors != (Vectors{R: 4, G: 2, B: 3}) {
		t.Fatalf("unexpected vectors after growth: %+v", character.Vectors)
	}
	if character.Level != 2 {
		t.Fatalf("expected level 2 after one advancement, got %d", character.Level)
	}
	if len(character.AdvancementHistory) != 1 || character.AdvancementHistory[0].Choice != AdvancementVectorGrowth {
		t.Fatalf("unexpected advancement history: %+v", character.AdvancementHistory)
	}
}

func TestAdvanceVectorGrowthRejectsWrongBudget(t *testing.T) {
	character, err := NewCharacter("hero", "Hero", Vectors{R: 3, G: 2, B: 2}, nil)
	if err != nil {
		t.Fatal(err)
	}
	if err := character.AdvanceVectorGrowth(Vectors{R: 3}); err == nil {
		t.Fatal("expected error when growth does not sum to AdvancementBudget")
	}
	if err := character.AdvanceVectorGrowth(Vectors{R: -1, G: 3}); err == nil {
		t.Fatal("expected error when growth is negative")
	}
}

func newAbility(id string, vector Vector, tier int, requirements map[Vector]int, timing Timing) Ability {
	return Ability{
		ID:           id,
		Name:         id,
		Vector:       vector,
		Tier:         tier,
		Requirements: requirements,
		Timing:       timing,
		Cost:         1,
		Range:        1,
		Duration:     "instant",
		Effects:      []string{"test effect"},
		Tags:         []string{"test"},
	}
}

func TestAdvanceNewAbilityGrantsWhenRequirementsMet(t *testing.T) {
	character, err := NewCharacter("hero", "Hero", Vectors{R: 3, G: 2, B: 2}, nil)
	if err != nil {
		t.Fatal(err)
	}
	ability := newAbility("strike", VectorR, 2, map[Vector]int{VectorR: 3}, TimingAction)
	if err := character.AdvanceNewAbility(ability); err != nil {
		t.Fatal(err)
	}
	if len(character.Abilities) != 1 || character.Abilities[0].ID != "strike" {
		t.Fatalf("expected granted ability on character, got %+v", character.Abilities)
	}
	if character.AdvancementHistory[0].Choice != AdvancementNewAbility {
		t.Fatalf("expected new_ability advancement recorded, got %+v", character.AdvancementHistory)
	}
}

func TestAdvanceNewAbilityRejectsWhenRequirementsNotMet(t *testing.T) {
	character, err := NewCharacter("hero", "Hero", Vectors{R: 1, G: 1, B: 1}, nil)
	if err != nil {
		t.Fatal(err)
	}
	ability := newAbility("strike", VectorR, 2, map[Vector]int{VectorR: 5}, TimingAction)
	if err := character.AdvanceNewAbility(ability); err == nil {
		t.Fatal("expected error when character does not meet ability requirements")
	}
}

func TestAdvanceNewReactionRequiresReactionTiming(t *testing.T) {
	character, err := NewCharacter("hero", "Hero", Vectors{R: 3, G: 2, B: 2}, nil)
	if err != nil {
		t.Fatal(err)
	}
	actionAbility := newAbility("dodge", VectorG, 1, nil, TimingAction)
	if err := character.AdvanceNewReaction(actionAbility); err == nil {
		t.Fatal("expected error when granting a non-reaction ability as new-reaction advancement")
	}
	reactionAbility := newAbility("counter", VectorG, 1, nil, TimingReaction)
	if err := character.AdvanceNewReaction(reactionAbility); err != nil {
		t.Fatal(err)
	}
	if character.AdvancementHistory[0].Choice != AdvancementNewReaction {
		t.Fatalf("expected new_reaction advancement recorded, got %+v", character.AdvancementHistory)
	}
}

func TestAdvanceNewStateOrManeuverAccessGrantsAbility(t *testing.T) {
	character, err := NewCharacter("hero", "Hero", Vectors{R: 3, G: 2, B: 2}, nil)
	if err != nil {
		t.Fatal(err)
	}
	ability := newAbility("breach", VectorR, 1, nil, TimingAction)
	if err := character.AdvanceNewStateOrManeuverAccess(ability); err != nil {
		t.Fatal(err)
	}
	if len(character.Abilities) != 1 {
		t.Fatalf("expected ability granted, got %+v", character.Abilities)
	}
	if character.AdvancementHistory[0].Choice != AdvancementNewStateAccess {
		t.Fatalf("expected new_state_or_maneuver_access advancement recorded, got %+v", character.AdvancementHistory)
	}
}

func TestAdvanceSpecializationIgnoresOneTierOfRequirementsForMatchingVector(t *testing.T) {
	character, err := NewCharacter("hero", "Hero", Vectors{R: 2, G: 1, B: 1}, nil)
	if err != nil {
		t.Fatal(err)
	}
	ability := newAbility("advanced-strike", VectorR, 2, map[Vector]int{VectorR: 3}, TimingAction)

	if err := character.AdvanceNewAbility(ability); err == nil {
		t.Fatal("expected requirements to block the ability before specialization")
	}

	if err := character.AdvanceSpecialization(VectorR); err != nil {
		t.Fatal(err)
	}
	if !character.Specializations[VectorR] {
		t.Fatal("expected character to be specialized in R")
	}

	if err := character.AdvanceNewAbility(ability); err != nil {
		t.Fatalf("expected specialization discount to satisfy requirements, got error: %v", err)
	}
}

func TestAdvanceSpecializationRejectsDuplicate(t *testing.T) {
	character, err := NewCharacter("hero", "Hero", Vectors{R: 3, G: 2, B: 2}, nil)
	if err != nil {
		t.Fatal(err)
	}
	if err := character.AdvanceSpecialization(VectorR); err != nil {
		t.Fatal(err)
	}
	if err := character.AdvanceSpecialization(VectorR); err == nil {
		t.Fatal("expected error when specializing in the same vector twice")
	}
}

func TestReplaceAbilityIsFreeFormPerU01(t *testing.T) {
	original := newAbility("guard", VectorB, 1, nil, TimingAction)
	character, err := NewCharacter("hero", "Hero", Vectors{R: 2, G: 2, B: 3}, []Ability{original})
	if err != nil {
		t.Fatal(err)
	}
	improved := original
	improved.Cost = 0
	improved.Duration = "1 scene"
	improved.Limits = []string{"once per scene", "no longer requires setup"}
	if err := character.ReplaceAbility(improved); err != nil {
		t.Fatal(err)
	}
	if character.Abilities[0].Cost != 0 || character.Abilities[0].Duration != "1 scene" {
		t.Fatalf("expected ability replaced with improved version, got %+v", character.Abilities[0])
	}
	if character.AdvancementHistory[0].Choice != AdvancementAbilityImprovement {
		t.Fatalf("expected ability_improvement advancement recorded, got %+v", character.AdvancementHistory)
	}
}

func TestReplaceAbilityRejectsUnknownID(t *testing.T) {
	character, err := NewCharacter("hero", "Hero", Vectors{R: 2, G: 2, B: 3}, nil)
	if err != nil {
		t.Fatal(err)
	}
	unknown := newAbility("ghost", VectorB, 1, nil, TimingAction)
	if err := character.ReplaceAbility(unknown); err == nil {
		t.Fatal("expected error when improving an ability the character does not have")
	}
}

func TestGrantResourceRequiresValidDefinition(t *testing.T) {
	character, err := NewCharacter("hero", "Hero", Vectors{R: 2, G: 2, B: 3}, nil)
	if err != nil {
		t.Fatal(err)
	}
	if err := character.GrantResource(ResourceDefinition{Name: "", Max: 3}); err == nil {
		t.Fatal("expected error for empty resource name")
	}
	if err := character.GrantResource(ResourceDefinition{Name: "reserve-charges", Max: 0}); err == nil {
		t.Fatal("expected error for non-positive resource max")
	}
}

func TestGrantResourceAddsCustomResourceAtFullValue(t *testing.T) {
	character, err := NewCharacter("hero", "Hero", Vectors{R: 2, G: 2, B: 3}, nil)
	if err != nil {
		t.Fatal(err)
	}
	definition := ResourceDefinition{Name: "reserve-charges", Max: 3}
	if err := character.GrantResource(definition); err != nil {
		t.Fatal(err)
	}
	if character.CustomResources["reserve-charges"] != 3 {
		t.Fatalf("expected granted resource at full value, got %+v", character.CustomResources)
	}
	if character.AdvancementHistory[0].Choice != AdvancementNewResource {
		t.Fatalf("expected new_resource advancement recorded, got %+v", character.AdvancementHistory)
	}
}

func TestAdvancementChoicesDoNotStackWithinLevelBoundaries(t *testing.T) {
	character, err := NewCharacter("hero", "Hero", Vectors{R: 3, G: 2, B: 2}, nil)
	if err != nil {
		t.Fatal(err)
	}
	if err := character.AdvanceVectorGrowth(Vectors{R: 2}); err != nil {
		t.Fatal(err)
	}
	if err := character.AdvanceSpecialization(VectorG); err != nil {
		t.Fatal(err)
	}
	if character.Level != 3 {
		t.Fatalf("expected two advancements to reach level 3, got %d", character.Level)
	}
	if len(character.AdvancementHistory) != 2 {
		t.Fatalf("expected two distinct advancement records, got %+v", character.AdvancementHistory)
	}
	if character.AdvancementHistory[0].Level == character.AdvancementHistory[1].Level {
		t.Fatalf("expected each advancement to land on its own level, got %+v", character.AdvancementHistory)
	}
}
