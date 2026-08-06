# Code structure

SMS-X is now a service-oriented codebase. The primary applications are no longer rooted in the repository-level `internal/` directory; each service and worker owns its executable entrypoint, application wiring, domain code, infrastructure adapters, and generated OpenAPI documentation.

## Top-level directories

- [`services/`](services/) — independently built HTTP services.
  - Every service follows the same broad shape:
    - `cmd/` — service process entrypoint.
    - `internal/app/` — dependency wiring, Gin setup, and route registration.
    - `internal/<domain>/` — handlers, DTOs, services, repositories, and domain helpers.
    - `internal/platform/` — adapters for Postgres, Redis, Elasticsearch, Kafka, MinIO, or SMTP as relevant.
    - `internal/middleware/` and `internal/shared/` — HTTP auth middleware and reusable service-local code.
    - `docs/` — Swagger files generated for that service.
  - `auth-service/` — login, refresh-token rotation, and logout.
  - `users-service/` — administrator-managed users.
  - `servers-service/` — server ownership, CRUD, filtering, and CSV export.
  - `agents-service/` — agent registration, API-key authentication, and metrics ingestion.
  - `jobs-service/` — asynchronous CSV import job creation and job status.
  - `monitoring-service/` — Elasticsearch-backed report aggregation and email delivery.

- [`workers/`](workers/) — asynchronous and scheduled processes, each with its own `cmd/`, `internal/`, and infrastructure adapters.
  - `metrics-consumer/` — consumes agent metric events and writes them to Elasticsearch.
  - `files-import-consumer/` — processes CSV import jobs from Kafka and records results.
  - `cron-scheduler/` — schedules recurring work such as daily reporting.

- [`agents/`](agents/) — installable monitoring agents, separate from the SMS-X control-plane services.

- [`frontend/`](frontend/) — React/Vite management console.

- [`migrations/`](migrations/) — ordered PostgreSQL schema migrations used by `golang-migrate`.

- [`scripts/`](scripts/) — operational scripts grouped by concern:
  - `elasticsearch/`, `kafka/`, and `minio/` provision local infrastructure.
  - `postgres/` contains development database helpers, including `seed-dev-users.sql`.
  - `seed_data/` creates sample CSV inputs for import testing.

- [`cmd/`](cmd/) — legacy monolith and simulation utilities retained for reference and targeted local testing. It is not the production service entrypoint layer. `cmd/simulation/single_server/` remains the sample server used for agent testing.

- [`internal/`](internal/) — legacy monolith application structure retained as a reference/scaffolding base while code has moved into service- and worker-local packages. New production behavior should be implemented in the owning service or worker rather than here.

- [`docs/`](docs/) — legacy, repository-wide generated Swagger output from the monolith. Service-local `docs/` directories are the authoritative API documentation; this directory can be removed once no external workflow still consumes it.

## Infrastructure and project files

- `docker-compose.infra.yml` — shared data stores and one-shot infrastructure setup containers.
- `docker-compose.yml` — service and worker containers plus Traefik routing.
- `docker-compose.dev.yml` — development-only port exposure.
- `docker-compose.client.yml` — sample client server and monitoring-agent stack.
- `Makefile` — environment initialization, serialized image builds, Compose lifecycle commands, and log shortcuts.
- `REQUIREMENT.md` — assignment requirements.
- `DESIGN.md` — earlier design notes; useful background, but not necessarily current implementation truth.
