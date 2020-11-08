# Changelog

This changelog adheres to semantic version according to [keepachangelog](https://keepachangelog.com/en/1.0.0/).

## Unreleased

## [0.7.1] - 2020-11-08

### Added
- `version` command to print version of `httpreq`.

## [0.7.0] - 2020-11-08

### Added
- Flag named `--verbose` to output debug logs.
- Flag named `--response-body-only` to only output the response body.
- HTTP verbs: HEAD, PUT, PATCH

### Changed
- Name of `--json` to `--body` in `httpreq post` command. It also changed the behavior to also
  accept a filename instead of just a string of content.
- Name of `--output-file` to `--output` only.
- Output errors on stderr file descriptor

### Removed
- Unused features:
    * `httpreq run`
    * `httpreq parse-url`

## [0.6.0] - 2019-12-17

### Added
- `--timeout (short version: -T)` flag
- `-H` as short version for `--header`

## [0.5.0] - 2019-10-07

### Fixed
- `--header` arguments did not end up in the requests

## [0.4.0] - 2019-09-15

### Added
- `env` section in spec files. Allows user to use define and use environment variables in headers and URLs.

## [0.3.0] - 2019-09-15

### Added
- `--sandbox` flag to run command.

## [0.2.1] - 2019-09-15

### Fixed
- Parsing URLs with query parameters.

## [0.2.0] - 2019-09-15

### Added
- `parse-url` command.

### Fixed
- Running get, post or delete with `--sandbox` now returns correct path.

## [0.1.0] - 2019-09-14
