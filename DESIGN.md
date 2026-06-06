# SMS-X System Design Document

---

# 1. System Overview

SMS-X is a backend system designed to manage and monitor the operational status of up to 10,000 servers.

Instead of manually checking server health, users interact with SMS-X via APIs to:
- Register and manage servers
- Monitor server health status
- Query server information
- Generate uptime reports
- Receive automated email notifications

The system is designed as a **modular monolith**, separating concerns into multiple components while remaining deployable as a single system. This allows for future evolution into a distributed microservices architecture.

---

# 2. Goals & Requirements

## 2.1 Functional Requirements

- Manage servers (CRUD operations)
- Periodically check server status
- Store and query server status
- Support filtering, pagination, and sorting
- Import/export server data via CSV
- Generate reports:
  - Total servers
  - Servers up/down
  - Uptime percentage
- Send daily email reports
- Provide authentication and authorization

## 2.2 Non-Functional Requirements

- Handle up to 10,000 servers
- Support high-frequency health checks (every 5 seconds)
- Ensure secure authentication using JWT
- Use appropriate storage technologies:
  - PostgreSQL (primary data)
  - Elasticsearch (logs & analytics)
  - Redis (token storage)
- Provide OpenAPI documentation
- Maintain modular and scalable architecture

---

# 3. High-Level Architecture

The system consists of the following components:

## 3.1 API Service
- Handles all client requests
- Implements business logic
- Provides REST APIs
- Communicates with PostgreSQL, Redis, and Elasticsearch

## 3.2 Server Checker Worker
- Runs continuously (every 5 seconds)
- Performs health checks on all servers
- Updates server status
- Writes logs to Elasticsearch

## 3.3 Email Worker
- Sends periodic reports (daily)
- Uses aggregated data from PostgreSQL and Elasticsearch

## 3.4 Simulation Service (Testing Only)
- Simulates 10,000 servers locally
- Used for load testing and validation

---

# 4. System Workflow

## 4.1 Server Health Check Flow

1. Worker retrieves all servers from PostgreSQL
2. Worker performs TCP connection checks (IP + Port)
3. For each server:
   - Update status in PostgreSQL
   - Send log entry to Elasticsearch
4. Repeat every 5 seconds

## 4.2 Reporting Flow

1. API receives report request
2. Query PostgreSQL for:
   - Total servers
   - Current status (up/down)
3. Query Elasticsearch for:
   - Historical uptime data
4. Aggregate results
5. Return response
6. Optionally send email report

## 4.3 Authentication Flow

1. User logs in with email/password
2. System validates credentials
3. Returns:
   - Access token (JWT)
   - Refresh token (stored in Redis)
4. Access token used for API requests
5. Refresh token used to renew access tokens

---

# 5. API Workflows

## Auth

- `POST /login`: Help user log in
   - Check if the body parameter got all required fields (email, password) -> 400 if invalid
   - Find if the user with provided username and password exist in Postgres, then get the role of the user as well (Could be admin/user) -> 401 if not exist, or if user exist, but password is wrong. The check is done by hash the input password using bcrypt function, then check if matches the hashed password stored in Postgres
   - Generate JWT access token (TTL 15 mins) and refresh token (TTL 7 days) for user, store refresh token in Redis (with same TTL of 7 days), then return in the response

- `POST /logout`: Log out an existing user
   - Check if the body is provided with refresh token -> 400 is not provided
   - Delete refresh token in Redis, if exist

- `POST /refresh`: Help user refresh their access token, using their refresh token, to help users in case of losing access token, or access token expired
   - Check if user provide any refresh token -> 400 if not provided
   - Go to Redis cache to check if this refresh token actually exist -> 401 if this refresh token is not in Redis
   - Delete the provided refresh token, for security purposes
   - Generate new access and refresh token to the user

Note: 

- JWT tokens in this application are signed using HS256 algorithm (i.e symmetric signing. The secret is stored securely in server by only provided from .env, not hardcoded). The registered claims contains information on expiration time, and when it is issued, and for private claims, user ID, role (admin/user), and token type (refresh/access) are stored
- For the future, asymmetric signing will be considered and researched into 

## Servers

- `GET /servers`: Retrieve servers information
   - Extract search query parameters from the URL. The parameters are for searching (by status, by server's protocol, and/or by server's name), pagination (Currently still using offset pagination), and sorting (by which field, which order) -> 400 if the parameters are fed with invalid value (E.g. offset with string and not number)
   - If some criterias above are not provided, there will be default values for it
   - Get the servers based on search, pagination, and ordering criteria, then returns back to the user

