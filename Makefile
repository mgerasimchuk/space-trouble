run: vendor-install
	docker-compose down -v --remove-orphans ; docker-compose up -d
run-rebuild: vendor-install
	docker-compose down -v --remove-orphans ; docker-compose up --build -d
vendor-install:
	@if [ -d "vendor" ]; then echo "Vendor folder already exists. Skip vendor installing."; else docker run --rm -v $(shell pwd):/app -w /app golang:1.20.0 /bin/bash -c "go mod tidy && go mod vendor"; fi
test:
	docker run -v $(shell pwd):/src golang:1.20.0 /bin/bash -c 'cd /src && GO111MODULE=on go test -mod vendor -covermode=count -coverprofile=cover.out -v ./... && go tool cover -html=cover.out -o=cover.html'
lint-architecture:
	docker run --rm -v $(shell pwd):/app -w /app golang:1.20.0 /bin/bash -c "go install github.com/fdaines/arch-go@v1.4.2 && arch-go -v"