package testutils_test

import (
	"testing"

	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/types/descriptorpb"

	"protomcp.org/protomcp/pkg/generator/testutils"
)

func TestNewFileDescriptor(t *testing.T) {
	name := "test.proto"
	pkg := "test.v1"
	goPkg := "github.com/example/test/v1"

	file := testutils.NewFileDescriptor(name, pkg, goPkg)

	if file.GetName() != name {
		t.Errorf("file name = %q, want %q", file.GetName(), name)
	}
	if file.GetPackage() != pkg {
		t.Errorf("package = %q, want %q", file.GetPackage(), pkg)
	}
	if file.GetOptions().GetGoPackage() != goPkg {
		t.Errorf("go package = %q, want %q", file.GetOptions().GetGoPackage(), goPkg)
	}
}

func TestNewField(t *testing.T) {
	name := "user_id"
	number := int32(1)
	fieldType := descriptorpb.FieldDescriptorProto_TYPE_STRING

	field := testutils.NewField(name, number, fieldType)

	if field.GetName() != name {
		t.Errorf("field name = %q, want %q", field.GetName(), name)
	}
	if field.GetNumber() != number {
		t.Errorf("field number = %d, want %d", field.GetNumber(), number)
	}
	if field.GetType() != fieldType {
		t.Errorf("field type = %v, want %v", field.GetType(), fieldType)
	}
	if field.GetLabel() != descriptorpb.FieldDescriptorProto_LABEL_OPTIONAL {
		t.Errorf("field label = %v, want LABEL_OPTIONAL", field.GetLabel())
	}
}

func TestNewMessage(t *testing.T) {
	name := "User"
	field1 := testutils.NewField("id", 1, descriptorpb.FieldDescriptorProto_TYPE_STRING)
	field2 := testutils.NewField("name", 2, descriptorpb.FieldDescriptorProto_TYPE_STRING)

	msg := testutils.NewMessage(name, field1, field2)

	if msg.GetName() != name {
		t.Errorf("message name = %q, want %q", msg.GetName(), name)
	}
	if len(msg.Field) != 2 {
		t.Fatalf("message has %d fields, want 2", len(msg.Field))
	}
	if msg.Field[0].GetName() != "id" {
		t.Errorf("first field name = %q, want %q", msg.Field[0].GetName(), "id")
	}
	if msg.Field[1].GetName() != "name" {
		t.Errorf("second field name = %q, want %q", msg.Field[1].GetName(), "name")
	}
}

func TestNewService(t *testing.T) {
	name := "UserService"
	method1 := testutils.NewMethod("GetUser", ".test.GetUserRequest", ".test.User")
	method2 := testutils.NewMethod("ListUsers", ".test.ListUsersRequest", ".test.ListUsersResponse")

	service := testutils.NewService(name, method1, method2)

	if service.GetName() != name {
		t.Errorf("service name = %q, want %q", service.GetName(), name)
	}
	if len(service.Method) != 2 {
		t.Fatalf("service has %d methods, want 2", len(service.Method))
	}
	if service.Method[0].GetName() != "GetUser" {
		t.Errorf("first method name = %q, want %q", service.Method[0].GetName(), "GetUser")
	}
	if service.Method[1].GetName() != "ListUsers" {
		t.Errorf("second method name = %q, want %q", service.Method[1].GetName(), "ListUsers")
	}
}

func TestNewMethod(t *testing.T) {
	name := "GetUser"
	inputType := ".test.GetUserRequest"
	outputType := ".test.User"

	method := testutils.NewMethod(name, inputType, outputType)

	if method.GetName() != name {
		t.Errorf("method name = %q, want %q", method.GetName(), name)
	}
	if method.GetInputType() != inputType {
		t.Errorf("input type = %q, want %q", method.GetInputType(), inputType)
	}
	if method.GetOutputType() != outputType {
		t.Errorf("output type = %q, want %q", method.GetOutputType(), outputType)
	}
}

func TestNewCodeGenRequest(t *testing.T) {
	file1 := testutils.NewFileDescriptor("file1.proto", "test.v1", "test/v1")
	file2 := testutils.NewFileDescriptor("file2.proto", "test.v2", "test/v2")

	req := testutils.NewCodeGenRequest(file1, file2)

	if len(req.ProtoFile) != 2 {
		t.Fatalf("request has %d proto files, want 2", len(req.ProtoFile))
	}
	if len(req.FileToGenerate) != 2 {
		t.Fatalf("request has %d files to generate, want 2", len(req.FileToGenerate))
	}
	if req.FileToGenerate[0] != "file1.proto" {
		t.Errorf("first file to generate = %q, want %q", req.FileToGenerate[0], "file1.proto")
	}
	if req.FileToGenerate[1] != "file2.proto" {
		t.Errorf("second file to generate = %q, want %q", req.FileToGenerate[1], "file2.proto")
	}
}

