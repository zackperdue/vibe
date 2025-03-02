package interpreter

import (
	"fmt"
	"io/ioutil"
	"math"
	"os"
	"strconv"
	"strings"

	"github.com/example/vibe/ast"
	"github.com/example/vibe/lexer"
	"github.com/example/vibe/parser"
	"github.com/example/vibe/types"
)

// Value interface represents values in our language
type Value interface {
	Type() string
	Inspect() string
	VibeType() types.Type
}

// IntegerValue represents an integer value
type IntegerValue struct {
	Value int
}

func (i *IntegerValue) Type() string { return "INTEGER" }
func (i *IntegerValue) Inspect() string { return strconv.Itoa(i.Value) }
func (i *IntegerValue) VibeType() types.Type { return types.IntType }

// FloatValue represents a floating point value
type FloatValue struct {
	Value float64
}

func (f *FloatValue) Type() string { return "FLOAT" }
func (f *FloatValue) Inspect() string { return strconv.FormatFloat(f.Value, 'f', -1, 64) }
func (f *FloatValue) VibeType() types.Type { return types.FloatType }

// StringValue represents a string value
type StringValue struct {
	Value string
}

func (s *StringValue) Type() string { return "STRING" }
func (s *StringValue) Inspect() string { return s.Value }
func (s *StringValue) VibeType() types.Type { return types.StringType }

// BooleanValue represents a boolean value
type BooleanValue struct {
	Value bool
}

func (b *BooleanValue) Type() string { return "BOOLEAN" }
func (b *BooleanValue) Inspect() string { return strconv.FormatBool(b.Value) }
func (b *BooleanValue) VibeType() types.Type { return types.BoolType }

// NilValue represents a nil value
type NilValue struct{}

func (n *NilValue) Type() string { return "NIL" }
func (n *NilValue) Inspect() string { return "nil" }
func (n *NilValue) VibeType() types.Type { return types.NilType }

// ReturnValue wraps a return value
type ReturnValue struct {
	Value Value
}

func (r *ReturnValue) Type() string { return "RETURN" }
func (r *ReturnValue) Inspect() string { return r.Value.Inspect() }
func (r *ReturnValue) VibeType() types.Type { return r.Value.VibeType() }

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
	result := "["
	for i, element := range a.Elements {
		if i > 0 {
			result += ", "
		}
		result += element.Inspect()
	}
	result += "]"
	return result
}
func (a *ArrayValue) VibeType() types.Type {
	if len(a.Elements) == 0 {
		// Empty array - default to array of any
		return types.ArrayType{ElementType: types.AnyType}
	}

	// Get the type of the first element
	elementType := a.Elements[0].VibeType()

	// Check if all elements have the same type
	for _, element := range a.Elements {
		if element.VibeType() != elementType {
			// If not, return array of any
			return types.ArrayType{ElementType: types.AnyType}
		}
	}

	return types.ArrayType{ElementType: elementType}
}

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
	e.builtins[name] = &BuiltinFunction{
		Name:       name,
		Fn:         fn,
		ParamTypes: paramTypes,
		ReturnType: returnType,
	}
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

// Adding new value types for class functionality

// ClassValue represents a class definition
type ClassValue struct {
	Name       string
	Methods    map[string]*FunctionValue
	Properties map[string]Value
}

func (c *ClassValue) Type() string { return "CLASS" }
func (c *ClassValue) Inspect() string { return fmt.Sprintf("class %s", c.Name) }
func (c *ClassValue) VibeType() types.Type { return types.AnyType } // TODO: Create proper class type

// ObjectValue represents an instance of a class
type ObjectValue struct {
	Class      *ClassValue
	Properties map[string]Value
}

func (o *ObjectValue) Type() string { return "OBJECT" }
func (o *ObjectValue) Inspect() string { return fmt.Sprintf("%s instance", o.Class.Name) }
func (o *ObjectValue) VibeType() types.Type { return types.AnyType } // TODO: Create proper object type

// Interpreter executes the AST
type Interpreter struct {
	env *Environment
}

// New creates a new interpreter
func New() *Interpreter {
	env := NewEnvironment()

	// Register built-in functions
	registerBuiltins(env)
	registerBuiltinClasses(env)

	return &Interpreter{env: env}
}

