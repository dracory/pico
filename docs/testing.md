# Testing

## Running Tests

```bash
# Run all tests
task test

# Run with coverage
task cover

# Run specific package
go test ./internal/config/
go test ./internal/app/
```

## Test Configuration

Tests use `config.New()` to create configs programmatically, avoiding environment variable dependencies:

```go
func newTestConfig() config.ConfigInterface {
    cfg := config.New()
    cfg.SetAppName("PicoTest")
    cfg.SetAppEnv("testing")
    cfg.SetAppHost("127.0.0.1")
    cfg.SetAppPort("8080")
    cfg.SetAppUrl("http://localhost:8080")
    cfg.SetAppDebug(true)
    cfg.SetDatabaseDriver("sqlite")
    cfg.SetDatabaseName(fmt.Sprintf("file:pico_test_%d?mode=memory&cache=shared", time.Now().UnixNano()))
    return cfg
}
```

### In-Memory SQLite

App tests use shared in-memory SQLite with unique DSNs per test to avoid collisions:

```go
fmt.Sprintf("file:pico_test_%d?mode=memory&cache=shared", time.Now().UnixNano())
```

## Config Tests

Config tests that require environment variables use `t.Setenv()` which automatically saves and restores env vars after each test:

```go
func setRequiredEnv(t *testing.T) {
    t.Helper()
    t.Setenv(KEY_APP_HOST, "localhost")
    t.Setenv(KEY_APP_PORT, "8080")
    t.Setenv(KEY_APP_ENVIRONMENT, "testing")
    t.Setenv(KEY_DB_DRIVER, "sqlite")
    t.Setenv(KEY_DB_DATABASE, ":memory:")
}
```

### Testing Validation Errors

```go
func TestNewFromEnv_MissingRequiredFields(t *testing.T) {
    t.Setenv(KEY_APP_HOST, "")
    t.Setenv(KEY_APP_PORT, "")
    // ...

    _, err := NewFromEnv()
    verr, ok := err.(env.ValidationError)
    if !ok {
        t.Fatalf("expected env.ValidationError, got %T", err)
    }
}
```

## App Tests

App tests create a full app instance with a real in-memory database:

```go
func TestNew_SetsDatabase(t *testing.T) {
    cfg := newTestConfig()

    a, err := New(cfg)
    if err != nil {
        t.Fatalf("app.New returned error: %v", err)
    }
    defer a.Close()

    if a.GetDatabase() == nil {
        t.Fatal("expected database to be non-nil")
    }
}
```

> **Important:** Always call `defer a.Close()` after `app.New()` to prevent database connection leaks.

## Environment-Aware Behavior

The router skips `LoggerMiddleware` and `RecoveryMiddleware` when the environment is `testing`:

```go
if app.GetConfig() != nil && !app.GetConfig().IsEnvTesting() {
    middlewares = append(middlewares,
        rtrMiddleware.LoggerMiddleware(),
        rtrMiddleware.RecoveryMiddleware(),
    )
}
```

This keeps test output clean and avoids middleware interference during testing.

## Coverage

Current coverage:

| Package | Coverage |
|---------|----------|
| `internal/config` | ~92% |
| `internal/app` | ~89% |
