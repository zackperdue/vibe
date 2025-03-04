package object

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/vibe-lang/vibe/ast"
)

// ObjectType represents the type of an object
type ObjectType string

const (
	INTEGER_OBJ      = "INTEGER"
	FLOAT_OBJ        = "FLOAT"
	BOOLEAN_OBJ      = "BOOLEAN"
	NULL_OBJ         = "NULL"
	RETURN_VALUE_OBJ = "RETURN_VALUE"
	ERROR_OBJ        = "ERROR"
	FUNCTION_OBJ     = "FUNCTION"
	STRING_OBJ       = "STRING"
	BUILTIN_OBJ      = "BUILTIN"
	ARRAY_OBJ        = "ARRAY"
	CLASS_OBJ        = "CLASS"
	OBJECT_OBJ       = "OBJECT"
)

// Object represents an object in our language
type Object interface {
	Type() ObjectType
	Inspect() string
}

// Integer represents an integer object
type Integer struct {
	Value int64
}

func (i *Integer) Type() ObjectType { return INTEGER_OBJ }
func (i *Integer) Inspect() string  { return strconv.FormatInt(i.Value, 10) }

// Float represents a float object
type Float struct {
	Value float64
}

func (f *Float) Type() ObjectType { return FLOAT_OBJ }
func (f *Float) Inspect() string  { return strconv.FormatFloat(f.Value, 'f', -1, 64) }

// Boolean represents a boolean object
type Boolean struct {
	Value bool
}

func (b *Boolean) Type() ObjectType { return BOOLEAN_OBJ }
func (b *Boolean) Inspect() string  { return strconv.FormatBool(b.Value) }

// Null represents a null object
type Null struct{}

func (n *Null) Type() ObjectType { return NULL_OBJ }
func (n *Null) Inspect() string  { return "null" }

// ReturnValue represents a return value
type ReturnValue struct {
	Value Object
}

func (rv *ReturnValue) Type() ObjectType { return RETURN_VALUE_OBJ }
func (rv *ReturnValue) Inspect() string  { return rv.Value.Inspect() }

// Error represents an error
type Error struct {
	Message string
}

func (e *Error) Type() ObjectType { return ERROR_OBJ }
func (e *Error) Inspect() string  { return "ERROR: " + e.Message }

// Function represents a function
type Function struct {
	Parameters []*ast.Parameter
	Body       *ast.BlockStmt
	Env        *Environment
}

func (f *Function) Type() ObjectType { return FUNCTION_OBJ }
func (f *Function) Inspect() string {
	var out strings.Builder
	params := []string{}
	for _, p := range f.Parameters {
		params = append(params, p.String())
	}
	out.WriteString("fn(")
	out.WriteString(strings.Join(params, ", "))
	out.WriteString(") {\n")
	out.WriteString(f.Body.String())
	out.WriteString("\n}")
	return out.String()
}

// String represents a string object
type String struct {
	Value string
}

func (s *String) Type() ObjectType { return STRING_OBJ }
func (s *String) Inspect() string  { return s.Value }

// BuiltinFunction represents a builtin function
type BuiltinFunction struct {
	Fn func(args ...Object) Object
}

func (b *BuiltinFunction) Type() ObjectType { return BUILTIN_OBJ }
func (b *BuiltinFunction) Inspect() string  { return "builtin function" }

// Array represents an array object
type Array struct {
	Elements []Object
}

func (a *Array) Type() ObjectType { return ARRAY_OBJ }
func (a *Array) Inspect() string {
	var out strings.Builder
	elements := []string{}
	for _, e := range a.Elements {
		elements = append(elements, e.Inspect())
	}
	out.WriteString("[")
	out.WriteString(strings.Join(elements, ", "))
	out.WriteString("]")
	return out.String()
}

// Class represents a class object
type Class struct {
	Name       string
	Methods    map[string]*Function
	Properties map[string]Object
}

func (c *Class) Type() ObjectType { return CLASS_OBJ }
func (c *Class) Inspect() string  { return fmt.Sprintf("class %s", c.Name) }

// ObjectInstance represents an object instance
type ObjectInstance struct {
	Class      *Class
	Properties map[string]Object
}

func (o *ObjectInstance) Type() ObjectType { return OBJECT_OBJ }
func (o *ObjectInstance) Inspect() string  { return fmt.Sprintf("%s instance", o.Class.Name) }