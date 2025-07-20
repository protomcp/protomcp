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

// GenerateFile generates Go code for a single proto file
func (g *Generator) GenerateFile(file *protogen.File, opts *GeneratorOptions) error {
	if file == nil {
		return errors.New("file cannot be nil")
	}

	if opts.NeedsTypes(file) {
		if err := g.generateTypesWithTemplate(file, opts); err != nil {
			return fmt.Errorf("failed to generate types for %s: %w", file.Desc.Path(), err)
		}
	}

	if opts.NeedsNoImpl(file) {
		if err := g.generateNoImplWithTemplate(file, opts); err != nil {
			return fmt.Errorf("failed to generate noImpl for %s: %w", file.Desc.Path(), err)
		}
	}

	return nil
}
