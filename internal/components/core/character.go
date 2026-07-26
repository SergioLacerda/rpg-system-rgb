package core

import "fmt"

// Character is a playable or non-playable RGB Core V2 entity: vectors,
// derived resources, active states, and declared abilities.
type Character struct {
	ID        string
	Name      string
	Vectors   Vectors
	Resources Resources
	States    map[State]bool
	Abilities []Ability
}

// NewCharacter builds a Character with resources derived from vectors and
// validates its abilities. It does not enforce the starting point budget
// (decision matrix D-002 closure evidence); callers that need a starting
// player character should validate the budget separately.
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

// AddState marks the character as holding the given state. Adding
// StateDowned or StateInjured clears StateHealthy.
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

// RemoveState clears the given state from the character.
func (character *Character) RemoveState(state State) {
	delete(character.States, state)
}

// HasState reports whether the character currently holds the given state.
func (character Character) HasState(state State) bool {
	return character.States[state]
}