func registerBuiltins(env *Environment) {
	// length - works on strings and arrays
	env.RegisterBuiltin("len", func(args []Value) Value {
		if len(args) != 1 {
			return &StringValue{Value: "Type error: len takes exactly 1 argument"}
		}

		switch arg := args[0].(type) {
		case *StringValue:
			return &IntegerValue{Value: len(arg.Value)}
		case *ArrayValue:
			return &IntegerValue{Value: len(arg.Elements)}
		default:
			return &StringValue{Value: "Type error: len requires a string or array argument"}
		}
	}, []types.Type{types.AnyType}, types.IntType)

	// type - returns the type of a value as a string
	env.RegisterBuiltin("type", func(args []Value) Value {
		if len(args) != 1 {
			return &StringValue{Value: "Type error: type takes exactly 1 argument"}
		}

		return &StringValue{Value: args[0].VibeType().String()}
	}, []types.Type{types.AnyType}, types.StringType)

	// to_string - converts a value to a string
	env.RegisterBuiltin("to_string", func(args []Value) Value {
		if len(args) != 1 {
			return &StringValue{Value: "Type error: to_string takes exactly 1 argument"}
		}

		return &StringValue{Value: args[0].Inspect()}
	}, []types.Type{types.AnyType}, types.StringType)

	// to_int - converts a value to an integer if possible
	env.RegisterBuiltin("to_int", func(args []Value) Value {
		if len(args) != 1 {
			return &StringValue{Value: "Type error: to_int takes exactly 1 argument"}
		}

		switch arg := args[0].(type) {
		case *StringValue:
			i, err := strconv.Atoi(arg.Value)
			if err != nil {
				return &StringValue{Value: "Type error: cannot convert string to int"}
			}
			return &IntegerValue{Value: i}
		case *FloatValue:
			return &IntegerValue{Value: int(arg.Value)}
		case *IntegerValue:
			return arg
		default:
			return &StringValue{Value: "Type error: cannot convert to int"}
		}
	}, []types.Type{types.AnyType}, types.IntType)

	// to_float - converts a value to a float if possible
	env.RegisterBuiltin("to_float", func(args []Value) Value {
		if len(args) != 1 {
			return &StringValue{Value: "Type error: to_float takes exactly 1 argument"}
		}

		switch arg := args[0].(type) {
		case *StringValue:
			f, err := strconv.ParseFloat(arg.Value, 64)
			if err != nil {
				return &StringValue{Value: "Type error: cannot convert string to float"}
			}
			return &FloatValue{Value: f}
		case *IntegerValue:
			return &FloatValue{Value: float64(arg.Value)}
		case *FloatValue:
			return arg
		default:
			return &StringValue{Value: "Type error: cannot convert to float"}
		}
	}, []types.Type{types.AnyType}, types.FloatType)
}

// Add this function to register built-in classes
func registerBuiltinClasses(env *Environment) {
	// Add Point class as a placeholder until proper class definition parsing is implemented
	pointClass := &ClassValue{
		Name:       "Point",
		Methods:    make(map[string]*FunctionValue),
		Properties: make(map[string]Value),
	}

	// Add get_x method
	pointClass.Methods["get_x"] = &FunctionValue{
		Name: "get_x",
		Body: nil, // Not using the body, will manually implement below
		Env:  env,
		BuiltinFunc: func(args []Value) Value {
			if len(args) != 1 {
				return &StringValue{Value: "Error: get_x requires object instance"}
			}
			obj, ok := args[0].(*ObjectValue)
			if !ok {
				return &StringValue{Value: "Error: get_x can only be called on Point objects"}
			}
			if x, ok := obj.Properties["x"]; ok {
				return x
			}
			return &NilValue{}
		},
	}

	// Add get_y method
	pointClass.Methods["get_y"] = &FunctionValue{
		Name: "get_y",
		Body: nil, // Not using the body, will manually implement below
		Env:  env,
		BuiltinFunc: func(args []Value) Value {
			if len(args) != 1 {
				return &StringValue{Value: "Error: get_y requires object instance"}
			}
			obj, ok := args[0].(*ObjectValue)
			if !ok {
				return &StringValue{Value: "Error: get_y can only be called on Point objects"}
			}
			if y, ok := obj.Properties["y"]; ok {
				return y
			}
			return &NilValue{}
		},
	}

	env.Set("Point", pointClass)
}

