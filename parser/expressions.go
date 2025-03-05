package parser

import (
	"fmt"
	"strconv"

	"github.com/vibe-lang/vibe/ast"
	"github.com/vibe-lang/vibe/lexer"
)

// parseExpression parses an expression with the given precedence
func (p *Parser) parseExpression(precedence int) ast.Node {
	// Handle prefix expressions
	var leftExp ast.Node

	switch p.curToken.Type {
	case lexer.IDENT:
		leftExp = &ast.Identifier{Name: p.curToken.Literal}
	case lexer.INT:
		value, err := strconv.ParseFloat(p.curToken.Literal, 64)
		if err != nil {
			p.addError(fmt.Sprintf("Could not parse %q as integer", p.curToken.Literal))
			return nil
		}
		leftExp = &ast.NumberLiteral{Value: value, IsInt: true}
	case lexer.FLOAT:
		value, err := strconv.ParseFloat(p.curToken.Literal, 64)
		if err != nil {
			p.addError(fmt.Sprintf("Could not parse %q as float", p.curToken.Literal))
			return nil
		}
		leftExp = &ast.NumberLiteral{Value: value, IsInt: false}
	case lexer.STRING:
		leftExp = &ast.StringLiteral{Value: p.curToken.Literal}
	case lexer.TRUE:
		leftExp = &ast.BooleanLiteral{Value: true}
	case lexer.FALSE:
		leftExp = &ast.BooleanLiteral{Value: false}
	case lexer.NIL:
		leftExp = &ast.NilLiteral{}
	case lexer.FUNCTION:
		// Handle anonymous function expressions by using the updated parseFunctionDefinition
		leftExp = p.parseFunctionDefinition()
		return leftExp // Return early as parseFunctionDefinition already advances tokens
	case lexer.LPAREN:
		p.nextToken() // Skip the opening parenthesis
		exp := p.parseExpression(ast.LOWEST)

		// Check for closing parenthesis but don't advance past it yet
		if !p.curTokenIs(lexer.RPAREN) {
			if !p.peekTokenIs(lexer.RPAREN) {
				p.addError(fmt.Sprintf("Expected next token to be %s, got %s instead at line %d, column %d",
					lexer.RPAREN, p.peekToken.Type, p.peekToken.Line, p.peekToken.Column))
				return nil
			}
			p.nextToken() // Move to the closing parenthesis
		}

		// Set leftExp to the parenthesized expression
		leftExp = exp

		// We're now at the closing parenthesis, but we don't advance past it yet
		// The nextToken() call after the switch will advance us past the closing parenthesis
		// Then the outer loop will handle any infix operators that follow
	case lexer.LBRACKET:
		leftExp = p.parseArrayLiteral()
		// parseArrayLiteral now advances past the closing bracket

		// Check if the current token is an opening bracket for an index expression
		if p.curTokenIs(lexer.LBRACKET) {
			leftExp = p.parseIndexExpression(leftExp)
		}

		return leftExp
	case lexer.MINUS, lexer.BANG:
		operator := p.curToken.Literal
		p.nextToken() // Consume the operator
		operand := p.parseExpression(ast.PREFIX)
		leftExp = &ast.UnaryExpr{Operator: operator, Right: operand}
	case lexer.SELF:
		leftExp = &ast.SelfExpr{}
	default:
		return nil
	}

	// Move to the next token after the prefix expression
	// Only advance if we're not already at the next token
	if !p.curTokenIs(lexer.LBRACKET) && !p.curTokenIs(lexer.LPAREN) && !p.curTokenIs(lexer.DOT) {
		p.nextToken()
	}

	// Special case for function calls
	if p.curTokenIs(lexer.LPAREN) {
		leftExp = p.parseCallExpression(leftExp)
		// parseCallExpression already advances the token past the closing parenthesis
	}

	// Now handle infix expressions
	for precedence < p.curPrecedence() && !p.curTokenIs(lexer.EOF) {
		if !isInfixOperator(p.curToken.Type) {
			break
		}

		switch p.curToken.Type {
		case lexer.PLUS, lexer.MINUS, lexer.ASTERISK, lexer.SLASH, lexer.MODULO, lexer.POWER,
			lexer.EQ, lexer.NOT_EQ, lexer.LT, lexer.GT, lexer.LT_EQ, lexer.GT_EQ,
			lexer.AND, lexer.OR:
			leftExp = p.parseBinaryExpression(leftExp)
		case lexer.LPAREN:
			leftExp = p.parseCallExpression(leftExp)
		case lexer.LBRACKET:
			leftExp = p.parseIndexExpression(leftExp)
		case lexer.DOT:
			leftExp = p.parseDotExpression(leftExp)
		default:
			return leftExp
		}
	}

	return leftExp
}

