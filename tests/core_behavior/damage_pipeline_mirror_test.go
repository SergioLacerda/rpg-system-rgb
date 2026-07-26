package corebehavior

import (
	"testing"

	"github.com/SergioLacerda/rpg-system-rgb/internal/components/core"
)

// mirrors: tests/features/damage/armor_shield_absorption.feature#Shield stays depleted across consecutive hits
func TestShieldStaysDepletedAcrossConsecutiveHitsFeatureExample(t *testing.T) {
	target, err := core.NewCharacter("target", "Target", core.Vectors{R: 2, G: 2, B: 3}, nil)
	if err != nil {
		t.Fatal(err)
	}
	target.Resources.CurrentShield = 0

	result, err := core.ApplyDamage(&target, core.DamageInput{Impact: 3})
	if err != nil {
		t.Fatal(err)
	}
	if result.ShieldAbsorbed != 0 {
		t.Fatalf("shield absorbed got %d want 0", result.ShieldAbsorbed)
	}
	if result.HealthDamage != 3 {
		t.Fatalf("health damage got %d want 3", result.HealthDamage)
	}
	if target.HasState(core.StateShielded) {
		t.Fatal("target must not have the shielded state once shield reserve is depleted")
	}
}

// mirrors: tests/features/damage/armor_shield_absorption.feature#Excess penetration does not increase damage
func TestExcessPenetrationDoesNotIncreaseDamageFeatureExample(t *testing.T) {
	target, err := core.NewCharacter("target", "Target", core.Vectors{R: 2, G: 2, B: 0}, nil)
	if err != nil {
		t.Fatal(err)
	}
	target.Resources.Armor = 2
	target.Resources.CurrentShield = 0

	result, err := core.ApplyDamage(&target, core.DamageInput{Impact: 4, Penetration: 6})
	if err != nil {
		t.Fatal(err)
	}
	if result.ArmorReduced != 0 {
		t.Fatalf("armor reduced got %d want 0", result.ArmorReduced)
	}
	if result.HealthDamage != 4 {
		t.Fatalf("health damage got %d want 4", result.HealthDamage)
	}
}

// mirrors: tests/features/damage/armor_shield_absorption.feature#Melee impact combines weapon and R
func TestMeleeImpactCombinesWeaponAndRFeatureExample(t *testing.T) {
	impact, err := core.ComputeImpact(core.ImpactSourceMelee, 2, 3)
	if err != nil {
		t.Fatal(err)
	}
	if impact != 5 {
		t.Fatalf("melee impact got %d want 5", impact)
	}
}

// mirrors: tests/features/damage/armor_shield_absorption.feature#Firearm impact uses the weapon value only
func TestFirearmImpactUsesWeaponValueOnlyFeatureExample(t *testing.T) {
	impact, err := core.ComputeImpact(core.ImpactSourceFirearm, 4, 3)
	if err != nil {
		t.Fatal(err)
	}
	if impact != 4 {
		t.Fatalf("firearm impact got %d want 4", impact)
	}
}

// mirrors: tests/features/damage/armor_shield_absorption.feature#Weapon category penetration defaults
func TestWeaponCategoryPenetrationDefaultsFeatureExamples(t *testing.T) {
	cases := []struct {
		category    core.WeaponCategory
		penetration int
	}{
		{category: core.WeaponCategoryPistol, penetration: 1},
		{category: core.WeaponCategorySMG, penetration: 2},
		{category: core.WeaponCategoryCarbine, penetration: 2},
		{category: core.WeaponCategoryAssaultRifle, penetration: 3},
		{category: core.WeaponCategorySniperRifle, penetration: 4},
		{category: core.WeaponCategoryMachineGun, penetration: 3},
	}
	for _, tc := range cases {
		got, err := tc.category.DefaultPenetration()
		if err != nil {
			t.Fatal(err)
		}
		if got != tc.penetration {
			t.Fatalf("category=%s penetration got %d want %d", tc.category, got, tc.penetration)
		}
	}
}
