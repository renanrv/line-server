build:
	./build.sh

run:
	./run.sh

$(GOBIN)/golangci-lint:
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(GOBIN) v2.1.1

.PHONY: tools
tools: $(GOBIN)/golangci-lint

lint: tools
	$(GOBIN)/golangci-lint run ./... --output.text.colors --timeout=5m
	gofmt -w .
	goimports -w .

gen:
	oapi-codegen --config docs/openapi/oapi-codegen.config.yaml docs/openapi/lineserver.openapi.yaml

test: unit-test cover
	go mod tidy

unit-test:
	mkdir -p .coverage
	rm -f .coverage/cover_unit.out
	go test -timeout=10s -race -benchmem -tags=unit -coverpkg=./... -coverprofile=".coverage/cover_unit.out" ./...

cover:
	rm -f .coverage/cover.out
	rm -f .coverage/cover_unit.main_filtered.out
	rm -f .coverage/cover_unit.filtered.out
    # Filter out the main.go and generated files from the unit test coverage report
	grep -v "main.go" .coverage/cover_unit.out > .coverage/cover_unit.main_filtered.out
	grep -v "server/server.gen.go" .coverage/cover_unit.main_filtered.out > .coverage/cover_unit.filtered.out
	cat .coverage/cover_unit.filtered.out >> .coverage/cover.out
	go tool cover -func=.coverage/cover.out

cover-report: cover
	go tool cover -html=.coverage/cover.out