# npi-2025-2026-example
Example repository for NPI subject to present usage of GH Actions/Docker/Docker Compose

## Notes Server

A simple Go-based HTTP server for storing, listing, and retrieving notes. The server supports three storage backends:

- **Ephemeral** (in-memory) - default
- **Redis**
- **Postgres**

The storage backend is controlled via environment variables.

## Features

- Create notes with auto-generated UUIDs
- Retrieve notes by ID
- List all notes
- Health check endpoint
- Multiple storage backends (ephemeral, Redis, Postgres)
- Dockerized application

## API Endpoints

### POST /notes
Create a new note.

**Request Body:**
```json
{
  "content": "Your note content here"
}
```

**Response:**
```json
{
  "id": "generated-uuid",
  "content": "Your note content here",
  "created_at": "2025-10-26T08:00:00Z"
}
```

### GET /notes
List all notes.

**Response:**
```json
[
  {
    "id": "uuid-1",
    "content": "Note 1",
    "created_at": "2025-10-26T08:00:00Z"
  },
  {
    "id": "uuid-2",
    "content": "Note 2",
    "created_at": "2025-10-26T08:01:00Z"
  }
]
```

### GET /notes/{id}
Retrieve a specific note by ID.

**Response:**
```json
{
  "id": "uuid-1",
  "content": "Note 1",
  "created_at": "2025-10-26T08:00:00Z"
}
```

### GET /health
Health check endpoint.

**Response:**
```json
{
  "status": "ok"
}
```

## Environment Variables

### General
- `STORAGE_BACKEND` - Storage backend to use (`ephemeral`, `redis`, or `postgres`). Default: `ephemeral`
- `PORT` - Port to run the server on. Default: `8080`

### Redis (when STORAGE_BACKEND=redis)
- `REDIS_ADDR` - Redis server address. Default: `localhost:6379`
- `REDIS_PASSWORD` - Redis password. Default: empty
- `REDIS_DB` - Redis database number. Default: `0`

### Postgres (when STORAGE_BACKEND=postgres)
- `POSTGRES_CONN` - PostgreSQL connection string. Default: `host=localhost port=5432 user=postgres password=postgres dbname=notes sslmode=disable`

## Running Locally

### Build and Run
```bash
# Build the application
go build -o notes-server .

# Run with ephemeral storage (default)
./notes-server

# Run with Redis storage
STORAGE_BACKEND=redis REDIS_ADDR=localhost:6379 ./notes-server

# Run with Postgres storage
STORAGE_BACKEND=postgres POSTGRES_CONN="host=localhost port=5432 user=postgres password=postgres dbname=notes sslmode=disable" ./notes-server
```

## Running with Docker

### Build Docker Image
```bash
docker build -t notes-server .
```

### Run Container
```bash
# Run with ephemeral storage
docker run -p 8080:8080 -e STORAGE_BACKEND=ephemeral notes-server

# Run with Redis storage (assuming Redis is running on host)
docker run -p 8080:8080 -e STORAGE_BACKEND=redis -e REDIS_ADDR=host.docker.internal:6379 notes-server

# Run with Postgres storage (assuming Postgres is running on host)
docker run -p 8080:8080 -e STORAGE_BACKEND=postgres -e POSTGRES_CONN="host=host.docker.internal port=5432 user=postgres password=postgres dbname=notes sslmode=disable" notes-server
```

## Testing

Run the tests:
```bash
go test ./...
```

## Examples

### Create a note
```bash
curl -X POST http://localhost:8080/notes \
  -H "Content-Type: application/json" \
  -d '{"content":"My first note"}'
```

### List all notes
```bash
curl http://localhost:8080/notes
```

### Get a specific note
```bash
curl http://localhost:8080/notes/{note-id}
```

### Health check
```bash
curl http://localhost:8080/health
```

