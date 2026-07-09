package routes

import (
	"net/http"
	"time"

	"project/internal/app"

	"github.com/dracory/rtr"
	rtrMiddleware "github.com/dracory/rtr/middlewares"
)

func globalMiddlewares(app app.AppInterface) []rtr.MiddlewareInterface {
	middlewares := []rtr.MiddlewareInterface{
		rtrMiddleware.CompressMiddleware(5, "text/html", "text/css"),
		rtrMiddleware.GetHead(),
		rtrMiddleware.CleanPathMiddleware(),
		rtrMiddleware.RedirectSlashesMiddleware(),
		rtrMiddleware.TimeoutMiddleware(30 * time.Second),
	}

	if app.GetConfig() != nil && !app.GetConfig().IsEnvTesting() {
		middlewares = append(middlewares,
			rtrMiddleware.LoggerMiddleware(),
			rtrMiddleware.RecoveryMiddleware(),
		)
	}

	return middlewares
}

func routes(app app.AppInterface) []rtr.RouteInterface {
	return []rtr.RouteInterface{
		rtr.Get("/", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"name":"pico","status":"running"}`))
		}),
	}
}

func Router(app app.AppInterface) rtr.RouterInterface {
	r := rtr.NewRouter()

	r.AddBeforeMiddlewares(globalMiddlewares(app))

	for _, route := range routes(app) {
		r.AddRoute(route)
	}

	return r
}
