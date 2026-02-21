# MLP Project Structure

Complete file structure of the MLP production-ready backend.

```
mlp/
├── cmd/
│   └── main.go                              # Application entry point, server initialization
│
├── internal/
│   ├── config/
│   │   └── config.go                        # Environment configuration loader
│   │
│   ├── middleware/
│   │   └── auth.go                          # JWT authentication middleware
│   │
│   ├── user/
│   │   ├── model.go                         # User data structures
│   │   ├── repository.go                    # User database operations
│   │   ├── service.go                       # User business logic
│   │   └── handler.go                       # User HTTP handlers
│   │
│   ├── lecture/
│   │   ├── model.go                         # Lecture data structures
│   │   ├── repository.go                    # Lecture database operations
│   │   ├── service.go                       # Lecture business logic
│   │   └── handler.go                       # Lecture HTTP handlers
│   │
│   ├── audio/
│   │   ├── model.go                         # Audio data structures
│   │   ├── repository.go                    # Audio database operations
│   │   ├── service.go                       # Audio business logic
│   │   └── handler.go                       # Audio HTTP handlers
│   │
│   ├── video/
│   │   ├── model.go                         # Video data structures
│   │   ├── repository.go                    # Video database operations
│   │   ├── service.go                       # Video business logic
│   │   └── handler.go                       # Video HTTP handlers
│   │
│   └── storage/
│       └── minio.go                         # MinIO client for object storage
│
├── migrations/
│   └── 001_init.sql                         # Initial database schema
│
├── lipsync-service/                         # Python service for lip sync
│   ├── app.py                               # Flask API server
│   ├── wav2lip.py                           # Sync.so client placeholder
│   ├── requirements.txt                     # Python dependencies
│   └── Dockerfile                           # Python service container
│
├── .env                                     # Environment variables (not in git)
├── .env.example                             # Environment template
├── .gitignore                               # Git ignore rules
├── .dockerignore                            # Docker ignore rules
│
├── docker-compose.yml                       # Multi-service orchestration
├── Dockerfile                               # Go application container
├── Makefile                                 # Build and run commands
│
├── go.mod                                   # Go dependencies
├── go.sum                                   # Go dependency checksums
│
├── MLP.postman_collection.json              # Postman API collection
├── MLP.postman_environment.json             # Postman local environment
├── MLP-Docker.postman_environment.json      # Postman docker environment
│
├── README.md                                # Main project documentation
└── POSTMAN.md                               # Postman collection guide
```

## File Descriptions

### Core Application

| File | Purpose |
|------|---------|
| `cmd/main.go` | Application entry point, initializes all services, sets up routes, starts HTTP server |
| `internal/config/config.go` | Loads and validates environment variables, provides typed configuration |
| `internal/middleware/auth.go` | JWT token validation, request authentication |
| `internal/storage/minio.go` | MinIO client wrapper for file uploads/downloads |

### Domain Modules

Each domain (user, lecture, audio, video) follows clean architecture with 4 layers:

| Layer | File | Responsibility |
|-------|------|----------------|
| **Model** | `model.go` | Data structures, request/response types |
| **Repository** | `repository.go` | Database queries, data access layer |
| **Service** | `service.go` | Business logic, validation, transformations |
| **Handler** | `handler.go` | HTTP request/response handling, routing |

### Database

| File | Purpose |
|------|---------|
| `migrations/001_init.sql` | Creates tables: users, lectures, audios, videos with indexes |

### Containerization

| File | Purpose |
|------|---------|
| `Dockerfile` | Multi-stage build for Go application |
| `docker-compose.yml` | Orchestrates: app, postgres, minio, lipsync services |
| `lipsync-service/Dockerfile` | Python Flask service container |

### API Testing

| File | Purpose |
|------|---------|
| `MLP.postman_collection.json` | Complete API collection with all endpoints |
| `MLP.postman_environment.json` | Environment variables for local testing |
| `MLP-Docker.postman_environment.json` | Environment variables for Docker testing |
| `POSTMAN.md` | Detailed guide for using Postman collection |

### Configuration

| File | Purpose |
|------|---------|
| `.env` | Actual environment variables (gitignored) |
| `.env.example` | Template for environment variables |
| `go.mod` | Go module dependencies |
| `Makefile` | Common tasks: run, build, migrate, docker commands |

## Architecture Overview

