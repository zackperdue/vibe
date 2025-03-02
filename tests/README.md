# Vibe Language Test Suite

This directory contains a comprehensive test suite for the Vibe programming language. The tests are organized by component to ensure thorough coverage of the language features and implementation.

## Test Structure

The test suite is organized into the following directories:

- **`lexer/`**: Tests for the lexical analyzer (tokenization)
- **`parser/`**: Tests for the parser (syntax analysis)
- **`interpreter/`**: Tests for the interpreter (execution)
- **`integration/`**: End-to-end tests that verify the complete pipeline

Each directory contains tests for its respective component, with comprehensive coverage of language features.

## Running Tests

To run all tests, use the Makefile at the root of the project:

```bash
make test
```

To run tests for a specific component:

```bash
make test-lexer
make test-parser
make test-interpreter
make test-integration
```

## Test Coverage

The test suite covers the following language features:

### Lexer Tests

- Basic tokenization
- Type declarations
- Operators and delimiters
- Number literals (integers, floats)
- String literals
- Keywords
- Line and column tracking
- Complex code samples

### Parser Tests

- Type declarations (simple types and generic types)
- Variable declarations with type annotations
- Array literals
- For loops
- If-else statements
- Function definitions
- Class definitions

### Interpreter Tests

- Evaluation of expressions (arithmetic, comparison, boolean)
- String concatenation
- Variable declarations and assignments
- Conditional expressions
- Return statements
- Function objects and application
- Closures
- Arrays and indexing
- For loops
- Error handling

### Integration Tests

- End-to-end tests from source code to execution
- Complex programs combining multiple language features
- Error handling across the entire pipeline

## Adding New Tests

When adding new tests, follow these guidelines:

1. Place the test in the appropriate directory based on the component being tested
2. Follow the naming convention of existing test functions
3. Ensure that the test verifies both the happy path and error cases
4. Update this README if you add tests for new language features

## Test Helper Functions

The test files include helper functions for common testing tasks:

- `testEval`: Evaluates code and returns the result
- `testIntegerObject`: Verifies that an object is an integer with the expected value
- `testBooleanObject`: Verifies that an object is a boolean with the expected value
- `testStringObject`: Verifies that an object is a string with the expected value
- `testNilObject`: Verifies that an object is nil