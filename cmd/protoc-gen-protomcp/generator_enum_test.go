package main

import (
	"strings"
	"testing"

	"google.golang.org/protobuf/types/descriptorpb"

	"protomcp.org/protomcp/pkg/generator/testutils"
)

// enumDefaultsTestCase represents a test case for GeneratorOptions enum defaults
type enumDefaultsTestCase struct {
	opts              *GeneratorOptions
	name              string
	wantEnumPattern   string
	wantGenerateEnums bool
}

// test runs the enum defaults test case
func (tc enumDefaultsTestCase) test(t *testing.T) {
	t.Helper()

	testutils.AssertEqual(t, tc.opts.GetGenerateEnums(), tc.wantGenerateEnums, "GetGenerateEnums()")

	// Only check pattern if not empty string expected
	if tc.wantEnumPattern != "" {
		testutils.AssertEqual(t, tc.opts.GetEnumPattern(), tc.wantEnumPattern, "GetEnumPattern()")
	} else {
		// For empty string, check that it returns the default
		testutils.AssertEqual(t, tc.opts.GetEnumPattern(), DefaultEnumPattern, "GetEnumPattern() default")
	}
}

func TestGeneratorOptions_EnumDefaults(t *testing.T) {
	tests := []enumDefaultsTestCase{
		{
			name:              "nil options uses defaults",
			opts:              nil,
			wantGenerateEnums: DefaultGenerateEnums,
			wantEnumPattern:   DefaultEnumPattern,
		},
		{
			name:              "empty options uses defaults",
			opts:              &GeneratorOptions{},
			wantGenerateEnums: false, // Zero value for bool
			wantEnumPattern:   "",    // Zero value for string
		},
		{
			name: "explicit values override defaults",
			opts: &GeneratorOptions{
				GenerateEnums: false,
				EnumPattern:   "E%",
			},
			wantGenerateEnums: false,
			wantEnumPattern:   "E%",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, tt.test)
	}
}

// hasEnumsTestCase represents a test case for GeneratorOptions.HasEnums
type hasEnumsTestCase struct {
	setupFile func() *descriptorpb.FileDescriptorProto
	name      string
	want      bool
}

// test runs the HasEnums test case
func (tc hasEnumsTestCase) test(t *testing.T) {
	t.Helper()

	var opts GeneratorOptions

	// For nil file test, we handle it directly
	if tc.name == "nil file" {
		testutils.AssertEqual(t, opts.HasEnums(nil), tc.want, "HasEnums(nil)")
		return
	}

	// For other tests, create the plugin
	plugin, err := testutils.NewPlugin(t, tc.setupFile())
	if err != nil {
		t.Fatalf("failed to create plugin: %v", err)
	}
	file := plugin.Files[0]

	testutils.AssertEqual(t, opts.HasEnums(file), tc.want, "HasEnums()")
}

func TestGeneratorOptions_HasEnums(t *testing.T) {
	tests := []hasEnumsTestCase{
		{
			name: "file with enums",
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
			want: true,
		},
		{
			name: "file without enums",
			setupFile: func() *descriptorpb.FileDescriptorProto {
				return testutils.NewFileDescriptor("test.proto", "test", "github.com/example/test")
			},
			want: false,
		},
		{
			name: "nil file",
			setupFile: func() *descriptorpb.FileDescriptorProto {
				return nil
			},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, tt.test)
	}
}

// enumGenerationTestCase represents a test case for enum generation
type enumGenerationTestCase struct {
	name        string
	setupFile   func() *descriptorpb.FileDescriptorProto
	opts        *GeneratorOptions
	wantEnums   []string
	wantConsts  []string
	wantMethods []string
	wantMaps    []string
	wantImports []string
}

