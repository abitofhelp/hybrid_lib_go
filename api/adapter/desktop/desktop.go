// SPDX-License-Identifier: BSD-3-Clause
// Copyright (c) 2025 Michael Gardner, A Bit of Help, Inc.
// Package: desktop
// Description: Desktop/console implementation with infrastructure wired

// Package desktop provides ready-to-use implementations for desktop applications.
//
// This package wires infrastructure (console output) with the application use cases,
// providing factory functions that return fully configured greeters.
//
// Architecture Notes:
//   - Platform-specific instantiation layer
//   - Wires Infrastructure adapters to Application use cases
//   - Provides factory functions for ready-to-use implementations
//   - Follows hexagonal architecture with DIP
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
//	    // Success - greeting was written to console
//	}
package desktop

import (
	"context"

	"github.com/abitofhelp/hybrid_lib_go/api"
	"github.com/abitofhelp/hybrid_lib_go/application/usecase"
	"github.com/abitofhelp/hybrid_lib_go/infrastructure/adapter"
)

// Greeter is a ready-to-use greeter with console output.
type Greeter struct {
	useCase *usecase.GreetUseCase[*adapter.ConsoleWriter]
}

// NewGreeter creates a new Greeter with console output.
// This is the recommended way to create a ready-to-use greeter for desktop apps.
func NewGreeter() *Greeter {
	writer := adapter.NewConsoleWriter()
	uc := usecase.NewGreetUseCase[*adapter.ConsoleWriter](writer)
	return &Greeter{useCase: uc}
}

// Execute performs the greet operation, writing output to console.
func (g *Greeter) Execute(ctx context.Context, cmd api.GreetCommand) api.Result[api.Unit] {
	return g.useCase.Execute(ctx, cmd)
}

// GreeterWithWriter creates a Greeter with a custom writer.
// Use this when you need to redirect output (e.g., to a buffer for testing).
func GreeterWithWriter[W api.WriterPort](writer W) *GreeterCustom[W] {
	uc := usecase.NewGreetUseCase[W](writer)
	return &GreeterCustom[W]{useCase: uc}
}

// GreeterCustom is a greeter with a custom writer type.
type GreeterCustom[W api.WriterPort] struct {
	useCase *usecase.GreetUseCase[W]
}

// Execute performs the greet operation with the custom writer.
func (g *GreeterCustom[W]) Execute(ctx context.Context, cmd api.GreetCommand) api.Result[api.Unit] {
	return g.useCase.Execute(ctx, cmd)
}
