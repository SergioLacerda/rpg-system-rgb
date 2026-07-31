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
	$(PYTHON) -m pip install -r docs-build/requirements-docs.txt

docs-pdf-install: docs-install FORCE
	$(PYTHON) -m pip install -r docs-build/requirements-docs-pdf.txt

docs-build: FORCE
	$(MKDOCS) build --strict -f $(MKDOCS_CONFIG)

docs-pdf: FORCE
	$(MKDOCS) build -f $(PDF_EN_CONFIG)
	$(MKDOCS) build -f $(PDF_PT_BR_CONFIG)
	mkdir -p "$(PDF_PUBLIC_DIR)"
	cp "$(PDF_BUILD_DIR)/en/$(PDF_BASENAME)-latest-en.pdf" "$(PDF_PUBLIC_DIR)/$(PDF_BASENAME)-latest-en.pdf"
	cp "$(PDF_BUILD_DIR)/pt-br/$(PDF_BASENAME)-latest-pt-br.pdf" "$(PDF_PUBLIC_DIR)/$(PDF_BASENAME)-latest-pt-br.pdf"
	cp "$(PDF_PUBLIC_DIR)/$(PDF_BASENAME)-latest-en.pdf" "$(PDF_PUBLIC_DIR)/$(PDF_BASENAME)-$(PDF_VERSION)-en.pdf"
	cp "$(PDF_PUBLIC_DIR)/$(PDF_BASENAME)-latest-pt-br.pdf" "$(PDF_PUBLIC_DIR)/$(PDF_BASENAME)-$(PDF_VERSION)-pt-br.pdf"

docs-preview: FORCE
	$(MKDOCS) serve --dev-addr $(MKDOCS_HOST):$(MKDOCS_PORT) -f $(MKDOCS_CONFIG)

landing-install: FORCE
	$(NPM) --prefix $(LANDING_DIR) install

landing-build: docs-build FORCE
	ASTRO_BASE="$(LANDING_BASE)" $(NPM) --prefix $(LANDING_DIR) run build

landing-preview: landing-build FORCE
	ASTRO_BASE="$(LANDING_BASE)" $(NPM) --prefix $(LANDING_DIR) run preview -- --host $(LANDING_HOST) --port $(LANDING_PORT)

preview: landing-preview FORCE

pdf-publish: FORCE
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
