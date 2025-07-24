package testutils

import (
	"testing"

	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/descriptorpb"
	"google.golang.org/protobuf/types/pluginpb"
)

// NewFileDescriptor creates a FileDescriptorProto with common defaults
func NewFileDescriptor(name, pkg, goPkg string) *descriptorpb.FileDescriptorProto {
	return &descriptorpb.FileDescriptorProto{
		Name:    proto.String(name),
		Package: proto.String(pkg),
		Options: &descriptorpb.FileOptions{
			GoPackage: proto.String(goPkg),
		},
	}
}

// NewField creates a FieldDescriptorProto with common defaults
func NewField(
	name string,
	number int32,
	fieldType descriptorpb.FieldDescriptorProto_Type,
) *descriptorpb.FieldDescriptorProto {
	return &descriptorpb.FieldDescriptorProto{
		Name:   proto.String(name),
		Number: proto.Int32(number),
		Type:   fieldType.Enum(),
		Label:  descriptorpb.FieldDescriptorProto_LABEL_OPTIONAL.Enum(),
	}
}

// NewMessage creates a DescriptorProto with the given name and fields
func NewMessage(name string, fields ...*descriptorpb.FieldDescriptorProto) *descriptorpb.DescriptorProto {
	return &descriptorpb.DescriptorProto{
		Name:  proto.String(name),
		Field: fields,
	}
}

// NewService creates a ServiceDescriptorProto with the given name and methods
func NewService(name string, methods ...*descriptorpb.MethodDescriptorProto) *descriptorpb.ServiceDescriptorProto {
	return &descriptorpb.ServiceDescriptorProto{
		Name:   proto.String(name),
		Method: methods,
	}
}

// NewMethod creates a MethodDescriptorProto
func NewMethod(name, inputType, outputType string) *descriptorpb.MethodDescriptorProto {
	return &descriptorpb.MethodDescriptorProto{
		Name:       proto.String(name),
		InputType:  proto.String(inputType),
		OutputType: proto.String(outputType),
	}
}

// NewCodeGenRequest creates a CodeGeneratorRequest with the given files
func NewCodeGenRequest(files ...*descriptorpb.FileDescriptorProto) *pluginpb.CodeGeneratorRequest {
	fileNames := make([]string, len(files))
	for i, f := range files {
		fileNames[i] = f.GetName()
	}
	return &pluginpb.CodeGeneratorRequest{
		ProtoFile:      files,
		FileToGenerate: fileNames,
	}
}

// NewPlugin creates a protogen.Plugin from a FileDescriptorProto
func NewPlugin(t *testing.T, protoFile *descriptorpb.FileDescriptorProto) (*protogen.Plugin, error) {
	t.Helper()

	req := NewCodeGenRequest(protoFile)
	plugin, err := protogen.Options{}.New(req)
	if err != nil {
		return nil, err
	}

	return plugin, nil
}

// NewEnumField creates a FieldDescriptorProto for an enum field
func NewEnumField(name string, number int32, typeName string) *descriptorpb.FieldDescriptorProto {
	return &descriptorpb.FieldDescriptorProto{
		Name:     proto.String(name),
		Number:   proto.Int32(number),
		Type:     descriptorpb.FieldDescriptorProto_TYPE_ENUM.Enum(),
		TypeName: proto.String(typeName),
		Label:    descriptorpb.FieldDescriptorProto_LABEL_OPTIONAL.Enum(),
	}
}

// NewEnum creates an EnumDescriptorProto with the given name and values
func NewEnum(name string, values ...*descriptorpb.EnumValueDescriptorProto) *descriptorpb.EnumDescriptorProto {
	return &descriptorpb.EnumDescriptorProto{
		Name:  proto.String(name),
		Value: values,
	}
}

// NewEnumValue creates an EnumValueDescriptorProto
func NewEnumValue(name string, number int32) *descriptorpb.EnumValueDescriptorProto {
	return &descriptorpb.EnumValueDescriptorProto{
		Name:   proto.String(name),
		Number: proto.Int32(number),
	}
}

// RunGenerator runs the generator with the given request and options
// This requires the generator parameter to avoid circular dependencies
func RunGenerator(
	t *testing.T,
	req *pluginpb.CodeGeneratorRequest,
	genFunc func(*protogen.Plugin) error,
) *pluginpb.CodeGeneratorResponse {
	t.Helper()

	plugin, err := protogen.Options{}.New(req)
	if err != nil {
		t.Fatalf("failed to create plugin: %v", err)
	}

	if err := genFunc(plugin); err != nil {
		t.Fatalf("generator error: %v", err)
	}

	return plugin.Response()
}

// GetGeneratedFileContent returns the content of a generated file by name
func GetGeneratedFileContent(t *testing.T, plugin *protogen.Plugin, filename string) (string, bool) {
	t.Helper()

	response := plugin.Response()
	for _, file := range response.File {
		if file.GetName() == filename {
			return file.GetContent(), true
		}
	}

	return "", false
}