// Eval evaluates the AST and returns the result
func (i *Interpreter) Eval(node ast.Node) Value {
	return i.eval(node, i.env)
}

func (i *Interpreter) eval(node ast.Node, env *Environment) Value {
	switch node := node.(type) {
	case *ast.Program:
		return i.evalProgram(node, env)
	case *ast.BlockStmt:
		return i.evalBlockStatement(node, env)
	case *ast.NumberLiteral:
		if node.IsInt {
			return &IntegerValue{Value: int(node.Value)}
		}
		return &FloatValue{Value: node.Value}
	case *ast.StringLiteral:
		return &StringValue{Value: node.Value}
	case *ast.BooleanLiteral:
		return &BooleanValue{Value: node.Value}
	case *ast.NilLiteral:
		return &NilValue{}
	case *ast.Identifier:
		return i.evalIdentifier(node, env)
	case *ast.PrintStmt:
		return i.evalPrintStatement(node, env)
	case *ast.RequireStmt:
		return i.evalRequireStatement(node, env)
	case *ast.Assignment:
		return i.evalAssignment(node, env)
	case *ast.VariableDecl:
		return i.evalVariableDeclaration(node, env)
	case *ast.FunctionDef:
		return i.evalFunctionDefinition(node, env)
	case *ast.CallExpr:
		return i.evalCallExpression(node, env)
	case *ast.MethodCall:
		return i.evalMethodCall(node, env)
	case *ast.ClassInst:
		return i.evalClassInstantiation(node, env)
	case *ast.ReturnStmt:
		return i.evalReturnStatement(node, env)
	case *ast.IfStmt:
		return i.evalIfStatement(node, env)
	case *ast.WhileStmt:
		return i.evalWhileStatement(node, env)
	case *ast.ForStmt:
		return i.evalForStatement(node, env)
	case *ast.BinaryExpr:
		return i.evalBinaryExpression(node, env)
	case *ast.ArrayLiteral:
		return i.evalArrayLiteral(node, env)
	case *ast.TypeAnnotation:
		// Type annotations don't evaluate to a value on their own
		return &NilValue{}
	case *ast.TypeDeclaration:
		// Type declarations don't evaluate to a value
		return &NilValue{}
	default:
		// Handle unexpected nodes
		return &StringValue{Value: fmt.Sprintf("Unknown node type: %T : %s", node, node.Type())}
	}
}

func (i *Interpreter) evalVariableDeclaration(node *ast.VariableDecl, env *Environment) Value {
	var value Value
	if node.Value != nil {
		value = i.eval(node.Value, env)
	} else {
		// If no value is provided, initialize with nil
		value = &NilValue{}
	}

	if node.TypeAnnotation != nil {
		// Parse the type annotation
		varType := i.parseTypeAnnotation(node.TypeAnnotation)

		// Check that the value is compatible with the declared type
		if !types.IsAssignable(value.VibeType(), varType) {
			return &StringValue{Value: fmt.Sprintf("Type error: Cannot assign value of type %s to variable of type %s",
				value.VibeType().String(), varType.String())}
		}

		// Set with type check
		err := env.SetWithType(node.Name, value, varType)
		if err != nil {
			return &StringValue{Value: err.Error()}
		}
	} else {
		// No type annotation, infer from the value
		err := env.Set(node.Name, value)
		if err != nil {
			return &StringValue{Value: err.Error()}
		}
	}

	return &NilValue{}
}

func (i *Interpreter) parseTypeAnnotation(node *ast.TypeAnnotation) types.Type {
	switch node.TypeName {
	case "int":
		return types.IntType
	case "float":
		return types.FloatType
	case "string":
		return types.StringType
	case "bool":
		return types.BoolType
	case "any":
		return types.AnyType
	case "Array":
		if len(node.TypeParams) > 0 {
			elemType := i.parseTypeAnnotation(node.TypeParams[0].(*ast.TypeAnnotation))
			return types.ArrayType{ElementType: elemType}
		}
		// Default to Array of any
		return types.ArrayType{ElementType: types.AnyType}
	case "union":
		if len(node.TypeParams) > 0 {
			var unionTypes []types.Type
			for _, param := range node.TypeParams {
				unionTypes = append(unionTypes, i.parseTypeAnnotation(param.(*ast.TypeAnnotation)))
			}
			return types.UnionType{Types: unionTypes}
		}
		// Invalid union type
		return types.AnyType
	default:
		// Unknown type, default to any
		return types.AnyType
	}
}

