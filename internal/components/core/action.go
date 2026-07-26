package core

import "fmt"

type Action struct {
	Actor             string
	Target            string
	PrimaryVector     Vector
	SecondaryVector   Vector
	Intent            string
	TargetStateChange State
	Procedure         Procedure
	Consequence       string
}

func (action Action) Validate() error {
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
