package lexer

import (
	"fmt"
)

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

// readChar reads the next character and advances our position in the input
func (l *Lexer) readChar() {
	if l.readPosition >= len(l.input) {
		l.ch = 0
	} else {
		l.ch = l.input[l.readPosition]
	}
	l.position = l.readPosition
	l.readPosition++

	// Update line and column information
	if l.ch == '\n' {
		l.line++
		l.column = 0
	} else {
		l.column++
	}
}

// peekChar looks at the next character without advancing our position
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
			tok = Token{Type: EQ, Literal: string(ch) + string(l.ch), Line: line, Column: column}
		} else {
			tok = NewToken(ASSIGN, l.ch, line, column)
		}
	case '+':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			tok = Token{Type: PLUS_ASSIGN, Literal: string(ch) + string(l.ch), Line: line, Column: column}
		} else {
			tok = NewToken(PLUS, l.ch, line, column)
		}
	case '-':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			tok = Token{Type: MINUS_ASSIGN, Literal: string(ch) + string(l.ch), Line: line, Column: column}
		} else {
			tok = NewToken(MINUS, l.ch, line, column)
		}
	case '!':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			tok = Token{Type: NOT_EQ, Literal: string(ch) + string(l.ch), Line: line, Column: column}
		} else {
			tok = NewToken(BANG, l.ch, line, column)
		}
	case '*':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			tok = Token{Type: MUL_ASSIGN, Literal: string(ch) + string(l.ch), Line: line, Column: column}
		} else if l.peekChar() == '*' {
			ch := l.ch
			l.readChar()
			tok = Token{Type: POWER, Literal: string(ch) + string(l.ch), Line: line, Column: column}
		} else {
			tok = NewToken(ASTERISK, l.ch, line, column)
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
			tok = Token{Type: DIV_ASSIGN, Literal: string(ch) + string(l.ch), Line: line, Column: column}
		} else {
			tok = NewToken(SLASH, l.ch, line, column)
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
			tok = Token{Type: MOD_ASSIGN, Literal: string(ch) + string(l.ch), Line: line, Column: column}
		} else {
			tok = NewToken(MODULO, l.ch, line, column)
		}
	case '<':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			tok = Token{Type: LT_EQ, Literal: string(ch) + string(l.ch), Line: line, Column: column}
		} else {
			tok = NewToken(LT, l.ch, line, column)
		}
	case '>':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			tok = Token{Type: GT_EQ, Literal: string(ch) + string(l.ch), Line: line, Column: column}
		} else {
			tok = NewToken(GT, l.ch, line, column)
		}
	case '&':
		if l.peekChar() == '&' {
			ch := l.ch
			l.readChar()
			tok = Token{Type: AND, Literal: string(ch) + string(l.ch), Line: line, Column: column}
		} else {
			tok = NewToken(ILLEGAL, l.ch, line, column)
		}
	case '|':
		if l.peekChar() == '|' {
			ch := l.ch
			l.readChar()
			tok = Token{Type: OR, Literal: string(ch) + string(l.ch), Line: line, Column: column}
		} else {
			tok = NewToken(ILLEGAL, l.ch, line, column)
		}
	case ';':
		tok = NewToken(SEMICOLON, l.ch, line, column)
	case ':':
		tok = NewToken(COLON, l.ch, line, column)
	case ',':
		tok = NewToken(COMMA, l.ch, line, column)
	case '.':
		tok = NewToken(DOT, l.ch, line, column)
	case '@':
		tok = NewToken(AT, l.ch, line, column)
	case '(':
		tok = NewToken(LPAREN, l.ch, line, column)
	case ')':
		tok = NewToken(RPAREN, l.ch, line, column)
	case '{':
		tok = NewToken(LBRACE, l.ch, line, column)
	case '}':
		tok = NewToken(RBRACE, l.ch, line, column)
	case '[':
		tok = NewToken(LBRACKET, l.ch, line, column)
	case ']':
		tok = NewToken(RBRACKET, l.ch, line, column)
	case '"':
		tok.Type = STRING
		tok.Literal = l.readString()
		tok.Line = line
		tok.Column = column
	case 0:
		tok.Literal = ""
		tok.Type = EOF
		tok.Line = line
		tok.Column = column
	default:
		if isLetter(l.ch) {
			tok.Literal = l.readIdentifier()
			tok.Type = LookupIdent(tok.Literal)
			tok.Line = line
			tok.Column = column
			return tok
		} else if isDigit(l.ch) {
			return l.readNumber()
		} else {
			tok = NewToken(ILLEGAL, l.ch, line, column)
		}
	}

	l.readChar()
	return tok
}

// readIdentifier reads an identifier and advances our position until it
// encounters a non-letter character
func (l *Lexer) readIdentifier() string {
	position := l.position
	for isLetter(l.ch) || isDigit(l.ch) {
		l.readChar()
	}
	return l.input[position:l.position]
}

// readNumber reads a number (integer or float) and advances our position until
// it encounters a non-number character
func (l *Lexer) readNumber() Token {
	position := l.position
	line := l.line
	column := l.column - 1 // Adjust for the first digit

	// Read the integer part
	for isDigit(l.ch) {
		l.readChar()
	}

	// Check if we have a decimal point followed by digits
	if l.ch == '.' && isDigit(l.peekChar()) {
		l.readChar() // consume the '.'
		for isDigit(l.ch) {
			l.readChar()
		}

		// Return a float token
		return Token{
			Type:    FLOAT,
			Literal: l.input[position:l.position],
			Line:    line,
			Column:  column,
		}
	}

	// Return an integer token
	return Token{
		Type:    INT,
		Literal: l.input[position:l.position],
		Line:    line,
		Column:  column,
	}
}

// readString reads a string and advances our position until it
// encounters the closing quote or EOF
func (l *Lexer) readString() string {
	l.readChar() // Skip the opening quote
	var result string

	for l.ch != '"' && l.ch != 0 {
		// Handle escape sequences
		if l.ch == '\\' {
			l.readChar() // Skip the backslash
			switch l.ch {
			case '"':
				result += string('"')
			case 'n':
				result += string('\n')
			case 't':
				result += string('\t')
			case 'r':
				result += string('\r')
			case '\\':
				result += string('\\')
			default:
				// Just add the character after the backslash
				result += string(l.ch)
			}
		} else {
			result += string(l.ch)
		}
		l.readChar()
	}

	// If we reached EOF without closing the string
	if l.ch == 0 {
		l.addError("Unterminated string literal")
	}

	return result
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

// Error returns an error message with line and column information
func (l *Lexer) Error(message string) string {
	return fmt.Sprintf("Error at line %d, column %d: %s", l.line, l.column, message)
}

// addError adds an error message for the current token
func (l *Lexer) addError(message string) {
	fmt.Printf("%s\n", l.Error(message))
}