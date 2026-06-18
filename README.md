# SMS-X (Server Management System)

SMS-X is a backend system for managing and monitoring up to 10,000 servers.

It allows users to:
- Register and manage servers (IP + port)
- Monitor server health in near real-time
- Query server status via APIs
- Generate uptime reports
- Receive email notifications

---

# Features

- Server CRUD (Create, Read, Update, Delete)
- Server health checking (every 5 seconds)
- CSV import/export
- JWT authentication (access + refresh tokens)
- Role-based access (admin / user)
- Elasticsearch-based uptime analytics
- Email reporting

---

# Tech Stack

- **Golang**
- **PostgreSQL** (primary database)
- **Elasticsearch** (log storage & analytics)
- **Redis** (refresh token storage)
- **Docker Compose** (infrastructure)
- **Swaggo** (OpenAPI documentation)

---

# Setup

## 1. Environment Variables

Create a `.env` file in the root directory (Name it `.env.docker` if you want to run the whole app in docker):

```
POSTGRES_USER=postgres
POSTGRES_PASSWORD=secret
POSTGRES_DB=postgres
POSTGRES_HOST=localhost
POSTGRES_PORT=5432

ELASTIC_USER=elastic
ELASTIC_PASSWORD=elastic123
KIBANA_PASSWORD=kibana123
ELASTIC_HOST=[http://127.0.0.1](http://127.0.0.1)
ELASTIC_PORT=9200
ELASTIC_STATUS_DATA_STREAM_SOURCE=server-status
ELASTIC_METRICS_DATA_STREAM_SOURCE=server-metric

REDIS_HOST=localhost
REDIS_PORT=6379
REDIS_PASSWORD=redis123

MAIL_SERVER=smtp.gmail.com
MAIL_PORT=587
MAIL_USER=[your_email@gmail.com](mailto:your_email@gmail.com)
MAIL_PASSWORD=your_app_password
MAIL_FROM_USER=[your_email@gmail.com](mailto:your_email@gmail.com)

JWT_SECRET=your_secret
````

---

## 2. How to run

### For dev environment (To setup infrastructure)

Setup:
```bash
make dev.up
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

**Note**: Elasticsearch may take some time to fully start.

### For prod simulation

Start up:
```bash
make prod.up
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
make prod.logs SERVICE=<service-name>
```

---

## 3. Run Applications (For dev mode)

### API Server

```bash
go run ./cmd/api
```

### Server Checker Worker

```bash
go run ./cmd/worker/servers_checker
```

* Runs every 5 seconds
* Checks server health
* Updates PostgreSQL
* Logs to Elasticsearch

### Email Worker

```bash
go run ./cmd/worker/emailer
```

* Sends daily reports via email

---

## 4. (Optional) Simulation

```bash
go run ./cmd/simulation
```

* Simulates opening 10k services on 10k ports on localhost
* For testing only (not production)

---

# Authentication

## Login

```http
POST /auth/login
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

---

## Refresh Token

```http
POST /auth/refresh
```

---

## Logout

```http
POST /auth/logout
```

---

# API Overview

## Public APIs

* `POST /auth/login`
* `POST /auth/refresh`
* `POST /auth/logout`

---

## Authenticated APIs

* `GET /servers`
* `GET /servers/{id}`
* `GET /users`

Supports:

* Filtering
* Pagination (offset-based)
* Sorting

---

## Admin APIs

* `POST /servers`

* `PATCH /servers/{id}`

* `DELETE /servers/{id}`

* `POST /servers/import` (CSV)

* `GET /servers/export` (CSV)

* `POST /servers/report`

---

# Reporting

## Generate Report

```http
POST /servers/report
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

---

# OpenAPI Docs

Generate docs:

```bash
swag init -g ./internal/app/api.go
```

Docs available in `/docs`, and can be viewed and tested at http://localhost:8080/swagger/index.html when running in local

---

# Notes

* Elasticsearch must be fully started before use
* Running in WSL may cause high disk usage issues
* Recommended to run on native OS

---

# Testing

Unit test coverage is limited in this version due to time constraints.

The codebase is structured to support testing with clear separation of:

* handler
* service
* repository

---

# Future Improvements

* Full RBAC system
* Redis caching layer
* Kafka integration
* Improved test coverage
* Horizontal scaling

---

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