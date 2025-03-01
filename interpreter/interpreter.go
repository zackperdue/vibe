package interpreter

import (
	"fmt"
	"strconv"

	"github.com/example/crystal/parser"
	"github.com/example/crystal/types"
)

// Value represents a runtime value
type Value interface {
	Type() string
	Inspect() string
	CrystalType() types.Type
}

// IntegerValue represents an integer value
type IntegerValue struct {
	Value int
}

func (i *IntegerValue) Type() string    { return "INTEGER" }
func (i *IntegerValue) Inspect() string { return strconv.Itoa(i.Value) }
func (i *IntegerValue) CrystalType() types.Type { return types.IntType }

// FloatValue represents a floating-point value
type FloatValue struct {
	Value float64
}

func (f *FloatValue) Type() string    { return "FLOAT" }
func (f *FloatValue) Inspect() string { return strconv.FormatFloat(f.Value, 'f', -1, 64) }
func (f *FloatValue) CrystalType() types.Type { return types.FloatType }

// StringValue represents a string value
type StringValue struct {
	Value string
}

func (s *StringValue) Type() string    { return "STRING" }
func (s *StringValue) Inspect() string { return s.Value }
func (s *StringValue) CrystalType() types.Type { return types.StringType }

// BooleanValue represents a boolean value
type BooleanValue struct {
	Value bool
}

func (b *BooleanValue) Type() string    { return "BOOLEAN" }
func (b *BooleanValue) Inspect() string { return strconv.FormatBool(b.Value) }
func (b *BooleanValue) CrystalType() types.Type { return types.BoolType }

// NilValue represents a nil value
type NilValue struct{}

func (n *NilValue) Type() string    { return "NIL" }
func (n *NilValue) Inspect() string { return "nil" }
func (n *NilValue) CrystalType() types.Type { return types.NilType }

// ReturnValue wraps a value being returned from a function
type ReturnValue struct {
	Value Value
}

func (r *ReturnValue) Type() string    { return "RETURN" }
func (r *ReturnValue) Inspect() string { return r.Value.Inspect() }
func (r *ReturnValue) CrystalType() types.Type { return r.Value.CrystalType() }

// FunctionValue represents a user-defined function
type FunctionValue struct {
	Name        string
	Parameters  []string
	ParameterTypes []types.Type
	ReturnType  types.Type
	Body        *parser.BlockStmt
	Env         *Environment
}

func (f *FunctionValue) Type() string { return "FUNCTION" }
func (f *FunctionValue) Inspect() string {
	return fmt.Sprintf("<function:%s>", f.Name)
}
func (f *FunctionValue) CrystalType() types.Type {
	paramTypes := make([]types.Type, len(f.ParameterTypes))
	copy(paramTypes, f.ParameterTypes)
	return types.FunctionType{
		ParameterTypes: paramTypes,
		ReturnType:     f.ReturnType,
	}
}

// ArrayValue represents an array
type ArrayValue struct {
	Elements []Value
	ElementType types.Type
}

func (a *ArrayValue) Type() string { return "ARRAY" }
func (a *ArrayValue) Inspect() string {
	result := "["
	for i, elem := range a.Elements {
		if i > 0 {
			result += ", "
		}
		result += elem.Inspect()
	}
	result += "]"
	return result
}
func (a *ArrayValue) CrystalType() types.Type {
	return types.ArrayType{ElementType: a.ElementType}
}

// Environment stores variable bindings
type Environment struct {
	store    map[string]Value
	types    map[string]types.Type
	outer    *Environment
	builtins map[string]BuiltinFunction
}

// NewEnvironment creates a new environment
func NewEnvironment() *Environment {
	return &Environment{
		store:    make(map[string]Value),
		types:    make(map[string]types.Type),
		builtins: make(map[string]BuiltinFunction),
	}
}

// NewEnclosedEnvironment creates a new environment with an outer environment
func NewEnclosedEnvironment(outer *Environment) *Environment {
	env := NewEnvironment()
	env.outer = outer
	env.builtins = outer.builtins
	return env
}

