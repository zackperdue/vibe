package parser_test

import (
	"testing"

	"github.com/vibe-lang/vibe/ast"
	"github.com/vibe-lang/vibe/lexer"
	"github.com/vibe-lang/vibe/parser"
)

// TestTypeDeclaration tests parsing of simple type declarations
func TestTypeDeclaration(t *testing.T) {
	input := `type StringAlias = string`

	l := lexer.New(input)
	program, errors := parser.Parse(l)

	if len(errors) != 0 {
		t.Errorf("parser has %d errors", len(errors))
		for _, msg := range errors {
			t.Errorf("parser error: %q", msg)
		}
		t.FailNow()
	}

	if len(program.Statements) != 1 {
		t.Fatalf("program has not enough statements. got=%d", len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.TypeDeclaration)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.TypeDeclaration. got=%T", program.Statements[0])
	}

	if stmt.Name != "StringAlias" {
		t.Errorf("stmt.Name not 'StringAlias'. got=%q", stmt.Name)
	}

	typeAnnotation, ok := stmt.TypeValue.(*ast.TypeAnnotation)
	if !ok {
		t.Fatalf("stmt.TypeValue is not ast.TypeAnnotation. got=%T", stmt.TypeValue)
	}

	if typeAnnotation.TypeName != "string" {
		t.Errorf("typeAnnotation.TypeName not 'string'. got=%q", typeAnnotation.TypeName)
	}
}

// TestGenericTypeDeclaration tests parsing of type declarations with generic types
func TestGenericTypeDeclaration(t *testing.T) {
	input := `type StringArray = Array<string>`

	l := lexer.New(input)
	program, errors := parser.Parse(l)

	if len(errors) != 0 {
		t.Errorf("parser has %d errors", len(errors))
		for _, msg := range errors {
			t.Errorf("parser error: %q", msg)
		}
		t.FailNow()
	}

	if len(program.Statements) != 1 {
		t.Fatalf("program has not enough statements. got=%d", len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.TypeDeclaration)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.TypeDeclaration. got=%T", program.Statements[0])
	}

	if stmt.Name != "StringArray" {
		t.Errorf("stmt.Name not 'StringArray'. got=%q", stmt.Name)
	}

	typeAnnotation, ok := stmt.TypeValue.(*ast.TypeAnnotation)
	if !ok {
		t.Fatalf("stmt.TypeValue is not ast.TypeAnnotation. got=%T", stmt.TypeValue)
	}

	if typeAnnotation.TypeName != "Array" {
		t.Errorf("typeAnnotation.TypeName not 'Array'. got=%q", typeAnnotation.TypeName)
	}

	if len(typeAnnotation.TypeParams) != 1 {
		t.Fatalf("length of typeAnnotation.TypeParams wrong. got=%d", len(typeAnnotation.TypeParams))
	}

	param, ok := typeAnnotation.TypeParams[0].(*ast.TypeAnnotation)
	if !ok {
		t.Fatalf("typeAnnotation.TypeParams[0] is not ast.TypeAnnotation. got=%T", typeAnnotation.TypeParams[0])
	}

	if param.TypeName != "string" {
		t.Errorf("param.TypeName not 'string'. got=%q", param.TypeName)
	}
}

// TestMultipleGenericTypeParameters tests parsing of types with multiple generic parameters
func TestMultipleGenericTypeParameters(t *testing.T) {
	input := `type Dictionary = Map<string, number>`

	l := lexer.New(input)
	program, errors := parser.Parse(l)

	if len(errors) != 0 {
		t.Errorf("parser has %d errors", len(errors))
		for _, msg := range errors {
			t.Errorf("parser error: %q", msg)
		}
		t.FailNow()
	}

	stmt, ok := program.Statements[0].(*ast.TypeDeclaration)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.TypeDeclaration. got=%T", program.Statements[0])
	}

	typeAnnotation, ok := stmt.TypeValue.(*ast.TypeAnnotation)
	if !ok {
		t.Fatalf("stmt.TypeValue is not ast.TypeAnnotation. got=%T", stmt.TypeValue)
	}

	if len(typeAnnotation.TypeParams) != 2 {
		t.Fatalf("length of typeAnnotation.TypeParams wrong. got=%d", len(typeAnnotation.TypeParams))
	}

	stringParam, ok := typeAnnotation.TypeParams[0].(*ast.TypeAnnotation)
	if !ok || stringParam.TypeName != "string" {
		t.Fatalf("First type parameter should be 'string'. got=%v", typeAnnotation.TypeParams[0])
	}

	numberParam, ok := typeAnnotation.TypeParams[1].(*ast.TypeAnnotation)
	if !ok || numberParam.TypeName != "number" {
		t.Fatalf("Second type parameter should be 'number'. got=%v", typeAnnotation.TypeParams[1])
	}
}

