// Package regex2json enables extracting data from text into JSON using just regular expressions.
//
// Expressions how to transform matched values into data are defined as capture groups' names.
// Expressions can consist from a series of operators, called one after the other.
package regex2json

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/tkuchiki/go-timezone"
)

// Op is the operator's function type.
// Operator can receive any type and return any type.
// It can error.
type Op = func(in any) (any, error)

type optionalType int

const (
	// Singleton for optional value. Optional value gets discarded eventually.
	optional optionalType = iota
)

var tz = timezone.New()

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

// TimeLayouts is a map of time layouts supported by [TimeOperator].
// RFC3339NanoZeros is the same as RFC3339Nano but without removing trailing zeros.
var TimeLayouts = map[string]string{
	"ANSIC":                   time.ANSIC,
	"UnixDate":                time.UnixDate,
	"RubyDate":                time.RubyDate,
	"RFC822":                  time.RFC822,
	"RFC822Z":                 time.RFC822Z,
	"RFC850":                  time.RFC850,
	"RFC1123":                 time.RFC1123,
	"RFC1123Z":                time.RFC1123Z,
	"RFC3339":                 time.RFC3339,
	"RFC3339Milli":            "2006-01-02T15:04:05.000Z07:00",
	"RFC3339Micro":            "2006-01-02T15:04:05.000000Z07:00",
	"RFC3339Nano":             time.RFC3339Nano,
	"RFC3339NanoZeros":        "2006-01-02T15:04:05.000000000Z07:00",
	"Kitchen":                 time.Kitchen,
	"Stamp":                   time.Stamp,
	"StampMilli":              time.StampMilli,
	"StampMicro":              time.StampMicro,
	"StampNano":               time.StampNano,
	"DateTime":                time.DateTime,
	"DateOnly":                time.DateOnly,
	"TimeOnly":                time.TimeOnly,
	"Nginx":                   "02/Jan/2006:15:04:05 -0700",
	"LogDateTime":             "2006/01/02 15:04:05",
	"LogDateOnly":             "2006/01/02",
	"LogDateTimeMicroseconds": "2006/01/02 15:04:05.000000",
	"LogTimeMicroseconds":     "15:04:05.000000",
}

// IntOperator returns the bool operator which parses input string
// into an int value using [strconv.ParseInt].
//
// It does not expect any arguments.
func IntOperator(args ...string) (Op, error) {
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
}

// FloatOperator returns the bool operator which parses input string
// into a float value using [strconv.ParseFloat].
//
// It does not expect any arguments.
func FloatOperator(args ...string) (Op, error) {
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
}

// BoolOperator returns the bool operator which parses input string
// into a bool value using [strconv.ParseBool].
//
// It does not expect any arguments.
func BoolOperator(args ...string) (Op, error) {
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
}

// ArrayOperator returns the array operator which wraps the input value
// into an array with one element, the input value.
//
// It does not expect any arguments.
func ArrayOperator(args ...string) (Op, error) {
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
}

// NullOperator returns the null operator which returns null if the input
// is an empty string.
//
// It does not expect any arguments.
func NullOperator(args ...string) (Op, error) {
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
			return nil, nil //nolint:nilnil
		}
		return s, nil
	}, nil
}

// OptionalOperator returns the optional operator which returns no value if the input
// is an empty string.
//
// It does not expect any arguments.
func OptionalOperator(args ...string) (Op, error) {
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
}

// ObjectOperator returns the object operator which constructs an (possibly nested)
// object based on provided path as arguments. E.g., calling it with arguments foo
// and bar will return an object {"foo": {"bar": <in>}}.
func ObjectOperator(args ...string) (Op, error) {
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
}

