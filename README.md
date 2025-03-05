# Vibe Programming Language

![Vibe Programming Language Logo](vibe.jpg)

Vibe is a custom programming language implementation written in Go, featuring a lexer, parser, interpreter, and type system. It provides a clean, expressive syntax with strong typing.

## Features

- Strong typing system (int, float, string, bool, arrays, functions)
- Function definitions and calls with type annotations
- Parentheses-free function calls for parameter-less functions
- Optional return types for functions
- Flexible function syntax (can omit both parentheses and return types)
- Arrays and array operations
- Control flow statements (if/else, while, for)
- Module system with `require` statements for code reuse
- Variables and assignment
- Arithmetic and logical operations
- Comprehensive error system with informative messages
- Interactive REPL mode
- File execution
- Comments with `#`

## Installation

### Prerequisites

- Go 1.20 or later

### Building from Source

1. Clone the repository:
   ```bash
   git clone https://github.com/zackperdue/vibe.git
   cd vibe
   ```

2. Build the project:
   ```bash
   go build -o vibe
   ```

3. (Optional) Add to your PATH:
   ```bash
   export PATH=$PATH:$(pwd)
   ```

## Usage

### Running a Vibe Program

Create a file with the `.vi` extension and run it:

```bash
./vibe path/to/program.vi
# or
go run main.go path/to/program.vi
```

If you omit the `.vi` extension, it will be added automatically:

```bash
./vibe path/to/program
```

### Interactive Mode (REPL)

Start the interactive REPL to test code snippets:

```bash
./vibe -i
# or
go run main.go -i
```

The Vibe REPL (Read-Eval-Print Loop) provides an interactive environment with several helpful features:

- **Command History**: Previous commands are saved and can be accessed using up/down arrow keys
- **Persistent History**: Command history is saved between sessions in `~/.vibe_history`
- **Multi-line Input**: Code blocks spanning multiple lines are supported
- **Syntax Highlighting**: Error messages and output use color coding for improved readability
- **Error Recovery**: The REPL recovers gracefully from errors without crashing
- **Tab Completion**: Basic tab completion for common keywords and functions
- **Special Commands**:
  - `exit` or `quit`: Exit the REPL
  - `clear`: Clear the screen

Example REPL session:
```
Welcome to the Vibe programming language!
Type 'exit' to quit
>> x = 5
=> 5
>> y = 10
=> 10
>> x + y
=> 15
>> if x > 3
..  puts "x is greater than 3"
.. end
x is greater than 3
=> nil
```

### Debug Mode

Run a program with debug output to see parsing and execution details:

```bash
./vibe -d path/to/program.vi
# or
go run main.go -d path/to/program.vi
```

## Language Syntax

### Hello World

```ruby
# Hello World in Vibe
puts "Hello, World!"
```

### Variables and Types

```ruby
# Variable declaration with automatic type inference
x = 42
name = "Vibe"
pi = 3.14159
isActive = true

# With type annotations
age: int = 30
message: string = "Hello"

# The type system will enforce type safety
# x = "string" # This would cause a type error
```

### Functions

```ruby
# Function definition with type annotations
def add(a: int, b: int): int do
  return a + b
end

# Function call
result = add(5, 10)
puts result  # Outputs: 15

# Functions without parameters can be defined without parentheses
def hello: string do
  return "Hello, World!"
end

# Functions without parameters can be called without parentheses
greeting = hello
puts greeting  # Outputs: Hello, World!

# Traditional parentheses syntax also works
greeting = hello()
puts greeting  # Outputs: Hello, World!

# Functions can also be defined without a return type
def log_message do
  puts "This is a log message"
  # No return value needed
end

# Call function without return type
log_message

# Combining both features: no parentheses and no return type
def simple_logger do
  puts "Log entry created"
end

# Call the function without parentheses
simple_logger
```

### Control Flow

```ruby
# If/else statements
if x > 10
  puts "x is greater than 10"
elsif x > 5
  puts "x is greater than 5 but not greater than 10"
else
  puts "x is not greater than 5"
end

# While loops
i = 0
while i < 5 do
  puts i
  i = i + 1
end

# For loops
numbers = [1, 2, 3, 4, 5]
for num in numbers do
  puts num
end
```

### Arrays

```ruby
# Array literals
numbers = [1, 2, 3, 4, 5]

# Accessing elements (zero-indexed)
first = numbers[0]  # 1

# Modifying elements
numbers[2] = 10  # [1, 2, 10, 4, 5]
```

### Modules and Require

Vibe supports a module system with the `require` statement to include code from other files:

```ruby
# In math_utils.vi
def add(a, b) do
  return a + b
end

def multiply(a, b) do
  return a * b
end
```

```ruby
# In main.vi
require "./math_utils"

result = add(5, 3)
puts "Result of add: " + result

product = multiply(4, 7)
puts "Result of multiply: " + product
```

### Error Types and Handling

Vibe has a comprehensive error system that helps developers identify and fix issues in their code. Errors are displayed with clear messages indicating the problem and, when available, location information. The language categorizes errors into several types:

#### Parser Errors

These occur during the parsing phase when the source code doesn't conform to the language's syntax:

```
Parser errors:
    Expected token type 'END', got 'EOF' at line 5, column 1
```

Parser errors include:
- Missing closing keywords (like missing 'end' in blocks)
- Unexpected tokens
- Invalid syntax in function definitions
- Mismatched parentheses or brackets
- Invalid expression syntax

#### Type Errors

These occur when attempting operations with incompatible types:

```ruby
# Assigning a string to an int variable
age: int = "thirty"  # Type error: Cannot assign value of type STRING to variable age of type INT

# Using an operator with incompatible types
"hello" - 5  # Type error: unsupported operator - for types STRING and INTEGER

# Passing wrong argument types
def greet(name: string) do
  puts "Hello, " + name
end

greet(42)  # Type error: Parameter 'name' of function 'greet' expects STRING, got INTEGER
```

