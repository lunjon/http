# HTTPReq

This is a simple command to execute HTTP request using you command line. Why not use curl? curl is awesome, but `httpreq` was only created for convenience and to simplify some requests.

[![Build Status](https://travis-ci.org/lunjon/httpreq.svg?branch=master)](https://travis-ci.org/lunjon/httpreq)

## Installation

Having Go 1.12+ installed run:

```sh
$ go get ./...
$ cd cmd/httpreq
$ go install
```

## Usage

### Sending requests

**Format**:
```sh
$ httpreq <method> <route> [options]
```

**Description**: Used to perform basic HTTP requests.
Supported verbs: GET, HEAD, POST, PUT, PATCH, DELETE

**Flags**:

- `-H, --header` (string): Specify a key/value pairs (`name=value` or `name:value`) to use as an HTTP header. They can be either a comma separated list of key/value pairs or specified using multiple times.
    * For instance: `--header h1=value1,h2=value2` and `--header h1:value1 --header h2=value2` will yield the same result.
- `-4/--aws-sigv4` (bool: Sign the request with AWS signature V4.
    * If the `--aws-profile` flag is given it tries to use the credentials for that profile, else it looks for the environment variables.
- `--aws-region` (string): The AWS region to use when signing the request. 
    * Default is `eu-west-1`
- `--aws-profile` (string): Use the AWS profile when signing the request.
    * Note that the profile must have credentials defined in the profile for it to work.
- `-T, --timeout` (duration): Specify request timeout.
    * Default is 10 seconds.
- `-s, --silent` (string): Suppress response output.
- `-v, --verbose` (string): Output debug logs.

**Examples**:

```sh
# GET http://localhost/api/test
$ httpreq get /api/test
GET      http://localhost/api/test
Status   200 OK
Elapsed  31.72 ms
{
    ...
}
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

### parse-url

**Format**:
```sh
$ httpreq <url>  [flags]
```

**Description**: Tries to parse the URL and outputs the full URL.

**Flags**:
- `-d/--details` (bool): Output detailed information.

**Examples**:
```sh
$ httpreq parse-url host.com/api/id
https://host.com/api/id
```

### Default headers

Default headers can be set by using an environment variable: `DEFAULT_HEADERS`.
The string should contain headers in the same format specified using the
`--header` flag, and multiple headers should be separated by a `|`.

## Important Notes

**URLs**: URLs used in httpreq support the different formats below:
- `/path` ==> `http://localhost/path`
- `:port/path`: ==> `http://localhost:port/path`
- `host.com[:port]/path` ==> `https://host.com[:port]/path`
- `http[s]://host.com[:port]/path` ==> `http[s]://host.com[:port]/path`
