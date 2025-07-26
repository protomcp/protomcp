package main

import (
	"testing"

	"google.golang.org/protobuf/types/descriptorpb"

	gengo "protomcp.org/protomcp/pkg/generator/gen-go"
	"protomcp.org/protomcp/pkg/generator/testutils"
)

// needsNoImplTestCase represents a test case for NeedsNoImpl
type needsNoImplTestCase struct {
	setupFile  func() *descriptorpb.FileDescriptorProto
	opts       *GeneratorOptions
	name       string
	wantNoImpl bool
}

// test runs the NeedsNoImpl test case
func (tc needsNoImplTestCase) test(t *testing.T) {
	t.Helper()

	// Create plugin from proto file
	plugin, err := testutils.NewPlugin(t, tc.setupFile())
	if err != nil {
		t.Fatalf("failed to create plugin: %v", err)
	}
	file := plugin.Files[0]

	// Test NeedsNoImpl
	got := tc.opts.NeedsNoImpl(file)
	if got != tc.wantNoImpl {
		t.Errorf("NeedsNoImpl() = %v, want %v", got, tc.wantNoImpl)
	}
}

func TestGeneratorOptions_NeedsNoImpl(t *testing.T) {
	tests := []needsNoImplTestCase{
		{
			name: "nil options with messages uses defaults",
			setupFile: func() *descriptorpb.FileDescriptorProto {
				file := testutils.NewFileDescriptor("test.proto", "test", "github.com/example/test")
				file.MessageType = []*descriptorpb.DescriptorProto{
					testutils.NewMessage("TestMessage",
						testutils.NewField("id", 1, descriptorpb.FieldDescriptorProto_TYPE_STRING),
					),
				}
				return file
			},
			opts:       nil,
			wantNoImpl: true,
		},
		{
			name: "noImpl disabled returns false",
			setupFile: func() *descriptorpb.FileDescriptorProto {
				file := testutils.NewFileDescriptor("test.proto", "test", "github.com/example/test")
				file.MessageType = []*descriptorpb.DescriptorProto{
					testutils.NewMessage("TestMessage",
						testutils.NewField("id", 1, descriptorpb.FieldDescriptorProto_TYPE_STRING),
					),
				}
				return file
			},
			opts: &GeneratorOptions{
				GenerateNoImpl:     false,
				GenerateInterfaces: true,
				GenerateServices:   true,
			},
			wantNoImpl: false,
		},
		{
			name: "noImpl enabled with messages",
			setupFile: func() *descriptorpb.FileDescriptorProto {
				file := testutils.NewFileDescriptor("test.proto", "test", "github.com/example/test")
				file.MessageType = []*descriptorpb.DescriptorProto{
					testutils.NewMessage("TestMessage",
						testutils.NewField("id", 1, descriptorpb.FieldDescriptorProto_TYPE_STRING),
					),
				}
				return file
			},
			opts: &GeneratorOptions{
				GenerateNoImpl:     true,
				GenerateInterfaces: true,
				GenerateServices:   true,
			},
			wantNoImpl: true,
		},
		{
			name: "noImpl enabled with services",
			setupFile: func() *descriptorpb.FileDescriptorProto {
				file := testutils.NewFileDescriptor("test.proto", "test", "github.com/example/test")
				// Add required message types for the service
				file.MessageType = []*descriptorpb.DescriptorProto{
					testutils.NewMessage("GetItemRequest",
						testutils.NewField("id", 1, descriptorpb.FieldDescriptorProto_TYPE_STRING),
					),
					testutils.NewMessage("GetItemResponse",
						testutils.NewField("item", 1, descriptorpb.FieldDescriptorProto_TYPE_STRING),
					),
				}
				file.Service = []*descriptorpb.ServiceDescriptorProto{
					testutils.NewService("TestService",
						testutils.NewMethod("GetItem", "GetItemRequest", "GetItemResponse"),
					),
				}
				return file
			},
			opts: &GeneratorOptions{
				GenerateNoImpl:     true,
				GenerateInterfaces: true,
				GenerateServices:   true,
			},
			wantNoImpl: true,
		},
		{
			name: "noImpl enabled with both messages and services",
			setupFile: func() *descriptorpb.FileDescriptorProto {
				file := testutils.NewFileDescriptor("test.proto", "test", "github.com/example/test")
				file.MessageType = []*descriptorpb.DescriptorProto{
					testutils.NewMessage("TestMessage",
						testutils.NewField("id", 1, descriptorpb.FieldDescriptorProto_TYPE_STRING),
					),
					testutils.NewMessage("GetItemRequest",
						testutils.NewField("id", 1, descriptorpb.FieldDescriptorProto_TYPE_STRING),
					),
					testutils.NewMessage("GetItemResponse",
						testutils.NewField("item", 1, descriptorpb.FieldDescriptorProto_TYPE_STRING),
					),
				}
				file.Service = []*descriptorpb.ServiceDescriptorProto{
					testutils.NewService("TestService",
						testutils.NewMethod("GetItem", "GetItemRequest", "GetItemResponse"),
					),
				}
				return file
			},
			opts: &GeneratorOptions{
				GenerateNoImpl:     true,
				GenerateInterfaces: true,
				GenerateServices:   true,
			},
			wantNoImpl: true,
		},
		{
			name: "noImpl enabled but no messages or services",
			setupFile: func() *descriptorpb.FileDescriptorProto {
				return testutils.NewFileDescriptor("test.proto", "test", "github.com/example/test")
			},
			opts: &GeneratorOptions{
				GenerateNoImpl:     true,
				GenerateInterfaces: true,
				GenerateServices:   true,
			},
			wantNoImpl: false,
		},
		{
			name: "noImpl enabled with only enums",
			setupFile: func() *descriptorpb.FileDescriptorProto {
				file := testutils.NewFileDescriptor("test.proto", "test", "github.com/example/test")
				file.EnumType = []*descriptorpb.EnumDescriptorProto{
					testutils.NewEnum("Status",
						testutils.NewEnumValue("STATUS_UNSPECIFIED", 0),
						testutils.NewEnumValue("STATUS_ACTIVE", 1),
					),
				}
				return file
			},
			opts: &GeneratorOptions{
				GenerateNoImpl:     true,
				GenerateInterfaces: true,
				GenerateServices:   true,
				GenerateEnums:      true,
			},
			wantNoImpl: false,
		},
		{
			name: "noImpl enabled but interfaces disabled",
			setupFile: func() *descriptorpb.FileDescriptorProto {
				file := testutils.NewFileDescriptor("test.proto", "test", "github.com/example/test")
				file.MessageType = []*descriptorpb.DescriptorProto{
					testutils.NewMessage("TestMessage",
						testutils.NewField("id", 1, descriptorpb.FieldDescriptorProto_TYPE_STRING),
					),
				}
				return file
			},
			opts: &GeneratorOptions{
				GenerateNoImpl:     true,
				GenerateInterfaces: false,
				GenerateServices:   true,
			},
			wantNoImpl: false,
		},
		{
			name: "noImpl enabled with messages and interfaces but services disabled",
			setupFile: func() *descriptorpb.FileDescriptorProto {
				file := testutils.NewFileDescriptor("test.proto", "test", "github.com/example/test")
				// Add required message types for the service
				file.MessageType = []*descriptorpb.DescriptorProto{
					testutils.NewMessage("GetItemRequest",
						testutils.NewField("id", 1, descriptorpb.FieldDescriptorProto_TYPE_STRING),
					),
					testutils.NewMessage("GetItemResponse",
						testutils.NewField("item", 1, descriptorpb.FieldDescriptorProto_TYPE_STRING),
					),
				}
				file.Service = []*descriptorpb.ServiceDescriptorProto{
					testutils.NewService("TestService",
						testutils.NewMethod("GetItem", "GetItemRequest", "GetItemResponse"),
					),
				}
				return file
			},
			opts: &GeneratorOptions{
				GenerateNoImpl:     true,
				GenerateInterfaces: true,
				GenerateServices:   false,
			},
			wantNoImpl: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, tt.test)
	}
}

