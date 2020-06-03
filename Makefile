FILES_TO_FMT      ?= $(shell find . -path ./vendor -prune -o -name '*.go' -print)

GOFLAGS   :=
LDFLAGS   :=

BUILDTIME := $(shell date -u +"%Y-%m-%dT%H:%M:%SZ")
GITSHA := $(shell git rev-parse --short HEAD 2>/dev/null)

ifndef VERSION
	VERSION := git-$(GITSHA)
endif

GOFLAGS += -trimpath

LDFLAGS += -X $(PKG)/version.version=$(VERSION)
LDFLAGS += -X $(PKG)/version.commit=$(GITSHA)
LDFLAGS += -X $(PKG)/version.buildTime=$(BUILDTIME)

.PHONY: build
build:
	@go build $(GOFLAGS) -ldflags "$(LDFLAGS)" -o ./bin ./...

.PHONY: deps
deps: ## Ensures fresh go.mod and go.sum.
	@go mod tidy
	@go mod verify

.PHONY: lint
lint:
	@if [ ! -f ./bin/golangci-lint ]; then \
		curl -sfL https://install.goreleaser.com/github.com/golangci/golangci-lint.sh | sh -s $(GOLANGCI_LINT_VERSION); \
	fi;
	@echo "golangci-lint checking..."
	@./bin/golangci-lint run --deadline=30m --enable=misspell --enable=gosec --enable=gofmt --enable=goimports ./cmd/... ./...
	@go vet ./...

.PHONY: format
format: ## Formats Go code
	@echo ">> formatting code"
	@gofmt -s -w $(FILES_TO_FMT)

.PHONY: test-unit
test-unit: ## Runs all Youtube Go unit tests
test-unit:
	@go test -v -cover

.PHONY: test-integration
test-integration: ## Runs all Youtube Go integration tests
test-integration:
	@go test -v -tags="integration" -cover