// makeContent generates code and returns the generated types file content
func (tc enumGenerationTestCase) makeContent(t *testing.T) string {
	t.Helper()

	// Create plugin
	plugin, err := testutils.NewPlugin(t, tc.setupFile())
	if err != nil {
		t.Fatalf("failed to create plugin: %v", err)
	}

	// Generate code
	gen := NewGenerator(plugin)
	file := plugin.Files[0]
	if err := gen.GenerateFile(file, tc.opts); err != nil {
		t.Fatalf("GenerateFile() error = %v", err)
	}

	// Get generated content
	content, ok := testutils.GetGeneratedFileContent(t, plugin, "github.com/example/test/test.types.go")
	if !ok {
		t.Fatal("generated file not found")
	}

	return content
}

// test runs the enum generation test case
func (tc enumGenerationTestCase) test(t *testing.T) {
	t.Helper()

	content := tc.makeContent(t)

	// Check enum types
	for _, expected := range tc.wantEnums {
		testutils.AssertContains(t, content, expected)
	}

	// Check constants
	for _, expected := range tc.wantConsts {
		testutils.AssertContains(t, content, expected)
	}

	// Check methods
	for _, expected := range tc.wantMethods {
		testutils.AssertContains(t, content, expected)
	}

	// Check maps
	for _, expected := range tc.wantMaps {
		testutils.AssertContains(t, content, expected)
	}

	// Check imports
	for _, expected := range tc.wantImports {
		testutils.AssertContains(t, content, expected)
	}
}