Common type errors include:
- Variable assignment type mismatches
- Function parameter type mismatches
- Function return type mismatches
- Operator type incompatibilities
- Type conversion failures
- Array index type errors (using non-integers as indices)

#### Name Resolution Errors

These occur when the interpreter cannot find a referenced variable, function, or other identifier:

```ruby
# Undefined variable
puts x  # Error: variable 'x' not found

# Using a variable before declaration
y = x + 5  # Error: variable 'x' not found

# Calling an undefined function
result = calculate_total()  # Error: function 'calculate_total' not found
```

#### Module Errors

These occur when working with the module system and the `require` statement:

```ruby
# File not found
require "./nonexistent_file"  # Error: could not load module: file not found

# Circular dependencies
# In file1.vi
require "./file2"
# In file2.vi
require "./file1"  # Error: circular module dependency detected
```

#### Array Errors

These are specific to array operations:

```ruby
# Index out of bounds
arr = [1, 2, 3]
element = arr[5]  # Error: index out of bounds: index 5 exceeds array length 3

# Invalid index type
arr = [1, 2, 3]
element = arr["one"]  # Type error: array index must be INTEGER, got STRING
```

#### Function Errors

These relate specifically to function definitions and calls:

```ruby
# Wrong number of arguments
def greet(name, age) do
  puts "Hello, " + name + "! You are " + age + " years old."
end

greet("Alice")  # Error: wrong number of arguments: expected 2, got 1

# Return type mismatch
def get_age(): int do
  return "thirty"  # Type error: function 'get_age' must return INT, got STRING
end
```

#### Runtime Errors

These occur during program execution:

```ruby
# Division by zero
x = 10 / 0  # Error: division by zero

# Accessing undefined variables
puts unknown_var  # Error: variable 'unknown_var' not found

# Object-oriented programming errors
point = null
point.x  # Error: Cannot call method on nil
```

Common runtime errors include:
- Division by zero
- Undefined variable access
- Method calls on null/nil objects
- Missing method errors
- Array index out of bounds
- Stack overflow from infinite recursion

#### How Errors are Displayed

In the REPL mode, errors are displayed immediately with descriptive messages and color-coding:
```
>> "hello" - 5
=> Type error: unsupported operator - for types STRING and INTEGER

>> if x > 10
..  puts "x is greater than 10"
..
=> Parser error: Expected token type 'END', got 'EOF' at line 3, column 1
```

In file execution mode, errors are displayed with the error type, message, and location information:
```
Error in file example.vi at line 15, column 10: division by zero
```

When running in debug mode, additional context about the error location and surrounding code may be provided:
```
Error in file example.vi at line 15, column 10: division by zero
  13 | y = 5
  14 | # This will cause an error
> 15 | result = x / 0
     |          ^^^^^
  16 | puts result
```

#### Debugging Tips

When encountering errors:
1. Check for syntax errors first (mismatched blocks, missing keywords)
2. Verify types match for all operations and assignments
3. Ensure all variables are defined before use
4. For runtime errors, add debug print statements to track variable values
Check out the example programs in the `examples/` directory:

- `hello.vi` - Basic hello world program
- `fibonacci.vi` - Fibonacci sequence calculator
- `minimal_fib.vi` - Minimal Fibonacci implementation
- `typed_program.vi` - Example demonstrating the type system
- `for_loop.vi` - Examples of for loops

## Test Suite

The `tests/` directory contains test files for various language features:

- Basic feature tests checking syntax and semantics
- Edge case tests for handling numeric operations, complex expressions
- Module tests for the `require` functionality
- Error handling tests

Run the test suite with:

```bash
./run_tests.sh    # Run all tests
./run_go_tests.sh # Run Go unit tests
```

## Project Structure

- `main.go` - Entry point for the language runtime
- `lexer/` - Tokenization of source code
- `parser/` - Converting tokens to AST
- `interpreter/` - Executing the AST
- `types/` - Type system implementation
- `examples/` - Example Vibe programs
- `tests/` - Test suite

## Development

### Development Workflow

When implementing new language features:

1. Update the lexer to recognize any new tokens
2. Extend the parser to handle the new syntax
3. Implement interpretation logic
4. Add type checking if applicable
5. Write tests to verify behavior
6. Create example programs demonstrating the feature

## Roadmap

Current development priorities:

1. Expanding the standard library
2. Improving error messages and debugging
3. Optimizing performance
4. Adding more complex data structures
5. Enhancing the module system with better error handling

## Recent Changes

- Added module system with `require` statement for including code from other files
- Implemented string concatenation with various types
- Created comprehensive test suite for edge cases
- Fixed parser issues with nested function definitions
- Improved error handling for file loading

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

Please ensure your code follows the existing style and includes appropriate tests.

## License

This project is licensed under the MIT License - see the LICENSE file for details.

## Acknowledgments

- Inspired by various programming language implementations
- Built with Go's excellent standard library

## Test Organization

Tests in this project follow Go's standard convention:

- **Component Tests**: Tests are co-located with the code they test in the same package
  - For example, lexer tests are in the `lexer` directory alongside the implementation
  - Parser tests are in the `parser` directory
  - Interpreter tests are in the `interpreter` directory

- **Integration Tests**: Tests that verify the interaction between multiple components are kept separate in the `tests/integration` directory

To run tests:

```bash
# Run all tests
./run_go_tests.sh

# Run tests for a specific package
go test ./lexer
go test ./parser
go test ./interpreter
```

Note: Some legacy tests still exist in the `tests/` directory and are gradually being moved to their implementation directories.