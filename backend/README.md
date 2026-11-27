# Brickwise Backend (Java)

Spring Boot 3 / Java 21 REST API.

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
| `DATABASE_URL` | _(empty)_ | Postgres / Supabase URI → use Postgres instead of SQLite |
| `DOCUMENTS_DIR` | `./data/documents` | Lease PDFs |
| `PROPERTY_IMAGES_DIR` | `./data/property-images` | Property photos |
| `QBO_*` | — | QuickBooks Online OAuth (optional) |
| `MAIL_HOST` | _(empty)_ | SMTP host; required to email REPS CSV exports |
| `MAIL_PORT` | `587` | SMTP port |
| `MAIL_USERNAME` / `MAIL_PASSWORD` | — | SMTP credentials |
| `MAIL_FROM` | `noreply@brickwise.local` | From address for REPS emails |

### Supabase (hosted Postgres)

REPS and the rest of the API persist through JDBC. Set `DATABASE_URL` to the Supabase Postgres connection string (do not commit the password):

```bash
export DATABASE_URL='postgresql://postgres:YOUR_PASSWORD@db.YOUR_PROJECT_REF.supabase.co:5432/postgres'
mvn spring-boot:run
```

Prefer the **direct** connection host (`db.<project-ref>.supabase.co`, port `5432`).  
If `sslmode` is omitted for a `*.supabase.co` / `*.supabase.com` host, the backend adds `sslmode=require` automatically.

On startup you should see `Using Postgres database (DATABASE_URL)`.  
`GET /health` returns `"database":"postgres"` when connected via `DATABASE_URL`, or `"database":"sqlite"` when falling back to the local file.

To use local SQLite again, unset `DATABASE_URL` and restart.

Copy `backend/.env.example` for a local template. Cloud / production hosts must inject `DATABASE_URL` as a secret — without it the Java API will not see Supabase data.

## API

Same paths as before under `/api/v1` (properties, residences, mortgages, time-entries, documents, images, dashboard, QBO).

REPS time tracking uses:

- `GET/POST /api/v1/time-entries`
- `GET/PUT/DELETE /api/v1/time-entries/{id}`
- `GET /api/v1/time-entries/export` (CSV download)
- `POST /api/v1/time-entries/export/email` (email CSV + summary; body: `{ "to": "...", "subject?": "...", filters… }`)

On first startup with an empty `time_entries` table, sample REPS rows (and a sample property if needed) are seeded automatically.
