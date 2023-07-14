package regex2json

import (
	"strings"
	"time"
)

// We determine if layout is not parsing year, month, or day by parsing layout with layout itself
// and seeing if any of those has been parsed as 0 or 1 (that is documented behavior of time.Parse
// when something is not being parsed). This works great for year (in layout it is 2006 so if it
// is parsed as 0 we know it is not parsing year) and day (in layout it is 2 so if it is parsed
// as 1 we know it) but not for month (because in layout is 1). So we have to modify month to
// February/2 in timestamps as well to see if the month is parsed as 1 or 2.
func timestampFromLayout(layout string) string {
	// There are few parts in layouts with special meaning so we remove those.
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
