# Upgrade Guide: v0.1.0 to v0.2.0

This guide describes the steps required to upgrade your Pico application from version `v0.1.0` to `v0.2.0`.

## Overview

Version `v0.2.0` includes dependency updates, internal maintenance updates, and removal of the external `tint` logging library in favor of Go's standard `log/slog` handlers.

## Dependency Changes

### Updated Dependencies

- `modernc.org/sqlite`: upgraded to `v1.58.0`
- `github.com/dracory/cdn`: upgraded to `v1.12.0`
- `modernc.org/libc`: upgraded to `v1.75.6`

### Removed Dependencies

- `github.com/lmittmann/tint`: Removed in favor of standard library `log/slog.NewTextHandler`.

## Upgrade Steps

1. **Update Go Dependencies**

   Run the following command in your project root to fetch updated dependencies and clean up unused modules:

   ```bash
   go get -u ./...
   go mod tidy
   ```

2. **Verify Configuration & Tests**

   Run the test suite to confirm everything compiles and works as expected:

   ```bash
   go test ./...
   ```
