package core

// Outcome classifies the margin of a resolved action, per the canonical
// margin table in docs/core/en/combat/attack_and_defense.md.
type Outcome string

const (
	// OutcomeStrongSuccess is a margin of 3 or more (canonical: "strong hit").
	OutcomeStrongSuccess Outcome = "strong_success"
	// OutcomeSuccess is a margin of 1 or 2 (canonical: "hit").
	OutcomeSuccess Outcome = "success"
	// OutcomeSuccessWithCost is a margin of exactly 0 (canonical: "hit with cost").
	OutcomeSuccessWithCost Outcome = "success_with_cost"
	// OutcomeFailureWithOpportunity is a margin of -1 or -2 (canonical: "miss with opportunity").
	OutcomeFailureWithOpportunity Outcome = "failure_with_opportunity"
	// OutcomeClearFailure is a margin of -3 or less (canonical: "clear miss").
	OutcomeClearFailure Outcome = "clear_failure"
)

// Resolution is the result of resolving an acting value against an opposing
// value through the margin-based engine (decision matrix D-001).
type Resolution struct {
	ActingValue   int
	Modifier      int
	OpposingValue int
	Margin        int
	Outcome       Outcome
}

// Resolve computes a margin-based resolution: acting value plus modifier,
// minus opposing value, classified into an Outcome.
func Resolve(actingValue, modifier, opposingValue int) Resolution {
	margin := actingValue + modifier - opposingValue
	return Resolution{
		ActingValue:   actingValue,
		Modifier:      modifier,
		OpposingValue: opposingValue,
		Margin:        margin,
		Outcome:       ClassifyMargin(margin),
	}
}

// ClassifyMargin maps a margin to its canonical Outcome band, per
// docs/core/en/combat/attack_and_defense.md.
func ClassifyMargin(margin int) Outcome {
	switch {
	case margin >= 3:
		return OutcomeStrongSuccess
	case margin >= 1:
		return OutcomeSuccess
	case margin == 0:
		return OutcomeSuccessWithCost
	case margin >= -2:
		return OutcomeFailureWithOpportunity
	default:
		return OutcomeClearFailure
	}
}

// Successful reports whether the outcome allows a following effect (such as
// damage) to start.
func (outcome Outcome) Successful() bool {
	return outcome == OutcomeStrongSuccess || outcome == OutcomeSuccess || outcome == OutcomeSuccessWithCost
}
