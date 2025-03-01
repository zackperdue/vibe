#!/bin/bash

# Define colors for output
GREEN='\033[0;32m'
RED='\033[0;31m'
ORANGE='\033[0;33m'
NC='\033[0m' # No Color

echo "Running tests for Vibe programming language..."
echo ""

# Run all tests with verbose output
go test -v ./... | grep -v "=== RUN" | grep -v "--- PASS" | grep -v "PASS" | grep -v "FAIL" | grep -v "ok" | grep -v "?"

# Get exit status
test_status=${PIPESTATUS[0]}

if [ $test_status -eq 0 ]; then
    echo -e "${GREEN}All tests passed!${NC}"
else
    echo -e "${RED}Tests failed with exit status: $test_status${NC}"
fi

exit $test_status