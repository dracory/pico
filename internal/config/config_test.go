package config

import (
	"os"
	"testing"
)

func setTestEnv(t *testing.T) {
	t.Helper()
	t.Setenv(KEY_APP_NAME, "PicoTest")
	t.Setenv(KEY_APP_ENVIRONMENT, APP_ENVIRONMENT_TESTING)
	t.Setenv(KEY_APP_DEBUG, "true")
	t.Setenv(KEY_APP_HOST, "127.0.0.1")
	t.Setenv(KEY_APP_PORT, "8080")
	t.Setenv(KEY_APP_URL, "http://localhost:8080")
	t.Setenv(KEY_DB_DRIVER, "sqlite")
	t.Setenv(KEY_DB_DATABASE, ":memory:")
}

func TestNew(t *testing.T) {
	cfg := New()
	if cfg == nil {
		t.Fatal("New() returned nil")
	}
	if cfg.GetAppName() != "" {
		t.Errorf("expected empty app name, got %q", cfg.GetAppName())
	}
}

func TestNewFromEnv(t *testing.T) {
	setTestEnv(t)

	cfg, err := NewFromEnv()
	if err != nil {
		t.Fatalf("NewFromEnv() error: %v", err)
	}
	if cfg == nil {
		t.Fatal("NewFromEnv() returned nil config")
	}

	if cfg.GetAppName() != "PicoTest" {
		t.Errorf("expected app name PicoTest, got %q", cfg.GetAppName())
	}
	if cfg.GetAppEnv() != APP_ENVIRONMENT_TESTING {
		t.Errorf("expected env testing, got %q", cfg.GetAppEnv())
	}
	if cfg.GetAppHost() != "127.0.0.1" {
		t.Errorf("expected host 127.0.0.1, got %q", cfg.GetAppHost())
	}
	if cfg.GetAppPort() != "8080" {
		t.Errorf("expected port 8080, got %q", cfg.GetAppPort())
	}
	if cfg.GetAppUrl() != "http://localhost:8080" {
		t.Errorf("expected url http://localhost:8080, got %q", cfg.GetAppUrl())
	}
	if !cfg.GetAppDebug() {
		t.Error("expected debug true")
	}
}

func TestNewFromEnvMissingRequired(t *testing.T) {
	os.Clearenv()

	_, err := NewFromEnv()
	if err == nil {
		t.Fatal("expected error when required env vars are missing")
	}
}

func TestEnvironmentChecks(t *testing.T) {
	tests := []struct {
		env  string
		check func(ConfigInterface) bool
		name  string
	}{
		{APP_ENVIRONMENT_DEVELOPMENT, func(c ConfigInterface) bool { return c.IsEnvDevelopment() }, "IsEnvDevelopment"},
		{APP_ENVIRONMENT_LOCAL, func(c ConfigInterface) bool { return c.IsEnvLocal() }, "IsEnvLocal"},
		{APP_ENVIRONMENT_PRODUCTION, func(c ConfigInterface) bool { return c.IsEnvProduction() }, "IsEnvProduction"},
		{APP_ENVIRONMENT_STAGING, func(c ConfigInterface) bool { return c.IsEnvStaging() }, "IsEnvStaging"},
		{APP_ENVIRONMENT_TESTING, func(c ConfigInterface) bool { return c.IsEnvTesting() }, "IsEnvTesting"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := New()
			cfg.SetAppEnv(tt.env)
			if !tt.check(cfg) {
				t.Errorf("%s() should return true for env %q", tt.name, tt.env)
			}
			cfg.SetAppEnv("other")
			if tt.check(cfg) {
				t.Errorf("%s() should return false for env other", tt.name)
			}
		})
	}
}

func TestDatabaseConfigFromEnv(t *testing.T) {
	setTestEnv(t)

	cfg, err := NewFromEnv()
	if err != nil {
		t.Fatalf("NewFromEnv() error: %v", err)
	}

	if cfg.GetDatabaseDriver() != "sqlite" {
		t.Errorf("expected driver sqlite, got %q", cfg.GetDatabaseDriver())
	}
	if cfg.GetDatabaseName() != ":memory:" {
		t.Errorf("expected database :memory:, got %q", cfg.GetDatabaseName())
	}
	if cfg.GetDatabaseMaxOpenConns() != 1 {
		t.Errorf("expected max open conns 1 for sqlite, got %d", cfg.GetDatabaseMaxOpenConns())
	}
	if cfg.GetDatabaseMaxIdleConns() != 1 {
		t.Errorf("expected max idle conns 1 for sqlite, got %d", cfg.GetDatabaseMaxIdleConns())
	}
	if cfg.GetDatabaseDefaultConnection() != "default" {
		t.Errorf("expected default connection 'default', got %q", cfg.GetDatabaseDefaultConnection())
	}
}

