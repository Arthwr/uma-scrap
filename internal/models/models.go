package models

import (
	"fmt"
	"sync"
)

// --- ENUMS / CONSTANTS ---
type EventType int

const (
	Trainee EventType = iota
	Support
	Scenario
)

// --- ENUM METHODS ---
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

// --- DOMAIN STRUCTS ---
type Event struct {
	UmaName   string
	EventName string
	URL       string
	Type      EventType
}

// --- AGGREGATE STORE ---
type EventStore struct {
	mux    sync.Mutex
	Events map[EventType]map[string][]Event
}

func NewEventStore() *EventStore {
	return &EventStore{
		Events: map[EventType]map[string][]Event{
			Trainee: make(map[string][]Event),
			Support: make(map[string][]Event),
		},
	}
}

func (s *EventStore) AddEvent(e Event) {
	s.mux.Lock()
	defer s.mux.Unlock()

	if _, ok := s.Events[e.Type]; !ok {
		s.Events[e.Type] = make(map[string][]Event)
	}
	s.Events[e.Type][e.UmaName] = append(s.Events[e.Type][e.UmaName], e)
}