// Get returns the value of a variable or builtin function
func (e *Environment) Get(name string) (Value, bool) {
	// Check builtins first
	if builtin, ok := e.builtins[name]; ok {
		return &builtin, true
	}

	// Then check regular variables
	obj, ok := e.store[name]
	if !ok && e.outer != nil {
		obj, ok = e.outer.Get(name)
	}
	return obj, ok
}

// Set sets the value of a variable
func (e *Environment) Set(name string, val Value) Value {
	e.store[name] = val
	return val
}

// SetWithType sets a variable with a specific type
func (e *Environment) SetWithType(name string, val Value, typ types.Type) (Value, error) {
	// Check if the variable already has a type
	if existingType, ok := e.GetType(name); ok {
		// Variable exists, check if type compatible
		if !types.IsAssignable(val.CrystalType(), existingType) {
			return nil, fmt.Errorf("type error: cannot assign %s to variable '%s' of type %s",
				val.CrystalType().String(), name, existingType.String())
		}
	} else {
		// New variable, set its type
		e.types[name] = typ
	}

	e.store[name] = val
	return val, nil
}

// GetType returns the type of a variable
func (e *Environment) GetType(name string) (types.Type, bool) {
	typ, ok := e.types[name]
	if !ok && e.outer != nil {
		typ, ok = e.outer.GetType(name)
	}
	return typ, ok
}

// RegisterBuiltin registers a builtin function
func (e *Environment) RegisterBuiltin(name string, fn BuiltinFunction) {
	e.builtins[name] = fn
}

// BuiltinFunction represents a built-in function
type BuiltinFunction struct {
	Name string
	Fn   func(args ...Value) Value
	ReturnType types.Type
	ParamTypes []types.Type
}

func (b *BuiltinFunction) Type() string    { return "BUILTIN" }
func (b *BuiltinFunction) Inspect() string { return fmt.Sprintf("<builtin:%s>", b.Name) }
func (b *BuiltinFunction) CrystalType() types.Type {
	paramTypes := make([]types.Type, len(b.ParamTypes))
	copy(paramTypes, b.ParamTypes)
	return types.FunctionType{
		ParameterTypes: paramTypes,
		ReturnType:     b.ReturnType,
	}
}

// Interpreter executes the AST
type Interpreter struct {
	env *Environment
	typeChecker *types.TypeChecker
}

// New creates a new interpreter
func New() *Interpreter {
	env := NewEnvironment()
	interpreter := &Interpreter{
		env: env,
		typeChecker: types.NewTypeChecker(),
	}

	// Register built-in functions
	interpreter.registerBuiltins()

	return interpreter
}

func (i *Interpreter) registerBuiltins() {
	i.env.RegisterBuiltin("len", BuiltinFunction{
		Name: "len",
		ReturnType: types.IntType,
		ParamTypes: []types.Type{types.AnyType},
		Fn: func(args ...Value) Value {
			if len(args) != 1 {
				return &StringValue{Value: "wrong number of arguments for len()"}
			}

			switch arg := args[0].(type) {
			case *StringValue:
				return &IntegerValue{Value: len(arg.Value)}
			case *ArrayValue:
				return &IntegerValue{Value: len(arg.Elements)}
			default:
				return &StringValue{Value: "argument to len() not supported"}
			}
		},
	})
}

// Eval evaluates a program
func (i *Interpreter) Eval(program parser.Node) Value {
	return i.eval(program, i.env)
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
		// Special case for true/false identifiers
		if node.Name == "true" {
			return &BooleanValue{Value: true}
		} else if node.Name == "false" {
			return &BooleanValue{Value: false}
		}
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
		return &NilValue{}
	}
}

