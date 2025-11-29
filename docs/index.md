# Hybrid_Lib_Go Documentation Index

**Version:** 1.0.1  
**Date:** November 28, 2025  
**SPDX-License-Identifier:** BSD-3-Clause
**License File:** See the LICENSE file in the project root.
**Copyright:** (c) 2025 Michael Gardner, A Bit of Help, Inc.  
**Status:** Released  

---

## Welcome

Welcome to the **Hybrid_Lib_Go** documentation. This Go 1.23+ library demonstrates professional hexagonal architecture with functional programming principles, static dependency injection via generics, and railway-oriented error handling.

This is a **library** template (not an application). It provides reusable business logic that consuming applications can import.

---

## Quick Navigation

### Getting Started

- **[Quick Start Guide](quick_start.md)** - Get up and running in minutes
  - Installation instructions
  - Library usage
  - Understanding the architecture
  - Running tests

### Formal Documentation

- **[Software Requirements Specification (SRS)](formal/software_requirements_specification.md)** - Complete requirements
  - Functional requirements
  - Non-functional requirements
  - System constraints

- **[Software Design Specification (SDS)](formal/software_design_specification.md)** - Architecture and design
  - 4-layer library hexagonal architecture
  - Static dependency injection via generics
  - Railway-oriented programming patterns
  - Component relationships

- **[Software Test Guide](formal/software_test_guide.md)** - Testing documentation
  - Test organization (unit/integration)
  - Running tests
  - Coverage procedures

---

## Architecture Overview

Hybrid_Lib_Go implements a **4-layer library hexagonal architecture** (also known as Ports and Adapters or Clean Architecture):

### Layer Structure

```
┌─────────────────────────────────────────────┐
│  API Layer (Public Facade)                  │  api/, api/adapter/desktop/
├─────────────────────────────────────────────┤
│  Infrastructure                             │  Driven Adapters (Console Writer)
├─────────────────────────────────────────────┤
│  Application                                │  Use Cases + Ports
├─────────────────────────────────────────────┤
│  Domain                                     │  Business Logic (ZERO dependencies)
└─────────────────────────────────────────────┘
```

### Library vs Application Architecture

| Library (4 layers)       | Application (5 layers)      |
|--------------------------|----------------------------|
| api/                     | bootstrap/                 |
| api/adapter/desktop/             | presentation/              |
| infrastructure/          | infrastructure/            |
| application/             | application/               |
| domain/                  | domain/                    |

### Key Principles

1. **Domain Isolation**: Domain layer has zero external dependencies
2. **API Facade**: api/ re-exports types, does NOT import infrastructure
3. **Platform Wiring**: api/adapter/desktop/ wires infrastructure to application
4. **Static Dispatch**: Dependency injection via generics (compile-time, zero overhead)
5. **Railway-Oriented**: Result monads for error handling (no panics across boundaries)
6. **Multi-Module Workspace**: go.work manages separate go.mod per layer

---

## Visual Documentation

### UML Diagrams

Located in `diagrams/` directory:

- **layer_dependencies.svg** - Shows 4-layer library dependency flow
- **api_reexport_pattern.svg** - API facade re-export pattern
- **package_structure.svg** - Actual package hierarchy
- **error_handling_flow.svg** - Error propagation through layers
- **static_dispatch_api.svg** - Generic vs interface comparison

All diagrams are generated from PlantUML sources (.puml files).

---

## Library Usage

### Basic Usage

```go
import (
    "context"
    "github.com/abitofhelp/hybrid_lib_go/api"
    "github.com/abitofhelp/hybrid_lib_go/api/adapter/desktop"
)

func main() {
    // Create a greeter with console output
    greeter := desktop.NewGreeter()

    // Execute greeting
    ctx := context.Background()
    result := greeter.Execute(ctx, api.NewGreetCommand("Alice"))

    if result.IsOk() {
        // Success - greeting was written to console
    } else {
        // Handle error
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

### Custom Writer

```go
// Use a custom writer for testing or other output destinations
writer := &MockWriter{Buffer: &bytes.Buffer{}}
greeter := desktop.GreeterWithWriter[*MockWriter](writer)
result := greeter.Execute(ctx, api.NewGreetCommand("Bob"))
output := writer.String() // "Hello, Bob!\n"
```

---

## Key Features

### Static Dependency Injection

Uses **generics** instead of interfaces for dependency injection:

```go
// Port definition (interface constraint)
type WriterPort interface {
    Write(ctx context.Context, message string) domerr.Result[model.Unit]
}

