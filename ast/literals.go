package ast

import (
	"fmt"
	"strings"
)

// NumberLiteral represents a number literal in the AST
type NumberLiteral struct {
	Value float64
	IsInt bool
}

func (n *NumberLiteral) Type() NodeType { return NumberNode }
func (n *NumberLiteral) String() string {
	if n.IsInt {
		return fmt.Sprintf("%d", int(n.Value))
	}
	return fmt.Sprintf("%f", n.Value)
}

// StringLiteral represents a string literal in the AST
type StringLiteral struct {
	Value string
}

func (s *StringLiteral) Type() NodeType { return StringNode }
func (s *StringLiteral) String() string { return fmt.Sprintf("String(%q)", s.Value) }

// BooleanLiteral represents a boolean literal in the AST
type BooleanLiteral struct {
	Value bool
}

func (b *BooleanLiteral) Type() NodeType { return BooleanNode }
func (b *BooleanLiteral) String() string { return fmt.Sprintf("Boolean(%t)", b.Value) }

// NilLiteral represents a nil literal in the AST
type NilLiteral struct{}

func (n *NilLiteral) Type() NodeType { return NilNode }
func (n *NilLiteral) String() string { return "Nil" }

// ArrayLiteral represents an array literal in the AST
type ArrayLiteral struct {
	Elements []Node
}

func (a *ArrayLiteral) Type() NodeType { return ArrayLiteralNode }
func (a *ArrayLiteral) String() string {
	elements := []string{}
	for _, el := range a.Elements {
		elements = append(elements, el.String())
	}
	return fmt.Sprintf("[%s]", strings.Join(elements, ", "))
}

// Identifier represents an identifier in the AST
type Identifier struct {
	Name string
}

func (i *Identifier) Type() NodeType { return IdentifierNode }
func (i *Identifier) String() string { return i.Name }

// SelfExpr represents a 'self' expression in the AST
type SelfExpr struct{}

func (s *SelfExpr) Type() NodeType { return SelfExprNode }
func (s *SelfExpr) String() string { return "self" }

// InstanceVar represents an instance variable (@var) in the AST
type InstanceVar struct {
	Name string
}

func (i *InstanceVar) Type() NodeType { return InstanceVarNode }
func (i *InstanceVar) String() string { return "@" + i.Name }