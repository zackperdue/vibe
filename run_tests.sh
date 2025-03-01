#!/bin/bash

# Change to the crystal directory
cd "$(dirname "$0")"

# Set up text colors for output
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

echo -e "${YELLOW}Running Crystal Language Test Suite${NC}"
echo "=================================="

# Function to run a test and report results
run_test() {
    local test_file=$1
    local test_name=$2

    echo -e "${YELLOW}Running test: ${test_name}${NC}"

    # Run the test and capture exit code
    ./crystal-lang "tests/$test_file" 2>&1
    local exit_code=$?

    # Check result
    if [ $exit_code -eq 0 ]; then
        echo -e "${GREEN}✓ Test passed${NC}"
    else
        echo -e "${RED}✗ Test failed with exit code $exit_code${NC}"
    fi

    echo "=================================="
}

# Run all tests
run_test "test_basic.crystal" "Basic Language Features"
run_test "test_control_flow.crystal" "Control Flow Structures"
run_test "test_functions.crystal" "Function Definitions and Calls"
run_test "test_edge_cases.crystal" "Edge Cases and Error Handling"
run_test "test_objects.crystal" "Object-Oriented Features"
run_test "test_arrays.crystal" "Arrays and Iteration"

# Run the example files as tests too
echo -e "${YELLOW}Running example files as tests:${NC}"
echo "=================================="

run_test "../examples/minimal_fib.crystal" "Minimal Fibonacci Example"
run_test "../examples/fibonacci.crystal" "Fibonacci Example"

echo -e "${YELLOW}All tests completed${NC}"