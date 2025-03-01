package lexer

import (
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
		if tok.Type == EOF {
			break
		}
	}

	// Basic validation - make sure we found some tokens
	if len(tokens) < 10 {
		t.Fatalf("Expected at least 10 tokens, got %d", len(tokens))
	}

	// Check for expected token types
	expectedTypes := map[TokenType]bool{
		IDENT:    false,
		INT:      false,
		FUNCTION: false,
		ASSIGN:   false,
		STRING:   false,
		LPAREN:   false,
		RPAREN:   false,
		LBRACE:   false,
		RBRACE:   false,
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
		if _, exists := expectedValues[tok.Literal]; exists {
			expectedValues[tok.Literal] = true
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
	if tok1.Type != STRING {
		t.Fatalf("Expected STRING token type, got %s", tok1.Type)
	}
	if tok1.Literal != "hello" {
		t.Fatalf("Expected string value 'hello', got %s", tok1.Literal)
	}

	// Second token should be a string
	tok2 := l.NextToken()
	if tok2.Type != STRING {
		t.Fatalf("Expected STRING token type, got %s", tok2.Type)
	}
	if tok2.Literal != "world" {
		t.Fatalf("Expected string value 'world', got %s", tok2.Literal)
	}
}

func TestNumberTokens(t *testing.T) {
	input := `5 10 3.14`

	l := New(input)

	// Test first number (5)
	tok := l.NextToken()
	if tok.Type != INT {
		t.Fatalf("Expected INT token type, got %s", tok.Type)
	}
	if tok.Literal != "5" {
		t.Fatalf("Expected value 5, got %s", tok.Literal)
	}

	// Test second number (10)
	tok = l.NextToken()
	if tok.Type != INT {
		t.Fatalf("Expected INT token type, got %s", tok.Type)
	}
	if tok.Literal != "10" {
		t.Fatalf("Expected value 10, got %s", tok.Literal)
	}

	// Test third number (3.14)
	tok = l.NextToken()
	if tok.Type != FLOAT {
		t.Fatalf("Expected FLOAT token type, got %s", tok.Type)
	}
	if tok.Literal != "3.14" {
		t.Fatalf("Expected value 3.14, got %s", tok.Literal)
	}
}

func TestKeywords(t *testing.T) {
	input := `def let var true false if else elsif return while nil print`

	l := New(input)

	expectedTokens := []TokenType{
		FUNCTION, LET, VAR, TRUE, FALSE, IF, ELSE, ELSIF, RETURN, WHILE, NIL, PRINT,
	}

	for i, expected := range expectedTokens {
		tok := l.NextToken()
		if tok.Type != expected {
			t.Fatalf("Token %d: expected %s, got %s", i, expected, tok.Type)
		}
	}
}

func TestOperators(t *testing.T) {
	input := `= + - ! * / < > == != <= >= && ||`

	l := New(input)

	expectedTokens := []TokenType{
		ASSIGN, PLUS, MINUS, BANG, ASTERISK, SLASH,
		LT, GT, EQ, NOT_EQ, LT_EQ, GT_EQ, AND, OR,
	}

	for i, expected := range expectedTokens {
		tok := l.NextToken()
		if tok.Type != expected {
			t.Fatalf("Token %d: expected %s, got %s", i, expected, tok.Type)
		}
	}
}

func TestDelimiters(t *testing.T) {
	input := `, ; : . ( ) { } [ ]`

	l := New(input)

	expectedTokens := []TokenType{
		COMMA, SEMICOLON, COLON, DOT, LPAREN, RPAREN, LBRACE, RBRACE, LBRACKET, RBRACKET,
	}

	for i, expected := range expectedTokens {
		tok := l.NextToken()
		if tok.Type != expected {
			t.Fatalf("Token %d: expected %s, got %s", i, expected, tok.Type)
		}
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
		if tok.Type == EOF {
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
		if tok.Type == FUNCTION && tok.Literal == "def" {
			defFound = true
			break
		}
	}

	if !defFound {
		t.Fatal("Could not find 'def' keyword token in collected tokens")
	}
}