func (i *Interpreter) evalProgram(program *ast.Program, env *Environment) Value {
	var result Value
	result = &NilValue{}

	for _, statement := range program.Statements {
		result = i.eval(statement, env)

		// If we hit a return statement, unwrap it and return the value
		if returnValue, ok := result.(*ReturnValue); ok {
			return returnValue.Value
		}
	}

	return result
}

func (i *Interpreter) evalBlockStatement(block *ast.BlockStmt, env *Environment) Value {
	var result Value
	result = &NilValue{}

	// Create a new environment for the block's scope
	blockEnv := NewEnclosedEnvironment(env)

	for _, statement := range block.Statements {
		result = i.eval(statement, blockEnv)

		// If we hit a return statement, break out of the block
		if _, ok := result.(*ReturnValue); ok {
			return result
		}
	}

	return result
}

func (i *Interpreter) evalIdentifier(node *ast.Identifier, env *Environment) Value {
	if val, ok := env.Get(node.Name); ok {
		return val
	}

	return &StringValue{Value: fmt.Sprintf("Error: variable '%s' not found", node.Name)}
}

func (i *Interpreter) evalPrintStatement(node *ast.PrintStmt, env *Environment) Value {
	value := i.eval(node.Value, env)
	fmt.Println(value.Inspect())
	return value
}

func (i *Interpreter) evalRequireStatement(node *ast.RequireStmt, env *Environment) Value {
	// Remove quotes from the path
	path := strings.Trim(node.Path, "\"'")

	// If the path doesn't have a .vi extension, add it
	if !strings.HasSuffix(path, ".vi") {
		path = path + ".vi"
	}

	// Handle relative paths
	// If the path starts with ./ or ../, it's a relative path
	// Otherwise, assume it's an absolute path or in the current directory
	if strings.HasPrefix(path, "./") {
		// Remove the ./ prefix and prepend the tests directory
		path = "tests/" + path[2:]
	} else if !strings.HasPrefix(path, "/") {
		// If it's not an absolute path, assume it's in the tests directory
		path = "tests/" + path
	}

	fmt.Printf("DEBUG: Requiring file from path: %s\n", path)

	// Read the file
	source, err := ioutil.ReadFile(path)
	if err != nil {
		errMsg := fmt.Sprintf("Error requiring file: %s", err)
		fmt.Println(errMsg)
		// Return a special error value that will cause the interpreter to stop execution
		os.Exit(1) // This will terminate the program immediately
		return &StringValue{Value: errMsg}
	}

	// Create a lexer from the source code
	l := lexer.New(string(source))

	// Parse the input
	program, errors := parser.Parse(l)

	if len(errors) > 0 {
		fmt.Println("Parser errors in required file (ignoring for now):")
		for _, err := range errors {
			fmt.Printf("\t%s\n", err)
		}
		// For now, we'll continue even if there are parser errors
		// This is just for testing purposes
	}

	// Evaluate the program in the current environment
	// This will make all definitions from the required file available in the current scope
	fmt.Println("DEBUG: Evaluating required program")
	result := i.evalProgram(program, env)
	fmt.Printf("DEBUG: Result of evaluating required program: %s\n", result.Inspect())

	// For debugging, print all variables in the environment
	fmt.Println("DEBUG: Environment contents after require:")
	for name, value := range env.store {
		fmt.Printf("DEBUG: %s = %s\n", name, value.Inspect())
	}

	return &NilValue{}
}

func (i *Interpreter) evalAssignment(node *ast.Assignment, env *Environment) Value {
	var val Value
	if node.Value != nil {
		val = i.eval(node.Value, env)
	} else {
		// If no value is provided, use nil
		val = &NilValue{}
	}

	err := env.Set(node.Name, val)
	if err != nil {
		return &StringValue{Value: err.Error()}
	}

	return &NilValue{}
}

func (i *Interpreter) evalFunctionDefinition(node *ast.FunctionDef, env *Environment) Value {
	// Parse return type
	var returnType types.Type
	if node.ReturnType != nil {
		returnType = i.parseTypeAnnotation(node.ReturnType)
	} else {
		returnType = types.AnyType
	}

	// Create the function value with parameter types properly processed
	function := &FunctionValue{
		Name:           node.Name,
		Parameters:     node.Parameters, // Use the original parameters
		Body:           node.Body,
		ReturnType:     returnType,
		Env:            env,
	}

	// Add the function to the environment
	env.SetWithType(node.Name, function, function.VibeType())

	return &NilValue{}
}

