# General
- Config sub-commands
  - [x] Init
  - [x] List
  - [x] Edit
- Request history
  - [ ] Write to file
  - [ ] Truncate to x entries

# Configuration
- [x] Format: toml?
- [x] Timeout
- Sections
  - [x] aliases
  - [ ] history (for config, not actual history)
- Complex alias: support having the alias value:
  - [x] string: use value as is
  - [ ] table: `name = { stage = "", prod = ""}`
    - usage: `{name.stage}`

# TUI
- [x] URL suggestions from alias file
- [x] Headers
- [x] Body
  - [x] Implement skip
  - [x] Implement editor
  - [x] Implement file
- [ ] Resend request
- [ ] Options
  - [ ] Timeout
  - [ ] Certificate
- [ ] Splash screen
  - [ ] New request
  - [ ] History

