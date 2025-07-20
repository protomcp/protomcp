package main

import (
	"strings"
	"testing"

	"google.golang.org/protobuf/types/descriptorpb"

	"protomcp.org/protomcp/pkg/generator/testutils"
)

// Setup functions for common test scenarios

// setupFileWithMessageField creates a proto file with a User message containing a Profile message field
func setupFileWithMessageField() *descriptorpb.FileDescriptorProto {
	protoFile := testutils.NewFileDescriptor("test.proto", "test", "github.com/example/test")

	// Define Profile message
	profileMsg := testutils.NewMessage("Profile",
		testutils.NewField("bio", 1, descriptorpb.FieldDescriptorProto_TYPE_STRING),
		testutils.NewField("avatar_url", 2, descriptorpb.FieldDescriptorProto_TYPE_STRING),
	)

	// Define User message with Profile field
	userField := testutils.NewField("profile", 2, descriptorpb.FieldDescriptorProto_TYPE_MESSAGE)
	userField.TypeName = func(s string) *string { return &s }(".test.Profile")

	userMsg := testutils.NewMessage("User",
		testutils.NewField("id", 1, descriptorpb.FieldDescriptorProto_TYPE_STRING),
		userField,
		testutils.NewField("name", 3, descriptorpb.FieldDescriptorProto_TYPE_STRING),
	)

	protoFile.MessageType = []*descriptorpb.DescriptorProto{profileMsg, userMsg}
	return protoFile
}

// setupFileWithRepeatedMessageField creates a proto file with repeated message fields
func setupFileWithRepeatedMessageField() *descriptorpb.FileDescriptorProto {
	protoFile := testutils.NewFileDescriptor("test.proto", "test", "github.com/example/test")

	// Define Address message
	addressMsg := testutils.NewMessage("Address",
		testutils.NewField("street", 1, descriptorpb.FieldDescriptorProto_TYPE_STRING),
		testutils.NewField("city", 2, descriptorpb.FieldDescriptorProto_TYPE_STRING),
	)

	// Define repeated address field
	addressesField := testutils.NewField("addresses", 3, descriptorpb.FieldDescriptorProto_TYPE_MESSAGE)
	addressesField.TypeName = func(s string) *string { return &s }(".test.Address")
	label := descriptorpb.FieldDescriptorProto_LABEL_REPEATED
	addressesField.Label = &label

	// Define Person message with repeated Address field
	personMsg := testutils.NewMessage("Person",
		testutils.NewField("name", 1, descriptorpb.FieldDescriptorProto_TYPE_STRING),
		testutils.NewField("age", 2, descriptorpb.FieldDescriptorProto_TYPE_INT32),
		addressesField,
	)

	protoFile.MessageType = []*descriptorpb.DescriptorProto{addressMsg, personMsg}
	return protoFile
}

// messageFieldTestCase represents a test case for message field types
type messageFieldTestCase struct {
	setupFile          func() *descriptorpb.FileDescriptorProto
	name               string
	messagePattern     string
	wantContains       []string
	wantNotContains    []string
	wantNoImplContains []string
	generateInterfaces bool
	generateNoImpl     bool
}

// test runs the message field test case
func (tc messageFieldTestCase) test(t *testing.T) {
	t.Helper()

	protoFile := tc.setupFile()
	req := testutils.NewCodeGenRequest(protoFile)
	opts := &GeneratorOptions{
		GenerateInterfaces: tc.generateInterfaces,
		GenerateNoImpl:     tc.generateNoImpl,
		InterfacePattern:   tc.messagePattern,
	}
	resp := runGenerator(t, req, opts)

	// Check generated files
	if !tc.generateInterfaces {
		testutils.AssertFileCount(t, resp, 0)
		return
	}

	testutils.AssertFileCount(t, resp, 2) // types.go and noimpl.go

	// Check types file
	tc.checkTypesFile(t, resp.File[0].GetContent())

	// Check noimpl file if enabled
	if tc.generateNoImpl && len(tc.wantNoImplContains) > 0 {
		tc.checkNoImplFile(t, resp.File[1].GetContent())
	}
}

// checkTypesFile validates the generated types file content
func (tc messageFieldTestCase) checkTypesFile(t *testing.T, content string) {
	t.Helper()

	for _, want := range tc.wantContains {
		testutils.AssertContains(t, content, want)
	}

	for _, notWant := range tc.wantNotContains {
		if strings.Contains(content, notWant) {
			t.Errorf("generated types file should not contain %q", notWant)
		}
	}
}

// checkNoImplFile validates the generated noimpl file content
func (tc messageFieldTestCase) checkNoImplFile(t *testing.T, content string) {
	t.Helper()

	for _, want := range tc.wantNoImplContains {
		testutils.AssertContains(t, content, want)
	}
}

// TestMessageFieldTypes tests that message fields generate interface types (not pointers)
// and that NoImpl returns nil for interfaces
func TestMessageFieldTypes(t *testing.T) {
	tests := []messageFieldTestCase{
		{
			name:               "message fields with interfaces enabled",
			setupFile:          setupFileWithMessageField,
			generateInterfaces: true,
			generateNoImpl:     true,
			messagePattern:     "I%",
			wantContains: []string{
				// Interface definitions
				"type IUser interface {",
				"type IProfile interface {",
				// Field getters return interfaces, not pointers
				"GetProfile() IProfile",
				"SetProfile(v IProfile) error",
			},
			wantNotContains: []string{
				// Should NOT have pointer-to-interface
				"GetProfile() *IProfile",
				"SetProfile(v *IProfile)",
			},
			wantNoImplContains: []string{
				// NoImpl returns nil for interface fields
				"func (*NoImplUser) GetProfile() IProfile {",
				"return nil",
			},
		},
		{
			name:               "message fields with custom pattern",
			setupFile:          setupFileWithMessageField,
			generateInterfaces: true,
			generateNoImpl:     true,
			messagePattern:     "%Interface",
			wantContains: []string{
				// Interface definitions with custom pattern
				"type UserInterface interface {",
				"type ProfileInterface interface {",
				// Field getters
				"GetProfile() ProfileInterface",
				"SetProfile(v ProfileInterface) error",
			},
			wantNotContains: []string{
				// Should NOT have pointer-to-interface
				"GetProfile() *ProfileInterface",
			},
			wantNoImplContains: []string{
				// NoImpl
				"func (*NoImplUser) GetProfile() ProfileInterface {",
				"return nil",
			},
		},
		{
			name:               "repeated message fields",
			setupFile:          setupFileWithRepeatedMessageField,
			generateInterfaces: true,
			generateNoImpl:     true,
			messagePattern:     "I%",
			wantContains: []string{
				// Check repeated message field
				"GetAddresses() []IAddress",
				"SetAddresses(v []IAddress) error",
			},
			wantNotContains: []string{
				// Should NOT have pointer in slice
				"[]*IAddress",
			},
			wantNoImplContains: []string{
				// NoImpl for repeated field
				"func (*NoImplPerson) GetAddresses() []IAddress {",
				"return nil",
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, tc.test)
	}
}
