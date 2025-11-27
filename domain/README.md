<!-- SPDX-License-Identifier: BSD-3-Clause -->

# Domain Layer

The innermost layer containing pure business logic with **zero external dependencies**.

## Responsibilities

- Define core business entities and value objects
- Implement business rules and validation
- Define error types and Result monad for functional error handling
- Remain completely isolated from infrastructure concerns

## Key Packages

- `error/` - Error types (ErrorKind, ErrorType) and Result[T] monad implementation
- `valueobject/` - Immutable value objects (Person, Option[T])
- `test/` - Reusable test framework

## Architectural Rules

- **No imports** from application, infrastructure, or api layers
- **No external dependencies** - only Go standard library
- All types should be immutable where possible
- Use Result[T] for operations that can fail (no panics)

## Example

```go
// Creating a Person value object with validation
result := valueobject.CreatePerson("Alice")
if result.IsOk() {
    person := result.Value()
    greeting := person.GreetingMessage() // "Hello, Alice!"
} else {
    // Handle validation error
    err := result.ErrorInfo()
    fmt.Println(err.Message)
}
```

## Result Monad

The domain provides a custom Result[T] monad with zero external dependencies:

```go
// Create results
ok := domerr.Ok[int](42)
err := domerr.Err[int](domerr.NewValidationError("invalid"))

// Query
if ok.IsOk() { ... }
if err.IsError() { ... }

// Extract values
value := ok.Value()           // Panics if error
errInfo := err.ErrorInfo()    // Panics if ok

// Safe alternatives
value := ok.UnwrapOr(0)       // Returns default if error
value := ok.UnwrapOrElse(fn)  // Computes default if error

// Functional operations
mapped := ok.Map(func(x int) int { return x * 2 })
chained := ok.AndThen(func(x int) Result[int] { return validate(x) })
```
