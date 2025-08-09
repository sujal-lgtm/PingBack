package services

import (
	"pingback/internal/models"
	"sync"
)

type Store struct {
	mu     sync.RWMutex
	events map[string]models.Event
}

func NewStore() *Store {
	return &Store{
		events: make(map[string]models.Event),
	}
}

func (s *Store) Save(event models.Event) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.events[event.ID] = event
}

func (s *Store) Get(id string) (models.Event, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	event, ok := s.events[id]
	return event, ok
}

func (s *Store) GetAll() []models.Event {
	s.mu.RLock()
	defer s.mu.RUnlock()

	events := make([]models.Event, 0, len(s.events))
	for _, e := range s.events {
		events = append(events, e)
	}
	return events
}

func (s *Store) GetByID(id string) (models.Event, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	e, ok := s.events[id]
	return e, ok
}

func (s *Store) Delete(id string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.events[id]; ok {
		delete(s.events, id)
		return true
	}
	return false
}