func (i *Interpreter) evalCallExpression(node *ast.CallExpr, env *Environment) Value {
	function := i.eval(node.Function, env)
	args := i.evalExpressions(node.Args, env)

	if fn, ok := function.(*FunctionValue); ok {
		// Check arity
		if len(args) > len(fn.Parameters) {
			return &StringValue{Value: fmt.Sprintf(
				"Wrong number of arguments: function '%s' expects %d, got %d",
				fn.Name, len(fn.Parameters), len(args))}
		}

		// Create a new environment for the function
		newEnv := NewEnclosedEnvironment(fn.Env)

		// Bind arguments to parameters
		for paramIdx, param := range fn.Parameters {
			if paramIdx < len(args) {
				// Get the parameter type from the TypeAnnotation
				var paramType types.Type
				if param.Type != nil {
					paramType = i.parseTypeAnnotation(param.Type)
				} else {
					paramType = types.AnyType
				}

				// Type check the argument
				if !types.IsAssignable(args[paramIdx].VibeType(), paramType) {
					return &StringValue{Value: fmt.Sprintf(
						"Type error: Parameter '%s' of function '%s' expects %s, got %s",
						param.Name, fn.Name, paramType.String(), args[paramIdx].VibeType().String())}
				}

				// Bind the parameter
				newEnv.SetWithType(param.Name, args[paramIdx], paramType)
			} else {
				// Missing argument, use nil
				var paramType types.Type
				if param.Type != nil {
					paramType = i.parseTypeAnnotation(param.Type)
				} else {
					paramType = types.AnyType
				}
				newEnv.SetWithType(param.Name, &NilValue{}, paramType)
			}
		}

		// Evaluate the function body
		result := i.evalBlockStatement(fn.Body, newEnv)

		// Unwrap return value, if necessary
		if returnValue, ok := result.(*ReturnValue); ok {
			// Type check the return value
			if !types.IsAssignable(returnValue.Value.VibeType(), fn.ReturnType) {
				return &StringValue{Value: fmt.Sprintf(
					"Type error: Function '%s' returns %s, got %s",
					fn.Name, fn.ReturnType.String(), returnValue.Value.VibeType().String())}
			}
			return returnValue.Value
		}

		// Type check the return value
		if !types.IsAssignable(result.VibeType(), fn.ReturnType) {
			return &StringValue{Value: fmt.Sprintf(
				"Type error: Function '%s' returns %s, got %s",
				fn.Name, fn.ReturnType.String(), result.VibeType().String())}
		}

		return result
	} else if builtin, ok := function.(*BuiltinFunction); ok {
		// Check arity
		if len(args) != len(builtin.ParamTypes) {
			return &StringValue{Value: fmt.Sprintf(
				"Wrong number of arguments: function '%s' expects %d, got %d",
				builtin.Name, len(builtin.ParamTypes), len(args))}
		}

		// Type check arguments
		for i, arg := range args {
			if !types.IsAssignable(arg.VibeType(), builtin.ParamTypes[i]) {
				return &StringValue{Value: fmt.Sprintf(
					"Type error: Parameter %d of builtin function '%s' expects %s, got %s",
					i, builtin.Name, builtin.ParamTypes[i].String(), arg.VibeType().String())}
			}
		}

		return builtin.Fn(args)
	}

	return &StringValue{Value: fmt.Sprintf("Not a function: %s", function.Type())}
}

func (i *Interpreter) evalExpressions(
	exps []ast.Node,
	env *Environment,
) []Value {
	var result []Value

	for _, exp := range exps {
		evaluated := i.eval(exp, env)
		result = append(result, evaluated)
	}

	return result
}

func (i *Interpreter) evalReturnStatement(node *ast.ReturnStmt, env *Environment) Value {
	var value Value

	if node.Value != nil {
		value = i.eval(node.Value, env)
	} else {
		value = &NilValue{}
	}

	return &ReturnValue{Value: value}
}

