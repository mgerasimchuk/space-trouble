run:
	docker-compose down -v --remove-orphans ; docker-compose up -d
run-rebuild:
	docker-compose down -v --remove-orphans ; docker-compose up --build -d
test:
	docker run -v ${PWD}:/src golang:1.16.3 /bin/bash -c 'cd /src && GO111MODULE=on go test -mod vendor -covermode=count -coverprofile=cover.out -v ./... && go tool cover -html=cover.out -o=cover.html'
