BINARY := cih-mr
BIN_DIR := bin

.PHONY: all build test lint install clean help

all: lint test build

build: ## Build the binary to bin/cih-mr
	go build -o $(BIN_DIR)/$(BINARY) .

test: ## Run all tests with race detection and coverage
	go test ./... -race -cover

lint: ## Run go vet across all packages
	go vet ./...

install: ## Install the binary via go install
	go install

clean: ## Remove the bin/ directory
	rm -rf $(BIN_DIR)

help: ## Show this help message
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-10s\033[0m %s\n", $$1, $$2}'
