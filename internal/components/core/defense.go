package core

import "fmt"

// Procedure names a specific defensive or offensive technique, per the
// defense procedure taxonomy in docs/core/en/combat/attack_and_defense.md
// (decision matrix D-004). Procedures are not collapsed into one generic
// defense formula: each declares its own vector owner.
type Procedure string

const (
	// ProcedureAttack applies direct R pressure against a target.
	ProcedureAttack Procedure = "attack"
	// ProcedureEvade uses G to avoid or alter contact.
	ProcedureEvade Procedure = "evade"
	// ProcedureReposition uses G to change range, cover, or engagement.
	ProcedureReposition Procedure = "reposition"
	// ProcedureBlock uses B, equipment, or an ability to intentionally receive pressure.
	ProcedureBlock Procedure = "block"
	// ProcedureSustain uses B to preserve a state or position under pressure.
	ProcedureSustain Procedure = "sustain"
	// ProcedureInterrupt uses R to stop an action by applying pressure first.
	ProcedureInterrupt Procedure = "interrupt"
	// ProcedureCounterpressure answers pressure with force after it lands.
	ProcedureCounterpressure Procedure = "counterpressure"
)

// Validate reports whether the procedure is a known Core V2 procedure.
func (procedure Procedure) Validate() error {
	switch procedure {
	case ProcedureAttack, ProcedureEvade, ProcedureReposition, ProcedureBlock,
		ProcedureSustain, ProcedureInterrupt, ProcedureCounterpressure:
		return nil
	default:
		return fmt.Errorf("unknown procedure %q", procedure)
	}
}

// Vector returns the RGB vector that owns this procedure.
func (procedure Procedure) Vector() Vector {
	switch procedure {
	case ProcedureEvade, ProcedureReposition:
		return VectorG
	case ProcedureBlock, ProcedureSustain:
		return VectorB
	case ProcedureAttack, ProcedureInterrupt, ProcedureCounterpressure:
		return VectorR
	default:
		return ""
	}
}

// DefenseValue returns the character's point value in the vector that owns
// the given procedure.
func DefenseValue(character Character, procedure Procedure) (int, error) {
	if err := procedure.Validate(); err != nil {
		return 0, err
	}
	vector := procedure.Vector()
	if vector == "" {
		return 0, fmt.Errorf("procedure %q has no vector owner", procedure)
	}
	return character.Vectors.Value(vector)
}
