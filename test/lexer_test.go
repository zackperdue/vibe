package test

import (
	"testing"

	"github.com/example/vibe/lexer"
)

func TestLexerTypeToken(t *testing.T) {
	// Sample type declaration code
	input := `
		type StringAlias = string
		type StringArray = Array<string>

		a: StringAlias = "Hello"
		b: StringArray = ["World", "!"]
	`

	// Create a lexer
	l := lexer.New(input)

	// Test tokens
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