package properties

import (
	"testing"

	"github.com/SergioLacerda/rpg-system-rgb/internal/components/core"
)

// mirrors: tests/features/combat/attack_evasion.feature#R greater than or equal to G always yields a successful outcome
func TestTeachingShorthandStaysConsistentWithMargins(t *testing.T) {
	for r := 0; r <= 10; r++ {
		for g := 0; g <= 10; g++ {
			if r < g {
				continue
			}
			resolution := core.Resolve(r, 0, g)
			if !resolution.Outcome.Successful() {
				t.Fatalf("R=%d G=%d margin=%d outcome=%s must be successful when R >= G", r, g, resolution.Margin, resolution.Outcome)
			}
		}
	}
}
