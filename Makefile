.PHONY: all dep lint vet test test-coverage build clean

# custom define
PROJECT := uptoc
MAINFILE := cmd/main.go

all: build

dep: ## Get the dependencies
	@go mod download

lint: ## Lint Golang files
	@golangci-lint --version
	@golangci-lint run -D errcheck

test: ## Run tests with coverage
	go test -coverprofile .coverprofile ./...
	go tool cover --func=.coverprofile

coverage-html: ## show coverage by the html
	go tool cover -html=.coverprofile

build: dep ## Build the binary file
	@go build -o build/bin/$(PROJECT) $(MAINFILE)

clean: ## Remove previous build
	@rm -rf ./build

help: ## Display this help screen
	@grep -h -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'
