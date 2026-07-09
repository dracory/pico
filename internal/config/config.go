package config

import (
	"github.com/dracory/env"
)

type envValidator = env.Validator

type ConfigInterface interface {
	AppConfigInterface
	DatabaseConfigInterface
}

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

	IsEnvDevelopment() bool
	IsEnvLocal() bool
	IsEnvProduction() bool
	IsEnvStaging() bool
	IsEnvTesting() bool
}

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

type DatabaseConfigInterface interface {
	SetDatabaseDefaultConnection(string)
	GetDatabaseDefaultConnection() string

	GetDatabaseConnections() []DatabaseConnectionConfigInterface
	GetDatabaseConnectionByName(name string) DatabaseConnectionConfigInterface
}

type configImplementation struct {
	appName  string
	appEnv   string
	appHost  string
	appPort  string
	appUrl   string
	appDebug bool

	databaseConnections map[string]DatabaseConnectionConfigInterface
}

func New() ConfigInterface {
	return &configImplementation{}
}

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

var _ ConfigInterface = (*configImplementation)(nil)

func (c *configImplementation) setAppConfig(s appSettings) {
	c.appName = s.name
	c.appUrl = s.url
	c.appHost = s.host
	c.appPort = s.port
	c.appEnv = s.env
	c.appDebug = s.debug
}

func (c *configImplementation) SetAppName(v string) { c.appName = v }
func (c *configImplementation) GetAppName() string  { return c.appName }

func (c *configImplementation) SetAppEnv(v string) { c.appEnv = v }
func (c *configImplementation) GetAppEnv() string  { return c.appEnv }

func (c *configImplementation) SetAppHost(v string) { c.appHost = v }
func (c *configImplementation) GetAppHost() string  { return c.appHost }

func (c *configImplementation) SetAppPort(v string) { c.appPort = v }
func (c *configImplementation) GetAppPort() string  { return c.appPort }

func (c *configImplementation) SetAppUrl(v string) { c.appUrl = v }
func (c *configImplementation) GetAppUrl() string  { return c.appUrl }

func (c *configImplementation) SetAppDebug(v bool) { c.appDebug = v }
func (c *configImplementation) GetAppDebug() bool  { return c.appDebug }

func (c *configImplementation) IsEnvDevelopment() bool { return c.appEnv == "development" }
func (c *configImplementation) IsEnvLocal() bool       { return c.appEnv == "local" }
func (c *configImplementation) IsEnvProduction() bool  { return c.appEnv == "production" }
func (c *configImplementation) IsEnvStaging() bool     { return c.appEnv == "staging" }
func (c *configImplementation) IsEnvTesting() bool     { return c.appEnv == "testing" }

func (c *configImplementation) setDatabaseConfig(conns map[string]DatabaseConnectionConfigInterface) {
	c.databaseConnections = conns
}

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
