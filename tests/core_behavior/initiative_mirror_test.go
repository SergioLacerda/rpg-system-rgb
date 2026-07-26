package corebehavior

import (
	"testing"

	"github.com/SergioLacerda/rpg-system-rgb/internal/components/core"
)

// mirrors: tests/features/encounters/initiative.feature#The faster character opens the round
func TestFasterCharacterOpensTheRoundFeatureExample(t *testing.T) {
	order := core.InitiativeOrder([]core.InitiativeEntry{
		{ActorID: "vanguard", G: 2},
		{ActorID: "runner", G: 5},
	})
	if len(order) != 2 || order[0] != "runner" || order[1] != "vanguard" {
		t.Fatalf("order got %v want [runner vanguard]", order)
	}
}

// mirrors: tests/features/encounters/initiative.feature#A surprise attack ignores initiative
func TestSurpriseAttackIgnoresInitiativeFeatureExample(t *testing.T) {
	order := core.InitiativeOrder([]core.InitiativeEntry{
		{ActorID: "vanguard", G: 2, Surprise: true},
		{ActorID: "runner", G: 5},
	})
	if len(order) != 2 || order[0] != "vanguard" || order[1] != "runner" {
		t.Fatalf("order got %v want [vanguard runner]", order)
	}
}
