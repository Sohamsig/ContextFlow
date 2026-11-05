# ContextFlow

### Open-Source Enterprise Context Infrastructure

**Connect your data. Govern your knowledge. Power your AI.**

ContextFlow is an open-source project focused on building reliable infrastructure for connecting enterprise data sources, managing contextual information, and making that information accessible to AI applications through controlled APIs.

The project explores the engineering foundations behind enterprise context infrastructure: data ingestion, persistence, reliability, governance, observability, and permission-aware retrieval.

> **Project status:** Early development — foundational infrastructure and database integration are being established.

---

## Table of Contents

- [Vision](#vision)
- [Core Objectives](#core-objectives)
- [Architecture](#architecture)
- [Technology Stack](#technology-stack)
- [Current Implementation](#current-implementation)
- [Project Structure](#project-structure)
- [Getting Started](#getting-started)
- [Environment Configuration](#environment-configuration)
- [Database Infrastructure](#database-infrastructure)
- [Backend API](#backend-api)
- [Reliability and Engineering Decisions](#reliability-and-engineering-decisions)
- [Development Progress](#development-progress)
- [Known Issues](#known-issues)
- [Security](#security)
- [Roadmap](#roadmap)
- [Contributing](#contributing)
- [License](#license)

---

## Vision

Modern AI applications need more than language models. They need reliable access to relevant, structured, and permission-aware information from multiple enterprise systems.

ContextFlow aims to provide an infrastructure layer between enterprise data sources and AI applications.

The long-term objective is to build a system that can:

- Connect to external data sources.
- Ingest and synchronize data reliably.
- Store contextual records and metadata.
- Track ingestion operations and execution state.
- Apply access controls and governance policies.
- Expose context through APIs.
- Monitor system health and operational failures.

ContextFlow is being developed incrementally, starting with a reliable backend and PostgreSQL foundation.

## Core Objectives

| Area | Objective |
|---|---|
| Data connectivity | Connect external enterprise data sources |
| Data ingestion | Build reliable ingestion and synchronization workflows |
| Context storage | Persist contextual records and metadata |
| Reliability | Support retries, checkpoints, and recovery |
| Governance | Introduce organizations, workspaces, and access policies |
| Context APIs | Provide controlled access to stored context |
| Observability | Monitor pipelines, errors, and system health |

## Architecture

The intended architecture is modular, with an initial Go backend and PostgreSQL persistence layer.

```text
             EXTERNAL DATA SOURCES
             ┌────────────────────┐
             │ PostgreSQL         │
             │ REST APIs          │
             │ Object Storage     │
             └─────────┬──────────┘
                       │
                       ▼
             ┌────────────────────┐
             │ Connector Layer    │
             │ Source Integration │
             └─────────┬──────────┘
                       │
                       ▼
             ┌────────────────────┐
             │ Ingestion Engine   │
             │ Workers / Retries  │
             │ Checkpoints        │
             └─────────┬──────────┘
                       │
                       ▼
             ┌────────────────────┐
             │ Context Storage    │
             │ PostgreSQL         │
             │ Records / Metadata │
             └─────────┬──────────┘
                       │
                       ▼
             ┌────────────────────┐
             │ Governance Layer   │
             │ RBAC / Policies    │
             │ Audit Events       │
             └─────────┬──────────┘
                       │
                       ▼
             ┌────────────────────┐
             │ Context API        │
             │ REST Endpoints     │
             └─────────┬──────────┘
                       │
                       ▼
                AI APPLICATIONS
```

*This represents the target architecture. Not every component has been implemented yet.*

### Architectural approach

The initial implementation follows a modular monolith approach.

This keeps the system easier to develop, test, debug, and operate while preserving clear boundaries between backend components.

The project will introduce additional infrastructure only when justified by actual requirements.

## Technology Stack

| Component | Technology | Purpose |
|---|---|---|
| Backend | Go | API and infrastructure services |
| HTTP Router | Chi | HTTP routing and middleware |
| Database | PostgreSQL 17 | Persistent storage |
| Database Driver | pgx v5 | PostgreSQL connectivity |
| Connection Pool | pgxpool | Database connection management |
| Containerization | Docker | Isolated infrastructure |
| Orchestration | Docker Compose | Local development environment |
| Frontend | React + TypeScript | Planned management dashboard |
| Version Control | Git + GitHub | Source control and collaboration |

## Current Implementation

### 1. Backend foundation

The backend is being developed in Go.

Implemented foundations include:

- Application entry point.
- Environment-based configuration.
- PostgreSQL connection-pool initialization.
- HTTP routing using Chi.
- Request logging middleware.
- Panic recovery middleware.
- Database-aware health endpoint.

### 2. PostgreSQL infrastructure

PostgreSQL 17 Alpine is configured through Docker Compose.

The database configuration includes:

- Dedicated database: `contextflow`
- Database user: `contextflow`
- Container: `contextflow-postgres`
- Internal PostgreSQL port: `5432`
- Persistent Docker volume: `contextflow_postgres_data`
- Container health check using `pg_isready`

The persistent volume allows database data to survive normal container restarts and recreation.

### 3. Configuration management

The backend reads its database connection string from an environment variable.

```go
type Config struct {
    DatabaseURL string
}

func Load() (*Config, error) {
    databaseURL := os.Getenv("DATABASE_URL")

    if databaseURL == "" {
        return nil, fmt.Errorf(
            "DATABASE_URL environment variable is required",
        )
    }

    return &Config{
        DatabaseURL: databaseURL,
    }, nil
}
```

This separates application configuration from source code and avoids hardcoding database credentials.

### 4. Database connection pooling

The PostgreSQL connection pool uses `pgxpool`.

Current configuration:

- Maximum connections: 10
- Minimum connections: 1
- Maximum connection lifetime: 1 hour
- Startup connectivity verification through `Ping`

```go
config.MaxConns = 10
config.MinConns = 1
config.MaxConnLifetime = time.Hour
```

The application verifies database connectivity before starting its HTTP server.

### 5. Health monitoring

The backend exposes:

```http
GET /api/v1/health
```

The endpoint is designed to report:

- Application status.
- Service name.
- Application version.
- Database connectivity.
- UTC timestamp.

Possible response when the database is reachable:

```json
{
  "status": "healthy",
  "service": "contextflow",
  "version": "0.1.0",
  "database": "connected",
  "timestamp": "2026-11-05T15:30:00Z"
}
```

The response above illustrates the intended healthy response format; successful end-to-end execution from the Windows host is still being verified.

---

## Project Structure

Current project organization:

```text
ContextFlow/
│
├── backend/
│   ├── cmd/
│   │   └── server/
│   │       └── main.go
│   │
│   ├── internal/
│   │   ├── config/
│   │   │   └── config.go
│   │   │
│   │   └── database/
│   │       └── postgres.go
│   │
│   ├── go.mod
│   └── go.sum
│
├── deployments/
│
├── docs/
│
├── frontend/
│
├── .env
├── .gitignore
├── docker-compose.yml
├── LICENSE
└── README.md
```

*Some directories are reserved for upcoming implementation.*

## Getting Started

### Prerequisites

Install the following tools:

- Go (compatible with the project's `go.mod` version)
- Docker Desktop
- Docker Compose
- Git
- PowerShell (Windows development environment)

### 1. Clone the repository

```powershell
git clone https://github.com/Sohamsig/ContextFlow.git

cd ContextFlow
```

### 2. Configure environment variables

Create a `.env` file in the project root.

```env
POSTGRES_USER=contextflow
POSTGRES_PASSWORD=replace_with_a_local_secret
POSTGRES_DB=contextflow
```

Use a unique local development password. Do not commit `.env` to Git.

### 3. Start PostgreSQL

```powershell
docker compose up -d
```

Check container status:

```powershell
docker compose ps
```

Wait until PostgreSQL reports:

```text
Up (healthy)
```

### 4. Verify database connectivity inside the container

```powershell
docker compose exec postgres psql -U contextflow -d contextflow -c "SELECT current_database(), current_user, version();"
```

### 5. Configure the backend connection

From the backend directory:

```powershell
$env:DATABASE_URL = "postgres://contextflow:YOUR_PASSWORD@127.0.0.1:5432/contextflow?sslmode=disable"
```

Replace `YOUR_PASSWORD` with the password configured for the PostgreSQL role.

### 6. Start the backend

```powershell
go run ./cmd/server
```

The intended startup message is:

```text
ContextFlow API starting on :8080
```

### 7. Verify the health endpoint

Once the backend starts successfully:

```powershell
Invoke-RestMethod http://localhost:8080/api/v1/health
```

---

## Environment Configuration

| Variable | Description |
|---|---|
| `POSTGRES_USER` | PostgreSQL username |
| `POSTGRES_PASSWORD` | PostgreSQL password |
| `POSTGRES_DB` | PostgreSQL database |
| `DATABASE_URL` | Backend PostgreSQL connection string |

Example `.env.example`:

```env
POSTGRES_USER=contextflow
POSTGRES_PASSWORD=change_me
POSTGRES_DB=contextflow
```

The actual `.env` file must remain local and must never be committed.

## Database Infrastructure

### Docker Compose

The current PostgreSQL service is configured with:

```yaml
services:
  postgres:
    image: postgres:17-alpine
    container_name: contextflow-postgres
    restart: unless-stopped

    environment:
      POSTGRES_USER: ${POSTGRES_USER}
      POSTGRES_PASSWORD: ${POSTGRES_PASSWORD}
      POSTGRES_DB: ${POSTGRES_DB}

    ports:
      - "5432:5432"

    volumes:
      - contextflow_postgres_data:/var/lib/postgresql/data

    healthcheck:
      test:
        [
          "CMD-SHELL",
          "pg_isready -U $${POSTGRES_USER} -d $${POSTGRES_DB}"
        ]
      interval: 5s
      timeout: 5s
      retries: 5

volumes:
  contextflow_postgres_data:
```

### Important database behavior

PostgreSQL initialization environment variables are applied when the data directory is first initialized.

Changing `POSTGRES_PASSWORD` in `.env` does not automatically change the password of an already-created PostgreSQL role.

For an existing database, update the role explicitly:

```sql
ALTER ROLE contextflow WITH PASSWORD 'new_password';
```

Never delete the persistent database volume just to resolve a password issue.

---

## Reliability and Engineering Decisions

ContextFlow is being designed around reliable data movement and predictable system behavior.

The following principles will guide implementation:

**Idempotent ingestion**

Repeated ingestion attempts should not create duplicate records when the same source data is processed again.

**Checkpoint management**

A checkpoint should advance only after the corresponding destination transaction commits successfully.

**Retry handling**

Transient failures should be retried using controlled retry policies rather than silently discarded.

**Database connection management**

The application should use connection pooling and bounded timeouts instead of opening an unrestricted database connection for every request.

**Separation of concerns**

Configuration, database access, API routing, ingestion, and governance should remain separate logical modules.

**Observability**

Important operations should produce useful logs and measurable operational information.

These are design principles and planned reliability guarantees; they are not claims that the complete ingestion engine has already been implemented.

---

## Development Progress

The following table records the verified progress to date.

| Component | Status | Details |
|---|---|---|
| GitHub repository | <span>Completed</span> | Public repository initialized |
| Go backend foundation | <span>Implemented</span> | Application entry point and packages |
| Environment configuration | <span>Implemented</span> | `DATABASE_URL` validation |
| Chi HTTP router | <span>Implemented</span> | Routing and middleware |
| PostgreSQL Docker image | <span>Running</span> | PostgreSQL 17 Alpine |
| Persistent database volume | <span>Configured</span> | Docker named volume |
| PostgreSQL health check | <span>Healthy</span> | Container reports healthy |
| Database role configuration | <span>Verified internally</span> | TCP authentication inside container |
| Go-to-PostgreSQL connection | <span>In progress</span> | Windows-host authentication issue |
| Health endpoint | <span>Implemented</span> | End-to-end verification pending |
| Ingestion engine | <span>Planned</span> | Workers, retries, checkpoints |
| Context data model | <span>Planned</span> | Records and metadata |
| Governance and RBAC | <span>Planned</span> | Access control and audit |
| Frontend dashboard | <span>Planned</span> | React + TypeScript |
| Production deployment | <span>Planned</span> | Deployment and operational hardening |

## Known Issues

### PostgreSQL authentication from the Windows host

The PostgreSQL container is healthy and accepts the configured password through TCP from inside the container.

However, the Go backend running directly on Windows currently encounters:

```text
failed SASL auth:
FATAL: password authentication failed for user "contextflow"
(SQLSTATE 28P01)
```

Additional investigation has established:

- Docker publishes container port `5432` to host port `5432`.
- Docker Desktop owns the Windows listener for that published port.
- PostgreSQL listens on `0.0.0.0:5432` inside the container.
- No `PG*` environment-variable overrides were found in the backend PowerShell session.
- A separate Windows `postgres.exe` process was also identified.

The exact cause of the host-to-container authentication failure has not yet been established.

**Current next step:** Inspect the Docker network configuration and test connectivity through the container network, isolating the application from Windows host port forwarding.

No database volume deletion or destructive database reset is required for this investigation.

## Security

Security is an important part of the intended architecture.

Current practices:

- Environment-based database credentials.
- `.env` excluded from version control.
- Persistent database storage.
- No hardcoded database passwords in Go source files.

Future security work:

- Authentication and authorization.
- Workspace-level access isolation.
- Role-based access control.
- Permission-aware context retrieval.
- Audit logging.
- Secure secret management.
- TLS and production database connection configuration.
- Input validation and request limits.

**Important:** The current development setup is not a production security deployment. The local database configuration and API require further hardening before production use.

## Roadmap

### Phase 1 — Infrastructure Foundation

- [x] Initialize GitHub repository.
- [x] Establish Go backend structure.
- [x] Configure Docker Compose.
- [x] Start PostgreSQL 17.
- [x] Configure persistent database storage.
- [x] Implement environment-based configuration.
- [x] Implement PostgreSQL connection-pool initialization.
- [x] Implement database-aware health handler.
- [ ] Resolve Windows-host PostgreSQL connectivity.
- [ ] Verify successful backend startup.
- [ ] Add automated backend and database tests.

### Phase 2 — Context Data Foundation

- [ ] Design context record schema.
- [ ] Implement database migrations.
- [ ] Build source registration.
- [ ] Implement initial PostgreSQL connector.
- [ ] Add record ingestion.
- [ ] Implement idempotent persistence.

### Phase 3 — Reliable Data Movement

- [ ] Build ingestion workers.
- [ ] Add retry policies.
- [ ] Implement checkpoints.
- [ ] Add execution tracking.
- [ ] Handle partial failures.
- [ ] Add recovery mechanisms.

### Phase 4 — Governance

- [ ] Introduce organizations and workspaces.
- [ ] Implement role-based access control.
- [ ] Add permission-aware retrieval.
- [ ] Implement audit events.

### Phase 5 — Context API and Dashboard

- [ ] Build context retrieval endpoints.
- [ ] Add filtering and pagination.
- [ ] Develop React + TypeScript dashboard.
- [ ] Display ingestion runs and system health.
- [ ] Add error and pipeline monitoring.

### Phase 6 — Production Readiness

- [ ] Automated testing and CI.
- [ ] Structured logging.
- [ ] Metrics and distributed tracing.
- [ ] Deployment configuration.
- [ ] Security review.
- [ ] Load and failure testing.
- [ ] Operational documentation.

---

## Contributing

Contributions, bug reports, design discussions, and suggestions are welcome.

1. Fork the repository.
2. Create a feature branch.
3. Make focused changes.
4. Add tests for new behavior.
5. Run the relevant checks.
6. Submit a pull request with a clear description.

Please avoid committing secrets, local environment files, database dumps, or generated build artifacts.

## License

This project is licensed under the terms of the MIT License. See [LICENSE](LICENSE) for details.

---

## Project Information

**Project:** ContextFlow  
**Category:** Open-source infrastructure / Developer tools / AI infrastructure  
**Primary language:** Go  
**Database:** PostgreSQL  
**Repository:** [github.com/Sohamsig/ContextFlow](https://github.com/Sohamsig/ContextFlow)

**ContextFlow — Building the infrastructure layer between enterprise data and AI applications.**
