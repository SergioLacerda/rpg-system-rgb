package core

import "fmt"

// Encounter is a scripted sequence of actions, optionally bound to a
// round-deadline Objective. An encounter can succeed through its objective
// without requiring every opposing character to be defeated.
type Encounter struct {
	ID      string
	Name    string
	Actions []Action
	// Objective declares the round-bound success condition for this
	// encounter. Nil means the encounter has no round-deadline objective.
	Objective *Objective
}

// Objective declares a round-deadline success condition for an encounter:
// a target character must hold a required state before the deadline round
// ends, or the encounter fails for the declared reason.
type Objective struct {
	DeadlineRounds int
	Target         string
	RequiredState  State
	FailureReason  string
}

// Validate reports whether the objective is well-formed.
func (objective Objective) Validate() error {
	if objective.DeadlineRounds < 1 {
		return fmt.Errorf("objective deadline rounds must be at least 1")
	}
	if objective.Target == "" {
		return fmt.Errorf("objective target must be non-empty")
	}
	if !validState(objective.RequiredState) {
		return fmt.Errorf("unknown objective required state %q", objective.RequiredState)
	}
	if objective.FailureReason == "" {
		return fmt.Errorf("objective failure reason must be non-empty")
	}
	return nil
}

// ObjectiveOutcome reports whether a declared Objective was met before its
// round deadline.
type ObjectiveOutcome struct {
	Declared      bool
	Succeeded     bool
	ResolvedRound int
	FailureReason string
}

// EncounterResult is the outcome of running an Encounter: each action's
// resolution, any steps that could not be resolved, and the objective
// outcome when one was declared.
type EncounterResult struct {
	ActionResults  []Resolution
	UndefinedSteps []string
	Objective      ObjectiveOutcome
}

// RunEncounter resolves every action in the encounter against the supplied
// characters, in list order, and evaluates the encounter's Objective
// against the round deadline when one is declared. An encounter can
// succeed through its objective without defeating every opposing
// character.
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
	if encounter.Objective != nil {
		if err := encounter.Objective.Validate(); err != nil {
			return EncounterResult{}, fmt.Errorf("encounter objective invalid: %w", err)
		}
	}

	result := EncounterResult{}
	currentRound := 0
	for _, action := range encounter.Actions {
		round, err := resolveAction(characters, action, &result)
		if round > currentRound {
			currentRound = round
		}
		if err != nil {
			result.UndefinedSteps = append(result.UndefinedSteps, err.Error())
		}
	}

	if encounter.Objective != nil {
		result.Objective = evaluateObjective(characters, *encounter.Objective, currentRound)
	}
	return result, nil
}

// resolveAction validates and resolves a single action against characters,
// appending its Resolution to result and applying its target state change
// on success. It returns the action's effective round, or an error
// describing why the step could not be resolved.
func resolveAction(characters map[string]*Character, action Action, result *EncounterResult) (int, error) {
	if err := action.Validate(); err != nil {
		return 0, err
	}
	round := action.Round
	if round == 0 {
		round = 1
	}

	actor := characters[action.Actor]
	if actor == nil {
		return round, fmt.Errorf("missing actor %s", action.Actor)
	}
	target := characters[action.Target]
	if target == nil {
		return round, fmt.Errorf("missing target %s", action.Target)
	}

	actingValue, err := actor.Vectors.Value(action.PrimaryVector)
	if err != nil {
		return round, err
	}
	opposingValue, err := DefenseValue(*target, action.Procedure)
	if err != nil {
		return round, err
	}

	resolution := Resolve(actingValue, 0, opposingValue)
	result.ActionResults = append(result.ActionResults, resolution)
	if resolution.Outcome.Successful() && action.TargetStateChange != "" {
		_ = target.AddState(action.TargetStateChange)
	}
	return round, nil
}

// evaluateObjective checks whether the objective's target holds the
// required state, and otherwise reports failure once the round clock has
// reached the objective's deadline.
func evaluateObjective(characters map[string]*Character, objective Objective, currentRound int) ObjectiveOutcome {
	outcome := ObjectiveOutcome{Declared: true, FailureReason: objective.FailureReason}
	target := characters[objective.Target]
	if target != nil && target.HasState(objective.RequiredState) {
		outcome.Succeeded = true
		outcome.ResolvedRound = min(currentRound, objective.DeadlineRounds)
		if outcome.ResolvedRound == 0 {
			outcome.ResolvedRound = 1
		}
		return outcome
	}
	outcome.ResolvedRound = max(currentRound, objective.DeadlineRounds)
	return outcome
}
