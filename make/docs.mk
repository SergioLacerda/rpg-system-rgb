.PHONY: docs-build docs-pdf docs-preview

docs-build: FORCE ## Build documentation Library with Go
	$(GOENV) $(GO) run ./cmd/rgb docs library --source "$(DOCS_SOURCE)" --out "$(LIBRARY_DIR)"

docs-pdf: FORCE ## Build and publish latest PDF downloads locally
	$(GOENV) $(GO) run ./cmd/rgb docs pdf --public-dir "$(PDF_PUBLIC_DIR)" --basename "$(PDF_BASENAME)" --version "$(PDF_VERSION)" --source-en "$(PDF_SRC_EN)" --source-pt-br "$(PDF_SRC_PT_BR)"
	$(MAKE) release-artifact-manifest
	$(MAKE) release-supply-chain

docs-preview: FORCE ## Serve documentation locally
	$(MAKE) docs-build
	@printf '%s\n' "Library built at $(LIBRARY_DIR)/index.html"
