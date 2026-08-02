package fixtures

import (
	"os"
	"regexp"
	"strconv"
	"strings"
	"testing"

	"github.com/SergioLacerda/rpg-system-rgb/internal/components/core"
)

var yamlCombatAttackBlock = regexp.MustCompile(
	`(?s)- id:\s*(\S+)\s*` +
		`defender_health:\s*(\d+)\s*` +
		`defender_shield:\s*(\d+)\s*` +
		`defender_armor:\s*(\d+)\s*` +
		`weapon_damage:\s*(\d+)\s*` +
		`penetration:\s*(\d+)\s*` +
		`expected_effective_armor:\s*(\d+)\s*` +
		`expected_remaining_damage:\s*(\d+)\s*` +
		`expected_remaining_shield:\s*(\d+)\s*` +
		`expected_health_damage:\s*(\d+)\s*` +
		`expected_remaining_health:\s*(\d+)`,
)

// combatAttackFixture is one parsed attack entry from
// examples/combat-example.yaml.
type combatAttackFixture struct {
	id                      string
	defenderHealth          int
	defenderShield          int
	defenderArmor           int
	weaponDamage            int
	penetration             int
	expectedEffectiveArmor  int
	expectedRemainingDamage int
	expectedRemainingShield int
	expectedHealthDamage    int
	expectedRemainingHealth int
}

// parseCombatExampleFixture reads examples/combat-example.yaml and extracts
// each attack block. It is a line-oriented reader rather than a full YAML
// parser, matching the fixture file's regular, hand-authored structure —
// the same approach parity_test.go uses for characters/core-v2.yaml.
func parseCombatExampleFixture(t *testing.T) []combatAttackFixture {
	t.Helper()
	body, err := os.ReadFile("examples/combat-example.yaml")
	if err != nil {
		t.Fatal(err)
	}

	blocks := strings.Split(string(body), "\n  - id:")
	var attacks []combatAttackFixture
	for i, block := range blocks {
		if i == 0 {
			continue // header/comment content before the first attack entry
		}
		match := yamlCombatAttackBlock.FindStringSubmatch("  - id:" + block)
		if match == nil {
			t.Fatalf("could not parse combat attack block starting with: %.60s", block)
		}
		ints := make([]int, 10)
		for idx, field := range match[2:] {
			value, convErr := strconv.Atoi(field)
			if convErr != nil {
				t.Fatalf("invalid integer field in attack %q: %v", match[1], convErr)
			}
			ints[idx] = value
		}
		attacks = append(attacks, combatAttackFixture{
			id:                      match[1],
			defenderHealth:          ints[0],
			defenderShield:          ints[1],
			defenderArmor:           ints[2],
			weaponDamage:            ints[3],
			penetration:             ints[4],
			expectedEffectiveArmor:  ints[5],
			expectedRemainingDamage: ints[6],
			expectedRemainingShield: ints[7],
			expectedHealthDamage:    ints[8],
			expectedRemainingHealth: ints[9],
		})
	}
	return attacks
}

// mirrors: docs/core/en/reference/combat_example.md (and its PT-br mirror)
// TestCombatExampleFixtureMatchesEngine runs every attack in
// examples/combat-example.yaml through the real core.ApplyDamage pipeline
// and asserts the result matches the fixture's expected_* fields. This is
// the guardrail against a repeat of the published combat example bug
// (health increasing after damage, or shield not absorbing damage the
// engine says it should absorb) — see
// .analysis/refined/20260801-core-rules-executable-examples-turn-contract.
func TestCombatExampleFixtureMatchesEngine(t *testing.T) {
	for _, attack := range parseCombatExampleFixture(t) {
		t.Run(attack.id, func(t *testing.T) {
			defender, err := core.NewCharacter(attack.id+"-defender", "Defender", core.Vectors{R: 1, G: 1, B: 1}, nil)
			if err != nil {
				t.Fatal(err)
			}
			defender.Resources.MaxHealth = attack.defenderHealth
			defender.Resources.CurrentHealth = attack.defenderHealth
			defender.Resources.MaxShield = attack.defenderShield
			defender.Resources.CurrentShield = attack.defenderShield
			defender.Resources.Armor = attack.defenderArmor

			result, err := core.ApplyDamage(&defender, core.DamageInput{
				Impact:      attack.weaponDamage,
				Penetration: attack.penetration,
			})
			if err != nil {
				t.Fatal(err)
			}

			effectiveArmor := max(attack.defenderArmor-attack.penetration, 0)
			remainingDamage := attack.weaponDamage - result.ArmorReduced

			if effectiveArmor != attack.expectedEffectiveArmor {
				t.Errorf("effective armor got %d want %d", effectiveArmor, attack.expectedEffectiveArmor)
			}
			if remainingDamage != attack.expectedRemainingDamage {
				t.Errorf("remaining damage got %d want %d", remainingDamage, attack.expectedRemainingDamage)
			}
			if result.HealthDamage != attack.expectedHealthDamage {
				t.Errorf("health damage got %d want %d", result.HealthDamage, attack.expectedHealthDamage)
			}
			if defender.Resources.CurrentShield != attack.expectedRemainingShield {
				t.Errorf("remaining shield got %d want %d", defender.Resources.CurrentShield, attack.expectedRemainingShield)
			}
			if defender.Resources.CurrentHealth != attack.expectedRemainingHealth {
				t.Errorf("remaining health got %d want %d", defender.Resources.CurrentHealth, attack.expectedRemainingHealth)
			}
			if defender.Resources.CurrentHealth > attack.defenderHealth {
				t.Fatalf("health increased after damage: started %d, ended %d", attack.defenderHealth, defender.Resources.CurrentHealth)
			}
		})
	}
}
