// Package main implements protoc-gen-protomcp, a protocol buffer compiler
// plugin that generates Go code for serving Protocol Buffer services through
// JSON-RPC 2.0, MCP (Model Context Protocol), and REST endpoints.
//
// # Installation
//
//	go install github.com/protomcp/protomcp/cmd/protoc-gen-protomcp@latest
//
// # Usage
//
//	protoc --protomcp_out=. --protomcp_opt=paths=source_relative service.proto
//
// # Generated Code
//
// For each service definition, the plugin generates:
//
//   - Service interfaces with protocol-agnostic methods
//   - JSON-RPC 2.0 dispatcher with method routing
//   - MCP dispatcher with tool/resource/prompt handlers
//   - REST dispatcher using google.api.http annotations
//   - Message types with JSON v2 and binary marshalling
//   - JSON Schema validators from buf.validate rules
//   - HTTP/2 and QUIC transport servers
//   - Mock implementations for testing
//
// # Options
//
// The plugin supports various options through --protomcp_opt:
//
//   - paths=source_relative: Use source-relative import paths
//   - module=<path>: Override module path detection
//   - logging=true: Enable debug logging during generation
//
// # Proto Annotations
//
// The plugin recognises:
//
//   - google.api.http: REST endpoint configuration
//   - buf.validate: Validation rules for messages
//   - protomcp.jsonrpc: JSON-RPC method options
//   - protomcp.mcp: MCP tool/resource definitions
//
// The generated code integrates with pkg/protomcp for runtime support.
package main
