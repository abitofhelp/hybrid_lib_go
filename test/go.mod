// SPDX-License-Identifier: BSD-3-Clause
// Copyright (c) 2025 Michael Gardner, A Bit of Help, Inc.

module github.com/abitofhelp/hybrid_lib_go/test

go 1.23.0

// Test module - CAN have external dependencies for integration tests
// This is separate from /src modules which must have ZERO external module dependencies

require (
	github.com/abitofhelp/hybrid_lib_go/api v0.0.0
	github.com/abitofhelp/hybrid_lib_go/api/adapter/desktop v0.0.0
	github.com/abitofhelp/hybrid_lib_go/application v0.0.0
	github.com/abitofhelp/hybrid_lib_go/domain v0.0.0
	github.com/stretchr/testify v1.11.1
)

require (
	github.com/abitofhelp/hybrid_lib_go/infrastructure v0.0.0 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

replace github.com/abitofhelp/hybrid_lib_go/api => ../api

replace github.com/abitofhelp/hybrid_lib_go/api/adapter/desktop => ../api/adapter/desktop

replace github.com/abitofhelp/hybrid_lib_go/domain => ../domain

replace github.com/abitofhelp/hybrid_lib_go/application => ../application

replace github.com/abitofhelp/hybrid_lib_go/infrastructure => ../infrastructure
