# 🌍 Worlds API

A REST API for managing virtual worlds and users. Includes persistence, event publishing, and scalable design considerations.

## 📂 Project Structure
```
worlds-api/
├── cmd/                # Application entrypoints (Cobra commands)
│   └── start.go        # `start` command to run the API
├── handler/            # HTTP handlers (request/response mapping, validation)
├── services/           # Business logic and orchestration
├── dal/                # Data access layer (Postgres, Redis)
├── models/             # Core domain models and DTOs
├── migrations/         # Database migrations (go-pg based)
├── test/
│   └── end2end/        # End-to-end tests
├── Makefile            # Common tasks (migrate, test, run-local, docker up)
├── docker-compose.yml  # Local development setup
└── Dockerfile          # Multi-stage docker file (builder + runner)
```

## 🚀 Setup / Run Instructions

### Prerequisites
- **Go 1.24** (same version as Dockerfile)
- **Docker & Docker Compose** (for local services)
- **Make** (for build automation)

### Local
```bash
make setup
make deps
make migrate-init
make migrate
make run-local
```

### Docker
```bash
make setup
make deps
make migrate-init
make migrate
make up
```

## 🧪 Testing

Ensure the app is running (`make run-local` or `make up`).

Tests will wipe the database and re-run migrations.

Run:
```bash
make test
```

Tests are currently end-to-end only. They reset the DB before execution.

**Future:** add unit tests and isolated test environments.

## 📖 Design Overview

### Layers

- **Handler** → HTTP layer; request validation, input/output mapping
- **Service** → Business rules and orchestration between layers
- **DAL (Data Access Layer)** → DB persistence in Postgres
- **Models** → DTOs and domain entities, separate from transport and persistence logic

This separation improves maintainability and testing.

### Infrastructure Choices

- **Postgres** → Authoritative persistence layer for worlds and users. Mature, reliable, and easy to extend with schemas
- **Redis** → 
  - Pub/Sub for event publishing (simulated message bus)
  - Tracking user world membership in-memory (SADD, SREM, SCARD)
- **Docker** → Standardized local/dev setup across environments

### Eventing

- Redis Pub/Sub simulates async message queues
- Event publishing happens in separate goroutines to keep API latency low
- Can be mocked in tests since its just an interface

### Middleware

- **Logging** → Structured logs with context fields
- **Metrics** → Planned for future (Prometheus/Grafana)
- **Authentication** → Stubbed, to be implemented later

### Toolset

- **Viper** → Configuration
- **go-playground/validator** → Request param/body validation
- **go-pg/migrations** → Database migrations

## 🔑 Architectural Decisions

### Redis for User World Membership

**Considerations:**

- **Durability:** Redis is not durable — data may be lost on restart. This is acceptable because membership can be re-synced from clients if we implement a periodic keep-alive route, which is useful to have anyway.

- **TTL for Cleanup:** Redis supports TTLs, which makes it well-suited for this ephemeral data. We can set a TTL to 3× the keep-alive interval so old entries are cleaned up automatically. At the moment, no TTL is in place, which would be problematic because user records would persist indefinitely. A proper client connection mechanism should ensure entries are removed when clients disconnect.

- **Memory Growth:** Another concern is Redis growing too quickly if not managed. In production, we would add metrics to monitor memory usage and track key counts. If Redis becomes a bottleneck, we can:
  1. Scale it vertically (larger instance)
  2. Or migrate this functionality to Postgres/MongoDB for more IOPS and persistence

### Active World Definition

The term "active worlds" was not explicitly defined in the specification, so I chose to assign it a semantic meaning: **a world is considered active if `user_count > 0`**.

This is a minimal definition that allows us to return meaningful results without overcomplicating the initial implementation. In the future, we can extend this to support richer client-side filtering, pagination, and more advanced search queries, but for now this provides a clear and simple baseline.

### Redis Pub/Sub for Events

Redis Pub/Sub is used to broadcast events (e.g., world created, world updated) in this implementation. It provides a lightweight way to simulate a message queue.

**Key Considerations:**

- **Durability & Features:** Redis Pub/Sub is not durable and lacks advanced features such as persistence, replay, or consumer groups. For production systems that require stronger delivery guarantees, we would migrate to a more robust message broker like Kafka, NATS, or RabbitMQ.

- **Sufficient for Now:** For this project and test scenarios, Redis Pub/Sub is more than enough. It's simple to set up, has no extra operational cost, and allows us to demonstrate event-driven design without introducing unnecessary complexity.

- **Mockability:** Since the API only depends on an abstract event publisher interface, Redis Pub/Sub can be easily replaced with a mock during testing.

- **Responsiveness:** Event publishing is offloaded to background goroutines, ensuring that API requests return quickly without being blocked by downstream consumers.


## 📡 API Endpoints

### Worlds Management

| Method | Endpoint | Description |
|--------|----------|-------------|
| `POST` | `/worlds` | Create a new world |
| `GET` | `/worlds` | List all worlds |
| `GET` | `/worlds/{id}` | Get world by ID |
| `PUT` | `/worlds/{id}` | Update world details |
| `POST` | `/worlds/{id}/join` | Join a specific world |
| `GET` | `/worlds/my-current` | Get current user's active world |

### User Management

| Method | Endpoint | Description |
|--------|----------|-------------|
| `POST` | `/user/{id}` | Create a new user |

### Base URL
```
http://localhost:8080
```

### Authentication
Most endpoints require authentication (currently stubbing the token exchange process, can implement supabase, firebase, auth0 later)

#### Current Implementation (Stubbed)
The authentication is currently stubbed for development purposes. The system expects a `Bearer` token in the `Authorization` header, but treats the token value as the user ID directly.

**Example:**
```bash
# Using curl
curl -H "Authorization: Bearer user123" http://localhost:8080/worlds/my-current

# Using Go
authHeaders := map[string]string{"Authorization": "Bearer " + userID}
```

**Note:** This is a temporary implementation for development. In production, proper JWT or OAuth tokens should be used instead of direct user IDs.

## Future Work

### Metrics
- Prometheus + Grafana dashboards
- Collect response latency percentiles, error rate by route, throughput

### Redis improvements
- Better Lua scripting (currently hard to debug)
- Consistent, atomic operations for join/leave world

### Testing
- Add unit tests (not just E2E)
- Spin up isolated test env with in-process httptest.Server

### CI/CD
- Run tests before building/pushing
- Containerize migrations and run them in CD

### Performance optimizations
- Cache GET /worlds with pagination
- Use Redis + CDN (e.g. Cloudflare) for high QPS reads

### Setup improvements
- Add setup.sh that waits for Postgres/Redis readiness

### Auth
- Implement real authentication + authorization

### API Documentation
- Swagger/OpenAPI with live Swagger UI
