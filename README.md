# http

A CLI program for sending basic HTTP requests. Features:
 - Simplified URL parsing, e.g. `:1234/path` becomes `http://localhost:1234/path`
 - Persistant URL aliases
 - AWS signature v4 support
 - Reading request body either a string or file, it can detect the difference

## Installation

Having Go 1.14+ installed run:

```sh
# Without cloning repository...
$ go get github.com/lunjon/http

# or in project root
$ go install
```

## Usage

### Sending requests

To get started use `http --help`.

**Examples**:

```sh
# POST http://localhost:1234/api/test 
$ http post :1234/api/test --json '{"field":"value"}'
GET      http://localhost:1234/api/test
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
