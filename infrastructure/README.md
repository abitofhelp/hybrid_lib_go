<!-- SPDX-License-Identifier: BSD-3-Clause -->

# Infrastructure Layer

Implements outbound ports with concrete adapters for external services.

## Responsibilities

- Implement outbound port interfaces defined in application layer
- Handle I/O operations (console, file, network, database)
- Convert infrastructure errors to domain error types
- Provide panic recovery at system boundaries

## Key Packages

- `adapter/` - Concrete implementations of outbound ports

## Architectural Rules

- **Can import**: domain, application layers
- **Cannot import**: api layer
- Must implement interfaces from `application/port/outbound/`
- Convert all errors to domain.error.ErrorType
- Recover from panics and convert to Result errors

## Example Adapter

```go
// ConsoleWriter implements outbound.WriterPort
type ConsoleWriter struct {
    w io.Writer
}

func (cw *ConsoleWriter) Write(ctx context.Context, msg string) domerr.Result[model.Unit] {
    // Panic recovery wrapper
    defer func() {
        if r := recover(); r != nil {
            result = domerr.Err[model.Unit](
                domerr.NewInfrastructureError("panic recovered"))
        }
    }()

    // Context cancellation check
    select {
    case <-ctx.Done():
        return domerr.Err[model.Unit](
            domerr.NewInfrastructureError("context cancelled"))
    default:
    }

    // Actual I/O operation
    _, err := fmt.Fprintln(cw.w, msg)
    if err != nil {
        return domerr.Err[model.Unit](
            domerr.NewInfrastructureError(err.Error()))
    }

    return domerr.Ok(model.UnitValue)
}
```

## Panic Recovery Pattern

All infrastructure adapters should recover from panics at the boundary:

```go
func (adapter *SomeAdapter) Operation() (result domerr.Result[T]) {
    defer func() {
        if r := recover(); r != nil {
            result = domerr.Err[T](
                domerr.NewInfrastructureError(fmt.Sprintf("panic: %v", r)))
        }
    }()

    // ... operation that might panic ...
}
```

This ensures that panics in third-party libraries or unexpected conditions
are converted to proper Result errors and don't crash the application.
