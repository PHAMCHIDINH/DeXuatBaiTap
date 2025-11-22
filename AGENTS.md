# Repository Guidelines

## Project Structure & Module Organization
- `TimMach_api/`: Gin-based Go API; feature handlers live under `modules/<feature>`, middleware in `middleware/`, configuration in `config/`, and database assets in `db/migrations` with sqlc output in `db/sqlc`.
- `TimMach_client/`: Vite + React + TypeScript client; API wrappers in `src/api`, shared UI in `src/components`, routing in `src/routes`, pages in `src/pages`, and shared types/hooks/utils in their respective folders.
- `ml-python/`: FastAPI service serving the heart-disease model (`main.py`); model artifacts reside in `training/`; Dockerfile used by compose.
- `docs/` and top-level architecture markdown files summarize domain context.

## Build, Test, and Development Commands
- Full stack: `docker-compose up --build` brings up Postgres, ML model, and the Go API.
- API: `cd TimMach_api && go run main.go` (expects DB + ML available). `go test ./...` for unit tests. Migrations via `make goose-up` / `make goose-down`; create new migration with `make goose-create name=add_table`; regenerate sqlc code with `make sqlc`.
- Frontend: `cd TimMach_client && npm install && npm run dev` for local dev; `npm run build` for production bundle; `npm run preview` to serve the build locally.
- ML service: `cd ml-python && pip install -r requirements.txt && uvicorn main:app --reload` for standalone runs (compose also builds it).

## Coding Style & Naming Conventions
- Go: run `gofmt` before committing; use CamelCase for exported identifiers, lowercase package names, and keep handlers/services organized per feature under `modules`. SQL column names remain snake_case to match migrations/sqlc models.
- TypeScript/React: prefer functional components; PascalCase component and page files; camelCase hooks/utilities; centralize HTTP clients in `src/api` and reuse shared types from `src/types`.
- SQL migrations: name files with a timestamp prefix (`YYYYMMDDHHMM_<topic>.sql`) and keep forward-only changes paired with corresponding sqlc regeneration.

## Testing Guidelines
- Go: place `_test.go` files next to implementations; favor table-driven cases and run `go test ./...` (add `-race` when feasible).
- Frontend: no test runner is configured yet; if adding, use Vitest + React Testing Library with colocated files like `Component.test.tsx` near the source.
- ML: add FastAPI/ML unit tests under `ml-python/tests` and validate `/predict` with representative payloads.

## Commit & Pull Request Guidelines
- Commit messages: concise, imperative verbs (e.g., “Add patient stats endpoint”); batch related changes; reference issue IDs when available.
- Pull requests: include summary and rationale, note schema/API changes, list commands/tests executed, attach UI screenshots/gifs for client updates, and call out migration/seed impacts.
- Secrets/config: load via environment variables (`DB_URL`, `JWT_SECRET`, `ML_BASE_URL`, `PORT`); keep `.env` files local and out of version control.***
