package main

import (
	"fmt"
	"os"
	"pratt-parser/src/lexer"
)

func main() {
	bytes, _ := os.ReadFile("./examples/00.lang")

	fmt.Printf("Source code: \n%s\n\nTokens:\n", string(bytes))

	tokens := lexer.Tokenise(string(bytes))

	for _, token := range tokens {
		token.Debug()
	}
}
