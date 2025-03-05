package lexer

import (
	"testing"
)

// TestNextToken tests the NextToken function of the lexer
func TestNextTokenDetailed(t *testing.T) {
	input := `five = 5;
ten = 10;
add = def(x, y) {
  x + y;
};
result = add(five, ten);
!-/*5;
5 < 10 > 5;
if (5 < 10) {
  return true;
} else {
  return false;
}
10 == 10;
10 != 9;
"hello world";
"hello" + "world";
[1, 2];
{"foo": "bar"};
for i in [1, 2, 3] {
  let x = i;
};
class Person {
  def init(name) {
    self.name = name;
  }

  def greet() {
    return "Hello, " + self.name;
  }
};
`

	tests := []struct {
		expectedType    TokenType
		expectedLiteral string
	}{
		{IDENT, "five"},
		{ASSIGN, "="},
		{INT, "5"},
		{SEMICOLON, ";"},
		{IDENT, "ten"},
		{ASSIGN, "="},
		{INT, "10"},
		{SEMICOLON, ";"},
		{IDENT, "add"},
		{ASSIGN, "="},
		{FUNCTION, "def"},
		{LPAREN, "("},
		{IDENT, "x"},
		{COMMA, ","},
		{IDENT, "y"},
		{RPAREN, ")"},
		{LBRACE, "{"},
		{IDENT, "x"},
		{PLUS, "+"},
		{IDENT, "y"},
		{SEMICOLON, ";"},
		{RBRACE, "}"},
		{SEMICOLON, ";"},
		{IDENT, "result"},
		{ASSIGN, "="},
		{IDENT, "add"},
		{LPAREN, "("},
		{IDENT, "five"},
		{COMMA, ","},
		{IDENT, "ten"},
		{RPAREN, ")"},
		{SEMICOLON, ";"},
		{BANG, "!"},
		{MINUS, "-"},
		{SLASH, "/"},
		{ASTERISK, "*"},
		{INT, "5"},
		{SEMICOLON, ";"},
		{INT, "5"},
		{LT, "<"},
		{INT, "10"},
		{GT, ">"},
		{INT, "5"},
		{SEMICOLON, ";"},
		{IF, "if"},
		{LPAREN, "("},
		{INT, "5"},
		{LT, "<"},
		{INT, "10"},
		{RPAREN, ")"},
		{LBRACE, "{"},
		{RETURN, "return"},
		{TRUE, "true"},
		{SEMICOLON, ";"},
		{RBRACE, "}"},
		{ELSE, "else"},
		{LBRACE, "{"},
		{RETURN, "return"},
		{FALSE, "false"},
		{SEMICOLON, ";"},
		{RBRACE, "}"},
		{INT, "10"},
		{EQ, "=="},
		{INT, "10"},
		{SEMICOLON, ";"},
		{INT, "10"},
		{NOT_EQ, "!="},
		{INT, "9"},
		{SEMICOLON, ";"},
		{STRING, "hello world"},
		{SEMICOLON, ";"},
		{STRING, "hello"},
		{PLUS, "+"},
		{STRING, "world"},
		{SEMICOLON, ";"},
		{LBRACKET, "["},
		{INT, "1"},
		{COMMA, ","},
		{INT, "2"},
		{RBRACKET, "]"},
		{SEMICOLON, ";"},
		{LBRACE, "{"},
		{STRING, "foo"},
		{COLON, ":"},
		{STRING, "bar"},
		{RBRACE, "}"},
		{SEMICOLON, ";"},
		{FOR, "for"},
		{IDENT, "i"},
		{IN, "in"},
		{LBRACKET, "["},
		{INT, "1"},
		{COMMA, ","},
		{INT, "2"},
		{COMMA, ","},
		{INT, "3"},
		{RBRACKET, "]"},
		{LBRACE, "{"},
		{IDENT, "let"},
		{IDENT, "x"},
		{ASSIGN, "="},
		{IDENT, "i"},
		{SEMICOLON, ";"},
		{RBRACE, "}"},
		{SEMICOLON, ";"},
		{CLASS, "class"},
		{IDENT, "Person"},
		{LBRACE, "{"},
		{FUNCTION, "def"},
		{IDENT, "init"},
		{LPAREN, "("},
		{IDENT, "name"},
		{RPAREN, ")"},
		{LBRACE, "{"},
		{SELF, "self"},
		{DOT, "."},
		{IDENT, "name"},
		{ASSIGN, "="},
		{IDENT, "name"},
		{SEMICOLON, ";"},
		{RBRACE, "}"},
		{FUNCTION, "def"},
		{IDENT, "greet"},
		{LPAREN, "("},
		{RPAREN, ")"},
		{LBRACE, "{"},
		{RETURN, "return"},
		{STRING, "Hello, "},
		{PLUS, "+"},
		{SELF, "self"},
		{DOT, "."},
		{IDENT, "name"},
		{SEMICOLON, ";"},
		{RBRACE, "}"},
		{RBRACE, "}"},
		{SEMICOLON, ";"},
		{EOF, ""},
	}

	l := New(input)

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

// TestReadCharDetailed tests the readChar method of the lexer
func TestReadCharDetailed(t *testing.T) {
	input := "hello"
	l := New(input)

	// Since readChar is not exported, we'll test it indirectly through NextToken
	// This assumes that the first call to NextToken reads all characters of "hello"
	tok := l.NextToken()

	if tok.Type != IDENT {
		t.Fatalf("Expected token type to be IDENT, got=%q", tok.Type)
	}

	if tok.Literal != "hello" {
		t.Fatalf("Expected token literal to be 'hello', got=%q", tok.Literal)
	}

	// The next token should be EOF
	tok = l.NextToken()
	if tok.Type != EOF {
		t.Fatalf("Expected token type to be EOF, got=%q", tok.Type)
	}
}

// TestReadIdentifierDetailed tests the readIdentifier method of the lexer
func TestReadIdentifierDetailed(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"abc", "abc"},
		{"abc123", "abc123"},
		{"abc_123", "abc_123"},
		{"_abc", "_abc"},
		{"_", "_"},
	}

	for i, tt := range tests {
		l := New(tt.input)
		tok := l.NextToken()

		if tok.Type != IDENT {
			t.Fatalf("tests[%d] - Expected token type to be IDENT, got=%q", i, tok.Type)
		}

		if tok.Literal != tt.expected {
			t.Fatalf("tests[%d] - Expected token literal to be %q, got=%q", i, tt.expected, tok.Literal)
		}
	}
}

