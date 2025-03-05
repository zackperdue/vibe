# Vibe Language Test Guidelines

## Test Organization

The Vibe language project follows Go's standard test organization practices:

1. **Co-located Tests**: Test files are placed in the same package as the code they test
   - Test files are named `<component>_test.go`
   - For example, `lexer_test.go` tests the functionality in `lexer.go`

2. **Integration Tests**: Tests that verify the interaction between multiple components are kept in `tests/integration`
   - These test the entire system from source code to execution

## Writing New Tests

When adding new features or modifying existing ones, follow these guidelines:

### Test File Naming

- Use descriptive names for test files: `<component>_test.go`
- For specialized tests, use `<component>_<feature>_test.go` (e.g., `lexer_keywords_test.go`)

### Test Function Naming

- Test functions should be named `Test<Feature>` (e.g., `TestLexerBasics`, `TestForLoopParsing`)
- For different aspects of the same feature, use `Test<Feature>_<Aspect>` (e.g., `TestAssignment_WithTypeAnnotation`)

### Test Structure

Each test should follow this structure:

1. **Setup**: Define input and expected output
2. **Execution**: Run the code being tested
3. **Verification**: Check if the result matches the expectation

Example:

```go
func TestLexerBasics(t *testing.T) {
    // Setup
    input := "x = 5"
    expected := []struct{
        expectedType TokenType
        expectedLiteral string
    }{
        {IDENT, "x"},
        {ASSIGN, "="},
        {INT, "5"},
        {EOF, ""},
    }

    // Execution
    l := New(input)

    // Verification
    for i, tt := range expected {
        tok := l.NextToken()
        if tok.Type != tt.expectedType {
            t.Fatalf("tests[%d] - tokentype wrong. expected=%q, got=%q",
                i, tt.expectedType, tok.Type)
        }
        if tok.Literal != tt.expectedLiteral {
            t.Fatalf("tests[%d] - literal wrong. expected=%q, got=%q",
                i, tt.expectedLiteral, tok.Literal)
        }
    }
}
```

### Table-Driven Tests

For testing multiple similar cases, use table-driven tests:

```go
func TestOperations(t *testing.T) {
    tests := []struct {
        input    string
        expected int
    }{
        {"5 + 5", 10},
        {"5 - 3", 2},
        {"2 * 4", 8},
        {"8 / 2", 4},
    }

    for _, tt := range tests {
        result := evaluate(tt.input)
        if result != tt.expected {
            t.Errorf("For input %q, expected %d, got %d",
                tt.input, tt.expected, result)
        }
    }
}
```

### Test Helpers

- Create helper functions for common operations
- Place them in the same test file if they're specific to those tests
- Consider creating a `testutil` package for helpers used across multiple packages

### Testing Edge Cases

Always include tests for edge cases:

- Empty input
- Invalid syntax
- Boundary conditions
- Error handling

### Testing Parser Changes

When modifying the parser:

1. Add tests for new syntax features
2. Ensure AST node structure is correctly generated
3. Test with both valid and invalid inputs
4. Update any affected integration tests

### Testing Interpreter Changes

When modifying the interpreter:

1. Test individual evaluation logic
2. Test type checking and error handling
3. Add integration tests for the full execution path

## Running Tests

```bash
# Run all tests
./run_go_tests.sh

# Run tests for a specific package
go test ./lexer
go test ./parser
go test ./interpreter

# Run a specific test
go test ./parser -run=TestForLoopParsing

# Run tests with verbose output
go test -v ./lexer
```

## Test Coverage

Aim for high test coverage:

```bash
# Generate coverage report
go test -coverprofile=coverage.out ./...

# View coverage in browser
go tool cover -html=coverage.out
```

## Continuous Integration

All tests are run in CI:

1. Tests must pass before merging
2. Coverage should not decrease
3. Both component tests and integration tests are run

By following these guidelines, we ensure that the Vibe language remains robust and reliable as it evolves.