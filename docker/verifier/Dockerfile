FROM golang:1.20.0 AS build

WORKDIR /app

COPY . /app

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 GO111MODULE=on go build -mod vendor -a -installsuffix cgo -o /bin/verifier ./cmd/verifier/main.go

# I can't make slim image because I should wait postgres in docker compose with bash...:(
# FROM scratch
# COPY --from=build /bin/verifier /app/verifier

RUN cp /bin/verifier /verifier

ENTRYPOINT ["/verifier"]
