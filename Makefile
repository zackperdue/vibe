.PHONY: test test-lexer test-parser test-interpreter test-integration build clean

# Default Go command
GO = go

# Project packages
PKG_LEXER = ./lexer
PKG_PARSER = ./parser
PKG_INTERPRETER = ./interpreter
PKG_TESTS = ./tests/...

# Path to main package
MAIN_PKG = ./cmd/vibe

# Build output
OUTPUT = vibe

# Build flags
BUILD_FLAGS = -v

# Test flags
TEST_FLAGS = -v

# Run all tests
test: test-lexer test-parser test-interpreter test-integration

# Run lexer tests
test-lexer:
	@echo "Running lexer tests..."
	$(GO) test $(TEST_FLAGS) $(PKG_LEXER)
	$(GO) test $(TEST_FLAGS) ./tests/lexer

# Run parser tests
test-parser:
	@echo "Running parser tests..."
	$(GO) test $(TEST_FLAGS) $(PKG_PARSER)
	$(GO) test $(TEST_FLAGS) ./tests/parser

# Run interpreter tests
test-interpreter:
	@echo "Running interpreter tests..."
	$(GO) test $(TEST_FLAGS) $(PKG_INTERPRETER)
	$(GO) test $(TEST_FLAGS) ./tests/interpreter

# Run integration tests
test-integration:
	@echo "Running integration tests..."
	$(GO) test $(TEST_FLAGS) ./tests/integration

# Build the binary
build:
	@echo "Building $(OUTPUT)..."
	$(GO) build $(BUILD_FLAGS) -o $(OUTPUT) $(MAIN_PKG)

# Clean build artifacts
clean:
	@echo "Cleaning up..."
	rm -f $(OUTPUT)
	$(GO) clean

# Install dependencies
deps:
	$(GO) mod tidy

# Run the vibe interpreter
run: build
	./$(OUTPUT)

# Show help
help:
	@echo "Available commands:"
	@echo "  make test              - Run all tests"
	@echo "  make test-lexer        - Run lexer tests"
	@echo "  make test-parser       - Run parser tests"
	@echo "  make test-interpreter  - Run interpreter tests"
	@echo "  make test-integration  - Run integration tests"
	@echo "  make build             - Build the binary"
	@echo "  make run               - Run the vibe interpreter"
	@echo "  make clean             - Clean build artifacts"
	@echo "  make deps              - Update dependencies"
	@echo "  make help              - Show this help message"