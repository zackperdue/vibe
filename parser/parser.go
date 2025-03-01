package parser

import (
	"fmt"
	"strconv"

	"github.com/example/crystal/lexer"
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

func (p *Parser) parseProgram() *Program {
	program := &Program{Statements: []Node{}}

	for p.curToken.Type != lexer.EOF {
		if p.curToken.Type == lexer.NEWLINE {
			p.nextToken()
			continue
		}

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
	case lexer.KEYWORD:
		switch p.curToken.Value {
		case "def":
			return p.parseFunctionDefinition()
		case "if":
			return p.parseIfStatement()
		case "while":
			return p.parseWhileStatement()
		case "return":
			return p.parseReturnStatement()
		case "puts", "print":
			return p.parsePrintStatement()
		case "type":
			return p.parseTypeDeclaration()
		}
	case lexer.IDENTIFIER:
		// Check if it's a variable declaration with type annotation
		if p.peekToken.Type == lexer.COLON {
			return p.parseVariableDeclaration()
		}

		if p.peekToken.Type == lexer.OPERATOR && p.peekToken.Value == "=" {
			return p.parseAssignment()
		}
		return p.parseExpressionStatement()
	}

	return p.parseExpressionStatement()
}

func (p *Parser) parseFunctionDefinition() Node {
	// Skip 'def' keyword
	p.nextToken()

	if p.curToken.Type != lexer.IDENTIFIER {
		p.errors = append(p.errors, fmt.Sprintf("Expected function name, got %s", p.curToken.Type))
		return nil
	}

	name := p.curToken.Value
	p.nextToken()

	if p.curToken.Type != lexer.LPAREN {
		p.errors = append(p.errors, fmt.Sprintf("Expected '(', got %s", p.curToken.Type))
		return nil
	}
	p.nextToken()

	// Parse parameters
	var parameters []string
	var paramTypes []*TypeAnnotation

	if p.curToken.Type != lexer.RPAREN {
		// Parse parameters with type annotations
		for {
			if p.curToken.Type != lexer.IDENTIFIER {
				p.errors = append(p.errors, fmt.Sprintf("Expected parameter name, got %s", p.curToken.Type))
				break
			}

			paramName := p.curToken.Value
			parameters = append(parameters, paramName)
			p.nextToken()

			// Check for type annotation
			var paramType *TypeAnnotation
			if p.curToken.Type == lexer.COLON {
				p.nextToken() // Skip ':'
				paramType = p.parseTypeAnnotation()
			} else {
				// Default to 'any' type if no annotation
				paramType = &TypeAnnotation{TypeName: "any"}
			}

			paramTypes = append(paramTypes, paramType)

			if p.curToken.Type != lexer.COMMA {
				break
			}

			p.nextToken() // Skip ','
		}
	}

	if p.curToken.Type != lexer.RPAREN {
		p.errors = append(p.errors, fmt.Sprintf("Expected ')', got %s", p.curToken.Type))
		return nil
	}
	p.nextToken()

	// Parse return type
	var returnType *TypeAnnotation
	if p.curToken.Type == lexer.ARROW {
		p.nextToken() // Skip '->'
		returnType = p.parseTypeAnnotation()
	} else {
		// Default return type is 'any'
		returnType = &TypeAnnotation{TypeName: "any"}
	}

	// Function body starts immediately, no need for brackets
	body := p.parseBlockStatement("end")

	return &FunctionDef{
		Name:       name,
		Parameters: parameters,
		ParamTypes: paramTypes,
		ReturnType: returnType,
		Body:       body,
	}
}

func (p *Parser) parseTypeAnnotation() *TypeAnnotation {
	var typeName string

	if p.curToken.Type == lexer.IDENTIFIER || p.curToken.Type == lexer.KEYWORD {
		typeName = p.curToken.Value
		p.nextToken()
	} else {
		p.errors = append(p.errors, fmt.Sprintf("Expected type name, got %s", p.curToken.Type))
		return nil
	}

	// Check for generic type parameters like Array<string>
	var typeParams []Node
	if p.curToken.Type == lexer.OPERATOR && p.curToken.Value == "<" {
		p.nextToken() // Skip '<'

		// Parse the type parameter(s)
		for p.curToken.Type != lexer.OPERATOR || p.curToken.Value != ">" {
			if p.curToken.Type == lexer.EOF {
				p.errors = append(p.errors, "Unexpected EOF while parsing type parameters")
				break
			}

			param := p.parseTypeAnnotation()
			if param != nil {
				typeParams = append(typeParams, param)
			}

			// Skip comma if present
			if p.curToken.Type == lexer.COMMA {
				p.nextToken()
			}
		}

		if p.curToken.Type == lexer.OPERATOR && p.curToken.Value == ">" {
			p.nextToken() // Skip '>'
		}
	}

	// Handle union types with |
	if p.curToken.Type == lexer.PIPE {
		p.nextToken() // Skip '|'
		rightType := p.parseTypeAnnotation()
		if rightType != nil {
			// Create a special union type annotation
			return &TypeAnnotation{
				TypeName: "union",
				TypeParams: []Node{
					&TypeAnnotation{TypeName: typeName, TypeParams: typeParams},
					rightType,
				},
			}
		}
	}

	return &TypeAnnotation{TypeName: typeName, TypeParams: typeParams}
}

func (p *Parser) parseTypeDeclaration() *TypeDeclaration {
	p.nextToken() // Skip 'type'

	if p.curToken.Type != lexer.IDENTIFIER {
		p.errors = append(p.errors, "Expected type name after 'type' keyword")
		return nil
	}

	name := p.curToken.Value
	p.nextToken()

	if p.curToken.Type != lexer.OPERATOR || p.curToken.Value != "=" {
		p.errors = append(p.errors, "Expected '=' after type name")
		return nil
	}

	p.nextToken() // Skip '='

	typeValue := p.parseTypeAnnotation()
	if typeValue == nil {
		return nil
	}

	return &TypeDeclaration{
		Name: name,
		TypeValue: typeValue,
	}
}

func (p *Parser) parseVariableDeclaration() Node {
	var name string
	var typeAnnotation *TypeAnnotation

	// We only support variable declarations starting with the variable name
	if p.curToken.Type == lexer.IDENTIFIER {
		// Get the variable name
		name = p.curToken.Value
		p.nextToken()

		if p.curToken.Type != lexer.COLON {
			p.errors = append(p.errors, "Expected ':' after variable name in declaration")
			return nil
		}

		p.nextToken() // Skip ':'
		typeAnnotation = p.parseTypeAnnotation()
	} else {
		p.errors = append(p.errors, "Expected variable name at the start of variable declaration")
		return nil
	}

	var value Node
	// Check for initialization
	if p.curToken.Type == lexer.OPERATOR && p.curToken.Value == "=" {
		p.nextToken() // Skip '='
		value = p.parseExpression()
	}

	return &VariableDecl{
		Name: name,
		TypeAnnotation: typeAnnotation,
		Value: value,
	}
}

func (p *Parser) parseIfStatement() Node {
	// Skip 'if' keyword
	p.nextToken()

	condition := p.parseExpression()

	// Look for the block start - no need for brackets
	for p.curToken.Type == lexer.NEWLINE {
		p.nextToken()
	}

	consequence := p.parseBlockStatement("end")

	// Check for elsif or else
	var elseIfBlocks []ElseIfBlock
	var alternative *BlockStmt

	for p.peekToken.Type == lexer.KEYWORD && (p.peekToken.Value == "elsif" || p.peekToken.Value == "else") {
		p.nextToken() // Move to elsif/else

		if p.curToken.Value == "elsif" {
			p.nextToken() // Move past elsif
			elsifCondition := p.parseExpression()

			// Look for block start
			for p.curToken.Type == lexer.NEWLINE {
				p.nextToken()
			}

			elsifBlock := p.parseBlockStatement("end")
			elseIfBlocks = append(elseIfBlocks, ElseIfBlock{
				Condition:   elsifCondition,
				Consequence: elsifBlock,
			})
		} else if p.curToken.Value == "else" {
			p.nextToken() // Move past else

			// Look for block start
			for p.curToken.Type == lexer.NEWLINE {
				p.nextToken()
			}

			alternative = p.parseBlockStatement("end")
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

	condition := p.parseExpression()

	// Look for the block start - no need for brackets
	for p.curToken.Type == lexer.NEWLINE {
		p.nextToken()
	}

	body := p.parseBlockStatement("end")

	return &WhileStmt{
		Condition: condition,
		Body:      body,
	}
}

func (p *Parser) parseReturnStatement() Node {
	// Skip 'return' keyword
	p.nextToken()

	// Check if return has no value
	if p.curToken.Type == lexer.NEWLINE {
		return &ReturnStmt{Value: nil}
	}

	value := p.parseExpression()
	return &ReturnStmt{Value: value}
}

func (p *Parser) parsePrintStatement() Node {
	isPuts := p.curToken.Value == "puts"
	// Skip 'puts' or 'print' keyword
	p.nextToken()

	value := p.parseExpression()

	// Create different print nodes based on the type
	if isPuts {
		return &PrintStmt{Value: value}
	} else {
		return &PrintStmt{Value: value}
	}
}

func (p *Parser) parseAssignment() Node {
	name := p.curToken.Value
	p.nextToken() // Move to '='
	p.nextToken() // Move past '='

	value := p.parseExpression()
	return &Assignment{Name: name, Value: value}
}

func (p *Parser) parseExpressionStatement() Node {
	return p.parseExpression()
}

func (p *Parser) parseExpression() Node {
	return p.parseBinaryExpression(0)
}

// Simple operator precedence table
func precedence(op string) int {
	switch op {
	case "||":
		return 1
	case "&&":
		return 2
	case "==", "!=", "<", ">", "<=", ">=":
		return 3
	case "+", "-":
		return 4
	case "*", "/":
		return 5
	default:
		return 0
	}
}

func (p *Parser) parseBinaryExpression(prec int) Node {
	left := p.parseUnaryExpression()

	for p.curToken.Type == lexer.OPERATOR && precedence(p.curToken.Value) > prec {
		operator := p.curToken.Value
		opPrec := precedence(operator)
		p.nextToken()
		right := p.parseBinaryExpression(opPrec)
		left = &BinaryExpr{
			Left:     left,
			Operator: operator,
			Right:    right,
		}
	}

	return left
}

func (p *Parser) parseUnaryExpression() Node {
	switch {
	case p.curToken.Type == lexer.INTEGER:
		val, err := strconv.Atoi(p.curToken.Value)
		if err != nil {
			p.errors = append(p.errors, fmt.Sprintf("Could not parse %q as integer: %s", p.curToken.Value, err))
			return nil
		}
		p.nextToken() // Advance past the integer token
		return &NumberLiteral{Value: float64(val), IsInt: true}
	case p.curToken.Type == lexer.FLOAT:
		val, err := strconv.ParseFloat(p.curToken.Value, 64)
		if err != nil {
			p.errors = append(p.errors, fmt.Sprintf("Could not parse %q as float: %s", p.curToken.Value, err))
			return nil
		}
		p.nextToken() // Advance past the float token
		return &NumberLiteral{Value: val, IsInt: false}
	case p.curToken.Type == lexer.STRING:
		strValue := p.curToken.Value
		p.nextToken() // Advance past the string token
		return &StringLiteral{Value: strValue}
	case p.curToken.Type == lexer.KEYWORD:
		switch p.curToken.Value {
		case "true":
			p.nextToken() // Advance past the 'true' token
			return &BooleanLiteral{Value: true}
		case "false":
			p.nextToken() // Advance past the 'false' token
			return &BooleanLiteral{Value: false}
		case "nil":
			p.nextToken() // Advance past the 'nil' token
			return &NilLiteral{}
		}
		// Fall through for other keywords
		fallthrough
	case p.curToken.Type == lexer.IDENTIFIER:
		identifier := &Identifier{Name: p.curToken.Value}
		p.nextToken() // Advance past the identifier token

		if p.curToken.Type == lexer.LPAREN {
			// Function call
			p.nextToken() // Move past '('

			var args []Node
			if p.curToken.Type != lexer.RPAREN {
				args = p.parseCallArguments()
			}

			if p.curToken.Type != lexer.RPAREN {
				p.errors = append(p.errors, fmt.Sprintf("Expected ')', got %s", p.curToken.Type))
				return nil
			}
			p.nextToken() // Move past ')'

			return &CallExpr{
				Function: identifier,
				Args:     args,
			}
		}

		return identifier
	}

	return nil
}

func (p *Parser) parseCallArguments() []Node {
	var args []Node

	args = append(args, p.parseExpression())
	p.nextToken()

	for p.curToken.Type == lexer.COMMA {
		p.nextToken() // Skip comma
		args = append(args, p.parseExpression())
		p.nextToken()
	}

	return args
}

func (p *Parser) parseBlockStatement(endToken string) *BlockStmt {
	block := &BlockStmt{Statements: []Node{}}

	p.nextToken() // Skip newline or opening token

	for p.curToken.Type != lexer.KEYWORD || p.curToken.Value != endToken {
		if p.curToken.Type == lexer.EOF {
			p.errors = append(p.errors, "Unexpected EOF, expected '" + endToken + "'")
			return block
		}

		if p.curToken.Type == lexer.NEWLINE {
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

// Parse transforms the tokens into an AST
func Parse(l *lexer.Lexer) (*Program, []string) {
	p := New(l)
	program := p.parseProgram()
	return program, p.Errors()
}