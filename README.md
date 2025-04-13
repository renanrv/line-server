# Line Server
Line Server is a network server that serves individual lines of an immutable text file over the network to clients using the following simple REST API:

 * `GET /lines/<line index>`
   * Returns an HTTP status of 200 and the text of the requested line
   * Returns HTTP 413 status if the requested line is beyond the end of the file

The API specification is defined in the OpenAPI 3.0 format and can be found in the [lineserver.openapi.yaml](docs/openapi/lineserver.openapi.yaml) file.

## Usage

### Prerequisites
* Go 1.24 or later

### Build the system
```bash
./build.sh
```
or
```bash
make build
```

### Run the server
```bash
./run.sh
```
or 
```bash
make run
```

The server will listen on port 8080 by default. You can change the host and port by setting the `HTTP_ADDR` environment variable.

### Call the REST API
```bash
curl -i -X GET http://localhost:8080/v0/lines/1
```

This will return the second line of the file (line index starts at 0). The server will return a 413 status if the requested line is beyond the end of the file.

### Run the tests
```bash
make test
```

## Third-party libraries or other tools

The OpenAPI file was used to generate code for the server-side implementation and the HTTP models, using the [oapi-codegen](https://github.com/oapi-codegen/oapi-codegen) command-line tool and library.
