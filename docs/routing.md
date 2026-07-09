# HTTP Routing

## Basic Routing

All routes are defined in `internal/routes/router.go`. Pico uses the [dracory/rtr](https://github.com/dracory/rtr) router, which provides a clean, expressive API for defining HTTP routes.

```go
func routes(app app.AppInterface) []rtr.RouteInterface {
    return []rtr.RouteInterface{
        rtr.Get("/", func(w http.ResponseWriter, r *http.Request) {
            w.Header().Set("Content-Type", "application/json")
            w.WriteHeader(http.StatusOK)
            w.Write([]byte(`{"status":"ok"}`))
        }),
    }
}
```

## Available Router Methods

The router allows you to register routes that respond to any HTTP verb:

```go
rtr.Get(uri, handler)
rtr.Post(uri, handler)
rtr.Put(uri, handler)
rtr.Patch(uri, handler)
rtr.Delete(uri, handler)
rtr.Options(uri, handler)
rtr.Head(uri, handler)
```

## Adding Routes

To add a new route, append it to the `routes()` function return slice:

```go
func routes(app app.AppInterface) []rtr.RouteInterface {
    return []rtr.RouteInterface{
        rtr.Get("/", homeHandler),
        rtr.Get("/users", usersListHandler),
        rtr.Post("/users", usersCreateHandler),
        rtr.Get("/users/{id}", usersShowHandler),
    }
}
```

## Route Handlers

Route handlers are standard `http.HandlerFunc` functions:

```go
func usersListHandler(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)
    w.Write([]byte(`[{"id":1,"name":"Alice"}]`))
}
```

## Middleware

### Global Middleware

Global middleware is applied to all routes and is defined in the `globalMiddlewares()` function:

```go
func globalMiddlewares(app app.AppInterface) []rtr.MiddlewareInterface {
    middlewares := []rtr.MiddlewareInterface{
        rtrMiddleware.CompressMiddleware(5, "text/html", "text/css"),
        rtrMiddleware.GetHead(),
        rtrMiddleware.CleanPathMiddleware(),
        rtrMiddleware.RedirectSlashesMiddleware(),
        rtrMiddleware.TimeoutMiddleware(30 * time.Second),
    }

    // Logger and Recovery are skipped in testing environment
    if app.GetConfig() != nil && !app.GetConfig().IsEnvTesting() {
        middlewares = append(middlewares,
            rtrMiddleware.LoggerMiddleware(),
            rtrMiddleware.RecoveryMiddleware(),
        )
    }

    return middlewares
}
```

### Available Middleware

Pico includes the following middleware from `dracory/rtr/middlewares`:

| Middleware | Description |
|-----------|-------------|
| `CompressMiddleware` | Gzip compression for responses |
| `GetHead` | Automatically adds HEAD routes for GET handlers |
| `CleanPathMiddleware` | Cleans double slashes from paths |
| `RedirectSlashesMiddleware` | Redirects trailing slashes |
| `TimeoutMiddleware` | Request timeout (default: 30s) |
| `LoggerMiddleware` | Structured request logging (skipped in testing) |
| `RecoveryMiddleware` | Panic recovery (skipped in testing) |

### Adding Custom Middleware

Add middleware to the global chain by appending to the slice:

```go
middlewares = append(middlewares, myCustomMiddleware())
```

## The Router Function

The `Router()` function wires everything together:

```go
func Router(app app.AppInterface) rtr.RouterInterface {
    r := rtr.NewRouter()
    r.AddBeforeMiddlewares(globalMiddlewares(app))
    for _, route := range routes(app) {
        r.AddRoute(route)
    }
    return r
}
```

This is called from `cmd/server/main.go` and passed to the web server:

```go
server, err := websrv.Start(websrv.Options{
    Host:    cfg.GetAppHost(),
    Port:    cfg.GetAppPort(),
    URL:     cfg.GetAppUrl(),
    Handler: routes.Router(a).ServeHTTP,
})
```
