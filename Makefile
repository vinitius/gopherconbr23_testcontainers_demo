help: ## Show Help
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

tools: ## Basic helper tools
	go install golang.org/x/tools/cmd/goimports@latest
	go install mvdan.cc/gofumpt@latest

fmt: ## Format files
	goimports -local gopherconbr.org/23/testcontainers/demo -w .
	gofumpt -l -w .

test: ## Run unit tests
	touch count.out
	go test -covermode=count -coverprofile=count.out ./...
	$(MAKE) coverage

coverage: ## Unit tests coverage
	go tool cover -func=count.out

lint: ## Run linter
	@golangci-lint -v run ./...
	@test -z "$$(golangci-lint run ./...)"

.PHONY: test