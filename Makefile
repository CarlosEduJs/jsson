BINARY_NAME=jsson
VERSION=$(shell git describe --tags --always --dirty 2>/dev/null | sed 's/^v//' || echo "dev")
COMMIT=$(shell git rev-parse --short HEAD 2>/dev/null || echo "none")
DATE=$(shell date -u +"%Y-%m-%dT%H:%M:%SZ" 2>/dev/null || echo "unknown")
LDFLAGS=-ldflags "-s -w -X jsson/internal/version.Version=$(VERSION) -X jsson/internal/version.Commit=$(COMMIT) -X jsson/internal/version.Date=$(DATE)"

.PHONY: build test test-e2e lint clean fmt run install help

build:
	go build $(LDFLAGS) -o bin/$(BINARY_NAME) ./cmd/jsson

test:
	go test -v -race -count=1 ./...

test-e2e:
	go test -tags=e2e -count=1 ./cmd/jsson/

lint:
	go vet ./...
	@if [ -n "$$(gofmt -l .)" ]; then \
		echo "Unformatted files:"; \
		gofmt -l .; \
		exit 1; \
	fi

clean:
	rm -rf bin/ dist/

fmt:
	gofmt -l -w .

run: build
	./bin/$(BINARY_NAME) $(ARGS)

install: build
	cp bin/$(BINARY_NAME) $(GOPATH)/bin/$(BINARY_NAME)

help:
	@echo "Targets:"
	@echo "  build      Build binary to bin/jsson with version ldflags"
	@echo "  test       Run unit tests"
	@echo "  test-e2e   Run e2e tests (golden files)"
	@echo "  lint       Run go vet + gofmt check"
	@echo "  clean      Remove bin/ and dist/"
	@echo "  fmt        Format all Go source files"
	@echo "  run        Build and run (pass ARGS=...)"
	@echo "  install    Copy binary to GOPATH/bin"