func TestDatabaseConfigNonSQLite(t *testing.T) {
	t.Setenv(KEY_APP_NAME, "PicoTest")
	t.Setenv(KEY_APP_ENVIRONMENT, APP_ENVIRONMENT_TESTING)
	t.Setenv(KEY_APP_HOST, "127.0.0.1")
	t.Setenv(KEY_APP_PORT, "8080")
	t.Setenv(KEY_DB_DRIVER, "postgres")
	t.Setenv(KEY_DB_HOST, "localhost")
	t.Setenv(KEY_DB_PORT, "5432")
	t.Setenv(KEY_DB_DATABASE, "picodb")
	t.Setenv(KEY_DB_USERNAME, "user")
	t.Setenv(KEY_DB_PASSWORD, "pass")

	cfg, err := NewFromEnv()
	if err != nil {
		t.Fatalf("NewFromEnv() error: %v", err)
	}

	if cfg.GetDatabaseDriver() != "postgres" {
		t.Errorf("expected driver postgres, got %q", cfg.GetDatabaseDriver())
	}
	if cfg.GetDatabaseMaxOpenConns() != 25 {
		t.Errorf("expected max open conns 25 for postgres, got %d", cfg.GetDatabaseMaxOpenConns())
	}
	if cfg.GetDatabaseMaxIdleConns() != 5 {
		t.Errorf("expected max idle conns 5 for postgres, got %d", cfg.GetDatabaseMaxIdleConns())
	}
	if cfg.GetDatabaseSSLMode() != "require" {
		t.Errorf("expected ssl mode require, got %q", cfg.GetDatabaseSSLMode())
	}
}

func TestSetters(t *testing.T) {
	cfg := New()

	cfg.SetAppName("MyApp")
	cfg.SetAppEnv("production")
	cfg.SetAppHost("0.0.0.0")
	cfg.SetAppPort("3000")
	cfg.SetAppUrl("https://myapp.com")
	cfg.SetAppDebug(false)

	if cfg.GetAppName() != "MyApp" {
		t.Errorf("expected MyApp, got %q", cfg.GetAppName())
	}
	if cfg.GetAppEnv() != "production" {
		t.Errorf("expected production, got %q", cfg.GetAppEnv())
	}
	if cfg.GetAppHost() != "0.0.0.0" {
		t.Errorf("expected 0.0.0.0, got %q", cfg.GetAppHost())
	}
	if cfg.GetAppPort() != "3000" {
		t.Errorf("expected 3000, got %q", cfg.GetAppPort())
	}
	if cfg.GetAppUrl() != "https://myapp.com" {
		t.Errorf("expected https://myapp.com, got %q", cfg.GetAppUrl())
	}
	if cfg.GetAppDebug() {
		t.Error("expected debug false")
	}
}

func TestDatabaseNeatConfig(t *testing.T) {
	setTestEnv(t)

	cfg, err := NewFromEnv()
	if err != nil {
		t.Fatalf("NewFromEnv() error: %v", err)
	}

	neatCfg := DatabaseNeatConfig(cfg)
	if neatCfg.Default != "default" {
		t.Errorf("expected default connection 'default', got %q", neatCfg.Default)
	}
	conn, ok := neatCfg.Connections["default"]
	if !ok {
		t.Fatal("expected default connection in neat config")
	}
	if conn.Driver != "sqlite" {
		t.Errorf("expected driver sqlite, got %q", conn.Driver)
	}
	if conn.Database != ":memory:" {
		t.Errorf("expected database :memory:, got %q", conn.Database)
	}
}

func TestDatabaseNeatConfigNil(t *testing.T) {
	neatCfg := DatabaseNeatConfig(nil)
	if neatCfg.Default != "" {
		t.Errorf("expected empty default for nil config, got %q", neatCfg.Default)
	}
}

func TestPortToInt(t *testing.T) {
	tests := []struct {
		port     string
		driver   string
		expected int
	}{
		{"3306", "mysql", 3306},
		{"5432", "postgres", 5432},
		{"", "mysql", 3306},
		{"", "postgres", 5432},
		{"", "sqlite", 0},
		{"invalid", "mysql", 0},
		{"1433", "sqlserver", 1433},
		{"1521", "oracle", 1521},
	}

	for _, tt := range tests {
		result := portToInt(tt.port, tt.driver)
		if result != tt.expected {
			t.Errorf("portToInt(%q, %q) = %d, expected %d", tt.port, tt.driver, result, tt.expected)
		}
	}
}

func TestGetVersion(t *testing.T) {
	v := GetVersion()
	if v == "" {
		t.Error("GetVersion() returned empty string")
	}
}
