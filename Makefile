GO ?= go
GOFMT ?= gofmt
GOLANGCI_LINT ?= ./bin/golangci-lint
LINT_PLUGIN ?= plugin/domainctor.so

.PHONY: test fmt lint vet tidy help lint-plugin tools

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
	$(GOFMT) -w $(shell find internal -name '*.go')

test:
	$(GO) test ./...

vet:
	$(GO) vet ./...

tidy:
	$(GO) mod tidy

lint: lint-plugin
	$(GOLANGCI_LINT) run ./...

lint-plugin:
	$(GO) build -buildmode=plugin -o $(LINT_PLUGIN) ./plugin/domainctor.go

tools:
	@mkdir -p ./bin
	GOBIN=$(CURDIR)/bin $(GO) install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@v2.8.0
