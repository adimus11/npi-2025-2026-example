package storage

import (
	"database/sql"
	"errors"
	"fmt"

	_ "github.com/lib/pq"
)

// PostgresStorage implements PostgreSQL-based storage
type PostgresStorage struct {
	db *sql.DB
}

// NewPostgresStorage creates a new Postgres storage instance
func NewPostgresStorage(connStr string) (*PostgresStorage, error) {
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Create table if it doesn't exist
	createTableSQL := `
		CREATE TABLE IF NOT EXISTS notes (
			id TEXT PRIMARY KEY,
			content TEXT NOT NULL,
			created_at TIMESTAMP NOT NULL
		)
	`
	if _, err := db.Exec(createTableSQL); err != nil {
		return nil, fmt.Errorf("failed to create table: %w", err)
	}

	return &PostgresStorage{db: db}, nil
}

// Store saves a note in Postgres
func (s *PostgresStorage) Store(note *Note) error {
	query := `
		INSERT INTO notes (id, content, created_at)
		VALUES ($1, $2, $3)
		ON CONFLICT (id) DO UPDATE
		SET content = EXCLUDED.content, created_at = EXCLUDED.created_at
	`
	_, err := s.db.Exec(query, note.ID, note.Content, note.CreatedAt)
	if err != nil {
		return fmt.Errorf("failed to store note: %w", err)
	}
	return nil
}

// Get retrieves a note by ID from Postgres
func (s *PostgresStorage) Get(id string) (*Note, error) {
	query := `SELECT id, content, created_at FROM notes WHERE id = $1`
	var note Note
	err := s.db.QueryRow(query, id).Scan(&note.ID, &note.Content, &note.CreatedAt)
	if err == sql.ErrNoRows {
		return nil, errors.New("note not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get note: %w", err)
	}
	return &note, nil
}

// List returns all notes from Postgres
func (s *PostgresStorage) List() ([]*Note, error) {
	query := `SELECT id, content, created_at FROM notes ORDER BY created_at DESC`
	rows, err := s.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to list notes: %w", err)
	}
	defer rows.Close()

	notes := make([]*Note, 0)
	for rows.Next() {
		var note Note
		if err := rows.Scan(&note.ID, &note.Content, &note.CreatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan note: %w", err)
		}
		notes = append(notes, &note)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating notes: %w", err)
	}

	return notes, nil
}

// Close closes the database connection
func (s *PostgresStorage) Close() error {
	return s.db.Close()
}
