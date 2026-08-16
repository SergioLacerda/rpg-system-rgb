.PHONY: tools-install actionlint-install shellcheck-install check-generated-drift check-publication-runtime check-governance-files lint-yaml lint-shell

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

check-governance-files: FORCE ## Validate required public OSS governance files exist
	scripts/ci/check-governance-files.sh

lint-yaml: actionlint-install FORCE ## Lint GitHub Actions workflows with actionlint
	"$(ACTIONLINT)"

lint-shell: shellcheck-install FORCE ## Lint CI shell scripts with shellcheck
	"$(SHELLCHECK)" scripts/ci/*.sh