- `GET /servers/:id`: Retrieve a server's information
   - Check if the id in the path parameter is a valid number -> Returns 400 if id is negative number, or not a number
   - Try retrieve that server in Postgres -> Returns 404 if not exist, otherwise, return the server info back to user in the response

- `POST /servers`: Create a server
   - User send the request with body content for server creation. The system checks if all necessary info (name, status, IP address, port, protocol) are provided, if not provided or invalid data -> Returns 400
   - The system try to create a record of this server. If the server already exist or some other error happens -> Returns 500, alongside with reason message (E.g.: Server already existed)

- `DELETE /servers/:id`: Delete a server
   - Check if given id is valid -> Returns 400 if invalid or negative number
   - Check if user exist in Postgres -> Returns 404 if not exist
   - Delete the record in Postgres, then returns 204

- `PATCH /servers/:id`: Edit a server
   - Check if given id or request body is valid -> Returns 400 if either is invalid
   - Check if the server with given ID exist -> 404 if not
   - Update the fields in the record in Postgres, then returns 200 with the newly updated record in the response's body

- `GET /servers/export`: Export servers info into a csv file, then sends that back to user. Working the same as `GET /servers` (Yes, even the pagination, sorting, etc.), but with extra of writing out to a csv file, then sends that to the user as part of the response body. Each line, from the header line of the csv, to each line of each server's info are written, and if something's wrong in between each write, it will return 500 as response code

- `POST /servers/import`: Import servers from a csv file
   - Check if the file is included in the body -> 400 if not yet
   - Try opening the file -> 500 if it can't be opened
   - For each line of the file (exclude header), they are validated (See if the line contains required fields, and if the values in correct type), then created in Postgres. If any of the steps above failed, the line is skipped, and failure count is incremented. Otherwise, success count is incremented
   - All the success and failure counts, together with list of successfully/failed inserted rows are returned in the response's body
   - Close the input csv file after (It is a defered function)

