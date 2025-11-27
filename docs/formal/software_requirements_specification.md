# Software Requirements Specification (SRS)

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

This Software Requirements Specification (SRS) describes the functional and non-functional requirements for Hybrid_Lib_Go, a professional Go 1.23+ library template demonstrating hexagonal architecture with functional programming principles.

### 1.2 Scope

Hybrid_Lib_Go provides:
- Professional 4-layer hexagonal architecture implementation
- Static dependency injection via Go generics
- Railway-oriented error handling with Result monads
- API facade pattern for clean public interface
- Architecture boundary enforcement tooling

### 1.3 Definitions

| Term | Definition |
|------|------------|
| **Library** | Reusable code package consumed by applications |
| **Consumer** | Application that imports and uses this library |
| **Result Monad** | Type encapsulating success (Ok) or failure (Err) |
| **Static Dispatch** | Compile-time method resolution via generics |
| **API Facade** | Public interface re-exporting internal types |
| **Composition Root** | Location where dependencies are wired together |

### 1.4 References

- Go 1.23 Language Specification
- Hexagonal Architecture (Alistair Cockburn)
- Railway-Oriented Programming (Scott Wlaschin)

---

## 2. Overall Description

### 2.1 Product Perspective

Hybrid_Lib_Go is a standalone library template implementing professional architectural patterns:

1. **Domain**: Pure business logic (ZERO external dependencies)
2. **Application**: Use cases, ports, commands
3. **Infrastructure**: Driven adapters
4. **API**: Public facade and composition roots

### 2.2 Product Features

1. **Hexagonal Architecture**: 4-layer library architecture
2. **Type-Safe DI**: Static dispatch via generics
3. **Functional Error Handling**: Result monad pattern
4. **Clean API**: Public facade hiding internal structure
5. **Architecture Enforcement**: Automated boundary validation

### 2.3 User Classes

| User Class | Description | Needs |
|------------|-------------|-------|
| **Library Consumers** | Developers using the library | Clean API, documentation, examples |
| **Library Maintainers** | Developers maintaining the library | Clear architecture, tests, tooling |
| **Template Users** | Developers creating new libraries | Branding scripts, patterns to follow |

### 2.4 Operating Environment

- **Runtime**: Go 1.23+ (generics and workspace support)
- **Build Tools**: Make, Go toolchain
- **Optional**: Python 3 (architecture validation)

### 2.5 Constraints

- Domain layer must have zero external module dependencies
- API layer must not import infrastructure directly
- All error handling must use Result monad (no panics across boundaries)

---

## 3. Functional Requirements

### 3.1 Domain Layer (FR-01)

**FR-01.1**: Implement Person value object with validation
- Name must not be empty
- Name must not exceed 100 characters
- Immutable after creation

**FR-01.2**: Implement ErrorType with error kinds
- ValidationError for input validation failures
- InfrastructureError for I/O and system failures

**FR-01.3**: Implement Result[T] monad
- Ok variant for successful values
- Err variant for error information
- Methods: IsOk(), IsError(), Value(), ErrorInfo()
- Functional operations: Map(), AndThen(), MapError(), etc.

**FR-01.4**: Implement Option[T] monad
- Some variant for present values
- None variant for absent values
- Methods: IsSome(), IsNone(), Value()

**FR-01.5**: Zero external dependencies
- No third-party imports in domain layer
- Only Go standard library allowed

### 3.2 Application Layer (FR-02)

**FR-02.1**: Implement GreetCommand input DTO
- Encapsulate name parameter
- Factory function for creation

**FR-02.2**: Define WriterPort output interface
- Write(ctx, message) returns Result[Unit]
- Context support for cancellation

**FR-02.3**: Define GreetPort input interface
- Execute(ctx, cmd) returns Result[Unit]
- Generic type parameter for writer

**FR-02.4**: Implement GreetUseCase generic use case
- Generic over WriterPort implementation
- Orchestrate domain validation and output
- Return Result[Unit]

**FR-02.5**: Implement Unit type
- Represents void return value
- Used in Result[Unit] for operations with no return value

**FR-02.6**: Re-export Domain.Error for outer layers
- Type aliases for error types
- Zero overhead re-exports

### 3.3 Infrastructure Layer (FR-03)

**FR-03.1**: Implement ConsoleWriter adapter
- Implement WriterPort interface
- Write to stdout by default
- Support custom io.Writer injection

**FR-03.2**: Implement panic recovery
- Catch panics at adapter boundary
- Convert to InfrastructureError Result

**FR-03.3**: Implement context cancellation
- Check context before I/O operations
- Return InfrastructureError on cancellation

### 3.4 API Layer (FR-04)

**FR-04.1**: Re-export domain types in api/
- Result, ErrorType, ErrorKind, Person
- ValidationError, InfrastructureError constants
- MaxNameLength constant

**FR-04.2**: Re-export application types in api/
- GreetCommand, WriterPort, Unit
- Factory functions

**FR-04.3**: Do NOT import infrastructure in api/
- Facade only re-exports types
- No infrastructure wiring

**FR-04.4**: Implement Greeter in api/adapter/desktop/
- Ready-to-use greeter with console output
- NewGreeter() factory function

**FR-04.5**: Implement GreeterWithWriter in api/adapter/desktop/
- Accept custom WriterPort implementation
- Generic over writer type

### 3.5 Testing (FR-05)

**FR-05.1**: Unit tests for domain layer
- Test Person validation
- Test Result monad operations

**FR-05.2**: Unit tests for application layer
- Test use case with mock writer
- Test error propagation

**FR-05.3**: Integration tests for API layer
- Test consumer usage patterns
- Test error handling

---

## 4. Non-Functional Requirements

### 4.1 Performance (NFR-01)

