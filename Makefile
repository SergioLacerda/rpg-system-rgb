GO ?= go
GOCACHE ?= /tmp/go-cache
GOENV := GOCACHE=$(GOCACHE)
GOLANGCI_VERSION ?= v2.11.4
GOLANGCI ?= $(shell command -v golangci-lint 2>/dev/null || echo $(HOME)/go/bin/golangci-lint)
GOVULNCHECK_VERSION ?= v1.7.0
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
PDF_REGRESSION_CANDIDATE_DIR ?= $(PDF_PUBLIC_DIR)
RELEASE_MANIFEST ?= $(PDF_PUBLIC_DIR)/$(PDF_BASENAME)-release-manifest.json
RELEASE_CHECKSUMS ?= $(PDF_PUBLIC_DIR)/SHA256SUMS
RELEASE_SBOM ?= $(PDF_PUBLIC_DIR)/$(PDF_BASENAME)-$(PDF_VERSION)-sbom.spdx.json
RELEASE_PROVENANCE ?= $(PDF_PUBLIC_DIR)/$(PDF_BASENAME)-$(PDF_VERSION)-provenance.json
PDF_SRC ?=
PDF_SRC_EN ?=
PDF_SRC_PT_BR ?=
PDF_VERSION ?= v2.0
SKILL_SOURCE_DIR ?= skills/specialist
SKILL_NAME ?= rgb-specialist
SKILL_VERSION ?= v0.1
SKILL_MANIFEST ?= $(PDF_PUBLIC_DIR)/$(SKILL_NAME)-skill-manifest.json
SKILL_CHECKSUMS ?= $(PDF_PUBLIC_DIR)/$(SKILL_NAME)-SHA256SUMS
COMPONENTS_COVER_THRESHOLD ?= 90
BUNDLES_COVER_THRESHOLD ?= 90
CORE_COVER_THRESHOLD ?= 90
FIXTURES_COVER_THRESHOLD ?= 90
MAKER_COVER_THRESHOLD ?= 90
SPECIALIST_COVER_THRESHOLD ?= 90
TOOLING_COVER_THRESHOLD ?= 90
PUBLICATION_COVER_THRESHOLD ?= 90
APP_COVER_THRESHOLD ?= 90
CLI_COVER_THRESHOLD ?= 90
TOOLING_CLI_COVER_THRESHOLD ?= 90
SKILLPKG_COVER_THRESHOLD ?= 90

.PHONY: help install review-structure check-fast check release-check FORCE

.DEFAULT_GOAL := help

FORCE:

##@ Meta & Composite Gates

help: FORCE ## List available targets
	@awk 'BEGIN {FS = ":.*?## "} \
	  /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } \
	  /^[a-zA-Z0-9_.-]+:.*?## / { printf "  \033[36m%-28s\033[0m %s\n", $$1, $$2 }' \
	  $(MAKEFILE_LIST)

install: landing-install FORCE ## Install local development dependencies

include make/go.mk
include make/docs.mk
include make/web.mk
include make/release.mk
include make/governance.mk

# Full base-structure review gate. See:
# docs/engineering/base-structure-review-workflow.md
review-structure: vet lint test test-arch validate cover-check FORCE ## Run the full base-structure review gate
	@echo "review-structure: all gates passed"

check-fast: fmt-check vet lint test test-arch validate check-generated-drift cover-check mutation-core check-publication-runtime check-governance-files FORCE ## Run the fast local and PR gate
	@echo "check-fast: all gates passed"

check: check-fast lint-web test-web landing-build lint-yaml lint-shell FORCE ## Run the full development and main gate
	@echo "check: all gates passed"

release-check: check docs-pdf pdf-editorial-check skill-package release-skill-check FORCE ## Run the full release gate
	@echo "release-check: all gates passed"
