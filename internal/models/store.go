package models

import (
	"fmt"
	"path/filepath"
	"sync"

	"github.com/arthwr/uma-scrap/internal/config"
)

type EventStore struct {
	mux        sync.Mutex
	Counts     map[EventType]int
	Events     map[EventType]map[string][]Event
	Collection map[string]EventDetails
}

type ExportedData struct {
	Counts     map[EventType]int                `json:"counts"`
	Events     map[EventType]map[string][]Event `json:"events"`
	Collection map[string]EventDetails          `json:"collection"`
}

func NewEventStore() *EventStore {
	return &EventStore{
		Events:     make(map[EventType]map[string][]Event),
		Counts:     make(map[EventType]int),
		Collection: make(map[string]EventDetails),
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

func (s *EventStore) AddEventDetail(eventName string, eventDetails EventDetails) {
	s.mux.Lock()
	defer s.mux.Unlock()

	s.Collection[eventName] = eventDetails
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

	data := ExportedData{
		Counts:     s.Counts,
		Events:     s.Events,
		Collection: s.Collection,
	}

	if err := writeJSON(path, data); err != nil {
		return err
	}

	return cleanupOldFiles(dir, config.DEF_EVENTS_FILENAME_PATTERN, config.MAX_STORAGE_KEEP)
}
