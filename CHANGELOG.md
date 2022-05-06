# Changelog

This changelog adheres to semantic version according to [keepachangelog](https://keepachangelog.com/en/1.0.0/).

## Unreleased

### Added
- "User-Agent" header (does not override any specified by user)
- `--no-follow-redirects` flag for disabling redirects in HTTP client

### Fixed
- Missing output on error

### Removed
- Brief output: `--brief`

### Changed
- Renamed client certificate flags to `--cert` and `--key`

## [0.11.0] - 2021-09-24

### Added
- MIME type detection in `--body` flag if file is given
- Remove alias flag: `alias --remove name`

### Changed
- Only use environment variables for AWS credentials
- Default timeout to 30 s
- Enforce alias name pattern

### Removed
- `--aws-profile` option
- `comp` command

## [0.10.0] - 2021-09-16

### Added
- Fail flag: `--fail`/`-f`
- Brief output: `--brief`
- Client certificate flags: `--cert-pub-file` and `--cert-key-file`

## [0.9.1] - 2021-04-29

### Added
- Shorthand to `--body` flag: `-B`
- Ability to read body from stdin
- `repeat` flag

### Changed
- Rename sub-command `url` to `alias`
- Change `--version`/`-V` into sub-command
- Rename shorthand `-V` for verbose to `-v`

### Removed
- Sub-command aliases

## [0.8.0] - 2021-03-23

### Added
- `gen` command to generate shell completion files
- `url` command to create URL aliases

### Changed
- Log to stderr

## [0.7.1] - 2020-11-08

### Added
- `version` option to print version of `http`

## [0.7.0] - 2020-11-08

### Added
- Flag named `--verbose` to output debug logs
- Flag named `--response-body-only` to only output the response body
- HTTP verbs: HEAD, PUT, PATCH

### Changed
- Name of `--json` to `--body` in `http post` command. It also changed the behavior to also
  accept a filename instead of just a string of content
- Name of `--output-file` to `--output` only
- Output errors on stderr file descriptor

### Removed
- Unused features:
    * `http run`
    * `http parse-url`

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
- `--sandbox` flag to run command

## [0.2.1] - 2019-09-15

### Fixed
- Parsing URLs with query parameters

## [0.2.0] - 2019-09-15

### Added
- `parse-url` command

### Fixed
- Running get, post or delete with `--sandbox` now returns correct path

## [0.1.0] - 2019-09-14
