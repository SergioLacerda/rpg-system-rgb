package core

import "fmt"

type Procedure string

const (
	ProcedureAttack          Procedure = "attack"
	ProcedureEvade           Procedure = "evade"
	ProcedureReposition      Procedure = "reposition"
	ProcedureBlock           Procedure = "block"
	ProcedureSustain         Procedure = "sustain"
	ProcedureInterrupt       Procedure = "interrupt"
	ProcedureCounterpressure Procedure = "counterpressure"
)

func (procedure Procedure) Validate() error {
	switch procedure {
	case ProcedureAttack, ProcedureEvade, ProcedureReposition, ProcedureBlock,
		ProcedureSustain, ProcedureInterrupt, ProcedureCounterpressure:
		return nil
	default:
		return fmt.Errorf("unknown procedure %q", procedure)
	}
}

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
