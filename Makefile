BINARY_NAME=gongfeng
VERSION=$(shell git describe --tags --always --dirty)
LDFLAGS=-ldflags "-X github.com/studyzy/gongfeng-cli/internal/cmd.Version=$(VERSION)"

.PHONY: build test lint fmt coverage install clean

build:
	go build $(LDFLAGS) -o $(BINARY_NAME) ./cmd/gongfeng/

test:
	go test ./...

lint:
	go vet ./...
	goimports -l .

fmt:
	gofmt -w .
	goimports -w .

coverage:
	go test ./... -coverprofile=coverage.out
	go tool cover -func=coverage.out

install:
	go install $(LDFLAGS) ./cmd/gongfeng/

clean:
	rm -f $(BINARY_NAME) coverage.out
