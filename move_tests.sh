#!/bin/bash

# Create backup directory
mkdir -p backup

# Function to extract test functions and move them
extract_tests() {
    local source_file=$1
    local dest_dir=$2
    local prefix=$3

    # Get the filename without path
    local filename=$(basename "$source_file")

    # Create a new file for the extracted tests
    local dest_file="${dest_dir}/${prefix}_${filename}"

    # Copy the file with the new name to the backup directory
    cp "$source_file" "backup/${prefix}_${filename}"

    # Extract the package name from the source file
    local package_name=$(grep "^package" "$source_file" | head -1 | awk '{print $2}')

    # Create the new file with the updated package name
    echo "package ${dest_dir##*/}" > "$dest_file"

    # Add imports from the source file, adjusting if needed
    grep -A 10 "^import" "$source_file" | grep -v "^func" | grep -v "^package" >> "$dest_file"

    # Extract the test functions from the source file
    grep -n "^func Test" "$source_file" | while read -r line; do
        # Get the line number and test name
        local line_num=$(echo "$line" | cut -d: -f1)
        local test_name=$(echo "$line" | cut -d: -f2 | sed 's/^func \([^(]*\).*/\1/')

        # Check if this test already exists in the destination directory
        if grep -q "^func ${test_name}" "${dest_dir}"/*_test.go 2>/dev/null; then
            # Test exists, rename it by adding the prefix
            local new_test_name="${prefix}${test_name}"

            # Extract the test function and any helper functions
            local next_func_line=$(grep -n "^func " "$source_file" | awk -F: -v start="$line_num" '$1 > start {print $1; exit}')
            if [ -z "$next_func_line" ]; then
                # No more functions, extract to the end of file
                tail -n +$line_num "$source_file" | sed "s/^func ${test_name}/func ${new_test_name}/" >> "$dest_file"
            else
                # Extract until the next function
                sed -n "${line_num},$(($next_func_line-1))p" "$source_file" | sed "s/^func ${test_name}/func ${new_test_name}/" >> "$dest_file"
            fi
        else
            # Test doesn't exist, extract it as is
            local next_func_line=$(grep -n "^func " "$source_file" | awk -F: -v start="$line_num" '$1 > start {print $1; exit}')
            if [ -z "$next_func_line" ]; then
                # No more functions, extract to the end of file
                tail -n +$line_num "$source_file" >> "$dest_file"
            else
                # Extract until the next function
                sed -n "${line_num},$(($next_func_line-1))p" "$source_file" >> "$dest_file"
            fi
        fi
    done

    # Remove any duplicate or empty imports
    local temp_file=$(mktemp)
    awk '!seen[$0]++' "$dest_file" > "$temp_file"
    mv "$temp_file" "$dest_file"

    # Fix imports if needed
    # Handle lexer imports when already in the lexer package
    if [ "${dest_dir##*/}" = "lexer" ]; then
        sed -i '' 's/"github.com\/vibe-lang\/vibe\/lexer"//g' "$dest_file"

        # Replace any "lexer." references with direct references
        sed -i '' 's/lexer\.//g' "$dest_file"
    fi

    # Handle parser imports when already in the parser package
    if [ "${dest_dir##*/}" = "parser" ]; then
        sed -i '' 's/"github.com\/vibe-lang\/vibe\/parser"//g' "$dest_file"

        # Replace any "parser." references with direct references
        sed -i '' 's/parser\.//g' "$dest_file"
    fi

    # Handle interpreter imports when already in the interpreter package
    if [ "${dest_dir##*/}" = "interpreter" ]; then
        sed -i '' 's/"github.com\/vibe-lang\/vibe\/interpreter"//g' "$dest_file"

        # Replace any "interpreter." references with direct references
        sed -i '' 's/interpreter\.//g' "$dest_file"
    fi

    # Clean up the file
    sed -i '' '/^import ()/d' "$dest_file"
    sed -i '' '/^$/d' "$dest_file"

    echo "Created $dest_file"
}

# Process lexer tests
extract_tests "./tests/lexer/lexer_test.go" "./lexer" "moved"
extract_tests "./tests/lexer/debug_test.go" "./lexer" "moved"

# Process parser tests
extract_tests "./tests/parser/parser_test.go" "./parser" "moved"

# Process interpreter tests
extract_tests "./tests/interpreter/interpreter_test.go" "./interpreter" "moved"

echo "Test files processed and moved successfully!"