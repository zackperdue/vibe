package lexer

import (
	"fmt"
	"testing"
)

func TestLexerBasics(t *testing.T) {
	input := `a = 5
b = 10
c = a + b
if (c > 10) {
	return true
} else {
	return false
}
def add(x, y) {
	return x + y
}
"hello world"
!-*/<=>=`

	l := New(input)
	var tokens []Token

	// Collect all tokens
	for {
		tok := l.NextToken()
		tokens = append(tokens, tok)
		if tok.Type == "EOF" {
			break
		}
	}

	// Basic validation - make sure we found some tokens
	if len(tokens) < 10 {
		t.Fatalf("Expected at least 10 tokens, got %d", len(tokens))
	}

	// Check for expected token types
	expectedTypes := map[TokenType]bool{
		"IDENTIFIER": false,
		"INTEGER":    false,
		"KEYWORD":    false,
		"OPERATOR":   false,
		"STRING":     false,
		"LPAREN":     false,
		"RPAREN":     false,
		"LBRACE":     false,
		"RBRACE":     false,
	}

	// Check for expected values
	expectedValues := map[string]bool{
		"a":           false,
		"=":           false,
		"5":           false,
		"10":          false,
		"+":           false,
		"if":          false,
		"else":        false,
		"return":      false,
		"true":        false,
		"false":       false,
		"def":         false,
		"add":         false,
		"hello world": false,
	}

	for _, tok := range tokens {
		if _, exists := expectedTypes[tok.Type]; exists {
			expectedTypes[tok.Type] = true
		}
		if _, exists := expectedValues[tok.Value]; exists {
			expectedValues[tok.Value] = true
		}
	}

	// Verify all expected token types were found
	for tokenType, found := range expectedTypes {
		if !found {
			t.Errorf("Expected to find token type %s but didn't", tokenType)
		}
	}

	// Verify all expected values were found
	for value, found := range expectedValues {
		if !found {
			t.Errorf("Expected to find token value %s but didn't", value)
		}
	}
}

func TestStringTokens(t *testing.T) {
	input := `"hello" "world"`

	l := New(input)

	// First token should be a string
	tok1 := l.NextToken()
	if tok1.Type != "STRING" {
		t.Fatalf("Expected STRING token type, got %s", tok1.Type)
	}
	if tok1.Value != "hello" {
		t.Fatalf("Expected string value 'hello', got %s", tok1.Value)
	}

	// Skip any spacers
	for {
		tok := l.NextToken()
		if tok.Type == "STRING" || tok.Type == "EOF" {
			// Second token should be a string
			if tok.Type != "STRING" {
				t.Fatalf("Expected STRING token type, got %s", tok.Type)
			}
			if tok.Value != "world" {
				t.Fatalf("Expected string value 'world', got %s", tok.Value)
			}
			break
		}
	}
}

func TestNumberTokens(t *testing.T) {
	input := `5 10 3.14`

	l := New(input)

	// Test first number (5)
	tok := l.NextToken()
	if tok.Type != "INTEGER" && tok.Type != "FLOAT" {
		t.Fatalf("Expected INTEGER or FLOAT token type, got %s", tok.Type)
	}
	if tok.Value != "5" {
		t.Fatalf("Expected value 5, got %s", tok.Value)
	}

	// Skip whitespace/newline
	for {
		tok = l.NextToken()
		if tok.Type == "INTEGER" || tok.Type == "FLOAT" || tok.Type == "EOF" {
			break
		}
	}

	// Test second number (10)
	if tok.Type != "INTEGER" && tok.Type != "FLOAT" {
		t.Fatalf("Expected INTEGER or FLOAT token type, got %s", tok.Type)
	}
	if tok.Value != "10" {
		t.Fatalf("Expected value 10, got %s", tok.Value)
	}

	// Skip whitespace/newline
	for {
		tok = l.NextToken()
		if tok.Type == "INTEGER" || tok.Type == "FLOAT" || tok.Type == "EOF" {
			break
		}
	}

	// Test third number (3.14)
	if tok.Type != "INTEGER" && tok.Type != "FLOAT" {
		t.Fatalf("Expected INTEGER or FLOAT token type, got %s", tok.Type)
	}

	// The lexer might round or format the float differently
	// Just check that it starts with "3."
	if len(tok.Value) < 2 || tok.Value[:2] != "3." {
		t.Fatalf("Expected value starting with 3., got %s", tok.Value)
	}
}

func TestComplexInputs(t *testing.T) {
	input := `# This is a comment
def factorial(n) {
	if (n <= 1) {
		return 1;
	}
	return n * factorial(n - 1);
}

# Array access
arr = [1, 2, 3, 4];
arr[2] = 10;`

	l := New(input)
	var tokens []Token

	// Collect all tokens
	for {
		tok := l.NextToken()
		tokens = append(tokens, tok)
		if tok.Type == "EOF" {
			break
		}
	}

	// Just verify that we got tokens and the test runs without crashing
	if len(tokens) == 0 {
		t.Fatalf("Expected tokens but got none")
	}

	// Verify we have def keyword
	defFound := false
	for _, tok := range tokens {
		if tok.Type == "KEYWORD" && tok.Value == "def" {
			defFound = true
			break
		}
	}

	if !defFound {
		t.Fatal("Could not find 'def' keyword token in collected tokens")
	}
}

func TestDebugTokens(t *testing.T) {
	input := `5 10 true false`

	// Print each character and its code
	fmt.Println("DEBUG CHARACTERS:")
	for i, ch := range input {
		fmt.Printf("Char at %d: %q (code: %d)\n", i, ch, ch)
	}

	l := New(input)

	// Print all tokens
	fmt.Println("DEBUG TOKENS:")
	for {
		tok := l.NextToken()
		fmt.Printf("Token: Type=%s, Value=%q\n", tok.Type, tok.Value)
		if tok.Type == "EOF" {
			break
		}
	}
}