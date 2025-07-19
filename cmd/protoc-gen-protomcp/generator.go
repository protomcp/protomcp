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
}

// GenerateFile generates Go code for a single proto file
func (g *Generator) GenerateFile(file *protogen.File, opts *GeneratorOptions) error {
	if file == nil {
		return errors.New("file cannot be nil")
	}

	// Handle nil options gracefully with defaults
	if opts == nil {
		opts = &GeneratorOptions{
			InterfacePattern:   DefaultInterfacePattern,
			GenerateInterfaces: DefaultGenerateInterfaces,
			GenerateServices:   DefaultGenerateServices,
		}
	}

	// Apply defaults for empty fields
	if opts.InterfacePattern == "" {
		opts.InterfacePattern = DefaultInterfacePattern
	}

	if opts.GenerateInterfaces || opts.GenerateServices {
		if err := g.generateTypesWithTemplate(file, opts); err != nil {
			return fmt.Errorf("failed to generate types for %s: %w", file.Desc.Path(), err)
		}
	}

	return nil
}
