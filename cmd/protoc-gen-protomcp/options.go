package main

import "google.golang.org/protobuf/compiler/protogen"

const (
	// DefaultGenerateInterfaces controls whether to generate interfaces for messages by default
	DefaultGenerateInterfaces = true
	// DefaultGenerateServices controls whether to generate service interfaces by default
	DefaultGenerateServices = true
	// DefaultGenerateNoImpl controls whether to generate NoImpl structs by default
	DefaultGenerateNoImpl = true
	// DefaultGenerateEnums controls whether to generate enum helper methods by default
	DefaultGenerateEnums = true
	// DefaultInterfacePattern is the default pattern for interface names
	DefaultInterfacePattern = "I%"
	// DefaultEnumPattern is the default pattern for enum type names
	DefaultEnumPattern = "%Enum"
)

// GeneratorOptions controls code generation
type GeneratorOptions struct {
	InterfacePattern   string // Pattern for interface names, e.g., "I%" or "%Interface"
	EnumPattern        string // Pattern for enum type names, e.g., "%Enum" or "E%"
	GenerateInterfaces bool
	GenerateServices   bool
	GenerateNoImpl     bool // Generate NoImpl structs for interfaces
	GenerateEnums      bool // Generate enum helper methods
}

// GetInterfacePattern returns the interface pattern, defaulting to DefaultInterfacePattern if not set
func (o *GeneratorOptions) GetInterfacePattern() string {
	if o == nil || o.InterfacePattern == "" {
		return DefaultInterfacePattern
	}
	return o.InterfacePattern
}

// GetGenerateInterfaces returns whether to generate interfaces, defaulting to DefaultGenerateInterfaces
func (o *GeneratorOptions) GetGenerateInterfaces() bool {
	if o == nil {
		return DefaultGenerateInterfaces
	}
	return o.GenerateInterfaces
}

// GetGenerateServices returns whether to generate services, defaulting to DefaultGenerateServices
func (o *GeneratorOptions) GetGenerateServices() bool {
	if o == nil {
		return DefaultGenerateServices
	}
	return o.GenerateServices
}

// GetGenerateNoImpl returns whether to generate NoImpl structs, defaulting to DefaultGenerateNoImpl
func (o *GeneratorOptions) GetGenerateNoImpl() bool {
	if o == nil {
		return DefaultGenerateNoImpl
	}
	return o.GenerateNoImpl
}

// GetGenerateEnums returns whether to generate enum helper methods, defaulting to DefaultGenerateEnums
func (o *GeneratorOptions) GetGenerateEnums() bool {
	if o == nil {
		return DefaultGenerateEnums
	}
	return o.GenerateEnums
}

// GetEnumPattern returns the enum pattern, defaulting to DefaultEnumPattern if not set
func (o *GeneratorOptions) GetEnumPattern() string {
	if o == nil || o.EnumPattern == "" {
		return DefaultEnumPattern
	}
	return o.EnumPattern
}

// HasMessages returns true if the file has message definitions
func (*GeneratorOptions) HasMessages(file *protogen.File) bool {
	return file != nil && len(file.Messages) > 0
}

// HasServices returns true if the file has service definitions
func (*GeneratorOptions) HasServices(file *protogen.File) bool {
	return file != nil && len(file.Services) > 0
}

// HasEnums returns true if the file has enum definitions
func (*GeneratorOptions) HasEnums(file *protogen.File) bool {
	return file != nil && len(file.Enums) > 0
}

// NeedsMessages returns true if the file needs message interface generation
func (o *GeneratorOptions) NeedsMessages(file *protogen.File) bool {
	return o.HasMessages(file) && o.GetGenerateInterfaces()
}

// NeedsServices returns true if the file needs service interface generation
func (o *GeneratorOptions) NeedsServices(file *protogen.File) bool {
	return o.HasServices(file) && o.GetGenerateServices()
}

// NeedsTypes returns true if the file needs type generation
func (o *GeneratorOptions) NeedsTypes(file *protogen.File) bool {
	return o.NeedsMessages(file) || o.NeedsServices(file) || o.NeedsEnums(file)
}

// NeedsNoImpl returns true if the file needs NoImpl generation
func (o *GeneratorOptions) NeedsNoImpl(file *protogen.File) bool {
	// NoImpl is only for messages and services, not for enums
	return (o.NeedsMessages(file) || o.NeedsServices(file)) && o.GetGenerateNoImpl()
}

// NeedsEnums returns true if the file needs enum helper generation
func (o *GeneratorOptions) NeedsEnums(file *protogen.File) bool {
	return o.HasEnums(file) && o.GetGenerateEnums()
}
