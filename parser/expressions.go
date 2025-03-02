package parser

import (
	"fmt"
	"strconv"

	"github.com/example/vibe/ast"
	"github.com/example/vibe/lexer"
)

// parseExpression parses an expression with the given precedence
func (p *Parser) parseExpression(precedence int) ast.Node {
	fmt.Printf("DEBUG parseExpression: current token: %s (%s) at line %d, column %d, peek token: %s (%s) at line %d, column %d\n",
		p.curToken.Type, p.curToken.Literal, p.curToken.Line, p.curToken.Column,
		p.peekToken.Type, p.peekToken.Literal, p.peekToken.Line, p.peekToken.Column)

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
		fmt.Println("DEBUG parseExpression: found LBRACKET, calling parseArrayLiteral")
		leftExp = p.parseArrayLiteral()
		// No need to advance token here as parseArrayLiteral handles it
		fmt.Printf("DEBUG parseExpression (after parseArrayLiteral): current token: %s (%s) at line %d, column %d, peek token: %s (%s) at line %d, column %d\n",
			p.curToken.Type, p.curToken.Literal, p.curToken.Line, p.curToken.Column,
			p.peekToken.Type, p.peekToken.Literal, p.peekToken.Line, p.peekToken.Column)
		// Continue to the infix expression handling below
		// Don't return here, let the outer loop handle any infix operators that follow
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
	p.nextToken()

	// Special case for function calls
	if p.curTokenIs(lexer.LPAREN) {
		fmt.Println("DEBUG parseExpression: found function call")
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
			p.nextToken() // Move past the closing bracket
		case lexer.DOT:
			leftExp = p.parseDotExpression(leftExp)
		default:
			return leftExp
		}
	}

	fmt.Printf("DEBUG parseExpression (end): current token: %s (%s) at line %d, column %d, peek token: %s (%s) at line %d, column %d\n",
		p.curToken.Type, p.curToken.Literal, p.curToken.Line, p.curToken.Column,
		p.peekToken.Type, p.peekToken.Literal, p.peekToken.Line, p.peekToken.Column)

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
	fmt.Printf("DEBUG parseCallExpression: current token: %s (%s) at line %d, column %d, peek token: %s (%s) at line %d, column %d\n",
		p.curToken.Type, p.curToken.Literal, p.curToken.Line, p.curToken.Column,
		p.peekToken.Type, p.peekToken.Literal, p.peekToken.Line, p.peekToken.Column)

	callExpr := &ast.CallExpr{
		Function: function,
		Args:     []ast.Node{},
	}

	// Skip the opening parenthesis
	p.nextToken()

	// Handle empty argument list
	if p.curTokenIs(lexer.RPAREN) {
		fmt.Println("DEBUG parseCallExpression: empty argument list")
		p.nextToken() // Skip closing parenthesis

		// Check if the next token is an infix operator
		if isInfixOperator(p.curToken.Type) {
			// If it is, parse the infix expression
			return p.parseBinaryExpression(callExpr)
		}

		return callExpr
	}

	// Parse first argument
	fmt.Println("DEBUG parseCallExpression: parsing first argument")
	arg := p.parseExpression(ast.LOWEST)
	callExpr.Args = append(callExpr.Args, arg)

	// Parse additional arguments
	for p.peekTokenIs(lexer.COMMA) {
		p.nextToken() // Skip comma
		p.nextToken() // Move to next argument

		fmt.Println("DEBUG parseCallExpression: parsing additional argument")
		arg := p.parseExpression(ast.LOWEST)
		callExpr.Args = append(callExpr.Args, arg)
	}

	// Expect closing parenthesis
	if !p.curTokenIs(lexer.RPAREN) {
		if p.peekTokenIs(lexer.RPAREN) {
			p.nextToken() // Move to the closing parenthesis
		} else {
			fmt.Printf("DEBUG parseCallExpression: expected closing parenthesis, got %s\n", p.peekToken.Type)
			p.addError(fmt.Sprintf("Expected closing parenthesis, got %s at line %d, column %d",
				p.peekToken.Type, p.peekToken.Line, p.peekToken.Column))
			return nil
		}
	}

	// Check if the next token after the closing parenthesis is an infix operator
	if isInfixOperator(p.peekToken.Type) {
		// Skip past the closing parenthesis
		p.nextToken()

		// Parse the infix expression
		return p.parseBinaryExpression(callExpr)
	}

	// Skip past the closing parenthesis
	p.nextToken()

	fmt.Printf("DEBUG parseCallExpression: after closing parenthesis: current token: %s (%s) at line %d, column %d, peek token: %s (%s) at line %d, column %d\n",
		p.curToken.Type, p.curToken.Literal, p.curToken.Line, p.curToken.Column,
		p.peekToken.Type, p.peekToken.Literal, p.peekToken.Line, p.peekToken.Column)

	// Check if the current token is an infix operator
	if isInfixOperator(p.curToken.Type) {
		// If it is, parse the infix expression
		return p.parseBinaryExpression(callExpr)
	}

	return callExpr
}

