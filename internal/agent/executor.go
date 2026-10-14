package agent

import (
	"fmt"

	"github.com/pwalters04/sre-investigation-agent/internal/tools"
)

var allowedActions = map[string]bool{
	"node_health":     true,
	"pod_health":      true,
	"events":          true,
	"workload_health": true,
}

func ExecuteDecision(
	investigation *Investigation,
	decision Decision,
) error {

	if decision.Stop {
		investigation.Stop = true
		investigation.Status = StatusCompleted
		investigation.Confidence = decision.Confidence
		investigation.Summary = decision.Reason
		investigation.NextAction = ""
		return nil
	}

	if !allowedActions[decision.NextAction] {
		return fmt.Errorf(
			"decision requested unauthorized action %q",
			decision.NextAction,
		)
	}

	if investigation.StepCount >= investigation.MaxSteps {
		investigation.Stop = true
		investigation.Status = StatusMaxSteps
		investigation.NextAction = ""
		return nil
	}

	investigation.NextAction = decision.NextAction
	investigation.Confidence = decision.Confidence

	result, err := tools.RunDiagnostic(decision.NextAction)
	if err != nil {
		return fmt.Errorf(
			"failed to execute diagnostic %q: %w",
			decision.NextAction,
			err,
		)
	}

	investigation.AddEvidence(
		result.Diagnostic,
		result.Results,
	)

	return nil
}
