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
PDF_BASENAME ?= rgb-system-core-v2
PDF_BUILD_DIR ?= $(LANDING_DIR)/.pdfbuild
PDF_EN_CONFIG ?= docs-build/mkdocs-pdf-en.yml
PDF_LOCALE ?= pt-br
PDF_PUBLIC_DIR ?= $(LANDING_DIR)/public/downloads
PDF_PT_BR_CONFIG ?= docs-build/mkdocs-pdf-pt-br.yml
MKDOCS_CONFIG ?= docs-build/mkdocs.yml
PDF_SRC ?=
PDF_VERSION ?= v0.2
# Coverage floor for cover-check. Deliberately low relative to the current
# total (see `make cover`) so it fails only on a real regression, not on
# day-to-day fluctuation. Raise it over time as coverage improves — never
# lower it without recording why in an ADR (mirrors the gocyclo ratchet).
COVER_THRESHOLD ?= 30

.PHONY: help install fmt test test-arch cover cover-check vet lint validate \
        generate bundle compile compile-full docs-install docs-pdf-install \
        docs-build docs-pdf docs-preview landing-install landing-check \
        landing-build landing-preview preview pdf-publish review-structure \
        go-file-size-report lint-web lint-web-fix test-web lint-yaml lint-shell FORCE

.DEFAULT_GOAL := help

FORCE:

