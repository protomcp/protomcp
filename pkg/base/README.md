# protomcp/pkg/base

[![Go Reference][godoc-badge]][godoc-link]
[![codecov][codecov-badge]][codecov-link]
[![Go Report Card][goreport-badge]][goreport-link]

Package base provides shared base types and interfaces used across all protocols
in the ProtoMCP ecosystem.

## Overview

This package serves as the foundation for the unified proto code generation
design, ensuring consistency between NanoRPC and ProtoMCP implementations. It
provides core abstractions that enable protocol-agnostic service
implementations.

## Core Components

- **Message Interfaces**: Base message types with JSON v2 and binary marshalling
  support
- **Dispatcher Interface**: Protocol-agnostic message routing
- **Validation Types**: Detailed field-level validation errors
- **Error Handling**: Common error variables and patterns
- **Context Utilities**: Type-safe context value storage

## Usage

```go
import "github.com/protomcp/protomcp/pkg/base"
```

## Design Principles

- Interface-first design to avoid tight coupling
- Protocol independence through abstraction
- Modern Go patterns (generics, encoding/json/v2)
- Integration with darvaza.org/core for production-ready utilities

[godoc-badge]: https://pkg.go.dev/badge/github.com/protomcp/protomcp/pkg/base.svg
[godoc-link]: https://pkg.go.dev/github.com/protomcp/protomcp/pkg/base
[codecov-badge]: https://codecov.io/gh/protomcp/protomcp/graph/badge.svg?flag=base
[codecov-link]: https://codecov.io/gh/protomcp/protomcp?flag=base
[goreport-badge]: https://goreportcard.com/badge/github.com/protomcp/protomcp
[goreport-link]: https://goreportcard.com/report/github.com/protomcp/protomcp
