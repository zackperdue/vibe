# Vibe Programming Language Development Guide

## Build/Test Commands

```bash
# Run all tests
make test

# Run specific test components
make test-lexer
make test-parser
make test-interpreter
make test-integration

# Run a single test
go test -v ./lexer -run TestLexer_SomeSpecificTest
go test -v ./parser -run TestParser_SomeSpecificTest
go test -v ./interpreter -run TestInterpreter_SomeSpecificTest

# Run the vibe interpreter
go run main.go path/to/program.vi

# Run REPL
go run main.go -i

# Format code
go fmt ./...
```

## Code Style Guidelines

- **Formatting**: Use standard Go formatting (gofmt)
- **Naming**: 
  - CamelCase for exported symbols, camelCase for internal ones
  - Keep AST node names descriptive (e.g., `FunctionLiteral`, `InfixExpression`)
- **Errors**: Return descriptive error messages, collect errors during parsing rather than failing on first error
- **Tests**: Write comprehensive tests for new features, update existing tests when modifying behavior
- **Documentation**: Document complex algorithms and parser/interpreter patterns
- **Architecture**:
  - Keep lexer, parser, and interpreter components separate
  - Follow the visitor pattern for AST traversal
  - Add new language features by updating lexer, parser, and interpreter in sync
  - Respect the type system in all interpreter operations