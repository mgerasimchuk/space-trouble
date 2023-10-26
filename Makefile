DOCKER_COMPOSE_FILENAME=deployments/docker-compose/docker-compose.yaml

start: vendor-install
	docker-compose -f ${DOCKER_COMPOSE_FILENAME} down -v --remove-orphans ; docker-compose -f ${DOCKER_COMPOSE_FILENAME} up -d
start-from-scratch: vendor-install
	docker-compose -f ${DOCKER_COMPOSE_FILENAME} down -v --remove-orphans ; docker-compose -f ${DOCKER_COMPOSE_FILENAME} up --build -d
stop:
	docker-compose -f ${DOCKER_COMPOSE_FILENAME} down
restart: stop start
logs-api:
	docker-compose -f ${DOCKER_COMPOSE_FILENAME} logs -f api
logs-verifier:
	docker-compose -f ${DOCKER_COMPOSE_FILENAME} logs -f verifier
test:
	docker run -v $(shell pwd):/app golang:1.20.0 /bin/bash -c 'cd /app && GO111MODULE=on go test -mod vendor -covermode=count -coverprofile=assets/coverage/coverage.out -v ./... && go tool cover -html=assets/coverage/coverage.out -o=assets/coverage/coverage.html'
lint-architecture:
	docker run --rm -v $(shell pwd):/app -w /app golang:1.20.0 /bin/bash -c "go install github.com/fdaines/arch-go@v1.4.2 && arch-go -v"
vendor-install:
	@if [ -d "vendor" ]; then echo "Vendor folder already exists. Skip vendor installing."; else docker run --rm -v $(shell pwd):/app -w /app golang:1.20.0 /bin/bash -c "go mod tidy && go mod vendor"; fi