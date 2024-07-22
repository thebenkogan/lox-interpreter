package main

import (
	"fmt"
	"os"

	"github.com/thebenkogan/lox-interpreter/internal/lexer"
	"github.com/thebenkogan/lox-interpreter/internal/parser"
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

	file, err := os.Open(args[2])
	if err != nil {
		return fmt.Errorf("Error opening file: %w", err)
	}
	tokens, errors := lexer.Tokenize(file)

	command := args[1]
	switch command {
	case "tokenize":
		for _, token := range tokens {
			fmt.Println(token.String())
		}
		if len(errors) > 0 {
			for _, error := range errors {
				fmt.Fprintln(os.Stderr, error.String())
			}
			os.Exit(65)
		}
	case "parse":
		if len(errors) > 0 {
			for _, error := range errors {
				fmt.Fprintln(os.Stderr, error.String())
			}
			os.Exit(65)
		}
		expr, err := parser.Parse(tokens)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %s\n", err.Error())
			os.Exit(65)
		}
		fmt.Println(expr.String())
	case "evaluate":
		if len(errors) > 0 {
			for _, error := range errors {
				fmt.Fprintln(os.Stderr, error.String())
			}
			os.Exit(65)
		}
		expr, err := parser.Parse(tokens)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %s\n", err.Error())
			os.Exit(65)
		}
		result, err := expr.Evaluate()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Runtime Error: %s\n", err.Error())
			os.Exit(70)
		}
		fmt.Println(result)
	default:
		return fmt.Errorf("Unknown command: %s\n", command)
	}

	return nil
}
