package core

import "fmt"

// ImpactSource identifies where an attack's impact value comes from, per
// docs/core/en/combat/damage_model.md.
type ImpactSource string

const (
	// ImpactSourceMelee impact combines the weapon value with the attacker's R.
	ImpactSourceMelee ImpactSource = "melee"
	// ImpactSourceFirearm impact uses the weapon value only.
	ImpactSourceFirearm ImpactSource = "firearm"
	// ImpactSourceExplosive impact uses the explosive profile value only.
	ImpactSourceExplosive ImpactSource = "explosive"
	// ImpactSourceAbility impact uses the ability's declared effect value only.
	ImpactSourceAbility ImpactSource = "ability"
)

// ComputeImpact returns the impact value entering the damage pipeline
// (before penetration, armor, and shield), per the default impact source
// rules in docs/core/en/combat/damage_model.md. weaponValue is the weapon,
// explosive profile, or ability effect value; attackerR is only added for
// melee impact.
func ComputeImpact(source ImpactSource, weaponValue, attackerR int) (int, error) {
	switch source {
	case ImpactSourceMelee:
		return weaponValue + attackerR, nil
	case ImpactSourceFirearm, ImpactSourceExplosive, ImpactSourceAbility:
		return weaponValue, nil
	default:
		return 0, fmt.Errorf("unknown impact source %q", source)
	}
}