// TimeOperator returns the time operator which parses the input string
// into a timestamp and then formats the timestamp back into a string.
//
// It accepts four arguments, in order:
//
//   - parsing layout (required)
//   - formatting layout (default RFC3339Milli)
//   - formatting location (default [time.UTC])
//   - parsing location (default [time.Local])
func TimeOperator(args ...string) (Op, error) {
	if len(args) == 0 {
		return nil, fmt.Errorf("missing parse layout argument")
	} else if len(args) > 4 { //nolint:gomnd
		return nil, fmt.Errorf("unexpected arguments: %s", strings.Join(args[4:], ", "))
	}
	parseLayout, ok := TimeLayouts[args[0]]
	if !ok {
		return nil, fmt.Errorf("unknown format: %s", args[0])
	}
	formatLayout := TimeLayouts["RFC3339Milli"] // Default.
	if len(args) > 1 {
		formatLayout, ok = TimeLayouts[args[1]]
		if !ok {
			return nil, fmt.Errorf("unknown format layout: %s", args[1])
		}
	}
	var err error
	formatLocation := time.UTC // Default.
	if len(args) > 2 {         //nolint:gomnd
		// Capture group names in Go support only a limited set of characters.
		// So we replace the first _ with / which is common in time zone names.
		// See: https://github.com/golang/go/issues/60784
		formatLocation, err = time.LoadLocation(strings.Replace(args[2], "_", "/", 1))
		if err != nil {
			return nil, err
		}
	}
	parseLocation := time.Local // Default.
	if len(args) > 3 {          //nolint:gomnd
		// Capture group names in Go support only a limited set of characters.
		// So we replace the first _ with / which is common in time zone names.
		// See: https://github.com/golang/go/issues/60784
		parseLocation, err = time.LoadLocation(strings.Replace(args[3], "_", "/", 1))
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
		t, err := time.ParseInLocation(parseLayout, s, parseLocation)
		if err != nil {
			return nil, fmt.Errorf(`unable to parse "%s" into time with layout "%s" (%s) in location "%s": %w`, s, parseLayout, args[0], parseLocation, err)
		}
		// Parsing might not succeed in using timezone abbreviation when present (when it does not match parseLocation).
		// In such case time.ParseInLocation uses a fabricated location with the given timezone abbreviation and a zero
		// offset. We try to obtain correct location from timezone abbreviation and parse again in that location.
		zone, offset := t.Zone()
		if t.Location() != parseLocation && offset == 0 {
			l, err := time.LoadLocation(zone)
			if err == nil {
				t, err = time.ParseInLocation(parseLayout, s, l)
				if err != nil {
					return nil, fmt.Errorf(`unable to parse "%s" into time with layout "%s" (%s) in location "%s": %w`, s, parseLayout, args[0], l, err)
				}
			} else {
				zones, err := tz.GetTimezones(zone)
				if err != nil {
					return nil, fmt.Errorf(`unable to parse "%s" into time with layout "%s" (%s): unable to parse timezone "%s": %w`, s, parseLayout, args[0], zone, err)
				}
				found := false
				for _, z := range zones {
					l, err := time.LoadLocation(z)
					if err == nil {
						t, err = time.ParseInLocation(parseLayout, s, l)
						if err != nil {
							return nil, fmt.Errorf(`unable to parse "%s" into time with layout "%s" (%s) in location "%s": %w`, s, parseLayout, args[0], l, err)
						}
						found = true
						break
					}
				}
				if !found {
					return nil, fmt.Errorf(`unable to parse "%s" into time with layout "%s" (%s): unable to parse timezone "%s"`, s, parseLayout, args[0], zone)
				}
			}
		}
		return t.In(formatLocation).Format(formatLayout), nil
	}, nil
}

// Library is a map of all supported operators.
var Library = map[string]func(args ...string) (Op, error){
	"int":      IntOperator,
	"float":    FloatOperator,
	"bool":     BoolOperator,
	"array":    ArrayOperator,
	"null":     NullOperator,
	"optional": OptionalOperator,
	"object":   ObjectOperator,
	"time":     TimeOperator,
}

// Expression is a compiled expression which can be applied on a value
// to transforms it by calling operators one after the other.
//
// The syntax of an expression consists of a series of operator names,
// each name followed by possible arguments. Operator's name and its
// arguments are separated by __ (double underscore). Operators (and
// their arguments) themselves are separated by ___ (triple underscore).
// The first operator is implicitly object and its name should not be
// provided.
//
// Example:
//
//	foo__bar___time__UnixDate
//
// Corresponds to:
//
//	object("foo", "bar")(time("UnixDate")(<in>))
//
// The object returns a function which transforms
// its input into object {"foo": {"bar": <in>}} and time returns a function
// which parses its input according to UnixDate layout and formats it back
// to (default) RFC3339Milli layout. The formatted time is thus stored in
// the object. E.g., for input "Fri Jun  9 22:21:17 CEST 2023" the output
// is {"foo": {"bar": "2023-06-09T20:21:17.000Z"}}.
type Expression struct {
	expression string
	fns        []Op
}

// Apply runs the Expression on the value and transforms it by calling operators
// one after the other. Expression always returns an object which is then merged
// into output.
//
// Merging of objects is recursive and is designed so that multiple Expressions
// can be applied using the same output, which collects results from those Expressions.
//
// Objects are merged by merging their fields. Arrays are merged by concatenating them.
// Object is merged with an array by merging the object with the first element of the
// array, when it is the first element is an object. Otherwise the object is prepended
// to the array. Array is merged with an object in a similar way, only merging is done
// with the last element of the array (if it is an object) or the object is appended
// to the array. Merging other values with the array prepends them. Merging an array
// with other values appends them.
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
							// We prepend the map to the slice.
							left[key] = append([]any{lv}, rv...)
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
							// We append the map to the slice.
							left[key] = append(lv, rv)
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

// String returns the original expression used to compile this Expression.
func (s Expression) String() string {
	return s.expression
}

// NewExpression compiles the expression into the Expression.
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
