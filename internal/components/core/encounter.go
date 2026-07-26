package core

import "fmt"

type Encounter struct {
	ID      string
	Name    string
	Actions []Action
}

type EncounterResult struct {
	ActionResults  []Resolution
	UndefinedSteps []string
}

func RunEncounter(characters map[string]*Character, encounter Encounter) (EncounterResult, error) {
	if encounter.ID == "" {
		return EncounterResult{}, fmt.Errorf("encounter ID must be non-empty")
	}
	if encounter.Name == "" {
		return EncounterResult{}, fmt.Errorf("encounter name must be non-empty")
	}
	if len(encounter.Actions) == 0 {
		return EncounterResult{}, fmt.Errorf("encounter actions must be non-empty")
	}

	result := EncounterResult{}
	for _, action := range encounter.Actions {
		if err := action.Validate(); err != nil {
			result.UndefinedSteps = append(result.UndefinedSteps, err.Error())
			continue
		}
		actor := characters[action.Actor]
		target := characters[action.Target]
		if actor == nil {
			result.UndefinedSteps = append(result.UndefinedSteps, fmt.Sprintf("missing actor %s", action.Actor))
			continue
		}
		if target == nil {
			result.UndefinedSteps = append(result.UndefinedSteps, fmt.Sprintf("missing target %s", action.Target))
			continue
		}

		actingValue, err := actor.Vectors.Value(action.PrimaryVector)
		if err != nil {
			result.UndefinedSteps = append(result.UndefinedSteps, err.Error())
			continue
		}
		opposingValue, err := DefenseValue(*target, action.Procedure)
		if err != nil {
			result.UndefinedSteps = append(result.UndefinedSteps, err.Error())
			continue
		}
		resolution := Resolve(actingValue, 0, opposingValue)
		result.ActionResults = append(result.ActionResults, resolution)
		if resolution.Outcome.Successful() && action.TargetStateChange != "" {
			_ = target.AddState(action.TargetStateChange)
		}
	}
	return result, nil
}
