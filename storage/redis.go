package storage

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/redis/go-redis/v9"
)

// RedisStorage implements Redis-based storage
type RedisStorage struct {
	client *redis.Client
	ctx    context.Context
}

// NewRedisStorage creates a new Redis storage instance
func NewRedisStorage(addr, password string, db int) (*RedisStorage, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})

	ctx := context.Background()
	if err := client.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("failed to connect to Redis: %w", err)
	}

	return &RedisStorage{
		client: client,
		ctx:    ctx,
	}, nil
}

// Store saves a note in Redis
func (s *RedisStorage) Store(note *Note) error {
	data, err := json.Marshal(note)
	if err != nil {
		return fmt.Errorf("failed to marshal note: %w", err)
	}

	key := fmt.Sprintf("note:%s", note.ID)
	if err := s.client.Set(s.ctx, key, data, 0).Err(); err != nil {
		return fmt.Errorf("failed to store note: %w", err)
	}

	return nil
}

// Get retrieves a note by ID from Redis
func (s *RedisStorage) Get(id string) (*Note, error) {
	key := fmt.Sprintf("note:%s", id)
	data, err := s.client.Get(s.ctx, key).Bytes()
	if err == redis.Nil {
		return nil, errors.New("note not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get note: %w", err)
	}

	var note Note
	if err := json.Unmarshal(data, &note); err != nil {
		return nil, fmt.Errorf("failed to unmarshal note: %w", err)
	}

	return &note, nil
}

// List returns all notes from Redis
func (s *RedisStorage) List() ([]*Note, error) {
	keys, err := s.client.Keys(s.ctx, "note:*").Result()
	if err != nil {
		return nil, fmt.Errorf("failed to list keys: %w", err)
	}

	notes := make([]*Note, 0, len(keys))
	for _, key := range keys {
		data, err := s.client.Get(s.ctx, key).Bytes()
		if err != nil {
			continue
		}

		var note Note
		if err := json.Unmarshal(data, &note); err != nil {
			continue
		}
		notes = append(notes, &note)
	}

	return notes, nil
}

// Close closes the Redis connection
func (s *RedisStorage) Close() error {
	return s.client.Close()
}
