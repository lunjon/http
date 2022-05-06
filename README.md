# http

A CLI for sending basic HTTP requests. Features:
 - Simplified URL parsing, e.g. `:1234/path` becomes `http://localhost:1234/path`
 - Persistant URL aliases
 - Convenient request body handling through stdin, file or flag
 - AWS signature v4 support

## Installation

Visit [releases](https://github.com/lunjon/http/releases/latest) and download
executable for your platform (if available) or use `go get github.com/lunjon/http`.

## Usage

To get started use `http --help`.

### Sending requests

**Examples**:

```sh
# POST http://localhost:1234/api/test 
$ http post :1234/api/test --body '{"field":"value"}'
...

# GET https://api.example/resources/abbccc-122333, using header X-User with value donald
$ http get api.example/resources/abbccc-122333 --header x-user=donald
...
```

### Default headers

Default headers can be set by using an environment variable: `DEFAULT_HEADERS`.
The string should contain headers in the same format specified using the
`--header` flag, and multiple headers should be separated by a `|`.

### URL alias

The `alias` is used to list, create and remove URL aliases:
 - List: `http alias`
 - Add: `http alias name url`
 - Remove: `http alias --remove name`

An alias can then be used in the request URL:
```sh
$ http get "{name}/api/path"
```

### Request body

Can be specified as:
- string: `http post http://example.com/api --body '{"name":"meow"}'`
- file: `http post http://example.com/api --body r.json`
- stdin: `cat file | http post http://example.com/api`

## Shell completion (WIP)

If using bash, zsh or any other bash-like shell you can use the [shell completion
script](./complete.sh):

```sh
$ source complete.sh
```
