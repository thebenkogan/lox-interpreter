package main

import (
	"fmt"
	"os"

	"github.com/codecrafters-io/interpreter-starter-go/internal/lexer"
)

func main() {
	if err := run(os.Args); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run(args []string) error {
	if len(args) < 3 {
		return fmt.Errorf("Usage: ./your_program.sh tokenize <filename>")
	}

	command := args[1]
	switch command {
	case "tokenize":
		file, err := os.Open(args[2])
		if err != nil {
			return fmt.Errorf("Error opening file: %w", err)
		}
		tokens, errors := lexer.Tokenize(file)
		for _, token := range tokens {
			fmt.Println(token.String())
		}
		if len(errors) > 0 {
			for _, error := range errors {
				fmt.Fprintln(os.Stderr, error.String())
			}
			os.Exit(65)
		}
	default:
		return fmt.Errorf("Unknown command: %s\n", command)
	}

	return nil
}
