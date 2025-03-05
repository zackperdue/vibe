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
		// Store current token position in case we need to recover
		startToken := p.curToken
		startPeekToken := p.peekToken

		// Try to parse the statement
		stmt := p.parseStatement()

		// If statement parsing failed but we're at an identifier followed by a parenthesis,
		// it might be a function call that failed to parse due to a missing closing parenthesis
		if stmt == nil && startToken.Type == lexer.IDENT && startPeekToken.Type == lexer.LPAREN {
			// Manually create a function call node
			funcName := startToken.Literal

			// Create a simple function call with no arguments
			stmt = &ast.CallExpr{
				Function: &ast.Identifier{Name: funcName},
				Args:     []ast.Node{},
			}

			// Skip past the function name and opening parenthesis
			p.nextToken() // to the opening parenthesis
			p.nextToken() // past the opening parenthesis

			// If we're not already at a closing parenthesis, we might have an argument
			if !p.curTokenIs(lexer.RPAREN) {
				// Try to parse one argument
				arg := p.parseExpression(ast.LOWEST)
				if arg != nil {
					stmt.(*ast.CallExpr).Args = append(stmt.(*ast.CallExpr).Args, arg)
				}

				// Now skip to the next statement
				for !p.curTokenIs(lexer.EOF) &&
					  !containsTokenType(p.curToken.Type, endTokens) &&
					  p.curToken.Line == startToken.Line {
					p.nextToken()
				}
			}
		}

		if stmt != nil {
			block.Statements = append(block.Statements, stmt)
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


	// Create an empty block for the body
	forStmt.Body = &ast.BlockStmt{Statements: []ast.Node{}}

	// Check if we're already at END (which means we missed DO)
	if p.curTokenIs(lexer.END) {
		// This is a syntax error - we need a DO before END
		p.addError(fmt.Sprintf("Expected 'do' before 'end' in for loop at line %d, column %d",
			p.curToken.Line, p.curToken.Column))
		return nil
	}

	// Check if current token is DO
	if !p.curTokenIs(lexer.DO) {
		p.addError(fmt.Sprintf("Expected 'do' after iterable, got %s at line %d, column %d",
			p.curToken.Type, p.curToken.Line, p.curToken.Column))
		return nil
	}


	// Check if the body is empty (do end)
	if p.peekTokenIs(lexer.END) {
		p.nextToken() // Move to END token
	} else {
		// Parse the body - this will advance to the END token
		p.nextToken() // Move past 'do' to the first statement in the body
		forStmt.Body = p.parseBlockStatements(lexer.END)
	}


	// After parsing the body, we should be at 'end'
	if !p.curTokenIs(lexer.END) {
		p.addError(fmt.Sprintf("Expected 'end' to close for loop, got %s at line %d, column %d",
			p.curToken.Type, p.curToken.Line, p.curToken.Column))
		return nil
	}

	p.nextToken() // Move past 'end'


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

	// At this point, we should already be at the closing parenthesis
	if !p.curTokenIs(lexer.RPAREN) {
		p.addError(fmt.Sprintf("Expected closing parenthesis, got %s at line %d, column %d",
			p.curToken.Type, p.curToken.Line, p.curToken.Column))
	} else {
		// Move past the closing parenthesis
		p.nextToken()
	}

	return parameters
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

			typeParams = append(typeParams, param)

			// If we have a comma, advance past it to the next parameter
			if p.curTokenIs(lexer.COMMA) {
				p.nextToken()
			} else if !p.curTokenIs(lexer.GT) {
				// If not a comma and not the closing '>', something is wrong
				p.addError(fmt.Sprintf("Expected ',' or '>', got %s", p.curToken.Type))
				break
			}
		}

		// Set the type parameters on the annotation
		typeAnnotation.TypeParams = typeParams

		// Advance past the closing '>'
		if p.curTokenIs(lexer.GT) {
			p.nextToken()
		}
	} else {
		// No generic parameters, just advance past the type name
		p.nextToken()
	}

	return typeAnnotation
}

// parseClassDefinition parses a class definition
func (p *Parser) parseClassDefinition() ast.Node {
	// Current token is 'class'

	// Expect the next token to be the class name
	if !p.expectPeek(lexer.IDENT) {
		return nil
	}

	className := p.curToken.Literal

	// Check for inheritance
	var parentClass string
	if p.peekTokenIs(lexer.INHERITS) {
		p.nextToken() // Advance to 'inherits'

		// Expect parent class name
		if !p.expectPeek(lexer.IDENT) {
			return nil
		}

		parentClass = p.curToken.Literal
	}

	// Expect 'do' keyword to start the class body
	if !p.expectPeek(lexer.DO) {
		return nil
	}

	// Parse class body
	p.nextToken() // Move past 'do'

	// Create a slice to hold the methods
	methods := []ast.Node{}

	// Parse methods until we reach 'end'
	for !p.curTokenIs(lexer.END) && !p.curTokenIs(lexer.EOF) {
		// Skip any non-method tokens
		if !p.curTokenIs(lexer.FUNCTION) {
			p.nextToken()
			continue
		}

		// Parse method definition
		method := p.parseFunctionDefinition()
		if method != nil {
			methods = append(methods, method)
		}
	}

	// Expect 'end' to close the class definition
	if !p.curTokenIs(lexer.END) {
		p.addError(fmt.Sprintf("Expected 'end' to close class definition, got %s instead at line %d, column %d",
			p.curToken.Type, p.curToken.Line, p.curToken.Column))
		return nil
	}

	// Consume the 'end' token
	p.nextToken()

	return &ast.ClassDef{
		Name:    className,
		Parent:  parentClass,
		Methods: methods,
	}
}

// parseFunctionDefinition parses a function definition in the form:
// def name(param1: type1, param2: type2, ...): returnType do
//   statements...
// end
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