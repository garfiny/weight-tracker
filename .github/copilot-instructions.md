## Purpose

Short, actionable guidance for AI coding agents working on this repository (small Go HTTP server implemented as a weight-tracking web application).

## Big picture / architecture
- Single binary Go web service. Entry point: `src/main.go`.
- Routing uses `github.com/gorilla/mux`; route registration lives in `src/main.go`.
- Handlers live in the `handlers` package under `src/handlers` (e.g. `src/handlers/handler.go`).
- Domain: weight-tracking web application. Expected high-level flow:
  - Client (web/mobile) -> HTTP API -> handlers -> storage layer -> JSON responses.
  - Handlers should accept and return JSON; responses set `Content-Type: application/json` (see `handlers.HandleRequest`).
  - Storage is not present in the current code; assume either an in-memory store for examples/tests or a small embedded DB (SQLite) if persistence is required.

## Key files to inspect first
- `src/main.go` — router setup, route registrations, server listen (port :8080).
- `src/handlers/handler.go` — example handler; sets Content-Type and encodes a simple JSON map.
- `src/handlers/handler.go` — current example handler; replace or extend with weight-specific handlers (see "API & data model" below).
- `go.mod` — module name `weight-tracker` (local import path is `weight-tracker/src/...`).
- `README.md` — run/build instructions.

## Developer workflows (concrete commands)
- Install/update deps: `go mod tidy` (keeps `go.sum` up to date).
- Run locally: `go run src/main.go` (server starts on :8080).
- Build binary: `go build -o bin/weight-tracker ./src`.
- Run tests (none exist yet): `go test ./...` (add package tests under `src/...`).
- Lint/format: `gofmt -w .` and `go vet ./...`.
- Debugging (local): use Delve if available, e.g. `dlv debug -- ./src`.

## Project-specific conventions & patterns
- Code lives under `src/` (not flat root); follow this layout when adding packages.
- Keep routing registration in `src/main.go`. Handlers should be pure functions that accept `(w http.ResponseWriter, r *http.Request)` and write the response.
- Handlers set `Content-Type` to `application/json` and use `json.NewEncoder(w).Encode(...)` to write responses (see `HandleRequest`).
- Use exported function names for handlers you intend to register (e.g. `HandleRequest`, `NewEndpoint`).
- Module path is `weight-tracker` in `go.mod` — if upstream/publishing changes, update imports to the new module path.

API & data model (repo-specific suggestions)
- This repo represents a weight-tracking app. Recommended minimal API:
  - GET    /weights           -> list weight entries (paginated optional)
  - POST   /weights           -> create a weight entry (JSON body)
  - GET    /weights/{id}      -> fetch single entry
  - PUT    /weights/{id}      -> update entry
  - DELETE /weights/{id}      -> delete entry

- Example Go data model (suggestion for `models/weight.go` or inside `handlers`):
  - type WeightEntry struct {
      ID     string `json:"id"`
      Date   string `json:"date"`   // ISO-8601 date e.g. 2025-10-18
      Weight float64 `json:"weight"`
      Notes  string  `json:"notes,omitempty"`
    }

- Storage assumptions (explicit):
  - Current codebase has no persistence. Reasonable assumptions for next steps:
    1) Start with an in-memory store (map[string]WeightEntry) for fast iteration and tests.
    2) If persistence needed, add SQLite with `github.com/mattn/go-sqlite3` or use a small file-backed JSON store.
  - Document which approach you pick in PRs and add migrations/seed data if using a DB.

## Tests: how to write them (concrete example)
- Use `net/http/httptest` to exercise handlers directly. Pattern:
  - Create `httptest.NewRecorder()` and `http.NewRequest()`.
  - Call handler: `handlers.HandleRequest(rec, req)` or mount router and call `r.ServeHTTP(rec, req)`.
  - Decode response JSON and assert expected fields (e.g., `message == "Hello, World!"`).

  Example test for weight handler (high level):
  - Create `httptest.NewRecorder()` and `http.NewRequest("POST", "/weights", bytes.NewReader(payload))`.
  - Call the handler or mount router and call `r.ServeHTTP(rec, req)`.
  - Assert status code, then decode JSON and assert `weight` and `date` fields match.

## Integration points & dependencies
- External dependency observed: `github.com/gorilla/mux` (router). Add other deps with `go get` and run `go mod tidy`.
- No external services (databases, cloud APIs) are wired in the current code — look for new packages in `src/` when adding integrations.

If adding persistence, common choices:
- SQLite (`github.com/mattn/go-sqlite3`) for a single-file DB and minimal ops.
- PostgreSQL (`github.com/lib/pq`) for production-like setups (update `go.mod` and config).

If adding auth in future, plan routes under `/api/v1/*` and follow same handler patterns.

## Small PR checklist for agents
- Ensure new imports are added to `go.mod` and `go.sum` via `go get` / `go mod tidy`.
- Run `gofmt -w` and `go vet ./...` before submitting changes.
- Add unit tests for handlers using `httptest` when changing behavior.
- Keep route registration centralized in `src/main.go`.

Additional checklist for weight features:
- Add a `models` or `store` package if you introduce persistence.
- Provide a simple seed script or test fixture for weight entries when adding DB support.
- Update `README.md` with API examples and sample curl/snippet responses.

## When in doubt
- Inspect `src/main.go` first to see how routes are wired and which handlers are in use.
- Follow the existing handler pattern (JSON map responses, Content-Type header).

Assumptions made in this file (documented):
- No persistence layer exists in the current repository; instructions above recommend an in-memory store for quick work and tests.
- The project is small and uses `mux` for routing; keep changes minimal and centralized in `src/` unless expanding.

---
If you want, I can also:
- add a starter unit test under `src/handlers/handler_test.go` using `httptest`.
- add a Makefile or GitHub Actions workflow for build/test.
Please tell me which to add next.
