package main

import (
	"bytes"
	"fmt"
	"strings"

	"google.golang.org/protobuf/compiler/protogen"

	gengo "protomcp.org/protomcp/pkg/generator/gen-go"
)

// generateTypesWithTemplate generates interface definitions using templates
func (g *Generator) generateTypesWithTemplate(file *protogen.File, opts *GeneratorOptions) error {
	filename := file.GeneratedFilenamePrefix + ".types.go"
	genFile := g.plugin.NewGeneratedFile(filename, file.GoImportPath)

	// Prepare template data
	data := g.prepareTemplateData(file, opts)

	// Process messages
	if opts.GenerateInterfaces {
		g.processMessages(genFile, file.Messages, opts, data)
	}

	// Process services
	if opts.GenerateServices {
		g.processServices(genFile, file.Services, opts, data)
	}

	// Render template
	var buf bytes.Buffer
	if err := templates.ExecuteTemplate(&buf, fileTemplate, data); err != nil {
		return fmt.Errorf("failed to execute template: %w", err)
	}

	// Write to generated file
	genFile.P(buf.String())

	return nil
}

// prepareTemplateData creates the base template data structure
func (g *Generator) prepareTemplateData(file *protogen.File, opts *GeneratorOptions) *TemplateData {
	// Collect imports as paths
	var stdPaths []string
	var thirdPartyPaths []string

	// Add context import if we need services (standard library)
	if g.NeedsServices(file, opts) {
		stdPaths = append(stdPaths, "context")
	}

	// Use GenerateImports to process and organize imports
	importGroups := gengo.GenerateImports(stdPaths, thirdPartyPaths)

	return &TemplateData{
		Package:      string(file.GoPackageName),
		SourceFile:   file.Desc.Path(),
		ImportGroups: importGroups,
		Messages:     make([]MessageData, 0, len(file.Messages)),
		Services:     make([]ServiceData, 0, len(file.Services)),
	}
}

// processMessages adds all messages to template data
func (g *Generator) processMessages(gen *protogen.GeneratedFile, messages []*protogen.Message,
	opts *GeneratorOptions, data *TemplateData) {
	for _, msg := range messages {
		if msg.Desc.IsMapEntry() {
			continue // Skip map entry messages
		}
		msgData := g.buildMessageData(gen, msg, opts)
		data.Messages = append(data.Messages, msgData)
	}
}

// processServices adds all services to template data
func (g *Generator) processServices(gen *protogen.GeneratedFile, services []*protogen.Service,
	opts *GeneratorOptions, data *TemplateData) {
	for _, svc := range services {
		svcData := g.buildServiceData(gen, svc, opts)
		data.Services = append(data.Services, svcData)
	}
}

// buildMessageData builds template data for a message
func (*Generator) buildMessageData(gen *protogen.GeneratedFile, msg *protogen.Message,
	opts *GeneratorOptions) MessageData {
	data := MessageData{
		Name:          msg.GoIdent.GoName,
		InterfaceName: gengo.InterfaceNameForMessage(msg, opts.InterfacePattern),
		Comment:       strings.TrimSpace(string(msg.Comments.Leading)),
		Fields:        make([]FieldData, 0, len(msg.Fields)),
		OneOfGroups:   make([]OneOfData, 0, len(OneOfGroups(msg))),
	}

	// Process regular fields
	processRegularFields(gen, msg, opts, &data)

	// Process `oneof` groups
	processOneOfGroups(gen, msg, opts, &data)

	return data
}

// processRegularFields adds regular (non-oneof) fields to message data
func processRegularFields(gen *protogen.GeneratedFile, msg *protogen.Message,
	opts *GeneratorOptions, data *MessageData) {
	for _, field := range msg.Fields {
		if field.Oneof != nil && !gengo.IsSyntheticOneOf(field.Oneof) {
			continue // Skip fields that are part of real `oneof` groups
		}

		fieldData := FieldData{
			Name:     string(field.Desc.Name()),
			GoName:   field.GoName,
			Type:     gengo.GoTypeForFieldWithPattern(gen, field, opts.InterfacePattern),
			Comment:  strings.TrimSpace(string(field.Comments.Leading)),
			Optional: field.Desc.HasOptionalKeyword(),
			IsOneOf:  false,
		}
		data.Fields = append(data.Fields, fieldData)
	}
}

// OneOfGroups returns the oneof groups for a message
func OneOfGroups(msg *protogen.Message) []*protogen.Oneof {
	// cspell:disable-next-line
	return msg.Oneofs // Using protogen API field name
}

// processOneOfGroups adds oneof groups to message data
func processOneOfGroups(gen *protogen.GeneratedFile, msg *protogen.Message, opts *GeneratorOptions, data *MessageData) {
	for _, oneof := range OneOfGroups(msg) {
		if gengo.IsSyntheticOneOf(oneof) {
			continue // Skip synthetic `oneof` groups
		}

		oneOfData := OneOfData{
			Name:   string(oneof.Desc.Name()),
			GoName: oneof.GoName,
			Fields: make([]OneOfFieldData, 0, len(oneof.Fields)),
		}

		for _, field := range oneof.Fields {
			fieldData := OneOfFieldData{
				Name:      string(field.Desc.Name()),
				GoName:    field.GoName,
				Type:      gengo.GoTypeForFieldWithPattern(gen, field, opts.InterfacePattern),
				OneOfName: string(oneof.Desc.Name()),
			}
			oneOfData.Fields = append(oneOfData.Fields, fieldData)
		}

		data.OneOfGroups = append(data.OneOfGroups, oneOfData)
	}
}

// buildServiceData builds template data for a service
func (*Generator) buildServiceData(_ *protogen.GeneratedFile, svc *protogen.Service,
	opts *GeneratorOptions) ServiceData {
	data := ServiceData{
		Name:          svc.GoName,
		InterfaceName: gengo.InterfaceNameForService(svc, opts.InterfacePattern),
		Comment:       strings.TrimSpace(string(svc.Comments.Leading)),
		Methods:       make([]MethodData, 0, len(svc.Methods)),
	}

	for _, method := range svc.Methods {
		methodData := MethodData{
			Name:         method.GoName,
			Comment:      strings.TrimSpace(string(method.Comments.Leading)),
			RequestType:  gengo.InterfaceNameForMessage(method.Input, opts.InterfacePattern),
			ResponseType: gengo.InterfaceNameForMessage(method.Output, opts.InterfacePattern),
		}
		data.Methods = append(data.Methods, methodData)
	}

	return data
}
