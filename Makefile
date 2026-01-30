GO ?= go
GOFMT ?= gofmt

.PHONY: test fmt lint lint-fix vet tidy help

help:
	@echo "Targets:"
	@echo "  fmt   - format Go files"
	@echo "  test  - run tests"
	@echo "  vet   - run go vet"
	@echo "  tidy  - run go mod tidy"
	@echo "  lint  - alias for vet"
	@echo "  lint-fix - run fmt (auto-fix)"

fmt:
	$(GOFMT) -w $(shell find internal -name '*.go')

test:
	$(GO) test ./...

vet:
	$(GO) vet ./...

tidy:
	$(GO) mod tidy

lint: vet

lint-fix: fmt
