package config

import (
	"strconv"
	"strings"
	"time"

	"github.com/dracory/neat/database/db"
)

func databaseConfig(env *envValidator) map[string]DatabaseConnectionConfigInterface {
	conn := defaultConnectionFromEnv(env)
	return map[string]DatabaseConnectionConfigInterface{conn.name: conn}
}

func defaultConnectionFromEnv(env *envValidator) *databaseConnectionSettings {
	driver := env.GetStringOrError(KEY_DEFAULT_DB_DRIVER, "select the database driver (e.g., sqlite, postgres)")
	name := defaultConnectionName
	host := env.GetString(KEY_DEFAULT_DB_HOST)
	port := env.GetString(KEY_DEFAULT_DB_PORT)
	dbName := env.GetStringOrError(KEY_DEFAULT_DB_DATABASE, "set the database name")
	user := env.GetString(KEY_DEFAULT_DB_USERNAME)
	pass := env.GetString(KEY_DEFAULT_DB_PASSWORD)
	dsn := env.GetString(KEY_DEFAULT_DB_DSN)
	prefix := env.GetString(KEY_DEFAULT_DB_PREFIX)

	maxOpenConns := env.GetIntOrDefault(KEY_DEFAULT_DB_MAX_OPEN_CONNS, 25)
	if isSQLiteLike(driver) {
		maxOpenConns = 1
	}

	maxIdleConns := env.GetIntOrDefault(KEY_DEFAULT_DB_MAX_IDLE_CONNS, 5)
	if isSQLiteLike(driver) {
		maxIdleConns = 1
	}

	connMaxLifetime := time.Duration(env.GetIntOrDefault(KEY_DEFAULT_DB_CONN_MAX_LIFETIME_SECONDS, 300)) * time.Second
	if isSQLiteLike(driver) {
		connMaxLifetime = 30 * time.Second
	}

	connMaxIdleTime := time.Duration(env.GetIntOrDefault(KEY_DEFAULT_DB_CONN_MAX_IDLE_TIME_SECONDS, 5)) * time.Second

	charset := env.GetStringOrDefault(KEY_DEFAULT_DB_CHARSET, "utf8mb4")
	timezone := env.GetStringOrDefault(KEY_DEFAULT_DB_TIMEZONE, "UTC")

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

func isSQLiteLike(driver string) bool {
	return driver == driverSQLite || driver == driverTurso
}

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
