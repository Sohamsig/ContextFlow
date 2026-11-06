# ContextFlow

**Open-source enterprise context infrastructure for AI applications.**

Connect your data. Govern your knowledge. Power your AI.

ContextFlow aims to connect enterprise data sources, ingest and organize information, enforce access policies, and provide reliable context to AI applications.

> **Project status:** Early development. The backend foundation is being implemented.

## Vision

AI applications need more than raw data. They need reliable, structured, permission-aware context.

ContextFlow is being designed to provide:

- Data source connectors.
- Reliable ingestion and processing.
- Context storage and metadata management.
- Access control and governance.
- Context retrieval APIs.
- Pipeline monitoring and observability.

These are project goals. Not all capabilities are implemented yet.

## Architecture

### Intended architecture

```text
External Data Sources
        |
        v
Connector Framework
        |
        v
Ingestion and Processing Engine
        |
        v
Context Storage
        |
        v
Governance and Access Control
        |
        v
Context Retrieval API
        |
        v
AI Applications and Agents
```

### Current backend foundation

```text
Go HTTP API (Chi)
        |
        v
Configuration Layer
        |
        v
PostgreSQL Connection Pool (pgx)
        |
        v
PostgreSQL 17
```

The initial implementation uses a modular Go backend and PostgreSQL.

## Technology Stack

| Component | Technology |
|---|---|
| Backend | Go |
| HTTP router | Chi |
| Database | PostgreSQL 17 |
| Database driver | pgx v5 |
| Connection pooling | pgxpool |
| Local infrastructure | Docker |
| Container orchestration | Docker Compose |
| Frontend | Planned: React and TypeScript |

## Repository Structure

```text
ContextFlow/
├── backend/
│   ├── cmd/
│   │   └── server/
│   │       └── main.go
│   ├── internal/
│   │   ├── config/
│   │   └── database/
│   └── go.mod
├── deployments/
├── docs/
├── frontend/
├── .env.example
├── .gitignore
├── docker-compose.yml
├── LICENSE
└── README.md
```

## Getting Started

### Prerequisites

- Go 1.26.4 or a compatible toolchain.
- Docker Desktop.
- Docker Compose.

### 1. Clone the repository

```powershell
git clone https://github.com/Sohamsig/ContextFlow.git
Set-Location .\ContextFlow
```

### 2. Configure environment variables

Create your local environment file:

```powershell
Copy-Item .env.example .env
```

Update the placeholder password in `.env` with a local development password.

Keep `.env` private. Never commit real credentials.

### 3. Start PostgreSQL

```powershell
docker compose up -d
```

Check the database container:

```powershell
docker compose ps
```

PostgreSQL uses a persistent Docker volume and a container health check.

### 4. Configure the Go backend

Set `DATABASE_URL` in the PowerShell terminal where you will run the API.

Replace `YOUR_PASSWORD` with the same password configured in `.env`.

```powershell
$env:DATABASE_URL = "postgres://contextflow:YOUR_PASSWORD@127.0.0.1:5432/contextflow?sslmode=disable"
```

The current Go configuration reads `DATABASE_URL` from the process environment. It does not automatically load `.env`.

### 5. Start the API

```powershell
Set-Location .\backend
go run .\cmd\server
```

The server is configured to listen on port `8080`.

### 6. Check API health

Open another PowerShell terminal:

```powershell
Invoke-RestMethod http://localhost:8080/api/v1/health
```

The health endpoint reports:

- Service status.
- Service name and version.
- Database connectivity.
- UTC timestamp.

The endpoint returns HTTP 200 when the database is reachable and HTTP 503 when the database check fails.

## API

### Health Check

`GET /api/v1/health`

The health endpoint checks PostgreSQL connectivity using the configured connection pool.

Example response when healthy:

```json
{
  "status": "healthy",
  "service": "contextflow",
  "version": "0.1.0",
  "database": "connected",
  "timestamp": "2026-11-06T12:00:00Z"
}
```

The timestamp above is illustrative.

## Engineering Principles

The project is being developed around these reliability principles:

- Explicit configuration.
- Environment-based secret management.
- Bounded database connection pools.
- Context-aware operations and timeouts.
- Idempotent ingestion.
- Checkpoint advancement after successful destination commits.
- Recoverable processing and observable failures.

These are engineering goals; the complete ingestion and recovery mechanisms are not implemented yet.

## Development Status

| Component | Status |
|---|---|
| Git repository and Go module | Initialized |
| Configuration package | Implemented |
| PostgreSQL connection pool | Implemented |
| Chi HTTP router | Implemented |
| Health endpoint | Implemented |
| PostgreSQL Docker service | Configured |
| Automated tests | Pending |
| Database migrations | Planned |
| Data connector framework | Planned |
| Ingestion workers and checkpoints | Planned |
| Context retrieval API | Planned |
| Governance and RBAC | Planned |
| React dashboard | Planned |
| Metrics and distributed tracing | Planned |

## Security

- Never commit `.env` or real credentials.
- Use `.env.example` for placeholder configuration.
- Use secure database authentication and TLS in deployed environments.
- Enforce authorization before exposing enterprise context.
- Avoid logging sensitive records or credentials.

The current configuration is for local development, not production deployment.

## Roadmap

1. Verify backend and PostgreSQL integration.
2. Add database migrations and automated tests.
3. Define the context record and metadata model.
4. Implement the first PostgreSQL source connector.
5. Build reliable ingestion, retries, and checkpoints.
6. Add permission-aware context retrieval.
7. Build pipeline monitoring and a React dashboard.
8. Introduce observability and deployment infrastructure as needed.

## Contributing

Contributions, bug reports, and technical discussions are welcome.

Please open an issue before making significant architectural changes.

## License

This project is licensed under the MIT License. See [LICENSE](LICENSE).
