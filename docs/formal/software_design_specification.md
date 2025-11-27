# Software Design Specification (SDS)

**Project:** Hybrid_Lib_Go - Go 1.23+ Library
**Version:** 1.0.0
**Date:** November 26, 2025
**SPDX-License-Identifier:** BSD-3-Clause
**License File:** See the LICENSE file in the project root.
**Copyright:** (c) 2025 Michael Gardner, A Bit of Help, Inc.
**Status:** Released

---

## 1. Introduction

### 1.1 Purpose

This Software Design Specification (SDS) describes the architectural design and detailed design of Hybrid_Lib_Go, a professional Go 1.23+ library demonstrating hexagonal architecture with functional programming principles.

### 1.2 Scope

This document covers:
- 4-layer library organization and dependencies
- Static dispatch via generics pattern
- Railway-oriented error handling design
- API facade and composition root patterns
- Module boundary enforcement

### 1.3 Definitions

| Term | Definition |
|------|------------|
| **Hexagonal Architecture** | Ports and Adapters pattern isolating business logic |
| **Static Dispatch** | Compile-time method resolution via generics |
| **Result Monad** | Functional error handling pattern (Ok/Err) |
| **API Facade** | Public interface layer re-exporting internal types |
| **Composition Root** | Location where dependencies are wired together |

---

## 2. Architectural Design

### 2.1 System Architecture Overview

Hybrid_Lib_Go implements a **4-layer hexagonal architecture** for libraries:

```
┌─────────────────────────────────────────────────────────────┐
│                    Consumer Application                      │
└─────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────┐
│  API Layer                                                   │
│  ┌─────────────────────┬───────────────────────────────────┐│
│  │  api/ (facade)      │  api/adapter/desktop/ (composition root)  ││
│  │  Re-exports types   │  Wires infrastructure             ││
│  │  App + Domain only  │  ALL layers                       ││
│  └─────────────────────┴───────────────────────────────────┘│
├─────────────────────────────────────────────────────────────┤
│  Infrastructure                                              │
│  Driven adapters (ConsoleWriter)                            │
│  Depends on: Application + Domain                           │
├─────────────────────────────────────────────────────────────┤
│  Application                                                 │
│  Use cases, ports, commands                                 │
│  Depends on: Domain ONLY                                    │
├─────────────────────────────────────────────────────────────┤
│  Domain                                                      │
│  Pure business logic, value objects, Result monad           │
│  Depends on: NOTHING (zero external dependencies)           │
└─────────────────────────────────────────────────────────────┘
```

### 2.2 Layer Descriptions

#### Domain Layer

**Purpose**: Pure business logic with zero external dependencies

**Components**:
- Value Objects: `domain/valueobject/person.go`
- Error Types: `domain/error/error.go`
- Result Monad: `domain/error/result.go`

**Characteristics**:
- ZERO external module dependencies
- All types are immutable
- Validation returns Result (no exceptions)
- Pure functions where possible

#### Application Layer

**Purpose**: Use case orchestration and port definitions

**Components**:
- Use Cases: `application/usecase/greet.go`
- Commands: `application/command/greet.go`
- Output Ports: `application/port/outbound/writer.go`
- Input Ports: `application/port/inbound/greet.go`
- Models: `application/model/unit.go`
- Error Re-exports: `application/error/error.go`

**Characteristics**:
- Depends ONLY on Domain
- Defines port interfaces (contracts)
- Generic use cases for static dispatch

#### Infrastructure Layer

**Purpose**: External system adapters

**Components**:
- Console Writer: `infrastructure/adapter/consolewriter.go`

**Characteristics**:
- Implements output ports from Application
- Uses Domain types directly
- Panic recovery at boundaries
- Context cancellation support

#### API Layer

**Purpose**: Public interface for library consumers

**Components**:
- Facade: `api/api.go` - Re-exports types
- Desktop Composition: `api/adapter/desktop/desktop.go` - Wires infrastructure

**Characteristics**:
- `api/` depends on Application + Domain (NOT Infrastructure)
- `api/adapter/desktop/` depends on ALL layers (composition root)
- Provides ready-to-use instances and custom wiring options

