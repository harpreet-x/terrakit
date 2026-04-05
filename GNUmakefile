PROVIDER_BINARY := terrakit
CLI_BINARY      := terrakit
REGISTRY        := registry.terraform.io
NAMESPACE       := harpreet-x
NAME            := terrakit
VERSION         := 0.1.0
OS_ARCH         := $(shell go env GOOS)_$(shell go env GOARCH)

PLUGIN_DIR := $(HOME)/.terraform.d/plugins/$(REGISTRY)/$(NAMESPACE)/$(NAME)/$(VERSION)/$(OS_ARCH)

.PHONY: build install test vet clean help

## help: list available targets
help:
	@grep -E '^##' GNUmakefile | sed 's/## //'

## build: compile the provider binary and the terrakit CLI
build:
	go build -o $(PROVIDER_BINARY) .
	go build -o $(CLI_BINARY) ./cmd/terrakit

## install: build provider + CLI and install both into PATH / plugin cache
install: build
	mkdir -p $(PLUGIN_DIR)
	cp $(PROVIDER_BINARY) $(PLUGIN_DIR)/
	@INSTALL_DIR=$$(go env GOPATH)/bin; \
	 mkdir -p "$$INSTALL_DIR"; \
	 cp $(CLI_BINARY) "$$INSTALL_DIR/$(CLI_BINARY)"; \
	 echo "Installed: $$INSTALL_DIR/$(CLI_BINARY)"

## test: run all unit tests
test:
	go test ./... -v -count=1

## vet: run go vet across all packages
vet:
	go vet ./...

## clean: remove compiled binaries
clean:
	rm -f $(PROVIDER_BINARY) $(CLI_BINARY)
