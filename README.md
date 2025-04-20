# Line Server

![CI](https://github.com/renanrv/line-server/actions/workflows/go-tests.yml/badge.svg)
[![Go Coverage](https://github.com/renanrv/line-server/wiki/coverage.svg)](https://raw.githack.com/wiki/renanrv/line-server/coverage.html)

Line Server is a network server that serves individual lines of an immutable text file over the network to clients using the following simple REST API:

 * `GET /lines/<line index>`
   * Returns an HTTP status of 200 and the text of the requested line
   * Returns HTTP 413 status if the requested line is beyond the end of the file

The API specification is defined in the OpenAPI 3.0 format and can be found in the [lineserver.openapi.yaml](docs/openapi/lineserver.openapi.yaml) file.

## How does the system work?

The system reads an immutable text file and serves individual lines to clients via a REST API. It uses an optional file index summary to optimize line retrieval for large files. If the index summary is unavailable, the system reads the file line by line. The server is implemented in Go and uses the `oapi-codegen` library to generate server-side code from the OpenAPI specification.

## How will the system perform with a 1 GB file? a 10 GB file? a 100 GB file?

- **1 GB file**: The system will perform efficiently, especially if a file index summary is available. Without the index, performance will degrade slightly as the file size increases due to line-by-line scanning.
- **10 GB file**: With a file index summary, the system will still perform well, as it can seek directly to the required line. Without the index, performance will degrade further due to the need to scan a larger file.
- **100 GB file**: The system will require sufficient memory and disk I/O capacity. With a file index summary, performance will remain acceptable. Without the index, performance will be significantly slower due to the need to scan the file sequentially.

## How will the system perform with 100 users? 10000 users? 1000000 users?

- **100 users**: The system can handle this load comfortably, assuming sufficient server resources.
- **10,000 users**: Performance will depend on the server's hardware and network capacity. Load balancing and horizontal scaling may be required.
- **1,000,000 users**: The system will require significant scaling, such as deploying multiple instances behind a load balancer and using caching mechanisms to reduce file access overhead.

## What documentation, websites, papers, etc., were consulted in doing this assignment?

- [Go Documentation](https://golang.org/doc/)
- [oapi-codegen GitHub Repository](https://github.com/oapi-codegen/oapi-codegen)
- [Docker Documentation](https://docs.docker.com/)
- Various articles and tutorials on file handling in Go

## What third-party libraries or other tools does the system use? What was the process to choose each library or framework used?

- **[oapi-codegen](https://github.com/oapi-codegen/oapi-codegen)**: Used to generate server-side code from the OpenAPI specification. Chosen for its compatibility with Go and ease of use.
- **[zerolog](https://github.com/rs/zerolog)**: Used for structured logging. Chosen for its performance and simplicity.
- **[errors](https://pkg.go.dev/github.com/pkg/errors)**: Used for enhanced error handling. Chosen for its ability to wrap errors with additional context.
- **[flag](https://github.com/namsral/flag)**: Used for command-line argument and environment variable parsing. Chosen for its flexibility and ease of use.
- **[testify](https://github.com/stretchr/testify)**: Used for unit testing. Chosen for its rich set of assertions and mocking capabilities.

## What was the estimated time spent on the exercise? What would be potential improvements and priorities if given unlimited additional time?

- **Time spent**: Approximately [X hours/days].
- **If more time were available**:
   1. Add automated end-to-end and performance tests.
   2. Add support for concurrent file access with locking mechanisms.
   3. Implement a watchdog mechanism to monitor file changes and update the index summary.
   4. Identify automatically the memory and disk I/O capacity of the server and adjust the maximum number of in-memory indexes accordingly.

## What are some critical observations or areas for improvement in the code?

- The code is functional and meets the requirements, but there are areas for improvement:
   - **Testing**: More test cases are needed to cover high-concurrency scenarios.
   - **Scalability**: Additional testing under extreme conditions is necessary to determine whether the current implementation can handle very high user loads effectively.

## Usage

### Line Server

#### Prerequisites
* Go 1.24 or later

#### Build the system
```bash
./build.sh
```
or
```bash
make build
```

#### Run the server
```bash
./run.sh
```
or 
```bash
make run
```

The server will listen on port 8080 by default. You can change the host and port by setting the `HTTP_ADDR` environment variable.

#### Call the REST API
```bash
curl -i -X GET http://localhost:8080/v0/lines/1
```

This will return the second line of the file (line index starts at 0). The server will return a 413 status if the requested line is beyond the end of the file.

#### Run the tests
```bash
make test
```

### File Generator Tool

The `file-generator` is an internal tool designed to generate large text files for testing purposes. It creates a file with a specified number of lines, where each line contains a unique string.

#### Generating a File

To generate a file, run the `file-generator` tool with the following command:

```bash
go run ./internal-tools/file-generator/main.go <number_of_lines> <output_file_path>
```

##### Command-Line Arguments
* `number_of_lines`: Specifies the number of lines to generate.
* `output_file_path`: Specifies the path to the output file.

##### Example
To generate a file with 1,000 lines at `./data/sample.txt`, run:

```bash
go run ./internal-tools/file-generator/main.go 1000 ./data/sample.txt
```

This will create a file named `sample.txt` in the `data` directory, containing 1,000 lines of unique strings.

##### Notes

* Ensure the `output` directory exists before running the tool, or the tool will fail to write the file.
* The tool is useful for testing the `Line Server` with large files.

## Benchmarks

To evaluate the performance of the `Line Server`, benchmarks were conducted using the [`oha`](https://github.com/hatoo/oha) tool. This tool is used to simulate load and concurrency scenarios, providing insights into the server's behavior under different conditions.

### Prerequisites

* Install the `oha` tool by following the instructions in its [GitHub repository](https://github.com/hatoo/oha).

### Running Benchmarks

You can use the `oha` tool to test the server's performance by sending multiple requests to the `Line Server` API. Below are examples of how to run benchmarks:

#### Example 1: Testing with 100 concurrent users and 1,000 total requests
```bash
oha -n 1000 -c 100 --disable-keepalive --latency-correction -m GET --rand-regex-url http://localhost:8080/v0/lines/[1-9]{1,8}
```

* `-n 1000`: Total number of requests to send.
* `-c 100`: Number of concurrent users.
* `--disable-keepalive`: Disables HTTP keep-alive connections.
* `--latency-correction`: Correct latency to avoid coordinated omission problem.
* `--rand-regex-url`: Randomly generates URLs based on the provided regex pattern.

#### Example 2: Testing with 10,000 total requests and 500 concurrent users

```bash
hey -n 10000 -c 500 --disable-keepalive --latency-correction -m GET --rand-regex-url http://localhost:8080/v0/lines/[1-9]{1,8}
```

### Results
Below is a table with the results of the benchmarks for different scenarios:

#### No in-memory index
| Scenario                   | Total Requests | Concurrent Users | Requests per Second | Average Latency (ms) | Error Rate (%) |
|----------------------------|----------------|------------------|----------------------|-----------------------|----------------|
| Small Load (1,000 reqs)    | 1,000          | 100              | [Insert Value]       | [Insert Value]        | [Insert Value] |
| Medium Load (100,000 reqs) | 100,000        | 10,000           | [Insert Value]       | [Insert Value]        | [Insert Value] |
| High Load (10,000,000 reqs)   | 10,000,000     | 1,000,000        | [Insert Value]       | [Insert Value]        | [Insert Value] |

#### With in-memory index
| Scenario                   | Total Requests | Concurrent Users | Requests per Second | Average Latency (ms) | Error Rate (%) |
|----------------------------|----------------|------------------|----------------------|-----------------------|----------------|
| Small Load (1,000 reqs)    | 1,000          | 100              | [Insert Value]       | [Insert Value]        | [Insert Value] |
| Medium Load (100,000 reqs) | 100,000        | 10,000           | [Insert Value]       | [Insert Value]        | [Insert Value] |
| High Load (10,000,000 reqs)   | 10,000,000     | 1,000,000        | [Insert Value]       | [Insert Value]        | [Insert Value] |

#### Notes
* Ensure the server is running before executing the benchmarks.
* Use different line indexes in the API endpoint to test various scenarios.
* For large-scale tests, consider running the benchmarks on a machine with sufficient resources to avoid client-side bottlenecks.