### 2.3 Module Structure

```
hybrid_lib_go/
├── go.work                    # Workspace definition
├── domain/
│   ├── go.mod                 # ZERO dependencies
│   ├── error/
│   │   ├── error.go           # ErrorType, ErrorKind
│   │   └── result.go          # Result[T] monad
│   └── valueobject/
│       └── person.go          # Person value object
├── application/
│   ├── go.mod                 # depends: domain
│   ├── command/
│   │   └── greet_command.go   # GreetCommand DTO
│   ├── error/
│   │   └── error.go           # Re-exports domain/error
│   ├── model/
│   │   └── unit.go            # Unit type (void)
│   ├── port/
│   │   ├── inbound/
│   │   │   └── greet_port.go  # GreetPort interface
│   │   └── outbound/
│   │       └── writer_port.go # WriterPort interface
│   └── usecase/
│       └── greet_usecase.go   # GreetUseCase[W]
├── infrastructure/
│   ├── go.mod                 # depends: application, domain
│   └── adapter/
│       └── console_writer.go  # ConsoleWriter
├── api/
│   ├── go.mod                 # depends: application, domain (NOT infrastructure)
│   ├── api.go                 # Public facade
│   └── desktop/
│       ├── go.mod             # depends: ALL
│       └── desktop.go         # Composition root
└── test/
    └── integration/           # API usage tests
```

---

## 3. Detailed Design

### 3.1 Domain Layer Design

**domain/error/result.go** (Result Monad):

```go
type Result[T any] struct {
    value    T
    err      ErrorType
    hasError bool
}

func Ok[T any](value T) Result[T]
func Err[T any](err ErrorType) Result[T]
func (r Result[T]) IsOk() bool
func (r Result[T]) IsError() bool
func (r Result[T]) Value() T
func (r Result[T]) ErrorInfo() ErrorType
```

**domain/valueobject/person.go**:

```go
const MaxNameLength = 100

type Person struct {
    name string  // Immutable
}

func CreatePerson(name string) domerr.Result[Person] {
    // Validation logic
    // Returns Ok(Person) or Err(ValidationError)
}

func (p Person) Name() string
func (p Person) GreetingMessage() string
```

**Design Decision - Error Type Naming**:

The `ErrorKind` and `ErrorType` names intentionally include the "Error" prefix despite the package name being `error` (which causes linter stutter warnings like `error.ErrorKind`). This is an intentional design decision because these types are re-exported through the API facade as `api.ErrorKind` and `api.ErrorType`, where the full names provide clarity to library consumers.

### 3.2 Application Layer Design

**application/usecase/greet_usecase.go** (Generic Use Case):

```go
type GreetUseCase[W outbound.WriterPort] struct {
    writer W
}

func NewGreetUseCase[W outbound.WriterPort](writer W) *GreetUseCase[W] {
    return &GreetUseCase[W]{writer: writer}
}

func (uc *GreetUseCase[W]) Execute(ctx context.Context, cmd command.GreetCommand) domerr.Result[model.Unit] {
    personResult := valueobject.CreatePerson(cmd.Name())
    if personResult.IsError() {
        return domerr.Err[model.Unit](personResult.ErrorInfo())
    }
    person := personResult.Value()
    return uc.writer.Write(ctx, person.GreetingMessage())
}
```

**application/port/outbound/writer_port.go**:

```go
type WriterPort interface {
    Write(ctx context.Context, message string) domerr.Result[model.Unit]
}
```

### 3.3 Infrastructure Layer Design

**infrastructure/adapter/console_writer.go**:

```go
type ConsoleWriter struct {
    w io.Writer
}

func NewConsoleWriter() *ConsoleWriter {
    return &ConsoleWriter{w: os.Stdout}
}

func (cw *ConsoleWriter) Write(ctx context.Context, message string) (result domerr.Result[model.Unit]) {
    // Panic recovery
    defer func() {
        if r := recover(); r != nil {
            result = domerr.Err[model.Unit](...)
        }
    }()

    // Context cancellation check
    select {
    case <-ctx.Done():
        return domerr.Err[model.Unit](...)
    default:
    }

    // Write operation
    _, err := fmt.Fprintln(cw.w, message)
    if err != nil {
        return domerr.Err[model.Unit](...)
    }
    return domerr.Ok(model.UnitValue)
}
```

