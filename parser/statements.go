package parser

import (
	"fmt"

	"github.com/example/vibe/ast"
	"github.com/example/vibe/lexer"
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
		fmt.Printf("DEBUG parseBlockStatements: current token: %s (%s) at line %d, column %d, peek token: %s (%s) at line %d, column %d\n",
			p.curToken.Type, p.curToken.Literal, p.curToken.Line, p.curToken.Column,
			p.peekToken.Type, p.peekToken.Literal, p.peekToken.Line, p.peekToken.Column)

		stmt := p.parseStatement()
		if stmt != nil {
			fmt.Printf("DEBUG: Added statement: %T\n", stmt)
			block.Statements = append(block.Statements, stmt)
		} else {
			fmt.Println("DEBUG: Statement was nil, not adding to block")
		}

		// If we've reached an end token, don't advance any further
		if containsTokenType(p.curToken.Type, endTokens) {
			break
		}

		// Advance to the next token
		p.nextToken()
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

	// Skip 'for' keyword and expect an identifier for iterator
	p.nextToken()
	if p.curToken.Type != lexer.IDENT {
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

	// After parsing the iterable expression, we need to check for the DO token
	if !p.curTokenIs(lexer.DO) {
		if !p.expectPeek(lexer.DO) {
			p.addError(fmt.Sprintf("Expected 'do' after iterable, got %s at line %d, column %d",
				p.peekToken.Type, p.peekToken.Line, p.peekToken.Column))
			return nil
		}
	}

	// Parse the body - this will advance to the END token
	p.nextToken() // Move past 'do' to the first statement in the body
	forStmt.Body = p.parseBlockStatements(lexer.END)

	// After parsing the body, we should be at 'end'
	if !p.curTokenIs(lexer.END) {
		p.addError(fmt.Sprintf("Expected 'end' to close for loop, got %s at line %d, column %d",
			p.curToken.Type, p.curToken.Line, p.curToken.Column))
		return nil
	}

	// Move past 'end'
	p.nextToken()

	return forStmt
}

// parseFunctionDefinition parses a function definition
func (p *Parser) parseFunctionDefinition() ast.Node {
	// Current token is 'function'

	// Move to function name
	if !p.expectPeek(lexer.IDENT) {
		p.addError(fmt.Sprintf("Expected function name, got %s", p.peekToken.Type))
		return nil
	}

	name := p.curToken.Literal

	// Parse parameters
	parameters := p.parseFunctionParameters()

	// Check for return type annotation
	var returnType *ast.TypeAnnotation
	if p.curTokenIs(lexer.COLON) {
		p.nextToken() // Skip ':'
		returnType = p.parseTypeAnnotation()
	}

	// Expect 'do' keyword
	if !p.expectPeek(lexer.DO) {
		p.addError(fmt.Sprintf("Expected 'do' after function definition, got %s", p.peekToken.Type))
		return nil
	}

	// Parse function body
	body := p.parseBlockStatements(lexer.END)

	// Ensure 'end' token is consumed
	if !p.curTokenIs(lexer.END) {
		p.addError(fmt.Sprintf("Expected 'end' to close function body, got %s", p.curToken.Type))
		return nil
	}

	p.nextToken() // Skip 'end'

	return &ast.FunctionDef{
		Name:       name,
		Parameters: parameters,
		ReturnType: returnType,
		Body:       body,
	}
}

// parseFunctionParameters parses function parameters with optional type annotations
func (p *Parser) parseFunctionParameters() []ast.Parameter {
	var parameters []ast.Parameter

	if !p.expectPeek(lexer.LPAREN) {
		return parameters
	}

	// Handle empty parameter list
	if p.peekTokenIs(lexer.RPAREN) {
		p.nextToken() // Move to the closing parenthesis
		p.nextToken() // Move past the closing parenthesis to the next token
		return parameters
	}

	p.nextToken() // Move past opening parenthesis to first parameter

	// Parse parameters
	for !p.curTokenIs(lexer.RPAREN) && !p.curTokenIs(lexer.EOF) {
		// Parse parameter name
		if !p.curTokenIs(lexer.IDENT) {
			p.addError(fmt.Sprintf("Expected parameter name, got %s", p.curToken.Type))
			return parameters
		}

		param := ast.Parameter{Name: p.curToken.Literal}
		p.nextToken() // Move past parameter name

		// Parse type annotation if present
		if p.curTokenIs(lexer.COLON) {
			p.nextToken() // Move to type name
			param.Type = p.parseTypeAnnotation()
			// Note: parseTypeAnnotation will have advanced the token to the position AFTER the type
		}

		parameters = append(parameters, param)

		// If the current token is a comma, move past it and continue parsing more parameters
		if p.curTokenIs(lexer.COMMA) {
			p.nextToken() // Skip comma
			continue
		}

		// At this point, we should be at a closing parenthesis
		break
	}

	// Check and handle closing parenthesis
	if !p.curTokenIs(lexer.RPAREN) {
		p.addError(fmt.Sprintf("Expected ',' or ')' after parameter, got %s", p.curToken.Type))
	} else {
		p.nextToken() // Skip the closing parenthesis, advancing to the next token (which could be a colon for return type or 'do')
	}

	return parameters
}

// parseTypeAnnotation parses a type annotation
func (p *Parser) parseTypeAnnotation() *ast.TypeAnnotation {
	if !p.curTokenIs(lexer.IDENT) {
		p.addError(fmt.Sprintf("Expected type name, got %s", p.curToken.Type))
		return &ast.TypeAnnotation{TypeName: "any"} // Default type
	}

	typeName := p.curToken.Literal

	// Create the type annotation
	typeAnnotation := &ast.TypeAnnotation{
		TypeName: typeName,
	}

	// Move past the type name
	p.nextToken()

	// Check for generic type parameters
	if p.curTokenIs(lexer.LT) {
		p.nextToken() // Skip '<'

		var typeParams []ast.Node
		// Parse the first type parameter
		typeParam := p.parseTypeAnnotation()
		typeParams = append(typeParams, typeParam)

		// Parse additional type parameters
		for p.curTokenIs(lexer.COMMA) {
			p.nextToken() // Skip ','
			typeParam := p.parseTypeAnnotation()
			typeParams = append(typeParams, typeParam)
		}

		typeAnnotation.TypeParams = typeParams

		// Expect closing '>'
		if !p.curTokenIs(lexer.GT) {
			p.addError(fmt.Sprintf("Expected closing '>' after type parameters, got %s", p.curToken.Type))
		} else {
			p.nextToken() // Move past the closing '>' to the next token (which is usually a comma or closing parenthesis)
		}
	}

	return typeAnnotation
}

// parseClassDefinition parses a class definition
func (p *Parser) parseClassDefinition() ast.Node {
	// Current token is 'class'

	// Move to class name
	if !p.expectPeek(lexer.IDENT) {
		p.addError(fmt.Sprintf("Expected class name, got %s", p.peekToken.Type))
		return nil
	}

	className := p.curToken.Literal

	// Check for inheritance
	var parentClass string
	if p.peekTokenIs(lexer.INHERITS) {
		p.nextToken() // Move to 'inherits'

		// Move to parent class name
		if !p.expectPeek(lexer.IDENT) {
			p.addError(fmt.Sprintf("Expected parent class name, got %s", p.peekToken.Type))
			return nil
		}

		parentClass = p.curToken.Literal
	}

	// Expect 'do' keyword
	if !p.expectPeek(lexer.DO) {
		p.addError(fmt.Sprintf("Expected 'do' after class definition, got %s", p.peekToken.Type))
		return nil
	}

	// Parse class body
	body := p.parseBlockStatements(lexer.END)

	// Ensure 'end' token is consumed
	if !p.curTokenIs(lexer.END) {
		p.addError(fmt.Sprintf("Expected 'end' to close class body, got %s", p.curToken.Type))
		return nil
	}

	p.nextToken() // Skip 'end'

	return &ast.ClassDef{
		Name:    className,
		Parent:  parentClass,
		Methods: body.Statements,
	}
}