// SPDX-License-Identifier: BSD-3-Clause
// Copyright (c) 2025 Michael Gardner, A Bit of Help, Inc.
// Package: api
// Description: Public API facade for hybrid_lib_go

// Package api provides the public API facade for the hybrid_lib_go library.
//
// This package re-exports domain types and application ports, providing a
// single entry point for library consumers. Infrastructure is hidden - use
// the desktop sub-package for ready-to-use implementations.
//
// Architecture Notes:
//   - Part of the API layer (public facade)
//   - Re-exports Domain types (Result, ErrorInfo, Person)
//   - Re-exports Application types (GreetPort, GreetCommand, Unit)
//   - Does NOT import Infrastructure (hidden implementation detail)
//   - Use api/adapter/desktop for ready-to-use greeter with console output
//
// Usage:
//
//	import "github.com/abitofhelp/hybrid_lib_go/api"
//	import "github.com/abitofhelp/hybrid_lib_go/api/adapter/desktop"
//
//	// Create a greeter with console output
//	greeter := desktop.NewGreeter()
//
//	// Execute greeting
//	result := greeter.Execute(ctx, api.NewGreetCommand("Alice"))
//	if result.IsOk() {
//	    // Success - greeting was written
//	}
package api

import (
	"github.com/abitofhelp/hybrid_lib_go/application/command"
	"github.com/abitofhelp/hybrid_lib_go/application/model"
	"github.com/abitofhelp/hybrid_lib_go/application/port/inbound"
	"github.com/abitofhelp/hybrid_lib_go/application/port/outbound"
	domerr "github.com/abitofhelp/hybrid_lib_go/domain/error"
	"github.com/abitofhelp/hybrid_lib_go/domain/valueobject"
)

// ============================================================================
// Domain Types (Re-exported)
// ============================================================================

// Result is a generic result type for operations that can fail.
// It follows the Railway-Oriented Programming pattern.
type Result[T any] = domerr.Result[T]

// ErrorType contains information about an error (kind + message).
type ErrorType = domerr.ErrorType

// ErrorKind represents the category of error.
type ErrorKind = domerr.ErrorKind

// Person is an immutable value object representing a person's name.
type Person = valueobject.Person

// Error kind constants
const (
	ValidationError     = domerr.ValidationError
	InfrastructureError = domerr.InfrastructureError
)

// Ok creates a successful Result containing the given value.
func Ok[T any](value T) Result[T] {
	return domerr.Ok(value)
}

// Err creates a failed Result containing the given error.
func Err[T any](err ErrorType) Result[T] {
	return domerr.Err[T](err)
}

// CreatePerson creates a new Person value object with validation.
func CreatePerson(name string) Result[Person] {
	return valueobject.CreatePerson(name)
}

// MaxNameLength is the maximum allowed length for a person's name.
const MaxNameLength = valueobject.MaxNameLength

// ============================================================================
// Application Types (Re-exported)
// ============================================================================

// Unit represents a void/unit type for operations that return no value.
type Unit = model.Unit

// GreetCommand is a command DTO for the greet use case.
type GreetCommand = command.GreetCommand

// NewGreetCommand creates a new GreetCommand with the given name.
func NewGreetCommand(name string) GreetCommand {
	return command.NewGreetCommand(name)
}

// GreetPort is the input port interface for the greet use case.
type GreetPort = inbound.GreetPort

// WriterPort is the output port interface for writing messages.
type WriterPort = outbound.WriterPort
