// SPDX-License-Identifier: BSD-3-Clause
// Copyright (c) 2025 Michael Gardner, A Bit of Help, Inc.

//go:build integration

// Package integration provides API integration tests.
//
// Integration tests verify the complete library flow through the public API.
// They test:
// - API facade usage
// - Factory functions
// - Error handling
// - Result monad operations
//
// Run with: go test -v -tags=integration ./test/integration/...
package integration

import (
	"bytes"
	"context"
	"testing"

	"github.com/abitofhelp/hybrid_lib_go/api"
	"github.com/abitofhelp/hybrid_lib_go/api/adapter/desktop"
	"github.com/abitofhelp/hybrid_lib_go/application/model"
	domerr "github.com/abitofhelp/hybrid_lib_go/domain/error"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)


// ============================================================================
// Mock Writer for Testing
// ============================================================================

// MockWriter captures output for testing.
type MockWriter struct {
	Buffer bytes.Buffer
}

// Write implements outbound.WriterPort.
func (w *MockWriter) Write(ctx context.Context, msg string) domerr.Result[model.Unit] {
	w.Buffer.WriteString(msg)
	return domerr.Ok(model.Unit{})
}

// String returns the captured output.
func (w *MockWriter) String() string {
	return w.Buffer.String()
}

// ============================================================================
// API Facade Tests
// ============================================================================

func TestGreeter_Execute_Success(t *testing.T) {
	// Arrange: Create greeter with mock writer
	writer := &MockWriter{}
	greeter := desktop.GreeterWithWriter[*MockWriter](writer)
	cmd := api.NewGreetCommand("Alice")
	ctx := context.Background()

	// Act
	result := greeter.Execute(ctx, cmd)

	// Assert
	require.True(t, result.IsOk(), "Expected successful result")
	assert.Contains(t, writer.String(), "Hello, Alice!")
}

func TestGreeter_Execute_EmptyName_ReturnsValidationError(t *testing.T) {
	// Arrange
	writer := &MockWriter{}
	greeter := desktop.GreeterWithWriter[*MockWriter](writer)
	cmd := api.NewGreetCommand("")
	ctx := context.Background()

	// Act
	result := greeter.Execute(ctx, cmd)

	// Assert
	require.True(t, result.IsError(), "Expected validation error")
	errInfo := result.ErrorInfo()
	assert.Equal(t, api.ValidationError, errInfo.Kind)
}

func TestGreeter_Execute_LongName_ReturnsValidationError(t *testing.T) {
	// Arrange
	writer := &MockWriter{}
	greeter := desktop.GreeterWithWriter[*MockWriter](writer)
	longName := make([]byte, api.MaxNameLength+1)
	for i := range longName {
		longName[i] = 'a'
	}
	cmd := api.NewGreetCommand(string(longName))
	ctx := context.Background()

	// Act
	result := greeter.Execute(ctx, cmd)

	// Assert
	require.True(t, result.IsError(), "Expected validation error for long name")
	errInfo := result.ErrorInfo()
	assert.Equal(t, api.ValidationError, errInfo.Kind)
}

// ============================================================================
// Domain Type Re-export Tests
// ============================================================================

func TestAPI_CreatePerson_Success(t *testing.T) {
	// Act
	result := api.CreatePerson("Bob")

	// Assert
	require.True(t, result.IsOk())
	person := result.Value()
	assert.Equal(t, "Bob", person.GetName())
}

func TestAPI_CreatePerson_EmptyName_ReturnsError(t *testing.T) {
	// Act
	result := api.CreatePerson("")

	// Assert
	require.True(t, result.IsError())
	assert.Equal(t, api.ValidationError, result.ErrorInfo().Kind)
}

func TestAPI_Result_Ok(t *testing.T) {
	// Act
	result := api.Ok("success")

	// Assert
	require.True(t, result.IsOk())
	assert.Equal(t, "success", result.Value())
}

func TestAPI_Result_Err(t *testing.T) {
	// Arrange
	err := api.ErrorType{
		Kind:    api.ValidationError,
		Message: "test error",
	}

	// Act
	result := api.Err[string](err)

	// Assert
	require.True(t, result.IsError())
	assert.Equal(t, "test error", result.ErrorInfo().Message)
}

// ============================================================================
// Console Greeter Tests (writes to stdout)
// ============================================================================

func TestDesktop_NewGreeter_Success(t *testing.T) {
	// This test verifies the default console greeter works.
	// Note: Output goes to stdout, can't easily capture here.
	greeter := desktop.NewGreeter()
	require.NotNil(t, greeter, "NewGreeter should return non-nil")

	// Execute - will print to console
	cmd := api.NewGreetCommand("IntegrationTest")
	ctx := context.Background()
	result := greeter.Execute(ctx, cmd)

	assert.True(t, result.IsOk(), "Console greeter should succeed")
}

// ============================================================================
// Multiple Names Test
// ============================================================================

func TestGreeter_Execute_MultipleNames(t *testing.T) {
	names := []string{"Alice", "Bob", "Charlie", "世界", "مرحبا"}

	for _, name := range names {
		t.Run(name, func(t *testing.T) {
			writer := &MockWriter{}
			greeter := desktop.GreeterWithWriter[*MockWriter](writer)
			cmd := api.NewGreetCommand(name)
			ctx := context.Background()

			result := greeter.Execute(ctx, cmd)

			require.True(t, result.IsOk(), "Should succeed for name: %s", name)
			assert.Contains(t, writer.String(), name)
		})
	}
}
