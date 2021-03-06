NAME := knope

.PHONY: explain
explain:
	### Welcome
	#
	#  ___   _  __    _  _______  _______  _______
	# |   | | ||  |  | ||       ||       ||       |
	# |   |_| ||   |_| ||   _   ||    _  ||    ___|
	# |      _||       ||  | |  ||   |_| ||   |___
	# |     |_ |  _    ||  |_|  ||    ___||    ___|
	# |    _  || | |   ||       ||   |    |   |___
	# |___| |_||_|  |__||_______||___|    |_______|
	#
	#
	### Installation
	#
	# $$ make all
	#
	### Targets
	@cat Makefile* | grep -E '^[a-zA-Z_-]+:.*?## .*$$' | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

GITCOMMIT := $(shell git rev-parse --short HEAD)

.PHONY: clean
clean: ## Clean the local dependencies
	rm -fr vendor

.PHONY: install
install: ## Install the local dependencies
	go get ./...

.PHONY: vet
vet: ## Vet the code
	go vet -v ./...

.PHONY: lint
lint: ## Lint the code
	golint -set_exit_status $(shell go list ./... | grep -v vendor)

.PHONY: build
build: ## Build the application
	go build .

.PHONY: static
static: ## Build the application
	CGO_ENABLED=0 go build -ldflags "-extldflags -static -X github.com/benmatselby/$(NAME)/version.GITCOMMIT=$(GITCOMMIT)" -o $(NAME) .

.PHONY: test
test: ## Run the unit tests
	go test ./... -coverprofile=coverage.out

.PHONY: test-cov
test-cov: test ## Run the unit tests with coverage
	go tool cover -html=coverage.out

.PHONY: all ## Run everything
all: clean install vet build test

.PHONY: static-all ## Run everything
static-all: clean install vet static test
