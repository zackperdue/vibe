package lexer_test

import (
	"fmt"
	"testing"

	"github.com/vibe-lang/vibe/lexer"
)

func TestLexerPositions(t *testing.T) {
	input := `let x = 5
y = 10`

	l := lexer.New(input)

	// Print each token with its position
	for {
		tok := l.NextToken()
		fmt.Printf("Token: %s, Type: %s, Line: %d, Column: %d\n",
			tok.Literal, tok.Type, tok.Line, tok.Column)

		// Break if we reach EOF
		if tok.Type == lexer.EOF {
			break
		}
	}
}