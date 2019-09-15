# Specification

This specification defines the structure of *request* files. [spec.yaml](./spec.yaml) and [spec.json](./spec.json) are examples of two spec files.

## Formats

The request files can be written as JSON or YAML files.

## Fields

### requests

The only required field in request files. It is an array of one or more [request items](#request-item).

#### request-item

Each request item contains the following fields (`*` are required):

- **id\*** (string): A global unique ID of the request. It should be a string that contains no whitespace.
- **url\*** (string): A valid URL conforming to the rules defined in [README](../README.md). Supports - environment variables.
- **headers** (map[string, string]): Unique headers for this particular request. If a global header is - defined, the value in here will override the global header value. Supports environment variables.
- **body** (dynamic): Required in POST requests. It can be of any structure.
- **aws** (dynamic): Used when sending requests that require AWS signiture V4. It can have either form below:
    - Only the value `true`. This will use default values; AWS region will be `eu-west-1` and environment credentials.
    - Two fields: `profile` and `region`. If profile is set, it will try to use that AWS profile. If only region is specified it will use that region and environment credentials.

### env

`env` is a global field that defines environment variables that will be expanded in the headers (both global and per request) and URLs of the requests. `env` should be of type `map[string, string]`.

### headers

`headers` is a global field that defines headers that will be used in all requests. If a header with the same name is defined in a request, the value in the request will be used (i.e. override the global header). `headers` should be of type `map[string, string]`.