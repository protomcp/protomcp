package main

import (
	"testing"

	"google.golang.org/protobuf/compiler/protogen"

	"protomcp.org/protomcp/pkg/generator/testutils"
)

func TestGetInterfacePattern(t *testing.T) {
	testNilOptionsReturnsDefault := func(t *testing.T) {
		var opts *GeneratorOptions
		got := opts.GetInterfacePattern()
		testutils.AssertEqual(t, got, DefaultInterfacePattern, "GetInterfacePattern()")
	}

	testEmptyPatternReturnsDefault := func(t *testing.T) {
		opts := &GeneratorOptions{InterfacePattern: ""}
		got := opts.GetInterfacePattern()
		testutils.AssertEqual(t, got, DefaultInterfacePattern, "GetInterfacePattern()")
	}

	testCustomPattern := func(t *testing.T) {
		opts := &GeneratorOptions{InterfacePattern: "%Interface"}
		got := opts.GetInterfacePattern()
		testutils.AssertEqual(t, got, "%Interface", "GetInterfacePattern()")
	}

	t.Run("nil options returns default", testNilOptionsReturnsDefault)
	t.Run("empty pattern returns default", testEmptyPatternReturnsDefault)
	t.Run("custom pattern", testCustomPattern)
}

func TestGettersWithNilOptions(t *testing.T) {
	testGetGenerateInterfaces := func(t *testing.T) {
		var opts *GeneratorOptions
		got := opts.GetGenerateInterfaces()
		testutils.AssertEqual(t, got, DefaultGenerateInterfaces, "GetGenerateInterfaces() with nil")
	}

	testGetGenerateServices := func(t *testing.T) {
		var opts *GeneratorOptions
		got := opts.GetGenerateServices()
		testutils.AssertEqual(t, got, DefaultGenerateServices, "GetGenerateServices() with nil")
	}

	testGetGenerateNoImpl := func(t *testing.T) {
		var opts *GeneratorOptions
		got := opts.GetGenerateNoImpl()
		testutils.AssertEqual(t, got, DefaultGenerateNoImpl, "GetGenerateNoImpl() with nil")
	}

	t.Run("GetGenerateInterfaces", testGetGenerateInterfaces)
	t.Run("GetGenerateServices", testGetGenerateServices)
	t.Run("GetGenerateNoImpl", testGetGenerateNoImpl)
}

func TestHasMethodsWithNilFile(t *testing.T) {
	testHasMessagesWithNilFile := func(t *testing.T) {
		opts := &GeneratorOptions{}
		got := opts.HasMessages(nil)
		testutils.AssertEqual(t, got, false, "HasMessages(nil)")
	}

	testHasServicesWithNilFile := func(t *testing.T) {
		opts := &GeneratorOptions{}
		got := opts.HasServices(nil)
		testutils.AssertEqual(t, got, false, "HasServices(nil)")
	}

	t.Run("HasMessages with nil file", testHasMessagesWithNilFile)
	t.Run("HasServices with nil file", testHasServicesWithNilFile)
}

func TestNeedsMethodsCoverage(t *testing.T) {
	testNeedsTypesWithOnlyMessages := func(t *testing.T) {
		opts := &GeneratorOptions{
			GenerateInterfaces: true,
			GenerateServices:   false,
		}

		file := &protogen.File{
			Messages: []*protogen.Message{{}}, // One message
			Services: []*protogen.Service{},   // No services
		}

		got := opts.NeedsTypes(file)
		testutils.AssertEqual(t, got, true, "NeedsTypes with only messages")
	}

	testNeedsTypesWithOnlyServices := func(t *testing.T) {
		opts := &GeneratorOptions{
			GenerateInterfaces: false,
			GenerateServices:   true,
		}

		file := &protogen.File{
			Messages: []*protogen.Message{},   // No messages
			Services: []*protogen.Service{{}}, // One service
		}

		got := opts.NeedsTypes(file)
		testutils.AssertEqual(t, got, true, "NeedsTypes with only services")
	}

	testNeedsNoImplWithFalseGenerateNoImpl := func(t *testing.T) {
		opts := &GeneratorOptions{
			GenerateInterfaces: true,
			GenerateNoImpl:     false,
		}

		file := &protogen.File{
			Messages: []*protogen.Message{{}}, // Has content
		}

		got := opts.NeedsNoImpl(file)
		testutils.AssertEqual(t, got, false, "NeedsNoImpl with GenerateNoImpl=false")
	}

	t.Run("NeedsTypes with only messages", testNeedsTypesWithOnlyMessages)
	t.Run("NeedsTypes with only services", testNeedsTypesWithOnlyServices)
	t.Run("NeedsNoImpl with false GenerateNoImpl", testNeedsNoImplWithFalseGenerateNoImpl)
}
