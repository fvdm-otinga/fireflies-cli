BINARY  := fireflies
PKG     := github.com/fvdm-otinga/fireflies-cli
VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo dev)

.PHONY: build test lint vet generate docs release-dry scrub clean tidy

build: ## Build the fireflies binary into ./$(BINARY)
	go build -trimpath -ldflags "-s -w -X main.Version=$(VERSION) -X main.Commit=none -X main.Date=unknown" -o $(BINARY) .

test: ## Run all tests with race detector
	go test -race ./...

vet: ## Run go vet
	go vet ./...

lint: ## Requires golangci-lint in PATH
	golangci-lint run ./...

generate: ## Regenerate genqlient types from schema + .graphql files
	go run github.com/Khan/genqlient

docs: build ## Generate per-command Markdown reference into docs/reference/
	go run -tags gendocs scripts/gen-docs.go

release-dry: ## goreleaser snapshot (no publish) — prints archive sizes when done
	goreleaser release --snapshot --clean --skip=publish,sbom
	@echo ""
	@echo "=== Archive sizes ==="
	@find dist/ -name '*.tar.gz' -o -name '*.zip' | sort | xargs -I{} sh -c 'printf "%s\t%s\n" "$$(du -sh {} | cut -f1)" "{}"'

scrub: ## Scrub testdata/fixtures of any leaked secrets (in-place)
	bash scripts/scrub-fixtures.sh testdata/fixtures

scrub-check: ## CI check: fail if any secrets are found in fixtures
	bash scripts/scrub-fixtures.sh --check testdata/fixtures

tidy: ## Run go mod tidy
	go mod tidy

clean: ## Remove build artifacts
	rm -rf $(BINARY) dist/
