// Package testutils provides common testing utilities for protomcp tests.
//
// This package contains helper functions and factories for creating protocol
// buffer definitions, running generators, and asserting on generated output.
// It simplifies the process of writing comprehensive tests for protocol buffer
// code generators.
//
// # Factory Functions
//
// The package provides factory functions for creating protocol buffer structures:
//
//   - NewFileDescriptor: Creates file descriptors with common defaults
//   - NewMessage: Creates message types with fields
//   - NewField: Creates field descriptors with proper types
//   - NewService: Creates service definitions
//   - NewMethod: Creates RPC method definitions
//   - NewEnum: Creates enum types with values
//   - NewCodeGenRequest: Creates code generator requests
//
// # Assertion Helpers
//
// Common assertion functions for testing generated code:
//
//   - AssertContains: Verifies generated content contains expected strings
//   - AssertFileCount: Validates the number of generated files
//   - AssertEqual: Generic equality assertion with clear error messages
//   - AssertSliceEqual: Compares string slices
//   - AssertSliceOfSlicesEqual: Compares nested string slices
//
// # Running Generators
//
// The RunGenerator function provides a convenient way to execute a generator
// function with a prepared request and capture its response for testing.
//
// # Example Usage
//
//	func TestServiceGeneration(t *testing.T) {
//	    // Create a proto file with a service
//	    file := NewFileDescriptor("api.proto", "api.v1", "github.com/example/api/v1")
//
//	    // Add a message type
//	    userMsg := NewMessage("User",
//	        NewField("id", 1, descriptorpb.FieldDescriptorProto_TYPE_STRING),
//	        NewField("name", 2, descriptorpb.FieldDescriptorProto_TYPE_STRING),
//	    )
//	    file.MessageType = append(file.MessageType, userMsg)
//
//	    // Add a service with methods
//	    service := NewService("UserService",
//	        NewMethod("GetUser", ".api.v1.GetUserRequest", ".api.v1.User"),
//	    )
//	    file.Service = append(file.Service, service)
//
//	    // Run the generator
//	    response := RunGenerator(t, NewCodeGenRequest(file), Generate)
//
//	    // Assert on the output
//	    AssertFileCount(t, response, 1)
//	    AssertContains(t, response.File[0].GetContent(), "type UserService interface")
//	}
package testutils
