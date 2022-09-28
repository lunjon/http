# Features

## General
- [ ] `http history run N`: runs an entry from history (0 latest, etc.)
- [ ] Rename --body flag to:
  - [ ] --data: read body as string
  - [ ] --data-stdin: read body from stdin
  - [ ] --data-file: read body from file
  - [ ] --data-form: read URL encoded body
- Option for specifying output format
  - [ ] table
  - [ ] json
  - [ ] none

## Configuration
- [ ] Section in README
- [x] Timeout
- Sections
  - [x] aliases
  - [ ] history (for config, not actual history)
    - [ ] `enabled = boolean`
    - [ ] `size = int`
- Complex alias: support having the alias value:
  - [x] string: use value as is
  - [ ] table: `name = { stage = "", prod = ""}`
    - usage: `{name.stage}`

## Serve
- [x] Graceful shutdown, handle `ctrl-c`
- [x] Rename to only `serve`
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
