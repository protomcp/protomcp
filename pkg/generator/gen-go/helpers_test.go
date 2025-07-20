package gengo

import (
	"testing"

	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/types/descriptorpb"

	"protomcp.org/protomcp/pkg/generator/testutils"
)

// interfaceNameTestCase represents a test case for interface naming
type interfaceNameTestCase struct {
	name           string
	typeName       string
	pattern        string
	expectedResult string
}

// noImplTestCase represents a test case for NoImpl naming
type noImplTestCase struct {
	name           string
	typeName       string
	pattern        string
	expectedNoImpl string
	isService      bool
}

// newServiceTestCase creates a test case for service interface naming
func newServiceTestCase(name, serviceName, pattern, expected string) interfaceNameTestCase {
	return interfaceNameTestCase{
		name:           name,
		typeName:       serviceName,
		pattern:        pattern,
		expectedResult: expected,
	}
}

// newMessageTestCase creates a test case for message interface naming
func newMessageTestCase(name, messageName, pattern, expected string) interfaceNameTestCase {
	return interfaceNameTestCase{
		name:           name,
		typeName:       messageName,
		pattern:        pattern,
		expectedResult: expected,
	}
}

// newNoImplTestCase creates a test case for NoImpl naming
func newNoImplTestCase(name, typeName, pattern, expected string, isService bool) noImplTestCase {
	return noImplTestCase{
		name:           name,
		typeName:       typeName,
		pattern:        pattern,
		expectedNoImpl: expected,
		isService:      isService,
	}
}

// runTestInterfaceNameForService runs a single test case for service interface naming
func runTestInterfaceNameForService(t *testing.T, tt interfaceNameTestCase) {
	t.Helper()

	// Create proto file with service
	protoFile := testutils.NewFileDescriptor("test.proto", "test", "github.com/example/test")
	serviceProto := &descriptorpb.ServiceDescriptorProto{
		Name: &tt.typeName,
	}
	protoFile.Service = []*descriptorpb.ServiceDescriptorProto{serviceProto}

	// Create request and run through protogen
	req := testutils.NewCodeGenRequest(protoFile)
	var result string

	testutils.RunGenerator(t, req, func(plugin *protogen.Plugin) error {
		for _, file := range plugin.Files {
			for _, svc := range file.Services {
				result = InterfaceNameForService(svc, tt.pattern)
			}
		}
		return nil
	})

	if result != tt.expectedResult {
		t.Errorf("InterfaceNameForService(%q, %q) = %q, want %q",
			tt.typeName, tt.pattern, result, tt.expectedResult)
	}
}

func TestInterfaceNameForService(t *testing.T) {
	tests := []interfaceNameTestCase{
		newServiceTestCase(
			"Service suffix already present with I% pattern",
			"CalculatorService",
			"I%",
			"ICalculatorService",
		),
		newServiceTestCase(
			"No Service suffix with I% pattern",
			"Calculator",
			"I%",
			"ICalculatorService",
		),
		newServiceTestCase(
			"Service suffix already present with %Interface pattern",
			"CalculatorService",
			"%Interface",
			"CalculatorServiceInterface",
		),
		newServiceTestCase(
			"No Service suffix with %Interface pattern",
			"Calculator",
			"%Interface",
			"CalculatorServiceInterface",
		),
		newServiceTestCase(
			"Empty pattern defaults to I% with Service suffix",
			"Calculator",
			"",
			"ICalculatorService",
		),
		newServiceTestCase(
			"Empty pattern with existing Service suffix",
			"CalculatorService",
			"",
			"ICalculatorService",
		),
		newServiceTestCase(
			"Multiple Service suffixes",
			"ServiceService",
			"I%",
			"IServiceService",
		),
		newServiceTestCase(
			"Service as complete name",
			"Service",
			"I%",
			"IService",
		),
		newServiceTestCase(
			"Bare % pattern with Calculator",
			"Calculator",
			"%",
			"CalculatorService",
		),
		newServiceTestCase(
			"Bare % pattern with existing Service suffix",
			"CalculatorService",
			"%",
			"CalculatorService",
		),
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			runTestInterfaceNameForService(t, tt)
		})
	}
}

