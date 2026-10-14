package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"

	"github.com/pwalters04/sre-investigation-agent/internal/agent"
	"github.com/pwalters04/sre-investigation-agent/internal/scenario"
)

func main() {
	scenarioPath := flag.String(
		"scenario",
		"",
		"path to an investigation scenario markdown file",
	)
	flag.Parse()

	if *scenarioPath == "" {
		log.Fatal(
			"scenario is required: --scenario <path>",
		)
	}

	testScenario, err := scenario.Load(*scenarioPath)
	if err != nil {
		log.Fatal(err)
	}

	investigation := agent.NewInvestigation(
		"INC-001",
		"cka",
		testScenario.Incident,
	)

	reasoner, err := agent.NewLLMReasoner()
	if err != nil {
		log.Fatal(err)
	}

	err = agent.Investigate(
		investigation,
		reasoner,
	)
	if err != nil {
		log.Fatal(err)
	}

	output, err := json.MarshalIndent(
		investigation,
		"",
		"  ",
	)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(string(output))
}
