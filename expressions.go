package main

import (
	"fmt"
	"strconv"
	"strings"
	"time"
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

var TimeFormats = map[string]string{
	"ANSIC":        time.ANSIC,
	"UnixDate":     time.UnixDate,
	"RubyDate":     time.RubyDate,
	"RFC822":       time.RFC822,
	"RFC822Z":      time.RFC822Z,
	"RFC850":       time.RFC850,
	"RFC1123":      time.RFC1123,
	"RFC1123Z":     time.RFC1123Z,
	"RFC3339":      time.RFC3339,
	"RFC3339Milli": "2006-01-02T15:04:05.000Z07:00",
	"RFC3339Micro": "2006-01-02T15:04:05.000000Z07:00",
	"RFC3339Nano":  time.RFC3339Nano,
	// RFC3339Nano without removing trailing zeros.
	"RFC3339NanoZeros": "2006-01-02T15:04:05.000000000Z07:00",
	"Kitchen":          time.Kitchen,
	"Stamp":            time.Stamp,
	"StampMilli":       time.StampMilli,
	"StampMicro":       time.StampMicro,
	"StampNano":        time.StampNano,
	"DateTime":         time.DateTime,
	"DateOnly":         time.DateOnly,
	"TimeOnly":         time.TimeOnly,
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
	},
	"object": func(args ...string) (Op, error) {
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
	"time": func(args ...string) (Op, error) {
		if len(args) == 0 {
			return nil, fmt.Errorf("missing format argument")
		} else if len(args) > 3 {
			return nil, fmt.Errorf("unexpected arguments: %s", strings.Join(args[3:], ", "))
		}
		parseLayout, ok := TimeFormats[args[0]]
		if !ok {
			return nil, fmt.Errorf("unknown format: %s", args[0])
		}
		formatLayout := TimeFormats["RFC3339Milli"] // Default.
		if len(args) > 1 {
			formatLayout, ok = TimeFormats[args[1]]
			if !ok {
				return nil, fmt.Errorf("unknown format: %s", args[1])
			}
		}
		var err error
		formatLocation := time.UTC // Default.
		if len(args) > 2 {
			// Capture group names in Go support only a limited set of characters.
			// So we replace the first _ with / which is common in time zone names.
			formatLocation, err = time.LoadLocation(strings.Replace(args[2], "_", "/", 1))
			if err != nil {
				return nil, err
			}
		}
		return func(in any) (any, error) {
			s, skip, err := toStringOrSkip(in)
			if err != nil {
				return nil, err
			}
			if skip {
				return in, nil
			}
			t, err := time.Parse(parseLayout, s)
			if err != nil {
				return nil, fmt.Errorf(`unable to parse "%s" into time with layout "%s" (%s): %w`, s, parseLayout, args[0], err)
			}
			return t.In(formatLocation).Format(formatLayout), nil
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
	// The first operator is the object, so we know the type of in.
	return s.merge(output, in.(map[string]any))
}

func (s Expression) merge(left map[string]any, right map[string]any) error {
	for key, rightValue := range right {
		leftValue, ok := left[key]
		if ok {
			switch lv := leftValue.(type) {
			case map[string]any:
				switch rv := rightValue.(type) {
				case map[string]any:
					// Left and right are maps. We merge them.
					err := s.merge(lv, rv)
					if err != nil {
						return fmt.Errorf("%s__%w", key, err)
					}
				case []any:
					if len(rv) == 0 {
						// Left is a map, right is an empty slice. We wrap the map into an slice.
						left[key] = []any{leftValue}
					} else {
						switch r := rv[0].(type) {
						case map[string]any:
							// Left is a map, right is a slice and the first element of right is a map.
							// We merge left into the first element of the slice and the slice is the result.
							err := s.merge(lv, r)
							if err != nil {
								return fmt.Errorf("%s__%w", key, err)
							}
							rv[0] = lv
							left[key] = rv
						default:
							// Left is a map, right is a slice and the first element of right is not a map.
							// We do not know how to merge a map with something which is not a map.
							return fmt.Errorf("%s: type mismatch", key)
						}
					}
				default:
					// Left is a map, right is not a map nor a slice. We do not know how to merge that.
					return fmt.Errorf("%s: type mismatch", key)
				}
			case []any:
				switch rv := rightValue.(type) {
				case map[string]any:
					if len(lv) == 0 {
						// Left is an empty slice, right is a map. We wrap the map into an slice.
						left[key] = []any{rightValue}
					} else {
						switch l := lv[len(lv)-1].(type) {
						case map[string]any:
							// Left is a slice and the last element of left is a map, right is a map.
							// We merge right into the last element of the slice and the slice is the result.
							err := s.merge(l, rv)
							if err != nil {
								return fmt.Errorf("%s__%w", key, err)
							}
							lv[len(lv)-1] = l
							left[key] = lv
						default:
							// Left is a slice and the last element of left is not a map, right is a map.
							// We do not know how to merge a map with something which is not a map.
							return fmt.Errorf("%s: type mismatch", key)
						}
					}
				case []any:
					// Left is a slice, right is a slice. We concatenate right to the end of left.
					left[key] = append(lv, rv...)
				default:
					// Left is a slice, right is not a map nor a slice. We append it to the end of left.
					left[key] = append(lv, rv)
				}
			default:
				switch rv := rightValue.(type) {
				case []any:
					// Left is not a map nor a slice, right is a slice. We prepend it to the start of right.
					left[key] = append([]any{lv}, rv...)
				default:
					return fmt.Errorf("%s: value already exist", key)
				}
			}
		} else {
			left[key] = rightValue
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
	// The first operator is implicitly the object. We make it explicit. We do not allow/support
	// optionally explicit first operator so that we can support "object" as field name in an object.
	// We also do not want to require that the first object operator should always be specified.
	chain[0] = "object__" + chain[0]

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
