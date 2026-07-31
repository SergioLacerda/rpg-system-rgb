.PHONY: FORCE

GO ?= go
GOCACHE ?= /tmp/go-cache
GOENV := GOCACHE=$(GOCACHE)
GOLANGCI ?= $(shell command -v golangci-lint 2>/dev/null || echo $(HOME)/go/bin/golangci-lint)
NPM ?= npm
PYTHON ?= python3
MKDOCS ?= $(PYTHON) -m mkdocs
MKDOCS_HOST ?= 127.0.0.1
MKDOCS_PORT ?= 8000
LANDING_DIR ?= web/landing
LANDING_HOST ?= 127.0.0.1
LANDING_PORT ?= 4324
LANDING_BASE ?= /rpg-system-rgb
# Coverage floor for cover-check. Deliberately low relative to the current
# total (see `make cover`) so it fails only on a real regression, not on
# day-to-day fluctuation. Raise it over time as coverage improves — never
# lower it without recording why in an ADR (mirrors the gocyclo ratchet).
COVER_THRESHOLD ?= 30

FORCE:

fmt: FORCE
	$(GOENV) $(GO) fmt ./...

hello: FORCE
	$(GOENV) $(GO) run ./cmd/rgb

test: FORCE
	$(GOENV) $(GO) test ./...

test-arch: FORCE
	$(GOENV) $(GO) test ./tests/architecture/ -v

cover: FORCE
	$(GOENV) $(GO) test ./... -coverprofile=/tmp/rgb-coverage.out
	$(GOENV) $(GO) tool cover -func=/tmp/rgb-coverage.out | tail -5

# cover-check fails the build if total statement coverage drops below
# COVER_THRESHOLD. The threshold starts low by design (see the variable's
# comment above) — it is a floor against regression, not a target.
cover-check: FORCE
	$(GOENV) $(GO) test ./... -coverprofile=/tmp/rgb-coverage.out
	@pct=$$($(GOENV) $(GO) tool cover -func=/tmp/rgb-coverage.out | tail -1 | grep -oE '[0-9]+\.[0-9]+'); \
	echo "total coverage: $$pct% (floor: $(COVER_THRESHOLD)%)"; \
	awk -v pct="$$pct" -v floor="$(COVER_THRESHOLD)" 'BEGIN { exit !(pct + 0 >= floor + 0) }' \
		|| { echo "coverage $$pct% is below the $(COVER_THRESHOLD)% floor"; exit 1; }

vet: FORCE
	$(GOENV) $(GO) vet ./...

lint: FORCE
	@if [ -x "$(GOLANGCI)" ]; then \
		$(GOENV) "$(GOLANGCI)" run ./...; \
	else \
		echo "golangci-lint not found; falling back to go vet"; \
		echo "install: go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@latest"; \
		$(GOENV) $(GO) vet ./...; \
	fi

validate: FORCE
	$(GOENV) $(GO) run ./cmd/rgb-tooling validate

generate: FORCE
	$(GOENV) $(GO) run ./cmd/rgb-tooling generate

compile: FORCE
	$(GOENV) $(GO) run ./cmd/rgb-compiler no-html

compile-full: FORCE
	$(GOENV) $(GO) run ./cmd/rgb-compiler all

docs-install: FORCE
	$(PYTHON) -m pip install -r requirements-docs.txt

docs-build: FORCE
	$(MKDOCS) build --strict

docs-preview: FORCE
	$(MKDOCS) serve --dev-addr $(MKDOCS_HOST):$(MKDOCS_PORT)

landing-install: FORCE
	$(NPM) --prefix $(LANDING_DIR) install

landing-build: docs-build FORCE
	ASTRO_BASE="$(LANDING_BASE)" $(NPM) --prefix $(LANDING_DIR) run build

landing-preview: landing-build FORCE
	ASTRO_BASE="$(LANDING_BASE)" $(NPM) --prefix $(LANDING_DIR) run preview -- --host $(LANDING_HOST) --port $(LANDING_PORT)

preview: landing-preview FORCE

# Full base-structure review gate. See:
# docs/engineering/base-structure-review-workflow.md
review-structure: vet lint test test-arch validate cover-check FORCE
	@echo "review-structure: all gates passed"

# go-file-size-report lists non-test .go files under cmd/ and internal/
# over 200 lines. Informational only — does not fail the build.
go-file-size-report: FORCE
	@files=$$(find cmd internal -type f -name '*.go' ! -name '*_test.go' | sort); \
	results=$$(for f in $$files; do \
		lines=$$(wc -l < "$$f" | tr -d ' '); \
		if [ "$$lines" -gt 200 ]; then printf "%s %s\n" "$$f" "$$lines"; fi; \
	done | sort -k2,2nr -k1,1); \
	if [ -n "$$results" ]; then printf "%s\n" "$$results"; else echo "none"; fi
