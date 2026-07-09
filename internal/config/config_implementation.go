package config

import (
	"github.com/dracory/env"
)

type envValidator = env.Validator

type configImplementation struct {
	appName  string
	appEnv   string
	appHost  string
	appPort  string
	appUrl   string
	appDebug bool

	databaseDefaultConnection      string
	databaseConnections            map[string]DatabaseConnectionConfigInterface
	databaseDriver                 string
	databaseHost                   string
	databasePort                   string
	databaseName                   string
	databaseUsername               string
	databasePassword               string
	databaseSSLMode                string
	databaseCharset                string
	databaseTimezone               string
	databaseDSN                    string
	databasePrefix                 string
	databaseMaxOpenConns           int
	databaseMaxIdleConns           int
	databaseConnMaxLifetimeSeconds int
	databaseConnMaxIdleTimeSeconds int
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

func (c *configImplementation) setDatabaseConfig(s databaseSettings) {
	c.databaseDefaultConnection = s.defaultConnection
	c.databaseConnections = s.connections
	c.databaseDriver = s.driver
	c.databaseHost = s.host
	c.databasePort = s.port
	c.databaseName = s.name
	c.databaseUsername = s.user
	c.databasePassword = s.pass
	c.databaseSSLMode = s.sslMode
	c.databaseCharset = s.charset
	c.databaseTimezone = s.timezone
	c.databaseDSN = s.dsn
	c.databasePrefix = s.prefix
	c.databaseMaxOpenConns = s.maxOpenConns
	c.databaseMaxIdleConns = s.maxIdleConns
	c.databaseConnMaxLifetimeSeconds = int(s.connMaxLifetime.Seconds())
	c.databaseConnMaxIdleTimeSeconds = int(s.connMaxIdleTime.Seconds())
}

func (c *configImplementation) defaultConnection() DatabaseConnectionConfigInterface {
	if c == nil {
		return nil
	}
	if conn, ok := c.databaseConnections[c.databaseDefaultConnection]; ok && conn != nil {
		return conn
	}
	return nil
}

func (c *configImplementation) SetDatabaseDefaultConnection(v string) { c.databaseDefaultConnection = v }
func (c *configImplementation) GetDatabaseDefaultConnection() string  { return c.databaseDefaultConnection }

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

func (c *configImplementation) SetDatabaseDriver(v string) { c.databaseDriver = v }
func (c *configImplementation) GetDatabaseDriver() string {
	if conn := c.defaultConnection(); conn != nil {
		return conn.GetDriver()
	}
	return c.databaseDriver
}

func (c *configImplementation) SetDatabaseHost(v string) { c.databaseHost = v }
func (c *configImplementation) GetDatabaseHost() string {
	if conn := c.defaultConnection(); conn != nil {
		return conn.GetHost()
	}
	return c.databaseHost
}

func (c *configImplementation) SetDatabasePort(v string) { c.databasePort = v }
func (c *configImplementation) GetDatabasePort() string {
	if conn := c.defaultConnection(); conn != nil {
		return conn.GetPort()
	}
	return c.databasePort
}

func (c *configImplementation) SetDatabaseName(v string) { c.databaseName = v }
func (c *configImplementation) GetDatabaseName() string {
	if conn := c.defaultConnection(); conn != nil {
		return conn.GetDatabase()
	}
	return c.databaseName
}

func (c *configImplementation) SetDatabaseUsername(v string) { c.databaseUsername = v }
func (c *configImplementation) GetDatabaseUsername() string {
	if conn := c.defaultConnection(); conn != nil {
		return conn.GetUsername()
	}
	return c.databaseUsername
}

func (c *configImplementation) SetDatabasePassword(v string) { c.databasePassword = v }
func (c *configImplementation) GetDatabasePassword() string {
	if conn := c.defaultConnection(); conn != nil {
		return conn.GetPassword()
	}
	return c.databasePassword
}

func (c *configImplementation) SetDatabaseSSLMode(v string) { c.databaseSSLMode = v }
func (c *configImplementation) GetDatabaseSSLMode() string {
	if conn := c.defaultConnection(); conn != nil {
		return conn.GetSSLMode()
	}
	return c.databaseSSLMode
}

func (c *configImplementation) SetDatabaseMaxOpenConns(v int) { c.databaseMaxOpenConns = v }
func (c *configImplementation) GetDatabaseMaxOpenConns() int   { return c.databaseMaxOpenConns }

func (c *configImplementation) SetDatabaseMaxIdleConns(v int) { c.databaseMaxIdleConns = v }
func (c *configImplementation) GetDatabaseMaxIdleConns() int   { return c.databaseMaxIdleConns }

func (c *configImplementation) SetDatabaseConnMaxLifetimeSeconds(v int) { c.databaseConnMaxLifetimeSeconds = v }
func (c *configImplementation) GetDatabaseConnMaxLifetimeSeconds() int   { return c.databaseConnMaxLifetimeSeconds }

func (c *configImplementation) SetDatabaseConnMaxIdleTimeSeconds(v int) { c.databaseConnMaxIdleTimeSeconds = v }
func (c *configImplementation) GetDatabaseConnMaxIdleTimeSeconds() int   { return c.databaseConnMaxIdleTimeSeconds }

func (c *configImplementation) SetDatabaseCharset(v string) { c.databaseCharset = v }
func (c *configImplementation) GetDatabaseCharset() string {
	if conn := c.defaultConnection(); conn != nil {
		return conn.GetCharset()
	}
	return c.databaseCharset
}

func (c *configImplementation) SetDatabaseTimezone(v string) { c.databaseTimezone = v }
func (c *configImplementation) GetDatabaseTimezone() string {
	if conn := c.defaultConnection(); conn != nil {
		return conn.GetTimezone()
	}
	return c.databaseTimezone
}

func (c *configImplementation) SetDatabaseDSN(v string) { c.databaseDSN = v }
func (c *configImplementation) GetDatabaseDSN() string {
	if conn := c.defaultConnection(); conn != nil {
		return conn.GetDSN()
	}
	return c.databaseDSN
}

func (c *configImplementation) SetDatabasePrefix(v string) { c.databasePrefix = v }
func (c *configImplementation) GetDatabasePrefix() string {
	if conn := c.defaultConnection(); conn != nil {
		return conn.GetPrefix()
	}
	return c.databasePrefix
}
