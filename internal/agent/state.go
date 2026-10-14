package agent

import (
	"encoding/json"
	"time"
)

type InvestigationStatus string

const (
	StatusInvestigating InvestigationStatus = "investigating"
	StatusCompleted     InvestigationStatus = "completed"
	StatusFailed        InvestigationStatus = "failed"
	StatusMaxSteps      InvestigationStatus = "max_steps_reached"
)

type Investigation struct {
	IncidentID string `json:"incident_id"`

	TargetCluster string `json:"target_cluster"`
	Symptom       string `json:"symptom"`

	Status InvestigationStatus `json:"status"`

	StepCount int `json:"step_count"`
	MaxSteps  int `json:"max_steps"`

	Evidence        []Evidence       `json:"evidence"`
	ExecutedActions []Action         `json:"executed_actions"`
	Decisions       []DecisionRecord `json:"decisions"`
	Hypotheses      []Hypothesis     `json:"hypotheses"`

	Confidence float64 `json:"confidence"`

	NextAction string `json:"next_action,omitempty"`
	Stop       bool   `json:"stop"`

	Summary string `json:"summary,omitempty"`

	StartedAt time.Time `json:"started_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Evidence struct {
	Diagnostic string          `json:"diagnostic"`
	Results    json.RawMessage `json:"results"`

	CollectedAt time.Time `json:"collected_at"`
}

type Action struct {
	Diagnostic string `json:"diagnostic"`
	Step       int    `json:"step"`

	ExecutedAt time.Time `json:"executed_at"`
}

type DecisionRecord struct {
	Step           int       `json:"step"`
	NextAction     string    `json:"next_action,omitempty"`
	Reason         string    `json:"reason"`
	EvidenceNeeded string    `json:"evidence_needed,omitempty"`
	Confidence     float64   `json:"confidence"`
	Stop           bool      `json:"stop"`
	DecidedAt      time.Time `json:"decided_at"`
}

type Hypothesis struct {
	Description string  `json:"description"`
	Confidence  float64 `json:"confidence"`

	Evidence []string `json:"evidence,omitempty"`
}

func NewInvestigation(
	incidentID string,
	targetCluster string,
	symptom string,
) *Investigation {

	now := time.Now().UTC()

	return &Investigation{
		IncidentID:    incidentID,
		TargetCluster: targetCluster,
		Symptom:       symptom,

		Status: StatusInvestigating,

		StepCount: 0,
		MaxSteps:  8,

		Evidence:        []Evidence{},
		ExecutedActions: []Action{},
		Decisions:       []DecisionRecord{},
		Hypotheses:      []Hypothesis{},

		Confidence: 0,
		Stop:       false,

		StartedAt: now,
		UpdatedAt: now,
	}
}

func (investigation *Investigation) AddEvidence(
	diagnostic string,
	results json.RawMessage,
) {
	now := time.Now().UTC()

	investigation.StepCount++

	investigation.Evidence = append(
		investigation.Evidence,
		Evidence{
			Diagnostic:  diagnostic,
			Results:     results,
			CollectedAt: now,
		},
	)

	investigation.ExecutedActions = append(
		investigation.ExecutedActions,
		Action{
			Diagnostic: diagnostic,
			Step:       investigation.StepCount,
			ExecutedAt: now,
		},
	)

	investigation.UpdatedAt = now

	if investigation.StepCount >= investigation.MaxSteps {
		investigation.Status = StatusMaxSteps
		investigation.Stop = true
	}
}

func (investigation *Investigation) AddDecision(
	decision Decision,
) {
	now := time.Now().UTC()

	investigation.Decisions = append(
		investigation.Decisions,
		DecisionRecord{
			Step:           investigation.StepCount + 1,
			NextAction:     decision.NextAction,
			Reason:         decision.Reason,
			EvidenceNeeded: decision.EvidenceNeeded,
			Confidence:     decision.Confidence,
			Stop:           decision.Stop,
			DecidedAt:      now,
		},
	)

	investigation.UpdatedAt = now
}
