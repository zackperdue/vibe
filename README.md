# Vibe Programming Language

![Vibe Programming Language Logo](vibe.jpg)

Vibe is a custom programming language implementation written in Go, featuring a lexer, parser, interpreter, and type system. It provides a clean, expressive syntax with strong typing.

## Features

- Strong typing system (int, float, string, bool, arrays, functions)
- Function definitions and calls with type annotations
- Parentheses-free function calls for parameter-less functions
- Arrays and array operations
- Control flow statements (if/else, while, for)
- Module system with `require` statements for code reuse
- Variables and assignment
- Arithmetic and logical operations
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

## Examples

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