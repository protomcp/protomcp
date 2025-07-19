# protomcp

> [!WARNING]
> This project is still in the design phase and not ready for production use.
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

## Key Features

- **Dual Protocol Support**: Generate unified endpoints for `JSON-RPC` 2.0
  and `MCP`.
- **Modern Transports**: `HTTP/2` and `QUIC` protocol support.
- **Interface-First Design**: Prioritizes interfaces over concrete structs
  for maximum modularity.
- **Schema Validation**: Integrated JSON Schema generation and validation.
- **Protobuf Integration**: Works with existing `.proto` service definitions.
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
├── cmd/                    # CLI tools and protoc plugins
├── pkg/protomcp/          # Core library code
├── internal/build/        # Build system and tooling
└── examples/              # Usage examples and demos
```

## Development Status

### 🚧 DESIGN PHASE - NOT PRODUCTION READY 🚧

This project is currently in active design and early development. Nothing is
stable yet and everything may change. The initial focus is on:

1. Core generator framework.
2. Interface design patterns.
3. JSON-RPC 2.0 foundation.
4. MCP protocol integration.
5. Transport layer abstraction.

**Do not use this in production environments.**

## Contributing

See [AGENT.md](AGENT.md) for development guidelines and build system details.

## License

This project is licensed under the MIT License—see the
[LICENCE.txt](LICENCE.txt) file for details.
