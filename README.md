# http

A CLI program for sending basic HTTP requests. Features:
 - Simplified URL parsing, e.g. `:1234/path` becomes `http://localhost:1234/path`
 - Persistant URL aliases
 - Convenient request body handling through stdin, file or flag
 - AWS signature v4 support

## Installation

Visit [releases](https://github.com/lunjon/http/releases/latest) and download
executable for your platform (if available).

## Usage

### Sending requests

To get started use `http --help`.

**Examples**:

```sh
# POST http://localhost:1234/api/test 
$ http post :1234/api/test --json '{"field":"value"}'
POST      http://localhost:1234/api/test
Status   201 OK
Elapsed  15.11 ms
{
    ...
}

# GET https://api.example/resources/abbccc-122333, using header X-User with value donald
$ http get https://api.example/resources/abbccc-122333 --header x-user=donald
GET      https://api.example/resources/abbccc-122333
Status   403 Forbidden
Elapsed  102.97 ms
```

### Default headers

Default headers can be set by using an environment variable: `DEFAULT_HEADERS`.
The string should contain headers in the same format specified using the
`--header` flag, and multiple headers should be separated by a `|`.

### URL alias

URL alias can be created with the `alias` sub-command like so:
 - List: `http alias`
 - Add:  `http alias <name> <url>`

An alias can then be used in the request URL:
```sh
$ http get "{name}/api/path"
```

### Request body

Can be specified by three methods:
- As string: `http post http://example.com/api --body '{"name":"meow"}'`
- As file: `http post http://example.com/api --body r.json`
- Pipe from stdin: `http post http://example.com/api < r.json`