func TestNewEnum(t *testing.T) {
	name := "Status"
	value1 := testutils.NewEnumValue("STATUS_UNKNOWN", 0)
	value2 := testutils.NewEnumValue("STATUS_ACTIVE", 1)
	value3 := testutils.NewEnumValue("STATUS_INACTIVE", 2)

	enum := testutils.NewEnum(name, value1, value2, value3)

	if enum.GetName() != name {
		t.Errorf("enum name = %q, want %q", enum.GetName(), name)
	}
	if len(enum.Value) != 3 {
		t.Fatalf("enum has %d values, want 3", len(enum.Value))
	}
	if enum.Value[0].GetName() != "STATUS_UNKNOWN" {
		t.Errorf("first value name = %q, want %q", enum.Value[0].GetName(), "STATUS_UNKNOWN")
	}
	if enum.Value[1].GetNumber() != 1 {
		t.Errorf("second value number = %d, want %d", enum.Value[1].GetNumber(), 1)
	}
}

func TestNewEnumValue(t *testing.T) {
	name := "STATUS_ACTIVE"
	number := int32(1)

	value := testutils.NewEnumValue(name, number)

	if value.GetName() != name {
		t.Errorf("enum value name = %q, want %q", value.GetName(), name)
	}
	if value.GetNumber() != number {
		t.Errorf("enum value number = %d, want %d", value.GetNumber(), number)
	}
}

// newEnumFieldTestCase represents a test case for NewEnumField
type newEnumFieldTestCase struct {
	name     string
	field    string
	typeName string
	number   int32
}

// test runs the NewEnumField test case
func (tc newEnumFieldTestCase) test(t *testing.T) {
	t.Helper()

	field := testutils.NewEnumField(tc.field, tc.number, tc.typeName)

	// Verify field properties
	testutils.AssertEqual(t, field.GetName(), tc.field, "field name")
	testutils.AssertEqual(t, field.GetNumber(), tc.number, "field number")
	testutils.AssertEqual(t, field.GetType(), descriptorpb.FieldDescriptorProto_TYPE_ENUM, "field type")
	testutils.AssertEqual(t, field.GetTypeName(), tc.typeName, "type name")
	testutils.AssertEqual(t, field.GetLabel(), descriptorpb.FieldDescriptorProto_LABEL_OPTIONAL, "field label")
}

func TestNewEnumField(t *testing.T) {
	tests := []newEnumFieldTestCase{
		{
			name:     "basic enum field",
			field:    "status",
			number:   1,
			typeName: ".test.Status",
		},
		{
			name:     "enum field with package",
			field:    "priority",
			number:   2,
			typeName: ".example.Priority",
		},
		{
			name:     "enum field with nested type",
			field:    "type",
			number:   3,
			typeName: ".test.Message.Type",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, tt.test)
	}
}

func TestNewPlugin(t *testing.T) {
	t.Run("success", testNewPluginSuccess)
	t.Run("error with nil descriptor", testNewPluginErrorNilDescriptor)
}

func testNewPluginSuccess(t *testing.T) {
	file := testutils.NewFileDescriptor("test.proto", "test.v1", "test/v1")

	// Add a message to make it more realistic
	msg := testutils.NewMessage("User",
		testutils.NewField("id", 1, descriptorpb.FieldDescriptorProto_TYPE_STRING),
	)
	file.MessageType = append(file.MessageType, msg)

	plugin, err := testutils.NewPlugin(t, file)
	if err != nil {
		t.Fatalf("failed to create plugin: %v", err)
	}

	if plugin == nil {
		t.Fatal("plugin is nil")
	}
	if len(plugin.Files) != 1 {
		t.Fatalf("plugin has %d files, want 1", len(plugin.Files))
	}
	if plugin.Files[0].Desc.Path() != "test.proto" {
		t.Errorf("file name = %q, want %q", plugin.Files[0].Desc.Path(), "test.proto")
	}
}

