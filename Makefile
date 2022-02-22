ROOT_DIR := $(shell dirname $(realpath $(firstword $(MAKEFILE_LIST))))

# go-get-tool will 'go get' any package $2 and install it to $1.
define go-get-tool
@[ -f $(1) ] || { \
set -e ;\
TMP_DIR=$$(mktemp -d) ;\
cd $$TMP_DIR ;\
go mod init tmp ;\
echo "Downloading $(2)" ;\
GOBIN=$(ROOT_DIR)/bin go get $(2) ;\
rm -rf $$TMP_DIR ;\
}
endef

## --------------------------------------
## Tooling Binaries
## --------------------------------------

GINKGO := $(shell pwd)/bin/ginkgo
ginkgo: ## Download ginkgo locally if necessary.
	$(call go-get-tool,$(GINKGO),github.com/onsi/ginkgo/ginkgo@v1.16.4)

GOLANGCI_LINT := $(shell pwd)/bin/golangci-lint
golangci-lint: ## Download golangci-lint locally if necessary.
	$(call go-get-tool,$(GOLANGCI_LINT),github.com/golangci/golangci-lint/cmd/golangci-lint@v1.42.1)

## --------------------------------------
## Linting and fixing linter errors
## --------------------------------------

.PHONY: lint
lint: golangci-lint ## Lint codebase
	$(GOLANGCI_LINT) run -v --fast=false

.PHONY: lint-fix
lint-fix: $(GOLANGCI_LINT) ## Lint the codebase and run auto-fixers if supported by the linter.
	GOLANGCI_LINT_EXTRA_ARGS=--fix $(MAKE) lint

## --------------------------------------
## Testing
## --------------------------------------

.PHONY: test
test: lint ## Run tests.
	GO_ENV=test go test -v ./... -coverprofile=cover.out
