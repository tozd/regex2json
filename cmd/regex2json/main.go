// regex2json reads lines from stdin, matching every line with the provided regexp.
// If line matches, values from captured named groups are mapped into output JSON
// which is then written out to stdout.
//
// Capture groups' names are compiled into Expressions and describe how are matched
// values mapped and transformed into output JSON. See [regex2json.Expression] for
// details on the syntax and [regex2json.Library] for available operators.
//
// Any failed expression is logged to stderr while the rest of the output JSON is still
// written out.
//
// If regexp can match multiple times per line, all matches are combined together
// into the same one JSON output per line.
//
// Usage:
//
//	regex2json <regexp>
//
// Example:
//
//	$ while true; do LC_TIME=C date; sleep 1; done | regex2json "(?P<date___time__UnixDate__RFC3339>.+)"
//	{"date":"2023-06-13T11:26:45Z"}
//	{"date":"2023-06-13T11:26:46Z"}
//	{"date":"2023-06-13T11:26:47Z"}
//
// Example:
//
// 	$ echo '192.168.0.100 - - [13/Jun/2023:13:15:13 +0000] "GET /index.html HTTP/1.1" 200 1234 "-" "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36"' | regex2json '^(?P<address>\S+) - (?P<user>\S+) \[(?P<time___time__Nginx__RFC3339>[\w:/]+\s[+\-]\d{4})\] "(?P<method>\S+)\s?(?P<url>\S+)?\s?(?P<http>\S+)?" (?P<status___int>\d{3}) (?P<size___int>\d+|-) "(?P<referrer>[^"]*)" "(?P<agent>[^"]*)"'
// 	{"address":"192.168.0.100","agent":"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36","http":"HTTP/1.1","method":"GET","referrer":"-","size":1234,"status":200,"time":"2023-06-13T13:15:13Z","url":"/index.html","user":"-"}

package main

import (
	"log"
	"os"
	"regexp"

	"gitlab.com/tozd/regex2json"
)

const (
	exitSuccess = 0
	exitFailure = 1
	// 2 is used when Golang runtime fails due to an unrecovered panic or an unexpected runtime condition.
)

func main() {
	errorLogger := log.New(os.Stderr, "error: ", 0)
	warnLogger := log.New(os.Stderr, "warning: ", 0)

	if len(os.Args) != 2 { //nolint:gomnd
		errorLogger.Printf("invalid number of arguments, got %d, expected 1", len(os.Args)-1)
		os.Exit(exitFailure)
	}

	r, err := regexp.Compile(os.Args[1])
	if err != nil {
		errorLogger.Printf("invalid regexp: %s", err)
		os.Exit(exitFailure)
	}

	err = regex2json.Transform(r, os.Stdin, os.Stdout, warnLogger)
	if err != nil {
		errorLogger.Printf("%s", err)
		os.Exit(exitFailure)
	}

	os.Exit(exitSuccess)
}
