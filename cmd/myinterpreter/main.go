package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"

	"github.com/thebenkogan/lox-interpreter/internal/interpreter"
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
	command := args[1]
	if command == "repl" {
		return repl()
	}

	file, err := os.Open(args[2])
	if err != nil {
		return fmt.Errorf("Error opening file: %w", err)
	}

	switch command {
	case "tokenize":
		tokens, lexerErr := lexer.Tokenize(file)
		for _, token := range tokens {
			fmt.Println(token.String())
		}
		if lexerErr != nil {
			fmt.Fprint(os.Stderr, lexerErr.Error())
			os.Exit(lexerErr.Code())
		}
	case "parse":
		tokens, lexerErr := lexer.Tokenize(file)
		if lexerErr != nil {
			fmt.Fprint(os.Stderr, lexerErr.Error())
			os.Exit(lexerErr.Code())
		}
		expr, err := parser.Parse(tokens)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %s\n", err.Error())
			os.Exit(65)
		}
		for _, statement := range expr {
			fmt.Println(statement.String())
		}
	case "execute":
		err := interpreter.Interpret(file)
		if err != nil {
			fmt.Fprint(os.Stderr, err.Error())
			os.Exit(err.Code())
		}
	default:
		return fmt.Errorf("Unknown command: %s\n", command)
	}

	return nil
}

func repl() error {
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("> ")
		line, err := reader.ReadString('\n')
		if err != nil {
			return err
		}
		err = interpreter.Interpret(bytes.NewBuffer([]byte(line)))
		if err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
		}
	}
}
