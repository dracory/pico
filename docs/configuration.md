# Configuration

## Introduction

All configuration for Pico is driven by environment variables, loaded from a `.env` file in the project root. The configuration system provides validated, typed access to application and database settings via the `ConfigInterface`.

## Environment Configuration

Pico uses the [dracory/env](https://github.com/dracory/env) library to load environment variables from `.env` automatically when `config.NewFromEnv()` is called.

### Loading Configuration

```go
cfg, err := config.NewFromEnv()
if err != nil {
    // handle validation error
}
```

If any required environment variable is missing, `NewFromEnv()` returns a `ValidationError` containing all validation errors:

```go
verr, ok := err.(env.ValidationError)
if ok {
    for _, e := range verr.Errors() {
        fmt.Println(e)
    }
}
```

### Creating a Config Programmatically

You can also create a config without environment variables, useful for testing:

```go
cfg := config.New()
cfg.SetAppName("MyApp")
cfg.SetAppEnv("production")
cfg.SetAppHost("0.0.0.0")
cfg.SetAppPort("3000")
cfg.SetDatabaseDriver("sqlite")
cfg.SetDatabaseName(":memory:")
```

## Application Configuration

| Variable | Default | Required | Description |
|----------|---------|----------|-------------|
| `APP_NAME` | `Pico` | No | Application name |
| `APP_ENV` | — | Yes | Environment: `local`, `development`, `staging`, `production`, `testing` |
| `APP_DEBUG` | `false` | No | Show detailed error messages |
| `APP_HOST` | — | Yes | Listen host (e.g., `0.0.0.0`) |
| `APP_PORT` | — | Yes | Listen port (e.g., `8080`) |
| `APP_URL` | `http://localhost:8080` | No | Application URL |

### Environment Detection

The config provides helper methods to check the current environment:

```go
cfg.IsEnvLocal()       // true when APP_ENV=local
cfg.IsEnvDevelopment() // true when APP_ENV=development
cfg.IsEnvStaging()     // true when APP_ENV=staging
cfg.IsEnvProduction()  // true when APP_ENV=production
cfg.IsEnvTesting()     // true when APP_ENV=testing
```

## Database Configuration

| Variable | Default | Required | Description |
|----------|---------|----------|-------------|
| `DB_DRIVER` | — | Yes | Database driver: `sqlite`, `turso`, `mysql`, `postgres` |
| `DB_DATABASE` | — | Yes | Database name, file path (SQLite), or libsql URL (Turso) |
| `DB_HOST` | — | Yes* | Database host (*not required for SQLite or Turso*) |
| `DB_PORT` | — | Yes* | Database port (*not required for SQLite or Turso*) |
| `DB_USERNAME` | — | Yes* | Database username (*not required for SQLite or Turso*) |
| `DB_PASSWORD` | — | Yes* | Database password (*not required for SQLite or Turso*) |
| `DB_SSL_MODE` | `require` | No | SSL mode (PostgreSQL only) |
| `DB_CHARSET` | `utf8mb4` | No | Character set (MySQL only) |
| `DB_TIMEZONE` | `UTC` | No | Database timezone |
| `DB_DSN` | — | No | Direct DSN override |
| `DB_PREFIX` | — | No | Table prefix |
| `DB_DEFAULT_CONNECTION` | `default` | No | Default connection name |

### Connection Pool Settings

| Variable | Default (non-SQLite) | Default (SQLite / Turso) | Description |
|----------|---------------------|------------------------|-------------|
| `DB_MAX_OPEN_CONNS` | `25` | `1` | Maximum open connections |
| `DB_MAX_IDLE_CONNS` | `5` | `1` | Maximum idle connections |
| `DB_CONN_MAX_LIFETIME_SECONDS` | `300` | `30` | Max connection lifetime (seconds) |
| `DB_CONN_MAX_IDLE_TIME_SECONDS` | `5` | `5` | Max idle time (seconds) |

SQLite and Turso use conservative pool settings to avoid writer contention. These can be overridden but the defaults are recommended.

### Example: SQLite

```env
DB_DRIVER=sqlite
DB_DATABASE=database.db
```

### Example: Turso

```env
DB_DRIVER=turso
DB_DATABASE=libsql://my-db.turso.io
DB_PASSWORD=your-auth-token
```

### Example: PostgreSQL

```env
DB_DRIVER=postgres
DB_HOST=localhost
DB_PORT=5432
DB_DATABASE=picodb
DB_USERNAME=user
DB_PASSWORD=secret
DB_SSL_MODE=require
```

### Example: MySQL

```env
DB_DRIVER=mysql
DB_HOST=localhost
DB_PORT=3306
DB_DATABASE=picodb
DB_USERNAME=user
DB_PASSWORD=secret
DB_CHARSET=utf8mb4
```

## Accessing Configuration Values

Configuration is accessed through the `ConfigInterface`, which is available from the app instance:

```go
a, err := app.New(cfg)

a.GetConfig().GetAppName()
a.GetConfig().GetAppEnv()
a.GetConfig().GetDatabaseDriver()
a.GetConfig().GetDatabaseName()
```

### Multi-Connection Support

Pico supports named database connections. The default connection is named `"default"` but can be renamed via `DB_DEFAULT_CONNECTION`:

```go
conns := cfg.GetDatabaseConnections()           // []DatabaseConnectionConfigInterface
conn := cfg.GetDatabaseConnectionByName("default") // DatabaseConnectionConfigInterface
```

Each connection exposes its own driver, host, port, database, username, password, SSL mode, charset, timezone, DSN, and prefix.
