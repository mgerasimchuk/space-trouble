DOCKER_COMPOSE_FILENAME=deployments/docker-compose/docker-compose.yaml

start: vendor-install
	docker-compose -f ${DOCKER_COMPOSE_FILENAME} up -d
start-from-scratch: stop
	 docker-compose --file ${DOCKER_COMPOSE_FILENAME} up --remove-orphans --renew-anon-volumes --force-recreate --build -d
stop:
	docker-compose -f ${DOCKER_COMPOSE_FILENAME} down
restart: stop start
logs-api:
	docker-compose -f ${DOCKER_COMPOSE_FILENAME} logs -f api
logs-verifier:
	docker-compose -f ${DOCKER_COMPOSE_FILENAME} logs -f verifier
test-unit: vendor-install
	docker run -v $(shell pwd):/app -w /app golang:1.20.0 /bin/bash \
	-c 'go test -covermode=count -coverprofile=assets/coverage/coverage.out -v ./internal/... && go tool cover -html=assets/coverage/coverage.out -o=assets/coverage/coverage.html'
test-integration-api: start
	if ! docker-compose -f ${DOCKER_COMPOSE_FILENAME} run test-integration-api; then \
		echo "\nLogs for api from docker-compose:" ;\
		docker-compose -f ${DOCKER_COMPOSE_FILENAME} logs api; \
		echo "\nServices status:" ;\
		docker-compose -f ${DOCKER_COMPOSE_FILENAME} ps ;\
		exit 255 ;\
	fi
lint-golangci: vendor-install
	docker run --rm -v $(shell pwd):/app -w /app golangci/golangci-lint:v1.55.1 \
	golangci-lint run --timeout 5m30s -v
lint-architecture:
	docker run --rm -v $(shell pwd):/app -w /app golang:1.20.0 /bin/bash \
	-c "go install github.com/fdaines/arch-go@v1.4.2 && arch-go -v"
vendor-install:
	@if [ -d "vendor" ]; then echo "Vendor folder already exists. Skip vendor installing."; else docker run --rm -v $(shell pwd):/app -w /app golang:1.20.0 /bin/bash -c "go mod tidy && go mod vendor"; fi