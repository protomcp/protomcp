# protomcp/pkg/generator

[![Go Reference][godoc-badge]][godoc-link]
[![codecov][codecov-badge]][codecov-link]
[![Go Report Card][goreport-badge]][goreport-link]

Package generator provides utilities for ProtoMCP code generation, including
debugging, tracing facilities, and testing utilities.

## Overview

This package serves as the foundation for building protocol buffer code
generators that produce Go code for JSON-RPC 2.0 and MCP (Model Context
Protocol) endpoints. It provides common utilities, test helpers, and
infrastructure needed by the protoc-gen-protomcp plugin.

## Sub-packages

### testutils

The `testutils` package provides comprehensive testing utilities:

- **Factory Functions**: Create protocol buffer structures for testing
- **Assertion Helpers**: Validate generated code output
- **Generator Testing**: Run generators with prepared requests

## Usage

### Testing Code Generators

```go
import "protomcp.org/protomcp/pkg/generator/testutils"

func TestMyGenerator(t *testing.T) {
    // Create a test proto file
    file := testutils.NewFileDescriptor(
        "api.proto", "api.v1", "github.com/example/api/v1")

    // Add a service
    service := testutils.NewService("UserService",
        testutils.NewMethod(
            "GetUser", ".api.v1.GetUserRequest", ".api.v1.User"),
    )
    file.Service = append(file.Service, service)

    // Run the generator
    req := testutils.NewCodeGenRequest(file)
    response := testutils.RunGenerator(t, req, MyGenerator)

    // Assert on output
    testutils.AssertFileCount(t, response, 1)
    content := response.File[0].GetContent()
    testutils.AssertContains(t, content, "type UserService interface")
}
```

## Features

- Protocol buffer factory functions for easy test setup
- Comprehensive assertion helpers for testing generated code
- Support for creating complex proto definitions programmatically
- Integration with protogen for generator development
- LazyBuffer utility for efficient string building

## LazyBuffer Utility

The package provides a LazyBuffer type that wraps strings.Builder but doesn't
return errors - hence "lazy". No more `_, _ =` error handling:

```go
import "protomcp.org/protomcp/pkg/generator"

// Create a buffer
var buf generator.LazyBuffer

// Append strings without error handling (empty strings are ignored)
buf.WriteString("Hello", " ", "World")  // No _, _ = needed!

// Append runes (including Unicode)
buf.WriteRunes(' ', '-', ' ', '🚀')

// Printf formatting without error handling
buf.Printf(" count: %d", 42)  // No _, _ = needed!

// Get the result
result := buf.String() // "Hello World - 🚀 count: 42"

// Method chaining for fluent API
message := buf.WriteString("Status: ").
    Printf("%d%%", 100).
    WriteRunes(' ', '✓').
    String() // "Status: 100% ✓"

// Compare with strings.Builder:
var sb strings.Builder
_, _ = sb.WriteString("Hello")  // Must handle error return
_, _ = sb.WriteRune(' ')         // Must handle error return
_, _ = fmt.Fprintf(&sb, "%d", 42)  // Must handle error return

// Bonus: nil-safe operations
var nilBuf *generator.LazyBuffer
nilBuf.WriteString("test") // Won't panic
nilBuf.WriteRunes('a', 'b') // Won't panic
str := nilBuf.String()      // Returns empty string
```

LazyBuffer is particularly useful in code generators where you frequently build
strings and don't want to clutter your code with error handling that will never
actually error.

[godoc-badge]: https://pkg.go.dev/badge/protomcp.org/protomcp/pkg/generator.svg
[godoc-link]: https://pkg.go.dev/protomcp.org/protomcp/pkg/generator
[codecov-badge]: https://codecov.io/gh/protomcp/protomcp/graph/badge.svg?flag=generator
[codecov-link]: https://codecov.io/gh/protomcp/protomcp?flag=generator
[goreport-badge]: https://goreportcard.com/badge/protomcp.org/protomcp
[goreport-link]: https://goreportcard.com/report/protomcp.org/protomcp
