package parser

import (
	"fmt"
	"strconv"

	"github.com/example/vibe/lexer"
)

// Different types of AST nodes
type NodeType string

const (
	ProgramNode      NodeType = "Program"
	NumberNode       NodeType = "Number"
	StringNode       NodeType = "String"
	IdentifierNode   NodeType = "Identifier"
	BinaryExprNode   NodeType = "BinaryExpr"
	CallExprNode     NodeType = "CallExpr"
	FunctionDefNode  NodeType = "FunctionDef"
	ReturnStmtNode   NodeType = "ReturnStmt"
	IfStmtNode       NodeType = "IfStmt"
	WhileStmtNode    NodeType = "WhileStmt"
	BlockStmtNode    NodeType = "BlockStmt"
	AssignmentNode   NodeType = "Assignment"
	VariableUseNode  NodeType = "VariableUse"
	BooleanNode      NodeType = "Boolean"
	NilNode          NodeType = "Nil"
	PrintStmtNode    NodeType = "PrintStmt"
	TypeAnnotationNode NodeType = "TypeAnnotation"
	TypeDeclarationNode NodeType = "TypeDeclaration"
	VariableDeclNode NodeType = "VariableDecl"
	UnaryExprNode    NodeType = "UnaryExpr"
	ArrayLiteralNode   NodeType = "ArrayLiteral"
	IndexExprNode    NodeType = "IndexExpr"
	DotExprNode      NodeType = "DotExpr"
)

// Operator precedence
const (
	LOWEST     = 1
	EQUALS     = 2  // ==
	LESSGREATER = 3  // > or <
	SUM        = 4  // +
	PRODUCT    = 5  // *
	PREFIX     = 6  // -X or !X
	CALL       = 7  // myFunction(X)
	INDEX      = 8  // array[index]
	DOT        = 9  // obj.property
)

// Node represents a node in the AST
type Node interface {
	Type() NodeType
	String() string
}

// Program is the root node of the AST
type Program struct {
	Statements []Node
}

func (p *Program) Type() NodeType { return ProgramNode }
func (p *Program) String() string {
	result := "Program {\n"
	for _, stmt := range p.Statements {
		result += "  " + stmt.String() + "\n"
	}
	result += "}"
	return result
}

// NumberLiteral represents a number literal
type NumberLiteral struct {
	Value float64
	IsInt bool
}

func (n *NumberLiteral) Type() NodeType { return NumberNode }
func (n *NumberLiteral) String() string {
	if n.IsInt {
		return fmt.Sprintf("Number(%d)", int(n.Value))
	}
	return fmt.Sprintf("Number(%f)", n.Value)
}

// StringLiteral represents a string literal
type StringLiteral struct {
	Value string
}

func (s *StringLiteral) Type() NodeType { return StringNode }
func (s *StringLiteral) String() string { return fmt.Sprintf("String(%q)", s.Value) }

// Identifier represents a variable or function name
type Identifier struct {
	Name string
}

func (i *Identifier) Type() NodeType { return IdentifierNode }
func (i *Identifier) String() string { return fmt.Sprintf("Identifier(%s)", i.Name) }

// BinaryExpr represents a binary expression (e.g. a + b)
type BinaryExpr struct {
	Left     Node
	Operator string
	Right    Node
}

func (b *BinaryExpr) Type() NodeType { return BinaryExprNode }
func (b *BinaryExpr) String() string {
	leftStr := "nil"
	rightStr := "nil"

	if b.Left != nil {
		leftStr = b.Left.String()
	}

	if b.Right != nil {
		rightStr = b.Right.String()
	}

	return fmt.Sprintf("BinaryExpr(%s %s %s)", leftStr, b.Operator, rightStr)
}

// CallExpr represents a function call
type CallExpr struct {
	Function Node
	Args     []Node
}

func (c *CallExpr) Type() NodeType { return CallExprNode }
func (c *CallExpr) String() string {
	result := fmt.Sprintf("CallExpr(%s, [", c.Function.String())
	for i, arg := range c.Args {
		if i > 0 {
			result += ", "
		}
		result += arg.String()
	}
	result += "])"
	return result
}

