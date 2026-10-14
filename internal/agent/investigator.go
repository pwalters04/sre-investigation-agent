package agent

import "fmt"

func Investigate(
	investigation *Investigation,
	reasoner Reasoner,
) error {

	for !investigation.Stop {

		if investigation.StepCount >= investigation.MaxSteps {
			investigation.Status = StatusMaxSteps
			investigation.Stop = true
			break
		}

		decision, err := reasoner.Decide(investigation)
		if err != nil {
			investigation.Status = StatusFailed
			investigation.Stop = true

			return fmt.Errorf(
				"reasoning failed: %w",
				err,
			)
		}

		investigation.AddDecision(decision)

		err = ExecuteDecision(
			investigation,
			decision,
		)
		if err != nil {
			investigation.Status = StatusFailed
			investigation.Stop = true

			return fmt.Errorf(
				"decision execution failed: %w",
				err,
			)
		}
	}

	return nil
}
