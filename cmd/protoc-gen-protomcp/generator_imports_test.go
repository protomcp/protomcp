package main

import (
	"testing"

	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/types/descriptorpb"

	"protomcp.org/protomcp/pkg/generator/testutils"
)

// Test constants
const (
	testServiceName = "TestService"
	testEnumName    = "TestEnum"
	testEnumValue   = "TEST_ENUM_UNSPECIFIED"
)

// hasContentTestCase represents a test case for HasMessages/HasServices
type hasContentTestCase struct {
	setupFile   func() *descriptorpb.FileDescriptorProto
	name        string
	hasMessages bool
	hasServices bool
}

// test runs the HasContent test case
func (tc hasContentTestCase) test(t *testing.T) {
	t.Helper()

	protoFile := tc.setupFile()
	req := testutils.NewCodeGenRequest(protoFile)
	gen := tc.newGenerator(t)

	testutils.RunGenerator(t, req, gen)
}

func (tc hasContentTestCase) newGenerator(t *testing.T) func(*protogen.Plugin) error {
	return func(plugin *protogen.Plugin) error {
		gen := NewGenerator(plugin)

		for _, file := range plugin.Files {
			testutils.AssertEqual(t, gen.HasMessages(file), tc.hasMessages, "HasMessages()")
			testutils.AssertEqual(t, gen.HasServices(file), tc.hasServices, "HasServices()")
		}
		return nil
	}
}

// needsTypesTestCase represents a test case for NeedsTypes
type needsTypesTestCase struct {
	setupFile  func() *descriptorpb.FileDescriptorProto
	opts       *GeneratorOptions
	name       string
	needsTypes bool
}

// test runs the NeedsTypes test case
func (tc needsTypesTestCase) test(t *testing.T) {
	t.Helper()

	protoFile := tc.setupFile()
	req := testutils.NewCodeGenRequest(protoFile)
	gen := tc.newGenerator(t)

	testutils.RunGenerator(t, req, gen)
}

func (tc needsTypesTestCase) newGenerator(t *testing.T) func(*protogen.Plugin) error {
	return func(plugin *protogen.Plugin) error {
		gen := NewGenerator(plugin)

		for _, file := range plugin.Files {
			testutils.AssertEqual(t, gen.NeedsTypes(file, tc.opts), tc.needsTypes, "NeedsTypes()")
		}
		return nil
	}
}

// templateDataTestCase represents a test case for prepareTemplateData
type templateDataTestCase struct {
	name      string
	setupFile func() *descriptorpb.FileDescriptorProto
	opts      *GeneratorOptions
	want      [][]string
}

// test runs the template data test case
func (tc templateDataTestCase) test(t *testing.T) {
	t.Helper()

	protoFile := tc.setupFile()
	req := testutils.NewCodeGenRequest(protoFile)
	gen := tc.newGenerator(t)

	testutils.RunGenerator(t, req, gen)
}

func (tc templateDataTestCase) newGenerator(t *testing.T) func(*protogen.Plugin) error {
	return func(plugin *protogen.Plugin) error {
		gen := NewGenerator(plugin)

		for _, file := range plugin.Files {
			data := gen.prepareTemplateData(file, tc.opts)
			got := extractImportPaths(data)
			testutils.AssertSliceOfSlicesEqual(t, got, tc.want, "imports")
		}
		return nil
	}
}

// extractImportPaths extracts the import paths from template data
func extractImportPaths(data *TemplateData) [][]string {
	result := make([][]string, 0, len(data.ImportGroups))
	for _, group := range data.ImportGroups {
		var paths []string
		for _, imp := range group {
			paths = append(paths, imp.Path)
		}
		result = append(result, paths)
	}
	return result
}

func TestGeneratorHasContent(t *testing.T) {
	tests := []hasContentTestCase{
		{
			name: "file with messages only",
			setupFile: func() *descriptorpb.FileDescriptorProto {
				file := testutils.NewFileDescriptor("test.proto", "test", "github.com/example/test")
				file.MessageType = []*descriptorpb.DescriptorProto{
					testutils.NewMessage("TestMessage"),
				}
				return file
			},
			hasMessages: true,
			hasServices: false,
		},
		{
			name: "file with services only",
			setupFile: func() *descriptorpb.FileDescriptorProto {
				file := testutils.NewFileDescriptor("test.proto", "test", "github.com/example/test")
				serviceName := testServiceName
				file.Service = []*descriptorpb.ServiceDescriptorProto{
					{Name: &serviceName},
				}
				return file
			},
			hasMessages: false,
			hasServices: true,
		},
		{
			name: "file with both messages and services",
			setupFile: func() *descriptorpb.FileDescriptorProto {
				file := testutils.NewFileDescriptor("test.proto", "test", "github.com/example/test")
				file.MessageType = []*descriptorpb.DescriptorProto{
					testutils.NewMessage("TestMessage"),
				}
				serviceName := testServiceName
				file.Service = []*descriptorpb.ServiceDescriptorProto{
					{Name: &serviceName},
				}
				return file
			},
			hasMessages: true,
			hasServices: true,
		},
		{
			name: "file with enums only",
			setupFile: func() *descriptorpb.FileDescriptorProto {
				file := testutils.NewFileDescriptor("test.proto", "test", "github.com/example/test")
				enumName := testEnumName
				valueName := testEnumValue
				var valueNumber int32
				file.EnumType = []*descriptorpb.EnumDescriptorProto{
					{
						Name: &enumName,
						Value: []*descriptorpb.EnumValueDescriptorProto{
							{Name: &valueName, Number: &valueNumber},
						},
					},
				}
				return file
			},
			hasMessages: false,
			hasServices: false,
		},
		{
			name: "empty file",
			setupFile: func() *descriptorpb.FileDescriptorProto {
				return testutils.NewFileDescriptor("test.proto", "test", "github.com/example/test")
			},
			hasMessages: false,
			hasServices: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, tt.test)
	}
}

