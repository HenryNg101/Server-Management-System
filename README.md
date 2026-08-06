# SMS-X — Server Management System

SMS-X manages server endpoints and their monitoring data at scale. It provides server inventory management, agent-based metrics ingestion, CSV import/export, uptime reporting, and role-based access through independently deployable Go services.

## Components

- **Services:** authentication, users, servers, agents, jobs, and monitoring.
- **Workers:** metrics ingestion, file import processing, and scheduled reporting.
- **Infrastructure:** PostgreSQL, Redis, Elasticsearch, Kafka, MinIO, Traefik, and Kibana.
- **Frontend:** a React/Vite operations console in [`frontend/`](frontend/).

## Prerequisites

- Docker Desktop with Docker Compose v2
- GNU Make (or run the equivalent `docker compose` commands)
- Go 1.25+ only when running tools such as migrations or Swagger generation locally

## Configuration and startup

Create a local `.env.docker` file in the repository root. It is intentionally ignored by Git; use the following development template and replace every placeholder before starting the stack:

```dotenv
# Docker networking (make env_init writes the active host IP on Windows/WSL)
HOST=127.0.0.1

# PostgreSQL
POSTGRES_USER=postgres
POSTGRES_PASSWORD=change-me
POSTGRES_DB=postgres
POSTGRES_HOST=postgres
POSTGRES_PORT=5432

# Elasticsearch and Kibana
ELASTIC_USER=elastic
ELASTIC_PASSWORD=change-me
KIBANA_PASSWORD=change-me
ELASTIC_HOST=http://elasticsearch
ELASTIC_PORT=9200
ELASTIC_STATUS_DATA_STREAM_SOURCE=server-status
ELASTIC_METRICS_DATA_STREAM_SOURCE=server-metrics

# Redis
REDIS_HOST=redis
REDIS_PORT=6379
REDIS_PASSWORD=change-me

# Kafka
KAFKA_HOST=kafka
KAFKA_PORT=9092

# MinIO
MINIO_ENDPOINT=minio:9000
MINIO_PUBLIC_ENDPOINT=host.docker.internal:9000
MINIO_ACCESS_KEY=change-me
MINIO_SECRET_KEY=change-me
MINIO_USE_SSL=false

# SMTP mail delivery
MAIL_SERVER=smtp.example.com
MAIL_PORT=587
MAIL_USER=notifications@example.com
MAIL_PASSWORD=change-me
MAIL_FROM_USER=notifications@example.com

# Authentication
JWT_SECRET=replace-with-a-long-random-secret
```

Start the development stack:

```bash
make dev.up ENV_FILE=.env.docker
```

Stop it:

```bash
make dev.down
```

`make dev.up` applies pending PostgreSQL migrations through the `postgres-migrate` Compose service before starting dependent application services. Use the service-specific Swagger documentation under `services/<service>/docs/` to explore APIs.

## PostgreSQL migrations

Migration SQL files live in [`migrations/`](migrations/). The system uses [`golang-migrate`](https://github.com/golang-migrate/migrate).

### Run through Docker Compose

When the database container is running, rerun all pending migrations with:

```bash
ENV_FILE=.env.docker docker compose --env-file .env.docker -f docker-compose.infra.yml run --rm postgres-migrate
```

### Run locally with the migrate CLI

Install the PostgreSQL-enabled CLI once:

```bash
go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
```

Set a connection string that reaches the forwarded PostgreSQL port. With the default development configuration, replace `<HOST>` with the value written to `HOST` in `.env.docker` by `make env_init`.

```bash
export DATABASE_URL='postgres://postgres:secret@<HOST>:5432/postgres?sslmode=disable'
```

On PowerShell:

```powershell
$env:DATABASE_URL = 'postgres://postgres:secret@<HOST>:5432/postgres?sslmode=disable'
```

Apply all pending migrations:

```bash
migrate -path migrations -database "$DATABASE_URL" up
```

Roll back only the most recently applied migration:

```bash
migrate -path migrations -database "$DATABASE_URL" down 1
```

Inspect the current migration version:

```bash
migrate -path migrations -database "$DATABASE_URL" version
```

`down` without a number rolls back every migration and is intentionally not recommended for normal development. If a migration fails and leaves the database marked dirty, inspect and correct the migration before using `migrate force <version>`.

### Development users

The initial database contains no users, so authentication cannot bootstrap itself through the API. After migrations have completed, seed local-only accounts:

```bash
psql "$DATABASE_URL" -f scripts/postgres/seed-dev-users.sql
```

This creates these accounts only if they do not already exist:

| Role | Email | Password |
| --- | --- | --- |
| Administrator | `admin@example.com` | `password` |
| User | `user@example.com` | `password` |

The seed script is for local development only. Change or remove these accounts before using any non-local database.

## Authentication and access

`POST /auth/login` returns short-lived access and refresh JWTs. Protected APIs require the access token as `Authorization: Bearer <token>`; refresh tokens can be rotated through `POST /auth/refresh` and invalidated through `POST /auth/logout`.

The backend remains the authorization authority. In the current services, all authenticated users can work with the server records they own; administrators can access all server records and the admin-only user, import-job, and reporting operations.

## API documentation

Each service owns its API definition:

- `services/auth-service/docs/`
- `services/users-service/docs/`
- `services/servers-service/docs/`
- `services/agents-service/docs/`
- `services/jobs-service/docs/`
- `services/monitoring-service/docs/`

To regenerate one service’s Swagger files, run this from that service directory:

```bash
swag init -g ./internal/app/api.go
```

## Useful commands

```bash
# Build the currently selected development images serially to reduce Docker resource pressure
make dev.build ENV_FILE=.env.docker

# Follow logs for one service
make prod.logs SERVICE=servers-service

# Start the sample monitoring-agent client stack
make client.up ENV_FILE_CLIENT=.env.client
```
