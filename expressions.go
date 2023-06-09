package main

import (
	"fmt"
	"strconv"
	"strings"
)

type Op = func(in any) (any, error)

type optionalType int

const (
	// Singleton for optional value. Optional value gets discarded eventually.
	optional optionalType = iota
)

func toStringOrSkip(in any) (string, bool, error) {
	s, ok := (in).(string)
	if !ok {
		if in == nil || in == optional {
			return "", true, nil
		}
		return "", false, fmt.Errorf("value is not a string")
	}
	return s, false, nil
}

var Library = map[string]func(args ...string) (Op, error){
	"int": func(args ...string) (Op, error) {
		if len(args) > 0 {
			return nil, fmt.Errorf("unexpected arguments: %s", strings.Join(args, ", "))
		}
		return func(in any) (any, error) {
			s, skip, err := toStringOrSkip(in)
			if err != nil {
				return nil, err
			}
			if skip {
				return in, nil
			}
			n, err := strconv.ParseInt(s, 10, 64)
			if err != nil {
				return nil, fmt.Errorf(`unable to parse "%s" into int: %w`, s, err)
			}
			return n, nil
		}, nil
	},
	"float": func(args ...string) (Op, error) {
		if len(args) > 0 {
			return nil, fmt.Errorf("unexpected arguments: %s", strings.Join(args, ", "))
		}
		return func(in any) (any, error) {
			s, skip, err := toStringOrSkip(in)
			if err != nil {
				return nil, err
			}
			if skip {
				return in, nil
			}
			f, err := strconv.ParseFloat(s, 64)
			if err != nil {
				return nil, fmt.Errorf(`unable to parse "%s" into float: %w`, s, err)
			}
			return f, nil
		}, nil
	},
	"bool": func(args ...string) (Op, error) {
		if len(args) > 0 {
			return nil, fmt.Errorf("unexpected arguments: %s", strings.Join(args, ", "))
		}
		return func(in any) (any, error) {
			s, skip, err := toStringOrSkip(in)
			if err != nil {
				return nil, err
			}
			if skip {
				return in, nil
			}
			b, err := strconv.ParseBool(s)
			if err != nil {
				return nil, fmt.Errorf(`unable to parse "%s" into bool: %w`, s, err)
			}
			return b, nil
		}, nil
	},
	"array": func(args ...string) (Op, error) {
		if len(args) > 0 {
			return nil, fmt.Errorf("unexpected arguments: %s", strings.Join(args, ", "))
		}
		return func(in any) (any, error) {
			// An opportunity to discard optional value.
			if in == optional {
				return []any{}, nil
			}
			return []any{in}, nil
		}, nil
	},
	"null": func(args ...string) (Op, error) {
		if len(args) > 0 {
			return nil, fmt.Errorf("unexpected arguments: %s", strings.Join(args, ", "))
		}
		return func(in any) (any, error) {
			s, skip, err := toStringOrSkip(in)
			if err != nil {
				return nil, err
			}
			if skip {
				return in, nil
			}
			if s == "" {
				return nil, nil
			}
			return s, nil
		}, nil
	},
	"optional": func(args ...string) (Op, error) {
		if len(args) > 0 {
			return nil, fmt.Errorf("unexpected arguments: %s", strings.Join(args, ", "))
		}
		return func(in any) (any, error) {
			s, skip, err := toStringOrSkip(in)
			if err != nil {
				return nil, err
			}
			if skip {
				return in, nil
			}
			if s == "" {
				return optional, nil
			}
			return s, nil
		}, nil
	}, "path": func(args ...string) (Op, error) {
		if len(args) == 0 {
			return nil, fmt.Errorf("missing path arguments")
		}
		return func(in any) (any, error) {
			// We discard optional value.
			if in == optional {
				return in, nil
			}

			res := map[string]any{}

			current := res
			for i, arg := range args {
				// Are we on the last segment of the path?
				if i == len(args)-1 {
					current[arg] = in
				} else {
					m := map[string]any{}
					current[arg] = m
					current = m
				}
			}

			return res, nil
		}, nil
	},
}

type Expression struct {
	expression string
	fns        []Op
}

func (s Expression) Apply(output map[string]any, value string) error {
	var in any = value
	var err error
	for _, f := range s.fns {
		in, err = f(in)
		if err != nil {
			return err
		}
	}
	// We discard optional value.
	if in == optional {
		return nil
	}
	// The first operator is the path, so we know the type of in.
	return s.merge(output, in.(map[string]any))
}

func (s Expression) merge(output map[string]any, update map[string]any) error {
	for key, value := range update {
		outputValue, ok := output[key]
		if ok {
			switch o := outputValue.(type) {
			case map[string]any:
				v, ok := value.(map[string]any)
				if ok {
					err := s.merge(o, v)
					if err != nil {
						return fmt.Errorf("%s__%w", key, err)
					}
				} else {
					return fmt.Errorf("%s: type mismatch", key)
				}
			case []any:
				v, ok := value.([]any)
				if ok {
					output[key] = append(o, v...)
				} else {
					return fmt.Errorf("%s: type mismatch", key)
				}
			default:
				return fmt.Errorf("%s: value already exist", key)
			}
		} else {
			output[key] = value
		}
	}

	return nil
}

func (s Expression) String() string {
	return s.expression
}

func NewExpression(expression string) (*Expression, error) {
	if expression == "" {
		return nil, fmt.Errorf(`empty expression`)
	}

	res := &Expression{
		expression: expression,
		fns:        make([]Op, 0),
	}

	chain := strings.Split(expression, "___")
	// The first operator is implicitly the path. We make it explicit. We do not allow/support
	// optionally explicit first operator so that we can support "path" as field name in an object.
	// We also do not want to require that the first path operator should always be specified.
	chain[0] = "path__" + chain[0]

	for _, c := range chain {
		ops := strings.Split(c, "__")
		functor, ok := Library[ops[0]]
		if !ok {
			return nil, fmt.Errorf(`unknown operator "%s" for expression "%s"`, ops[0], expression)
		}
		f, err := functor(ops[1:]...)
		if err != nil {
			return nil, fmt.Errorf(`compiling operator "%s" for expression "%s": %w`, ops[0], expression, err)
		}
		// We prepend the new operator, so that in Apply we call from the last to the first operator.
		res.fns = append([]Op{f}, res.fns...)
	}

	return res, nil
}
