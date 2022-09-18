# Features

## General
- Request history
  - [x] Write to file
  - Sub-commands:
    - [x] `http history`: base form, lists history
    - [x] `http history clear`: clears history
    - [ ] `http history run N`: runs an entry from history (0 latest, etc.)
  - [ ] Configuration section: `[history]`
    - [ ] `enabled = boolean`
    - [ ] `size = int`
  - [ ] Truncate to x entries
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
- Complex alias: support having the alias value:
  - [x] string: use value as is
  - [ ] table: `name = { stage = "", prod = ""}`
    - usage: `{name.stage}`

## Serve
- [ ] Rename to only `serve`
- Options
  - [ ] Flag for current status
    - Shows req/s
  - [ ] Response status code
  - [ ] Response body
  - [ ] Response headers
  - [ ] Server certificate
  - [ ] Summary

## TUI
- [x] URL suggestions from alias file
- [x] Headers
- [x] Body
  - [x] Implement skip
  - [x] Implement editor
  - [x] Implement file
- Options
  - [ ] Timeout
  - [ ] Certificate
- [ ] Splash screen
  - [ ] New request
  - [ ] History

# Fixes
