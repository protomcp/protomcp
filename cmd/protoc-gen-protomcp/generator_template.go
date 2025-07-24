package main

import (
	"bytes"
	"fmt"
	"strings"

	"darvaza.org/core"
	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/reflect/protoreflect"

	gengo "protomcp.org/protomcp/pkg/generator/gen-go"
)

// generateTypesWithTemplate generates interface definitions using templates
func (g *Generator) generateTypesWithTemplate(file *protogen.File, opts *GeneratorOptions) error {
	if file == nil || g == nil {
		return core.ErrInvalid
	}

	filename := file.GeneratedFilenamePrefix + ".types.go"
	genFile := g.plugin.NewGeneratedFile(filename, file.GoImportPath)

	// Prepare template data
	data := g.prepareTemplateData(file, opts)

	// Process messages
	if opts.GetGenerateInterfaces() {
		g.processMessages(genFile, file.Messages, opts, data)
	}

	// Process services
	if opts.GetGenerateServices() {
		g.processServices(genFile, file.Services, opts, data)
	}

	// Process enums
	if opts.NeedsEnums(file) {
		g.processEnums(genFile, file.Enums, opts, data)
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

// generateNoImplWithTemplate generates NoImpl structs in a separate file
func (g *Generator) generateNoImplWithTemplate(file *protogen.File, opts *GeneratorOptions) error {
	if file == nil {
		return core.Wrap(core.ErrInvalid, "file")
	}
	filename := file.GeneratedFilenamePrefix + ".noimpl.go"
	genFile := g.plugin.NewGeneratedFile(filename, file.GoImportPath)

	// Prepare template data
	data := g.prepareNoImplTemplateData(file, opts)

	// Process messages and services
	if opts.GetGenerateInterfaces() {
		g.processMessages(genFile, file.Messages, opts, data)
	}
	if opts.GetGenerateServices() {
		g.processServices(genFile, file.Services, opts, data)
	}

	// Render template
	var buf bytes.Buffer
	if err := templates.ExecuteTemplate(&buf, noImplFileTemplate, data); err != nil {
		return core.Wrapf(err, "failed to execute template %q for file %q", noImplFileTemplate, filename)
	}

	// Write to generated file
	genFile.P(buf.String())

	return nil
}

// prepareTemplateData creates the base template data structure
func (*Generator) prepareTemplateData(file *protogen.File, opts *GeneratorOptions) *TemplateData {
	// Collect imports as paths
	var stdPaths []string
	var thirdPartyPaths []string

	// Add context import if we need services (standard library)
	if opts.NeedsServices(file) {
		stdPaths = append(stdPaths, "context")
	}

	// Add errors import if we need enums (standard library)
	if opts.NeedsEnums(file) {
		stdPaths = append(stdPaths, "errors")
	}
	// Use GenerateImports to process and organize imports
	importGroups := gengo.GenerateImports(stdPaths, thirdPartyPaths)

	return &TemplateData{
		Package:      string(file.GoPackageName),
		SourceFile:   file.Desc.Path(),
		ImportGroups: importGroups,
		Messages:     make([]MessageData, 0, len(file.Messages)),
		Services:     make([]ServiceData, 0, len(file.Services)),
		Enums:        make([]EnumData, 0, len(file.Enums)),
		NoImpl:       opts.GetGenerateNoImpl(),
	}
}

// prepareNoImplTemplateData creates template data for NoImpl file
func (*Generator) prepareNoImplTemplateData(file *protogen.File, opts *GeneratorOptions) *TemplateData {
	// Collect imports as paths
	var stdPaths []string
	var thirdPartyPaths []string

	// Add context import if we have services (standard library)
	if opts.HasServices(file) && opts.GetGenerateServices() {
		stdPaths = append(stdPaths, "context")
	}

	// Add required imports for NoImpl only if we have content
	if opts.NeedsNoImpl(file) {
		thirdPartyPaths = append(thirdPartyPaths, "darvaza.org/core")
		if opts.HasMessages(file) && opts.GetGenerateInterfaces() {
			thirdPartyPaths = append(thirdPartyPaths, "google.golang.org/protobuf/reflect/protoreflect")
		}
	}

	// Use GenerateImports to process and organize imports
	importGroups := gengo.GenerateImports(stdPaths, thirdPartyPaths)

	return &TemplateData{
		Package:      string(file.GoPackageName),
		SourceFile:   file.Desc.Path(),
		ImportGroups: importGroups,
		Messages:     make([]MessageData, 0, len(file.Messages)),
		Services:     make([]ServiceData, 0, len(file.Services)),
		Enums:        make([]EnumData, 0, len(file.Enums)),
		NoImpl:       true,
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
		InterfaceName: gengo.InterfaceNameForMessage(msg, opts.GetInterfacePattern()),
		NoImplName:    "NoImpl" + msg.GoIdent.GoName,
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
			Name:      string(field.Desc.Name()),
			GoName:    field.GoName,
			Type:      gengo.GoTypeForFieldWithPattern(gen, field, opts.GetInterfacePattern()),
			Comment:   strings.TrimSpace(string(field.Comments.Leading)),
			Optional:  field.Desc.HasOptionalKeyword(),
			IsOneOf:   false,
			IsMessage: field.Desc.Kind() == protoreflect.MessageKind,
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
				Type:      gengo.GoTypeForFieldWithPattern(gen, field, opts.GetInterfacePattern()),
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
		InterfaceName: gengo.InterfaceNameForService(svc, opts.GetInterfacePattern()),
		NoImplName:    "NoImpl" + svc.GoName,
		Comment:       strings.TrimSpace(string(svc.Comments.Leading)),
		Methods:       make([]MethodData, 0, len(svc.Methods)),
	}

	for _, method := range svc.Methods {
		methodData := MethodData{
			Name:         method.GoName,
			Comment:      strings.TrimSpace(string(method.Comments.Leading)),
			RequestType:  gengo.InterfaceNameForMessage(method.Input, opts.GetInterfacePattern()),
			ResponseType: gengo.InterfaceNameForMessage(method.Output, opts.GetInterfacePattern()),
		}
		data.Methods = append(data.Methods, methodData)
	}

	return data
}

// processEnums adds all enums to template data
func (g *Generator) processEnums(gen *protogen.GeneratedFile, enums []*protogen.Enum,
	opts *GeneratorOptions, data *TemplateData) {
	for _, enum := range enums {
		enumData := g.buildEnumData(gen, enum, opts)
		data.Enums = append(data.Enums, enumData)
	}
}

// buildEnumData builds template data for an enum
func (*Generator) buildEnumData(_ *protogen.GeneratedFile, enum *protogen.Enum, opts *GeneratorOptions) EnumData {
	enumName := gengo.EnumNameFor(enum, opts.GetEnumPattern())

	data := EnumData{
		Name:        enumName,
		NamePrivate: strings.ToLower(enumName[:1]) + enumName[1:],
		Comment:     strings.TrimSpace(string(enum.Comments.Leading)),
		Values:      make([]EnumValueData, 0, len(enum.Values)),
	}

	originalPrefix := enum.GoIdent.GoName + "_"

	for _, value := range enum.Values {
		unprefixedName, _ := strings.CutPrefix(value.GoIdent.GoName, originalPrefix)
		name := data.Name + "_" + unprefixedName

		valueData := EnumValueData{
			Name:           name,
			UnprefixedName: unprefixedName,
			Comment:        strings.TrimSpace(string(value.Comments.Leading)),
			Number:         int32(value.Desc.Number()),
		}
		data.Values = append(data.Values, valueData)
	}

	return data
}
