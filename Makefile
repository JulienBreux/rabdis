METRICS_PORT=9090
RABDIS_CONFIG_FILE_LOCAL=rabdis.yaml.dist
RABDIS_CONFIG_FILE=/etc/rabdis.yaml

generate:
	go generate

lint:
	golangci-lint run

test:
	go test -count=1 -cover -coverprofile=coverage.out -v ./...

coverage:
	go tool cover -func=coverage.out

coverage-html:
	go tool cover -html=coverage.out

clean:
	rm -Rf ./bin
	rm -Rf coverage.out

build: clean
	mkdir -p ./bin
	go build -o ./bin ./cmd/rabdis

build-image:
	docker build -t ghcr.io/julienbreux/rabdis:latest .

run: build
	./bin/rabdis

run-container:
	docker run --rm -v $(PWD)/${RABDIS_CONFIG_FILE_LOCAL}:${RABDIS_CONFIG_FILE} -e RABDIS_CONFIG_FILE=${RABDIS_CONFIG_FILE} -p ${METRICS_PORT}:${METRICS_PORT} julienbreux/rabdis:latest

env-up:
	docker-compose up --detach --force-recreate --remove-orphans

env-down:
	docker-compose down --rmi all --volumes --remove-orphans

env-logs:
	docker-compose logs --follow

.PHONY: generate lint test coverage coverage-html clean build build-image run run-container env-up env-down env-logs
