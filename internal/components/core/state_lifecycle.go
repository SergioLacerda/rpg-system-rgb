package core

// StateLifecycle documents how a Core V2 state enters and leaves play. It
// closes the decision matrix D-006 closure evidence: every Core state must
// have at least one producing source and a defined way it is cleared.
// ProducedBy and ClearedBy are descriptive, not executable: some states are
// produced or cleared by a named Procedure, others by a resource threshold
// (the damage pipeline) or by a declared Action.TargetStateChange.
type StateLifecycle struct {
	ProducedBy string
	ClearedBy  string
}

var stateLifecycles = map[State]StateLifecycle{
	StateHealthy: {
		ProducedBy: "character creation",
		ClearedBy:  "automatically cleared when injured or downed is added (Character.AddState)",
	},
	StateInjured: {
		ProducedBy: "damage pipeline (ApplyDamage) when health damage occurs and health remains above zero",
		ClearedBy:  "stabilize procedure (B)",
	},
	StateDowned: {
		ProducedBy: "damage pipeline (ApplyDamage) when health reaches zero",
		ClearedBy:  "recovery outside the encounter scope",
	},
	StateExposed: {
		ProducedBy: "declared Action.TargetStateChange",
		ClearedBy:  "reposition procedure (G) into cover",
	},
	StateCovered: {
		ProducedBy: "reposition procedure (G)",
		ClearedBy:  "declared Action.TargetStateChange back to exposed",
	},
	StateEngaged: {
		ProducedBy: "declared Action.TargetStateChange",
		ClearedBy:  "reposition procedure (G) to distant",
	},
	StateDistant: {
		ProducedBy: "reposition procedure (G)",
		ClearedBy:  "declared Action.TargetStateChange back to engaged",
	},
	StateSuppressed: {
		ProducedBy: "counterpressure or attack procedure (R)",
		ClearedBy:  "reposition procedure (G) or recovery",
	},
	StateDisarmed: {
		ProducedBy: "attack or counterpressure procedure (R)",
		ClearedBy:  "declared Action.TargetStateChange or scene reset",
	},
	StateGuarded: {
		ProducedBy: "sustain procedure (B)",
		ClearedBy:  "declared Action.TargetStateChange when protection ends",
	},
	StateShielded: {
		ProducedBy: "character creation, when the B-derived shield reserve is positive",
		ClearedBy:  "automatically cleared when shield reserve reaches zero (damage pipeline)",
	},
	StateVulnerable: {
		ProducedBy: "counterpressure procedure (R)",
		ClearedBy:  "declared Action.TargetStateChange or recovery",
	},
	StateStabilized: {
		ProducedBy: "stabilize procedure (B)",
		ClearedBy:  "declared Action.TargetStateChange or scene reset",
	},
	StateHidden: {
		ProducedBy: "evade procedure (G)",
		ClearedBy:  "detection (Action.TargetStateChange to detected)",
	},
	StateDetected: {
		ProducedBy: "declared Action.TargetStateChange",
		ClearedBy:  "evade or reposition procedure (G) back to hidden",
	},
}

// Lifecycle returns the declared production and clearing sources for a
// state, and whether the state has a declared lifecycle at all.
func Lifecycle(state State) (StateLifecycle, bool) {
	lifecycle, ok := stateLifecycles[state]
	return lifecycle, ok
}