- `POST /servers/report`: Get servers statuses, then send emails to the emails in email list
   - By default, this API will take all servers information within the last 1 day, then take top 10 servers with worst uptime to report in the output, and no email will be sent. However, users can tweak this by adding optional parameters in the body for start and end time for the report, top N servers they want to be shown up in the response/email, and the emails to be sent to
   - If any of the above extra body parameter is provided, checks will be done to validate input data -> Returns 400 if invalid
   - The number of servers that are up and down are retrieved from Postgres, to get the current amount of servers that are up or down (Which probably could be customized for more flexibility in the future)
   - For the uptime calculation, a range query to Elasticsearch is applied to get status log data from `start` to `end` time, then a bucket aggregation is done to calculate uptime for each server, where inside, uptime is calculated by the average of status field (Since it's 0 and 1 for up and down, you can calculate just with avg, for example, is server is up and down 1 time -> Uptime is (0 + 1) / 2 = 0.5 -> 50% uptime), then the buckets are sorted based on uptime, from lowest to highest, then output is limited to top N buckets only
   - A HTML report then is built from the above results, then they are sent to the emails in the given emails list using SMTP. Specifically, I use SMTP's PLAIN auth for authentication, by using a testing gmail account with a Google's app password for the account, send to the mailing server, then the mail server directs those messages to the emails in the email list (Currently I'm using server `smtp.gmail.com:587` for testing, but in the future, I could try for more)

## Users

- `/POST users`: Create an user
- `/GET users`: Get list of users. Normally I would use it for the daily report to the servers administrator, by getting all users, then send email to them about the servers statuses for each 24 hours (Which is the job of the email worker)

---

# 6. API scopes

APIs are divided into two scopes, based on user's roles: `admin` and `user`. For APIs that involve read, and not having too much side effects, those would be APIs for everyone. Otherwise, they are for only users with `admin` role. Of course, all of these APIs, aside from the authentication/authorization APIs, needs JWT token for authentication (i.e authentication/authorization APIs have no scope, as they are public APIs). Anyway, here are the scopes of the APIs:

- For both user and admin roles:
   - `GET /servers`
   - `GET /servers/:id`
   - `GET /users`

- Only admin:
   - `POST /servers`
   - `PATCH /servers/:id`
   - `DELETE /servers/:id`
   - `POST /servers/import`
   - `POST /users`
   - `GET /servers/export`
   - `POST /servers/report`

---

# 7. Data Model

## 7.1 PostgreSQL

### Users
- ID
- Name
- Email (unique)
- Password (bcrypt hashed)
- Role (admin/user)
- CreatedAt

### Servers
Represents a monitorable service endpoint.

- ID
- Name (unique)
- Status (boolean)
- IPv4 Address
- Port
- Protocol
- CreatedAt
- LastUpdated

### Memberships
- Many-to-many relationship between Users and Servers
- Designed for future RBAC expansion
- Not utilized yet in current version

---

## 7.2 Elasticsearch

Used for high-volume log storage.

### Data Stream
- Pattern: `server-status*`

### Fields
- `@timestamp` (date)
- `server_id` (long)
- `status` (boolean)

### Purpose
- Store health check logs
- Support aggregation queries (uptime calculation)

---

# 🔧 8. Design Decisions (Expanded)

## 8.1 Server as Service Endpoint (IP + Port)

Instead of modeling a “server” as a physical machine, SMS-X models it as a **network-accessible service endpoint (IP + Port + Protocol)**.

### Rationale:

* A single machine can host multiple services (HTTP, SSH, FTP, etc.)
* A machine may be reachable (ICMP ping succeeds) but critical services may be down
* Monitoring at service-level provides **more granular and actionable insights**

### Tradeoff:

* Increases number of monitored entities
* Requires clearer definition of “server” in system context

---

## 8.2 Elasticsearch for High-Volume Time-Series Logs

Health checks generate extremely high write throughput:

* 10,000 servers × every 5 seconds
  → ~2,000 writes/sec
  → ~170M records/day

### Decision:

Use **Elasticsearch Data Streams** instead of PostgreSQL

### Why:

* Optimized for append-only time-series data
* Supports horizontal scaling
* Native aggregation support
* Built-in lifecycle management (ILM)

### Tradeoff:

* Eventual consistency
* More complex setup vs relational DB

---

## 8.3 Worker-Based Health Checking

Health checking is implemented as a **separate worker process**

### Benefits:

* Decouples monitoring from API
* Enables independent scaling
* Prevents API latency from being affected by network operations

### Design Consideration:

* Worker uses concurrency (worker pool pattern)
* Future improvement: Backpressure control

---

## 8.4 Modular Monolith Architecture

The system is structured as a **modular monolith**:

* Single codebase
* Multiple executables (API, workers)

### Benefits:

* Easier to develop and deploy
* Clear separation of concerns
* Avoids premature microservices complexity

### Future Path:

* Components can be extracted into microservices:

  * API service
  * Health check service
  * Reporting service

---

## 8.5 Pagination Strategy (Offset-Based)

Offset pagination is used for simplicity.

### Pros:

* Easy to implement
* Works well for small datasets

### Cons:

* Inefficient for large offsets (O(n) scan)
* Inconsistent results if data changes frequently

### Future Improvement:

* Cursor-based pagination for better performance

---

## 8.6 Synchronous vs Asynchronous Processing

Current system performs most operations synchronously.

### Limitation:

* Heavy operations (reporting, email, import/export) can block requests

### Future Direction:

* Introduce asynchronous processing via queue (Kafka/Redis Stream)
* Convert heavy operations into background jobs

---

# 9. Performance Considerations

## 9.1 High Throughput vs Low Latency

The system prioritizes **throughput over latency**, especially for:

* health checks
* log ingestion
* reporting

---

## 9.2 Bounded Concurrency with Worker Pool

The server checker worker uses a worker pool pattern to control concurrency when performing health checks.

### Implementation

- A fixed number of workers (e.g., 200) are spawned
- Servers are distributed as jobs through a buffered channel
- Workers continuously pull jobs and perform health checks
- Results are collected through a result channel

### Rationale

- Prevents spawning unbounded goroutines (e.g., 10,000 concurrent checks)
- Avoids exhausting system resources (CPU, memory, network sockets)
- Provides stable and predictable throughput

### Tradeoffs

- Fixed worker count may not fully utilize system resources under varying load
- Requires tuning based on environment

### Future Improvements

- Dynamic worker scaling based on system load
- Rate limiting per target host
- Retry and backoff strategies for failed checks

---

## 9.3 Database Performance

### PostgreSQL:

* Stores only current state (not logs)
* Indexed on frequently queried fields (name, status)

### Elasticsearch:

* Handles heavy aggregation workloads
* Reduces pressure on relational DB

---

## 9.4 Batch Processing Opportunities

Current implementation processes many operations individually.

### Improvements:

* Batch DB updates
* Use Elasticsearch Bulk API
* Reduce network overhead

---

## 9.5 Caching (Planned)

Redis can be used for:

* Frequently accessed server lists
* Report results
* Metadata (counts)

### Benefit:

* Reduce repeated expensive queries

---

## 9.6 Heavy Operation Handling

Operations like:

* reporting
* CSV import/export
* email sending

Should be:

* moved to background jobs
* processed asynchronously

---

# 10. Authentication & Security

- JWT-based authentication
- Access + refresh token mechanism
- Refresh tokens stored in Redis with TTL
- Passwords hashed using bcrypt

-- 

# 11. Deployment & Environment Strategy

## 11.1 Multi-Environment Configuration

System supports multiple environments:

* Local development (.env)
* Docker-based environment (.env.docker)

### Key Difference:

* Service discovery (localhost vs container names)

---

## 11.2 Docker-Based Infrastructure

Docker Compose is used to provision:

* PostgreSQL
* Redis
* Elasticsearch
* Kibana

### Benefits:

* Reproducible environments
* Easy onboarding
* Isolation of dependencies

---

## 11.3 Application Deployment

Application services are separated from infrastructure:

* `docker-compose.infra.yml` → databases
* `docker-compose.yml` → application

### Makefile simplifies:

* environment startup
* logs
* teardown

---

## 11.4 Elasticsearch Initialization

Custom scripts:

* create index template
* initialize data stream

Ensures:

* consistent mappings
* correct ingestion behavior

---

## 11.5 Limitations

* No CI/CD pipeline yet
* No containerized app images (Go binaries run directly)
* No orchestration (e.g., Kubernetes)

---

# 12. Tradeoffs & Limitations

## 12.1 Limited Test Coverage

* Full unit testing not implemented due to time constraints
* However, architecture supports testability (layered design)

---

## 12.2 Simplified Authorization

* Only supports admin/user roles
* No fine-grained access control

---

## 12.3 Synchronous Processing

* Heavy operations executed inline
* Can impact performance under load

---

## 12.4 Elasticsearch Usage

* Basic setup without:

  * ILM policies
  * shard tuning
  * replication tuning

---

## 12.5 No Backpressure Handling

* Worker may overload system under extreme scale
* No queue or buffering layer

---

## 12.6 Offset Pagination Limitation

* Performance degradation for large datasets

---

## 12.7 Email Delivery Constraints

* Relies on SMTP (Gmail)
* Limited throughput and rate limits

---

# 13. Future Improvements

## 13.1 Event-Driven Architecture (Kafka)

Introduce Kafka to:

* decouple components
* handle spikes in workload
* enable asynchronous processing

### Example:

* health check → Kafka → consumers update DB + ES

---

## 13.2 Change Data Capture (CDC - Debezium)

* Capture DB changes automatically
* Stream updates into Kafka
* Keep systems in sync without tight coupling

---

## 13.3 Pre-Aggregation for Reporting

* Store daily uptime summaries
* Avoid expensive real-time aggregations

---

## 13.4 Background Job System

Convert heavy APIs into async jobs:

* reporting
* import/export
* email sending

---

## 13.5 Horizontal Scaling

* multiple API instances behind load balancer
* distributed workers

---

## 13.6 Advanced Caching

* Redis caching for:

  * queries
  * reports
  * counts

---

## 13.7 Full RBAC System

* Utilize Membership table
* Fine-grained permissions per server

---

## 13.8 Observability

Add:

* Metrics (Prometheus)
* Tracing
* Structured logging
* More advanced servers check modelling (Agent based, ICMP, etc.)

---

## 13.9 Improved Elasticsearch Strategy

* Index Lifecycle Management (ILM)
* Hot-warm architecture
* Optimized shard allocation

---

# 14. Conclusion

SMS-X provides a scalable and modular solution for monitoring server health at scale.

The system balances:
- Simplicity (monolith)
- Scalability (Elasticsearch, workers)
- Extensibility (modular structure)

This design serves as a strong foundation for future evolution into a distributed system. It prioritizes correctness and clarity for an MVP, while outlining a clear path toward high-throughput, distributed architecture.

<!--TODO: Add architecture and sequence diagrams here-->