// parseIndexExpression parses an array index expression
func (p *Parser) parseIndexExpression(array ast.Node) ast.Node {
	indexExpr := &ast.IndexExpr{
		Array: array,
	}

	p.nextToken() // Skip '['

	indexExpr.Index = p.parseExpression(ast.LOWEST)

	if !p.expectPeek(lexer.RBRACKET) {
		return nil
	}

	return indexExpr
}

// parseDotExpression parses a dot notation expression
func (p *Parser) parseDotExpression(object ast.Node) ast.Node {
	p.nextToken() // Skip '.'

	if p.curToken.Type != lexer.IDENT {
		p.addError(fmt.Sprintf("Expected property name after '.', got %s", p.curToken.Type))
		return nil
	}

	// Check if this is a method call
	if p.peekTokenIs(lexer.LPAREN) {
		methodCall := &ast.MethodCall{
			Object: object,
			Method: p.curToken.Literal,
			Args:   []ast.Node{},
		}

		p.nextToken() // Skip to '('
		p.nextToken() // Skip '('

		// Parse arguments if any
		if p.curToken.Type != lexer.RPAREN {
			// Parse first argument
			arg := p.parseExpression(ast.LOWEST)
			methodCall.Args = append(methodCall.Args, arg)

			// Parse additional arguments
			for p.peekTokenIs(lexer.COMMA) {
				p.nextToken() // Skip comma
				p.nextToken() // Move to next argument

				arg := p.parseExpression(ast.LOWEST)
				methodCall.Args = append(methodCall.Args, arg)
			}
		}

		if !p.expectPeek(lexer.RPAREN) {
			return nil
		}

		return methodCall
	}

	// Simple property access
	dotExpr := &ast.DotExpr{
		Object:   object,
		Property: p.curToken.Literal,
	}

	return dotExpr
}

