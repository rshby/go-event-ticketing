BINARY := bin/go-event-ticketing

.PHONY: all fmt lint test build clean

all: fmt lint test build

fmt:
	gofmt -s -w .
	goimports -w .

lint:
	go vet ./...
	golangci-lint run ./...

test:
	go test -race -cover ./...

build:
	go build -o $(BINARY) .

clean:
	rm -rf bin
