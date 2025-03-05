package interpreter

import (
	"fmt"

	"github.com/vibe-lang/vibe/types"
)

// Environment wraps the symbol table for variables and functions
type Environment struct {
	store    map[string]Value
	types    map[string]types.Type
	outer    *Environment
	builtins map[string]*BuiltinFunction
}

// NewEnvironment creates a new environment
func NewEnvironment() *Environment {
	s := make(map[string]Value)
	t := make(map[string]types.Type)
	b := make(map[string]*BuiltinFunction)
	return &Environment{store: s, types: t, builtins: b, outer: nil}
}

// NewEnclosedEnvironment creates a new environment with an outer environment
func NewEnclosedEnvironment(outer *Environment) *Environment {
	env := NewEnvironment()
	env.outer = outer
	// Copy builtins from outer
	for name, builtin := range outer.builtins {
		env.builtins[name] = builtin
	}
	return env
}

// Get retrieves a value from the environment
func (e *Environment) Get(name string) (Value, bool) {
	// Check for builtins first
	if builtin, ok := e.builtins[name]; ok {
		return builtin, true
	}

	// Then check variables
	obj, ok := e.store[name]
	if !ok && e.outer != nil {
		obj, ok = e.outer.Get(name)
	}
	return obj, ok
}

// Set sets a value in the environment
func (e *Environment) Set(name string, val Value) error {
	// Check if a value with this name already exists and has a type
	existingType, hasType := e.types[name]
	if hasType {
		// Validate that the new value is compatible with the type
		if !types.IsAssignable(val.VibeType(), existingType) {
			return fmt.Errorf("Type error: Cannot assign value of type %s to variable %s of type %s",
				val.VibeType().String(), name, existingType.String())
		}
	}

	e.store[name] = val
	return nil
}

// SetWithType sets a value with a type annotation
func (e *Environment) SetWithType(name string, val Value, typ types.Type) error {
	// Validate that the value is compatible with the type
	if !types.IsAssignable(val.VibeType(), typ) {
		return fmt.Errorf("Type error: Cannot assign value of type %s to variable %s of type %s",
			val.VibeType().String(), name, typ.String())
	}

	e.store[name] = val
	e.types[name] = typ
	return nil
}

// RegisterBuiltin registers a built-in function
func (e *Environment) RegisterBuiltin(name string, fn func(args []Value) Value, paramTypes []types.Type, returnType types.Type) {
	builtin := &BuiltinFunction{
		Fn:         fn,
		ParamTypes: paramTypes,
		ReturnType: returnType,
	}
	e.builtins[name] = builtin
}