// needsEnumsTestCase represents a test case for NeedsEnums
type needsEnumsTestCase struct {
	setupFile func() *descriptorpb.FileDescriptorProto
	opts      *GeneratorOptions
	name      string
	wantEnums bool
}

// test runs the NeedsEnums test case
func (tc needsEnumsTestCase) test(t *testing.T) {
	t.Helper()

	// Create plugin from proto file
	plugin, err := testutils.NewPlugin(t, tc.setupFile())
	if err != nil {
		t.Fatalf("failed to create plugin: %v", err)
	}
	file := plugin.Files[0]

	// Test NeedsEnums
	got := tc.opts.NeedsEnums(file)
	if got != tc.wantEnums {
		t.Errorf("NeedsEnums() = %v, want %v", got, tc.wantEnums)
	}
}

func TestGeneratorOptions_NeedsEnums(t *testing.T) {
	tests := []needsEnumsTestCase{
		{
			name: "nil options uses defaults",
			setupFile: func() *descriptorpb.FileDescriptorProto {
				file := testutils.NewFileDescriptor("test.proto", "test", "github.com/example/test")
				file.EnumType = []*descriptorpb.EnumDescriptorProto{
					testutils.NewEnum("Status",
						testutils.NewEnumValue("STATUS_UNSPECIFIED", 0),
						testutils.NewEnumValue("STATUS_ACTIVE", 1),
					),
				}
				return file
			},
			opts:      nil,
			wantEnums: true, // Default is now true
		},
		{
			name: "enums disabled returns false",
			setupFile: func() *descriptorpb.FileDescriptorProto {
				file := testutils.NewFileDescriptor("test.proto", "test", "github.com/example/test")
				file.EnumType = []*descriptorpb.EnumDescriptorProto{
					testutils.NewEnum("Status",
						testutils.NewEnumValue("STATUS_UNSPECIFIED", 0),
						testutils.NewEnumValue("STATUS_ACTIVE", 1),
					),
				}
				return file
			},
			opts: &GeneratorOptions{
				GenerateEnums: false,
			},
			wantEnums: false,
		},
		{
			name: "enums enabled with enums",
			setupFile: func() *descriptorpb.FileDescriptorProto {
				file := testutils.NewFileDescriptor("test.proto", "test", "github.com/example/test")
				file.EnumType = []*descriptorpb.EnumDescriptorProto{
					testutils.NewEnum("Status",
						testutils.NewEnumValue("STATUS_UNSPECIFIED", 0),
						testutils.NewEnumValue("STATUS_ACTIVE", 1),
					),
				}
				return file
			},
			opts: &GeneratorOptions{
				GenerateEnums: true,
			},
			wantEnums: true,
		},
		{
			name: "enums enabled but no enums in file",
			setupFile: func() *descriptorpb.FileDescriptorProto {
				file := testutils.NewFileDescriptor("test.proto", "test", "github.com/example/test")
				file.MessageType = []*descriptorpb.DescriptorProto{
					testutils.NewMessage("TestMessage",
						testutils.NewField("id", 1, descriptorpb.FieldDescriptorProto_TYPE_STRING),
					),
				}
				return file
			},
			opts: &GeneratorOptions{
				GenerateEnums: true,
			},
			wantEnums: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, tt.test)
	}
}

