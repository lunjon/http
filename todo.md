# Features

## General
- Update/change flags for body:
  - [ ] Rename `--body` to `--data` for sending string
  - [ ] Add `--data-file` for sending data from file
  - [ ] Add `--data-form` for sending URL encoded
  - [ ] Add `--data-stdin` for reading body from stdin
- Certificates
  - [ ] Add option for password
  - [ ] Add option for specifying certificate format
    - [ ] Make `--key` optional
- Option for specifying output format
  - [ ] Integrate with `--display`
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
