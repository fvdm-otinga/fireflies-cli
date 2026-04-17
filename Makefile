BINARY := fireflies
PKG    := github.com/fvdm-otinga/fireflies-cli
VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo dev)

.PHONY: build test lint vet generate release-dry clean tidy

build: ## Build the fireflies binary into ./$(BINARY)
	go build -trimpath -ldflags "-s -w -X main.Version=$(VERSION)" -o $(BINARY) .

test: ## Run all tests with race detector
	go test -race ./...

vet:
	go vet ./...

lint: ## Requires golangci-lint in PATH
	golangci-lint run ./...

generate: ## Regenerate genqlient types from schema + .graphql files
	go run github.com/Khan/genqlient

tidy:
	go mod tidy

release-dry: ## Run goreleaser in snapshot/dry-run mode
	goreleaser release --snapshot --clean --skip=publish

clean:
	rm -rf $(BINARY) dist/
