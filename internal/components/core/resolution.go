package core

type Outcome string

const (
	OutcomeStrongSuccess          Outcome = "strong_success"
	OutcomeSuccess                Outcome = "success"
	OutcomeSuccessWithCost        Outcome = "success_with_cost"
	OutcomeFailureWithOpportunity Outcome = "failure_with_opportunity"
	OutcomeClearFailure           Outcome = "clear_failure"
)

type Resolution struct {
	ActingValue   int
	Modifier      int
	OpposingValue int
	Margin        int
	Outcome       Outcome
}

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

func ClassifyMargin(margin int) Outcome {
	switch {
	case margin >= 4:
		return OutcomeStrongSuccess
	case margin >= 1:
		return OutcomeSuccess
	case margin == 0:
		return OutcomeSuccessWithCost
	case margin >= -3:
		return OutcomeFailureWithOpportunity
	default:
		return OutcomeClearFailure
	}
}

func (outcome Outcome) Successful() bool {
	return outcome == OutcomeStrongSuccess || outcome == OutcomeSuccess || outcome == OutcomeSuccessWithCost
}