// Helper function to check if a token type is an infix operator
func isInfixOperator(tokenType lexer.TokenType) bool {
	switch tokenType {
	case lexer.PLUS, lexer.MINUS, lexer.ASTERISK, lexer.SLASH, lexer.MODULO, lexer.POWER,
		lexer.EQ, lexer.NOT_EQ, lexer.LT, lexer.GT, lexer.LT_EQ, lexer.GT_EQ,
		lexer.AND, lexer.OR, lexer.LPAREN, lexer.LBRACKET, lexer.DOT:
		return true
	default:
		return false
	}
}

// Get precedence for operators
func (p *Parser) peekPrecedence() int {
	switch p.peekToken.Type {
	case lexer.EQ, lexer.NOT_EQ:
		return ast.EQUALS
	case lexer.LT, lexer.GT, lexer.LT_EQ, lexer.GT_EQ:
		return ast.LESSGREATER
	case lexer.PLUS, lexer.MINUS:
		return ast.SUM
	case lexer.ASTERISK, lexer.SLASH, lexer.MODULO:
		return ast.PRODUCT
	case lexer.POWER:
		return ast.POWER
	case lexer.LPAREN:
		return ast.CALL
	case lexer.LBRACKET:
		return ast.INDEX
	case lexer.DOT:
		return ast.DOT
	default:
		return ast.LOWEST
	}
}

func (p *Parser) curPrecedence() int {
	switch p.curToken.Type {
	case lexer.EQ, lexer.NOT_EQ:
		return ast.EQUALS
	case lexer.LT, lexer.GT, lexer.LT_EQ, lexer.GT_EQ:
		return ast.LESSGREATER
	case lexer.PLUS, lexer.MINUS:
		return ast.SUM
	case lexer.ASTERISK, lexer.SLASH, lexer.MODULO:
		return ast.PRODUCT
	case lexer.POWER:
		return ast.POWER
	case lexer.LPAREN:
		return ast.CALL
	case lexer.LBRACKET:
		return ast.INDEX
	case lexer.DOT:
		return ast.DOT
	default:
		return ast.LOWEST
	}
}

// parseBinaryExpression parses a binary expression
func (p *Parser) parseBinaryExpression(left ast.Node) ast.Node {
	operator := p.curToken.Literal
	precedence := p.curPrecedence()

	p.nextToken()

	right := p.parseExpression(precedence)

	return &ast.BinaryExpr{
		Left:     left,
		Operator: operator,
		Right:    right,
	}
}

// parseCallExpression parses a function call
func (p *Parser) parseCallExpression(function ast.Node) ast.Node {
	callExpr := &ast.CallExpr{
		Function: function,
		Args:     []ast.Node{},
	}

	// Skip the opening parenthesis
	p.nextToken()

	// Handle empty argument list
	if p.curTokenIs(lexer.RPAREN) {
		p.nextToken() // Skip closing parenthesis
		return callExpr
	}

	// Parse first argument
	arg := p.parseExpression(ast.LOWEST)
	callExpr.Args = append(callExpr.Args, arg)

	// Parse additional arguments
	for p.peekTokenIs(lexer.COMMA) {
		p.nextToken() // Skip comma
		p.nextToken() // Move to next argument

		arg := p.parseExpression(ast.LOWEST)
		callExpr.Args = append(callExpr.Args, arg)
	}

	// Expect closing parenthesis
	if !p.expectPeek(lexer.RPAREN) {
		p.addError(fmt.Sprintf("Expected next token to be ), got %s instead at line %d, column %d",
			p.peekToken.Type, p.peekToken.Line, p.peekToken.Column))
		return nil
	}

	return callExpr
}

