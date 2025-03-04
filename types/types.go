package types

import (
	"fmt"

	"github.com/vibe-lang/vibe/ast"
)

// Type represents a type in the language
type Type interface {
	String() string
}

// IntType represents the integer type
var IntType = SimpleType{"int"}

// FloatType represents the floating point type
var FloatType = SimpleType{"float"}

// StringType represents the string type
var StringType = SimpleType{"string"}

// BoolType represents the boolean type
var BoolType = SimpleType{"bool"}

// NilType represents the nil type
var NilType = SimpleType{"nil"}

// AnyType represents any type
var AnyType = SimpleType{"any"}

// SimpleType represents a basic type
type SimpleType struct {
	Name string
}

func (t SimpleType) String() string {
	return t.Name
}

// ArrayType represents an array type
type ArrayType struct {
	ElementType Type
}

func (t ArrayType) String() string {
	return fmt.Sprintf("Array<%s>", t.ElementType.String())
}

// FunctionType represents a function type
type FunctionType struct {
	ParameterTypes []Type
	ReturnType     Type
}

func (t FunctionType) String() string {
	result := "def("
	for i, param := range t.ParameterTypes {
		if i > 0 {
			result += ", "
		}
		result += param.String()
	}
	result += ") -> " + t.ReturnType.String()
	return result
}

// UnionType represents a union of types
type UnionType struct {
	Types []Type
}

func (t UnionType) String() string {
	result := "union["
	for i, typ := range t.Types {
		if i > 0 {
			result += ", "
		}
		result += typ.String()
	}
	result += "]"
	return result
}

// IsAssignable determines if a value of type src can be assigned to a variable of type dst
func IsAssignable(src, dst Type) bool {
	// Any type can be assigned to any
	if _, ok := dst.(SimpleType); ok && dst.String() == "any" {
		return true
	}

	// Same type is always assignable
	if src.String() == dst.String() {
		return true
	}

	// Nil can be assigned to any non-primitive type
	if _, ok := src.(SimpleType); ok && src.String() == "nil" {
		if _, ok := dst.(SimpleType); ok {
			// Only allow nil to be assigned to specific primitive types
			switch dst.String() {
			case "int", "float", "string", "bool":
				return false
			}
		}
		return true
	}

	// For union types, the source must be assignable to at least one of the union types
	if unionType, ok := dst.(UnionType); ok {
		for _, t := range unionType.Types {
			if IsAssignable(src, t) {
				return true
			}
		}
		return false
	}

	// Number type coercion: int -> float
	if srcSimple, ok := src.(SimpleType); ok && srcSimple.Name == "int" {
		if dstSimple, ok := dst.(SimpleType); ok && dstSimple.Name == "float" {
			return true
		}
	}

	// Array type compatibility
	if srcArray, ok := src.(ArrayType); ok {
		if dstArray, ok := dst.(ArrayType); ok {
			// Check element type compatibility
			return IsAssignable(srcArray.ElementType, dstArray.ElementType)
		}
	}

	// Function type compatibility
	if srcFunc, ok := src.(FunctionType); ok {
		if dstFunc, ok := dst.(FunctionType); ok {
			// Return type must be assignable
			if !IsAssignable(srcFunc.ReturnType, dstFunc.ReturnType) {
				return false
			}

			// Parameter count must match
			if len(srcFunc.ParameterTypes) != len(dstFunc.ParameterTypes) {
				return false
			}

			// Parameters must be assignable in reverse (contravariant)
			for i := range srcFunc.ParameterTypes {
				if !IsAssignable(dstFunc.ParameterTypes[i], srcFunc.ParameterTypes[i]) {
					return false
				}
			}

			return true
		}
	}

	return false
}

// TypeChecker provides type checking functionality
type TypeChecker struct {
	types map[string]Type
}

// NewTypeChecker creates a new type checker
func NewTypeChecker() *TypeChecker {
	return &TypeChecker{
		types: make(map[string]Type),
	}
}

// CheckType verifies that a node's type matches the expected type
func (tc *TypeChecker) CheckType(node ast.Node, expected Type) (Type, error) {
	actual, err := tc.InferType(node)
	if err != nil {
		return nil, err
	}

	if !IsAssignable(actual, expected) {
		return nil, fmt.Errorf("type error: cannot use %s as %s", actual.String(), expected.String())
	}

	return actual, nil
}

// InferType determines the type of a node
func (tc *TypeChecker) InferType(node ast.Node) (Type, error) {
	switch n := node.(type) {
	case *ast.NumberLiteral:
		if n.IsInt {
			return IntType, nil
		}
		return FloatType, nil
	case *ast.StringLiteral:
		return StringType, nil
	case *ast.BooleanLiteral:
		return BoolType, nil
	case *ast.NilLiteral:
		return NilType, nil
	default:
		return AnyType, nil
	}
}