// FunctionDef represents a function definition
type FunctionDef struct {
	Name       string
	Parameters []string
	ParamTypes []*TypeAnnotation
	ReturnType *TypeAnnotation
	Body       *BlockStmt
}

func (f *FunctionDef) Type() NodeType { return FunctionDefNode }
func (f *FunctionDef) String() string {
	result := fmt.Sprintf("FunctionDef(%s, [", f.Name)
	for i, param := range f.Parameters {
		if i > 0 {
			result += ", "
		}
		result += param
	}
	result += "], " + f.Body.String() + ")"
	return result
}

// ReturnStmt represents a return statement
type ReturnStmt struct {
	Value Node
}

func (r *ReturnStmt) Type() NodeType { return ReturnStmtNode }
func (r *ReturnStmt) String() string {
	if r.Value == nil {
		return "ReturnStmt(nil)"
	}
	return fmt.Sprintf("ReturnStmt(%s)", r.Value.String())
}

// IfStmt represents an if statement
type IfStmt struct {
	Condition     Node
	Consequence   *BlockStmt
	Alternative   *BlockStmt
	ElseIfBlocks  []ElseIfBlock
}

type ElseIfBlock struct {
	Condition   Node
	Consequence *BlockStmt
}

func (i *IfStmt) Type() NodeType { return IfStmtNode }
func (i *IfStmt) String() string {
	result := fmt.Sprintf("IfStmt(%s, %s", i.Condition.String(), i.Consequence.String())
	for _, elseIf := range i.ElseIfBlocks {
		result += fmt.Sprintf(", ElseIf(%s, %s)", elseIf.Condition.String(), elseIf.Consequence.String())
	}
	if i.Alternative != nil {
		result += fmt.Sprintf(", Else(%s)", i.Alternative.String())
	}
	result += ")"
	return result
}

// WhileStmt represents a while loop
type WhileStmt struct {
	Condition Node
	Body      *BlockStmt
}

func (w *WhileStmt) Type() NodeType { return WhileStmtNode }
func (w *WhileStmt) String() string {
	return fmt.Sprintf("WhileStmt(%s, %s)", w.Condition.String(), w.Body.String())
}

// BlockStmt represents a block of statements
type BlockStmt struct {
	Statements []Node
}

func (b *BlockStmt) Type() NodeType { return BlockStmtNode }
func (b *BlockStmt) String() string {
	result := "Block {\n"
	for _, stmt := range b.Statements {
		result += "  " + stmt.String() + "\n"
	}
	result += "}"
	return result
}

// Assignment represents a variable assignment
type Assignment struct {
	Name  string
	Value Node
}

func (a *Assignment) Type() NodeType { return AssignmentNode }
func (a *Assignment) String() string {
	return fmt.Sprintf("Assignment(%s = %s)", a.Name, a.Value.String())
}

// BooleanLiteral represents a boolean value
type BooleanLiteral struct {
	Value bool
}

func (b *BooleanLiteral) Type() NodeType { return BooleanNode }
func (b *BooleanLiteral) String() string { return fmt.Sprintf("Boolean(%t)", b.Value) }

// NilLiteral represents a nil value
type NilLiteral struct{}

func (n *NilLiteral) Type() NodeType { return NilNode }
func (n *NilLiteral) String() string { return "Nil" }

// PrintStmt represents a print statement
type PrintStmt struct {
	Value Node
}

func (p *PrintStmt) Type() NodeType { return PrintStmtNode }
func (p *PrintStmt) String() string { return fmt.Sprintf("PrintStmt(%s)", p.Value.String()) }

// TypeAnnotation represents a type annotation
type TypeAnnotation struct {
	TypeName   string
	TypeParams []Node // For generic types like Array<string>
}

func (t *TypeAnnotation) Type() NodeType { return TypeAnnotationNode }
func (t *TypeAnnotation) String() string {
	if len(t.TypeParams) == 0 {
		return fmt.Sprintf("Type(%s)", t.TypeName)
	}

	params := ""
	for i, param := range t.TypeParams {
		if i > 0 {
			params += ", "
		}
		params += param.String()
	}

	return fmt.Sprintf("Type(%s<%s>)", t.TypeName, params)
}