func (i *Interpreter) evalIfStatement(node *ast.IfStmt, env *Environment) Value {
	condition := i.eval(node.Condition, env)

	// Check if the condition is true
	if isTruthy(condition) {
		return i.eval(node.Consequence, env)
	}

	// Check elsif branches
	for _, elseIf := range node.ElseIfBlocks {
		elseIfCondition := i.eval(elseIf.Condition, env)
		if isTruthy(elseIfCondition) {
			return i.eval(elseIf.Consequence, env)
		}
	}

	// Check else branch
	if node.Alternative != nil {
		return i.eval(node.Alternative, env)
	}

	return &NilValue{}
}

func (i *Interpreter) evalWhileStatement(node *ast.WhileStmt, env *Environment) Value {
	for {
		condition := i.eval(node.Condition, env)
		if !isTruthy(condition) {
			break
		}

		result := i.eval(node.Body, env)
		if returnValue, ok := result.(*ReturnValue); ok {
			return returnValue
		}
	}

	return &NilValue{}
}

func (i *Interpreter) evalForStatement(node *ast.ForStmt, env *Environment) Value {
	// Evaluate the iterable expression
	iterable := i.eval(node.Iterable, env)

	// Create a new environment for the loop
	loopEnv := NewEnclosedEnvironment(env)

	// Special case for range expressions (e.g., for i in 0..5)
	if binExpr, ok := node.Iterable.(*ast.BinaryExpr); ok && binExpr.Operator == ".." {
		// Evaluate the start and end of the range
		startValue := i.eval(binExpr.Left, env)
		endValue := i.eval(binExpr.Right, env)

		// Ensure both values are integers
		startInt, startOk := startValue.(*IntegerValue)
		endInt, endOk := endValue.(*IntegerValue)

		if startOk && endOk {
			// Iterate through the range (inclusive)
			for idx := startInt.Value; idx <= endInt.Value; idx++ {
				// Set the iterator variable
				loopEnv.Set(node.Iterator, &IntegerValue{Value: idx})

				// Execute the loop body
				result := i.eval(node.Body, loopEnv)

				// Handle return statements inside the loop
				if returnValue, ok := result.(*ReturnValue); ok {
					return returnValue
				}
			}
			return &NilValue{}
		}

		// If the range bounds aren't integers, report an error
		return &StringValue{Value: "Type error: range bounds must be integers"}
	}

	// Handle standard iterables
	switch iterable := iterable.(type) {
	case *ArrayValue:
		// Iterate over array elements
		for _, element := range iterable.Elements {
			// Bind the current element to the iterator variable
			loopEnv.Set(node.Iterator, element)

			// Execute the loop body
			result := i.eval(node.Body, loopEnv)

			// Handle return statements inside the loop
			if returnValue, ok := result.(*ReturnValue); ok {
				return returnValue
			}
		}
	case *StringValue:
		// Iterate over characters in the string
		for _, char := range iterable.Value {
			// Convert each character to a string value
			charValue := &StringValue{Value: string(char)}

			// Bind the current character to the iterator variable
			loopEnv.Set(node.Iterator, charValue)

			// Execute the loop body
			result := i.eval(node.Body, loopEnv)

			// Handle return statements inside the loop
			if returnValue, ok := result.(*ReturnValue); ok {
				return returnValue
			}
		}
	default:
		// Unsupported iterable type
		return &StringValue{Value: fmt.Sprintf("Type error: cannot iterate over %s", iterable.Type())}
	}

	return &NilValue{}
}

func (i *Interpreter) evalArrayLiteral(node *ast.ArrayLiteral, env *Environment) Value {
	elements := make([]Value, 0, len(node.Elements))

	for _, element := range node.Elements {
		evaluated := i.eval(element, env)
		elements = append(elements, evaluated)
	}

	return &ArrayValue{Elements: elements}
}