// parseIndexExpression parses an array index expression
func (p *Parser) parseIndexExpression(array ast.Node) ast.Node {
	indexExpr := &ast.IndexExpr{
		Array: array,
	}

	// Skip the opening bracket
	p.nextToken()

	// Parse the index expression
	indexExpr.Index = p.parseExpression(ast.LOWEST)

	// Check for closing bracket
	if !p.expectPeek(lexer.RBRACKET) {
		p.addError(fmt.Sprintf("Expected next token to be ], got %s instead at line %d, column %d",
			p.peekToken.Type, p.peekToken.Line, p.peekToken.Column))
		return nil
	}

	return indexExpr
}

// parseDotExpression parses a dot notation expression
func (p *Parser) parseDotExpression(object ast.Node) ast.Node {
	p.nextToken() // consume the dot and move to the identifier

	if p.curToken.Type != lexer.IDENT {
		p.addError(fmt.Sprintf("Expected property name after '.', got %s", p.curToken.Type))
		return nil
	}

	propertyOrMethodName := p.curToken.Literal

	// Check if it's a method call
	if p.peekTokenIs(lexer.LPAREN) {
		// It's a method call
		methodCall := &ast.MethodCall{
			Object: object,
			Method: propertyOrMethodName,
			Args:   []ast.Node{},
		}

		p.nextToken() // consume IDENT and move to LPAREN
		p.nextToken() // consume LPAREN

		// Handle empty argument list
		if p.curTokenIs(lexer.RPAREN) {
			p.nextToken() // consume RPAREN
			return methodCall
		}

		// Parse first argument
		arg := p.parseExpression(ast.LOWEST)
		if arg != nil {
			methodCall.Args = append(methodCall.Args, arg)
		}

		// Parse additional arguments
		for p.peekTokenIs(lexer.COMMA) {
			p.nextToken() // consume current token
			p.nextToken() // consume comma and move to the next argument
			arg := p.parseExpression(ast.LOWEST)
			if arg != nil {
				methodCall.Args = append(methodCall.Args, arg)
			}
		}

		// Ensure we have a closing parenthesis
		if !p.expectPeek(lexer.RPAREN) {
			p.addError(fmt.Sprintf("Expected closing parenthesis, got %s", p.peekToken.Type))
			return nil
		}

		return methodCall
	}

	// It's a property access
	dotExpr := &ast.DotExpr{
		Object:   object,
		Property: propertyOrMethodName,
	}

	p.nextToken() // Move past the property name
	return dotExpr
}

// parseArrayLiteral parses an array literal expression
func (p *Parser) parseArrayLiteral() ast.Node {
	arrayLit := &ast.ArrayLiteral{
		Elements: []ast.Node{},
	}

	// Check for empty array []
	if p.peekTokenIs(lexer.RBRACKET) {
		p.nextToken() // Move past the opening bracket to the closing bracket
		p.nextToken() // Advance past the closing bracket
		return arrayLit
	}

	// Parse the first element
	p.nextToken() // Move past the opening bracket
	firstElement := p.parseExpression(ast.LOWEST)
	if firstElement == nil {
		return nil
	}
	arrayLit.Elements = append(arrayLit.Elements, firstElement)

	// Parse additional elements (if any)
	for p.curTokenIs(lexer.COMMA) {
		p.nextToken() // Move past the comma

		// Allow trailing comma [1, 2, ]
		if p.curTokenIs(lexer.RBRACKET) {
			p.nextToken() // Advance past the closing bracket
			return arrayLit
		}

		element := p.parseExpression(ast.LOWEST)
		if element == nil {
			return nil
		}
		arrayLit.Elements = append(arrayLit.Elements, element)
	}

	// Check if we're already at the closing bracket
	if p.curTokenIs(lexer.RBRACKET) {
		p.nextToken() // Advance past the closing bracket
		return arrayLit
	}

	// Otherwise, expect the next token to be the closing bracket
	if !p.expectPeek(lexer.RBRACKET) {
		p.addError(fmt.Sprintf("Expected ']' after array elements, got %s at line %d, column %d",
			p.peekToken.Type, p.peekToken.Line, p.peekToken.Column))
		return nil
	}

	// expectPeek already advanced past the closing bracket
	return arrayLit
}