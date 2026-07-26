package properties

import (
	"testing"

	"github.com/SergioLacerda/rpg-system-rgb/internal/components/core"
)

func TestDamageFinalValueIsNeverNegative(t *testing.T) {
	for impact := 0; impact <= 12; impact++ {
		for penetration := 0; penetration <= 6; penetration++ {
			for armor := 0; armor <= 6; armor++ {
				for shield := 0; shield <= 6; shield++ {
					target, err := core.NewCharacter("target", "Target", core.Vectors{R: 2, G: 2, B: shield}, nil)
					if err != nil {
						t.Fatal(err)
					}
					target.Resources.Armor = armor
					target.Resources.CurrentShield = shield

					result, err := core.ApplyDamage(&target, core.DamageInput{Impact: impact, Penetration: penetration})
					if err != nil {
						t.Fatal(err)
					}
					if result.HealthDamage < 0 {
						t.Fatalf("negative health damage for impact=%d penetration=%d armor=%d shield=%d", impact, penetration, armor, shield)
					}
					if target.Resources.CurrentShield < 0 {
						t.Fatalf("negative shield for impact=%d penetration=%d armor=%d shield=%d", impact, penetration, armor, shield)
					}
					if target.Resources.CurrentHealth < 0 {
						t.Fatalf("negative health for impact=%d penetration=%d armor=%d shield=%d", impact, penetration, armor, shield)
					}
				}
			}
		}
	}
}

func TestArmorNeverAbsorbsMoreThanIncomingImpact(t *testing.T) {
	for impact := 0; impact <= 12; impact++ {
		for penetration := 0; penetration <= 6; penetration++ {
			for armor := 0; armor <= 12; armor++ {
				target, err := core.NewCharacter("target", "Target", core.Vectors{R: 2, G: 2, B: 0}, nil)
				if err != nil {
					t.Fatal(err)
				}
				target.Resources.Armor = armor

				result, err := core.ApplyDamage(&target, core.DamageInput{Impact: impact, Penetration: penetration})
				if err != nil {
					t.Fatal(err)
				}
				if result.ArmorReduced > impact {
					t.Fatalf("armor reduced %d exceeds impact %d", result.ArmorReduced, impact)
				}
			}
		}
	}
}

func TestMoreArmorCannotIncreaseHealthDamage(t *testing.T) {
	for impact := 0; impact <= 12; impact++ {
		for penetration := 0; penetration <= 6; penetration++ {
			previousDamage := impact + 1
			for armor := 0; armor <= 12; armor++ {
				target, err := core.NewCharacter("target", "Target", core.Vectors{R: 2, G: 2, B: 0}, nil)
				if err != nil {
					t.Fatal(err)
				}
				target.Resources.Armor = armor

				result, err := core.ApplyDamage(&target, core.DamageInput{Impact: impact, Penetration: penetration})
				if err != nil {
					t.Fatal(err)
				}
				if result.HealthDamage > previousDamage {
					t.Fatalf("health damage increased with armor: impact=%d penetration=%d armor=%d damage=%d previous=%d", impact, penetration, armor, result.HealthDamage, previousDamage)
				}
				previousDamage = result.HealthDamage
			}
		}
	}
}
