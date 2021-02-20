GO ?= go1.16
export GO

.PHONY: build
build: fmt ## Build truce into bin folder
	@mkdir -p bin
	@echo "Building Truce from source."
	@$(GO) build -o bin/truce ./cmd/truce/...

.PHONY: install
install: ## Install truce globally
	$(GO) install ./cmd/truce/...

.PHONY: test
test: build
	$(GO) test -count 5 -race ./...

.PHONY: fmt
fmt: deps ## Run go fmt -s and cue fmt all over the shop
	@cue fmt ./cue/...
	@gofmt -s -w $(shell find . -name "*.go")

require-buildtools:
	@$(GO) version >/dev/null || (echo "${GO} is currently required. Try 'make install-go'." && exit 2)

.PHONY: install-go
install-go: ## Install latest beta version of Go (requires a stable version of Go).
	go get golang.org/dl/$(GO)
	$(GO) download

.PHONY: deps
deps: require-buildtools ## Download dependencies
	@$(GO) mod download
	@$(GO) get cuelang.org/go/cue
	@$(GO) mod tidy

.PHONY: examplegen ## Re-generate example project
examplegen: cleanExample build ## generate example directory cue services
	@bin/truce gen example/service.cue

# http://marmelab.com/blog/2016/02/29/auto-documented-makefile.html
.PHONY: help
help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

cleanExample:
	@rm example/{types,server,client}.go 2> /dev/null || true

.DEFAULT_GOAL := help
