# AST Interface Issues and Fix Strategy

## Current Issues

The Vibe project's AST interface has the following issues:

1. **Interface Mismatch**: The project has interfaces `ast.Node`, `ast.Statement`, and `ast.Expression` with distinct hierarchies, but the parser tests are not using them correctly:
   - `ast.Statement` extends `ast.Node` by adding a `statementNode()` method
   - `ast.Expression` extends `ast.Node` by adding an `expressionNode()` method
   - Parser tests are incorrectly using ast.Node where ast.Statement or ast.Expression is expected

2. **Node Type Changes**: The AST has undergone changes, with certain node types being renamed:
   - `ast.LetStatement` → `ast.VariableDecl`
   - `ast.ForLoopStatement` → `ast.ForStmt`
   - `ast.ReturnStatement` → `ast.ReturnStmt`

3. **Field Name Changes**: Field names in AST nodes have changed:
   - `Name.Value` → `Name` (string)
   - `ReturnValue` → `Value`
   - Changes in field types and accessor methods

4. **TokenLiteral Method**: Tests are referencing a `TokenLiteral()` method that is no longer part of the interface

## Fix Strategy

### 1. Update AST Interface Usage in Tests

- Modify all tests to use `ast.Node` as the base type for AST nodes
- Add type assertions for specific node types as needed (e.g., `node.(*ast.VariableDecl)`)
- Replace references to old node types with new ones (e.g., `ast.LetStatement` → `ast.VariableDecl`)

### 2. Add Missing Test Helper Functions

- Add implementations for `testIntegerLiteral`, `testLiteralExpression`, etc. based on the current AST structure
- These functions should handle type assertions and value checking properly

### 3. Fix Type Assertions

- Replace incorrect type assertions like `stmt.(ast.Statement)` with proper node type assertions
- Ensure that all AST node references match the current AST structure

### 4. Update Field References

- Update all field references to match the current AST node structure
- For example, change `Name.Value` to `Name` for `VariableDecl` nodes

### 5. Replace TokenLiteral With Type()

- Replace references to `TokenLiteral()` with `Type()` method checks where appropriate
- Update comparison strings to match the NodeType constants defined in the AST package

## Implementation Steps

1. First, fix the declarations_test.go file as it has fewer references to update
2. Then fix expressions_test.go by updating expression-related AST node references
3. Finally, fix statements_test.go with the updates we've identified

This approach will address the main interface issues while keeping the tests aligned with the current AST structure.