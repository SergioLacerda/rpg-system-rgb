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

// StartingPointBudget is the canonical number of points a starting player
// character distributes among R, G, and B, per
// docs/core/en/core/character_creation.md.
const StartingPointBudget = 7

// NewCharacter builds a Character with resources derived from vectors and
// validates its abilities. It does not enforce the starting point budget:
// it is also used to build NPCs and fixtures whose vectors intentionally
// fall outside that budget. Use NewStartingCharacter for player characters
// at creation time.
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

// NewStartingCharacter builds a Character like NewCharacter, but additionally
// enforces the starting point budget (R + G + B must equal
// StartingPointBudget), per docs/core/en/core/character_creation.md.
func NewStartingCharacter(id, name string, vectors Vectors, abilities []Ability) (Character, error) {
	if total := vectors.R + vectors.G + vectors.B; total != StartingPointBudget {
		return Character{}, fmt.Errorf("starting character must distribute exactly %d points, got %d", StartingPointBudget, total)
	}
	return NewCharacter(id, name, vectors, abilities)
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
