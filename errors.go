package regex2json

import (
	"errors"
)

var (
	ErrUnexpectedType       = errors.New("unexpected type")
	ErrTypeMismatch         = errors.New("type mismatch")
	ErrUnexpectedArgument   = errors.New("unexpected argument")
	ErrInvalidValue         = errors.New("invalid value")
	ErrValueAlreadyExist    = errors.New("value already exist")
	ErrMissingArgument      = errors.New("missing argument")
	ErrInvalidCaptureGroup  = errors.New("invalid capture group")
	ErrCompilingExpressions = errors.New("compiling expressions")
	ErrEmptyExpression      = errors.New("empty expression")
	ErrEmptyOperator        = errors.New("empty operator")
	ErrInvalidOperator      = errors.New("invalid operator")
	ErrCompilingOperator    = errors.New("compiling operator")
)
