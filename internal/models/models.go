package models

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"
)

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
	UmaName   string
	EventName string
	URL       string
	Type      EventType
}

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

func (s *EventStore) ExportJSON(dir, filename string) error {
	if dir == "" {
		return fmt.Errorf("directory path cannot be empty")
	}
	if filename == "" {
		return fmt.Errorf("filename cannot be empty")
	}
	if filepath.Ext(filename) != ".json" {
		return fmt.Errorf("filename must have a .json extension")
	}
	if err := os.MkdirAll(dir, os.FileMode(0755)); err != nil {
		return err
	}

	path := filepath.Join(dir, filename)

	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	enc := json.NewEncoder(file)
	enc.SetIndent("", "  ")
	return enc.Encode(s.Events)
}
