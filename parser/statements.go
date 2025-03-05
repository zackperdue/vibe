package parser

import (
	"fmt"

	"github.com/vibe-lang/vibe/ast"
	"github.com/vibe-lang/vibe/lexer"
)

// parseIfStatement parses an if statement
func (p *Parser) parseIfStatement() ast.Node {
	stmt := &ast.IfStmt{}

	p.nextToken()
	stmt.Condition = p.parseExpression(ast.LOWEST)

	if !p.curTokenIs(lexer.DO) {
		if !p.expectPeek(lexer.DO) {
			p.addError(fmt.Sprintf("Expected 'do' after if condition, got %s", p.peekToken.Type))
			return nil
		}
	}

	p.nextToken()
	stmt.Consequence = p.parseBlockStatements(lexer.END, lexer.ELSIF, lexer.ELSE)

	for p.curTokenIs(lexer.ELSIF) {
		p.nextToken()
		elsifCondition := p.parseExpression(ast.LOWEST)

		if !p.curTokenIs(lexer.DO) {
			if !p.expectPeek(lexer.DO) {
				p.addError(fmt.Sprintf("Expected 'do' after elsif condition, got %s", p.peekToken.Type))
				return nil
			}
		}

		p.nextToken()
		elsifConsequence := p.parseBlockStatements(lexer.END, lexer.ELSIF, lexer.ELSE)

		stmt.ElseIfBlocks = append(stmt.ElseIfBlocks, ast.ElseIfBlock{
			Condition:   elsifCondition,
			Consequence: elsifConsequence,
		})
	}

	if p.curTokenIs(lexer.ELSE) {
		p.nextToken()
		stmt.Alternative = p.parseBlockStatements(lexer.END)
	}

	if !p.curTokenIs(lexer.END) {
		p.addError(fmt.Sprintf("Expected 'end' to close if statement, got %s", p.curToken.Type))
		return nil
	}

	p.nextToken()
	return stmt
}

// parseBlockStatements parses a block of statements until one of the end tokens is reached
func (p *Parser) parseBlockStatements(endTokens ...lexer.TokenType) *ast.BlockStmt {
	block := &ast.BlockStmt{
		Statements: []ast.Node{},
	}

	// Continue parsing statements until we hit one of the end tokens or EOF
	for !p.curTokenIs(lexer.EOF) && !containsTokenType(p.curToken.Type, endTokens) {
		// Report error if we encounter a semicolon
		if p.curTokenIs(lexer.SEMICOLON) {
			p.addError(fmt.Sprintf("Unexpected semicolon at line %d, column %d. Vibe syntax does not allow semicolons.",
				p.curToken.Line, p.curToken.Column))
			p.nextToken() // Skip the semicolon to continue parsing
			continue
		}

		// Store current token position in case we need to recover
		startToken := p.curToken
		startPeekToken := p.peekToken

		var stmt ast.Node

		// Special case for assignments
		if p.curTokenIs(lexer.IDENT) && p.peekTokenIs(lexer.ASSIGN) {
			stmt = p.parseAssignment()
		} else if p.curTokenIs(lexer.IDENT) && (p.peekTokenIs(lexer.PLUS_ASSIGN) ||
			p.peekTokenIs(lexer.MINUS_ASSIGN) ||
			p.peekTokenIs(lexer.MUL_ASSIGN) ||
			p.peekTokenIs(lexer.DIV_ASSIGN) ||
			p.peekTokenIs(lexer.MOD_ASSIGN)) {
			// Handle compound assignments
			stmt = p.parseCompoundAssignment()
		} else if p.curTokenIs(lexer.AT) && p.peekTokenIs(lexer.IDENT) {
			// Special case for instance variable assignments
			stmt = p.parseInstanceVarAssignment()
		} else if p.curTokenIs(lexer.IDENT) && p.peekTokenIs(lexer.LPAREN) {
			// Special case for function calls
			expr := p.parseExpression(ast.LOWEST)
			if expr != nil {
				stmt = &ast.ExpressionStatement{Expression: expr}
			}
		} else {
			// Try to parse the statement
			stmt = p.parseStatement()
		}

		// If statement parsing failed, try to recover
		if stmt == nil {
			// If we're still at the same token, we need to manually advance to avoid an infinite loop
			if p.curToken == startToken && p.peekToken == startPeekToken {
				p.addError(fmt.Sprintf("Failed to parse statement at line %d, column %d. Skipping.",
					p.curToken.Line, p.curToken.Column))
				p.nextToken() // Force advance to the next token
				continue
			}
		} else {
			block.Statements = append(block.Statements, stmt)
		}
	}

	return block
}

