package storage

import (
	"errors"
	"sync"
)

// EphemeralStorage implements in-memory storage
type EphemeralStorage struct {
	mu    sync.RWMutex
	notes map[string]*Note
}

// NewEphemeralStorage creates a new ephemeral storage instance
func NewEphemeralStorage() *EphemeralStorage {
	return &EphemeralStorage{
		notes: make(map[string]*Note),
	}
}

// Store saves a note in memory
func (s *EphemeralStorage) Store(note *Note) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.notes[note.ID] = note
	return nil
}

// Get retrieves a note by ID
func (s *EphemeralStorage) Get(id string) (*Note, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	note, ok := s.notes[id]
	if !ok {
		return nil, errors.New("note not found")
	}
	return note, nil
}

// List returns all notes
func (s *EphemeralStorage) List() ([]*Note, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	notes := make([]*Note, 0, len(s.notes))
	for _, note := range s.notes {
		notes = append(notes, note)
	}
	return notes, nil
}

// Close does nothing for ephemeral storage
func (s *EphemeralStorage) Close() error {
	return nil
}
