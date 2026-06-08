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

![Architecture diagram](https://www.plantuml.com/plantuml/png/XPJ1Rjim38RlUWeYbtL0Wvvs6R2XMHki0x8XZLjqLzLcRA6ob8bK1cFOku-oQgF9sZWNM-97ykSlEPV4Ed1ihMB35g7uNin_mesfc_aAzsXX4Sh6C9OS0ohr3bQwyv6LnIq3UmX2CbGc266yIyIYP1z8wVI0AslGSTekEc9iuT57L-dGgPIXNMqHPhbf1cRmHoc0qhSxxoGLPelrDoWmx4s9Cz04iZu4KX2bLOFbapmVV917GeUjGtpPQcDVKlr6NgVMbMRzg4bqhJrn7R2uNVNSzPU3wD9gObCIQh4e5ohwJjPcCmXc6wmCzR7-JQc_IfNMnoeDVRU6BBq7qZhvzEcjdyIocs0SOz2vnycCPtu-_mp9nOzmSoTDD_Wh8Z6SNMtkhyzbkVE1ps4HkNcl4YVyq6fC8R6FS4fWogvmXwv2LneygMve9RvAcmtgwk8X64QWtfJUwP5P5iBSHWmTM5yJdCKTcnOXzE97sTbcGGls0U42zjeQAPJ0RK1gWFqJEfgQuzKh1LlhVvzC7FY3gOPvh27-dxZV9Na3GwmvqaYDasTZSZGPC3d0oqi-hNly0jxIJh5jTUYEmInkKTS1wpH5FfHySiTekjy25-Wgm5vTFhl9KNvudhX8JeZdVoA_7J_wJNG4AEex-024fNSLgJM6hGyjGOFQSWDbWaabKQcAEger3Z7KgsJTgQ_eC6X66cSO2rpAMspjFm00)

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

---

# 4. API Workflows

## 4.1 Auth

- `POST /login`: Help user log in
   ```mermaid
   sequenceDiagram

      actor Client
      participant Handler as Gin Handler
      participant Service as Auth Service
      participant UserRepo as User Repository
      participant RedisRepo as Redis Repository
      participant DB as PostgreSQL
      participant Redis as Redis

      %% =========================
      %% LOGIN WORKFLOW
      %% =========================

      Client->>Handler: POST /login (email, password)

      alt Invalid request body
         Handler-->>Client: 400 Bad Request
      else Valid request
         Handler->>Service: ValidateUser(email, password)

         Service->>UserRepo: FindByEmail(email)
         UserRepo->>DB: SELECT user
         DB-->>UserRepo: User record
         UserRepo-->>Service: User

         Note right of Service: Password verified using bcrypt

         alt Invalid credentials
               Service-->>Handler: error
               Handler-->>Client: 401 Unauthorized
         else Valid credentials
               Service-->>Handler: user

               Handler->>Handler: Generate Access Token (15m expiration time)

               Note right of Handler: JWT (HS256)<br>user_id, role, exp, iat

               Handler->>Service: GenerateRefreshToken(userID, role)

               Service->>Service: Generate Refresh Token (7d expiration time)

               Note right of Service: Refresh token for session continuation

               Service->>RedisRepo: StoreToken(refresh:<token>)
               RedisRepo->>Redis: SET with TTL (7d)

               Note right of Redis: key = refresh:<token><br>value = user info

               Redis-->>RedisRepo: OK
               RedisRepo-->>Service: OK

               Service-->>Handler: refresh token
               Handler-->>Client: 200 OK (access + refresh)
         end
      end
   ```

- `POST /logout`: Log out an existing user
   ```mermaid
   sequenceDiagram

   actor Client
   participant Handler as Gin Handler
   participant Service as Auth Service
   participant UserRepo as User Repository
   participant RedisRepo as Redis Repository
   participant DB as PostgreSQL
   participant Redis as Redis

   %% =========================
   %% LOGOUT WORKFLOW
   %% =========================

   Client->>Handler: POST /logout (refreshToken)

   alt Invalid request
      Handler-->>Client: 400 Bad Request
   else Valid request
      Handler->>Service: DeleteOldRefreshToken(token)
      Service->>RedisRepo: DEL refresh:<token>
      RedisRepo->>Redis: DEL key

      Note right of Redis: Logout = token invalidation<br>(stateless JWT + stateful refresh)

      alt Deletion error
         Service-->>Handler: error
         Handler-->>Client: 500 Internal Server Error
      else Success
         Service-->>Handler: OK
         Handler-->>Client: 200 OK
      end
   end
   ```

- `POST /refresh`: Help user refresh their access token, using their refresh token, to help users in case of losing access token, or access token expired
   ```mermaid
   sequenceDiagram

   actor Client
   participant Handler as Gin Handler
   participant Service as Auth Service
   participant UserRepo as User Repository
   participant RedisRepo as Redis Repository
   participant DB as PostgreSQL
   participant Redis as Redis

   %% =========================
   %% REFRESH WORKFLOW
   %% =========================

   Client->>Handler: POST /refresh (refreshToken)

   alt Missing/invalid body
      Handler-->>Client: 400 Bad Request
   else Valid request
      Handler->>Service: GetUserFromRefreshToken(token)

      Service->>RedisRepo: GET refresh:<token>
      RedisRepo->>Redis: GET key

      alt Token not found
         Redis-->>RedisRepo: nil
         RedisRepo-->>Service: error
         Service-->>Handler: error
         Handler-->>Client: 401 Unauthorized
      else Token valid
         Redis-->>RedisRepo: user data
         RedisRepo-->>Service: user data
         Service-->>Handler: userInfo

         Handler->>Service: DeleteOldRefreshToken(token)
         Service->>RedisRepo: DEL refresh:<token>
         RedisRepo->>Redis: DEL key

         Note right of Handler: Refresh token rotation<br>(prevents replay attacks)

         Handler->>Handler: Generate Access Token (15m)

         Handler->>Service: GenerateRefreshToken(userID, role)
         Service->>RedisRepo: StoreToken(new token)
         RedisRepo->>Redis: SET with TTL

         Handler-->>Client: 200 OK (new access + refresh)
      end
   end
   ```

Note: 

- JWT tokens in this application are signed using HS256 algorithm (i.e symmetric signing. The secret is stored securely in server by only provided from .env, not hardcoded). The registered claims contains information on expiration time, and when it is issued, and for private claims, user ID, role (admin/user), and token type (refresh/access) are stored
- For the future, asymmetric signing will be considered and researched into 

## 4.2 Servers

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

- `GET /servers/export`:
   ```mermaid
   sequenceDiagram

   actor Client
   participant Middleware as Auth Middleware
   participant Handler as Gin Handler
   participant Service as Server Service
   participant Repo as Server Repository
   participant DB as PostgreSQL

   Client->>Middleware: GET /servers/export (JWT, query)

   alt Invalid JWT / Not admin
      Middleware-->>Client: 401 Unauthorized
   else Authenticated
      Middleware->>Handler: Forward request

      Handler->>Handler: ParseServersQuery()

      alt Invalid query params
         Handler-->>Client: 400 Bad Request
      else Valid query

         Handler->>Service: GetServers(query)
         Service->>Repo: FindAll(query)

         Repo->>Repo: Build SQL query (filters, sort, pagination)

         Note right of Repo: Filters: status, protocol, name (ILIKE)<br>Safe sort whitelist<br>Prevent SQL injection<br>Count before pagination<br>Offset pagination (performance tradeoff)

         Repo->>DB: SELECT + COUNT
         DB-->>Repo: rows + total
         Repo-->>Service: servers, total
         Service-->>Handler: PaginatedServers

         Handler->>Handler: Set headers (text/csv)

         %% ✅ FIXED NOTE (removed problematic symbols)
         Note right of Handler: Content-Disposition<br>attachment filename servers.csv

         Handler->>Handler: Write CSV header

         loop For each server
               Handler->>Handler: Write CSV row
         end

         Handler->>Handler: Flush writer

         alt CSV write error
               Handler-->>Client: 500 Internal Server Error
         else Success
               Handler-->>Client: 200 OK (CSV stream)
         end
      end
   end

   Note over Handler: Large exports are memory-bound<br>Future: stream from DB or async export job
   ```

- `POST /servers/import`: Import servers from a csv file
   ```mermaid
   sequenceDiagram

   actor Client
   participant Middleware as Auth Middleware
   participant Handler as Gin Handler
   participant Service as Server Service
   participant Repo as Server Repository
   participant DB as PostgreSQL

   Client->>Middleware: POST /servers/import (CSV, JWT)

   alt Invalid / missing JWT
      Middleware-->>Client: 401 Unauthorized
   else Authenticated
      Middleware->>Handler: Forward request

      Handler->>Handler: Extract file from form-data

      alt File missing
         Handler-->>Client: 400 Bad Request
      else File present
         Handler->>Handler: Open file stream

         alt Cannot open file
               Handler-->>Client: 500 Internal Server Error
         else File opened

               Handler->>Service: ImportServers(fileReader)

               Service->>Service: csv.ReadAll()
               Note right of Service: Reads entire file into memory<br>Not optimal for large files

               alt Invalid / empty CSV
                  Service-->>Handler: error
                  Handler-->>Client: 500 Internal Server Error
               else Valid CSV

                  Service->>Service: Extract headers

                  loop For each row
                     Service->>Service: mapRow(headers, row)
                     Service->>Service: parseServer(record)

                     Note right of Service: Validation:<br>- name required<br>- status (bool)<br>- IPv4 format<br>- port (0–65535)<br>- default protocol = tcp

                     alt Validation failed
                           Service->>Service: Append to failures[]
                     else Valid record
                           Service->>Repo: Create(server)

                           Repo->>DB: INSERT (ON CONFLICT DO NOTHING)
                           Note right of Repo: Prevent duplicate server names

                           alt Insert error
                              Repo-->>Service: error
                              Service->>Service: Append to failures[]
                           else Duplicate
                              Repo-->>Service: conflict
                              Service->>Service: Append to failures[]
                           else Success
                              Repo-->>Service: created server
                              Service->>Service: Append to successes[]
                           end
                     end
                  end

                  Service-->>Handler: ImportServersResponse

                  Note right of Service: Response includes:<br>- success count<br>- failure count<br>- success list<br>- failure list

                  Handler-->>Client: 200 OK (JSON result)
               end
         end
      end
   end

   Note over Service,Repo: Current limitations:<br>- csv.ReadAll() is memory-heavy<br>- Sequential inserts<br><br>Future improvements:<br>- Stream parsing<br>- Batch inserts<br>- Async processing
   ```

- `POST /servers/report`: Get servers statuses, then send emails to the emails in email list
   - By default, this API will take all servers information within the last 1 day, then take top 10 servers with worst uptime to report in the output, and no email will be sent. However, users can tweak this by adding optional parameters in the body for start and end time for the report, top N servers they want to be shown up in the response/email, and the emails to be sent to
   - If any of the above extra body parameter is provided, checks will be done to validate input data -> Returns 400 if invalid
   - The number of servers that are up and down are retrieved from Postgres, to get the current amount of servers that are up or down (Which probably could be customized for more flexibility in the future)
   - For the uptime calculation, a range query to Elasticsearch is applied to get status log data from `start` to `end` time, then a bucket aggregation is done to calculate uptime for each server, where inside, uptime is calculated by the average of status field (Since it's 0 and 1 for up and down, you can calculate just with avg, for example, is server is up and down 1 time -> Uptime is (0 + 1) / 2 = 0.5 -> 50% uptime), then the buckets are sorted based on uptime, from lowest to highest, then output is limited to top N buckets only
   - A HTML report then is built from the above results, then they are sent to the emails in the given emails list using SMTP. Specifically, I use SMTP's PLAIN auth for authentication, by using a testing gmail account with a Google's app password for the account, send to the mailing server, then the mail server directs those messages to the emails in the email list (Currently I'm using server `smtp.gmail.com:587` for testing, but in the future, I could try for more)
   ```mermaid
   sequenceDiagram

   actor Client
   participant Middleware as Auth Middleware
   participant Handler as Gin Handler
   participant Service as Server Service
   participant Repo as Server Repository
   participant ESRepo as Elasticsearch Repo
   participant DB as PostgreSQL
   participant ES as Elasticsearch
   participant Mailer as Email Utility
   participant SMTP as SMTP Server

   %% =========================
   %% REQUEST
   %% =========================

   Client->>Middleware: POST /servers/report (filters, emails)

   alt Invalid / missing JWT
   Middleware-->>Client: 401 Unauthorized
   else Authenticated
      Middleware->>Handler: Forward request

      alt Invalid request body / params
         Handler-->>Client: 400 Bad Request
      else Valid request

         Handler->>Handler: Parse time range (default last 24h)
         Handler->>Handler: Validate TopN

         Handler->>Service: SendReports(start, end, topN, emails)

         %% =========================
         %% POSTGRES STATS
         %% =========================
         Service->>Repo: GetStats()
         Repo->>DB: COUNT total servers
         Repo->>DB: COUNT where status = up
         DB-->>Repo: total, up
         Repo-->>Service: total, up, down

         %% =========================
         %% ELASTICSEARCH AGGREGATION
         %% =========================
         Service->>ESRepo: GetDailyUptime(start, end, topN)

         ESRepo->>ES: Query logs (time range + aggregation)

         Note right of ESRepo: Aggregation per server<br>Average status field<br>Top N lowest uptime

         ES-->>ESRepo: Aggregated buckets
         ESRepo-->>Service: uptime map

         %% =========================
         %% BUILD REPORT
         %% =========================
         Service->>Service: Build Report object

         alt No email provided
               Service-->>Handler: Report
               Handler-->>Client: 200 OK (report only)

         else Emails provided

               Service->>Service: Build HTML content

               Note right of Service: HTML report with stats<br>and uptime summary

               Service->>Mailer: Send(to, subject, html)

               Mailer->>SMTP: Send email

               alt Email sending failed
                  Mailer-->>Service: error
                  Service-->>Handler: error
                  Handler-->>Client: 500 Internal Server Error
               else Success
                  SMTP-->>Mailer: OK
                  Mailer-->>Service: OK
                  Service-->>Handler: Report
                  Handler-->>Client: 200 OK (report + email sent)
               end
         end
      end
   end

   %% =========================
   %% GLOBAL NOTES
   %% =========================
   Note over Service,Mailer: Email sending is synchronous<br>Future improvement: async via queue

   Note over ESRepo,ES: Uses aggregation query<br>Efficient for large log datasets
   ```

## 4.3 Users

- `/POST users`: Create an user
- `/GET users`: Get list of users. Normally I would use it for the daily report to the servers administrator, by getting all users, then send email to them about the servers statuses for each 24 hours (Which is the job of the email worker)

---

# 5. API scopes

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

# 6. Data Model

## 6.1 PostgreSQL

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

## 6.2 Elasticsearch

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

# 🔧 7. Design Decisions

## 7.1 Server as Service Endpoint (IP + Port)

Instead of modeling a “server” as a physical machine, SMS-X models it as a **network-accessible service endpoint (IP + Port + Protocol)**.

### Rationale:

* A single machine can host multiple services (HTTP, SSH, FTP, etc.)
* A machine may be reachable (ICMP ping succeeds) but critical services may be down
* Monitoring at service-level provides **more granular and actionable insights**

### Tradeoff:

* Increases number of monitored entities
* Requires clearer definition of “server” in system context

---

## 7.2 Elasticsearch for High-Volume Time-Series Logs

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

## 7.3 Worker-Based Health Checking

Health checking is implemented as a **separate worker process**

### Benefits:

* Decouples monitoring from API
* Enables independent scaling
* Prevents API latency from being affected by network operations

### Design Consideration:

* Worker uses concurrency (worker pool pattern)
* Future improvement: Backpressure control

---

## 7.4 Modular Monolith Architecture

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

## 7.5 Pagination Strategy (Offset-Based)

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

## 7.6 Synchronous vs Asynchronous Processing

Current system performs most operations synchronously.

### Limitation:

* Heavy operations (reporting, email, import/export) can block requests

### Future Direction:

* Introduce asynchronous processing via queue (Kafka/Redis Stream)
* Convert heavy operations into background jobs

---

# 8. Performance Considerations

## 8.1 High Throughput vs Low Latency

The system prioritizes **throughput over latency**, especially for:

* health checks
* log ingestion
* reporting

---

## 8.2 Bounded Concurrency with Worker Pool

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

## 8.3 Database Performance

### PostgreSQL:

* Stores only current state (not logs)
* Indexed on frequently queried fields (name, status)

### Elasticsearch:

* Handles heavy aggregation workloads
* Reduces pressure on relational DB

---

## 8.4 Batch Processing Opportunities

Current implementation processes many operations individually.

### Improvements:

* Batch DB updates
* Use Elasticsearch Bulk API
* Reduce network overhead

---

## 8.5 Caching (Planned)

Redis can be used for:

* Frequently accessed server lists
* Report results
* Metadata (counts)

### Benefit:

* Reduce repeated expensive queries

---

## 8.6 Heavy Operation Handling

Operations like:

* reporting
* CSV import/export
* email sending

Should be:

* moved to background jobs
* processed asynchronously

---

# 9. Authentication & Security

- JWT-based authentication
- Access + refresh token mechanism
- Refresh tokens stored in Redis with TTL
- Passwords hashed using bcrypt

--

# 10. Testing

Unit testings were done using different methods for each component of the codebase:
   - Handler: For HTTP testing, `httptest` package from `net/http` is used. And testify is used for test assertions
   - Repository: Running mock containers using `testcontainers-go` package, to create mock DB for DB operations testings
   - Functionalities: Just do normal testing, using `*testing.T` type for showing failed result

Also, mock codes are used, for components that need extra dependencies. While it's hard to fully cover every cases in the codebase, I'm trying to get as much as I could

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