package core

import "fmt"

// Stabilize applies the sustain procedure's recovery effect: it clears
// StateInjured and marks the character StateStabilized, closing the
// recovery closure evidence in decision matrix D-006.
func Stabilize(character *Character) error {
	if character == nil {
		return fmt.Errorf("character must be non-nil")
	}
	character.RemoveState(StateInjured)
	return character.AddState(StateStabilized)
}
