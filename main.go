package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"regexp"
)

const (
	exitSuccess = 0
	exitFailure = 1
	// 2 is used when Golang runtime fails due to an unrecovered panic or an unexpected runtime condition.
)

func compileExpressions(r *regexp.Regexp) ([]*Expression, error) {
	expressions := make([]*Expression, 0)

	for i, expression := range r.SubexpNames() {
		// We skip the entire expression match at 0 index.
		if i == 0 {
			expressions = append(expressions, nil)
			continue
		}

		if expression == "" {
			return nil, fmt.Errorf("capture group without expression")
		}

		s, err := NewExpression(expression)
		if err != nil {
			return nil, err
		}
		expressions = append(expressions, s)
	}

	return expressions, nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintf(os.Stderr, "error: invalid number of arguments, got %d, expected 1\n", len(os.Args)-1)
		os.Exit(exitFailure)
	}

	r, err := regexp.Compile(os.Args[1])
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: invalid regexp: %s\n", err)
		os.Exit(exitFailure)
	}

	expressions, err := compileExpressions(r)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: compiling expressions: %s\n", err)
		os.Exit(exitFailure)
	}

	encoder := json.NewEncoder(os.Stdout)
	encoder.SetEscapeHTML(false)

	scanner := bufio.NewScanner(os.Stdin)

	res := true
	for res {
		res = scanner.Scan()
		line := scanner.Bytes()
		if len(line) > 0 {
			output := map[string]any{}

			matches := r.FindAllSubmatch(line, -1)
			for _, match := range matches {
				for i, value := range match {
					// Nil expressions we skip.
					if expressions[i] == nil {
						continue
					}

					v := string(value)

					err := expressions[i].Apply(output, v)
					if err != nil {
						fmt.Fprintf(os.Stderr, "warning: failed to apply expression \"%s\" for value \"%s\" and line \"%s\": %s\n", expressions[i].String(), v, line, err)
					}
				}
			}

			err := encoder.Encode(output)
			if err != nil {
				fmt.Fprintf(os.Stderr, "error: failed to write json: %s\n", err)
				os.Exit(exitFailure)
			}
		}
	}

	os.Exit(exitSuccess)
}
