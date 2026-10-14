package agent

type Decision struct {
	NextAction     string  `json:"next_action,omitempty"`
	Reason         string  `json:"reason"`
	EvidenceNeeded string  `json:"evidence_needed,omitempty"`
	Confidence     float64 `json:"confidence"`
	Stop           bool    `json:"stop"`
}