help: FORCE ## List available targets
	@grep -E '^[a-zA-Z_-]+:.*?## ' $(MAKEFILE_LIST) | sort | \
		awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-22s\033[0m %s\n", $$1, $$2}'

fmt: FORCE ## Format Go code
	$(GOENV) $(GO) fmt ./...

install: docs-install landing-install FORCE ## Install local development dependencies

test: FORCE ## Run all Go tests
	$(GOENV) $(GO) test ./...

test-arch: FORCE ## Run architecture tests
	$(GOENV) $(GO) test ./tests/architecture/ -v

cover: FORCE ## Generate and print Go coverage summary
	$(GOENV) $(GO) test ./... -coverprofile=/tmp/rgb-coverage.out
	$(GOENV) $(GO) tool cover -func=/tmp/rgb-coverage.out | tail -5

# cover-check fails the build if total statement coverage drops below
# COVER_THRESHOLD. The threshold starts low by design (see the variable's
# comment above) — it is a floor against regression, not a target.
cover-check: FORCE ## Enforce the Go coverage floor
	$(GOENV) $(GO) test ./... -coverprofile=/tmp/rgb-coverage.out
	@pct=$$($(GOENV) $(GO) tool cover -func=/tmp/rgb-coverage.out | tail -1 | grep -oE '[0-9]+\.[0-9]+'); \
	echo "total coverage: $$pct% (floor: $(COVER_THRESHOLD)%)"; \
	awk -v pct="$$pct" -v floor="$(COVER_THRESHOLD)" 'BEGIN { exit !(pct + 0 >= floor + 0) }' \
		|| { echo "coverage $$pct% is below the $(COVER_THRESHOLD)% floor"; exit 1; }

vet: FORCE ## Run go vet
	$(GOENV) $(GO) vet ./...

lint: FORCE ## Run golangci-lint, falling back to go vet when unavailable
	@if [ -x "$(GOLANGCI)" ]; then \
		$(GOENV) "$(GOLANGCI)" run ./...; \
	else \
		echo "golangci-lint not found; falling back to go vet"; \
		echo "install: go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@latest"; \
		$(GOENV) $(GO) vet ./...; \
	fi

validate: FORCE ## Validate generated RGB content
	$(GOENV) $(GO) run ./cmd/rgb-tooling validate

generate: FORCE ## Generate RGB content artifacts
	$(GOENV) $(GO) run ./cmd/rgb-tooling generate

bundle: FORCE ## Bundle RGB content artifacts
	$(GOENV) $(GO) run ./cmd/rgb-tooling bundle

compile: FORCE ## Compile RGB content without HTML output
	$(GOENV) $(GO) run ./cmd/rgb-compiler no-html

compile-full: FORCE ## Compile all RGB content outputs
	$(GOENV) $(GO) run ./cmd/rgb-compiler all

docs-install: FORCE ## Install MkDocs dependencies
	$(PYTHON) -m pip install -r docs-build/requirements-docs.txt

docs-pdf-install: docs-install FORCE ## Install MkDocs PDF dependencies
	$(PYTHON) -m pip install -r docs-build/requirements-docs-pdf.txt

docs-build: FORCE ## Build documentation with MkDocs strict mode
	$(MKDOCS) build --strict -f $(MKDOCS_CONFIG)

docs-pdf: FORCE ## Build and publish latest PDF downloads locally
	$(MKDOCS) build -f $(PDF_EN_CONFIG)
	$(MKDOCS) build -f $(PDF_PT_BR_CONFIG)
	mkdir -p "$(PDF_PUBLIC_DIR)"
	cp "$(PDF_BUILD_DIR)/en/$(PDF_BASENAME)-latest-en.pdf" "$(PDF_PUBLIC_DIR)/$(PDF_BASENAME)-latest-en.pdf"
	cp "$(PDF_BUILD_DIR)/pt-br/$(PDF_BASENAME)-latest-pt-br.pdf" "$(PDF_PUBLIC_DIR)/$(PDF_BASENAME)-latest-pt-br.pdf"
	cp "$(PDF_PUBLIC_DIR)/$(PDF_BASENAME)-latest-en.pdf" "$(PDF_PUBLIC_DIR)/$(PDF_BASENAME)-$(PDF_VERSION)-en.pdf"
	cp "$(PDF_PUBLIC_DIR)/$(PDF_BASENAME)-latest-pt-br.pdf" "$(PDF_PUBLIC_DIR)/$(PDF_BASENAME)-$(PDF_VERSION)-pt-br.pdf"

docs-preview: FORCE ## Serve documentation locally
	$(MKDOCS) serve --dev-addr $(MKDOCS_HOST):$(MKDOCS_PORT) -f $(MKDOCS_CONFIG)

landing-install: FORCE ## Install landing page npm dependencies
	$(NPM) --prefix $(LANDING_DIR) install

landing-check: landing-install FORCE ## Run Astro TypeScript checks for the landing page
	$(NPM) --prefix $(LANDING_DIR) run check

landing-build: docs-build FORCE ## Build the landing page
	ASTRO_BASE="$(LANDING_BASE)" $(NPM) --prefix $(LANDING_DIR) run build

landing-preview: landing-build FORCE ## Preview the landing page locally
	ASTRO_BASE="$(LANDING_BASE)" $(NPM) --prefix $(LANDING_DIR) run preview -- --host $(LANDING_HOST) --port $(LANDING_PORT)

preview: landing-preview FORCE ## Preview the full published site locally

pdf-publish: FORCE ## Publish a provided PDF into landing downloads
	@test -n "$(PDF_SRC)" || { echo "PDF_SRC is required"; exit 1; }
	@test -f "$(PDF_SRC)" || { echo "PDF_SRC not found: $(PDF_SRC)"; exit 1; }
	@test "$(PDF_LOCALE)" = "pt-br" -o "$(PDF_LOCALE)" = "en" || { echo "PDF_LOCALE must be pt-br or en"; exit 1; }
	mkdir -p "$(PDF_PUBLIC_DIR)"
	cp "$(PDF_SRC)" "$(PDF_PUBLIC_DIR)/$(PDF_BASENAME)-$(PDF_VERSION)-$(PDF_LOCALE).pdf"
	cp "$(PDF_SRC)" "$(PDF_PUBLIC_DIR)/$(PDF_BASENAME)-latest-$(PDF_LOCALE).pdf"
	@printf '%s\n' "published $(PDF_PUBLIC_DIR)/$(PDF_BASENAME)-$(PDF_VERSION)-$(PDF_LOCALE).pdf"
	@printf '%s\n' "updated $(PDF_PUBLIC_DIR)/$(PDF_BASENAME)-latest-$(PDF_LOCALE).pdf"

# Full base-structure review gate. See:
# docs/engineering/base-structure-review-workflow.md
review-structure: vet lint test test-arch validate cover-check FORCE ## Run the full base-structure review gate
	@echo "review-structure: all gates passed"

# go-file-size-report lists non-test .go files under cmd/ and internal/
# over 200 lines. Informational only — does not fail the build.
go-file-size-report: FORCE ## Report large non-test Go source files
	@files=$$(find cmd internal -type f -name '*.go' ! -name '*_test.go' | sort); \
	results=$$(for f in $$files; do \
		lines=$$(wc -l < "$$f" | tr -d ' '); \
		if [ "$$lines" -gt 200 ]; then printf "%s %s\n" "$$f" "$$lines"; fi; \
	done | sort -k2,2nr -k1,1); \
	if [ -n "$$results" ]; then printf "%s\n" "$$results"; else echo "none"; fi

lint-web: landing-install FORCE ## Run landing page lint/type checks
	$(NPM) --prefix $(LANDING_DIR) run lint

lint-web-fix: landing-install FORCE ## Format landing page files and rerun lint/type checks
	$(NPM) --prefix $(LANDING_DIR) run lint:fix

test-web: landing-install FORCE ## Run landing page unit tests with coverage
	$(NPM) --prefix $(LANDING_DIR) run test

lint-yaml: FORCE ## Lint GitHub Actions workflows with actionlint
	@if command -v actionlint >/dev/null 2>&1; then \
		actionlint; \
	else \
		echo "actionlint not found; install: https://github.com/rhysd/actionlint"; \
		exit 1; \
	fi

lint-shell: FORCE ## Lint CI shell scripts with shellcheck
	@if command -v shellcheck >/dev/null 2>&1; then \
		shellcheck scripts/ci/*.sh; \
	else \
		echo "shellcheck not found; install: https://www.shellcheck.net"; \
		exit 1; \
	fi