func testNewPluginErrorNilDescriptor(t *testing.T) {
	// Create an invalid file descriptor that will cause protogen.New to fail
	plugin, err := testutils.NewPlugin(t, nil)
	if err == nil {
		t.Fatal("expected error for nil file descriptor, got nil")
	}
	if plugin != nil {
		t.Fatal("expected nil plugin for error case")
	}
}

func TestRunGenerator(t *testing.T) {
	// Create a simple test file
	file := testutils.NewFileDescriptor("test.proto", "test.v1", "test/v1")

	// Add the message types that the service references
	getItemRequest := testutils.NewMessage("GetItemRequest",
		testutils.NewField("id", 1, descriptorpb.FieldDescriptorProto_TYPE_STRING),
	)
	item := testutils.NewMessage("Item",
		testutils.NewField("id", 1, descriptorpb.FieldDescriptorProto_TYPE_STRING),
		testutils.NewField("name", 2, descriptorpb.FieldDescriptorProto_TYPE_STRING),
	)
	file.MessageType = append(file.MessageType, getItemRequest, item)

	service := testutils.NewService("TestService",
		testutils.NewMethod("GetItem", ".test.v1.GetItemRequest", ".test.v1.Item"),
	)
	file.Service = append(file.Service, service)

	// Create a simple generator function
	var generatorCalled bool
	genFunc := func(plugin *protogen.Plugin) error {
		generatorCalled = true
		// Generate a simple file
		g := plugin.NewGeneratedFile("test.pb.go", "test/v1")
		g.P("// Generated code")
		g.P("package v1")
		return nil
	}

	// Run the generator
	req := testutils.NewCodeGenRequest(file)
	response := testutils.RunGenerator(t, req, genFunc)

	if !generatorCalled {
		t.Error("generator function was not called")
	}
	if len(response.File) != 1 {
		t.Fatalf("response has %d files, want 1", len(response.File))
	}
	if response.File[0].GetName() != "test.pb.go" {
		t.Errorf("generated file name = %q, want %q", response.File[0].GetName(), "test.pb.go")
	}
}

// getGeneratedFileContentTestCase represents a test case for GetGeneratedFileContent
type getGeneratedFileContentTestCase struct {
	setupPlugin func(t *testing.T) *protogen.Plugin
	name        string
	filename    string
	expectFound bool
}

// test runs the GetGeneratedFileContent test case
func (tc getGeneratedFileContentTestCase) test(t *testing.T) {
	t.Helper()

	plugin := tc.setupPlugin(t)

	content, ok := testutils.GetGeneratedFileContent(t, plugin, tc.filename)
	testutils.AssertEqual(t, ok, tc.expectFound, "file found")

	if tc.expectFound && tc.filename == "test.go" {
		expectedContent := "package test\n\n// Generated content\n"
		testutils.AssertEqual(t, content, expectedContent, "file content")
	}
}

// setupPluginWithGeneratedFile creates a plugin with a single generated file "test.go"
// containing basic package declaration and a comment.
func setupPluginWithGeneratedFile(t *testing.T) *protogen.Plugin {
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

	// Simulate generating a file
	genFile := plugin.NewGeneratedFile("test.go", protogen.GoImportPath(protoFile.GetOptions().GetGoPackage()))
	genFile.P("package test")
	genFile.P("")
	genFile.P("// Generated content")

	return plugin
}

// newGeneratedFileContentTestCase creates a new test case for GetGeneratedFileContent
func newGeneratedFileContentTestCase(
	name, filename string, expectFound bool,
	setupPlugin func(t *testing.T) *protogen.Plugin,
) getGeneratedFileContentTestCase {
	return getGeneratedFileContentTestCase{
		setupPlugin: setupPlugin,
		name:        name,
		filename:    filename,
		expectFound: expectFound,
	}
}

func TestGetGeneratedFileContent(t *testing.T) {
	setupPlugin := setupPluginWithGeneratedFile
	tests := []getGeneratedFileContentTestCase{
		newGeneratedFileContentTestCase("existing file", "test.go", true, setupPlugin),
		newGeneratedFileContentTestCase("non-existent file", "non-existent.go", false, setupPlugin),
		newGeneratedFileContentTestCase("empty filename", "", false, setupPlugin),
	}

	for _, tt := range tests {
		t.Run(tt.name, tt.test)
	}
}

// multipleFilesTestCase represents a test case for multiple file generation
type multipleFilesTestCase struct {
	name    string
	content string
}

