# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

## [0.10.0] - 2023-07-14

### Added

- Add `DayMonthYearTime` and `DayMonthYearTimeMilli` time layouts.

## [0.9.0] - 2023-07-13

### Added

- Add `DayMonthTime`, `DayMonthTimeMilli`, `WeekDayMonthDayTime`, and
  `WeekDayMonthDayTimeMilli` time layouts.
- If time layout is missing year, month, or day those are filled with the current time.

## [0.8.0] - 2023-07-13

### Added

- Add `DateTimeMilli` time layout.

## [0.7.0] - 2023-07-12

### Changed

- Do not output empty objects.

## [0.6.0] - 2023-07-12

### Added

- Add `LogPostgreSQL` time layout.

## [0.5.0] - 2023-07-11

### Added

- Add `json` operator to parse JSON strings.

## [0.4.0] - 2023-07-11

### Added

- Add `ISO8601`, `ISO8601Milli`,`ISO8601Micro`, `ISO8601Nano`, and `ISO8601NanoZeros`
  time layouts.

## [0.3.0] - 2023-07-10

### Added

- Add `LogDateTime`, `LogDateOnly`, `LogDateTimeMicroseconds`, and `LogTimeMicroseconds`
  time layouts to parse standard Go log timestamps.

## [0.2.0] - 2023-07-09

### Changed

- Lines not matching the regexp are now written to stderr instead of being ignored.

## [0.1.0] - 2023-06-13

### Added

- First public release.

[unreleased]: https://gitlab.com/tozd/regex2json/-/compare/v0.10.0...main
[0.10.0]: https://gitlab.com/tozd/regex2json/-/compare/v0.9.0...v0.10.0
[0.9.0]: https://gitlab.com/tozd/regex2json/-/compare/v0.8.0...v0.9.0
[0.8.0]: https://gitlab.com/tozd/regex2json/-/compare/v0.7.0...v0.8.0
[0.7.0]: https://gitlab.com/tozd/regex2json/-/compare/v0.6.0...v0.7.0
[0.6.0]: https://gitlab.com/tozd/regex2json/-/compare/v0.5.0...v0.6.0
[0.5.0]: https://gitlab.com/tozd/regex2json/-/compare/v0.4.0...v0.5.0
[0.4.0]: https://gitlab.com/tozd/regex2json/-/compare/v0.3.0...v0.4.0
[0.3.0]: https://gitlab.com/tozd/regex2json/-/compare/v0.2.0...v0.3.0
[0.2.0]: https://gitlab.com/tozd/regex2json/-/compare/v0.1.0...v0.2.0
[0.1.0]: https://gitlab.com/tozd/regex2json/-/tags/v0.1.0

<!-- markdownlint-disable-file MD024 -->
