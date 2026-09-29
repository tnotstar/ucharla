BINARY_NAME ?= ucharla
BUILD_DIR   ?= bin
BIN         := $(BUILD_DIR)/$(BINARY_NAME)

GO          ?= go
GOFMT       ?= gofmt

.PHONY: all fmt vet check build run clean help

all: build

fmt: ## Format Go source code
	$(GOFMT) -s -w .

vet: ## Run go vet static analysis
	$(GO) vet ./...

check: fmt vet ## Run formatting and static analysis

build: ## Build the binary
	@mkdir -p $(BUILD_DIR)
	$(GO) build -o $(BIN) .

run: build ## Build and run the application (pass args with ARGS="...")
	./$(BIN) $(ARGS)

clean: ## Clean build artifacts
	$(GO) clean
	rm -rf $(BUILD_DIR)

help: ## Show this help message
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "  %-8s %s\n", $$1, $$2}'
