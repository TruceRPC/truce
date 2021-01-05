.PHONY: deps
deps: ## Download dependencies
	@go mod download
	@go get cuelang.org/go/cue

.PHONY: build
build: trucegen ## Build truce into bin folder
	@mkdir -p bin
	go build -o bin/truce ./cmd/truce/...

.PHONY: fmt
fmt: deps ## Run go fmt -s and cue fmt all over the shop
	@cue fmt truce.cue
	@gofmt -s -w $(shell find . -name "*.go")

.PHONY: trucegen
trucegen: fmt ## Generate truce specification
	@echo "Generating truce.go..."
	@go run internal/cmd/trucegen/main.go

.PHONY: examplegen
examplegen: trucegen build ## generate example directory cue services
	@rm example/{types,server,client}.go 2> /dev/null || true
	bin/truce -src example/service.cue gen types > example/types.go
	bin/truce -src example/service.cue gen server > example/server.go
	bin/truce -src example/service.cue gen client > example/client.go

# http://marmelab.com/blog/2016/02/29/auto-documented-makefile.html
.PHONY: help
help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

.DEFAULT_GOAL := help
