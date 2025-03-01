#!/bin/bash

# Change to the script directory (project root)
cd "$(dirname "$0")"

# Set up text colors for output
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

echo -e "${YELLOW}Running Vibe Language Test Suite${NC}"
echo "=================================="

# Variables to track test results
TOTAL_TESTS=0
PASSED_TESTS=0
FAILED_TESTS=0

# Function to run a test and report results
run_test() {
    local test_file=$1
    local test_name=$2

    ((TOTAL_TESTS++))

    echo -e "${BLUE}Running test: ${test_name}${NC}"
    echo -e "File: ${test_file}"

    # Run the test and capture output and exit code
    OUTPUT=$(go run main.go "${test_file}" 2>&1)
    local exit_code=$?

    # Check result
    if [ $exit_code -eq 0 ]; then
        echo -e "${GREEN}✓ Test passed${NC}"
        ((PASSED_TESTS++))
    else
        echo -e "${RED}✗ Test failed with exit code $exit_code${NC}"
        echo -e "${RED}Output:${NC}"
        echo "$OUTPUT"
        ((FAILED_TESTS++))
    fi

    echo "=================================="
}

# Run Go unit tests first
echo -e "${YELLOW}Running Go unit tests:${NC}"
./run_go_tests.sh
go_test_exit=$?

if [ $go_test_exit -ne 0 ]; then
    echo -e "${RED}⚠ Go unit tests failed. Continuing with functional tests...${NC}"
    echo "=================================="
fi

# Run all functional tests
echo -e "${YELLOW}Running functional tests:${NC}"
echo "=================================="

run_test "tests/test_basic.vi" "Basic Language Features"
run_test "tests/test_control_flow.vi" "Control Flow Structures"
run_test "tests/test_functions.vi" "Function Definitions and Calls"
run_test "tests/test_edge_cases.vi" "Edge Cases and Error Handling"
run_test "tests/test_objects.vi" "Object-Oriented Features"
run_test "tests/test_arrays.vi" "Arrays and Iteration"

# Run the example files as tests too
echo -e "${YELLOW}Running example files as tests:${NC}"
echo "=================================="

run_test "examples/minimal_fib.vi" "Minimal Fibonacci Example"
run_test "examples/fibonacci.vi" "Fibonacci Example"
run_test "examples/hello.vi" "Hello World Example"
run_test "examples/simple.vi" "Simple Program Example"

# Print summary
echo -e "${YELLOW}Test Summary:${NC}"
echo "=================================="
echo -e "Total tests: ${TOTAL_TESTS}"
echo -e "${GREEN}Passed: ${PASSED_TESTS}${NC}"
if [ $FAILED_TESTS -gt 0 ]; then
    echo -e "${RED}Failed: ${FAILED_TESTS}${NC}"
    echo "=================================="
    echo -e "${RED}⚠ Some tests failed!${NC}"
    exit 1
else
    echo -e "${GREEN}All tests passed successfully!${NC}"
    exit 0
fi