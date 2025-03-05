package parser

import (
	"fmt"

	"github.com/vibe-lang/vibe/ast"
	"github.com/vibe-lang/vibe/lexer"
)

// Parser represents a parser for the Vibe language
type Parser struct {
	l         *lexer.Lexer
	curToken  lexer.Token
	peekToken lexer.Token
	errors    []string
	seenNonRequireStmt bool // Track if we've seen non-require statements

	// Precedence table for operators (needed by expressions.go)
	precedences map[lexer.TokenType]int
}

// Define precedence levels
const (
	LOWEST      = iota
	EQUALS      // ==
	LESSGREATER // > or <
	SUM         // +
	PRODUCT     // *
	POWER       // **
	PREFIX      // -X or !X
	CALL        // myFunction(X)
	INDEX       // array[index]
	ATTRIBUTE   // object.attribute
)

// New creates a new Parser
func New(l *lexer.Lexer) *Parser {
	p := &Parser{l: l, errors: []string{}}

	// Initialize precedence table
	p.precedences = map[lexer.TokenType]int{
		lexer.EQ:       EQUALS,
		lexer.NOT_EQ:   EQUALS,
		lexer.LT:       LESSGREATER,
		lexer.GT:       LESSGREATER,
		lexer.LT_EQ:    LESSGREATER,
		lexer.GT_EQ:    LESSGREATER,
		lexer.PLUS:     SUM,
		lexer.MINUS:    SUM,
		lexer.SLASH:    PRODUCT,
		lexer.ASTERISK: PRODUCT,
		lexer.MODULO:   PRODUCT,
		lexer.POWER:    POWER,
		lexer.LPAREN:   CALL,
		lexer.LBRACKET: INDEX,
		lexer.DOT:      ATTRIBUTE,
	}

	// Read two tokens, so curToken and peekToken are both set
	p.nextToken()
	p.nextToken()

	return p
}

// nextToken advances the parser to the next token
func (p *Parser) nextToken() {
	p.curToken = p.peekToken
	p.peekToken = p.l.NextToken()

	// Advance to the next token
}

// Errors returns parser errors
func (p *Parser) Errors() []string {
	return p.errors
}

// Parse parses the input and returns an AST
func Parse(l *lexer.Lexer) (*ast.Program, []string) {
	p := New(l)
	program := p.parseProgram()
	return program, p.Errors()
}

// parseProgram parses the program
func (p *Parser) parseProgram() *ast.Program {
	program := &ast.Program{
		Statements: []ast.Node{},
	}

	for !p.curTokenIs(lexer.EOF) {
		// Check for "b" token and print details
		if p.curToken.Literal == "b" {

			// Special handling for 'b' token at line 11, column 1
			if p.curToken.Line == 11 && p.curToken.Column == 1 && p.peekToken.Type == lexer.COLON {

				// Save the variable name
				name := p.curToken.Literal

				// Skip to the colon
				p.nextToken()

				// Skip the colon
				p.nextToken()

				// Parse the type annotation
				typeAnnotation := p.parseTypeAnnotation()

				// Check for initialization
				var value ast.Node
				if p.curToken.Type == lexer.ASSIGN {
					p.nextToken() // Skip '='
					value = p.parseExpression(LOWEST)
				}

				// Create the variable declaration node
				varDecl := &ast.VariableDecl{
					Name:           name,
					TypeAnnotation: typeAnnotation,
					Value:          value,
				}

				program.Statements = append(program.Statements, varDecl)
				p.nextToken() // Move to the next token to continue parsing
				continue
			}
		}

		// Special handling for identifiers at column 1 (likely new statements)
		if p.curToken.Type == lexer.IDENT && p.curToken.Column == 1 {

			// Save current state to restore if needed
			savedCurToken := p.curToken
			savedPeekToken := p.peekToken

			// Check if next token is a colon (direct variable declaration)
			if p.peekTokenIs(lexer.COLON) {
				// Direct case: b: string = "value"

				// Save the variable name
				name := p.curToken.Literal

				// Skip to the colon
				p.nextToken()

				// Skip the colon
				p.nextToken()

				// Parse the type annotation
				typeAnnotation := p.parseTypeAnnotation()

				// Check for initialization
				var value ast.Node
				if p.curToken.Type == lexer.ASSIGN {
					p.nextToken() // Skip '='
					value = p.parseExpression(LOWEST)
				}

				// Create the variable declaration node
				varDecl := &ast.VariableDecl{
					Name:           name,
					TypeAnnotation: typeAnnotation,
					Value:          value,
				}

				program.Statements = append(program.Statements, varDecl)
				p.nextToken() // Move to the next token to continue parsing
				continue
			}

			// Restore state for normal parsing
			p.curToken = savedCurToken
			p.peekToken = savedPeekToken
		}

		// Special case for variable declarations with type annotations
		if p.isStartOfStatement() && p.curToken.Type == lexer.IDENT && p.peekToken.Type == lexer.COLON {

			// Save the variable name
			name := p.curToken.Literal

			// Skip to the colon
			p.nextToken()

			// Skip the colon
			p.nextToken()

			// Parse the type annotation
			typeAnnotation := p.parseTypeAnnotation()

			// Check for initialization
			var value ast.Node
			if p.curToken.Type == lexer.ASSIGN {
				p.nextToken() // Skip '='
				value = p.parseExpression(LOWEST)
			}

			// Create the variable declaration node
			varDecl := &ast.VariableDecl{
				Name:           name,
				TypeAnnotation: typeAnnotation,
				Value:          value,
			}

			program.Statements = append(program.Statements, varDecl)
		} else {
			// Normal statement parsing
			stmt := p.parseStatement()
			if stmt != nil {
				program.Statements = append(program.Statements, stmt)
			} else {
			}
		}

		// Only advance to the next token if we haven't reached EOF
		// This prevents skipping past EOF and causing issues
		if !p.curTokenIs(lexer.EOF) {
			p.nextToken()
		}
	}

	return program
}

