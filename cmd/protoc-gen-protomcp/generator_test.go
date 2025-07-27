package main

import (
	"fmt"
	"testing"

	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/types/descriptorpb"
	"google.golang.org/protobuf/types/pluginpb"

	"protomcp.org/protomcp/pkg/generator/testutils"
)

// runGenerator runs the generator with the given request and options
func runGenerator(
	t *testing.T,
	req *pluginpb.CodeGeneratorRequest,
	opts *GeneratorOptions,
) *pluginpb.CodeGeneratorResponse {
	return testutils.RunGenerator(t, req, func(plugin *protogen.Plugin) error {
		return generateFiles(plugin, opts)
	})
}

func generateFiles(plugin *protogen.Plugin, opts *GeneratorOptions) error {
	gen := NewGenerator(plugin)
	for _, file := range plugin.Files {
		if !file.Generate {
			continue
		}
		if err := gen.GenerateFile(file, opts); err != nil {
			return err
		}
	}
	return nil
}

// TestGenerateSimpleMessage tests generating an interface for a simple message
func TestGenerateSimpleMessage(t *testing.T) {
	// Create a simple proto file with one message
	protoFile := testutils.NewFileDescriptor("test.proto", "test", "github.com/example/test")
	protoFile.MessageType = []*descriptorpb.DescriptorProto{
		testutils.NewMessage("SimpleMessage",
			testutils.NewField("id", 1, descriptorpb.FieldDescriptorProto_TYPE_STRING),
			testutils.NewField("count", 2, descriptorpb.FieldDescriptorProto_TYPE_INT32),
		),
	}

	// Create and run generator
	req := testutils.NewCodeGenRequest(protoFile)
	opts := &GeneratorOptions{
		GenerateInterfaces: true,
	}
	resp := runGenerator(t, req, opts)

	// Check generated files
	testutils.AssertFileCount(t, resp, 1)

	genFile := resp.File[0]
	expectedName := "github.com/example/test/test.types.go"
	if genFile.GetName() != expectedName {
		t.Errorf("generated file name = %q, want %q", genFile.GetName(), expectedName)
	}

	content := genFile.GetContent()

	// Check expected content
	expectedStrings := []string{
		"type ISimpleMessage interface {", // Interface declaration
		"GetId() string",                  // Getter for id
		"GetCount() int32",                // Getter for count
		"SetId(v string) error",           // Setter for id
		"SetCount(v int32) error",         // Setter for count
	}

	for _, expected := range expectedStrings {
		testutils.AssertContains(t, content, expected)
	}
}

// TestGenerateWithoutInterfaces tests that no interfaces are generated when disabled
func TestGenerateWithoutInterfaces(t *testing.T) {
	// Create a simple proto file
	protoFile := testutils.NewFileDescriptor("test.proto", "test", "github.com/example/test")
	protoFile.MessageType = []*descriptorpb.DescriptorProto{
		testutils.NewMessage("SimpleMessage"),
	}

	// Create and run generator
	req := testutils.NewCodeGenRequest(protoFile)
	opts := &GeneratorOptions{
		GenerateInterfaces: false, // Disabled
	}
	resp := runGenerator(t, req, opts)

	// Check that no files were generated
	testutils.AssertFileCount(t, resp, 0)
}

// fieldTypeTestCase represents a test case for field type generation
type fieldTypeTestCase struct {
	name       string
	wantGetter string
	wantSetter string
	fieldType  descriptorpb.FieldDescriptorProto_Type
}

// newFieldTypeTestCase creates a new field type test case
func newFieldTypeTestCase(
	name string,
	fieldType descriptorpb.FieldDescriptorProto_Type,
	goType string,
) fieldTypeTestCase {
	return fieldTypeTestCase{
		name:       name,
		fieldType:  fieldType,
		wantGetter: fmt.Sprintf("GetValue() %s", goType),
		wantSetter: fmt.Sprintf("SetValue(v %s) error", goType),
	}
}

