.PHONY: build run test clean deps lint docker

BINARY_NAME=brahma
BUILD_DIR=bin
GO=go

VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
BUILD_TIME=$(shell date -u '+%Y-%m-%d_%H:%M:%S')
LDFLAGS=-ldflags "-X main.Version=$(VERSION) -X main.BuildTime=$(BUILD_TIME)"

build:
	@mkdir -p $(BUILD_DIR)
	$(GO) build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME) ./cmd/brahma

run: build
	./$(BUILD_DIR)/$(BINARY_NAME) -config config.json

test:
	$(GO) test -v -race -cover ./...

clean:
	rm -rf $(BUILD_DIR)
	$(GO) clean

deps:
	$(GO) mod download
	$(GO) mod tidy

lint:
	golangci-lint run ./...

docker:
	docker build -t $(BINARY_NAME):$(VERSION) .

docker-run:
	docker run -p 8080:8080 -v $(PWD)/config.json:/app/config.json $(BINARY_NAME):$(VERSION)
