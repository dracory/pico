# Pico Documentation

## Introduction

Pico is a minimal Go web application starter template — "Laravel Lumen" for the Dracory ecosystem. It provides the essentials to build fast microservices and APIs: configuration, routing, and database access.

## Documentation

- **[Installation](installation.md)** — Requirements, setup, and project structure
- **[Configuration](configuration.md)** — Environment variables, app and database config
- **[The Application](application.md)** — AppInterface, dependency injection, lifecycle
- **[HTTP Routing](routing.md)** — Route definitions, handlers, and middleware
- **[Database](database.md)** — SQLite/MySQL/PostgreSQL, connection pooling, multi-connection
- **[Testing](testing.md)** — Test patterns, in-memory SQLite, coverage
- **[Deployment](deployment.md)** — Building, Docker, process management, health checks

## Quick Start

```bash
git clone https://github.com/dracory/pico
cd pico
cp .env.example .env
go run ./cmd/server
```

Server starts at `http://localhost:8080`.
