package ast

import (
	"fmt"
	"strings"
)

// BlockStmt represents a block of statements in the AST
type BlockStmt struct {
	Statements []Node
}

func (b *BlockStmt) Type() NodeType { return BlockStmtNode }
func (b *BlockStmt) String() string {
	stmts := []string{}
	for _, stmt := range b.Statements {
		stmts = append(stmts, stmt.String())
	}
	return fmt.Sprintf("Block { %s }", strings.Join(stmts, "; "))
}

// ReturnStmt represents a return statement in the AST
type ReturnStmt struct {
	Value Node
}

func (r *ReturnStmt) Type() NodeType { return ReturnStmtNode }
func (r *ReturnStmt) String() string {
	if r.Value == nil {
		return "return"
	}
	return fmt.Sprintf("return %s", r.Value.String())
}

// ElseIfBlock represents an else-if branch in an if statement
type ElseIfBlock struct {
	Condition   Node
	Consequence *BlockStmt
}

// IfStmt represents an if statement in the AST
type IfStmt struct {
	Condition     Node
	Consequence   *BlockStmt
	Alternative   *BlockStmt
	ElseIfBlocks  []ElseIfBlock
}

func (i *IfStmt) Type() NodeType { return IfStmtNode }
func (i *IfStmt) String() string {
	var result string
	result = fmt.Sprintf("if %s %s", i.Condition.String(), i.Consequence.String())

	for _, elseIf := range i.ElseIfBlocks {
		result += fmt.Sprintf(" elsif %s %s", elseIf.Condition.String(), elseIf.Consequence.String())
	}

	if i.Alternative != nil {
		result += fmt.Sprintf(" else %s", i.Alternative.String())
	}

	return result
}

// WhileStmt represents a while statement in the AST
type WhileStmt struct {
	Condition Node
	Body      *BlockStmt
}

func (w *WhileStmt) Type() NodeType { return WhileStmtNode }
func (w *WhileStmt) String() string {
	return fmt.Sprintf("while %s %s", w.Condition.String(), w.Body.String())
}

// ForStmt represents a for statement in the AST
type ForStmt struct {
	Iterator  string     // The variable that will hold each element
	Iterable  Node       // The expression to iterate over
	Body      *BlockStmt
}

func (f *ForStmt) Type() NodeType { return ForStmtNode }
func (f *ForStmt) String() string {
	return fmt.Sprintf("for %s in %s %s", f.Iterator, f.Iterable.String(), f.Body.String())
}

// PrintStmt represents a print statement in the AST
type PrintStmt struct {
	Value Node
}

func (p *PrintStmt) Type() NodeType { return PrintStmtNode }
func (p *PrintStmt) String() string {
	return fmt.Sprintf("puts %s", p.Value.String())
}

// RequireStmt represents a require statement in the AST
type RequireStmt struct {
	Path string
}

func (r *RequireStmt) Type() NodeType { return RequireStmtNode }
func (r *RequireStmt) String() string {
	return fmt.Sprintf("require %q", r.Path)
}

// Assignment represents an assignment statement in the AST
type Assignment struct {
	Name  string
	Value Node
}

func (a *Assignment) Type() NodeType { return AssignmentNode }
func (a *Assignment) String() string {
	if a.Value == nil {
		return fmt.Sprintf("%s = nil", a.Name)
	}
	return fmt.Sprintf("%s = %s", a.Name, a.Value.String())
}