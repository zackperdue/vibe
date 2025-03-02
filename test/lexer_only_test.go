package test

import (
	"testing"

	"github.com/example/vibe/lexer"
)

func TestLexerTypeTokenOnly(t *testing.T) {
	// Sample type declaration code
	input := `type StringAlias = string`

	// Create a lexer
	l := lexer.New(input)

	// Check if the TYPE token is recognized
	tok := l.NextToken()
	if tok.Type != lexer.TYPE {
		t.Fatalf("Expected token type to be TYPE, got %s", tok.Type)
	}
	if tok.Literal != "type" {
		t.Fatalf("Expected token literal to be 'type', got %s", tok.Literal)
	}

	// Check the identifier
	tok = l.NextToken()
	if tok.Type != lexer.IDENT {
		t.Fatalf("Expected token type to be IDENT, got %s", tok.Type)
	}
	if tok.Literal != "StringAlias" {
		t.Fatalf("Expected token literal to be 'StringAlias', got %s", tok.Literal)
	}

	// Check the equals sign
	tok = l.NextToken()
	if tok.Type != lexer.ASSIGN {
		t.Fatalf("Expected token type to be ASSIGN, got %s", tok.Type)
	}
	if tok.Literal != "=" {
		t.Fatalf("Expected token literal to be '=', got %s", tok.Literal)
	}

	// Check the type name
	tok = l.NextToken()
	if tok.Type != lexer.IDENT {
		t.Fatalf("Expected token type to be IDENT, got %s", tok.Type)
	}
	if tok.Literal != "string" {
		t.Fatalf("Expected token literal to be 'string', got %s", tok.Literal)
	}

	// Check EOF
	tok = l.NextToken()
	if tok.Type != lexer.EOF {
		t.Fatalf("Expected token type to be EOF, got %s", tok.Type)
	}
}