// New function to evaluate a variable declaration with type
func (i *Interpreter) evalVariableDeclaration(node *parser.VariableDecl, env *Environment) Value {
	var value Value = &NilValue{}

	// Evaluate the initial value if provided
	if node.Value != nil {
		value = i.eval(node.Value, env)
	}

	// Convert parser.TypeAnnotation to types.Type
	varType := i.parseTypeAnnotation(node.TypeAnnotation)

	// Check if the value's type is compatible with the declared type
	if !types.IsAssignable(value.CrystalType(), varType) {
		fmt.Printf("Type error: Cannot assign value of type %s to variable of type %s\n",
			value.CrystalType().String(), varType.String())
		return &NilValue{}
	}

	// Set the variable with its type
	result, err := env.SetWithType(node.Name, value, varType)
	if err != nil {
		fmt.Println(err)
		return &NilValue{}
	}

	return result
}

// Helper function to convert parser.TypeAnnotation to types.Type
func (i *Interpreter) parseTypeAnnotation(typeAnnotation *parser.TypeAnnotation) types.Type {
	if typeAnnotation == nil {
		return types.AnyType
	}

	switch typeAnnotation.TypeName {
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
		// Handle array type with element type
		if len(typeAnnotation.TypeParams) > 0 {
			if param, ok := typeAnnotation.TypeParams[0].(*parser.TypeAnnotation); ok {
				elementType := i.parseTypeAnnotation(param)
				return types.ArrayType{ElementType: elementType}
			}
		}
		// Default to Array<any> if no element type specified
		return types.ArrayType{ElementType: types.AnyType}
	case "union":
		// Handle union types
		if len(typeAnnotation.TypeParams) > 0 {
			unionTypes := make([]types.Type, 0, len(typeAnnotation.TypeParams))
			for _, param := range typeAnnotation.TypeParams {
				if typeParam, ok := param.(*parser.TypeAnnotation); ok {
					unionTypes = append(unionTypes, i.parseTypeAnnotation(typeParam))
				}
			}
			return types.UnionType{Types: unionTypes}
		}
	}

	// Default to any type if unknown
	return types.AnyType
}

// Evaluate a program by evaluating each statement sequentially
func (i *Interpreter) evalProgram(program *parser.Program, env *Environment) Value {
	var result Value = &NilValue{}

	for _, statement := range program.Statements {
		result = i.eval(statement, env)

		// If the statement returns a value, we need to unwrap it
		if returnValue, ok := result.(*ReturnValue); ok {
			return returnValue.Value
		}
	}

	return result
}

// Evaluate a block statement
func (i *Interpreter) evalBlockStatement(block *parser.BlockStmt, env *Environment) Value {
	var result Value = &NilValue{}

	for idx, statement := range block.Statements {
		// For the last statement, if it's not a return statement and not a control flow statement,
		// it's the implicit return value
		isLastStatement := idx == len(block.Statements)-1

		result = i.eval(statement, env)

		// Handle explicit return statements
		if result != nil && result.Type() == "RETURN" {
			return result
		}

		// For the last statement, treat it as the implicit return if it's an expression
		if isLastStatement && result != nil {
			// Don't wrap nil values, assignments, or declarations as returns
			switch statement.(type) {
			case *parser.Assignment, *parser.VariableDecl, *parser.PrintStmt, *parser.TypeDeclaration:
				// Keep the result as is - these aren't expression values to return
			default:
				// For expressions, wrap as a return value but don't modify the type
				if _, ok := result.(*ReturnValue); !ok {
					result = &ReturnValue{Value: result}
				}
			}
		}
	}

	return result
}

// Evaluate an identifier (variable lookup)
func (i *Interpreter) evalIdentifier(node *parser.Identifier, env *Environment) Value {
	val, ok := env.Get(node.Name)
	if !ok {
		return &NilValue{}
	}
	return val
}

// Evaluate a print statement
func (i *Interpreter) evalPrintStatement(node *parser.PrintStmt, env *Environment) Value {
	val := i.eval(node.Value, env)
	fmt.Println(val.Inspect())
	return val
}

