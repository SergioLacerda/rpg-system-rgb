package core

import "fmt"

type Character struct {
	ID        string
	Name      string
	Vectors   Vectors
	Resources Resources
	States    map[State]bool
	Abilities []Ability
}

func NewCharacter(id, name string, vectors Vectors, abilities []Ability) (Character, error) {
	if id == "" {
		return Character{}, fmt.Errorf("character ID must be non-empty")
	}
	if name == "" {
		return Character{}, fmt.Errorf("character name must be non-empty")
	}
	if err := vectors.Validate(); err != nil {
		return Character{}, err
	}
	for _, ability := range abilities {
		if err := ability.Validate(); err != nil {
			return Character{}, fmt.Errorf("%s ability invalid: %w", id, err)
		}
	}
	character := Character{
		ID:        id,
		Name:      name,
		Vectors:   vectors,
		Resources: DeriveResources(vectors),
		States: map[State]bool{
			StateHealthy: true,
		},
		Abilities: abilities,
	}
	if character.Resources.CurrentShield > 0 {
		character.States[StateShielded] = true
	}
	return character, nil
}

func (character *Character) AddState(state State) error {
	if !validState(state) {
		return fmt.Errorf("unknown state %q", state)
	}
	if character.States == nil {
		character.States = map[State]bool{}
	}
	character.States[state] = true
	if state == StateDowned {
		delete(character.States, StateHealthy)
	}
	if state == StateInjured {
		delete(character.States, StateHealthy)
	}
	return nil
}

func (character *Character) RemoveState(state State) {
	delete(character.States, state)
}

func (character Character) HasState(state State) bool {
	return character.States[state]
}