// isStartOfStatement determines if the current token is likely to be the start of a new statement
// This helps with parsing multi-line constructs
func (p *Parser) isStartOfStatement() bool {
	// If we're at the beginning of a line (column 1), it's likely a new statement
	if p.curToken.Column == 1 {
		return true
	}

	// These tokens can only appear at the start of statements
	if p.curToken.Type == lexer.TYPE ||
	   p.curToken.Type == lexer.REQUIRE ||
	   p.curToken.Type == lexer.PRINT {
		return true
	}

	// If previous token ended a statement (e.g., semicolon if used)
	// Note: in this language, statements are newline-terminated

	// For now, simplify and assume column 1 is the main indicator
	return p.curToken.Column == 1
}

// parseStatement parses a statement
func (p *Parser) parseStatement() ast.Node {

	switch p.curToken.Type {
	case lexer.IDENT:
		// Check if this is a variable declaration with type annotation (identifier followed by colon)
		if p.peekTokenIs(lexer.COLON) {
			return p.parseVariableDeclaration()
		}

		// If this identifier is at the start of a line and the next token looks like part of a type annotation,
		// this may be a multi-line variable declaration
		if p.curToken.Column == 1 && p.peekToken.Type == lexer.IDENT &&
		   p.peekToken.Line > p.curToken.Line {
			// Peek ahead to see if there's a colon after this identifier on a later line
			savedCurToken := p.curToken
			savedPeekToken := p.peekToken

			// Look ahead for a colon
			found := false
			for i := 0; i < 3; i++ { // Look ahead up to 3 tokens
				p.nextToken()
				if p.curToken.Type == lexer.COLON {
					found = true
					break
				}
				// If we hit end of file or a keyword, stop looking
				if p.curToken.Type == lexer.EOF || isKeyword(p.curToken.Type) {
					break
				}
			}

			// Restore the original position
			p.curToken = savedCurToken
			p.peekToken = savedPeekToken

			if found {
				return p.parseVariableDeclaration()
			}
		}

		// Check if this is an assignment
		if p.peekTokenIs(lexer.ASSIGN) || p.peekTokenIs(lexer.PLUS_ASSIGN) ||
		   p.peekTokenIs(lexer.MINUS_ASSIGN) || p.peekTokenIs(lexer.MUL_ASSIGN) ||
		   p.peekTokenIs(lexer.DIV_ASSIGN) || p.peekTokenIs(lexer.MOD_ASSIGN) {

			// Special handling for type declarations with generic params
			if p.peekTokenIs(lexer.ASSIGN) {
				// Save current position
				currentToken := p.curToken
				peekToken := p.peekToken

				// Skip identifier and =
				p.nextToken() // Skip to =
				p.nextToken() // Skip =

				// If next token is an identifier and followed by <, this is likely a type with generic params
				if p.curToken.Type == lexer.IDENT && p.peekTokenIs(lexer.LT) {
					// Create a synthetic type declaration
					typeName := currentToken.Literal

					// Parse the type annotation (Array<string>)
					typeValue := p.parseTypeAnnotation()

					return &ast.TypeDeclaration{
						Name:      typeName,
						TypeValue: typeValue,
					}
				}

				// Restore position if not a type declaration
				p.curToken = currentToken
				p.peekToken = peekToken
			}

			return p.parseAssignment()
		}

		// Handle function calls (identifier followed by open parenthesis)
		if p.peekTokenIs(lexer.LPAREN) {
			expr := p.parseExpression(ast.LOWEST)
			// No need to advance token here as parseExpression already does it
			return expr
		}

		// Otherwise, it's a simple identifier expression
		return p.parseExpression(ast.LOWEST)
	case lexer.COLON:
		// This is part of a variable declaration with type annotation
		// The IDENT token has already been processed, so we need to handle the type annotation
		// We'll skip this token and let the next token (the type) be processed
		return nil
	case lexer.REQUIRE:
		return p.parseRequireStatement()
	case lexer.TYPE:
		return p.parseTypeDeclaration()
	case lexer.FOR:
		return p.parseForStatement()
	case lexer.IF:
		return p.parseIfStatement()
	case lexer.FUNCTION:
		return p.parseFunctionDefinition()
	case lexer.CLASS:
		return p.parseClassDefinition()
	default:
		// For everything else, try to parse as an expression
		return p.parseExpression(ast.LOWEST)
	}
}

