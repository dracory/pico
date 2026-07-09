# Deployment

## Building

```bash
task build
```

This produces `bin/pico.exe`. For cross-compilation:

```bash
GOOS=linux GOARCH=amd64 go build -o bin/pico ./cmd/server
```

## Environment Configuration

Set environment variables for your production environment:

```env
APP_NAME=MyApp
APP_ENV=production
APP_DEBUG=false
APP_HOST=0.0.0.0
APP_PORT=8080
APP_URL=https://myapp.com

DB_DRIVER=postgres
DB_HOST=db.internal
DB_PORT=5432
DB_DATABASE=picodb
DB_USERNAME=pico
DB_PASSWORD=secure-password
DB_SSL_MODE=require
```

> **Security:** Never commit `.env` to version control. Use secrets management or environment injection in your deployment platform.

## Docker

Create a `Dockerfile`:

```dockerfile
FROM golang:1.26 AS builder
WORKDIR /app
COPY . .
RUN CGO_ENABLED=0 go build -o /pico ./cmd/server

FROM alpine:latest
RUN apk add --no-cache ca-certificates
WORKDIR /app
COPY --from=builder /pico .
COPY .env.example .env
EXPOSE 8080
CMD ["./pico"]
```

Build and run:

```bash
docker build -t pico .
docker run -p 8080:8080 --env-file .env pico
```

> **Note:** Pico uses `modernc.org/sqlite` (pure Go), so `CGO_ENABLED=0` works for SQLite deployments.

## Graceful Shutdown

Pico handles `SIGINT` and `SIGTERM` for graceful shutdown:

```go
sigs := make(chan os.Signal, 1)
signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

<-sigs
// Server shuts down with 10s timeout
server.Shutdown(shutdownCtx)
```

This ensures in-flight requests complete before the process exits. The database connection is closed via `defer a.Close()`.

## Health Check

The default route `/` returns a JSON health check:

```json
{"name":"pico","status":"running"}
```

Use this for load balancer health checks:

```bash
curl http://localhost:8080/
```

## Process Management

### Systemd

```ini
[Unit]
Description=Pico Web Server
After=network.target

[Service]
Type=simple
User=pico
WorkingDirectory=/opt/pico
ExecStart=/opt/pico/pico
Restart=always
RestartSec=5
EnvironmentFile=/opt/pico/.env

[Install]
WantedBy=multi-user.target
```

### Behind a Reverse Proxy

When running behind Nginx, Caddy, or a load balancer, set `APP_HOST` to `0.0.0.0` and configure the proxy to forward to `APP_PORT`:

```nginx
server {
    listen 80;
    server_name myapp.com;

    location / {
        proxy_pass http://127.0.0.1:8080;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }
}
```

## Database Considerations

### SQLite

For production with SQLite:
- Use a file path, not `:memory:`
- Consider WAL mode for better read concurrency
- Keep `DB_MAX_OPEN_CONNS=1` to avoid writer contention

### Turso

For production with Turso:
- Set the auth token via `DB_PASSWORD`
- Turso is ideal for edge deployments and serverless environments
- Same pool constraints as SQLite (MaxOpenConns=1)

```env
DB_DRIVER=turso
DB_DATABASE=libsql://my-db.turso.io
DB_PASSWORD=your-auth-token
```

### PostgreSQL / MySQL

For production:
- Use connection pooling with appropriate limits
- Set `DB_SSL_MODE=require` for encrypted connections
- Monitor pool exhaustion
