package models

import "fmt"

type EventType int

const (
	Trainee EventType = iota
	Support
	Scenario
)

func (et EventType) String() string {
	switch et {
	case Trainee:
		return "Trainee"
	case Support:
		return "Support"
	case Scenario:
		return "Scenario"
	default:
		return "Unknown"
	}
}

func (et EventType) MarshalText() ([]byte, error) {
	return []byte(et.String()), nil
}

func EventTypeFromString(s string) (EventType, error) {
	switch s {
	case "Trainee":
		return Trainee, nil
	case "Support":
		return Support, nil
	case "Scenario":
		return Scenario, nil
	default:
		return -1, fmt.Errorf("unknown EventType: %s", s)
	}
}

type Event struct {
	UmaName   string    `json:"uma_name"`
	EventName string    `json:"event_name"`
	URL       string    `json:"url"`
	Type      EventType `json:"type"`
}
