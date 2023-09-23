package regex2json

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"regexp"
)

// CompileExpressions compiles all names of named capture groups into a slice of Expressions.
// The Expression at index 0 is nil and should be skipped as it corresponds to the entire regexp match.
func CompileExpressions(r *regexp.Regexp) ([]*Expression, error) {
	expressions := make([]*Expression, 0)

	for i, expression := range r.SubexpNames() {
		// We skip the entire regexp match at index 0.
		if i == 0 {
			expressions = append(expressions, nil)
			continue
		}

		if expression == "" {
			return nil, fmt.Errorf("%w: expression missing", ErrInvalidCaptureGroup)
		}

		s, err := NewExpression(expression)
		if err != nil {
			return nil, err
		}
		expressions = append(expressions, s)
	}

	return expressions, nil
}

// Transform reads lines from in, matching every line with regexp r. If line matches, values from
// captured named groups are mapped into output JSON which is then written out to matched writer.
// If the line does not match, it is written to unmatched writer.
//
// Capture groups' names are compiled into Expressions and describe how are matched values mapped
// and transformed into output JSON. See [Expression] for details on the syntax and [Library] for
// available operators.
//
// If logger is provided, any error (e.g., a failed expression) is logged to it while the rest
// of the output JSON is still written out.
// If logger is not provided, the error is returned as error of the function, aborting the transformation.
//
// If regexp r can match multiple times per line, all matches are combined together into
// the same ome JSON output per line.
func Transform(r *regexp.Regexp, in io.Reader, matched, unmatched io.Writer, logger *log.Logger) error {
	expressions, err := CompileExpressions(r)
	if err != nil {
		return fmt.Errorf("%w: %w", ErrCompilingExpressions, err)
	}

	encoder := json.NewEncoder(matched)
	encoder.SetEscapeHTML(false)

	scanner := bufio.NewScanner(in)

	res := true
	for res {
		res = scanner.Scan()
		line := scanner.Bytes()
		if len(line) > 0 {
			output := map[string]any{}

			matches := r.FindAllSubmatch(line, -1)
			if len(matches) == 0 {
				_, err := unmatched.Write(append(line, '\n'))
				if err != nil {
					if logger != nil {
						logger.Printf(`failed to write unmatched line "%s": %s`, line, err)
					} else {
						return fmt.Errorf(`failed to write unmatched line "%s": %w`, line, err)
					}
				}
				continue
			}

			for _, match := range matches {
				for i, value := range match {
					// Nil expressions we skip.
					if expressions[i] == nil {
						continue
					}

					v := string(value)

					err := expressions[i].Apply(output, v)
					if err != nil {
						if logger != nil {
							logger.Printf(`failed to apply expression "%s" for value "%s" and line "%s": %s`, expressions[i].String(), v, line, err)
						} else {
							return fmt.Errorf(`failed to apply expression "%s" for value "%s" and line "%s": %w`, expressions[i].String(), v, line, err)
						}
					}
				}
			}

			// We do not output empty objects.
			if len(output) == 0 {
				continue
			}

			err := encoder.Encode(output)
			if err != nil {
				return fmt.Errorf("failed to write json: %w", err)
			}
		}
	}

	return nil
}
