package interpreter

import (
	"fmt"
	"strconv"

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

// FunctionValue represents a function
type FunctionValue struct {
	Name           string
	Parameters     []string
	ParameterTypes []types.Type
	Body           *parser.BlockStmt
	ReturnType     types.Type
	Env            *Environment
}

func (f *FunctionValue) Type() string { return "FUNCTION" }
func (f *FunctionValue) Inspect() string {
	return fmt.Sprintf("function %s", f.Name)
}
func (f *FunctionValue) VibeType() types.Type {
	return types.FunctionType{
		ParameterTypes: f.ParameterTypes,
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

// Interpreter executes the AST
type Interpreter struct {
	env *Environment
}

// New creates a new interpreter
func New() *Interpreter {
	env := NewEnvironment()

	// Register built-in functions
	registerBuiltins(env)

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

// Eval evaluates the AST and returns the result
func (i *Interpreter) Eval(node parser.Node) Value {
	return i.eval(node, i.env)
}

func (i *Interpreter) eval(node parser.Node, env *Environment) Value {
	switch node := node.(type) {
	case *parser.Program:
		return i.evalProgram(node, env)
	case *parser.BlockStmt:
		return i.evalBlockStatement(node, env)
	case *parser.NumberLiteral:
		if node.IsInt {
			return &IntegerValue{Value: int(node.Value)}
		}
		return &FloatValue{Value: node.Value}
	case *parser.StringLiteral:
		return &StringValue{Value: node.Value}
	case *parser.BooleanLiteral:
		return &BooleanValue{Value: node.Value}
	case *parser.NilLiteral:
		return &NilValue{}
	case *parser.Identifier:
		return i.evalIdentifier(node, env)
	case *parser.PrintStmt:
		return i.evalPrintStatement(node, env)
	case *parser.Assignment:
		return i.evalAssignment(node, env)
	case *parser.VariableDecl:
		return i.evalVariableDeclaration(node, env)
	case *parser.FunctionDef:
		return i.evalFunctionDefinition(node, env)
	case *parser.CallExpr:
		return i.evalCallExpression(node, env)
	case *parser.ReturnStmt:
		return i.evalReturnStatement(node, env)
	case *parser.IfStmt:
		return i.evalIfStatement(node, env)
	case *parser.WhileStmt:
		return i.evalWhileStatement(node, env)
	case *parser.BinaryExpr:
		return i.evalBinaryExpression(node, env)
	case *parser.TypeAnnotation:
		// Type annotations don't evaluate to a value on their own
		return &NilValue{}
	case *parser.TypeDeclaration:
		// Type declarations don't evaluate to a value
		return &NilValue{}
	default:
		// Handle unexpected nodes
		return &StringValue{Value: fmt.Sprintf("Unknown node type: %T", node)}
	}
}

func (i *Interpreter) evalVariableDeclaration(node *parser.VariableDecl, env *Environment) Value {
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

func (i *Interpreter) parseTypeAnnotation(node *parser.TypeAnnotation) types.Type {
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
			elemType := i.parseTypeAnnotation(node.TypeParams[0].(*parser.TypeAnnotation))
			return types.ArrayType{ElementType: elemType}
		}
		// Default to Array of any
		return types.ArrayType{ElementType: types.AnyType}
	case "union":
		if len(node.TypeParams) > 0 {
			var unionTypes []types.Type
			for _, param := range node.TypeParams {
				unionTypes = append(unionTypes, i.parseTypeAnnotation(param.(*parser.TypeAnnotation)))
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

func (i *Interpreter) evalProgram(program *parser.Program, env *Environment) Value {
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

func (i *Interpreter) evalBlockStatement(block *parser.BlockStmt, env *Environment) Value {
	var result Value
	result = &NilValue{}

	for _, statement := range block.Statements {
		result = i.eval(statement, env)

		// If we hit a return statement, break execution and return it up
		if result.Type() == "RETURN" {
			return result
		}
	}

	// The last statement in a block is treated as the implicit return value,
	// but only for expression statements, not declarations or control flow
	if len(block.Statements) > 0 {
		lastStmt := block.Statements[len(block.Statements)-1]
		switch lastStmt.(type) {
		case *parser.Assignment, *parser.VariableDecl, *parser.PrintStmt, *parser.TypeDeclaration:
			// These statements don't produce a value to return
			return &NilValue{}
		case *parser.IfStmt, *parser.WhileStmt, *parser.FunctionDef:
			// Control flow statements are handled separately
			return result
		default:
			// For expression statements, return the result
			return result
		}
	}

	return result
}

func (i *Interpreter) evalIdentifier(node *parser.Identifier, env *Environment) Value {
	if val, ok := env.Get(node.Name); ok {
		return val
	}

	return &StringValue{Value: fmt.Sprintf("Error: variable '%s' not found", node.Name)}
}

func (i *Interpreter) evalPrintStatement(node *parser.PrintStmt, env *Environment) Value {
	value := i.eval(node.Value, env)
	fmt.Println(value.Inspect())
	return &NilValue{}
}

func (i *Interpreter) evalAssignment(node *parser.Assignment, env *Environment) Value {
	val := i.eval(node.Value, env)

	err := env.Set(node.Name, val)
	if err != nil {
		return &StringValue{Value: err.Error()}
	}

	return &NilValue{}
}

func (i *Interpreter) evalFunctionDefinition(node *parser.FunctionDef, env *Environment) Value {
	// Parse parameter types
	paramTypes := make([]types.Type, len(node.ParamTypes))
	for j, paramType := range node.ParamTypes {
		paramTypes[j] = i.parseTypeAnnotation(paramType)
	}

	// Parse return type
	var returnType types.Type
	if node.ReturnType != nil {
		returnType = i.parseTypeAnnotation(node.ReturnType)
	} else {
		returnType = types.AnyType
	}

	// Create the function value
	function := &FunctionValue{
		Name:           node.Name,
		Parameters:     node.Parameters,
		ParameterTypes: paramTypes,
		Body:           node.Body,
		ReturnType:     returnType,
		Env:            env,
	}

	// Add the function to the environment
	env.SetWithType(node.Name, function, function.VibeType())

	return &NilValue{}
}

func (i *Interpreter) evalCallExpression(node *parser.CallExpr, env *Environment) Value {
	// Evaluate the function expression
	function := i.eval(node.Function, env)

	// Evaluate arguments
	args := make([]Value, len(node.Args))
	for idx, arg := range node.Args {
		args[idx] = i.eval(arg, env)
	}

	// Call the function
	switch fn := function.(type) {
	case *BuiltinFunction:
		// Type check arguments for built-in functions
		for idx, arg := range args {
			if idx < len(fn.ParamTypes) && !types.IsAssignable(arg.VibeType(), fn.ParamTypes[idx]) {
				return &StringValue{Value: fmt.Sprintf(
					"Type error: Parameter %d of %s expects %s, got %s",
					idx+1, fn.Name, fn.ParamTypes[idx].String(), arg.VibeType().String())}
			}
		}
		return fn.Fn(args)
	case *FunctionValue:
		// Create a new environment for the function call
		newEnv := NewEnclosedEnvironment(fn.Env)

		// Bind parameters to arguments
		for paramIdx, param := range fn.Parameters {
			if paramIdx < len(args) {
				// Type check the argument
				if paramIdx < len(fn.ParameterTypes) && !types.IsAssignable(args[paramIdx].VibeType(), fn.ParameterTypes[paramIdx]) {
					return &StringValue{Value: fmt.Sprintf(
						"Type error: Parameter '%s' of function '%s' expects %s, got %s",
						param, fn.Name, fn.ParameterTypes[paramIdx].String(), args[paramIdx].VibeType().String())}
				}

				// Bind the parameter
				newEnv.SetWithType(param, args[paramIdx], fn.ParameterTypes[paramIdx])
			} else {
				// Missing argument, use nil
				newEnv.SetWithType(param, &NilValue{}, fn.ParameterTypes[paramIdx])
			}
		}

		// Evaluate the function body in the new environment
		evaluated := i.eval(fn.Body, newEnv)

		// Unwrap return value if needed
		var returnValue Value
		if rv, ok := evaluated.(*ReturnValue); ok {
			returnValue = rv.Value
		} else {
			// Implicit return
			returnValue = evaluated
		}

		// Check that the return value is compatible with the declared return type
		if !types.IsAssignable(returnValue.VibeType(), fn.ReturnType) {
			return &StringValue{Value: fmt.Sprintf(
				"Type error: Function '%s' is declared to return %s but returned %s",
				fn.Name, fn.ReturnType.String(), returnValue.VibeType().String())}
		}

		return returnValue
	default:
		return &StringValue{Value: fmt.Sprintf("Type error: %s is not a function", function.Inspect())}
	}
}

func (i *Interpreter) evalReturnStatement(node *parser.ReturnStmt, env *Environment) Value {
	var value Value

	if node.Value != nil {
		value = i.eval(node.Value, env)
	} else {
		value = &NilValue{}
	}

	return &ReturnValue{Value: value}
}

func (i *Interpreter) evalIfStatement(node *parser.IfStmt, env *Environment) Value {
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

func (i *Interpreter) evalWhileStatement(node *parser.WhileStmt, env *Environment) Value {
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

func (i *Interpreter) evalBinaryExpression(node *parser.BinaryExpr, env *Environment) Value {
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