// Helper function to check if a token type is in a list
func containsTokenType(tokenType lexer.TokenType, tokenTypes []lexer.TokenType) bool {
	for _, t := range tokenTypes {
		if tokenType == t {
			return true
		}
	}
	return false
}

// parseWhileStatement parses a while loop statement
func (p *Parser) parseWhileStatement() ast.Node {
	// Skip 'while' keyword
	p.nextToken()

	// Parse condition
	condition := p.parseExpression(ast.LOWEST)

	// Optional 'do' keyword
	if p.peekTokenIs(lexer.DO) {
		p.nextToken() // Skip to 'do'
		p.nextToken() // Skip 'do'
	}

	// Parse body
	body := p.parseBlockStatements(lexer.END)

	// Ensure 'end' token is consumed
	if p.curTokenIs(lexer.END) {
		p.nextToken() // Skip 'end'
	} else {
		p.addError(fmt.Sprintf("Expected 'end' to close while loop, got %s", p.curToken.Type))
	}

	return &ast.WhileStmt{
		Condition: condition,
		Body:      body,
	}
}

// parseForStatement parses a for loop statement
func (p *Parser) parseForStatement() ast.Node {
	// Current token is 'for'
	forStmt := &ast.ForStmt{}

	// Skip 'for' keyword
	p.nextToken()

	// Expect an identifier for iterator
	if !p.curTokenIs(lexer.IDENT) {
		p.addError(fmt.Sprintf("Expected identifier for iterator, got %s", p.curToken.Type))
		return nil
	}
	forStmt.Iterator = p.curToken.Literal

	// Move to 'in' token
	if !p.expectPeek(lexer.IN) {
		p.addError(fmt.Sprintf("Expected 'in' after iterator, got %s", p.peekToken.Type))
		return nil
	}

	// Move past 'in' to start of iterable expression
	p.nextToken()

	// Parse the iterable expression
	forStmt.Iterable = p.parseExpression(ast.LOWEST)
	if forStmt.Iterable == nil {
		return nil
	}

	// Expect 'do' keyword after the iterable
	if !p.expectPeek(lexer.DO) {
		p.addError(fmt.Sprintf("Expected 'do' after iterable, got %s at line %d, column %d",
			p.peekToken.Type, p.peekToken.Line, p.peekToken.Column))
		return nil
	}

	// Move past 'do' to the first statement in the body
	p.nextToken()

	// Parse the body statements until we reach 'end'
	forStmt.Body = p.parseBlockStatements(lexer.END)

	// After parsing the body, we should be at 'end'
	if p.curTokenIs(lexer.END) {
		// Move past 'end'
		p.nextToken()
	} else {
		p.addError(fmt.Sprintf("Expected 'end' to close for loop, got %s at line %d, column %d",
			p.curToken.Type, p.curToken.Line, p.curToken.Column))
	}

	return forStmt
}

// parseFunctionParameters parses function parameters
func (p *Parser) parseFunctionParameters() []ast.Parameter {
	parameters := []ast.Parameter{}

	// Check for empty parameter list
	if p.peekTokenIs(lexer.RPAREN) {
		p.nextToken() // consume the closing parenthesis
		p.nextToken() // move past the closing parenthesis
		return parameters
	}

	// Move to the first parameter name
	p.nextToken()

	// Parse first parameter
	if !p.curTokenIs(lexer.IDENT) {
		msg := fmt.Sprintf("Expected parameter name to be an identifier, got %s instead at line %d, column %d",
			p.curToken.Type, p.curToken.Line, p.curToken.Column)
		p.errors = append(p.errors, msg)
		return parameters
	}

	// Create first parameter
	param := ast.Parameter{
		Name: p.curToken.Literal,
	}

	// Check for type annotation for first parameter
	if p.peekTokenIs(lexer.COLON) {
		p.nextToken() // consume colon
		p.nextToken() // move to type name

		param.Type = &ast.TypeAnnotation{
			TypeName: p.curToken.Literal,
		}

		// Move to the next token after type name
		p.nextToken()
	} else {
		// Move to the next token (comma or closing parenthesis)
		p.nextToken()
	}

	parameters = append(parameters, param)

	// Parse additional parameters
	for p.curTokenIs(lexer.COMMA) {
		p.nextToken() // move to parameter name

		if !p.curTokenIs(lexer.IDENT) {
			msg := fmt.Sprintf("Expected parameter name to be an identifier, got %s instead at line %d, column %d",
				p.curToken.Type, p.curToken.Line, p.curToken.Column)
			p.errors = append(p.errors, msg)
			break
		}

		// Create parameter
		nextParam := ast.Parameter{
			Name: p.curToken.Literal,
		}

		// Check for type annotation
		if p.peekTokenIs(lexer.COLON) {
			p.nextToken() // consume colon
			p.nextToken() // move to type name

			nextParam.Type = &ast.TypeAnnotation{
				TypeName: p.curToken.Literal,
			}

			// Move to the next token after type name
			p.nextToken()
		} else {
			// Move to the next token (comma or closing parenthesis)
			p.nextToken()
		}

		parameters = append(parameters, nextParam)
	}

	// Check for closing parenthesis
	if !p.curTokenIs(lexer.RPAREN) {
		msg := fmt.Sprintf("Expected next token to be %s, got %s instead at line %d, column %d",
			lexer.RPAREN, p.curToken.Type, p.curToken.Line, p.curToken.Column)
		p.addError(msg)

		// Try to recover by skipping tokens until we find a closing parenthesis or reach another significant token
		for !p.curTokenIs(lexer.RPAREN) && !p.curTokenIs(lexer.EOF) &&
			!p.curTokenIs(lexer.COLON) && !p.curTokenIs(lexer.DO) {
			p.nextToken()
		}
	}

	// If we successfully found the closing parenthesis, advance past it
	if p.curTokenIs(lexer.RPAREN) {
		p.nextToken() // Move past the closing parenthesis
	}

	return parameters
}

