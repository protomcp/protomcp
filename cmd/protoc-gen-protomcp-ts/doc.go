// Package main implements protoc-gen-protomcp-ts, a protocol buffer compiler
// plugin that generates TypeScript client code for services defined using
// Protocol Buffers. The generated clients support JSON-RPC 2.0, MCP (Model
// Context Protocol), and REST protocols.
//
// # Installation
//
//	go install protomcp.org/protomcp/cmd/protoc-gen-protomcp-ts@latest
//
// # Usage
//
//	protoc --protomcp-ts_out=. --protomcp-ts_opt=paths=source_relative service.proto
//
// # Generated Code
//
// For each service definition, the plugin generates:
//
//   - TypeScript interfaces for all messages with full type safety
//   - Enum definitions with string literal types
//   - Discriminated unions for protobuf `oneof` fields
//   - Type guards for runtime type checking
//   - JSON-RPC client with async/await support
//   - MCP client for AI assistant integration
//   - REST client using fetch API
//   - Validation functions from buf.validate rules
//   - Mock clients for testing
//
// # Type Mapping
//
// Protocol Buffer types map to TypeScript as follows:
//
//   - Scalar types: number, string, boolean, Uint8Array
//   - Enums: String literal union types
//   - Messages: TypeScript interfaces
//   - Repeated: Arrays
//   - Maps: Record<K, V> or index signatures
//   - `Oneof`: Discriminated unions
//   - Timestamps: Date | string
//   - Any: unknown with runtime checks
//
// # Options
//
// The plugin supports various options through --protomcp-ts_opt:
//
//   - paths=source_relative: Use source-relative import paths
//   - emit_unpopulated_fields=true: Include undefined fields
//   - use_proto_names=true: Use proto field names (not camelCase)
//   - esm=true: Generate ES modules (default: CommonJS)
//
// # Client Features
//
// Generated clients include:
//
//   - Automatic request/response validation
//   - Configurable transports and interceptors
//   - Retry logic with exponential backoff
//   - Request batching for JSON-RPC
//   - TypeScript 5.0+ features
//
// The generated code has zero runtime dependencies by default, with optional
// integrations for popular libraries.
package main
