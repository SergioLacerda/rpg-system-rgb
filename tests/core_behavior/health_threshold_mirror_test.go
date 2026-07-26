package corebehavior

import (
	"testing"

	"github.com/SergioLacerda/rpg-system-rgb/internal/components/core"
)

// mirrors: tests/features/damage/armor_shield_absorption.feature#Damage equal to remaining health downs the target
func TestDamageEqualToRemainingHealthDownsTheTargetFeatureExample(t *testing.T) {
	target, err := core.NewCharacter("target", "Target", core.Vectors{R: 0, G: 0, B: 0}, nil)
	if err != nil {
		t.Fatal(err)
	}
	target.Resources.CurrentHealth = 3
	target.Resources.CurrentShield = 0

	result, err := core.ApplyDamage(&target, core.DamageInput{Impact: 3})
	if err != nil {
		t.Fatal(err)
	}
	if result.HealthDamage != 3 {
		t.Fatalf("health damage got %d want 3", result.HealthDamage)
	}
	if !target.HasState(core.StateDowned) {
		t.Fatal("target must become downed")
	}
	if target.HasState(core.StateHealthy) {
		t.Fatal("target must not remain healthy")
	}
}

// mirrors: tests/features/damage/armor_shield_absorption.feature#Fully mitigated damage does not injure
func TestFullyMitigatedDamageDoesNotInjureFeatureExample(t *testing.T) {
	target, err := core.NewCharacter("target", "Target", core.Vectors{R: 0, G: 0, B: 0}, nil)
	if err != nil {
		t.Fatal(err)
	}
	target.Resources.Armor = 5
	target.Resources.CurrentShield = 0

	result, err := core.ApplyDamage(&target, core.DamageInput{Impact: 4, Penetration: 0})
	if err != nil {
		t.Fatal(err)
	}
	if result.HealthDamage != 0 {
		t.Fatalf("health damage got %d want 0", result.HealthDamage)
	}
	if !target.HasState(core.StateHealthy) {
		t.Fatal("target must remain healthy")
	}
}
