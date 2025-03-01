#!/bin/bash

# Set up text colors for output
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

echo -e "${YELLOW}Running Go Tests for Crystal Implementation${NC}"
echo "=================================="

# Run tests for each package
run_package_tests() {
    local package=$1
    echo -e "${YELLOW}Testing package: $package${NC}"
    go test ./$package -v

    if [ $? -eq 0 ]; then
        echo -e "${GREEN}✓ Tests passed for $package${NC}"
    else
        echo -e "${RED}✗ Tests failed for $package${NC}"
    fi

    echo "=================================="
}

# Run tests for each package
run_package_tests "lexer"
run_package_tests "parser"
run_package_tests "interpreter"

echo -e "${YELLOW}All Go tests completed${NC}"