func TestGeneratorNeedsTypes(t *testing.T) {
	tests := []needsTypesTestCase{
		{
			name: "messages with interfaces enabled",
			setupFile: func() *descriptorpb.FileDescriptorProto {
				file := testutils.NewFileDescriptor("test.proto", "test", "github.com/example/test")
				file.MessageType = []*descriptorpb.DescriptorProto{
					testutils.NewMessage("TestMessage"),
				}
				return file
			},
			opts: &GeneratorOptions{
				GenerateInterfaces: true,
				GenerateServices:   false,
			},
			needsTypes: true,
		},
		{
			name: "messages with interfaces disabled",
			setupFile: func() *descriptorpb.FileDescriptorProto {
				file := testutils.NewFileDescriptor("test.proto", "test", "github.com/example/test")
				file.MessageType = []*descriptorpb.DescriptorProto{
					testutils.NewMessage("TestMessage"),
				}
				return file
			},
			opts: &GeneratorOptions{
				GenerateInterfaces: false,
				GenerateServices:   false,
			},
			needsTypes: false,
		},
		{
			name: "services with services enabled",
			setupFile: func() *descriptorpb.FileDescriptorProto {
				file := testutils.NewFileDescriptor("test.proto", "test", "github.com/example/test")
				serviceName := testServiceName
				file.Service = []*descriptorpb.ServiceDescriptorProto{
					{Name: &serviceName},
				}
				return file
			},
			opts: &GeneratorOptions{
				GenerateInterfaces: false,
				GenerateServices:   true,
			},
			needsTypes: true,
		},
		{
			name: "enums only with everything enabled",
			setupFile: func() *descriptorpb.FileDescriptorProto {
				file := testutils.NewFileDescriptor("test.proto", "test", "github.com/example/test")
				enumName := testEnumName
				valueName := testEnumValue
				var valueNumber int32
				file.EnumType = []*descriptorpb.EnumDescriptorProto{
					{
						Name: &enumName,
						Value: []*descriptorpb.EnumValueDescriptorProto{
							{Name: &valueName, Number: &valueNumber},
						},
					},
				}
				return file
			},
			opts: &GeneratorOptions{
				GenerateInterfaces: true,
				GenerateServices:   true,
			},
			needsTypes: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, tt.test)
	}
}

// generatorImportsTestCases returns test cases for generator imports
func generatorImportsTestCases() []templateDataTestCase {
	return []templateDataTestCase{
		{
			name: "messages only - types file",
			setupFile: func() *descriptorpb.FileDescriptorProto {
				file := testutils.NewFileDescriptor("test.proto", "test", "github.com/example/test")
				file.MessageType = []*descriptorpb.DescriptorProto{
					testutils.NewMessage("TestMessage"),
				}
				return file
			},
			opts: &GeneratorOptions{
				GenerateInterfaces: true,
				GenerateServices:   false,
			},
			want: [][]string{},
		},
		{
			name: "services only - types file",
			setupFile: func() *descriptorpb.FileDescriptorProto {
				file := testutils.NewFileDescriptor("test.proto", "test", "github.com/example/test")
				serviceName := testServiceName
				file.Service = []*descriptorpb.ServiceDescriptorProto{
					{Name: &serviceName},
				}
				return file
			},
			opts: &GeneratorOptions{
				GenerateInterfaces: false,
				GenerateServices:   true,
			},
			want: [][]string{
				{"context"},
			},
		},
		{
			name: "messages and services - types file",
			setupFile: func() *descriptorpb.FileDescriptorProto {
				file := testutils.NewFileDescriptor("test.proto", "test", "github.com/example/test")
				file.MessageType = []*descriptorpb.DescriptorProto{
					testutils.NewMessage("TestMessage"),
				}
				serviceName := testServiceName
				file.Service = []*descriptorpb.ServiceDescriptorProto{
					{Name: &serviceName},
				}
				return file
			},
			opts: &GeneratorOptions{
				GenerateInterfaces: true,
				GenerateServices:   true,
			},
			want: [][]string{
				{"context"},
			},
		},
		{
			name: "enums only - no imports",
			setupFile: func() *descriptorpb.FileDescriptorProto {
				file := testutils.NewFileDescriptor("test.proto", "test", "github.com/example/test")
				enumName := testEnumName
				valueName := testEnumValue
				var valueNumber int32
				file.EnumType = []*descriptorpb.EnumDescriptorProto{
					{
						Name: &enumName,
						Value: []*descriptorpb.EnumValueDescriptorProto{
							{Name: &valueName, Number: &valueNumber},
						},
					},
				}
				return file
			},
			opts: &GeneratorOptions{
				GenerateInterfaces: true,
				GenerateServices:   true,
			},
			want: [][]string{},
		},
	}
}

func TestGeneratorImports(t *testing.T) {
	tests := generatorImportsTestCases()

	for _, tt := range tests {
		t.Run(tt.name, tt.test)
	}
}
