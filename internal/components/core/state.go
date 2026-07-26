package core

type State string

const (
	StateHealthy    State = "healthy"
	StateInjured    State = "injured"
	StateDowned     State = "downed"
	StateExposed    State = "exposed"
	StateCovered    State = "covered"
	StateEngaged    State = "engaged"
	StateDistant    State = "distant"
	StateSuppressed State = "suppressed"
	StateDisarmed   State = "disarmed"
	StateGuarded    State = "guarded"
	StateShielded   State = "shielded"
	StateVulnerable State = "vulnerable"
	StateStabilized State = "stabilized"
	StateHidden     State = "hidden"
	StateDetected   State = "detected"
)

func validState(state State) bool {
	switch state {
	case StateHealthy, StateInjured, StateDowned, StateExposed, StateCovered,
		StateEngaged, StateDistant, StateSuppressed, StateDisarmed, StateGuarded,
		StateShielded, StateVulnerable, StateStabilized, StateHidden, StateDetected:
		return true
	default:
		return false
	}
}
