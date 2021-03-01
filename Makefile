METRICS_PORT=9090
RABDIS_CONFIG_FILE_LOCAL=rabdis.yaml.dist
RABDIS_CONFIG_FILE=/etc/rabdis.yaml

generate: ## Run go generate
	go generate

lint: ## Lint code
	golangci-lint run

test: ## Test packages
	go test -count=1 -cover -coverprofile=coverage.out -v ./...

coverage: ## Test coverage with default output
	go tool cover -func=coverage.out

coverage-html: ## Test coverage with html output
	go tool cover -html=coverage.html

clean: ## Clean project
	rm -Rf ./bin
	rm -Rf coverage.out

build: clean ## Build local binary
	mkdir -p ./bin
	go build -o ./bin ./cmd/rabdis

build-image: ## Build local image
	docker build -t ghcr.io/julienbreux/rabdis:latest .

run: build ## Run local binary
	./bin/rabdis

run-container: ## Run prepared local container
	docker run --rm -v $(PWD)/${RABDIS_CONFIG_FILE_LOCAL}:${RABDIS_CONFIG_FILE} -e RABDIS_CONFIG_FILE=${RABDIS_CONFIG_FILE} -p ${METRICS_PORT}:${METRICS_PORT} julienbreux/rabdis:latest

env-up: ## Up local environment
	docker-compose up --detach --force-recreate --remove-orphans

env-down: ## Down local environment
	docker-compose down --rmi all --volumes --remove-orphans

env-logs: ## Read and follow logs of local environment
	docker-compose logs --follow

help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

.PHONY: generate lint test coverage coverage-html clean build build-image run run-container env-up env-down env-logs help
