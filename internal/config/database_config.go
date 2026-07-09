package config

import (
	"strconv"
	"strings"
	"time"

	"github.com/dracory/neat/database/db"
)

// databaseConfig reads database configuration from environment variables
// and returns a map of connection configurations keyed by connection name.
func databaseConfig(env *envValidator) map[string]DatabaseConnectionConfigInterface {
	conn := defaultConnectionFromEnv(env)
	return map[string]DatabaseConnectionConfigInterface{conn.name: conn}
}

// defaultConnectionFromEnv reads the default database connection from environment variables.
func defaultConnectionFromEnv(env *envValidator) *databaseConnectionSettings {
	// Database Driver
	//
	// The database driver to use for the application.
	// Supported values: sqlite, turso, postgres, mysql
	driver := env.GetStringOrError(KEY_DEFAULT_DB_DRIVER, "select the database driver (e.g., sqlite, postgres)")
	name := defaultConnectionName

	// Database Host
	//
	// The hostname or IP address of the database server.
	// Not required when using sqlite.
	host := env.GetString(KEY_DEFAULT_DB_HOST)

	// Database Port
	//
	// The port the database server is listening on.
	// Common defaults: postgres=5432, mysql=3306
	// Not required when using sqlite.
	port := env.GetString(KEY_DEFAULT_DB_PORT)

	// Database Name
	//
	// The name of the database to connect to.
	// For sqlite, this is the file path (e.g., ./database.db or :memory:)
	dbName := env.GetStringOrError(KEY_DEFAULT_DB_DATABASE, "set the database name")

	// Database Username
	//
	// The username for authenticating with the database server.
	// Not required when using sqlite.
	user := env.GetString(KEY_DEFAULT_DB_USERNAME)

	// Database Password
	//
	// The password for authenticating with the database server.
	// Not required when using sqlite.
	pass := env.GetString(KEY_DEFAULT_DB_PASSWORD)

	// Direct DSN override
	//
	// Optional driver-specific connection string. When provided, it takes
	// precedence over the individual host/port/username/password fields.
	dsn := env.GetString(KEY_DEFAULT_DB_DSN)

	// Table prefix
	//
	// Optional prefix applied to table names by the ORM/query layer.
	prefix := env.GetString(KEY_DEFAULT_DB_PREFIX)

	// Connection Pool - Max Open Connections
	//
	// Maximum number of open connections to the database.
	// SQLite should stay at 1 to avoid concurrent write issues.
	// For postgres/mysql, 25 is a reasonable default for most apps.
	maxOpenConns := env.GetIntOrDefault(KEY_DEFAULT_DB_MAX_OPEN_CONNS, 25)
	if isSQLiteLike(driver) {
		maxOpenConns = 1
	}

	// Connection Pool - Max Idle Connections
	//
	// Maximum number of idle connections kept in the pool.
	// Should be less than or equal to MaxOpenConns.
	maxIdleConns := env.GetIntOrDefault(KEY_DEFAULT_DB_MAX_IDLE_CONNS, 5)
	if isSQLiteLike(driver) {
		maxIdleConns = 1
	}

	// Connection Pool - Max Connection Lifetime
	//
	// Maximum time a connection may be reused. Connections older than this
	// are closed and replaced. 0 means no limit.
	// Unit: seconds. Default: 300 (5 minutes)
	connMaxLifetime := time.Duration(env.GetIntOrDefault(KEY_DEFAULT_DB_CONN_MAX_LIFETIME_SECONDS, 300)) * time.Second
	if isSQLiteLike(driver) {
		connMaxLifetime = 30 * time.Second
	}

	// Connection Pool - Max Connection Idle Time
	//
	// Maximum time a connection may be idle before being closed.
	// 0 means no limit.
	// Unit: seconds. Default: 5
	connMaxIdleTime := time.Duration(env.GetIntOrDefault(KEY_DEFAULT_DB_CONN_MAX_IDLE_TIME_SECONDS, 5)) * time.Second

	// Database Charset
	//
	// Character set for the database connection. Only used for MySQL.
	// Example: utf8mb4, utf8
	charset := env.GetStringOrDefault(KEY_DEFAULT_DB_CHARSET, "utf8mb4")

	// Database Timezone
	//
	// Timezone for the database connection.
	// Example: UTC, America/New_York, Europe/London
	timezone := env.GetStringOrDefault(KEY_DEFAULT_DB_TIMEZONE, "UTC")

	// SSL mode default for non-SQLite drivers
	sslMode := env.GetStringOrDefault(KEY_DEFAULT_DB_SSL_MODE, "require")
	if isSQLiteLike(driver) {
		sslMode = ""
	}

	if !isSQLiteLike(driver) {
		env.RequireWhen(true, KEY_DEFAULT_DB_HOST, "required when `DB_DRIVER` is not sqlite", host)
		env.RequireWhen(true, KEY_DEFAULT_DB_PORT, "required when `DB_DRIVER` is not sqlite", port)
		env.RequireWhen(true, KEY_DEFAULT_DB_USERNAME, "required when `DB_DRIVER` is not sqlite", user)
		env.RequireWhen(true, KEY_DEFAULT_DB_PASSWORD, "required when `DB_DRIVER` is not sqlite", pass)
	}

	return &databaseConnectionSettings{
		name:                name,
		isDefault:           true,
		driver:              driver,
		host:                host,
		port:                port,
		database:            dbName,
		username:            user,
		password:            pass,
		sslMode:             sslMode,
		charset:             charset,
		timezone:            timezone,
		dsn:                 dsn,
		prefix:              prefix,
		maxOpenConns:        maxOpenConns,
		maxIdleConns:        maxIdleConns,
		connMaxLifetimeSecs: int(connMaxLifetime.Seconds()),
		connMaxIdleTimeSecs: int(connMaxIdleTime.Seconds()),
	}
}

