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

### get, post and delete

**Format**:
```sh
$ httpreq <method> <route> [options]
```

**Description**: Used to perform 

**Flags**:

- `--header` (string): Specify a key/value pairs (`name=value` or `name:value`) to use as an HTTP header. They can be either a comma separated list of key/value pairs or specified using multiple times.
    * For instance: `--header h1=value1,h2=value2` and `--header h1:value1 --header h2=value2` will yield the same result.
- `-4/--aws-sigv4` (bool: Sign the request with AWS signature V4.
    * If the `--aws-profile` flag is given it tries to use the credentials for that profile, else it looks for the environment variables.
- `--aws-region` (string): The AWS region to use when signing the request. 
    * Default is `eu-west-1`
- `--aws-profile` (string): Use the AWS profile when signing the request.
    * Note that the profile must have credentials defined in the profile for it to work.
- `--output-file` (string): If there was any response body, output the content to the given file.
    * If not set, it outputs the content to stdout.
- `--sandbox` (bool): Run the request to a local server that only echo request information.

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

### run

**Format**:
```sh
$ httpreq run <file> [flags]
```

**Description**: `httpreq` provide a command called `run` for running requests from a file. These files, lets call them *spec* files, are written as JSON or YAML files in a special format. The total specification for such files can be found in `docs/spec.json` and `docs/spec.yaml` respectively.

**Flags**:
- `--sandbox` (bool): Run the request to a local server that only echo request information.

**Examples**:
```yaml
requests:
    - 
        name: example request
        method: get
        url: https://api.example.com/path
```

### sandbox

**Format**:
```sh
$ httpreq sandbox 
```

**Description**: Start a local server at port 8118 (can be changed using `--port`). It will block the program.

**Flags**:
- `-p/--port` (int): Start the server on this port instead of default (8118).

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

## Important Notes

**URLs**: URLs used in httpreq support the different formats below:
- `/path` ==> `http://localhost/path`
- `:port/path`: ==> `http://localhost:port/path`
- `host.com[:port]/path` ==> `https://host.com[:port]/path`
- `http[s]://host.com[:port]/path` ==> `http[s]://host.com[:port]/path`

## TODO

- **Variable support in spec files**: It would be nice to define global variables, e.g. an API url, and use them in the requests