// parseArrayLiteral parses an array literal
func (p *Parser) parseArrayLiteral() ast.Node {
	fmt.Printf("DEBUG parseArrayLiteral: current token: %s (%s) at line %d, column %d, peek token: %s (%s) at line %d, column %d\n",
		p.curToken.Type, p.curToken.Literal, p.curToken.Line, p.curToken.Column,
		p.peekToken.Type, p.peekToken.Literal, p.peekToken.Line, p.peekToken.Column)

	// Create the array literal - cur token is [
	arrayLit := &ast.ArrayLiteral{
		Elements: []ast.Node{},
	}

	// Skip the opening bracket [
	p.nextToken()
	fmt.Printf("DEBUG parseArrayLiteral (after skip [): current token: %s (%s) at line %d, column %d, peek token: %s (%s) at line %d, column %d\n",
		p.curToken.Type, p.curToken.Literal, p.curToken.Line, p.curToken.Column,
		p.peekToken.Type, p.peekToken.Literal, p.peekToken.Line, p.peekToken.Column)

	// Handle empty array case: []
	if p.curTokenIs(lexer.RBRACKET) {
		fmt.Println("DEBUG parseArrayLiteral: empty array case")
		p.nextToken() // Move past the closing bracket
		fmt.Printf("DEBUG parseArrayLiteral (after empty array): current token: %s (%s) at line %d, column %d, peek token: %s (%s) at line %d, column %d\n",
			p.curToken.Type, p.curToken.Literal, p.curToken.Line, p.curToken.Column,
			p.peekToken.Type, p.peekToken.Literal, p.peekToken.Line, p.peekToken.Column)
		return arrayLit
	}

	// Parse the first element
	fmt.Println("DEBUG parseArrayLiteral: parsing first element")
	firstElement := p.parseExpression(ast.LOWEST)
	if firstElement == nil {
		fmt.Println("DEBUG parseArrayLiteral: first element is nil")
		return nil
	}
	arrayLit.Elements = append(arrayLit.Elements, firstElement)
	fmt.Printf("DEBUG parseArrayLiteral (after first element): current token: %s (%s) at line %d, column %d, peek token: %s (%s) at line %d, column %d\n",
		p.curToken.Type, p.curToken.Literal, p.curToken.Line, p.curToken.Column,
		p.peekToken.Type, p.peekToken.Literal, p.peekToken.Line, p.peekToken.Column)

	// Parse additional elements (if any)
	for p.curTokenIs(lexer.COMMA) {
		fmt.Println("DEBUG parseArrayLiteral: found comma, parsing additional element")
		p.nextToken() // Move past the comma
		fmt.Printf("DEBUG parseArrayLiteral (after comma): current token: %s (%s) at line %d, column %d, peek token: %s (%s) at line %d, column %d\n",
			p.curToken.Type, p.curToken.Literal, p.curToken.Line, p.curToken.Column,
			p.peekToken.Type, p.peekToken.Literal, p.peekToken.Line, p.peekToken.Column)

		// Allow trailing comma [1, 2, ]
		if p.curTokenIs(lexer.RBRACKET) {
			fmt.Println("DEBUG parseArrayLiteral: found trailing comma")
			p.nextToken() // Move past the closing bracket
			fmt.Printf("DEBUG parseArrayLiteral (after trailing comma): current token: %s (%s) at line %d, column %d, peek token: %s (%s) at line %d, column %d\n",
				p.curToken.Type, p.curToken.Literal, p.curToken.Line, p.curToken.Column,
				p.peekToken.Type, p.peekToken.Literal, p.peekToken.Line, p.peekToken.Column)
			return arrayLit
		}

		element := p.parseExpression(ast.LOWEST)
		if element == nil {
			fmt.Println("DEBUG parseArrayLiteral: additional element is nil")
			return nil
		}
		arrayLit.Elements = append(arrayLit.Elements, element)
		fmt.Printf("DEBUG parseArrayLiteral (after additional element): current token: %s (%s) at line %d, column %d, peek token: %s (%s) at line %d, column %d\n",
			p.curToken.Type, p.curToken.Literal, p.curToken.Line, p.curToken.Column,
			p.peekToken.Type, p.peekToken.Literal, p.peekToken.Line, p.peekToken.Column)
	}

	// Check if we're already at the closing bracket
	if p.curTokenIs(lexer.RBRACKET) {
		fmt.Println("DEBUG parseArrayLiteral: at closing bracket")
		p.nextToken() // Move past the closing bracket
		fmt.Printf("DEBUG parseArrayLiteral (after closing bracket): current token: %s (%s) at line %d, column %d, peek token: %s (%s) at line %d, column %d\n",
			p.curToken.Type, p.curToken.Literal, p.curToken.Line, p.curToken.Column,
			p.peekToken.Type, p.peekToken.Literal, p.peekToken.Line, p.peekToken.Column)
		return arrayLit
	}

	// Otherwise, expect the next token to be the closing bracket
	fmt.Println("DEBUG parseArrayLiteral: expecting closing bracket")
	if !p.expectPeek(lexer.RBRACKET) {
		p.addError(fmt.Sprintf("Expected ']' after array elements, got %s at line %d, column %d",
			p.peekToken.Type, p.peekToken.Line, p.peekToken.Column))
		return nil
	}

	// We've already advanced the token in expectPeek
	fmt.Printf("DEBUG parseArrayLiteral (after expectPeek): current token: %s (%s) at line %d, column %d, peek token: %s (%s) at line %d, column %d\n",
		p.curToken.Type, p.curToken.Literal, p.curToken.Line, p.curToken.Column,
		p.peekToken.Type, p.peekToken.Literal, p.peekToken.Line, p.peekToken.Column)
	return arrayLit
}