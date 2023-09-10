help: ## Show Help
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

tools: ## Basic helper tools
	go install golang.org/x/tools/cmd/goimports@latest
	go install mvdan.cc/gofumpt@latest

fmt: ## Format files
	goimports -local gopherconbr.org/23/testcontainers/demo -w .
	gofumpt -l -w .

test: ## Run tests
	go test -v ./...

lint: ## Run linter
	docker run --rm -v $(PWD):/app -w /app golangci/golangci-lint:v1.52.2 golangci-lint run -v

.PHONY: test