// Evaluate a variable assignment
func (i *Interpreter) evalAssignment(node *parser.Assignment, env *Environment) Value {
	val := i.eval(node.Value, env)

	// Check if variable has a declared type
	if varType, ok := env.GetType(node.Name); ok {
		// Check if types are compatible
		if !types.IsAssignable(val.CrystalType(), varType) {
			fmt.Printf("Type error: Cannot assign value of type %s to variable '%s' of type %s\n",
				val.CrystalType().String(), node.Name, varType.String())
			return &NilValue{}
		}
	}

	env.Set(node.Name, val)
	return val
}

// Evaluate a function definition
func (i *Interpreter) evalFunctionDefinition(node *parser.FunctionDef, env *Environment) Value {
	// Convert parameter types
	var paramTypes []types.Type
	for _, typeAnnotation := range node.ParamTypes {
		paramType := i.parseTypeAnnotation(typeAnnotation)
		paramTypes = append(paramTypes, paramType)
	}

	// Convert return type
	returnType := i.parseTypeAnnotation(node.ReturnType)

	fn := &FunctionValue{
		Name:        node.Name,
		Parameters:  node.Parameters,
		ParameterTypes: paramTypes,
		ReturnType:  returnType,
		Body:        node.Body,
		Env:         env,
	}

	env.Set(node.Name, fn)
	return fn
}

// Evaluate a function call
func (i *Interpreter) evalCallExpression(call *parser.CallExpr, env *Environment) Value {
	function := i.eval(call.Function, env)
	args := i.evalExpressions(call.Args, env)

	// Check if it's a builtin function
	if builtin, ok := function.(*BuiltinFunction); ok {
		// Type check parameters for builtin functions
		if len(builtin.ParamTypes) > 0 && len(args) != len(builtin.ParamTypes) {
			fmt.Printf("Type error: Function %s expects %d arguments, got %d\n",
				builtin.Name, len(builtin.ParamTypes), len(args))
			return &NilValue{}
		}

		for idx, arg := range args {
			if idx < len(builtin.ParamTypes) && !types.IsAssignable(arg.CrystalType(), builtin.ParamTypes[idx]) {
				fmt.Printf("Type error: Parameter %d of function %s expects type %s, got %s\n",
					idx+1, builtin.Name, builtin.ParamTypes[idx].String(), arg.CrystalType().String())
				return &NilValue{}
			}
		}

		return builtin.Fn(args...)
	}

	// Otherwise, it's a user-defined function
	if fn, ok := function.(*FunctionValue); ok {
		// Type check parameters
		if len(fn.ParameterTypes) > 0 && len(args) != len(fn.Parameters) {
			fmt.Printf("Type error: Function %s expects %d arguments, got %d\n",
				fn.Name, len(fn.Parameters), len(args))
			return &NilValue{}
		}

		// Create a new environment for the function call
		functionEnv := NewEnclosedEnvironment(fn.Env)

		// Bind parameters to arguments
		for paramIdx, param := range fn.Parameters {
			if paramIdx < len(args) {
				// Type check the argument
				if paramIdx < len(fn.ParameterTypes) && !types.IsAssignable(args[paramIdx].CrystalType(), fn.ParameterTypes[paramIdx]) {
					fmt.Printf("Type error: Parameter %s of function %s expects type %s, got %s\n",
						param, fn.Name, fn.ParameterTypes[paramIdx].String(), args[paramIdx].CrystalType().String())
					return &NilValue{}
				}

				// Set parameter with its type
				_, err := functionEnv.SetWithType(param, args[paramIdx], fn.ParameterTypes[paramIdx])
				if err != nil {
					fmt.Println(err)
					return &NilValue{}
				}
			} else {
				// Default value for missing arguments
				_, err := functionEnv.SetWithType(param, &NilValue{}, fn.ParameterTypes[paramIdx])
				if err != nil {
					fmt.Println(err)
					return &NilValue{}
				}
			}
		}

		// Evaluate the function body
		evaluated := i.eval(fn.Body, functionEnv)

		// Unwrap the return value if necessary
		var returnValue Value
		if retVal, ok := evaluated.(*ReturnValue); ok {
			returnValue = retVal.Value
		} else {
			returnValue = evaluated
		}

		// Type check the return value
		if !types.IsAssignable(returnValue.CrystalType(), fn.ReturnType) {
			fmt.Printf("Type error: Function %s should return type %s, got %s\n",
				fn.Name, fn.ReturnType.String(), returnValue.CrystalType().String())
			return &NilValue{}
		}

		return returnValue
	}

	return &StringValue{Value: "not a function: " + function.Type()}
}

