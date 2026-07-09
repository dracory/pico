# Database

## Introduction

Pico uses the [dracory/neat](https://github.com/dracory/neat) database layer, which provides a thin abstraction over `database/sql` with connection pooling, multi-driver support, and query building.

## Supported Drivers

| Driver | Import Path | Default Port |
|--------|-------------|-------------|
| SQLite | `modernc.org/sqlite` | — |
| Turso | `tursodatabase/libsql-client-go` | — |
| MySQL | (via neat) | `3306` |
| MariaDB | (via neat) | `3306` |
| PostgreSQL | (via neat) | `5432` |
| Oracle | (via neat) | `1521` |
| SQL Server | (via neat) | `1433` |

SQLite uses `modernc.org/sqlite`, a pure Go implementation — no CGO required.
Turso uses `libsql-client-go` for remote access to Turso Cloud over the libSQL wire protocol — also pure Go, no CGO required.

## Configuration

Database configuration is loaded from environment variables. See the [Configuration](configuration.md#database-configuration) docs for the full variable reference.

### SQLite

```env
DB_DRIVER=sqlite
DB_DATABASE=database.db
```

For in-memory testing:

```env
DB_DRIVER=sqlite
DB_DATABASE=:memory:
```

### Turso

Turso is a SQLite edge database. It connects remotely to Turso Cloud via the libSQL wire protocol:

```env
DB_DRIVER=turso
DB_DATABASE=libsql://my-db.turso.io
DB_PASSWORD=your-auth-token
```

Neat builds the connection string from these components — no manual DSN needed. The `DB_DATABASE` field holds the libsql URL, and `DB_PASSWORD` holds the auth token. Get both from the [Turso dashboard](https://app.turso.tech).

For your specific database:
```env
DB_DRIVER=turso
DB_DATABASE=libsql://pico-sinevia.aws-eu-west-1.turso.io
DB_PASSWORD=your-auth-token
```

Turso inherits SQLite's single-writer limitation, so the same pool constraints apply (MaxOpenConns=1, MaxIdleConns=1). Host, port, and username are not required.

### PostgreSQL

```env
DB_DRIVER=postgres
DB_HOST=localhost
DB_PORT=5432
DB_DATABASE=picodb
DB_USERNAME=user
DB_PASSWORD=secret
DB_SSL_MODE=require
```

### MySQL

```env
DB_DRIVER=mysql
DB_HOST=localhost
DB_PORT=3306
DB_DATABASE=picodb
DB_USERNAME=user
DB_PASSWORD=secret
DB_CHARSET=utf8mb4
```

## Connection Pooling

Pico automatically configures connection pools based on the driver:

| Setting | SQLite / Turso | MySQL / PostgreSQL |
|---------|-----------------|-------------------|
| Max Open Connections | 1 | 25 |
| Max Idle Connections | 1 | 5 |
| Connection Max Lifetime | 30s | 300s |
| Connection Max Idle Time | 5s | 5s |

SQLite and Turso use a single connection to prevent writer contention. These values can be overridden via environment variables:

```env
DB_MAX_OPEN_CONNS=10
DB_MAX_IDLE_CONNS=5
DB_CONN_MAX_LIFETIME_SECONDS=600
DB_CONN_MAX_IDLE_TIME_SECONDS=60
```

> **Note:** For SQLite and Turso, `DB_MAX_OPEN_CONNS` and `DB_MAX_IDLE_CONNS` are always overridden to 1 regardless of env var settings, to prevent concurrent writer issues.

## Accessing the Database

The database is available from the app instance via two interfaces:

### Standard `*sql.DB`

```go
db := app.GetDatabase() // *sql.DB
rows, err := db.Query("SELECT id, name FROM users")
```

### Neat Database

```go
neatDB := app.GetNeatDatabase() // *neatdatabase.Database
```

The neat database provides higher-level query building and connection management.

## Multi-Connection Support

Pico supports named database connections. The default connection is named `"default"` but can be customized:

```env
DB_DEFAULT_CONNECTION=primary
```

Access connections programmatically:

```go
conns := cfg.GetDatabaseConnections()
conn := cfg.GetDatabaseConnectionByName("default")

conn.GetDriver()    // "sqlite"
conn.GetHost()      // "localhost"
conn.GetDatabase()  // "picodb"
```

## DSN Override

If you need full control over the connection string, set `DB_DSN`:

```env
DB_DSN=postgres://user:pass@localhost:5432/picodb?sslmode=require
```

When `DB_DSN` is set, it takes precedence over individual connection parameters.

## Table Prefix

Use `DB_PREFIX` to prefix all tables:

```env
DB_PREFIX=pico_
```

## Graceful Shutdown

The database is automatically closed when the app shuts down. In `main.go`:

```go
defer func() {
    if err := a.Close(); err != nil {
        slog.Error("Failed to close app", "error", err)
    }
}()
```

`app.Close()` closes the neat database connection and nilifies the references. Calling `Close()` multiple times is safe.
