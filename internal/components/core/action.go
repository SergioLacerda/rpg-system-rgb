package core

import "fmt"

// Action is one declared step of an encounter: an actor applying a vector
// against a target through a named procedure, with its intended state
// change and narrative consequence.
type Action struct {
	Actor             string
	Target            string
	PrimaryVector     Vector
	SecondaryVector   Vector
	Intent            string
	TargetStateChange State
	Procedure         Procedure
	Consequence       string
	// Round is the 1-indexed encounter round this action occurs in. Zero
	// means unspecified and is treated as round 1.
	Round int
}

// Validate reports whether the action satisfies the minimal action
// classification contract: actor, target, vectors, intent, state change,
// procedure, and consequence must all be well-formed.
func (action Action) Validate() error {
	if action.Round < 0 {
		return fmt.Errorf("action round must be non-negative")
	}
	if action.Actor == "" {
		return fmt.Errorf("action actor must be non-empty")
	}
	if action.Target == "" {
		return fmt.Errorf("action target must be non-empty")
	}
	if err := action.PrimaryVector.Validate(); err != nil {
		return fmt.Errorf("primary vector invalid: %w", err)
	}
	if action.SecondaryVector != "" {
		if err := action.SecondaryVector.Validate(); err != nil {
			return fmt.Errorf("secondary vector invalid: %w", err)
		}
	}
	if action.Intent == "" {
		return fmt.Errorf("action intent must be non-empty")
	}
	if action.TargetStateChange != "" && !validState(action.TargetStateChange) {
		return fmt.Errorf("unknown target state change %q", action.TargetStateChange)
	}
	if err := action.Procedure.Validate(); err != nil {
		return err
	}
	if action.Consequence == "" {
		return fmt.Errorf("action consequence must be non-empty")
	}
	return nil
}
