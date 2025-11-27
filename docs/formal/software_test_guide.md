# Software Test Guide

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

This Software Test Guide describes the testing strategy, organization, and procedures for Hybrid_Lib_Go library.

### 1.2 Scope

This document covers:
- Test organization and structure
- Test execution procedures
- Coverage requirements
- Testing patterns for hexagonal architecture

### 1.3 Test Framework

- **Unit Tests**: Go standard testing package with testify/assert
- **Integration Tests**: Go testing with build tags
- **Architecture Validation**: Python arch_guard.py script

---

## 2. Test Organization

### 2.1 Test Structure

```
hybrid_lib_go/
├── domain/
│   ├── error/
│   │   └── result_test.go          # Result monad tests
│   └── valueobject/
│       └── person_test.go          # Person validation tests
├── application/
│   └── usecase/
│       └── greet_usecase_test.go   # Use case tests
├── infrastructure/
│   └── adapter/
│       └── console_writer_test.go  # Adapter tests
├── test/
│   ├── integration/
│   │   └── greet_flow_test.go      # API integration tests
│   └── python/
│       ├── conftest.py             # Pytest fixtures
│       └── test_arch_guard_go.py   # Architecture validation tests
└── scripts/
    └── arch_guard/
        └── arch_guard.py           # Architecture validator
```

### 2.2 Test Types

| Type | Location | Build Tag | Purpose |
|------|----------|-----------|---------|
| Unit | Co-located `*_test.go` | None | Test individual components |
| Integration | `test/integration/` | `integration` | Test API usage patterns |
| Architecture | `test/python/` | N/A (pytest) | Validate layer boundaries |

---

## 3. Running Tests

### 3.1 Make Targets

```bash
# Run all tests (unit + integration)
make test

# Run unit tests only
make test-unit

# Run integration tests only
make test-integration

# Run all tests with verbose output
make test-all

# Run tests with coverage
make test-coverage
```

### 3.2 Direct Go Commands

```bash
# Run all unit tests
go test ./domain/... ./application/... ./infrastructure/... ./api/...

# Run specific test
go test -v -run TestGreeter_Execute_Success ./test/integration/...

# Run integration tests
go test -v -tags=integration ./test/integration/...

# Run with coverage
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

### 3.3 Architecture Validation

```bash
# Run architecture validation
make check-arch

# Or directly
python3 scripts/arch_guard/arch_guard.py

# Run Python architecture tests
cd test/python && pytest -v test_arch_guard_go.py
```

---

## 4. Test Patterns

### 4.1 Domain Layer Testing

**Pattern**: Pure unit tests with no mocks

```go
func TestCreatePerson_EmptyName_ReturnsValidationError(t *testing.T) {
    result := valueobject.CreatePerson("")

    assert.True(t, result.IsError())
    assert.Equal(t, domerr.ValidationError, result.ErrorInfo().Kind)
}

func TestCreatePerson_ValidName_ReturnsOk(t *testing.T) {
    result := valueobject.CreatePerson("Alice")

    assert.True(t, result.IsOk())
    assert.Equal(t, "Alice", result.Value().Name())
}
```

**Key Points**:
- No mocks needed (pure functions)
- Test business rules and validation
- Test Result monad behavior

### 4.2 Application Layer Testing

**Pattern**: Mock output ports

```go
// MockWriter implements WriterPort for testing
type MockWriter struct {
    Buffer  bytes.Buffer
    Err     error
}

func (w *MockWriter) Write(ctx context.Context, msg string) domerr.Result[model.Unit] {
    if w.Err != nil {
        return domerr.Err[model.Unit](domerr.NewInfrastructureError(w.Err.Error()))
    }
    w.Buffer.WriteString(msg)
    return domerr.Ok(model.UnitValue)
}

func TestGreetUseCase_Execute_Success(t *testing.T) {
    writer := &MockWriter{}
    uc := usecase.NewGreetUseCase[*MockWriter](writer)
    cmd := command.NewGreetCommand("Alice")

    result := uc.Execute(context.Background(), cmd)

    assert.True(t, result.IsOk())
    assert.Contains(t, writer.Buffer.String(), "Hello, Alice!")
}
```

**Key Points**:
- Mock only output ports
- Use real domain objects
- Test orchestration logic

### 4.3 Infrastructure Layer Testing

**Pattern**: Test with real I/O or injected writers

```go
func TestConsoleWriter_Write_Success(t *testing.T) {
    var buf bytes.Buffer
    writer := adapter.NewConsoleWriterWithOutput(&buf)

    result := writer.Write(context.Background(), "Hello, Test!")

    assert.True(t, result.IsOk())
    assert.Equal(t, "Hello, Test!\n", buf.String())
}

func TestConsoleWriter_Write_ContextCancelled(t *testing.T) {
    ctx, cancel := context.WithCancel(context.Background())
    cancel() // Cancel immediately

    writer := adapter.NewConsoleWriter()
    result := writer.Write(ctx, "Should not write")

    assert.True(t, result.IsError())
    assert.Equal(t, domerr.InfrastructureError, result.ErrorInfo().Kind)
}
```

**Key Points**:
- Test context cancellation
- Test panic recovery
- Inject test writers where possible

### 4.4 API Integration Testing

**Pattern**: Test consumer usage patterns

```go
//go:build integration

