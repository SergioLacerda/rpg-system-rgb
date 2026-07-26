package properties

import (
	"testing"

	"github.com/SergioLacerda/rpg-system-rgb/internal/components/core"
)

func TestIncreasingGDoesNotReduceEvadeValue(t *testing.T) {
	previous := -1
	for g := 0; g <= 10; g++ {
		character, err := core.NewCharacter("runner", "Runner", core.Vectors{R: 1, G: g, B: 1}, nil)
		if err != nil {
			t.Fatal(err)
		}
		value, err := core.DefenseValue(character, core.ProcedureEvade)
		if err != nil {
			t.Fatal(err)
		}
		if value < previous {
			t.Fatalf("evade value decreased when G increased: g=%d value=%d previous=%d", g, value, previous)
		}
		previous = value
	}
}

func TestProcedureOwnersStayDistinct(t *testing.T) {
	if core.ProcedureEvade.Vector() != core.VectorG {
		t.Fatal("evade must be owned by G")
	}
	if core.ProcedureBlock.Vector() != core.VectorB {
		t.Fatal("block must be owned by B")
	}
	if core.ProcedureInterrupt.Vector() != core.VectorR {
		t.Fatal("interrupt must be owned by R")
	}
}
