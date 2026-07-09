package routes

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"project/internal/app"
	"project/internal/config"

	"github.com/dracory/rtr"
	_ "modernc.org/sqlite"
)

func newTestApp() app.AppInterface {
	cfg := config.New()
	cfg.SetAppEnv("testing")
	cfg.SetAppName("PicoTest")
	cfg.SetAppHost("127.0.0.1")
	cfg.SetAppPort("8080")
	cfg.SetAppUrl("http://localhost:8080")
	cfg.SetDatabaseDriver("sqlite")
	cfg.SetDatabaseName(":memory:")

	a, err := app.New(cfg)
	if err != nil {
		panic("failed to create test app: " + err.Error())
	}
	return a
}

func TestRouterNotNil(t *testing.T) {
	a := newTestApp()
	r := Router(a)
	if r == nil {
		t.Fatal("Router() returned nil")
	}
}

func TestRouterHasRoutes(t *testing.T) {
	a := newTestApp()
	r := Router(a)
	if len(r.GetRoutes()) == 0 {
		t.Fatal("Router() has no routes")
	}
}

func TestRootEndpoint(t *testing.T) {
	a := newTestApp()
	r := Router(a)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()

	r.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", rec.Code)
	}

	body := rec.Body.String()
	if !strings.Contains(body, "pico") {
		t.Errorf("expected body to contain 'pico', got %q", body)
	}
	if !strings.Contains(body, "running") {
		t.Errorf("expected body to contain 'running', got %q", body)
	}
}

func TestGlobalMiddlewaresTesting(t *testing.T) {
	mw := globalMiddlewares(newTestApp())
	if len(mw) == 0 {
		t.Fatal("expected at least some middlewares")
	}

	for _, m := range mw {
		if m == nil {
			t.Fatal("found nil middleware")
		}
	}
}

func TestRoutesList(t *testing.T) {
	a := newTestApp()
	rs := routes(a)
	if len(rs) == 0 {
		t.Fatal("routes() returned empty list")
	}

	for _, r := range rs {
		if r == nil {
			t.Fatal("found nil route")
		}
	}
}

func TestRouterImplementsInterface(t *testing.T) {
	a := newTestApp()
	r := Router(a)

	var _ rtr.RouterInterface = r
}