**NFR-01.1**: Zero runtime overhead for DI
- Static dispatch eliminates vtable lookups
- Method calls must be inlinable

**NFR-01.2**: No heap allocation for Result values
- Result struct should be stack-allocatable
- No boxing of primitive values

### 4.2 Reliability (NFR-02)

**NFR-02.1**: No panics across layer boundaries
- All panics caught at infrastructure boundary
- Converted to Result errors

**NFR-02.2**: Explicit error handling
- All fallible operations return Result
- No silent failures

### 4.3 Maintainability (NFR-03)

**NFR-03.1**: Enforce module boundaries
- go.mod prevents forbidden imports
- arch_guard.py validates at build time

**NFR-03.2**: Test coverage
- Domain layer: 90%+ coverage
- Application layer: 90%+ coverage
- Infrastructure layer: 70%+ coverage

### 4.4 Portability (NFR-04)

**NFR-04.1**: Go 1.23+ compatibility
- Use generics (Go 1.18+)
- Use workspaces (Go 1.18+)

**NFR-04.2**: Cross-platform support
- No platform-specific code in domain/application
- Infrastructure adapters may be platform-specific

### 4.5 Documentation (NFR-05)

**NFR-05.1**: Comprehensive documentation
- README with usage examples
- API reference
- Architecture diagrams

**NFR-05.2**: Code documentation
- SPDX headers on all source files
- GoDoc comments on public APIs

---

## 5. System Constraints

### 5.1 Architectural Constraints

| ID | Constraint | Enforced By |
|----|------------|-------------|
| SC-01 | Domain has zero external dependencies | go.mod |
| SC-02 | Application depends only on Domain | go.mod |
| SC-03 | Infrastructure depends on Application + Domain | go.mod |
| SC-04 | api/ depends on Application + Domain (NOT Infrastructure) | go.mod, arch_guard.py |
| SC-05 | api/adapter/desktop/ is composition root (can depend on all) | go.mod |
| SC-06 | All errors use Result monad | Code review |
| SC-07 | No panics across boundaries | Panic recovery pattern |

### 5.2 Dependency Rules

```
api/adapter/desktop/  ────────────────────────────────┐
    │                                                 │
    ▼                                                 ▼
   api/ ─────────────────┐                   infrastructure/
    │                    │                            │
    │                    ▼                            │
    │              application/ ◄─────────────────────┘
    │                    │
    ▼                    ▼
              ─────► domain/ ◄─────
```

---

## 6. Verification and Validation

### 6.1 Test Matrix

| Requirement | Test Type | Test Location | Status |
|-------------|-----------|---------------|--------|
| FR-01 (Domain) | Unit | domain/error/result_test.go, domain/valueobject/person_test.go | Pass |
| FR-02 (Application) | Unit | application/usecase/main_test.go | Pass |
| FR-03 (Infrastructure) | Unit | infrastructure/adapter/main_test.go | Pass |
| FR-04 (API) | Integration | test/integration/greet_flow_test.go | Pass |
| FR-05 (Testing) | Meta | Verified | Pass |

### 6.2 Architecture Validation

```bash
# Validate architecture boundaries
make check-arch

# Expected output: ✓ Architecture validation PASSED
```

### 6.3 Build Verification

```bash
# Build all modules
make build

# Run all tests
make test-all
```

---

## 7. Traceability Matrix

### 7.1 Requirements to Components

| Requirement | Component(s) |
|-------------|--------------|
| FR-01.1 | domain/valueobject/person.go |
| FR-01.2 | domain/error/error.go |
| FR-01.3 | domain/error/result.go |
| FR-01.4 | domain/valueobject/option.go |
| FR-02.1 | application/command/greet.go |
| FR-02.2 | application/port/outbound/writer.go |
| FR-02.3 | application/port/inbound/greet.go |
| FR-02.4 | application/usecase/greet.go |
| FR-02.5 | application/model/unit.go |
| FR-02.6 | application/error/error.go |
| FR-03.1-3 | infrastructure/adapter/consolewriter.go |
| FR-04.1-3 | api/api.go |
| FR-04.4-5 | api/adapter/desktop/desktop.go |

### 7.2 Requirements to Tests

| Requirement | Test File(s) |
|-------------|--------------|
| FR-01 | domain/valueobject/person_test.go, domain/error/result_test.go |
| FR-02 | application/usecase/main_test.go |
| FR-03 | infrastructure/adapter/main_test.go |
| FR-04 | test/integration/greet_flow_test.go |

---

## 8. Appendices

### Appendix A: Layer Summary

| Layer | Dependencies | Purpose | Test Type |
|-------|--------------|---------|-----------|
| Domain | NONE | Business logic | Unit |
| Application | Domain | Use cases, ports | Unit |
| Infrastructure | App + Domain | Adapters | Unit/Integration |
| API | App + Domain | Public facade | Integration |
| API/adapter/desktop | ALL | Composition root | Integration |

### Appendix B: API Usage Example

```go
import (
    "context"
    "github.com/abitofhelp/hybrid_lib_go/api"
    "github.com/abitofhelp/hybrid_lib_go/api/adapter/desktop"
)

func main() {
    greeter := desktop.NewGreeter()
    ctx := context.Background()

    result := greeter.Execute(ctx, api.NewGreetCommand("Alice"))

    if result.IsOk() {
        // Success
    } else {
        switch result.ErrorInfo().Kind {
        case api.ValidationError:
            // Handle validation error
        case api.InfrastructureError:
            // Handle infrastructure error
        }
    }
}
```

---

**Document History**

| Version | Date | Author | Changes |
|---------|------|--------|---------|
| 1.0.0 | 2025-11-26 | Michael Gardner | Initial library version |
