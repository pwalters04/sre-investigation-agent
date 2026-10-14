# SRE Investigation Agent

A proof-of-concept agentic SRE investigation system built in Go.

The agent investigates Kubernetes incidents by reasoning over structured diagnostic evidence and selecting from a constrained set of read-only tools. Go manages state, guardrails, execution flow, and stopping conditions, while Ansible provides deterministic Kubernetes diagnostics.

## Prerequisites

To run the agent, you need:

- Go
- Ansible
- `kubectl`
- Access to a Kubernetes cluster
- An OpenAI API key

The LLM reasoner uses the OpenAI API, so the API key must be available as an environment variable before starting the agent.

```bash
export OPENAI_API_KEY="your-api-key"
```

You can confirm that the variable is set with:

```bash
echo $OPENAI_API_KEY
```

Do not commit the API key to the repository or store it directly in the Go source code.

## Current Architecture

```text
Incident / Scenario
        ↓
Investigation State
        ↓
LLM Reasoner
        ↓
Structured Decision
        ↓
Go Guardrails / Allowlist
        ↓
Diagnostic Executor
        ↓
Ansible
        ↓
Kubernetes
        ↓
Structured Evidence
        ↓
Investigation State
        ↓
Reason Again or Stop
```

## Current Diagnostics

The agent can currently select:

- `node_health` — node readiness, pressure conditions, capacity, and scheduled pods.
- `pod_health` — pod and container health, resources, probes, runtime state, init containers, and pod events.
- `events` — cluster-wide Kubernetes events.
- `workload_health` — Deployment, StatefulSet, and DaemonSet controller health.

All diagnostics are read-only and protected by allowlists in both Go and Ansible.

## Agent Capabilities

The current POC includes:

- Investigation state and evidence accumulation.
- LLM-based diagnostic selection.
- Multi-step reasoning loop.
- Structured decisions.
- Decision and execution history.
- Maximum investigation step limits.
- Tool authorization guardrails.
- Confidence tracking.
- Final investigation summary.
- Scenario-based incident input.

## Scenario Testing

Scenarios can be stored as Markdown files:

```text
scenarios/
├── healthy-cluster.md
└── image-pull-backoff.md
```

Run a scenario with:

```bash
go run ./cmd/investigator \
  --scenario scenarios/healthy-cluster.md
```

The scenario provides the incident context to the agent while expected evidence and outcomes can later be used for automated evaluations.

## Current Goal

The project is evolving toward a reusable SRE investigation agent that can accept real incident and alert messages, autonomously gather evidence, diagnose likely causes, and stop when sufficient evidence has been collected.

Planned next steps include additional failure scenarios, automated evals, pod log diagnostics, improved confidence handling, tracing, and deployment.