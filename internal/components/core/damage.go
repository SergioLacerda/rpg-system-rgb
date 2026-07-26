package core

import "fmt"

type DamageInput struct {
	Impact      int
	Penetration int
}

type DamageStep struct {
	Name  string
	Value int
}

type DamageResult struct {
	Steps          []DamageStep
	ArmorReduced   int
	ShieldAbsorbed int
	HealthDamage   int
}

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

func min(left, right int) int {
	if left < right {
		return left
	}
	return right
}

func max(left, right int) int {
	if left > right {
		return left
	}
	return right
}
