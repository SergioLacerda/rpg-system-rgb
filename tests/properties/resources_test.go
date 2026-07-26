package properties

import (
	"testing"

	"github.com/SergioLacerda/rpg-system-rgb/internal/components/core"
)

func TestResourceDerivationIsDeterministic(t *testing.T) {
	for r := 0; r <= 8; r++ {
		for b := 0; b <= 8; b++ {
			resources := core.DeriveResources(core.Vectors{R: r, G: 1, B: b})
			expectedHealth := 4 + r + b
			if resources.MaxHealth != expectedHealth || resources.CurrentHealth != expectedHealth {
				t.Fatalf("health derivation got max=%d current=%d want=%d", resources.MaxHealth, resources.CurrentHealth, expectedHealth)
			}
			if resources.MaxShield != b || resources.CurrentShield != b {
				t.Fatalf("shield derivation got max=%d current=%d want=%d", resources.MaxShield, resources.CurrentShield, b)
			}
		}
	}
}

func TestHealthDoesNotDependOnG(t *testing.T) {
	base := core.DeriveResources(core.Vectors{R: 2, G: 0, B: 3})
	for g := 1; g <= 8; g++ {
		resources := core.DeriveResources(core.Vectors{R: 2, G: g, B: 3})
		if resources.MaxHealth != base.MaxHealth {
			t.Fatalf("G changed health: g=%d got=%d want=%d", g, resources.MaxHealth, base.MaxHealth)
		}
	}
}
