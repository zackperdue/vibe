package parser

import (
	"fmt"
	"strconv"

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

// Parse is a helper function that creates a new parser and parses the program
func Parse(l *lexer.Lexer) (*ast.Program, []string) {
	p := New(l)
	program := p.parseProgram()
	return program, p.Errors()
}

// parseProgram parses the entire program
func (p *Parser) parseProgram() *ast.Program {
	fmt.Printf("DEBUG: Starting to parse program\n")
	program := &ast.Program{
		Statements: []ast.Node{},
	}

	for !p.curTokenIs(lexer.EOF) {
		fmt.Printf("DEBUG: Parsing next statement, current token: '%s' (Type: %s, Line: %d, Col: %d)\n",
			p.curToken.Literal, p.curToken.Type, p.curToken.Line, p.curToken.Column)

		// Skip semicolons between statements
		if p.curTokenIs(lexer.SEMICOLON) {
			p.nextToken()
			continue
		}

		// Skip over closing parentheses, which might occur after a function call
		if p.curTokenIs(lexer.RPAREN) {
			fmt.Printf("DEBUG: Skipping over closing parenthesis\n")
			p.nextToken()
			continue
		}

		fmt.Printf("DEBUG: Parsing regular statement\n")
		stmt := p.parseStatement()
		if stmt != nil {
			// Add to the program's statements
			program.Statements = append(program.Statements, stmt)
			fmt.Printf("DEBUG: Added statement of type %T\n", stmt)
		} else {
			fmt.Printf("DEBUG: ⚠️ Statement was nil, skipping\n")
		}

		// Move to the next token to start parsing the next statement
		// This should skip any tokens that weren't consumed by the parser
		if p.curTokenIs(lexer.SEMICOLON) {
			p.nextToken() // Skip the semicolon before moving to the next statement
		}
	}

	fmt.Printf("DEBUG: Finished parsing program, found %d statements\n", len(program.Statements))
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
	case lexer.LET:
		return p.parseLetStatement()
	case lexer.RETURN:
		return p.parseReturnStatement()
	case lexer.IF:
		return p.parseIfStatement()
	case lexer.WHILE:
		return p.parseWhileStatement()
	case lexer.FOR:
		return p.parseForStatement()
	case lexer.FUNCTION:
		return p.parseFunctionDefinition()
	case lexer.CLASS:
		return p.parseClassDefinition()
	case lexer.REQUIRE:
		return p.parseRequireStatement()
	default:
		// Check for variable declaration with type annotation (x: int = 5)
		if p.curTokenIs(lexer.IDENT) && p.peekTokenIs(lexer.COLON) {
			return p.parseVariableDeclaration()
		}
		// Check for assignment (a = 1)
		if p.curTokenIs(lexer.IDENT) && p.peekTokenIs(lexer.ASSIGN) {
			return p.parseAssignment()
		}
		// Otherwise, it's an expression statement
		return p.parseExpressionStatement()
	}
}

// parseExpressionStatement parses an expression statement
func (p *Parser) parseExpressionStatement() ast.Node {
	stmt := &ast.ExpressionStatement{
		Expression: p.parseExpression(ast.LOWEST),
	}

	// Optional semicolon
	if p.peekTokenIs(lexer.SEMICOLON) {
		p.nextToken()
	}

	return stmt
}

// parseLetStatement parses a let statement
func (p *Parser) parseLetStatement() ast.Node {
	stmt := &ast.VariableDecl{
		Name: "",
	}

	// Skip 'let' keyword
	p.nextToken()

	// Parse variable name
	if !p.curTokenIs(lexer.IDENT) {
		p.addError(fmt.Sprintf("Expected identifier after 'let', got %s", p.curToken.Type))
		return nil
	}

	stmt.Name = p.curToken.Literal

	// Parse optional type annotation
	if p.peekTokenIs(lexer.COLON) {
		p.nextToken() // Skip colon
		p.nextToken() // Move to type name
		stmt.TypeAnnotation = p.parseTypeAnnotation()
	}

	// Expect equals sign
	if !p.expectPeek(lexer.ASSIGN) {
		return nil
	}

	// Skip equals sign
	p.nextToken()

	// Parse value expression
	stmt.Value = p.parseExpression(ast.LOWEST)

	// Optional semicolon
	if p.peekTokenIs(lexer.SEMICOLON) {
		p.nextToken()
	}

	return stmt
}

// parseReturnStatement parses a return statement
func (p *Parser) parseReturnStatement() ast.Node {
	stmt := &ast.ReturnStmt{}

	// Skip 'return' keyword
	p.nextToken()

	// Parse return value (if any)
	if !p.curTokenIs(lexer.SEMICOLON) {
		stmt.Value = p.parseExpression(ast.LOWEST)
	}

	// Optional semicolon
	if p.peekTokenIs(lexer.SEMICOLON) {
		p.nextToken()
	}

	return stmt
}

