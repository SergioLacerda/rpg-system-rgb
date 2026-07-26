package corebehavior

import (
	"testing"

	"github.com/SergioLacerda/rpg-system-rgb/internal/components/core"
)

// mirrors: tests/features/combat/movement.feature#Ground movement per turn
func TestGroundMovementPerTurnFeatureExamples(t *testing.T) {
	cases := []struct {
		g, meters int
	}{
		{g: 0, meters: 0},
		{g: 3, meters: 6},
		{g: 5, meters: 10},
	}
	for _, tc := range cases {
		if got := core.MovementDistance(tc.g); got != tc.meters {
			t.Fatalf("G=%d movement got %d want %d", tc.g, got, tc.meters)
		}
	}
}

// mirrors: tests/features/combat/movement.feature#Aerial movement uses the same formula
func TestAerialMovementUsesSameFormulaFeatureExample(t *testing.T) {
	if got := core.MovementDistance(3); got != 6 {
		t.Fatalf("aerial movement got %d want 6", got)
	}
}
