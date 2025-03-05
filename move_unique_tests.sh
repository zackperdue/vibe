#!/bin/bash

# This script identifies and moves unique tests from tests/ to their respective implementation directories

# Create a backup directory structure
mkdir -p backup/tests/lexer backup/tests/parser backup/tests/interpreter

# Function to check if a test exists in destination directory
test_exists_in_dest() {
    local test_name=$1
    local dest_dir=$2

    if grep -q "^func ${test_name}" "${dest_dir}"/*_test.go 2>/dev/null; then
        return 0 # Test exists (true)
    else
        return 1 # Test doesn't exist (false)
    fi
}

# Function to extract test functions from a file
extract_unique_tests() {
    local source_file=$1
    local dest_dir=$2
    local prefix=$3

    # Get the filename without path
    local filename=$(basename "$source_file")
    local dirname=$(dirname "$source_file")
    local source_pkg=${dirname##*/}

    echo "Processing $source_file -> ${dest_dir}/${source_pkg}_${filename}"

    # Create a destination file
    local dest_file="${dest_dir}/${source_pkg}_${filename}"

    # Copy the original file to backup
    cp "$source_file" "backup/$source_file"

    # Start with package declaration
    echo "package ${dest_dir##*/}" > "$dest_file"

    # Add imports from source file
    grep -A 20 "^import (" "$source_file" | grep -v "^func" | grep -v "^package" >> "$dest_file"

    # Fix potential duplicate closing parenthesis in imports
    sed -i '' 's/)\s*)/)/g' "$dest_file"

    # Get all test functions
    local unique_tests_found=0

    # Create a temp file to hold extracted functions
    local temp_functions=$(mktemp)

    # Extract test function names and their line numbers
    grep -n "^func Test" "$source_file" | while read -r line; do
        local line_num=$(echo "$line" | cut -d: -f1)
        local test_line=$(echo "$line" | cut -d: -f2-)
        local test_name=$(echo "$test_line" | sed 's/^func \([^(]*\).*/\1/')

        # Check if test already exists in destination
        if test_exists_in_dest "$test_name" "$dest_dir"; then
            echo "  - Skipping $test_name (already exists in $dest_dir)"
            continue
        fi

        # Count unique tests found
        unique_tests_found=$((unique_tests_found + 1))

        # Find the next function's line number or EOF
        local next_func_line=$(grep -n "^func " "$source_file" | awk -F: -v start="$line_num" '$1 > start {print $1; exit}')

        if [ -z "$next_func_line" ]; then
            # No more functions, extract to end of file
            echo "  + Adding $test_name (to end of file)"
            tail -n +$line_num "$source_file" | sed "s/^func ${test_name}/func ${prefix}${test_name}/" >> "$temp_functions"
        else
            # Extract up to the next function
            echo "  + Adding $test_name (lines $line_num-$((next_func_line-1)))"
            sed -n "${line_num},$(($next_func_line-1))p" "$source_file" | sed "s/^func ${test_name}/func ${prefix}${test_name}/" >> "$temp_functions"
        fi
    done

    # Append functions to the destination file
    cat "$temp_functions" >> "$dest_file"
    rm "$temp_functions"

    # Fix imports
    if [ "$dest_dir" = "lexer" ]; then
        sed -i '' 's/"github.com\/vibe-lang\/vibe\/lexer"//g' "$dest_file"
        sed -i '' 's/lexer\.//' "$dest_file"
    elif [ "$dest_dir" = "parser" ]; then
        sed -i '' 's/"github.com\/vibe-lang\/vibe\/parser"//g' "$dest_file"
        sed -i '' 's/parser\.//' "$dest_file"
    elif [ "$dest_dir" = "interpreter" ]; then
        sed -i '' 's/"github.com\/vibe-lang\/vibe\/interpreter"//g' "$dest_file"
        sed -i '' 's/interpreter\.//' "$dest_file"
    fi

    # Clean up empty imports
    sed -i '' '/^import ()/d' "$dest_file"

    # Count unique tests
    local tests_moved=$(grep -c "^func ${prefix}Test" "$dest_file")

    # If no unique tests were found, remove the file
    if [ "$tests_moved" -eq 0 ]; then
        echo "No unique tests found in $source_file"
        rm "$dest_file"
    else
        echo "Created $dest_file with $tests_moved unique tests"
    fi
}

# Process lexer tests
echo "Processing lexer tests..."
extract_unique_tests "tests/lexer/lexer_test.go" "lexer" "Moved"
extract_unique_tests "tests/lexer/debug_test.go" "lexer" "Moved"

# Process parser tests
echo "Processing parser tests..."
extract_unique_tests "tests/parser/parser_test.go" "parser" "Moved"

# Process interpreter tests
echo "Processing interpreter tests..."
extract_unique_tests "tests/interpreter/interpreter_test.go" "interpreter" "Moved"

echo "Test migration complete!"
echo "Unique tests have been moved to their respective implementation directories."
echo "Original files are backed up in the backup/ directory."