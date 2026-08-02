GO ?= go
GOCACHE ?= /tmp/go-cache
GOENV := GOCACHE=$(GOCACHE)
GOLANGCI_VERSION ?= v2.11.4
GOLANGCI ?= $(shell command -v golangci-lint 2>/dev/null || echo $(HOME)/go/bin/golangci-lint)
ACTIONLINT_VERSION ?= 1.7.9
SHELLCHECK_VERSION ?= 0.10.0
TOOLS_DIR ?= /tmp/rgb-system-tools
TOOLS_BIN_DIR := $(TOOLS_DIR)/bin
ACTIONLINT ?= $(TOOLS_BIN_DIR)/actionlint
SHELLCHECK ?= $(TOOLS_BIN_DIR)/shellcheck
NPM ?= npm
LANDING_DIR ?= web/landing
LANDING_HOST ?= 127.0.0.1
LANDING_PORT ?= 4324
LANDING_BASE ?= /rpg-system-rgb
DOCS_SOURCE ?= docs
LIBRARY_DIR ?= $(LANDING_DIR)/public/library
PDF_BASENAME ?= rgb-system-core-v2
PDF_LOCALE ?= pt-br
PDF_PUBLIC_DIR ?= $(LANDING_DIR)/public/downloads
RELEASE_MANIFEST ?= $(PDF_PUBLIC_DIR)/$(PDF_BASENAME)-release-manifest.json
RELEASE_CHECKSUMS ?= $(PDF_PUBLIC_DIR)/SHA256SUMS
PDF_SRC ?=
PDF_SRC_EN ?=
PDF_SRC_PT_BR ?=
PDF_VERSION ?= v0.2
# Coverage floor for cover-check. Deliberately low relative to the current
# total (see `make cover`) so it fails only on a real regression, not on
# day-to-day fluctuation. Raise it over time as coverage improves — never
# lower it without recording why in an ADR (mirrors the gocyclo ratchet).
COVER_THRESHOLD ?= 30

.PHONY: help install fmt fmt-check test test-arch cover cover-check vet lint validate \
        generate bundle \
        tools-install actionlint-install shellcheck-install \
        docs-build docs-pdf docs-preview landing-install landing-check \
        landing-build landing-preview preview pdf-publish review-structure \
        go-file-size-report lint-web lint-web-fix test-web lint-yaml lint-shell \
        check-fast check check-generated-drift check-publication-runtime release-artifact-manifest \
        release-artifact-check release-check check-governance-files FORCE

.DEFAULT_GOAL := help

FORCE:

