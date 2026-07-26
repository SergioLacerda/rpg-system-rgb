package corebehavior

import (
	"testing"

	"github.com/SergioLacerda/rpg-system-rgb/internal/components/core"
)

// mirrors: tests/features/combat/grapple.feature#Grapple resolution
func TestGrappleResolutionFeatureExamples(t *testing.T) {
	cases := []struct {
		attackerG, attackerR, targetG, targetR int
		succeeds                               bool
	}{
		{attackerG: 3, attackerR: 3, targetG: 3, targetR: 3, succeeds: true},
		{attackerG: 4, attackerR: 2, targetG: 3, targetR: 3, succeeds: false},
		{attackerG: 2, attackerR: 4, targetG: 3, targetR: 3, succeeds: false},
		{attackerG: 4, attackerR: 4, targetG: 3, targetR: 3, succeeds: true},
	}
	for _, tc := range cases {
		got := core.ResolveGrapple(tc.attackerG, tc.attackerR, tc.targetG, tc.targetR)
		if got != tc.succeeds {
			t.Fatalf("attacker G%d/R%d vs target G%d/R%d: got %v want %v",
				tc.attackerG, tc.attackerR, tc.targetG, tc.targetR, got, tc.succeeds)
		}
	}
}
