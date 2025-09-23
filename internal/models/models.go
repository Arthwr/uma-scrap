package models

import (
	"fmt"
	"path/filepath"
	"sync"

	"github.com/arthwr/uma-scrap/internal/config"
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
	Counts map[EventType]int
}

func NewEventStore() *EventStore {
	return &EventStore{
		Events: map[EventType]map[string][]Event{
			Trainee: make(map[string][]Event),
			Support: make(map[string][]Event),
		},
		Counts: make(map[EventType]int),
	}
}

func (s *EventStore) AddEvent(e Event) {
	s.mux.Lock()
	defer s.mux.Unlock()

	if _, ok := s.Events[e.Type]; !ok {
		s.Events[e.Type] = make(map[string][]Event)
	}

	s.Events[e.Type][e.UmaName] = append(s.Events[e.Type][e.UmaName], e)
	s.Counts[e.Type]++
}

func (s *EventStore) ExportJSON(dir string) error {
	if dir == "" {
		return fmt.Errorf("directory path cannot be empty")
	}

	if err := ensureDir(dir); err != nil {
		return err
	}

	filename := makeFilename()
	path := filepath.Join(dir, filename)

	if err := writeJSON(path, s.Events); err != nil {
		return err
	}

	return cleanupOldFiles(dir, config.DEF_EVENTS_FILENAME_PATTERN, config.MAX_STORAGE_KEEP)
}
