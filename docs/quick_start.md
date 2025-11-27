# Hybrid_Lib_Go Quick Start Guide

**Version:** 1.0.0
**Date:** November 26, 2025
**SPDX-License-Identifier:** BSD-3-Clause
**License File:** See the LICENSE file in the project root.
**Copyright:** (c) 2025 Michael Gardner, A Bit of Help, Inc.
**Status:** Released

---

## Table of Contents

- [Installation](#installation)
- [Library Usage](#library-usage)
- [Understanding the Architecture](#understanding-the-architecture)
- [Making Your First Change](#making-your-first-change)
- [Running Tests](#running-tests)
- [Build Targets](#build-targets)
- [Common Issues](#common-issues)
- [Next Steps](#next-steps)

---

## Installation

### Prerequisites

- **Go**: Version 1.23+ (workspace and generics support required)
- **Make**: GNU Make for build automation
- **Python 3**: For architecture validation (optional)
- **Java 11+**: For PlantUML diagram generation (optional)

### Clone and Build

```bash
# Clone the repository
git clone https://github.com/abitofhelp/hybrid_lib_go.git
cd hybrid_lib_go

# Build all modules
make build

# Or build directly with Go
go build ./...
```

### Verify Installation

```bash
# Run unit tests
make test

# Run integration tests
make test-integration
```

**Success!** You've built the library and verified it works.

---

## Library Usage

Hybrid_Lib_Go is a **library**, not an application. You import it into your own application.

### Basic Usage

```go
package main

import (
    "context"
    "fmt"
    "os"

    "github.com/abitofhelp/hybrid_lib_go/api"
    "github.com/abitofhelp/hybrid_lib_go/api/adapter/desktop"
)

func main() {
    // Create a greeter with console output
    greeter := desktop.NewGreeter()

    // Create context and command
    ctx := context.Background()
    cmd := api.NewGreetCommand("Alice")

    // Execute greeting
    result := greeter.Execute(ctx, cmd)

    // Handle result
    if result.IsOk() {
        fmt.Println("Greeting sent successfully!")
        os.Exit(0)
    } else {
        errInfo := result.ErrorInfo()
        switch errInfo.Kind {
        case api.ValidationError:
            fmt.Fprintf(os.Stderr, "Invalid input: %s\n", errInfo.Message)
        case api.InfrastructureError:
            fmt.Fprintf(os.Stderr, "System error: %s\n", errInfo.Message)
        }
        os.Exit(1)
    }
}
```

### Custom Writer for Testing

```go
package myapp_test

import (
    "bytes"
    "context"
    "testing"

    "github.com/abitofhelp/hybrid_lib_go/api"
    "github.com/abitofhelp/hybrid_lib_go/api/adapter/desktop"
    "github.com/abitofhelp/hybrid_lib_go/application/model"
    domerr "github.com/abitofhelp/hybrid_lib_go/domain/error"
)

// MockWriter captures output for testing
type MockWriter struct {
    Buffer bytes.Buffer
}

func (w *MockWriter) Write(ctx context.Context, msg string) domerr.Result[model.Unit] {
    w.Buffer.WriteString(msg)
    return domerr.Ok(model.Unit{})
}

func (w *MockWriter) String() string {
    return w.Buffer.String()
}

func TestGreeting(t *testing.T) {
    writer := &MockWriter{}
    greeter := desktop.GreeterWithWriter[*MockWriter](writer)

    ctx := context.Background()
    result := greeter.Execute(ctx, api.NewGreetCommand("Bob"))

    if !result.IsOk() {
        t.Fatalf("Expected success, got error")
    }

    if !strings.Contains(writer.String(), "Hello, Bob!") {
        t.Errorf("Expected greeting, got: %s", writer.String())
    }
}
```

### Direct Domain Usage

```go
// Use domain types directly
result := api.CreatePerson("Alice")
if result.IsOk() {
    person := result.Value()
    greeting := person.GreetingMessage()
    fmt.Println(greeting) // "Hello, Alice!"
}
```

---

## Understanding the Architecture

Hybrid_Lib_Go demonstrates **4-layer library hexagonal architecture**:

### Layer Overview

```
┌─────────────────────────────────────────────┐
│  API Layer (Public Facade)                  │  ← api/, api/adapter/desktop/
├─────────────────────────────────────────────┤
│  Infrastructure (Adapters)                  │  ← Technical implementations
├─────────────────────────────────────────────┤
│  Application (Use Cases + Ports)            │  ← Orchestration layer
├─────────────────────────────────────────────┤
│  Domain (Business Logic)                    │  ← Pure business rules (ZERO dependencies)
└─────────────────────────────────────────────┘
```

### Key Architectural Principles

1. **Domain has zero dependencies** - Pure business logic
2. **API layer does NOT import infrastructure** - Uses re-exports only
3. **Platform wiring in api/adapter/desktop/** - Creates ready-to-use instances
4. **Static dependency injection** - Via generics (compile-time wiring)
5. **Railway-oriented programming** - Result monads for error handling
6. **Multi-module workspace** - go.work manages separate go.mod per layer

### Request Flow Example

```
Consumer App
    ↓
api/adapter/desktop.Greeter (ready-to-use)
    ↓
application/usecase.GreetUseCase (validates via Domain)
    ↓
domain/valueobject.Person (business rules)
    ↓
infrastructure/adapter.ConsoleWriter (output)
    ↓
Result[Unit] (success or error)
```

---

## Making Your First Change

Let's modify the greeting message:

### Step 1: Locate the Domain Logic

```bash
# Open the Person value object
# File: domain/valueobject/person.go
```

### Step 2: Modify Greeting Format

Find the `GreetingMessage()` method and modify it. The business logic is pure and has no dependencies.

### Step 3: Rebuild and Test

```bash
# Rebuild
make build

# Run tests to ensure nothing broke
make test-all
```

**Best Practice**: Always run tests after making changes.

---

## Running Tests

Hybrid_Lib_Go includes comprehensive testing:

### Test Organization

- **Unit Tests**: Domain and application logic
- **Integration Tests**: API usage verification

### Run All Tests

```bash
# Run entire test suite
make test-all
```

### Run Specific Test Suites

```bash
# Unit tests only (fast)
make test

# Integration tests only
make test-integration
```

**Test Framework**: testify for assertions.

---

## Build Targets

### Building

```bash
make build              # Build all modules
```

### Testing

```bash
make test               # Run unit tests
make test-all           # Run all tests (unit + integration)
make test-integration   # Integration tests only
```

### Quality & Architecture

```bash
make check-arch         # Validate architecture boundaries
make fmt                # Format code with gofmt
make lint               # Run golangci-lint
```

### Cleaning

```bash
make clean              # Clean build artifacts
```

---

## Common Issues

### Q: Build fails with "go.work not found"

**A:** Ensure you're in the project root directory:

```bash
pwd
# Should show: .../hybrid_lib_go
```

### Q: "package not found" errors

**A:** Sync the workspace:

```bash
go work sync
```

### Q: Architecture validation warnings appear

**A:** The `make check-arch` target validates layer boundaries:

```bash
# View architecture validation
make check-arch
```

Violations indicate forbidden imports (e.g., API importing Infrastructure).

### Q: How do I run a single test?

**A:** Use go test directly:

```bash
# Run specific test
go test -v -run TestGreeter_Execute_Success ./test/integration/...

# Run with build tag
go test -v -tags=integration ./test/integration/...
```

---

## Next Steps

### Explore the Architecture

- **[Software Design Specification](formal/software_design_specification.md)** - Deep dive into architecture
- **[Architecture Diagrams](diagrams/)** - Visual documentation
- **[Layer Dependencies](diagrams/layer_dependencies.svg)** - See dependency flow

### Read the Source Code

Start with the API facade:

```bash
# See how types are re-exported
cat api/api.go

# See how infrastructure is wired
cat api/adapter/desktop/desktop.go
```

Then explore each layer:

```bash
# Domain (pure business logic)
ls domain/

# Application (use cases and ports)
ls application/

# Infrastructure (adapters)
ls infrastructure/
```

### Study the Test Suite

```bash
# See how tests are organized
ls -R test/
```

### Understand Static Dispatch

```go
// Generic type with interface constraint
type GreetUseCase[W outbound.WriterPort] struct {
    writer W
}

// Concrete type known at compile time
uc := usecase.NewGreetUseCase[*adapter.ConsoleWriter](writer)

// Method call is statically dispatched (no vtable)
uc.Execute(ctx, cmd)
```

### Add Your Own Use Case

Follow the pattern:

1. **Domain**: Create value objects/entities (`domain/valueobject/`)
2. **Application**: Define command, use case, ports (`application/`)
3. **Infrastructure**: Implement adapters (`infrastructure/adapter/`)
4. **API**: Re-export types (`api/api.go`)
5. **API/desktop**: Wire infrastructure (`api/adapter/desktop/desktop.go`)
6. **Tests**: Add unit/integration tests (`test/`)

---

## Documentation Index

- **[Main Documentation Hub](index.md)** - All documentation links
- **[Software Requirements Specification](formal/software_requirements_specification.md)** - Requirements
- **[Software Design Specification](formal/software_design_specification.md)** - Architecture
- **[Software Test Guide](formal/software_test_guide.md)** - Testing guide

---

## Support

For questions or issues:

- **Issues**: [GitHub Issues](https://github.com/abitofhelp/hybrid_lib_go/issues)
- **Documentation**: See `docs/` directory

---

## License

Hybrid_Lib_Go is licensed under the BSD-3-Clause License.
Copyright (c) 2025 Michael Gardner, A Bit of Help, Inc.

See [LICENSE](../LICENSE) for full license text.
