# Features

## General
- [ ] Update/change flags for body:
  - [ ] Rename `--body` to `--data` for sending string
  - [ ] Add `--data-file` for sending data from file
  - [ ] Add `--data-form` for sending URL encoded
  - [ ] Add `--data-stdin` for reading body from stdin
- Option for specifying output format
  - [ ] table
  - [ ] json
  - [ ] none

## Configuration
- [ ] Section in README
- Complex alias: support having the alias value:
  - [x] string: use value as is
  - [ ] table: `name = { stage = "", prod = ""}`
    - usage: `{name.stage}`

## Serve
- Options
  - [x] Flag for current status
  - [ ] Summary
  - [ ] Response status code
  - [ ] Response body
  - [ ] Response headers
  - [ ] Server certificate

## TUI
- [x] URL suggestions from alias file
- [x] Headers
- [x] Body
  - [x] Implement skip
  - [x] Implement editor
  - [x] Implement file
