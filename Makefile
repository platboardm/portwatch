# Makefile for portwatch
BINARY   := portwatch
CMD      := ./cmd/portwatch
VERSION  ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo dev)
LDFLAGS  := -ldflags "-X main.version=$(VERSION)"

.PHONY: all build test lint clean run

all: build

build:
	go build $(LDFLAGS) -o $(BINARY) $(CMD)

test:
	go test ./...

test-short:
	go test -short ./...

test-race:
	go test -race ./...

lint:
	golangci-lint run ./...

clean:
	rm -f $(BINARY)

run: build
	./$(BINARY) --config internal/config/example_config.yaml
