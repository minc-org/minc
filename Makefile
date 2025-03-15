all: install

SHELL := /bin/bash -o pipefail

# Go and compilation related variables
VERSION ?= $(shell git describe --tags --dirty | tr -d v)
BUILD_DIR ?= out
SOURCE_DIRS = cmd pkg test
RELEASE_DIR ?= release

GOOS ?= $(shell go env GOOS)
GOARCH ?= $(shell go env GOARCH)

GOPATH ?= $(shell go env GOPATH)
ORG := github.com/minc-org

SOURCES := $(shell git ls-files '*.go' ":^vendor")
SOURCES := $(SOURCES) go.mod go.sum Makefile

LDFLAGS := -X github.com/minc-org/minc/pkg/constants.Version=$(VERSION) -extldflags='-static' -s -w $(GO_LDFLAGS)

# Add default target
.PHONY: default
default: install

# Create and update the vendor directory
.PHONY: vendor
vendor:
	go mod tidy
	go mod vendor

.PHONY: install
install: $(SOURCES)
	go install -ldflags="$(LDFLAGS)" ./cmd/minc

$(BUILD_DIR)/macos-amd64/minc_darwin_amd64: $(SOURCES)
	GOARCH=amd64 GOOS=darwin go build -ldflags="$(LDFLAGS)" -o $@ ./cmd/minc

$(BUILD_DIR)/macos-arm64/minc_darwin_arm64: $(SOURCES)
	GOARCH=arm64 GOOS=darwin go build -ldflags="$(LDFLAGS)" -o $@ ./cmd/minc

$(BUILD_DIR)/linux-amd64/minc_linux_amd64: $(SOURCES)
	GOOS=linux GOARCH=amd64 go build -ldflags="$(LDFLAGS)" -o $@ ./cmd/minc

$(BUILD_DIR)/linux-arm64/minc_linux_arm64: $(SOURCES)
	GOOS=linux GOARCH=arm64 go build -ldflags="$(LDFLAGS)" -o $@ ./cmd/minc

$(BUILD_DIR)/windows-amd64/minc.exe: $(SOURCES)
	GOARCH=amd64 GOOS=windows go build -ldflags="$(LDFLAGS)" -o $@ ./cmd/minc


.PHONY: cross ## Cross compiles all binaries
cross: $(BUILD_DIR)/macos-amd64/minc_darwin_amd64 $(BUILD_DIR)/macos-arm64/minc_darwin_arm64 $(BUILD_DIR)/linux-amd64/minc_linux_amd64 $(BUILD_DIR)/linux-arm64/minc_linux_arm64 $(BUILD_DIR)/windows-amd64/minc.exe

.PHONY: release ## Put all binary to release folder
release: cross
	mkdir -p $(BUILD_DIR)/release
	cp $(BUILD_DIR)/macos-amd64/minc_darwin_amd64 $(BUILD_DIR)/macos-arm64/minc_darwin_arm64 $(BUILD_DIR)/linux-amd64/minc_linux_amd64 $(BUILD_DIR)/linux-arm64/minc_linux_arm64 $(BUILD_DIR)/windows-amd64/minc.exe $(BUILD_DIR)/release

.PHONY: clean ## Remove all build artifacts
clean:
	rm -rf $(BUILD_DIR)
	rm -f $(GOPATH)/bin/minp

.PHONY: fmt
fmt:
	go fmt ./...