// parseFunctionDefinition parses a function definition
func (p *Parser) parseFunctionDefinition() ast.Node {
	// Skip the 'def' keyword
	p.nextToken()

	var funcName string
	var isAnonymous bool

	// Check if it's an anonymous function (if the current token is a left parenthesis)
	if p.curTokenIs(lexer.LPAREN) {
		isAnonymous = true
	} else if p.curTokenIs(lexer.IDENT) {
		// Regular named function
		funcName = p.curToken.Literal
		p.nextToken() // Move to the next token
	} else {
		msg := fmt.Sprintf("Expected function name to be an identifier, got %s instead at line %d, column %d",
			p.curToken.Type, p.curToken.Line, p.curToken.Column)
		p.errors = append(p.errors, msg)
		return nil
	}

	// At this point, we should be at a left parenthesis
	if !p.curTokenIs(lexer.LPAREN) {
		msg := fmt.Sprintf("Expected '(' after function name, got %s instead at line %d, column %d",
			p.curToken.Type, p.curToken.Line, p.curToken.Column)
		p.errors = append(p.errors, msg)
		return nil
	}

	// Parse function parameters
	parameters := p.parseFunctionParameters()

	// For named functions, check that all parameters have type annotations
	if !isAnonymous {
		for _, param := range parameters {
			if param.Type == nil {
				msg := fmt.Sprintf("Missing type annotation for parameter '%s' at line %d, column %d",
					param.Name, p.curToken.Line, p.curToken.Column)
				p.errors = append(p.errors, msg)
			}
		}
	}

	// Check for return type annotation - only required for named functions
	var returnType *ast.TypeAnnotation
	if p.curTokenIs(lexer.COLON) {
		p.nextToken() // consume colon

		if !p.curTokenIs(lexer.IDENT) {
			msg := fmt.Sprintf("Expected return type to be an identifier, got %s instead at line %d, column %d",
				p.curToken.Type, p.curToken.Line, p.curToken.Column)
			p.errors = append(p.errors, msg)
		} else {
			returnType = &ast.TypeAnnotation{
				TypeName: p.curToken.Literal,
			}
		}

		p.nextToken() // move past the return type
	} else if !isAnonymous {
		// Return type annotation is required for named functions only
		msg := fmt.Sprintf("Missing return type annotation for function '%s' at line %d, column %d",
			funcName, p.curToken.Line, p.curToken.Column)
		p.errors = append(p.errors, msg)
	}

	// Check for 'do' keyword to start function body
	if !p.curTokenIs(lexer.DO) {
		msg := fmt.Sprintf("Expected 'do' after function declaration, got %s instead at line %d, column %d",
			p.curToken.Type, p.curToken.Line, p.curToken.Column)
		p.errors = append(p.errors, msg)
		return nil
	}

	// Parse function body
	p.nextToken() // move past 'do'
	body := p.parseBlockStatements(lexer.END)

	// Check for 'end' keyword to close function definition
	if !p.curTokenIs(lexer.END) {
		msg := fmt.Sprintf("Expected 'end' to close function definition, got %s instead at line %d, column %d",
			p.curToken.Type, p.curToken.Line, p.curToken.Column)
		p.errors = append(p.errors, msg)
		return nil
	}

	// Consume the 'end' token
	p.nextToken()

	return &ast.FunctionDef{
		Name:       funcName,
		Parameters: parameters,
		ReturnType: returnType,
		Body:       body,
	}
}