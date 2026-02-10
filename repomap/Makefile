# Makefile for repomap
BINARY_NAME=repomap
VERSION := $(shell git describe --tags --always --dirty)
BUILD_DIR=build
LDFLAGS=-ldflags "-X main.version=$(VERSION)"

.PHONY: all build clean test lint release

all: build

build:
	@echo "Building $(BINARY_NAME)..."
	go build $(LDFLAGS) -o $(BINARY_NAME) ./cmd/repomap

clean:
	@echo "Cleaning up..."
	rm -f $(BINARY_NAME)
	rm -rf $(BUILD_DIR)
	rm -f repomap.json repomap.xml

test:
	@echo "Running tests..."
	go test ./...

lint:
	@echo "Linting..."
	go vet ./...

release: clean test
	@echo "Building release binaries..."
	mkdir -p $(BUILD_DIR)
	# Linux
	GOOS=linux GOARCH=amd64 go build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-linux-amd64 ./cmd/repomap
	GOOS=linux GOARCH=arm64 go build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-linux-arm64 ./cmd/repomap
	# macOS
	GOOS=darwin GOARCH=amd64 go build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-amd64 ./cmd/repomap
	GOOS=darwin GOARCH=arm64 go build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-arm64 ./cmd/repomap
	# Windows
	GOOS=windows GOARCH=amd64 go build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-windows-amd64.exe ./cmd/repomap
	
	@echo "Generating checksums..."
	cd $(BUILD_DIR) && sha256sum * > checksums.txt
	@echo "Release build complete in $(BUILD_DIR)/"