// databaseConnectionSettings represents a single database connection.
type databaseConnectionSettings struct {
	name                string
	isDefault           bool
	driver              string
	host                string
	port                string
	database            string
	username            string
	password            string
	sslMode             string
	charset             string
	timezone            string
	dsn                 string
	prefix              string
	maxOpenConns        int
	maxIdleConns        int
	connMaxLifetimeSecs int
	connMaxIdleTimeSecs int
}

func (c *databaseConnectionSettings) GetName() string                { return c.name }
func (c *databaseConnectionSettings) GetDriver() string              { return c.driver }
func (c *databaseConnectionSettings) GetHost() string                { return c.host }
func (c *databaseConnectionSettings) GetPort() string                { return c.port }
func (c *databaseConnectionSettings) GetDatabase() string            { return c.database }
func (c *databaseConnectionSettings) GetUsername() string            { return c.username }
func (c *databaseConnectionSettings) GetPassword() string            { return c.password }
func (c *databaseConnectionSettings) GetSSLMode() string             { return c.sslMode }
func (c *databaseConnectionSettings) GetCharset() string             { return c.charset }
func (c *databaseConnectionSettings) GetTimezone() string            { return c.timezone }
func (c *databaseConnectionSettings) GetDSN() string                 { return c.dsn }
func (c *databaseConnectionSettings) GetPrefix() string              { return c.prefix }
func (c *databaseConnectionSettings) GetMaxOpenConns() int           { return c.maxOpenConns }
func (c *databaseConnectionSettings) GetMaxIdleConns() int           { return c.maxIdleConns }
func (c *databaseConnectionSettings) GetConnMaxLifetimeSeconds() int { return c.connMaxLifetimeSecs }
func (c *databaseConnectionSettings) GetConnMaxIdleTimeSeconds() int { return c.connMaxIdleTimeSecs }

// DatabaseNeatConfig converts a ConfigInterface into a neat ORM db.DBConfig.
// It maps all configured connections and the connection pool settings.
func DatabaseNeatConfig(cfg ConfigInterface) db.DBConfig {
	if cfg == nil {
		return db.DBConfig{}
	}

	defaultConnection := cfg.GetDatabaseDefaultConnection()
	if defaultConnection == "" {
		defaultConnection = defaultConnectionName
	}

	connections := make(map[string]db.ConnectionConfig)
	for _, conn := range cfg.GetDatabaseConnections() {
		if conn == nil {
			continue
		}
		connections[conn.GetName()] = connectionNeatConfig(conn)
	}

	pool := db.PoolConfig{
		QueryTimeout: 30,
	}
	if conn := cfg.GetDatabaseConnectionByName(defaultConnection); conn != nil {
		pool.MaxOpenConns = conn.GetMaxOpenConns()
		pool.MaxIdleConns = conn.GetMaxIdleConns()
		pool.ConnMaxLifetime = conn.GetConnMaxLifetimeSeconds()
		pool.ConnMaxIdleTime = conn.GetConnMaxIdleTimeSeconds()
	}

	return db.DBConfig{
		Default:     defaultConnection,
		Connections: connections,
		Pool:        pool,
	}
}

// connectionNeatConfig converts a single DatabaseConnectionConfigInterface
// into a neat ORM db.ConnectionConfig.
func connectionNeatConfig(conn DatabaseConnectionConfigInterface) db.ConnectionConfig {
	if conn == nil {
		return db.ConnectionConfig{}
	}

	driver := strings.ToLower(strings.TrimSpace(conn.GetDriver()))

	return db.ConnectionConfig{
		Driver:   driver,
		Dsn:      conn.GetDSN(),
		Host:     conn.GetHost(),
		Database: conn.GetDatabase(),
		Username: conn.GetUsername(),
		Password: conn.GetPassword(),
		Charset:  conn.GetCharset(),
		SSLMode:  conn.GetSSLMode(),
		Timezone: conn.GetTimezone(),
		Prefix:   conn.GetPrefix(),
		Port:     portToInt(conn.GetPort(), driver),
	}
}

// isSQLiteLike returns true when the driver is SQLite or Turso (libsql),
// which share the same file-based concurrency constraints.
func isSQLiteLike(driver string) bool {
	return driver == driverSQLite || driver == driverTurso
}

// portToInt converts a string port to an int, returning driver-specific
// defaults when the port is empty.
func portToInt(port, driver string) int {
	port = strings.TrimSpace(port)
	if port == "" {
		switch driver {
		case "mysql":
			return 3306
		case "postgres":
			return 5432
		case "sqlserver":
			return 1433
		case "oracle":
			return 1521
		default:
			return 0
		}
	}

	v, err := strconv.Atoi(port)
	if err != nil {
		return 0
	}
	return v
}
