# API Gateway

RESTful API backend for the SEO Tech Platform built with Go and Gin framework.

## Features

- 🔐 JWT authentication
- 👤 User management
- 📁 Project management
- 🔍 Audit orchestration
- 📊 Results aggregation
- 🗃️ PostgreSQL database with GORM

## Development

```bash
# Install dependencies
go mod download

# Run locally
go run cmd/server/main.go

# Build
go build -o bin/api-server ./cmd/server

# Run tests
go test ./...
```

## API Endpoints

### Authentication

- `POST /api/v1/auth/register` - Register new user
- `POST /api/v1/auth/login` - Login

### Projects (Protected)

- `GET /api/v1/projects` - List all projects
- `POST /api/v1/projects` - Create project
- `GET /api/v1/projects/:id` - Get project details
- `PUT /api/v1/projects/:id` - Update project
- `DELETE /api/v1/projects/:id` - Delete project

### Audits (Protected)

- `POST /api/v1/audits/start` - Start new audit
- `GET /api/v1/audits/:id` - Get audit details
- `GET /api/v1/audits/project/:project_id` - List audits for project
- `GET /api/v1/audits/:id/pages` - Get page results

## Configuration

Set environment variables:

- `DB_HOST`, `DB_PORT`, `DB_NAME`, `DB_USER`, `DB_PASSWORD`
- `REDIS_URL`
- `JWT_SECRET`
- `CRAWLER_SERVICE_URL`
- `ANALYZER_SERVICE_URL`
