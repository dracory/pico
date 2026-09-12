# Upgrade Guide: v0.1.0 to v0.2.0

This guide describes the steps required to upgrade your Pico application from version `v0.1.0` to `v0.2.0`.

## Overview

Version `v0.2.0` includes dependency updates and internal maintenance updates.

## Dependency Updates

The following dependencies have been updated:

- `modernc.org/sqlite`: upgraded to `v1.58.0`
- `github.com/dracory/cdn`: upgraded to `v1.12.0`
- `modernc.org/libc`: upgraded to `v1.75.6`

## Upgrade Steps

1. **Update Go Dependencies**

   Run the following command in your project root to fetch the latest updated dependencies:

   ```bash
   go get -u ./...
   go mod tidy
   ```

2. **Verify Configuration & Tests**

   Run the test suite to confirm everything compiles and works as expected:

   ```bash
   go test ./...
   ```
