package config

import (
	"github.com/dracory/env"
)

// envValidator is a local alias for env.Validator for use in config loaders.
type envValidator = env.Validator

// ============================================================================
// Main Config Interface
// ============================================================================

// ConfigInterface defines the contract for application configuration.
// It composes all domain-specific configuration interfaces.
type ConfigInterface interface {
	AppConfigInterface
	DatabaseConfigInterface
}

// ============================================================================
// App Config Interface
// ============================================================================

// AppConfigInterface defines application-level configuration methods.
type AppConfigInterface interface {
	SetAppName(string)
	GetAppName() string

	SetAppEnv(string)
	GetAppEnv() string

	SetAppHost(string)
	GetAppHost() string

	SetAppPort(string)
	GetAppPort() string

	SetAppUrl(string)
	GetAppUrl() string

	SetAppDebug(bool)
	GetAppDebug() bool

	// Environment helpers
	IsEnvDevelopment() bool
	IsEnvLocal() bool
	IsEnvProduction() bool
	IsEnvStaging() bool
	IsEnvTesting() bool
}

// ============================================================================
// Database Connection Config Interface
// ============================================================================

// DatabaseConnectionConfigInterface defines a single database connection.
// It provides a Laravel-like connection configuration while keeping
// compatibility with the existing single-database getters.
type DatabaseConnectionConfigInterface interface {
	GetName() string
	GetDriver() string
	GetHost() string
	GetPort() string
	GetDatabase() string
	GetUsername() string
	GetPassword() string
	GetSSLMode() string
	GetCharset() string
	GetTimezone() string
	GetDSN() string
	GetPrefix() string
	GetMaxOpenConns() int
	GetMaxIdleConns() int
	GetConnMaxLifetimeSeconds() int
	GetConnMaxIdleTimeSeconds() int
}

// ============================================================================
// Database Config Interface
// ============================================================================

// DatabaseConfigInterface defines database configuration methods.
// It supports multi-connection configuration via DatabaseConnectionConfigInterface.
type DatabaseConfigInterface interface {
	SetDatabaseDefaultConnection(string)
	GetDatabaseDefaultConnection() string

	GetDatabaseConnections() []DatabaseConnectionConfigInterface
	GetDatabaseConnectionByName(name string) DatabaseConnectionConfigInterface
}

// configImplementation holds all configuration values.
type configImplementation struct {
	// App configuration
	appName  string
	appEnv   string
	appHost  string
	appPort  string
	appUrl   string
	appDebug bool

	// Database configuration
	databaseConnections map[string]DatabaseConnectionConfigInterface
}

// New constructs a new configuration instance.
func New() ConfigInterface {
	return &configImplementation{}
}

// NewFromEnv constructs a configuration instance populated from environment variables.
func NewFromEnv() (ConfigInterface, error) {
	env.Load(".env")

	v := &envValidator{}
	cfg := &configImplementation{}

	cfg.setAppConfig(appConfig(v))
	cfg.setDatabaseConfig(databaseConfig(v))

	if err := v.Err(); err != nil {
		return nil, err
	}

	return cfg, nil
}

// Ensure configImplementation satisfies ConfigInterface
var _ ConfigInterface = (*configImplementation)(nil)

// ============================================================================
// App Config Implementation
// ============================================================================

func (c *configImplementation) setAppConfig(s appSettings) {
	c.appName = s.name
	c.appUrl = s.url
	c.appHost = s.host
	c.appPort = s.port
	c.appEnv = s.env
	c.appDebug = s.debug
}

func (c *configImplementation) SetAppName(v string) { c.appName = v }

func (c *configImplementation) GetAppName() string { return c.appName }

func (c *configImplementation) SetAppEnv(v string) { c.appEnv = v }

func (c *configImplementation) GetAppEnv() string { return c.appEnv }

func (c *configImplementation) SetAppHost(v string) { c.appHost = v }

func (c *configImplementation) GetAppHost() string { return c.appHost }

func (c *configImplementation) SetAppPort(v string) { c.appPort = v }

func (c *configImplementation) GetAppPort() string { return c.appPort }

func (c *configImplementation) SetAppUrl(v string) { c.appUrl = v }

func (c *configImplementation) GetAppUrl() string { return c.appUrl }

func (c *configImplementation) SetAppDebug(v bool) { c.appDebug = v }

func (c *configImplementation) GetAppDebug() bool { return c.appDebug }

func (c *configImplementation) IsEnvDevelopment() bool { return c.appEnv == "development" }

func (c *configImplementation) IsEnvLocal() bool { return c.appEnv == "local" }

func (c *configImplementation) IsEnvProduction() bool { return c.appEnv == "production" }

func (c *configImplementation) IsEnvStaging() bool { return c.appEnv == "staging" }

func (c *configImplementation) IsEnvTesting() bool { return c.appEnv == "testing" }

// ============================================================================
// Database Config Implementation
// ============================================================================

func (c *configImplementation) setDatabaseConfig(conns map[string]DatabaseConnectionConfigInterface) {
	c.databaseConnections = conns
}

// defaultConnection returns the default database connection configuration.
// It falls back to nil if no connections are configured.
func (c *configImplementation) defaultConnection() DatabaseConnectionConfigInterface {
	if c == nil {
		return nil
	}
	for _, conn := range c.databaseConnections {
		if s, ok := conn.(*databaseConnectionSettings); ok && s.isDefault {
			return conn
		}
	}
	return nil
}

// ensureDefaultConnection returns the default connection, creating one if none exists.
func (c *configImplementation) ensureDefaultConnection() *databaseConnectionSettings {
	if c.databaseConnections == nil {
		c.databaseConnections = map[string]DatabaseConnectionConfigInterface{}
	}
	for _, conn := range c.databaseConnections {
		if s, ok := conn.(*databaseConnectionSettings); ok && s.isDefault {
			return s
		}
	}
	s := &databaseConnectionSettings{name: defaultConnectionName, isDefault: true}
	c.databaseConnections[defaultConnectionName] = s
	return s
}

func (c *configImplementation) SetDatabaseDefaultConnection(v string) {
	if v == "" {
		return
	}
	for name, conn := range c.databaseConnections {
		s, ok := conn.(*databaseConnectionSettings)
		if !ok || !s.isDefault {
			continue
		}
		if name == v {
			return
		}
		s.name = v
		c.databaseConnections[v] = conn
		delete(c.databaseConnections, name)
		return
	}
}
func (c *configImplementation) GetDatabaseDefaultConnection() string {
	if conn := c.defaultConnection(); conn != nil {
		return conn.GetName()
	}
	return ""
}

func (c *configImplementation) GetDatabaseConnections() []DatabaseConnectionConfigInterface {
	if c == nil || len(c.databaseConnections) == 0 {
		return nil
	}
	conns := make([]DatabaseConnectionConfigInterface, 0, len(c.databaseConnections))
	for _, conn := range c.databaseConnections {
		conns = append(conns, conn)
	}
	return conns
}

func (c *configImplementation) GetDatabaseConnectionByName(name string) DatabaseConnectionConfigInterface {
	if c == nil || c.databaseConnections == nil {
		return nil
	}
	return c.databaseConnections[name]
}
