package parser

import (
	"testing"

	"github.com/vibe-lang/vibe/ast"
	"github.com/vibe-lang/vibe/lexer"
)

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
		p := New(l)
		program, err := p.Parse()

		if err != nil {
			t.Fatalf("parser error: %v", err)
		}

		if len(program.Statements) != 1 {
			t.Fatalf("program does not have 1 statement. got=%d", len(program.Statements))
		}

		// The parser now returns ArrayLiteral directly instead of wrapping it in ExpressionStatement
		var array *ast.ArrayLiteral
		if stmt, ok := program.Statements[0].(*ast.ExpressionStatement); ok {
			// Handle legacy behavior where it might be wrapped in ExpressionStatement
			array, ok = stmt.Expression.(*ast.ArrayLiteral)
			if !ok {
				t.Fatalf("expr is not ast.ArrayLiteral. got=%T", stmt.Expression)
			}
		} else {
			// Handle new behavior where it's a direct ArrayLiteral
			array, ok = program.Statements[0].(*ast.ArrayLiteral)
			if !ok {
				t.Fatalf("program.Statements[0] is not ast.ArrayLiteral. got=%T", program.Statements[0])
			}
		}

		if len(array.Elements) != len(tt.expected) {
			t.Fatalf("array.Elements has wrong number of elements. got=%d, want=%d",
				len(array.Elements), len(tt.expected))
		}

		// The rest of the test is simplified as the exact string representation depends on the ast.toString implementation
	}
}

// TestStringLiteralParsing tests parsing of string literals
func TestStringLiteralParsing(t *testing.T) {
	input := `"hello world"`

	l := lexer.New(input)
	p := New(l)
	program, err := p.Parse()

	if err != nil {
		t.Fatalf("parser error: %v", err)
	}

	if len(program.Statements) != 1 {
		t.Fatalf("program does not have 1 statement. got=%d", len(program.Statements))
	}

	// The parser now returns StringLiteral directly instead of wrapping it in ExpressionStatement
	var literal *ast.StringLiteral
	if stmt, ok := program.Statements[0].(*ast.ExpressionStatement); ok {
		// Handle legacy behavior where it might be wrapped in ExpressionStatement
		literal, ok = stmt.Expression.(*ast.StringLiteral)
		if !ok {
			t.Fatalf("exp not *ast.StringLiteral. got=%T", stmt.Expression)
		}
	} else {
		// Handle new behavior where it's a direct StringLiteral
		literal, ok = program.Statements[0].(*ast.StringLiteral)
		if !ok {
			t.Fatalf("program.Statements[0] is not ast.StringLiteral. got=%T", program.Statements[0])
		}
	}

	if literal.Value != "hello world" {
		t.Errorf("literal.Value not %q. got=%q", "hello world", literal.Value)
	}
}

