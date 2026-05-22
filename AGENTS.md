# Sales Insight — Agent Quick-Start

## Project Overview

Monorepo: Go backend (API + sync engine) + React/Vite frontend. Docker Compose is the primary dev environment. PostgreSQL 15 is the only external dependency.

**Key docs:** `docs/context/ESTADO_ACTUAL.md` (current state), `docs/context/ARCHITECTURE.md` (full architecture). Read `ESTADO_ACTUAL.md` before starting any work.

## Repository Structure

```
backend/
  cmd/api/main.go          # Only real entrypoint (API HTTP + scheduler)
  cmd/worker/              # EMPTY — placeholder for future extraction
  internal/
    platform/              # Infra: config, db, logger, scheduler, security
    merchant/              # Complete: model, service, handler, sqlc repo + tests
    auth/                  # Complete: OAuth/JWT service, handler, middleware
    clover/                # Complete: REST client, OAuth client, rate limiter, DTOs
    sync/                  # Complete: ETL engine (orchestrator + per-entity extractors)
    token_cache/           # In-memory token cache (RWMutex)
    orders/                # Model + interface + STUB repo
    products/              # Model + interface + STUB repo
    employees/             # Model + interface + STUB repo
    payments/              # Model + interface + STUB repo
    analytics/             # EMPTY (Phase 5)
    dashboard/             # EMPTY (Phase 5)
  migrations/              # golang-migrate SQL files
  sqlc.yaml                # One entry per feature package
frontend/                  # Vite + React + Tailwind starter (no app code yet)
```

## Essential Commands

```bash
# Full stack (requires backend/.env with real/semi-real values)
docker compose up -d --build

# Just database (needed before running backend locally)
make db-up                  # or: docker compose up -d db

# Run migrations
make migrate-up             # or: migrate -path backend/migrations -database "$DATABASE_URL" up

# Backend only (local)
cd backend && go run ./cmd/api

# Build backend binary
cd backend && go build -o main ./cmd/api

# Run all backend tests
cd backend && go test ./...

# Frontend dev server
cd frontend && npm run dev

# Frontend build
cd frontend && npm run build
```

## Development Setup (Local, non-Docker)

1. **Database:** `make db-up` starts PostgreSQL 15 on port 5432.
2. **Env:** Copy `backend/.env.example` → `backend/.env` and fill in values.
   - `ENCRYPTION_KEY` must be exactly 32 characters.
   - `JWT_SECRET` is required.
   - `CLOVER_CLIENT_ID` / `CLOVER_CLIENT_SECRET` are needed for OAuth/sync.
3. **Migrations:** Install `golang-migrate` CLI, then run migrations against `postgres://sales_insight:sales_insight_secret@localhost:5432/sales_insight?sslmode=disable`.
4. **Backend:** `cd backend && go run ./cmd/api`
5. **Frontend:** `cd frontend && npm install && npm run dev`

## Backend Architecture Rules

- **Feature-based layout:** Each domain lives in `internal/{feature}/` with `model.go`, `repository.go` (interface), `service.go`, `handler.go`, and optionally `sqlc/`.
- **No ORM:** Use [sqlc](https://docs.sqlc.dev/). Each feature that needs DB access has its own `sqlc.yaml` entrypoint. Run `sqlc generate` from `backend/`.
- **Manual DI only:** All wiring happens in `cmd/api/main.go`. Do not introduce Wire, FX, Dig, or samber/do without discussion.
- **Repository pattern:** `feature/repository.go` defines the interface; `feature/sqlc/repository.go` is the concrete implementation. Business logic depends on interfaces.
- **Tests use hand-written mocks:** See `internal/merchant/service_test.go` for the pattern. No testify/mock is currently used.

## Codegen

- **sqlc:** Edit `internal/{feature}/sqlc/queries.sql`, then run `sqlc generate` in `backend/`. Generated files (`db.go`, `models.go`, `queries.sql.go`) are per-package and must not be hand-edited.
- **sqlc quirks:**
  - `sqlc.yaml` reads schema from `migrations/`, so every package’s generated `models.go` contains structs for **all** tables, not just its own queries.
  - Each package also gets its own `DBTX` interface (`db.go`). This is expected duplication.

## Known State / Deuda Técnica

Read `docs/context/ESTADO_ACTUAL.md` for the full list. The most relevant for agents:

- **Many domain repos are stubs:** `orders/`, `products/`, `employees/`, `payments/` have model + interface + stub repo returning `ErrNotImplemented`. The sync engine compiles and runs because it injects stubs.
- **analytics/ and dashboard/ are empty** (Phase 5).
- **First sync cursor is zero time:** The scheduler will pull all historical data on first run. See `platform/scheduler/cron.go` TODO.
- **Category summary calculation is a nil placeholder** (`sync/transform.go`).

## Environment & Config Gotchas

- **Config loader searches for `.env` recursively** from the working directory up to root. This means `go run ./cmd/api` works whether you run it from `backend/` or `backend/cmd/api/`.
- **Config uses manual field mapping** (not Viper `Unmarshal`) because `mapstructure:",squash"` does not read env vars into nested structs reliably.
- **Graceful shutdown differs by env:** Development calls `e.Close()` (immediate, avoids TIME_WAIT); production calls `e.Shutdown()` (10s timeout).
- **A `.env` file already exists in `backend/`** — verify it has real values before running locally; do not assume it is only a template.