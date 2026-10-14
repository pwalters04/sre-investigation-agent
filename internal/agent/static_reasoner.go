package agent

type StaticReasoner struct {
	Decisions []Decision
	index     int
}

func (reasoner *StaticReasoner) Decide(
	investigation *Investigation,
) (Decision, error) {

	if reasoner.index >= len(reasoner.Decisions) {
		return Decision{
			Reason:     "No additional diagnostic decisions are available.",
			Confidence: investigation.Confidence,
			Stop:       true,
		}, nil
	}

	decision := reasoner.Decisions[reasoner.index]
	reasoner.index++

	return decision, nil
}
