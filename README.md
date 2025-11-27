# Hybrid_Lib_Go - Go Library with Hexagonal Architecture

**Version:** 1.0.0
**Date:** November 26, 2025
**Copyright:** (c) 2025 Michael Gardner, A Bit of Help, Inc.
**License:** BSD-3-Clause

## Overview

A **professional Go library** demonstrating **hybrid DDD/Clean/Hexagonal architecture** with **strict module boundaries** enforced via Go workspaces and **functional programming** principles using custom **domain-level Result monads** (ZERO external dependencies in domain layer).

> **Library Template:** This project serves as a **template for enterprise Go library development**. Use the included `scripts/brand_project/brand_project.py` script to generate a new project from this template with your own project name, module paths, and branding.

This is a **reusable library** showcasing:
- **4-Layer Hexagonal Architecture** (Domain, Application, Infrastructure, API)
- **Strict Module Boundaries** via go.work and separate go.mod per layer
- **Static Dispatch via Generics** (zero-overhead dependency injection)
- **Railway-Oriented Programming** with Result monads (no panics across boundaries)
- **API Facade Pattern** for clean public interface
- **Multi-Module Workspace** (compiler-enforced boundaries)

## Architecture

### Module Structure

**Strict boundaries enforced by Go modules:**

```
hybrid_lib_go/
├── go.work                          # Workspace definition (manages all modules)
├── domain/                          # Module: Pure business logic (ZERO external dependencies)
│   └── go.mod                       # ZERO external dependencies - custom Result types
├── application/                     # Module: Use cases and ports
│   └── go.mod                       # Depends ONLY on domain
├── infrastructure/                  # Module: Driven adapters
│   └── go.mod                       # Depends on application + domain
├── api/                             # Module: Public facade (re-exports types)
│   ├── go.mod                       # Depends on application + domain (NOT infrastructure)
│   └── adapter/
│       └── desktop/                 # Sub-module: Composition root
│           └── go.mod               # Depends on ALL modules (wires infrastructure)
└── test/
    └── integration/                 # Integration tests for API usage
```

### Key Architectural Rules

**4-Layer Library Architecture:**

| Layer | Dependencies | Purpose |
|-------|-------------|---------|
| `domain/` | NONE | Pure business logic, value objects, Result monad |
| `application/` | domain | Use cases, ports, commands |
| `infrastructure/` | application, domain | Adapters (ConsoleWriter, etc.) |
| `api/` | application, domain | Public facade, re-exports types |
| `api/adapter/desktop/` | ALL | Composition root, wires infrastructure |

**Critical Boundary Rules:**
- **api/** re-exports types but does NOT import infrastructure
- **api/adapter/desktop/** (composition root) CAN import infrastructure
- **Domain** has ZERO external dependencies
- All dependencies flow INWARD toward Domain

### Library Usage

**Basic Usage (with convenience component):**

```go
import (
    "context"
    "github.com/abitofhelp/hybrid_lib_go/api"
    "github.com/abitofhelp/hybrid_lib_go/api/adapter/desktop"
)

func main() {
    // Create ready-to-use greeter
    greeter := desktop.NewGreeter()

    // Execute greeting
    ctx := context.Background()
    result := greeter.Execute(ctx, api.NewGreetCommand("Alice"))

    if result.IsOk() {
        // Greeting was written to console
    } else {
        errInfo := result.ErrorInfo()
        switch errInfo.Kind {
        case api.ValidationError:
            // Handle validation error
        case api.InfrastructureError:
            // Handle infrastructure error
        }
    }
}
```

**Custom Writer (for testing or different output):**

```go
import (
    "context"
    "github.com/abitofhelp/hybrid_lib_go/api"
    "github.com/abitofhelp/hybrid_lib_go/api/adapter/desktop"
)

// Implement WriterPort interface
type MyWriter struct { /* ... */ }

func (w *MyWriter) Write(ctx context.Context, msg string) api.Result[api.Unit] {
    // Custom write logic
    return api.Ok(api.Unit{})
}

func main() {
    writer := &MyWriter{}
    greeter := desktop.GreeterWithWriter(writer)
    result := greeter.Execute(ctx, api.NewGreetCommand("Bob"))
}
```

### Dependency Injection Pattern

**Static Dispatch via Generics:**

```go
// Port interface defines the contract
type WriterPort interface {
    Write(ctx context.Context, message string) domerr.Result[model.Unit]
}

// Generic use case with interface constraint
type GreetUseCase[W outbound.WriterPort] struct {
    writer W
}

