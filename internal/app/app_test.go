package app

import (
	"fmt"
	"testing"
	"time"

	"project/internal/config"

	_ "modernc.org/sqlite"
)

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

func TestNew_NilConfig(t *testing.T) {
	_, err := New(nil)
	if err == nil {
		t.Fatal("expected error for nil config")
	}
}

func TestNew_SetsDefaultLogger(t *testing.T) {
	cfg := newTestConfig()

	a, err := New(cfg)
	if err != nil {
		t.Fatalf("app.New returned error: %v", err)
	}
	defer a.Close()

	if a.GetLogger() == nil {
		t.Fatal("expected app logger to be non-nil right after app.New")
	}
}

func TestNew_SetsConfig(t *testing.T) {
	cfg := newTestConfig()

	a, err := New(cfg)
	if err != nil {
		t.Fatalf("app.New returned error: %v", err)
	}
	defer a.Close()

	if a.GetConfig() == nil {
		t.Fatal("expected config to be non-nil")
	}
	if a.GetConfig().GetAppName() != "PicoTest" {
		t.Errorf("expected app name PicoTest, got %q", a.GetConfig().GetAppName())
	}
}

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

func TestNew_SetsNeatDatabase(t *testing.T) {
	cfg := newTestConfig()

	a, err := New(cfg)
	if err != nil {
		t.Fatalf("app.New returned error: %v", err)
	}
	defer a.Close()

	if a.GetNeatDatabase() == nil {
		t.Fatal("expected neat database to be non-nil")
	}
}

func TestClose_NilReceiverDoesNotPanic(t *testing.T) {
	var a *appImplementation
	if err := a.Close(); err != nil {
		t.Fatalf("expected nil error, got: %v", err)
	}
}

func TestClose_NilDatabaseDoesNotPanic(t *testing.T) {
	a := &appImplementation{}
	if err := a.Close(); err != nil {
		t.Fatalf("expected nil error, got: %v", err)
	}

	if err := a.Close(); err != nil {
		t.Fatalf("expected nil error on repeated close, got: %v", err)
	}
}

func TestClose_ClosesNeatDatabase(t *testing.T) {
	cfg := newTestConfig()

	a, err := New(cfg)
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}

	impl, ok := a.(*appImplementation)
	if !ok {
		t.Fatal("expected *appImplementation")
	}

	if impl.neatDB == nil {
		t.Fatal("expected neatDB to be set")
	}
	if impl.db == nil {
		t.Fatal("expected sql db to be set")
	}

	if err := a.Close(); err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}

	if impl.neatDB != nil {
		t.Error("expected neatDB to be nil after close")
	}
	if impl.db != nil {
		t.Error("expected sql db to be nil after close")
	}
}

func TestSetLogger(t *testing.T) {
	a := &appImplementation{}

	if a.GetLogger() != nil {
		t.Fatal("expected nil logger before SetLogger")
	}

	a.SetLogger(nil)
	if a.GetLogger() != nil {
		t.Fatal("expected nil logger after SetLogger(nil)")
	}
}

func TestSetConfig(t *testing.T) {
	a := &appImplementation{}

	a.SetConfig(newTestConfig())
	if a.GetConfig() == nil {
		t.Fatal("GetConfig() returned nil after SetConfig")
	}
	if a.GetConfig().GetAppName() != "PicoTest" {
		t.Errorf("expected app name PicoTest, got %q", a.GetConfig().GetAppName())
	}
}
