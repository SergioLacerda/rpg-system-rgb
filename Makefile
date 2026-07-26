.PHONY: fmt hello test validate

GO ?= go
GOCACHE ?= /tmp/go-cache
GOENV := GOCACHE=$(GOCACHE)

fmt:
	$(GOENV) $(GO) fmt ./...

hello:
	$(GOENV) $(GO) run ./cmd/rgb

test:
	$(GOENV) $(GO) test ./...

validate:
	$(GOENV) $(GO) run scripts/validate_semantic_docs.go