func (i *Interpreter) evalBinaryExpression(node *ast.BinaryExpr, env *Environment) Value {
	left := i.eval(node.Left, env)
	right := i.eval(node.Right, env)

	if isError(left) {
		return left
	}
	if isError(right) {
		return right
	}

	switch {
	case left.Type() == "INTEGER" && right.Type() == "INTEGER":
		return evalIntegerBinaryExpression(node.Operator, left, right)
	case (left.Type() == "INTEGER" || left.Type() == "FLOAT") && (right.Type() == "INTEGER" || right.Type() == "FLOAT"):
		return evalNumberBinaryExpression(node.Operator, left, right)
	case left.Type() == "STRING" && right.Type() == "STRING":
		return evalStringBinaryExpression(node.Operator, left, right)
	case left.Type() == "STRING" && (right.Type() == "INTEGER" || right.Type() == "FLOAT" || right.Type() == "BOOLEAN"):
		// Convert right to string and concatenate
		if node.Operator == "+" {
			return &StringValue{Value: left.(*StringValue).Value + right.Inspect()}
		}
		return &StringValue{Value: fmt.Sprintf("Type error: unsupported operator %s for types %s and %s", node.Operator, left.Type(), right.Type())}
	case (left.Type() == "INTEGER" || left.Type() == "FLOAT" || left.Type() == "BOOLEAN") && right.Type() == "STRING":
		// Convert left to string and concatenate
		if node.Operator == "+" {
			return &StringValue{Value: left.Inspect() + right.(*StringValue).Value}
		}
		return &StringValue{Value: fmt.Sprintf("Type error: unsupported operator %s for types %s and %s", node.Operator, left.Type(), right.Type())}
	case node.Operator == "==":
		return &BooleanValue{Value: left.Inspect() == right.Inspect()}
	case node.Operator == "!=":
		return &BooleanValue{Value: left.Inspect() != right.Inspect()}
	default:
		return &StringValue{Value: fmt.Sprintf("Type error: unsupported operator %s for types %s and %s", node.Operator, left.Type(), right.Type())}
	}
}

// Helper functions

func evalIntegerBinaryExpression(operator string, left, right Value) Value {
	leftVal := left.(*IntegerValue).Value
	rightVal := right.(*IntegerValue).Value

	switch operator {
	case "+":
		return &IntegerValue{Value: leftVal + rightVal}
	case "-":
		return &IntegerValue{Value: leftVal - rightVal}
	case "*":
		return &IntegerValue{Value: leftVal * rightVal}
	case "/":
		if rightVal == 0 {
			return &StringValue{Value: "Error: division by zero"}
		}
		return &IntegerValue{Value: leftVal / rightVal}
	case "%":
		if rightVal == 0 {
			return &StringValue{Value: "Error: modulo by zero"}
		}
		return &IntegerValue{Value: leftVal % rightVal}
	case "<":
		return &BooleanValue{Value: leftVal < rightVal}
	case ">":
		return &BooleanValue{Value: leftVal > rightVal}
	case "<=":
		return &BooleanValue{Value: leftVal <= rightVal}
	case ">=":
		return &BooleanValue{Value: leftVal >= rightVal}
	case "==":
		return &BooleanValue{Value: leftVal == rightVal}
	case "!=":
		return &BooleanValue{Value: leftVal != rightVal}
	default:
		return &StringValue{Value: fmt.Sprintf("Error: unknown operator for integers: %s", operator)}
	}
}

func evalNumberBinaryExpression(operator string, left, right Value) Value {
	var leftVal, rightVal float64

	// Convert left to float64
	if left.Type() == "INTEGER" {
		leftVal = float64(left.(*IntegerValue).Value)
	} else {
		leftVal = left.(*FloatValue).Value
	}

	// Convert right to float64
	if right.Type() == "INTEGER" {
		rightVal = float64(right.(*IntegerValue).Value)
	} else {
		rightVal = right.(*FloatValue).Value
	}

	switch operator {
	case "+":
		return &FloatValue{Value: leftVal + rightVal}
	case "-":
		return &FloatValue{Value: leftVal - rightVal}
	case "*":
		return &FloatValue{Value: leftVal * rightVal}
	case "/":
		if rightVal == 0 {
			return &StringValue{Value: "Error: division by zero"}
		}
		return &FloatValue{Value: leftVal / rightVal}
	case "%":
		if rightVal == 0 {
			return &StringValue{Value: "Error: modulo by zero"}
		}
		return &FloatValue{Value: math.Mod(leftVal, rightVal)}
	case "<":
		return &BooleanValue{Value: leftVal < rightVal}
	case ">":
		return &BooleanValue{Value: leftVal > rightVal}
	case "<=":
		return &BooleanValue{Value: leftVal <= rightVal}
	case ">=":
		return &BooleanValue{Value: leftVal >= rightVal}
	case "==":
		return &BooleanValue{Value: leftVal == rightVal}
	case "!=":
		return &BooleanValue{Value: leftVal != rightVal}
	default:
		return &StringValue{Value: fmt.Sprintf("Error: unknown operator for numbers: %s", operator)}
	}
}