// TypeDeclaration represents a type declaration (type aliases and interfaces)
type TypeDeclaration struct {
	Name      string
	TypeValue Node // Could be a TypeAnnotation or another structure
}

func (t *TypeDeclaration) Type() NodeType { return TypeDeclarationNode }
func (t *TypeDeclaration) String() string {
	return fmt.Sprintf("TypeDecl(%s = %s)", t.Name, t.TypeValue.String())
}

// VariableDecl represents a variable declaration with a type
type VariableDecl struct {
	Name           string
	TypeAnnotation *TypeAnnotation
	Value          Node // Initial value (can be nil)
}

func (v *VariableDecl) Type() NodeType { return VariableDeclNode }
func (v *VariableDecl) String() string {
	initialValue := "nil"
	if v.Value != nil {
		initialValue = v.Value.String()
	}
	return fmt.Sprintf("VarDecl(%s: %s = %s)", v.Name, v.TypeAnnotation.String(), initialValue)
}

// UnaryExpr represents a unary expression like !x or -5
type UnaryExpr struct {
	Operator string
	Right    Node
}

func (u *UnaryExpr) Type() NodeType { return UnaryExprNode }
func (u *UnaryExpr) String() string {
	if u.Right == nil {
		return fmt.Sprintf("(%s<nil>)", u.Operator)
	}
	return fmt.Sprintf("(%s%s)", u.Operator, u.Right.String())
}

// ArrayLiteral represents an array literal
type ArrayLiteral struct {
	Elements []Node
}

func (a *ArrayLiteral) Type() NodeType { return ArrayLiteralNode }
func (a *ArrayLiteral) String() string {
	result := "["
	for i, elem := range a.Elements {
		if i > 0 {
			result += ", "
		}
		result += elem.String()
	}
	result += "]"
	return result
}

// IndexExpr represents an index expression
type IndexExpr struct {
	Array Node
	Index Node
}

func (i *IndexExpr) Type() NodeType { return IndexExprNode }
func (i *IndexExpr) String() string {
	return fmt.Sprintf("%s[%s]", i.Array.String(), i.Index.String())
}

// DotExpr represents a dot expression
type DotExpr struct {
	Object   Node
	Property string
}

func (d *DotExpr) Type() NodeType { return DotExprNode }
func (d *DotExpr) String() string {
	return fmt.Sprintf("%s.%s", d.Object.String(), d.Property)
}

// Parser parses tokens into an AST
type Parser struct {
	l         *lexer.Lexer
	curToken  lexer.Token
	peekToken lexer.Token
	errors    []string
}

// New creates a new parser
func New(l *lexer.Lexer) *Parser {
	p := &Parser{l: l}
	// Read two tokens, so curToken and peekToken are both set
	p.nextToken()
	p.nextToken()
	return p
}

func (p *Parser) nextToken() {
	p.curToken = p.peekToken
	p.peekToken = p.l.NextToken()
}

func (p *Parser) Errors() []string {
	return p.errors
}

// Parse function creates a new parser and parses the program
func Parse(l *lexer.Lexer) (*Program, []string) {
	p := New(l)
	program := p.parseProgram()
	return program, p.errors
}

func (p *Parser) parseProgram() *Program {
	program := &Program{Statements: []Node{}}

	for p.curToken.Type != lexer.EOF {
		stmt := p.parseStatement()
		if stmt != nil {
			program.Statements = append(program.Statements, stmt)
		}
		p.nextToken()
	}

	return program
}

func (p *Parser) parseStatement() Node {
	switch p.curToken.Type {
	case lexer.FUNCTION:
		return p.parseFunctionDefinition()
	case lexer.IF:
		return p.parseIfStatement()
	case lexer.WHILE:
		return p.parseWhileStatement()
	case lexer.RETURN:
		return p.parseReturnStatement()
	case lexer.PRINT:
		return p.parsePrintStatement()
	case lexer.LET, lexer.VAR:
		return p.parseVariableDeclaration()
	case lexer.IDENT:
		// Check if it's an assignment
		if p.peekToken.Type == lexer.ASSIGN {
			return p.parseAssignment()
		}
		return p.parseExpressionStatement()
	default:
		return p.parseExpressionStatement()
	}
}

