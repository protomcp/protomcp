// Package protomcp is the root package for the ProtoMCP project, a protoc
// generator framework that creates unified JSON-RPC 2.0 and MCP (Model Context
// Protocol) endpoints from Protocol Buffer service definitions.
//
// ProtoMCP generates Go code that serves both protocols over HTTP/2 and QUIC
// transports, emphasising interface-based design for modularity. The generated
// code prioritises Go interfaces over concrete protobuf types, enabling
// protocol-agnostic service implementations.
//
// # Key Features
//
//   - Unified protocol support: JSON-RPC 2.0 and MCP from single proto definitions
//   - Interface-first architecture for maximum flexibility
//   - Modern transport support: HTTP/2 and QUIC
//   - Comprehensive validation with JSON Schema
//   - TypeScript client generation
//   - Integration with darvaza.org utilities for production readiness
//
// # Project Structure
//
//   - cmd/protoc-gen-protomcp: Protocol buffer compiler plugin for Go code generation
//   - cmd/protoc-gen-protomcp-ts: Protocol buffer compiler plugin for TypeScript generation
//   - pkg/protomcp: Core library with protocol implementations and utilities
//
// This is the parent module that coordinates the build and testing of all
// submodules in the project.
package protomcp
