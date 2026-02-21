# MLP - Production-Ready Backend

A production-ready backend API built with Go, following clean architecture principles.

## Features

- Clean Architecture structure
- PostgreSQL with pgx driver
- MinIO for object storage
- JWT authentication
- REST API with chi router
- Structured logging with zap
- Docker & Docker Compose ready
- Database migrations
- Graceful shutdown

## Tech Stack

- **Language**: Go 1.22+
- **Web Framework**: Chi Router
- **Database**: PostgreSQL 16
- **Object Storage**: MinIO
- **Authentication**: JWT
- **Logging**: Zap
- **Containerization**: Docker

## Project Structure

```
mlp/
├── cmd/
│   └── main.go                 # Application entry point
├── internal/
│   ├── config/
│   │   └── config.go           # Configuration management
│   ├── middleware/
│   │   └── auth.go             # JWT authentication middleware
│   ├── user/
│   │   ├── handler.go          # HTTP handlers
│   │   ├── service.go          # Business logic
│   │   ├── repository.go       # Database access
│   │   └── model.go            # Data models
│   ├── lecture/
│   ├── audio/
│   ├── video/
│   └── storage/
│       └── minio.go            # MinIO client
├── migrations/
│   └── 001_init.sql            # Database migrations
├── lipsync-service/            # Python service placeholder
├── .env                        # Environment variables
├── .env.example                # Environment variables template
├── docker-compose.yml          # Docker services configuration
├── Dockerfile                  # Application container
├── Makefile                    # Build and run commands
└── go.mod                      # Go dependencies
```

## Getting Started

### Prerequisites

- Go 1.22+
- Docker & Docker Compose
- PostgreSQL 16 (if running locally)
- MinIO (if running locally)

### Installation

1. Clone the repository
2. Copy environment variables:

```bash
cp .env.example .env
```

3. Update `.env` with your configuration

### Running with Docker

Start all services:

```bash
make docker-up
```

Stop all services:

```bash
make docker-down
```

View logs:

```bash
make docker-logs
```

### Running Locally

1. Install dependencies:

```bash
make deps
```

2. Start PostgreSQL and MinIO (via Docker):

```bash
docker-compose up db minio -d
```

3. Run migrations:

```bash
make migrate
```

4. Run the application:

```bash
make run
```

## API Endpoints

### Health Check

```
GET /api/v1/health
```

### Authentication

```
POST /api/v1/auth/register
POST /api/v1/auth/login
```

### Lectures (Authenticated)

```
POST /api/v1/lectures
GET  /api/v1/lectures
GET  /api/v1/lectures/{id}
```

### Audios (Authenticated)

```
GET /api/v1/audios/{id}
GET /api/v1/audios/lecture/{lecture_id}
```

### Videos (Authenticated)

```
GET /api/v1/videos/{id}
GET /api/v1/videos/audio/{audio_id}
```

## Environment Variables

See `.env.example` for all available configuration options.

Key variables:

- `APP_PORT`: Server port (default: 8080)
- `DB_HOST`: PostgreSQL host
- `DB_PORT`: PostgreSQL port
- `DB_USER`: PostgreSQL user
- `DB_PASSWORD`: PostgreSQL password
- `DB_NAME`: Database name
- `JWT_SECRET`: JWT signing secret
- `JWT_EXPIRATION`: Token expiration time
- `MINIO_ENDPOINT`: MinIO server endpoint
- `MINIO_ACCESS_KEY`: MinIO access key
- `MINIO_SECRET_KEY`: MinIO secret key
- `MINIO_BUCKET`: MinIO bucket name

## Database Schema

### Users

- `id` (UUID, Primary Key)
- `email` (VARCHAR, Unique)
- `password_hash` (VARCHAR)
- `created_at` (TIMESTAMP)

### Lectures

- `id` (UUID, Primary Key)
- `user_id` (UUID, Foreign Key)
- `topic` (VARCHAR)
- `content` (TEXT)
- `created_at` (TIMESTAMP)

### Audios

- `id` (UUID, Primary Key)
- `lecture_id` (UUID, Foreign Key)
- `url` (TEXT)
- `created_at` (TIMESTAMP)

### Videos

- `id` (UUID, Primary Key)
- `audio_id` (UUID, Foreign Key)
- `url` (TEXT)
- `created_at` (TIMESTAMP)

## Makefile Commands

```bash
make run            # Run the application locally
make build          # Build the application binary
make migrate        # Run database migrations
make docker-up      # Start all services with docker-compose
make docker-down    # Stop all services
make docker-logs    # Show docker logs
make clean          # Clean build artifacts and stop docker
make test           # Run tests
make deps           # Download and tidy dependencies
make fmt            # Format code
make lint           # Run linter
```

## Services

### PostgreSQL

- Port: 5432
- Admin UI: N/A

### MinIO

- API Port: 9000
- Console Port: 9001
- Access console at: http://localhost:9001

### Lipsync Service

- Port: 5000
- Python-based service for future lip sync integration

## Development

### Adding New Migrations

Create a new SQL file in the `migrations/` directory with format `XXX_description.sql`.

### Code Style

Format your code:

```bash
make fmt
```

Run linter:

```bash
make lint
```

## Production Considerations

- Update JWT secret in production
- Use strong database passwords
- Enable SSL for MinIO in production
- Configure proper CORS settings
- Set appropriate rate limits
- Enable request logging
- Use environment-specific configurations
- Implement proper backup strategies

## License

MIT