func (p *Parser) parseFunctionDefinition() Node {
	funcDef := &FunctionDef{}

	// Function name
	p.nextToken()
	if p.curToken.Type != lexer.IDENT {
		p.errors = append(p.errors, fmt.Sprintf("Expected function name, got %s", p.curToken.Type))
		return nil
	}
	funcDef.Name = p.curToken.Literal

	// Parameters
	p.nextToken()
	if p.curToken.Type != lexer.LPAREN {
		p.errors = append(p.errors, fmt.Sprintf("Expected '(' after function name, got %s", p.curToken.Type))
		return nil
	}

	funcDef.Parameters, funcDef.ParamTypes = p.parseFunctionParameters()

	// Check for return type annotation
	if p.curToken.Type == lexer.COLON {
		p.nextToken()
		funcDef.ReturnType = p.parseTypeAnnotation()
	}

	// Function body
	if p.curToken.Type != lexer.LBRACE {
		p.errors = append(p.errors, fmt.Sprintf("Expected '{' to start function body, got %s", p.curToken.Type))
		return nil
	}

	funcDef.Body = p.parseBlockStatement()

	return funcDef
}

func (p *Parser) parseTypeAnnotation() *TypeAnnotation {
	typeAnnotation := &TypeAnnotation{}
	var typeName string

	if p.curToken.Type == lexer.IDENT || p.curToken.Type == lexer.FUNCTION ||
	   p.curToken.Type == lexer.TRUE || p.curToken.Type == lexer.FALSE ||
	   p.curToken.Type == lexer.NIL {
		typeName = p.curToken.Literal
		p.nextToken()
	} else {
		p.errors = append(p.errors, fmt.Sprintf("Expected type name, got %s", p.curToken.Type))
		return nil
	}

	typeAnnotation.TypeName = typeName

	// Check for generic type parameters like Array<string>
	var typeParams []Node
	if p.curToken.Type == lexer.LT {
		p.nextToken() // Skip '<'

		// Parse the type parameter(s)
		for p.curToken.Type != lexer.GT {
			if p.curToken.Type == lexer.EOF {
				p.errors = append(p.errors, "Unexpected EOF while parsing type parameters")
				break
			}

			// Parse the type parameter (which is another type annotation)
			paramType := p.parseTypeAnnotation()
			if paramType != nil {
				typeParams = append(typeParams, paramType)
			}

			// Check for comma
			if p.curToken.Type == lexer.COMMA {
				p.nextToken() // Skip ','
			}
		}

		if p.curToken.Type == lexer.GT {
			p.nextToken() // Skip '>'
		}
	}

	// Handle union types with |
	if p.curToken.Type == lexer.OR {
		p.nextToken() // Skip '|'
		rightType := p.parseTypeAnnotation()
		if rightType != nil {
			// Create a union type
			unionType := &TypeAnnotation{
				TypeName: "union",
				TypeParams: []Node{
					typeAnnotation,
					rightType,
				},
			}
			return unionType
		}
	}

	typeAnnotation.TypeParams = typeParams
	return typeAnnotation
}

func (p *Parser) parseTypeDeclaration() *TypeDeclaration {
	p.nextToken() // Skip 'type'

	if p.curToken.Type != lexer.IDENT {
		p.errors = append(p.errors, "Expected type name after 'type' keyword")
		return nil
	}

	name := p.curToken.Literal
	p.nextToken()

	if p.curToken.Type != lexer.ASSIGN {
		p.errors = append(p.errors, "Expected '=' after type name")
		return nil
	}

	p.nextToken() // Skip '='
	typeValue := p.parseTypeAnnotation()

	return &TypeDeclaration{
		Name:      name,
		TypeValue: typeValue,
	}
}

