package scenario

import (
	"fmt"
	"os"
	"strings"
)

type Scenario struct {
	Path     string
	Incident string
}

func Load(path string) (*Scenario, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf(
			"failed to read scenario file %q: %w",
			path,
			err,
		)
	}

	content := string(data)

	incident, err := extractSection(
		content,
		"# Incident",
	)
	if err != nil {
		return nil, fmt.Errorf(
			"failed to load incident from scenario %q: %w",
			path,
			err,
		)
	}

	return &Scenario{
		Path:     path,
		Incident: incident,
	}, nil
}

func extractSection(
	content string,
	heading string,
) (string, error) {

	start := strings.Index(content, heading)
	if start == -1 {
		return "", fmt.Errorf(
			"section %q not found",
			heading,
		)
	}

	section := content[start+len(heading):]

	section = strings.TrimSpace(section)

	if nextHeading := strings.Index(section, "\n# "); nextHeading != -1 {
		section = section[:nextHeading]
	}

	section = strings.TrimSpace(section)

	if section == "" {
		return "", fmt.Errorf(
			"section %q is empty",
			heading,
		)
	}

	return section, nil
}