// Evaluate a list of expressions
func (i *Interpreter) evalExpressions(exps []parser.Node, env *Environment) []Value {
	var result []Value

	for _, exp := range exps {
		evaluated := i.eval(exp, env)
		result = append(result, evaluated)
	}

	return result
}

// Evaluate a return statement
func (i *Interpreter) evalReturnStatement(node *parser.ReturnStmt, env *Environment) Value {
	if node.Value == nil {
		return &ReturnValue{Value: &NilValue{}}
	}

	val := i.eval(node.Value, env)
	return &ReturnValue{Value: val}
}

// Evaluate an if statement
func (i *Interpreter) evalIfStatement(node *parser.IfStmt, env *Environment) Value {
	condition := i.eval(node.Condition, env)

	if isTruthy(condition) {
		result := i.eval(node.Consequence, env)
		// If it's a return value (explicit or implicit), unwrap and return it
		if ret, ok := result.(*ReturnValue); ok {
			return ret
		}
		return result
	} else if len(node.ElseIfBlocks) > 0 {
		// Check each elsif block
		for _, elseIf := range node.ElseIfBlocks {
			if isTruthy(i.eval(elseIf.Condition, env)) {
				result := i.eval(elseIf.Consequence, env)
				// If it's a return value (explicit or implicit), unwrap and return it
				if ret, ok := result.(*ReturnValue); ok {
					return ret
				}
				return result
			}
		}
	}

	if node.Alternative != nil {
		result := i.eval(node.Alternative, env)
		// If it's a return value (explicit or implicit), unwrap and return it
		if ret, ok := result.(*ReturnValue); ok {
			return ret
		}
		return result
	}

	return &NilValue{}
}

// Evaluate a while statement
func (i *Interpreter) evalWhileStatement(node *parser.WhileStmt, env *Environment) Value {
	var result Value = &NilValue{}

	for {
		condition := i.eval(node.Condition, env)
		if !isTruthy(condition) {
			break
		}

		result = i.eval(node.Body, env)

		// Check for a return statement within the loop
		if returnValue, ok := result.(*ReturnValue); ok {
			return returnValue
		}
	}

	return result
}