func (p *Parser) parseVariableDeclaration() Node {
	var name string
	var typeAnnotation *TypeAnnotation

	// We only support variable declarations starting with the variable name
	if p.curToken.Type == lexer.IDENT {
		// Get the variable name
		name = p.curToken.Literal
		p.nextToken()

		// Check for type annotation
		if p.curToken.Type == lexer.COLON {
			p.nextToken() // Skip ':'
			typeAnnotation = p.parseTypeAnnotation()
		}
	} else {
		p.errors = append(p.errors, fmt.Sprintf("Expected variable name, got %s", p.curToken.Type))
		return nil
	}

	var value Node
	// Check for initialization
	if p.curToken.Type == lexer.ASSIGN {
		p.nextToken() // Skip '='
		value = p.parseExpression(LOWEST)
	}

	return &VariableDecl{
		Name:           name,
		TypeAnnotation: typeAnnotation,
		Value:          value,
	}
}

func (p *Parser) parseIfStatement() Node {
	// Skip 'if' keyword
	p.nextToken()

	condition := p.parseExpression(LOWEST)

	// Look for the block start
	consequence := p.parseBlockStatement()

	// Check for elsif or else
	var elseIfBlocks []ElseIfBlock
	var alternative *BlockStmt

	for p.peekToken.Type == lexer.ELSIF || p.peekToken.Type == lexer.ELSE {
		p.nextToken() // Move to elsif/else

		if p.curToken.Type == lexer.ELSIF {
			p.nextToken() // Move past elsif
			elsifCondition := p.parseExpression(LOWEST)

			// Look for block start
			elsifBlock := p.parseBlockStatement()
			elseIfBlocks = append(elseIfBlocks, ElseIfBlock{
				Condition:   elsifCondition,
				Consequence: elsifBlock,
			})
		} else if p.curToken.Type == lexer.ELSE {
			p.nextToken() // Move past else

			// Look for block start
			alternative = p.parseBlockStatement()
			break // Must be the last block
		}
	}

	return &IfStmt{
		Condition:    condition,
		Consequence:  consequence,
		ElseIfBlocks: elseIfBlocks,
		Alternative:  alternative,
	}
}

func (p *Parser) parseWhileStatement() Node {
	// Skip 'while' keyword
	p.nextToken()

	condition := p.parseExpression(LOWEST)

	// Look for the block start
	body := p.parseBlockStatement()

	return &WhileStmt{
		Condition: condition,
		Body:      body,
	}
}

func (p *Parser) parseReturnStatement() Node {
	// Skip 'return' keyword
	p.nextToken()

	// Check if return has no value
	if p.curToken.Type == lexer.SEMICOLON {
		return &ReturnStmt{Value: nil}
	}

	value := p.parseExpression(LOWEST)
	return &ReturnStmt{Value: value}
}

func (p *Parser) parsePrintStatement() Node {
	isPuts := p.curToken.Literal == "puts"
	// Skip 'puts' or 'print' keyword
	p.nextToken()

	value := p.parseExpression(LOWEST)

	// Create different print nodes based on the type
	if isPuts {
		return &PrintStmt{Value: value}
	} else {
		return &PrintStmt{Value: value}
	}
}

func (p *Parser) parseAssignment() Node {
	name := p.curToken.Literal
	p.nextToken() // Move to '='
	p.nextToken() // Move past '='

	value := p.parseExpression(LOWEST)
	return &Assignment{Name: name, Value: value}
}

func (p *Parser) parseExpressionStatement() Node {
	return p.parseExpression(0)
}

