package config

// ============================================================================
// Environment Constants
// ============================================================================

const APP_ENVIRONMENT_DEVELOPMENT = "development"
const APP_ENVIRONMENT_LOCAL = "local"
const APP_ENVIRONMENT_PRODUCTION = "production"
const APP_ENVIRONMENT_STAGING = "staging"
const APP_ENVIRONMENT_TESTING = "testing"

// ============================================================================
// Database Driver Constants
// ============================================================================

const driverSQLite = "sqlite"
const driverTurso = "turso"
const defaultConnectionName = "default"

// ============================================================================
// App Environment Variable Keys
// ============================================================================

const KEY_APP_DEBUG = "APP_DEBUG"
const KEY_APP_ENVIRONMENT = "APP_ENV"
const KEY_APP_NAME = "APP_NAME"
const KEY_APP_URL = "APP_URL"
const KEY_APP_HOST = "APP_HOST"
const KEY_APP_PORT = "APP_PORT"

// ============================================================================
// Database Environment Variable Keys
// ============================================================================

// Keys for the default database connection
const KEY_DEFAULT_DB_DRIVER = "DB_DRIVER"
const KEY_DEFAULT_DB_HOST = "DB_HOST"
const KEY_DEFAULT_DB_PORT = "DB_PORT"
const KEY_DEFAULT_DB_DATABASE = "DB_DATABASE"
const KEY_DEFAULT_DB_USERNAME = "DB_USERNAME"
const KEY_DEFAULT_DB_PASSWORD = "DB_PASSWORD"
const KEY_DEFAULT_DB_SSL_MODE = "DB_SSL_MODE"
const KEY_DEFAULT_DB_CHARSET = "DB_CHARSET"
const KEY_DEFAULT_DB_TIMEZONE = "DB_TIMEZONE"
const KEY_DEFAULT_DB_DSN = "DB_DSN"
const KEY_DEFAULT_DB_PREFIX = "DB_PREFIX"
const KEY_DEFAULT_DB_MAX_OPEN_CONNS = "DB_MAX_OPEN_CONNS"
const KEY_DEFAULT_DB_MAX_IDLE_CONNS = "DB_MAX_IDLE_CONNS"
const KEY_DEFAULT_DB_CONN_MAX_LIFETIME_SECONDS = "DB_CONN_MAX_LIFETIME_SECONDS"
const KEY_DEFAULT_DB_CONN_MAX_IDLE_TIME_SECONDS = "DB_CONN_MAX_IDLE_TIME_SECONDS"
