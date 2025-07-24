package main

import (
	"embed"
	"text/template"

	gengo "protomcp.org/protomcp/pkg/generator/gen-go"
)

//go:embed templates/*.tmpl
var templateFS embed.FS

var templates = template.Must(template.New("").ParseFS(templateFS, "templates/*.tmpl"))

// Template names
const (
	fileTemplate          = "file.tmpl"
	messageTemplate       = "message.tmpl"
	fieldTemplate         = "field.tmpl"
	serviceTemplate       = "service.tmpl"
	rpcTemplate           = "rpc.tmpl"
	enumTemplate          = "enum.tmpl"
	noImplFileTemplate    = "noImpl_file.tmpl"
	noImplMessageTemplate = "noImpl_message.tmpl"
	noImplServiceTemplate = "noImpl_service.tmpl"
)

// TemplateData holds the data for rendering templates
type TemplateData struct {
	Package      string
	SourceFile   string
	ImportGroups [][]gengo.Import // Groups of imports, separated by blank lines
	Messages     []MessageData
	Services     []ServiceData
	Enums        []EnumData
	NoImpl       bool // Whether to generate NoImpl structs
}

// MessageData holds data for a message interface
type MessageData struct {
	Name          string
	InterfaceName string
	NoImplName    string // Name for NoImpl struct
	Comment       string
	Fields        []FieldData
	OneOfGroups   []OneOfData
}

// FieldData holds data for a field
type FieldData struct {
	Name      string
	GoName    string
	Type      string
	Comment   string
	OneOfName string
	Optional  bool
	IsOneOf   bool
	IsMessage bool // true if field is a message type (interface)
}

// OneOfData holds data for a `oneof` group
type OneOfData struct {
	Name   string
	GoName string
	Fields []OneOfFieldData
}

// OneOfFieldData holds data for a field in a `oneof` group
type OneOfFieldData struct {
	Name      string
	GoName    string
	Type      string
	OneOfName string // Name of the parent `oneof` group
}

// ServiceData holds data for a service interface
type ServiceData struct {
	Name          string
	InterfaceName string
	NoImplName    string // Name for NoImpl struct
	Comment       string
	Methods       []MethodData
}

// MethodData holds data for an RPC method
type MethodData struct {
	Name         string
	Comment      string
	RequestType  string
	ResponseType string
}

// EnumData holds data for an enum
type EnumData struct {
	// Generated name with pattern applied (e.g., StatusEnum)
	Name string
	// Private name for variables (e.g., statusEnum)
	NamePrivate string
	Comment     string
	Values      []EnumValueData
}

// EnumValueData holds data for an enum value
type EnumValueData struct {
	// Full constant name (e.g., StatusEnum_STATUS_PENDING)
	Name string
	// Unprefixed name for maps (e.g., STATUS_PENDING)
	UnprefixedName string
	Comment        string
	Number         int32
}
