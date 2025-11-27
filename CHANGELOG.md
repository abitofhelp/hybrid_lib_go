# Changelog

**Version:** 1.0.0
**Date:** November 26, 2025
**SPDX-License-Identifier:** BSD-3-Clause
**License File:** See the LICENSE file in the project root.
**Copyright:** (c) 2025 Michael Gardner, A Bit of Help, Inc.
**Status:** Released


All notable changes to Hybrid_Lib_Go will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

---

## [1.0.0] - 2025-11-26

_First stable release - Professional Go 1.23+ library template demonstrating hybrid DDD/Clean/Hexagonal architecture with functional programming principles._

### Added

#### Architecture
- **4-Layer Library Hexagonal Architecture**: Domain, Application, Infrastructure, API
- **API Facade Pattern**: api/ re-exports types without importing infrastructure
- **Desktop Adapter**: api/adapter/desktop/ wires infrastructure as composition root
- **Static Dispatch via Generics**: Zero-overhead dependency injection using Go generics
- **Port Abstraction**: Inbound (GreetPort) and Outbound (WriterPort) port interfaces
- **Architecture Guard**: Python script (arch_guard.py) for automated boundary validation
- **Context Support**: context.Context for cancellation and deadline propagation

#### Domain Layer
- **Result[T] Monad**: Railway-oriented programming for explicit error handling
  - `Ok[T]()`/`Err[T]()` constructors
  - `IsOk()`/`IsError()` predicates
  - `Value()`/`ErrorInfo()` accessors with precondition checks
  - `UnwrapOr()`, `UnwrapOrElse()`, `Expect()` safe alternatives
  - `Map()`, `MapTo()`, `AndThen()`, `AndThenTo()` functional operations
  - `MapError()`, `Fallback()`, `FallbackWith()` error handling
  - `Recover()`, `RecoverWith()` error recovery
  - `Tap()` for side effects
- **Person Value Object**: Immutable with validation via smart constructor
- **Error Types**: ValidationError and InfrastructureError with descriptive messages

#### Application Layer
- **GreetUseCase[W WriterPort]**: Generic use case with static dispatch
- **GreetCommand DTO**: Immutable command for crossing layer boundaries
- **Unit Model**: Type-safe void/unit representation
- **Error Re-export**: application/error re-exports domain/error for outer layers (API, adapters)

#### Infrastructure Layer
- **ConsoleWriter**: Adapter implementing WriterPort
  - Panic recovery via defer/recover
  - Context cancellation support
  - io.Writer injection for testability
- `NewWriter(io.Writer)`: Factory for custom output targets
- `NewConsoleWriter()`: Convenience factory for stdout
- `NewStderrWriter()`: Factory for stderr output

#### API Layer (Public Facade)
- **api/api.go**: Re-exports types from domain and application
  - Does NOT import infrastructure (maintains architectural boundary)
  - Zero-overhead type aliases for Result, Person, GreetCommand, etc.
  - Factory functions for creating domain/application objects
