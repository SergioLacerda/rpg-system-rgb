package core

// State is a tactical condition a character can hold, per the minimal Core
// V2 state model (decision matrix D-006).
type State string

const (
	// StateHealthy means the character has taken no health-affecting damage.
	StateHealthy State = "healthy"
	// StateInjured means the character has taken health damage but is not downed.
	StateInjured State = "injured"
	// StateDowned means the character has no health remaining.
	StateDowned State = "downed"
	// StateExposed means the character lacks cover or concealment.
	StateExposed State = "exposed"
	// StateCovered means the character gained cover or concealment.
	StateCovered State = "covered"
	// StateEngaged means the character is in close contact with a threat.
	StateEngaged State = "engaged"
	// StateDistant means the character has increased range from a threat.
	StateDistant State = "distant"
	// StateSuppressed means the character is under sustained pressure that limits action.
	StateSuppressed State = "suppressed"
	// StateDisarmed means the character has lost a weapon or means of R pressure.
	StateDisarmed State = "disarmed"
	// StateGuarded means the character is protected by a sustaining ally or effect.
	StateGuarded State = "guarded"
	// StateShielded means the character has shield reserve remaining.
	StateShielded State = "shielded"
	// StateVulnerable means the character's pressure source has been weakened by R.
	StateVulnerable State = "vulnerable"
	// StateStabilized means the character's continuity has been restored under pressure.
	StateStabilized State = "stabilized"
	// StateHidden means the character is not currently detected.
	StateHidden State = "hidden"
	// StateDetected means the character has been found despite concealment.
	StateDetected State = "detected"
)

// allStates is the single source of truth for the Core V2 state model.
// validState and AllStates both derive from it so the valid-state check and
// the enumerable list can never drift apart.
var allStates = []State{
	StateHealthy, StateInjured, StateDowned, StateExposed, StateCovered,
	StateEngaged, StateDistant, StateSuppressed, StateDisarmed, StateGuarded,
	StateShielded, StateVulnerable, StateStabilized, StateHidden, StateDetected,
}

// AllStates returns every state in the Core V2 state model.
func AllStates() []State {
	states := make([]State, len(allStates))
	copy(states, allStates)
	return states
}

func validState(state State) bool {
	for _, candidate := range allStates {
		if candidate == state {
			return true
		}
	}
	return false
}
