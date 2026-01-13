# Crawler Service

High-performance web crawler built with Golang and Colly framework.

## Features

- 🚀 Fast concurrent crawling
- 🔄 Redis-based job queue
- 📊 HTML parsing and data extraction
- 🎯 Configurable crawl depth and limits
- 📝 Structured logging

## Development

```bash
# Install dependencies
go mod download

# Run locally
go run cmd/worker/main.go

# Build
go build -o bin/crawler-worker ./cmd/worker

# Run tests
go test ./...
```

## Configuration

Set environment variables in `.env`:

- `REDIS_URL`: Redis connection string
- `DB_HOST`, `DB_PORT`, `DB_NAME`, `DB_USER`, `DB_PASSWORD`: Database configuration
- `LOG_LEVEL`: Logging level (debug, info, warn, error)

## Architecture

```
crawler-service/
├── cmd/worker/        - Entry point for crawler worker
├── internal/
│   ├── engine/        - Core crawler logic using Colly
│   ├── parser/        - HTML parsing utilities
│   └── queue/         - Redis queue implementation
└── pkg/               - Shared packages (config, logger)
```
