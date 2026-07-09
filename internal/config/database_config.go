package config

import (
	"strconv"
	"strings"
	"time"

	"github.com/dracory/neat/database/db"
)

func databaseConfig(env *envValidator) databaseSettings {
	driver := env.GetStringOrError(KEY_DB_DRIVER, "select the database driver (e.g., sqlite, postgres)")
	defaultConnection := env.GetStringOrDefault(KEY_DB_DEFAULT_CONNECTION, "default")
	host := env.GetString(KEY_DB_HOST)
	port := env.GetString(KEY_DB_PORT)
	name := env.GetStringOrError(KEY_DB_DATABASE, "set the database name")
	user := env.GetString(KEY_DB_USERNAME)
	pass := env.GetString(KEY_DB_PASSWORD)
	dsn := env.GetString(KEY_DB_DSN)
	prefix := env.GetString(KEY_DB_PREFIX)

	maxOpenConns := env.GetIntOrDefault(KEY_DB_MAX_OPEN_CONNS, 25)
	if driver == driverSQLite {
		maxOpenConns = 1
	}

	maxIdleConns := env.GetIntOrDefault(KEY_DB_MAX_IDLE_CONNS, 5)
	if driver == driverSQLite {
		maxIdleConns = 1
	}

	connMaxLifetime := time.Duration(env.GetIntOrDefault(KEY_DB_CONN_MAX_LIFETIME_SECONDS, 300)) * time.Second
	if driver == driverSQLite {
		connMaxLifetime = 30 * time.Second
	}

	connMaxIdleTime := time.Duration(env.GetIntOrDefault(KEY_DB_CONN_MAX_IDLE_TIME_SECONDS, 5)) * time.Second

	charset := env.GetStringOrDefault(KEY_DB_CHARSET, "utf8mb4")
	timezone := env.GetStringOrDefault(KEY_DB_TIMEZONE, "UTC")

	sslMode := env.GetStringOrDefault(KEY_DB_SSL_MODE, "require")
	if driver == driverSQLite {
		sslMode = ""
	}

	if driver != driverSQLite {
		env.RequireWhen(true, KEY_DB_HOST, "required when `DB_DRIVER` is not sqlite", host)
		env.RequireWhen(true, KEY_DB_PORT, "required when `DB_DRIVER` is not sqlite", port)
		env.RequireWhen(true, KEY_DB_USERNAME, "required when `DB_DRIVER` is not sqlite", user)
		env.RequireWhen(true, KEY_DB_PASSWORD, "required when `DB_DRIVER` is not sqlite", pass)
	}

	defaultConn := databaseConnectionSettings{
		name:     defaultConnection,
		driver:   driver,
		host:     host,
		port:     port,
		database: name,
		username: user,
		password: pass,
		sslMode:  sslMode,
		charset:  charset,
		timezone: timezone,
		dsn:      dsn,
		prefix:   prefix,
	}

	connections := map[string]DatabaseConnectionConfigInterface{
		defaultConnection: &defaultConn,
	}

	return databaseSettings{
		defaultConnection: defaultConnection,
		connections:       connections,
		driver:            driver,
		host:              host,
		port:              port,
		name:              name,
		user:              user,
		pass:              pass,
		sslMode:           sslMode,
		maxOpenConns:      maxOpenConns,
		maxIdleConns:      maxIdleConns,
		connMaxLifetime:   connMaxLifetime,
		connMaxIdleTime:   connMaxIdleTime,
		charset:           charset,
		timezone:          timezone,
		dsn:               dsn,
		prefix:            prefix,
	}
}

type databaseConnectionSettings struct {
	name     string
	driver   string
	host     string
	port     string
	database string
	username string
	password string
	sslMode  string
	charset  string
	timezone string
	dsn      string
	prefix   string
}

type databaseSettings struct {
	defaultConnection string
	connections       map[string]DatabaseConnectionConfigInterface
	driver            string
	host              string
	port              string
	name              string
	user              string
	pass              string
	sslMode           string
	maxOpenConns      int
	maxIdleConns      int
	connMaxLifetime   time.Duration
	connMaxIdleTime   time.Duration
	charset           string
	timezone          string
	dsn               string
	prefix            string
}

func (c *databaseConnectionSettings) GetName() string     { return c.name }
func (c *databaseConnectionSettings) GetDriver() string   { return c.driver }
func (c *databaseConnectionSettings) GetHost() string     { return c.host }
func (c *databaseConnectionSettings) GetPort() string     { return c.port }
func (c *databaseConnectionSettings) GetDatabase() string { return c.database }
func (c *databaseConnectionSettings) GetUsername() string { return c.username }
func (c *databaseConnectionSettings) GetPassword() string { return c.password }
func (c *databaseConnectionSettings) GetSSLMode() string  { return c.sslMode }
func (c *databaseConnectionSettings) GetCharset() string  { return c.charset }
func (c *databaseConnectionSettings) GetTimezone() string { return c.timezone }
func (c *databaseConnectionSettings) GetDSN() string      { return c.dsn }
func (c *databaseConnectionSettings) GetPrefix() string   { return c.prefix }

func DatabaseNeatConfig(cfg ConfigInterface) db.DBConfig {
	if cfg == nil {
		return db.DBConfig{}
	}

	defaultConnection := cfg.GetDatabaseDefaultConnection()
	if defaultConnection == "" {
		defaultConnection = "default"
	}

	connections := make(map[string]db.ConnectionConfig)
	for _, conn := range cfg.GetDatabaseConnections() {
		if conn == nil {
			continue
		}
		connections[conn.GetName()] = connectionNeatConfig(conn)
	}

	if _, ok := connections[defaultConnection]; !ok {
		connections[defaultConnection] = connectionNeatConfig(&databaseConnectionSettings{
			name:     defaultConnection,
			driver:   cfg.GetDatabaseDriver(),
			host:     cfg.GetDatabaseHost(),
			port:     cfg.GetDatabasePort(),
			database: cfg.GetDatabaseName(),
			username: cfg.GetDatabaseUsername(),
			password: cfg.GetDatabasePassword(),
			sslMode:  cfg.GetDatabaseSSLMode(),
			charset:  cfg.GetDatabaseCharset(),
			timezone: cfg.GetDatabaseTimezone(),
			dsn:      cfg.GetDatabaseDSN(),
			prefix:   cfg.GetDatabasePrefix(),
		})
	}

	pool := db.PoolConfig{
		MaxOpenConns:    cfg.GetDatabaseMaxOpenConns(),
		MaxIdleConns:    cfg.GetDatabaseMaxIdleConns(),
		ConnMaxLifetime: cfg.GetDatabaseConnMaxLifetimeSeconds(),
		ConnMaxIdleTime: cfg.GetDatabaseConnMaxIdleTimeSeconds(),
		QueryTimeout:    30,
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
