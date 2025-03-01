package parser

import (
	"fmt"
	"strconv"
	"strings"

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
	ForStmtNode      NodeType = "ForStmt"
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

	// Class-related node types
	ClassDefNode      NodeType = "ClassDef"      // For class definitions
	MethodDefNode     NodeType = "MethodDef"     // For method definitions
	ClassInstNode     NodeType = "ClassInst"     // For class instantiation (new)
	MethodCallNode    NodeType = "MethodCall"    // For method calls
	SelfExprNode      NodeType = "SelfExpr"      // For self expressions
	SuperCallNode     NodeType = "SuperCall"     // For super calls
	InstanceVarNode   NodeType = "InstanceVar"   // For instance variables (@name)
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
func (i *Identifier) String() string { return i.Name }

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

// TypeAnnotation represents a type annotation
type TypeAnnotation struct {
	TypeName    string
	GenericType *TypeAnnotation
	TypeParams  []Node // For generic types like Array<string>
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

// Parameter represents a function parameter with its type
type Parameter struct {
	Name string
	Type *TypeAnnotation
}

func (p *Parameter) String() string {
	if p.Type == nil {
		return p.Name
	}
	return p.Name + ": " + p.Type.String()
}

// FunctionDef represents a function definition
type FunctionDef struct {
	Name       string
	Parameters []Parameter
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
		result += param.String()
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
	if a.Value == nil {
		return fmt.Sprintf("Assignment(%s = nil)", a.Name)
	}
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
func (p *PrintStmt) String() string {
	if p.Value == nil {
		return "PrintStmt(nil)"
	}
	return fmt.Sprintf("PrintStmt(%s)", p.Value.String())
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

// ForStmt represents a for loop with iterator
type ForStmt struct {
	Iterator  string     // The variable that will hold each element
	Iterable  Node       // The expression to iterate over
	Body      *BlockStmt
}

func (f *ForStmt) Type() NodeType { return ForStmtNode }
func (f *ForStmt) String() string {
	iterableStr := "<nil>"
	if f.Iterable != nil {
		iterableStr = f.Iterable.String()
	}

	bodyStr := "<nil>"
	if f.Body != nil {
		bodyStr = f.Body.String()
	}

	return fmt.Sprintf("ForStmt(%s in %s, %s)", f.Iterator, iterableStr, bodyStr)
}

// MethodCall represents a method call expression
type MethodCall struct {
	Object Node   // The object on which the method is called
	Method string // The name of the method
	Args   []Node // Arguments passed to the method
}

// Type returns the type of the node
func (m *MethodCall) Type() NodeType {
	return MethodCallNode
}

// String returns a string representation of the method call
func (m *MethodCall) String() string {
	var args []string
	for _, arg := range m.Args {
		args = append(args, arg.String())
	}

	return fmt.Sprintf("%s.%s(%s)", m.Object.String(), m.Method, strings.Join(args, ", "))
}

// SelfExpr represents a 'self' expression in a method
type SelfExpr struct{}

// Type returns the type of the node
func (s *SelfExpr) Type() NodeType {
	return SelfExprNode
}

// String returns a string representation of the self expression
func (s *SelfExpr) String() string {
	return "self"
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

	// Debug print
	fmt.Printf("DEBUG: Parsed %d statements\n", len(program.Statements))
	for i, stmt := range program.Statements {
		if stmt != nil {
			fmt.Printf("DEBUG: Statement %d: %T - %s\n", i, stmt, stmt.String())
		} else {
			fmt.Printf("DEBUG: Statement %d: nil\n", i)
		}
	}

	return program, p.errors
}

func (p *Parser) parseProgram() *Program {
	fmt.Println("DEBUG: Starting to parse program")
	program := &Program{}
	program.Statements = []Node{}

	for p.curToken.Type != lexer.EOF {
		fmt.Printf("DEBUG: parseProgram - current token: %s, literal: %s, peek token: %s, literal: %s\n",
			p.curToken.Type, p.curToken.Literal, p.peekToken.Type, p.peekToken.Literal)

		// Handle class inheritance pattern
		if p.curToken.Type == lexer.IDENT && p.peekToken.Type == lexer.INHERITS {
			fmt.Println("DEBUG: parseProgram - Detected class inheritance pattern, handling it specially")
			className := p.curToken.Literal
			p.nextToken() // consume IDENT
			p.nextToken() // consume INHERITS
			parentClass := p.curToken.Literal
			p.nextToken() // consume parent class name

			// Parse the class definition
			classDef := p.parseClassDefinition(className, parentClass)
			if classDef != nil {
				program.Statements = append(program.Statements, classDef)
				fmt.Printf("DEBUG: parseProgram - Successfully added CLASS definition: %s\n", classDef.String())
			}
			continue
		}

		// Handle regular class definitions
		if p.curToken.Type == lexer.CLASS {
			fmt.Println("DEBUG: parseProgram - Detected class definition")
			classDef := p.parseClassDefinition("", "")
			if classDef != nil {
				program.Statements = append(program.Statements, classDef)
				fmt.Printf("DEBUG: parseProgram - Successfully added CLASS definition: %s\n", classDef.String())
			}

			fmt.Printf("DEBUG: parseProgram - After class definition, current token: %s, peek token: %s\n",
				p.curToken.Type, p.peekToken.Type)

			continue
		}

		// Handle assignment pattern
		if p.curToken.Type == lexer.IDENT {
			fmt.Printf("DEBUG: parseProgram - Found identifier: %s, peek token: %s\n", p.curToken.Literal, p.peekToken.Type)

			if p.peekToken.Type == lexer.ASSIGN {
				fmt.Printf("DEBUG: parseProgram - Detected assignment to variable: %s\n", p.curToken.Literal)

				// Save the variable name
				name := p.curToken.Literal

				// Skip to '='
				p.nextToken()
				fmt.Printf("DEBUG: parseProgram - Now at =, current token: %s\n", p.curToken.Type)

				// Skip '='
				p.nextToken()
				fmt.Printf("DEBUG: parseProgram - Now after =, current token: %s, literal: %s\n", p.curToken.Type, p.curToken.Literal)

				// Parse the right side of the assignment
				value := p.parseExpression(LOWEST)
				if value == nil {
					fmt.Println("DEBUG: parseProgram - Failed to parse right side of assignment, skipping to next statement")
					// Skip to the next statement
					for p.curToken.Type != lexer.SEMICOLON && p.curToken.Type != lexer.EOF {
						p.nextToken()
					}
					continue
				}

				fmt.Printf("DEBUG: parseProgram - Parsed right side of assignment: %T - %s\n", value, value.String())

				// Create and add the assignment node
				assignment := &Assignment{
					Name:  name,
					Value: value,
				}

				program.Statements = append(program.Statements, assignment)
				fmt.Printf("DEBUG: parseProgram - added assignment: %s\n", assignment.String())

				continue
			}
		}

		stmt := p.parseStatement()
		if stmt != nil {
			program.Statements = append(program.Statements, stmt)
			fmt.Printf("DEBUG: parseProgram - added statement: %T - %s\n", stmt, stmt.String())
		} else {
			fmt.Printf("DEBUG: parseProgram - statement was nil, skipping token: %s\n", p.curToken.Type)
		}

		// Only advance to the next token if we haven't reached EOF
		if p.curToken.Type != lexer.EOF {
			p.nextToken()
		}
	}

	fmt.Printf("DEBUG: Parsed %d statements\n", len(program.Statements))
	for i, stmt := range program.Statements {
		fmt.Printf("DEBUG: Statement %d: %T - %s\n", i, stmt, stmt.String())
	}

	return program
}

func (p *Parser) parseStatement() Node {
	fmt.Printf("DEBUG: parseStatement - current token: %s, literal: %s\n", p.curToken.Type, p.curToken.Literal)

	switch p.curToken.Type {
	case lexer.IDENT:
		// Check if this is an assignment
		if p.peekToken.Type == lexer.ASSIGN {
			return p.parseAssignment()
		}
		return p.parseExpressionStatement()
	case lexer.ASSIGN:
		// If we encounter an assignment operator directly, we need to skip it
		// This can happen when parsing multiple assignments in sequence
		return nil
	case lexer.RETURN:
		return p.parseReturnStatement()
	case lexer.PRINT:
		return p.parsePrintStatement()
	case lexer.IF:
		return p.parseIfStatement()
	case lexer.FUNCTION:
		return p.parseFunctionDefinition()
	case lexer.FOR:
		fmt.Println("DEBUG: Detected FOR token in parseStatement, calling parseForStatement")
		return p.parseForStatement()
	case lexer.WHILE:
		return p.parseWhileStatement()
	case lexer.CLASS:
		return p.parseClassDefinition("", "")
	case lexer.SUPER:
		return p.parseSuperCall()
	case lexer.IN, lexer.DO, lexer.END:
		// These tokens are part of control structures and should be handled by their respective parsers
		fmt.Printf("DEBUG: Skipping token %s as it should be handled by its control structure parser\n", p.curToken.Type)
		return nil
	case lexer.AT:
		// Handle @ symbol (instance variables)
		return p.parseInstanceVariable()
	case lexer.ILLEGAL:
		// Special handling for any illegal tokens
		return nil
	case lexer.INT, lexer.FLOAT, lexer.STRING, lexer.TRUE, lexer.FALSE, lexer.NIL,
		lexer.LPAREN, lexer.LBRACKET, lexer.LBRACE, lexer.MINUS, lexer.BANG:
		return p.parseExpressionStatement()
	default:
		return nil
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

	funcDef.Parameters = p.parseFunctionParameters()

	// Check for return type annotation with : syntax
	if p.curToken.Type == lexer.COLON {
		p.nextToken()
		funcDef.ReturnType = p.parseTypeAnnotation()
	} else {
		// Default return type is "int"
		funcDef.ReturnType = &TypeAnnotation{TypeName: "int"}
	}

	// Check for 'do' keyword
	if p.curToken.Type != lexer.DO {
		p.errors = append(p.errors, fmt.Sprintf("Expected 'do' after function definition, got %s", p.curToken.Type))
	} else {
		p.nextToken() // Skip 'do'
	}

	// Parse function body
	funcDef.Body = &BlockStmt{Statements: []Node{}}

	// Parse statements until we see 'end' or EOF
	for p.curToken.Type != lexer.END && p.curToken.Type != lexer.EOF {
		if p.curToken.Type == lexer.SEMICOLON {
			p.nextToken()
			continue
		}

		stmt := p.parseStatement()
		if stmt != nil {
			funcDef.Body.Statements = append(funcDef.Body.Statements, stmt)
		}
		p.nextToken()
	}

	// Check that we found the 'end' keyword
	if p.curToken.Type != lexer.END {
		p.errors = append(p.errors, "Expected 'end' to close function body")
	} else {
		p.nextToken() // Skip the 'end'
	}

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
	ifStmt := &IfStmt{}

	// Skip 'if' keyword
	p.nextToken()

	// Parse condition
	ifStmt.Condition = p.parseExpression(LOWEST)

	// No opening brace for if statements anymore
	// Parse the consequence block directly
	ifStmt.Consequence = &BlockStmt{Statements: []Node{}}

	// Parse statements until we see 'else', 'elsif', 'end', or EOF
	for p.peekToken.Type != lexer.ELSE && p.peekToken.Type != lexer.ELSIF && p.peekToken.Type != lexer.END && p.peekToken.Type != lexer.EOF {
		p.nextToken()

		if p.curToken.Type == lexer.SEMICOLON {
			continue
		}

		stmt := p.parseStatement()
		if stmt != nil {
			ifStmt.Consequence.Statements = append(ifStmt.Consequence.Statements, stmt)
		}
	}

	// Check for 'else' or 'elsif'
	if p.peekToken.Type == lexer.ELSE || p.peekToken.Type == lexer.ELSIF {
		p.nextToken() // Move to 'else' or 'elsif'

		// Check if it's 'elsif' or just 'else'
		if p.curToken.Type == lexer.ELSIF {
			// This is an 'elsif'
			elseIfBlock := ElseIfBlock{
				Condition:   nil,
				Consequence: nil,
			}

			// Parse the condition
			p.nextToken()
			elseIfBlock.Condition = p.parseExpression(LOWEST)

			// Parse the consequence statements
			elseIfBlock.Consequence = &BlockStmt{Statements: []Node{}}

			// Parse statements until we see 'else', 'elsif', 'end', or EOF
			for p.peekToken.Type != lexer.ELSE && p.peekToken.Type != lexer.ELSIF && p.peekToken.Type != lexer.END && p.peekToken.Type != lexer.EOF {
				p.nextToken()

				if p.curToken.Type == lexer.SEMICOLON {
					continue
				}

				stmt := p.parseStatement()
				if stmt != nil {
					elseIfBlock.Consequence.Statements = append(elseIfBlock.Consequence.Statements, stmt)
				}
			}

			ifStmt.ElseIfBlocks = append(ifStmt.ElseIfBlocks, elseIfBlock)

			// Recursively parse any additional 'elsif' or 'else' blocks
			if p.peekToken.Type == lexer.ELSE || p.peekToken.Type == lexer.ELSIF {
				// Create a temporary if statement to parse the remainder
				tempIf, ok := p.parseIfStatement().(*IfStmt)
				if ok {
					// Transfer any elsif blocks
					ifStmt.ElseIfBlocks = append(ifStmt.ElseIfBlocks, tempIf.ElseIfBlocks...)

					// Transfer the else block if there is one
					ifStmt.Alternative = tempIf.Alternative
				}

				// Skip past the 'end' token since it was consumed by the recursive call
				return ifStmt
			}
		} else if p.curToken.Type == lexer.ELSE {
			// This is just an 'else'
			// Skip 'else' keyword
			p.nextToken()

			// Parse the alternative block
			ifStmt.Alternative = &BlockStmt{Statements: []Node{}}

			// Parse statements until we see 'end' or EOF
			for p.peekToken.Type != lexer.END && p.peekToken.Type != lexer.EOF {
				p.nextToken()

				if p.curToken.Type == lexer.SEMICOLON {
					continue
				}

				stmt := p.parseStatement()
				if stmt != nil {
					ifStmt.Alternative.Statements = append(ifStmt.Alternative.Statements, stmt)
				}
			}
		}
	}

	// Consume the 'end' token
	if p.peekToken.Type == lexer.END {
		p.nextToken() // Move to 'end'
		p.nextToken() // Skip 'end'
	} else {
		p.errors = append(p.errors, "Expected 'end' to close if statement")
	}

	return ifStmt
}

func (p *Parser) parseWhileStatement() Node {
	// Skip 'while' keyword
	p.nextToken()

	condition := p.parseExpression(LOWEST)

	// Check for 'do' keyword
	if p.curToken.Type != lexer.DO {
		p.errors = append(p.errors, fmt.Sprintf("Expected 'do' after while condition, got %s", p.curToken.Type))
	} else {
		p.nextToken() // Skip 'do'
	}

	// Parse while loop body directly
	body := &BlockStmt{Statements: []Node{}}

	// Parse statements until we see 'end' or EOF
	for p.peekToken.Type != lexer.END && p.peekToken.Type != lexer.EOF {
		p.nextToken()

		if p.curToken.Type == lexer.SEMICOLON {
			continue
		}

		stmt := p.parseStatement()
		if stmt != nil {
			body.Statements = append(body.Statements, stmt)
		}
	}

	// Consume the 'end' token
	if p.peekToken.Type == lexer.END {
		p.nextToken() // Move to 'end'
		p.nextToken() // Skip 'end'
	} else {
		p.errors = append(p.errors, "Expected 'end' to close while loop")
	}

	return &WhileStmt{
		Condition: condition,
		Body:      body,
	}
}

func (p *Parser) parseReturnStatement() Node {
	// Skip 'return' keyword
	p.nextToken()

	// Check if return has no value
	if p.curToken.Type == lexer.SEMICOLON || p.curToken.Type == lexer.EOF {
		return &ReturnStmt{Value: nil}
	}

	value := p.parseExpression(LOWEST)

	// Create a ReturnStmt node
	return &ReturnStmt{Value: value}
}

func (p *Parser) parsePrintStatement() Node {
	stmt := &PrintStmt{}

	// Skip 'print' or 'puts' keyword
	p.nextToken()

	// Check if it's the print(expr) syntax with parentheses
	if p.curToken.Type == lexer.LPAREN {
		// Skip '('
		p.nextToken()

		// Parse the expression to print
		stmt.Value = p.parseExpression(LOWEST)

		// Skip to check for the closing paren
		if p.peekToken.Type == lexer.RPAREN {
			p.nextToken()
		}

		// Ensure we have a closing parenthesis
		if p.curToken.Type != lexer.RPAREN {
			p.errors = append(p.errors, fmt.Sprintf("Expected ')' after print value, got %s", p.curToken.Type))
			return nil
		}

		// Skip ')'
		p.nextToken()
	} else {
		// It's the puts expr syntax without parentheses
		// Parse the expression to print
		stmt.Value = p.parseExpression(LOWEST)
	}

	return stmt
}

func (p *Parser) parseAssignment() Node {
	debugf("parseAssignment - at token: %s", p.curToken.Type)

	// Save the variable name
	name := p.curToken.Literal

	// Skip to '='
	p.nextToken()

	// Skip '='
	p.nextToken()

	// Parse the right side of the assignment
	value := p.parseExpression(LOWEST)
	if value == nil {
		fmt.Println("DEBUG: parseAssignment - Failed to parse right side of assignment")
		return nil
	}

	// Create and return the assignment node
	assignment := &Assignment{
		Name:  name,
		Value: value,
	}

	debugf("parseAssignment - created assignment: %s = %s", name, value.String())
	return assignment
}

func (p *Parser) parseExpressionStatement() Node {
	return p.parseExpression(0)
}

func (p *Parser) parseExpression(precedence int) Node {
	fmt.Printf("DEBUG: parseExpression - at token: %s, literal: %s\n", p.curToken.Type, p.curToken.Literal)

	// Check for instance variables (@name)
	if p.curToken.Type == lexer.AT {
		return p.parseInstanceVariable()
	}

	// Check for self keyword
	if p.curToken.Type == lexer.SELF {
		return p.parseSelfExpr()
	}

	// Check for super call
	if p.curToken.Type == lexer.SUPER {
		return p.parseSuperCall()
	}

	// Continue with the existing prefix/infix expression parsing
	var leftExp Node

	// Prefix expressions
	switch p.curToken.Type {
	case lexer.IDENT:
		leftExp = &Identifier{Name: p.curToken.Literal}

		// Check for generic type parameter like Box<Int>
		if p.peekToken.Type == lexer.LT {
			ident := p.curToken.Literal
			p.nextToken() // Skip to '<'

			// Now we're at '<'
			p.nextToken() // Skip to the type parameter

			// Parse the type parameter
			if p.curToken.Type != lexer.IDENT {
				p.errors = append(p.errors, fmt.Sprintf("Expected type parameter, got %s", p.curToken.Type))
				return nil
			}

			typeParam := &Identifier{Name: p.curToken.Literal}

			// Create a binary expression to represent the generic type
			leftExp = &BinaryExpr{
				Left:     &Identifier{Name: ident},
				Operator: "<",
				Right:    typeParam,
			}

			// Skip to '>'
			p.nextToken()
			if p.curToken.Type != lexer.GT {
				p.errors = append(p.errors, fmt.Sprintf("Expected '>' after type parameter, got %s", p.curToken.Type))
				return nil
			}
		}

	case lexer.INT:
		value, err := strconv.ParseFloat(p.curToken.Literal, 64)
		if err != nil {
			p.errors = append(p.errors, fmt.Sprintf("Could not parse %q as integer", p.curToken.Literal))
			return nil
		}
		leftExp = &NumberLiteral{Value: value, IsInt: true}
	case lexer.FLOAT:
		value, err := strconv.ParseFloat(p.curToken.Literal, 64)
		if err != nil {
			p.errors = append(p.errors, fmt.Sprintf("Could not parse %q as float", p.curToken.Literal))
			return nil
		}
		leftExp = &NumberLiteral{Value: value, IsInt: false}
	case lexer.STRING:
		leftExp = &StringLiteral{Value: p.curToken.Literal}
	case lexer.TRUE:
		leftExp = &BooleanLiteral{Value: true}
	case lexer.FALSE:
		leftExp = &BooleanLiteral{Value: false}
	case lexer.NIL:
		leftExp = &NilLiteral{}
	case lexer.LPAREN:
		p.nextToken() // Consume '('
		leftExp = p.parseExpression(LOWEST)

		if p.peekToken.Type != lexer.RPAREN {
			p.errors = append(p.errors, fmt.Sprintf("Expected ')', got %s", p.peekToken.Type))
			return nil
		}
		p.nextToken() // Consume ')'
	case lexer.LBRACKET:
		leftExp = p.parseArrayLiteral()
	case lexer.MINUS, lexer.BANG:
		operator := p.curToken.Literal
		p.nextToken() // Consume the operator
		operand := p.parseExpression(PREFIX)
		leftExp = &UnaryExpr{Operator: operator, Right: operand}
	default:
		return nil
	}

	// Skip to the next token
	p.nextToken()

	// Now parse any infix expressions
	for precedence < p.curPrecedence() && p.curToken.Type != lexer.EOF {
		switch p.curToken.Type {
		case lexer.PLUS, lexer.MINUS, lexer.ASTERISK, lexer.SLASH, lexer.MODULO,
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
	case lexer.PLUS, lexer.MINUS, lexer.ASTERISK, lexer.SLASH, lexer.MODULO,
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
	case lexer.ASTERISK, lexer.SLASH, lexer.MODULO:
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
	case lexer.ASTERISK, lexer.SLASH, lexer.MODULO:
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

	// First, check if we're at the opening brace
	if p.curToken.Type == lexer.LBRACE {
		p.nextToken() // Skip '{'
	} else {
		// If we're not at an opening brace, don't try to parse a block
		p.errors = append(p.errors, fmt.Sprintf("Expected '{' to start block, got %s", p.curToken.Type))
		return block
	}

	for p.curToken.Type != lexer.RBRACE && p.curToken.Type != lexer.EOF {
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

	if p.curToken.Type != lexer.RBRACE {
		p.errors = append(p.errors, "Unexpected EOF, expected '}'")
	}

	return block
}

// parseFunctionParameters parses function parameters with optional type annotations
func (p *Parser) parseFunctionParameters() []Parameter {
	var parameters []Parameter

	p.nextToken() // Skip '('

	// Handle empty parameter list
	if p.curToken.Type == lexer.RPAREN {
		p.nextToken() // Skip ')'
		return parameters
	}

	// Parse first parameter
	if p.curToken.Type != lexer.IDENT {
		p.errors = append(p.errors, fmt.Sprintf("Expected parameter name, got %s", p.curToken.Type))
		// Try to recover
		for p.curToken.Type != lexer.RPAREN && p.curToken.Type != lexer.EOF {
			p.nextToken()
		}
		if p.curToken.Type == lexer.RPAREN {
			p.nextToken() // Skip ')'
		}
		return parameters
	}

	// Save the parameter name
	paramName := p.curToken.Literal
	param := Parameter{Name: paramName}

	// Check for type annotation
	p.nextToken()
	if p.curToken.Type == lexer.COLON {
		p.nextToken() // Skip ':'
		// Parse type
		if p.curToken.Type != lexer.IDENT {
			p.errors = append(p.errors, fmt.Sprintf("Expected type name, got %s", p.curToken.Type))
			param.Type = &TypeAnnotation{TypeName: "any"}
		} else {
			param.Type = &TypeAnnotation{TypeName: p.curToken.Literal}
			p.nextToken() // Move past type name
		}
	} else {
		// No type annotation, add a default
		param.Type = &TypeAnnotation{TypeName: "any"}
	}

	parameters = append(parameters, param)

	// Parse additional parameters
	for p.curToken.Type == lexer.COMMA {
		p.nextToken() // Skip ','

		// Parse parameter name
		if p.curToken.Type != lexer.IDENT {
			p.errors = append(p.errors, fmt.Sprintf("Expected parameter name after comma, got %s", p.curToken.Type))
			// Try to recover
			for p.curToken.Type != lexer.RPAREN && p.curToken.Type != lexer.EOF {
				p.nextToken()
			}
			if p.curToken.Type == lexer.RPAREN {
				p.nextToken() // Skip ')'
			}
			return parameters
		}

		// Save the parameter name
		paramName = p.curToken.Literal
		param = Parameter{Name: paramName}

		// Check for type annotation
		p.nextToken()
		if p.curToken.Type == lexer.COLON {
			p.nextToken() // Skip ':'
			// Parse type
			if p.curToken.Type != lexer.IDENT {
				p.errors = append(p.errors, fmt.Sprintf("Expected type name, got %s", p.curToken.Type))
				param.Type = &TypeAnnotation{TypeName: "any"}
			} else {
				param.Type = &TypeAnnotation{TypeName: p.curToken.Literal}
				p.nextToken() // Move past type name
			}
		} else {
			// No type annotation, add a default
			param.Type = &TypeAnnotation{TypeName: "any"}
		}

		parameters = append(parameters, param)
	}

	// Make sure we've reached the closing parenthesis
	if p.curToken.Type != lexer.RPAREN {
		p.errors = append(p.errors, fmt.Sprintf("Expected ')' after parameters, got %s", p.curToken.Type))
		// Try to recover
		for p.curToken.Type != lexer.RPAREN && p.curToken.Type != lexer.EOF {
			p.nextToken()
		}
	}

	// Skip the closing parenthesis
	if p.curToken.Type == lexer.RPAREN {
		p.nextToken() // Skip ')'
	}

	return parameters
}

func (p *Parser) parseArrayLiteral() Node {
	arrayLit := &ArrayLiteral{Elements: []Node{}}

	// We're already at '[', skip to the first element
	p.nextToken()

	// Empty array case
	if p.curToken.Type == lexer.RBRACKET {
		p.nextToken() // Skip ']'
		return arrayLit
	}

	// Parse first element
	firstElement := p.parseExpression(LOWEST)
	if firstElement != nil {
		arrayLit.Elements = append(arrayLit.Elements, firstElement)
	}

	// Parse remaining elements
	for p.peekToken.Type == lexer.COMMA {
		p.nextToken() // Move to the comma
		p.nextToken() // Move past the comma

		// Handle trailing comma
		if p.curToken.Type == lexer.RBRACKET {
			break
		}

		element := p.parseExpression(LOWEST)
		if element != nil {
			arrayLit.Elements = append(arrayLit.Elements, element)
		}
	}

	// Check for closing bracket
	if p.peekToken.Type != lexer.RBRACKET {
		p.errors = append(p.errors, fmt.Sprintf("Expected ']', got %s", p.peekToken.Type))
		return nil
	}

	p.nextToken() // Move to ']'
	p.nextToken() // Skip the closing bracket

	return arrayLit
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

func (p *Parser) parseDotExpression(left Node) Node {
	debugf("parseDotExpression - at token: %s, left: %s", p.curToken.Type, left.String())

	// Skip the '.' token
	p.nextToken()

	// Next token should be the method name or 'new'
	if p.curToken.Type != lexer.IDENT && p.curToken.Type != lexer.NEW {
		p.addError(fmt.Sprintf("Expected method name or 'new' after '.', got %s", p.curToken.Type))
		return nil
	}

	// If it's 'new', parse it as a class instantiation
	if p.curToken.Type == lexer.NEW {
		return p.parseClassInstantiation(left)
	}

	// Otherwise it's a method call
	methodCall := &MethodCall{
		Object: left,
		Method: p.curToken.Literal,
		Args:   []Node{},
	}

	// Skip method name
	p.nextToken()

	// Check for opening parenthesis
	if p.curToken.Type != lexer.LPAREN {
		p.addError(fmt.Sprintf("Expected '(' after method name, got %s", p.curToken.Type))
		return nil
	}

	// Skip '('
	p.nextToken()

	// Parse arguments if any
	if p.curToken.Type != lexer.RPAREN {
		// Parse the first argument
		arg := p.parseExpression(LOWEST)
		if arg != nil {
			methodCall.Args = append(methodCall.Args, arg)
		}

		// Parse additional arguments
		for p.peekToken.Type == lexer.COMMA {
			p.nextToken() // Skip the comma
			p.nextToken() // Move to the next argument

			arg = p.parseExpression(LOWEST)
			if arg != nil {
				methodCall.Args = append(methodCall.Args, arg)
			}
		}
	}

	// Check for closing parenthesis
	if p.curToken.Type != lexer.RPAREN {
		p.addError(fmt.Sprintf("Expected ')' to close arguments, got %s", p.curToken.Type))
		return nil
	}

	// Skip ')'
	p.nextToken()

	return methodCall
}

// parseClassInstantiation parses a class instantiation (ClassName.new(...))
func (p *Parser) parseClassInstantiation(left Node) Node {
	fmt.Printf("DEBUG: parseClassInstantiation - at token: %s, literal: %s\n", p.curToken.Type, p.curToken.Literal)

	// Create a ClassInst node
	classInst := &ClassInst{
		Class:    left,
		Args:     []Node{},
		TypeArgs: []Node{},
	}

	// Skip the 'new' token
	p.nextToken()
	fmt.Printf("DEBUG: parseClassInstantiation - after 'new', token: %s, literal: %s\n", p.curToken.Type, p.curToken.Literal)

	// Check for opening parenthesis
	if p.curToken.Type != lexer.LPAREN {
		p.errors = append(p.errors, fmt.Sprintf("Expected '(' after 'new', got %s", p.curToken.Type))
		return nil
	}

	// Skip '('
	p.nextToken()
	fmt.Printf("DEBUG: parseClassInstantiation - after '(', token: %s, literal: %s\n", p.curToken.Type, p.curToken.Literal)

	// Parse arguments if any
	if p.curToken.Type != lexer.RPAREN {
		// Parse the first argument
		arg := p.parseExpression(LOWEST)
		if arg != nil {
			classInst.Args = append(classInst.Args, arg)
			fmt.Printf("DEBUG: parseClassInstantiation - added first argument: %s\n", arg.String())
		}

		// Parse additional arguments
		for p.curToken.Type == lexer.COMMA {
			fmt.Printf("DEBUG: parseClassInstantiation - found comma, token: %s\n", p.curToken.Type)

			// Skip the comma
			p.nextToken()

			fmt.Printf("DEBUG: parseClassInstantiation - after comma, token: %s, literal: %s\n", p.curToken.Type, p.curToken.Literal)

			// Parse the next argument
			arg := p.parseExpression(LOWEST)
			if arg != nil {
				classInst.Args = append(classInst.Args, arg)
				fmt.Printf("DEBUG: parseClassInstantiation - added additional argument: %s\n", arg.String())
			}
		}
	}

	// Check for closing parenthesis
	if p.curToken.Type != lexer.RPAREN {
		p.errors = append(p.errors, fmt.Sprintf("Expected ')' to close arguments, got %s", p.curToken.Type))
		return nil
	}

	// Skip ')'
	p.nextToken()
	fmt.Printf("DEBUG: parseClassInstantiation - after ')', token: %s, literal: %s\n", p.curToken.Type, p.curToken.Literal)

	fmt.Printf("DEBUG: parseClassInstantiation - created class instantiation: %s\n", classInst.String())
	return classInst
}

// Helper method to parse an array element
func (p *Parser) parseArrayElement() Node {
	// Handle different element types
	switch p.curToken.Type {
	case lexer.INT:
		value, err := strconv.ParseFloat(p.curToken.Literal, 64)
		if err != nil {
			p.errors = append(p.errors, fmt.Sprintf("Could not parse %q as integer", p.curToken.Literal))
			return nil
		}
		p.nextToken() // Move past the number
		return &NumberLiteral{Value: value, IsInt: true}

	case lexer.FLOAT:
		value, err := strconv.ParseFloat(p.curToken.Literal, 64)
		if err != nil {
			p.errors = append(p.errors, fmt.Sprintf("Could not parse %q as float", p.curToken.Literal))
			return nil
		}
		p.nextToken() // Move past the number
		return &NumberLiteral{Value: value, IsInt: false}

	case lexer.STRING:
		literal := p.curToken.Literal
		p.nextToken() // Move past the string
		return &StringLiteral{Value: literal}

	case lexer.IDENT:
		name := p.curToken.Literal
		p.nextToken() // Move past the identifier
		return &Identifier{Name: name}

	case lexer.TRUE:
		p.nextToken() // Move past 'true'
		return &BooleanLiteral{Value: true}

	case lexer.FALSE:
		p.nextToken() // Move past 'false'
		return &BooleanLiteral{Value: false}

	case lexer.NIL:
		p.nextToken() // Move past 'nil'
		return &NilLiteral{}

	default:
		p.errors = append(p.errors, fmt.Sprintf("Unexpected token in array: %s", p.curToken.Type))
		return nil
	}
}

func debugf(format string, args ...interface{}) {
	fmt.Printf("DEBUG: "+format+"\n", args...)
}

func (p *Parser) addError(msg string) {
	p.errors = append(p.errors, msg)
}

func (p *Parser) parseSelfExpr() Node {
	debugf("parseSelfExpr - at token: %s", p.curToken.Type)

	if p.curToken.Type != lexer.SELF {
		p.addError(fmt.Sprintf("Expected 'self' keyword, got %s", p.curToken.Type))
		return nil
	}

	return &SelfExpr{}
}

// ClassInst represents a class instantiation expression
type ClassInst struct {
	Class    Node   // The class being instantiated
	Args     []Node // Arguments passed to the constructor
	TypeArgs []Node // Type arguments for generic classes
}

// Type returns the type of the node
func (c *ClassInst) Type() NodeType {
	return ClassInstNode
}

// String returns a string representation of the class instantiation
func (c *ClassInst) String() string {
	var args []string
	for _, arg := range c.Args {
		args = append(args, arg.String())
	}

	return fmt.Sprintf("ClassInst(%s, args[%s])", c.Class.String(), strings.Join(args, ", "))
}

// ClassDef represents a class definition
type ClassDef struct {
	Name       string            // The name of the class
	Parent     string            // The parent class (if any)
	Methods    []Node            // Methods defined in the class
	Fields     []struct {        // Fields defined in the class
		Name          string
		TypeAnnotation struct {
			TypeName string
		}
	}
	TypeParams []string          // Type parameters for generic classes
}

// Type returns the type of the node
func (c *ClassDef) Type() NodeType {
	return ClassDefNode
}

// String returns a string representation of the class definition
func (c *ClassDef) String() string {
	var methods []string
	for _, method := range c.Methods {
		methods = append(methods, method.String())
	}

	var fields []string
	for _, field := range c.Fields {
		fields = append(fields, field.Name)
	}

	parent := c.Parent
	if parent == "" {
		parent = "none"
	}

	return fmt.Sprintf("ClassDef(%s, parent=%s, methods=%s, fields=%s)",
		c.Name, parent, strings.Join(methods, ", "), strings.Join(fields, ", "))
}

// parseForStatement parses a for loop statement
func (p *Parser) parseForStatement() Node {
	fmt.Printf("DEBUG: parseForStatement - at token: %s, literal: %s\n", p.curToken.Type, p.curToken.Literal)

	// Skip 'for' token
	p.nextToken()

	// Parse iterator variable
	if p.curToken.Type != lexer.IDENT {
		p.errors = append(p.errors, fmt.Sprintf("Expected iterator variable name, got %s", p.curToken.Type))
		return nil
	}

	iterator := p.curToken.Literal
	p.nextToken()

	// Check for 'in' keyword
	if p.curToken.Type != lexer.IN {
		p.errors = append(p.errors, fmt.Sprintf("Expected 'in' after iterator variable, got %s", p.curToken.Type))
		return nil
	}

	// Skip 'in' token
	p.nextToken()

	// Parse iterable expression
	iterable := p.parseExpression(LOWEST)

	// Check for 'do' keyword
	if p.curToken.Type != lexer.DO {
		p.errors = append(p.errors, fmt.Sprintf("Expected 'do' after iterable expression, got %s", p.curToken.Type))
	} else {
		p.nextToken() // Skip 'do'
	}

	// Parse for loop body
	body := &BlockStmt{Statements: []Node{}}

	// Parse statements until we see 'end' or EOF
	for p.curToken.Type != lexer.END && p.curToken.Type != lexer.EOF {
		if p.curToken.Type == lexer.SEMICOLON {
			p.nextToken()
			continue
		}

		stmt := p.parseStatement()
		if stmt != nil {
			body.Statements = append(body.Statements, stmt)
		}
		p.nextToken()
	}

	// Skip 'end' token
	if p.curToken.Type == lexer.END {
		p.nextToken()
	}

	return &ForStmt{
		Iterator: iterator,
		Iterable: iterable,
		Body:     body,
	}
}

// parseInstanceVariable parses an instance variable (@name)
func (p *Parser) parseInstanceVariable() Node {
	fmt.Printf("DEBUG: parseInstanceVariable - at token: %s, literal: %s\n", p.curToken.Type, p.curToken.Literal)

	// Skip '@' token
	if p.curToken.Type == lexer.AT {
		p.nextToken()
	}

	// Parse variable name
	if p.curToken.Type != lexer.IDENT {
		p.errors = append(p.errors, fmt.Sprintf("Expected instance variable name after @, got %s", p.curToken.Type))
		return nil
	}

	name := p.curToken.Literal
	p.nextToken()

	return &Identifier{Name: "@" + name}
}

// parseSuperCall parses a super call (super.method(...) or super(...))
func (p *Parser) parseSuperCall() Node {
	fmt.Printf("DEBUG: parseSuperCall - at token: %s, literal: %s\n", p.curToken.Type, p.curToken.Literal)

	// Skip 'super' token
	p.nextToken()

	// Check if this is a direct super call (super(...)) or a method call (super.method(...))
	if p.curToken.Type == lexer.LPAREN {
		// Direct super call
		// Skip '(' token
		p.nextToken()

		// Create a method call node for the constructor
		methodCall := &MethodCall{
			Object: &Identifier{Name: "super"},
			Method: "initialize", // Implicit constructor call
			Args:   []Node{},
		}

		// Parse arguments if any
		if p.curToken.Type != lexer.RPAREN {
			// Parse first argument
			arg := p.parseExpression(LOWEST)
			if arg != nil {
				methodCall.Args = append(methodCall.Args, arg)
			}

			// Parse additional arguments
			for p.peekToken.Type == lexer.COMMA {
				p.nextToken() // Skip comma
				p.nextToken() // Move to the next argument

				arg = p.parseExpression(LOWEST)
				if arg != nil {
					methodCall.Args = append(methodCall.Args, arg)
				}
			}
		}

		// Check for closing parenthesis
		if p.curToken.Type != lexer.RPAREN {
			p.errors = append(p.errors, fmt.Sprintf("Expected ')' to close arguments, got %s", p.curToken.Type))
			return nil
		}

		// Skip ')' token
		p.nextToken()

		return methodCall
	} else if p.curToken.Type == lexer.DOT {
		// Method call on super (super.method(...))
		// Skip '.' token
		p.nextToken()

		// Parse method name
		if p.curToken.Type != lexer.IDENT {
			p.errors = append(p.errors, fmt.Sprintf("Expected method name after 'super.', got %s", p.curToken.Type))
			return nil
		}

		methodName := p.curToken.Literal
		p.nextToken()

		// Check for opening parenthesis
		if p.curToken.Type != lexer.LPAREN {
			p.errors = append(p.errors, fmt.Sprintf("Expected '(' after method name, got %s", p.curToken.Type))
			return nil
		}

		// Skip '(' token
		p.nextToken()

		// Create a method call node
		methodCall := &MethodCall{
			Object: &Identifier{Name: "super"},
			Method: methodName,
			Args:   []Node{},
		}

		// Parse arguments if any
		if p.curToken.Type != lexer.RPAREN {
			// Parse first argument
			arg := p.parseExpression(LOWEST)
			if arg != nil {
				methodCall.Args = append(methodCall.Args, arg)
			}

			// Parse additional arguments
			for p.peekToken.Type == lexer.COMMA {
				p.nextToken() // Skip comma
				p.nextToken() // Move to the next argument

				arg = p.parseExpression(LOWEST)
				if arg != nil {
					methodCall.Args = append(methodCall.Args, arg)
				}
			}
		}

		// Check for closing parenthesis
		if p.curToken.Type != lexer.RPAREN {
			p.errors = append(p.errors, fmt.Sprintf("Expected ')' to close arguments, got %s", p.curToken.Type))
			return nil
		}

		// Skip ')' token
		p.nextToken()

		return methodCall
	} else {
		p.errors = append(p.errors, fmt.Sprintf("Expected '(' or '.' after 'super', got %s", p.curToken.Type))
		return nil
	}
}

// parseBinaryExpression parses a binary expression
func (p *Parser) parseBinaryExpression(left Node) Node {
	fmt.Printf("DEBUG: parseBinaryExpression - at token: %s, literal: %s\n", p.curToken.Type, p.curToken.Literal)

	// The current token is the operator
	operator := p.curToken.Literal
	precedence := p.curPrecedence()

	// Advance past the operator
	p.nextToken()

	// Parse the right side of the expression
	right := p.parseExpression(precedence)

	// Create and return a binary expression node
	return &BinaryExpr{
		Left:     left,
		Operator: operator,
		Right:    right,
	}
}

// parseClassDefinition parses a class definition
func (p *Parser) parseClassDefinition(className string, parentClass string) Node {
	fmt.Printf("DEBUG: parseClassDefinition - at token: %s, literal: %s\n", p.curToken.Type, p.curToken.Literal)

	classDef := &ClassDef{
		Name:       className,
		Parent:     parentClass,
		Methods:    []Node{},
		TypeParams: []string{},
	}

	// If className is empty, we need to parse it from the token stream
	if className == "" {
		// Skip the 'class' token if we're at it
		if p.curToken.Type == lexer.CLASS {
			p.nextToken()
		}

		// Parse the class name
		if p.curToken.Type != lexer.IDENT {
			p.errors = append(p.errors, fmt.Sprintf("Expected class name, got %s", p.curToken.Type))
			return nil
		}

		classDef.Name = p.curToken.Literal
		p.nextToken()
	}

	// Check for generic type parameters
	if p.curToken.Type == lexer.LT {
		p.nextToken() // Skip '<'

		// Parse type parameters
		for p.curToken.Type != lexer.GT && p.curToken.Type != lexer.EOF {
			if p.curToken.Type == lexer.IDENT {
				classDef.TypeParams = append(classDef.TypeParams, p.curToken.Literal)
			}

			p.nextToken()

			// Skip comma if present
			if p.curToken.Type == lexer.COMMA {
				p.nextToken()
			}
		}

		// Skip '>'
		if p.curToken.Type == lexer.GT {
			p.nextToken()
		}
	}

	// Check for parent class
	if p.curToken.Type == lexer.INHERITS {
		p.nextToken() // Skip 'inherits'

		if p.curToken.Type != lexer.IDENT {
			p.errors = append(p.errors, fmt.Sprintf("Expected parent class name, got %s", p.curToken.Type))
			return nil
		}

		classDef.Parent = p.curToken.Literal
		p.nextToken()
	}

	// Parse class body (methods and fields)
	for p.curToken.Type != lexer.END && p.curToken.Type != lexer.EOF {
		// Parse method definitions
		if p.curToken.Type == lexer.FUNCTION {
			p.nextToken() // Skip 'def'

			// Check if this is a class method (self.method)
			isClassMethod := false
			if p.curToken.Type == lexer.SELF {
				isClassMethod = true
				p.nextToken() // Skip 'self'

				// Skip the dot
				if p.curToken.Type != lexer.DOT {
					p.errors = append(p.errors, fmt.Sprintf("Expected '.' after 'self', got %s", p.curToken.Type))
					// Skip to next statement
					for p.curToken.Type != lexer.FUNCTION && p.curToken.Type != lexer.END && p.curToken.Type != lexer.EOF {
						p.nextToken()
					}
					continue
				}
				p.nextToken() // Skip '.'
			}

			// Parse method name
			if p.curToken.Type != lexer.IDENT {
				p.errors = append(p.errors, fmt.Sprintf("Expected method name, got %s", p.curToken.Type))
				// Skip to next statement
				for p.curToken.Type != lexer.FUNCTION && p.curToken.Type != lexer.END && p.curToken.Type != lexer.EOF {
					p.nextToken()
				}
				continue
			}

			methodName := p.curToken.Literal
			p.nextToken()

			// Parse method parameters
			var parameters []Parameter
			if p.curToken.Type == lexer.LPAREN {
				parameters = p.parseFunctionParameters()
			}

			// Parse return type if present
			var returnType *struct{ TypeName string }
			if p.curToken.Type == lexer.COLON {
				p.nextToken() // Skip ':'
				if p.curToken.Type == lexer.IDENT {
					returnType = &struct{ TypeName string }{TypeName: p.curToken.Literal}
					p.nextToken() // Skip the type name
				}
			}

			// Parse method body
			var methodBody *BlockStmt

			// Create a method block
			methodBody = &BlockStmt{Statements: []Node{}}

			// Parse method body statements until we see 'end' or the next 'def' or class 'end'
			for p.curToken.Type != lexer.END && p.curToken.Type != lexer.FUNCTION && p.curToken.Type != lexer.EOF {
				if p.curToken.Type == lexer.SEMICOLON {
					p.nextToken()
					continue
				}

				stmt := p.parseStatement()
				if stmt != nil {
					methodBody.Statements = append(methodBody.Statements, stmt)
				}

				// Only advance if we're not at the end of the method
				if p.curToken.Type != lexer.END && p.curToken.Type != lexer.FUNCTION && p.curToken.Type != lexer.EOF {
					p.nextToken()
				}
			}

			// Create and add the method definition
			method := &MethodDef{
				Name:          methodName,
				Parameters:    convertToNodeSlice(parameters),
				Body:          methodBody,
				IsClassMethod: isClassMethod,
				ReturnType:    returnType,
			}

			classDef.Methods = append(classDef.Methods, method)

			// If we found 'end', it's the end of this method, skip it
			if p.curToken.Type == lexer.END {
				p.nextToken()
			}
		} else if p.curToken.Type == lexer.IDENT {
			// This might be a field definition
			fieldName := p.curToken.Literal
			// Skip the field name
			p.nextToken()

			// Check if it's a field definition with type annotation
			if p.curToken.Type == lexer.COLON {
				p.nextToken() // Skip ':'

				// Parse the type
				var typeName string
				if p.curToken.Type == lexer.IDENT {
					typeName = p.curToken.Literal
					p.nextToken()
				} else {
					p.errors = append(p.errors, fmt.Sprintf("Expected type name after ':', got %s", p.curToken.Type))
					// Skip to next statement
					for p.curToken.Type != lexer.FUNCTION && p.curToken.Type != lexer.END && p.curToken.Type != lexer.EOF {
						p.nextToken()
					}
					continue
				}

				// Add field to class definition
				field := struct {
					Name          string
					TypeAnnotation struct {
						TypeName string
					}
				}{
					Name: fieldName,
				}
				field.TypeAnnotation.TypeName = typeName
				classDef.Fields = append(classDef.Fields, field)
			} else {
				// Not a field definition, just skip
				p.nextToken()
			}
		} else {
			// Skip other tokens for now
			p.nextToken()
		}
	}

	// Skip 'end' token for the class
	if p.curToken.Type == lexer.END {
		fmt.Printf("DEBUG: parseClassDefinition - Found class end, moving to next token\n")
		p.nextToken()
	}

	fmt.Printf("DEBUG: parseClassDefinition - Finished parsing class, current token: %s, literal: %s\n",
		p.curToken.Type, p.curToken.Literal)

	return classDef
}

// Helper function to convert Parameter slice to Node slice
func convertToNodeSlice(params []Parameter) []Node {
	var nodes []Node
	for _, param := range params {
		// Create a simple node for each parameter - you might want a more sophisticated conversion
		nodes = append(nodes, &Identifier{Name: param.Name})
	}
	return nodes
}

// MethodDef represents a method definition
type MethodDef struct {
	Name          string   // The name of the method
	Parameters    []Node   // The parameters of the method
	Body          Node     // The body of the method
	IsClassMethod bool     // Whether this is a class method (static)
	ReturnType    *struct {
		TypeName string
	}
}

// Type returns the type of the node
func (m *MethodDef) Type() NodeType {
	return MethodDefNode
}

// String returns a string representation of the method definition
func (m *MethodDef) String() string {
	var params []string
	for _, param := range m.Parameters {
		params = append(params, param.String())
	}

	methodType := "instance"
	if m.IsClassMethod {
		methodType = "class"
	}

	returnType := "void"
	if m.ReturnType != nil {
		returnType = m.ReturnType.TypeName
	}

	return fmt.Sprintf("MethodDef(%s, %s, params=[%s], returnType=%s)",
		m.Name, methodType, strings.Join(params, ", "), returnType)
}