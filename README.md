# Convert text to JSON using only regular expressions

[![pkg.go.dev](https://pkg.go.dev/badge/gitlab.com/tozd/regex2json)](https://pkg.go.dev/gitlab.com/tozd/regex2json)
[![Go Report Card](https://goreportcard.com/badge/gitlab.com/tozd/regex2json)](https://goreportcard.com/report/gitlab.com/tozd/regex2json)
[![pipeline status](https://gitlab.com/tozd/regex2json/badges/main/pipeline.svg?ignore_skipped=true)](https://gitlab.com/tozd/regex2json/-/pipelines)
[![coverage report](https://gitlab.com/tozd/regex2json/badges/main/coverage.svg)](https://gitlab.com/tozd/regex2json/-/graphs/main/charts)

Main motivation for this tool is to convert traditional text-based (and line-based) logs to JSON
for programs which do not support JSON logs themselves.
It can be used in online manner (pipelining output of the program into regex2json) or offline manner
(to process logs stored in files). But the tool is more general and can enable any workflow where
you prefer operating on JSON instead of text. It works especially great when combined with
[jq](https://jqlang.github.io/jq/).

Features:

- Reads stdin line by line, converting each line to JSON to stdout.
- Supports transformations of matched capture groups by specifying the transformation as capture group's name.
- Transformation consists of a series of operators (e.g., parsing numbers, timestamps, creating arrays and objects).
- Supports regexp matching a line multiple times, combining all matches into one JSON.

## Installation

This is a tool implemented in Go. You can use `go install` to install the latest stable (released) version:

```sh
go install gitlab.com/tozd/regex2json/cmd/regex2json@latest
```

[Releases page](https://gitlab.com/tozd/regex2json/-/releases)
contains a list of stable versions. Each includes:

- Statically compiled binaries.
- Docker images.

To install the latest development version (`main` branch):

```sh
go install gitlab.com/tozd/regex2json/cmd/regex2json@main
```

## Usage

regex2json reads lines from stdin, matching every line with the provided regexp.
If line matches, values from captured named groups are mapped into output JSON
which is then written out to stdout.

Capture groups' names are compiled into Expressions and describe how are matched
values mapped and transformed into output JSON. See
[Expression](https://pkg.go.dev/gitlab.com/tozd/regex2json#Expression)
for details on the syntax and
[Library](https://pkg.go.dev/gitlab.com/tozd/regex2json#Library)
for available operators.

Any failed expression is logged to stderr while the rest of the output JSON is still
written out.

If regexp can match multiple times per line, all matches are combined together
into the same one JSON output per line.

Usage:

```sh
regex2json <regexp>
```

Example:

```sh
$ while true; do LC_TIME=C date; sleep 1; done | regex2json "(?P<date___time__UnixDate__RFC3339>.+)"
{"date":"2023-06-13T11:26:45Z"}
{"date":"2023-06-13T11:26:46Z"}
{"date":"2023-06-13T11:26:47Z"}
```

### As a package

See full package documentation on [pkg.go.dev](https://pkg.go.dev/badge/gitlab.com/tozd/regex2json)
on using regex2json as a Go package.

## Contributing

Feel free to make a merge-request add more time layouts and/or operators.

## Related projects

- [jc](https://github.com/kellyjonbrazil/jc) – jc enables the same idea of converting text-based output of
  programs into JSON, but its focus is to support popular programs out of the box. regex2json enables quick
  transformations by providing a regexp with expressions how captured groups are transformed into JSON.

## GitHub mirror

There is also a [read-only GitHub mirror available](https://github.com/tozd/regex2json),
if you need to fork the project there.