func (p *Parser) parseExpression(precedence int) Node {
	// Prefix parsing functions
	var left Node

	switch p.curToken.Type {
	case lexer.IDENT:
		left = &Identifier{Name: p.curToken.Literal}
	case lexer.INT:
		value, err := strconv.ParseInt(p.curToken.Literal, 0, 64)
		if err != nil {
			p.errors = append(p.errors, fmt.Sprintf("Could not parse %s as integer", p.curToken.Literal))
			return nil
		}
		left = &NumberLiteral{Value: float64(value), IsInt: true}
	case lexer.FLOAT:
		value, err := strconv.ParseFloat(p.curToken.Literal, 64)
		if err != nil {
			p.errors = append(p.errors, fmt.Sprintf("Could not parse %s as float", p.curToken.Literal))
			return nil
		}
		left = &NumberLiteral{Value: value, IsInt: false}
	case lexer.STRING:
		left = &StringLiteral{Value: p.curToken.Literal}
	case lexer.TRUE:
		left = &BooleanLiteral{Value: true}
	case lexer.FALSE:
		left = &BooleanLiteral{Value: false}
	case lexer.NIL:
		left = &NilLiteral{}
	case lexer.LPAREN:
		p.nextToken()
		left = p.parseExpression(LOWEST)
		if p.curToken.Type != lexer.RPAREN {
			p.errors = append(p.errors, fmt.Sprintf("Expected ), got %s", p.curToken.Type))
			return nil
		}
	case lexer.LBRACKET:
		left = p.parseArrayLiteral()
	case lexer.PLUS, lexer.MINUS, lexer.BANG:
		operator := p.curToken.Literal
		p.nextToken()
		operand := p.parseExpression(PREFIX)
		left = &UnaryExpr{
			Operator: operator,
			Right:    operand,
		}
	default:
		p.errors = append(p.errors, fmt.Sprintf("No prefix parser for %s", p.curToken.Type))
		return nil
	}

	// Infix parsing functions
	for precedence < p.peekPrecedence() {
		if !isInfixOperator(p.peekToken.Type) {
			return left
		}

		p.nextToken()

		switch p.curToken.Type {
		case lexer.PLUS, lexer.MINUS, lexer.ASTERISK, lexer.SLASH,
				lexer.EQ, lexer.NOT_EQ, lexer.LT, lexer.GT, lexer.LT_EQ, lexer.GT_EQ,
				lexer.AND, lexer.OR:
			operator := p.curToken.Literal
			rightPrecedence := p.curPrecedence()
			p.nextToken()
			right := p.parseExpression(rightPrecedence)
			left = &BinaryExpr{
				Left:     left,
				Operator: operator,
				Right:    right,
			}
		case lexer.LPAREN:
			left = p.parseCallExpression(left)
		case lexer.LBRACKET:
			left = p.parseIndexExpression(left)
		case lexer.DOT:
			left = p.parseDotExpression(left)
		default:
			return left
		}
	}

	return left
}

// Helper function to check if a token type is an infix operator
func isInfixOperator(tokenType lexer.TokenType) bool {
	switch tokenType {
	case lexer.PLUS, lexer.MINUS, lexer.ASTERISK, lexer.SLASH,
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
		return EQUALS
	case lexer.LT, lexer.GT, lexer.LT_EQ, lexer.GT_EQ:
		return LESSGREATER
	case lexer.PLUS, lexer.MINUS:
		return SUM
	case lexer.ASTERISK, lexer.SLASH:
		return PRODUCT
	case lexer.LPAREN:
		return CALL
	case lexer.LBRACKET:
		return INDEX
	case lexer.DOT:
		return DOT
	default:
		return LOWEST
	}
}

func (p *Parser) curPrecedence() int {
	switch p.curToken.Type {
	case lexer.EQ, lexer.NOT_EQ:
		return EQUALS
	case lexer.LT, lexer.GT, lexer.LT_EQ, lexer.GT_EQ:
		return LESSGREATER
	case lexer.PLUS, lexer.MINUS:
		return SUM
	case lexer.ASTERISK, lexer.SLASH:
		return PRODUCT
	case lexer.LPAREN:
		return CALL
	case lexer.LBRACKET:
		return INDEX
	case lexer.DOT:
		return DOT
	default:
		return LOWEST
	}
}

func (p *Parser) parseBlockStatement() *BlockStmt {
	block := &BlockStmt{Statements: []Node{}}

	p.nextToken() // Skip opening token

	for p.curToken.Type != lexer.RBRACE {
		if p.curToken.Type == lexer.EOF {
			p.errors = append(p.errors, "Unexpected EOF, expected '}'")
			return block
		}

		if p.curToken.Type == lexer.SEMICOLON {
			p.nextToken()
			continue
		}

		stmt := p.parseStatement()
		if stmt != nil {
			block.Statements = append(block.Statements, stmt)
		}
		p.nextToken()
	}

	return block
}

