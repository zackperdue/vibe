package ast

import (
	"fmt"
	"strings"
)

// Parameter represents a function parameter in the AST
type Parameter struct {
	Name string
	Type *TypeAnnotation
}

func (p *Parameter) String() string {
	if p.Type == nil {
		return p.Name
	}
	return fmt.Sprintf("%s: %s", p.Name, p.Type.String())
}

// FunctionDef represents a function definition in the AST
type FunctionDef struct {
	Name       string
	Parameters []Parameter
	ReturnType *TypeAnnotation
	Body       *BlockStmt
}

func (f *FunctionDef) Type() NodeType { return FunctionDefNode }
func (f *FunctionDef) String() string {
	params := []string{}
	for _, p := range f.Parameters {
		params = append(params, p.String())
	}

	returnTypeStr := ""
	if f.ReturnType != nil {
		returnTypeStr = ": " + f.ReturnType.String()
	}

	return fmt.Sprintf("def %s(%s)%s %s",
		f.Name,
		strings.Join(params, ", "),
		returnTypeStr,
		f.Body.String())
}

// TypeAnnotation represents a type annotation in the AST
type TypeAnnotation struct {
	TypeName    string
	GenericType *TypeAnnotation
	TypeParams  []Node // For generic types like Array<string>
}

func (t *TypeAnnotation) Type() NodeType { return TypeAnnotationNode }
func (t *TypeAnnotation) String() string {
	if t.GenericType != nil {
		return fmt.Sprintf("%s<%s>", t.TypeName, t.GenericType.String())
	}

	if len(t.TypeParams) > 0 {
		params := []string{}
		for _, p := range t.TypeParams {
			params = append(params, p.String())
		}
		return fmt.Sprintf("%s<%s>", t.TypeName, strings.Join(params, ", "))
	}

	return t.TypeName
}

// TypeDeclaration represents a type declaration in the AST
type TypeDeclaration struct {
	Name      string
	TypeValue Node // Could be a TypeAnnotation or another structure
}

func (t *TypeDeclaration) Type() NodeType { return TypeDeclarationNode }
func (t *TypeDeclaration) String() string {
	return fmt.Sprintf("type %s = %s", t.Name, t.TypeValue.String())
}

// VariableDecl represents a variable declaration in the AST
type VariableDecl struct {
	Name           string
	TypeAnnotation *TypeAnnotation
	Value          Node // Initial value (can be nil)
}

func (v *VariableDecl) Type() NodeType { return VariableDeclNode }
func (v *VariableDecl) String() string {
	typeAnnotation := ""
	if v.TypeAnnotation != nil {
		typeAnnotation = ": " + v.TypeAnnotation.String()
	}

	if v.Value == nil {
		return fmt.Sprintf("%s%s", v.Name, typeAnnotation)
	}

	return fmt.Sprintf("%s%s = %s", v.Name, typeAnnotation, v.Value.String())
}

// ClassDef represents a class definition in the AST
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

func (c *ClassDef) Type() NodeType { return ClassDefNode }
func (c *ClassDef) String() string {
	var result string

	if c.Parent != "" {
		result = fmt.Sprintf("class %s inherits %s", c.Name, c.Parent)
	} else {
		result = fmt.Sprintf("class %s", c.Name)
	}

	if len(c.TypeParams) > 0 {
		result += fmt.Sprintf("<%s>", strings.Join(c.TypeParams, ", "))
	}

	methods := []string{}
	for _, method := range c.Methods {
		methods = append(methods, method.String())
	}

	if len(methods) > 0 {
		result += " { " + strings.Join(methods, "; ") + " }"
	} else {
		result += " {}"
	}

	return result
}