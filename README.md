# SMS-X (Server Management System)

SMS-X is a system used for managing and monitoring up to 10,000 servers.

It allows users to:
- Register and manage servers (IP + port)
- Monitor server health in near real-time
- Query server status via APIs
- Generate uptime reports
- Receive email notifications



# Features

- Server CRUD (Create, Read, Update, Delete)
- Server health checking (every 5 seconds)
- CSV import/export
- JWT authentication (access + refresh tokens)
- Role-based access (admin / user)
- Elasticsearch-based uptime analytics
- Email reporting
- Agent-based monitoring to collect metrics from client's servers



# Tech Stack

- **Golang**
- **PostgreSQL** (primary database)
- **Elasticsearch** (log storage & analytics)
- **Redis** (refresh token storage)
- **Docker Compose** (infrastructure)
- **Swaggo** (OpenAPI documentation)



# Setup

## 1. Environment Variables

Create an environment variable file in the root directory:

```
POSTGRES_USER=postgres
POSTGRES_PASSWORD=...
POSTGRES_DB=postgres
POSTGRES_HOST=postgres
POSTGRES_PORT=5432

ELASTIC_USER=elastic
ELASTIC_PASSWORD=...
KIBANA_PASSWORD=...
ELASTIC_HOST=http://elasticsearch
ELASTIC_PORT=9200
ELASTIC_STATUS_DATA_STREAM_SOURCE=server-status
ELASTIC_METRICS_DATA_STREAM_SOURCE=server-metric

REDIS_HOST=redis
REDIS_PORT=6379
REDIS_PASSWORD=...

# Information to setup mail server
MAIL_SERVER=smtp.gmail.com
MAIL_PORT=587
MAIL_USER=[your_email@gmail.com](mailto:your_email@gmail.com)
MAIL_PASSWORD=your_app_password
MAIL_FROM_USER=[your_email@gmail.com](mailto:your_email@gmail.com)

JWT_SECRET=your_secret

# Client stuff. For testing monitoring agent on sample server
CLIENT_AGENT_SERVER_ID=1
CLIENT_AGENT_API_SERVER = http://host.docker.internal:8080

# Generated API key when you 
CLIENT_AGENT_API_KEY = ...
````

## 2. How to run

### For dev environment

Setup:
```bash
make dev.up ENV_FILE="Your env file's path"
````

Tear down:
```bash
make dev.down
````

Services:

* PostgreSQL
* Redis
* Elasticsearch
* Kibana
* API Server: Gin application to expose all APIs for usage
* Server Checker Worker: A worker that runs every 5 seconds to checks server's reachability using ping, then updates information to Elasticsearch, and logs to Elasticsearch for append-only logs
* Email worker: Sends daily reports via email

**Note**: Elasticsearch may take some time to fully start.

### For prod simulation

Prod has quite similar settings, with the exception of exposing minimal service ports and disable debugging information for security purposes

Start up:
```bash
make prod.up ENV_FILE="Your env file's path"
```

Tear down:
```bash
make prod.down
```

Check logs of the whole system:
```bash
make prod.logs
```

Check logs of specific services:
```bash
make prod.logs SERVICE="service name to check logs"
```

## 3. (Optional) Simulation

```bash
go run ./cmd/simulation
```

* Simulates opening 10k services on 10k ports on localhost
* For testing only (not production)

# Authentication

## Login

```http
POST /api/v1/login
```

Body:

```json
{
  "email": "admin@example.com",
  "password": "password"
}
```

Response:

* Access token
* Refresh token

## Refresh Token

```http
POST /api/v1/refresh
```

## Logout

```http
POST /api/v1/logout
```

# API Overview

You can see all documented APIs through Swagger OpenAPI docs by accessing `/swagger/index.html` after the application is up

## Public APIs

* `POST /api/v1/login`
* `POST /api/v1/refresh`
* `POST /api/v1/logout`



## Authenticated APIs

* `GET /api/v1/servers`
* `GET /api/v1/servers/{id}`
* `GET /api/v1/users`

Supports:

* Filtering
* Pagination (offset-based)
* Sorting



## Admin APIs

* `POST /api/v1/servers`

* `PATCH /api/v1/servers/{id}`

* `DELETE /api/v1/servers/{id}`

* `POST /api/v1/servers/import` (CSV)

* `GET /api/v1/servers/export` (CSV)

* `POST /api/v1/servers/report`



# Reporting

## Generate Report

```http
POST /api/v1/servers/report
```

Body:

```json
{
  "start_time": "2026-05-01T00:00:00Z",
  "end_time": "2026-05-02T00:00:00Z",
  "top_n": 10,
  "emails": ["admin@example.com"]
}
```

Returns:

* Total servers
* Servers up/down
* Uptime per server

Optionally sends email report.



# OpenAPI Docs

Generate docs:

```bash
swag init -g ./internal/app/api.go
```

Docs available in `/docs`, and can be viewed and tested at http://localhost:8080/swagger/index.html when running in local



# Notes

* Elasticsearch must be fully started (healthy status) before testing



# Testing

Unit test coverage is limited in this version due to time constraints.

The codebase is structured to support testing with clear separation of:

* handler
* service
* repository



# Future Improvements

* Full RBAC system
* Redis caching layer
* Kafka integration
* Improved test coverage
* Horizontal scaling



# Project Structure

```
/cmd
/internal
/migrations
/docs
docker-compose.*.yml
Makefile
go.mod
go.sum
```