// enumPatternTestCase represents a test case for enum pattern naming
type enumPatternTestCase struct {
	name     string
	enumName string
	pattern  string
	wantName string
}

// test runs the enum pattern test case
func (tc enumPatternTestCase) test(t *testing.T) {
	t.Helper()

	// Create a mock enum
	file := testutils.NewFileDescriptor("test.proto", "test", "github.com/example/test")
	file.EnumType = []*descriptorpb.EnumDescriptorProto{
		testutils.NewEnum(tc.enumName,
			testutils.NewEnumValue(tc.enumName+"_UNSPECIFIED", 0),
		),
	}

	plugin, err := testutils.NewPlugin(t, file)
	if err != nil {
		t.Fatalf("failed to create plugin: %v", err)
	}
	enum := plugin.Files[0].Enums[0]

	// Apply pattern
	got := gengo.EnumNameFor(enum, tc.pattern)
	if got != tc.wantName {
		t.Errorf("EnumNameFor(%q, %q) = %q, want %q", tc.enumName, tc.pattern, got, tc.wantName)
	}
}

func TestEnumNamePattern(t *testing.T) {
	tests := []enumPatternTestCase{
		{
			name:     "default pattern with Status",
			enumName: "Status",
			pattern:  "%Enum",
			wantName: "StatusEnum",
		},
		{
			name:     "prefix pattern",
			enumName: "Status",
			pattern:  "E%",
			wantName: "EStatus",
		},
		{
			name:     "no pattern",
			enumName: "Status",
			pattern:  "%",
			wantName: "Status",
		},
		{
			name:     "empty pattern defaults to no change",
			enumName: "Status",
			pattern:  "",
			wantName: "Status",
		},
		{
			name:     "complex pattern",
			enumName: "FileType",
			pattern:  "My%Enum",
			wantName: "MyFileTypeEnum",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, tt.test)
	}
}
