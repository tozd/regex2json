package r2j

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"regexp"
)

func CompileExpressions(r *regexp.Regexp) ([]*Expression, error) {
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

func Transform(r *regexp.Regexp, in io.Reader, out io.Writer, logger *log.Logger) error {
	expressions, err := CompileExpressions(r)
	if err != nil {
		return fmt.Errorf("compiling expressions: %w", err)
	}

	encoder := json.NewEncoder(out)
	encoder.SetEscapeHTML(false)

	scanner := bufio.NewScanner(in)

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
						logger.Printf(`failed to apply expression "%s" for value "%s" and line "%s": %s`, expressions[i].String(), v, line, err)
					}
				}
			}

			err := encoder.Encode(output)
			if err != nil {
				return fmt.Errorf("failed to write json: %w", err)
			}
		}
	}

	return nil
}
