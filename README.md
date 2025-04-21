# Line Server

![CI](https://github.com/renanrv/line-server/actions/workflows/go-tests.yml/badge.svg)
[![Go Coverage](https://github.com/renanrv/line-server/wiki/coverage.svg)](https://raw.githack.com/wiki/renanrv/line-server/coverage.html)

Line Server is a network server that serves individual lines of an immutable text file over the network to clients using the following simple REST API:

 * `GET /lines/<line index>`
   * Returns an HTTP status of 200 and the text of the requested line
   * Returns HTTP 413 status if the requested line is beyond the end of the file

The API specification is defined in the OpenAPI 3.0 format and can be found in the [lineserver.openapi.yaml](docs/openapi/lineserver.openapi.yaml) file.

## How does the system work?

The system reads an immutable text file and serves individual lines to clients via a REST API. 
It uses an optional file index summary to optimize line retrieval for large files. 
If the index summary is unavailable, the system reads the file line by line.


The server is designed to handle large files efficiently by using an in-memory index that maps line numbers to file offsets. 
This allows for quick access to specific lines without the need to read the entire file sequentially.
With its default configuration, the server will check the available memory of the server in order to use 70% of it for the in-memory index.
With the total memory value and the estimated memory usage per index entry (8 bytes for *int* key + 8 bytes for *int64* value, map overhead excluded),
it is possible to calculate the maximum number of in-memory indexes that can be created.
This logic can be found in the [index_to_memory.go](pkg/fileprocessing/index_to_memory.go) file.


The server is implemented in Go and uses the `oapi-codegen` library to generate server-side code from the OpenAPI specification.

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

- **Time spent**: Approximately 16 hours.
- **If more time were available**:
   1. Add automated end-to-end and performance tests.
   2. Assess support for concurrent file access with locking mechanisms.
   3. Assess dumping the in-memory index to disk for persistence and queries in case of scenarios with low memory available.
   4. Implement a watchdog mechanism to monitor file changes and update the index summary.

## What are some critical observations or areas for improvement in the code?

- The code is functional and meets the requirements, but there are areas for improvement:
   - **Testing**: More test cases are needed to cover high-concurrency scenarios.
   - **Scalability**: Additional testing under extreme conditions is necessary to determine whether the current implementation can handle very high user loads effectively.

## Usage

### Line Server

#### Build and run the system

##### Environment variables

The following environment variables are supported by the Line Server:

| Variable Name         | Default Value           | Description                                                                 |
|-----------------------|-------------------------|-----------------------------------------------------------------------------|
| `HTTP_ADDR`           | `:8080`                | The address that will expose the server API.                               |
| `DEBUG_ADDR`          | `:8081`                | The address for debug and metrics.                                         |
| `FILE_PATH`           | `./data/sample_100.txt`| The path to the file that will be used to read the lines.                  |
| `MAX_INDEXES`         | `0`                    | The maximum number of indexes to generate. `0` uses all available memory. Negative values disable in-memory index generation. |
| `CORS_ALLOWED_ORIGINS`| `http://localhost:8080`| Comma-separated list of allowed origins for CORS.                          |
| `LOG_LEVEL`           | `1`                    | The log level for the server. `0` for debug, `1` for info, `2` for warning, `3` for error. |

###### Notes
* You can override these variables by setting them in your environment or passing them as flags when running the server.
* For Docker or Docker Compose, these variables can be set using the -e flag or in the docker-compose.yml file.

##### Local setup

1. Prerequisites:

* Go 1.24 or later

2. Install Go dependencies:
```bash
go mod tidy
```

3. Build the project:
```bash
./build.sh
```
or
```bash
make build
```

4. Run the server
```bash
./run.sh
```
or
```bash
make run
```

The server will listen on port 8080 by default. You can change the host and port by setting the `HTTP_ADDR` environment variable.

##### Using Docker
1. Build the Docker image:

```bash
docker build -t line-server .
```

2. Run the Docker container:

```bash
docker run -p 8080:8080 -v ./data:/app/data -it line-server
```

This command maps the `data` directory from your host machine to the container, allowing the server to access files in that directory.

The server will be accessible at `http://localhost:8080`.

###### Notes
* Ensure Docker are installed on your system.
* You can modify the environment variables using `-e` flags with the docker run command with the format `-e FILE_PATH=./data/sample_1000.txt`

##### Using Docker Compose

1. Build and run the Docker container using Docker Compose:

```bash
docker-compose up --build
```

The server will be accessible at `http://localhost:8080`.

###### Notes
* Ensure Docker and Docker Compose are installed on your system.
* You can modify the environment variables in the docker-compose.yml file.

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

The `file-generator` is an internal tool designed to generate large text files for testing purposes. It creates a file with a specified size in GB, where each line contains a unique string.

#### Generating a File

To generate a file, run the `file-generator` tool with the following command:

##### Local setup
```bash
go run ./internal-tools/file-generator/main.go <file_size_gb> <output_file_path>
```

###### Command-Line Arguments
* `file_size_gb`: Specifies the size in GB for the generated file.
* `output_file_path`: Specifies the path to the output file.

###### Example
To generate a file with 1 GB at `./data/sample.txt`, run:

```bash
go run ./internal-tools/file-generator/main.go 1 ./data/sample.txt
```

This will create a file named `sample.txt` in the `data` directory, containing all the lines of unique strings that fit within 1 GB.

##### Using Docker

```bash
docker run -v ./data:/app/data -it line-server ./file-generator <file_size_gb> <output_file_path>
```

###### Example
To generate a file with 1 GB at `./data/sample.txt`, run:

```bash
docker run -v ./data:/app/data -it line-server ./file-generator 1 ./data/sample.txt
```

This will create a file named `sample.txt` in the `data` directory, containing all the lines of unique strings that fit within 1 GB.


##### Using Docker Compose

```bash
docker-compose run --rm -it line-server ./file-generator <file_size_gb> <output_file_path>
```

###### Example
To generate a file with 1 GB at `./data/sample.txt`, run:

```bash
docker-compose run --rm -it line-server ./file-generator 1 ./data/sample.txt
```

This will create a file named `sample.txt` in the `data` directory, containing all the lines of unique strings that fit within 1 GB.

##### Notes

* Ensure the `output` directory exists before running the tool, or the tool will fail to write the file.
* The tool is useful for testing the `Line Server` with large files.

## Benchmarks

To evaluate the performance of the `Line Server`, benchmarks were conducted using the [`k6`](https://k6.io/) tool. This tool is used to simulate load and concurrency scenarios, providing insights into the server's behavior under different conditions.

### Prerequisites

* Install the `k6` tool by following the instructions in its [official documentation](https://grafana.com/docs/k6/latest/).

### Running Benchmarks

You can use the `k6` tool to test the server's performance by sending multiple requests to the `Line Server` API. Below are examples of how to run benchmarks:

#### Example 1: Testing with 100 concurrent users and 1,000 total requests

1. Create a JavaScript file named `script.js` with the following content:

```javascript
import http from 'k6/http';
import { check } from 'k6';

export const options = {
  vus: 100,
  iterations: 1000,
};

export default function () {
  const randomId = Math.floor(Math.random() * 77489496); // 0 to 77489495 inclusive
  const url = `http://localhost:8080/v0/lines/${randomId}`;
  const res = http.get(url, { timeout: '600s'});

  check(res, {
    'status is 200': (r) => r.status === 200,
  });
}
```

2. Run the benchmark using the `k6` command:
```bash
k6 run script.js
```

* `vus 100`: Number of virtual users (concurrent users).
* `iterations 1000`: Total number of requests to send.

### Results

The benchmarks were conducted on the following machine:

- **Operating System**: macOS
- **Processor**: Apple M1 (8 cores)
- **Memory**: 16 GB

Below is a table with the results of the benchmarks for different scenarios:

#### No in-memory index (1 GB file)
| Scenario                 | Total Requests | Concurrent Users | Requests per Second | Average Latency (s) | p95 Latency (s) | Error Rate (%) |
|--------------------------|-------------|----------------|---------------------|---------------------|-----------------|----------------|
| Small Load (100 reqs)    | 100         | 10             | 3.96                | 2.44                | 7.67            | 0.00           |
| Small Load (1,000 reqs)  | 1,000       | 100            | 5.38                | 18.12               | 33.58           | 0.00           |
| Medium Load (10,000 reqs) | 10,000      | 1,000          | 5.27                | 157                 | 324             | 1.74           |

#### With in-memory index (1 GB file) - all lines indexed
| Scenario                 | Total Requests | Concurrent Users | Requests per Second | Average Latency (s) | p95 Latency (s) | Error Rate (%) |
|--------------------------|--------------|------------------|---------------------|---------------------|-----------------|----------------|
| Small Load (100 reqs)    | 100         | 10               | 3507.79             | 0.00259             | 0.00428         | 0.00           |
| Small Load (1,000 reqs)  | 1,000        | 100              | 15648.72            | 0.00564             | 0.01222         | 0.00           |
| Medium Load (10,000 reqs) | 10,000       | 1,000            | 5131.21           | 0.17901             | 0.30976         | 0.85           |

#### With in-memory index (4 GB file) - half of the lines indexed
| Scenario                 | Total Requests | Concurrent Users | Requests per Second | Average Latency (s) | p95 Latency (s) | Error Rate (%) |
|--------------------------|--------------|------------------|---------------------|---------------------|-----------------|----------------|
| Small Load (100 reqs)    | 100         | 10               | 780.84              | 0.01113             | 0.03539         | 0.00           |
| Small Load (1,000 reqs)  | 1,000        | 100              | 768.04              | 0.12446             | 0.29438         | 0.00           |
| Medium Load (10,000 reqs) | 10,000       | 1,000            | 891.41              | 1.07                | 2.81            | 1.82           |


#### Notes
* Ensure the server is running before executing the benchmarks.
* Use different line indexes in the API endpoint to test various scenarios.
* For large-scale tests, consider running the benchmarks on a machine with sufficient resources to avoid client-side bottlenecks.
