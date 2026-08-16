.PHONY: pdf-publish release-artifact-manifest release-artifact-check release-supply-chain release-supply-chain-check pdf-editorial-check

pdf-publish: FORCE ## Publish a provided PDF into landing downloads
	@test -n "$(PDF_SRC)" || { echo "PDF_SRC is required"; exit 1; }
	@test -f "$(PDF_SRC)" || { echo "PDF_SRC not found: $(PDF_SRC)"; exit 1; }
	@test "$(PDF_LOCALE)" = "pt-br" -o "$(PDF_LOCALE)" = "en" || { echo "PDF_LOCALE must be pt-br or en"; exit 1; }
	mkdir -p "$(PDF_PUBLIC_DIR)"
	cp "$(PDF_SRC)" "$(PDF_PUBLIC_DIR)/$(PDF_BASENAME)-$(PDF_VERSION)-$(PDF_LOCALE).pdf"
	cp "$(PDF_SRC)" "$(PDF_PUBLIC_DIR)/$(PDF_BASENAME)-latest-$(PDF_LOCALE).pdf"
	@printf '%s\n' "published $(PDF_PUBLIC_DIR)/$(PDF_BASENAME)-$(PDF_VERSION)-$(PDF_LOCALE).pdf"
	@printf '%s\n' "updated $(PDF_PUBLIC_DIR)/$(PDF_BASENAME)-latest-$(PDF_LOCALE).pdf"

release-artifact-manifest: FORCE ## Write release PDF manifest and checksums
	$(GOENV) $(GO) run ./cmd/rgb release manifest --public-dir "$(PDF_PUBLIC_DIR)" --basename "$(PDF_BASENAME)" --version "$(PDF_VERSION)" --manifest "$(RELEASE_MANIFEST)" --checksums "$(RELEASE_CHECKSUMS)"

release-artifact-check: FORCE ## Validate release PDF manifest and checksums
	$(GOENV) $(GO) run ./cmd/rgb release check --public-dir "$(PDF_PUBLIC_DIR)" --basename "$(PDF_BASENAME)" --version "$(PDF_VERSION)" --manifest "$(RELEASE_MANIFEST)" --checksums "$(RELEASE_CHECKSUMS)"

release-supply-chain: FORCE ## Write release SBOM and provenance metadata
	scripts/ci/write-release-supply-chain.sh "$(PDF_PUBLIC_DIR)" "$(PDF_BASENAME)" "$(PDF_VERSION)" "$(RELEASE_MANIFEST)" "$(RELEASE_CHECKSUMS)" "$(RELEASE_SBOM)" "$(RELEASE_PROVENANCE)"

release-supply-chain-check: FORCE ## Validate release SBOM and provenance metadata
	scripts/ci/check-release-supply-chain.sh "$(PDF_PUBLIC_DIR)" "$(PDF_BASENAME)" "$(PDF_VERSION)" "$(RELEASE_MANIFEST)" "$(RELEASE_CHECKSUMS)" "$(RELEASE_SBOM)" "$(RELEASE_PROVENANCE)"

pdf-editorial-check: release-artifact-check release-supply-chain-check FORCE ## Validate PDF editorial smoke, raster, and supply-chain checks
	@echo "pdf-editorial-check: all gates passed"
