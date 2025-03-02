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
	RequireStmtNode  NodeType = "RequireStmt"

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
	condStr := "<nil>"
	if w.Condition != nil {
		condStr = w.Condition.String()
	}

	bodyStr := "<nil>"
	if w.Body != nil {
		bodyStr = w.Body.String()
	}

	return fmt.Sprintf("WhileStmt(%s, %s)", condStr, bodyStr)
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
	seenNonRequireStmt bool // Track if we've seen non-require statements
}

// New creates a new parser
func New(l *lexer.Lexer) *Parser {
	p := &Parser{l: l, seenNonRequireStmt: false}
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
	program := &Program{}
	program.Statements = []Node{}

	fmt.Println("DEBUG: Starting to parse program")

	// Track latest identifier for possible assignments across tokens
	var lastIdent string
	var expectingAssignment bool
	var expectingTypeAnnotation bool
	var typeAnnotation *TypeAnnotation

	for p.curToken.Type != lexer.EOF {
		fmt.Printf("DEBUG: parseProgram - current token: %s, literal: %s, peek token: %s, literal: %s\n",
			p.curToken.Type, p.curToken.Literal, p.peekToken.Type, p.peekToken.Literal)

		// Special handling for class blocks
		if p.curToken.Type == lexer.CLASS || (p.peekToken.Type == lexer.INHERITS && p.curToken.Type == lexer.IDENT) {
			// ... existing code for class handling ...
			// For now, just skip over the class definition to avoid infinite loop
			// Skip 'class' token
			if p.curToken.Type == lexer.CLASS {
				p.nextToken()
			}

			// Skip class name
			if p.curToken.Type == lexer.IDENT {
				p.nextToken()
			}

			// Skip 'inherits' and parent class if present
			if p.curToken.Type == lexer.INHERITS {
				p.nextToken() // skip 'inherits'
				p.nextToken() // skip parent class name
			}

			// Skip until we reach 'end' at the proper nesting level
			depth := 0
			for {
				if p.curToken.Type == lexer.FUNCTION || p.curToken.Type == lexer.IF || p.curToken.Type == lexer.CLASS {
					depth++
				} else if p.curToken.Type == lexer.END {
					depth--
					if depth < 0 {
						break // We've found the end of the class definition
					}
				}

				// Check for end of class at top level
				if depth == 0 && (p.curToken.Type == lexer.CLASS ||
					p.curToken.Type == lexer.EOF) {
					break
				}

				p.nextToken()
			}

			p.nextToken() // Skip the final 'end' token
			fmt.Printf("DEBUG: parseProgram - After skipping class definition, current token: %s, peek token: %s\n",
				p.curToken.Type, p.peekToken.Type)
			continue
		}

		// Check for variable declaration with type annotation (a: string = "hello")
		if p.curToken.Type == lexer.IDENT && p.peekToken.Type == lexer.COLON {
			lastIdent = p.curToken.Literal
			expectingTypeAnnotation = true
			p.nextToken() // Move to COLON token
			p.nextToken() // Move past COLON to the type

			// Parse the type annotation
			typeAnnotation = p.parseTypeAnnotation()

			// If the next token is '=', we also have a value
			if p.curToken.Type == lexer.ASSIGN {
				expectingAssignment = true
				p.nextToken() // Move past ASSIGN to the expression
				fmt.Printf("DEBUG: parseProgram - Recognized variable declaration with type for '%s', now at token: %s\n",
					lastIdent, p.curToken.Type)
			} else {
				// Handle the case where there's no assignment (just a declaration)
				varDecl := &VariableDecl{
					Name:           lastIdent,
					TypeAnnotation: typeAnnotation,
					Value:          nil, // No initial value
				}
				program.Statements = append(program.Statements, varDecl)
				fmt.Printf("DEBUG: parseProgram - added variable declaration: %s\n", varDecl.String())
				expectingTypeAnnotation = false
				expectingAssignment = false
				lastIdent = ""
				typeAnnotation = nil
				p.seenNonRequireStmt = true // Mark that we've seen a non-require statement
			}
		} else if p.curToken.Type == lexer.IDENT && p.peekToken.Type == lexer.ASSIGN {
			// Regular assignment without type annotation
			lastIdent = p.curToken.Literal
			expectingAssignment = true
			p.nextToken() // Move to ASSIGN token
			p.nextToken() // Move past ASSIGN to the expression
			fmt.Printf("DEBUG: parseProgram - Recognized assignment to variable '%s', now at token: %s\n",
				lastIdent, p.curToken.Type)
		}

		stmt := p.parseStatement()
		if stmt != nil {
			// If we were expecting an assignment with a type annotation
			if expectingAssignment && expectingTypeAnnotation {
				varDecl := &VariableDecl{
					Name:           lastIdent,
					TypeAnnotation: typeAnnotation,
					Value:          stmt,
				}
				program.Statements = append(program.Statements, varDecl)
				fmt.Printf("DEBUG: parseProgram - added variable declaration with value: %s\n", varDecl.String())
				expectingAssignment = false
				expectingTypeAnnotation = false
				lastIdent = ""
				typeAnnotation = nil
				p.seenNonRequireStmt = true // Mark that we've seen a non-require statement
			} else if expectingAssignment {
				// Regular assignment without type annotation
				assignment := &Assignment{
					Name:  lastIdent,
					Value: stmt,
				}
				program.Statements = append(program.Statements, assignment)
				fmt.Printf("DEBUG: parseProgram - added assignment: %s\n", assignment.String())
				expectingAssignment = false
				lastIdent = ""
				p.seenNonRequireStmt = true // Mark that we've seen a non-require statement
			} else {
				program.Statements = append(program.Statements, stmt)
				fmt.Printf("DEBUG: parseProgram - added statement: %T - %s\n", stmt, stmt.String())

				// Only set the flag if this is not a require statement
				if stmt.Type() != RequireStmtNode {
					p.seenNonRequireStmt = true
				}
			}
		} else if p.curToken.Type != lexer.EOF {
			// If statement is nil and we're not at EOF, skip this token
			fmt.Printf("DEBUG: parseProgram - statement was nil, skipping token: %s\n", p.curToken.Type)
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
		if p.peekToken.Type == lexer.ASSIGN || p.peekToken.Type == lexer.PLUS_ASSIGN ||
		   p.peekToken.Type == lexer.MINUS_ASSIGN || p.peekToken.Type == lexer.MUL_ASSIGN ||
		   p.peekToken.Type == lexer.DIV_ASSIGN || p.peekToken.Type == lexer.MOD_ASSIGN {
			return p.parseCompoundAssignment()
		}
		return p.parseExpressionStatement()
	case lexer.ASSIGN, lexer.PLUS_ASSIGN, lexer.MINUS_ASSIGN, lexer.MUL_ASSIGN, lexer.DIV_ASSIGN, lexer.MOD_ASSIGN:
		// If we encounter an assignment operator directly, we need to skip it
		// This can happen when parsing multiple assignments in sequence
		return nil
	case lexer.RETURN:
		return p.parseReturnStatement()
	case lexer.PRINT:
		fmt.Printf("DEBUG: parseStatement - detected print token, calling parsePrintStatement\n")
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
	case lexer.REQUIRE:
		fmt.Println("DEBUG: Detected REQUIRE token in parseStatement, calling parseRequireStatement")
		return p.parseRequireStatement()
	case lexer.CLASS:
		fmt.Println("DEBUG: Detected CLASS token in parseStatement, calling parseClassDefinition")
		return p.parseClassDefinition()
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

	// Handle function definition without parentheses (no parameters)
	if p.curToken.Type == lexer.COLON {
		// Function has no parameters and uses the no-parentheses syntax
		funcDef.Parameters = []Parameter{}

		// Parse return type
		p.nextToken()
		funcDef.ReturnType = p.parseTypeAnnotation()
	} else if p.curToken.Type == lexer.LPAREN {
		// Traditional function with parameters in parentheses
		funcDef.Parameters = p.parseFunctionParameters()

		// Check for return type annotation with : syntax
		if p.curToken.Type == lexer.COLON {
			p.nextToken()
			funcDef.ReturnType = p.parseTypeAnnotation()
		} else {
			// Default return type is "int"
			funcDef.ReturnType = &TypeAnnotation{TypeName: "int"}
		}
	} else {
		p.errors = append(p.errors, fmt.Sprintf("Expected '(' or ':' after function name, got %s", p.curToken.Type))
		return nil
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
	fmt.Printf("DEBUG: parseWhileStatement - starting at token: %s\n", p.curToken.Type)

	// Skip 'while' keyword
	p.nextToken()

	fmt.Printf("DEBUG: parseWhileStatement - after skipping 'while', at token: %s, peek: %s\n",
		p.curToken.Type, p.peekToken.Type)

	// If we have an identifier followed by a comparison operator, handle it specially
	if p.curToken.Type == lexer.IDENT &&
		(p.peekToken.Type == lexer.LT || p.peekToken.Type == lexer.GT ||
		 p.peekToken.Type == lexer.LT_EQ || p.peekToken.Type == lexer.GT_EQ ||
		 p.peekToken.Type == lexer.EQ || p.peekToken.Type == lexer.NOT_EQ) {

		fmt.Printf("DEBUG: parseWhileStatement - detected comparison expression\n")

		// Create the left side of the comparison
		left := &Identifier{Name: p.curToken.Literal}

		// Move to the comparison operator
		p.nextToken()
		operator := p.curToken.Literal

		// Move to the right side
		p.nextToken()

		// Parse the right side
		var right Node
		switch p.curToken.Type {
		case lexer.INT:
			value, _ := strconv.ParseFloat(p.curToken.Literal, 64)
			right = &NumberLiteral{Value: value, IsInt: true}
		case lexer.FLOAT:
			value, _ := strconv.ParseFloat(p.curToken.Literal, 64)
			right = &NumberLiteral{Value: value, IsInt: false}
		case lexer.IDENT:
			right = &Identifier{Name: p.curToken.Literal}
		default:
			p.errors = append(p.errors, fmt.Sprintf("Expected number or identifier after comparison operator, got %s", p.curToken.Type))
			right = &NumberLiteral{Value: 0, IsInt: true} // Default to avoid nil
		}

		// Create the condition as a binary expression
		condition := &BinaryExpr{
			Left:     left,
			Operator: operator,
			Right:    right,
		}

		fmt.Printf("DEBUG: parseWhileStatement - created condition: %s\n", condition.String())

		// Move to the next token (should be 'do')
		p.nextToken()

		// Check for 'do' keyword
		if p.curToken.Type != lexer.DO {
			p.errors = append(p.errors, fmt.Sprintf("Expected 'do' after while condition, got %s", p.curToken.Type))
			// Try to find it in the next token
			if p.peekToken.Type == lexer.DO {
				p.nextToken() // Move to 'do'
			}
		}

		// Skip 'do' if we're on it
		if p.curToken.Type == lexer.DO {
			p.nextToken()
		}

		// Parse while loop body directly
		body := &BlockStmt{Statements: []Node{}}

		// Parse statements until we see 'end' or EOF
		for p.curToken.Type != lexer.END && p.curToken.Type != lexer.EOF {
			fmt.Printf("DEBUG: parseWhileStatement - parsing statement in body, token: %s\n", p.curToken.Type)

			var stmt Node

			// Handle print statements
			if p.curToken.Type == lexer.PRINT {
				fmt.Printf("DEBUG: parseWhileStatement - detected print statement\n")
				// Skip 'print' or 'puts' keyword
				p.nextToken()

				// Parse the expression to print
				expr := p.parseExpression(LOWEST)

				// Create a print statement
				stmt = &PrintStmt{Value: expr}
				fmt.Printf("DEBUG: parseWhileStatement - created print statement: %s\n", stmt.String())
			} else if p.curToken.Type == lexer.IDENT &&
				(p.peekToken.Type == lexer.ASSIGN || p.peekToken.Type == lexer.PLUS_ASSIGN ||
				 p.peekToken.Type == lexer.MINUS_ASSIGN || p.peekToken.Type == lexer.MUL_ASSIGN ||
				 p.peekToken.Type == lexer.DIV_ASSIGN || p.peekToken.Type == lexer.MOD_ASSIGN) {
				// Handle assignments
				stmt = p.parseCompoundAssignment()
			} else {
				// Handle other statements
				stmt = p.parseStatement()
			}

			if stmt != nil {
				fmt.Printf("DEBUG: parseWhileStatement - added statement to body: %T\n", stmt)
				body.Statements = append(body.Statements, stmt)
			} else {
				fmt.Printf("DEBUG: parseWhileStatement - statement was nil, skipping\n")
			}

			// Move to the next token
			if p.curToken.Type != lexer.END && p.curToken.Type != lexer.EOF {
				p.nextToken()
			}
		}

		// Skip the 'end' token if present
		if p.curToken.Type == lexer.END {
			p.nextToken()
		} else {
			p.errors = append(p.errors, "Expected 'end' to close while loop")
		}

		return &WhileStmt{
			Condition: condition,
			Body:      body,
		}
	} else {
		// Fall back to the regular expression parsing for other cases
		condition := p.parseExpression(LOWEST)
		if condition == nil {
			p.errors = append(p.errors, "Invalid or missing condition in while statement")
			condition = &BooleanLiteral{Value: false} // Default to false to avoid nil pointer
		}

		// Check for 'do' keyword
		if p.curToken.Type != lexer.DO {
			p.errors = append(p.errors, fmt.Sprintf("Expected 'do' after while condition, got %s", p.curToken.Type))
			// Try to find it in the next token
			if p.peekToken.Type == lexer.DO {
				p.nextToken() // Move to 'do'
			}
		}

		// Skip 'do' if we're on it
		if p.curToken.Type == lexer.DO {
			p.nextToken()
		}

		// Parse while loop body directly
		body := &BlockStmt{Statements: []Node{}}

		// Parse statements until we see 'end' or EOF
		for p.curToken.Type != lexer.END && p.curToken.Type != lexer.EOF {
			fmt.Printf("DEBUG: parseWhileStatement - parsing statement in body, token: %s\n", p.curToken.Type)

			var stmt Node

			// Handle print statements
			if p.curToken.Type == lexer.PRINT {
				fmt.Printf("DEBUG: parseWhileStatement - detected print statement\n")
				// Skip 'print' or 'puts' keyword
				p.nextToken()

				// Parse the expression to print
				expr := p.parseExpression(LOWEST)

				// Create a print statement
				stmt = &PrintStmt{Value: expr}
				fmt.Printf("DEBUG: parseWhileStatement - created print statement: %s\n", stmt.String())
			} else if p.curToken.Type == lexer.IDENT &&
				(p.peekToken.Type == lexer.ASSIGN || p.peekToken.Type == lexer.PLUS_ASSIGN ||
				 p.peekToken.Type == lexer.MINUS_ASSIGN || p.peekToken.Type == lexer.MUL_ASSIGN ||
				 p.peekToken.Type == lexer.DIV_ASSIGN || p.peekToken.Type == lexer.MOD_ASSIGN) {
				// Handle assignments
				stmt = p.parseCompoundAssignment()
			} else {
				// Handle other statements
				stmt = p.parseStatement()
			}

			if stmt != nil {
				fmt.Printf("DEBUG: parseWhileStatement - added statement to body: %T\n", stmt)
				body.Statements = append(body.Statements, stmt)
			} else {
				fmt.Printf("DEBUG: parseWhileStatement - statement was nil, skipping\n")
			}

			// Move to the next token
			if p.curToken.Type != lexer.END && p.curToken.Type != lexer.EOF {
				p.nextToken()
			}
		}

		// Skip the 'end' token if present
		if p.curToken.Type == lexer.END {
			p.nextToken()
		} else {
			p.errors = append(p.errors, "Expected 'end' to close while loop")
		}

		return &WhileStmt{
			Condition: condition,
			Body:      body,
		}
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
	fmt.Printf("DEBUG: parsePrintStatement - starting at token: %s\n", p.curToken.Type)

	stmt := &PrintStmt{}

	// Skip 'print' or 'puts' keyword
	p.nextToken()

	fmt.Printf("DEBUG: parsePrintStatement - after skipping 'print', at token: %s\n", p.curToken.Type)

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

		fmt.Printf("DEBUG: parsePrintStatement - created print statement: %s\n", stmt.String())
	}

	return stmt
}

func (p *Parser) parseCompoundAssignment() Node {
	debugf("parseCompoundAssignment - at token: %s", p.curToken.Type)

	// Save the variable name
	name := p.curToken.Literal

	// Skip to the assignment operator
	p.nextToken()

	// Remember the assignment operator
	operator := p.curToken.Type

	// Skip the assignment operator
	p.nextToken()

	var value Node

	// For compound assignments, create a binary expression
	if operator != lexer.ASSIGN {
		// Get the left side (the variable)
		left := &Identifier{Name: name}

		// Determine the binary operator based on the compound assignment
		var binOp string
		switch operator {
		case lexer.PLUS_ASSIGN:
			binOp = "+"
		case lexer.MINUS_ASSIGN:
			binOp = "-"
		case lexer.MUL_ASSIGN:
			binOp = "*"
		case lexer.DIV_ASSIGN:
			binOp = "/"
		case lexer.MOD_ASSIGN:
			binOp = "%"
		}

		// Parse the right-hand expression
		right := p.parseExpression(LOWEST)
		if right == nil {
			fmt.Println("DEBUG: parseCompoundAssignment - Failed to parse right side of assignment")
			return nil
		}

		// Create a binary expression for the operation
		value = &BinaryExpr{
			Left:     left,
			Operator: binOp,
			Right:    right,
		}
	} else {
		// For regular assignment, just parse the expression
		value = p.parseExpression(LOWEST)
		if value == nil {
			fmt.Println("DEBUG: parseCompoundAssignment - Failed to parse right side of assignment")
			return nil
		}
	}

	// Create and return the assignment node
	assignment := &Assignment{
		Name:  name,
		Value: value,
	}

	debugf("parseCompoundAssignment - created assignment: %s = %s", name, value.String())
	return assignment
}

func (p *Parser) parseExpressionStatement() Node {
	return p.parseExpression(0)
}

func (p *Parser) parseExpression(precedence int) Node {
	fmt.Printf("DEBUG: parseExpression - at token: %s, literal: %s\n", p.curToken.Type, p.curToken.Literal)
	fmt.Printf("DEBUG: parseExpression - precedence: %d, peek token: %s\n", precedence, p.peekToken.Type)

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
			// Only handle as generic type if we're not in a comparison context
			// Check if the token after '<' is an identifier (type name)
			if p.peekTokenIs(lexer.LT) && p.peekTokenIs(lexer.IDENT) {
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
		}

		// Enable function calls without parentheses for functions with no parameters
		// Only do this if we're not in a context where the identifier might be used for something else
		// like an assignment target, a property name, etc.
		// We can infer this is a function call if we're at the end of an expression
		if !isInfixOperator(p.peekToken.Type) &&
		   p.peekToken.Type != lexer.LPAREN &&
		   p.peekToken.Type != lexer.LBRACKET &&
		   p.peekToken.Type != lexer.DOT &&
		   p.peekToken.Type != lexer.ASSIGN &&
		   p.peekToken.Type != lexer.PLUS_ASSIGN &&
		   p.peekToken.Type != lexer.MINUS_ASSIGN &&
		   p.peekToken.Type != lexer.MUL_ASSIGN &&
		   p.peekToken.Type != lexer.DIV_ASSIGN &&
		   p.peekToken.Type != lexer.MOD_ASSIGN {
			// Create a CallExpr with empty args
			leftExp = &CallExpr{
				Function: leftExp,
				Args:     []Node{},
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

	// Skip to the next token, but only if not DO, as we need to preserve this for statements
	if p.peekToken.Type != lexer.DO {
		p.nextToken()
	}

	// Now parse any infix expressions
	for precedence < p.curPrecedence() && p.curToken.Type != lexer.EOF {
		// Don't proceed with infix parsing if the next token is DO
		if p.peekToken.Type == lexer.DO {
			break
		}

		fmt.Printf("DEBUG: parseExpression - infix - current token: %s, precedence: %d, curPrecedence: %d\n",
			p.curToken.Type, precedence, p.curPrecedence())

		switch p.curToken.Type {
		case lexer.PLUS, lexer.MINUS, lexer.ASTERISK, lexer.SLASH, lexer.MODULO,
			lexer.EQ, lexer.NOT_EQ, lexer.LT, lexer.GT, lexer.LT_EQ, lexer.GT_EQ,
			lexer.AND, lexer.OR:
			fmt.Printf("DEBUG: parseExpression - calling parseBinaryExpression with operator: %s\n", p.curToken.Literal)
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
			lexer.AND, lexer.OR:
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
	fmt.Printf("DEBUG: parseArrayLiteral - after '[', current token: %s, literal: %s\n", p.curToken.Type, p.curToken.Literal)

	// Empty array case
	if p.curToken.Type == lexer.RBRACKET {
		// Skip ']'
		p.nextToken()
		fmt.Printf("DEBUG: parseArrayLiteral - empty array, current token after ]: %s, peek: %s\n",
			p.curToken.Type, p.peekToken.Type)
		return arrayLit
	}

	// Parse first element
	firstElement := p.parseExpression(LOWEST)
	if firstElement != nil {
		arrayLit.Elements = append(arrayLit.Elements, firstElement)
	}

	// Parse remaining elements
	for p.curToken.Type == lexer.COMMA {
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
	if p.curToken.Type != lexer.RBRACKET {
		p.errors = append(p.errors, fmt.Sprintf("Expected ']', got %s", p.curToken.Type))
		return nil
	}

	// Important: We need to check if the next token is DO before consuming it
	if p.peekToken.Type == lexer.DO {
		fmt.Printf("DEBUG: parseArrayLiteral - detected DO after array, preserving it\n")
		// We want to move past the ']' but not consume any token after that
		p.nextToken() // Moves to the DO token
		return arrayLit
	}

	// Skip the closing bracket
	p.nextToken()
	fmt.Printf("DEBUG: parseArrayLiteral - array with %d elements, current token after ]: %s, peek: %s\n",
		len(arrayLit.Elements), p.curToken.Type, p.peekToken.Type)

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
		p.errors = append(p.errors, fmt.Sprintf("Expected ')', got %s", p.peekToken.Type))
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

	// Check for range operator '..' in for loops (e.g. 0..5)
	if p.peekToken.Type == lexer.DOT {
		// We have a '..' range operator
		// Skip first '.' token
		p.nextToken()
		// Skip second '.' token
		p.nextToken()

		// Parse the end of the range
		end := p.parseExpression(LOWEST)

		// Create a range expression (represented as a binary expression with '..' operator)
		return &BinaryExpr{
			Left:     left,
			Operator: "..",
			Right:    end,
		}
	}

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

			arg := p.parseExpression(LOWEST)
			if arg != nil {
				methodCall.Args = append(methodCall.Args, arg)
				fmt.Printf("DEBUG: parseClassInstantiation - added additional argument: %s\n", arg.String())
			}
		}

		// Ensure we are at the closing parenthesis
		if p.peekToken.Type == lexer.RPAREN {
			p.nextToken()
		} else {
			p.errors = append(p.errors, fmt.Sprintf("Expected ')' to close arguments, got %s", p.peekToken.Type))
			return nil
		}

		p.nextToken() // Skip the closing parenthesis
		return methodCall
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

	// Create the ClassInst node
	classInst := &ClassInst{
		Token:     p.curToken,
		Class:     left,
		Arguments: []Node{},
	}

	// Check for opening parenthesis
	if !p.expectPeek(lexer.LPAREN) {
		return nil
	}
	fmt.Printf("DEBUG: parseClassInstantiation - after 'new', token: %s, literal: %s\n", p.curToken.Type, p.curToken.Literal)

	// Skip '('
	p.nextToken()
	fmt.Printf("DEBUG: parseClassInstantiation - after '(', token: %s, literal: %s\n", p.curToken.Type, p.curToken.Literal)

	// Handle empty arguments list
	if p.curToken.Type == lexer.RPAREN {
		p.nextToken() // Skip ')'
		fmt.Printf("DEBUG: parseClassInstantiation - empty args, after ')', token: %s, literal: %s\n", p.curToken.Type, p.curToken.Literal)
		return classInst
	}

	// Parse first argument
	arg := p.parseExpression(LOWEST)
	if arg != nil {
		classInst.Arguments = append(classInst.Arguments, arg)
		fmt.Printf("DEBUG: parseClassInstantiation - added first argument: %s\n", arg.String())
	}

	// Parse additional arguments
	for p.curToken.Type == lexer.COMMA {
		p.nextToken() // Skip comma and move to next arg
		fmt.Printf("DEBUG: parseClassInstantiation - after comma, token: %s, literal: %s\n", p.curToken.Type, p.curToken.Literal)

		arg := p.parseExpression(LOWEST)
		if arg != nil {
			classInst.Arguments = append(classInst.Arguments, arg)
			fmt.Printf("DEBUG: parseClassInstantiation - added additional argument: %s\n", arg.String())
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
	fmt.Printf("DEBUG: parseClassInstantiation - created class instantiation with %d args\n", len(classInst.Arguments))

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
	Token      lexer.Token
	Class      Node
	Arguments  []Node
	TypeArgs   []Node
}

// Type returns the type of the node
func (c *ClassInst) Type() NodeType {
	return ClassInstNode
}

// String returns a string representation of the class instantiation
func (c *ClassInst) String() string {
	var args []string
	for _, arg := range c.Arguments {
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
	// Create new ForStmt node
	stmt := &ForStmt{}

	// Skip the 'for' token
	p.nextToken()

	// Parse iterator (variable name)
	if p.curToken.Type != lexer.IDENT {
		p.errors = append(p.errors, fmt.Sprintf("Expected identifier for iterator, got %s", p.curToken.Type))
		return nil
	}
	stmt.Iterator = p.curToken.Literal

	// Expect 'in' token
	p.nextToken()
	if p.curToken.Type != lexer.IN {
		p.errors = append(p.errors, fmt.Sprintf("Expected 'in' after iterator, got %s", p.curToken.Type))
		return nil
	}

	// Skip 'in' token
	p.nextToken()

	// Parse the iterable expression
	stmt.Iterable = p.parseExpression(LOWEST)

	// Special handling for array literals which might have consumed the 'do' token
	// due to how parseArrayLiteral works
	if p.curToken.Type == lexer.END {
		// We've somehow skipped the 'do' token and landed on 'end' directly
		// This means the loop body is empty, so we'll create an empty body
		stmt.Body = &BlockStmt{Statements: []Node{}}

		// Skip the 'end' token
		p.nextToken()
		return stmt
	}

	// The 'do' keyword is optional
	if p.curToken.Type == lexer.DO {
		// Skip 'do'
		p.nextToken()
	} else if p.peekToken.Type == lexer.DO {
		// Move to and skip 'do'
		p.nextToken()
		p.nextToken()
	}
	// No error if 'do' is not present - proceed with parsing the body

	// Create a new block for the body
	bodyBlock := &BlockStmt{Statements: []Node{}}

	// Parse statements until we reach 'end'
	for p.curToken.Type != lexer.END && p.curToken.Type != lexer.EOF {
		stmt := p.parseStatement()
		if stmt != nil {
			bodyBlock.Statements = append(bodyBlock.Statements, stmt)
		}

		// Check if we've reached the end token after parsing a statement
		if p.curToken.Type == lexer.END {
			break
		}

		p.nextToken()
	}

	// Set the body
	stmt.Body = bodyBlock

	// Skip the 'end' token if present
	if p.curToken.Type == lexer.END {
		p.nextToken()
	} else if p.curToken.Type == lexer.EOF {
		p.errors = append(p.errors, "Expected 'end' at the end of the for statement")
	}

	return stmt
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

				arg := p.parseExpression(LOWEST)
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

		// Skip '('
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

				arg := p.parseExpression(LOWEST)
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

		// Skip ')'
		p.nextToken()

		return methodCall
	} else {
		p.errors = append(p.errors, fmt.Sprintf("Expected '(' or '.' after 'super', got %s", p.curToken.Type))
		return nil
	}
}

// parseBinaryExpression parses a binary expression
func (p *Parser) parseBinaryExpression(left Node) Node {
	// Save references to the current token (operator token)
	operator := p.curToken.Literal
	fmt.Printf("DEBUG: parseBinaryExpression - at token: %s, literal: %s\n", p.curToken.Type, p.curToken.Literal)
	fmt.Printf("DEBUG: parseBinaryExpression - left: %s\n", left.String())
	fmt.Printf("DEBUG: parseBinaryExpression - operator: %s, precedence: %d\n", operator, p.curPrecedence())

	// Remember precedence of the operator
	precedence := p.curPrecedence()

	// Skip the operator token
	p.nextToken()
	fmt.Printf("DEBUG: parseBinaryExpression - now at token: %s, literal: %s\n", p.curToken.Type, p.curToken.Literal)

	// Parse the right-hand-side expression
	var right Node

	// Handle the case where an identifier might be automatically converted to a function call
	// without parentheses. We need to check if it's an identifier before automatic conversion.
	if p.curToken.Type == lexer.IDENT {
		// Create an identifier node first
		identNode := &Identifier{Name: p.curToken.Literal}
		p.nextToken() // Consume the identifier

		// Check if it should be treated as a function call without parentheses
		if !isInfixOperator(p.curToken.Type) &&
		   p.curToken.Type != lexer.LPAREN &&
		   p.curToken.Type != lexer.LBRACKET &&
		   p.curToken.Type != lexer.DOT &&
		   p.curToken.Type != lexer.ASSIGN &&
		   p.curToken.Type != lexer.PLUS_ASSIGN &&
		   p.curToken.Type != lexer.MINUS_ASSIGN &&
		   p.curToken.Type != lexer.MUL_ASSIGN &&
		   p.curToken.Type != lexer.DIV_ASSIGN &&
		   p.curToken.Type != lexer.MOD_ASSIGN {
			// Create a CallExpr with empty args
			right = &CallExpr{
				Function: identNode,
				Args:     []Node{},
			}
		} else {
			// Just use it as a regular identifier
			right = identNode
		}
	} else {
		// Regular expression parsing
		right = p.parseExpression(precedence)
	}

	if right == nil {
		return nil
	}

	fmt.Printf("DEBUG: parseBinaryExpression - right: %s\n", right.String())

	// Create and return a binary expression
	expr := &BinaryExpr{
		Left:     left,
		Operator: operator,
		Right:    right,
	}
	fmt.Printf("DEBUG: parseBinaryExpression - created expression: %s\n", expr.String())

	return expr
}

func (p *Parser) expectPeek(t lexer.TokenType) bool {
	if p.peekToken.Type == t {
		p.nextToken()
		return true
	}
	p.errors = append(p.errors, fmt.Sprintf("Expected %s, got %s", t, p.peekToken.Type))
	return false
}

func (p *Parser) peekTokenIs(t lexer.TokenType) bool {
	return p.peekToken.Type == t
}


// RequireStmt represents a require statement
type RequireStmt struct {
	Path string
}

func (r *RequireStmt) Type() NodeType { return RequireStmtNode }
func (r *RequireStmt) String() string {
	return fmt.Sprintf("RequireStmt(%s)", r.Path)
}

func (p *Parser) parseRequireStatement() Node {
	fmt.Printf("DEBUG: parseRequireStatement - starting at token: %s\n", p.curToken.Type)

	// Check if we've already seen non-require statements
	if p.seenNonRequireStmt {
		p.errors = append(p.errors, "Error: 'require' statements must appear at the top of the file")
	}

	// Skip 'require' keyword
	p.nextToken()

	// Parse path string
	if p.curToken.Type != lexer.STRING {
		p.errors = append(p.errors, fmt.Sprintf("Expected string path after 'require', got %s", p.curToken.Type))
		return nil
	}

	path := p.curToken.Literal
	p.nextToken() // Move past the string

	return &RequireStmt{
		Path: path,
	}
}

func (p *Parser) parseClassDefinition() Node {
	fmt.Printf("DEBUG: parseClassDefinition - starting at token: %s\n", p.curToken.Type)

	// Skip 'class' keyword
	p.nextToken()

	// Get class name
	if p.curToken.Type != lexer.IDENT {
		p.errors = append(p.errors, fmt.Sprintf("Expected class name after 'class', got %s", p.curToken.Type))
		return nil
	}

	className := p.curToken.Literal
	p.nextToken()

	// Check for inheritance
	var parentClass string
	if p.curToken.Type == lexer.INHERITS {
		p.nextToken() // Skip 'inherits'

		if p.curToken.Type != lexer.IDENT {
			p.errors = append(p.errors, fmt.Sprintf("Expected parent class name after 'inherits', got %s", p.curToken.Type))
			return nil
		}

		parentClass = p.curToken.Literal
		p.nextToken()
	}

	// Check for 'do' keyword
	if p.curToken.Type != lexer.DO {
		p.errors = append(p.errors, fmt.Sprintf("Expected 'do' after class declaration, got %s", p.curToken.Type))
		return nil
	}

	// Skip 'do' keyword
	p.nextToken()

	// Parse methods and instance variables
	methods := []Node{}

	// Parse methods and instance variables
	for p.curToken.Type != lexer.END && p.curToken.Type != lexer.EOF {
		var stmt Node

		if p.curToken.Type == lexer.FUNCTION {
			// Parse method definition
			stmt = p.parseFunctionDefinition()
			if stmt != nil {
				methods = append(methods, stmt)
			}
		} else if p.curToken.Type == lexer.AT {
			// Parse instance variable
			stmt = p.parseInstanceVariable()
		} else {
			// Parse other statements
			stmt = p.parseStatement()
		}

		// Move to the next token if not at the end
		if p.curToken.Type != lexer.END && p.curToken.Type != lexer.EOF {
			p.nextToken()
		}
	}

	// Skip the 'end' token
	if p.curToken.Type == lexer.END {
		p.nextToken()
	} else {
		p.errors = append(p.errors, "Expected 'end' to close class definition")
	}

	return &ClassDef{
		Name:    className,
		Parent:  parentClass,
		Methods: methods,
	}
}