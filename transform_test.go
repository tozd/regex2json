package regex2json_test

import (
	"bytes"
	"fmt"
	"log"
	"regexp"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"gitlab.com/tozd/regex2json"
)

func TestTransform(t *testing.T) {
	for i, tt := range Tests {
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			// We have to use some prefix so that no line is an empty line.
			re := "test"
			value := "test"
			for _, exp := range tt.Exps {
				re += ":"
				value += ":"
				re += "(?P<" + exp.Expression + ">.*)"
				value += exp.Value
			}
			value += "\n"
			r, err := regexp.Compile(re)
			require.NoError(t, err)
			in := bytes.Buffer{}
			_, err = in.WriteString(value)
			require.NoError(t, err)
			out := bytes.Buffer{}
			outerr := bytes.Buffer{}
			l := bytes.Buffer{}
			warnLogger := log.New(&l, "warning: ", 0)
			err = regex2json.Transform(r, &in, &out, &outerr, warnLogger)
			require.NoError(t, err)
			assert.Equal(t, tt.Expected+"\n", out.String())
			assert.Equal(t, "", outerr.String())
			lString := l.String()
			if lString != "" {
				for i, logLine := range strings.Split(strings.TrimRight(lString, "\n"), "\n") {
					if i < len(tt.Errors) {
						assert.True(t, strings.HasSuffix(logLine, tt.Errors[i]), "expected: %s, got: %s", tt.Errors[i], logLine)
					} else {
						assert.Fail(t, fmt.Sprintf("unexpected log message: %s", logLine))
					}
				}
			} else if len(tt.Errors) > 0 {
				assert.Fail(t, "expected log messages")
			}
		})
	}
}

func TestUnmatchedTransform(t *testing.T) {
	r := regexp.MustCompile(`test`)
	in := bytes.Buffer{}
	_, err := in.WriteString(`foobar`)
	require.NoError(t, err)
	out := bytes.Buffer{}
	outerr := bytes.Buffer{}
	l := bytes.Buffer{}
	warnLogger := log.New(&l, "warning: ", 0)
	err = regex2json.Transform(r, &in, &out, &outerr, warnLogger)
	require.NoError(t, err)
	assert.Equal(t, "", out.String())
	assert.Equal(t, "foobar\n", outerr.String())
	assert.Equal(t, "", l.String())
}
