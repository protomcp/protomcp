# AGENT.md

This file provides guidance to AI agents when working with code in this
repository. For developers and general project information, please refer to
[README.md](README.md) first.

> [!WARNING]
> This project is in DESIGN PHASE ONLY. Nothing is stable, everything may
> change dramatically. Do not consider any current implementation as final.

## Repository Overview

`protomcp` is a protoc generator framework for creating unified JSON-RPC 2.0 and
MCP (Model Context Protocol) endpoints from Protocol Buffer service definitions.
It generates Go code that serves both protocols over HTTP/2 and QUIC transports,
emphasizing interface-based design for modularity.

## Prerequisites

Before starting development, ensure you have:

- Go 1.23 or later installed (check with `go version`).
- `make` command available (usually pre-installed on Unix systems).
- Protocol Buffers compiler (`protoc`) installed.
- `pnpm` for JavaScript/TypeScript tooling (preferred over npm).
- Git configured for proper line endings.

## Common Development Commands

```bash
# Full build cycle (get deps, generate, tidy, build)
make all

# Run tests
make test

# Run tests with coverage
make test GOTEST_FLAGS="-cover"

# Run tests with verbose output and coverage
make test GOTEST_FLAGS="-v -cover"

# Generate coverage reports
make coverage

# Generate Codecov configuration and upload scripts
make codecov

# Format code and tidy dependencies (run before committing)
make tidy

# Clean build artifacts
make clean

# Update dependencies
make up

# Run go:generate directives
make generate

# Generate example code from proto files
make examples

# Lint example proto files
make examples-lint
```

## Build System Features

### Multi-Module Support

The project uses a sophisticated build system inherited from nanorpc that
handles multiple Go modules:

- **Root module**: `protomcp.org/protomcp`
- **Submodules**: Package-specific modules as needed.
- **Dynamic rules**: Generated via `internal/build/gen_mk.sh`.
- **Dependency tracking**: Handles inter-module dependencies.

### Tool Integration

The build system includes comprehensive tooling:

#### Linting and Quality

- **golangci-lint**: Go code linting with version selection.
- **revive**: Additional Go linting with custom rules.
- **buf**: Protocol buffer linting and generation.
- **markdownlint**: Markdown formatting and style checking.
- **shellcheck**: Shell script analysis.
- **cspell**: Spell checking for documentation and code.
- **languagetool**: Grammar checking for Markdown files.

All Go tools (golangci-lint, revive, buf) are managed via `go run` for
consistent versioning without manual installation.

#### Coverage and Testing

- **Coverage collection**: Automated across all modules.
- **Codecov integration**: Multi-module coverage reporting.
- **Test execution**: Parallel testing with dependency management.

#### Development Tools

- **Whitespace fixing**: Automated trailing whitespace removal.
- **EOF handling**: Ensures files end with newlines.
- **Dynamic tool detection**: Tools auto-detected via pnpx.

### Configuration Files

Tool configurations are stored in `internal/build/`:

- `markdownlint.json`: Markdown linting rules (80-char lines)
- `cspell.json`: Spell checking dictionary and rules
- `languagetool.cfg`: Grammar checking configuration
- `revive.toml`: Go linting rules and thresholds

## Project Architecture

### Core Components

- **protoc generator**: Protocol Buffer service definition parser and code
  generator
- **JSON-RPC 2.0 layer**: Built on sourcegraph/jsonrpc2 for proven
  reliability
- **MCP protocol layer**: Implementation of Anthropic's Model Context
  Protocol
- **Transport abstraction**: HTTP/2 and QUIC transport support
- **Interface framework**: Service interfaces independent of protobuf
  concrete types
- **Validation system**: JSON Schema generation and request/response
  validation

### Design Principles

- **Interface-First**: Prioritize interfaces over structs to avoid tight
  coupling
- **Modular Architecture**: Clear separation between protocol, transport, and
  business logic
- **Schema Validation**: Comprehensive JSON Schema validation for type safety
- **Transport Agnostic**: Services work across HTTP/2 and QUIC transports
- **Protocol Unification**: Single service definition serves both JSON-RPC
  2.0 and MCP

### Key Dependencies

- **sourcegraph/jsonrpc2**: Core JSON-RPC 2.0 implementation
- **Model Context Protocol**: Anthropic's MCP for AI assistant integration
- **HTTP/2 & QUIC**: Modern transport protocols for performance
- **JSON Schema**: Validation and type safety
- **Protocol Buffers**: Service definition source of truth

### Example Proto Files

The `proto/examples/` directory contains example protobuf definitions for
testing:

- **Structure**: Follows buf's recommended layout (e.g., `calculator/v1/`)
- **Configuration**:
  - `buf.yaml`: Lint configuration using STANDARD rules
  - `buf.gen.yaml`: Code generation configuration
- **Usage**:

  ```bash
  # Generate Go code from examples
  make examples

  # Lint proto files
  make examples-lint
  ```

The generated code demonstrates the protoc-gen-protomcp output with:

- Interface definitions for messages
- Service interfaces with context support
- Configurable interface naming patterns

## Development Workflow

### Before Starting Work

