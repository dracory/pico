package app

import (
	"testing"

	"project/internal/config"
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
	cfg.SetDatabaseName(":memory:")
	return cfg
}

func TestNewNilConfig(t *testing.T) {
	_, err := New(nil)
	if err == nil {
		t.Fatal("expected error for nil config")
	}
}

func TestAppSettersAndGetters(t *testing.T) {
	a := &appImplementation{}

	a.SetConfig(newTestConfig())
	if a.GetConfig() == nil {
		t.Fatal("GetConfig() returned nil after SetConfig")
	}
	if a.GetConfig().GetAppName() != "PicoTest" {
		t.Errorf("expected app name PicoTest, got %q", a.GetConfig().GetAppName())
	}
}

func TestAppLogger(t *testing.T) {
	a := &appImplementation{}

	if a.GetLogger() != nil {
		t.Fatal("expected nil logger before SetLogger")
	}

	a.SetLogger(nil)
	if a.GetLogger() != nil {
		t.Fatal("expected nil logger after SetLogger(nil)")
	}
}

func TestAppCloseNil(t *testing.T) {
	var a *appImplementation
	err := a.Close()
	if err != nil {
		t.Errorf("Close() on nil app should return nil, got %v", err)
	}
}

func TestAppCloseNoDB(t *testing.T) {
	a := &appImplementation{}
	err := a.Close()
	if err != nil {
		t.Errorf("Close() with no database should return nil, got %v", err)
	}
}