// Evaluate a binary expression
func (i *Interpreter) evalBinaryExpression(node *parser.BinaryExpr, env *Environment) Value {
	// Handle the case where Left is nil (parser limitation)
	var left Value
	if node.Left == nil {
		left = &NilValue{}
	} else {
		left = i.eval(node.Left, env)
	}

	right := i.eval(node.Right, env)

	// Special case for unary operations where the left operand is nil
	if left.Type() == "NIL" {
		if right.Type() == "INTEGER" {
			rightVal := right.(*IntegerValue).Value
			switch node.Operator {
			case "+":
				// Unary plus - just return the value
				return &IntegerValue{Value: rightVal}
			case "-":
				// Unary minus - negate the value
				return &IntegerValue{Value: -rightVal}
			}
		}
	}

	// String concatenation - convert any value to string if one operand is a string
	if node.Operator == "+" && (left.Type() == "STRING" || right.Type() == "STRING") {
		leftStr := left.Inspect()
		rightStr := right.Inspect()
		return &StringValue{Value: leftStr + rightStr}
	}

	// Integer operations
	if left.Type() == "INTEGER" && right.Type() == "INTEGER" {
		leftVal := left.(*IntegerValue).Value
		rightVal := right.(*IntegerValue).Value

		switch node.Operator {
		case "+":
			return &IntegerValue{Value: leftVal + rightVal}
		case "-":
			return &IntegerValue{Value: leftVal - rightVal}
		case "*":
			return &IntegerValue{Value: leftVal * rightVal}
		case "/":
			if rightVal == 0 {
				return &StringValue{Value: "division by zero"}
			}
			return &IntegerValue{Value: leftVal / rightVal}
		case "==":
			return &BooleanValue{Value: leftVal == rightVal}
		case "!=":
			return &BooleanValue{Value: leftVal != rightVal}
		case "<":
			return &BooleanValue{Value: leftVal < rightVal}
		case ">":
			return &BooleanValue{Value: leftVal > rightVal}
		case "<=":
			return &BooleanValue{Value: leftVal <= rightVal}
		case ">=":
			return &BooleanValue{Value: leftVal >= rightVal}
		}
	}

	// Float operations
	if (left.Type() == "INTEGER" || left.Type() == "FLOAT") &&
	   (right.Type() == "INTEGER" || right.Type() == "FLOAT") {
		var leftVal, rightVal float64

		if left.Type() == "INTEGER" {
			leftVal = float64(left.(*IntegerValue).Value)
		} else {
			leftVal = left.(*FloatValue).Value
		}

		if right.Type() == "INTEGER" {
			rightVal = float64(right.(*IntegerValue).Value)
		} else {
			rightVal = right.(*FloatValue).Value
		}

		switch node.Operator {
		case "+":
			return &FloatValue{Value: leftVal + rightVal}
		case "-":
			return &FloatValue{Value: leftVal - rightVal}
		case "*":
			return &FloatValue{Value: leftVal * rightVal}
		case "/":
			if rightVal == 0 {
				return &StringValue{Value: "division by zero"}
			}
			return &FloatValue{Value: leftVal / rightVal}
		case "==":
			return &BooleanValue{Value: leftVal == rightVal}
		case "!=":
			return &BooleanValue{Value: leftVal != rightVal}
		case "<":
			return &BooleanValue{Value: leftVal < rightVal}
		case ">":
			return &BooleanValue{Value: leftVal > rightVal}
		case "<=":
			return &BooleanValue{Value: leftVal <= rightVal}
		case ">=":
			return &BooleanValue{Value: leftVal >= rightVal}
		}
	}

	// String operations
	if left.Type() == "STRING" && right.Type() == "STRING" {
		leftVal := left.(*StringValue).Value
		rightVal := right.(*StringValue).Value

		switch node.Operator {
		case "==":
			return &BooleanValue{Value: leftVal == rightVal}
		case "!=":
			return &BooleanValue{Value: leftVal != rightVal}
		}
	}

	// Boolean operations
	if left.Type() == "BOOLEAN" && right.Type() == "BOOLEAN" {
		leftVal := left.(*BooleanValue).Value
		rightVal := right.(*BooleanValue).Value

		switch node.Operator {
		case "&&":
			return &BooleanValue{Value: leftVal && rightVal}
		case "||":
			return &BooleanValue{Value: leftVal || rightVal}
		case "==":
			return &BooleanValue{Value: leftVal == rightVal}
		case "!=":
			return &BooleanValue{Value: leftVal != rightVal}
		}
	}

	// Unsupported operation
	return &StringValue{Value: fmt.Sprintf("unsupported operation: %s %s %s",
		left.Type(), node.Operator, right.Type())}
}

// Helper function to determine if a value is truthy
func isTruthy(val Value) bool {
	switch val := val.(type) {
	case *NilValue:
		return false
	case *BooleanValue:
		return val.Value
	case *IntegerValue:
		return val.Value != 0
	case *FloatValue:
		return val.Value != 0
	case *StringValue:
		return val.Value != ""
	default:
		return true
	}
}