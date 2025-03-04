package parser

import (
	"fmt"

	"github.com/vibe-lang/vibe/ast"
	"github.com/vibe-lang/vibe/lexer"
)

// parseVariableDeclaration parses a variable declaration with type annotation
func (p *Parser) parseVariableDeclaration() ast.Node {
	fmt.Println("DEBUG: parseVariableDeclaration called")

	// Get the variable name
	name := p.curToken.Literal
	fmt.Println("DEBUG: Variable name:", name)
	p.nextToken()

	// Check for type annotation
	var typeAnnotation *ast.TypeAnnotation
	if p.curToken.Type == lexer.COLON {
		fmt.Println("DEBUG: Found colon, parsing type annotation")
		p.nextToken() // Skip ':'
		typeAnnotation = p.parseTypeAnnotation()
		fmt.Println("DEBUG: Type annotation:", typeAnnotation.String())
	}

	// Check for initialization
	var value ast.Node
	if p.curToken.Type == lexer.ASSIGN {
		fmt.Println("DEBUG: Found assignment, parsing value")
		p.nextToken() // Skip '='
		value = p.parseExpression(ast.LOWEST)
		if value != nil {
			fmt.Println("DEBUG: Value:", value.String())
		}
	}

	result := &ast.VariableDecl{
		Name:           name,
		TypeAnnotation: typeAnnotation,
		Value:          value,
	}
	fmt.Println("DEBUG: Created VariableDecl node:", result.String())
	return result
}

// parseTypeDeclaration parses a type declaration
func (p *Parser) parseTypeDeclaration() ast.Node {
	p.nextToken() // Skip 'type'

	if p.curToken.Type != lexer.IDENT {
		p.addError("Expected type name after 'type' keyword")
		return nil
	}

	name := p.curToken.Literal
	p.nextToken()

	if p.curToken.Type != lexer.ASSIGN {
		p.addError("Expected '=' after type name")
		return nil
	}

	p.nextToken() // Skip '='
	typeValue := p.parseTypeAnnotation()

	// Create the type declaration node
	return &ast.TypeDeclaration{
		Name:      name,
		TypeValue: typeValue,
	}
}

// parseClassInstantiation parses a class instantiation expression
func (p *Parser) parseClassInstantiation(class ast.Node) ast.Node {
	classInst := &ast.ClassInst{
		Class:     class,
		Arguments: []ast.Node{},
	}

	// Skip to '('
	if !p.expectPeek(lexer.LPAREN) {
		return nil
	}

	// Skip '('
	p.nextToken()

	// Handle empty arguments list
	if p.curToken.Type == lexer.RPAREN {
		p.nextToken() // Skip ')'
		return classInst
	}

	// Parse first argument
	arg := p.parseExpression(ast.LOWEST)
	classInst.Arguments = append(classInst.Arguments, arg)

	// Parse additional arguments
	for p.peekTokenIs(lexer.COMMA) {
		p.nextToken() // Skip comma
		p.nextToken() // Move to next argument

		arg := p.parseExpression(ast.LOWEST)
		classInst.Arguments = append(classInst.Arguments, arg)
	}

	if !p.expectPeek(lexer.RPAREN) {
		return nil
	}

	return classInst
}

// parseSuperCall parses a super call
func (p *Parser) parseSuperCall() ast.Node {
	// Skip 'super' token
	p.nextToken()

	// Check if this is a direct super call (super(...)) or a method call (super.method(...))
	if p.curToken.Type == lexer.LPAREN {
		// Direct super call
		superCall := &ast.SuperCall{
			Args: []ast.Node{},
		}

		// Skip '(' token
		p.nextToken()

		// Parse arguments if any
		if p.curToken.Type != lexer.RPAREN {
			// Parse first argument
			arg := p.parseExpression(ast.LOWEST)
			superCall.Args = append(superCall.Args, arg)

			// Parse additional arguments
			for p.peekTokenIs(lexer.COMMA) {
				p.nextToken() // Skip comma
				p.nextToken() // Move to next argument

				arg := p.parseExpression(ast.LOWEST)
				superCall.Args = append(superCall.Args, arg)
			}
		}

		if !p.expectPeek(lexer.RPAREN) {
			return nil
		}

		return superCall
	} else if p.curToken.Type == lexer.DOT {
		// Method call on super (super.method(...))
		p.nextToken() // Skip '.'

		if p.curToken.Type != lexer.IDENT {
			p.addError(fmt.Sprintf("Expected method name after 'super.', got %s", p.curToken.Type))
			return nil
		}

		methodName := p.curToken.Literal
		p.nextToken()

		// Check for opening parenthesis
		if p.curToken.Type != lexer.LPAREN {
			p.addError(fmt.Sprintf("Expected '(' after method name, got %s", p.curToken.Type))
			return nil
		}

		superCall := &ast.SuperCall{
			Method: methodName,
			Args:   []ast.Node{},
		}

		// Skip '('
		p.nextToken()

		// Parse arguments if any
		if p.curToken.Type != lexer.RPAREN {
			// Parse first argument
			arg := p.parseExpression(ast.LOWEST)
			superCall.Args = append(superCall.Args, arg)

			// Parse additional arguments
			for p.peekTokenIs(lexer.COMMA) {
				p.nextToken() // Skip comma
				p.nextToken() // Move to next argument

				arg := p.parseExpression(ast.LOWEST)
				superCall.Args = append(superCall.Args, arg)
			}
		}

		if !p.expectPeek(lexer.RPAREN) {
			return nil
		}

		return superCall
	} else {
		p.addError(fmt.Sprintf("Expected '(' or '.' after 'super', got %s", p.curToken.Type))
		return nil
	}
}