package lexer

import (
	"fmt"
)

// TokenType identifies the type of token
type TokenType string

// Token represents a lexical token
type Token struct {
	Type    TokenType
	Literal string
	Line    int
	Column  int
}

// Define TokenTypes
const (
	ILLEGAL = "ILLEGAL" // Illegal token
	EOF     = "EOF"     // End of file

	// Identifiers and literals
	IDENT  = "IDENT"  // Variable and function names
	INT    = "INT"    // Integer literals
	FLOAT  = "FLOAT"  // Floating point literals
	STRING = "STRING" // String literals

	// Operators
	ASSIGN   = "="
	PLUS     = "+"
	MINUS    = "-"
	BANG     = "!"
	ASTERISK = "*"
	SLASH    = "/"
	MODULO   = "%"

	LT = "<"
	GT = ">"

	EQ     = "=="
	NOT_EQ = "!="
	LT_EQ  = "<="
	GT_EQ  = ">="

	AND = "&&"
	OR  = "||"

	// Delimiters
	COMMA     = ","
	SEMICOLON = ";"
	COLON     = ":"
	DOT       = "."
	AT        = "@"  // For instance variables

	LPAREN   = "("
	RPAREN   = ")"
	LBRACE   = "{"
	RBRACE   = "}"
	LBRACKET = "["
	RBRACKET = "]"

	// Keywords
	FUNCTION = "FUNCTION"
	LET      = "LET"
	VAR      = "VAR"
	TRUE     = "TRUE"
	FALSE    = "FALSE"
	IF       = "IF"
	ELSE     = "ELSE"
	ELSIF    = "ELSIF"
	RETURN   = "RETURN"
	WHILE    = "WHILE"
	FOR      = "FOR"
	IN       = "IN"
	NIL      = "NIL"
	PRINT    = "PRINT"
	END      = "END"
	DO       = "DO"

	// Class-related keywords
	CLASS    = "CLASS"
	INHERITS = "INHERITS"
	SELF     = "SELF"
	SUPER    = "SUPER"
	NEW      = "NEW"

	// Compound assignment operators
	PLUS_ASSIGN   = "+="
	MINUS_ASSIGN  = "-="
	MUL_ASSIGN    = "*="
	DIV_ASSIGN    = "/="
	MOD_ASSIGN    = "%="
)

// keywords maps strings to their keyword TokenType
var keywords = map[string]TokenType{
	"def":      FUNCTION,
	"let":      LET,
	"var":      VAR,
	"true":     TRUE,
	"false":    FALSE,
	"if":       IF,
	"else":     ELSE,
	"elsif":    ELSIF,
	"return":   RETURN,
	"while":    WHILE,
	"for":      FOR,
	"in":       IN,
	"nil":      NIL,
	"print":    PRINT,
	"puts":     PRINT,
	"end":      END,
	"do":       DO,

	// Class-related keywords
	"class":    CLASS,
	"inherits": INHERITS,
	"self":     SELF,
	"super":    SUPER,
	"new":      NEW,
}

// Lexer analyzes the input and breaks it up into tokens
type Lexer struct {
	input        string
	position     int  // current position in input (points to current char)
	readPosition int  // current reading position in input (after current char)
	ch           byte // current character being examined
	line         int  // current line number
	column       int  // current column number
}

// New creates a new Lexer
func New(input string) *Lexer {
	l := &Lexer{input: input, line: 1, column: 0}
	l.readChar()
	return l
}

// readChar reads the next character and advances the position in the input string
func (l *Lexer) readChar() {
	if l.readPosition >= len(l.input) {
		l.ch = 0
	} else {
		l.ch = l.input[l.readPosition]
	}
	l.position = l.readPosition
	l.readPosition++
	l.column++

	// If we just read a newline, increment line counter and reset column
	if l.ch == '\n' {
		l.line++
		l.column = 0
	}
}

// peekChar looks at the next character without advancing the position
func (l *Lexer) peekChar() byte {
	if l.readPosition >= len(l.input) {
		return 0
	}
	return l.input[l.readPosition]
}

