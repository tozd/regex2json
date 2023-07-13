package regex2json

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

var testLayouts = map[string]string{
	"ANSIC":       time.ANSIC,
	"UnixDate":    time.UnixDate,
	"RubyDate":    time.RubyDate,
	"RFC822":      time.RFC822,
	"RFC822Z":     time.RFC822Z,
	"RFC850":      time.RFC850,
	"RFC1123":     time.RFC1123,
	"RFC1123Z":    time.RFC1123Z,
	"RFC3339":     time.RFC3339,
	"RFC3339Nano": time.RFC3339Nano,
	"Kitchen":     time.Kitchen,
	"Stamp":       time.Stamp,
	"StampMilli":  time.StampMilli,
	"StampMicro":  time.StampMicro,
	"StampNano":   time.StampNano,
	"DateTime":    time.DateTime,
	"DateOnly":    time.DateOnly,
	"TimeOnly":    time.TimeOnly,
}

func TestLayoutsWithoutYear(t *testing.T) {
	assert.Equal(t, map[string]bool{
		"Kitchen":    true,
		"Stamp":      true,
		"StampMicro": true,
		"StampMilli": true,
		"StampNano":  true,
		"TimeOnly":   true,
	}, layoutsWithoutYear(testLayouts))
}

func TestLayoutsWithoutMonth(t *testing.T) {
	assert.Equal(t, map[string]bool{
		"Kitchen":  true,
		"TimeOnly": true,
	}, layoutsWithoutMonth(testLayouts))
}

func TestLayoutsWithoutDay(t *testing.T) {
	assert.Equal(t, map[string]bool{
		"Kitchen":  true,
		"TimeOnly": true,
	}, layoutsWithoutDay(testLayouts))
}
