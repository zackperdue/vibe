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
	fmt.Printf("DEBUG ENV: Looking up variable '%s' in environment\n", name)

	// First check the current environment
	obj, ok := e.store[name]
	if ok {
		fmt.Printf("DEBUG ENV: ✅ Found variable '%s' = %s in current environment\n",
			name, obj.Inspect())
		return obj, true
	}

	// Check if it's a built-in function
	if builtin, ok := e.builtins[name]; ok {
		fmt.Printf("DEBUG ENV: ✅ Found built-in function '%s'\n", name)
		return builtin, true
	}

	// If not found and we have an outer environment, look there
	if e.outer != nil {
		fmt.Printf("DEBUG ENV: Variable '%s' not found in current environment, checking outer\n", name)
		obj, ok := e.outer.Get(name)
		if ok {
			fmt.Printf("DEBUG ENV: ✅ Found variable '%s' = %s in outer environment\n",
				name, obj.Inspect())
		} else {
			fmt.Printf("DEBUG ENV: ⚠️ Variable '%s' not found in outer environment\n", name)
		}
		return obj, ok
	}

	fmt.Printf("DEBUG ENV: ⚠️ Variable '%s' not found in any environment\n", name)
	return nil, false
}

// Set sets a value in the environment
func (e *Environment) Set(name string, val Value) error {
	fmt.Printf("DEBUG ENV: Setting variable '%s' = %s (type: %s) in environment\n",
		name, val.Inspect(), val.Type())

	// Check if the value is nil
	if val == nil {
		fmt.Printf("DEBUG ENV: ⚠️ Attempted to set nil value for '%s'\n", name)
		return fmt.Errorf("cannot set nil value for %s", name)
	}

	// Set the value in the current environment
	e.store[name] = val

	// Also set the type information (inferred from the value)
	e.types[name] = val.VibeType()

	fmt.Printf("DEBUG ENV: ✅ Variable '%s' set in environment with type %s\n",
		name, val.VibeType().String())
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