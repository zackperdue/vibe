package ast

// NodeType identifies the type of node in the AST
type NodeType string

const (
	// Program nodes
	ProgramNode      NodeType = "Program"

	// Expression nodes
	NumberNode       NodeType = "Number"
	StringNode       NodeType = "String"
	IdentifierNode   NodeType = "Identifier"
	BinaryExprNode   NodeType = "BinaryExpr"
	UnaryExprNode    NodeType = "UnaryExpr"
	CallExprNode     NodeType = "CallExpr"
	IndexExprNode    NodeType = "IndexExpr"
	DotExprNode      NodeType = "DotExpr"
	BooleanNode      NodeType = "Boolean"
	NilNode          NodeType = "Nil"
	ArrayLiteralNode NodeType = "ArrayLiteral"
	VariableUseNode  NodeType = "VariableUse"
	SelfExprNode     NodeType = "SelfExpr"
	SuperCallNode    NodeType = "SuperCall"
	InstanceVarNode  NodeType = "InstanceVar"

	// Statement nodes
	BlockStmtNode    NodeType = "BlockStmt"
	ReturnStmtNode   NodeType = "ReturnStmt"
	IfStmtNode       NodeType = "IfStmt"
	WhileStmtNode    NodeType = "WhileStmt"
	ForStmtNode      NodeType = "ForStmt"
	PrintStmtNode    NodeType = "PrintStmt"
	RequireStmtNode  NodeType = "RequireStmt"

	// Declaration nodes
	AssignmentNode   NodeType = "Assignment"
	FunctionDefNode  NodeType = "FunctionDef"
	VariableDeclNode NodeType = "VariableDecl"
	ClassDefNode     NodeType = "ClassDef"
	MethodDefNode    NodeType = "MethodDef"
	ClassInstNode    NodeType = "ClassInst"
	MethodCallNode   NodeType = "MethodCall"

	// Type-related nodes
	TypeAnnotationNode NodeType = "TypeAnnotation"
	TypeDeclarationNode NodeType = "TypeDeclaration"
)

// Node is the interface that all AST nodes must implement
type Node interface {
	Type() NodeType
	String() string
}

// Program is the root node of every AST
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

// Operator precedences
const (
	LOWEST     = 1
	EQUALS     = 2  // ==
	LESSGREATER = 3  // > or <
	SUM        = 4  // +
	PRODUCT    = 5  // *
	POWER      = 6  // **
	PREFIX     = 7  // -X or !X
	CALL       = 8  // myFunction(X)
	INDEX      = 9  // array[index]
	DOT        = 10  // obj.property
)