// parseAssignment parses an assignment statement
func (p *Parser) parseAssignment() ast.Node {
	fmt.Printf("DEBUG PARSER: parseAssignment starting with token: %+v\n", p.curToken)

	// Parse the left side (identifier)
	left := &ast.Identifier{Name: p.curToken.Literal}

	// Skip to the equals sign
	p.nextToken()

	// Skip the equals sign
	p.nextToken()

	fmt.Printf("DEBUG PARSER: After equals, current token is: %+v\n", p.curToken)

	// Parse the right side (value)
	var right ast.Node

	// Special case for integer literals
	if p.curTokenIs(lexer.INT) {
		fmt.Printf("DEBUG PARSER: Found integer literal: %s\n", p.curToken.Literal)
		value, err := strconv.ParseFloat(p.curToken.Literal, 64)
		if err != nil {
			p.addError(fmt.Sprintf("Could not parse %q as integer", p.curToken.Literal))
			return nil
		}
		right = &ast.NumberLiteral{Value: value, IsInt: true}
	} else {
		right = p.parseExpression(ast.LOWEST)
	}

	// Create the assignment node
	assignment := &ast.Assignment{
		Name:  left.Name,
		Value: right,
	}

	fmt.Printf("DEBUG PARSER: Parsed assignment %s = %v\n", left.Name, right)

	// Optional semicolon
	if p.peekTokenIs(lexer.SEMICOLON) {
		p.nextToken()
	}

	return assignment
}

// parseClassDefinition parses a class definition
func (p *Parser) parseClassDefinition() ast.Node {
	class := &ast.ClassDef{
		Name:    "",
		Methods: []ast.Node{},
	}

	// Skip 'class' keyword
	p.nextToken()

	// Parse class name
	if !p.curTokenIs(lexer.IDENT) {
		p.addError(fmt.Sprintf("Expected identifier after 'class', got %s", p.curToken.Type))
		return nil
	}

	class.Name = p.curToken.Literal

	// Expect opening brace
	if !p.expectPeek(lexer.LBRACE) {
		return nil
	}

	// Skip opening brace
	p.nextToken()

	// Parse methods
	for !p.curTokenIs(lexer.RBRACE) && !p.curTokenIs(lexer.EOF) {
		// Parse method
		if p.curTokenIs(lexer.FUNCTION) {
			method := p.parseFunctionDefinition()
			if method != nil {
				class.Methods = append(class.Methods, method)
			}
		} else {
			p.addError(fmt.Sprintf("Expected method definition, got %s", p.curToken.Type))
			return nil
		}
	}

	// Skip closing brace
	if !p.curTokenIs(lexer.RBRACE) {
		p.addError(fmt.Sprintf("Expected closing brace, got %s", p.curToken.Type))
		return nil
	}
	p.nextToken()

	return class
}

// parseRequireStatement parses a require statement
func (p *Parser) parseRequireStatement() ast.Node {
	stmt := &ast.RequireStmt{}

	// Skip 'require' keyword
	p.nextToken()

	// Parse module name
	if !p.curTokenIs(lexer.STRING) {
		p.addError(fmt.Sprintf("Expected string after 'require', got %s", p.curToken.Type))
		return nil
	}

	stmt.Path = p.curToken.Literal

	// Optional semicolon
	if p.peekTokenIs(lexer.SEMICOLON) {
		p.nextToken()
	}

	// Mark that we've seen a require statement
	p.seenNonRequireStmt = true

	return stmt
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

// Helper method to read an integer value from the current token
func (p *Parser) readInt() int64 {
	val, err := strconv.ParseInt(p.curToken.Literal, 0, 64)
	if err != nil {
		msg := fmt.Sprintf("could not parse %q as integer", p.curToken.Literal)
		p.errors = append(p.errors, msg)
		return 0
	}
	return val
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

// parseTypeAnnotation parses a type annotation
func (p *Parser) parseTypeAnnotation() *ast.TypeAnnotation {
	// A type annotation must start with an identifier (the type name)
	if !p.curTokenIs(lexer.IDENT) {
		p.addError(fmt.Sprintf("Expected type name, got %s", p.curToken.Type))
		return &ast.TypeAnnotation{TypeName: "any"} // Default to any type
	}

	// Create the type annotation with the current token as the type name
	typeAnnotation := &ast.TypeAnnotation{
		TypeName: p.curToken.Literal,
	}

	// Check if there's a generic type parameter list starting with '<'
	if p.peekTokenIs(lexer.LT) {
		p.nextToken() // Advance to '<'
		p.nextToken() // Advance to first type parameter

		// Initialize the type parameters slice
		var typeParams []ast.Node

		// Parse type parameters until we hit the closing '>'
		for !p.curTokenIs(lexer.GT) && !p.curTokenIs(lexer.EOF) {
			// Parse the type parameter, which could be a simple type or another generic type
			var param ast.Node

			if p.curTokenIs(lexer.IDENT) {
				// Check if this is a generic type (has a '<' after it)
				if p.peekTokenIs(lexer.LT) {
					// This is a nested generic type
					param = p.parseTypeAnnotation() // This will handle the nested generic
				} else {
					// Simple type
					param = &ast.TypeAnnotation{TypeName: p.curToken.Literal}
					p.nextToken() // Advance past the identifier
				}
			} else {
				p.addError(fmt.Sprintf("Expected type parameter name, got %s", p.curToken.Type))
				break
			}

			// Add the parameter to our list
			typeParams = append(typeParams, param)

			// If next token is a comma, skip it and continue
			if p.peekTokenIs(lexer.COMMA) {
				p.nextToken() // Advance to the comma
				p.nextToken() // Advance past the comma
			}
		}

		// Store the type parameters in the type annotation
		typeAnnotation.TypeParams = typeParams

		// Advance past the closing '>'
		if p.curTokenIs(lexer.GT) {
			p.nextToken()
		}
	} else {
		// If there are no type parameters, advance past the type name
		p.nextToken()
	}

	return typeAnnotation
}
