# Remaining Tasks for Vibe Project Reorganization

## Parser Package Fixes

The following parser tests need to be fixed due to changes in the AST structure:

1. **TestLetStatements**
   - Update to handle direct nodes instead of ExpressionStatement wrappers
   - Fix expectations for statement structure

2. **TestReturnStatements**
   - Fix to handle `ReturnStmt` instead of `ReturnStatement`
   - Update value checking logic

3. **TestForLoopStatement**
   - Update to handle `ForStmt` instead of `ForLoopStatement`
   - Fix body statement checking

4. **TestParseIndexExpression**
   - Fix index expression validation
   - Handle changes in Program.String() output format

5. **TestIfExpression Body Check**
   - Update consequence and alternative checking logic
   - Handle direct nodes instead of ExpressionStatement wrappers

## Test File Moving

After fixing the parser issues, the following files should be moved from the tests directory to their respective implementation directories:

1. **Lexer Tests**
   - Move unique tests from `tests/lexer/lexer_test.go` to `lexer/moved_lexer_test.go`
   - Move unique tests from `tests/lexer/debug_test.go` to `lexer/moved_debug_test.go`
   - Ensure no test function names conflict

2. **Parser Tests**
   - Move unique tests from `tests/parser/parser_test.go` to appropriate test files in the parser directory
   - Ensure tests work with current AST structure

3. **Interpreter Tests**
   - Move unique tests from `tests/interpreter/interpreter_test.go` to appropriate test files in the interpreter directory
   - Check for function conflicts with existing tests

## Integration Tests

Integration tests should remain in the `tests/integration` directory since they test the interaction between multiple components.

## Documentation Updates

1. **Update README.md**
   - Document the new test organization
   - Explain the rationale behind keeping integration tests separate

2. **Update run_go_tests.sh**
   - Once all tests are passing, remove the temporary exclusions
   - Update documentation in the script

3. **Create Test Guidelines**
   - Document best practices for adding new tests
   - Explain the Go testing conventions being followed

## Long-Term Improvements

1. **Consistent AST Structure**
   - Ensure consistent handling of statements and expressions
   - Document the AST structure for future developers

2. **Test Helper Functions**
   - Create shared test helpers to avoid duplication
   - Ensure consistent test patterns across packages

3. **CI Integration**
   - Update any CI scripts to use the new test organization
   - Ensure CI runs all tests, including integration tests