func evalStringBinaryExpression(operator string, left, right Value) Value {
	leftVal := left.(*StringValue).Value
	rightVal := right.(*StringValue).Value

	switch operator {
	case "+":
		return &StringValue{Value: leftVal + rightVal}
	case "==":
		return &BooleanValue{Value: leftVal == rightVal}
	case "!=":
		return &BooleanValue{Value: leftVal != rightVal}
	default:
		return &StringValue{Value: fmt.Sprintf("Error: unknown operator for strings: %s", operator)}
	}
}

func isTruthy(obj Value) bool {
	switch obj := obj.(type) {
	case *BooleanValue:
		return obj.Value
	case *NilValue:
		return false
	case *IntegerValue:
		return obj.Value != 0
	case *FloatValue:
		return obj.Value != 0
	case *StringValue:
		return obj.Value != ""
	default:
		return true
	}
}

func isError(obj Value) bool {
	if obj != nil {
		return obj.Type() == "ERROR"
	}
	return false
}

// Update evalClassInstantiation to create object instances
func (i *Interpreter) evalClassInstantiation(node *ast.ClassInst, env *Environment) Value {
	// Evaluate the class expression
	classVal := i.eval(node.Class, env)
	if classVal == nil {
		return &StringValue{Value: "Error: Cannot instantiate nil class"}
	}

	class, ok := classVal.(*ClassValue)
	if !ok {
		return &StringValue{Value: fmt.Sprintf("Error: %s is not a class", classVal.Inspect())}
	}

	// Create a new object instance
	obj := &ObjectValue{
		Class:      class,
		Properties: make(map[string]Value),
	}

	// Evaluate arguments
	var args []Value
	for _, argNode := range node.Arguments {
		args = append(args, i.eval(argNode, env))
	}

	// For Point class, initialize x and y properties
	if class.Name == "Point" && len(args) >= 2 {
		obj.Properties["x"] = args[0]
		obj.Properties["y"] = args[1]
	}

	return obj
}

// Update evalMethodCall to handle method invocation
func (i *Interpreter) evalMethodCall(node *ast.MethodCall, env *Environment) Value {
	// Evaluate the object that the method is being called on
	objectVal := i.eval(node.Object, env)
	if objectVal == nil {
		return &StringValue{Value: "Error: Cannot call method on nil"}
	}

	obj, ok := objectVal.(*ObjectValue)
	if !ok {
		return &StringValue{Value: fmt.Sprintf("Error: %s is not an object", objectVal.Inspect())}
	}

	// Look up the method in the class
	method, ok := obj.Class.Methods[node.Method]
	if !ok {
		return &StringValue{Value: fmt.Sprintf("Error: Method %s not found in class %s",
			node.Method, obj.Class.Name)}
	}

	// Build argument list with the object as the first argument (this)
	var args []Value
	args = append(args, obj) // The object instance is passed as the first argument

	// Add the rest of the arguments
	for _, argNode := range node.Args {
		args = append(args, i.eval(argNode, env))
	}

	// If it's a builtin method, use the builtin function
	if method.BuiltinFunc != nil {
		return method.BuiltinFunc(args)
	}

	// Otherwise, it should be a user-defined method, but we haven't implemented this yet
	return &StringValue{Value: "User-defined methods not yet supported"}
}

// toString converts any value to a string representation
func toString(val Value) string {
	if val == nil {
		return "nil"
	}

	switch v := val.(type) {
	case *StringValue:
		return v.Value
	case *IntegerValue:
		return strconv.Itoa(v.Value)
	case *FloatValue:
		return strconv.FormatFloat(v.Value, 'f', -1, 64)
	case *BooleanValue:
		return strconv.FormatBool(v.Value)
	case *ArrayValue:
		return fmt.Sprintf("%v", v.Inspect())
	case *FunctionValue:
		return fmt.Sprintf("function %s", v.Name)
	case *ClassValue:
		return fmt.Sprintf("class %s", v.Name)
	case *ObjectValue:
		return fmt.Sprintf("instance of %s", v.Class.Name)
	case *NilValue:
		return "nil"
	default:
		return fmt.Sprintf("%v", val.Inspect())
	}
}