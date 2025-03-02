package lexer_test

import (
	"testing"

	"github.com/example/vibe/lexer"
)

// TestBasicTokenization tests the lexer's ability to tokenize basic constructs
func TestBasicTokenization(t *testing.T) {
	input := `a = 5
b = 10
c = a + b
if (c > 10) {
	return true
} else {
	return false
}
"hello world"
for x in [1, 2, 3] do
	print(x)
end`

	tests := []struct {
		expectedType    lexer.TokenType
		expectedLiteral string
	}{
		{lexer.IDENT, "a"},
		{lexer.ASSIGN, "="},
		{lexer.INT, "5"},
		{lexer.IDENT, "b"},
		{lexer.ASSIGN, "="},
		{lexer.INT, "10"},
		{lexer.IDENT, "c"},
		{lexer.ASSIGN, "="},
		{lexer.IDENT, "a"},
		{lexer.PLUS, "+"},
		{lexer.IDENT, "b"},
		{lexer.IF, "if"},
		{lexer.LPAREN, "("},
		{lexer.IDENT, "c"},
		{lexer.GT, ">"},
		{lexer.INT, "10"},
		{lexer.RPAREN, ")"},
		{lexer.LBRACE, "{"},
		{lexer.RETURN, "return"},
		{lexer.TRUE, "true"},
		{lexer.RBRACE, "}"},
		{lexer.ELSE, "else"},
		{lexer.LBRACE, "{"},
		{lexer.RETURN, "return"},
		{lexer.FALSE, "false"},
		{lexer.RBRACE, "}"},
		{lexer.STRING, "hello world"},
		{lexer.FOR, "for"},
		{lexer.IDENT, "x"},
		{lexer.IN, "in"},
		{lexer.LBRACKET, "["},
		{lexer.INT, "1"},
		{lexer.COMMA, ","},
		{lexer.INT, "2"},
		{lexer.COMMA, ","},
		{lexer.INT, "3"},
		{lexer.RBRACKET, "]"},
		{lexer.DO, "do"},
		{lexer.IDENT, "print"},
		{lexer.LPAREN, "("},
		{lexer.IDENT, "x"},
		{lexer.RPAREN, ")"},
		{lexer.END, "end"},
		{lexer.EOF, ""},
	}

	l := lexer.New(input)

	for i, tt := range tests {
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

// TestTypeDeclarations tests the lexer's ability to tokenize type declarations
func TestTypeDeclarations(t *testing.T) {
	input := `
		type StringAlias = string
		type StringArray = Array<string>
        type NumberMap = Map<string, number>

		a: StringAlias = "Hello"
		b: StringArray = ["World", "!"]
	`

	tests := []struct {
		expectedType    lexer.TokenType
		expectedLiteral string
	}{
		{lexer.TYPE, "type"},
		{lexer.IDENT, "StringAlias"},
		{lexer.ASSIGN, "="},
		{lexer.IDENT, "string"},

		{lexer.TYPE, "type"},
		{lexer.IDENT, "StringArray"},
		{lexer.ASSIGN, "="},
		{lexer.IDENT, "Array"},
		{lexer.LT, "<"},
		{lexer.IDENT, "string"},
		{lexer.GT, ">"},

		{lexer.TYPE, "type"},
		{lexer.IDENT, "NumberMap"},
		{lexer.ASSIGN, "="},
		{lexer.IDENT, "Map"},
		{lexer.LT, "<"},
		{lexer.IDENT, "string"},
		{lexer.COMMA, ","},
		{lexer.IDENT, "number"},
		{lexer.GT, ">"},

		{lexer.IDENT, "a"},
		{lexer.COLON, ":"},
		{lexer.IDENT, "StringAlias"},
		{lexer.ASSIGN, "="},
		{lexer.STRING, "Hello"},

		{lexer.IDENT, "b"},
		{lexer.COLON, ":"},
		{lexer.IDENT, "StringArray"},
		{lexer.ASSIGN, "="},
		{lexer.LBRACKET, "["},
		{lexer.STRING, "World"},
		{lexer.COMMA, ","},
		{lexer.STRING, "!"},
		{lexer.RBRACKET, "]"},

		{lexer.EOF, ""},
	}

	l := lexer.New(input)

	for i, tt := range tests {
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

// TestOperators tests the lexer's ability to tokenize operators
func TestOperators(t *testing.T) {
	input := `+ - * / %
    == != < > <= >=
    && || !
    += -= *= /= %=`

	tests := []struct {
		expectedType    lexer.TokenType
		expectedLiteral string
	}{
		{lexer.PLUS, "+"},
		{lexer.MINUS, "-"},
		{lexer.ASTERISK, "*"},
		{lexer.SLASH, "/"},
		{lexer.MODULO, "%"},

		{lexer.EQ, "=="},
		{lexer.NOT_EQ, "!="},
		{lexer.LT, "<"},
		{lexer.GT, ">"},
		{lexer.LT_EQ, "<="},
		{lexer.GT_EQ, ">="},

		{lexer.AND, "&&"},
		{lexer.OR, "||"},
		{lexer.BANG, "!"},

		{lexer.PLUS_ASSIGN, "+="},
		{lexer.MINUS_ASSIGN, "-="},
		{lexer.MUL_ASSIGN, "*="},
		{lexer.DIV_ASSIGN, "/="},
		{lexer.MOD_ASSIGN, "%="},

		{lexer.EOF, ""},
	}

	l := lexer.New(input)

	for i, tt := range tests {
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

// TestNumbers tests the lexer's ability to tokenize integer and floating-point numbers
func TestNumbers(t *testing.T) {
	input := `5 10 3.14 42.0 0.123 123.456`

	tests := []struct {
		expectedType    lexer.TokenType
		expectedLiteral string
	}{
		{lexer.INT, "5"},
		{lexer.INT, "10"},
		{lexer.FLOAT, "3.14"},
		{lexer.FLOAT, "42.0"},
		{lexer.FLOAT, "0.123"},
		{lexer.FLOAT, "123.456"},
		{lexer.EOF, ""},
	}

	l := lexer.New(input)

	for i, tt := range tests {
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

// TestStrings tests the lexer's ability to tokenize string literals
func TestStrings(t *testing.T) {
	input := `"hello" "world" "hello world" "" "123" "!@#$%^&*()"`

	tests := []struct {
		expectedType    lexer.TokenType
		expectedLiteral string
	}{
		{lexer.STRING, "hello"},
		{lexer.STRING, "world"},
		{lexer.STRING, "hello world"},
		{lexer.STRING, ""},
		{lexer.STRING, "123"},
		{lexer.STRING, "!@#$%^&*()"},
		{lexer.EOF, ""},
	}

	l := lexer.New(input)

	for i, tt := range tests {
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

// TestKeywords tests the lexer's ability to tokenize keywords
func TestKeywords(t *testing.T) {
	input := `if else elsif for while do end def return class inherits self super new let var true false nil type require puts in`

	tests := []struct {
		expectedType    lexer.TokenType
		expectedLiteral string
	}{
		{lexer.IF, "if"},
		{lexer.ELSE, "else"},
		{lexer.ELSIF, "elsif"},
		{lexer.FOR, "for"},
		{lexer.WHILE, "while"},
		{lexer.DO, "do"},
		{lexer.END, "end"},
		{lexer.FUNCTION, "def"},
		{lexer.RETURN, "return"},
		{lexer.CLASS, "class"},
		{lexer.INHERITS, "inherits"},
		{lexer.SELF, "self"},
		{lexer.SUPER, "super"},
		{lexer.NEW, "new"},
		{lexer.LET, "let"},
		{lexer.VAR, "var"},
		{lexer.TRUE, "true"},
		{lexer.FALSE, "false"},
		{lexer.NIL, "nil"},
		{lexer.TYPE, "type"},
		{lexer.REQUIRE, "require"},
		{lexer.PRINT, "puts"},
		{lexer.IN, "in"},
		{lexer.EOF, ""},
	}

	l := lexer.New(input)

	for i, tt := range tests {
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

// TestLineAndColumnTracking tests that the lexer correctly tracks line and column positions
func TestLineAndColumnTracking(t *testing.T) {
	input := `let x = 5
y = 10`

	tests := []struct {
		expectedType    lexer.TokenType
		expectedLiteral string
		expectedLine    int
		expectedColumn  int
	}{
		{lexer.LET, "let", 1, 1},
		{lexer.IDENT, "x", 1, 5},
		{lexer.ASSIGN, "=", 1, 7},
		{lexer.INT, "5", 1, 8},
		{lexer.IDENT, "y", 2, 1},
		{lexer.ASSIGN, "=", 2, 3},
		{lexer.INT, "10", 2, 4},
		{lexer.EOF, "", 2, 7},
	}

	l := lexer.New(input)

	for i, tt := range tests {
		tok := l.NextToken()

		if tok.Type != tt.expectedType {
			t.Fatalf("tests[%d] - tokentype wrong. expected=%q, got=%q",
				i, tt.expectedType, tok.Type)
		}

		if tok.Literal != tt.expectedLiteral {
			t.Fatalf("tests[%d] - literal wrong. expected=%q, got=%q",
				i, tt.expectedLiteral, tok.Literal)
		}

		if tok.Line != tt.expectedLine {
			t.Fatalf("tests[%d] - line wrong. expected=%d, got=%d",
				i, tt.expectedLine, tok.Line)
		}

		if tok.Column != tt.expectedColumn {
			t.Fatalf("tests[%d] - column wrong. expected=%d, got=%d",
				i, tt.expectedColumn, tok.Column)
		}
	}
}

// TestComplexCode tests the lexer with a more complex code sample
func TestComplexCode(t *testing.T) {
	input := `
# This is a comment
function fibonacci(n: number): number do
    if n <= 1 do
        return n
    end

    return fibonacci(n - 1) + fibonacci(n - 2)
end

class Point inherits Object do
    function init(x: number, y: number) do
        @x = x
        @y = y
    end

    function distance(other: Point): number do
        dx = @x - other.x
        dy = @y - other.y
        return (dx * dx + dy * dy) ** 0.5
    end
end

p1 = new Point(3, 4)
p2 = new Point(0, 0)
print("Distance: " + p1.distance(p2).to_string())
`

	l := lexer.New(input)

	// Just test that we can tokenize without errors
	tokenCount := 0
	for {
		tok := l.NextToken()
		tokenCount++

		if tok.Type == lexer.EOF {
			break
		}

		// Check that no tokens have empty type (except EOF)
		if tok.Type == "" {
			t.Fatalf("Token has empty type at line %d, column %d", tok.Line, tok.Column)
		}
	}

	// We should have tokenized a significant number of tokens
	if tokenCount < 50 {
		t.Fatalf("Expected at least 50 tokens, got %d", tokenCount)
	}
}