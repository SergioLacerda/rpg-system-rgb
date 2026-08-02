package fixtures

import (
	"os"
	"regexp"
	"strconv"
	"strings"
	"testing"
)

// lastNumberOnLine finds "<field> = ..." (possibly an arithmetic expression
// like "3 - 1 = 2") and returns the last integer on that line — the actual
// resulting value, regardless of whether the line shows its work.
func lastNumberOnLine(t *testing.T, text, field string) int {
	t.Helper()
	pattern := regexp.MustCompile(`(?m)^` + regexp.QuoteMeta(field) + ` = .*(-?\d+)\s*$`)
	match := pattern.FindStringSubmatch(text)
	if match == nil {
		t.Fatalf("field %q not found in document text", field)
	}
	value, err := strconv.Atoi(match[1])
	if err != nil {
		t.Fatalf("field %q has non-integer value %q: %v", field, match[1], err)
	}
	return value
}

// mirrors: docs/core/en/reference/combat_example.md (and its PT-br mirror)
// TestCombatExampleDocMatchesFixture parses the numeric facts out of the
// published combat example and cross-checks them against
// examples/combat-example.yaml — the same fixture
// TestCombatExampleFixtureMatchesEngine validates against core.ApplyDamage.
// This is the guardrail against the doc drifting back to hand-typed,
// unchecked numbers (the exact defect this fixture was introduced to fix):
// if someone edits the Markdown prose without updating the fixture (or vice
// versa), this test fails.
func TestCombatExampleDocMatchesFixture(t *testing.T) {
	fixtureAttacks := parseCombatExampleFixture(t)
	if len(fixtureAttacks) != 2 {
		t.Fatalf("expected 2 fixture attacks, got %d", len(fixtureAttacks))
	}
	attack1, attack2 := fixtureAttacks[0], fixtureAttacks[1]

	body, err := os.ReadFile("../../docs/core/en/reference/combat_example.md")
	if err != nil {
		t.Fatal(err)
	}
	doc := string(body)

	secondAttackIdx := strings.Index(doc, "## Second Attack")
	if secondAttackIdx == -1 {
		t.Fatal(`"## Second Attack" section not found in combat_example.md`)
	}
	firstAttackText := doc[:secondAttackIdx]
	secondAttackText := doc[secondAttackIdx:]

	checkAttack := func(t *testing.T, text string, want combatAttackFixture) {
		t.Helper()
		if got := lastNumberOnLine(t, text, "Effective Armor"); got != want.expectedEffectiveArmor {
			t.Errorf("Effective Armor: doc has %d, fixture wants %d", got, want.expectedEffectiveArmor)
		}
		if got := lastNumberOnLine(t, text, "Remaining Shield"); got != want.expectedRemainingShield {
			t.Errorf("Remaining Shield: doc has %d, fixture wants %d", got, want.expectedRemainingShield)
		}
		if got := lastNumberOnLine(t, text, "Health Damage"); got != want.expectedHealthDamage {
			t.Errorf("Health Damage: doc has %d, fixture wants %d", got, want.expectedHealthDamage)
		}
		got := lastNumberOnLine(t, text, "Remaining Health")
		if got != want.expectedRemainingHealth {
			t.Errorf("Remaining Health: doc has %d, fixture wants %d", got, want.expectedRemainingHealth)
		}
		if got > want.defenderHealth {
			t.Errorf("Remaining Health %d exceeds starting health %d — health must never increase from damage", got, want.defenderHealth)
		}
	}

	t.Run("first attack", func(t *testing.T) { checkAttack(t, firstAttackText, attack1) })
	t.Run("second attack", func(t *testing.T) { checkAttack(t, secondAttackText, attack2) })
}