// api/adapter/desktop wires with concrete types
func NewGreeter() *Greeter {
    writer := adapter.NewConsoleWriter()
    uc := usecase.NewGreetUseCase[*adapter.ConsoleWriter](writer)
    return &Greeter{useCase: uc}
}
```

**Benefits:**
- Zero runtime overhead (no vtable lookups)
- Type-safe (verified at compile time)
- Static dispatch (compiler knows exact types)
- Inlining potential (optimizer can inline method calls)

## Error Handling: Railway-Oriented Programming

**NO PANICS across layer boundaries.** All errors propagate via domain Result monad:

```go
// Domain defines custom Result[T] monad (ZERO external dependencies)
func Execute(cmd GreetCommand) domerr.Result[model.Unit] {
    personResult := valueobject.CreatePerson(cmd.Name())

    if personResult.IsError() {
        return domerr.Err[model.Unit](personResult.ErrorInfo())
    }

    person := personResult.Value()
    return writer.Write(ctx, person.GreetingMessage())
}
```

**Error Types:**
- `ValidationError` - Invalid input (empty name, name too long)
- `InfrastructureError` - I/O failures, system errors

## Building

### Prerequisites

- **Go 1.23+** (for workspace and generics support)
- **Make** (for build automation)
- **Python 3** (optional, for architecture validation)

### Build Commands

```bash
# Build all library modules
make build

# Run all tests (unit + integration)
make test

# Run unit tests only
make test-unit

# Run integration tests
make test-integration

# Validate architecture boundaries
make check-arch

# Format code
make fmt

# Run linter
make lint

# Clean artifacts
make clean
```

## Testing

```bash
# Run all tests
make test-all

# Run with coverage
make test-coverage

# Run specific test
go test -v -run TestGreeter_Execute_Success ./test/integration/...
```

**Test Structure:**
- **Unit tests**: Co-located with code (`*_test.go`)
- **Integration tests**: `test/integration/` with `//go:build integration` tag

## API Reference

### Types (via `api` package)

| Type | Description |
|------|-------------|
| `Result[T]` | Result monad (Ok or Error) |
| `ErrorType` | Error information struct |
| `ErrorKind` | Error category (Validation, Infrastructure) |
| `Person` | Domain value object |
| `GreetCommand` | Input command |
| `WriterPort` | Output port interface |
| `Unit` | Void return type |

### Constants

| Constant | Description |
|----------|-------------|
| `api.ValidationError` | Error kind for validation failures |
| `api.InfrastructureError` | Error kind for I/O failures |
| `api.MaxNameLength` | Maximum allowed name length (100) |

### Functions

| Function | Description |
|----------|-------------|
| `api.NewGreetCommand(name)` | Create a greet command |
| `api.CreatePerson(name)` | Create a Person value object |
| `api.Ok[T](value)` | Create successful Result |
| `api.Err[T](error)` | Create error Result |
| `desktop.NewGreeter()` | Create ready-to-use greeter |
| `desktop.GreeterWithWriter(w)` | Create greeter with custom writer |

## Module Boundaries

**Enforced by go.mod dependencies:**

- **domain**: ZERO external dependencies (custom Result types)
- **application**: domain ONLY
- **infrastructure**: application + domain
- **api**: application + domain (NOT infrastructure)
- **api/adapter/desktop**: ALL modules (composition root)

**Compiler enforces these rules** - attempting to import forbidden packages results in build errors.

## Documentation

- **[Quick Start Guide](docs/quick_start.md)** - Get up and running
- **[Documentation Index](docs/index.md)** - All documentation links
- **[Software Design Specification](docs/formal/software_design_specification.md)** - Architecture details
- **[Architecture Diagrams](docs/diagrams/)** - Visual documentation

## Creating a New Project

Use the `brand_project.py` script to create a new library from this template:

```bash
cd scripts
python3 -m brand_project \
    --old-project hybrid_lib_go \
    --new-project my_library \
    --old-org abitofhelp \
    --new-org mycompany \
    --source /path/to/hybrid_lib_go \
    --target /path/to/my_library
```

**What gets updated:**
- Project name throughout all files
- GitHub organization/username in module paths
- Copyright holder information
- All `go.mod` module paths
- Import statements in Go source files

## Standards Compliance

This project follows:
- **Go Language Standards** (`~/.claude/agents/go.md`)
- **Architecture Standards** (`~/.claude/agents/architecture.md`)
- **Functional Programming Standards** (`~/.claude/agents/functional.md`)

### Key Standards Applied

1. **SPDX Headers:** All `.go` files have SPDX license headers
2. **Result Monads:** All fallible operations return `Result[T]`
3. **No Panics:** Errors are values, recovery patterns at boundaries
4. **Module Boundaries:** Compiler-enforced via go.mod
5. **Static Dispatch:** Generic types with interface constraints
6. **API Facade:** Clean public interface via api/ package

## License

BSD-3-Clause - See LICENSE file in project root.

## Author

Michael Gardner
A Bit of Help, Inc.
https://github.com/abitofhelp