1. **Understand protocols**: Review JSON-RPC 2.0 and MCP specifications
2. **Check dependencies**: Understand sourcegraph/jsonrpc2 patterns
3. **Review nanorpc**: Study inherited build system patterns
4. **Interface design**: Plan interface-first architecture

### Code Quality Standards

The project enforces quality through:

- **Go standards**: Standard Go conventions and formatting
- **Field alignment**: Structs optimized for memory efficiency

  ```bash
  # Fix field alignment issues (exclude generated files like *.pb.go)
  GOXTOOLS="golang.org/x/tools/go/analysis/passes"
  FA="$GOXTOOLS/fieldalignment/cmd/fieldalignment"
  # Only run on hand-written files, not generated ones
  go run "$FA@latest" -fix <files>

  # For test files with complex types, create a temporary file:
  # 1. Copy struct definitions to a temp.go file with simplified types
  # 2. Run fieldalignment on the temp file
  # 3. Apply the suggested field ordering to the test files
  # 4. Remove the temp file
  ```

- **Interface patterns**: Prefer interfaces over concrete types
- **Validation**: JSON Schema validation for all external inputs
- **Testing**: Comprehensive unit and integration tests
- **Documentation**: All public APIs must be documented

### Protocol Implementation Guidelines

When implementing protocol features:

1. **Interface-first**: Define interfaces before implementations
2. **Schema validation**: Generate and validate JSON schemas
3. **Transport agnostic**: Services must work across transports
4. **Error handling**: Consistent error patterns across protocols
5. **Type safety**: Leverage Go's type system for correctness

## Testing Guidelines

### Test Structure

- **Interface testing**: Test against interfaces, not implementations
- **Protocol compliance**: Verify JSON-RPC 2.0 and MCP spec compliance
- **Transport testing**: Test across HTTP/2 and QUIC transports
- **Schema validation**: Test JSON Schema validation behaviour
- **Integration tests**: End-to-end protocol and transport testing

### Running Tests

```bash
# Run all tests
make test

# Run tests with race detection
make test GOTEST_FLAGS="-race"

# Run specific tests
make test GOTEST_FLAGS="-run TestSpecific"

# Generate coverage
make coverage

# Test specific module
make test-protomcp
```

## Important Notes

### Build System

- Go 1.23 is the minimum required version
- The Makefile dynamically generates rules for submodules
- Tool versions are selected based on Go version
- All tools are auto-detected with fallback to no-op

### Protocol Considerations

- Both JSON-RPC 2.0 and MCP are JSON-based protocols
- Transport must be abstracted to support HTTP/2 and QUIC
- Interface design is critical for modularity
- Schema validation ensures type safety across protocol boundaries

### Development Environment

- Always use `pnpm` instead of `npm` for JavaScript/TypeScript tooling
- Protocol buffer files define the source of truth for services
- Generated code should not be manually edited
- Use `make generate` after service definition changes

## Pre-commit Checklist

1. **ALWAYS run `make tidy` first** - Fix ALL issues before committing:
   - Go code formatting and whitespace clean-up
   - Markdown files checked with markdownlint and cspell
   - Shell scripts checked with shellcheck
   - Protocol buffer files regenerated if needed
2. **Verify all tests pass** with `make test`
3. **Check coverage** with `make coverage` if adding new code
4. **Update documentation** if changing public APIs
5. **Run `make generate`** if service definitions changed

## Git Usage Guidelines

**CRITICAL**: Always follow these git practices to avoid accidental commits:

1. **NEVER use bulk operations** - Always explicitly specify files:

   ```bash
   # CORRECT - explicitly specify files
   git add file1.go file2.go
   git commit -s file1.go file2.go -m "commit message"

   # WRONG - bulk staging/committing
   git add .
   git add -A
   git add -u
   git commit -s -m "commit message"
   git commit -a -m "commit message"
   ```

2. **Use `-s` when doing commits** - Don't take credit for the work

3. **Check what you're committing**:

   ```bash
   git status --porcelain  # Check current state
   git diff --cached       # Review staged changes before committing
   ```

4. **Atomic commits** - Each commit should contain only related changes for a
   single purpose

## Troubleshooting

### Common Issues

1. **Protocol buffer compilation**:
   - Ensure `protoc` is installed and in PATH
   - Verify import paths in proto files
   - Check service definitions are valid

2. **Module dependencies**:
   - Run `make tidy` to fix go.mod issues
   - Check that replace directives are correct
   - Verify inter-module dependencies

3. **Tool detection failures**:
   - Install tools globally with `pnpm install -g <tool>`
   - Check that pnpx is available and functional
   - Tools fall back to no-op if not found

4. **Coverage issues**:
   - Ensure all modules have test files
   - Check that `.tmp/index` exists
   - Use `GOTEST_FLAGS` for additional test configuration

### Getting Help

- Check existing issues and documentation
- Review JSON-RPC 2.0 and MCP specifications
- Study sourcegraph/jsonrpc2 patterns
- Examine nanorpc build system for reference

This project focuses on providing unified, interface-based protocol endpoints
that bridge Protocol Buffer service definitions with modern JSON-based protocols
and transport layers.
