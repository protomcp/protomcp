// Package gengo provides reusable helpers for Go code generation.
package gengo

import (
	"strings"

	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/reflect/protoreflect"

	"darvaza.org/core"
)

// IsSyntheticOneOf returns true if the `oneof` is a synthetic `oneof` for proto3 optional
func IsSyntheticOneOf(oneof *protogen.Oneof) bool {
	// Synthetic `oneof` groups for optional fields have names starting with "_"
	return strings.HasPrefix(string(oneof.Desc.Name()), "_")
}

// GoTypeForFieldWithPattern returns the Go type string for a proto field with interface pattern support
func GoTypeForFieldWithPattern(gen *protogen.GeneratedFile, field *protogen.Field, pattern string) string {
	// Build the base type
	baseType := getBaseTypeWithPattern(gen, field, pattern)

	// Handle repeated fields
	if field.Desc.IsList() {
		return "[]" + baseType
	}

	// Handle map fields
	if field.Desc.IsMap() {
		keyType := GoTypeForFieldWithPattern(gen, field.Message.Fields[0], pattern)
		valueType := GoTypeForFieldWithPattern(gen, field.Message.Fields[1], pattern)
		return "map[" + keyType + "]" + valueType
	}

	// Handle optional fields (proto3 optional or proto2 optional)
	if field.Desc.HasOptionalKeyword() && field.Desc.Kind() == protoreflect.MessageKind {
		return baseType // Already has pointer
	}

	return baseType
}

// getBaseTypeWithPattern returns the base Go type for a field without modifiers, applying interface pattern
func getBaseTypeWithPattern(gen *protogen.GeneratedFile, field *protogen.Field, pattern string) string {
	// Handle scalar types
	if typ := scalarGoType(field.Desc.Kind()); typ != "" {
		return typ
	}

	// Handle complex types
	switch field.Desc.Kind() {
	case protoreflect.MessageKind:
		if field.Message != nil {
			// Apply interface pattern to message types
			interfaceName := InterfaceNameForMessage(field.Message, pattern)
			return gen.QualifiedGoIdent(protogen.GoIdent{
				GoImportPath: field.Message.GoIdent.GoImportPath,
				GoName:       interfaceName,
			})
		}
		return "*any"
	case protoreflect.EnumKind:
		if field.Enum != nil {
			return gen.QualifiedGoIdent(field.Enum.GoIdent)
		}
		return "int32"
	default:
		return "any"
	}
}

var scalarTypeMap = map[protoreflect.Kind]string{
	protoreflect.BoolKind:     "bool",
	protoreflect.Int32Kind:    "int32",
	protoreflect.Sint32Kind:   "int32",
	protoreflect.Sfixed32Kind: "int32",
	protoreflect.Uint32Kind:   "uint32",
	protoreflect.Fixed32Kind:  "uint32",
	protoreflect.Int64Kind:    "int64",
	protoreflect.Sint64Kind:   "int64",
	protoreflect.Sfixed64Kind: "int64",
	protoreflect.Uint64Kind:   "uint64",
	protoreflect.Fixed64Kind:  "uint64",
	protoreflect.FloatKind:    "float32",
	protoreflect.DoubleKind:   "float64",
	protoreflect.StringKind:   "string",
	protoreflect.BytesKind:    "[]byte",
}

// scalarGoType returns the Go type for scalar proto kinds
func scalarGoType(kind protoreflect.Kind) string {
	return scalarTypeMap[kind]
}

// InterfaceNameForMessage returns the interface name for a message
func InterfaceNameForMessage(msg *protogen.Message, pattern string) string {
	if pattern == "" {
		pattern = "I%" // Default pattern
	}
	return strings.ReplaceAll(pattern, "%", msg.GoIdent.GoName)
}

// Import represents an import statement
type Import struct {
	Path  string
	Alias string
}

// GenerateImports processes multiple groups of import paths, deduplicating and sorting each group.
// Empty groups are skipped, and groups are meant to be separated by empty lines in the output.
func GenerateImports(groups ...[]string) [][]Import {
	result := make([][]Import, 0, len(groups))

	for _, paths := range groups {
		if len(paths) == 0 {
			continue
		}

		// Deduplicate and sort
		uniquePaths := core.SliceUnique(paths)
		core.SliceSortOrdered(uniquePaths)

		// Convert to Import structs
		imports := make([]Import, 0, len(uniquePaths))
		for _, path := range uniquePaths {
			imports = append(imports, Import{Path: path})
		}

		if len(imports) > 0 {
			result = append(result, imports)
		}
	}

	return result
}
