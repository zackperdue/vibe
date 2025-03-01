package types

import "fmt"

// Type represents a type in the Crystal language
type Type interface {
	String() string
	Equal(Type) bool
}

// BasicType represents primitive types
type BasicType string

const (
	IntType    BasicType = "int"
	FloatType  BasicType = "float"
	StringType BasicType = "string"
	BoolType   BasicType = "bool"
	NilType    BasicType = "nil"
	AnyType    BasicType = "any"
)

func (b BasicType) String() string {
	return string(b)
}

func (b BasicType) Equal(other Type) bool {
	if otherBasic, ok := other.(BasicType); ok {
		return b == otherBasic
	}
	return false
}

// ArrayType represents array types
type ArrayType struct {
	ElementType Type
}

func (a ArrayType) String() string {
	return fmt.Sprintf("Array<%s>", a.ElementType.String())
}

func (a ArrayType) Equal(other Type) bool {
	if otherArray, ok := other.(ArrayType); ok {
		return a.ElementType.Equal(otherArray.ElementType)
	}
	return false
}

// FunctionType represents function types
type FunctionType struct {
	ParameterTypes []Type
	ReturnType     Type
}

func (f FunctionType) String() string {
	params := ""
	for i, paramType := range f.ParameterTypes {
		if i > 0 {
			params += ", "
		}
		params += paramType.String()
	}
	return fmt.Sprintf("(%s) => %s", params, f.ReturnType.String())
}

func (f FunctionType) Equal(other Type) bool {
	if otherFunc, ok := other.(FunctionType); ok {
		if len(f.ParameterTypes) != len(otherFunc.ParameterTypes) {
			return false
		}
		for i, paramType := range f.ParameterTypes {
			if !paramType.Equal(otherFunc.ParameterTypes[i]) {
				return false
			}
		}
		return f.ReturnType.Equal(otherFunc.ReturnType)
	}
	return false
}

// ObjectType represents object types with properties
type ObjectType struct {
	Properties map[string]Type
}

func (o ObjectType) String() string {
	result := "{"
	i := 0
	for name, propType := range o.Properties {
		if i > 0 {
			result += ", "
		}
		result += fmt.Sprintf("%s: %s", name, propType.String())
		i++
	}
	result += "}"
	return result
}

func (o ObjectType) Equal(other Type) bool {
	if otherObj, ok := other.(ObjectType); ok {
		if len(o.Properties) != len(otherObj.Properties) {
			return false
		}
		for name, propType := range o.Properties {
			otherPropType, ok := otherObj.Properties[name]
			if !ok || !propType.Equal(otherPropType) {
				return false
			}
		}
		return true
	}
	return false
}

// UnionType represents union types (A | B)
type UnionType struct {
	Types []Type
}

func (u UnionType) String() string {
	result := ""
	for i, t := range u.Types {
		if i > 0 {
			result += " | "
		}
		result += t.String()
	}
	return result
}

func (u UnionType) Equal(other Type) bool {
	if otherUnion, ok := other.(UnionType); ok {
		if len(u.Types) != len(otherUnion.Types) {
			return false
		}
		// This is a simplistic implementation; a more accurate one would check
		// all permutations since union types can be in any order
		for i, t := range u.Types {
			if !t.Equal(otherUnion.Types[i]) {
				return false
			}
		}
		return true
	}
	return false
}

// TypeChecker handles type checking operations
type TypeChecker struct {
	// Add fields as needed
}

// NewTypeChecker creates a new TypeChecker
func NewTypeChecker() *TypeChecker {
	return &TypeChecker{}
}

// IsAssignable checks if a value of fromType can be assigned to a variable of toType
func IsAssignable(fromType Type, toType Type) bool {
	// Any type can be assigned to Any
	if toType == AnyType {
		return true
	}

	// Exact match
	if fromType.Equal(toType) {
		return true
	}

	// Number type conversions
	if fromType == IntType && toType == FloatType {
		return true
	}

	// Nil can be assigned to any object/array/function type
	if fromType == NilType {
		switch toType.(type) {
		case ArrayType, ObjectType, FunctionType:
			return true
		}
	}

	// Union type compatibilities
	if unionType, ok := toType.(UnionType); ok {
		for _, t := range unionType.Types {
			if IsAssignable(fromType, t) {
				return true
			}
		}
	}

	if unionType, ok := fromType.(UnionType); ok {
		// All types in the union must be assignable to the target type
		for _, t := range unionType.Types {
			if !IsAssignable(t, toType) {
				return false
			}
		}
		return true
	}

	return false
}
