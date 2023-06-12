package main

import (
	"log"
	"os"
	"regexp"

	"gitlab.com/tozd/regex2json/regex2json"
)

const (
	exitSuccess = 0
	exitFailure = 1
	// 2 is used when Golang runtime fails due to an unrecovered panic or an unexpected runtime condition.
)

func main() {
	errorLogger := log.New(os.Stderr, "error: ", 0)
	warnLogger := log.New(os.Stderr, "warning: ", 0)

	if len(os.Args) != 2 {
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