// TestReadNumberDetailed tests the readNumber method of the lexer
func TestReadNumberDetailed(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"5", "5"},
		{"123", "123"},
		{"0", "0"},
		{"9876543210", "9876543210"},
	}

	for i, tt := range tests {
		l := New(tt.input)
		tok := l.NextToken()

		if tok.Type != INT {
			t.Fatalf("tests[%d] - Expected token type to be INT, got=%q", i, tok.Type)
		}

		if tok.Literal != tt.expected {
			t.Fatalf("tests[%d] - Expected token literal to be %q, got=%q", i, tt.expected, tok.Literal)
		}
	}
}

// TestReadStringDetailed tests the readString method of the lexer
func TestReadStringDetailed(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{`"hello"`, "hello"},
		{`"hello world"`, "hello world"},
		{`"123"`, "123"},
		{`""`, ""},
		{`"hello \"world\""`, `hello "world"`},
	}

	for i, tt := range tests {
		l := New(tt.input)
		tok := l.NextToken()

		if tok.Type != STRING {
			t.Fatalf("tests[%d] - Expected token type to be STRING, got=%q", i, tok.Type)
		}

		if tok.Literal != tt.expected {
			t.Fatalf("tests[%d] - Expected token literal to be %q, got=%q", i, tt.expected, tok.Literal)
		}
	}
}

// TestSkipWhitespaceDetailed tests the skipWhitespace method of the lexer
func TestSkipWhitespaceDetailed(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"  5  ", "5"},
		{"\t\n123\r", "123"},
		{" \n\t\r0 ", "0"},
	}

	for i, tt := range tests {
		l := New(tt.input)
		tok := l.NextToken()

		if tok.Type != INT {
			t.Fatalf("tests[%d] - Expected token type to be INT, got=%q", i, tok.Type)
		}

		if tok.Literal != tt.expected {
			t.Fatalf("tests[%d] - Expected token literal to be %q, got=%q", i, tt.expected, tok.Literal)
		}
	}
}