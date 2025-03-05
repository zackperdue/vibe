package parser

import (
	"fmt"

	"github.com/vibe-lang/vibe/ast"
)

// Parse is a test helper that parses the input into a program
func (p *Parser) Parse() (*ast.Program, error) {
	program := p.parseProgram()
	if len(p.errors) > 0 {
		return program, fmt.Errorf("parser errors: %v", p.errors)
	}
	return program, nil
}