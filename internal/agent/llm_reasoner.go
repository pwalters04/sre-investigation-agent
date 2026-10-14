package agent

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/agenticgokit/agenticgokit/v1beta"
	_ "github.com/agenticgokit/agenticgokit/plugins/llm/openai"
)

type LLMReasoner struct {
	agent v1beta.Agent
}

func NewLLMReasoner() (*LLMReasoner, error) {

	config := &v1beta.Config{
		Name: "sre-investigation-reasoner",
		SystemPrompt: `
You are an SRE investigation reasoning agent.

Your job is to examine the current investigation state and decide
what diagnostic evidence is needed for the next step.

You may only select one of these diagnostic actions:

- node_health
- pod_health
- events
- workload_health

You must never request shell commands, kubectl commands,
remediation, mutation, deletion, restart, scaling, or configuration changes.

Return only valid JSON matching this structure:

{
  "next_action": "node_health | pod_health | events | workload_health",
  "reason": "brief reason for the decision",
  "evidence_needed": "what evidence this diagnostic should provide",
  "confidence": 0.0,
  "stop": false
}

When sufficient evidence exists to complete the investigation,
return:

{
  "reason": "brief conclusion based on the available evidence",
  "confidence": 0.0,
  "stop": true
}

Do not include markdown.
Do not include code fences.
Do not include any text outside the JSON object.
`,

		Timeout: 60 * time.Second,

		LLM: v1beta.LLMConfig{
			Provider:    "openai",
			Model:       "gpt-4.1-mini",
			Temperature: 0.1,
			MaxTokens:   500,
		},
	}

	llmAgent, err := v1beta.NewBuilder(
		"sre-investigation-reasoner",
	).
		WithConfig(config).
		Build()

	if err != nil {
		return nil, fmt.Errorf(
			"failed to build LLM reasoner: %w",
			err,
		)
	}

	return &LLMReasoner{
		agent: llmAgent,
	}, nil
}

func (reasoner *LLMReasoner) Decide(
	investigation *Investigation,
) (Decision, error) {

	investigationJSON, err := json.MarshalIndent(
		investigation,
		"",
		"  ",
	)
	if err != nil {
		return Decision{}, fmt.Errorf(
			"failed to serialize investigation: %w",
			err,
		)
	}

	prompt := fmt.Sprintf(
		`Review the following Kubernetes investigation state and decide the next action.

Investigation:

%s`,
		string(investigationJSON),
	)

	result, err := reasoner.agent.Run(
		context.Background(),
		prompt,
	)
	if err != nil {
		return Decision{}, fmt.Errorf(
			"LLM reasoning failed: %w",
			err,
		)
	}

	var decision Decision

	if err := json.Unmarshal(
		[]byte(result.Content),
		&decision,
	); err != nil {
		return Decision{}, fmt.Errorf(
			"failed to parse LLM decision: %w\nresponse: %s",
			err,
			result.Content,
		)
	}

	return decision, nil
}