// NextToken returns the next token in the input
func (l *Lexer) NextToken() Token {
	var tok Token

	l.skipWhitespace()

	// Remember the starting position of the token
	line := l.line
	column := l.column

	switch l.ch {
	case '=':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			tok = Token{Type: EQ, Literal: string(ch) + string(l.ch)}
		} else {
			tok = newToken(ASSIGN, l.ch)
		}
	case '+':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			tok = Token{Type: PLUS_ASSIGN, Literal: string(ch) + string(l.ch)}
		} else {
			tok = newToken(PLUS, l.ch)
		}
	case '-':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			tok = Token{Type: MINUS_ASSIGN, Literal: string(ch) + string(l.ch)}
		} else {
			tok = newToken(MINUS, l.ch)
		}
	case '!':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			tok = Token{Type: NOT_EQ, Literal: string(ch) + string(l.ch)}
		} else {
			tok = newToken(BANG, l.ch)
		}
	case '*':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			tok = Token{Type: MUL_ASSIGN, Literal: string(ch) + string(l.ch)}
		} else {
			tok = newToken(ASTERISK, l.ch)
		}
	case '/':
		// Check for comments
		if l.peekChar() == '/' {
			// Skip the rest of the line (comment)
			l.readChar() // consume the second '/'
			for l.ch != '\n' && l.ch != 0 {
				l.readChar()
			}
			return l.NextToken() // Get the next valid token
		} else if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			tok = Token{Type: DIV_ASSIGN, Literal: string(ch) + string(l.ch)}
		} else {
			tok = newToken(SLASH, l.ch)
		}
	case '#':
		// Skip the rest of the line (comment)
		for l.ch != '\n' && l.ch != 0 {
			l.readChar()
		}
		return l.NextToken() // Get the next valid token
	case '%':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			tok = Token{Type: MOD_ASSIGN, Literal: string(ch) + string(l.ch)}
		} else {
			tok = newToken(MODULO, l.ch)
		}
	case '<':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			tok = Token{Type: LT_EQ, Literal: string(ch) + string(l.ch)}
		} else {
			tok = newToken(LT, l.ch)
		}
	case '>':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			tok = Token{Type: GT_EQ, Literal: string(ch) + string(l.ch)}
		} else {
			tok = newToken(GT, l.ch)
		}
	case '&':
		if l.peekChar() == '&' {
			ch := l.ch
			l.readChar()
			tok = Token{Type: AND, Literal: string(ch) + string(l.ch)}
		} else {
			tok = newToken(ILLEGAL, l.ch)
		}
	case '|':
		if l.peekChar() == '|' {
			ch := l.ch
			l.readChar()
			tok = Token{Type: OR, Literal: string(ch) + string(l.ch)}
		} else {
			tok = newToken(ILLEGAL, l.ch)
		}
	case ',':
		tok = newToken(COMMA, l.ch)
	case ';':
		tok = newToken(SEMICOLON, l.ch)
	case ':':
		tok = newToken(COLON, l.ch)
	case '.':
		tok = newToken(DOT, l.ch)
	case '@':
		tok = newToken(AT, l.ch)
	case '(':
		tok = newToken(LPAREN, l.ch)
	case ')':
		tok = newToken(RPAREN, l.ch)
	case '{':
		tok = newToken(LBRACE, l.ch)
	case '}':
		tok = newToken(RBRACE, l.ch)
	case '[':
		tok = newToken(LBRACKET, l.ch)
	case ']':
		tok = newToken(RBRACKET, l.ch)
	case '"':
		tok.Type = STRING
		tok.Literal = l.readString()
	case 0:
		tok.Type = EOF
		tok.Literal = ""
	default:
		if isLetter(l.ch) {
			tok.Literal = l.readIdentifier()
			tok.Type = lookupIdent(tok.Literal)
			fmt.Printf("DEBUG: NextToken - identifier: %s, token type: %s\n", tok.Literal, tok.Type)
			tok.Line = line
			tok.Column = column
			return tok
		} else if isDigit(l.ch) {
			numToken := l.readNumber()
			numToken.Line = line
			numToken.Column = column
			return numToken
		} else {
			tok = newToken(ILLEGAL, l.ch)
		}
	}

	tok.Line = line
	tok.Column = column

	l.readChar()
	fmt.Printf("DEBUG: NextToken - token type: %s, literal: %s\n", tok.Type, tok.Literal)
	return tok
}

// readIdentifier reads in an identifier and advances the lexer position
func (l *Lexer) readIdentifier() string {
	position := l.position
	for isLetter(l.ch) || isDigit(l.ch) {
		l.readChar()
	}
	return l.input[position:l.position]
}

// readNumber reads in a number (integer or float) and advances the lexer position
func (l *Lexer) readNumber() Token {
	position := l.position
	isFloat := false

	// Read digits before decimal point
	for isDigit(l.ch) {
		l.readChar()
	}

	// Check for decimal point
	if l.ch == '.' && isDigit(l.peekChar()) {
		isFloat = true
		l.readChar() // consume the dot

		// Read digits after decimal point
		for isDigit(l.ch) {
			l.readChar()
		}
	}

	// Get the numeric string
	numStr := l.input[position:l.position]

	// Create token with appropriate type
	var tok Token
	if isFloat {
		tok = Token{Type: FLOAT, Literal: numStr, Line: l.line, Column: l.column - len(numStr)}
	} else {
		tok = Token{Type: INT, Literal: numStr, Line: l.line, Column: l.column - len(numStr)}
	}

	return tok
}

// readString reads a string literal
func (l *Lexer) readString() string {
	position := l.position + 1 // Skip the opening quote

	for {
		l.readChar()
		if l.ch == '"' || l.ch == 0 {
			break
		}
	}

	return l.input[position:l.position]
}

// skipWhitespace skips any whitespace characters
func (l *Lexer) skipWhitespace() {
	for l.ch == ' ' || l.ch == '\t' || l.ch == '\n' || l.ch == '\r' {
		l.readChar()
	}
}

// isLetter returns true if the character is a letter or underscore
func isLetter(ch byte) bool {
	return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z' || ch == '_'
}

// isDigit returns true if the character is a digit
func isDigit(ch byte) bool {
	return '0' <= ch && ch <= '9'
}

// lookupIdent checks if an identifier is a keyword
func lookupIdent(ident string) TokenType {
	fmt.Printf("DEBUG: lookupIdent checking: %s\n", ident)
	if tok, ok := keywords[ident]; ok {
		fmt.Printf("DEBUG: lookupIdent found keyword: %s -> %s\n", ident, tok)
		return tok
	}
	fmt.Printf("DEBUG: lookupIdent not a keyword: %s -> IDENT\n", ident)
	return IDENT
}

// newToken creates a new token
func newToken(tokenType TokenType, ch byte) Token {
	return Token{Type: tokenType, Literal: string(ch)}
}

// Error returns a formatted lexer error with position information
func (l *Lexer) Error(message string) string {
	return fmt.Sprintf("Lexer error at line %d, column %d: %s", l.line, l.column, message)
}