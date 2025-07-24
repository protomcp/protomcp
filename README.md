# protomcp

[![Go Reference][godoc-badge]][godoc-link]
[![codecov][codecov-badge]][codecov-link]
[![Go Report Card][goreport-badge]][goreport-link]

> [!WARNING]
> This project is in early development and not ready for production use.
> APIs and architecture may change significantly.

A modular `protoc` generator framework for creating combined `JSON-RPC` 2.0 and
`MCP` (Model Context Protocol) endpoints from `.proto` service definitions,
supporting `HTTP/2` and `QUIC` transport protocols.

## Overview

`protomcp` generates Go code that provides unified interfaces for services
defined in Protocol Buffer (`.proto`) files, enabling them to serve both
`JSON-RPC` 2.0 and Anthropic's Model Context Protocol (MCP) over modern
transport protocols. The generator prioritizes interface-based design for
modularity and loose coupling.

## Quick Example

Given a protobuf service definition:

```protobuf
service Calculator {
    rpc Add(AddRequest) returns (AddResponse);
}

message AddRequest {
    double a = 1;
    double b = 2;
}

enum Status {
    STATUS_UNSPECIFIED = 0;
    STATUS_ACTIVE = 1;
    STATUS_INACTIVE = 2;
}
```

ProtoMCP generates Go interfaces and enum helpers:

```go
// Generated enum with %Enum pattern (avoids clash with .pb.go)
type StatusEnum int32

const (
    StatusEnum_UNSPECIFIED StatusEnum = 0
    StatusEnum_ACTIVE      StatusEnum = 1
    StatusEnum_INACTIVE    StatusEnum = 2
)

// Helper methods for enum
func (x StatusEnum) String() string         // Returns "UNSPECIFIED", etc.
func (x StatusEnum) IsValid() bool          // Validates enum value
func (x StatusEnum) MarshalText() ([]byte, error)
func (x *StatusEnum) UnmarshalText(text []byte) error  // Decode from text

// Generated interfaces
type IAddRequest interface {
    GetA() float64
    SetA(v float64) error
    GetB() float64
    SetB(v float64) error
}

type ICalculatorService interface {
    Add(ctx context.Context, req IAddRequest) (IAddResponse, error)
}
```

## Key Features

### Currently Implemented

- **Interface Generation**: Template-based generation of Go interfaces from
  protobuf messages.
- **Service Interfaces**: Context-aware service method signatures.
- **Enum Helpers**: Generate type-safe enum types with String(), IsValid(),
  and text marshaling methods.
- **Configurable Naming**: Support for custom interface naming patterns (e.g.,
  `I%`, `%Interface`) and enum patterns (e.g., `%Enum`, `E%`).
- **Proto3 Support**: Full support for all field types including optional
  fields.
- **Test Infrastructure**: Comprehensive test utilities and factory functions.

### Planned Features

- **Dual Protocol Support**: Generate unified endpoints for `JSON-RPC` 2.0
  and `MCP`.
- **Modern Transports**: `HTTP/2` and `QUIC` protocol support.
- **Schema Validation**: Integrated JSON Schema generation and validation.
- **TypeScript Generation**: Client SDK generation with full type safety.
- **Sourcegraph jsonrpc2**: Built on proven JSON-RPC 2.0 foundation.

## Architecture

The project follows a modular architecture with clear separation of concerns:

- **Generator Package**: Core `protoc` plugin for code generation.
- **Transport Layer**: `HTTP/2` and `QUIC` transport implementations.
- **Protocol Layer**: `JSON-RPC` 2.0 and `MCP` protocol handlers.
- **Validation Layer**: JSON Schema validation and type safety.
- **Interface Layer**: Service interfaces independent of concrete protobuf
  types.

## Dependencies

- **Go 1.23+**: Modern Go version for latest features.
- **sourcegraph/jsonrpc2**: Proven JSON-RPC 2.0 implementation.
- **Protocol Buffers**: For service definition parsing.
- **JSON Schema**: For request/response validation.
- **HTTP/2 & QUIC**: Modern transport protocol support.

## Project Structure

```text
protomcp/
├── cmd/
│   └── protoc-gen-protomcp/    # Template-based protoc plugin
├── internal/
│   ├── gen-go/                 # Go code generation helpers
│   ├── testutils/              # Test factory functions
│   └── build/                  # Build system and tooling
├── pkg/base/                   # Base types (currently minimal)
├── proto/examples/             # Example proto files
└── docs/                       # Documentation
```

## Development Status

### 🚧 EARLY DEVELOPMENT - NOT PRODUCTION READY 🚧

This project has a working protoc plugin that generates Go interfaces from
protobuf definitions:

**Completed:**

- ✅ Template-based code generator (`protoc-gen-protomcp`)
- ✅ Interface generation for messages with getters/setters
- ✅ Service interface generation with context support
- ✅ Enum helper generation with validation and text marshaling
- ✅ Configurable interface and enum naming patterns
- ✅ Comprehensive test infrastructure

**In Progress:**

- 🚧 Validation helper generation (buf.validate integration)
- 🚧 Protocol dispatcher generation (JSON-RPC 2.0, MCP)
- 🚧 Message implementation generation

**Planned:**

- ⏳ TypeScript client generation
- ⏳ Transport layer implementation (HTTP/2, QUIC)
- ⏳ Production features (metrics, logging, auth)

**Do not use this in production environments.**

## Contributing

See [AGENT.md](AGENT.md) for development guidelines and build system details.

## License

This project is licensed under the MIT License—see the
[LICENCE.txt](LICENCE.txt) file for details.

[godoc-badge]: https://pkg.go.dev/badge/protomcp.org/protomcp.svg
[godoc-link]: https://pkg.go.dev/protomcp.org/protomcp
[codecov-badge]: https://codecov.io/gh/protomcp/protomcp/graph/badge.svg?flag=root
[codecov-link]: https://codecov.io/gh/protomcp/protomcp?flag=root
[goreport-badge]: https://goreportcard.com/badge/protomcp.org/protomcp
[goreport-link]: https://goreportcard.com/report/protomcp.org/protomcp
