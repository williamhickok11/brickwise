# Brickwise Backend (Java)

Spring Boot 3 / Java 21 REST API. Converted from the Go backend in `../backend-go/`.

See **[LEARNING.md](./LEARNING.md)** for C# /.NET → Java notes.

## Run

```bash
mvn spring-boot:run
```

Or: `mvn -q -DskipTests package && java -jar target/brickwise-backend-0.1.0-SNAPSHOT.jar`

## Config

| Variable | Default | Notes |
|---|---|---|
| `PORT` | `8080` | HTTP port |
| `DB_PATH` | `./brickwise.db` | SQLite path when `DATABASE_URL` unset |
| `DATABASE_URL` | _(empty)_ | Postgres URI → use Postgres instead of SQLite |
| `DOCUMENTS_DIR` | `./data/documents` | Lease PDFs |
| `PROPERTY_IMAGES_DIR` | `./data/property-images` | Property photos |
| `QBO_*` | — | QuickBooks Online OAuth (optional) |

## API

Same paths as before under `/api/v1` (properties, residences, mortgages, time-entries, documents, images, dashboard, QBO).
