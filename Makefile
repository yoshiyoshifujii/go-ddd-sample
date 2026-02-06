GO ?= go
GOFMT ?= gofmt
GOLANGCI_LINT ?= ./bin/golangci-lint
LINT_PLUGIN ?= plugin/domainctor.so
CACHE_DIR ?= $(CURDIR)/.cache
GOCACHE ?= $(CACHE_DIR)/go-build
GOMODCACHE ?= $(CACHE_DIR)/go/pkg/mod
GOPATH ?= $(CACHE_DIR)/go
GOLANGCI_LINT_CACHE ?= $(CACHE_DIR)/golangci-lint
GOENV ?= GOCACHE=$(GOCACHE) GOMODCACHE=$(GOMODCACHE) GOPATH=$(GOPATH)

.PHONY: test fmt lint vet tidy help lint-plugin tools cache

help:
	@echo "Targets:"
	@echo "  fmt   - format Go files"
	@echo "  test  - run tests"
	@echo "  vet   - run go vet"
	@echo "  tidy  - run go mod tidy"
	@echo "  lint  - run golangci-lint"
	@echo "  lint-plugin - build golangci-lint plugin"
	@echo "  tools - install required tools into ./bin"

fmt:
	@mkdir -p $(GOCACHE) $(GOMODCACHE) $(GOPATH)
	$(GOFMT) -w $(shell find internal -name '*.go')

test:
	@mkdir -p $(GOCACHE) $(GOMODCACHE) $(GOPATH)
	$(GOENV) $(GO) test ./...

vet:
	@mkdir -p $(GOCACHE) $(GOMODCACHE) $(GOPATH)
	$(GOENV) $(GO) vet ./...

tidy:
	@mkdir -p $(GOCACHE) $(GOMODCACHE) $(GOPATH)
	$(GOENV) $(GO) mod tidy

lint: cache lint-plugin
	$(GOENV) GOLANGCI_LINT_CACHE=$(GOLANGCI_LINT_CACHE) $(GOLANGCI_LINT) run ./...

lint-plugin: cache
	$(GOENV) $(GO) build -buildmode=plugin -o $(LINT_PLUGIN) ./plugin/domainctor.go

cache:
	@mkdir -p $(GOCACHE) $(GOMODCACHE) $(GOPATH) $(GOLANGCI_LINT_CACHE)

tools:
	@mkdir -p ./bin
	cd tools && $(GOENV) GOBIN=$(CURDIR)/bin $(GO) install github.com/golangci/golangci-lint/v2/cmd/golangci-lint
