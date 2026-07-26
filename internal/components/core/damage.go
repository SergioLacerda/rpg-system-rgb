package core

import "fmt"

// DamageInput is the raw impact and penetration entering the damage
// pipeline, per docs/core/en/combat/damage_model.md.
type DamageInput struct {
	Impact      int
	Penetration int
}

// DamageStep records one named stage of the damage pipeline, in resolution order.
type DamageStep struct {
	Name  string
	Value int
}

// DamageResult is the outcome of running the damage pipeline: the ordered
// steps and the totals applied at each mitigation layer.
type DamageResult struct {
	Steps          []DamageStep
	ArmorReduced   int
	ShieldAbsorbed int
	HealthDamage   int
}

// ApplyDamage runs the canonical damage pipeline against target: penetration
// reduces armor, armor reduces impact, shield absorbs what armor did not,
// and the remainder becomes health damage and a state consequence.
func ApplyDamage(target *Character, input DamageInput) (DamageResult, error) {
	if target == nil {
		return DamageResult{}, fmt.Errorf("target must be non-nil")
	}
	if input.Impact < 0 {
		return DamageResult{}, fmt.Errorf("impact must be non-negative")
	}
	if input.Penetration < 0 {
		return DamageResult{}, fmt.Errorf("penetration must be non-negative")
	}

	effectiveArmor := max(target.Resources.Armor-input.Penetration, 0)
	armorReduced := min(effectiveArmor, input.Impact)
	afterArmor := input.Impact - armorReduced
	shieldAbsorbed := min(target.Resources.CurrentShield, afterArmor)
	healthDamage := afterArmor - shieldAbsorbed

	target.Resources.CurrentShield -= shieldAbsorbed
	target.Resources.CurrentHealth = max(target.Resources.CurrentHealth-healthDamage, 0)
	if target.Resources.CurrentShield == 0 {
		target.RemoveState(StateShielded)
	}
	if target.Resources.CurrentHealth == 0 {
		_ = target.AddState(StateDowned)
	} else if healthDamage > 0 {
		_ = target.AddState(StateInjured)
	}

	return DamageResult{
		Steps: []DamageStep{
			{Name: "impact", Value: input.Impact},
			{Name: "penetration", Value: input.Penetration},
			{Name: "armor_reduction", Value: armorReduced},
			{Name: "shield_absorption", Value: shieldAbsorbed},
			{Name: "health_consequence", Value: healthDamage},
		},
		ArmorReduced:   armorReduced,
		ShieldAbsorbed: shieldAbsorbed,
		HealthDamage:   healthDamage,
	}, nil
}
