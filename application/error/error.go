// SPDX-License-Identifier: BSD-3-Clause
// Copyright (c) 2025 Michael Gardner, A Bit of Help, Inc.
// Package: error
// Description: Error type facade for outer layers

// Package error re-exports Domain.Error types for use by outer layers.
// This maintains the architectural boundary: API/Adapters -> Application -> Domain.
//
// Architecture Notes:
//   - Part of the APPLICATION layer (orchestration/contract layer)
//   - Re-exports Domain error types without wrapping (zero overhead)
//   - Allows outer layers (API facade, adapters) to access error types without depending on Domain
//   - Infrastructure may use domain.Error directly (it's allowed to depend on domain)
//
// Why This Exists:
//
// In our hybrid architecture:
//   - Domain is the only shareable layer across applications
//   - Application/Infrastructure/API are library-specific
//   - Outer layers (API facade, consuming apps) should NOT depend on Domain directly
//   - Application acts as the contract/facade layer for outer consumers
//
// Usage (outer layers - API, adapters, consuming apps):
//
//	import apperr "github.com/abitofhelp/hybrid_lib_go/application/error"
//	// NOT: import domerr "github.com/abitofhelp/hybrid_lib_go/domain/error"
//
//	switch err.Kind {
//	case apperr.ValidationError:
//	    // Handle validation error
//	case apperr.InfrastructureError:
//	    // Handle infrastructure error
//	}
package error

import domerr "github.com/abitofhelp/hybrid_lib_go/domain/error"

// Re-export error types from domain (zero overhead type aliases)

// ErrorKind represents categories of errors (re-exported from domain)
type ErrorKind = domerr.ErrorKind

// Error kind constants (re-exported from domain)
const (
	ValidationError     = domerr.ValidationError
	InfrastructureError = domerr.InfrastructureError
)

// ErrorType is the concrete error type (re-exported from domain)
type ErrorType = domerr.ErrorType

// Result is the Result monad type (re-exported from domain)
// Outer layers use this type but do not create Results
// (Results are created by Application layer and passed to outer layers)
type Result[T any] = domerr.Result[T]

// Constructor functions (re-exported from domain)
var (
	NewValidationError     = domerr.NewValidationError
	NewInfrastructureError = domerr.NewInfrastructureError
)