// test runs the multiple files test case
func (tc multipleFilesTestCase) test(t *testing.T, plugin *protogen.Plugin) {
	t.Helper()

	content, ok := testutils.GetGeneratedFileContent(t, plugin, tc.name)
	if !ok {
		t.Fatalf("Expected to find file %s", tc.name)
	}
	expectedContent := tc.content + "\n"
	testutils.AssertEqual(t, content, expectedContent, "file content")
}

func TestGetGeneratedFileContent_MultipleFiles(t *testing.T) {
	// Create a proto file
	protoFile := testutils.NewFileDescriptor("test.proto", "test", "github.com/example/test")

	// Create plugin
	plugin, err := testutils.NewPlugin(t, protoFile)
	if err != nil {
		t.Fatalf("failed to create plugin: %v", err)
	}

	// Define test cases
	files := []multipleFilesTestCase{
		{"file1.go", "package test\n\n// File 1"},
		{"file2.go", "package test\n\n// File 2"},
		{"file3.go", "package test\n\n// File 3"},
	}

	// Generate files
	for _, f := range files {
		genFile := plugin.NewGeneratedFile(f.name, protogen.GoImportPath(protoFile.GetOptions().GetGoPackage()))
		genFile.P(f.content)
	}

	// Test retrieving each file
	for _, f := range files {
		t.Run(f.name, func(t *testing.T) {
			f.test(t, plugin)
		})
	}
}

func TestCompleteExample(t *testing.T) {
	// This test demonstrates a complete example of using the factory functions
	// to create a proto file with messages, enums, and services

	// Create the file descriptor
	file := testutils.NewFileDescriptor("api.proto", "api.v1", "github.com/example/api/v1")

	// Create an enum
	statusEnum := testutils.NewEnum("UserStatus",
		testutils.NewEnumValue("USER_STATUS_UNKNOWN", 0),
		testutils.NewEnumValue("USER_STATUS_ACTIVE", 1),
		testutils.NewEnumValue("USER_STATUS_INACTIVE", 2),
	)
	file.EnumType = append(file.EnumType, statusEnum)

	// Create message types
	userMsg := testutils.NewMessage("User",
		testutils.NewField("id", 1, descriptorpb.FieldDescriptorProto_TYPE_STRING),
		testutils.NewField("name", 2, descriptorpb.FieldDescriptorProto_TYPE_STRING),
		testutils.NewField("email", 3, descriptorpb.FieldDescriptorProto_TYPE_STRING),
		testutils.NewField("status", 4, descriptorpb.FieldDescriptorProto_TYPE_ENUM),
	)
	file.MessageType = append(file.MessageType, userMsg)

	getUserReq := testutils.NewMessage("GetUserRequest",
		testutils.NewField("id", 1, descriptorpb.FieldDescriptorProto_TYPE_STRING),
	)
	file.MessageType = append(file.MessageType, getUserReq)

	listUsersReq := testutils.NewMessage("ListUsersRequest",
		testutils.NewField("page_size", 1, descriptorpb.FieldDescriptorProto_TYPE_INT32),
		testutils.NewField("page_token", 2, descriptorpb.FieldDescriptorProto_TYPE_STRING),
	)
	file.MessageType = append(file.MessageType, listUsersReq)

	listUsersResponse := testutils.NewMessage("ListUsersResponse",
		testutils.NewField("users", 1, descriptorpb.FieldDescriptorProto_TYPE_MESSAGE),
		testutils.NewField("next_page_token", 2, descriptorpb.FieldDescriptorProto_TYPE_STRING),
	)
	file.MessageType = append(file.MessageType, listUsersResponse)

	// Create a service
	userService := testutils.NewService("UserService",
		testutils.NewMethod("GetUser", ".api.v1.GetUserRequest", ".api.v1.User"),
		testutils.NewMethod("ListUsers", ".api.v1.ListUsersRequest", ".api.v1.ListUsersResponse"),
	)
	file.Service = append(file.Service, userService)

	// Verify the structure
	if len(file.EnumType) != 1 {
		t.Errorf("file has %d enums, want 1", len(file.EnumType))
	}
	if len(file.MessageType) != 4 {
		t.Errorf("file has %d messages, want 4", len(file.MessageType))
	}
	if len(file.Service) != 1 {
		t.Errorf("file has %d services, want 1", len(file.Service))
	}
	if len(file.Service[0].Method) != 2 {
		t.Errorf("service has %d methods, want 2", len(file.Service[0].Method))
	}
}
