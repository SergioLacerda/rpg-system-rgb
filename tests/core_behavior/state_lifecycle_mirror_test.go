package corebehavior

import (
	"testing"

	"github.com/SergioLacerda/rpg-system-rgb/internal/components/core"
)

// mirrors: tests/features/core/state_lifecycle.feature#No orphan states exist in the core model
func TestNoOrphanStatesExistInTheCoreModelFeatureExample(t *testing.T) {
	for _, state := range core.AllStates() {
		lifecycle, ok := core.Lifecycle(state)
		if !ok {
			t.Fatalf("state %s has no declared lifecycle entry", state)
		}
		if lifecycle.ProducedBy == "" {
			t.Fatalf("state %s does not declare a producing source", state)
		}
		if lifecycle.ClearedBy == "" {
			t.Fatalf("state %s does not declare how it is cleared", state)
		}
	}
}

// mirrors: tests/features/core/state_lifecycle.feature#Stabilizing an injured character
func TestStabilizingAnInjuredCharacterFeatureExample(t *testing.T) {
	character, err := core.NewCharacter("wounded", "Wounded", core.Vectors{R: 2, G: 2, B: 2}, nil)
	if err != nil {
		t.Fatal(err)
	}
	character.Resources.CurrentHealth = 5
	if err := character.AddState(core.StateInjured); err != nil {
		t.Fatal(err)
	}

	if err := core.Stabilize(&character); err != nil {
		t.Fatal(err)
	}
	if !character.HasState(core.StateStabilized) {
		t.Fatal("character must become stabilized")
	}
	if character.HasState(core.StateInjured) {
		t.Fatal("character must no longer be injured")
	}
}