// runTestInterfaceNameForMessage runs a single test case for message interface naming
func runTestInterfaceNameForMessage(t *testing.T, tt interfaceNameTestCase) {
	t.Helper()

	// Create proto file with message
	protoFile := testutils.NewFileDescriptor("test.proto", "test", "github.com/example/test")
	protoFile.MessageType = []*descriptorpb.DescriptorProto{
		testutils.NewMessage(tt.typeName),
	}

	// Create request and run through protogen
	req := testutils.NewCodeGenRequest(protoFile)
	var result string

	testutils.RunGenerator(t, req, func(plugin *protogen.Plugin) error {
		for _, file := range plugin.Files {
			for _, msg := range file.Messages {
				result = InterfaceNameForMessage(msg, tt.pattern)
			}
		}
		return nil
	})

	if result != tt.expectedResult {
		t.Errorf("InterfaceNameForMessage(%q, %q) = %q, want %q",
			tt.typeName, tt.pattern, result, tt.expectedResult)
	}
}

func TestInterfaceNameForMessage(t *testing.T) {
	tests := []interfaceNameTestCase{
		newMessageTestCase(
			"I% pattern",
			"AddRequest",
			"I%",
			"IAddRequest",
		),
		newMessageTestCase(
			"%Interface pattern",
			"AddRequest",
			"%Interface",
			"AddRequestInterface",
		),
		newMessageTestCase(
			"Empty pattern defaults to I%",
			"AddRequest",
			"",
			"IAddRequest",
		),
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			runTestInterfaceNameForMessage(t, tt)
		})
	}
}

// createProtoFileForTest creates a proto file descriptor for testing
// revive:disable-next-line:flag-parameter
func createProtoFileForTest(typeName string, isService bool) *descriptorpb.FileDescriptorProto {
	protoFile := testutils.NewFileDescriptor("test.proto", "test", "github.com/example/test")

	if isService {
		serviceProto := &descriptorpb.ServiceDescriptorProto{
			Name: &typeName,
		}
		protoFile.Service = []*descriptorpb.ServiceDescriptorProto{serviceProto}
	} else {
		protoFile.MessageType = []*descriptorpb.DescriptorProto{
			testutils.NewMessage(typeName),
		}
	}

	return protoFile
}

// extractNoImplNameForService extracts the NoImpl name for a service from the plugin
func extractNoImplNameForService(plugin *protogen.Plugin) string {
	for _, file := range plugin.Files {
		if len(file.Services) > 0 {
			return "NoImpl" + file.Services[0].GoName
		}
	}
	return ""
}

// extractNoImplNameForMessage extracts the NoImpl name for a message from the plugin
func extractNoImplNameForMessage(plugin *protogen.Plugin) string {
	for _, file := range plugin.Files {
		if len(file.Messages) > 0 {
			return "NoImpl" + file.Messages[0].GoIdent.GoName
		}
	}
	return ""
}

// runTestNoImplNaming runs a single test case for NoImpl naming
func runTestNoImplNaming(t *testing.T, tt noImplTestCase) {
	t.Helper()

	protoFile := createProtoFileForTest(tt.typeName, tt.isService)
	req := testutils.NewCodeGenRequest(protoFile)

	var actualNoImpl string
	testutils.RunGenerator(t, req, func(plugin *protogen.Plugin) error {
		if tt.isService {
			actualNoImpl = extractNoImplNameForService(plugin)
		} else {
			actualNoImpl = extractNoImplNameForMessage(plugin)
		}
		return nil
	})

	if actualNoImpl != tt.expectedNoImpl {
		typeDesc := "message"
		if tt.isService {
			typeDesc = "service"
		}
		t.Errorf("NoImpl name for %s %q = %q, want %q",
			typeDesc, tt.typeName, actualNoImpl, tt.expectedNoImpl)
	}
}

func TestNoImplNaming(t *testing.T) {
	tests := []noImplTestCase{
		newNoImplTestCase(
			"Message with I% pattern",
			"AddRequest",
			"I%",
			"NoImplAddRequest",
			false,
		),
		newNoImplTestCase(
			"Service with I% pattern",
			"Calculator",
			"I%",
			"NoImplCalculator",
			true,
		),
		newNoImplTestCase(
			"Service with Service suffix",
			"CalculatorService",
			"I%",
			"NoImplCalculatorService",
			true,
		),
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			runTestNoImplNaming(t, tt)
		})
	}
}