// TestNestedGenericTypes tests parsing of types with nested generic parameters
func TestNestedGenericTypes(t *testing.T) {
	input := `type NestedArray = Array<Array<string>>`

	l := lexer.New(input)
	program, errors := parser.Parse(l)

	if len(errors) != 0 {
		t.Errorf("parser has %d errors", len(errors))
		for _, msg := range errors {
			t.Errorf("parser error: %q", msg)
		}
		t.FailNow()
	}

	stmt, ok := program.Statements[0].(*ast.TypeDeclaration)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.TypeDeclaration. got=%T", program.Statements[0])
	}

	typeAnnotation, ok := stmt.TypeValue.(*ast.TypeAnnotation)
	if !ok {
		t.Fatalf("stmt.TypeValue is not ast.TypeAnnotation. got=%T", stmt.TypeValue)
	}

	if len(typeAnnotation.TypeParams) != 1 {
		t.Fatalf("length of typeAnnotation.TypeParams wrong. got=%d", len(typeAnnotation.TypeParams))
	}

	innerType, ok := typeAnnotation.TypeParams[0].(*ast.TypeAnnotation)
	if !ok || innerType.TypeName != "Array" {
		t.Fatalf("Inner type should be 'Array'. got=%v", typeAnnotation.TypeParams[0])
	}

	if len(innerType.TypeParams) != 1 {
		t.Fatalf("Inner Array should have 1 type parameter. got=%d", len(innerType.TypeParams))
	}

	stringParam, ok := innerType.TypeParams[0].(*ast.TypeAnnotation)
	if !ok || stringParam.TypeName != "string" {
		t.Fatalf("Innermost type parameter should be 'string'. got=%v", innerType.TypeParams[0])
	}
}

// TestVariableDeclarationWithTypeAnnotation tests parsing of variable declarations with type annotations
func TestVariableDeclarationWithTypeAnnotation(t *testing.T) {
	input := `a: string = "hello"`

	l := lexer.New(input)
	program, errors := parser.Parse(l)

	if len(errors) != 0 {
		t.Errorf("parser has %d errors", len(errors))
		for _, msg := range errors {
			t.Errorf("parser error: %q", msg)
		}
		t.FailNow()
	}

	if len(program.Statements) != 1 {
		t.Fatalf("program has not enough statements. got=%d", len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.VariableDecl)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.VariableDecl. got=%T", program.Statements[0])
	}

	if stmt.Name != "a" {
		t.Errorf("stmt.Name not 'a'. got=%q", stmt.Name)
	}

	typeAnnotation := stmt.TypeAnnotation
	if typeAnnotation == nil {
		t.Fatalf("stmt.TypeAnnotation is nil")
	}

	if typeAnnotation.TypeName != "string" {
		t.Errorf("typeAnnotation.TypeName not 'string'. got=%q", typeAnnotation.TypeName)
	}

	stringLiteral, ok := stmt.Value.(*ast.StringLiteral)
	if !ok {
		t.Fatalf("stmt.Value is not ast.StringLiteral. got=%T", stmt.Value)
	}

	if stringLiteral.Value != "hello" {
		t.Errorf("stringLiteral.Value not 'hello'. got=%q", stringLiteral.Value)
	}
}

// TestArrayLiteralParsing tests parsing of array literals
func TestArrayLiteralParsing(t *testing.T) {
	tests := []struct {
		input    string
		expected []string // string representation of expected values
	}{
		{`[]`, []string{}},
		{`[1]`, []string{"1"}},
		{`[1, 2]`, []string{"1", "2"}},
		{`["a", "b", "c"]`, []string{`String("a")`, `String("b")`, `String("c")`}},
		{`[1, "a", true]`, []string{"1", `String("a")`, "Boolean(true)"}},
		{`[1,]`, []string{"1"}}, // trailing comma
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		program, errors := parser.Parse(l)

		if len(errors) != 0 {
			t.Errorf("parser has %d errors for input %q", len(errors), tt.input)
			for _, msg := range errors {
				t.Errorf("parser error: %q", msg)
			}
			continue
		}

		if len(program.Statements) != 1 {
			t.Fatalf("program should have 1 statement for input %q. got=%d",
				tt.input, len(program.Statements))
		}

		expr, ok := program.Statements[0].(*ast.ArrayLiteral)
		if !ok {
			t.Fatalf("program.Statements[0] is not ast.ArrayLiteral for input %q. got=%T",
				tt.input, program.Statements[0])
		}

		if len(expr.Elements) != len(tt.expected) {
			t.Fatalf("array has wrong number of elements for input %q. got=%d, want=%d",
				tt.input, len(expr.Elements), len(tt.expected))
		}

		// Simple string representation check
		for i, element := range expr.Elements {
			if element.String() != tt.expected[i] {
				t.Errorf("element %d has wrong value for input %q. got=%s, want=%s",
					i, tt.input, element.String(), tt.expected[i])
			}
		}
	}
}

