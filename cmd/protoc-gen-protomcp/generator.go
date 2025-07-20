package main

import (
	"errors"
	"fmt"

	"google.golang.org/protobuf/compiler/protogen"
)

const (
	// SupportedFeatures declares which protobuf features this generator supports
	// Bit 1: Support for proto3 optional fields
	SupportedFeatures = 1

	// DefaultGenerateInterfaces controls whether to generate interfaces for messages by default
	DefaultGenerateInterfaces = true
	// DefaultGenerateServices controls whether to generate service interfaces by default
	DefaultGenerateServices = true
	// DefaultInterfacePattern is the default pattern for interface names
	DefaultInterfacePattern = "I%"
)

// Generator handles Go code generation for ProtoMCP
type Generator struct {
	plugin *protogen.Plugin
}

// NewGenerator creates a new Generator
func NewGenerator(plugin *protogen.Plugin) *Generator {
	if plugin == nil {
		return nil
	}

	return &Generator{
		plugin: plugin,
	}
}

// GeneratorOptions controls code generation
type GeneratorOptions struct {
	InterfacePattern   string // Pattern for interface names, e.g., "I%" or "%Interface"
	GenerateInterfaces bool
	GenerateServices   bool
	GenerateNoImpl     bool // Generate NoImpl structs for interfaces
}

// HasMessages returns true if the file has message definitions
func (*Generator) HasMessages(file *protogen.File) bool {
	return file != nil && len(file.Messages) > 0
}

// HasServices returns true if the file has service definitions
func (*Generator) HasServices(file *protogen.File) bool {
	return file != nil && len(file.Services) > 0
}

// NeedsMessages returns true if the file needs message interface generation
func (g *Generator) NeedsMessages(file *protogen.File, opts *GeneratorOptions) bool {
	return g.HasMessages(file) && opts.GenerateInterfaces
}

// NeedsServices returns true if the file needs service interface generation
func (g *Generator) NeedsServices(file *protogen.File, opts *GeneratorOptions) bool {
	return g.HasServices(file) && opts.GenerateServices
}

// NeedsTypes returns true if the file needs type generation
func (g *Generator) NeedsTypes(file *protogen.File, opts *GeneratorOptions) bool {
	return g.NeedsMessages(file, opts) || g.NeedsServices(file, opts)
}

// NeedsNoImpl returns true if the file needs NoImpl generation
func (g *Generator) NeedsNoImpl(file *protogen.File, opts *GeneratorOptions) bool {
	return opts.GenerateNoImpl && g.NeedsTypes(file, opts)
}

// normalizeOptions ensures options have valid defaults
func (*Generator) normalizeOptions(opts *GeneratorOptions) *GeneratorOptions {
	if opts == nil {
		return &GeneratorOptions{
			InterfacePattern:   DefaultInterfacePattern,
			GenerateInterfaces: DefaultGenerateInterfaces,
			GenerateServices:   DefaultGenerateServices,
		}
	}

	// Apply defaults for empty fields
	if opts.InterfacePattern == "" {
		opts.InterfacePattern = DefaultInterfacePattern
	}

	return opts
}

// GenerateFile generates Go code for a single proto file
func (g *Generator) GenerateFile(file *protogen.File, opts *GeneratorOptions) error {
	if file == nil {
		return errors.New("file cannot be nil")
	}

	opts = g.normalizeOptions(opts)

	if g.NeedsTypes(file, opts) {
		if err := g.generateTypesWithTemplate(file, opts); err != nil {
			return fmt.Errorf("failed to generate types for %s: %w", file.Desc.Path(), err)
		}
	}

	if g.NeedsNoImpl(file, opts) {
		if err := g.generateNoImplWithTemplate(file, opts); err != nil {
			return fmt.Errorf("failed to generate noImpl for %s: %w", file.Desc.Path(), err)
		}
	}

	return nil
}
