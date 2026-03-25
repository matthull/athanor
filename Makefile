.PHONY: build install clean test lint fmt vet check

BINARY := ath
BUILD_DIR := .
INSTALL_DIR := $(HOME)/.local/bin

VERSION := $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
COMMIT := $(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")
BUILD_TIME := $(shell date -u +"%Y-%m-%dT%H:%M:%SZ")

LDFLAGS := -X github.com/matthull/athanor/internal/cli.Version=$(VERSION) \
           -X github.com/matthull/athanor/internal/cli.Commit=$(COMMIT) \
           -X github.com/matthull/athanor/internal/cli.BuildTime=$(BUILD_TIME)

build:
	go build -ldflags "$(LDFLAGS)" -o $(BUILD_DIR)/$(BINARY) ./cmd/ath

install: build
	install -m 755 $(BUILD_DIR)/$(BINARY) $(INSTALL_DIR)/$(BINARY)
	@echo "Installed $(BINARY) to $(INSTALL_DIR)/$(BINARY)"

clean:
	rm -f $(BUILD_DIR)/$(BINARY)
	go clean -testcache

test:
	go test ./... -v

test-short:
	go test ./... -v -short

lint:
	golangci-lint run ./...

fmt:
	gofmt -w .
	goimports -w .

vet:
	go vet ./...

# Run all checks (what CI would run)
check: fmt vet lint test
