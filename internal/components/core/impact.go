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

// WeaponCategory identifies a firearm class with a canonical penetration
// value, per docs/core/en/weapons/mechanics/penetration.md. This table is
// narrower than, and not a 1:1 join with, the finer-grained weapon list in
// docs/core/en/weapons/categories/firearms.md's damage table.
type WeaponCategory string

const (
	// WeaponCategoryPistol is a lightweight sidearm.
	WeaponCategoryPistol WeaponCategory = "pistol"
	// WeaponCategorySMG is a high-rate-of-fire short-to-medium range weapon.
	WeaponCategorySMG WeaponCategory = "smg"
	// WeaponCategoryCarbine is a lightweight rifle.
	WeaponCategoryCarbine WeaponCategory = "carbine"
	// WeaponCategoryAssaultRifle balances damage and capacity.
	WeaponCategoryAssaultRifle WeaponCategory = "assault_rifle"
	// WeaponCategorySniperRifle is a high-damage precision weapon.
	WeaponCategorySniperRifle WeaponCategory = "sniper_rifle"
	// WeaponCategoryMachineGun provides sustained suppressive fire.
	WeaponCategoryMachineGun WeaponCategory = "machine_gun"
)

// DefaultPenetration returns the category's canonical penetration value,
// per the Weapon Penetration Table in
// docs/core/en/weapons/mechanics/penetration.md.
func (category WeaponCategory) DefaultPenetration() (int, error) {
	switch category {
	case WeaponCategoryPistol:
		return 1, nil
	case WeaponCategorySMG:
		return 2, nil
	case WeaponCategoryCarbine:
		return 2, nil
	case WeaponCategoryAssaultRifle:
		return 3, nil
	case WeaponCategorySniperRifle:
		return 4, nil
	case WeaponCategoryMachineGun:
		return 3, nil
	default:
		return 0, fmt.Errorf("unknown weapon category %q", category)
	}
}