// Generic use case with static dispatch
type GreetUseCase[W WriterPort] struct {
    writer W
}

// Wiring in api/adapter/desktop (compile-time resolution)
writer := adapter.NewConsoleWriter()
uc := usecase.NewGreetUseCase[*adapter.ConsoleWriter](writer)
```

**Benefits**:
- Zero runtime overhead (no vtable lookups)
- Full inlining potential
- Compile-time type safety
- All method calls devirtualized

### Railway-Oriented Programming

All errors propagate via **Result monad** (no panics across boundaries):

```go
// Use case returns Result[Unit]
func (uc *GreetUseCase[W]) Execute(ctx context.Context, cmd command.GreetCommand) domerr.Result[model.Unit] {
    personResult := valueobject.CreatePerson(cmd.Name())

    if personResult.IsError() {
        return domerr.Err[model.Unit](personResult.ErrorInfo())
    }

    person := personResult.Value()
    return uc.writer.Write(ctx, person.GreetingMessage())
}
```

### API Facade Pattern

**Pattern**: api/ re-exports types from domain and application without importing infrastructure:

```go
// api/api.go - Public facade
import (
    "github.com/abitofhelp/hybrid_lib_go/application/command"
    domerr "github.com/abitofhelp/hybrid_lib_go/domain/error"
    "github.com/abitofhelp/hybrid_lib_go/domain/valueobject"
)

// Re-exported types
type Result[T any] = domerr.Result[T]
type Person = valueobject.Person
type GreetCommand = command.GreetCommand

// Factory functions
func NewGreetCommand(name string) GreetCommand {
    return command.NewGreetCommand(name)
}
```

Infrastructure is wired in platform-specific sub-packages (api/adapter/desktop/).

---

## Directory Structure

```
hybrid_lib_go/
├── api/                       # Public facade
│   ├── go.mod                 # Depends on application + domain (NOT infrastructure)
│   ├── api.go                 # Re-exports types
│   └── adapter/
│       └── desktop/           # Platform-specific wiring (composition root)
│           ├── go.mod         # Depends on all layers
│           └── desktop.go     # Creates ready-to-use greeter
├── domain/                    # Pure business logic
│   ├── go.mod                 # ZERO external dependencies
│   ├── error/                 # Result monad, error types
│   └── valueobject/           # Immutable value objects
├── application/               # Use cases + ports
│   ├── go.mod                 # Depends only on domain
│   ├── command/               # Input DTOs
│   ├── error/                 # Error helpers
│   ├── model/                 # Output DTOs (Unit)
│   ├── port/
│   │   ├── inbound/           # Use case interfaces
│   │   └── outbound/          # Infrastructure interfaces
│   └── usecase/               # Use case orchestration
├── infrastructure/            # Adapters (driven)
│   ├── go.mod                 # Depends on application + domain
│   └── adapter/               # Console writer
├── test/
│   └── integration/           # API integration tests
├── docs/
│   ├── formal/                # SRS, SDS, Test Guide
│   ├── diagrams/              # UML diagrams
│   ├── quick_start.md         # Getting started
│   └── index.md               # This file
├── scripts/
│   └── arch_guard/            # Architecture validation
├── go.work                    # Workspace definition
├── Makefile                   # Build automation
└── README.md                  # Project overview
```

---

## Build System

### Make Targets

**Building**:
```bash
make build              # Build all modules
```

**Testing**:
```bash
make test               # Run unit tests
make test-all           # Run all tests (unit + integration)
make test-integration   # Integration tests only
```

**Quality**:
```bash
make check-arch         # Architecture validation
make fmt                # Format code
make lint               # Run linter
```

---

## Dependencies

### Runtime Dependencies

- **None**: Domain layer has zero external dependencies

### Development Dependencies

- **testify** (v1.11.1): Testing assertions (test module only)

### Build Requirements

- **Go**: 1.23+ (workspace and generics support)
- **Make**: GNU Make for build automation
- **Python 3**: For architecture validation (arch_guard.py)
- **Java 11+**: For PlantUML diagram generation (optional)

---

## License

Hybrid_Lib_Go is licensed under the **BSD-3-Clause License**.

Copyright (c) 2025 Michael Gardner, A Bit of Help, Inc.

See [LICENSE](../LICENSE) for full license text.

---

## Project Links

- **GitHub**: https://github.com/abitofhelp/hybrid_lib_go
- **Author**: Michael Gardner (https://github.com/abitofhelp)
- **Company**: A Bit of Help, Inc.

---

**Last Updated**: November 26, 2025