func TestGreeter_Execute_Success(t *testing.T) {
    writer := &MockWriter{}
    greeter := desktop.GreeterWithWriter[*MockWriter](writer)
    ctx := context.Background()

    result := greeter.Execute(ctx, api.NewGreetCommand("Alice"))

    assert.True(t, result.IsOk())
    assert.Contains(t, writer.String(), "Hello, Alice!")
}

func TestGreeter_Execute_EmptyName_ReturnsValidationError(t *testing.T) {
    writer := &MockWriter{}
    greeter := desktop.GreeterWithWriter[*MockWriter](writer)
    ctx := context.Background()

    result := greeter.Execute(ctx, api.NewGreetCommand(""))

    assert.True(t, result.IsError())
    assert.Equal(t, api.ValidationError, result.ErrorInfo().Kind)
}
```

**Key Points**:
- Use `//go:build integration` tag
- Test through public API
- Verify error handling

---

## 5. Coverage Requirements

### 5.1 Per-Layer Coverage Targets

| Layer | Target | Rationale |
|-------|--------|-----------|
| Domain | 90%+ | Core business logic must be thoroughly tested |
| Application | 90%+ | Use cases orchestrate critical flows |
| Infrastructure | 70%+ | Some I/O paths harder to test |
| API | Integration | Tested via integration tests |

### 5.2 Running Coverage

```bash
# Full coverage report
make test-coverage

# Per-layer coverage
go test -coverprofile=coverage.out ./domain/...
go tool cover -func=coverage.out
```

### 5.3 Coverage Report

```bash
# Generate HTML report
go tool cover -html=coverage/coverage.out -o coverage/coverage.html

# View text summary
cat coverage/coverage_summary.txt
```

---

## 6. Architecture Testing

### 6.1 arch_guard.py Validation

The `arch_guard.py` script validates:
- go.mod dependencies are correct
- Source files don't import forbidden packages
- API layer doesn't import infrastructure

```bash
# Run validation
make check-arch

# Expected output:
# ✓ go.mod Configuration: VALID
# ✓ Source File Dependencies: 0 violation(s)
# ✓ Architecture validation PASSED
```

### 6.2 Python Test Suite

```bash
# Run architecture test suite
cd test/python
pytest -v test_arch_guard_go.py

# Test specific rule
pytest -v test_arch_guard_go.py::test_api_cannot_import_infrastructure
```

### 6.3 Key Rules Tested

| Rule | Test |
|------|------|
| Domain has no deps | `test_domain_has_no_dependencies` |
| App depends on Domain only | `test_application_depends_on_domain` |
| api/ cannot import infra | `test_api_cannot_import_infrastructure` |
| api/adapter/desktop/ can import all | `test_api_desktop_can_import_all` |

---

## 7. Test Data

### 7.1 Valid Test Cases

| Input | Expected Result |
|-------|-----------------|
| `"Alice"` | Ok, "Hello, Alice!" |
| `"Bob Smith"` | Ok, "Hello, Bob Smith!" |
| `"世界"` | Ok, "Hello, 世界!" |

### 7.2 Error Test Cases

| Input | Expected Error |
|-------|----------------|
| `""` (empty) | ValidationError: name cannot be empty |
| 101+ chars | ValidationError: name exceeds max length |
| Cancelled context | InfrastructureError: context cancelled |

---

## 8. Continuous Integration

### 8.1 CI Pipeline Steps

1. **Build**: `make build`
2. **Unit Tests**: `make test-unit`
3. **Integration Tests**: `make test-integration`
4. **Architecture Validation**: `make check-arch`
5. **Linting**: `make lint`

### 8.2 Pre-commit Checks

```bash
# Full validation before commit
make check-arch && make test-all && make lint
```

---

## 9. Troubleshooting

### 9.1 Common Issues

**Test fails with "package not found"**
```bash
# Sync workspace
go work sync
```

**Architecture validation fails**
```bash
# Check specific violation
make check-arch
# Look for "FORBIDDEN_LATERAL_DEPENDENCY" or "ILLEGAL_LAYER_DEPENDENCY"
```

**Integration tests not running**
```bash
# Ensure build tag is used
go test -tags=integration ./test/integration/...
```

### 9.2 Debug Mode

```bash
# Run single test with verbose output
go test -v -run TestSpecificTest ./path/to/package/...

# Run with race detector
go test -race ./...
```

---

## 10. Test Maintenance

### 10.1 Adding New Tests

1. **Unit tests**: Add `*_test.go` in same package
2. **Integration tests**: Add to `test/integration/` with build tag
3. **Update coverage**: Ensure new code has tests

### 10.2 Test Naming Convention

```
Test{Component}_{Method}_{Scenario}

Examples:
- TestCreatePerson_EmptyName_ReturnsValidationError
- TestGreeter_Execute_Success
- TestConsoleWriter_Write_ContextCancelled
```

---

**Document History**

| Version | Date | Author | Changes |
|---------|------|--------|---------|
| 1.0.0 | 2025-11-26 | Michael Gardner | Initial library version |
