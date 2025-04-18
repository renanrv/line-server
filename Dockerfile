# build the binary we want to run
FROM golang:1.24-bullseye AS builder
COPY . /src
WORKDIR /src

RUN go build -trimpath -o ./bin/server ./cmd/server
RUN go build -trimpath -o ./bin/file-generator ./internal-tools/file-generator

# copy over the binary we built to the image we're going to run
FROM debian:buster-slim
COPY --from=builder /src/bin/server /app/server
COPY --from=builder /src/bin/file-generator /app/file-generator
RUN chmod +x /app/server
RUN chmod +x /app/file-generator
WORKDIR /app
RUN apt-get update && apt-get install -y ca-certificates
CMD ["./server"]
