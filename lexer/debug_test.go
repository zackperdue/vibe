package lexer

import (
	"fmt"
	"testing"
)

func TestFalseToken(t *testing.T) {
	input := "false"  // Using double quotes instead of backticks

	// Print each character and its code
	fmt.Println("DEBUG INPUT CHARACTERS:")
	for i, ch := range input {
		fmt.Printf("Char at %d: %q (code: %d)\n", i, ch, ch)
	}

	l := New(input)

	// Print each character as it's read
	fmt.Println("DEBUG FALSE TOKEN:")

	// Print the token
	tok := l.NextToken()
	fmt.Printf("Token: Type=%s, Value=%q\n", tok.Type, tok.Literal)

	if tok.Literal != "false" {
		t.Errorf("Expected token value to be 'false', got %q", tok.Literal)
	}

	if tok.Type != FALSE {
		t.Errorf("Expected token type to be FALSE, got %s", tok.Type)
	}
}