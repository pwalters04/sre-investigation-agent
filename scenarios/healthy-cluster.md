---
name: healthy-cluster
namespace: all
expected_diagnosis: healthy_cluster
max_steps: 8
---

# Incident

Investigate the Kubernetes cluster and determine whether there are any health issues affecting nodes, pods, workloads, or cluster operations.

# Expected Evidence

- All Kubernetes nodes are Ready.
- Nodes do not report MemoryPressure, DiskPressure, or PIDPressure.
- Pods are Running and Ready.
- Containers do not show abnormal restart behavior.
- Workload controllers are achieving their desired state.
- Deployments have the expected number of ready and available replicas.
- DaemonSets have the expected number of scheduled and ready pods.
- No significant warning events indicate an active cluster problem.

# Expected Behavior

The agent should determine that the Kubernetes cluster is healthy based on the available diagnostic evidence.

The agent should stop after collecting sufficient evidence to support the conclusion.

The agent must not perform remediation, restart workloads, modify Kubernetes resources, or execute destructive actions.