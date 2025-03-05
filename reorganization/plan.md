# Vibe Project Reorganization Plan

## Current Structure Issues

The Vibe project currently has a mixed approach to tests organization:

1. Some tests are co-located with their implementation files (e.g., in `interpreter/`, `parser/`, `lexer/`, `object/` directories)
2. Other tests are in a separate `tests/` directory organized by component (e.g., `tests/interpreter/`, `tests/parser/`)
3. There are integration tests in a dedicated `tests/integration/` directory
4. The parser package has build errors related to AST interface implementation

## Recommended Approach from Documentation

According to our documentation in `making-updates.mdc`:

> Tests in Go projects should follow the standard Go convention:
> - Test files should be placed in the same package as the code they test
> - Test files should be named with the `_test.go` suffix
> - For example, for a file named `lexer.go`, tests should be in `lexer_test.go` in the same directory

## Planned Changes

### 1. Fix the AST Interface Issues in `parser` Package

The parser package has build errors related to the AST interface implementation. We've identified the following issues:

#### AST Interface Mismatch
- The project has interfaces `ast.Node`, `ast.Statement`, and `ast.Expression` with distinct hierarchies
- Parser tests are incorrectly using `ast.Node` where `ast.Statement` or `ast.Expression` is expected
- Several AST node types have been renamed (e.g., `ast.LetStatement` → `ast.VariableDecl`)
- Field names and structure have changed in AST nodes

#### Fix Strategy
1. Update AST interface usage in tests to correctly use the current node types
2. Add proper type assertions and implement missing helper functions
3. Update field references to match the current AST node structure
4. Replace obsolete method references with current ones

### 2. Test Organization

#### Move Component-Specific Tests to Implementation Directories

Move tests from the separate `tests/` directory to their respective implementation directories:

- `tests/parser/*_test.go` → `parser/`
- `tests/lexer/*_test.go` → `lexer/`
- `tests/interpreter/*_test.go` → `interpreter/`
- `tests/object/*_test.go` → `object/`

Ensure there are no duplications or conflicts with existing tests.

#### Keep Integration Tests Separate

Integration tests should remain separate as they test the interaction between multiple components:

- Keep `tests/integration/` directory for integration tests
- Ensure they use public APIs of the components they test

### 3. Update Test Content

When moving tests, ensure they:
- Are properly imported and reference the correct packages
- Follow Go's standard testing practices
- Use table-driven tests where appropriate
- Test both normal operation and error conditions

### 4. Update Test Scripts

Update any test running scripts to reflect the new organization:
- `run_go_tests.sh`

### 5. Documentation Updates

Update any documentation that references tests to reflect the new organization.

## Expected Benefits

This reorganization will:
1. Align the project with Go's standard conventions
2. Make it easier to find tests for specific implementations
3. Ensure tests are updated when implementations change
4. Improve the development workflow by keeping related code together

## Implementation Approach

1. First, fix the AST interface implementation issues in the parser package
2. Move one component's tests at a time, starting with the most stable components
3. Run tests after each move to ensure functionality is preserved
4. Update build and test scripts once all moves are complete
5. Document the new organization for future contributors

## Implementation Order

1. Fix parser package AST interface issues
   - Update declarations_test.go
   - Update expressions_test.go
   - Update statements_test.go
   - Update for_loop_test.go
   - Update parser_test.go

2. Move component-specific tests (one by one)
   - Move lexer tests
   - Move object tests
   - Move interpreter tests
   - Move parser tests (after fixing interface issues)

3. Update test scripts and documentation

## Timeline

This reorganization should be completed in the following phases:
1. Fix parser package build errors
2. Move component-specific tests (one component at a time)
3. Update test scripts
4. Update documentation
