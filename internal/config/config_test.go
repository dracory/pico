package config

import (
	"testing"

	"github.com/dracory/env"
)

func setRequiredEnv(t *testing.T) {
	t.Helper()
	t.Setenv(KEY_APP_HOST, "localhost")
	t.Setenv(KEY_APP_PORT, "8080")
	t.Setenv(KEY_APP_ENVIRONMENT, "testing")
	t.Setenv(KEY_DEFAULT_DB_DRIVER, "sqlite")
	t.Setenv(KEY_DEFAULT_DB_DATABASE, ":memory:")
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

func TestNewFromEnv_Success(t *testing.T) {
	setRequiredEnv(t)
	cfg, err := NewFromEnv()
	if err != nil {
		t.Fatalf("NewFromEnv() failed: %v", err)
	}
	if cfg == nil {
		t.Fatal("NewFromEnv() returned nil config")
	}
	if cfg.GetAppHost() != "localhost" {
		t.Errorf("expected host=localhost, got %s", cfg.GetAppHost())
	}
	if cfg.GetAppPort() != "8080" {
		t.Errorf("expected port=8080, got %s", cfg.GetAppPort())
	}
	if cfg.GetDatabaseConnectionByName("default").GetDriver() != "sqlite" {
		t.Errorf("expected driver=sqlite, got %s", cfg.GetDatabaseConnectionByName("default").GetDriver())
	}
}

func TestNewFromEnv_MissingRequiredFields(t *testing.T) {
	t.Setenv(KEY_APP_HOST, "")
	t.Setenv(KEY_APP_PORT, "")
	t.Setenv(KEY_APP_ENVIRONMENT, "")
	t.Setenv(KEY_DEFAULT_DB_DRIVER, "")
	t.Setenv(KEY_DEFAULT_DB_DATABASE, "")

	_, err := NewFromEnv()
	if err == nil {
		t.Fatal("NewFromEnv() should fail with missing required fields")
	}

	verr, ok := err.(env.ValidationError)
	if !ok {
		t.Fatalf("expected env.ValidationError, got %T", err)
	}

	if len(verr.Errors()) == 0 {
		t.Error("expected validation errors, got none")
	}
}

func TestNewFromEnv_PostgresMissingConnectionDetails(t *testing.T) {
	t.Setenv(KEY_APP_HOST, "localhost")
	t.Setenv(KEY_APP_PORT, "8080")
	t.Setenv(KEY_APP_ENVIRONMENT, "testing")
	t.Setenv(KEY_DEFAULT_DB_DRIVER, "postgres")
	t.Setenv(KEY_DEFAULT_DB_DATABASE, "testdb")
	_, err := NewFromEnv()
	if err == nil {
		t.Fatal("NewFromEnv() should fail when postgres driver missing connection details")
	}

	verr, ok := err.(env.ValidationError)
	if !ok {
		t.Fatalf("expected env.ValidationError, got %T", err)
	}

	if len(verr.Errors()) < 4 {
		t.Errorf("expected at least 4 validation errors for postgres, got %d", len(verr.Errors()))
	}
}

func TestNewFromEnv_AppDebugMode(t *testing.T) {
	setRequiredEnv(t)
	t.Setenv(KEY_APP_DEBUG, "true")
	cfg, err := NewFromEnv()
	if err != nil {
		t.Fatalf("NewFromEnv() failed: %v", err)
	}

	if !cfg.GetAppDebug() {
		t.Error("GetAppDebug() = false, want true")
	}
}

func TestNewFromEnv_AppNameAndUrl(t *testing.T) {
	setRequiredEnv(t)
	t.Setenv(KEY_APP_NAME, "TestApp")
	t.Setenv(KEY_APP_URL, "http://test.example.com")
	cfg, err := NewFromEnv()
	if err != nil {
		t.Fatalf("NewFromEnv() failed: %v", err)
	}

	if cfg.GetAppName() != "TestApp" {
		t.Errorf("GetAppName() = %q, want TestApp", cfg.GetAppName())
	}
	if cfg.GetAppUrl() != "http://test.example.com" {
		t.Errorf("GetAppUrl() = %q, want http://test.example.com", cfg.GetAppUrl())
	}
}

func TestEnvironmentChecks(t *testing.T) {
	tests := []struct {
		env   string
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

func TestDatabaseConfig_SQLite(t *testing.T) {
	setRequiredEnv(t)
	cfg, err := NewFromEnv()
	if err != nil {
		t.Fatalf("NewFromEnv() error: %v", err)
	}

	conn := cfg.GetDatabaseConnectionByName("default")
	if conn.GetDriver() != "sqlite" {
		t.Errorf("expected driver sqlite, got %q", conn.GetDriver())
	}
	if conn.GetDatabase() != ":memory:" {
		t.Errorf("expected database :memory:, got %q", conn.GetDatabase())
	}
	if conn.GetMaxOpenConns() != 1 {
		t.Errorf("expected max open conns 1 for sqlite, got %d", conn.GetMaxOpenConns())
	}
	if conn.GetMaxIdleConns() != 1 {
		t.Errorf("expected max idle conns 1 for sqlite, got %d", conn.GetMaxIdleConns())
	}
	if cfg.GetDatabaseDefaultConnection() != "default" {
		t.Errorf("expected default connection 'default', got %q", cfg.GetDatabaseDefaultConnection())
	}
}

func TestDatabaseConfig_Postgres(t *testing.T) {
	t.Setenv(KEY_APP_HOST, "localhost")
	t.Setenv(KEY_APP_PORT, "8080")
	t.Setenv(KEY_APP_ENVIRONMENT, "testing")
	t.Setenv(KEY_DEFAULT_DB_DRIVER, "postgres")
	t.Setenv(KEY_DEFAULT_DB_HOST, "localhost")
	t.Setenv(KEY_DEFAULT_DB_PORT, "5432")
	t.Setenv(KEY_DEFAULT_DB_DATABASE, "picodb")
	t.Setenv(KEY_DEFAULT_DB_USERNAME, "user")
	t.Setenv(KEY_DEFAULT_DB_PASSWORD, "pass")
	cfg, err := NewFromEnv()
	if err != nil {
		t.Fatalf("NewFromEnv() error: %v", err)
	}

	conn := cfg.GetDatabaseConnectionByName("default")
	if conn.GetDriver() != "postgres" {
		t.Errorf("expected driver postgres, got %q", conn.GetDriver())
	}
	if conn.GetMaxOpenConns() != 25 {
		t.Errorf("expected max open conns 25 for postgres, got %d", conn.GetMaxOpenConns())
	}
	if conn.GetMaxIdleConns() != 5 {
		t.Errorf("expected max idle conns 5 for postgres, got %d", conn.GetMaxIdleConns())
	}
	if conn.GetSSLMode() != "require" {
		t.Errorf("expected ssl mode require, got %q", conn.GetSSLMode())
	}
}

func TestDatabaseConfig_CustomPoolSettings(t *testing.T) {
	setRequiredEnv(t)
	t.Setenv(KEY_DEFAULT_DB_MAX_OPEN_CONNS, "10")
	t.Setenv(KEY_DEFAULT_DB_MAX_IDLE_CONNS, "5")
	t.Setenv(KEY_DEFAULT_DB_CONN_MAX_LIFETIME_SECONDS, "600")
	t.Setenv(KEY_DEFAULT_DB_CONN_MAX_IDLE_TIME_SECONDS, "60")
	t.Setenv(KEY_DEFAULT_DB_CHARSET, "utf8")
	t.Setenv(KEY_DEFAULT_DB_TIMEZONE, "America/New_York")
	cfg, err := NewFromEnv()
	if err != nil {
		t.Fatalf("NewFromEnv() failed: %v", err)
	}

	conn := cfg.GetDatabaseConnectionByName("default")
	if conn.GetMaxOpenConns() != 1 {
		t.Errorf("expected max open conns 1 for sqlite override, got %d", conn.GetMaxOpenConns())
	}
	if conn.GetMaxIdleConns() != 1 {
		t.Errorf("expected max idle conns 1 for sqlite override, got %d", conn.GetMaxIdleConns())
	}
	if conn.GetConnMaxLifetimeSeconds() != 30 {
		t.Errorf("expected conn max lifetime 30 for sqlite override, got %d", conn.GetConnMaxLifetimeSeconds())
	}
	if conn.GetConnMaxIdleTimeSeconds() != 60 {
		t.Errorf("expected conn max idle time 60, got %d", conn.GetConnMaxIdleTimeSeconds())
	}
	if conn.GetCharset() != "utf8" {
		t.Errorf("expected charset utf8, got %q", conn.GetCharset())
	}
	if conn.GetTimezone() != "America/New_York" {
		t.Errorf("expected timezone America/New_York, got %q", conn.GetTimezone())
	}
}

func TestDatabaseConfig_Connections(t *testing.T) {
	setRequiredEnv(t)
	cfg, err := NewFromEnv()
	if err != nil {
		t.Fatalf("NewFromEnv() failed: %v", err)
	}

	conns := cfg.GetDatabaseConnections()
	if len(conns) == 0 {
		t.Fatal("expected at least one connection")
	}

	conn := cfg.GetDatabaseConnectionByName("default")
	if conn == nil {
		t.Fatal("expected default connection to exist")
	}
	if conn.GetName() != "default" {
		t.Errorf("expected name 'default', got %q", conn.GetName())
	}
	if conn.GetDriver() != "sqlite" {
		t.Errorf("expected driver sqlite, got %q", conn.GetDriver())
	}

	missing := cfg.GetDatabaseConnectionByName("nonexistent")
	if missing != nil {
		t.Error("expected nil for nonexistent connection")
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
	setRequiredEnv(t)
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
