# Learning the Brickwise Java API

This folder is a Java 17 / Spring Boot 3.5 copy of one slice of the Go API. The Go server in `backend/` stays the running app and the spec. This server uses the same SQLite file.

Java 17 is the language level. The compiler is set to 17 in `pom.xml`, so Java 21-only syntax will not compile. Records, text blocks, and `instanceof` pattern matching are already in 17. Virtual threads are not.

## Run it

Install a JDK 17 and point `JAVA_HOME` at it. Homebrew's JDK is not put on your `PATH` for you:

```bash
brew install openjdk@17
export JAVA_HOME="$(brew --prefix openjdk@17)/libexec/openjdk.jdk/Contents/Home"
```

From this directory:

```bash
./mvnw spring-boot:run
```

The first run downloads Maven. The API listens on http://localhost:8081.

```bash
curl http://localhost:8081/health
curl http://localhost:8081/api/v1/properties
```

`./mvnw test` runs the controller test.

The Vue app still talks to the Go API on port 8080. To try this server for a session, change the proxy target in `frontend/vite.config.ts` from `http://localhost:8080` to `http://localhost:8081` and restart `npm run dev`. Only `GET /api/v1/properties` is implemented here, so other screens will 404 until you add them.

## How one request moves

`GET /api/v1/properties`

1. `PropertyController.listProperties` reads the HTTP request and returns the list. Spring turns that list into JSON.
2. `PropertyService.listProperties` is where rules will go. Today it only calls the repository.
3. `PropertyRepository.findAll` runs the SQL with `JdbcTemplate` and maps each row to a `Property` record.

That is the Go stack in `backend/README.md`: handler, service, repository. Copy this package when you add the next resource.

## Names you already know

| C# / ASP.NET | Go in this repo | Spring in this folder |
|---|---|---|
| `Program.cs` | `backend/cmd/server/main.go` | `BrickwiseApplication` |
| `[ApiController]` | `internal/handlers` | `@RestController` |
| `[HttpGet]` | `r.Get(...)` | `@GetMapping` |
| constructor injection | `NewPropertyService(repo)` | constructor plus `@Service` |
| `record` | struct plus json tags | `record` |
| Dapper | `database/sql` in `internal/repository` | `JdbcTemplate` |
| `appsettings.json` | env vars | `src/main/resources/application.yml` |
| `.csproj` | `go.mod` | `pom.xml` |
| `dotnet test` | `go test ./...` | `./mvnw test` |

`pom.xml` is the project file. `./mvnw` is the wrapper, so you do not install Maven yourself. Dependencies download into a local cache the first time you build, the way `go mod download` does.

## JSON shape

The Vue client expects snake_case (`property_type`, `created_at`). Jackson's default is camelCase. `application.yml` sets `spring.jackson.property-naming-strategy: SNAKE_CASE` for the whole app, so the next record you add picks it up too.

SQLite stores `created_at` as `2026-03-30 20:17:18` with no timezone. `PropertyRepository.parseTimestamp` reads that as UTC and Jackson writes `2026-03-30T20:17:18Z`, which is the RFC3339 shape the Go API emits. An empty `value_updated_at` stays JSON `null`.

An empty property list is `[]`. The Go handler returns `null` in that case, because a nil slice encodes that way.

## The database

`application.yml` points at `../backend/brickwise.db`. The Go server creates the tables. Start the Go API once if that file or the `properties` table is missing. This Java process does not change the schema. Once you add `POST`, it will insert rows into that same file, so avoid running a write from both servers at the same moment. SQLite wants a single writer, which is why the pool size is 1.

Run `./mvnw spring-boot:run` from `backend-java`. The SQLite path is relative to the working directory.

## Your next endpoints

Each exercise copies the property package and adds one new idea. The Go files are the spec. Keep the URL and the JSON field names the same, or the Vue client will not understand the response.

1. `GET /api/v1/properties/{id}` — path variables and 404. Spec: `GetProperty` in `backend/internal/handlers/property.go` and `GetByID` in `backend/internal/repository/repository.go`. Add `findById` returning `Optional<Property>`. In the controller, `@GetMapping("/{id}")` and `@PathVariable int id`. When the optional is empty, respond with 404. `ResponseStatusException` with `HttpStatus.NOT_FOUND` is the small version of that.
2. `POST /api/v1/properties` — a request body, validation, and status 201. Spec: `CreateProperty`. The method argument is `@RequestBody CreatePropertyRequest request`, another record. Return `ResponseEntity.status(HttpStatus.CREATED).body(created)`.
3. `PUT /api/v1/properties/{id}` and `DELETE /api/v1/properties/{id}` — the rest of the property resource. Delete returns 204 with an empty body.
4. Residences, time entries, then mortgages. New package, same three types: controller, service, repository. Leave dashboard and QuickBooks on the Go server. OAuth is not the lesson.
5. After those reads and writes feel boring, try Spring Data JPA on a new table. This first slice uses `JdbcTemplate` on purpose: the Go repositories are hand-written SQL, and Hibernate plus SQLite hides the SQL you are trying to learn.

The controller test uses `@WebMvcTest` and `@MockitoBean`. `@WebMvcTest` builds only the web layer, and `@MockitoBean` replaces `PropertyService` so the test never opens SQLite. On Spring Boot 3.3 and earlier the same annotation is `@MockBean`.

## A request, in the order you will type it

For exercise 1, the new controller method looks like this. The body is yours to write.

```java
@GetMapping("/{id}")
public Property getProperty(@PathVariable int id) {
    return propertyService.getProperty(id);
}
```

`@PathVariable` is `[FromRoute]`. `@RequestBody` is `[FromBody]`.
