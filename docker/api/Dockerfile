FROM golang:1.16.3 AS build

COPY . /src

WORKDIR /src
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 GO111MODULE=on go build -mod vendor -a -installsuffix cgo -o /bin/api ./cmd/api/main.go

# I can't make slim image because I should wait postgres in docker compose with bash...:(
# FROM scratch
# COPY --from=build /bin/api /app/api
# COPY /src/internal/infrastructure/config/config.yml /app/config.yml

RUN mkdir /app &&\
    cp /bin/api /app/api &&\
    cp /src/internal/infrastructure/config/config.yml /app/config.yml

ENV CONFIG_FILE=/app/config.yml

WORKDIR /app

EXPOSE 8080

ENTRYPOINT ["./api"]