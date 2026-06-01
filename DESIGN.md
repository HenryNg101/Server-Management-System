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

# 5. Data Model

## 5.1 PostgreSQL

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

## 5.2 Elasticsearch

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

# 6. Key Design Decisions

## 6.1 Server = IP + Port

Instead of modeling servers as physical machines, the system models **service endpoints**.

Reason:
- A machine may be up while services are down
- Enables fine-grained monitoring

---

## 6.2 Use of Elasticsearch for Logs

Health checks generate high-frequency data:
- ~10,000 servers
- Every 5 seconds
- Millions of records per day

PostgreSQL is not suitable for this workload.

Solution:
- Use Elasticsearch for append-only logs
- Efficient for aggregation queries

---

## 6.3 Worker-Based Health Checking

A dedicated worker handles health checks.

Benefits:
- Decouples monitoring from API
- Supports concurrency
- Easier to scale independently

---

## 6.4 Modular Monolith Architecture

System is split into:
- API
- Workers
- Simulation

Benefits:
- Simpler deployment
- Clear separation of concerns
- Future-ready for microservices

---

## 6.5 Offset Pagination

Used in server listing APIs.

Pros:
- Simple to implement

Cons:
- Less efficient for large datasets

This will be considered further during the development, based on how this is used

---

# 7. Scaling Considerations

## Challenges

- High-frequency health checks
- Large log volume
- Concurrent network operations

## Solutions

- Use goroutines for concurrent checks
- Use Elasticsearch for log storage
- Keep PostgreSQL for current state only

## Future Scaling

- Introduce Kafka for event streaming
- Horizontal scaling of workers
- Load balancing API servers

---

# 8. Performance Considerations

- Concurrent health checks using goroutines
- Elasticsearch aggregations for efficient queries
- Minimal blocking in API layer

Potential improvements:
- Redis caching
- Cursor-based pagination
- Batch processing optimizations

---

# 9. Alternative Designs Considered

## 9.1 Store Logs in PostgreSQL

Rejected:
- High write volume
- Performance degradation

---

## 9.2 Model Server as IP Only

Rejected:
- Cannot detect service-level failures

---

## 9.3 Cron-Based Scheduling

Alternative:
- Use OS cron jobs

Chosen:
- Long-running worker with ticker

Reason:
- More control
- Simpler integration

---

# 10. Authentication & Security

- JWT-based authentication
- Access + refresh token mechanism
- Refresh tokens stored in Redis with TTL
- Passwords hashed using bcrypt

---

# 11. Reporting

Reports include:
- Total servers
- Servers up/down
- Uptime percentage

Data sources:
- PostgreSQL (current state)
- Elasticsearch (historical data)

Supports:
- API-based retrieval
- Email delivery

---

# 12. Tradeoffs & Limitations

- Limited unit test coverage due to time constraints
- Simplified authorization (admin/user only)
- Basic Elasticsearch configuration
- Cron implemented using ticker
- Membership system not fully utilized
- Redis caching not implemented

---

# 13. Future Improvements

- Full RBAC implementation
- Kafka integration
- Distributed worker system
- Redis caching layer
- Improved test coverage
- Horizontal scaling

---

# 14. Conclusion

SMS-X provides a scalable and modular solution for monitoring server health at scale.

The system balances:
- Simplicity (monolith)
- Scalability (Elasticsearch, workers)
- Extensibility (modular structure)

It serves as a strong foundation for future evolution into a distributed system.

<!--TODO: Add architecture and sequence diagrams here-->