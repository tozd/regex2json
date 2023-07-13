package regex2json

import (
	"strings"
	"time"
)

func timestampFromLayout(layout string) string {
	timestamp := strings.ReplaceAll(layout, "_", "")
	timestamp = strings.ReplaceAll(timestamp, "Z07:00", "Z")
	timestamp = strings.ReplaceAll(timestamp, "Z0700", "Z")
	// We first change all months to February.
	timestamp = strings.ReplaceAll(timestamp, "1", "2")
	// But we might have changed 15 to 25, so we revert it back.
	timestamp = strings.ReplaceAll(timestamp, "25", "15")
	timestamp = strings.ReplaceAll(timestamp, "Jan", "Feb")
	timestamp = strings.ReplaceAll(timestamp, "January", "February")
	return timestamp
}

func layoutsWithoutYear(layouts map[string]string) map[string]bool {
	output := map[string]bool{}

	for name, layout := range layouts {
		timestamp := timestampFromLayout(layout)
		t, err := time.Parse(layout, timestamp)
		if err != nil {
			panic(err)
		}
		if t.Year() == 0 {
			output[name] = true
		}
	}

	return output
}

func layoutsWithoutMonth(layouts map[string]string) map[string]bool {
	output := map[string]bool{}

	for name, layout := range layouts {
		timestamp := timestampFromLayout(layout)
		t, err := time.Parse(layout, timestamp)
		if err != nil {
			panic(err)
		}
		if t.Month() == 1 {
			output[name] = true
		}
	}

	return output
}

func layoutsWithoutDay(layouts map[string]string) map[string]bool {
	output := map[string]bool{}

	for name, layout := range layouts {
		timestamp := timestampFromLayout(layout)
		t, err := time.Parse(layout, timestamp)
		if err != nil {
			panic(err)
		}
		if t.Day() == 1 {
			output[name] = true
		}
	}

	return output
}