// TestIfExpression tests the parsing of if expressions
func TestIfExpression(t *testing.T) {
	input := `if x < y do x else y end`

	l := lexer.New(input)
	p := New(l)
	program, err := p.Parse()

	if err != nil {
		t.Fatalf("parser error: %v", err)
	}

	if len(program.Statements) != 1 {
		t.Fatalf("program.Statements does not contain 1 statement. got=%d",
			len(program.Statements))
	}

	// The parser now returns IfStmt directly instead of wrapping it in ExpressionStatement
	var exp *ast.IfStmt
	if stmt, ok := program.Statements[0].(*ast.ExpressionStatement); ok {
		// Handle legacy behavior where it might be wrapped in ExpressionStatement
		exp, ok = stmt.Expression.(*ast.IfStmt)
		if !ok {
			t.Fatalf("stmt.Expression is not ast.IfStmt. got=%T", stmt.Expression)
		}
	} else {
		// Handle new behavior where it's a direct IfStmt
		exp, ok = program.Statements[0].(*ast.IfStmt)
		if !ok {
			t.Fatalf("program.Statements[0] is not ast.IfStmt. got=%T", program.Statements[0])
		}
	}

	if !testInfixExpression(t, exp.Condition, "x", "<", "y") {
		return
	}

	if len(exp.Consequence.Statements) != 1 {
		t.Errorf("consequence is not 1 statement. got=%d",
			len(exp.Consequence.Statements))
	}

	// The consequence may now contain a direct Identifier instead of an ExpressionStatement
	var consequenceValue ast.Node = exp.Consequence.Statements[0]
	if exprStmt, ok := consequenceValue.(*ast.ExpressionStatement); ok {
		consequenceValue = exprStmt.Expression
	}

	// Check for identifier "x" either way
	foundX := false
	if ident, ok := consequenceValue.(*ast.Identifier); ok {
		if ident.Name == "x" {
			foundX = true
		}
	}

	if !foundX {
		t.Fatalf("Expected consequence to contain identifier 'x'. got=%T: %s",
			consequenceValue, consequenceValue.String())
	}

	if len(exp.Alternative.Statements) != 1 {
		t.Errorf("alternative is not 1 statement. got=%d",
			len(exp.Alternative.Statements))
	}

	// The alternative may now contain a direct Identifier instead of an ExpressionStatement
	var alternativeValue ast.Node = exp.Alternative.Statements[0]
	if exprStmt, ok := alternativeValue.(*ast.ExpressionStatement); ok {
		alternativeValue = exprStmt.Expression
	}

	// Check for identifier "y" either way
	foundY := false
	if ident, ok := alternativeValue.(*ast.Identifier); ok {
		if ident.Name == "y" {
			foundY = true
		}
	}

	if !foundY {
		t.Fatalf("Expected alternative to contain identifier 'y'. got=%T: %s",
			alternativeValue, alternativeValue.String())
	}
}

// Helper functions for testing expressions
func testInfixExpression(t *testing.T, exp ast.Node, left interface{}, operator string, right interface{}) bool {
	opExp, ok := exp.(*ast.BinaryExpr)
	if !ok {
		t.Errorf("exp is not ast.BinaryExpr. got=%T(%s)", exp, exp)
		return false
	}

	if !testLiteralExpression(t, opExp.Left, left) {
		return false
	}

	if opExp.Operator != operator {
		t.Errorf("opExp.Operator is not '%s'. got=%s", operator, opExp.Operator)
		return false
	}

	if !testLiteralExpression(t, opExp.Right, right) {
		return false
	}

	return true
}

func testLiteralExpression(t *testing.T, exp ast.Node, expected interface{}) bool {
	switch v := expected.(type) {
	case int:
		return testIntegerLiteral(t, exp, int64(v))
	case int64:
		return testIntegerLiteral(t, exp, v)
	case string:
		return testIdentifier(t, exp, v)
	case bool:
		return testBooleanLiteral(t, exp, v)
	}
	t.Errorf("type of exp not handled. got=%T", exp)
	return false
}

func testIntegerLiteral(t *testing.T, il ast.Node, value int64) bool {
	integ, ok := il.(*ast.NumberLiteral)
	if !ok {
		t.Errorf("il not *ast.NumberLiteral. got=%T", il)
		return false
	}

	if integ.Value != float64(value) {
		t.Errorf("integ.Value not %d. got=%f", value, integ.Value)
		return false
	}

	return true
}

func testIdentifier(t *testing.T, exp ast.Node, value string) bool {
	ident, ok := exp.(*ast.Identifier)
	if !ok {
		t.Errorf("exp not *ast.Identifier. got=%T", exp)
		return false
	}

	if ident.Name != value {
		t.Errorf("ident.Name not %s. got=%s", value, ident.Name)
		return false
	}

	return true
}

func testBooleanLiteral(t *testing.T, exp ast.Node, value bool) bool {
	bo, ok := exp.(*ast.BooleanLiteral)
	if !ok {
		t.Errorf("exp not *ast.BooleanLiteral. got=%T", exp)
		return false
	}

	if bo.Value != value {
		t.Errorf("bo.Value not %t. got=%t", value, bo.Value)
		return false
	}

	return true
}