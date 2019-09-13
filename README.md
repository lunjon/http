# HTTPReq

This is a simple command to execute HTTP request using you command line. Why not use curl? curl is awesome, but `httpreq` was only created for convenience and to simplify some requests.

## Installation

Having Go 1.12+ installed run:

```sh
$ go get ./...
$ cd cmd/httpreq
$ go install
```

## Usage

### get, post and delete

`httpreq` is really easy to use. The program use the the following format for commands:

```sh
$ httpreq <method> <route> [options]
```

For `httpreq <method> ...`commands there are some common flags:

- **--header**: Specify a key/value pairs (`name=value` or `name:value`) to use as an HTTP header. They can be either a comma separated list of key/value pairs or specified using multiple times.
    * For instance: `--header h1=value1,h2=value2` and `--header h1:value1 --header h2=value2` will yield the same result.
- **-4/--aws-sigv4**: Sign the request with AWS signature V4.
    * If the `--aws-profile` flag is given it tries to use the credentials for that profile, else it looks for the environment variables.
- **--aws-region**: The AWS region to use when signing the request. 
    * Default is `eu-west-1`
- **--aws-profile**: Use the AWS profile when signing the request.
    * Note that the profile must have credentials defined in the profile for it to work.
- **--output-file**: If there was any response body, output the content to the given file.
    * If not set, it outputs the content to stdout.
- **--sandbox**: Run the request to a local server and echo request information.

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

### run

`httpreq` provide a command called `run` for running requests from a file. These files, lets call them *spec* files, are written as JSON or YAML files in a special format. The total specification for such files can be found in `docs/spec.json` and `docs/spec.yaml` respectively.

An example spec file can be:

```yaml
requests:
    - 
        name: example request
        method: get
        url: https://api.example.com/path
```

### sandbox

Start a local server at port 8118 (can be changed using `--port`). It will block the program.

## TODO

- **Variable support in spec files**: It would be nice to define global variables, e.g. an API url, and use them in the requests
- **Request reference**: Support referencing requests result.
