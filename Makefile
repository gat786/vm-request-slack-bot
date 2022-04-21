.DEFAULT_GOAL := help

TARGETS := build/dispatcher build/requests build/notify


build/%: lambdas/*/%.go
	@mkdir -p build
	go build -o $@ $<

.PHONY : compile
compile: $(TARGETS) ## compile

.PHONY : gomod_tidy
gomod_tidy: ## run go mod tidy
	go mod tidy

.PHONY : gofmt
gofmt: ## run go fmt
	go fmt -x ./...

.PHONY : help
help: ## show this help
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'
