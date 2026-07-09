# Pico

A minimal Go web application starter template — "Laravel Lumen" for the Dracory ecosystem.

Pico is a stripped-down version of [Blueprint](https://github.com/dracory/blueprint), keeping only the essentials: **config**, **rtr** (router), and **neat** (database). No stores, no tasks, no schedulers, no CMS, no auth — just a fast, clean foundation for microservices and APIs.

## What's Included

- **Config** — Environment-based configuration with validation (app + database only)
- **Router** — High-performance `rtr` router with minimal middleware chain
- **Database** — `neat` ORM with SQLite/MySQL/PostgreSQL support and connection pooling
- **Graceful Shutdown** — Signal-based shutdown with timeout

## What's NOT Included (vs Blueprint)

- No stores (no userstore, blogstore, cmsstore, etc.)
- No authentication or session management
- No background tasks or schedulers
- No CMS or admin panel
- No email system
- No encryption/vault
- No caches (memory or file)

## Installation

```bash
git clone https://github.com/dracory/pico
cd pico
cp .env.example .env
```

## Local Development

```bash
go run ./cmd/server
```

The server starts at `http://localhost:8080` by default.

## Environment Variables

| Variable | Default | Description |
|----------|---------|-------------|
| `APP_NAME` | `Pico` | Application name |
| `APP_ENV` | — | Environment: `local`, `development`, `staging`, `production`, `testing` |
| `APP_DEBUG` | `false` | Show detailed errors |
| `APP_HOST` | — | Listen host (e.g., `0.0.0.0`) |
| `APP_PORT` | — | Listen port (e.g., `8080`) |
| `APP_URL` | `http://localhost:8080` | Application URL |
| `DB_DRIVER` | — | Database driver: `sqlite`, `turso`, `mysql`, `mariadb`, `postgres`, `oracle`, `sqlserver` |
| `DB_DATABASE` | — | Database name, file path (SQLite), or libsql URL (Turso) |
| `DB_HOST` | — | Database host (not required for SQLite) |
| `DB_PORT` | — | Database port (not required for SQLite) |
| `DB_USERNAME` | — | Database username (not required for SQLite) |
| `DB_PASSWORD` | — | Database password (not required for SQLite) |
| `DB_SSL_MODE` | `require` | SSL mode (PostgreSQL only) |
| `DB_CHARSET` | `utf8mb4` | Character set (MySQL only) |
| `DB_TIMEZONE` | `UTC` | Database timezone |
| `DB_DSN` | — | Direct DSN override (optional) |
| `DB_PREFIX` | — | Table prefix (optional) |
| `DB_MAX_OPEN_CONNS` | `25` | Max open connections (SQLite/Turso: `1`) |
| `DB_MAX_IDLE_CONNS` | `5` | Max idle connections (SQLite/Turso: `1`) |
| `DB_CONN_MAX_LIFETIME_SECONDS` | `300` | Max connection lifetime in seconds (SQLite/Turso: `30`) |
| `DB_CONN_MAX_IDLE_TIME_SECONDS` | `5` | Max connection idle time in seconds |

## Project Structure

```
pico/
├── cmd/
│   └── server/
│       └── main.go          # Entry point
├── internal/
│   ├── app/
│   │   ├── app_interface.go     # AppInterface (DI container)
│   │   ├── app_implementation.go # App implementation
│   │   └── database_open.go     # Database connection
│   ├── config/
│   │   ├── config.go            # ConfigInterface + implementation
│   │   ├── app_config.go        # App config loader
│   │   ├── database_config.go   # Database config loader
│   │   ├── constants.go         # Environment variable keys
│   │   └── version.go           # Version
│   └── routes/
│       └── router.go            # Router + minimal middleware
├── .env.example
├── go.mod
└── README.md
```

## License

AGPL-3.0