// TestForLoopParsing tests parsing of for loops
func TestForLoopParsing(t *testing.T) {
	tests := []struct {
		input           string
		expectedIter    string
		expectedElements []string
		expectedBodyLen int
	}{
		{
			`for x in [1] do
end`,
			"x",
			[]string{"1"},
			0,
		},
		{
			`for i in [1, 2, 3] do
				print(i)
			end`,
			"i",
			[]string{"1", "2", "3"},
			1,
		},
		{
			`for item in items do
				print(item)
				count += 1
			end`,
			"item",
			[]string{"items"},
			2,
		},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		program, errors := parser.Parse(l)

		if len(errors) != 0 {
			t.Errorf("parser has %d errors for input %q", len(errors), tt.input)
			for _, msg := range errors {
				t.Errorf("parser error: %q", msg)
			}
			continue
		}

		if len(program.Statements) != 1 {
			t.Fatalf("program should have 1 statement for input %q. got=%d",
				tt.input, len(program.Statements))
		}

		forStmt, ok := program.Statements[0].(*ast.ForStmt)
		if !ok {
			t.Fatalf("program.Statements[0] is not ast.ForStmt for input %q. got=%T",
				tt.input, program.Statements[0])
		}

		if forStmt.Iterator != tt.expectedIter {
			t.Errorf("iterator is not %q for input %q. got=%q",
				tt.expectedIter, tt.input, forStmt.Iterator)
		}

		// Check the iterable
		arrayLit, ok := forStmt.Iterable.(*ast.ArrayLiteral)
		if !ok && len(tt.expectedElements) > 0 {
			// If it's not an array literal but an identifier (for the third test case)
			ident, ok := forStmt.Iterable.(*ast.Identifier)
			if !ok {
				t.Fatalf("iterable is not ast.ArrayLiteral or ast.Identifier for input %q. got=%T",
					tt.input, forStmt.Iterable)
			}
			if ident.Name != tt.expectedElements[0] {
				t.Errorf("iterable identifier is not %q for input %q. got=%q",
					tt.expectedElements[0], tt.input, ident.Name)
			}
		} else if ok {
			// It's an array literal
			if len(arrayLit.Elements) != len(tt.expectedElements) {
				t.Errorf("array has wrong number of elements for input %q. got=%d, want=%d",
					tt.input, len(arrayLit.Elements), len(tt.expectedElements))
			} else {
				// Check array elements
				for i, element := range arrayLit.Elements {
					if element.String() != tt.expectedElements[i] {
						t.Errorf("element %d has wrong value for input %q. got=%s, want=%s",
							i, tt.input, element.String(), tt.expectedElements[i])
					}
				}
			}
		}

		// Check the body length
		if len(forStmt.Body.Statements) != tt.expectedBodyLen {
			t.Errorf("body has wrong number of statements for input %q. got=%d, want=%d",
				tt.input, len(forStmt.Body.Statements), tt.expectedBodyLen)
		}
	}
}

// TestIfStatementParsing tests parsing of if statements
func TestIfStatementParsing(t *testing.T) {
	tests := []struct {
		input          string
		conditionStr   string
		consequenceLen int
		hasElse        bool
		alternativeLen int
	}{
		{
			`if x > 5 do
				print("x is greater than 5")
			end`,
			"(x > 5)",
			1,
			false,
			0,
		},
		{
			`if x > 5 do
				print("x is greater than 5")
			else
				print("x is not greater than 5")
			end`,
			"(x > 5)",
			1,
			true,
			1,
		},
		{
			`if x > 5 do
				print("x is greater than 5")
			elsif x < 0 do
				print("x is negative")
			else
				print("x is between 0 and 5")
			end`,
			"(x > 5)",
			1,
			true,
			1,
		},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		program, errors := parser.Parse(l)

		if len(errors) != 0 {
			t.Errorf("parser has %d errors for input %q", len(errors), tt.input)
			for _, msg := range errors {
				t.Errorf("parser error: %q", msg)
			}
			continue
		}

		if len(program.Statements) != 1 {
			t.Fatalf("program should have 1 statement for input %q. got=%d",
				tt.input, len(program.Statements))
		}

		ifStmt, ok := program.Statements[0].(*ast.IfStmt)
		if !ok {
			t.Fatalf("program.Statements[0] is not ast.IfStmt for input %q. got=%T",
				tt.input, program.Statements[0])
		}

		// Check condition
		if ifStmt.Condition.String() != tt.conditionStr {
			t.Errorf("condition is not %q for input %q. got=%q",
				tt.conditionStr, tt.input, ifStmt.Condition.String())
		}

		// Check consequence
		if len(ifStmt.Consequence.Statements) != tt.consequenceLen {
			t.Errorf("consequence has wrong length for input %q. got=%d, want=%d",
				tt.input, len(ifStmt.Consequence.Statements), tt.consequenceLen)
		}

		// Check if there's an else
		if tt.hasElse && ifStmt.Alternative == nil {
			t.Errorf("expected alternative but got nil for input %q", tt.input)
		} else if !tt.hasElse && ifStmt.Alternative != nil {
			t.Errorf("expected no alternative but got one for input %q", tt.input)
		}

		// Check alternative length if it exists
		if ifStmt.Alternative != nil && len(ifStmt.Alternative.Statements) != tt.alternativeLen {
			t.Errorf("alternative has wrong length for input %q. got=%d, want=%d",
				tt.input, len(ifStmt.Alternative.Statements), tt.alternativeLen)
		}
	}
}

