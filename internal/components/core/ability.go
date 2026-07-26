package core

import "fmt"

// Timing declares whether an ability is used as an action or a reaction.
type Timing string

const (
	// TimingAction is used during the actor's own turn.
	TimingAction Timing = "action"
	// TimingReaction is used in response to another actor's action.
	TimingReaction Timing = "reaction"
)

// Ability is the minimal governed contract every RGB Core V2 ability must
// declare (decision matrix D-007): vector ownership, cost, timing, and
// effect are explicit rather than implied.
type Ability struct {
	ID           string
	Name         string
	Vector       Vector
	Tier         int
	Requirements map[Vector]int
	Timing       Timing
	Cost         int
	Range        int
	Duration     string
	Effects      []string
	Limits       []string
	Tags         []string
}

// Validate reports whether the ability satisfies the minimal ability
// contract required before it can be used by a character or the Specialist.
func (ability Ability) Validate() error {
	if ability.ID == "" {
		return fmt.Errorf("ability ID must be non-empty")
	}
	if ability.Name == "" {
		return fmt.Errorf("ability name must be non-empty")
	}
	if err := ability.Vector.Validate(); err != nil {
		return err
	}
	if ability.Tier < 1 {
		return fmt.Errorf("ability tier must be at least 1")
	}
	for vector, value := range ability.Requirements {
		if err := vector.Validate(); err != nil {
			return fmt.Errorf("requirement vector invalid: %w", err)
		}
		if value < 0 {
			return fmt.Errorf("requirement for %s must be non-negative", vector)
		}
	}
	switch ability.Timing {
	case TimingAction, TimingReaction:
	default:
		return fmt.Errorf("unknown ability timing %q", ability.Timing)
	}
	if ability.Cost < 0 {
		return fmt.Errorf("ability cost must be non-negative")
	}
	if ability.Range < 0 {
		return fmt.Errorf("ability range must be non-negative")
	}
	if ability.Duration == "" {
		return fmt.Errorf("ability duration must be non-empty")
	}
	if len(ability.Effects) == 0 {
		return fmt.Errorf("ability effects must be non-empty")
	}
	if len(ability.Tags) == 0 {
		return fmt.Errorf("ability tags must be non-empty")
	}
	return nil
}