- **api/adapter/desktop/**: Platform-specific composition root
  - Wires ConsoleWriter → GreetUseCase → Greeter
  - Provides ready-to-use `NewGreeter()` for console output
  - Provides `GreeterWithWriter[W]()` for custom writers

#### Testing
- **Domain Unit Tests**: 42 assertions in 2 test functions
  - Result[T] monad tests (19 assertions)
  - Person value object tests (23 assertions)
- **API Integration Tests**: 9 tests verifying consumer usage patterns
  - Valid input scenarios
  - Invalid input scenarios (empty, too long)
  - Error type propagation
- **Test Infrastructure**: testify assertions, mock writers

#### Documentation
- **Formal Documentation**:
  - Software Requirements Specification (SRS)
  - Software Design Specification (SDS)
  - Software Test Guide (STG)
- **PlantUML Diagrams** (5 diagrams with SVG):
  - `layer_dependencies.puml` - 4-layer library architecture
  - `package_structure.puml` - Package organization
  - `error_handling_flow.puml` - Railway-oriented flow
  - `static_dispatch_api.puml` - Static vs dynamic dispatch
  - `api_reexport_pattern.puml` - API facade re-export pattern
- **Guides**:
  - Quick Start Guide
  - Architecture Mapping Guide
  - Ports Mapping Guide
- **Source Code Documentation**: Comprehensive godoc for all packages

#### Build System
- **Makefile**: Comprehensive build automation
- **go.work**: Workspace for multi-module project
- **Separate go.mod per layer**: Clear dependency boundaries

#### Release Tooling
- **Release Script Validation**: Makefile target and documentation link validation
- **brand_project.py**: Graceful handling of non-existent output directories
- **SPDX Headers**: All layer READMEs include BSD-3-Clause license identifier

### Changed

#### Port Naming Conventions
- Renamed `application/port/inward` to `application/port/inbound` (driving adapters)
- Renamed `application/port/outward` to `application/port/outbound` (driven adapters)
- Updated all imports, documentation, and PlantUML diagrams

#### API Layer Structure
- Restructured `api/desktop` to `api/adapter/desktop`
- Consistent with adapter pattern used in infrastructure layer
- Aligns with application project's `presentation/adapter/cli/command` structure

#### Tooling
- Fixed arch_guard.py go.mod parsing to handle multi-line `require (...)` blocks
- arch_guard.py validates library-specific rules: api cannot import infrastructure

#### Documentation
- Updated all documentation for library architecture (4-layer, API facade)
- All formal documentation and diagrams updated with current paths

### Fixed

- **.gitattributes**: Replaced placeholder prose with actual git configuration
- **architecture_mapping.md**: Corrected file paths to match current project structure
- **ports_mapping.md**: Corrected file paths to match current project structure

### Architecture Patterns

- **API Facade Pattern**: api/ re-exports types without importing infrastructure
- **Static Dependency Injection**: Generic structs with interface constraints (compile-time DI)
- **Result Monad**: Railway-oriented error handling (no panics across boundaries)
- **Composition Root**: api/adapter/desktop/ wires all dependencies together
- **Ports & Adapters**: Application defines ports, Infrastructure implements adapters

### Technical Details

```go
// API Facade - re-exports types without infrastructure dependency
import (
    "github.com/abitofhelp/hybrid_lib_go/api"
    "github.com/abitofhelp/hybrid_lib_go/api/adapter/desktop"
)

// Ready-to-use greeter with console output
greeter := desktop.NewGreeter()
result := greeter.Execute(ctx, api.NewGreetCommand("Alice"))

// Custom writer for testing or other output
greeter := desktop.GreeterWithWriter[*MockWriter](mockWriter)
```

```go
// Static dispatch - concrete types known at compile time
consoleWriter := adapter.NewConsoleWriter()
uc := usecase.NewGreetUseCase[*adapter.ConsoleWriter](consoleWriter)

// All calls are statically dispatched - no vtable lookup
result := uc.Execute(ctx, cmd)
```

```go
// Railway-oriented error handling
result := valueobject.CreatePerson(name)
if result.IsError() {
    return domerr.Err[model.Unit](result.ErrorInfo())
}
person := result.Value()
```

### Compatibility

- **Go Version**: 1.23+ (generics with type inference)
- **Platforms**: Linux, macOS, BSD, Windows
- **Dependencies**: testify (test only)

### Notes

This is the initial release of Hybrid_Lib_Go, a Go library template companion to hybrid_app_go.
Both projects demonstrate identical architectural patterns with different entry points:

| Library (4 layers) | Application (5 layers) |
|--------------------|------------------------|
| api/ (facade)      | bootstrap/ (main)      |
| api/adapter/desktop/ | presentation/        |
| infrastructure/    | infrastructure/        |
| application/       | application/           |
| domain/            | domain/                |

Key differences:
- Library uses API facade pattern (api/ + api/adapter/desktop/)
- Application uses CLI command pattern (presentation/ + bootstrap/)
- Both use static dispatch via generics
- Both use railway-oriented programming with Result monads

---

## Release Notes Format

Each release will document changes in these categories:

- **Added** - New features
- **Changed** - Changes to existing functionality
- **Deprecated** - Soon-to-be-removed features
- **Removed** - Removed features
- **Fixed** - Bug fixes
- **Security** - Security vulnerability fixes

---

## License & Copyright

- **License**: BSD-3-Clause
- **Copyright**: (c) 2025 Michael Gardner, A Bit of Help, Inc.
- **SPDX-License-Identifier**: BSD-3-Clause

[1.0.0]: https://github.com/abitofhelp/hybrid_lib_go/releases/tag/v1.0.0
