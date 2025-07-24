# Generator Options Reference

The `protoc-gen-protomcp` generator supports various options to customize the
generated code. Options can be passed via the protoc command line or through
buf.gen.yaml configuration.

## Command Line Options

### Basic Usage

```bash
protoc --protomcp_out=. --protomcp_opt=paths=source_relative myfile.proto
```

### Available Options

#### `interfaces` (default: true)

Generate interface definitions for protobuf messages.

```bash
--protomcp_opt=interfaces=true
```

#### `services` (default: true)

Generate interface definitions for protobuf services.

```bash
--protomcp_opt=services=true
```

#### `noimpl` (default: true)

Generate NoImpl structs that implement the interfaces with stub methods.

```bash
--protomcp_opt=noimpl=true
```

#### `enums` (default: true)

Generate helper types for protobuf enums with validation and text marshaling.

```bash
--protomcp_opt=enums=true
```

When enabled, generates:

- Type-safe enum types (e.g., `StatusEnum` for enum `Status`)
- `String()` method for string representation
- `IsValid()` method for validation
- `MarshalText()` and `UnmarshalText()` for JSON/text encoding/decoding

#### `interface_pattern` (default: "I%")

Pattern for generated interface names. Use `%` as placeholder for the original
name.

Examples:

- `I%` → `IMessage` (prefix with "I")
- `%Interface` → `MessageInterface` (suffix with "Interface")
- `%` → `Message` (no modification)

```bash
--protomcp_opt=interface_pattern=I%
```

#### `enum_pattern` (default: "%Enum")

Pattern for generated enum type names. Use `%` as placeholder for the original
name.

Examples:

- `%Enum` → `StatusEnum` (suffix with "Enum", avoids clash with .pb.go)
- `E%` → `EStatus` (prefix with "E")
- `%Type` → `StatusType` (suffix with "Type")

```bash
--protomcp_opt=enum_pattern=%Enum
```

## buf.gen.yaml Configuration

```yaml
version: v2
plugins:
  - plugin: protomcp
    out: .
    opt:
      - paths=source_relative
      - interfaces=true
      - services=true
      - noimpl=true
      - enums=true
      - interface_pattern=I%
      - enum_pattern=%Enum
```

## Generated Files

The generator creates the following files:

### `.types.go` Files

Contains:

- Interface definitions for messages (if `interfaces=true`)
- Interface definitions for services (if `services=true`)
- Enum types with helper methods (if `enums=true`)

### `.noimpl.go` Files

Contains:

- NoImpl structs for messages (if `noimpl=true` and `interfaces=true`)
- NoImpl structs for services (if `noimpl=true` and `services=true`)

Note: Enum types don't have NoImpl structs as they are value types, not
interfaces.

## Examples

### Disable Enum Generation

```bash
protoc --protomcp_out=. \
  --protomcp_opt=paths=source_relative,enums=false \
  myfile.proto
```

### Custom Naming Patterns

```bash
protoc --protomcp_out=. \
  --protomcp_opt=paths=source_relative,\
interface_pattern=%Interface,enum_pattern=E% \
  myfile.proto
```

### Minimal Generation (Interfaces Only)

```bash
protoc --protomcp_out=. \
  --protomcp_opt=paths=source_relative,services=false,noimpl=false,enums=false \
  myfile.proto
```

## Pattern Examples

Given a proto file:

```protobuf
message User {
  string id = 1;
}

enum Status {
  STATUS_UNSPECIFIED = 0;
  STATUS_ACTIVE = 1;
}

service UserService {
  rpc GetUser(GetUserRequest) returns (GetUserResponse);
}
```

Different patterns produce:

| Type | Pattern | Result |
|------|---------|--------|
| Message Interface | `I%` | `IUser` |
| Message Interface | `%Interface` | `UserInterface` |
| Service Interface | `I%` | `IUserService` |
| Service Interface | `%Service` | `UserServiceService` |
| Enum Type | `%Enum` | `StatusEnum` |
| Enum Type | `E%` | `EStatus` |
| NoImpl Struct | N/A | `NoImplUser`, `NoImplUserService` |
