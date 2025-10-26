package storage

import "time"

// Note represents a simple note with id, content and timestamp
type Note struct {
	ID        string    `json:"id"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
}

// Storage defines the interface for note storage backends
type Storage interface {
	Store(note *Note) error
	Get(id string) (*Note, error)
	List() ([]*Note, error)
	Close() error
}