func TestGenerator_EnumGeneration(t *testing.T) {
	tests := []enumGenerationTestCase{
		{
			name: "comprehensive enum generation",
			setupFile: func() *descriptorpb.FileDescriptorProto {
				file := testutils.NewFileDescriptor("test.proto", "test", "github.com/example/test")
				file.EnumType = []*descriptorpb.EnumDescriptorProto{
					testutils.NewEnum("Status",
						testutils.NewEnumValue("STATUS_UNSPECIFIED", 0),
						testutils.NewEnumValue("STATUS_ACTIVE", 1),
						testutils.NewEnumValue("STATUS_INACTIVE", 2),
					),
					testutils.NewEnum("Priority",
						testutils.NewEnumValue("PRIORITY_UNSPECIFIED", 0),
						testutils.NewEnumValue("PRIORITY_LOW", 1),
						testutils.NewEnumValue("PRIORITY_MEDIUM", 2),
						testutils.NewEnumValue("PRIORITY_HIGH", 3),
					),
				}
				return file
			},
			opts: &GeneratorOptions{
				GenerateEnums: true,
				EnumPattern:   "%Enum",
			},
			wantEnums: []string{
				"type StatusEnum int32",
				"type PriorityEnum int32",
			},
			wantConsts: []string{
				"StatusEnum_STATUS_UNSPECIFIED StatusEnum = 0",
				"StatusEnum_STATUS_ACTIVE      StatusEnum = 1",
				"StatusEnum_STATUS_INACTIVE    StatusEnum = 2",
				"PriorityEnum_PRIORITY_UNSPECIFIED PriorityEnum = 0",
				"PriorityEnum_PRIORITY_LOW         PriorityEnum = 1",
				"PriorityEnum_PRIORITY_MEDIUM      PriorityEnum = 2",
				"PriorityEnum_PRIORITY_HIGH        PriorityEnum = 3",
			},
			wantMethods: []string{
				"func (x StatusEnum) String() string",
				"func (x StatusEnum) IsValid() bool",
				"func (x StatusEnum) MarshalText() ([]byte, error)",
				"func (x *StatusEnum) UnmarshalText(text []byte) error",
				"func (x PriorityEnum) String() string",
				"func (x PriorityEnum) IsValid() bool",
				"func (x PriorityEnum) MarshalText() ([]byte, error)",
				"func (x *PriorityEnum) UnmarshalText(text []byte) error",
			},
			wantMaps: []string{
				"var statusEnum_name = map[int32]string{",
				"var statusEnum_value = map[string]int32{",
				"var priorityEnum_name = map[int32]string{",
				"var priorityEnum_value = map[string]int32{",
			},
			wantImports: []string{
				`"errors"`,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, tt.test)
	}
}

// enumPatternsTestCase represents a test case for enum pattern generation
type enumPatternsTestCase struct {
	name            string
	pattern         string
	enumName        string
	wantTypeName    string
	wantConstPrefix string
}

// makeContent generates code and returns the generated types file content
func (tc enumPatternsTestCase) makeContent(t *testing.T) string {
	t.Helper()

	// Create proto file with enum
	protoFile := testutils.NewFileDescriptor("test.proto", "test", "github.com/example/test")
	protoFile.EnumType = []*descriptorpb.EnumDescriptorProto{
		testutils.NewEnum(tc.enumName,
			testutils.NewEnumValue(tc.enumName+"_UNSPECIFIED", 0),
			testutils.NewEnumValue(tc.enumName+"_ACTIVE", 1),
		),
	}

	// Create plugin
	plugin, err := testutils.NewPlugin(t, protoFile)
	if err != nil {
		t.Fatalf("failed to create plugin: %v", err)
	}

	// Generate code with custom pattern
	gen := NewGenerator(plugin)
	opts := &GeneratorOptions{
		GenerateEnums: true,
		EnumPattern:   tc.pattern,
	}

	file := plugin.Files[0]
	if err := gen.GenerateFile(file, opts); err != nil {
		t.Fatalf("GenerateFile() error = %v", err)
	}

	// Get generated content
	content, ok := testutils.GetGeneratedFileContent(t, plugin, "github.com/example/test/test.types.go")
	if !ok {
		t.Fatal("generated file not found")
	}

	return content
}

// test runs the enum patterns test case
func (tc enumPatternsTestCase) test(t *testing.T) {
	t.Helper()

	content := tc.makeContent(t)

	// Check type name
	expectedType := "type " + tc.wantTypeName + " int32"
	testutils.AssertContains(t, content, expectedType)

	// Check constant - need to include the full enum name prefix
	expectedConst := tc.wantConstPrefix + tc.enumName + "_UNSPECIFIED"
	testutils.AssertContains(t, content, expectedConst)
}

func TestGenerator_EnumPatterns(t *testing.T) {
	tests := []enumPatternsTestCase{
		{
			name:            "suffix pattern (default)",
			pattern:         "%Enum",
			enumName:        "Status",
			wantTypeName:    "StatusEnum",
			wantConstPrefix: "StatusEnum_",
		},
		{
			name:            "prefix pattern",
			pattern:         "E%",
			enumName:        "Status",
			wantTypeName:    "EStatus",
			wantConstPrefix: "EStatus_",
		},
		{
			name:            "custom suffix",
			pattern:         "%Type",
			enumName:        "Priority",
			wantTypeName:    "PriorityType",
			wantConstPrefix: "PriorityType_",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, tt.test)
	}
}

// enumDisabledTestCase represents a test case for disabled enum generation
type enumDisabledTestCase struct {
	setupFile     func() *descriptorpb.FileDescriptorProto
	opts          *GeneratorOptions
	name          string
	wantFileCount int
}

// test runs the enum disabled test case
func (tc enumDisabledTestCase) test(t *testing.T) {
	t.Helper()

	// Create plugin
	plugin, err := testutils.NewPlugin(t, tc.setupFile())
	if err != nil {
		t.Fatalf("failed to create plugin: %v", err)
	}

	// Generate code
	gen := NewGenerator(plugin)
	file := plugin.Files[0]
	if err := gen.GenerateFile(file, tc.opts); err != nil {
		t.Fatalf("GenerateFile() error = %v", err)
	}

	// Check generated file count
	testutils.AssertFileCount(t, plugin.Response(), tc.wantFileCount)
}

func TestGenerator_EnumDisabled(t *testing.T) {
	tests := []enumDisabledTestCase{
		{
			name: "enum generation disabled",
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
			wantFileCount: 0, // No files should be generated
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, tt.test)
	}
}

// enumWithMessagesTestCase represents a test case for enum generation with messages and services
type enumWithMessagesTestCase struct {
	name              string
	setupFile         func() *descriptorpb.FileDescriptorProto
	opts              *GeneratorOptions
	wantInTypes       []string
	wantInNoImpl      []string
	doNotWantInNoImpl []string
}

// makeContent generates code and returns both types and noimpl file contents
func (tc enumWithMessagesTestCase) makeContent(t *testing.T) (typesContent, noimplContent string) {
	t.Helper()

	// Create plugin
	plugin, err := testutils.NewPlugin(t, tc.setupFile())
	if err != nil {
		t.Fatalf("failed to create plugin: %v", err)
	}

	// Generate code
	gen := NewGenerator(plugin)
	file := plugin.Files[0]
	if err := gen.GenerateFile(file, tc.opts); err != nil {
		t.Fatalf("GenerateFile() error = %v", err)
	}

	// Get generated content
	typesContent, ok := testutils.GetGeneratedFileContent(t, plugin, "github.com/example/test/test.types.go")
	if !ok {
		t.Fatal("types file not found")
	}

	noimplContent, ok = testutils.GetGeneratedFileContent(t, plugin, "github.com/example/test/test.noimpl.go")
	if !ok {
		t.Fatal("noimpl file not found")
	}

	return typesContent, noimplContent
}

// test runs the enum with messages test case
func (tc enumWithMessagesTestCase) test(t *testing.T) {
	t.Helper()

	typesContent, noimplContent := tc.makeContent(t)

	// Check types file
	for _, expected := range tc.wantInTypes {
		testutils.AssertContains(t, typesContent, expected)
	}

	// Check noimpl file
	for _, expected := range tc.wantInNoImpl {
		testutils.AssertContains(t, noimplContent, expected)
	}

	// Check things that should NOT be in noimpl
	for _, notExpected := range tc.doNotWantInNoImpl {
		if strings.Contains(noimplContent, notExpected) {
			t.Errorf("NoImpl file should not contain %q", notExpected)
		}
	}
}

func TestGenerator_EnumWithMessagesAndServices(t *testing.T) {
	tests := []enumWithMessagesTestCase{
		{
			name: "enum with messages and services",
			setupFile: func() *descriptorpb.FileDescriptorProto {
				file := testutils.NewFileDescriptor("test.proto", "test", "github.com/example/test")

				// Add enum
				file.EnumType = []*descriptorpb.EnumDescriptorProto{
					testutils.NewEnum("Status",
						testutils.NewEnumValue("STATUS_UNSPECIFIED", 0),
						testutils.NewEnumValue("STATUS_ACTIVE", 1),
					),
				}

				// Add message
				file.MessageType = []*descriptorpb.DescriptorProto{
					testutils.NewMessage("Item",
						testutils.NewField("id", 1, descriptorpb.FieldDescriptorProto_TYPE_STRING),
						testutils.NewEnumField("status", 2, ".test.Status"),
					),
				}

				// Add service
				file.Service = []*descriptorpb.ServiceDescriptorProto{
					testutils.NewService("ItemService",
						testutils.NewMethod("GetItem", "Item", "Item"),
					),
				}

				return file
			},
			opts: &GeneratorOptions{
				GenerateInterfaces: true,
				GenerateServices:   true,
				GenerateEnums:      true,
				GenerateNoImpl:     true,
			},
			wantInTypes: []string{
				"type StatusEnum int32",
				"type IItem interface",
				"type IItemService interface",
			},
			wantInNoImpl: []string{
				"type NoImplItem struct",
			},
			doNotWantInNoImpl: []string{
				"StatusEnum",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, tt.test)
	}
}
