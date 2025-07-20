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
	cache := make(map[string]string)
	return goTypeForFieldWithPatternCached(gen, field, pattern, cache)
}

func goTypeForFieldWithPatternCached(
	gen *protogen.GeneratedFile, field *protogen.Field, pattern string, cache map[string]string,
) string {
	// Create cache key based on field's full name and pattern
	cacheKey := string(field.Desc.FullName()) + ":" + pattern
	if cached, ok := cache[cacheKey]; ok {
		return cached
	}

	// Mark as processing to detect cycles
	cache[cacheKey] = "any" // Fallback for circular references

	// Build the base type
	baseType := getBaseTypeWithPattern(gen, field, pattern)

	// Handle special field types
	var result string
	if field.Desc.IsList() {
		result = "[]" + baseType
	} else if field.Desc.IsMap() {
		keyType := goTypeForFieldWithPatternCached(gen, field.Message.Fields[0], pattern, cache)
		valueType := goTypeForFieldWithPatternCached(gen, field.Message.Fields[1], pattern, cache)
		result = "map[" + keyType + "]" + valueType
	} else {
		// Default case: return base type (handles optional fields and regular fields)
		result = baseType
	}

	// Cache the final result
	cache[cacheKey] = result
	return result
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
	return applyInterfacePattern(msg.GoIdent.GoName, pattern)
}

// InterfaceNameForService returns the interface name for a service based on the pattern
func InterfaceNameForService(svc *protogen.Service, pattern string) string {
	name := svc.GoName

	// For services, ensure we have "Service" suffix unless already present
	if !strings.HasSuffix(name, "Service") {
		name = name + "Service"
	}

	return applyInterfacePattern(name, pattern)
}

// applyInterfacePattern applies the naming pattern to generate interface names
func applyInterfacePattern(name, pattern string) string {
	if pattern == "" {
		pattern = "I%" // Default pattern
	}
	return strings.ReplaceAll(pattern, "%", name)
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
