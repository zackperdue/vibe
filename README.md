# Vibe Programming Language

![Vibe Programming Language Logo](vibe.jpg)

Vibe is a custom programming language implementation written in Go, featuring a lexer, parser, interpreter, and type system. It provides a clean, expressive syntax with strong typing.

## Features

- Strong typing system (int, float, string, bool, arrays, functions)
- Function definitions and calls with type annotations
- Arrays and array operations
- Control flow statements (if/else)
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
   git clone https://github.com/example/vibe.git
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

# The type system will enforce type safety
# x = "string" # This would cause a type error
```

### Functions

```ruby
# Function definition with type annotations
fn add(a: int, b: int): int {
  return a + b
}

# Function call
result = add(5, 10)
puts result  # Outputs: 15
```

### Control Flow

```ruby
# If/else statements
if x > 10 {
  puts "x is greater than 10"
} else {
  puts "x is not greater than 10"
}
```

### Arrays

```ruby
# Array literals
numbers = [1, 2, 3, 4, 5]

# Accessing elements (zero-indexed)
first = numbers[0]  # 1
```

## Examples

Check out the example programs in the `examples/` directory:

- `hello.vi` - Basic hello world program
- `fibonacci.vi` - Fibonacci sequence calculator
- `minimal_fib.vi` - Minimal Fibonacci implementation
- `typed_program.vi` - Example demonstrating the type system

## Project Structure

- `main.go` - Entry point for the language runtime
- `lexer/` - Tokenization of source code
- `parser/` - Converting tokens to AST
- `interpreter/` - Executing the AST
- `types/` - Type system implementation
- `examples/` - Example Vibe programs
- `tests/` - Test suite

## Development

### Running Tests

```bash
./run_tests.sh    # Run all tests
./run_go_tests.sh # Run Go unit tests
```

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
5. Implementing modules/imports

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