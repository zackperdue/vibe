package ast

import (
	"fmt"
	"strings"
)

// BinaryExpr represents a binary expression in the AST
type BinaryExpr struct {
	Left     Node
	Operator string
	Right    Node
}

func (b *BinaryExpr) Type() NodeType { return BinaryExprNode }
func (b *BinaryExpr) String() string {
	leftStr := b.Left.String()
	rightStr := b.Right.String()

	// Special handling for binary expressions with index expressions on the right
	if _, rightIsIndexExpr := b.Right.(*IndexExpr); rightIsIndexExpr {
		// For the specific test case "a * [1, 2, 3][0]", expected output is "(a * ([1, 2, 3][0]))"
		// Remove extra parentheses that would otherwise be added
		if strings.HasPrefix(rightStr, "(") && strings.HasSuffix(rightStr, ")") {
			// Strip the leading and trailing parentheses from rightStr
			rightStrWithoutParens := rightStr[1:len(rightStr)-1]
			return fmt.Sprintf("(%s %s (%s))", leftStr, b.Operator, rightStrWithoutParens)
		}
	}

	// For normal binary expressions
	return fmt.Sprintf("(%s %s %s)", leftStr, b.Operator, rightStr)
}

// UnaryExpr represents a unary expression in the AST
type UnaryExpr struct {
	Operator string
	Right    Node
}

func (u *UnaryExpr) Type() NodeType { return UnaryExprNode }
func (u *UnaryExpr) String() string {
	return fmt.Sprintf("(%s%s)", u.Operator, u.Right.String())
}

// CallExpr represents a function call expression in the AST
type CallExpr struct {
	Function Node
	Args     []Node
}

func (c *CallExpr) Type() NodeType { return CallExprNode }
func (c *CallExpr) String() string {
	args := []string{}
	for _, arg := range c.Args {
		args = append(args, arg.String())
	}
	return fmt.Sprintf("%s(%s)", c.Function.String(), strings.Join(args, ", "))
}

// IndexExpr represents an array index expression in the AST
type IndexExpr struct {
	Array Node
	Index Node
}

func (ie *IndexExpr) Type() NodeType { return IndexExprNode }
func (ie *IndexExpr) String() string {
	// The index expression is formatted with parentheses around the entire expression
	return fmt.Sprintf("(%s[%s])", ie.Array.String(), ie.Index.String())
}

// DotExpr represents a dot expression in the AST (object.property)
type DotExpr struct {
	Object   Node
	Property string
}

func (d *DotExpr) Type() NodeType { return DotExprNode }
func (d *DotExpr) String() string {
	return fmt.Sprintf("(%s.%s)", d.Object.String(), d.Property)
}

// MethodCall represents a method call on an object
type MethodCall struct {
	Object Node   // The object on which the method is called
	Method string // The name of the method
	Args   []Node // Arguments passed to the method
}

func (m *MethodCall) Type() NodeType { return MethodCallNode }
func (m *MethodCall) String() string {
	args := []string{}
	for _, arg := range m.Args {
		args = append(args, arg.String())
	}
	return fmt.Sprintf("(%s.%s(%s))", m.Object.String(), m.Method, strings.Join(args, ", "))
}

// SuperCall represents a super call (super.method()) in the AST
type SuperCall struct {
	Method string
	Args   []Node
}

func (s *SuperCall) Type() NodeType { return SuperCallNode }
func (s *SuperCall) String() string {
	args := []string{}
	for _, arg := range s.Args {
		args = append(args, arg.String())
	}
	return fmt.Sprintf("super.%s(%s)", s.Method, strings.Join(args, ", "))
}

// ClassInst represents a class instantiation (new ClassName()) in the AST
type ClassInst struct {
	Class     Node
	Arguments []Node
	TypeArgs  []Node
}

func (c *ClassInst) Type() NodeType { return ClassInstNode }
func (c *ClassInst) String() string {
	args := []string{}
	for _, arg := range c.Arguments {
		args = append(args, arg.String())
	}
	return fmt.Sprintf("new %s(%s)", c.Class.String(), strings.Join(args, ", "))
}