help: FORCE ## List available targets
	@grep -E '^[a-zA-Z_-]+:.*?## ' $(MAKEFILE_LIST) | sort | \
		awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-22s\033[0m %s\n", $$1, $$2}'

fmt: FORCE ## Format Go code
	$(GOENV) $(GO) fmt ./...

fmt-check: FORCE ## Check Go formatting without modifying files
	@files=$$(git ls-files '*.go' | xargs gofmt -l); \
	if [ -n "$$files" ]; then \
		echo "Go files need formatting:"; \
		printf '%s\n' "$$files"; \
		exit 1; \
	fi

install: landing-install FORCE ## Install local development dependencies

tools-install: actionlint-install shellcheck-install FORCE ## Install pinned local validation tools

actionlint-install: FORCE ## Install pinned actionlint into TOOLS_DIR
	@if [ ! -x "$(ACTIONLINT)" ]; then \
		scripts/ci/install-actionlint.sh "$(ACTIONLINT_VERSION)" "$(TOOLS_BIN_DIR)"; \
	fi
	@"$(ACTIONLINT)" -version

shellcheck-install: FORCE ## Install pinned shellcheck into TOOLS_DIR
	@if [ ! -x "$(SHELLCHECK)" ]; then \
		scripts/ci/install-shellcheck.sh "$(SHELLCHECK_VERSION)" "$(TOOLS_BIN_DIR)"; \
	fi
	@"$(SHELLCHECK)" --version

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
		echo "install: go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@$(GOLANGCI_VERSION)"; \
		$(GOENV) $(GO) vet ./...; \
	fi

validate: FORCE ## Validate generated RGB content
	$(GOENV) $(GO) run ./cmd/rgb validate

generate: FORCE ## Generate RGB content artifacts
	$(GOENV) $(GO) run ./cmd/rgb generate

bundle: FORCE ## Bundle RGB content artifacts
	$(GOENV) $(GO) run ./cmd/rgb bundle

docs-build: FORCE ## Build documentation Library with Go
	$(GOENV) $(GO) run ./cmd/rgb docs library --source "$(DOCS_SOURCE)" --out "$(LIBRARY_DIR)"

docs-pdf: FORCE ## Build and publish latest PDF downloads locally
	$(GOENV) $(GO) run ./cmd/rgb docs pdf --public-dir "$(PDF_PUBLIC_DIR)" --basename "$(PDF_BASENAME)" --version "$(PDF_VERSION)" --source-en "$(PDF_SRC_EN)" --source-pt-br "$(PDF_SRC_PT_BR)"
	$(MAKE) release-artifact-manifest

docs-preview: FORCE ## Serve documentation locally
	$(MAKE) docs-build
	@printf '%s\n' "Library built at $(LIBRARY_DIR)/index.html"

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

check-fast: fmt-check vet lint test test-arch validate cover-check FORCE ## Run the fast local and PR gate
	@echo "check-fast: all gates passed"

check: check-fast lint-web test-web landing-build lint-yaml lint-shell check-generated-drift check-publication-runtime check-governance-files FORCE ## Run the full development and main gate
	@echo "check: all gates passed"

check-generated-drift: FORCE ## Regenerate versioned artifacts and fail on drift
	scripts/ci/check-generated-drift.sh

check-publication-runtime: FORCE ## Reject retired publication runtime terms in active surfaces
	@p1='[Pp]yth'"on"; p2='[Mm]k[Dd]ocs'; p3='setup-py'"thon"; p4='requirements-'"docs"; p5='package-ecosystem: "p'"ip\""; \
	pattern="$$p1|$$p2|$$p3|$$p4|$$p5"; \
	if rg -n "$$pattern" . \
		--glob '!.analysis/**' \
		--glob '!.sdd/**' \
		--glob '!.strategist/**' \
		--glob '!docs/adr/**' \
		--glob '!docs/plans/**' \
		--glob '!web/landing/package-lock.json'; then \
		echo "retired publication runtime reference found"; \
		exit 1; \
	fi

release-artifact-manifest: FORCE ## Write release PDF manifest and checksums
	$(GOENV) $(GO) run ./cmd/rgb release manifest --public-dir "$(PDF_PUBLIC_DIR)" --basename "$(PDF_BASENAME)" --version "$(PDF_VERSION)" --manifest "$(RELEASE_MANIFEST)" --checksums "$(RELEASE_CHECKSUMS)"

release-artifact-check: FORCE ## Validate release PDF manifest and checksums
	$(GOENV) $(GO) run ./cmd/rgb release check --public-dir "$(PDF_PUBLIC_DIR)" --basename "$(PDF_BASENAME)" --version "$(PDF_VERSION)" --manifest "$(RELEASE_MANIFEST)" --checksums "$(RELEASE_CHECKSUMS)"

pdf-editorial-check: release-artifact-check FORCE ## Validate PDF editorial smoke and raster checks
	@echo "pdf-editorial-check: all gates passed"

.PHONY: pdf-editorial-check

release-check: check docs-pdf pdf-editorial-check FORCE ## Run the full release gate
	@echo "release-check: all gates passed"

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

lint-yaml: actionlint-install FORCE ## Lint GitHub Actions workflows with actionlint
	"$(ACTIONLINT)"

lint-shell: shellcheck-install FORCE ## Lint CI shell scripts with shellcheck
	"$(SHELLCHECK)" scripts/ci/*.sh

check-governance-files: FORCE ## Validate required public OSS governance files exist
	scripts/ci/check-governance-files.sh
