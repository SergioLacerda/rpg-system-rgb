.PHONY: fmt hello test test-arch cover lint vet validate generate compile compile-full review-structure

GO ?= go
GOCACHE ?= /tmp/go-cache
GOENV := GOCACHE=$(GOCACHE)
GOLANGCI ?= $(shell command -v golangci-lint 2>/dev/null || echo $(HOME)/go/bin/golangci-lint)

fmt:
	$(GOENV) $(GO) fmt ./...

hello:
	$(GOENV) $(GO) run ./cmd/rgb

test:
	$(GOENV) $(GO) test ./...

test-arch:
	$(GOENV) $(GO) test ./tests/architecture/ -v

cover:
	$(GOENV) $(GO) test ./... -coverprofile=/tmp/rgb-coverage.out
	$(GOENV) $(GO) tool cover -func=/tmp/rgb-coverage.out | tail -5

vet:
	$(GOENV) $(GO) vet ./...

lint:
	@if [ -x "$(GOLANGCI)" ]; then \
		$(GOENV) "$(GOLANGCI)" run ./...; \
	else \
		echo "golangci-lint not found; falling back to go vet"; \
		echo "install: go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@latest"; \
		$(GOENV) $(GO) vet ./...; \
	fi

validate:
	$(GOENV) $(GO) run ./cmd/rgb-tooling validate

generate:
	$(GOENV) $(GO) run ./cmd/rgb-tooling generate

compile:
	$(GOENV) $(GO) run ./cmd/rgb-compiler no-html

compile-full:
	$(GOENV) $(GO) run ./cmd/rgb-compiler all

# Full base-structure review gate. See:
# docs/engineering/base-structure-review-workflow.md
review-structure: vet lint test test-arch validate
	@echo "review-structure: all gates passed"