```
┌─────────────────────────────────────────────────────────────┐
│                     HTTP Requests (REST API)                 │
└──────────────────────┬──────────────────────────────────────┘
                       │
┌──────────────────────▼──────────────────────────────────────┐
│                    Handlers (HTTP Layer)                     │
│  - Parse requests                                            │
│  - Validate input                                            │
│  - Call services                                             │
│  - Format responses                                          │
└──────────────────────┬──────────────────────────────────────┘
                       │
┌──────────────────────▼──────────────────────────────────────┐
│                  Services (Business Logic)                   │
│  - Domain logic                                              │
│  - Validation rules                                          │
│  - Orchestration                                             │
│  - Call repositories                                         │
└──────────────────────┬──────────────────────────────────────┘
                       │
┌──────────────────────▼──────────────────────────────────────┐
│               Repositories (Data Access)                     │
│  - Database queries                                          │
│  - CRUD operations                                           │
│  - Data mapping                                              │
└──────────────────────┬──────────────────────────────────────┘
                       │
┌──────────────────────▼──────────────────────────────────────┐
│                  Database (PostgreSQL)                       │
│  - Users, Lectures, Audios, Videos tables                   │
└──────────────────────────────────────────────────────────────┘
```

## Dependencies Flow

```
main.go
  ├── config → Load environment variables
  ├── database → Connect to PostgreSQL
  ├── minio → Initialize object storage
  │
  ├── user
  │   ├── repository (needs: db)
  │   ├── service (needs: repository, jwt config)
  │   └── handler (needs: service, logger)
  │
  ├── lecture
  │   ├── repository (needs: db)
  │   ├── service (needs: repository)
  │   └── handler (needs: service, logger)
  │
  ├── audio
  │   ├── repository (needs: db)
  │   ├── service (needs: repository)
  │   └── handler (needs: service, logger)
  │
  ├── video
  │   ├── repository (needs: db)
  │   ├── service (needs: repository)
  │   └── handler (needs: service, logger)
  │
  └── middleware
      └── auth (needs: jwt secret, logger)
```

## Data Flow Example: Create Lecture

```
1. Client sends POST /api/v1/lectures
   └─> Request: { "topic": "ML", "content": "..." }

2. Middleware validates JWT token
   └─> Extracts user_id from token

3. Handler (lecture/handler.go)
   └─> Parses request body
   └─> Calls service.Create(ctx, user_id, req)

4. Service (lecture/service.go)
   └─> Creates Lecture model with UUID
   └─> Calls repository.Create(ctx, lecture)

5. Repository (lecture/repository.go)
   └─> Executes SQL INSERT
   └─> Returns error or success

6. Handler formats response
   └─> Returns 201 Created with lecture JSON
```

## Key Design Patterns

1. **Clean Architecture**: Separation of concerns across layers
2. **Dependency Injection**: Services receive dependencies via constructors
3. **Repository Pattern**: Abstract data access
4. **Middleware Chain**: Composable request processing
5. **Context Propagation**: Request context flows through all layers
6. **Interface-Based Design**: Repository and Service interfaces

## Environment Variables

All configuration is externalized in `.env`:

- **App**: Port
- **Database**: Host, port, credentials
- **JWT**: Secret, expiration
- **MinIO**: Endpoint, credentials, bucket
- **APIs**: Gemini, TTS, Sync keys

## Docker Services

When running `docker-compose up`:

1. **postgres** (port 5432): Database with auto-migrations
2. **minio** (ports 9000, 9001): Object storage with console
3. **lipsync** (port 5000): Python service for future integration
4. **app** (port 8080): Go backend API

## Testing Flow

1. Import Postman collection
2. Run Health Check
3. Register new user
4. Login (auto-saves JWT)
5. Create lecture (authenticated)
6. Get lectures, audios, videos

See [POSTMAN.md](POSTMAN.md) for complete testing guide.

## Code Statistics

- **Go Files**: 21
- **SQL Files**: 1
- **Python Files**: 2
- **Total Lines**: ~2,500+ (excluding dependencies)
- **API Endpoints**: 10
- **Database Tables**: 4

## Next Steps

To extend this project:

1. Add file upload endpoints (use MinIO client)
2. Implement audio/video generation logic
3. Add pagination to list endpoints
4. Implement soft deletes
5. Add unit tests
6. Add integration tests
7. Set up CI/CD pipeline
8. Add API rate limiting
9. Implement WebSocket for real-time updates
10. Add OpenAPI/Swagger documentation
