# Installation

## Requirements

The Pico framework has a few system requirements:

- **Go** >= 1.26
- **SQLite** (for local development) or **MySQL** / **PostgreSQL** (for production)

## Installing Pico

Clone the repository and install dependencies:

```bash
git clone https://github.com/dracory/pico
cd pico
go mod download
```

## Configuration

Copy the example environment file and configure it:

```bash
cp .env.example .env
```

All configuration options are stored in the `.env` file. See the [Configuration](configuration.md) documentation for details.

## Serving Your Application

### Development Server

To start the server with hot reload (requires [Air](https://github.com/air-verse/air)):

```bash
task dev
```

Or without hot reload:

```bash
go run ./cmd/server
```

The server starts at `http://localhost:8080` by default.

### Building

Build a binary:

```bash
task build
```

The binary is output to `bin/pico.exe`.

## Project Structure

```
pico/
├── cmd/
│   └── server/
│       └── main.go               # Application entry point
├── internal/
│   ├── app/
│   │   ├── app_interface.go      # AppInterface (DI container)
│   │   ├── app_implementation.go # App implementation
│   │   └── database_open.go      # Database connection
│   ├── config/
│   │   ├── config_interfaces.go     # ConfigInterface
│   │   ├── config_implementation.go # Config implementation
│   │   ├── app_config.go            # App config loader
│   │   ├── database_config.go       # Database config loader
│   │   ├── constants.go             # Environment variable keys
│   │   └── version.go               # Version
│   └── routes/
│       └── router.go              # Router + middleware
├── .env.example
├── go.mod
└── taskfile.yml
```

## Compatibility

Pico is a stripped-down version of [Blueprint](https://github.com/dracory/blueprint). It shares the same core libraries (`dracory/rtr`, `dracory/neat`, `dracory/env`) but removes all store, CMS, auth, and task infrastructure.

If your project outgrows Pico, you can migrate to Blueprint for the full feature set.
