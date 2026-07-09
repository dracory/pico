# The Application

## Introduction

The `AppInterface` is the central dependency injection container in Pico. It provides access to the config, logger, and database connections throughout the application lifecycle.

## AppInterface

```go
type AppInterface interface {
    Close() error

    GetLogger() *slog.Logger
    SetLogger(l *slog.Logger)

    GetConfig() config.ConfigInterface
    SetConfig(c config.ConfigInterface)

    GetDatabase() *sql.DB
    SetDatabase(db *sql.DB)

    GetNeatDatabase() *neatdatabase.Database
    SetNeatDatabase(db *neatdatabase.Database)
}
```

## Creating an Application

```go
cfg, err := config.NewFromEnv()
if err != nil {
    log.Fatal(err)
}

a, err := app.New(cfg)
if err != nil {
    log.Fatal(err)
}
defer a.Close()
```

`app.New()` performs the following:

1. **Validates** the config is non-nil
2. **Creates a logger** — `slog` with colored output via [tint](https://github.com/lmittmann/tint)
3. **Opens the database** — using `config.DatabaseNeatConfig()` to build the neat database config
4. **Returns** an `AppInterface` with all services initialized

## Accessing Services

### Logger

```go
logger := a.GetLogger()
logger.Info("Server started", "host", cfg.GetAppHost(), "port", cfg.GetAppPort())
```

The default logger writes to `stdout` with colored output. Replace it with your own:

```go
a.SetLogger(slog.New(slog.NewJSONHandler(os.Stdout, nil)))
```

### Config

```go
cfg := a.GetConfig()
cfg.GetAppName()
cfg.GetAppEnv()
cfg.IsEnvProduction()
```

See [Configuration](configuration.md) for full details.

### Database

```go
db := a.GetDatabase()      // *sql.DB — standard database/sql
neatDB := a.GetNeatDatabase() // *neatdatabase.Database — neat ORM
```

See [Database](database.md) for full details.

## Shutdown

Call `Close()` to release database connections:

```go
defer a.Close()
```

`Close()` is safe to call on:
- **Nil receiver** — returns `nil` without panic
- **Already closed app** — returns `nil` without error
- **App with no database** — returns `nil` without error

## Lifecycle in main.go

```go
func main() {
    // 1. Load config
    cfg, err := config.NewFromEnv()

    // 2. Create app (opens database, creates logger)
    a, err := app.New(cfg)
    defer a.Close()

    // 3. Start web server with router
    server, err := websrv.Start(websrv.Options{
        Host:    cfg.GetAppHost(),
        Port:    cfg.GetAppPort(),
        URL:     cfg.GetAppUrl(),
        Handler: routes.Router(a).ServeHTTP,
    })

    // 4. Wait for shutdown signal
    <-sigs

    // 5. Graceful server shutdown
    server.Shutdown(shutdownCtx)
}
```

## Extending the Application

To add a new service to the app:

1. Add the field and interface methods to `app_implementation.go` and `app_interface.go`
2. Initialize it in `app.New()`
3. Clean it up in `Close()`

```go
type appImplementation struct {
    cfg    config.ConfigInterface
    neatDB *neatdatabase.Database
    db     *sql.DB
    logger *slog.Logger
    cache  *cache.Cache // new service
}
```