### 3.4 API Layer Design

**api/api.go** (Public Facade):

```go
package api

// Re-exported types (zero overhead aliases)
type Result[T any] = domerr.Result[T]
type ErrorType = domerr.ErrorType
type Person = valueobject.Person
type GreetCommand = command.GreetCommand
type WriterPort = outbound.WriterPort
type Unit = model.Unit

// Constants
const (
    ValidationError     = domerr.ValidationError
    InfrastructureError = domerr.InfrastructureError
    MaxNameLength       = valueobject.MaxNameLength
)

// Factory functions
func Ok[T any](value T) Result[T]
func Err[T any](err ErrorType) Result[T]
func CreatePerson(name string) Result[Person]
func NewGreetCommand(name string) GreetCommand
```

**api/adapter/desktop/desktop.go** (Composition Root):

```go
package desktop

type Greeter struct {
    useCase *usecase.GreetUseCase[*adapter.ConsoleWriter]
}

func NewGreeter() *Greeter {
    writer := adapter.NewConsoleWriter()
    uc := usecase.NewGreetUseCase[*adapter.ConsoleWriter](writer)
    return &Greeter{useCase: uc}
}

func (g *Greeter) Execute(ctx context.Context, cmd api.GreetCommand) api.Result[api.Unit] {
    return g.useCase.Execute(ctx, cmd)
}

// Custom writer support
func GreeterWithWriter[W api.WriterPort](writer W) *GreeterCustom[W]
```

---

## 4. Design Patterns

### 4.1 Static Dispatch Pattern

**Purpose**: Zero-overhead dependency injection

**Implementation**:
```go
// Generic type with interface constraint
type GreetUseCase[W WriterPort] struct {
    writer W  // Concrete type parameter
}

// Instantiation with concrete type
uc := NewGreetUseCase[*adapter.ConsoleWriter](writer)

// Method calls are statically dispatched
uc.Execute(ctx, cmd)  // Compiler knows exact type
```

**Benefits**:
- No vtable lookups (zero runtime overhead)
- Full inlining potential
- Compile-time type verification

### 4.2 API Facade Pattern

**Purpose**: Clean public interface hiding internal structure

**Implementation**:
- `api/` re-exports types via type aliases
- `api/adapter/desktop/` wires infrastructure
- Consumers import only `api` and `api/adapter/desktop`

**Key Rule**: `api/` does NOT import infrastructure

### 4.3 Railway-Oriented Error Handling

**Purpose**: Explicit error propagation without exceptions

**Pattern**:
```go
func Execute(...) Result[Unit] {
    result1 := operation1()
    if result1.IsError() {
        return Err[Unit](result1.ErrorInfo())
    }

    result2 := operation2(result1.Value())
    if result2.IsError() {
        return Err[Unit](result2.ErrorInfo())
    }

    return Ok(Unit{})
}
```

### 4.4 Panic Recovery Pattern

**Purpose**: Convert panics to Result errors at boundaries

**Implementation**:
```go
func (cw *ConsoleWriter) Write(...) (result Result[Unit]) {
    defer func() {
        if r := recover(); r != nil {
            result = Err[Unit](NewInfrastructureError(...))
        }
    }()
    // ... operation that might panic
}
```

---

## 5. Data Flow

### 5.1 Success Path

```
Consumer
    │
    ▼
api/adapter/desktop.Greeter.Execute(ctx, api.NewGreetCommand("Alice"))
    │
    ▼
usecase.GreetUseCase[*ConsoleWriter].Execute(ctx, cmd)
    │
    ├─► valueobject.CreatePerson("Alice") → Ok(Person)
    │
    └─► cw.Write(ctx, "Hello, Alice!") → Ok(Unit)
    │
    ▼
Result[Unit]{value: Unit{}, hasError: false}
```

### 5.2 Validation Error Path

