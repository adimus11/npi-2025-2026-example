package storage

import (
	"testing"
	"time"
)

func TestEphemeralStorage(t *testing.T) {
	store := NewEphemeralStorage()
	defer store.Close()

	// Test Store
	note := &Note{
		ID:        "test-1",
		Content:   "Test note",
		CreatedAt: time.Now(),
	}

	err := store.Store(note)
	if err != nil {
		t.Fatalf("Failed to store note: %v", err)
	}

	// Test Get
	retrieved, err := store.Get("test-1")
	if err != nil {
		t.Fatalf("Failed to get note: %v", err)
	}

	if retrieved.ID != note.ID || retrieved.Content != note.Content {
		t.Errorf("Retrieved note doesn't match. Got %+v, want %+v", retrieved, note)
	}

	// Test Get non-existent
	_, err = store.Get("non-existent")
	if err == nil {
		t.Error("Expected error for non-existent note, got nil")
	}

	// Test List
	note2 := &Note{
		ID:        "test-2",
		Content:   "Another test note",
		CreatedAt: time.Now(),
	}
	store.Store(note2)

	notes, err := store.List()
	if err != nil {
		t.Fatalf("Failed to list notes: %v", err)
	}

	if len(notes) != 2 {
		t.Errorf("Expected 2 notes, got %d", len(notes))
	}
}
