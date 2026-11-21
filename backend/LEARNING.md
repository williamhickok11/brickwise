# Brickwise Java learning notes (for C# /.NET developers)

This backend was converted from Go to Java (Spring Boot 21) as a pet project for learning Java. Read this alongside the code — package-level Javadoc also calls out .NET analogies.

## Mental model map

| Concept you know (.NET) | Here (Java / Spring) |
|---|---|
| `Program.cs` + minimal hosting | `BrickwiseApplication` + `SpringApplication.run` |
| `appsettings.json` / env vars | `application.yml` + `BrickwiseProperties` (`@ConfigurationProperties`) |
| `services.AddScoped<T>()` | `@Service` / `@Repository` / `@Component` + constructor injection |
| `[ApiController]` + `[Route]` | `@RestController` + `@RequestMapping` |
| `IActionResult` / `Results.Ok` | return the object (Jackson serializes) or `ResponseEntity` |
| Exception middleware / filters | `@RestControllerAdvice` (`GlobalExceptionHandler`) |
| Dapper | `JdbcTemplate` + `RowMapper` |
| EF Core / Npgsql | Not used — explicit SQL keeps the learning path close to the Go version |
| `record` / POCOs | Java classes + getters/setters (or `record` for immutable DTOs) |
| `T?` / nullable reference types | `Optional<T>` for “maybe missing”; boxed types (`Integer`) for nullables |
| `HttpClient` | `java.net.http.HttpClient` (used in `QboClient`) |
| `IOptions<T>` | `@ConfigurationProperties` + `@EnableConfigurationProperties` |

## Layering (same as the old Go API)

```
web (controllers) → service → repository → JDBC / files
```

Routes and JSON field names stay snake_case so the Vue frontend does not need a rewrite.

## Domain vocabulary (unchanged)

- **Property** — a building/parcel you manage
- **Residence** — tenant / occupant history for a property (not a separate “Tenant” type yet)
- **Document** — one lease PDF per residence (file on disk + metadata row)
- **Leases as first-class entities** — not built yet; caller will pick later features

## How to run

```bash
cd backend
mvn spring-boot:run
```

Health: `GET http://localhost:8080/health`  
API base: `http://localhost:8080/api/v1`

Leave `DATABASE_URL` unset for SQLite (`./brickwise.db`). Set it for Postgres/Supabase.

## Where to start reading

1. `BrickwiseApplication.java` — boot
2. `config/DatabaseConfig.java` — SQLite vs Postgres
3. `web/PropertyController.java` — a full CRUD controller
4. `service/PropertyService.java` + `repository/PropertyRepository.java` — business vs data
5. `error/GlobalExceptionHandler.java` — consistent `{code,message}` errors
6. `qbo/QboClient.java` — external HTTP without a third-party SDK

The previous Go implementation is preserved under `../backend-go/` for side-by-side comparison.
