# Brickwise Backend API

A clean, idiomatic Go REST API with proper layering and production-ready patterns.

## Architecture

The API follows a clean layered architecture:

```
┌─────────────┐
│   Handlers  │  HTTP layer - request/response handling
├─────────────┤
│   Service   │  Business logic - validation, orchestration
├─────────────┤
│ Repository  │  Data access - database abstraction
├─────────────┤
│     DB      │  SQLite (dev) / Postgres (production)
└─────────────┘
```

### Layers

- **Handlers** (`internal/handlers/`): HTTP request/response handling, JSON encoding/decoding
- **Service** (`internal/service/`): Business logic, validation orchestration
- **Repository** (`internal/repository/`): Data access abstraction (SQLite/Postgres compatible)
- **Models** (`internal/models/`): Domain models and DTOs
- **Errors** (`internal/errors/`): Structured error handling
- **Validator** (`internal/validator/`): Input validation
- **Logger** (`internal/logger/`): Structured logging

## Features

- ✅ **chi router** - Lightweight, idiomatic HTTP router
- ✅ **Clean layering** - Separation of concerns (handlers → service → repository)
- ✅ **Database abstraction** - Works with both SQLite (dev) and Postgres (production)
- ✅ **Structured errors** - Consistent API error responses
- ✅ **Input validation** - Field-level validation with clear error messages
- ✅ **Structured logging** - JSON/text logging with context
- ✅ **Graceful shutdown** - Proper server shutdown handling
- ✅ **Tests** - Unit tests for handlers and services

## Setup

1. Install dependencies:
```bash
go mod tidy
```

2. Run with SQLite (default):
```bash
go run cmd/server/main.go
```

3. Run with Postgres:
```bash
DATABASE_URL="postgres://user:pass@localhost/dbname?sslmode=disable" go run cmd/server/main.go
```

## Environment Variables

| Variable | Description | Default |
|---|---|---|
| `PORT` | HTTP listen port | `8080` |
| `DATABASE_URL` | Postgres connection string (omit for SQLite) | _(empty → SQLite)_ |
| `DB_PATH` | SQLite file path | `./brickwise.db` |
| `LOG_FORMAT` | `json` or text | text |
| `QBO_CLIENT_ID` | Intuit app client ID | _(required for QBO)_ |
| `QBO_CLIENT_SECRET` | Intuit app client secret | _(required for QBO)_ |
| `QBO_REDIRECT_URI` | OAuth redirect URI registered in Intuit | `http://localhost:8080/api/v1/integrations/qbo/callback` |
| `QBO_FRONTEND_REDIRECT` | Where to send the browser after OAuth | `http://localhost:5173/settings/quickbooks` |
| `QBO_ENVIRONMENT` | `sandbox` or `production` | `sandbox` |
| `DOCUMENTS_DIR` | Local directory for lease PDFs | `./data/documents` |

### QuickBooks Online

1. Create an app at [developer.intuit.com](https://developer.intuit.com/).
2. Enable the Accounting API scope.
3. Add the redirect URI above.
4. Ask your bookkeeper to use **one Class per property** and code transactions to that Class.
5. In Brickwise: **QuickBooks** → Connect → map each property to its Class → Refresh.

Endpoints:

- `GET /api/v1/integrations/qbo/status`
- `GET /api/v1/integrations/qbo/connect`
- `GET /api/v1/integrations/qbo/callback`
- `DELETE /api/v1/integrations/qbo/disconnect`
- `GET /api/v1/integrations/qbo/classes`
- `PUT /api/v1/integrations/qbo/mappings`
- `POST /api/v1/integrations/qbo/sync`

## API Endpoints

- `GET /health` - Health check
- `GET /api/v1/properties` - List all properties
- `GET /api/v1/properties/{id}` - Get property by ID
- `POST /api/v1/properties` - Create property
- `PUT /api/v1/properties/{id}` - Update property
- `DELETE /api/v1/properties/{id}` - Delete property
- `GET /api/v1/dashboard/summary` - Portfolio summary (includes QBO actuals when connected)

## Design Tradeoffs

### Router: chi vs net/http
- **Choice**: chi router
- **Reason**: Lightweight, idiomatic, better middleware support than raw net/http
- **Tradeoff**: Adds dependency, but minimal and well-maintained

### Database: SQLite vs Postgres abstraction
- **Choice**: Abstracted repository with SQLite/Postgres support
- **Reason**: SQLite for dev/testing, Postgres for production
- **Tradeoff**: Some SQL differences (RETURNING clause) handled via detection
- **Alternative**: Could use sqlx or gorm, but adds complexity

### Error Handling: Structured errors
- **Choice**: Custom APIError type
- **Reason**: Consistent API responses, easy to extend
- **Tradeoff**: Manual error wrapping vs using errors.Is/As (simpler for this scale)

### Validation: Custom vs library
- **Choice**: Custom validation functions
- **Reason**: Simple, no external dependencies, easy to understand
- **Tradeoff**: More code vs using validator library (but avoids dependency)

### Logging: slog vs structured library
- **Choice**: Standard library slog
- **Reason**: Built-in, structured, production-ready
- **Tradeoff**: Newer API (Go 1.21+), but standard library

### Testing: Unit tests only
- **Choice**: Handler and service unit tests with mocks
- **Reason**: Fast, isolated, covers business logic
- **Tradeoff**: No integration tests (add if needed for complex queries)

## Testing

Run tests:
```bash
go test ./...
```

Run with coverage:
```bash
go test -cover ./...
```

## Production Considerations

1. **Database**: Use Postgres in production (set `DATABASE_URL`)
2. **CORS**: Configure allowed origins (currently `*` for dev)
3. **Logging**: Set `LOG_FORMAT=json` for structured logs
4. **Connection Pooling**: Already configured (max 25 connections)
5. **Timeouts**: Already configured (15s read, 60s write, 60s idle)
6. **Graceful Shutdown**: Implemented with 5s timeout
7. **QuickBooks tokens**: Stored in DB; protect the API before production use

## Future Enhancements

- Add integration tests
- Add request ID middleware (already added)
- Add rate limiting middleware
- Add authentication/authorization
- Add database migrations (consider migrate or golang-migrate)
- Add OpenAPI/Swagger documentation
