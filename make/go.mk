.PHONY: fmt fmt-check test test-arch cover cover-check vet lint validate generate bundle go-file-size-report mutation-core

fmt: FORCE ## Format Go code
	$(GOENV) $(GO) fmt ./...

fmt-check: FORCE ## Check Go formatting without modifying files
	@files=$$(git ls-files '*.go' | xargs gofmt -l); \
	if [ -n "$$files" ]; then \
		echo "Go files need formatting:"; \
		printf '%s\n' "$$files"; \
		exit 1; \
	fi

test: FORCE ## Run all Go tests
	$(GOENV) $(GO) test ./...

test-arch: FORCE ## Run architecture tests
	$(GOENV) $(GO) test ./tests/architecture/ -v

cover: FORCE ## Generate and print Go coverage summary
	$(GOENV) $(GO) test ./... -coverprofile=/tmp/rgb-coverage.out
	$(GOENV) $(GO) tool cover -func=/tmp/rgb-coverage.out | tail -5

cover-check: FORCE ## Enforce per-package Go coverage floors
	COMPONENTS_COVER_THRESHOLD="$(COMPONENTS_COVER_THRESHOLD)" \
	BUNDLES_COVER_THRESHOLD="$(BUNDLES_COVER_THRESHOLD)" \
	CORE_COVER_THRESHOLD="$(CORE_COVER_THRESHOLD)" \
	FIXTURES_COVER_THRESHOLD="$(FIXTURES_COVER_THRESHOLD)" \
	MAKER_COVER_THRESHOLD="$(MAKER_COVER_THRESHOLD)" \
	SPECIALIST_COVER_THRESHOLD="$(SPECIALIST_COVER_THRESHOLD)" \
	TOOLING_COVER_THRESHOLD="$(TOOLING_COVER_THRESHOLD)" \
	PUBLICATION_COVER_THRESHOLD="$(PUBLICATION_COVER_THRESHOLD)" \
	APP_COVER_THRESHOLD="$(APP_COVER_THRESHOLD)" \
	CLI_COVER_THRESHOLD="$(CLI_COVER_THRESHOLD)" \
	TOOLING_CLI_COVER_THRESHOLD="$(TOOLING_CLI_COVER_THRESHOLD)" \
	GO="$(GO)" GOCACHE="$(GOCACHE)" scripts/ci/check-go-coverage.sh

mutation-core: FORCE ## Run mutation smoke checks for internal/components/core
	GO="$(GO)" GOCACHE="$(GOCACHE)" scripts/ci/mutation-core.sh

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

# go-file-size-report lists non-test .go files under cmd/ and internal/
# over 200 lines. Informational only — does not fail the build.
go-file-size-report: FORCE ## Report large non-test Go source files
	@files=$$(find cmd internal -type f -name '*.go' ! -name '*_test.go' | sort); \
	results=$$(for f in $$files; do \
		lines=$$(wc -l < "$$f" | tr -d ' '); \
		if [ "$$lines" -gt 200 ]; then printf "%s %s\n" "$$f" "$$lines"; fi; \
	done | sort -k2,2nr -k1,1); \
	if [ -n "$$results" ]; then printf "%s\n" "$$results"; else echo "none"; fi
