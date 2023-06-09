package main

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type ExpValue struct {
	Expression string
	Value      string
}

func TestExpressions(t *testing.T) {
	tests := []struct {
		Exps     []ExpValue
		Expected string
	}{
		{[]ExpValue{{"foobar", "x"}}, `{"foobar":"x"}`},
		{[]ExpValue{{"foo", "x"}, {"bar", "y"}}, `{"bar":"y","foo":"x"}`},
		{[]ExpValue{{"nested__foo", "x"}, {"nested__bar", "y"}}, `{"nested":{"bar":"y","foo":"x"}}`},
		{[]ExpValue{{"nested__foo", "x"}, {"nested__foo", "y"}}, `nested__foo: value already exist`},
		{[]ExpValue{{"foo___array", "x"}, {"foo___array", "y"}}, `{"foo":["x","y"]}`},
		{[]ExpValue{{"foo___array", "x"}, {"foo", "y"}}, `foo: type mismatch`},
		{[]ExpValue{{"nested__foo___array", "x"}, {"nested__foo___array", "y"}}, `{"nested":{"foo":["x","y"]}}`},
		{[]ExpValue{{"foobar___bool", "true"}}, `{"foobar":true}`},
		{[]ExpValue{{"foobar___int", "42"}}, `{"foobar":42}`},
		{[]ExpValue{{"foobar___float", "42.1"}}, `{"foobar":42.1}`},
		{[]ExpValue{{"nested__foo___array___path__a", "x"}, {"nested__foo___array___path__b", "y"}}, `{"nested":{"foo":[{"a":"x"},{"b":"y"}]}}`},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			output := map[string]any{}
			for _, ev := range tt.Exps {
				e, err := NewExpression(ev.Expression)
				require.NoError(t, err)
				err = e.Apply(output, ev.Value)
				if err != nil {
					assert.Equal(t, tt.Expected, err.Error())
					return
				}
			}
			j, err := json.Marshal(output)
			require.NoError(t, err)
			assert.Equal(t, tt.Expected, string(j))
		})
	}
}
