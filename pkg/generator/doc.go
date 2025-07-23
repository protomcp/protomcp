// Package generator provides utilities for ProtoMCP code generation,
// including debugging and tracing facilities.
//
// This package serves as a foundation for building protocol buffer code
// generators that produce Go code for JSON-RPC 2.0 and MCP (Model Context
// Protocol) endpoints. It provides common utilities, test helpers, and
// infrastructure needed by the protoc-gen-protomcp plugin.
//
// # Overview
//
// The generator package includes:
//
//   - Test utilities for validating generated code
//   - Protocol buffer factory functions for testing
//   - Assertion helpers for common test scenarios
//   - Infrastructure for future tracing and debugging capabilities
//
// # Sub-packages
//
//   - testutils: Testing utilities and factories for protobuf structures
//
// # Usage
//
// This package is primarily used by the protoc-gen-protomcp plugin and
// its tests. It provides the foundation for:
//
//   - Creating test protocol buffer definitions
//   - Validating generated code output
//   - Testing code generation logic
//
// # Example
//
// Using the testutils package to create test protocol buffers:
//
//	import "protomcp.org/protomcp/pkg/generator/testutils"
//
//	func TestGenerator(t *testing.T) {
//	    // Create a test service definition
//	    file := testutils.NewFileDescriptor("test.proto", "test.v1", "test/v1")
//	    service := testutils.NewService("TestService",
//	        testutils.NewMethod("GetUser", ".test.v1.GetUserRequest", ".test.v1.GetUserResponse"),
//	    )
//	    file.Service = append(file.Service, service)
//
//	    // Run the generator
//	    response := testutils.RunGenerator(t, testutils.NewCodeGenRequest(file), myGenerator)
//
//	    // Validate output
//	    testutils.AssertFileCount(t, response, 1)
//	    testutils.AssertContains(t, response.File[0].GetContent(), "type TestService interface")
//	}
package generator
