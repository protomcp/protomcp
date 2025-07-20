package main

import (
	"testing"

	"google.golang.org/protobuf/types/descriptorpb"

	"protomcp.org/protomcp/pkg/generator/testutils"
)

// createProtoFileWithService creates a proto file with a service and required message types
func createProtoFileWithService() *descriptorpb.FileDescriptorProto {
	protoFile := testutils.NewFileDescriptor("test.proto", "test", "github.com/example/test")
	// Add the message types that the service methods reference
	protoFile.MessageType = []*descriptorpb.DescriptorProto{
		testutils.NewMessage("GetItemRequest",
			testutils.NewField("id", 1, descriptorpb.FieldDescriptorProto_TYPE_STRING),
		),
		testutils.NewMessage("GetItemResponse",
			testutils.NewField("item", 1, descriptorpb.FieldDescriptorProto_TYPE_STRING),
		),
	}
	protoFile.Service = []*descriptorpb.ServiceDescriptorProto{
		testutils.NewService("TestService",
			testutils.NewMethod("GetItem", "GetItemRequest", "GetItemResponse"),
		),
	}
	return protoFile
}

// TestNeedsNoImpl tests the NeedsNoImpl method
func TestNeedsNoImpl(t *testing.T) {
	t.Run("disabled when GenerateNoImpl is false", testNeedsNoImplDisabled)
	t.Run("disabled when GenerateInterfaces is false", testNeedsNoImplNoInterfaces)
	t.Run("enabled with messages", testNeedsNoImplWithMessages)
	t.Run("enabled with services", testNeedsNoImplWithServices)
	t.Run("disabled when GenerateServices is false", testNeedsNoImplNoServices)
	t.Run("disabled with empty file", testNeedsNoImplEmptyFile)
}

func testNeedsNoImplDisabled(t *testing.T) {
	protoFile := testutils.NewFileDescriptor("test.proto", "test", "github.com/example/test")
	protoFile.MessageType = []*descriptorpb.DescriptorProto{
		testutils.NewMessage("TestMessage",
			testutils.NewField("id", 1, descriptorpb.FieldDescriptorProto_TYPE_STRING),
		),
	}

	plugin, err := testutils.NewPlugin(t, protoFile)
	if err != nil {
		t.Fatalf("failed to create plugin: %v", err)
	}
	gen := NewGenerator(plugin)
	file := plugin.Files[0]

	opts := &GeneratorOptions{
		GenerateInterfaces: true,
		GenerateNoImpl:     false,
	}

	if gen.NeedsNoImpl(file, opts) {
		t.Error("NeedsNoImpl should return false when GenerateNoImpl is false")
	}
}

func testNeedsNoImplNoInterfaces(t *testing.T) {
	protoFile := testutils.NewFileDescriptor("test.proto", "test", "github.com/example/test")
	protoFile.MessageType = []*descriptorpb.DescriptorProto{
		testutils.NewMessage("TestMessage",
			testutils.NewField("id", 1, descriptorpb.FieldDescriptorProto_TYPE_STRING),
		),
	}

	plugin, err := testutils.NewPlugin(t, protoFile)
	if err != nil {
		t.Fatalf("failed to create plugin: %v", err)
	}
	gen := NewGenerator(plugin)
	file := plugin.Files[0]

	opts := &GeneratorOptions{
		GenerateInterfaces: false,
		GenerateNoImpl:     true,
	}

	if gen.NeedsNoImpl(file, opts) {
		t.Error("NeedsNoImpl should return false when GenerateInterfaces is false")
	}
}

func testNeedsNoImplWithMessages(t *testing.T) {
	protoFile := testutils.NewFileDescriptor("test.proto", "test", "github.com/example/test")
	protoFile.MessageType = []*descriptorpb.DescriptorProto{
		testutils.NewMessage("TestMessage",
			testutils.NewField("id", 1, descriptorpb.FieldDescriptorProto_TYPE_STRING),
		),
	}

	plugin, err := testutils.NewPlugin(t, protoFile)
	if err != nil {
		t.Fatalf("failed to create plugin: %v", err)
	}
	gen := NewGenerator(plugin)
	file := plugin.Files[0]

	opts := &GeneratorOptions{
		GenerateInterfaces: true,
		GenerateNoImpl:     true,
	}

	if !gen.NeedsNoImpl(file, opts) {
		t.Error("NeedsNoImpl should return true with messages and both options enabled")
	}
}

func testNeedsNoImplWithServices(t *testing.T) {
	protoFile := createProtoFileWithService()

	plugin, err := testutils.NewPlugin(t, protoFile)
	if err != nil {
		t.Fatalf("failed to create plugin: %v", err)
	}
	gen := NewGenerator(plugin)
	file := plugin.Files[0]

	opts := &GeneratorOptions{
		GenerateServices: true,
		GenerateNoImpl:   true,
	}

	if !gen.NeedsNoImpl(file, opts) {
		t.Error("NeedsNoImpl should return true with services and both options enabled")
	}
}

func testNeedsNoImplNoServices(t *testing.T) {
	protoFile := createProtoFileWithService()

	plugin, err := testutils.NewPlugin(t, protoFile)
	if err != nil {
		t.Fatalf("failed to create plugin: %v", err)
	}
	gen := NewGenerator(plugin)
	file := plugin.Files[0]

	opts := &GeneratorOptions{
		GenerateServices: false,
		GenerateNoImpl:   true,
	}

	if gen.NeedsNoImpl(file, opts) {
		t.Error("NeedsNoImpl should return false when GenerateServices is false")
	}
}

func testNeedsNoImplEmptyFile(t *testing.T) {
	protoFile := testutils.NewFileDescriptor("test.proto", "test", "github.com/example/test")
	// No messages or services

	plugin, err := testutils.NewPlugin(t, protoFile)
	if err != nil {
		t.Fatalf("failed to create plugin: %v", err)
	}
	gen := NewGenerator(plugin)
	file := plugin.Files[0]

	opts := &GeneratorOptions{
		GenerateInterfaces: true,
		GenerateServices:   true,
		GenerateNoImpl:     true,
	}

	if gen.NeedsNoImpl(file, opts) {
		t.Error("NeedsNoImpl should return false for empty file")
	}
}
