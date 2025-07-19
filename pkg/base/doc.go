// Package base provides shared base types and interfaces used across all
// protocols in the ProtoMCP ecosystem. This package serves as the foundation
// for the unified proto code generation design, ensuring consistency between
// NanoRPC and ProtoMCP implementations.
//
// # Core Components
//
//   - Message interfaces with JSON v2 and binary marshalling support
//   - Dispatcher interface for protocol-agnostic message routing
//   - Validation error types with detailed field-level information
//   - Common error variables and handling patterns
//   - Context utilities for type-safe value storage
//
// # Design Principles
//
// The base package follows these key principles:
//
//   - Interface-first design to avoid tight coupling
//   - Protocol independence through abstraction
//   - Modern Go patterns (generics, encoding/json/v2)
//   - Integration with darvaza.org/core for production-ready utilities
//
// All protocol-specific implementations (JSON-RPC, MCP, REST) build upon these
// base types, ensuring a consistent API across different protocols while
// allowing for protocol-specific optimisations.
package base
