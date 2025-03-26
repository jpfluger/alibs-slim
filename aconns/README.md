# aconns

A set of standardized interfaces and structures for managing communications in [`golang`](https://go.dev/) projects. `aconns` provides a unified way to handle inputs and outputs across various backends — including databases, web services, file systems, and other I/O systems.

It builds on the foundation of [`alibs-slim`](https://github.com/jpfluger/alibs-slim), and is designed to be composable, minimal, and dependency-aware. The core module provides abstractions, while specific implementations live in submodules.

## Project Structure

This project follows a multi-module design to keep dependencies clean and modular:

- **Core**
  - `github.com/jpfluger/alibs-slim/aconns` — Shared types and interfaces for I/O connection abstraction.
- **Drivers**
  - `github.com/jpfluger/alibs-slim/aconns/adb-pg` — PostgreSQL connection implementation.
  - `github.com/jpfluger/alibs-slim/aconns/adb-mysql` — MySQL connection implementation.
  - Etc...
- **Global Access**
  - `github.com/jpfluger/alibs-slim/aconns/g-aconns` — Loads and registers all available driver connectors.

## Features

- 🔌 Common interfaces for connecting to diverse systems
- 📦 Self-contained driver modules with their own `go.mod` files
- 🔍 Minimal core dependencies
- 🧩 Works seamlessly with `alibs-slim` and the [Echo Framework](https://echo.labstack.com/)
- ⚙️ Designed for testability, extensibility, and clarity

## Using `aconns`

This project uses [Go Modules](https://go.dev/ref/mod). If you're using Go 1.16 or newer, module support is on by default. If needed, enable it explicitly:

```bash
export GO111MODULE=on
```

### Importing Modules

Import only the components you need:

```go
// Core interfaces
import "github.com/jpfluger/alibs-slim/aconns"

// PostgreSQL support
import "github.com/jpfluger/alibs-slim/aconns/adb-pg"

// Register all drivers at once
// import "github.com/jpfluger/alibs-slim/aconns/g-aconns"
```
