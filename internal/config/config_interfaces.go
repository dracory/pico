package config

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
}

type DatabaseConfigInterface interface {
	SetDatabaseDriver(string)
	GetDatabaseDriver() string

	SetDatabaseHost(string)
	GetDatabaseHost() string

	SetDatabasePort(string)
	GetDatabasePort() string

	SetDatabaseName(string)
	GetDatabaseName() string

	SetDatabaseUsername(string)
	GetDatabaseUsername() string

	SetDatabasePassword(string)
	GetDatabasePassword() string

	SetDatabaseSSLMode(string)
	GetDatabaseSSLMode() string

	SetDatabaseMaxOpenConns(int)
	GetDatabaseMaxOpenConns() int

	SetDatabaseMaxIdleConns(int)
	GetDatabaseMaxIdleConns() int

	SetDatabaseConnMaxLifetimeSeconds(int)
	GetDatabaseConnMaxLifetimeSeconds() int

	SetDatabaseConnMaxIdleTimeSeconds(int)
	GetDatabaseConnMaxIdleTimeSeconds() int

	SetDatabaseCharset(string)
	GetDatabaseCharset() string

	SetDatabaseTimezone(string)
	GetDatabaseTimezone() string

	SetDatabaseDSN(string)
	GetDatabaseDSN() string

	SetDatabasePrefix(string)
	GetDatabasePrefix() string

	SetDatabaseDefaultConnection(string)
	GetDatabaseDefaultConnection() string

	GetDatabaseConnections() []DatabaseConnectionConfigInterface
	GetDatabaseConnectionByName(name string) DatabaseConnectionConfigInterface
}
