# HTTPReq

A CLI program for sending basic HTTP requests. Features:
 - Simplified URL parsing, e.g. `:1234/path` becomes `http://localhost:1234/path`
 - AWS signature v4 support
 - Reading request body either a string or file, it can detect the difference

## Installation

Having Go 1.13+ installed run:

```sh
# Without cloning repository
$ go get github.com/lunjon/httpreq

# Install from project root
$ go install
```

## Usage

### Sending requests

**Format**:
```sh
$ httpreq <method> <url> [options]
```

**Flags**:

- `-H, --header` (string): Specify a key/value pairs (`name=value` or `name:value`) to use as an HTTP header.
  They can be either a comma separated list of key/value pairs or specified using multiple times.
    * For instance: `--header h1=value1,h2=value2` and `--header h1:value1 --header h2=value2` will yield the same result.
- `-4/--aws-sigv4` (bool: Sign the request with AWS signature V4.
    * If the `--aws-profile` flag is given it tries to use the credentials for that profile, else it looks for the environment variables.
- `--aws-region` (string): The AWS region to use when signing the request. 
    * Default is `eu-west-1`
    * Note that the profile must have credentials defined in the profile for it to work.
- `-T, --timeout` (duration): Specify request timeout.
    * Default is 10 seconds.
- `-s, --silent` (string): Suppress response output.
- `-v, --verbose` (string): Output debug logs.

**Examples**:

```sh
# POST http://localhost:1234/api/test 
$ httpreq post :1234/api/test --json '{"field":"value"}'
GET      http://localhost:1234/api/test
Status   201 OK
Elapsed  15.11 ms
{
    ...
}

# GET https://api.example/resources/abbccc-122333, using header X-User with value donald
$ httpreq get https://api.example/resources/abbccc-122333 --header x-user=donald
GET      https://api.example/resources/abbccc-122333
Status   403 Forbidden
Elapsed  102.97 ms
```

### Default headers

Default headers can be set by using an environment variable: `DEFAULT_HEADERS`.
The string should contain headers in the same format specified using the
`--header` flag, and multiple headers should be separated by a `|`.
