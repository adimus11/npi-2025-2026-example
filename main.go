package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/adimus11/npi-2025-2026-example/storage"
	"github.com/google/uuid"
)

type Server struct {
	storage storage.Storage
}

type CreateNoteRequest struct {
	Content string `json:"content"`
}

func main() {
	storageBackend := os.Getenv("STORAGE_BACKEND")
	if storageBackend == "" {
		storageBackend = "ephemeral"
	}

	var store storage.Storage
	var err error

	switch storageBackend {
	case "ephemeral":
		log.Println("Using ephemeral (in-memory) storage")
		store = storage.NewEphemeralStorage()
	case "redis":
		log.Println("Using Redis storage")
		redisAddr := os.Getenv("REDIS_ADDR")
		if redisAddr == "" {
			redisAddr = "localhost:6379"
		}
		redisPassword := os.Getenv("REDIS_PASSWORD")
		redisDB := 0
		if dbStr := os.Getenv("REDIS_DB"); dbStr != "" {
			redisDB, err = strconv.Atoi(dbStr)
			if err != nil {
				log.Fatalf("Invalid REDIS_DB: %v", err)
			}
		}
		store, err = storage.NewRedisStorage(redisAddr, redisPassword, redisDB)
		if err != nil {
			log.Fatalf("Failed to initialize Redis storage: %v", err)
		}
	case "postgres":
		log.Println("Using Postgres storage")
		connStr := os.Getenv("POSTGRES_CONN")
		if connStr == "" {
			connStr = "host=localhost port=5432 user=postgres password=postgres dbname=notes sslmode=disable"
		}
		store, err = storage.NewPostgresStorage(connStr)
		if err != nil {
			log.Fatalf("Failed to initialize Postgres storage: %v", err)
		}
	default:
		log.Fatalf("Unknown storage backend: %s", storageBackend)
	}

	defer store.Close()

	server := &Server{storage: store}

	http.HandleFunc("/notes", server.handleNotes)
	http.HandleFunc("/notes/", server.handleNote)
	http.HandleFunc("/health", handleHealth)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server starting on port %s", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}

func (s *Server) handleNotes(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		s.listNotes(w, r)
	case http.MethodPost:
		s.createNote(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func (s *Server) handleNote(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	id := r.URL.Path[len("/notes/"):]
	if id == "" {
		http.Error(w, "Note ID required", http.StatusBadRequest)
		return
	}

	s.getNote(w, r, id)
}

func (s *Server) createNote(w http.ResponseWriter, r *http.Request) {
	var req CreateNoteRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.Content == "" {
		http.Error(w, "Content is required", http.StatusBadRequest)
		return
	}

	note := &storage.Note{
		ID:        uuid.New().String(),
		Content:   req.Content,
		CreatedAt: time.Now(),
	}

	if err := s.storage.Store(note); err != nil {
		log.Printf("Error storing note: %v", err)
		http.Error(w, "Failed to store note", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(note)
}

func (s *Server) getNote(w http.ResponseWriter, r *http.Request, id string) {
	note, err := s.storage.Get(id)
	if err != nil {
		if err.Error() == "note not found" {
			http.Error(w, "Note not found", http.StatusNotFound)
			return
		}
		log.Printf("Error getting note: %v", err)
		http.Error(w, "Failed to get note", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(note)
}

func (s *Server) listNotes(w http.ResponseWriter, r *http.Request) {
	notes, err := s.storage.List()
	if err != nil {
		log.Printf("Error listing notes: %v", err)
		http.Error(w, "Failed to list notes", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(notes)
}

func handleHealth(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, `{"status":"ok"}`)
}