// TestFunctionDefinitionParsing tests parsing of function definitions
func TestFunctionDefinitionParsing(t *testing.T) {
	tests := []struct {
		input        string
		name         string
		paramCount   int
		hasReturnType bool
		returnType   string
		bodyLength   int
	}{
		{
			input:        "def add(x: number, y: number): number do\n\t\t\t\treturn x + y\n\t\t\tend",
			name:         "add",
			paramCount:   2,
			hasReturnType: true,
			returnType:   "number",
			bodyLength:   1,
		},
		{
			input:        "def greet(name: string): string do\n\t\t\t\treturn \"Hello, \" + name\n\t\t\tend",
			name:         "greet",
			paramCount:   1,
			hasReturnType: true,
			returnType:   "string",
			bodyLength:   1,
		},
		{
			input:        "def empty(): void do\n\t\t\tend",
			name:         "empty",
			paramCount:   0,
			hasReturnType: true,
			returnType:   "void",
			bodyLength:   0,
		},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		program, errors := parser.Parse(l)

		if len(errors) != 0 {
			t.Errorf("parser has %d errors for input %q", len(errors), tt.input)
			for _, msg := range errors {
				t.Errorf("parser error: %q", msg)
			}
			continue
		}

		if len(program.Statements) != 1 {
			t.Fatalf("program.Statements does not contain 1 statement. got=%d",
				len(program.Statements))
		}

		stmt, ok := program.Statements[0].(*ast.FunctionDef)
		if !ok {
			t.Fatalf("program.Statements[0] is not ast.FunctionDef. got=%T",
				program.Statements[0])
		}

		if stmt.Name != tt.name {
			t.Errorf("stmt.Name not '%s'. got=%s", tt.name, stmt.Name)
		}

		if len(stmt.Parameters) != tt.paramCount {
			t.Errorf("wrong number of parameters. expected=%d, got=%d",
				tt.paramCount, len(stmt.Parameters))
		}

		if (stmt.ReturnType != nil) != tt.hasReturnType {
			t.Errorf("Return type presence is not %v. got=%v",
				tt.hasReturnType, stmt.ReturnType != nil)
		}

		if stmt.ReturnType != nil && stmt.ReturnType.TypeName != tt.returnType {
			t.Errorf("Return type not '%s'. got=%s", tt.returnType, stmt.ReturnType.TypeName)
		}

		if len(stmt.Body.Statements) != tt.bodyLength {
			t.Errorf("function body has wrong number of statements. expected=%d, got=%d",
				tt.bodyLength, len(stmt.Body.Statements))
		}
	}
}

// TestClassDefinitionParsing tests parsing of class definitions
func TestClassDefinitionParsing(t *testing.T) {
	input := `class Point inherits Object do
		def init(x: number, y: number): void do
			@x = x
			@y = y
		end

		def distance(other: Point): number do
			dx = @x - other.x
			return dx
		end
	end`

	l := lexer.New(input)
	program, errors := parser.Parse(l)

	if len(errors) != 0 {
		t.Errorf("parser has %d errors", len(errors))
		for _, msg := range errors {
			t.Errorf("parser error: %q", msg)
		}
		t.FailNow()
	}

	if len(program.Statements) != 1 {
		t.Fatalf("program.Statements does not contain 1 statement. got=%d",
			len(program.Statements))
	}

	classDef, ok := program.Statements[0].(*ast.ClassDef)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.ClassDef. got=%T",
			program.Statements[0])
	}

	if classDef.Name != "Point" {
		t.Errorf("classDef.Name not 'Point'. got=%q", classDef.Name)
	}

	if classDef.Parent != "Object" {
		t.Errorf("classDef.Parent not 'Object'. got=%q", classDef.Parent)
	}

	if len(classDef.Methods) != 2 {
		t.Fatalf("classDef.Methods does not contain 2 methods. got=%d",
			len(classDef.Methods))
	}
}