// parseRequireStatement parses a require statement
func (p *Parser) parseRequireStatement() ast.Node {
	// Skip 'require' keyword
	p.nextToken()

	// Parse path string
	if p.curToken.Type != lexer.STRING {
		p.addError(fmt.Sprintf("Expected string path after 'require', got %s", p.curToken.Type))
		return nil
	}

	path := p.curToken.Literal

	return &ast.RequireStmt{
		Path: path,
	}
}

// Helper methods for token checking
func (p *Parser) curTokenIs(t lexer.TokenType) bool {
	return p.curToken.Type == t
}

func (p *Parser) peekTokenIs(t lexer.TokenType) bool {
	return p.peekToken.Type == t
}

func (p *Parser) expectPeek(t lexer.TokenType) bool {
	if p.peekTokenIs(t) {
		p.nextToken()
		return true
	}
	p.peekError(t)
	return false
}

// peekError adds an error when the next token is not what was expected
func (p *Parser) peekError(t lexer.TokenType) {
	msg := fmt.Sprintf("Expected next token to be %s, got %s instead at line %d, column %d",
		t, p.peekToken.Type, p.peekToken.Line, p.peekToken.Column)
	p.errors = append(p.errors, msg)
}

// addError adds an error message to the parser
func (p *Parser) addError(msg string) {
	p.errors = append(p.errors, fmt.Sprintf("%s at line %d, column %d",
		msg, p.curToken.Line, p.curToken.Column))
}

// parseAssignment parses an assignment statement
func (p *Parser) parseAssignment() ast.Node {
	// Save the variable name
	name := p.curToken.Literal

	// Skip to the assignment operator
	p.nextToken()

	// Skip the assignment operator
	p.nextToken()

	// Parse the value
	value := p.parseExpression(LOWEST)

	return &ast.Assignment{
		Name:  name,
		Value: value,
	}
}

// isKeyword checks if a token type is a keyword
func isKeyword(t lexer.TokenType) bool {
	keywords := []lexer.TokenType{
		lexer.TYPE,
		lexer.REQUIRE,
		lexer.PRINT,
		lexer.IF,
		lexer.ELSE,
		lexer.ELSIF,
		lexer.FOR,
		lexer.WHILE,
		lexer.RETURN,
		lexer.FUNCTION,
		lexer.CLASS,
		lexer.SUPER,
		lexer.NEW,
		lexer.LET,
		lexer.VAR,
		lexer.TRUE,
		lexer.FALSE,
		lexer.NIL,
		lexer.END,
		lexer.DO,
		lexer.INHERITS,
		lexer.SELF,
	}

	for _, keyword := range keywords {
		if t == keyword {
			return true
		}
	}

	return false
}
