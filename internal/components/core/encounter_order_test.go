package core_test

import (
	"testing"

	"github.com/SergioLacerda/rpg-system-rgb/internal/components/core"
)

// TestRunEncounterResolvesActionsInDeclaredOrder locks the RunEncounter
// contract stated in its own doc comment: actions resolve "in list order,"
// so result.ActionResults must line up with encounter.Actions by
// declaration position. The two actions below use different actors,
// targets, and procedures so their margins are distinguishable — a
// processing-order defect (e.g. resolving actions in reverse) would put
// the wrong margin at each index.
func TestRunEncounterResolvesActionsInDeclaredOrder(t *testing.T) {
	alpha, err := core.NewCharacter("alpha", "Alpha", core.Vectors{R: 5, G: 1, B: 1}, nil)
	if err != nil {
		t.Fatalf("building alpha: %v", err)
	}
	bravo, err := core.NewCharacter("bravo", "Bravo", core.Vectors{R: 1, G: 4, B: 1}, nil)
	if err != nil {
		t.Fatalf("building bravo: %v", err)
	}
	characters := map[string]*core.Character{"alpha": &alpha, "bravo": &bravo}

	encounter := core.Encounter{
		ID:   "order-check",
		Name: "Declared Order Check",
		Actions: []core.Action{
			{
				Actor: "alpha", Target: "bravo",
				PrimaryVector: core.VectorR, Intent: "press", Procedure: core.ProcedureAttack,
				Consequence: "alpha presses bravo", Round: 1,
			},
			{
				Actor: "bravo", Target: "alpha",
				PrimaryVector: core.VectorG, Intent: "reposition", Procedure: core.ProcedureEvade,
				Consequence: "bravo repositions on alpha", Round: 1,
			},
		},
	}

	result, err := core.RunEncounter(characters, encounter)
	if err != nil {
		t.Fatalf("RunEncounter failed: %v", err)
	}
	if len(result.ActionResults) != 2 {
		t.Fatalf("expected 2 action results, got %d", len(result.ActionResults))
	}

	wantFirstMargin := alpha.Vectors.R - bravo.Vectors.R  // alpha's attack (R) vs bravo's R defense
	wantSecondMargin := bravo.Vectors.G - alpha.Vectors.G // bravo's reposition (G) vs alpha's G defense

	if got := result.ActionResults[0].Margin; got != wantFirstMargin {
		t.Fatalf("ActionResults[0] (alpha's declared-first attack) margin = %d, want %d", got, wantFirstMargin)
	}
	if got := result.ActionResults[1].Margin; got != wantSecondMargin {
		t.Fatalf("ActionResults[1] (bravo's declared-second reposition) margin = %d, want %d", got, wantSecondMargin)
	}
}
