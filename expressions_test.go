package regex2json_test

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"gitlab.com/tozd/regex2json"
)

type ExpValue struct {
	Expression string
	Value      string
}

var Tests = []struct {
	Exps     []ExpValue
	Expected string
	Errors   []string
}{
	{[]ExpValue{{"foobar", "x"}}, `{"foobar":"x"}`, []string{}},
	{[]ExpValue{{"foo", "x"}, {"bar", "y"}}, `{"bar":"y","foo":"x"}`, []string{}},
	{[]ExpValue{{"nested__foo", "x"}, {"nested__bar", "y"}}, `{"nested":{"bar":"y","foo":"x"}}`, []string{}},
	{[]ExpValue{{"nested__foo", "x"}, {"nested__foo", "y"}}, `{"nested":{"foo":"x"}}`, []string{`nested__foo: value already exist`}},
	{[]ExpValue{{"foo___array", "x"}, {"foo___array", "y"}}, `{"foo":["x","y"]}`, []string{}},
	{[]ExpValue{{"foo___array", "x"}, {"foo", "y"}}, `{"foo":["x","y"]}`, []string{}},
	{[]ExpValue{{"foo", "x"}, {"foo___array", "y"}}, `{"foo":["x","y"]}`, []string{}},
	{[]ExpValue{{"nested__foo___array", "x"}, {"nested__foo___array", "y"}}, `{"nested":{"foo":["x","y"]}}`, []string{}},
	{[]ExpValue{{"foobar___bool", "true"}}, `{"foobar":true}`, []string{}},
	{[]ExpValue{{"foobar___int", "42"}}, `{"foobar":42}`, []string{}},
	{[]ExpValue{{"foobar___float", "42.1"}}, `{"foobar":42.1}`, []string{}},
	{[]ExpValue{{"foobar___null", ""}}, `{"foobar":null}`, []string{}},
	{[]ExpValue{{"foobar___optional", ""}}, `{}`, []string{}},
	{[]ExpValue{{"nested__foo___array___optional", ""}, {"nested__foo___array___optional", "y"}}, `{"nested":{"foo":["y"]}}`, []string{}},
	{[]ExpValue{{"nested__foo___array___null", ""}, {"nested__foo___array___optional", "y"}}, `{"nested":{"foo":[null,"y"]}}`, []string{}},
	{[]ExpValue{{"nested__foo___array___object__a", "x"}, {"nested__foo___array___object__b", "y"}}, `{"nested":{"foo":[{"a":"x"},{"b":"y"}]}}`, []string{}},
	{[]ExpValue{{"nested__foo___array___object__a", "x"}, {"nested__foo___object__b", "y"}}, `{"nested":{"foo":[{"a":"x","b":"y"}]}}`, []string{}},
	{[]ExpValue{{"nested__foo___object__b", "y"}, {"nested__foo___array___object__a", "x"}}, `{"nested":{"foo":[{"a":"x","b":"y"}]}}`, []string{}},
	{[]ExpValue{{"foo", "x"}, {"foo___array___optional", ""}}, `{"foo":["x"]}`, []string{}},
	{[]ExpValue{{"foo", "x"}, {"foo___array___int", "1"}}, `{"foo":["x",1]}`, []string{}},
	{[]ExpValue{{"foo___array___int", "1"}, {"foo", "x"}}, `{"foo":[1,"x"]}`, []string{}},
	{[]ExpValue{{"foo__bar", "x"}, {"foo___array___int", "1"}}, `{"foo":[{"bar":"x"},1]}`, []string{}},
	{[]ExpValue{{"foo__bar", "x"}, {"foo___int", "1"}}, `{"foo":{"bar":"x"}}`, []string{`foo: type mismatch`}},
	{[]ExpValue{{"foo___array___optional", ""}, {"foo", "x"}}, `{"foo":["x"]}`, []string{}},
	{[]ExpValue{{"foo___array___int", "1"}, {"foo__bar", "x"}}, `{"foo":[1,{"bar":"x"}]}`, []string{}},
	{[]ExpValue{{"foo___time__UnixDate", "Fri Jun  9 22:21:17 CEST 2023"}}, `{"foo":"2023-06-09T20:21:17.000Z"}`, []string{}},
	{[]ExpValue{{"foo___time__UnixDate__DateTime", "Fri Jun  9 22:21:17 CEST 2023"}}, `{"foo":"2023-06-09 20:21:17"}`, []string{}},
	{[]ExpValue{{"foo___time__UnixDate__DateTime__UTC__UTC", "Fri Jun  9 22:21:17 CEST 2023"}}, `{"foo":"2023-06-09 20:21:17"}`, []string{}},
	{[]ExpValue{{"foo___time__UnixDate__DateTime", "Fri Jun  9 22:21:17 MST 2023"}}, `{"foo":"2023-06-10 05:21:17"}`, []string{}},
	{[]ExpValue{{"foo___time__UnixDate__DateTime__Europe_Ljubljana", "Fri Jun  9 22:21:17 CEST 2023"}}, `{"foo":"2023-06-09 22:21:17"}`, []string{}},
	{[]ExpValue{{"foo___time__DateTime__UnixDate__UTC__Europe_Ljubljana", "2023-06-09 22:21:17"}}, `{"foo":"Fri Jun  9 20:21:17 UTC 2023"}`, []string{}},
	{[]ExpValue{{"obj___json", `{"x":1,"y":"v"}`}}, `{"obj":{"x":1,"y":"v"}}`, []string{}},
	{[]ExpValue{{"___json", `{"x":1,"y":"v"}`}}, `{"x":1,"y":"v"}`, []string{}},
	{[]ExpValue{{"obj___json___optional", ``}}, `{}`, []string{}},
	{[]ExpValue{{"___json___optional", ``}}, `{}`, []string{}},
}

func TestExpression(t *testing.T) {
	for i, tt := range Tests {
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			output := map[string]any{}
			errI := 0
			for _, ev := range tt.Exps {
				e, err := regex2json.NewExpression(ev.Expression)
				require.NoError(t, err)
				err = e.Apply(output, ev.Value)
				if err != nil && errI < len(tt.Errors) {
					assert.EqualError(t, err, tt.Errors[errI])
					errI++
					continue
				}
			}
			j, err := json.Marshal(output)
			require.NoError(t, err)
			assert.Equal(t, tt.Expected, string(j))
		})
	}
}
