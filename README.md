# HTTPReq

This is a simple command to execute HTTP request using you command line. Why not use curl? curl is awesome, but this is only created for convenience and to simplify some requests.

## Installation

You need Go 1.12 for this to work.

```sh
$ go get ./...
$ cd cmd/httpreq
$ go install
```

## Usage

`httpreq` is really easy to use. The program use the the following format for commands:

```sh
$ httpreq <method> <route> [options]
```

The are common and specific flags for the different methods. The common are:

- **--header**: Specify a key/value pairs (`name=value` or `name:value`) to use as an HTTP header. They can be either a comma separated list of key/value pairs or specified using multiple times.
    * For instance: `--header h1=value1,h2=value2` and `--header h1:value1 --header h2=value2` will yield the same result.

### Examples

Below are some examples with a comment above each command that shows the corresponding request.

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