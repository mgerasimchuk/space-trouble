FROM golang:1.20.0 AS build

WORKDIR /app

COPY . /app

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 GO111MODULE=on go build -mod vendor -a -installsuffix cgo -o /bin/api ./cmd/api/main.go

# I can't make slim image because I should wait postgres in docker compose with bash...:(
# FROM scratch
# COPY --from=build /bin/api /app/api

RUN cp /bin/api /api

EXPOSE 8080

ENTRYPOINT ["/api"]