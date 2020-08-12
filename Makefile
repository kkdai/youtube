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

## help: Show makefile commands
.PHONY: help
help: Makefile
	@echo "---- Project: kkdai/youtube ----"
	@echo " Usage: make COMMAND"
	@echo
	@echo " Management Commands:"
	@sed -n 's/^##//p' $< | column -t -s ':' |  sed -e 's/^/ /'
	@echo

## build: Build project
.PHONY: build
build:
	@go build $(GOFLAGS) -ldflags "$(LDFLAGS)" -o ./bin ./cmd/youtubedr

## deps: Ensures fresh go.mod and go.sum
.PHONY: deps
deps:
	@go mod tidy
	@go mod verify

## lint: Run golangci-lint check
.PHONY: lint
lint:
	@if [ ! -f ./bin/golangci-lint ]; then \
		curl -sfL https://install.goreleaser.com/github.com/golangci/golangci-lint.sh | sh -s $(GOLANGCI_LINT_VERSION); \
	fi;
	@echo "golangci-lint checking..."
	@./bin/golangci-lint run --deadline=30m --enable=misspell --enable=gosec --enable=gofmt --enable=goimports --enable=golint ./cmd/... ./...
	@go vet ./...

## format: Formats Go code
.PHONY: format
format:
	@echo ">> formatting code"
	@gofmt -s -w $(FILES_TO_FMT)

## test-unit: Run all Youtube Go unit tests
.PHONY: test-unit
test-unit:
	@go test -v -cover . ./pkg/...


## test-integration: Run all Youtube Go integration tests
.PHONY: test-integration
test-integration:
	echo 'mode: atomic' > coverage.out
	go list ./... | xargs -n1 -I{} sh -c 'go test -race -tags=integration -covermode=atomic -coverprofile=coverage.tmp -coverpkg $(go list ./... | tr "\n" ",") {} && tail -n +2 coverage.tmp >> coverage.out || exit 255'
	rm coverage.tmp
