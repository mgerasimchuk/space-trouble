FROM golang:1.20.0 AS build

COPY . /src

WORKDIR /src
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 GO111MODULE=on go build -mod vendor -a -installsuffix cgo -o /bin/verifier ./cmd/verifier/main.go

# I can't make slim image because I should wait postgres in docker compose with bash...:(
# FROM scratch
# COPY --from=build /bin/verifier /app/verifier
# COPY /src/internal/infrastructure/config/config.yml /app/config.yml

RUN mkdir /app &&\
    cp /bin/verifier /app/verifier &&\
    cp /src/internal/infrastructure/config/config.yml /app/config.yml

ENV CONFIG_FILE=/app/config.yml

WORKDIR /app

ENTRYPOINT ["./verifier"]
