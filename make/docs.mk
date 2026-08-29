##@ Docs & Publication

.PHONY: docs-build docs-check docs-pdf docs-preview pdf-author pdf-author-unit-test skill-package

docs-build: FORCE ## Build documentation Library with Go
	$(GOENV) $(GO) run ./cmd/rgb docs library --source "$(DOCS_SOURCE)" --out "$(LIBRARY_DIR)"

docs-check: docs-build FORCE ## Validate generated public documentation artifacts
	$(GOENV) $(GO) run ./cmd/rgb docs check --library "$(LIBRARY_DIR)" --public-dir "$(PDF_PUBLIC_DIR)" --basename "$(PDF_BASENAME)" --version "$(PDF_VERSION)"

docs-pdf: FORCE ## Build and publish latest PDF downloads locally
	$(GOENV) $(GO) run ./cmd/rgb docs pdf --public-dir "$(PDF_PUBLIC_DIR)" --basename "$(PDF_BASENAME)" --version "$(PDF_VERSION)" --source-en "$(PDF_SRC_EN)" --source-pt-br "$(PDF_SRC_PT_BR)"
	$(MAKE) release-artifact-manifest
	$(MAKE) release-supply-chain

pdf-author: FORCE ## Manually author a reviewed PDF candidate with ADR-016 tooling
	tools/pdfbuild/build-pdf.sh "$(PDF_LOCALE)" "$(PDF_VERSION)"

pdf-author-unit-test: FORCE ## Run manual PDF authoring unit checks
	tools/pdfbuild/test-combine-html

skill-package: FORCE ## Build and publish the RGB Specialist skill .zip (ADR-013)
	$(GOENV) $(GO) run ./cmd/rgb docs skill --source "$(SKILL_SOURCE_DIR)" --out "$(PDF_PUBLIC_DIR)" --name "$(SKILL_NAME)" --version "$(SKILL_VERSION)"
	$(MAKE) release-skill-manifest

docs-preview: FORCE ## Serve documentation locally
	$(MAKE) docs-build
	@printf '%s\n' "Library built at $(LIBRARY_DIR)/index.html"

# --- Namespaced aliases (additive, non-breaking) ---
.PHONY: docs.build docs.check docs.pdf docs.preview docs.pdf-author docs.pdf-author-unit-test docs.skill-package
docs.build: docs-build
docs.check: docs-check
docs.pdf: docs-pdf
docs.preview: docs-preview
docs.pdf-author: pdf-author
docs.pdf-author-unit-test: pdf-author-unit-test
docs.skill-package: skill-package
