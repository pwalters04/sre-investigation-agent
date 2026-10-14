package agent

type Reasoner interface {
	Decide(investigation *Investigation) (Decision, error)
}
