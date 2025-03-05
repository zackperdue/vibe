package interpreter

import (
	"fmt"
	"strconv"

	"github.com/vibe-lang/vibe/ast"
	"github.com/vibe-lang/vibe/types"
)

// Value is the interface that all values in the interpreter must implement
type Value interface {
	Type() string
	Inspect() string
	VibeType() types.Type
}

// IntegerValue represents an integer in the interpreter
type IntegerValue struct {
	Value int64
}

func (i *IntegerValue) Type() string { return "INTEGER" }
func (i *IntegerValue) Inspect() string {
	return strconv.FormatInt(i.Value, 10)
}
func (i *IntegerValue) VibeType() types.Type { return types.IntType }

// FloatValue represents a floating point number in the interpreter
type FloatValue struct {
	Value float64
}

func (f *FloatValue) Type() string { return "FLOAT" }
func (f *FloatValue) Inspect() string {
	return strconv.FormatFloat(f.Value, 'f', -1, 64)
}
func (f *FloatValue) VibeType() types.Type { return types.FloatType }

// StringValue represents a string in the interpreter
type StringValue struct {
	Value string
}

func (s *StringValue) Type() string { return "STRING" }
func (s *StringValue) Inspect() string {
	return s.Value
}
func (s *StringValue) VibeType() types.Type { return types.StringType }

// BooleanValue represents a boolean in the interpreter
type BooleanValue struct {
	Value bool
}

func (b *BooleanValue) Type() string { return "BOOLEAN" }
func (b *BooleanValue) Inspect() string {
	return strconv.FormatBool(b.Value)
}
func (b *BooleanValue) VibeType() types.Type { return types.BoolType }

// NilValue represents a nil value in the interpreter
type NilValue struct{}

func (n *NilValue) Type() string       { return "NIL" }
func (n *NilValue) Inspect() string    { return "nil" }
func (n *NilValue) VibeType() types.Type { return types.NilType }

// ReturnValue represents a return value from a function
type ReturnValue struct {
	Value Value
}

func (r *ReturnValue) Type() string { return "RETURN_VALUE" }
func (r *ReturnValue) Inspect() string {
	return r.Value.Inspect()
}
func (r *ReturnValue) VibeType() types.Type {
	return r.Value.VibeType()
}

// FunctionValue represents a function in the interpreter
type FunctionValue struct {
	Name        string
	Parameters  []ast.Parameter
	Body        *ast.BlockStmt
	Env         *Environment
	ReturnType  types.Type
	BuiltinFunc func(args []Value) Value
}

func (f *FunctionValue) Type() string { return "FUNCTION" }
func (f *FunctionValue) Inspect() string {
	return fmt.Sprintf("function %s", f.Name)
}
func (f *FunctionValue) VibeType() types.Type {
	// Directly use the return type
	return types.FunctionType{
		ParameterTypes: []types.Type{types.AnyType}, // Simplified for now
		ReturnType:     f.ReturnType,
	}
}

// ArrayValue represents an array of values
type ArrayValue struct {
	Elements []Value
}

func (a *ArrayValue) Type() string { return "ARRAY" }
func (a *ArrayValue) Inspect() string {
	var out string
	elements := []string{}
	for _, e := range a.Elements {
		elements = append(elements, e.Inspect())
	}
	out = "[" + fmt.Sprintf("%s", elements) + "]"
	return out
}
func (a *ArrayValue) VibeType() types.Type {
	if len(a.Elements) == 0 {
		return types.ArrayType{ElementType: types.AnyType}
	}
	return types.ArrayType{ElementType: a.Elements[0].VibeType()}
}

// Index returns the element at the given index
func (a *ArrayValue) Index(index Value) Value {
	intIndex, ok := index.(*IntegerValue)
	if !ok {
		return &StringValue{Value: fmt.Sprintf("index must be an integer, got %s", index.Type())}
	}

	idx := intIndex.Value
	if idx < 0 || idx >= int64(len(a.Elements)) {
		return &NilValue{}
	}

	return a.Elements[idx]
}

// BuiltinFunction represents a built-in function
type BuiltinFunction struct {
	Name       string
	Fn         func(args []Value) Value
	ParamTypes []types.Type
	ReturnType types.Type
}

func (b *BuiltinFunction) Type() string { return "BUILTIN" }
func (b *BuiltinFunction) Inspect() string { return "builtin function: " + b.Name }
func (b *BuiltinFunction) VibeType() types.Type {
	return types.FunctionType{
		ParameterTypes: b.ParamTypes,
		ReturnType:     b.ReturnType,
	}
}

// ObjectValue represents an object in the interpreter
type ObjectValue struct {
	Class      string
	Properties map[string]Value
	Methods    map[string]*FunctionValue
}

func (o *ObjectValue) Type() string { return "OBJECT" }
func (o *ObjectValue) Inspect() string {
	return fmt.Sprintf("object: %s", o.Class)
}
func (o *ObjectValue) VibeType() types.Type { return types.AnyType } // TODO: Create proper object type

// ClassValue represents a class in the interpreter
type ClassValue struct {
	Name       string
	Properties map[string]types.Type
	Methods    map[string]*FunctionValue
}

func (c *ClassValue) Type() string { return "CLASS" }
func (c *ClassValue) Inspect() string {
	return fmt.Sprintf("class: %s", c.Name)
}
func (c *ClassValue) VibeType() types.Type { return types.AnyType } // TODO: Create proper class type