```
Consumer
    │
    ▼
greeter.Execute(ctx, api.NewGreetCommand(""))
    │
    ▼
usecase.Execute(ctx, cmd)
    │
    └─► valueobject.CreatePerson("") → Err(ValidationError)
    │
    ▼
Result[Unit]{err: {Kind: ValidationError, Message: "..."}, hasError: true}
```

### 5.3 Infrastructure Error Path

```
Consumer
    │
    ▼
greeter.Execute(ctx, api.NewGreetCommand("Alice"))
    │
    ▼
usecase.Execute(ctx, cmd)
    │
    ├─► CreatePerson("Alice") → Ok(Person)
    │
    └─► cw.Write(ctx, "Hello, Alice!") → [I/O Error or Panic]
    │
    ▼
Result[Unit]{err: {Kind: InfrastructureError, Message: "..."}, hasError: true}
```

---

## 6. Module Dependencies

### 6.1 Dependency Matrix

| Module | domain | application | infrastructure | api | api/adapter/desktop |
|--------|--------|-------------|----------------|-----|-------------|
| domain | - | | | | |
| application | X | - | | | |
| infrastructure | X | X | - | | |
| api | X | X | | - | |
| api/adapter/desktop | X | X | X | X | - |

### 6.2 go.mod Dependencies

**domain/go.mod**: No require statements (zero dependencies)

**application/go.mod**:
```go
require github.com/abitofhelp/hybrid_lib_go/domain v0.0.0
```

**infrastructure/go.mod**:
```go
require (
    github.com/abitofhelp/hybrid_lib_go/application v0.0.0
    github.com/abitofhelp/hybrid_lib_go/domain v0.0.0
)
```

**api/go.mod**:
```go
require (
    github.com/abitofhelp/hybrid_lib_go/application v0.0.0
    github.com/abitofhelp/hybrid_lib_go/domain v0.0.0
)
// Note: NO infrastructure dependency
```

**api/adapter/desktop/go.mod**:
```go
require (
    github.com/abitofhelp/hybrid_lib_go/api v0.0.0
    github.com/abitofhelp/hybrid_lib_go/application v0.0.0
    github.com/abitofhelp/hybrid_lib_go/domain v0.0.0
    github.com/abitofhelp/hybrid_lib_go/infrastructure v0.0.0
)
```

---

## 7. Testing Design

### 7.1 Test Organization

| Test Type | Location | Purpose |
|-----------|----------|---------|
| Unit | Co-located `*_test.go` | Test individual components |
| Integration | `test/integration/` | Test API usage patterns |

### 7.2 Test Strategy

- **Domain**: Pure unit tests, no mocks needed
- **Application**: Test use cases with mock writers
- **Infrastructure**: Test with real I/O where practical
- **API**: Integration tests verifying consumer patterns

---

## 8. Architecture Enforcement

### 8.1 Enforcement Mechanisms

1. **go.mod**: Compiler rejects forbidden imports
2. **arch_guard.py**: Validates dependencies at build time
3. **make check-arch**: CI/CD validation target

### 8.2 Key Rules Enforced

- Domain has zero external dependencies
- Application depends only on Domain
- api/ does NOT import infrastructure
- api/adapter/desktop/ is the only composition root

---

## 9. Appendices

### Appendix A: UML Diagrams

- `docs/diagrams/layer_dependencies.svg` - Layer dependency flow
- `docs/diagrams/package_structure.svg` - Package hierarchy
- `docs/diagrams/error_handling_flow.svg` - Error propagation
- `docs/diagrams/static_dispatch_api.svg` - Generic vs interface comparison
- `docs/diagrams/api_reexport_pattern.svg` - API facade pattern

### Appendix B: References

- [Hexagonal Architecture](https://alistair.cockburn.us/hexagonal-architecture/)
- [Railway-Oriented Programming](https://fsharpforfunandprofit.com/rop/)
- [Go Generics](https://go.dev/doc/tutorial/generics)

---

**Document History**

| Version | Date | Author | Changes |
|---------|------|--------|---------|
| 1.0.0 | 2025-11-26 | Michael Gardner | Initial library version |
