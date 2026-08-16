.PHONY: landing-install landing-check landing-build landing-preview preview lint-web lint-web-fix test-web

landing-install: FORCE ## Install landing page npm dependencies
	$(NPM) --prefix $(LANDING_DIR) install

landing-check: landing-install FORCE ## Run Astro TypeScript checks for the landing page
	$(NPM) --prefix $(LANDING_DIR) run check

landing-build: docs-build FORCE ## Build the landing page
	ASTRO_BASE="$(LANDING_BASE)" $(NPM) --prefix $(LANDING_DIR) run build

landing-preview: landing-build FORCE ## Preview the landing page locally
	ASTRO_BASE="$(LANDING_BASE)" $(NPM) --prefix $(LANDING_DIR) run preview -- --host $(LANDING_HOST) --port $(LANDING_PORT)

preview: landing-preview FORCE ## Preview the full published site locally

lint-web: landing-install FORCE ## Run landing page lint/type checks
	$(NPM) --prefix $(LANDING_DIR) run lint

lint-web-fix: landing-install FORCE ## Format landing page files and rerun lint/type checks
	$(NPM) --prefix $(LANDING_DIR) run lint:fix

test-web: landing-install FORCE ## Run landing page unit tests with coverage
	$(NPM) --prefix $(LANDING_DIR) run test
