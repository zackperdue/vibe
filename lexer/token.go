package lexer

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
	POWER    = "**"   // Exponentiation operator

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
	REQUIRE  = "REQUIRE"
	TYPE     = "TYPE"  // Type declarations

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
	"puts":     PRINT,
	"end":      END,
	"do":       DO,
	"require":  REQUIRE,
	"type":     TYPE,
	"class":    CLASS,
	"inherits": INHERITS,
	"self":     SELF,
	"super":    SUPER,
	"new":      NEW,
}

// LookupIdent checks if a given identifier is a keyword
func LookupIdent(ident string) TokenType {
	if tok, ok := keywords[ident]; ok {
		return tok
	}
	return IDENT
}

// Helper function to create a new token
func NewToken(tokenType TokenType, ch byte, line, column int) Token {
	return Token{Type: tokenType, Literal: string(ch), Line: line, Column: column}
}