// test runs the test case
func (tc fieldTypeTestCase) test(t *testing.T) {
	protoFile := testutils.NewFileDescriptor("test.proto", "test", "github.com/example/test")
	protoFile.MessageType = []*descriptorpb.DescriptorProto{
		testutils.NewMessage("TestMessage",
			testutils.NewField("value", 1, tc.fieldType),
		),
	}

	req := testutils.NewCodeGenRequest(protoFile)
	opts := &GeneratorOptions{
		GenerateInterfaces: true,
	}
	resp := runGenerator(t, req, opts)

	testutils.AssertFileCount(t, resp, 1)

	content := resp.File[0].GetContent()

	testutils.AssertContains(t, content, tc.wantGetter)
	testutils.AssertContains(t, content, tc.wantSetter)
}

// TestFieldTypes tests generation of various field types
func TestFieldTypes(t *testing.T) {
	testCases := []fieldTypeTestCase{
		newFieldTypeTestCase("double", descriptorpb.FieldDescriptorProto_TYPE_DOUBLE, "float64"),
		newFieldTypeTestCase("float", descriptorpb.FieldDescriptorProto_TYPE_FLOAT, "float32"),
		newFieldTypeTestCase("int64", descriptorpb.FieldDescriptorProto_TYPE_INT64, "int64"),
		newFieldTypeTestCase("uint64", descriptorpb.FieldDescriptorProto_TYPE_UINT64, "uint64"),
		newFieldTypeTestCase("bool", descriptorpb.FieldDescriptorProto_TYPE_BOOL, "bool"),
		newFieldTypeTestCase("string", descriptorpb.FieldDescriptorProto_TYPE_STRING, "string"),
		newFieldTypeTestCase("bytes", descriptorpb.FieldDescriptorProto_TYPE_BYTES, "[]byte"),
	}

	for _, tc := range testCases {
		t.Run(tc.name, tc.test)
	}
}

// createCalculatorProtoFile creates a proto file with a Calculator service
func createCalculatorProtoFile() *descriptorpb.FileDescriptorProto {
	protoFile := testutils.NewFileDescriptor("calculator.proto", "calculator", "github.com/example/calculator")
	protoFile.MessageType = []*descriptorpb.DescriptorProto{
		testutils.NewMessage("AddRequest",
			testutils.NewField("a", 1, descriptorpb.FieldDescriptorProto_TYPE_DOUBLE),
			testutils.NewField("b", 2, descriptorpb.FieldDescriptorProto_TYPE_DOUBLE),
		),
		testutils.NewMessage("AddResponse",
			testutils.NewField("result", 1, descriptorpb.FieldDescriptorProto_TYPE_DOUBLE),
		),
	}
	protoFile.Service = []*descriptorpb.ServiceDescriptorProto{
		testutils.NewService("Calculator",
			testutils.NewMethod("Add", ".calculator.AddRequest", ".calculator.AddResponse"),
		),
	}
	return protoFile
}

// TestGenerateServiceInterface tests generating an interface for a service
func TestGenerateServiceInterface(t *testing.T) {
	protoFile := createCalculatorProtoFile()

	req := testutils.NewCodeGenRequest(protoFile)
	opts := &GeneratorOptions{
		GenerateInterfaces: true,
		GenerateServices:   true,
	}
	resp := runGenerator(t, req, opts)

	testutils.AssertFileCount(t, resp, 1)

	content := resp.File[0].GetContent()

	// Check expected content
	expectedStrings := []string{
		// Service interface declaration with I% pattern
		"type ICalculatorService interface {",
		// Add method signature
		"Add(ctx context.Context, req IAddRequest) (IAddResponse, error)",
		// Context import
		`"context"`,
	}

	for _, expected := range expectedStrings {
		testutils.AssertContains(t, content, expected)
	}
}