// parseFunctionParameters parses function parameters with optional type annotations
func (p *Parser) parseFunctionParameters() ([]string, []*TypeAnnotation) {
	var parameters []string
	var paramTypes []*TypeAnnotation

	p.nextToken() // Skip '('

	// Handle empty parameter list
	if p.curToken.Type == lexer.RPAREN {
		p.nextToken() // Skip ')'
		return parameters, paramTypes
	}

	// Read first parameter
	if p.curToken.Type == lexer.IDENT {
		parameters = append(parameters, p.curToken.Literal)

		// Check for type annotation
		p.nextToken()
		if p.curToken.Type == lexer.COLON {
			p.nextToken() // Skip ':'
			paramTypes = append(paramTypes, p.parseTypeAnnotation())
		} else {
			// No type annotation, add nil
			paramTypes = append(paramTypes, nil)
		}

		p.nextToken() // Move to ',' or ')'
	}

	// Read other parameters
	for p.curToken.Type == lexer.COMMA {
		p.nextToken() // Skip ','

		if p.curToken.Type != lexer.IDENT {
			p.errors = append(p.errors, fmt.Sprintf("Expected parameter name, got %s", p.curToken.Type))
			break
		}

		parameters = append(parameters, p.curToken.Literal)

		// Check for type annotation
		p.nextToken()
		if p.curToken.Type == lexer.COLON {
			p.nextToken() // Skip ':'
			paramTypes = append(paramTypes, p.parseTypeAnnotation())
		} else {
			// No type annotation, add nil
			paramTypes = append(paramTypes, nil)
		}

		p.nextToken() // Move to ',' or ')'
	}

	if p.curToken.Type != lexer.RPAREN {
		p.errors = append(p.errors, fmt.Sprintf("Expected ')' after parameters, got %s", p.curToken.Type))
	} else {
		p.nextToken() // Skip ')'
	}

	return parameters, paramTypes
}

func (p *Parser) parseArrayLiteral() Node {
	// Skip '['
	p.nextToken()

	var elements []Node

	// Handle empty array
	if p.curToken.Type == lexer.RBRACKET {
		p.nextToken() // Skip ']'
		return &ArrayLiteral{Elements: elements}
	}

	// Parse first element
	element := p.parseExpression(LOWEST)
	elements = append(elements, element)

	// Parse remaining elements
	for p.curToken.Type == lexer.COMMA {
		p.nextToken() // Skip ','
		element = p.parseExpression(LOWEST)
		elements = append(elements, element)
	}

	if p.curToken.Type != lexer.RBRACKET {
		p.errors = append(p.errors, fmt.Sprintf("Expected ']', got %s", p.curToken.Type))
		return nil
	}

	p.nextToken() // Skip ']'
	return &ArrayLiteral{Elements: elements}
}

func (p *Parser) parseCallExpression(function Node) Node {
	// Skip '('
	p.nextToken()

	var args []Node

	// Handle empty argument list
	if p.curToken.Type == lexer.RPAREN {
		p.nextToken() // Skip ')'
		return &CallExpr{Function: function, Args: args}
	}

	// Parse first argument
	arg := p.parseExpression(LOWEST)
	args = append(args, arg)

	// Parse remaining arguments
	for p.curToken.Type == lexer.COMMA {
		p.nextToken() // Skip ','
		arg = p.parseExpression(LOWEST)
		args = append(args, arg)
	}

	if p.curToken.Type != lexer.RPAREN {
		p.errors = append(p.errors, fmt.Sprintf("Expected ')', got %s", p.curToken.Type))
		return nil
	}

	p.nextToken() // Skip ')'
	return &CallExpr{Function: function, Args: args}
}

func (p *Parser) parseIndexExpression(array Node) Node {
	// Skip '['
	p.nextToken()

	index := p.parseExpression(LOWEST)

	if p.curToken.Type != lexer.RBRACKET {
		p.errors = append(p.errors, fmt.Sprintf("Expected ']', got %s", p.curToken.Type))
		return nil
	}

	p.nextToken() // Skip ']'
	return &IndexExpr{Array: array, Index: index}
}

func (p *Parser) parseDotExpression(object Node) Node {
	// Skip '.'
	p.nextToken()

	if p.curToken.Type != lexer.IDENT {
		p.errors = append(p.errors, fmt.Sprintf("Expected identifier after '.', got %s", p.curToken.Type))
		return nil
	}

	property := p.curToken.Literal
	p.nextToken() // Skip property name

	return &DotExpr{Object: object, Property: property}
}