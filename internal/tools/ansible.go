package tools

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
)

var allowedDiagnostics = map[string]bool{
	"node_health":     true,
	"pod_health":      true,
	"events":          true,
	"workload_health": true,
}

type DiagnosticResult struct {
	Diagnostic string          `json:"diagnostic"`
	Results    json.RawMessage `json:"results"`
}

type NodeHealthRaw struct {
	Nodes  []Node  `json:"nodes"`
	Pods   []Pod   `json:"pods"`
	Events []Event `json:"events"`
}

type Node struct {
	Metadata Metadata   `json:"metadata"`
	Status   NodeStatus `json:"status"`
}

type Pod struct {
	Metadata Metadata  `json:"metadata"`
	Spec     PodSpec   `json:"spec"`
	Status   PodStatus `json:"status"`
}

type Event struct {
	Reason         string         `json:"reason"`
	Message        string         `json:"message"`
	Type           string         `json:"type"`
	Count          int            `json:"count"`
	InvolvedObject InvolvedObject `json:"involvedObject"`
}

type Metadata struct {
	Name      string `json:"name"`
	Namespace string `json:"namespace"`
}

type NodeStatus struct {
	Addresses   []NodeAddress   `json:"addresses"`
	Conditions  []NodeCondition `json:"conditions"`
	Allocatable ResourceList    `json:"allocatable"`
}

type NodeAddress struct {
	Type    string `json:"type"`
	Address string `json:"address"`
}

type NodeCondition struct {
	Type    string `json:"type"`
	Status  string `json:"status"`
	Reason  string `json:"reason"`
	Message string `json:"message"`
}

type ResourceList struct {
	CPU    string `json:"cpu"`
	Memory string `json:"memory"`
	Pods   string `json:"pods"`
}

type PodSpec struct {
	NodeName   string      `json:"nodeName"`
	Containers []Container `json:"containers"`
}

type Container struct {
	Name      string    `json:"name"`
	Resources Resources `json:"resources"`
}

type Resources struct {
	Requests map[string]string `json:"requests"`
	Limits   map[string]string `json:"limits"`
}

type PodStatus struct {
	Phase string `json:"phase"`
}

type InvolvedObject struct {
	Kind string `json:"kind"`
	Name string `json:"name"`
}

type NodeHealthResult struct {
	Diagnostic string `json:"diagnostic"`
	Results    struct {
		Nodes []NodeHealth `json:"nodes"`
	} `json:"results"`
}

type NodeHealth struct {
	Name              string          `json:"name"`
	IP                string          `json:"ip"`
	Status            string          `json:"status"`
	PodCount          int             `json:"pod_count"`
	NonTerminatedPods []string        `json:"non_terminated_pods"`
	Conditions        []NodeCondition `json:"conditions"`
	Allocatable       ResourceList    `json:"allocatable"`
	Events            []NodeEvent     `json:"events"`
}

type NodeEvent struct {
	Type    string `json:"type"`
	Reason  string `json:"reason"`
	Message string `json:"message"`
	Count   int    `json:"count"`
}

func RunDiagnostic(action string) (DiagnosticResult, error) {
	if !allowedDiagnostics[action] {
		return DiagnosticResult{},
			fmt.Errorf("diagnostic action %q is not allowed", action)
	}

	tempFile, err := os.CreateTemp("", "sre-diagnostic-*.json")
	if err != nil {
		return DiagnosticResult{},
			fmt.Errorf("failed to create diagnostic output file: %w", err)
	}

	outputPath := tempFile.Name()

	if err := tempFile.Close(); err != nil {
		return DiagnosticResult{},
			fmt.Errorf("failed to close diagnostic output file: %w", err)
	}

	defer os.Remove(outputPath)

	cmd := exec.Command(
		"ansible-playbook",
		"ansible/cluster_health.yaml",
		"-e",
		fmt.Sprintf("diagnostic_action=%s", action),
		"-e",
		fmt.Sprintf("evidence_output=%s", outputPath),
		"-vv",
	)

	ansibleOutput, err := cmd.CombinedOutput()
	if err != nil {
		return DiagnosticResult{},
			fmt.Errorf(
				"diagnostic %q failed: %w\n%s",
				action,
				err,
				string(ansibleOutput),
			)
	}

	data, err := os.ReadFile(outputPath)
	if err != nil {
		return DiagnosticResult{},
			fmt.Errorf("failed to read diagnostic output: %w", err)
	}

	var result DiagnosticResult

	if err := json.Unmarshal(data, &result); err != nil {
		return DiagnosticResult{},
			fmt.Errorf("failed to parse diagnostic JSON: %w", err)
	}

	return result, nil
}

func NormalizeNodeHealth(result DiagnosticResult) (NodeHealthResult, error) {
	var raw NodeHealthRaw

	if err := json.Unmarshal(result.Results, &raw); err != nil {
		return NodeHealthResult{}, fmt.Errorf(
			"failed to parse node health results: %w",
			err,
		)
	}

	var normalized NodeHealthResult
	normalized.Diagnostic = result.Diagnostic

	for _, node := range raw.Nodes {
		nodeHealth := NodeHealth{
			Name:        node.Metadata.Name,
			IP:          getNodeIP(node),
			Status:      getNodeStatus(node),
			Conditions:  node.Status.Conditions,
			Allocatable: node.Status.Allocatable,
		}

		for _, pod := range raw.Pods {
			if pod.Spec.NodeName != node.Metadata.Name {
				continue
			}

			if isTerminalPod(pod.Status.Phase) {
				continue
			}

			nodeHealth.NonTerminatedPods = append(
				nodeHealth.NonTerminatedPods,
				fmt.Sprintf(
					"%s/%s",
					pod.Metadata.Namespace,
					pod.Metadata.Name,
				),
			)
		}

		nodeHealth.PodCount = len(nodeHealth.NonTerminatedPods)

		for _, event := range raw.Events {
			if event.InvolvedObject.Kind != "Node" {
				continue
			}

			if event.InvolvedObject.Name != node.Metadata.Name {
				continue
			}

			nodeHealth.Events = append(
				nodeHealth.Events,
				NodeEvent{
					Type:    event.Type,
					Reason:  event.Reason,
					Message: event.Message,
					Count:   event.Count,
				},
			)
		}

		normalized.Results.Nodes = append(
			normalized.Results.Nodes,
			nodeHealth,
		)
	}

	return normalized, nil
}

func getNodeIP(node Node) string {
	for _, address := range node.Status.Addresses {
		if address.Type == "InternalIP" {
			return address.Address
		}
	}

	return "unknown"
}

func getNodeStatus(node Node) string {
	for _, condition := range node.Status.Conditions {
		if condition.Type == "Ready" {
			if condition.Status == "True" {
				return "Ready"
			}

			return "NotReady"
		}
	}

	return "Unknown"
}

func isTerminalPod(phase string) bool {
	return phase == "Succeeded" || phase == "Failed"
}
