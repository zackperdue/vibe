#!/bin/bash

# Colors for output
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

echo -e "${BLUE}Running Vibe Language Go tests...${NC}"
echo -e "${BLUE}=================================${NC}"

# Initialize counters
TOTAL_PACKAGES=0
PASSED_PACKAGES=0
TOTAL_TESTS=0
PASSED_TESTS=0
FAILED_TESTS=0

# Store start time
START_TIME=$(date +%s)

# Get a list of all Go packages with test files
PACKAGES=$(find . -name "*_test.go" -not -path "*/vendor/*" -not -path "*/.git/*" | xargs -n1 dirname | sort -u)

# Run tests for each package
for pkg in $PACKAGES; do
    ((TOTAL_PACKAGES++))
    pkg_name=$(basename $pkg)
    echo -e "\n${YELLOW}Testing package: ${pkg}${NC}"

    # Run tests with verbose flag and capture output
    TEST_OUTPUT=$(cd $pkg && go test -v 2>&1)
    TEST_RESULT=$?

    # Extract test counts from the output
    PKG_TESTS=$(echo "$TEST_OUTPUT" | grep -c "^=== RUN")
    TOTAL_TESTS=$((TOTAL_TESTS + PKG_TESTS))

    # Count passed and failed tests
    PKG_PASSED=$(echo "$TEST_OUTPUT" | grep -c "^--- PASS")
    PASSED_TESTS=$((PASSED_TESTS + PKG_PASSED))

    PKG_FAILED=$(echo "$TEST_OUTPUT" | grep -c "^--- FAIL")
    FAILED_TESTS=$((FAILED_TESTS + PKG_FAILED))

    # Print test output
    echo "$TEST_OUTPUT"

    # Print package result
    if [ $TEST_RESULT -eq 0 ]; then
        echo -e "${GREEN}✓ Package ${pkg_name} tests passed${NC}"
        ((PASSED_PACKAGES++))
    else
        echo -e "${RED}✗ Package ${pkg_name} had test failures${NC}"
    fi
done

# Calculate elapsed time
END_TIME=$(date +%s)
ELAPSED=$((END_TIME - START_TIME))

# Print summary
echo -e "\n${BLUE}Test Summary${NC}"
echo -e "${BLUE}============${NC}"
echo -e "Packages: ${PASSED_PACKAGES}/${TOTAL_PACKAGES} passed"
echo -e "Tests:    ${PASSED_TESTS}/${TOTAL_TESTS} passed (${FAILED_TESTS} failed)"
echo -e "Time:     ${ELAPSED} seconds"

# Set exit status
if [ $PASSED_PACKAGES -eq $TOTAL_PACKAGES ]; then
    echo -e "\n${GREEN}All tests passed!${NC}"
    exit 0
else
    echo -e "\n${RED}Some tests failed.${NC}"
    exit 1
fi