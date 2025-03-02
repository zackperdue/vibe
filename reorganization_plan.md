# Vibe Programming Language Reorganization Plan

## Current Issues

The current codebase has several large files that are difficult to maintain:
- `parser/parser.go` (67KB, 2465 lines)
- `interpreter/interpreter.go` (53KB, 1681 lines)

These large files make it difficult to:
1. Track down bugs
2. Make targeted edits
3. Understand the codebase
4. Collaborate effectively

## Proposed Directory Structure

```
vibe/
├── lexer/
│   ├── lexer.go
│   ├── token.go
│   └── lexer_test.go
├── parser/
│   ├── ast/
│   │   ├── ast.go (node interfaces and common types)
│   │   ├── expressions.go (expression nodes)
│   │   └── statements.go (statement nodes)
│   ├── parser.go (core parser functionality)
│   ├── expression_parser.go (expression parsing)
│   ├── statement_parser.go (statement parsing)
│   ├── function_parser.go (function-related parsing)
│   ├── class_parser.go (class-related parsing)
│   └── parser_test.go
├── interpreter/
│   ├── environment.go (environment and scope handling)
│   ├── interpreter.go (core interpreter functionality)
│   ├── values.go (value types)
│   ├── evaluator.go (statement evaluation)
│   ├── expression_evaluator.go (expression evaluation)
│   ├── function_evaluator.go (function evaluation)
│   ├── class_evaluator.go (class evaluation)
│   └── interpreter_test.go
├── types/
│   └── types.go
└── main.go
```

## Implementation Strategy

Since we can't easily refactor the existing codebase without breaking functionality, we should follow this approach:

1. Create a new branch for the reorganization
2. Create the new directory structure
3. Move code from the existing files to the new files, one piece at a time
4. Update imports and references
5. Run tests after each change to ensure functionality is preserved
6. Once all code is moved and tests pass, merge the branch

## Specific Reorganization Tasks

### Parser Reorganization

1. Create `parser/ast/` directory with:
   - `ast.go`: Node interface, NodeType, and basic AST nodes
   - `expressions.go`: Expression-related AST nodes
   - `statements.go`: Statement-related AST nodes

2. Split `parser/parser.go` into:
   - `parser.go`: Core parser functionality (initialization, token handling)
   - `expression_parser.go`: Expression parsing functions
   - `statement_parser.go`: Statement parsing functions
   - `function_parser.go`: Function-related parsing
   - `class_parser.go`: Class-related parsing

### Interpreter Reorganization

1. Split `interpreter/interpreter.go` into:
   - `environment.go`: Environment and scope handling
   - `values.go`: Value types and interfaces
   - `interpreter.go`: Core interpreter functionality
   - `evaluator.go`: Statement evaluation
   - `expression_evaluator.go`: Expression evaluation
   - `function_evaluator.go`: Function evaluation
   - `class_evaluator.go`: Class evaluation

## Benefits of Reorganization

1. **Improved Maintainability**: Smaller, focused files are easier to understand and maintain
2. **Better Bug Tracking**: Issues can be isolated to specific functionality
3. **Easier Collaboration**: Team members can work on different parts of the codebase without conflicts
4. **Enhanced Testability**: Focused functionality is easier to test
5. **Clearer Code Organization**: Code is organized by functionality, making it easier to find and understand

## Example: Function-Related Code

As an example, all function-related code would be organized as follows:
- AST nodes in `parser/ast/expressions.go` and `parser/ast/statements.go`
- Parsing logic in `parser/function_parser.go`
- Evaluation logic in `interpreter/function_evaluator.go`

This makes it much easier to understand and modify function-related functionality.

## Next Steps

1. Create a new branch for the reorganization
2. Set up the new directory structure
3. Begin moving code, starting with the AST nodes
4. Update imports and references
5. Run tests frequently to ensure functionality is preserved
6. Document the new structure
7. Merge when complete